package status

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEmptyRoadmapJSON(t *testing.T) {
	dir := project(t, "", nil)
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if rep.Outcome != "report" || !rep.Complete || ExitCode(rep) != 0 {
		t.Fatalf("empty report %+v exit %d", rep.Fatal, ExitCode(rep))
	}
	if rep.Summary.VerifiedPercent.Text() != "N/A" || rep.Summary.Tasks.Percent.Text() != "N/A" {
		t.Fatalf("empty denominators must be N/A, got %s %s", rep.Summary.VerifiedPercent.Text(), rep.Summary.Tasks.Percent.Text())
	}
	got := string(JSON(rep))
	var doc map[string]any
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "quarter") || strings.Contains(got, "as_of") || strings.Contains(got, "inferred") {
		t.Fatalf("unexpected timing fields: %s", got)
	}
	raw, _ := json.Marshal(doc["summary"])
	if strings.Contains(string(raw), "100") {
		t.Fatalf("empty summary must not become 100: %s", raw)
	}
	again := string(JSON(Evaluate(Option{Root: dir, CWD: dir})))
	if got != again {
		t.Fatal("JSON is not deterministic")
	}
}

func TestHistoryHumanReviewDoesNotRedirect(t *testing.T) {
	dir := project(t, exampleRoadmap, map[string]string{
		"specs/001-original/spec.md":     "# spec\n",
		"specs/001-original/plan.md":     "# plan\n",
		"specs/001-original/tasks.md":    "- [ ] Keep the original note [note](note.md)\n",
		"specs/001-original/note.md":     "ok\n",
		"specs/002-replacement/spec.md":  "# spec\n",
		"specs/002-replacement/plan.md":  "# plan\n",
		"specs/002-replacement/tasks.md": "- [x] Record the replacement [out](out.md)\n",
		"specs/002-replacement/out.md":   "ok\n",
	})
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if rep.Outcome != "report" {
		t.Fatalf("fatal %v", rep.Diagnostics)
	}
	if rep.Summary.ActiveItems != 2 || rep.Summary.VerifiedActiveItems != 1 || rep.Summary.SupersededItems != 1 {
		t.Fatalf("summary %+v", rep.Summary)
	}
	if rep.Summary.VerifiedPercent.Text() != "50" {
		t.Fatalf("percent %s", rep.Summary.VerifiedPercent.Text())
	}
	if rep.Summary.HumanReviewCount != 1 || rep.Summary.BlockedCount != 1 {
		t.Fatalf("review counts blocked=%d review=%d", rep.Summary.BlockedCount, rep.Summary.HumanReviewCount)
	}
	var follow Dependency
	for _, it := range rep.Items {
		if it.ID == "001" && !it.Active {
			if it.Status != "implemented" || it.SupersededBy != "002" {
				t.Fatalf("historical item %+v", it)
			}
		}
		if it.ID == "003" {
			if len(it.Dependencies) != 1 {
				t.Fatal(it.Dependencies)
			}
			follow = it.Dependencies[0]
		}
	}
	if follow.ID != "001" || follow.Satisfied || follow.State != "human_review" || follow.ReplacementID != "002" {
		t.Fatalf("dependency was rewritten or satisfied: %+v", follow)
	}
	if ExitCode(rep) != 0 {
		t.Fatalf("warnings must not fail the default error threshold, diags %+v", rep.Diagnostics)
	}
	warn := Evaluate(Option{Root: dir, CWD: dir, FailOn: "warning"})
	if ExitCode(warn) != 1 {
		t.Fatal("warning threshold should fail")
	}
}

func TestUnknownRawStatusAndRounding(t *testing.T) {
	body := "### 001 — One [status: Planned]\n" +
		"### 002 — Two [status: verified]\n" +
		"### 003 — Three [status: verified]\n"
	dir := project(t, body, map[string]string{
		"specs/a/spec.md": "", "specs/a/plan.md": "", "specs/a/tasks.md": "- [x] A\n",
	})
	// No spec dirs, so verified items miss required files. Use statuses only by omitting stage files
	// and accepting those diagnostics. Rounding is covered below with a direct percent check.
	if percentOf(1, 3).Text() != "33.33" || percentOf(2, 3).Text() != "66.67" {
		t.Fatalf("rounding %s %s", percentOf(1, 3).Text(), percentOf(2, 3).Text())
	}
	rep := Evaluate(Option{Root: dir, CWD: dir})
	found := false
	for _, it := range rep.Items {
		if it.ID == "001" {
			found = true
			if it.Status != "Planned" || it.StatusKnown || !it.Active {
				t.Fatalf("raw status changed: %+v", it)
			}
		}
	}
	if !found {
		t.Fatal("missing item")
	}
	seen := false
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_STATUS_UNKNOWN" {
			seen = true
		}
	}
	if !seen {
		t.Fatal("expected unknown warning")
	}
}

func TestInvalidSupersessionStaysActive(t *testing.T) {
	body := "### 001 — Kept [status: verified]\n- superseded-by: 999\n"
	dir := project(t, body, map[string]string{
		"specs/001/spec.md": "", "specs/001/plan.md": "", "specs/001/tasks.md": "- [x] Done\n",
	})
	// No spec dir on the item, so stage files are not consulted. The missing spec link is an error.
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if len(rep.Items) != 1 || !rep.Items[0].Active || rep.Summary.SupersededItems != 0 || rep.Summary.InvalidSupersessionItems != 1 {
		t.Fatalf("invalid supersession changed metrics: active=%v summary=%+v", rep.Items[0].Active, rep.Summary)
	}
	if rep.Items[0].Status != "verified" || rep.Summary.VerifiedActiveItems != 1 {
		t.Fatal("verified label was rewritten")
	}
}

func TestSupersessionCycleAndChain(t *testing.T) {
	body := strings.Join([]string{
		"### 001 — A [status: planned]",
		"- superseded-by: 002",
		"### 002 — B [status: planned]",
		"- superseded-by: 001",
		"### 010 — Origin [status: verified]",
		"- superseded-by: 011",
		"### 011 — Mid [status: planned]",
		"- superseded-by: 099",
		"### 020 — Old [status: implemented]",
		"- superseded-by: 021",
		"### 021 — Mid2 [status: planned]",
		"- superseded-by: 022",
		"### 022 — New [status: verified]",
	}, "\n") + "\n"
	dir := project(t, body, nil)
	rep := Evaluate(Option{Root: dir, CWD: dir})
	by := map[string]Item{}
	for _, it := range rep.Items {
		by[it.ID] = it
	}
	if !by["001"].Active || !by["002"].Active {
		t.Fatal("cycle members were excluded")
	}
	if !by["010"].Active || !by["011"].Active {
		t.Fatalf("invalid chain excluded rows: 010 %v 011 %v", by["010"].Active, by["011"].Active)
	}
	if by["010"].Supersession.Valid == nil || *by["010"].Supersession.Valid {
		t.Fatal("invalid chain marked valid")
	}
	if by["020"].Active || by["021"].Active || !by["022"].Active {
		t.Fatalf("valid chain active flags 020=%v 021=%v 022=%v", by["020"].Active, by["021"].Active, by["022"].Active)
	}
	if by["020"].Supersession.TerminalID != "022" {
		t.Fatalf("terminal %q", by["020"].Supersession.TerminalID)
	}
	if by["020"].Status != "implemented" || by["022"].Status != "verified" {
		t.Fatal("raw statuses changed")
	}
	if rep.Summary.VerifiedActiveItems != 2 {
		t.Fatalf("010 stays active verified and 022 is active verified, got %+v", rep.Summary)
	}
}

func TestTaskWeightsAndFences(t *testing.T) {
	body := strings.Join([]string{
		"### 001 — A [status: planned]",
		"- Spec dir: specs/a",
		"### 002 — B [status: planned]",
		"- Spec dir: specs/b",
		"### 003 — Shared [status: deferred]",
		"- Spec dir: specs/a",
	}, "\n") + "\n"
	tasksA := "- [x] One\n```\n- [x] Fenced\n```\n"
	tasksB := strings.Repeat("- [ ] Pad\n", 9)
	dir := project(t, body, map[string]string{
		"specs/a/tasks.md": tasksA,
		"specs/b/tasks.md": tasksB,
	})
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if rep.Summary.Tasks.Checked != 1 || rep.Summary.Tasks.Total != 10 || rep.Summary.Tasks.FilesRead != 2 {
		t.Fatalf("aggregate %+v", rep.Summary.Tasks)
	}
	if rep.Summary.Tasks.Percent.Text() != "10" {
		t.Fatalf("percent %s, not a mean of item percentages", rep.Summary.Tasks.Percent.Text())
	}
	for _, it := range rep.Items {
		if it.ID == "003" && it.Tasks.Total != 1 {
			t.Fatalf("shared file not attributed to inactive row: %+v", it.Tasks)
		}
	}
}

func TestCheckedVersusUncheckedTarget(t *testing.T) {
	body := "### 001 — A [status: planned]\n- Spec dir: specs/a\n"
	tasks := "- [x] Need [file](missing-checked.md)\n- [ ] Later [file](missing-later.md)\n- [x] URL [web](https://example.test/a)\n- [ ] Code `not-a-link.md`\n"
	dir := project(t, body, map[string]string{"specs/a/tasks.md": tasks})
	rep := Evaluate(Option{Root: dir, CWD: dir})
	var checkedErr, uncheckedInfo, unsupported int
	for _, d := range rep.Diagnostics {
		switch {
		case d.Code == "SS_TASK_TARGET_MISSING" && d.Severity == "error":
			checkedErr++
		case d.Code == "SS_TASK_TARGET_MISSING" && d.Severity == "info":
			uncheckedInfo++
		case d.Code == "SS_REFERENCE_UNSUPPORTED":
			unsupported++
		}
	}
	if checkedErr != 1 || uncheckedInfo != 1 || unsupported != 1 {
		t.Fatalf("checked=%d unchecked=%d unsupported=%d diags=%v", checkedErr, uncheckedInfo, unsupported, rep.Diagnostics)
	}
}

func TestMalformedLayoutIsFatal(t *testing.T) {
	for _, body := range []string{
		"| 001 | planned |\n| [status: planned] |\n",
		"### 01 — Short [status: planned]\n",
		"### 001 — Bad [status: planned]\n- Description: one\n- Description: two\n",
		"\xff\xfe not utf-8",
	} {
		dir := project(t, body, nil)
		rep := Evaluate(Option{Root: dir, CWD: dir, FailOn: "none"})
		if rep.Outcome != "fatal" || rep.Complete || rep.Summary != nil || ExitCode(rep) != 2 {
			t.Fatalf("layout %q was not fatal: outcome=%s summary=%v", body, rep.Outcome, rep.Summary)
		}
		if len(rep.Items) != 0 {
			t.Fatal("fatal report published items")
		}
	}
}

func TestScaleSharedTaskFile(t *testing.T) {
	var body strings.Builder
	for i := 1; i <= 100; i++ {
		fmt.Fprintf(&body, "### %03d — Item %d [status: planned]\n- Spec dir: specs/item\n", i, i)
	}
	var tasks strings.Builder
	tasks.WriteString("- [x] First\n")
	for i := 1; i < 10000; i++ {
		tasks.WriteString("- [ ] Rest\n")
	}
	dir := project(t, body.String(), map[string]string{"specs/item/tasks.md": tasks.String()})
	start := time.Now()
	rep := Evaluate(Option{Root: dir, CWD: dir})
	elapsed := time.Since(start)
	t.Logf("scale elapsed %s", elapsed)
	if rep.Summary == nil || rep.Summary.TotalItems != 100 || rep.Summary.Tasks.Total != 10000 || rep.Summary.Tasks.Checked != 1 || rep.Summary.Tasks.FilesRead != 1 {
		t.Fatalf("scale summary %+v elapsed %s fatal %v", rep.Summary, elapsed, rep.Diagnostics)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("scale exceeded sanity bound: %s", elapsed)
	}
}

func TestDeferredAreNotActive(t *testing.T) {
	body := "### 001 — Done [status: deferred]\n### 002 — Left [status: abandoned]\n"
	dir := project(t, body, nil)
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if rep.Summary.ActiveItems != 0 || rep.Summary.VerifiedPercent.ok {
		t.Fatalf("inactive roadmap got a percentage: %+v", rep.Summary)
	}
}

func project(t *testing.T, roadmap string, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".specify", "memory", "roadmap.md"), roadmap)
	for name, body := range files {
		writeFile(t, filepath.Join(dir, filepath.FromSlash(name)), body)
	}
	return dir
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

const exampleRoadmap = "### 001 — Preserve original feature [status: implemented]\n" +
	"- Description: Historical original feature.\n" +
	"- Outcome: Original outcome remains recorded.\n" +
	"- Scope: Original scope is retained.\n" +
	"- Spec dir: specs/001-original\n" +
	"- superseded-by: 002\n" +
	"\n" +
	"### 002 — Replacement feature [status: verified]\n" +
	"- Description: Replacement scope can differ.\n" +
	"- Outcome: Recorded as verified, not independently checked by this CLI.\n" +
	"- Scope: Replacement scope.\n" +
	"- Spec dir: specs/002-replacement\n" +
	"\n" +
	"### 003 — Follow-on work [status: planned]\n" +
	"- Description: Still refers to the original requirement.\n" +
	"- Outcome: Requires a scope decision.\n" +
	"- Scope: Follow-on.\n" +
	"- Depends on: 001\n"
