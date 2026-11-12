package extension

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/tosin2013/spec-status/internal/status"
)

// FR-017, contracts/fixtures.md: every skill fixture under testdata/skills
// must still produce the signals its EXPECTED.md lists, so a CLI change can
// never silently invalidate a hand-written expectation.
//
// Signals line grammar (clauses separated by ";"):
//   item <ID> [status=<s>] [tasks=<checked>/<total>] [state=<task state>]
//   diag <CODE> [item=<ID>]
//   dep <ID>-><DEPID> state=<dependency state>
//   summary [active=<n>] [verified=<n>] [blocked=<n>] [human_review=<n>]
//   outcome=<report|fatal>

var (
	signalsLineRE = regexp.MustCompile(`(?m)^- Signals[^:]*:\s*(.+)$`)
	depClauseRE   = regexp.MustCompile(`^dep (\d{3})->(\d{3}) state=(\S+)$`)
)

func TestSkillFixturesProduceTheirSignals(t *testing.T) {
	cases, _ := filepath.Glob(filepath.Join(repoRoot, "testdata", "skills", "*", "EXPECTED.md"))
	if len(cases) == 0 {
		t.Skip("no skill fixtures yet")
	}
	for _, exp := range cases {
		dir := filepath.Dir(exp)
		t.Run(filepath.Base(dir), func(t *testing.T) {
			m := signalsLineRE.FindStringSubmatch(readFile(t, exp))
			if m == nil {
				t.Fatal("EXPECTED.md has no '- Signals:' line")
			}
			abs, err := filepath.Abs(dir)
			if err != nil {
				t.Fatal(err)
			}
			rep := status.Evaluate(status.Option{Root: abs, CWD: abs, FailOn: "none"})
			wantFatal := false
			for _, clause := range strings.Split(m[1], ";") {
				clause = strings.TrimSpace(strings.Trim(strings.TrimSpace(clause), "`"))
				if clause == "" {
					continue
				}
				switch {
				case strings.HasPrefix(clause, "outcome="):
					want := strings.TrimPrefix(clause, "outcome=")
					wantFatal = want == "fatal"
					if rep.Outcome != want {
						t.Errorf("outcome %s, want %s", rep.Outcome, want)
					}
				case strings.HasPrefix(clause, "item "):
					checkItem(t, rep, strings.Fields(clause)[1:])
				case strings.HasPrefix(clause, "diag "):
					checkDiag(t, rep, strings.Fields(clause)[1:])
				case strings.HasPrefix(clause, "dep "):
					checkDep(t, rep, clause)
				case strings.HasPrefix(clause, "summary "):
					checkSummary(t, rep, strings.Fields(clause)[1:])
				default:
					t.Errorf("unknown signal clause %q", clause)
				}
			}
			if rep.Outcome == "fatal" && !wantFatal {
				t.Errorf("fixture is fatal (%+v) but its signals do not say outcome=fatal", rep.Fatal)
			}
		})
	}
}

func findItem(rep status.Report, id string) *status.Item {
	for i := range rep.Items {
		if rep.Items[i].ID == id {
			return &rep.Items[i]
		}
	}
	return nil
}

func checkItem(t *testing.T, rep status.Report, f []string) {
	t.Helper()
	if len(f) == 0 {
		t.Error("item clause needs an ID")
		return
	}
	it := findItem(rep, f[0])
	if it == nil {
		t.Errorf("item %s not in report", f[0])
		return
	}
	for _, kv := range f[1:] {
		k, v, _ := strings.Cut(kv, "=")
		switch k {
		case "status":
			if it.Status != v {
				t.Errorf("item %s status %q, want %q", it.ID, it.Status, v)
			}
		case "tasks":
			got := strconv.Itoa(it.Tasks.Checked) + "/" + strconv.Itoa(it.Tasks.Total)
			if got != v {
				t.Errorf("item %s tasks %s, want %s", it.ID, got, v)
			}
		case "state":
			if it.Tasks.State != v {
				t.Errorf("item %s task state %q, want %q", it.ID, it.Tasks.State, v)
			}
		default:
			t.Errorf("unknown item key %q", k)
		}
	}
}

func checkDiag(t *testing.T, rep status.Report, f []string) {
	t.Helper()
	code, item := f[0], ""
	for _, kv := range f[1:] {
		if k, v, ok := strings.Cut(kv, "="); ok && k == "item" {
			item = v
		}
	}
	for _, d := range rep.Diagnostics {
		if d.Code == code && (item == "" || d.ItemID == item) {
			return
		}
	}
	var got []string
	for _, d := range rep.Diagnostics {
		got = append(got, d.Code+"@"+d.ItemID)
	}
	t.Errorf("diagnostic %s (item %q) not found; got %v", code, item, got)
}

func checkDep(t *testing.T, rep status.Report, clause string) {
	t.Helper()
	m := depClauseRE.FindStringSubmatch(clause)
	if m == nil {
		t.Errorf("bad dep clause %q", clause)
		return
	}
	it := findItem(rep, m[1])
	if it == nil {
		t.Errorf("item %s not in report", m[1])
		return
	}
	for _, d := range it.Dependencies {
		if d.ID == m[2] {
			if d.State != m[3] {
				t.Errorf("dependency %s->%s state %q, want %q", m[1], m[2], d.State, m[3])
			}
			return
		}
	}
	t.Errorf("item %s has no dependency on %s", m[1], m[2])
}

func checkSummary(t *testing.T, rep status.Report, f []string) {
	t.Helper()
	if rep.Summary == nil {
		t.Error("summary clause on a report without a summary")
		return
	}
	got := map[string]int{
		"active":       rep.Summary.ActiveItems,
		"verified":     rep.Summary.VerifiedActiveItems,
		"blocked":      rep.Summary.BlockedCount,
		"human_review": rep.Summary.HumanReviewCount,
	}
	for _, kv := range f {
		k, v, _ := strings.Cut(kv, "=")
		g, ok := got[k]
		if !ok {
			t.Errorf("unknown summary key %q", k)
			continue
		}
		if strconv.Itoa(g) != v {
			t.Errorf("summary %s %d, want %s", k, g, v)
		}
	}
}
