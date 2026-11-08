package conformance

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

var (
	driftProject   = flag.String("project", "", "drift check: upstream project (spec-kit or speckit-roadmap)")
	driftCandidate = flag.String("candidate", "", "drift check: path to a downloaded candidate template")
)

// pinnedTemplates maps each project to its current pin's template files.
var pinnedTemplates = map[string][]string{
	"spec-kit":        {"spec-kit/v1.2.0/templates/tasks-template.md"},
	"speckit-roadmap": {"speckit-roadmap/v0.1.0/templates/roadmap-template.md"},
}

const inventoryPath = "specs/002-upstream-format-conformance/contracts/shape-inventory.md"

type inventoryRow struct {
	key, status, fixtures, section string
}

// loadInventory reads the shape inventory tables, keyed by section and key.
func loadInventory(t *testing.T) []inventoryRow {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot, inventoryPath))
	if err != nil {
		t.Fatal(err)
	}
	var rows []inventoryRow
	section := ""
	for _, line := range strings.Split(string(data), "\n") {
		switch {
		case strings.HasPrefix(line, "## Roadmap"):
			section = "roadmap"
		case strings.HasPrefix(line, "## Tasks"):
			section = "tasks"
		case strings.HasPrefix(line, "## "):
			section = ""
		}
		if section == "" || !strings.HasPrefix(line, "| `") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		if len(cells) < 5 {
			t.Fatalf("inventory row has %d cells: %s", len(cells), line)
		}
		rows = append(rows, inventoryRow{
			key:      strings.Trim(strings.TrimSpace(cells[0]), "`"),
			status:   strings.TrimSpace(cells[len(cells)-2]),
			fixtures: strings.TrimSpace(cells[len(cells)-1]),
			section:  section,
		})
	}
	return rows
}

// fixtureKeys returns, for each fixture, the shape keys its roadmap and tasks
// files use, by section.
func fixtureKeys(t *testing.T) map[string]map[string]map[string]bool {
	t.Helper()
	out := map[string]map[string]map[string]bool{}
	for _, dir := range fixtureDirs(t) {
		name := filepath.Base(dir)
		out[name] = map[string]map[string]bool{"roadmap": {}, "tasks": {}}
		_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			section := ""
			switch info.Name() {
			case "roadmap.md":
				section = "roadmap"
			case "tasks.md":
				section = "tasks"
			default:
				return nil
			}
			data, rerr := os.ReadFile(p)
			if rerr != nil {
				return rerr
			}
			for k := range shapeKeys(string(data)) {
				out[name][section][k] = true
			}
			return nil
		})
	}
	return out
}

// TestPinnedInventory checks the shape inventory against the pinned
// templates and the conformance fixtures (FR-003, SC-002).
func TestPinnedInventory(t *testing.T) {
	rows := loadInventory(t)
	listed := map[string]inventoryRow{}
	for _, r := range rows {
		listed[r.section+" "+r.key] = r
	}
	sectionOf := map[string]string{"spec-kit": "tasks", "speckit-roadmap": "roadmap"}
	for project, files := range pinnedTemplates {
		for _, f := range files {
			data, err := os.ReadFile(filepath.Join(repoRoot, "testdata", "upstream", f))
			if err != nil {
				t.Fatal(err)
			}
			for _, k := range sortedShapeKeys(shapeKeys(string(data))) {
				r, ok := listed[sectionOf[project]+" "+k]
				if !ok {
					t.Errorf("%s: pinned template uses %q, which the inventory does not list", f, k)
				} else if r.status != "supported" {
					t.Errorf("%s: %q is produced by a pinned template but listed as %q", f, k, r.status)
				}
			}
		}
	}
	fx := fixtureKeys(t)
	for name, sections := range fx {
		for section, keys := range sections {
			for k := range keys {
				if _, ok := listed[section+" "+k]; !ok {
					t.Errorf("fixture %s uses %s shape %q, which the inventory does not list", name, section, k)
				}
			}
		}
	}
	for _, r := range rows {
		names := fixtureClaims(r, fx)
		if names == nil {
			continue
		}
		if len(names) == 0 && (r.status == "supported" || r.status == "hand-edit") {
			t.Errorf("inventory %s %q (%s) names no fixture", r.section, r.key, r.status)
		}
		for _, n := range names {
			if fx[n] == nil {
				t.Errorf("inventory %s %q names unknown fixture %q", r.section, r.key, n)
			} else if !fx[n][r.section][r.key] {
				t.Errorf("inventory %s %q claims fixture %q, which does not use it", r.section, r.key, n)
			}
		}
	}
}

// fixtureClaims expands a row's Fixtures cell. It returns nil when the row
// defers to other tests.
func fixtureClaims(r inventoryRow, fx map[string]map[string]map[string]bool) []string {
	cell := r.fixtures
	if strings.HasPrefix(cell, "covered by") {
		return nil
	}
	if strings.HasPrefix(cell, "all ") {
		var names []string
		for n, s := range fx {
			if len(s[r.section]) > 0 && (r.section != "tasks" || s["tasks"]["task-line"] || s["tasks"]["task-checked"]) {
				names = append(names, n)
			}
		}
		sort.Strings(names)
		return names
	}
	var names []string
	for _, n := range strings.Split(cell, ",") {
		if n = strings.TrimSpace(n); n != "" {
			names = append(names, n)
		}
	}
	return names
}

// driftChange is one difference between a pinned and a candidate template.
type driftChange struct {
	Kind, Key string
	Example   shape
}

func (c driftChange) String() string {
	return fmt.Sprintf("%s: %s", c.Kind, c.Key)
}

// compareShapes reports shape keys added in, or removed from, a candidate.
func compareShapes(pinned, candidate map[string]shape) []driftChange {
	var out []driftChange
	for _, k := range sortedShapeKeys(candidate) {
		if _, ok := pinned[k]; !ok {
			out = append(out, driftChange{"added", k, candidate[k]})
		}
	}
	for _, k := range sortedShapeKeys(pinned) {
		if _, ok := candidate[k]; !ok {
			out = append(out, driftChange{"removed", k, pinned[k]})
		}
	}
	return out
}

// resolved reads a project's resolutions.md and returns the changes that
// have both a fixture and a decision ID.
func resolved(t *testing.T, project string) map[string]bool {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot, "testdata", "upstream", project, "resolutions.md"))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		cells := strings.Split(strings.Trim(line, "| "), "|")
		if len(cells) != 3 || strings.HasPrefix(strings.TrimSpace(cells[0]), "shape change") || strings.HasPrefix(strings.TrimSpace(cells[0]), "---") {
			continue
		}
		change, fixture, decision := strings.TrimSpace(cells[0]), strings.TrimSpace(cells[1]), strings.TrimSpace(cells[2])
		if change != "" && fixture != "" && decision != "" {
			out[strings.Trim(change, "`")] = true
		}
	}
	return out
}

func pinnedShapes(t *testing.T, project string) map[string]shape {
	t.Helper()
	files, ok := pinnedTemplates[project]
	if !ok {
		t.Fatalf("unknown project %q; want spec-kit or speckit-roadmap", project)
	}
	data, err := os.ReadFile(filepath.Join(repoRoot, "testdata", "upstream", files[0]))
	if err != nil {
		t.Fatal(err)
	}
	return shapeKeys(string(data))
}

// TestDrift compares a downloaded candidate template with the current pin
// (FR-008, FR-009). It reads local files only and is skipped unless
// -candidate is given:
//
//	go test ./internal/conformance -run TestDrift -v -args -project=spec-kit -candidate=PATH
func TestDrift(t *testing.T) {
	if *driftCandidate == "" {
		t.Skip("no -candidate given; see specs/002-upstream-format-conformance/contracts/drift-report.md")
	}
	data, err := os.ReadFile(*driftCandidate)
	if err != nil {
		t.Fatal(err)
	}
	changes := compareShapes(pinnedShapes(t, *driftProject), shapeKeys(string(data)))
	done := resolved(t, *driftProject)
	unresolved := 0
	for _, c := range changes {
		state := "unresolved"
		if done[c.String()] {
			state = "resolved"
		} else {
			unresolved++
		}
		t.Logf("%-8s %-30s line %d: %s (%s)", c.Kind+":", c.Key, c.Example.Line, strings.TrimSpace(c.Example.Text), state)
	}
	if unresolved > 0 {
		t.Errorf("%d unresolved shape changes; the %s pin cannot move until each has a fixture and a decision in testdata/upstream/%s/resolutions.md", unresolved, *driftProject, *driftProject)
	}
}

// TestDriftDetectsSimulatedChange runs the comparison on modified copies of
// the pinned templates (SC-004).
func TestDriftDetectsSimulatedChange(t *testing.T) {
	cases := []struct {
		project, file string
		want          []string
	}{
		{"spec-kit", "tasks-template-changed.md", []string{"added: task-marker:[P1]"}},
		{"speckit-roadmap", "roadmap-template-changed.md", []string{"added: field:Spec directory", "removed: field:Spec dir"}},
	}
	for _, c := range cases {
		data, err := os.ReadFile(filepath.Join(repoRoot, "testdata", "conformance-drift", c.file))
		if err != nil {
			t.Fatal(err)
		}
		got := map[string]bool{}
		for _, ch := range compareShapes(pinnedShapes(t, c.project), shapeKeys(string(data))) {
			got[ch.String()] = true
		}
		done := resolved(t, c.project)
		for _, w := range c.want {
			if !got[w] {
				t.Errorf("%s: drift check did not report %q (got %v)", c.file, w, got)
			}
			if done[w] {
				t.Errorf("%s: %q is unexpectedly marked resolved", c.file, w)
			}
		}
	}
}
