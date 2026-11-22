package human

import (
	"github.com/tosin2013/spec-status/internal/status"
	"path/filepath"
	"strings"
	"testing"
)

func TestInventoryDeterminism(t *testing.T) {
	for _, name := range inventoryFixtures {
		root := twoModeRoot(name)
		a, b := status.EvaluateInventory(status.InventoryOption{Root: root}), status.EvaluateInventory(status.InventoryOption{Root: root})
		ja, jb := string(status.InventoryJSON(a)), string(status.InventoryJSON(b))
		if ja != jb {
			t.Fatalf("%s: status.JSON differs between runs", name)
		}
		ta, tb := InventoryText(a, View{Width: 80}), InventoryText(b, View{Width: 80})
		if ta != tb {
			t.Fatalf("%s: text differs between runs", name)
		}
		abs, _ := filepath.Abs(root)
		if strings.Contains(ja, abs) || strings.Contains(ta, abs) || strings.Contains(ja, "as_of") || strings.Contains(ja, "timestamp\":") {
			t.Fatalf("%s: absolute path or timestamp in output", name)
		}
	}
}
