package status

import "testing"

// Cases from the fifth acceptance review of 001 (7 November 2026). Owner
// decision: a narrow tasks dialect that measures indentation from the
// enclosing list item and fails visibly where container rules would decide.

func TestEmptyCommentsAreFatal(t *testing.T) { // B1
	for _, c := range []string{"<!-->", "<!--->"} {
		if code := fatalCode(t, c+"\n\n"+base+"\n<!-- end -->\n", nil); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("roadmap %q: %s", c, code)
		}
		tasks := "- [ ] real\n" + c + "\n- [x] b\n<!-- -->\n"
		if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("tasks %q: %s", c, code)
		}
	}
}

func TestBlockQuotedStructureInTasksIsFatal(t *testing.T) { // B2
	for _, tasks := range []string{
		"- [ ] real\n\n> " + fence + "\n> - [x] example\n> " + fence + "\n",
		"- [ ] real\n\n> <!--\n> - [x] example\n> -->\n",
		"- [ ] real\n\n>     - [x] example\n",
	} {
		if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("%q: %s", tasks, code)
		}
	}
	prose := "- [ ] real\n\n> A quoted note, not a task.\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": prose})})
	if rep.Outcome != "report" || rep.Summary.Tasks.Total != 1 {
		t.Fatalf("quoted prose: %s %+v", rep.Outcome, rep.Fatal)
	}
}

func TestCheckboxFourPastItemIsFatal(t *testing.T) { // B3
	for _, tasks := range []string{
		"- [ ] a\n\n      - [x] code\n",
		"1. [ ] a\n\n       - [x] code\n",
		"- [ ] a\n      - [x] continuation text\n",
	} {
		if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks}); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("%q: %s", tasks, code)
		}
	}
}

func TestNestedTasksAfterIndentedParagraph(t *testing.T) { // B4
	for tasks, want := range map[string][2]int{
		"- [ ] a\n\n  details\n\n    - [x] b nested\n":                       {1, 2},
		"- [ ] T001 Setup\n  - [x] T001a sub\n  - [ ] T001b sub\n":           {1, 3},
		"1. [ ] a\n\n   more about a\n\n   - [x] b\n2. [ ] c\n":              {1, 3},
		"## Phase 1\n\n- [x] T001\n- [ ] T002\n\n## Phase 2\n\n- [ ] T003\n": {1, 3},
		"- [ ] a\n\nparagraph that ends the list\n\n- [x] b\n":               {1, 2},
	} {
		rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks})})
		if rep.Outcome != "report" || rep.Summary.Tasks.Checked != want[0] || rep.Summary.Tasks.Total != want[1] {
			t.Fatalf("%q: outcome %s tasks %+v, want %d/%d", tasks, rep.Outcome, rep.Summary, want[0], want[1])
		}
	}
}

func TestRawHTMLBlocksAreFatal(t *testing.T) { // B5
	for _, tag := range []string{"<style>", "<pre>", "<script type=x>", "<TEXTAREA>"} {
		body := "## Archive\n\n" + tag + "\n### 002 — Hidden [status: planned]\n\n## x\n</x>\n" + base
		if code := fatalCode(t, body, nil); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("roadmap %q: %s", tag, code)
		}
		if code := fatalCode(t, oneSpec, map[string]string{"specs/a/tasks.md": tag + "\n- [x] t\n"}); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("tasks %q: %s", tag, code)
		}
	}
	ok := "## Notes\n\n<details>\n<summary>More</summary>\n\nVisible text.\n\n</details>\n\n" + base
	if rep := Evaluate(Option{Root: project(t, ok, nil), FailOn: "none"}); rep.Outcome != "report" {
		t.Fatalf("<details> rejected: %+v", rep.Fatal)
	}
}

func TestTwoStatusMarkersAreFatal(t *testing.T) { // N1
	if code := fatalCode(t, "### 001 — A [status: planned] [status: verified]\n", nil); code != "SS_ITEM_MALFORMED" {
		t.Fatalf("code %s", code)
	}
}
