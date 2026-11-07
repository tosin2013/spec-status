package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tosin2013/spec-status/internal/status"
)

func TestInvocationAndNoWrite(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, ".specify", "memory", "roadmap.md"), "### 001 — Work [status: planned]\n")
	before := snapshot(t, dir)
	var out, errb bytes.Buffer
	code := Run([]string{"status", "--root", dir, "--json"}, dir, &out, &errb)
	if code != 0 {
		t.Fatalf("exit %d stderr %s", code, errb.String())
	}
	if errb.Len() != 0 {
		t.Fatalf("json stderr %s", errb.String())
	}
	if !json.Valid(out.Bytes()) {
		t.Fatalf("invalid json %s", out.String())
	}
	if !bytes.HasSuffix(out.Bytes(), []byte("\n")) {
		t.Fatal("missing trailing newline")
	}
	dec := json.NewDecoder(bytes.NewReader(out.Bytes()))
	var doc map[string]any
	if err := dec.Decode(&doc); err != nil {
		t.Fatal(err)
	}
	if dec.More() {
		t.Fatal("more than one JSON document")
	}
	after := snapshot(t, dir)
	if before != after {
		t.Fatalf("status wrote or removed files\nbefore %s\nafter %s", before, after)
	}
	out.Reset()
	errb.Reset()
	code = Run([]string{"status", "--root", dir}, dir, &out, &errb)
	if code != 0 || !strings.Contains(out.String(), "recorded verified") {
		t.Fatalf("text exit %d\n%s\n%s", code, out.String(), errb.String())
	}
}

func TestFatalJSONSurvivesNone(t *testing.T) {
	var out, errb bytes.Buffer
	code := Run([]string{"status", "--json", "--fail-on", "none", "--root", filepath.Join(t.TempDir(), "missing")}, t.TempDir(), &out, &errb)
	if code != 2 || errb.Len() != 0 {
		t.Fatalf("exit %d stderr %q stdout %s", code, errb.String(), out.String())
	}
	var doc struct {
		Outcome  string         `json:"outcome"`
		Complete bool           `json:"complete"`
		Summary  any            `json:"summary"`
		Error    map[string]any `json:"error"`
	}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Outcome != "fatal" || doc.Complete || doc.Summary != nil || doc.Error["code"] == "" {
		t.Fatalf("fatal envelope %+v", doc)
	}
}

func TestUnknownOptionAndNoQuarterFlag(t *testing.T) {
	var out, errb bytes.Buffer
	if code := Run([]string{"status", "--as-of", "2026-01-01"}, t.TempDir(), &out, &errb); code != 2 {
		t.Fatalf("as-of exit %d", code)
	}
	if strings.Contains(out.String(), "quarter") {
		t.Fatal(out.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"--json", "status", "--quarter"}, t.TempDir(), &out, &errb); code != 2 {
		t.Fatalf("quarter exit %d", code)
	}
	if !json.Valid(out.Bytes()) || errb.Len() != 0 {
		t.Fatalf("stdout %s stderr %s", out.String(), errb.String())
	}
}

func TestHelpDoesNotNeedAProject(t *testing.T) {
	var out, errb bytes.Buffer
	if code := Run([]string{"--help"}, t.TempDir(), &out, &errb); code != 0 {
		t.Fatal(code)
	}
	if !strings.Contains(out.String(), "development build") || strings.Contains(out.String(), "--as-of") {
		t.Fatal(out.String())
	}
	out.Reset()
	if code := Run([]string{"--version"}, t.TempDir(), &out, &errb); code != 0 {
		t.Fatal(code)
	}
	if !strings.Contains(out.String(), "not a release") {
		t.Fatal(out.String())
	}
}

func TestContradictionStaysVerified(t *testing.T) {
	dir := t.TempDir()
	body := "### 001 — Need [status: planned]\n\n### 002 — Claim [status: verified]\n- Depends on: 001\n- Spec dir: specs/claim\n"
	mustWrite(t, filepath.Join(dir, ".specify", "memory", "roadmap.md"), body)
	mustWrite(t, filepath.Join(dir, "specs", "claim", "spec.md"), "# spec\n")
	mustWrite(t, filepath.Join(dir, "specs", "claim", "plan.md"), "# plan\n")
	mustWrite(t, filepath.Join(dir, "specs", "claim", "tasks.md"), "- [ ] Not yet\n")
	rep := status.Evaluate(status.Option{Root: dir, CWD: dir})
	if rep.Summary == nil || rep.Summary.VerifiedActiveItems != 1 {
		t.Fatalf("summary %+v", rep.Summary)
	}
	found := false
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_VERIFIED_DEPENDENCY_CONTRADICTION" {
			found = true
		}
	}
	if !found || status.ExitCode(rep) != 1 {
		t.Fatalf("contradiction exit %d diags %v", status.ExitCode(rep), rep.Diagnostics)
	}
}

func TestSymlinkRoadmapDoesNotLeak(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	secretPath := filepath.Join(outside, "secret.md")
	if err := os.WriteFile(secretPath, []byte("TOPSECRET-VALUE\n### 001 — X [status: verified]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".specify", "memory"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(secretPath, filepath.Join(root, ".specify", "memory", "roadmap.md")); err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	code := Run([]string{"status", "--root", root}, root, &out, &errb)
	combined := out.String() + errb.String()
	if strings.Contains(combined, "TOPSECRET-VALUE") {
		t.Fatal("symlink target contents were reported")
	}
	if code != 2 {
		t.Fatalf("expected fatal symlink read, exit %d output %s", code, combined)
	}
}

func TestAdversarialTaskFIFO(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, ".specify", "memory", "roadmap.md"), "### 001 — A [status: planned]\n- Spec dir: specs/a\n")
	mustWrite(t, filepath.Join(dir, "specs", "a", "tasks.md"), "- [x] Wait [pipe](pipe)\n")
	if err := syscallMkfifo(filepath.Join(dir, "specs", "a", "pipe")); err != nil {
		t.Fatal(err)
	}
	done := make(chan int, 1)
	go func() {
		var out, errb bytes.Buffer
		done <- Run([]string{"status", "--root", dir}, dir, &out, &errb)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("status blocked while checking a FIFO task target")
	}
}

func snapshot(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		b.WriteString(rel)
		b.WriteByte(' ')
		b.WriteString(info.Mode().String())
		b.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
