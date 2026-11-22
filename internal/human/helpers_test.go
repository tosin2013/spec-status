package human

import (
	"github.com/tosin2013/spec-status/internal/schema"
	"os"
	"path/filepath"
	"testing"
)

// project writes a temporary Spec Kit project (copied from the internal/status
// tests, where the evaluation tests keep their own copy).
func project(t *testing.T, roadmap string, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".specify", "memory", "roadmap.md"), roadmap)
	for name, body := range files {
		writeFile(t, filepath.Join(dir, filepath.FromSlash(name)), body)
	}
	return dir
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

const exampleRoadmap = "### 001 — Preserve original feature [status: implemented]\n" +
	"- Description: Historical original feature.\n" +
	"- Outcome: Original outcome remains recorded.\n" +
	"- Scope: Original scope is retained.\n" +
	"- Spec dir: specs/001-original\n" +
	"- superseded-by: 002\n" +
	"\n" +
	"### 002 — Replacement feature [status: verified]\n" +
	"- Description: Replacement scope can differ.\n" +
	"- Outcome: Recorded as verified, not independently checked by this CLI.\n" +
	"- Scope: Replacement scope.\n" +
	"- Spec dir: specs/002-replacement\n" +
	"\n" +
	"### 003 — Follow-on work [status: planned]\n" +
	"- Description: Still refers to the original requirement.\n" +
	"- Outcome: Requires a scope decision.\n" +
	"- Scope: Follow-on.\n" +
	"- Depends on: 001\n"

func boolPtr(v bool) *bool { return &v }

var inventoryFixtures = []string{"no-roadmap", "brownfield", "no-specs", "with-roadmap"}

func twoModeRoot(name string) string {
	return filepath.Join("..", "..", "testdata", "two-mode", name)
}

type recognitionExpect struct {
	Exit              int     `json:"exit"`
	Code              string  `json:"code"`
	Message           *string `json:"message"`
	TotalItems        *int    `json:"total_items"`
	ExplicitEmptyLine *bool   `json:"explicit_empty_line"`
}

const explicitEmptyLine = "roadmap explicitly empty (Planned Specs heading, no items)"

func percentOf(num, den int) schema.Percent { return schema.PercentOf(num, den) }

const taskScope = "unique_linked_task_files"
