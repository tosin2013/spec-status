package conformance

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var (
	driftProject   = flag.String("project", "", "drift check: upstream project, as named in testdata/upstream/pins.json")
	driftCandidate = flag.String("candidate", "", "drift check: path to a downloaded candidate template")
	driftTemplate  = flag.Int("template", 0, "drift check: index of the project's template in testdata/upstream/pins.json")
)

// pinsFile is the single list of upstream pins, shared with
// scripts/check-upstream-drift.sh and the drift workflow (feature 013).
const pinsFile = "testdata/upstream/pins.json"

type pinTemplate struct {
	Path   string `json:"path"`   // repository path at the tag
	Pinned string `json:"pinned"` // copy under testdata/upstream/
}

type pin struct {
	Name      string        `json:"name"`
	Repo      string        `json:"repo"`
	Tag       string        `json:"tag"`
	Templates []pinTemplate `json:"templates"`
}

// pinnedTemplates maps each project to its current pin's template files,
// read from pinsFile.
var pinnedTemplates = func() map[string][]string {
	pins, err := loadPins()
	if err != nil {
		panic(err)
	}
	m := map[string][]string{}
	for _, p := range pins {
		for _, t := range p.Templates {
			m[p.Name] = append(m[p.Name], t.Pinned)
		}
	}
	return m
}()

func loadPins() ([]pin, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot, pinsFile))
	if err != nil {
		return nil, err
	}
	var doc struct {
		Projects []pin `json:"projects"`
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("%s: %w", pinsFile, err)
	}
	return doc.Projects, nil
}

var repoPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)

// TestPinsFile checks the shared pin list against the pinned copies.
func TestPinsFile(t *testing.T) {
	pins, err := loadPins()
	if err != nil {
		t.Fatal(err)
	}
	if len(pins) == 0 {
		t.Fatal("no pins")
	}
	seen := map[string]bool{}
	for _, p := range pins {
		if p.Name == "" || seen[p.Name] {
			t.Errorf("empty or duplicate project name %q", p.Name)
		}
		seen[p.Name] = true
		if !repoPattern.MatchString(p.Repo) {
			t.Errorf("%s: repo %q is not owner/name", p.Name, p.Repo)
		}
		if p.Tag == "" || len(p.Templates) == 0 {
			t.Errorf("%s: tag and templates are required", p.Name)
		}
		if _, err := os.Stat(filepath.Join(repoRoot, "testdata", "upstream", p.Name, "resolutions.md")); err != nil {
			t.Errorf("%s: no resolutions.md: %v", p.Name, err)
		}
		for _, tm := range p.Templates {
			if !strings.HasPrefix(tm.Pinned, p.Name+"/"+p.Tag+"/") {
				t.Errorf("%s: pinned copy %q is not under %s/%s/", p.Name, tm.Pinned, p.Name, p.Tag)
			}
			if tm.Path == "" || strings.HasPrefix(tm.Path, "/") || strings.Contains(tm.Path, "..") {
				t.Errorf("%s: repository path %q must be relative", p.Name, tm.Path)
			}
			if _, err := os.Stat(filepath.Join(repoRoot, "testdata", "upstream", filepath.FromSlash(tm.Pinned))); err != nil {
				t.Errorf("%s: pinned copy missing: %v", p.Name, err)
			}
		}
	}
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

// pinnedShapes returns the shape keys of a project's pinned template; index
// selects the template in pins.json, so each candidate is compared with its
// own pinned copy.
func pinnedShapes(t *testing.T, project string, index int) map[string]shape {
	t.Helper()
	files, ok := pinnedTemplates[project]
	if !ok {
		names := make([]string, 0, len(pinnedTemplates))
		for n := range pinnedTemplates {
			names = append(names, n)
		}
		sort.Strings(names)
		t.Fatalf("unknown project %q; want one of %s (testdata/upstream/pins.json)", project, strings.Join(names, ", "))
	}
	if index < 0 || index >= len(files) {
		t.Fatalf("project %q has %d template(s); -template=%d is out of range", project, len(files), index)
	}
	data, err := os.ReadFile(filepath.Join(repoRoot, "testdata", "upstream", files[index]))
	if err != nil {
		t.Fatal(err)
	}
	return shapeKeys(string(data))
}

// TestDrift compares a downloaded candidate template with the current pin
// (FR-008, FR-009). It reads local files only and is skipped unless
// -candidate is given:
//
//	go test ./internal/conformance -run TestDrift -v -args -project=spec-kit -candidate=PATH [-template=N]
func TestDrift(t *testing.T) {
	if *driftCandidate == "" {
		t.Skip("no -candidate given; see specs/002-upstream-format-conformance/contracts/drift-report.md")
	}
	data, err := os.ReadFile(*driftCandidate)
	if err != nil {
		t.Fatal(err)
	}
	changes := compareShapes(pinnedShapes(t, *driftProject, *driftTemplate), shapeKeys(string(data)))
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
		for _, ch := range compareShapes(pinnedShapes(t, c.project, 0), shapeKeys(string(data))) {
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
