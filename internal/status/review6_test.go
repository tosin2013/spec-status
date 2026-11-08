package status

import (
	"reflect"
	"testing"
)

// Cases from the sixth acceptance review of 001 (7 November 2026).

func TestWrappedFieldValuesAreJoined(t *testing.T) { // D1
	body := "### 002 — Two [status: verified]\n### 003 — Three [status: verified]\n### 004 — Four [status: undecided]\n" +
		"### 005 — Five [status: planned]\n- Depends on: 002 (verified), 003\n  (verified), 004 (undecided).\n" +
		"- Description: first line\n  second line\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	it := itemByID(t, rep, "005")
	if got := depIDs(it); !reflect.DeepEqual(got, []string{"002", "003", "004"}) {
		t.Fatalf("deps %v", got)
	}
	if rep.Summary.BlockedCount != 1 || it.Description != "first line second line" {
		t.Fatalf("blocked %d description %q", rep.Summary.BlockedCount, it.Description)
	}
	if code := fatalCode(t, base+"- Depends on:\n  - 002\n### 002 — B [status: planned]\n", nil); code != "SS_ITEM_MALFORMED" {
		t.Fatalf("list under Depends on: %s", code)
	}
}

func TestHTMLBlocksHideStructure(t *testing.T) { // D2, D5
	tasks := "<details>\n<summary>Phase 1 (done)</summary>\n- [x] T001 a\n\n</details>\n\n- [ ] T003 c\n"
	if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("tasks in HTML block: %s", code)
	}
	ok := "<details>\n<summary>Phase 1 (done)</summary>\n\n- [x] T001 a\n\n</details>\n\n- [ ] T003 c\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": ok})})
	if rep.Outcome != "report" || rep.Summary.Tasks.Checked != 1 || rep.Summary.Tasks.Total != 2 {
		t.Fatalf("details with blank line: %s %+v", rep.Outcome, rep.Summary)
	}
	for _, opener := range []string{"<a id=\"i4\"/>", "<div>", "</details>"} {
		body := opener + "\n### 004 — Hidden [status: verified]\n\n### 005 — Five [status: planned]\n- Depends on: 004\n"
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("roadmap %q: %s", opener, code)
		}
	}
	anchor := "<a id=\"i4\"></a>\n\n" + base
	if rep := Evaluate(Option{Root: project(t, anchor, nil), FailOn: "none"}); rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("anchor with blank line: %+v", rep.Fatal)
	}
}

func TestBlockConstructsCloseListItems(t *testing.T) { // D3
	for _, tasks := range []string{
		"- [x] a\n***\n    - [ ] b\n",
		"- [x] a\n---\n    - [ ] b\n",
		"- [x] a\n```\nx\n```\n    - [ ] b\n",
		"- [x] a\n> q\n    - [x] e\n",
		"   - [ ] f\n    - [x] e\n",
	} {
		if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("%q: %s", tasks, code)
		}
	}
}

func TestRemainingGFMListRules(t *testing.T) { // D4
	for _, tasks := range []string{
		"Text\n2. [x] b\n",
		"- [ ] T001\n  2. [x] sub\n",
		"-     [x] a\n",
		"- - [x] a\n",
	} {
		if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("%q: %s", tasks, code)
		}
	}
	for tasks, want := range map[string][2]int{
		"Text\n1. [x] b\n2. [ ] c\n":               {1, 2},
		"1. [x] a\n2. [ ] b\n3. [ ] c\n":           {1, 3},
		"- [ ] T001\n  1. [x] sub\n  2. [ ] sub\n": {1, 3},
	} {
		rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks})})
		if rep.Outcome != "report" || rep.Summary.Tasks.Checked != want[0] || rep.Summary.Tasks.Total != want[1] {
			t.Fatalf("%q: %s %+v, want %d/%d", tasks, rep.Outcome, rep.Summary, want[0], want[1])
		}
	}
}

func TestTasksFatalCitesTasksFile(t *testing.T) {
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": "- [x] a\n```\n"}), FailOn: "none"})
	if rep.Outcome != "fatal" || rep.Fatal.Path != "specs/a/tasks.md" || rep.Fatal.Line != 2 {
		t.Fatalf("fatal %+v", rep.Fatal)
	}
}
