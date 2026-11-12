// Package extension holds repository checks for the spec-status Spec Kit
// extension in extensions/spec-status (feature 005). It contains only tests
// and uses only the standard library; the spec binary does not import it.
package extension

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// repoRoot is the module root, two levels above this package.
var repoRoot = filepath.Join("..", "..")

func extPath(parts ...string) string {
	return filepath.Join(append([]string{repoRoot, "extensions", "spec-status"}, parts...)...)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

// sectionBody returns the text under the first "## heading" (exact match)
// up to the next "## " heading. It returns "" when the heading is absent.
func sectionBody(text, heading string) string {
	lines := strings.Split(text, "\n")
	start := -1
	for i, l := range lines {
		if h, ok := level2Heading(l); ok && h == heading {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return ""
	}
	var out []string
	for _, l := range lines[start:] {
		if _, ok := level2Heading(l); ok {
			break
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

// level2Heading reports whether line is a Markdown "## " heading, which may be
// indented by up to three spaces, and returns its text.
func level2Heading(line string) (string, bool) {
	trimmed := strings.TrimLeft(line, " ")
	if len(line)-len(trimmed) > 3 || !strings.HasPrefix(trimmed, "## ") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")), true
}

// splitRow splits a Markdown table row into trimmed cells.
func splitRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// readTable returns the rows of the first Markdown table in text, keyed by
// header cell text.
func readTable(text string) []map[string]string {
	var header []string
	var rows []map[string]string
	inTable := false
	for _, l := range strings.Split(text, "\n") {
		t := strings.TrimSpace(l)
		if !strings.HasPrefix(t, "|") {
			if inTable {
				break
			}
			continue
		}
		cells := splitRow(t)
		switch {
		case header == nil:
			header, inTable = cells, true
		case isDivider(cells):
		default:
			row := map[string]string{}
			for i, h := range header {
				if i < len(cells) {
					row[h] = cells[i]
				}
			}
			rows = append(rows, row)
		}
	}
	return rows
}

func isDivider(cells []string) bool {
	for _, c := range cells {
		if strings.Trim(c, "-: ") != "" {
			return false
		}
	}
	return true
}

// tableUnder returns the rows of the first table under "## heading".
func tableUnder(t *testing.T, path, heading string) []map[string]string {
	t.Helper()
	body := sectionBody(readFile(t, path), heading)
	if body == "" {
		t.Fatalf("%s: no section %q", path, heading)
	}
	rows := readTable(body)
	if len(rows) == 0 {
		t.Fatalf("%s: section %q has no table rows", path, heading)
	}
	return rows
}

var backtickRE = regexp.MustCompile("`([^`]+)`")

// stripTicks removes surrounding backticks and whitespace from a cell.
func stripTicks(s string) string {
	return strings.TrimSpace(strings.Trim(strings.TrimSpace(s), "`"))
}

// manifestCommand is one provides.commands entry of extension.yml.
type manifestCommand struct{ name, file string }

var (
	nameLineRE = regexp.MustCompile(`^\s*-\s*name:\s*"?([^"#]+?)"?\s*$`)
	fileLineRE = regexp.MustCompile(`^\s*file:\s*"?([^"#]+?)"?\s*$`)
)

// manifestCommands reads provides.commands from the fixed-format manifest
// with a line scan (contracts/extension-manifest.md), not a YAML parser.
func manifestCommands(t *testing.T, path string) []manifestCommand {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	var out []manifestCommand
	inProvides, inCommands := false, false
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := sc.Text()
		switch {
		case strings.HasPrefix(l, "provides:"):
			inProvides = true
			continue
		case inProvides && len(l) > 0 && l[0] != ' ' && l[0] != '#':
			inProvides, inCommands = false, false
		}
		if !inProvides {
			continue
		}
		if strings.TrimSpace(l) == "commands:" {
			inCommands = true
			continue
		}
		if inCommands && strings.HasPrefix(l, "  ") && !strings.HasPrefix(l, "   ") && strings.HasSuffix(strings.TrimSpace(l), ":") {
			inCommands = false // another provides.* key
		}
		if !inCommands {
			continue
		}
		if m := nameLineRE.FindStringSubmatch(l); m != nil {
			out = append(out, manifestCommand{name: strings.TrimSpace(m[1])})
		} else if m := fileLineRE.FindStringSubmatch(l); m != nil && len(out) > 0 {
			out[len(out)-1].file = strings.TrimSpace(m[1])
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

var slashRefRE = regexp.MustCompile(`/speckit[.-][a-z0-9.-]*[a-z0-9]`)

// slashRefs returns every /speckit… command reference in text, deduplicated
// and sorted.
func slashRefs(text string) []string {
	seen := map[string]bool{}
	for _, m := range slashRefRE.FindAllString(text, -1) {
		seen[m] = true
	}
	var out []string
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// spec-status commands, and the commands they may link to (FR-008).
var (
	extensionCommands = []string{
		"speckit.spec-status.review",
		"speckit.spec-status.explain",
		"speckit.spec-status.reconcile",
		"speckit.spec-status.unblock",
		"speckit.spec-status.verify-readiness",
	}
	allowedLinks = func() map[string]bool {
		m := map[string]bool{}
		for _, c := range extensionCommands {
			m["/"+c] = true
		}
		for _, c := range []string{"specify", "plan", "tasks", "implement", "analyze", "clarify"} {
			m["/speckit-"+c] = true
		}
		for _, c := range []string{"write", "brief", "debrief", "sync"} {
			m["/speckit.roadmap."+c] = true
		}
		return m
	}()
)

func setOf(xs []string) map[string]bool {
	m := map[string]bool{}
	for _, x := range xs {
		m[x] = true
	}
	return m
}

func sortedKeys(m map[string]bool) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// diff returns items in a but not b.
func diff(a, b map[string]bool) []string {
	var out []string
	for k := range a {
		if !b[k] {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}
