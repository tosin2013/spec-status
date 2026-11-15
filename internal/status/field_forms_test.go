package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Feature 008: real-world field forms and small defect fixes
// (specs/008-field-forms-and-fixes/contracts/changes.md). Expectations are
// hand-written in testdata/field-forms/*/expected-status.json.

type fieldFormExpect struct {
	Exit  int    `json:"exit"`
	Code  string `json:"code"`
	Items map[string]struct {
		SupersededBy *string  `json:"superseded_by"`
		Dependencies []string `json:"dependencies"`
	} `json:"items"`
	Diagnostics []struct {
		Item string `json:"item"`
		Code string `json:"code"`
	} `json:"diagnostics"`
}

func fieldFormsRoot(name string) string {
	return filepath.Join("..", "..", "testdata", "field-forms", name)
}

func TestFieldFormsFixtures(t *testing.T) {
	dirs, err := filepath.Glob(fieldFormsRoot("*"))
	if err != nil {
		t.Fatal(err)
	}
	ran := 0
	for _, dir := range dirs {
		raw, err := os.ReadFile(filepath.Join(dir, "expected-status.json"))
		if err != nil {
			continue // README.md and other files
		}
		ran++
		name := filepath.Base(dir)
		t.Run(name, func(t *testing.T) {
			var want fieldFormExpect
			if err := json.Unmarshal(raw, &want); err != nil {
				t.Fatal(err)
			}
			rep := Evaluate(Option{Root: dir})
			if got := ExitCode(rep); got != want.Exit {
				t.Fatalf("exit %d, want %d (fatal %+v, diags %+v)", got, want.Exit, rep.Fatal, rep.Diagnostics)
			}
			if want.Code != "" {
				if rep.Fatal == nil || rep.Fatal.Code != want.Code {
					t.Fatalf("fatal %+v, want %s", rep.Fatal, want.Code)
				}
				return
			}
			for id, w := range want.Items {
				var it *Item
				for i := range rep.Items {
					if rep.Items[i].ID == id {
						it = &rep.Items[i]
					}
				}
				if it == nil {
					t.Fatalf("item %s missing", id)
				}
				gotSup := it.SupersededBy
				if (w.SupersededBy == nil && gotSup != "") || (w.SupersededBy != nil && gotSup != *w.SupersededBy) {
					t.Errorf("%s superseded_by %q, want %v", id, gotSup, w.SupersededBy)
				}
				var gotDeps []string
				for _, d := range it.Dependencies {
					gotDeps = append(gotDeps, d.ID)
				}
				if strings.Join(gotDeps, ",") != strings.Join(w.Dependencies, ",") {
					t.Errorf("%s dependencies %v, want %v", id, gotDeps, w.Dependencies)
				}
			}
			var got, exp []string
			for _, d := range rep.Diagnostics {
				got = append(got, d.ItemID+" "+d.Code)
			}
			for _, d := range want.Diagnostics {
				exp = append(exp, d.Item+" "+d.Code)
			}
			sort.Strings(got)
			sort.Strings(exp)
			if strings.Join(got, "|") != strings.Join(exp, "|") {
				t.Errorf("diagnostics\n got %v\nwant %v", got, exp)
			}
		})
	}
	if ran < 9 {
		t.Fatalf("only %d field-form fixtures ran", ran)
	}
}

// US1: an accepted spelling gives the same JSON as superseded-by on the same line.
func TestSupersededBySpellingsMatchCanonical(t *testing.T) {
	body := func(label string) string {
		return "### 001 — A [status: planned]\n- Description: d\n- Outcome: o\n- Scope: s\n" + label + " 002\n\n### 002 — B [status: planned]\n- Description: d\n- Outcome: o\n- Scope: s\n"
	}
	canon := project(t, body("- superseded-by:"), nil)
	want := string(JSON(Evaluate(Option{Root: canon, CWD: canon})))
	for _, label := range []string{"- **Superseded by:**", "- **Superseded by**:", "- Superseded by:", "- superseded by:", "- **superseded by:**"} {
		dir := project(t, body(label), nil)
		if got := string(JSON(Evaluate(Option{Root: dir, CWD: dir}))); got != want {
			t.Errorf("%q: JSON differs from superseded-by\n got %s\nwant %s", label, got, want)
		}
	}
	// Prose with the label outside a field line stays the C-10 fatal.
	dir := project(t, "### 001 — A [status: planned]\n- Description: d\n- Outcome: o\n- Scope: s\n- Notes: replaced, Superseded by: 002\n", nil)
	if rep := Evaluate(Option{Root: dir, CWD: dir}); rep.Fatal == nil || rep.Fatal.Code != "SS_ITEM_MALFORMED" {
		t.Fatalf("prose label: %+v", rep.Fatal)
	}
}

// US2 (SC-002): the raw pinned template, copied at test time, exits 0 with
// exactly two placeholder findings at the template's placeholder lines.
func TestRawTemplatePlaceholders(t *testing.T) {
	tpl, err := os.ReadFile(filepath.Join("..", "..", "testdata", "upstream", "speckit-roadmap", "v0.1.0", "templates", "roadmap-template.md"))
	if err != nil {
		t.Fatal(err)
	}
	dir := project(t, string(tpl), nil)
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if rep.Fatal != nil || ExitCode(rep) != 0 {
		t.Fatalf("raw template: fatal %+v exit %d diags %+v", rep.Fatal, ExitCode(rep), rep.Diagnostics)
	}
	var lines []int
	for _, d := range rep.Diagnostics {
		if d.Severity == "error" {
			t.Errorf("error from template text: %+v", d)
		}
		if d.Code == "SS_FIELD_PLACEHOLDER" {
			if d.Message != "Field holds a template placeholder and is treated as not set." || d.ItemID != "001" {
				t.Errorf("placeholder finding %+v", d)
			}
			lines = append(lines, d.Line)
		}
	}
	sort.Ints(lines)
	if len(lines) != 2 || lines[0] != 68 || lines[1] != 71 {
		t.Fatalf("placeholder lines %v, want [68 71]", lines)
	}
}

// US2: the placeholder rule and its valid-list exception (C-08).
func TestPlaceholderValue(t *testing.T) {
	for v, want := range map[string]bool{
		"[x]": true, "[]": true, "_x_": true, "*x*": true, "[specs/001-.../ — set once the spec exists; otherwise omit]": true,
		"_set once the spec exists_": true, "[x].": true,
		"**x**": false, "__x__": false, "_x*": false, "[x](y)": false, "x": false, "[a] [b]": false, "": false,
	} {
		if _, got := placeholderValue(v); got != want {
			t.Errorf("placeholderValue(%q) = %v, want %v", v, got, want)
		}
	}
	for _, dep := range []string{"[002]", "_002, 003_", "[none]", "*002*"} {
		dir := project(t, "### 001 — A [status: planned]\n- Description: d\n- Outcome: o\n- Scope: s\n- Depends on: "+dep+"\n\n### 002 — B [status: planned]\n- Description: d\n- Outcome: o\n- Scope: s\n\n### 003 — C [status: planned]\n- Description: d\n- Outcome: o\n- Scope: s\n", nil)
		rep := Evaluate(Option{Root: dir, CWD: dir})
		if !hasCode(rep, "SS_DEPENDENCY_MALFORMED", "001") || hasCode(rep, "SS_FIELD_PLACEHOLDER", "001") {
			t.Errorf("Depends on %q: a valid enclosed list must stay malformed, got %+v", dep, rep.Diagnostics)
		}
	}
}

// US3: duplicate message and location; lists without repeats unchanged.
func TestDuplicateDependencyFinding(t *testing.T) {
	rep := Evaluate(Option{Root: fieldFormsRoot("duplicate-deps")})
	var dup []Diagnostic
	for _, d := range rep.Diagnostics {
		if d.Code == "SS_DEPENDENCY_DUPLICATE" {
			dup = append(dup, d)
		}
	}
	if len(dup) != 1 || dup[0].Severity != "warning" || dup[0].Message != "Dependency list names the same item more than once; it is counted once." || dup[0].Line != 5 {
		t.Fatalf("duplicate finding %+v", dup)
	}
	dir := project(t, "### 001 — A [status: planned]\n- Description: d\n- Outcome: o\n- Scope: s\n- Depends on: 002, 003\n\n### 002 — B [status: planned]\n- Description: d\n- Outcome: o\n- Scope: s\n\n### 003 — C [status: planned]\n- Description: d\n- Outcome: o\n- Scope: s\n", nil)
	if hasCode(Evaluate(Option{Root: dir, CWD: dir}), "SS_DEPENDENCY_DUPLICATE", "001") {
		t.Fatal("duplicate reported for a list without repeats")
	}
}

// A bracketed real path in Spec dir is a placeholder too (no path is guessed
// from markup); a specced item then still reports the missing spec link as
// an error, so it stays visible (contracts/changes.md §1, Doomer lead).
func TestBracketedSpecDirIsPlaceholder(t *testing.T) {
	dir := project(t, "### 001 — A [status: specced]\n- Description: d\n- Outcome: o\n- Scope: s\n- Spec dir: [specs/001-x/]\n", map[string]string{"specs/001-x/spec.md": "# s\n"})
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if !hasCode(rep, "SS_FIELD_PLACEHOLDER", "001") || !hasCode(rep, "SS_SPEC_REFERENCE_MISSING", "001") || ExitCode(rep) != 1 {
		t.Fatalf("bracketed Spec dir: exit %d diags %+v", ExitCode(rep), rep.Diagnostics)
	}
	if rep.Items[0].SpecDir != "" || rep.Items[0].Tasks.State != "not_linked" {
		t.Fatalf("placeholder must leave Spec dir unset: %+v", rep.Items[0])
	}
}
