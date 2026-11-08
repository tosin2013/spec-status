// Package conformance holds test-only checks for feature 002: hand-counted
// conformance fixtures, the pinned shape inventory, and the upstream drift
// check. Every file in this package is a _test.go file, so none of it is
// compiled into the spec binary (constitution Article II).
package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/tosin2013/spec-status/internal/status"
)

// expected is a fixture's hand-counted result. Its format is
// specs/002-upstream-format-conformance/contracts/expected-result.md. Only
// the fields present are compared.
type expected struct {
	Pins        []string       `json:"pins"`
	Origin      string         `json:"origin"`
	CountedBy   string         `json:"counted_by"`
	Outcome     string         `json:"outcome"`
	Exit        *int           `json:"exit"`
	Items       []expectedItem `json:"items"`
	Summary     *expectedSum   `json:"summary"`
	Diagnostics []expectedDiag `json:"diagnostics"`
	Fatal       *expectedFatal `json:"fatal"`
}

type expectedItem struct {
	ID           string         `json:"id"`
	Status       string         `json:"status"`
	Active       *bool          `json:"active"`
	Dependencies *[]expectedDep `json:"dependencies"`
	Tasks        *expectedTasks `json:"tasks"`
	// SpecDir and SupersededBy are compared when present; JSON null means
	// "none". They check rule (b) of D-2026-11-07-09 directly.
	SpecDir      json.RawMessage `json:"spec_dir"`
	SupersededBy json.RawMessage `json:"superseded_by"`
}

type expectedDep struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

type expectedTasks struct {
	Checked   *int `json:"checked"`
	Total     *int `json:"total"`
	FilesRead *int `json:"files_read"`
}

type expectedSum struct {
	TotalItems          *int           `json:"total_items"`
	ActiveItems         *int           `json:"active_items"`
	VerifiedActiveItems *int           `json:"verified_active_items"`
	BlockedCount        *int           `json:"blocked_count"`
	HumanReviewCount    *int           `json:"human_review_count"`
	Tasks               *expectedTasks `json:"tasks"`
}

type expectedDiag struct {
	Code string `json:"code"`
	Item string `json:"item"`
}

type expectedFatal struct {
	Code string `json:"code"`
	Path string `json:"path"`
	Line *int   `json:"line"`
}

func loadExpected(path string) (expected, error) {
	var e expected
	data, err := os.ReadFile(path)
	if err != nil {
		return e, err
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&e); err != nil {
		return e, fmt.Errorf("%s: %w", path, err)
	}
	if e.Outcome == "" || e.Exit == nil {
		return e, fmt.Errorf("%s: \"outcome\" and \"exit\" are required", path)
	}
	if e.CountedBy == "" {
		return e, fmt.Errorf("%s: \"counted_by\" is required; expected results are hand counts", path)
	}
	return e, nil
}

// compare returns every difference between a hand count and a report.
func compare(e expected, rep status.Report) []string {
	var diffs []string
	add := func(format string, args ...any) { diffs = append(diffs, fmt.Sprintf(format, args...)) }
	if rep.Outcome != e.Outcome {
		add("outcome: got %q, want %q", rep.Outcome, e.Outcome)
	}
	if got := status.ExitCode(rep); got != *e.Exit {
		add("exit: got %d, want %d", got, *e.Exit)
	}
	if e.Outcome == "fatal" {
		if e.Fatal != nil {
			if rep.Fatal == nil {
				add("fatal: missing")
			} else {
				if rep.Fatal.Code != e.Fatal.Code {
					add("fatal code: got %s, want %s", rep.Fatal.Code, e.Fatal.Code)
				}
				if e.Fatal.Path != "" && rep.Fatal.Path != e.Fatal.Path {
					add("fatal path: got %s, want %s", rep.Fatal.Path, e.Fatal.Path)
				}
				if e.Fatal.Line != nil && rep.Fatal.Line != *e.Fatal.Line {
					add("fatal line: got %d, want %d", rep.Fatal.Line, *e.Fatal.Line)
				}
			}
		}
		return diffs
	}
	if e.Items != nil {
		if len(rep.Items) != len(e.Items) {
			add("items: got %d, want %d", len(rep.Items), len(e.Items))
		}
		for i := 0; i < len(e.Items) && i < len(rep.Items); i++ {
			diffs = append(diffs, compareItem(e.Items[i], rep.Items[i])...)
		}
	}
	if e.Summary != nil && rep.Summary != nil {
		s, w := rep.Summary, e.Summary
		intEq(&diffs, "summary.total_items", s.TotalItems, w.TotalItems)
		intEq(&diffs, "summary.active_items", s.ActiveItems, w.ActiveItems)
		intEq(&diffs, "summary.verified_active_items", s.VerifiedActiveItems, w.VerifiedActiveItems)
		intEq(&diffs, "summary.blocked_count", s.BlockedCount, w.BlockedCount)
		intEq(&diffs, "summary.human_review_count", s.HumanReviewCount, w.HumanReviewCount)
		if w.Tasks != nil {
			intEq(&diffs, "summary.tasks.checked", s.Tasks.Checked, w.Tasks.Checked)
			intEq(&diffs, "summary.tasks.total", s.Tasks.Total, w.Tasks.Total)
			intEq(&diffs, "summary.tasks.files_read", s.Tasks.FilesRead, w.Tasks.FilesRead)
		}
	} else if e.Summary != nil {
		add("summary: missing from report")
	}
	if e.Diagnostics != nil {
		got := map[string]int{}
		for _, d := range rep.Diagnostics {
			got[d.ItemID+" "+d.Code]++
		}
		want := map[string]int{}
		for _, d := range e.Diagnostics {
			want[d.Item+" "+d.Code]++
		}
		for _, k := range sortedKeys(want, got) {
			if got[k] != want[k] {
				add("diagnostic %q: got %d, want %d", k, got[k], want[k])
			}
		}
	}
	return diffs
}

func compareItem(w expectedItem, it status.Item) []string {
	var diffs []string
	add := func(format string, args ...any) {
		diffs = append(diffs, fmt.Sprintf("item %s: ", w.ID)+fmt.Sprintf(format, args...))
	}
	if it.ID != w.ID {
		add("id: got %s (items compare in document order)", it.ID)
		return diffs
	}
	if it.Status != w.Status {
		add("status: got %q, want %q", it.Status, w.Status)
	}
	if w.Active != nil && it.Active != *w.Active {
		add("active: got %v, want %v", it.Active, *w.Active)
	}
	if w.Dependencies != nil {
		var got, want []string
		for _, d := range it.Dependencies {
			got = append(got, d.ID+"="+d.State)
		}
		for _, d := range *w.Dependencies {
			want = append(want, d.ID+"="+d.State)
		}
		if strings.Join(got, ",") != strings.Join(want, ",") {
			add("dependencies: got [%s], want [%s]", strings.Join(got, ","), strings.Join(want, ","))
		}
	}
	if want, ok, err := optionalString(w.SpecDir); err != nil {
		add("spec_dir: %v", err)
	} else if ok && it.SpecDir != want {
		add("spec_dir: got %q, want %q", it.SpecDir, want)
	}
	if want, ok, err := optionalString(w.SupersededBy); err != nil {
		add("superseded_by: %v", err)
	} else if ok && it.SupersededBy != want {
		add("superseded_by: got %q, want %q", it.SupersededBy, want)
	}
	if w.Tasks != nil {
		if w.Tasks.Checked != nil && it.Tasks.Checked != *w.Tasks.Checked {
			add("tasks.checked: got %d, want %d", it.Tasks.Checked, *w.Tasks.Checked)
		}
		if w.Tasks.Total != nil && it.Tasks.Total != *w.Tasks.Total {
			add("tasks.total: got %d, want %d", it.Tasks.Total, *w.Tasks.Total)
		}
	}
	return diffs
}

func intEq(diffs *[]string, name string, got int, want *int) {
	if want != nil && got != *want {
		*diffs = append(*diffs, fmt.Sprintf("%s: got %d, want %d", name, got, *want))
	}
}

func sortedKeys(maps ...map[string]int) []string {
	seen := map[string]bool{}
	var keys []string
	for _, m := range maps {
		for k := range m {
			if !seen[k] {
				seen[k] = true
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)
	return keys
}

// optionalString decodes an optional JSON string; null decodes to "".
func optionalString(raw json.RawMessage) (string, bool, error) {
	if len(raw) == 0 {
		return "", false, nil
	}
	if string(raw) == "null" {
		return "", true, nil
	}
	var v string
	err := json.Unmarshal(raw, &v)
	return v, true, err
}
