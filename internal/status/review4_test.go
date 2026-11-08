package status

import "testing"

// Cases from the fourth acceptance review of 001 (7 November 2026): fence and
// comment boundaries that differed from CommonMark. Where the parser and a
// renderer could disagree, the report now fails visibly.

const fence = "```"

func TestFenceCloserMustMatchOpenerIndent(t *testing.T) { // D1
	body := "### 001 — Real [status: planned]\n\n" + fence + "markdown\n    " + fence + "\n### 002 — Example [status: verified]\n" + fence + "\n\nNotes.\n    " + fence + "\n"
	if code := fatalCode(t, body, nil); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("code %s", code)
	}
	nested := "````markdown\n" + fence + "\n### 999 — Example [status: verified]\n" + fence + "\n````\n" + base
	rep := Evaluate(Option{Root: project(t, nested, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("longer outer fence: %s items=%d", rep.Outcome, len(rep.Items))
	}
}

func TestCommentBoundaries(t *testing.T) { // D2, D3, N1
	inItem := base + "- Description: x\n  <!-- note\n\n### 002 — Next [status: planned]\n  -->\n"
	if code := fatalCode(t, inItem, nil); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("comment in list item: %s", code)
	}
	indented := "    <!--\n" + base + "    -->\n"
	rep := Evaluate(Option{Root: project(t, indented, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("indented <!-- hid an item: %s items=%d", rep.Outcome, len(rep.Items))
	}
	trailing := base + "<!-- n --> - Depends on: 002\n### 002 — B [status: planned]\n"
	if code := fatalCode(t, trailing, nil); code != "SS_ITEM_MALFORMED" {
		t.Fatalf("text after comment: %s", code)
	}
	top := "<!--\n### 002 — Draft [status: verified]\n\nstill a comment\n-->\n" + base
	rep = Evaluate(Option{Root: project(t, top, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("top-level comment: %s items=%d", rep.Outcome, len(rep.Items))
	}
}

func TestTaskFenceAndCommentBoundaries(t *testing.T) { // D4, D5
	closer := "- [ ] T001\n" + fence + "\n    " + fence + "\n- [x] a\n- [x] b\n    " + fence + "\n"
	if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": closer}); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("task fence closer: %s", code)
	}
	comment := "- [ ] T001\n  <!-- old\n- [x] T002\n- [x] T003\n  -->\n"
	if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": comment}); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("task comment in list item: %s", code)
	}
}

func TestIndentedCodeInTasksIsNotTasks(t *testing.T) { // D6, N3
	tasks := "Example:\n\n    - [x] T999 sample\n    " + fence + "\n\n- [ ] T001\n  - [x] nested child\n\n    - [ ] continued in the list\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks})})
	if rep.Outcome != "report" || rep.Summary.Tasks.Total != 3 || rep.Summary.Tasks.Checked != 1 {
		t.Fatalf("tasks %+v (outcome %s), want 1/3", rep.Summary, rep.Outcome)
	}
}
