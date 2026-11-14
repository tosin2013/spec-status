package status

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/tosin2013/spec-status/internal/safeio"
)

// Feature 006 US3: inventory safety and determinism. Symlinks, FIFOs,
// sensitive names and limits are created at test time.

func invProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".specify", ".gitkeep"), "")
	return dir
}

func folderNamed(t *testing.T, rep InventoryReport, name string) Folder {
	t.Helper()
	for _, f := range rep.Folders {
		if f.Name == name {
			return f
		}
	}
	t.Fatalf("folder %q not listed: %+v", name, rep.Folders)
	return Folder{}
}

func hasInvDiag(rep InventoryReport, code, path string) bool {
	for _, d := range rep.Diagnostics {
		if d.Code == code && d.Path == path {
			return true
		}
	}
	return false
}

func TestInventorySymlinksAreNotFollowed(t *testing.T) {
	dir := invProject(t)
	outside := t.TempDir()
	writeFile(t, filepath.Join(outside, "tasks.md"), "- [x] T001 outside\n")
	writeFile(t, filepath.Join(dir, "specs", "001-real", "spec.md"), "# Real\n")
	if err := os.Symlink(outside, filepath.Join(dir, "specs", "002-linked")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "tasks.md"), filepath.Join(dir, "specs", "001-real", "tasks.md")); err != nil {
		t.Fatal(err)
	}
	var read []string
	inventoryAfterRead = func(p string) { read = append(read, p) }
	defer func() { inventoryAfterRead = nil }()
	rep := EvaluateInventory(InventoryOption{Root: dir})
	for _, p := range read {
		if strings.HasPrefix(p, outside) || strings.HasSuffix(p, "tasks.md") {
			t.Fatalf("symlinked content was read: %v", read)
		}
	}
	linked := folderNamed(t, rep, "002-linked")
	if !linked.Skipped || linked.Feature || linked.Tasks != nil || !hasInvDiag(rep, "SS_PATH_UNSAFE", "specs/002-linked") {
		t.Fatalf("symlinked folder: %+v", linked)
	}
	real := folderNamed(t, rep, "001-real")
	if real.Tasks == nil || real.Tasks.State != "skipped" || real.Tasks.Reason != "unsafe" || !hasInvDiag(rep, "SS_PATH_UNSAFE", "specs/001-real/tasks.md") {
		t.Fatalf("symlinked tasks.md: %+v", real.Tasks)
	}
	if rep.Summary.SkippedEntries != 1 || rep.Summary.Tasks.FilesSkipped != 1 || InventoryExitCode(rep) != 1 {
		t.Fatalf("summary %+v exit %d", rep.Summary, InventoryExitCode(rep))
	}
}

func TestInventoryFIFOTasksDoesNotBlock(t *testing.T) {
	dir := invProject(t)
	writeFile(t, filepath.Join(dir, "specs", "001-a", "spec.md"), "# A\n")
	if err := syscall.Mkfifo(filepath.Join(dir, "specs", "001-a", "tasks.md"), 0o644); err != nil {
		t.Skip("mkfifo unavailable")
	}
	done := make(chan InventoryReport, 1)
	go func() { done <- EvaluateInventory(InventoryOption{Root: dir}) }()
	select {
	case rep := <-done:
		f := folderNamed(t, rep, "001-a")
		if f.Tasks.State != "skipped" || f.Tasks.Reason != "unsafe" || !hasInvDiag(rep, "SS_PATH_UNSAFE", "specs/001-a/tasks.md") {
			t.Fatalf("FIFO tasks.md: %+v", f.Tasks)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("inventory blocked on a FIFO tasks.md")
	}
}

func TestInventorySensitiveNamesAreNotOpened(t *testing.T) {
	dir := invProject(t)
	writeFile(t, filepath.Join(dir, "specs", "001-a", "spec.md"), "# A\n")
	// Sensitive-named FIFOs would block if opened.
	if err := os.Mkdir(filepath.Join(dir, "specs", "002-client-secret"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(dir, "specs", "003-credentials"), 0o644); err != nil {
		t.Skip("mkfifo unavailable")
	}
	rep := EvaluateInventory(InventoryOption{Root: dir})
	for _, name := range []string{"002-client-secret", "003-credentials"} {
		f := folderNamed(t, rep, name)
		if !f.Skipped || f.Feature || !hasInvDiag(rep, "SS_REFERENCE_SENSITIVE_SKIPPED", "specs/"+name) {
			t.Fatalf("%s: %+v", name, f)
		}
	}
}

func TestInventoryUnusableSpecsDirectory(t *testing.T) {
	for _, kind := range []string{"symlink", "file"} {
		dir := invProject(t)
		switch kind {
		case "symlink":
			real := t.TempDir()
			writeFile(t, filepath.Join(real, "001-a", "spec.md"), "# A\n")
			if err := os.Symlink(real, filepath.Join(dir, "specs")); err != nil {
				t.Fatal(err)
			}
		case "file":
			writeFile(t, filepath.Join(dir, "specs"), "not a directory\n")
		}
		rep := EvaluateInventory(InventoryOption{Root: dir})
		if rep.Outcome != "report" || rep.SpecsState != "skipped" || len(rep.Folders) != 0 || !hasInvDiag(rep, "SS_PATH_UNSAFE", "specs") {
			t.Fatalf("%s specs: %+v", kind, rep)
		}
	}
}

func TestInventoryEntryLimitIsFatal(t *testing.T) {
	dir := invProject(t)
	for i := 0; i <= safeio.MaxDirEntries; i++ {
		if err := os.MkdirAll(filepath.Join(dir, "specs", strconv.Itoa(1000+i)+"-f"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	rep := EvaluateInventory(InventoryOption{Root: dir})
	if rep.Fatal == nil || rep.Fatal.Code != "SS_INPUT_LIMIT" || len(rep.Folders) != 0 || InventoryExitCode(rep) != 2 {
		t.Fatalf("limit: %+v", rep.Fatal)
	}
	if doc := string(InventoryJSON(rep)); !strings.Contains(doc, `"folders":[]`) || !strings.Contains(doc, `"summary":null`) {
		t.Fatalf("limit JSON: %s", doc)
	}
}

func TestInventoryTaskTotalLimitIsFatal(t *testing.T) {
	dir := invProject(t)
	line := "- [ ] T x\n"
	body := strings.Repeat(line, maxTasks/2+1)
	writeFile(t, filepath.Join(dir, "specs", "001-a", "tasks.md"), body)
	writeFile(t, filepath.Join(dir, "specs", "002-b", "tasks.md"), body)
	rep := EvaluateInventory(InventoryOption{Root: dir})
	if rep.Fatal == nil || rep.Fatal.Code != "SS_INPUT_LIMIT" {
		t.Fatalf("task total limit: %+v", rep.Fatal)
	}
}

func TestInventoryDeterminism(t *testing.T) {
	for _, name := range inventoryFixtures {
		root := twoModeRoot(name)
		a, b := EvaluateInventory(InventoryOption{Root: root}), EvaluateInventory(InventoryOption{Root: root})
		ja, jb := string(InventoryJSON(a)), string(InventoryJSON(b))
		if ja != jb {
			t.Fatalf("%s: JSON differs between runs", name)
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
