package status

import "strings"

// This file holds the checks that make unsupported item-like input fail
// visibly instead of disappearing (constitution Art. III).

// leadingIndent returns the indentation width, counting a tab as four columns.
func leadingIndent(line string) (width, bytes int) {
	for bytes < len(line) {
		switch line[bytes] {
		case ' ':
			width++
		case '\t':
			width += 4
		default:
			return width, bytes
		}
		bytes++
	}
	return width, bytes
}

// openFence recognizes a code fence opening: at most three columns of
// indentation, then three or more backticks or tildes. A backtick fence's info
// string cannot contain a backtick, so "```x```" is inline code, not a fence.
func openFence(line string) (n int, ch byte, ok bool) {
	width, b := leadingIndent(line)
	if width > 3 {
		return 0, 0, false
	}
	rest := line[b:]
	if len(rest) < 3 || (rest[0] != '`' && rest[0] != '~') {
		return 0, 0, false
	}
	ch = rest[0]
	for n < len(rest) && rest[n] == ch {
		n++
	}
	if n < 3 {
		return 0, 0, false
	}
	if ch == '`' && strings.IndexByte(rest[n:], '`') >= 0 {
		return 0, 0, false
	}
	return n, ch, true
}

// closesFence reports whether line closes a fence opened with n of ch.
func closesFence(line string, ch byte, n int) bool {
	width, b := leadingIndent(line)
	if width > 3 {
		return false
	}
	rest := strings.TrimRight(line[b:], " \t")
	return len(rest) >= n && strings.Trim(rest, string(ch)) == ""
}

// heading recognizes an ATX heading with at most three columns of indentation.
func heading(line string) (level int, indent int, text string, ok bool) {
	width, b := leadingIndent(line)
	if width > 3 {
		return 0, 0, "", false
	}
	rest := line[b:]
	for level < len(rest) && rest[level] == '#' {
		level++
	}
	if level == 0 || level > 6 {
		return 0, 0, "", false
	}
	if level < len(rest) && rest[level] != ' ' && rest[level] != '\t' {
		return 0, 0, "", false
	}
	return level, width, strings.TrimSpace(rest[level:]), true
}

// startsWithItemID reports whether s begins with exactly three digits.
func startsWithItemID(s string) bool {
	if len(s) < 3 || !isDigit3(s[:3]) {
		return false
	}
	return len(s) == 3 || s[3] < '0' || s[3] > '9'
}

// itemLikeLine reports a line that looks like a roadmap item but is not the
// supported "### NNN — Title [status: token]" heading at column zero.
// matchItem must be tried first.
func itemLikeLine(line string) bool {
	if level, _, text, ok := heading(line); ok {
		if strings.Contains(text, "[status:") {
			return true
		}
		if level == 3 {
			// Existing rule: any level-three heading that starts with a digit.
			return text != "" && text[0] >= '0' && text[0] <= '9'
		}
		return startsWithItemID(text)
	}
	return strings.Contains(withoutCodeSpans(line), "[status:")
}

// tableItemRow reports a pipe-table row whose first cell is an item ID or
// which carries a status marker.
func tableItemRow(line string) bool {
	trim := strings.TrimSpace(line)
	if !strings.HasPrefix(trim, "|") {
		return false
	}
	if strings.Contains(trim, "[status:") {
		return true
	}
	cells := strings.Split(trim[1:], "|")
	return len(cells) > 1 && isDigit3(strings.TrimSpace(cells[0])) && len(strings.TrimSpace(cells[0])) == 3
}

// withoutCodeSpans removes `inline code` so a documented example such as
// `### NNN — Title [status: token]` is prose, not an item.
func withoutCodeSpans(line string) string {
	var b strings.Builder
	in := false
	for i := 0; i < len(line); i++ {
		if line[i] == '`' {
			in = !in
			continue
		}
		if !in {
			b.WriteByte(line[i])
		}
	}
	return b.String()
}

// fieldKeys are the normalized forms of the recognized field labels.
var fieldKeys = map[string]bool{
	"description": true, "outcome": true, "scope": true, "scope(in)": true,
	"dependson": true, "specdir": true, "supersededby": true,
}

// fieldLike reports whether line is a list item whose label is a recognized
// field under some other spelling: different case, spacing, hyphenation,
// list marker, or indentation. matchField must be tried first; a fieldLike
// line that matchField rejected is a near miss.
func fieldLike(line string) bool {
	_, b := leadingIndent(line)
	rest := line[b:]
	switch {
	case len(rest) > 1 && (rest[0] == '-' || rest[0] == '*' || rest[0] == '+') && (rest[1] == ' ' || rest[1] == '\t'):
		rest = rest[2:]
	default:
		j := 0
		for j < len(rest) && j < 9 && rest[j] >= '0' && rest[j] <= '9' {
			j++
		}
		if j == 0 || j+1 >= len(rest) || (rest[j] != '.' && rest[j] != ')') || (rest[j+1] != ' ' && rest[j+1] != '\t') {
			return false
		}
		rest = rest[j+2:]
	}
	colon := strings.IndexByte(rest, ':')
	if colon < 0 || colon > 40 {
		return false
	}
	label := strings.ToLower(rest[:colon])
	label = strings.NewReplacer("*", "", "_", "", "-", "", " ", "", "\t", "").Replace(label)
	return fieldKeys[label]
}
