package cli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// Feature 006: spec inventory command (contracts/inventory-cli.md).

func inventoryFixture(name string) string {
	return filepath.Join("..", "..", "testdata", "two-mode", name)
}

func runCLI(args ...string) (int, string, string) {
	var out, errb bytes.Buffer
	code := Run(args, ".", &out, &errb)
	return code, out.String(), errb.String()
}

func TestInventoryCommandExitCodes(t *testing.T) {
	code, out, errs := runCLI("inventory", "--root", inventoryFixture("no-roadmap"), "--color=never")
	if code != 0 || errs != "" || !strings.Contains(out, "Spec Inventory") || !strings.Contains(out, "001-alpha") {
		t.Fatalf("human: exit %d stderr %q\n%s", code, errs, out)
	}
	code, out, _ = runCLI("inventory", "--root", inventoryFixture("no-roadmap"), "--json")
	var doc map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &doc) != nil || doc["inventory_schema_version"] != "1.0" {
		t.Fatalf("json: exit %d\n%s", code, out)
	}
	if code, _, _ = runCLI("inventory", "--root", inventoryFixture("no-roadmap"), "--json", "--fail-on", "info"); code != 1 {
		t.Fatalf("--fail-on info exit %d, want 1", code)
	}
	code, out, errs = runCLI("inventory", "--root", t.TempDir())
	if code != 2 || out != "" || !strings.HasPrefix(errs, "spec: SS_ROOT_UNAVAILABLE: ") {
		t.Fatalf("fatal: exit %d out %q err %q", code, out, errs)
	}
}

func TestInventoryRejectsRoadmapOnlyFlags(t *testing.T) {
	root := inventoryFixture("no-roadmap")
	for _, args := range [][]string{
		{"--roadmap", "x.md"}, {"--status", "planned"}, {"--blocked"}, {"--human-review"},
		{"--active"}, {"--historical"}, {"--explain"},
	} {
		name := strings.TrimPrefix(args[0], "--")
		full := append([]string{"inventory", "--root", root}, args...)
		code, _, errs := runCLI(full...)
		if code != 2 || !strings.Contains(errs, "SS_INVOCATION_INVALID") || !strings.Contains(errs, "--"+name+" is not used by inventory") {
			t.Errorf("%v: exit %d err %q", args, code, errs)
		}
		code, out, _ := runCLI(append(full, "--json", "--fail-on", "warning")...)
		var doc map[string]any
		if code != 2 || json.Unmarshal([]byte(out), &doc) != nil {
			t.Fatalf("%v --json: exit %d out %q", args, code, out)
		}
		cfg, _ := doc["configuration"].(map[string]any)
		if doc["inventory_schema_version"] != "1.0" || doc["outcome"] != "fatal" || cfg["fail_on"] != "warning" {
			t.Errorf("%v --json: not an inventory-shaped fatal with the requested fail_on: %s", args, out)
		}
	}
}

func TestInventoryInvocationErrors(t *testing.T) {
	root := inventoryFixture("no-roadmap")
	for _, args := range [][]string{
		{"inventory", "--root", root, "--version"},
		{"inventory", "--root", root, "--root", root},
		{"inventory", "--root", root, "--fail-on", "bogus"},
		{"inventory", "extra-arg", "--root", root},
	} {
		code, _, errs := runCLI(args...)
		if code != 2 || !strings.Contains(errs, "SS_INVOCATION_INVALID") {
			t.Errorf("%v: exit %d err %q", args, code, errs)
		}
		code, out, _ := runCLI(append(args, "--json")...)
		var doc map[string]any
		if code != 2 || json.Unmarshal([]byte(out), &doc) != nil || doc["inventory_schema_version"] != "1.0" {
			t.Errorf("%v --json: exit %d out %q", args, code, out)
		}
	}
}

func TestStatusNeverRunsInventory(t *testing.T) {
	code, out, errs := runCLI("status", "--root", inventoryFixture("no-roadmap"))
	if code != 2 || out != "" || !strings.Contains(errs, "SS_ROADMAP_UNAVAILABLE") || !strings.Contains(errs, "spec inventory") || strings.Contains(errs, "001-alpha") {
		t.Fatalf("status without roadmap: exit %d out %q err %q", code, out, errs)
	}
	code, out, _ = runCLI("status", "--root", inventoryFixture("no-roadmap"), "--json")
	if code != 2 || strings.Contains(out, "inventory_schema_version") || !strings.Contains(out, `"schema_version":"1.0"`) {
		t.Fatalf("status --json without roadmap: %s", out)
	}
}

func TestInventoryHelpAndCompletion(t *testing.T) {
	_, out, _ := runCLI("--help")
	if !strings.Contains(out, "inventory") {
		t.Errorf("root help lacks inventory:\n%s", out)
	}
	_, out, _ = runCLI("inventory", "--help")
	if !strings.Contains(out, "does not read the roadmap") || !strings.Contains(out, "not verification") {
		t.Errorf("inventory help:\n%s", out)
	}
	// Completion scripts are dynamic: the shell asks the binary for names.
	_, out, _ = runCLI("__complete", "")
	if !strings.Contains(out, "inventory\tList Spec Kit feature folders without a roadmap") {
		t.Errorf("completion lacks inventory:\n%s", out)
	}
}

func TestInventoryCommandDetection(t *testing.T) {
	root := inventoryFixture("brownfield")
	for _, args := range [][]string{
		{"--json", "inventory", "--root", root},
		{"inventory", "--json", "--root", root},
		{"--root", root, "--fail-on", "warning", "--json", "inventory"},
	} {
		code, out, _ := runCLI(args...)
		if code != 0 || !strings.Contains(out, `"inventory_schema_version":"1.0"`) {
			t.Errorf("%v: exit %d out %q", args, code, out)
		}
	}
	// A --root value named "inventory" is not the command.
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "inventory", ".specify", "memory", "roadmap.md"), "### 001 — Work [status: planned]\n")
	var out, errb bytes.Buffer
	code := Run([]string{"status", "--root", "inventory", "--json", "--bogus"}, dir, &out, &errb)
	if code != 2 || strings.Contains(out.String(), "inventory_schema_version") || !strings.Contains(out.String(), `"schema_version":"1.0"`) {
		t.Fatalf("status with --root inventory: exit %d out %s", code, out.String())
	}
	if got := requestedFailOn([]string{"inventory", "--fail-on", "--json"}); got != "" {
		t.Errorf("missing --fail-on value read as %q", got)
	}
	if got := requestedFailOn([]string{"inventory", "--root", "--fail-on", "--fail-on=none"}); got != "none" {
		t.Errorf("--root value \"--fail-on\" taken as the flag; got %q", got)
	}
}
