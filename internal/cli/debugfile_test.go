package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestCompletionDebugFileIsNotTouched(t *testing.T) {
	bin := builtSpec(t)
	repo := moduleRoot(t)
	example := filepath.Join(repo, "testdata", "example-roadmap")
	invocations := [][]string{
		{"--help"},
		{"status", "--help"},
		{"help", "status"},
		{"completion", "bash"},
		{"completion", "zsh"},
		{"completion", "fish"},
		{"completion", "powershell"},
		{"status", "--root", example},
		{"status", "--root", filepath.Join(t.TempDir(), "missing-root")},
		{"status", "--json", "--root", filepath.Join(t.TempDir(), "missing-root")},
		{"status", "--not-a-flag"},
		{"nope"},
		{"__complete", ""},
		{"__complete", "status", ""},
		{"__complete", "nosuch", ""},
		{"__completeNoDesc", "status", "--not-a-real-flag", ""},
	}

	for _, mode := range []string{"missing", "sentinel"} {
		debugDir := filepath.Join(t.TempDir(), mode)
		if err := os.MkdirAll(debugDir, 0o755); err != nil {
			t.Fatal(err)
		}
		debugPath := filepath.Join(debugDir, "cobra-debug.log")
		const sentinel = "SENTINEL-DO-NOT-CHANGE\n"
		if mode == "sentinel" {
			if err := os.WriteFile(debugPath, []byte(sentinel), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		for _, args := range invocations {
			before := snapshot(t, debugDir)
			stdout, stderr, code := runSpec(t, bin, args, debugPath)
			after := snapshot(t, debugDir)
			if before != after {
				t.Fatalf("%s %v changed %s\nbefore %s\nafter %s", mode, args, debugDir, before, after)
			}
			if mode == "missing" {
				if _, err := os.Stat(debugPath); !os.IsNotExist(err) {
					t.Fatalf("%v created %s", args, debugPath)
				}
			} else {
				got, err := os.ReadFile(debugPath)
				if err != nil || string(got) != sentinel {
					t.Fatalf("%v modified sentinel: %q %v", args, got, err)
				}
			}
			if bytes.Contains(stdout, []byte{0x1b}) || bytes.Contains(stderr, []byte{0x1b}) {
				t.Fatalf("%v emitted ANSI", args)
			}
			stdout2, stderr2, code2 := runSpec(t, bin, args, debugPath)
			if code2 != code || !bytes.Equal(stdout2, stdout) || !bytes.Equal(stderr2, stderr) {
				t.Fatalf("%v output not stable\ncode %d/%d\nstdout %q / %q\nstderr %q / %q", args, code, code2, stdout, stdout2, stderr, stderr2)
			}
		}
	}
}

func TestCompletionDebugEnvRestoredAcrossConcurrentRun(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sentinel.log")
	const sentinel = "SENTINEL-DO-NOT-CHANGE\n"
	if err := os.WriteFile(path, []byte(sentinel), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv(completionDebugEnv, path)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var out, errb bytes.Buffer
			_ = Run([]string{"__complete", "nosuch", ""}, t.TempDir(), &out, &errb)
		}()
	}
	wg.Wait()
	if got := os.Getenv(completionDebugEnv); got != path {
		t.Fatalf("environment restored to %q", got)
	}
	body, err := os.ReadFile(path)
	if err != nil || string(body) != sentinel {
		t.Fatalf("concurrent Run modified sentinel: %q %v", body, err)
	}
}

func TestGeneratedCompletionScriptDoesNotRunShellWrites(t *testing.T) {
	var out, errb bytes.Buffer
	if code := Run([]string{"completion", "bash"}, t.TempDir(), &out, &errb); code != 0 {
		t.Fatal(code, errb.String())
	}
	// The script is text. It is not executed here. Upstream shell helpers mention the debug variable.
	if !strings.Contains(out.String(), "BASH_COMP_DEBUG_FILE") {
		t.Fatal("bash completion script no longer mentions BASH_COMP_DEBUG_FILE; update the shell-versus-process note")
	}
}

func builtSpec(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "spec")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/spec")
	cmd.Dir = moduleRoot(t)
	cmd.Env = envWith("CGO_ENABLED", "0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	return bin
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(dir, "..", ".."))
}

func runSpec(t *testing.T, bin string, args []string, debugPath string) ([]byte, []byte, int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Env = envWith(completionDebugEnv, debugPath)
	cmd.Dir = t.TempDir()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatal(err)
		}
		code = exit.ExitCode()
	}
	return stdout.Bytes(), stderr.Bytes(), code
}

func envWith(key, val string) []string {
	prefix := key + "="
	out := make([]string, 0)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, prefix) {
			continue
		}
		out = append(out, entry)
	}
	return append(out, prefix+val)
}
