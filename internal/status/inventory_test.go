package status

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Feature 006 US2: spec inventory (contracts/inventory-json.md).
// Expected documents are hand-counted in testdata/two-mode/*/expected-inventory.json.

var inventoryFixtures = []string{"no-roadmap", "brownfield", "no-specs", "with-roadmap"}

func expectedInventory(t *testing.T, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(twoModeRoot(name), "expected-inventory.json"))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, raw); err != nil {
		t.Fatal(err)
	}
	buf.WriteByte('\n')
	return buf.Bytes()
}

func TestInventoryFixturesMatchHandCount(t *testing.T) {
	for _, name := range inventoryFixtures {
		t.Run(name, func(t *testing.T) {
			rep := EvaluateInventory(InventoryOption{Root: twoModeRoot(name)})
			got := InventoryJSON(rep)
			want := expectedInventory(t, name)
			if !bytes.Equal(got, want) {
				t.Fatalf("inventory JSON differs\n got: %s\nwant: %s", got, want)
			}
			if code := InventoryExitCode(rep); code != 0 {
				t.Fatalf("exit %d, want 0", code)
			}
		})
	}
}

func TestInventoryFailOnInfo(t *testing.T) {
	for _, name := range []string{"no-roadmap", "no-specs"} {
		rep := EvaluateInventory(InventoryOption{Root: twoModeRoot(name), FailOn: "info"})
		if code := InventoryExitCode(rep); code != 1 {
			t.Fatalf("%s: --fail-on info exit %d, want 1", name, code)
		}
		if !strings.Contains(string(InventoryJSON(rep)), `"fail_on":"info"`) {
			t.Fatalf("%s: fail_on not reported", name)
		}
	}
	rep := EvaluateInventory(InventoryOption{Root: twoModeRoot("brownfield"), FailOn: "info"})
	if code := InventoryExitCode(rep); code != 0 {
		t.Fatalf("brownfield --fail-on info exit %d, want 0", code)
	}
}

func TestInventoryNeverReadsRoadmap(t *testing.T) {
	var read []string
	inventoryAfterRead = func(path string) {
		// The seam also fires after a directory listing; count file reads only.
		if !strings.HasSuffix(filepath.ToSlash(path), "/specs") {
			read = append(read, filepath.ToSlash(path))
		}
	}
	defer func() { inventoryAfterRead = nil }()
	rep := EvaluateInventory(InventoryOption{Root: twoModeRoot("with-roadmap")})
	if rep.RoadmapFile != "present" {
		t.Fatalf("roadmap_file %q", rep.RoadmapFile)
	}
	for _, p := range read {
		if strings.HasSuffix(p, "roadmap.md") {
			t.Fatalf("inventory read the roadmap: %v", read)
		}
	}
	// The brownfield checklist must never be read either.
	read = nil
	EvaluateInventory(InventoryOption{Root: twoModeRoot("brownfield")})
	for _, p := range read {
		if strings.Contains(p, "checklists") || strings.HasSuffix(p, "spec.md") || strings.HasSuffix(p, "plan.md") {
			t.Fatalf("inventory read %s", p)
		}
	}
	if len(read) != 1 {
		t.Fatalf("brownfield reads %v, want only tasks.md", read)
	}
}

func TestInventoryUnreliableTasksIsFatal(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".specify", ".gitkeep"), "")
	writeFile(t, filepath.Join(dir, "specs", "001-a", "tasks.md"), "# Tasks\n\n```\n- [ ] T001 hidden\n")
	rep := EvaluateInventory(InventoryOption{Root: dir})
	if rep.Outcome != "fatal" || rep.Fatal == nil || InventoryExitCode(rep) != 2 {
		t.Fatalf("want fatal, got %+v", rep.Fatal)
	}
	if rep.Fatal.Path != "specs/001-a/tasks.md" || rep.Fatal.Line == 0 {
		t.Fatalf("fatal must cite the tasks file and line: %+v", rep.Fatal)
	}
	var doc map[string]any
	if err := json.Unmarshal(InventoryJSON(rep), &doc); err != nil {
		t.Fatal(err)
	}
	if doc["summary"] != nil || len(doc["folders"].([]any)) != 0 || doc["error"] == nil || doc["complete"] != false {
		t.Fatalf("fatal document shape: %v", doc)
	}
	if ExitCode := InventoryExitCode(EvaluateInventory(InventoryOption{Root: dir, FailOn: "none"})); ExitCode != 2 {
		t.Fatal("--fail-on none masked a fatal")
	}
}

func TestInventoryRootDiscovery(t *testing.T) {
	root, err := filepath.Abs(twoModeRoot("brownfield"))
	if err != nil {
		t.Fatal(err)
	}
	rep := EvaluateInventory(InventoryOption{CWD: filepath.Join(root, "specs", "001-piranha-server-management-ui")})
	if rep.Fatal != nil || rep.Summary == nil || rep.Summary.FeatureFolders != 1 {
		t.Fatalf("root discovery from a subfolder: %+v", rep.Fatal)
	}
	missing := EvaluateInventory(InventoryOption{Root: t.TempDir()})
	if missing.Fatal == nil || missing.Fatal.Code != "SS_ROOT_UNAVAILABLE" || InventoryExitCode(missing) != 2 {
		t.Fatalf("root without .specify: %+v", missing.Fatal)
	}
}

func TestInventoryHumanText(t *testing.T) {
	rep := EvaluateInventory(InventoryOption{Root: twoModeRoot("no-roadmap")})
	text := InventoryText(rep, View{Width: 80})
	for _, want := range []string{
		"Spec Inventory", "inventory (roadmap not read; no status inferred)", "3/5 (60%)",
		"001-alpha", "20261101-093000-beta", "not a feature folder",
		"Checked boxes are not verification.", "SS_FOLDER_NOT_FEATURE",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("human text lacks %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "A roadmap exists") {
		t.Error("roadmap line shown without a roadmap")
	}
	withRoadmap := InventoryText(EvaluateInventory(InventoryOption{Root: twoModeRoot("with-roadmap")}), View{Width: 80})
	if !strings.Contains(withRoadmap, "A roadmap exists; spec status reports it.") {
		t.Errorf("roadmap line missing:\n%s", withRoadmap)
	}
	noSpecs := InventoryText(EvaluateInventory(InventoryOption{Root: twoModeRoot("no-specs")}), View{Width: 80})
	if !strings.Contains(noSpecs, "source specs (missing)") {
		t.Errorf("missing specs not shown:\n%s", noSpecs)
	}
	fatal := EvaluateInventory(InventoryOption{Root: t.TempDir()})
	if got := InventoryText(fatal, View{Width: 80}); !strings.HasPrefix(got, "spec: SS_ROOT_UNAVAILABLE: ") {
		t.Errorf("fatal text %q", got)
	}
}
