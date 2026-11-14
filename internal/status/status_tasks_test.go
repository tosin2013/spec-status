package status

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// Feature 007: status/tasks consistency diagnostics
// (specs/007-status-tasks-diagnostics/contracts/diagnostics.md).

const (
	msgAhead = "Recorded status is implemented or verified, but the linked tasks file has unchecked tasks."
	msgEarly = "Every task in the linked tasks file is checked, but the recorded status is before implemented."
)

// stItem is one roadmap item for these tests. tasks "" means no tasks.md;
// otherwise a string of x (checked) and . (unchecked), one checkbox each.
type stItem struct {
	id, status, specDir, tasks, extra string
	noTasks                           bool
}

func stProject(t *testing.T, items []stItem) (string, map[string]int) {
	t.Helper()
	var b strings.Builder
	files := map[string]string{}
	lines := map[string]int{}
	line := 1
	for _, it := range items {
		lines[it.id] = line
		body := fmt.Sprintf("### %s — Item %s [status: %s]\n- Description: d\n- Outcome: o\n- Scope: s\n", it.id, it.id, it.status)
		if it.specDir != "" {
			body += "- Spec dir: " + it.specDir + "\n"
			files[it.specDir+"/spec.md"] = "# Spec\n"
			files[it.specDir+"/plan.md"] = "# Plan\n"
			if !it.noTasks {
				var tb strings.Builder
				tb.WriteString("# Tasks\n\n")
				for i, c := range it.tasks {
					mark := " "
					if c == 'x' {
						mark = "x"
					}
					fmt.Fprintf(&tb, "- [%s] T%03d task\n", mark, i+1)
				}
				files[it.specDir+"/tasks.md"] = tb.String()
			}
		}
		body += it.extra + "\n"
		b.WriteString(body)
		line += strings.Count(body, "\n")
	}
	return project(t, b.String(), files), lines
}

func stFindings(rep Report) []Diagnostic {
	var out []Diagnostic
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_STATUS_AHEAD_OF_TASKS" || d.Code == "SS_TASKS_DONE_STATUS_EARLY" {
			out = append(out, d)
		}
	}
	return out
}

func hasCode(rep Report, code, item string) bool {
	for _, d := range rep.Diagnostics {
		if d.Code == code && d.ItemID == item {
			return true
		}
	}
	return false
}

func TestStatusTasksDecisionTable(t *testing.T) {
	cases := []struct {
		name  string
		item  stItem
		code  string // "" means no finding
		rLine int    // expected related line in tasks.md (first checkbox is line 3)
	}{
		{"implemented 1/3", stItem{status: "implemented", tasks: "x.."}, "SS_STATUS_AHEAD_OF_TASKS", 4},
		{"verified 1/2", stItem{status: "verified", tasks: "x."}, "SS_STATUS_AHEAD_OF_TASKS", 4},
		{"verified 0/2", stItem{status: "verified", tasks: ".."}, "SS_STATUS_AHEAD_OF_TASKS", 3},
		{"implemented 3/3", stItem{status: "implemented", tasks: "xxx"}, "", 0},
		{"implemented 0/0", stItem{status: "implemented", tasks: ""}, "", 0},
		{"in-progress 2/2", stItem{status: "in-progress", tasks: "xx"}, "SS_TASKS_DONE_STATUS_EARLY", 3},
		{"planned 3/3", stItem{status: "planned", tasks: "xxx"}, "SS_TASKS_DONE_STATUS_EARLY", 3},
		{"undecided 1/1", stItem{status: "undecided", tasks: "x"}, "SS_TASKS_DONE_STATUS_EARLY", 3},
		{"needs-info 1/1", stItem{status: "needs-info", tasks: "x"}, "SS_TASKS_DONE_STATUS_EARLY", 3},
		{"specced 1/2", stItem{status: "specced", tasks: "x."}, "", 0},
		{"planned 0/0", stItem{status: "planned", tasks: ""}, "", 0},
		{"deferred 2/2", stItem{status: "deferred", tasks: "xx"}, "", 0},
		{"abandoned 1/3", stItem{status: "abandoned", tasks: "x.."}, "", 0},
		{"unknown done 2/2", stItem{status: "done", tasks: "xx"}, "", 0},
		{"implemented, tasks.md missing", stItem{status: "implemented", noTasks: true}, "", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.item.id, c.item.specDir = "001", "specs/001-x"
			dir, lines := stProject(t, []stItem{c.item})
			rep := Evaluate(Option{Root: dir, CWD: dir})
			if rep.Fatal != nil {
				t.Fatalf("fatal %+v", rep.Fatal)
			}
			got := stFindings(rep)
			if c.code == "" {
				if len(got) != 0 {
					t.Fatalf("unexpected findings %+v", got)
				}
				return
			}
			if len(got) != 1 {
				t.Fatalf("findings %+v, want one %s", got, c.code)
			}
			d := got[0]
			wantSev, wantMsg := "warning", msgAhead
			if c.code == "SS_TASKS_DONE_STATUS_EARLY" {
				wantSev, wantMsg = "info", msgEarly
			}
			if d.Code != c.code || d.Severity != wantSev || d.Message != wantMsg || d.ItemID != "001" {
				t.Fatalf("finding %+v", d)
			}
			if d.Path != ".specify/memory/roadmap.md" || d.Line != lines["001"] {
				t.Fatalf("source %s:%d, want roadmap:%d", d.Path, d.Line, lines["001"])
			}
			if len(d.Related) != 1 || d.Related[0].Path != "specs/001-x/tasks.md" || d.Related[0].Line != c.rLine {
				t.Fatalf("related %+v, want specs/001-x/tasks.md:%d", d.Related, c.rLine)
			}
		})
	}
}

func TestStatusTasksExclusions(t *testing.T) {
	t.Run("unknown status keeps only SS_STATUS_UNKNOWN", func(t *testing.T) {
		dir, _ := stProject(t, []stItem{{id: "001", status: "done", specDir: "specs/001-x", tasks: "xx"}})
		rep := Evaluate(Option{Root: dir, CWD: dir})
		if !hasCode(rep, "SS_STATUS_UNKNOWN", "001") || len(stFindings(rep)) != 0 {
			t.Fatalf("diagnostics %+v", rep.Diagnostics)
		}
	})
	t.Run("missing tasks.md keeps only the stage finding", func(t *testing.T) {
		dir, _ := stProject(t, []stItem{{id: "001", status: "implemented", specDir: "specs/001-x", noTasks: true}})
		rep := Evaluate(Option{Root: dir, CWD: dir})
		if !hasCode(rep, "SS_STAGE_FILE_MISSING", "001") || len(stFindings(rep)) != 0 {
			t.Fatalf("diagnostics %+v", rep.Diagnostics)
		}
	})
	t.Run("validly superseded item is history", func(t *testing.T) {
		dir, _ := stProject(t, []stItem{
			{id: "001", status: "implemented", specDir: "specs/001-x", tasks: "x..", extra: "- superseded-by: 002"},
			{id: "002", status: "planned"},
		})
		rep := Evaluate(Option{Root: dir, CWD: dir})
		if rep.Fatal != nil || len(stFindings(rep)) != 0 {
			t.Fatalf("superseded: fatal %+v findings %+v", rep.Fatal, stFindings(rep))
		}
	})
	t.Run("two items sharing one tasks file each get a finding", func(t *testing.T) {
		dir, _ := stProject(t, []stItem{
			{id: "001", status: "implemented", specDir: "specs/001-x", tasks: "x."},
			{id: "002", status: "verified", specDir: "specs/001-x", tasks: "x."},
		})
		rep := Evaluate(Option{Root: dir, CWD: dir})
		if !hasCode(rep, "SS_STATUS_AHEAD_OF_TASKS", "001") || !hasCode(rep, "SS_STATUS_AHEAD_OF_TASKS", "002") || len(stFindings(rep)) != 2 {
			t.Fatalf("shared: %+v", stFindings(rep))
		}
	})
	t.Run("summary counts are unchanged by the findings", func(t *testing.T) {
		dir, _ := stProject(t, []stItem{{id: "001", status: "implemented", specDir: "specs/001-x", tasks: "x.."}})
		rep := Evaluate(Option{Root: dir, CWD: dir})
		s := rep.Summary
		if s.TotalItems != 1 || s.ActiveItems != 1 || s.VerifiedActiveItems != 0 || s.BlockedCount != 0 || s.Tasks.Checked != 1 || s.Tasks.Total != 3 {
			t.Fatalf("summary %+v", s)
		}
	})
}

func TestStatusTasksFixtures(t *testing.T) {
	skills := func(name string) string { return filepath.Join("..", "..", "testdata", "skills", name) }
	t.Run("implemented-unchecked: ahead, warning thresholds", func(t *testing.T) {
		rep := Evaluate(Option{Root: skills("implemented-unchecked")})
		got := stFindings(rep)
		if len(got) != 1 || got[0].Code != "SS_STATUS_AHEAD_OF_TASKS" || got[0].ItemID != "001" {
			t.Fatalf("findings %+v", got)
		}
		if len(got[0].Related) != 1 || got[0].Related[0].Path != "specs/001-import/tasks.md" || got[0].Related[0].Line == 0 {
			t.Fatalf("related %+v", got[0].Related)
		}
		for failOn, want := range map[string]int{"error": 0, "warning": 1, "info": 1, "none": 0} {
			if code := ExitCode(Evaluate(Option{Root: skills("implemented-unchecked"), FailOn: failOn})); code != want {
				t.Errorf("--fail-on %s exit %d, want %d", failOn, code, want)
			}
		}
		if !strings.Contains(string(JSON(rep)), `"code":"SS_STATUS_AHEAD_OF_TASKS","severity":"warning"`) {
			t.Error("JSON lacks the diagnostic")
		}
	})
	t.Run("tasks-done-status-planned: early, info thresholds", func(t *testing.T) {
		rep := Evaluate(Option{Root: skills("tasks-done-status-planned")})
		got := stFindings(rep)
		if len(got) != 1 || got[0].Code != "SS_TASKS_DONE_STATUS_EARLY" || got[0].ItemID != "001" {
			t.Fatalf("findings %+v", got)
		}
		for failOn, want := range map[string]int{"error": 0, "warning": 0, "info": 1} {
			if code := ExitCode(Evaluate(Option{Root: skills("tasks-done-status-planned"), FailOn: failOn})); code != want {
				t.Errorf("--fail-on %s exit %d, want %d", failOn, code, want)
			}
		}
	})
	t.Run("roles-mixed: 002 ahead, 003 early", func(t *testing.T) {
		rep := Evaluate(Option{Root: skills("roles-mixed")})
		if !hasCode(rep, "SS_STATUS_AHEAD_OF_TASKS", "002") || !hasCode(rep, "SS_TASKS_DONE_STATUS_EARLY", "003") || len(stFindings(rep)) != 2 {
			t.Fatalf("findings %+v", stFindings(rep))
		}
	})
}

// TestStatusTasksRelatedSkipsLinelessChecks: related evidence uses the first
// matching checkbox that has a recorded line (Doomer lead, feature 007).
func TestStatusTasksRelatedSkipsLinelessChecks(t *testing.T) {
	it := Item{ID: "001", Status: "planned", Active: true, Line: 7, Tasks: ItemTasks{
		State: "read", Checked: 2, Total: 2, Source: "specs/001-x/tasks.md",
		Checks: []TaskCheck{{Line: 0, Checked: true}, {Line: 5, Checked: true}},
	}}
	d := statusTasksDiag(it, ".specify/memory/roadmap.md")
	if d == nil || len(d.Related) != 1 || d.Related[0].Line != 5 {
		t.Fatalf("early related %+v", d)
	}
	it.Status, it.Tasks.Checked = "implemented", 1
	it.Tasks.Checks = []TaskCheck{{Line: 0, Checked: false}, {Line: 4, Checked: true}, {Line: 6, Checked: false}}
	d = statusTasksDiag(it, ".specify/memory/roadmap.md")
	if d == nil || len(d.Related) != 1 || d.Related[0].Line != 6 {
		t.Fatalf("ahead related %+v", d)
	}
}
