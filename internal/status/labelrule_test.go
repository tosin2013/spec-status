package status

import (
	"reflect"
	"testing"
)

// Line-level label rule (owner direction, 8 November 2026): a visible line
// that contains a Depends on, Spec dir, or superseded-by label with a colon
// must be an accepted field line under an item, whatever the layout.

const labelBase = "### 001 — A [status: planned]\n- Outcome: x\n"

func TestMetricLabelOutsideFieldLineIsFatal(t *testing.T) {
	for _, body := range []string{
		// Layouts from earlier reviews, now caught by one rule.
		labelBase + "> Depends on: 002\n" + twoMore,
		"| Topic | Note |\n|---|---|\n| auth | Depends on: 002 |\n\n" + labelBase + twoMore,
		"<!-- intro --> Depends on: 002\n\n" + labelBase + twoMore,
		"<span>Spec dir: specs/001-a</span>\n\n" + labelBase + twoMore,
		labelBase + "- Notes: rollout depends on: vendor approval\n" + twoMore,
		labelBase + "- Depends on: 002\n  Spec dir: specs/001-a\n" + twoMore,
		labelBase + "- Notes: see below\n  <details>\n  superseded-by: 002\n  </details>\n" + twoMore,
		"Release note — Depends on: 002\n\n" + labelBase + twoMore,
		"# Roadmap\n\nItem 001 is **superseded by:** 002.\n\n" + labelBase + twoMore,
		labelBase + "- Notes: x\n  - Depends on 002 — Spec dir: specs/a\n" + twoMore,
		// Lookalike and spacing variants.
		"Depends  on ： 002\n\n" + labelBase + twoMore,
		"Spec dir: specs/a\n\n" + labelBase + twoMore,
	} {
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: code %s", body, code)
		}
	}
}

func TestMetricLabelInHiddenOrProseTextIsAllowed(t *testing.T) {
	for _, body := range []string{
		"```\n- Depends on: 002\n```\n\n" + labelBase + twoMore,
		"<!--\n- Depends on: 002\nSpec dir: specs/x\n-->\n\n" + labelBase + twoMore,
		"Write the field as `Depends on: 001`.\n\n" + labelBase + twoMore,
		"Intro <!-- Depends on: 002 --> text.\n\n" + labelBase + twoMore,
		"- Depends on whether the vendor ships v2\n\n" + labelBase + twoMore,
		labelBase + "- Notes: independs on: nothing; specdirectory: x\n" + twoMore,
	} {
		if rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"}); rep.Outcome != "report" {
			t.Fatalf("%q: became fatal: %+v", body, rep.Fatal)
		}
	}
}

func TestAcceptedFieldLinesStillRead(t *testing.T) {
	body := "### 001 — A [status: planned]\n- **Depends on:** 002,\n  003\n- Spec dir: specs/001-a\n" +
		"### 002 — B [status: verified]\n- superseded-by: 003\n### 003 — C [status: planned]\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if rep.Outcome != "report" {
		t.Fatalf("fatal: %+v", rep.Fatal)
	}
	if got := depIDs(itemByID(t, rep, "001")); !reflect.DeepEqual(got, []string{"002", "003"}) {
		t.Fatalf("deps %v", got)
	}
}
