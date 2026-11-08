package status

import "testing"

// Cases from the T026 acceptance review of 001 (8 November 2026), judged
// against D-2026-11-07-09. B1: a list or quote nested under a metric field
// was ignored after a blank line, comment, or quote.

func TestNestedBlockUnderMetricFieldIsFatal(t *testing.T) {
	head := "### 001 — A [status: verified]\n- Description: x\n\n### 003 — C [status: undecided]\n- Description: z\n\n### 002 — B [status: verified]\n- Description: y\n"
	for _, tail := range []string{
		"- **Depends on:** 001\n\n  - 003\n",
		"- **Depends on:** 001\n  > 003\n",
		"- **Depends on:** 001\n\n  > 003\n",
		"- **Depends on:** 001\n  <!-- c -->\n  - 003\n",
		"- **Depends on:** 001\n\n  1. 003\n",
		"- **Depends on:**\n  001\n\n  - 003\n",
		"- Spec dir: specs/002-b\n\n  - specs/other\n",
		"- superseded-by: 001\n\n  * 003\n",
		"- **Depends on:** 001\n  - 003\n",
	} {
		if code := fatalCode(t, head+tail, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: code %s", tail, code)
		}
	}
}

func TestNestedBlockUnderDescriptiveFieldStillAllowed(t *testing.T) {
	body := "### 001 — A [status: planned]\n- Notes: list\n\n  - one\n  > quoted\n- Depends on: none\n"
	if rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"}); rep.Outcome != "report" {
		t.Fatalf("fatal: %+v", rep.Fatal)
	}
}
