package status

import "testing"

// Cases from the eighth acceptance review of 001 (7 November 2026), the first
// judged against D-2026-11-07-09. D1 to D3 let a Depends on, Spec dir, or
// superseded-by value be silently misread (rule (b)); each now fails visibly.

func TestBlockStartsAreNotFieldContinuations(t *testing.T) { // D1
	head := "### 001 — Base [status: verified]\n- **Scope (in):** Core.\n"
	for _, line := range []string{
		" ### 002 — Old plan  [status: planned]",
		"\t### 002 — Old plan  [status: planned]",
		"   ### 002 — Old plan  [status: planned]",
	} {
		body := head + line + "\n- **Spec dir:** specs/002-old/\n- **superseded-by:** 003\n### 003 — New [status: planned]\n"
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: %s", line, code)
		}
	}
	for _, line := range []string{" ## Section", " ---", " <details>"} {
		body := head + line + "\n- **Spec dir:** specs/x\n"
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: %s", line, code)
		}
	}
	ok := head + "  wrapped scope text that continues the value\n- **Spec dir:** specs/a\n"
	rep := Evaluate(Option{Root: project(t, ok, map[string]string{"specs/a/spec.md": "x\n", "specs/a/plan.md": "x\n", "specs/a/tasks.md": "- [x] t\n"}), FailOn: "none"})
	if rep.Outcome != "report" || itemByID(t, rep, "001").Scope != "Core. wrapped scope text that continues the value" || itemByID(t, rep, "001").SpecDir != "specs/a" {
		t.Fatalf("ordinary continuation broken: %+v %+v", rep.Fatal, rep.Items)
	}
}

func TestHTMLBlocksUnderItems(t *testing.T) { // D2
	for _, body := range []string{
		base + "- **Description:** x\n\n <details><summary>Old</summary>\n- **superseded-by:** 003\n- **Spec dir:** specs/old/\n\n </details>\n### 003 — C [status: planned]\n",
		base + "- **Description:** x\n  <details>\n### 004 — Hidden [status: verified]\n- **Depends on:** 001\n\n  </details>\n",
	} {
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: %s", body, code)
		}
	}
	ok := base + "- **Notes:** see below\n\n  <details>\n  <summary>History</summary>\n  Older notes.\n  </details>\n"
	if rep := Evaluate(Option{Root: project(t, ok, nil), FailOn: "none"}); rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("details with prose under an item: %+v", rep.Fatal)
	}
}

func TestMetricValueAfterBlankLineIsFatal(t *testing.T) { // D3
	body := "### 001 — A [status: verified]\n### 002 — B [status: planned]\n- **Depends on:** 001\n\n  003\n### 003 — C [status: planned]\n"
	if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
		t.Fatalf("code %s", code)
	}
	ok := "### 001 — A [status: verified]\n### 002 — B [status: planned]\n- **Depends on:** 001\n\n- **Notes:** after a blank line\n\n  more notes after another blank\n"
	if rep := Evaluate(Option{Root: project(t, ok, nil), FailOn: "none"}); rep.Outcome != "report" {
		t.Fatalf("notes after blank lines rejected: %+v", rep.Fatal)
	}
}

func TestStruckThroughItemHeading(t *testing.T) { // D4
	body := "### ~~006 — Old idea~~  [status: abandoned]\n- **Description:** Superseded by 007.\n### 007 — New idea [status: planned]\n- **Depends on:** 006\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 2 {
		t.Fatalf("struck heading: %s %+v", rep.Outcome, rep.Fatal)
	}
	it := itemByID(t, rep, "006")
	if it.Title != "Old idea" || it.Status != "abandoned" || it.Active {
		t.Fatalf("struck item %+v", it)
	}
	for _, bad := range []string{"### ~~006 — Old idea  [status: abandoned]", "### ~~006 — Old~~ idea~~  [status: abandoned]"} {
		if code := fatalCode(t, bad+"\n", nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: %s", bad, code)
		}
	}
}
