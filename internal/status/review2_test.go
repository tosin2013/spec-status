package status

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Cases from the second acceptance review of 001 (7 November 2026), plus the
// coverage gaps it named. Owner decisions: the prefix rule for field labels,
// narrowed item-like headings, and human review for any supplied
// superseded-by line.

func depStates(t *testing.T, rep Report, id string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, d := range itemByID(t, rep, id).Dependencies {
		out[d.ID] = d.State
	}
	return out
}

func TestPrefixRuleCatchesColonlessAndWrappedLabels(t *testing.T) { // D1, D3
	for _, line := range []string{
		"- superseded-by 002",
		"- superseded-by — 002",
		"- superseded-by= 002",
		"- Superseded by 002",
		"- Depends on 002",
		"- **Depends on** 002",
		"- `Depends on:` 002",
		"> - Depends on: 002",
		"> Depends on: 002",
		"- Spec dir specs/a",
		"  - Depends on: 002",
	} {
		body := "### 001 — Old [status: verified]\n" + line + "\n### 002 — New [status: verified]\n"
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: code %s", line, code)
		}
	}
}

func TestPrefixRuleLeavesProseAlone(t *testing.T) {
	body := "# Overview\n- Scope: the whole product.\n- Description: a roadmap.\n" +
		"### 001 — A [status: planned]\n" +
		"- Scope creep is a risk.\n- Outcomes vary.\n- Descriptions follow.\n" +
		"- **Notes:** depends on nothing.\n  - Scope: nested note.\n  - Outcome: nested note.\n" +
		"- **Scope (out):** later.\n- Risk: low\n" +
		"  003–005 carry the obligation.\n- 100 - 200 ms budget\n" +
		"                                                                    deep prose line\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("prose rejected: %s %+v", rep.Outcome, rep.Fatal)
	}
}

func TestNarrowedItemHeadings(t *testing.T) {
	prose := "## 100 days plan\n### 2026 goals\n### 1. Introduction\n### 001 Title without dash\n" + base
	rep := Evaluate(Option{Root: project(t, prose, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("numeric headings rejected: %s %+v", rep.Outcome, rep.Fatal)
	}
	for _, line := range []string{"## 002 — Two", "### 002 - Two", "#### 002 – Two", "### 01 — Short [status: planned]", "- 002 — Bullet [verified]"} {
		if code := fatalCode(t, base+line+"\n", nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: code %s", line, code)
		}
	}
}

func TestTablesWithoutOuterPipes(t *testing.T) { // D2
	for _, body := range []string{
		"ID | Title | Status\n--- | --- | ---\n001 | Parser | verified\n",
		"Title | Status\n:--- | ---:\nParser | [status: verified]\n",
		"001 | Parser\n--- | ---\n",
	} {
		if code := fatalCode(t, body, nil); code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("%q: code %s", body, code)
		}
	}
	ok := base + "Name | Value\n--- | ---\nowner | team\n\n- Description: after the table\n"
	if rep := Evaluate(Option{Root: project(t, ok, nil), FailOn: "none"}); rep.Outcome != "report" || itemByID(t, rep, "001").Description != "after the table" {
		t.Fatalf("ordinary table broke parsing: %+v", rep.Fatal)
	}
}

func TestHTMLCommentsAreNotLive(t *testing.T) { // D5
	body := "<!--\n### 002 — Draft [status: verified]\n-->\n" + base + "<!-- one line ### 003 — X [status: planned] -->\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("commented item counted: %s items=%d", rep.Outcome, len(rep.Items))
	}
	if code := fatalCode(t, base+"<!--\n### 002 — Hidden [status: planned]\n", nil); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("unclosed comment: %s", code)
	}
	tasks := "- [x] one\n<!--\n- [ ] commented\n-->\n- [ ] two\n"
	rep = Evaluate(Option{Root: project(t, "### 001 — A [status: planned]\n- Spec dir: specs/a\n", map[string]string{"specs/a/tasks.md": tasks})})
	if rep.Summary.Tasks.Total != 2 {
		t.Fatalf("commented task counted: %+v", rep.Summary.Tasks)
	}
}

func TestSetextHeadingEndsItem(t *testing.T) { // D6
	for _, u := range []string{"=====", "-----"} {
		body := base + "\nOther\n" + u + "\n- Depends on: 002\n### 002 — B [status: planned]\n"
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%s: code %s", u, code)
		}
	}
}

func TestCycleLabelsOnlyCycleEdges(t *testing.T) { // D7
	body := "### 001 — A [status: planned]\n- Depends on: 002, 003, 005\n" +
		"### 002 — B [status: planned]\n- Depends on: 001\n" +
		"### 003 — Old [status: verified]\n- superseded-by: 004\n" +
		"### 004 — New [status: verified]\n" +
		"### 005 — Done [status: verified]\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	want := map[string]string{"002": "cycle", "003": "human_review", "005": "satisfied"}
	if got := depStates(t, rep, "001"); !reflect.DeepEqual(got, want) {
		t.Fatalf("001 deps %v, want %v", got, want)
	}
	if rep.Summary.HumanReviewCount != 1 {
		t.Fatalf("human review hidden by cycle: %d", rep.Summary.HumanReviewCount)
	}
	cycles := 0
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_DEPENDENCY_CYCLE" {
			cycles++
		}
	}
	if cycles != 2 {
		t.Fatalf("cycle diagnostics %d, want 2 (001→002, 002→001)", cycles)
	}
}

func TestInvalidSupersessionStillNeedsHumanReview(t *testing.T) { // D8
	body := "### 001 — A [status: verified]\n- superseded-by: 002\n" +
		"### 002 — B [status: verified]\n- superseded-by: 099\n" +
		"### 003 — Dep [status: planned]\n- Depends on: 001\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if st := depStates(t, rep, "003")["001"]; st != "human_review" {
		t.Fatalf("dependency on item with invalid chain is %q", st)
	}
	if !itemByID(t, rep, "001").Active {
		t.Fatal("invalid chain must leave 001 active (D-02A)")
	}
}

func TestDiscoveryStopsAtUnusableSpecify(t *testing.T) { // D9
	outer := t.TempDir()
	writeFile(t, filepath.Join(outer, ".specify", "memory", "roadmap.md"), base)
	for _, kind := range []string{"file", "symlink"} {
		inner := filepath.Join(outer, kind)
		if err := os.MkdirAll(inner, 0o755); err != nil {
			t.Fatal(err)
		}
		marker := filepath.Join(inner, ".specify")
		if kind == "file" {
			writeFile(t, marker, "x")
		} else if err := os.Symlink(filepath.Join(outer, ".specify"), marker); err != nil {
			t.Fatal(err)
		}
		rep := Evaluate(Option{CWD: inner})
		if rep.Outcome != "fatal" || rep.Fatal.Code != "SS_ROOT_UNAVAILABLE" {
			t.Fatalf("%s .specify: discovery widened to outer project (%s)", kind, rep.Outcome)
		}
	}
}

func TestSchemeAndAbsoluteReferences(t *testing.T) { // D10
	tasks := "- [x] [a](file:spec.md)\n- [x] [b](javascript:a)\n- [x] [c](/etc/hosts)\n- [x] [d](spec.md)\n"
	rep := Evaluate(Option{Root: project(t, "### 001 — A [status: planned]\n- Spec dir: specs/a\n", map[string]string{
		"specs/a/tasks.md": tasks, "specs/a/spec.md": "x\n",
	}), FailOn: "none"})
	unsupported := 0
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_TASK_TARGET_MISSING" {
			t.Fatalf("scheme or absolute link checked as local: %+v", d)
		}
		if d.Code == "SS_REFERENCE_UNSUPPORTED" {
			unsupported++
		}
	}
	if unsupported != 3 {
		t.Fatalf("unsupported %d, want 3", unsupported)
	}
	url := Evaluate(Option{Root: project(t, "### 001 — A [status: specced]\n- Spec dir: https://x.org/specs\n", nil), FailOn: "none"})
	if !hasDiag(url, "001", "SS_PATH_UNSAFE") || itemByID(t, url, "001").Tasks.State == "read" {
		t.Fatalf("URL spec dir not rejected: %v", diagCodesAll(url))
	}
}

// Coverage the review found missing.

func TestDuplicateIDsAndAmbiguousDependencies(t *testing.T) {
	body := "### 001 — A [status: verified]\n### 001 — A again [status: planned]\n### 002 — B [status: planned]\n- Depends on: 001\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if rep.Summary.TotalItems != 3 {
		t.Fatalf("duplicate occurrences not retained: %d", rep.Summary.TotalItems)
	}
	dups := 0
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_DUPLICATE_ID" {
			dups++
		}
	}
	if dups == 0 || depStates(t, rep, "002")["001"] != "ambiguous" {
		t.Fatalf("dups=%d state=%q", dups, depStates(t, rep, "002")["001"])
	}
}

func TestDependencyMissingSelfCycle(t *testing.T) {
	body := "### 001 — A [status: planned]\n- Depends on: 001, 099, 002\n### 002 — B [status: planned]\n- Depends on: 001\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	want := map[string]string{"001": "self", "099": "missing", "002": "cycle"}
	if got := depStates(t, rep, "001"); !reflect.DeepEqual(got, want) {
		t.Fatalf("deps %v, want %v", got, want)
	}
	for _, code := range []string{"SS_DEPENDENCY_SELF", "SS_DEPENDENCY_MISSING", "SS_DEPENDENCY_CYCLE"} {
		if !hasDiag(rep, "001", code) {
			t.Fatalf("missing %s", code)
		}
	}
}

func TestBOMAndLineEndings(t *testing.T) {
	lf := "### 001 — A [status: planned]\n- Description: d\n### 002 — B [status: planned]\n- Depends on: 001\n"
	want := JSON(Evaluate(Option{Root: project(t, lf, nil)}))
	for name, body := range map[string]string{
		"bom":  "\xef\xbb\xbf" + lf,
		"crlf": "### 001 — A [status: planned]\r\n- Description: d\r\n### 002 — B [status: planned]\r\n- Depends on: 001\r\n",
		"cr":   "### 001 — A [status: planned]\r- Description: d\r### 002 — B [status: planned]\r- Depends on: 001\r",
	} {
		if got := JSON(Evaluate(Option{Root: project(t, body, nil)})); string(got) != string(want) {
			t.Fatalf("%s: JSON differs from LF input", name)
		}
	}
}
