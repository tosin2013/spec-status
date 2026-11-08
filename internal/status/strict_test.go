package status

import (
	"strings"
	"testing"
)

// Cases from the 2026-11-07 acceptance review of 001. Each one used to drop
// item-like input without a diagnostic; each must now fail visibly.

func fatalCode(t *testing.T, body string, files map[string]string) string {
	t.Helper()
	rep := Evaluate(Option{Root: project(t, body, files), FailOn: "none"})
	if rep.Outcome != "fatal" || rep.Fatal == nil || ExitCode(rep) != 2 {
		t.Fatalf("expected fatal for %q, got outcome=%s items=%d diags=%v", body, rep.Outcome, len(rep.Items), diagCodesAll(rep))
	}
	return rep.Fatal.Code
}

func diagCodesAll(rep Report) []string {
	var out []string
	for _, d := range rep.Diagnostics {
		out = append(out, d.ItemID+":"+d.Code)
	}
	return out
}

const base = "### 001 — Base [status: verified]\n"

func TestItemLikeHeadingsAndRowsAreFatal(t *testing.T) { // review D1
	for _, line := range []string{
		"## 002 — Level two [status: planned]",
		"#### 003 — Level four [status: planned]",
		"   ### 004 — Indented [status: planned]",
		"## 005 — No status marker",
		"| 006 | Table item | planned |",
		"|006|Table|planned|",
		"- 007 — Bullet item [status: planned]",
		"**008 — Bold item [status: planned]**",
	} {
		code := fatalCode(t, base+line+"\n", nil)
		if code != "SS_ITEM_MALFORMED" && code != "SS_LAYOUT_UNSUPPORTED" {
			t.Fatalf("%q: code %s", line, code)
		}
	}
}

func TestOrdinaryHeadingsAndProseStayProse(t *testing.T) {
	body := "# Roadmap\n## Vision & End States\n- Ship it.\n| Column | Value |\n|---|---|\n| a | b |\n## 2026 goals\n### Core\n" + base +
		"- Description: Uses the `### NNN — Title [status: token]` shape.\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 1 {
		t.Fatalf("prose became fatal or an item: %s %+v", rep.Outcome, rep.Fatal)
	}
}

func TestHeadingEndsItemSoFieldsCannotAttachAbove(t *testing.T) { // review D1
	code := fatalCode(t, base+"### Group\n- Depends on: 001\n", nil)
	if code != "SS_ITEM_MALFORMED" {
		t.Fatalf("code %s", code)
	}
	if code := fatalCode(t, "- Spec dir: specs/a\n"+base, nil); code != "SS_ITEM_MALFORMED" {
		t.Fatalf("field before any item: %s", code)
	}
}

func TestNearMissFieldLabelsAreFatal(t *testing.T) { // review D3
	for _, line := range []string{
		"- Superseded-by: 002",
		"- superseded by: 002",
		"- depends on: 002",
		"- Depends On: 002",
		"* Depends on: 002",
		"-  Depends on: 002",
		"  - Depends on: 002",
		"1. Depends on: 002",
		"- **depends on:** 002",
		"- Spec-dir: specs/a",
		"- SCOPE: x",
	} {
		body := base + line + "\n### 002 — Next [status: verified]\n"
		if code := fatalCode(t, body, nil); code != "SS_ITEM_MALFORMED" {
			t.Fatalf("%q: code %s", line, code)
		}
	}
	// Other extension fields and free bullets remain descriptive text.
	ok := base + "- **Scope (out):** x\n- **Governed by:** C-01\n- **Notes:** depends on: nothing\n"
	if rep := Evaluate(Option{Root: project(t, ok, nil), FailOn: "none"}); rep.Outcome != "report" {
		t.Fatalf("descriptive bullets became fatal: %+v", rep.Fatal)
	}
}

func TestSupersededTypoCannotEarnCredit(t *testing.T) { // review D3 repro
	body := "### 001 — Original [status: verified]\n- Superseded-by: 002\n" +
		"### 002 — Replacement [status: verified]\n### 003 — Follow [status: planned]\n- Depends on: 001\n"
	fatalCode(t, body, nil)
}

func TestFencesFollowCommonMarkOpeners(t *testing.T) { // review D2
	inline := "```x```\n" + base + "### 002 — Two [status: planned]\n"
	rep := Evaluate(Option{Root: project(t, inline, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 2 {
		t.Fatalf("inline code hid an item: %s items=%d", rep.Outcome, len(rep.Items))
	}
	indented := base + "    ```\n### 002 — Two [status: planned]\n"
	rep = Evaluate(Option{Root: project(t, indented, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 2 {
		t.Fatalf("indented code hid an item: %s items=%d", rep.Outcome, len(rep.Items))
	}
	fenced := base + "~~~~\n### 999 — Example [status: planned]\n~~~\n~~~~~\n### 002 — Two [status: planned]\n"
	rep = Evaluate(Option{Root: project(t, fenced, nil), FailOn: "none"})
	if rep.Outcome != "report" || len(rep.Items) != 2 {
		t.Fatalf("fence close handling: %s items=%d", rep.Outcome, len(rep.Items))
	}
	if code := fatalCode(t, base+"```\n### 002 — Two [status: planned]\n", nil); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("unclosed fence: %s", code)
	}
}

func TestTaskFencesAndOrderedCheckboxes(t *testing.T) { // review D2
	body := "### 001 — A [status: planned]\n- Spec dir: specs/a\n"
	tasks := "- [x] one\n```js```\n- [ ] two\n1. [ ] three\n2) [x] four\n    ```\n    - [x] example in a fence\n    ```\n- [ ] five\n"
	rep := Evaluate(Option{Root: project(t, body, map[string]string{"specs/a/tasks.md": tasks})})
	if rep.Outcome != "report" || rep.Summary.Tasks.Total != 5 || rep.Summary.Tasks.Checked != 2 {
		t.Fatalf("tasks %+v, want 2/5", rep.Summary.Tasks)
	}
	if code := fatalCode(t, body, map[string]string{"specs/a/tasks.md": "- [x] one\n```\n- [ ] hidden\n"}); code != "SS_LAYOUT_UNSUPPORTED" {
		t.Fatalf("unclosed task fence: %s", code)
	}
}

func TestAggregateCheckboxBudget(t *testing.T) { // review D5
	var b strings.Builder
	files := map[string]string{}
	per := maxTasks/2 + 1
	page := strings.Repeat("- [ ] t\n", per)
	for _, name := range []string{"a", "b"} {
		b.WriteString("### 00" + map[string]string{"a": "1", "b": "2"}[name] + " — " + name + " [status: planned]\n- Spec dir: specs/" + name + "\n")
		files["specs/"+name+"/tasks.md"] = page
	}
	rep := Evaluate(Option{Root: project(t, b.String(), files), FailOn: "none"})
	if rep.Outcome != "fatal" || rep.Fatal.Code != "SS_INPUT_LIMIT" {
		t.Fatalf("aggregate %d checkboxes accepted: %s", 2*per, rep.Outcome)
	}
}

func TestTaskLinkForms(t *testing.T) { // review D6
	body := "### 001 — A [status: planned]\n- Spec dir: specs/a\n"
	tasks := "- [x] Titled [t](exists.md \"Title\")\n" +
		"- [x] Angle [t](<exists.md> 'Title')\n" +
		"- [x] Reference [doc][r1]\n" +
		"- [x] Autolink <missing-auto.md>\n" +
		"- [x] Junk after path [t](exists.md junk)\n" +
		"\n[r1]: missing-ref.md\n"
	rep := Evaluate(Option{Root: project(t, body, map[string]string{
		"specs/a/tasks.md":  tasks,
		"specs/a/exists.md": "x\n",
	}), FailOn: "none"})
	var missing, unsupported int
	for _, d := range rep.Diagnostics {
		switch d.Code {
		case "SS_TASK_TARGET_MISSING":
			missing++
		case "SS_REFERENCE_UNSUPPORTED":
			unsupported++
		}
	}
	if missing != 0 || unsupported != 3 {
		t.Fatalf("missing=%d unsupported=%d, want 0 and 3: %v", missing, unsupported, diagCodesAll(rep))
	}
}
