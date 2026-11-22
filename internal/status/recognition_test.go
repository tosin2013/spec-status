package status

import (
	"encoding/json"
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
