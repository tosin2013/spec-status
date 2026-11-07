package status

import (
	"strings"
	"testing"
)

func TestFiltersAreViewsOverTheFullResult(t *testing.T) {
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
	full := Text(rep, View{Width: 80})
	cases := []View{
		{Width: 80, Statuses: []string{"implemented"}},
		{Width: 80, Statuses: []string{"planned", "verified"}},
		{Width: 80, Blocked: true},
		{Width: 80, HumanReview: true},
		{Width: 80, Active: true},
		{Width: 80, Historical: true},
		{Width: 80, Statuses: []string{"planned"}, Blocked: true},
		{Width: 80, Active: true, Historical: true},
		{Width: 80, Statuses: []string{"WIP"}},
		{Width: 80, Explain: true},
	}
	for _, view := range cases {
		got := Text(rep, view)
		if !strings.Contains(got, "recorded verified   1/2 (50%)") {
			t.Fatalf("view %+v changed the full summary\n%s", view, got)
		}
		if !strings.Contains(got, "SS_DEPENDENCY_HUMAN_REVIEW") {
			t.Fatalf("view %+v hid a diagnostic\n%s", view, got)
		}
		if strings.Contains(got, "\x1b") {
			t.Fatal("plain view emitted color")
		}
		if got != Text(rep, view) {
			t.Fatal("filtered text is not deterministic")
		}
	}
	if strings.Contains(full, "shown ") {
		t.Fatal("unfiltered text should not label a view")
	}
	historical := Text(rep, View{Width: 80, Historical: true})
	if !strings.Contains(historical, "shown 1 of 3 items") || !strings.Contains(historical, "001") || strings.Contains(historical, "\n002 ") {
		t.Fatalf("historical view\n%s", historical)
	}
	active := Text(rep, View{Width: 80, Active: true, Statuses: []string{"implemented"}})
	if !strings.Contains(active, "shown 0 of 3 items") || !strings.Contains(active, "no item matches this view") {
		t.Fatalf("no-match view\n%s", active)
	}
	if ExitCode(rep) != 0 {
		t.Fatal("a no-match view must not invent a finding")
	}
	both := Text(rep, View{Width: 80, Active: true, Historical: true})
	if !strings.Contains(both, "shown 0 of 3 items") {
		t.Fatalf("active and historical together\n%s", both)
	}
	or := Text(rep, View{Width: 80, Statuses: []string{"implemented", "verified"}})
	if !strings.Contains(or, "shown 2 of 3 items") || !strings.Contains(or, "001  implemented") || !strings.Contains(or, "002  verified") || strings.Contains(or, "003  planned") {
		t.Fatalf("status OR\n%s", or)
	}
}

func TestExplainCitesParsedLocationsAndSaysWhenMissing(t *testing.T) {
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
	text := Text(rep, View{Width: 80, Explain: true, HumanReview: true})
	if !strings.Contains(text, "shown 1 of 3 items") {
		t.Fatalf("explain filter\n%s", text)
	}
	if !strings.Contains(text, "explain: dependency 001 is human_review at .specify/memory/roadmap.md:18") {
		t.Fatalf("missing dependency location\n%s", text)
	}
	if strings.Contains(text, "explain: excluded") {
		t.Fatalf("human-review view should not explain the superseded row\n%s", text)
	}
	all := Text(rep, View{Width: 120, Explain: true})
	if !strings.Contains(all, "explain: excluded from active and recorded-verified counts; superseded-by 002 at .specify/memory/roadmap.md:6") {
		t.Fatalf("missing supersession location\n%s", all)
	}
	if !strings.Contains(all, "explain: unchecked checkbox at specs/001-original/tasks.md:1") {
		t.Fatalf("missing unchecked task location\n%s", all)
	}
	if !strings.Contains(all, "explain: checked checkbox at specs/002-replacement/tasks.md:1") {
		t.Fatalf("missing checked task location\n%s", all)
	}
	if !strings.Contains(all, "explain: no task file linked; checkbox line evidence unavailable") {
		t.Fatalf("planned item should say the task line is unavailable\n%s", all)
	}
	missing := explainLines(Item{
		ID:           "009",
		Path:         "bad\x1b[31m.md",
		SupersededBy: "010",
		Supersession: Supersession{Supplied: true, Valid: boolPtr(true), Line: 0},
		Tasks:        ItemTasks{State: "not_linked"},
		Dependencies: []Dependency{{ID: "001", State: "unsatisfied", Line: 0}},
	})
	joined := strings.Join(missing, "\n")
	if strings.Contains(joined, "\x1b") {
		t.Fatalf("explain kept a terminal control: %q", joined)
	}
	if !strings.Contains(joined, "field line unavailable") || !strings.Contains(joined, "source line unavailable") {
		t.Fatalf("missing locations were invented: %s", joined)
	}
	raw := string(JSON(rep))
	if strings.Contains(raw, "explain:") || strings.Contains(raw, "shown ") {
		t.Fatal("JSON gained view text")
	}
}

func TestHistoricalBlockedViewDoesNotShrinkActiveBlockers(t *testing.T) {
	body := "### 001 — Old [status: implemented]\n" +
		"- superseded-by: 002\n" +
		"- Depends on: 003\n\n" +
		"### 002 — New [status: planned]\n\n" +
		"### 003 — Later [status: planned]\n"
	dir := project(t, body, nil)
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if rep.Summary.BlockedCount != 0 {
		t.Fatalf("historical dependency must stay out of the active blocker count: %+v", rep.Summary)
	}
	text := Text(rep, View{Width: 80, Blocked: true, Historical: true})
	if !strings.Contains(text, "active blockers     0") || !strings.Contains(text, "shown 1 of 3 items") || !strings.Contains(text, "001") {
		t.Fatalf("historical blocked view\n%s", text)
	}
}

func TestUnknownStatusMatchesExactly(t *testing.T) {
	body := "### 001 — Odd [status: WIP]\n\n### 002 — Known [status: planned]\n"
	dir := project(t, body, nil)
	rep := Evaluate(Option{Root: dir, CWD: dir})
	match := Text(rep, View{Width: 80, Statuses: []string{"WIP"}})
	if !strings.Contains(match, "shown 1 of 2 items") || !strings.Contains(match, "WIP") {
		t.Fatalf("unknown token\n%s", match)
	}
	miss := Text(rep, View{Width: 80, Statuses: []string{"wip"}})
	if !strings.Contains(miss, "shown 0 of 2 items") || !strings.Contains(miss, "no item matches this view") {
		t.Fatalf("case fold\n%s", miss)
	}
	if ExitCode(rep) != 0 && rep.Summary.StatusCounts["WIP"] != 1 {
		t.Fatalf("raw status lost: %+v", rep.Summary)
	}
}

func boolPtr(v bool) *bool { return &v }
