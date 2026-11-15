package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Feature 009 (roadmap 011): every README block of spec output is labeled
// with the command that produces it and must match real output byte for
// byte (specs/009-readme-sample-test/contracts/readme-samples.md).

const sampleLabelPrefix = "<!-- spec-sample:"

type readmeSample struct {
	labelLine int
	args      []string
	blockLine int
	body      string
}

// sampleFlags are the flags a README sample may use (research R3).
var sampleFlags = map[string]bool{
	"--color=never": true, "--root": true, "--status": true, "--blocked": true,
	"--human-review": true, "--active": true, "--historical": true, "--explain": true, "--fail-on": true,
}

var sampleValueFlags = map[string]bool{"--root": true, "--status": true, "--fail-on": true}

// parseReadmeSamples finds labeled samples and reports every rule violation.
func parseReadmeSamples(md string) ([]readmeSample, []string) {
	lines := strings.Split(md, "\n")
	var samples []readmeSample
	var problems []string
	labeled := map[int]bool{} // opening-fence line numbers that have a label
	inFence := false
	var pending *readmeSample
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		n := i + 1
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(line, "```") {
			if inFence {
				inFence = false
				continue
			}
			inFence = true
			info := strings.TrimSpace(strings.TrimPrefix(line, "```"))
			if pending != nil {
				if info != "" {
					problems = append(problems, fmt.Sprintf("README line %d: label is followed by a fenced block with info string %q; samples use a bare fence", pending.labelLine, info))
					pending = nil
					continue
				}
				var body strings.Builder
				j := i + 1
				for ; j < len(lines) && !strings.HasPrefix(lines[j], "```"); j++ {
					body.WriteString(lines[j])
					body.WriteByte('\n')
				}
				if j == len(lines) {
					problems = append(problems, fmt.Sprintf("README line %d: sample block is not closed", n))
				}
				pending.blockLine, pending.body = n, body.String()
				samples = append(samples, *pending)
				labeled[n] = true
				pending = nil
				continue
			}
			if info == "" && i+1 < len(lines) && !labeled[n] {
				first := lines[i+1]
				if strings.HasPrefix(first, "╭") || strings.HasPrefix(first, "spec: ") {
					problems = append(problems, fmt.Sprintf("README line %d: block looks like spec output but has no %s … --> label", n, sampleLabelPrefix))
				}
			}
			continue
		}
		if inFence {
			continue
		}
		if strings.HasPrefix(trim, sampleLabelPrefix) {
			if pending != nil {
				problems = append(problems, fmt.Sprintf("README line %d: label has no fenced block after it", pending.labelLine))
			}
			if !strings.HasSuffix(trim, "-->") {
				problems = append(problems, fmt.Sprintf("README line %d: label is not closed with -->", n))
				pending = nil
				continue
			}
			cmd := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trim, sampleLabelPrefix), "-->"))
			if msg := checkSampleArgs(cmd); msg != "" {
				problems = append(problems, fmt.Sprintf("README line %d: %s", n, msg))
				pending = nil
				continue
			}
			pending = &readmeSample{labelLine: n, args: strings.Fields(cmd)}
			continue
		}
		if pending != nil && trim != "" {
			problems = append(problems, fmt.Sprintf("README line %d: label has no fenced block after it", pending.labelLine))
			pending = nil
		}
	}
	if pending != nil {
		problems = append(problems, fmt.Sprintf("README line %d: label has no fenced block after it", pending.labelLine))
	}
	return samples, problems
}

// checkSampleArgs enforces the label safety rules (research R3).
func checkSampleArgs(cmd string) string {
	if strings.ContainsAny(cmd, "\"'`") {
		return "label must not contain quote characters"
	}
	args := strings.Fields(cmd)
	if len(args) == 0 || (args[0] != "status" && args[0] != "inventory") {
		return "label must start with status or inventory"
	}
	root := ""
	for i := 1; i < len(args); i++ {
		a := args[i]
		if !sampleFlags[a] {
			return "unsupported option " + a + " in label"
		}
		if sampleValueFlags[a] {
			if i+1 >= len(args) {
				return a + " needs a value in label"
			}
			if a == "--root" {
				root = args[i+1]
			}
			i++
		}
	}
	if root == "" {
		return "label needs --root testdata/<fixture>"
	}
	if !strings.HasPrefix(root, "testdata/") || strings.Contains(root, "..") {
		return "label root must be under testdata/ without .."
	}
	return ""
}

func TestReadmeSamples(t *testing.T) {
	repo, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(repo, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	samples, problems := parseReadmeSamples(string(raw))
	for _, p := range problems {
		t.Error(p)
	}
	if len(samples) == 0 {
		t.Fatal("README has no labeled spec output samples")
	}
	withBuild(t, resolveBuild("", "0.0.0-dev", "unknown", nil, false))
	t.Setenv("COLUMNS", "80")
	for _, s := range samples {
		var out, errb bytes.Buffer
		code := Run(s.args, repo, &out, &errb)
		got := out.String()
		switch {
		case code == 2:
			// Fatal runs print plain text to stderr only.
			got = errb.String()
		case errb.Len() != 0:
			t.Errorf("README line %d (spec %s): unexpected stderr on exit %d: %q", s.blockLine, strings.Join(s.args, " "), code, errb.String())
		}
		if strings.TrimSpace(got) == "" || strings.TrimSpace(s.body) == "" {
			// Never let an empty sample match empty output.
			t.Errorf("README line %d (spec %s): empty sample or empty output (exit %d)", s.blockLine, strings.Join(s.args, " "), code)
			continue
		}
		if got != s.body {
			t.Errorf("README line %d (spec %s) differs from real output\n--- README\n%s--- spec output\n%s", s.blockLine, strings.Join(s.args, " "), s.body, got)
		}
	}
}

func TestReadmeSampleGuards(t *testing.T) {
	block := "```\n╭──╮\n```\n"
	cases := map[string]string{
		"unlabeled output":      "text\n\n" + block,
		"unlabeled fatal":       "```\nspec: SS_X: m\n```\n",
		"label without block":   "<!-- spec-sample: status --root testdata/x -->\n\nprose\n",
		"label at end":          "<!-- spec-sample: status --root testdata/x -->\n",
		"bad subcommand":        "<!-- spec-sample: version --root testdata/x -->\n" + block,
		"missing root":          "<!-- spec-sample: status --color=never -->\n" + block,
		"root outside testdata": "<!-- spec-sample: status --root specs/x -->\n" + block,
		"dot-dot root":          "<!-- spec-sample: status --root testdata/../x -->\n" + block,
		"unsupported flag":      "<!-- spec-sample: status --json --root testdata/x -->\n" + block,
		"quote in label":        "<!-- spec-sample: status --root \"testdata/x\" -->\n" + block,
		"info string after":     "<!-- spec-sample: status --root testdata/x -->\n```text\n╭\n```\n",
	}
	for name, md := range cases {
		if _, problems := parseReadmeSamples(md); len(problems) == 0 {
			t.Errorf("%s: no problem reported for:\n%s", name, md)
		}
	}
	ok := "<!-- spec-sample: status --color=never --root testdata/example-roadmap --human-review --explain -->\n\n```\n╭──╮\n```\n\n```sh\nspec status\n```\n"
	samples, problems := parseReadmeSamples(ok)
	if len(problems) != 0 || len(samples) != 1 || samples[0].body != "╭──╮\n" || len(samples[0].args) != 6 {
		t.Fatalf("valid sample: %+v %v", samples, problems)
	}
}
