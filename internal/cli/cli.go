// Package cli is the spec command surface.
// There is no config file, network client, or write command.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/tosin2013/spec-status/internal/status"
)

// Overridden with -ldflags -X when a builder stamps a commit.
// Defaults describe a development build, not a release.
var (
	Version = "0.0.0-dev"
	Commit  = "unknown"
)

// Run executes one invocation and returns the process exit code.
func Run(args []string, cwd string, stdout, stderr io.Writer) int {
	jsonMode := false
	for _, a := range args {
		if a == "--json" {
			jsonMode = true
			break
		}
	}
	cmd, opt, err := parse(args)
	if err != nil {
		return fail(jsonMode, "SS_INVOCATION_INVALID", err.Error(), stdout, stderr)
	}
	switch cmd {
	case "help":
		fmt.Fprint(stdout, helpText)
		return 0
	case "version":
		fmt.Fprintf(stdout, "spec %s\ncommit: %s\ngo: %s\ntarget: %s/%s\nlabel: development build, not a release\n",
			Version, Commit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return 0
	}
	rep := status.Evaluate(status.Option{Root: opt.root, Roadmap: opt.roadmap, FailOn: opt.failOn, CWD: cwd})
	if jsonMode {
		fmt.Fprint(stdout, string(status.JSON(rep)))
		return status.ExitCode(rep)
	}
	if rep.Outcome == "fatal" {
		fmt.Fprint(stderr, status.Text(rep))
		return 2
	}
	fmt.Fprint(stdout, status.Text(rep))
	return status.ExitCode(rep)
}

type options struct {
	root, roadmap, failOn string
}

func parse(args []string) (string, options, error) {
	if len(args) == 0 {
		return "help", options{}, nil
	}
	var opt options
	opt.failOn = "error"
	mode := ""
	seen := map[string]string{}
	take := func(name, val string) error {
		if prev, ok := seen[name]; ok {
			if prev != val {
				return fmt.Errorf("Option %s is repeated.", name)
			}
			return fmt.Errorf("Option %s is repeated.", name)
		}
		seen[name] = val
		return nil
	}
	i := 0
	if args[0] != "" && !strings.HasPrefix(args[0], "-") {
		mode = args[0]
		i = 1
	}
	help, version := false, false
	for i < len(args) {
		a := args[i]
		switch {
		case a == "--help" || a == "-h":
			help = true
			i++
		case a == "--version" || a == "-V":
			version = true
			i++
		case a == "--json":
			i++
		case a == "--fail-on" || a == "--root" || a == "--roadmap":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				return "", options{}, fmt.Errorf("Option %s needs a value.", a)
			}
			val := args[i+1]
			if err := take(a, val); err != nil {
				return "", options{}, err
			}
			switch a {
			case "--fail-on":
				switch val {
				case "info", "warning", "error", "none":
					opt.failOn = val
				default:
					return "", options{}, fmt.Errorf("Threshold %q is not info, warning, error, or none.", val)
				}
			case "--root":
				opt.root = val
			case "--roadmap":
				opt.roadmap = val
			}
			i += 2
		default:
			if strings.HasPrefix(a, "-") {
				return "", options{}, fmt.Errorf("Unknown option %s.", a)
			}
			return "", options{}, fmt.Errorf("Unknown argument %s.", a)
		}
	}
	if help && version {
		return "", options{}, fmt.Errorf("Help and version cannot be combined.")
	}
	if (help || version) && hasJSON(args) {
		return "", options{}, fmt.Errorf("Help and version cannot be combined with --json.")
	}
	if mode == "" {
		if help {
			return "help", opt, nil
		}
		if version {
			return "version", opt, nil
		}
		return "", options{}, fmt.Errorf("Usage is spec status. Run spec --help.")
	}
	if mode != "status" {
		return "", options{}, fmt.Errorf("Unknown command %s.", mode)
	}
	if help && (opt.root != "" || opt.roadmap != "" || seen["--fail-on"] != "" || hasJSON(args)) {
		return "", options{}, fmt.Errorf("Help cannot be combined with status options.")
	}
	if help {
		return "help", opt, nil
	}
	if version {
		return "", options{}, fmt.Errorf("Version cannot be combined with status.")
	}
	return "status", opt, nil
}

func hasJSON(args []string) bool {
	for _, a := range args {
		if a == "--json" {
			return true
		}
	}
	return false
}

func fail(jsonMode bool, code, message string, stdout, stderr io.Writer) int {
	if jsonMode {
		fmt.Fprint(stdout, invocationJSON(code, message))
		return 2
	}
	fmt.Fprintf(stderr, "spec: %s: %s\n", code, message)
	return 2
}

func invocationJSON(code, message string) string {
	// Manual document so invocation failures do not depend on a project root.
	msg := jsonEscape(message)
	return `{"schema_version":"1.0","outcome":"fatal","complete":false,"source":{"roadmap":null,"input_format":null},"configuration":{"fail_on":"error","dependency_policy":"verified","task_scope":"unique_linked_task_files"},"summary":null,"items":[],"diagnostics":[{"code":"` + code + `","severity":"error","message":"` + msg + `","source":null,"item_id":null,"related":[]}],"error":{"code":"` + code + `","message":"` + msg + `","source":null}}` + "\n"
}

func jsonEscape(s string) string {
	b, err := json.Marshal(s)
	if err != nil || len(b) < 2 {
		return "invalid message"
	}
	return string(b[1 : len(b)-1])
}

const helpText = `spec — read-only status for a Spec Kit roadmap (development build)

Usage:
  spec status [--json] [--fail-on info|warning|error|none] [--root PATH] [--roadmap PATH]
  spec --help
  spec --version

spec status reads .specify/memory/roadmap.md under the selected project root.
It reports recorded statuses, task checkboxes, dependency blockers, and broken references.
It does not run tests, Git, or the network, and it does not change files.

The default root is the nearest ancestor directory that contains .specify.
--root is relative to the current directory. --roadmap is relative to that root.

Exit codes:
  0  no finding at or above the threshold (not proof the work is done)
  1  findings at or above the threshold
  2  invocation, configuration, required input, or fatal parse failure

The default threshold is error. Findings are still printed at every threshold.
Inferred quarters are not calculated. There is no release binary in this build.
`
