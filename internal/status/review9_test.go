package status

import (
	"reflect"
	"testing"
)

// Cases from the ninth acceptance review of 001 (7 November 2026), judged
// against D-2026-11-07-09. D1 to D4 dropped a metric value silently.

const twoMore = "### 002 — B [status: planned]\n### 003 — C [status: planned]\n"

func TestInlineTagInsideWrappedValue(t *testing.T) { // D1
	for _, tag := range []string{"<br>", "<span>", "</span>", "<img src=\"dep.png\">"} {
		body := "### 001 — A [status: planned]\n- **Depends on:** 002 (verified)\n  " + tag + "\n  003 (planned)\n" + twoMore
		rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
		if rep.Outcome != "report" || !hasDiag(rep, "001", "SS_DEPENDENCY_MALFORMED") {
			t.Fatalf("%q: outcome %s diags %v", tag, rep.Outcome, diagCodesAll(rep))
		}
	}
	sup := "### 001 — A [status: verified]\n- **superseded-by:** 002\n  </span>\n  003\n" + twoMore
	if rep := Evaluate(Option{Root: project(t, sup, nil), FailOn: "none"}); !hasDiag(rep, "001", "SS_SUPERSESSION_MALFORMED") {
		t.Fatalf("superseded-by with an inline tag: %v", diagCodesAll(rep))
	}
}

func TestTableEndsAtListMarker(t *testing.T) { // D2
	body := "### 001 — A [status: planned]\n- Notes: matrix\n  | a | b |\n  |---|---|\n- Depends on: 002 (see C-01 | ADR-3)\n" + twoMore
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if got := depIDs(itemByID(t, rep, "001")); !reflect.DeepEqual(got, []string{"002"}) {
		t.Fatalf("deps %v (outcome %s %+v)", got, rep.Outcome, rep.Fatal)
	}
}

func TestTextAfterMultiLineCommentIsFatal(t *testing.T) { // D3
	body := "### 001 — A [status: planned]\n- Depends on: 002\n<!-- note\nold --> - Depends on: 003\n" + twoMore
	if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
		t.Fatalf("code %s", code)
	}
}

func TestBlockSeparatesMetricValue(t *testing.T) { // D4
	for _, sep := range []string{"  <!-- added later -->", "  > quoted", "  <details>"} {
		body := "### 001 — A [status: planned]\n- Depends on: 002\n" + sep + "\n  , 003\n" + twoMore
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: %s", sep, code)
		}
	}
	body := "### 001 — A [status: planned]\n- **Depends on:** 002\n  # 003\n" + twoMore
	if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
		t.Fatalf("indented heading: %s", code)
	}
}

func TestMetricLabelAsHeadingIsFatal(t *testing.T) { // D5
	for _, body := range []string{
		"### 001 — A [status: planned]\n#### Depends on: 003\n" + twoMore,
		"# Roadmap\n#### Spec dir: specs/zzz\n" + base,
	} {
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: %s", body, code)
		}
	}
}

func TestDependsOnProseIsNotAField(t *testing.T) { // D6
	ok := base + "- **Notes:** coordination\n  - Depends on the auth team shipping first\n" +
		"## Open Questions\n- Depends on whether the vendor ships v2\n- Superseded by events, maybe\n"
	if rep := Evaluate(Option{Root: project(t, ok, nil), FailOn: "none"}); rep.Outcome != "report" {
		t.Fatalf("prose rejected: %+v", rep.Fatal)
	}
	for _, line := range []string{"- Depends on 002", "- Spec dir specs/a", "- Depends on none", "- superseded by 002", "- Depends on"} {
		if code := fatalCode(t, base+line+"\n### 002 — B [status: planned]\n", nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: %s", line, code)
		}
	}
}

func TestEmptyCheckboxCounts(t *testing.T) { // D7
	tasks := "- [ ] T001 a\n- [ ]\n- [x]\n- [X] T002 b\n"
	rep := Evaluate(Option{Root: project(t, oneSpec, map[string]string{"specs/a/tasks.md": tasks})})
	if rep.Summary.Tasks.Total != 4 || rep.Summary.Tasks.Checked != 2 {
		t.Fatalf("tasks %+v", rep.Summary.Tasks)
	}
}
