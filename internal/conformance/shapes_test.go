package conformance

import (
	"regexp"
	"sort"
	"strings"
)

// This classifier is deliberately independent of internal/status: it names
// the constructs a file uses (research.md R5) so the shape inventory and the
// drift check do not rely on the parser they are meant to check.

// shape is one classified line.
type shape struct {
	Key  string
	Line int
	Text string
}

var (
	reHeading    = regexp.MustCompile(`^ {0,3}(#{1,6})(?:[ \t]|$)`)
	reItem       = regexp.MustCompile(`^### [0-9]{3} — .*\[status: [^\]]+\]\s*$`)
	reTask       = regexp.MustCompile(`^(\s*)[-*+] \[([ xX])\] (.*)$`)
	reTaskID     = regexp.MustCompile(`^T[0-9]+\b`)
	reMarker     = regexp.MustCompile(`\[(P|US[0-9]+|Story)\]`)
	reBracket    = regexp.MustCompile(`\[[^\]]+\]`)
	reBoldField  = regexp.MustCompile(`^- \*\*([^*:]+):\*\*`)
	rePlainField = regexp.MustCompile(`^- ([A-Za-z][A-Za-z ()-]*):(?:\s|$)`)
	reOrdered    = regexp.MustCompile(`^\s*[0-9]+[.)] `)
	reBullet     = regexp.MustCompile(`^\s*[-*+] `)
	reThematic   = regexp.MustCompile(`^ {0,3}([-*_])(?:\s*[-*_]){2,}\s*$`)
	reHTML       = regexp.MustCompile(`^ {0,3}</?[A-Za-z]`)
	reDelimiter  = regexp.MustCompile(`^\s*\|?\s*:?-+:?\s*(\|\s*:?-+:?\s*)+\|?\s*$`)
)

// classify returns the shapes used by a file, one entry per occurrence.
func classify(text string) []shape {
	var out []shape
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	emit := func(key string, i int) { out = append(out, shape{Key: key, Line: i + 1, Text: lines[i]}) }
	inFence, inComment, inItem, afterField := false, false, false, false
	fenceMark := ""
	start := 0
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for j := 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == "---" {
				emit("front-matter", 0)
				start = j + 1
				break
			}
		}
	}
	for i := start; i < len(lines); i++ {
		line := lines[i]
		trim := strings.TrimSpace(line)
		if inFence {
			if strings.HasPrefix(trim, fenceMark) && strings.Trim(trim, fenceMark[:1]) == "" {
				inFence = false
			}
			continue
		}
		if inComment {
			if strings.Contains(line, "-->") {
				inComment = false
			}
			continue
		}
		if trim == "" {
			afterField = false
			continue
		}
		if afterField && (line[0] == ' ' || line[0] == '\t') && !reBullet.MatchString(line) {
			emit("field-continuation", i)
			continue
		}
		afterField = false
		switch {
		case strings.HasPrefix(trim, "```") || strings.HasPrefix(trim, "~~~"):
			emit("fence", i)
			inFence, fenceMark = true, trim[:3]
		case strings.HasPrefix(trim, "<!--"):
			emit("html-comment", i)
			if !strings.Contains(trim[4:], "-->") {
				inComment = true
			}
		case reItem.MatchString(line):
			emit("item-heading", i)
			inItem = true
		case reHeading.MatchString(line):
			level := len(reHeading.FindStringSubmatch(line)[1])
			emit("heading:"+string(rune('0'+level)), i)
			inItem = false
		case reThematic.MatchString(line):
			emit("thematic-break", i)
			inItem = false
		case reTask.MatchString(line):
			m := reTask.FindStringSubmatch(line)
			key := "task-line"
			if m[2] != " " {
				key = "task-checked"
			}
			emit(key, i)
			if m[1] != "" {
				emit("task-nested", i)
			}
			body := m[3]
			if reTaskID.MatchString(body) {
				emit("task-id:T###", i)
			}
			for _, mk := range reMarker.FindAllStringSubmatch(body, -1) {
				name := mk[1]
				if strings.HasPrefix(name, "US") && name != "Story" {
					name = "USn"
				}
				emit("task-marker:["+name+"]", i)
			}
			if rest := reMarker.ReplaceAllString(body, ""); reBracket.MatchString(rest) {
				emit("task-text-brackets", i)
			}
		case strings.HasPrefix(trim, ">"):
			emit("block-quote", i)
		case strings.Contains(trim, "|") && (reDelimiter.MatchString(line) || (i+1 < len(lines) && reDelimiter.MatchString(lines[i+1]))):
			emit("table", i)
		case reHTML.MatchString(line):
			emit("html-block", i)
		case inItem && (reBoldField.MatchString(line) || rePlainField.MatchString(line)):
			if m := reBoldField.FindStringSubmatch(line); m != nil {
				emit("field:"+strings.TrimSpace(m[1]), i)
				emit("field-style:bold", i)
			} else {
				emit("field:"+strings.TrimSpace(rePlainField.FindStringSubmatch(line)[1]), i)
				emit("field-style:plain", i)
			}
			afterField = true
		case reOrdered.MatchString(line):
			emit("ordered-list", i)
		case reBullet.MatchString(line):
			emit("bullet", i)
		}
	}
	return out
}

// shapeKeys returns the distinct shape keys a file uses, with the first line
// where each appears.
func shapeKeys(text string) map[string]shape {
	keys := map[string]shape{}
	for _, s := range classify(text) {
		if _, seen := keys[s.Key]; !seen {
			keys[s.Key] = s
		}
	}
	return keys
}

func sortedShapeKeys(m map[string]shape) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
