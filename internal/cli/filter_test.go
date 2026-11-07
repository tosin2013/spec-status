package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tosin2013/spec-status/internal/status"
)

func TestFilterFlagsAndExplain(t *testing.T) {
	dir := exampleDir(t)
	base := runArgs(t, []string{"status", "--root", dir, "--color=never"}, nil)
	if base.code != 0 {
		t.Fatalf("base %d %s", base.code, base.stderr)
	}
	review := runArgs(t, []string{"status", "--root", dir, "--color=never", "--human-review", "--explain"}, nil)
	if review.code != 0 || review.stderr != "" {
		t.Fatalf("review code %d stderr %s", review.code, review.stderr)
	}
	if !strings.Contains(review.stdout, "shown 1 of 3 items") || !strings.Contains(review.stdout, "recorded verified   1/2 (50%)") {
		t.Fatalf("review view changed totals\n%s", review.stdout)
	}
	if !strings.Contains(review.stdout, "explain: dependency 001 is human_review at .specify/memory/roadmap.md:18") {
		t.Fatalf("explain location\n%s", review.stdout)
	}
	if !strings.Contains(review.stdout, "SS_SPEC_REFERENCE_MISSING") || !strings.Contains(review.stdout, "SS_DEPENDENCY_HUMAN_REVIEW") {
		t.Fatalf("diagnostics were filtered\n%s", review.stdout)
	}
	again := runArgs(t, []string{"status", "--root", dir, "--color=never", "--human-review", "--explain"}, nil)
	if again.stdout != review.stdout || again.code != review.code {
		t.Fatal("explain output is not stable")
	}
	styled := runArgs(t, []string{"status", "--root", dir, "--color=always", "--explain"}, nil)
	if !strings.Contains(styled.stdout, "\x1b") || strings.Contains(stripCmdANSI(styled.stdout), "\x1b") {
		t.Fatal("explain color was not confined to styling")
	}
	plainPipe := runArgs(t, []string{"status", "--root", dir, "--explain"}, []string{"NO_COLOR=1", "TERM=dumb"})
	if strings.Contains(plainPipe.stdout, "\x1b") || !strings.Contains(plainPipe.stdout, "explain: excluded") {
		t.Fatalf("plain explain\n%s", plainPipe.stdout)
	}
	if !strings.Contains(plainPipe.stdout, "specs/001-original/tasks.md:1") {
		t.Fatalf("task evidence\n%s", plainPipe.stdout)
	}

	none := runArgs(t, []string{"status", "--root", dir, "--color=never", "--status", "not-a-status"}, nil)
	if none.code != 0 || !strings.Contains(none.stdout, "shown 0 of 3 items") || !strings.Contains(none.stdout, "no item matches this view") {
		t.Fatalf("no-match code %d\n%s", none.code, none.stdout)
	}
	hidden := runArgs(t, []string{"status", "--root", dir, "--color=never", "--status", "verified", "--fail-on", "warning"}, nil)
	if hidden.code != 1 || !strings.Contains(hidden.stdout, "shown 1 of 3 items") || !strings.Contains(hidden.stdout, "SS_DEPENDENCY_HUMAN_REVIEW") {
		t.Fatalf("hidden warning must still fail: code %d\n%s", hidden.code, hidden.stdout)
	}
	combined := runArgs(t, []string{"status", "--root", dir, "--color=never", "--status", "planned", "--blocked", "--active"}, nil)
	if combined.code != 0 || !strings.Contains(combined.stdout, "shown 1 of 3 items") || !strings.Contains(combined.stdout, "003") {
		t.Fatalf("combined\n%s", combined.stdout)
	}
	either := runArgs(t, []string{"status", "--root", dir, "--color=never", "--status", "implemented", "--status", "verified"}, nil)
	if either.code != 0 || !strings.Contains(either.stdout, "shown 2 of 3 items") || strings.Contains(either.stdout, "003  planned") {
		t.Fatalf("status OR\n%s", either.stdout)
	}
	historical := runArgs(t, []string{"status", "--root", dir, "--color=never", "--historical", "--status", "implemented"}, nil)
	if !strings.Contains(historical.stdout, "shown 1 of 3 items") || strings.Contains(historical.stdout, "002  verified") {
		t.Fatalf("historical status\n%s", historical.stdout)
	}
}

func TestFilterJSONStaysFullAndFatalStaysOneDocument(t *testing.T) {
	dir := exampleDir(t)
	plain := runArgs(t, []string{"status", "--json", "--root", dir}, nil)
	viewed := runArgs(t, []string{"status", "--json", "--root", dir, "--blocked", "--human-review", "--explain", "--status", "planned", "--active"}, nil)
	if plain.code != viewed.code || plain.stdout != viewed.stdout || viewed.stderr != "" || !json.Valid([]byte(viewed.stdout)) {
		t.Fatalf("JSON changed under a view: codes %d %d", plain.code, viewed.code)
	}
	if strings.Contains(viewed.stdout, "explain:") || strings.Contains(viewed.stdout, "shown ") {
		t.Fatal("schema 1.0 gained view fields")
	}
	fatal := runArgs(t, []string{"status", "--json", "--blocked", "--explain", "--root", filepath.Join(t.TempDir(), "missing")}, nil)
	if fatal.code != 2 || fatal.stderr != "" || !json.Valid([]byte(fatal.stdout)) || strings.Count(fatal.stdout, "\n") != 1 {
		t.Fatalf("fatal JSON code %d stderr %q docs %d", fatal.code, fatal.stderr, strings.Count(fatal.stdout, "\n"))
	}
	repeat := runArgs(t, []string{"status", "--status", "planned", "--status", "planned"}, nil)
	if repeat.code != 2 || repeat.stdout != "" || !strings.Contains(repeat.stderr, "SS_INVOCATION_INVALID") {
		t.Fatalf("repeat status %d %s %s", repeat.code, repeat.stdout, repeat.stderr)
	}
	repeatFlag := runArgs(t, []string{"status", "--blocked", "--blocked"}, nil)
	if repeatFlag.code != 2 || !strings.Contains(repeatFlag.stderr, "SS_INVOCATION_INVALID") {
		t.Fatalf("repeat blocked %d %s", repeatFlag.code, repeatFlag.stderr)
	}
	jsonRepeat := runArgs(t, []string{"--json", "status", "--explain", "--explain"}, nil)
	if jsonRepeat.code != 2 || jsonRepeat.stderr != "" || !json.Valid([]byte(jsonRepeat.stdout)) {
		t.Fatalf("json repeat %d %s %s", jsonRepeat.code, jsonRepeat.stdout, jsonRepeat.stderr)
	}
}

func TestFilterDoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	body := "### 001 — Need [status: planned]\n\n### 002 — Claim [status: verified]\n- Depends on: 001\n- Spec dir: specs/claim\n"
	mustWrite(t, filepath.Join(dir, ".specify", "memory", "roadmap.md"), body)
	mustWrite(t, filepath.Join(dir, "specs", "claim", "spec.md"), "# spec\n")
	mustWrite(t, filepath.Join(dir, "specs", "claim", "plan.md"), "# plan\n")
	mustWrite(t, filepath.Join(dir, "specs", "claim", "tasks.md"), "- [ ] Not yet\n")
	before := snapshot(t, dir)
	var out, errb bytes.Buffer
	code := Run([]string{"status", "--root", dir, "--color=never", "--status", "planned", "--explain"}, dir, &out, &errb)
	if code != 1 {
		t.Fatalf("hidden contradiction exit %d\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), "shown 1 of 2 items") || !strings.Contains(out.String(), "SS_VERIFIED_DEPENDENCY_CONTRADICTION") {
		t.Fatalf("threshold item was suppressed\n%s", out.String())
	}
	if !strings.Contains(out.String(), "explain: no task file linked") {
		t.Fatalf("missing task evidence\n%s", out.String())
	}
	after := snapshot(t, dir)
	if before != after {
		t.Fatal("filter run changed the tree")
	}
	rep := status.Evaluate(status.Option{Root: dir, CWD: dir})
	if rep.Summary == nil || rep.Summary.VerifiedActiveItems != 1 {
		t.Fatalf("full summary changed %+v", rep.Summary)
	}
}

func TestHelpMentionsFilters(t *testing.T) {
	var out, errb bytes.Buffer
	if code := Run([]string{"status", "--help"}, t.TempDir(), &out, &errb); code != 0 {
		t.Fatal(errb.String())
	}
	for _, phrase := range []string{"--status", "--blocked", "--human-review", "--active", "--historical", "--explain"} {
		if !strings.Contains(out.String(), phrase) {
			t.Fatalf("help missing %s\n%s", phrase, out.String())
		}
	}
	if _, err := os.Stat(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Fatal("help created a file")
	}
}
