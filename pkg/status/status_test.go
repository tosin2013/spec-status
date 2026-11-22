package status_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/tosin2013/spec-status/internal/cli"
	"github.com/tosin2013/spec-status/pkg/status"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

// fixtureRoots returns every directory under testdata that directly holds a
// non-symlink .specify directory, plus the repository itself and a
// directory with no project (a fatal report).
func fixtureRoots(t *testing.T) []string {
	t.Helper()
	repo := repoRoot(t)
	var roots []string
	err := filepath.Walk(filepath.Join(repo, "testdata"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".specify" {
			if fi, lerr := os.Lstat(path); lerr == nil && fi.Mode()&os.ModeSymlink == 0 {
				roots = append(roots, filepath.Dir(path))
			}
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) < 40 {
		t.Fatalf("found %d fixture roots, want at least 40; is the walk broken?", len(roots))
	}
	sort.Strings(roots)
	return append(roots, repo, t.TempDir())
}

func runCLI(t *testing.T, args ...string) ([]byte, int) {
	t.Helper()
	var out, errb bytes.Buffer
	code := cli.Run(args, repoRoot(t), &out, &errb)
	return out.Bytes(), code
}

func TestLoadExampleRoadmap(t *testing.T) {
	rep, err := status.Load(filepath.Join(repoRoot(t), "testdata", "example-roadmap"), status.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if status.SchemaVersion != "1.0" || rep.SchemaVersion != status.SchemaVersion {
		t.Fatalf("schema version %q / %q", rep.SchemaVersion, status.SchemaVersion)
	}
	if rep.Outcome != "report" || !rep.Complete || rep.Summary == nil || rep.Error != nil {
		t.Fatalf("outcome %q complete %v", rep.Outcome, rep.Complete)
	}
	var got []string
	for _, it := range rep.Items {
		got = append(got, it.ID+" "+it.Status+" "+it.Tasks.State)
		if (it.Tasks.State == "read") != (it.Tasks.Phases != nil) {
			t.Errorf("%s: phases nil=%v with state %s", it.ID, it.Tasks.Phases == nil, it.Tasks.State)
		}
	}
	if want := "001 implemented read|002 verified read|003 planned not_linked"; strings.Join(got, "|") != want {
		t.Fatalf("items %q, want %q", strings.Join(got, "|"), want)
	}
	s := rep.Summary
	if s.TotalItems != 3 || s.ActiveItems != 2 || s.VerifiedActiveItems != 1 || s.BlockedCount != 1 || s.HumanReviewCount != 1 {
		t.Fatalf("summary %+v", *s)
	}
	if !s.VerifiedPercent.Valid() || s.VerifiedPercent.Hundredths() != 5000 || s.VerifiedPercent.Text() != "50" {
		t.Fatalf("verified percent %s", s.VerifiedPercent.Text())
	}
	if status.ExitCode(rep) != 0 {
		t.Fatalf("exit %d", status.ExitCode(rep))
	}
}

func TestLoadFatalAndInvalid(t *testing.T) {
	repo := repoRoot(t)
	empty := t.TempDir()
	if err := os.Mkdir(filepath.Join(empty, ".specify"), 0o755); err != nil {
		t.Fatal(err)
	}
	rep, err := status.Load(empty, status.Options{FailOn: "none"})
	if err != nil {
		t.Fatalf("fatal roadmap returned a Go error: %v", err)
	}
	if rep.Outcome != "fatal" || rep.Error == nil || rep.Error.Code != "SS_ROADMAP_UNAVAILABLE" || rep.Summary != nil {
		t.Fatalf("got outcome %q error %+v", rep.Outcome, rep.Error)
	}
	if status.ExitCode(rep) != 2 {
		t.Fatal("fail-on none hid a fatal report")
	}

	if _, err := status.Load("", status.Options{}); !errors.Is(err, status.ErrInvalidInput) || !strings.Contains(err.Error(), "root") {
		t.Fatalf("empty root: %v", err)
	}
	if _, err := status.Load(repo, status.Options{FailOn: "loud"}); !errors.Is(err, status.ErrInvalidInput) || !strings.Contains(err.Error(), "FailOn") {
		t.Fatalf("bad fail-on: %v", err)
	}

	rep, err = status.Load(filepath.Join(repo, "testdata", "example-roadmap"), status.Options{Roadmap: filepath.Join(repo, "README.md")})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Outcome != "fatal" || rep.Error == nil || rep.Error.Code != "SS_PATH_UNSAFE" {
		t.Fatalf("roadmap outside root: outcome %q error %+v", rep.Outcome, rep.Error)
	}
}

// TestMatchesCLI is the no-drift guarantee: for every fixture root and every
// threshold, the library's bytes and exit code equal the CLI's.
func TestMatchesCLI(t *testing.T) {
	n := 0
	for _, root := range fixtureRoots(t) {
		for _, failOn := range []string{"error", "warning", "info", "none"} {
			want, code := runCLI(t, "status", "--json", "--root", root, "--fail-on", failOn)
			rep, err := status.Load(root, status.Options{FailOn: failOn})
			if err != nil {
				t.Fatalf("%s %s: %v", root, failOn, err)
			}
			if got := status.Marshal(rep); !bytes.Equal(got, want) {
				t.Errorf("%s --fail-on %s: bytes differ\nlibrary: %s\ncli:     %s", root, failOn, got, want)
			}
			if got := status.ExitCode(rep); got != code {
				t.Errorf("%s --fail-on %s: exit %d, cli %d", root, failOn, got, code)
			}
			n++
		}
	}
	t.Logf("compared %d root/threshold pairs", n)
}

// TestRoundTrip decodes the CLI's JSON into Report and re-encodes it.
func TestRoundTrip(t *testing.T) {
	for _, root := range fixtureRoots(t) {
		want, _ := runCLI(t, "status", "--json", "--root", root)
		var rep status.Report
		if err := json.Unmarshal(want, &rep); err != nil {
			t.Fatalf("%s: %v", root, err)
		}
		if got := status.Marshal(rep); !bytes.Equal(got, want) {
			t.Errorf("%s: round trip differs\ngot:  %s\nwant: %s", root, got, want)
		}
	}
}

func TestPercentJSON(t *testing.T) {
	for _, in := range []string{"null", "0", "50", "100", "41.67", "0.05", "12.5"} {
		var p status.Percent
		if err := json.Unmarshal([]byte(in), &p); err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		out, _ := json.Marshal(p)
		if string(out) != in {
			t.Errorf("%s round-trips to %s", in, out)
		}
	}
	for _, in := range []string{`"41"`, "1.234", "-1", "41.10", "41.", "01", "1e2", "true"} {
		var p status.Percent
		if err := json.Unmarshal([]byte(in), &p); err == nil {
			t.Errorf("%s accepted as %s", in, p.Text())
		}
	}
}

func TestConcurrentLoads(t *testing.T) {
	repo := repoRoot(t)
	roots := []string{
		filepath.Join(repo, "testdata", "example-roadmap"),
		filepath.Join(repo, "testdata", "phases"),
		t.TempDir(), // fatal: no project
	}
	want := map[string][]byte{}
	for _, r := range roots {
		rep, err := status.Load(r, status.Options{})
		if err != nil {
			t.Fatal(err)
		}
		want[r] = status.Marshal(rep)
	}
	var wg sync.WaitGroup
	errs := make(chan string, 64)
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 5; i++ {
				for _, r := range roots {
					rep, err := status.Load(r, status.Options{})
					if err != nil {
						errs <- err.Error()
						return
					}
					if !bytes.Equal(status.Marshal(rep), want[r]) {
						errs <- "concurrent result differs for " + r
						return
					}
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Error(e)
	}
}
