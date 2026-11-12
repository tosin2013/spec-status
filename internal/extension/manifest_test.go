package extension

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// FR-001, FR-003, FR-016, contracts/extension-manifest.md.

func TestManifestDeclaresExactlyTheFiveCommands(t *testing.T) {
	cmds := manifestCommands(t, extPath("extension.yml"))
	got := map[string]bool{}
	for _, c := range cmds {
		got[c.name] = true
		if want := "commands/" + c.name + ".md"; c.file != want {
			t.Errorf("%s: file %q, want %q", c.name, c.file, want)
		}
		if _, err := os.Stat(extPath(filepath.FromSlash(c.file))); err != nil {
			t.Errorf("%s: %v", c.name, err)
		}
	}
	if missing, extra := diff(setOf(extensionCommands), got), diff(got, setOf(extensionCommands)); len(missing)+len(extra) > 0 {
		t.Errorf("manifest commands: missing %v, extra %v", missing, extra)
	}
	files, _ := filepath.Glob(extPath("commands", "*.md"))
	for _, f := range files {
		if name := strings.TrimSuffix(filepath.Base(f), ".md"); !got[name] {
			t.Errorf("%s is not declared in extension.yml", filepath.Base(f))
		}
	}
}

func TestManifestHasNoRequiredExtensionsScriptsOrMandatoryHooks(t *testing.T) {
	text := readFile(t, extPath("extension.yml"))
	for _, l := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(l)
		switch {
		case strings.HasPrefix(l, "  extensions:"):
			t.Error("requires.extensions is not allowed: partners are optional (FR-026)")
		case strings.HasPrefix(l, "  scripts:"):
			t.Error("provides.scripts is not allowed")
		case trim == "optional: false":
			t.Error("mandatory hooks are not allowed (FR-003)")
		}
	}
	if !regexp.MustCompile(`(?m)^schema_version:\s*'?1\.0'?\s*$`).MatchString(text) {
		t.Error("schema_version must be '1.0'")
	}
	if !regexp.MustCompile(`(?m)^\s+id:\s*"?spec-status"?\s*$`).MatchString(text) {
		t.Error("extension.id must be spec-status")
	}
	if !regexp.MustCompile(`(?m)^\s+speckit_version:\s*'>=1\.2\.0'\s*$`).MatchString(text) {
		t.Error("requires.speckit_version must be '>=1.2.0'")
	}
}

var (
	manifestVersionRE = regexp.MustCompile(`(?m)^\s+version:\s*"?([0-9][^"\s]*)"?\s*$`)
	readmeVersionRE   = regexp.MustCompile(`(?m)^# spec-status extension v([0-9][^\s]*)\s*$`)
)

// manifestVersion returns extension.version: the first "version:" line
// inside the extension: block.
func manifestVersion(t *testing.T) string {
	t.Helper()
	text := readFile(t, extPath("extension.yml"))
	block := text
	if i := strings.Index(text, "\nextension:"); i >= 0 {
		block = text[i:]
		if j := strings.Index(block[1:], "\nrequires:"); j >= 0 {
			block = block[:j+1]
		}
	}
	m := manifestVersionRE.FindStringSubmatch(block)
	if m == nil {
		t.Fatal("extension.version not found")
	}
	return m[1]
}

func TestManifestVersionMatchesReadme(t *testing.T) {
	v := manifestVersion(t)
	m := readmeVersionRE.FindStringSubmatch(readFile(t, extPath("README.md")))
	if m == nil {
		t.Fatal("extensions/spec-status/README.md needs a '# spec-status extension v<version>' heading")
	}
	if m[1] != v {
		t.Errorf("README heading v%s, extension.yml version %s", m[1], v)
	}
}
