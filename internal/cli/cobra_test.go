package cli

import (
	"bytes"
	"debug/elf"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestHelpVersionCompletionAreDiscoverable(t *testing.T) {
	var out, errb bytes.Buffer
	if code := Run(nil, t.TempDir(), &out, &errb); code != 0 {
		t.Fatalf("bare exit %d\n%s%s", code, out.String(), errb.String())
	}
	if errb.Len() != 0 {
		t.Fatalf("help wrote stderr %s", errb.String())
	}
	for _, phrase := range []string{"development build", "status", "completion", "version", "spec status --root testdata/example-roadmap"} {
		if !strings.Contains(out.String(), phrase) {
			t.Fatalf("help missing %q\n%s", phrase, out.String())
		}
	}
	if strings.Contains(out.String(), "--as-of") || strings.Contains(out.String(), "\x1b[") {
		t.Fatal("help advertises a date flag or uses ANSI")
	}
	again, againErr := bytes.Buffer{}, bytes.Buffer{}
	if code := Run(nil, t.TempDir(), &again, &againErr); code != 0 || again.String() != out.String() {
		t.Fatal("help is not deterministic")
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"status", "--help"}, t.TempDir(), &out, &errb); code != 0 || errb.Len() != 0 {
		t.Fatalf("status help %d %s", code, errb.String())
	}
	for _, phrase := range []string{"--json", "--root", "--roadmap", "--fail-on", "spec status --json --fail-on warning"} {
		if !strings.Contains(out.String(), phrase) {
			t.Fatalf("status help missing %q\n%s", phrase, out.String())
		}
	}
	if strings.Contains(out.String(), "\x1b[") {
		t.Fatal("status help uses ANSI")
	}

	out.Reset()
	if code := Run([]string{"version"}, t.TempDir(), &out, &errb); code != 0 || !strings.Contains(out.String(), "not a release") {
		t.Fatalf("version command %d %s", code, out.String())
	}
	flagOut := bytes.Buffer{}
	if code := Run([]string{"--version"}, t.TempDir(), &flagOut, &errb); code != 0 || flagOut.String() != out.String() {
		t.Fatal("version flag and command differ")
	}
}

func TestCompletionIsStdoutOnly(t *testing.T) {
	dir := t.TempDir()
	before := snapshot(t, dir)
	shells := []string{"bash", "zsh", "fish", "powershell"}
	var first string
	for _, shell := range shells {
		var out, errb bytes.Buffer
		code := Run([]string{"completion", shell}, dir, &out, &errb)
		if code != 0 || errb.Len() != 0 || out.Len() == 0 {
			t.Fatalf("%s completion exit %d stderr %q stdout %d", shell, code, errb.String(), out.Len())
		}
		if strings.Contains(out.String(), "\x1b[") {
			t.Fatalf("%s completion uses ANSI", shell)
		}
		if !strings.Contains(out.String(), "spec") {
			t.Fatalf("%s completion does not name the command", shell)
		}
		again := bytes.Buffer{}
		if code := Run([]string{"completion", shell}, dir, &again, &errb); code != 0 || again.String() != out.String() {
			t.Fatalf("%s completion is not deterministic", shell)
		}
		if shell == "bash" {
			first = out.String()
		}
	}
	if !strings.Contains(first, "complete") {
		t.Fatal("bash completion script has no complete helper")
	}
	if snapshot(t, dir) != before {
		t.Fatal("completion changed the working directory")
	}
	var out, errb bytes.Buffer
	if code := Run([]string{"completion", "nu"}, dir, &out, &errb); code != 2 || out.Len() != 0 {
		t.Fatalf("unknown shell exit %d stdout %s stderr %s", code, out.String(), errb.String())
	}
	if strings.Count(errb.String(), "SS_INVOCATION_INVALID") != 1 {
		t.Fatalf("duplicate or missing error: %s", errb.String())
	}
}

func TestUnknownArgsAndJSONFatalStayClean(t *testing.T) {
	cases := [][]string{
		{"status", "--as-of", "2026-01-01"},
		{"nope"},
		{"status", "--fail-on", "sometimes"},
		{"status", "--root"},
		{"status", "--root", "a", "--root", "b"},
		{"status", "--version"},
	}
	for _, args := range cases {
		var out, errb bytes.Buffer
		code := Run(args, t.TempDir(), &out, &errb)
		if code != 2 || out.Len() != 0 || strings.Count(errb.String(), "\n") != 1 {
			t.Fatalf("%v exit %d\nstdout %q\nstderr %q", args, code, out.String(), errb.String())
		}
		if strings.Contains(out.String()+errb.String(), "\x1b[") {
			t.Fatal("ANSI in error")
		}
	}
	var out, errb bytes.Buffer
	args := []string{"--json", "status", "--quarter"}
	if code := Run(args, t.TempDir(), &out, &errb); code != 2 || errb.Len() != 0 || !json.Valid(out.Bytes()) {
		t.Fatalf("json unknown exit %d stderr %q stdout %s", code, errb.String(), out.String())
	}
	again := bytes.Buffer{}
	if code := Run(args, t.TempDir(), &again, &errb); code != 2 || again.String() != out.String() {
		t.Fatal("json fatal is not deterministic")
	}
	if strings.Contains(out.String(), "\x1b[") || strings.Count(out.String(), "\n") != 1 {
		t.Fatalf("json envelope polluted: %s", out.String())
	}
}

func TestStaticBuildHasNoInterpreter(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("ELF static check runs on linux")
	}
	bin := filepath.Join(t.TempDir(), "spec")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/spec")
	cmd.Dir = filepath.Clean(filepath.Join("..", ".."))
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	f, err := elf.Open(bin)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_INTERP {
			t.Fatal("CGO_ENABLED=0 binary still has a dynamic interpreter")
		}
	}
}
