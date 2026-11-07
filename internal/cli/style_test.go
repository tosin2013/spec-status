package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestColorModesAndPlainJSON(t *testing.T) {
	dir := exampleDir(t)
	plain := runArgs(t, []string{"status", "--root", dir, "--color=never"}, nil)
	if strings.Contains(plain.stdout, "\x1b") || plain.code != 0 {
		t.Fatalf("plain code %d\n%s", plain.code, plain.stdout)
	}
	styled := runArgs(t, []string{"status", "--root", dir, "--color=always"}, nil)
	if styled.code != 0 || !strings.Contains(styled.stdout, "\x1b") {
		t.Fatalf("styled code %d ansi %v", styled.code, strings.Contains(styled.stdout, "\x1b"))
	}
	if stripCmdANSI(plain.stdout) == "" || !strings.Contains(stripCmdANSI(styled.stdout), "recorded verified") {
		t.Fatal("missing recorded verified label")
	}
	again := runArgs(t, []string{"status", "--root", dir, "--color=always"}, nil)
	if again.stdout != styled.stdout || again.stderr != styled.stderr || again.code != styled.code {
		t.Fatal("styled output is not deterministic")
	}
	auto := runArgs(t, []string{"status", "--root", dir}, nil)
	if strings.Contains(auto.stdout, "\x1b") {
		t.Fatal("piped auto mode emitted ANSI")
	}
	noColor := runArgs(t, []string{"status", "--root", dir, "--color=auto"}, []string{"NO_COLOR=1", "TERM=xterm-256color"})
	if strings.Contains(noColor.stdout, "\x1b") {
		t.Fatal("NO_COLOR did not disable auto color")
	}
	dumb := runArgs(t, []string{"status", "--root", dir}, []string{"TERM=dumb"})
	if strings.Contains(dumb.stdout, "\x1b") {
		t.Fatal("TERM=dumb emitted ANSI")
	}
	forced := runArgs(t, []string{"status", "--root", dir, "--color=always"}, []string{"NO_COLOR=1", "TERM=dumb"})
	if !strings.Contains(forced.stdout, "\x1b") {
		t.Fatal("explicit --color=always stayed plain")
	}
	js := runArgs(t, []string{"status", "--root", dir, "--json", "--color=always"}, nil)
	if js.code != 0 || strings.Contains(js.stdout, "\x1b") || js.stderr != "" {
		t.Fatalf("json code %d stderr %q", js.code, js.stderr)
	}
	js2 := runArgs(t, []string{"status", "--root", dir, "--json"}, nil)
	if js.stdout != js2.stdout {
		t.Fatal("color flag changed JSON bytes")
	}
	ver := runArgs(t, []string{"version", "--color=always"}, nil)
	comp := runArgs(t, []string{"completion", "bash"}, nil)
	if strings.Contains(ver.stdout, "\x1b[") || strings.Contains(comp.stdout, "\x1b[") {
		t.Fatal("version or completion was styled")
	}
}

func TestNarrowColumnsAndUntrustedTitle(t *testing.T) {
	dir := t.TempDir()
	title := "你好世界\x1b[31mPWN\x1b[0m\x1b]8;;http://evil.example\x1b\\END"
	body := "### 001 — " + title + " [status: planned]\n"
	mustWrite(t, filepath.Join(dir, ".specify", "memory", "roadmap.md"), body)
	narrow := runArgs(t, []string{"status", "--root", dir, "--color=always"}, []string{"COLUMNS=24"})
	if narrow.code != 0 {
		t.Fatal(narrow.stderr)
	}
	if strings.Contains(narrow.stdout, "\x1b[31m") || strings.Contains(narrow.stdout, "http://evil.example") || strings.Contains(narrow.stdout, "\x1b]8") {
		t.Fatalf("untrusted controls rendered: %q", narrow.stdout)
	}
	plain := runArgs(t, []string{"status", "--root", dir, "--color=never"}, []string{"COLUMNS=24"})
	if stripCmdANSI(narrow.stdout) != plain.stdout {
		t.Fatalf("narrow color and plain facts differ\ncolor %s\nplain %s", stripCmdANSI(narrow.stdout), plain.stdout)
	}
	js := runArgs(t, []string{"status", "--root", dir, "--json"}, nil)
	if !strings.Contains(js.stdout, `\u001b`) {
		t.Fatalf("JSON did not keep the raw title escape: %s", js.stdout)
	}
	if js.code != 0 {
		t.Fatal(js.code)
	}
}

func TestColorDoesNotWriteOrChangeExits(t *testing.T) {
	dir := exampleDir(t)
	before := snapshot(t, dir)
	ok := runArgs(t, []string{"status", "--root", dir, "--color=always"}, nil)
	bad := runArgs(t, []string{"status", "--root", filepath.Join(dir, "missing"), "--color=always", "--fail-on", "none"}, nil)
	if ok.code != 0 || bad.code != 2 {
		t.Fatalf("exits ok %d bad %d", ok.code, bad.code)
	}
	if strings.Contains(bad.stderr, "\x1b") {
		t.Fatal("fatal text was colored")
	}
	if snapshot(t, dir) != before {
		t.Fatal("styled status wrote files")
	}
}

type cmdResult struct {
	stdout, stderr string
	code           int
}

func runArgs(t *testing.T, args []string, extra []string) cmdResult {
	t.Helper()
	var out, errb bytes.Buffer
	// Run in-process so --color=always does not need a TTY.
	// These environment changes are process-global. Tests in this package are not parallel.
	old := map[string]string{}
	for _, e := range extra {
		parts := strings.SplitN(e, "=", 2)
		old[parts[0]] = os.Getenv(parts[0])
		if err := os.Setenv(parts[0], parts[1]); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		for k, v := range old {
			if v == "" {
				_ = os.Unsetenv(k)
			} else {
				_ = os.Setenv(k, v)
			}
		}
	})
	code := Run(args, t.TempDir(), &out, &errb)
	return cmdResult{stdout: out.String(), stderr: errb.String(), code: code}
}

func exampleDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(moduleRoot(t), "testdata", "example-roadmap")
}

func stripCmdANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == 0x1b {
			i++
			if i < len(s) && s[i] == '[' {
				i++
				for i < len(s) && (s[i] < 0x40 || s[i] > 0x7e) {
					i++
				}
				if i < len(s) {
					i++
				}
				continue
			}
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
