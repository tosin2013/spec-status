package extension

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// FR-004, FR-008, FR-016, FR-027, contracts/command-format.md.

var (
	frontmatterRE = regexp.MustCompile(`(?s)\A---\n(.*?)\n---\n`)
	descriptionRE = regexp.MustCompile(`(?m)^description:\s*\S`)
	firstStepRE   = regexp.MustCompile(`(?m)^1\.\s.*$`)
)

func commandFiles(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob(extPath("commands", "*.md"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no command files in %s", extPath("commands"))
	}
	return files
}

func TestCommandFilesHaveRequiredStructure(t *testing.T) {
	for _, f := range commandFiles(t) {
		name := filepath.Base(f)
		text := readFile(t, f)
		fm := frontmatterRE.FindStringSubmatch(text)
		if fm == nil || !descriptionRE.MatchString(fm[1]) {
			t.Errorf("%s: frontmatter with description is required", name)
		}
		for _, h := range []string{"User Input", "Goal", "Operating Constraints", "Outline", "Follow-ups"} {
			if sectionBody(text, h) == "" {
				t.Errorf("%s: missing or empty ## %s", name, h)
			}
		}
		if !strings.Contains(sectionBody(text, "User Input"), "$ARGUMENTS") {
			t.Errorf("%s: ## User Input must contain $ARGUMENTS", name)
		}
		if !strings.Contains(sectionBody(text, "Operating Constraints"), "reference/rules.md") {
			t.Errorf("%s: ## Operating Constraints must point to reference/rules.md", name)
		}
		if step := firstStepRE.FindString(sectionBody(text, "Outline")); !strings.Contains(step, "spec version") {
			t.Errorf("%s: step 1 of ## Outline must run spec version, got %q", name, step)
		}
	}
	if fm := frontmatterRE.FindStringSubmatch(readFile(t, extPath("commands", "speckit.spec-status.review.md"))); fm == nil || !strings.Contains(fm[1], "argument-hint:") {
		t.Error("review: frontmatter needs argument-hint for roles")
	}
}

func TestCommandLinksResolve(t *testing.T) {
	for _, f := range commandFiles(t) {
		for _, ref := range slashRefs(readFile(t, f)) {
			if !allowedLinks[ref] {
				t.Errorf("%s links %s, which is not an allowed command (partner commands belong only in reference/roles.md)", filepath.Base(f), ref)
			}
		}
	}
	for _, ref := range slashRefs(readFile(t, extPath("reference", "rules.md"))) {
		if !allowedLinks[ref] {
			t.Errorf("rules.md links %s, which is not an allowed command", ref)
		}
	}
}

func TestEveryCommandIsReachableFromReview(t *testing.T) {
	review := readFile(t, extPath("commands", "speckit.spec-status.review.md"))
	for _, c := range extensionCommands[1:] {
		if !strings.Contains(sectionBody(review, "Follow-ups"), "/"+c) {
			t.Errorf("review ## Follow-ups must name /%s", c)
		}
	}
}
