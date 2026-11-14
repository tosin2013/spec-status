package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
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

// Build identity on the command surface (feature 004, contracts/build-identity.md).
func withBuild(t *testing.T, b buildIdentity) {
	t.Helper()
	buildOverride = &b
	t.Cleanup(func() { buildOverride = nil })
}

func TestVersionTextPerBuildKind(t *testing.T) {
	const sha = "0123456789abcdef0123456789abcdef01234567"
	tail := "go: " + runtime.Version() + "\ntarget: " + runtime.GOOS + "/" + runtime.GOARCH + "\n"
	for _, c := range []struct {
		b    buildIdentity
		want string
	}{
		{resolveBuild("", "0.0.0-dev", "unknown", nil, false),
			"spec 0.0.0-dev\ncommit: unknown\n" + tail + "label: development build, not a release\n"},
		{resolveBuild("release", "0.1.0", sha, nil, false),
			"spec 0.1.0\ncommit: " + sha + "\n" + tail + "label: release build\n"},
		{resolveBuild("", "0.0.0-dev", "unknown", &debug.BuildInfo{Main: debug.Module{Version: "v0.1.0", Sum: "h1:a="}}, true),
			"spec 0.1.0\ncommit: unknown\n" + tail + "label: module build (built from the module version, not a release archive)\n"},
	} {
		withBuild(t, c.b)
		if got := versionText(); got != c.want {
			t.Errorf("%s:\ngot  %q\nwant %q", c.b.kind, got, c.want)
		}
	}
}

func TestHelpSentencePerBuildKind(t *testing.T) {
	keep := "It does not change files, use the network, or infer quarters."
	for _, c := range []struct {
		b         buildIdentity
		want, not string
	}{
		{resolveBuild("", "0.0.0-dev", "unknown", nil, false), "This is a development build. " + keep, "release"},
		{resolveBuild("release", "0.1.0", "c", nil, false), "This is release 0.1.0. " + keep, "development build"},
		{resolveBuild("", "0.0.0-dev", "unknown", &debug.BuildInfo{Main: debug.Module{Version: "v0.1.0", Sum: "h1:a="}}, true), "This is a module build of 0.1.0, not a release archive. " + keep, "development build"},
	} {
		withBuild(t, c.b)
		var out, errb bytes.Buffer
		if code := Run([]string{"--help"}, t.TempDir(), &out, &errb); code != 0 {
			t.Fatal(code)
		}
		if !strings.Contains(out.String(), c.want) || strings.Contains(out.String(), c.not) {
			t.Errorf("%s help:\n%s", c.b.kind, out.String())
		}
	}
}

func TestStatusHeadingAndJSONPerBuildKind(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, ".specify", "memory", "roadmap.md"), "### 001 — Work [status: planned]\n")
	run := func(args ...string) string {
		var out, errb bytes.Buffer
		Run(args, dir, &out, &errb)
		return out.String()
	}
	withBuild(t, resolveBuild("", "0.0.0-dev", "unknown", nil, false))
	devText, devJSON := run("status", "--root", dir, "--color=never"), run("status", "--root", dir, "--json")
	withBuild(t, resolveBuild("release", "0.1.0", "c", nil, false))
	relText, relJSON := run("status", "--root", dir, "--color=never"), run("status", "--root", dir, "--json")
	if !strings.Contains(devText, "Spec Status  development build") || !strings.Contains(relText, "Spec Status  release 0.1.0") {
		t.Fatalf("headings:\n%s\n%s", devText, relText)
	}
	if devJSON != relJSON || devJSON == "" {
		t.Fatalf("JSON differs across build kinds:\n%s\n%s", devJSON, relJSON)
	}
}

// Feature 008 US5: --json=false, empty --root/--roadmap, fail_on in
// status invocation JSON (specs/008-field-forms-and-fixes/contracts/changes.md §5).
func TestInvocationEdgeCases(t *testing.T) {
	run := func(args ...string) (int, string, string) {
		var out, errb bytes.Buffer
		code := Run(args, ".", &out, &errb)
		return code, out.String(), errb.String()
	}
	code, out, errs := run("status", "--json=false", "--bogus")
	if code != 2 || out != "" || !strings.HasPrefix(errs, "spec: SS_INVOCATION_INVALID: ") {
		t.Errorf("--json=false: exit %d out %q err %q", code, out, errs)
	}
	code, out, _ = run("status", "--json=true", "--bogus")
	if code != 2 || !json.Valid([]byte(out)) {
		t.Errorf("--json=true: exit %d out %q", code, out)
	}
	for _, args := range [][]string{
		{"status", "--root", ""}, {"status", "--root="}, {"status", "--roadmap", ""}, {"status", "--roadmap="}, {"inventory", "--root", ""},
	} {
		name := "--root"
		if strings.Contains(strings.Join(args, " "), "roadmap") {
			name = "--roadmap"
		}
		code, _, errs := run(args...)
		if code != 2 || !strings.Contains(errs, "SS_INVOCATION_INVALID") || !strings.Contains(errs, name+" needs a non-empty value") {
			t.Errorf("%v: exit %d err %q", args, code, errs)
		}
	}
	_, out, _ = run("status", "--json", "--fail-on", "warning", "--bogus")
	if !strings.Contains(out, `"fail_on":"warning"`) || !strings.Contains(out, `"schema_version":"1.0"`) {
		t.Errorf("requested fail_on: %s", out)
	}
	_, out, _ = run("status", "--json", "--fail-on", "bogus")
	if !strings.Contains(out, `"fail_on":"error"`) {
		t.Errorf("invalid fail_on: %s", out)
	}
	// A --root value spelled --json does not turn JSON mode on.
	code, out, errs = run("status", "--root", "--json", "--bogus")
	if code != 2 || out != "" || !strings.Contains(errs, "SS_INVOCATION_INVALID") {
		t.Errorf("--root --json: exit %d out %q err %q", code, out, errs)
	}
}
