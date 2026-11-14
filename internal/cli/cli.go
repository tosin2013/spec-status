// Package cli is the command surface.
//
// Command, help, and flag handling use Cobra. The owner approved that
// dependency on 7 November 2026. Roadmap parsing and root-confined reads
// stay in the standard library, in internal/status and internal/safeio.
// Cobra is not a styling, network, or telemetry library, and this package
// does not add those.
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/tosin2013/spec-status/internal/status"
)

// Overridden with -ldflags -X when a builder stamps a commit.
// Defaults describe a development build, not a release.
// Kind is set to "release" only by the release workflow
// (.github/workflows/release.yml); see build.go.
var (
	Version = "0.0.0-dev"
	Commit  = "unknown"
	Kind    = ""
)

func init() {
	// An empty message disables Cobra's Windows explorer pause.
	cobra.MousetrapHelpText = ""
}

// completionDebugEnv is read by Cobra's hidden __complete error path.
// CompDebug opens it with O_APPEND|O_CREATE when set. Cobra exposes no hook
// to disable that write.
const completionDebugEnv = "BASH_COMP_DEBUG_FILE"

// completionDebugMu serializes the process-global environment around Run.
// Concurrent Run calls in one process wait here so one call cannot restore
// the variable while another is still inside Cobra. The product binary calls
// Run once from main. A caller must not assume the variable stays visible to
// other goroutines for the duration of Run.
var completionDebugMu sync.Mutex

// Run executes one invocation and returns the process exit code.
func Run(args []string, cwd string, stdout, stderr io.Writer) int {
	completionDebugMu.Lock()
	defer completionDebugMu.Unlock()
	restore, hideErr := hideCompletionDebugFile()
	if hideErr != nil {
		return failFor(args, jsonRequested(args), "SS_INVOCATION_INVALID", "Completion debug file setting could not be ignored.", stdout, stderr)
	}
	defer restore()
	jsonMode := jsonRequested(args)
	if msg := precheck(args); msg != "" {
		return failFor(args, jsonMode, "SS_INVOCATION_INVALID", msg, stdout, stderr)
	}
	root := newRoot(cwd)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs(args)
	err := root.Execute()
	var coded silentExit
	if errors.As(err, &coded) {
		return coded.code
	}
	if err == nil {
		return 0
	}
	return failFor(args, jsonMode, "SS_INVOCATION_INVALID", oneLine(err.Error()), stdout, stderr)
}

// hideCompletionDebugFile unsets BASH_COMP_DEBUG_FILE for this process and
// returns a function that restores the previous value. Cobra reads the
// variable with os.Getenv during hidden completion errors; clearing it is
// the supported way to keep that path from creating or appending a log.
func hideCompletionDebugFile() (func(), error) {
	prev, ok := os.LookupEnv(completionDebugEnv)
	if !ok {
		return func() {}, nil
	}
	if err := os.Unsetenv(completionDebugEnv); err != nil {
		return nil, err
	}
	return func() {
		_ = os.Setenv(completionDebugEnv, prev)
	}, nil
}

func newRoot(cwd string) *cobra.Command {
	opts := &runOpts{
		cwd:         cwd,
		failOn:      "error",
		color:       "auto",
		failOnFlag:  singleString{name: "fail-on"},
		rootFlag:    singleString{name: "root"},
		roadmapFlag: singleString{name: "roadmap"},
		colorFlag:   singleString{name: "color"},
		statusFlag:  statusList{name: "status"},
	}
	root := &cobra.Command{
		Use:           "spec",
		Short:         "Read-only status for a Spec Kit roadmap",
		Long:          rootLong(),
		Example:       rootExamples,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.version {
				fmt.Fprint(cmd.OutOrStdout(), versionText())
				return nil
			}
			return cmd.Help()
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true
	root.SetHelpTemplate(helpTemplate)
	root.Flags().BoolVarP(&opts.version, "version", "V", false, "Print build identity and exit")
	root.PersistentFlags().BoolVar(&opts.json, "json", false, "Print one JSON document on stdout")
	root.PersistentFlags().Var(&opts.failOnFlag, "fail-on", "Exit 1 at this severity: info, warning, error, or none")
	root.PersistentFlags().Var(&opts.rootFlag, "root", "Project directory that contains .specify")
	root.PersistentFlags().Var(&opts.roadmapFlag, "roadmap", "Roadmap file inside the selected root")
	root.PersistentFlags().Var(&opts.colorFlag, "color", "Human text color: auto, always, or never")
	root.PersistentFlags().Var(&opts.statusFlag, "status", "Exact raw status to show; repeat for any listed status")
	root.PersistentFlags().BoolVar(&opts.blocked, "blocked", false, "Show items with an unsatisfied dependency, including historical items")
	root.PersistentFlags().BoolVar(&opts.humanReview, "human-review", false, "Show items that need human review of a superseded dependency")
	root.PersistentFlags().BoolVar(&opts.active, "active", false, "Show active items")
	root.PersistentFlags().BoolVar(&opts.historical, "historical", false, "Show historical items, including valid supersession")
	root.PersistentFlags().BoolVar(&opts.explain, "explain", false, "Show mechanical file and line reasons for shown items")

	root.AddCommand(newStatus(opts), newInventory(opts), newVersion(), newCompletion())
	return root
}

type runOpts struct {
	cwd         string
	json        bool
	version     bool
	failOn      string
	color       string
	blocked     bool
	humanReview bool
	active      bool
	historical  bool
	explain     bool
	failOnFlag  singleString
	rootFlag    singleString
	roadmapFlag singleString
	colorFlag   singleString
	statusFlag  statusList
}

func (o *runOpts) prepare() {
	if o.failOnFlag.set {
		o.failOn = o.failOnFlag.value
	}
	if o.failOn == "" {
		o.failOn = "error"
	}
	if o.colorFlag.set {
		o.color = o.colorFlag.value
	}
	if o.color == "" {
		o.color = "auto"
	}
}

func newStatus(opts *runOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show recorded roadmap status",
		Long: `Read .specify/memory/roadmap.md and explicitly linked spec files.
Report recorded statuses, task checkboxes, dependency blockers, and broken references.
This command does not run tests or Git, and it does not change files.`,
		Example: `  spec status
  spec status --root testdata/example-roadmap
  spec status --json --fail-on warning
  spec status --color=never
  spec status --human-review --explain
  spec status --status planned --blocked`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.version {
				return errors.New("version cannot be combined with status")
			}
			opts.prepare()
			rep := status.Evaluate(status.Option{
				Root:    opts.rootFlag.value,
				Roadmap: opts.roadmapFlag.value,
				FailOn:  opts.failOn,
				CWD:     opts.cwd,
			})
			if opts.json {
				fmt.Fprint(cmd.OutOrStdout(), string(status.JSON(rep)))
				return silentExit{status.ExitCode(rep)}
			}
			view := status.View{
				Width: textColumns(), Color: false,
				Statuses: opts.statusFlag.values, Blocked: opts.blocked, HumanReview: opts.humanReview,
				Active: opts.active, Historical: opts.historical, Explain: opts.explain,
				Build: currentBuild().heading,
			}
			if rep.Outcome == "fatal" {
				fmt.Fprint(cmd.ErrOrStderr(), status.Text(rep, view))
				return silentExit{2}
			}
			view.Color = wantColor(cmd.OutOrStdout(), opts.color)
			fmt.Fprint(cmd.OutOrStdout(), status.Text(rep, view))
			code := status.ExitCode(rep)
			if code == 0 {
				return nil
			}
			return silentExit{code}
		},
	}
	cmd.ValidArgsFunction = cobra.NoFileCompletions
	return cmd
}

// inventoryRejected are roadmap-report options that spec inventory does not use.
var inventoryRejected = []string{"roadmap", "status", "blocked", "human-review", "active", "historical", "explain"}

func newInventory(opts *runOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inventory",
		Short: "List Spec Kit feature folders without a roadmap",
		Long: `List the feature folders directly under specs/, whether spec.md, plan.md and
tasks.md are present, and the task checkboxes in tasks.md.
This command does not read the roadmap and does not infer status, priority, or
dependencies. Checked boxes are not verification. It does not run tests or Git,
and it does not change files.`,
		Example: `  spec inventory
  spec inventory --json
  spec inventory --root path/to/project --fail-on warning`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.version {
				return errors.New("version cannot be combined with inventory")
			}
			for _, name := range inventoryRejected {
				if f := cmd.Flags().Lookup(name); f != nil && f.Changed {
					return errors.New("--" + name + " is not used by inventory")
				}
			}
			opts.prepare()
			rep := status.EvaluateInventory(status.InventoryOption{
				Root:   opts.rootFlag.value,
				FailOn: opts.failOn,
				CWD:    opts.cwd,
			})
			if opts.json {
				fmt.Fprint(cmd.OutOrStdout(), string(status.InventoryJSON(rep)))
				return silentExit{status.InventoryExitCode(rep)}
			}
			view := status.View{Width: textColumns(), Build: currentBuild().heading}
			if rep.Outcome == "fatal" {
				fmt.Fprint(cmd.ErrOrStderr(), status.InventoryText(rep, view))
				return silentExit{2}
			}
			view.Color = wantColor(cmd.OutOrStdout(), opts.color)
			fmt.Fprint(cmd.OutOrStdout(), status.InventoryText(rep, view))
			if code := status.InventoryExitCode(rep); code != 0 {
				return silentExit{code}
			}
			return nil
		},
	}
	cmd.ValidArgsFunction = cobra.NoFileCompletions
	return cmd
}

func newVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print build identity and exit",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprint(cmd.OutOrStdout(), versionText())
			return nil
		},
	}
}

func newCompletion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Print a shell completion script to stdout",
		Long: `Print a shell completion script to stdout. This command does not install the script or modify files.

The script is shell text from Cobra. spec does not run it. If you later source that script, the shell—not spec—may append to BASH_COMP_DEBUG_FILE when that variable is set in the shell.

  spec completion bash
  spec completion zsh
  spec completion fish
  spec completion powershell`,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fmt.Errorf("choose a shell: bash, zsh, fish, or powershell")
		},
	}
	addShell := func(name, short string, gen func(*cobra.Command) error) {
		cmd.AddCommand(&cobra.Command{
			Use:               name,
			Short:             short,
			Args:              cobra.NoArgs,
			ValidArgsFunction: cobra.NoFileCompletions,
			RunE: func(c *cobra.Command, _ []string) error {
				return gen(c)
			},
		})
	}
	addShell("bash", "Print a bash completion script to stdout", func(c *cobra.Command) error {
		return c.Root().GenBashCompletionV2(c.OutOrStdout(), true)
	})
	addShell("zsh", "Print a zsh completion script to stdout", func(c *cobra.Command) error {
		return c.Root().GenZshCompletion(c.OutOrStdout())
	})
	addShell("fish", "Print a fish completion script to stdout", func(c *cobra.Command) error {
		return c.Root().GenFishCompletion(c.OutOrStdout(), true)
	})
	addShell("powershell", "Print a PowerShell completion script to stdout", func(c *cobra.Command) error {
		return c.Root().GenPowerShellCompletionWithDesc(c.OutOrStdout())
	})
	return cmd
}

func wantColor(w io.Writer, mode string) bool {
	switch mode {
	case "always":
		return true
	case "never":
		return false
	default:
		if _, ok := os.LookupEnv("NO_COLOR"); ok {
			return false
		}
		if os.Getenv("TERM") == "dumb" {
			return false
		}
		return writerIsTerminal(w)
	}
}

func versionText() string {
	b := currentBuild()
	return fmt.Sprintf("spec %s\ncommit: %s\ngo: %s\ntarget: %s/%s\nlabel: %s\n",
		b.version, b.commit, runtime.Version(), runtime.GOOS, runtime.GOARCH, b.label)
}

type silentExit struct{ code int }

func (e silentExit) Error() string { return "" }

type singleString struct {
	name  string
	value string
	set   bool
}

func (s *singleString) Set(v string) error {
	if s.set {
		return fmt.Errorf("option --%s is repeated", s.name)
	}
	switch s.name {
	case "root", "roadmap":
		if v == "" {
			return fmt.Errorf("--%s needs a non-empty value", s.name)
		}
	case "fail-on":
		switch v {
		case "info", "warning", "error", "none":
		default:
			return fmt.Errorf("threshold %q is not info, warning, error, or none", v)
		}
	case "color":
		switch v {
		case "auto", "always", "never":
		default:
			return fmt.Errorf("color %q is not auto, always, or never", v)
		}
	}
	s.set = true
	s.value = v
	return nil
}

func (s *singleString) String() string { return s.value }
func (s *singleString) Type() string   { return "string" }

// Ensure the flag type satisfies pflag at compile time.
var _ pflag.Value = (*singleString)(nil)

type statusList struct {
	name   string
	values []string
}

func (s *statusList) Set(v string) error {
	if v == "" {
		return fmt.Errorf("option --%s requires a status token", s.name)
	}
	for _, existing := range s.values {
		if existing == v {
			return fmt.Errorf("option --%s %s is repeated", s.name, v)
		}
	}
	s.values = append(s.values, v)
	return nil
}

func (s *statusList) String() string { return strings.Join(s.values, ",") }
func (s *statusList) Type() string   { return "string" }

var _ pflag.Value = (*statusList)(nil)

func precheck(args []string) string {
	if flag := repeated(args); flag != "" {
		return "Option " + flag + " is repeated."
	}
	jsonMode := jsonRequested(args)
	help := countFlag(args, "--help")+countFlag(args, "-h") > 0
	version := countFlag(args, "--version")+countFlag(args, "-V") > 0
	statusCmd := hasCommand(args, "status")
	inventoryCmd := commandName(args) == "inventory"
	completion := hasCommand(args, "completion")
	versionCmd := hasCommand(args, "version")
	if help && (version || versionCmd) {
		return "Help and version cannot be combined."
	}
	if version && statusCmd {
		return "Version cannot be combined with status."
	}
	if version && inventoryCmd {
		return "Version cannot be combined with inventory."
	}
	if jsonMode && (help || version || versionCmd || completion) {
		return "JSON output cannot be combined with help, version, or completion."
	}
	if jsonMode && !statusCmd && !inventoryCmd {
		return "JSON output requires the status or inventory command."
	}
	return ""
}

func repeated(args []string) string {
	seen := map[string]int{}
	for _, a := range args {
		name := a
		if i := strings.IndexByte(name, '='); i >= 0 {
			name = name[:i]
		}
		switch name {
		case "--json", "--help", "-h", "--version", "-V", "--root", "--roadmap", "--fail-on", "--color", "--blocked", "--human-review", "--active", "--historical", "--explain":
			seen[name]++
			if seen[name] > 1 {
				return name
			}
		}
	}
	return ""
}

func countFlag(args []string, name string) int {
	n := 0
	for _, a := range args {
		got := a
		if i := strings.IndexByte(got, '='); i >= 0 {
			got = got[:i]
		}
		if got == name {
			n++
		}
	}
	return n
}

func hasCommand(args []string, name string) bool {
	for _, a := range args {
		if a == "--" {
			return false
		}
		if !strings.HasPrefix(a, "-") && a == name {
			return true
		}
	}
	return false
}

func oneLine(msg string) string {
	msg = strings.TrimSpace(msg)
	msg = strings.TrimPrefix(msg, "Error:")
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return "Invalid invocation."
	}
	return strings.Join(strings.Fields(msg), " ")
}

// failFor prints an invocation error. Under spec inventory --json it is an
// inventory-shaped document carrying the requested --fail-on when that value
// is valid (feature 006, research R9); otherwise the status-shaped document.
func failFor(args []string, jsonMode bool, code, message string, stdout, stderr io.Writer) int {
	if jsonMode && commandName(args) == "inventory" {
		fmt.Fprint(stdout, string(status.InventoryInvocationJSON(requestedFailOn(args), message)))
		return 2
	}
	return fail(jsonMode, requestedFailOn(args), code, message, stdout, stderr)
}

// valueFlags take a separate value token when written without "=", the way
// pflag consumes them.
var valueFlags = map[string]bool{"--root": true, "--roadmap": true, "--fail-on": true, "--color": true, "--status": true}

// walkArgs visits every flag with its value and returns the first positional
// token (the command), skipping the values of valueFlags, so a --root
// directory named "inventory" is not taken for the command. It stops at "--".
func walkArgs(args []string, flag func(name, value string, hasValue bool)) string {
	command := ""
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			break
		}
		if !strings.HasPrefix(a, "-") {
			if command == "" {
				command = a
			}
			continue
		}
		name, value, hasValue := a, "", false
		if j := strings.IndexByte(a, '='); j >= 0 {
			name, value, hasValue = a[:j], a[j+1:], true
		} else if valueFlags[a] && i+1 < len(args) {
			i++
			value, hasValue = args[i], true
		}
		if flag != nil {
			flag(name, value, hasValue)
		}
	}
	return command
}

// commandName is the first positional argument, or "".
func commandName(args []string) string { return walkArgs(args, nil) }

// requestedFailOn returns the first --fail-on value in args when it is
// valid, else "" (a missing or invalid value is not a requested threshold).
func requestedFailOn(args []string) string {
	got, seen := "", false
	walkArgs(args, func(name, value string, hasValue bool) {
		if name != "--fail-on" || seen {
			return
		}
		seen = true
		switch value {
		case "info", "warning", "error", "none":
			if hasValue {
				got = value
			}
		}
	})
	return got
}

// jsonRequested reports whether the last --json flag asks for JSON:
// "--json" and true values do, "--json=false" and false values do not, and
// an unparsable value does not (Cobra then reports it as plain text). A
// --root value spelled --json is skipped by walkArgs (feature 008).
func jsonRequested(args []string) bool {
	on := false
	walkArgs(args, func(name, value string, hasValue bool) {
		if name != "--json" {
			return
		}
		if !hasValue {
			on = true
			return
		}
		b, err := strconv.ParseBool(value)
		on = err == nil && b
	})
	return on
}

func fail(jsonMode bool, failOn, code, message string, stdout, stderr io.Writer) int {
	if jsonMode {
		fmt.Fprint(stdout, invocationJSON(failOn, code, message))
		return 2
	}
	fmt.Fprintf(stderr, "spec: %s: %s\n", code, message)
	return 2
}

func invocationJSON(failOn, code, message string) string {
	if failOn == "" {
		failOn = "error"
	}
	msg := jsonEscape(message)
	return `{"schema_version":"1.0","outcome":"fatal","complete":false,"source":{"roadmap":null,"input_format":null},"configuration":{"fail_on":"` + failOn + `","dependency_policy":"verified","task_scope":"unique_linked_task_files"},"summary":null,"items":[],"diagnostics":[{"code":"` + code + `","severity":"error","message":"` + msg + `","source":null,"item_id":null,"related":[]}],"error":{"code":"` + code + `","message":"` + msg + `","source":null}}` + "\n"
}

func jsonEscape(s string) string {
	b, err := json.Marshal(s)
	if err != nil || len(b) < 2 {
		return "invalid message"
	}
	return string(b[1 : len(b)-1])
}

// rootLong is the root help text; its first sentence of the second paragraph
// names the build kind (feature 004, contracts/build-identity.md).
func rootLong() string {
	return "Spec Status reads a Spec Kit roadmap and reports what the files record.\n\n" +
		currentBuild().helpSentence + " It does not change files, use the network, or infer quarters.\n" +
		"Exit 0 means no finding reached the threshold. It does not mean the work is done."
}

const rootExamples = `  spec status
  spec status --root testdata/example-roadmap
  spec status --json
  spec inventory
  spec completion bash
  spec version`

const helpTemplate = `{{with .Long}}{{. | trimTrailingWhitespaces}}

{{end}}Usage:
{{- if .Runnable}}
  {{.UseLine}}{{end}}
{{- if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}
{{- if .HasExample}}

Examples:
{{.Example}}{{end}}
{{- if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if or .IsAvailableCommand (eq .Name "help")}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}
{{- if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
{{- if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
{{- if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more about a command.{{end}}
`
