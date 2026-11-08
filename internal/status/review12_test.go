package status

import "testing"

// Cases from the third T026 acceptance review round of 001 (8 November 2026).

func TestItemAfterHTMLTagLineIsFatal(t *testing.T) { // B1
	for _, body := range []string{
		"Demo Roadmap\n============\n<span>\n### 001 — A [status: planned]\n- superseded-by: 002\n\n### 002 — B [status: planned]\n",
		"# Roadmap\n\n    example\n<span>\n### 001 — A [status: planned]\n- Depends on: 002\n\n### 002 — B [status: planned]\n",
		"# Roadmap\n\n<img src=\"a.png\" alt=\"x > y\">\n### 001 — A [status: planned]\n- Depends on: 002\n\n### 002 — B [status: planned]\n",
		"# Roadmap\n\n</div>\n- Depends on: 002\n\n### 001 — A [status: planned]\n",
	} {
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: code %s", body, code)
		}
	}
	ok := "# Roadmap\n\n<span>intro</span>\n\n### 001 — A [status: planned]\n- Depends on: none\n"
	if rep := Evaluate(Option{Root: project(t, ok, nil), FailOn: "none"}); rep.Outcome != "report" {
		t.Fatalf("tag line followed by a blank line became fatal: %+v", rep.Fatal)
	}
}

func TestMarkupInDependsOnIsMalformed(t *testing.T) { // B2
	tail := "\n### 002 — B [status: planned]\n### 003 — C [status: verified]\n### 004 — D [status: planned]\n"
	for _, v := range []string{
		`002 ([see](x "), 003 ("))`,
		`002 (a<!--), 003 (b-->)`,
		`002 (<span title="x), 004 (y">z</span>)`,
		"002 (`x), 003 (y`)",
		`002 \, 003`,
		`002 &amp; 003`,
	} {
		body := "### 001 — A [status: planned]\n- Depends on: " + v + tail
		rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
		if rep.Outcome != "report" || !hasDiag(rep, "001", "SS_DEPENDENCY_MALFORMED") {
			t.Fatalf("%q: outcome %s diags %v", v, rep.Outcome, diagCodesAll(rep))
		}
	}
}

func TestEscapedAndLookalikeColonsAreLabels(t *testing.T) { // N1, N2
	for _, line := range []string{
		"  superseded-by\\: 004",
		"  superseded-by∶ 004",
		"  Depends on꞉ 004",
	} {
		body := "### 001 — A [status: planned]\n- Description: x\n" + line + "\n### 004 — D [status: planned]\n"
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: code %s", line, code)
		}
	}
}
