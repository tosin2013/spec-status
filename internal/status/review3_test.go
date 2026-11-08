package status

import (
	"os"
	"path/filepath"
	"testing"
)

// Task-file cases from the third acceptance review of 001 (7 November 2026).

const oneSpec = "### 001 — A [status: planned]\n- Spec dir: specs/a\n"

func TestCheckboxSpellingsLikeGFM(t *testing.T) { // D3
	tasks := "-  [x] two spaces\n-\t[x] tab\n- [ ] parent\n    ```\n    - [x] example\n    ```\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks})})
	if rep.Outcome != "report" || rep.Summary.Tasks.Total != 3 || rep.Summary.Tasks.Checked != 2 {
		t.Fatalf("tasks %+v (outcome %s), want 2/3", rep.Summary, rep.Outcome)
	}
	// Narrow tasks dialect: block-quoted checkboxes are fatal, not counted.
	for _, q := range []string{"> - [x] quoted\n", "> > * [ ] nested quote\n"} {
		if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": q}); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("%q: code %s", q, code)
		}
	}
}

func TestUncheckableTaskLinksAreReported(t *testing.T) { // D4
	tasks := "- [x] Write [the `[x]` parser](parser.md)\n" +
		"- [x] Write [docs]\n" +
		"- [x] Wrapped task that names\n  [a file](wrapped.md) on the next line\n" +
		"\n[docs]: docs.md\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}), FailOn: "none"})
	byLine := map[int]bool{}
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_REFERENCE_UNSUPPORTED" {
			byLine[d.Line] = true
		}
	}
	for _, line := range []int{1, 2, 4} {
		if !byLine[line] {
			t.Fatalf("line %d link not reported: %v", line, diagCodesAll(rep))
		}
	}
	plain := "- [x] Uses [brackets] as text and [x] marks\n"
	rep = Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": plain}), FailOn: "none"})
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_REFERENCE_UNSUPPORTED" {
			t.Fatalf("plain brackets reported as a link: %+v", d)
		}
	}
}

func TestOneTasksFileUnderTwoPathsCountsOnce(t *testing.T) { // D5
	body := "### 001 — A [status: planned]\n- Spec dir: specs/a\n### 002 — B [status: planned]\n- Spec dir: specs/b\n"
	dir := project(t, body, map[string]string{"specs/a/tasks.md": "- [x] one\n- [ ] two\n"})
	if err := os.MkdirAll(filepath.Join(dir, "specs", "b"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(filepath.Join(dir, "specs", "a", "tasks.md"), filepath.Join(dir, "specs", "b", "tasks.md")); err != nil {
		t.Skipf("hard links unavailable: %v", err)
	}
	rep := Evaluate(Option{Root: dir})
	if rep.Summary.Tasks.Total != 2 || rep.Summary.Tasks.FilesRead != 1 {
		t.Fatalf("hard-linked tasks file counted twice: %+v", rep.Summary.Tasks)
	}
}
