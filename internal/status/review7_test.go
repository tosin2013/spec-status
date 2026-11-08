package status

import "testing"

// Cases from the seventh acceptance review of 001 (7 November 2026). Each
// produced a count that renderers do not show, with no diagnostic; each now
// fails visibly.

func TestLoneTagAfterHeadingStartsHTMLBlock(t *testing.T) { // B1
	roadmap := "## Phase 2\n<img src=\"phase2.png\">\n### 002 — Billing [status: verified]\n- Spec dir: specs/b\n\n" +
		"### 003 — Reports [status: planned]\n- Depends on: 002\n"
	if code := fatalCode(t, roadmap, nil); code != "SS_ITEM_MALFORMED" {
		t.Fatalf("roadmap: %s", code)
	}
	tasks := "## Phase 1\n<a id=\"phase-1\"/>\n- [x] T001\n- [x] T002\n"
	if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("tasks: %s", code)
	}
	ok := "## Phase 1\n\n<a id=\"phase-1\"></a>\n\n- [x] T001\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": ok})})
	if rep.Outcome != "report" || rep.Summary.Tasks.Total != 1 {
		t.Fatalf("anchor with blank lines: %s %+v", rep.Outcome, rep.Fatal)
	}
}

func TestParagraphAfterFenceClosesTask(t *testing.T) { // B2
	tasks := "- [ ] T001 Run the migration\n  " + fence + "\n  make migrate\n  " + fence + "\nThen verify:\n\n    - [x] T002 Check row counts\n"
	if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("code %s", code)
	}
	ok := "- [ ] T001 Run the migration\n  " + fence + "\n  make migrate\n  " + fence + "\n  Then verify the counts.\n- [x] T002 Check row counts\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": ok})})
	if rep.Outcome != "report" || rep.Summary.Tasks.Checked != 1 || rep.Summary.Tasks.Total != 2 {
		t.Fatalf("indented follow-up text: %s %+v", rep.Outcome, rep.Summary)
	}
}

func TestOtherHTMLBlockTypesAreFatal(t *testing.T) { // B3
	for _, opener := range []string{"<?", "<?xml version=\"1.0\"?>", "<![CDATA[", "<!DOCTYPE html>"} {
		if code := fatalCode(t, opener+"\n\n### 002 — Hidden [status: verified]\n- Spec dir: specs/b\n", nil); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("roadmap %q: %s", opener, code)
		}
		if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": opener + "\n\n- [x] t\n"}); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("tasks %q: %s", opener, code)
		}
	}
}

func TestTabsAdvanceToTabStops(t *testing.T) { // B4
	if w, _ := leadingIndent(" \t- x"); w != 4 {
		t.Fatalf("space+tab width %d, want 4", w)
	}
	if w, _ := leadingIndent("  \t\tx"); w != 8 {
		t.Fatalf("two spaces+two tabs width %d, want 8", w)
	}
	tasks := "  2. [x] Parent\n \t- [ ] Child\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}), FailOn: "none"})
	if rep.Outcome == "report" && rep.Summary.Tasks.Total == 2 {
		t.Fatalf("tab-indented child counted as a nested task: %+v", rep.Summary.Tasks)
	}
}

func TestMarkerClosesItemsItDoesNotReach(t *testing.T) { // B5
	tasks := "- [x] T001 a\n 1.  [ ] T002 b\n    - [x] T003 c\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}), FailOn: "none"})
	if rep.Outcome == "report" && rep.Summary.Tasks.Total == 3 {
		t.Fatalf("stale item stack counted T003: %+v", rep.Summary.Tasks)
	}
	nested := "1. [ ] a\n   - [x] b\n- [ ] c\n  - [x] d\n"
	rep = Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": nested})})
	if rep.Outcome != "report" || rep.Summary.Tasks.Checked != 2 || rep.Summary.Tasks.Total != 4 {
		t.Fatalf("ordinary nesting: %s %+v", rep.Outcome, rep.Summary)
	}
}

func TestBareMetricLabelUnderItemIsFatal(t *testing.T) { // review observation
	for _, line := range []string{"superseded-by: 002", "  superseded-by: 002", "   Depends on: 002", "  **Spec dir:** specs/x"} {
		body := base + "- Description: x\n\n" + line + "\n### 002 — B [status: verified]\n"
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: %s", line, code)
		}
	}
	prose := base + "- Notes: this work\n  depends on the vendor timeline.\n"
	if rep := Evaluate(Option{Root: project(t, prose, nil), FailOn: "none"}); rep.Outcome != "report" {
		t.Fatalf("prose continuation rejected: %+v", rep.Fatal)
	}
}
