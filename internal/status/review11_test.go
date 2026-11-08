package status

import "testing"

// Cases from the final T026 acceptance review of 001 (8 November 2026).
// Closed field body (D-2026-11-08-04): under a metric field, an indented
// line either joins the value as paragraph text or is fatal.

func TestClosedMetricFieldBody(t *testing.T) {
	head := "### 002 — B [status: verified]\n### 003 — C [status: planned]\n### 001 — A [status: planned]\n"
	for _, tail := range []string{
		"- **Depends on:** 002\n      <!-- c --> 003\n",                              // B-1
		"- **Depends on:** 002\n\t\t<!-- c --> 003\n",                                // B-1, tabs
		"- **Depends on:** 002\n      <div>and 003</div>\n",                          // B-1
		"- superseded-by: 002\n      <!-- was --> 003\n",                             // B-1
		"- Spec dir: specs/001-a\n      <p>/v2</p>\n",                                // B-1
		"- **Depends on:** 002\n    <!-- c --> 003\n",                                // B-2
		"- **Depends on:** 002\n\n    <!-- c --> 003\n",                              // B-2
		"- **Depends on:** 002\n\n  <details><summary>also</summary>003</details>\n", // B-3
		"- **Spec dir:** specs/001-a\n  <p>moved to specs/009-a</p>\n\n- **Description:** x\n",
		"- **Depends on:** 002\n    <div>003</div>\n",
		"- **Depends on:** 002\n  ```\n  003\n  ```\n", // nested code (owner ruling)
		"- **Depends on:** 002\n  <!-- note -->\n",
	} {
		if code := fatalCode(t, head+tail, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: code %s", tail, code)
		}
	}
}

func TestMetricContinuationStillJoins(t *testing.T) {
	body := "### 002 — B [status: verified]\n### 003 — C [status: verified]\n### 001 — A [status: planned]\n- **Depends on:** 002,\n  003\n- Notes: x\n  <!-- fine under a descriptive field -->\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if rep.Outcome != "report" {
		t.Fatalf("fatal: %+v", rep.Fatal)
	}
	if got := depIDs(itemByID(t, rep, "001")); len(got) != 2 {
		t.Fatalf("deps %v", got)
	}
}
