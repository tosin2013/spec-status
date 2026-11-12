package extension

import (
	"path/filepath"
	"strings"
	"testing"
)

// FR-015, contracts/reference-files.md: reference/output.md must match the
// diagnostics and JSON contracts exactly, so the skills never drift.

func contractPath(name string) string {
	return filepath.Join(repoRoot, "specs", "001-spec-status", name)
}

func TestOutputReferenceCodesMatchDiagnosticsContract(t *testing.T) {
	want := map[string]string{}
	for _, r := range tableUnder(t, contractPath("diagnostics.md"), "Registry") {
		code := stripTicks(r["Code"])
		if strings.HasPrefix(code, "SS_") {
			want[code] = r["Severity / class"]
		}
	}
	if len(want) == 0 {
		t.Fatal("no SS_ codes found in diagnostics.md")
	}
	got := map[string]string{}
	for _, r := range tableUnder(t, extPath("reference", "output.md"), "Diagnostic codes") {
		got[stripTicks(r["Code"])] = r["Severity"]
	}
	for code, sev := range want {
		g, ok := got[code]
		switch {
		case !ok:
			t.Errorf("output.md is missing %s", code)
		case g != sev:
			t.Errorf("%s severity: output.md %q, contract %q", code, g, sev)
		}
	}
	for code := range got {
		if _, ok := want[code]; !ok {
			t.Errorf("output.md lists %s, which is not in diagnostics.md", code)
		}
	}
	for _, r := range tableUnder(t, extPath("reference", "output.md"), "Diagnostic codes") {
		for _, ref := range slashRefs(r["Fix"]) {
			if !allowedLinks[ref] {
				t.Errorf("%s fix names %s, which is not an allowed command", stripTicks(r["Code"]), ref)
			}
		}
	}
}

func TestOutputReferenceHasEveryJSONField(t *testing.T) {
	body := sectionBody(readFile(t, extPath("reference", "output.md")), "JSON fields")
	if body == "" {
		t.Fatal("output.md has no ## JSON fields section")
	}
	listed := map[string]bool{}
	for _, m := range backtickRE.FindAllStringSubmatch(body, -1) {
		listed[m[1]] = true
	}
	for _, heading := range []string{"Root document (required fields)", "Summary object", "Item object"} {
		for _, r := range tableUnder(t, contractPath("json-contract.md"), heading) {
			for _, f := range strings.Split(r["Field"], ",") {
				f = stripTicks(f)
				if f != "" && !listed[f] {
					t.Errorf("output.md JSON fields is missing %q (json-contract.md %s)", f, heading)
				}
			}
		}
	}
}

func TestOutputReferenceStatesAndStatuses(t *testing.T) {
	cases := map[string][]string{
		"Dependency states": {"satisfied", "unfinished", "missing", "ambiguous", "self", "cycle", "human_review"},
		"Task states":       {"read", "missing", "skipped", "not_linked"},
		"Status tokens":     {"undecided", "needs-info", "planned", "specced", "in-progress", "implemented", "verified", "deferred", "abandoned"},
	}
	for heading, want := range cases {
		got := map[string]bool{}
		for _, r := range tableUnder(t, extPath("reference", "output.md"), heading) {
			// The first column holds the value.
			got[stripTicks(firstCell(r, heading, t))] = true
		}
		if missing, extra := diff(setOf(want), got), diff(got, setOf(want)); len(missing)+len(extra) > 0 {
			t.Errorf("%s: missing %v, extra %v", heading, missing, extra)
		}
	}
}

// firstCell returns the value of the table's first column, named by the
// reference-files contract: "State" for state tables, "Status" for tokens.
func firstCell(r map[string]string, heading string, t *testing.T) string {
	t.Helper()
	key := "State"
	if heading == "Status tokens" {
		key = "Status"
	}
	v, ok := r[key]
	if !ok {
		t.Fatalf("%s table needs a %q column", heading, key)
	}
	return v
}
