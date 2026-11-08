package conformance

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/tosin2013/spec-status/internal/status"
)

// repoRoot is the module root, two levels above this package.
var repoRoot = filepath.Join("..", "..")

// fixtureDirs returns every conformance fixture: each directory under
// testdata/conformance/, plus existing fixtures that carry an expected.json.
func fixtureDirs(t *testing.T) []string {
	t.Helper()
	var dirs []string
	entries, err := os.ReadDir(filepath.Join(repoRoot, "testdata", "conformance"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join(repoRoot, "testdata", "conformance", e.Name()))
		}
	}
	for _, extra := range []string{"speckit-roadmap-upstream"} {
		dir := filepath.Join(repoRoot, "testdata", extra)
		if _, err := os.Stat(filepath.Join(dir, "expected.json")); err == nil {
			dirs = append(dirs, dir)
		}
	}
	sort.Strings(dirs)
	return dirs
}

// TestConformance evaluates every fixture and compares it with its
// hand-counted expected.json (FR-004 to FR-007).
func TestConformance(t *testing.T) {
	dirs := fixtureDirs(t)
	if len(dirs) == 0 {
		t.Fatal("no conformance fixtures found")
	}
	for _, dir := range dirs {
		dir := dir
		t.Run(filepath.Base(dir), func(t *testing.T) {
			want, err := loadExpected(filepath.Join(dir, "expected.json"))
			if err != nil {
				t.Fatal(err)
			}
			abs, err := filepath.Abs(dir)
			if err != nil {
				t.Fatal(err)
			}
			rep := status.Evaluate(status.Option{Root: abs, CWD: abs})
			for _, d := range compare(want, rep) {
				t.Errorf("%s: %s", filepath.Base(dir), d)
			}
		})
	}
}
