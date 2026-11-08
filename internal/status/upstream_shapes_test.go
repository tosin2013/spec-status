package status

import (
	"os"
	"path/filepath"
	"testing"
)

// Edge cases named in specs/002-upstream-format-conformance/spec.md for the
// pinned Spec Kit tasks template (github/spec-kit v1.2.0).

func TestTasksFrontMatterCreatesNoTasks(t *testing.T) {
	tasks := "---\n\ndescription: \"Task list template for feature implementation\"\n---\n\n# Tasks: X\n\n- [ ] T001 One\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks})})
	if rep.Outcome != "report" || rep.Summary.Tasks.Total != 1 {
		t.Fatalf("front matter changed the count: %s %+v", rep.Outcome, rep.Summary)
	}
}

func TestTaskMarkersAreNotLinks(t *testing.T) {
	tasks := "- [ ] T010 [P] [US1] Contract test for [endpoint] in tests/contract/test_[name].py\n" +
		"- [x] T012 [P] [US1] Create [Entity1] model in src/models/[entity1].py\n" +
		"- [ ] T014 [US1] Implement [Service] in src/services/[service].py (depends on T012, T013)\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}), FailOn: "info"})
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_REFERENCE_UNSUPPORTED" || d.Code == "SS_TASK_TARGET_MISSING" {
			t.Fatalf("task marker read as a link: %+v", d)
		}
	}
	if rep.Summary.Tasks.Checked != 1 || rep.Summary.Tasks.Total != 3 {
		t.Fatalf("tasks %+v", rep.Summary.Tasks)
	}
	withDef := tasks + "\n[P]: parallel.md\n"
	rep = Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": withDef}), FailOn: "none"})
	if !hasDiag(rep, "001", "SS_REFERENCE_UNSUPPORTED") {
		t.Fatal("a [P] marker with a matching reference definition must be reported")
	}
}

// The pinned template itself, unfilled, must read without a fatal and count
// every task line.
func TestPinnedSpecKitTasksTemplateCounts(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "upstream", "spec-kit", "v1.2.0", "templates", "tasks-template.md"))
	if err != nil {
		t.Fatal(err)
	}
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": string(data)}), FailOn: "none"})
	if rep.Outcome != "report" || rep.Summary.Tasks.Total != 34 || rep.Summary.Tasks.Checked != 0 {
		t.Fatalf("pinned tasks template: %s %+v %+v", rep.Outcome, rep.Fatal, rep.Summary)
	}
}

// Review 5 N2: a field label spelled with Cyrillic or Greek lookalike
// letters must not slip past the metric-field rules (D-2026-11-07-09 (b)).
func TestLetterLookalikesInMetricLabels(t *testing.T) {
	for _, line := range []string{
		"  - Dеpends on: 009",
		"- Spеc dir: specs/x",
		"  supеrsеdеd-by: 002",
		"- Ѕpec dir: specs/x",
	} {
		body := base + "- Description: x\n" + line + "\n### 002 — B [status: planned]\n"
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: %s", line, code)
		}
	}
	title := "### 001 — Проект на русском [status: planned]\n- Description: Описание.\n"
	if rep := Evaluate(Option{Root: project(t, title, nil), FailOn: "none"}); rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("Cyrillic prose rejected: %+v", rep.Fatal)
	}
}
