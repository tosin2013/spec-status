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
			width += 4 - width%4 // advance to the next tab stop
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
	return openFenceIndent(line, 3)
}

// openFenceIndent is openFence with a caller-chosen indentation limit. Task
// files allow any indentation, because a fence inside a nested list item is
// indented past three columns.
func openFenceIndent(line string, maxIndent int) (n int, ch byte, ok bool) {
	width, b := leadingIndent(line)
	if width > maxIndent {
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
	_, b := leadingIndent(line)
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

// lookalikes maps Unicode spaces, dashes, and invisible characters to their
// ASCII forms, so "superseded‑by" (U+2011) or "Depends\u00a0on" cannot slip
// past the item and field checks. It is used for recognition only; stored
// text is never rewritten.
var lookalikes = strings.NewReplacer(
	"\u00a0", " ", "\u1680", " ", "\u2000", " ", "\u2001", " ", "\u2002", " ", "\u2003", " ",
	"\u2004", " ", "\u2005", " ", "\u2006", " ", "\u2007", " ", "\u2008", " ", "\u2009", " ",
	"\u200a", " ", "\u202f", " ", "\u205f", " ", "\u3000", " ",
	"\u2010", "-", "\u2011", "-", "\u2012", "-", "\u2212", "-", "\ufe63", "-", "\uff0d", "-",
	"\u200b", "", "\u200c", "", "\u200d", "", "\u2060", "", "\ufeff", "", "\u00ad", "",
	"\uff1a", ":", "\u2236", ":", "\ua789", ":", "\ufe13", ":", "\ufe55", ":", "\u02f8", ":",
	// Cyrillic and Greek letters that look like Latin ones, so a label such
	// as "Dеpends on" (Cyrillic е) is still recognized as a field label.
	"\u0430", "a", "\u0435", "e", "\u043e", "o", "\u0440", "p", "\u0441", "c",
	"\u0443", "y", "\u0445", "x", "\u0456", "i", "\u0458", "j", "\u0455", "s",
	"\u0410", "A", "\u0412", "B", "\u0415", "E", "\u041a", "K", "\u041c", "M",
	"\u041d", "H", "\u041e", "O", "\u0420", "P", "\u0421", "C", "\u0422", "T",
	"\u0425", "X", "\u0406", "I", "\u0405", "S",
	"\u03bf", "o", "\u03b1", "a", "\u03b5", "e", "\u03b9", "i", "\u03ba", "k",
	"\u039f", "O", "\u0391", "A", "\u0392", "B", "\u0395", "E", "\u0397", "H",
	"\u0399", "I", "\u039a", "K", "\u039c", "M", "\u039d", "N", "\u03a1", "P",
	"\u03a4", "T", "\u03a7", "X", "\u03a5", "Y", "\u0396", "Z",
)

func normalizeLookalikes(s string) string { return lookalikes.Replace(s) }

// idDash reports text that starts like an item: three digits, a dash, a space.
// Headings accept an em dash, en dash, or hyphen; other lines accept only an
// em or en dash, so prose such as "100 - 200 ms" stays prose.
func idDash(text string, hyphen bool) bool {
	if len(text) < 4 || !isDigit3(text[:3]) {
		return false
	}
	rest := strings.TrimLeft(text[3:], " \t")
	dashes := []string{"—", "–"}
	if hyphen {
		dashes = append(dashes, "-")
	}
	for _, d := range dashes {
		if strings.HasPrefix(rest, d) {
			after := rest[len(d):]
			// A dash then a space, as in "001 — Title". "003–005" is a range.
			return after == "" || after[0] == ' ' || after[0] == '\t'
		}
	}
	return false
}

// stripMarkers removes leading indentation, block-quote markers, one list
// marker, and leading emphasis or code characters. It reports whether a list
// or quote marker was present.
func stripMarkers(line string) (rest string, marked bool) {
	rest = strings.TrimLeft(line, " \t")
	for strings.HasPrefix(rest, ">") {
		rest = strings.TrimLeft(rest[1:], " \t")
		marked = true
	}
	switch {
	case len(rest) > 1 && (rest[0] == '-' || rest[0] == '*' || rest[0] == '+') && (rest[1] == ' ' || rest[1] == '\t'):
		rest = strings.TrimLeft(rest[2:], " \t")
		marked = true
	default:
		j := 0
		for j < len(rest) && j < 9 && rest[j] >= '0' && rest[j] <= '9' {
			j++
		}
		if j > 0 && j+1 < len(rest) && (rest[j] == '.' || rest[j] == ')') && (rest[j+1] == ' ' || rest[j+1] == '\t') {
			rest = strings.TrimLeft(rest[j+2:], " \t")
			marked = true
		}
	}
	return strings.TrimLeft(rest, "*_`"), marked
}

// itemLikeLine reports a line that looks like a roadmap item but is not the
// supported "### NNN — Title [status: token]" heading at column zero.
// matchItem must be tried first.
func itemLikeLine(line string) bool {
	line = normalizeLookalikes(line)
	if _, _, text, ok := heading(line); ok {
		return strings.Contains(text, "[status:") || idDash(text, true)
	}
	if strings.Contains(withoutCodeSpans(line), "[status:") {
		return true
	}
	rest, _ := stripMarkers(line)
	return idDash(rest, false)
}

// tableItemRow reports a pipe-table row whose first cell is an item ID or
// which carries a status marker.
func tableItemRow(line string) bool {
	trim := strings.TrimSpace(line)
	if !strings.Contains(trim, "|") {
		return false
	}
	if strings.Contains(trim, "[status:") {
		return true
	}
	cells := tableCells(trim)
	return len(cells) > 1 && len(cells[0]) == 3 && isDigit3(cells[0])
}

func tableCells(trim string) []string {
	trim = strings.TrimSuffix(strings.TrimPrefix(trim, "|"), "|")
	cells := strings.Split(trim, "|")
	for i := range cells {
		cells[i] = strings.Trim(strings.TrimSpace(cells[i]), "*_`")
	}
	return cells
}

// delimiterRow recognizes a GFM table delimiter row such as "--- | :---:".
func delimiterRow(line string) bool {
	trim := strings.TrimSpace(line)
	if !strings.Contains(trim, "|") {
		return false
	}
	for _, c := range tableCells(trim) {
		c = strings.TrimSuffix(strings.TrimPrefix(c, ":"), ":")
		if c == "" || strings.Trim(c, "-") != "" {
			return false
		}
	}
	return true
}

// setextUnderline recognizes "===" or "---" under a paragraph line.
func setextUnderline(line string) bool {
	width, b := leadingIndent(line)
	if width > 3 {
		return false
	}
	rest := strings.TrimRight(line[b:], " \t")
	return rest != "" && (strings.Trim(rest, "=") == "" || (len(rest) >= 2 && strings.Trim(rest, "-") == ""))
}

// paragraphLine reports a non-blank line that could carry a setext underline.
func paragraphLine(line string) bool {
	if strings.TrimSpace(line) == "" {
		return false
	}
	if width, _ := leadingIndent(line); width > 3 {
		return false
	}
	if _, _, _, ok := heading(line); ok {
		return false
	}
	_, marked := stripMarkers(line)
	return !marked
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

// metricFields change counts or relationships; their labels are strict
// everywhere. The descriptive fields only carry text.
var metricFields = map[string]bool{"Depends on": true, "Spec dir": true, "superseded-by": true}

// fieldPrefix applies the owner-approved prefix rule (7 November 2026): a
// list or quote line whose first words, after removing markup, case, and
// punctuation, name a field. matchField must be tried first; a fieldPrefix
// line that matchField rejected is a near miss.
//
// The multi-word metric labels match on a word boundary, so "- Depends on
// 002" is caught. The one-word descriptive labels match only when a
// separator follows, so "- Scope creep is a risk" stays prose.
func fieldPrefix(line string) (field string, ok bool) {
	rest, marked := stripMarkers(normalizeLookalikes(line))
	if !marked {
		return "", false
	}
	var b strings.Builder
	space := false
	for i := 0; i < len(rest) && b.Len() < 48; i++ {
		c := rest[i]
		switch {
		case c == '*' || c == '`':
			continue
		case c == '-' || c == '_' || c == ' ' || c == '\t':
			if !space && b.Len() > 0 {
				b.WriteByte(' ')
			}
			space = true
			continue
		case c >= 'A' && c <= 'Z':
			c += 'a' - 'A'
		}
		space = false
		b.WriteByte(c)
	}
	norm := b.String()
	boundary := func(p string) bool {
		if !strings.HasPrefix(norm, p) {
			return false
		}
		if len(norm) == len(p) {
			return true
		}
		c := norm[len(p)]
		if c >= 'a' && c <= 'z' {
			return false
		}
		// Without a colon, the words are a field label only when a value
		// follows: a digit, a path, "none", or a separator. "Depends on
		// whether the vendor ships" is prose (review 9 D6).
		after := strings.TrimLeft(norm[len(p):], " ")
		if after == "" {
			return true
		}
		if f := after[0]; !(f >= 'a' && f <= 'z') {
			return true
		}
		word := strings.Fields(after)[0]
		return strings.Contains(word, "/") || strings.TrimRight(word, ".,;") == "none"
	}
	for _, m := range []struct{ prefix, field string }{
		{"depends on", "Depends on"}, {"dependson", "Depends on"},
		{"spec dir", "Spec dir"}, {"specdir", "Spec dir"},
		{"superseded by", "superseded-by"}, {"supersededby", "superseded-by"},
	} {
		if boundary(m.prefix) {
			return m.field, true
		}
	}
	for _, m := range []struct{ prefix, field string }{
		{"description", "Description"}, {"outcome", "Outcome"}, {"scope (in)", "Scope"}, {"scope(in)", "Scope"}, {"scope", "Scope"},
	} {
		if strings.HasPrefix(norm, m.prefix) {
			after := strings.TrimLeft(norm[len(m.prefix):], " ")
			if after == "" || strings.ContainsAny(after[:1], ":=") || strings.HasPrefix(after, "—") || strings.HasPrefix(after, "–") {
				return m.field, true
			}
		}
	}
	return "", false
}

// descriptiveLabels are speckit-roadmap fields that carry text only. With
// the recognized fields they form the closed set of labels allowed at the
// top level of an item body.
var descriptiveLabels = []string{"Scope (out)", "Governed by", "Addresses", "Notes"}

// descriptiveField reports "- Label: text" for a descriptive label, in the
// same plain and bold spellings matchField accepts.
func descriptiveField(line string) bool {
	if !strings.HasPrefix(line, "- ") {
		return false
	}
	body := line[2:]
	for _, l := range descriptiveLabels {
		for _, p := range []string{l + ":", "**" + l + ":**", "**" + l + "**:"} {
			if strings.HasPrefix(body, p) {
				return true
			}
		}
	}
	return false
}

// thematicBreak recognizes "---", "***", or "___" (spaces allowed).
func thematicBreak(line string) bool {
	width, b := leadingIndent(line)
	if width > 3 {
		return false
	}
	rest := strings.ReplaceAll(strings.ReplaceAll(line[b:], " ", ""), "\t", "")
	if len(rest) < 3 {
		return false
	}
	return strings.Trim(rest, "-") == "" || strings.Trim(rest, "*") == "" || strings.Trim(rest, "_") == ""
}
