package extension

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// FR-021, FR-022, FR-025, FR-027, contracts/reference-files.md.

var reviewSections = []string{
	"progress", "blockers", "human-review", "diagnostics", "ready-to-start",
	"unchecked-tasks", "missing-artifacts", "evidence-gaps",
	"dependency-chains", "supersession", "constraints",
}

var roleNames = []string{"default", "pm", "developer", "qa", "architect"}

func ticked(cell string) []string {
	var out []string
	for _, m := range backtickRE.FindAllStringSubmatch(cell, -1) {
		out = append(out, m[1])
	}
	return out
}

func TestRolesUseKnownSectionsAndAllowedFollowUps(t *testing.T) {
	path := extPath("reference", "roles.md")
	got := map[string]bool{}
	for _, r := range tableUnder(t, path, "Review sections") {
		got[stripTicks(r["Section"])] = true
	}
	if missing, extra := diff(setOf(reviewSections), got), diff(got, setOf(reviewSections)); len(missing)+len(extra) > 0 {
		t.Errorf("Review sections: missing %v, extra %v", missing, extra)
	}
	roles := map[string]bool{}
	for _, r := range tableUnder(t, path, "Roles") {
		role := stripTicks(r["Role"])
		roles[role] = true
		lead := ticked(r["Leads with"])
		if len(lead) == 0 {
			t.Errorf("role %s: Leads with lists no sections", role)
		}
		for _, s := range lead {
			if !got[s] {
				t.Errorf("role %s leads with unknown section %q", role, s)
			}
		}
		fu := slashRefs(r["First follow-ups"])
		if len(fu) == 0 {
			t.Errorf("role %s: no first follow-ups", role)
		}
		for _, ref := range fu {
			if !allowedLinks[ref] {
				t.Errorf("role %s follow-up %s is not an allowed command", role, ref)
			}
		}
		for _, ref := range slashRefs(r["Role checks"]) {
			if !allowedLinks[ref] {
				t.Errorf("role %s checks name %s, which is not an allowed command", role, ref)
			}
		}
	}
	if missing, extra := diff(setOf(roleNames), roles), diff(roles, setOf(roleNames)); len(missing)+len(extra) > 0 {
		t.Errorf("Roles: missing %v, extra %v", missing, extra)
	}
}

var partnerCmdRE = regexp.MustCompile(`^speckit\.([a-z0-9-]+)\.([a-z0-9-]+)$`)

func TestPartnerRowsAreWellFormedAndConfined(t *testing.T) {
	path := extPath("reference", "roles.md")
	cols := []string{"Role", "Extension", "Version checked", "Catalog", "Install", "Requirements", "Commands", "Adds"}
	partnerIDs := map[string]bool{}
	for _, r := range tableUnder(t, path, "Partners") {
		for _, c := range cols {
			if strings.TrimSpace(r[c]) == "" {
				t.Errorf("partner row %v: empty %q", r, c)
			}
		}
		role, id := stripTicks(r["Role"]), stripTicks(r["Extension"])
		if !setOf(roleNames[1:])[role] {
			t.Errorf("partner %s: unknown role %q", id, role)
		}
		partnerIDs[id] = true
		cmds := ticked(r["Commands"])
		if len(cmds) == 0 {
			t.Errorf("partner %s lists no commands", id)
		}
		for _, c := range cmds {
			m := partnerCmdRE.FindStringSubmatch(c)
			if m == nil || m[1] != id {
				t.Errorf("partner %s command %q must be speckit.%s.<name>", id, c, id)
			}
		}
		if strings.Contains(r["Catalog"], "community") && !strings.Contains(strings.ToLower(r["Install"]+r["Requirements"]+r["Adds"]), "not vetted") {
			t.Errorf("community partner %s must carry a 'not vetted by Spec Kit' note", id)
		}
	}
	// Partner commands must never be hard links in command files.
	files, _ := filepath.Glob(extPath("commands", "*.md"))
	for _, f := range files {
		text := readFile(t, f)
		for id := range partnerIDs {
			if strings.Contains(text, "speckit."+id+".") {
				t.Errorf("%s names a %s command; partners belong only in reference/roles.md", filepath.Base(f), id)
			}
		}
	}
}
