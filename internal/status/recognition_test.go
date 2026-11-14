package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Feature 006 US1: roadmap recognition (contracts/roadmap-recognition.md).
// Expectations are hand-written in testdata/two-mode/*/expected-status.json.

type recognitionExpect struct {
	Exit              int     `json:"exit"`
	Code              string  `json:"code"`
	Message           *string `json:"message"`
	TotalItems        *int    `json:"total_items"`
	ExplicitEmptyLine *bool   `json:"explicit_empty_line"`
}

const explicitEmptyLine = "roadmap explicitly empty (Planned Specs heading, no items)"

func twoModeRoot(name string) string {
	return filepath.Join("..", "..", "testdata", "two-mode", name)
}

func TestRoadmapRecognitionFixtures(t *testing.T) {
	for _, name := range []string{
		"unrecognized-empty", "unrecognized-whitespace", "unrecognized-narrative",
		"planned-in-fence", "planned-wrong-form", "explicit-empty", "items-no-planned",
		"narrative-with-label", "no-roadmap",
	} {
		t.Run(name, func(t *testing.T) {
			root := twoModeRoot(name)
			raw, err := os.ReadFile(filepath.Join(root, "expected-status.json"))
			if err != nil {
				t.Fatal(err)
			}
			var want recognitionExpect
			if err := json.Unmarshal(raw, &want); err != nil {
				t.Fatal(err)
			}
			rep := Evaluate(Option{Root: root})
			if got := ExitCode(rep); got != want.Exit {
				t.Fatalf("exit %d, want %d (fatal %+v)", got, want.Exit, rep.Fatal)
			}
			var doc map[string]any
			if err := json.Unmarshal(JSON(rep), &doc); err != nil {
				t.Fatal(err)
			}
			if want.Code != "" {
				if rep.Fatal == nil || rep.Fatal.Code != want.Code {
					t.Fatalf("fatal %+v, want %s", rep.Fatal, want.Code)
				}
				if want.Message != nil && rep.Fatal.Message != *want.Message {
					t.Fatalf("message %q, want %q", rep.Fatal.Message, *want.Message)
				}
				errDoc, _ := doc["error"].(map[string]any)
				if errDoc == nil || errDoc["code"] != want.Code {
					t.Fatalf("JSON error %v", doc["error"])
				}
				text := Text(rep, View{Width: 80})
				if !strings.HasPrefix(text, "spec: "+want.Code+": ") {
					t.Fatalf("fatal text %q", text)
				}
				// A fatal is never masked by --fail-on none.
				none := Evaluate(Option{Root: root, FailOn: "none"})
				if ExitCode(none) != 2 {
					t.Fatal("--fail-on none masked a fatal")
				}
				return
			}
			if rep.Fatal != nil {
				t.Fatalf("unexpected fatal %+v", rep.Fatal)
			}
			summary, _ := doc["summary"].(map[string]any)
			if want.TotalItems != nil && int(summary["total_items"].(float64)) != *want.TotalItems {
				t.Fatalf("total_items %v, want %d", summary["total_items"], *want.TotalItems)
			}
			text := Text(rep, View{Width: 80})
			if want.ExplicitEmptyLine != nil && strings.Contains(text, explicitEmptyLine) != *want.ExplicitEmptyLine {
				t.Fatalf("explicit-empty line present=%v, want %v:\n%s", !*want.ExplicitEmptyLine, *want.ExplicitEmptyLine, text)
			}
		})
	}
}

func TestExplicitEmptyRoadmapHasNullDenominators(t *testing.T) {
	root := twoModeRoot("explicit-empty")
	rep := Evaluate(Option{Root: root})
	if rep.Outcome != "report" || !rep.Complete || ExitCode(rep) != 0 {
		t.Fatalf("explicit empty: %+v exit %d", rep.Fatal, ExitCode(rep))
	}
	if rep.Summary.VerifiedPercent.Text() != "N/A" || rep.Summary.Tasks.Percent.Text() != "N/A" {
		t.Fatalf("denominators must be N/A, got %s %s", rep.Summary.VerifiedPercent.Text(), rep.Summary.Tasks.Percent.Text())
	}
	got := string(JSON(rep))
	if strings.Contains(got, "quarter") || strings.Contains(got, "as_of") || strings.Contains(got, "inferred") {
		t.Fatalf("unexpected timing fields: %s", got)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(doc["summary"])
	if strings.Contains(string(raw), "100") {
		t.Fatalf("empty summary must not become 100: %s", raw)
	}
	if again := string(JSON(Evaluate(Option{Root: root}))); got != again {
		t.Fatal("JSON is not deterministic")
	}
}

func TestUnrecognizedThroughExplicitRoadmapFlag(t *testing.T) {
	dir := project(t, "### 001 — Real [status: planned]\n- Description: x\n- Outcome: y\n- Scope: z\n", map[string]string{
		"docs/notes.md": "# Notes\n\nJust prose.\n",
	})
	rep := Evaluate(Option{Root: dir, CWD: dir, Roadmap: "docs/notes.md"})
	if rep.Fatal == nil || rep.Fatal.Code != "SS_ROADMAP_UNRECOGNIZED" || ExitCode(rep) != 2 {
		t.Fatalf("--roadmap unrecognized: %+v", rep.Fatal)
	}
}

func TestPlannedSpecsHeadingTrailingWhitespace(t *testing.T) {
	dir := project(t, "# R\n\n## Planned Specs \t\n", nil)
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if rep.Fatal != nil || ExitCode(rep) != 0 {
		t.Fatalf("trailing whitespace heading not recognized: %+v", rep.Fatal)
	}
}

func TestItemsHideExplicitEmptyLine(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "example-roadmap")
	text := Text(Evaluate(Option{Root: root}), View{Width: 80})
	if strings.Contains(text, explicitEmptyLine) {
		t.Fatal("explicit-empty line shown for a roadmap with items")
	}
}
