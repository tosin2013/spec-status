package status

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func diagCodes(rep Report, itemID string) []string {
	var out []string
	for _, d := range rep.Diagnostics {
		if d.ItemID == itemID {
			out = append(out, d.Code)
		}
	}
	return out
}

func hasDiag(rep Report, itemID, code string) bool {
	for _, c := range diagCodes(rep, itemID) {
		if c == code {
			return true
		}
	}
	return false
}

func itemByID(t *testing.T, rep Report, id string) Item {
	t.Helper()
	for _, it := range rep.Items {
		if it.ID == id {
			return it
		}
	}
	t.Fatalf("item %s not in report", id)
	return Item{}
}

func TestBoldFieldLabelsMatchPlainLabels(t *testing.T) {
	files := map[string]string{"specs/a/spec.md": "# A\n"}
	plain := "### 001 — A [status: specced]\n" +
		"- Description: Desc.\n- Outcome: Out.\n- Scope: In.\n- Spec dir: specs/a\n" +
		"### 002 — B [status: planned]\n- Depends on: 001\n"
	boldColon := "### 001 — A  [status: specced]\n" +
		"- **Description:** Desc.\n- **Outcome:** Out.\n- **Scope (in):** In.\n- **Scope (out):** Not this.\n- **Spec dir:** specs/a/\n" +
		"### 002 — B  [status: planned]\n- **Depends on:** 001.\n- **Notes:** prose\n"
	boldOutside := "### 001 — A [status: specced]\n" +
		"- **Description**: Desc.\n- **Outcome**: Out.\n- **Scope**: In.\n- **Spec dir**: `specs/a`\n" +
		"### 002 — B [status: planned]\n- **Depends on**: 001\n"
	want := Evaluate(Option{Root: project(t, plain, files), FailOn: "none"})
	for _, body := range []string{boldColon, boldOutside} {
		got := Evaluate(Option{Root: project(t, body, files), FailOn: "none"})
		if got.Outcome != "report" {
			t.Fatalf("not a report: %+v", got.Fatal)
		}
		for _, id := range []string{"001", "002"} {
			w, g := itemByID(t, want, id), itemByID(t, got, id)
			if w.Description != g.Description || w.Outcome != g.Outcome || w.Scope != g.Scope || w.SpecDir != g.SpecDir {
				t.Fatalf("item %s fields differ:\nplain %+v\nbold  %+v", id, w, g)
			}
			if !reflect.DeepEqual(depIDs(w), depIDs(g)) {
				t.Fatalf("item %s deps differ: %v vs %v", id, depIDs(w), depIDs(g))
			}
		}
		if !reflect.DeepEqual(diagCodes(want, "001"), diagCodes(got, "001")) || !reflect.DeepEqual(diagCodes(want, "002"), diagCodes(got, "002")) {
			t.Fatalf("diagnostics differ for %q:\nplain %+v\nbold  %+v", body, want.Diagnostics, got.Diagnostics)
		}
	}
}

func depIDs(it Item) []string {
	var out []string
	for _, d := range it.Dependencies {
		out = append(out, d.ID)
	}
	return out
}

func TestScopeAndScopeInRepeatIsFatal(t *testing.T) {
	body := "### 001 — A [status: planned]\n- Scope: one\n- **Scope (in):** two\n"
	rep := Evaluate(Option{Root: project(t, body, nil), FailOn: "none"})
	if rep.Outcome != "fatal" || rep.Fatal == nil || rep.Fatal.Code != "SS_ITEM_MALFORMED" {
		t.Fatalf("expected fatal repeat, got outcome %s %+v", rep.Outcome, rep.Fatal)
	}
}

func TestDependencyListNormalization(t *testing.T) {
	head := "### 001 — A [status: verified]\n### 002 — B [status: verified]\n### 003 — C [status: planned]\n"
	cases := []struct {
		value     string
		want      []string
		malformed bool
	}{
		{"none.", nil, false},
		{"None", nil, false},
		{"001, 002.", []string{"001", "002"}, false},
		{"001 (verified), 002 (verified).", []string{"001", "002"}, false},
		{"001, 002-brief (verified)", nil, true},
		{"001 and 002", nil, true},
		{"", nil, true},
	}
	for _, c := range cases {
		rep := Evaluate(Option{Root: project(t, head+"- Depends on: "+c.value+"\n", nil), FailOn: "none"})
		if rep.Outcome != "report" {
			t.Fatalf("%q: not a report: %+v", c.value, rep.Fatal)
		}
		it := itemByID(t, rep, "003")
		if got := hasDiag(rep, "003", "SS_DEPENDENCY_MALFORMED"); got != c.malformed {
			t.Fatalf("%q: malformed=%v, want %v (%v)", c.value, got, c.malformed, diagCodes(rep, "003"))
		}
		if !reflect.DeepEqual(depIDs(it), c.want) {
			t.Fatalf("%q: deps %v, want %v", c.value, depIDs(it), c.want)
		}
	}
}

func TestSupersededByAcceptsTrailingPeriodOnly(t *testing.T) {
	head := "### 001 — A [status: implemented]\n"
	ok := Evaluate(Option{Root: project(t, head+"- **superseded-by:** 002.\n### 002 — B [status: planned]\n", nil), FailOn: "none"})
	if hasDiag(ok, "001", "SS_SUPERSESSION_MALFORMED") || itemByID(t, ok, "001").Active {
		t.Fatalf("trailing period should still supersede: %+v", ok.Diagnostics)
	}
	bad := Evaluate(Option{Root: project(t, head+"- superseded-by: 002 (verified)\n### 002 — B [status: planned]\n", nil), FailOn: "none"})
	if !hasDiag(bad, "001", "SS_SUPERSESSION_MALFORMED") {
		t.Fatalf("annotated supersession must stay malformed: %+v", bad.Diagnostics)
	}
}

// The fixture is the speckit-roadmap v0.1.0 repository's own roadmap and specs,
// copied unmodified. See testdata/speckit-roadmap-upstream/NOTICE.md.
func TestUpstreamSpeckitRoadmapFixture(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "speckit-roadmap-upstream")
	rep := Evaluate(Option{Root: root})
	if rep.Outcome != "report" || !rep.Complete {
		t.Fatalf("upstream fixture not a complete report: %+v", rep.Fatal)
	}
	if rep.Summary.TotalItems != 7 {
		t.Fatalf("items %d, want 7", rep.Summary.TotalItems)
	}
	for _, id := range []string{"002", "003", "004", "005"} {
		it := itemByID(t, rep, id)
		if !strings.HasPrefix(it.SpecDir, "specs/00") {
			t.Fatalf("%s spec dir not read: %q", id, it.SpecDir)
		}
		for _, code := range []string{"SS_SPEC_REFERENCE_MISSING", "SS_SPEC_DIRECTORY_MISSING", "SS_STAGE_FILE_MISSING"} {
			if hasDiag(rep, id, code) {
				t.Fatalf("%s reported %s although its spec files exist", id, code)
			}
		}
	}
	if !hasDiag(rep, "001", "SS_SPEC_REFERENCE_MISSING") {
		t.Fatal("001 has no spec dir; that mechanical finding must remain")
	}
	if !hasDiag(rep, "004", "SS_DEPENDENCY_MALFORMED") {
		t.Fatal("004 lists 003-brief; it must be flagged, not dropped")
	}
	if !hasDiag(rep, "007", "SS_DEPENDENCY_BLOCKED") || rep.Summary.BlockedCount < 1 {
		t.Fatalf("007 depends on abandoned 006; blocked %d", rep.Summary.BlockedCount)
	}
	if got := depIDs(itemByID(t, rep, "007")); !reflect.DeepEqual(got, []string{"002", "003", "004", "005", "006"}) {
		t.Fatalf("007 deps %v", got)
	}
	if rep.Summary.Tasks.FilesRead != 4 || rep.Summary.Tasks.Total == 0 {
		t.Fatalf("linked tasks not read: %+v", rep.Summary.Tasks)
	}
}
