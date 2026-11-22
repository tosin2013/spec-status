package human

import (
	"github.com/tosin2013/spec-status/internal/status"
	"strings"
	"testing"
)

func TestInventoryHumanText(t *testing.T) {
	rep := status.EvaluateInventory(status.InventoryOption{Root: twoModeRoot("no-roadmap")})
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
	withRoadmap := InventoryText(status.EvaluateInventory(status.InventoryOption{Root: twoModeRoot("with-roadmap")}), View{Width: 80})
	if !strings.Contains(withRoadmap, "A roadmap exists; spec status reports it.") {
		t.Errorf("roadmap line missing:\n%s", withRoadmap)
	}
	noSpecs := InventoryText(status.EvaluateInventory(status.InventoryOption{Root: twoModeRoot("no-specs")}), View{Width: 80})
	if !strings.Contains(noSpecs, "source specs (missing)") {
		t.Errorf("missing specs not shown:\n%s", noSpecs)
	}
	fatal := status.EvaluateInventory(status.InventoryOption{Root: t.TempDir()})
	if got := InventoryText(fatal, View{Width: 80}); !strings.HasPrefix(got, "spec: SS_ROOT_UNAVAILABLE: ") {
		t.Errorf("fatal text %q", got)
	}
}
