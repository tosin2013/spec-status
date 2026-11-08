package status

import (
	"strings"
	"unicode/utf8"
)

type rawItem struct {
	ID          string
	Title       string
	Status      string
	Line        int
	Description string
	Outcome     string
	Scope       string
	Deps        []string
	DepsLine    int
	DepsBad     bool
	SpecDir     string
	SpecLine    int
	HasSpec     bool
	SuperRaw    string
	SuperLine   int
	HasSuper    bool
	SuperBad    bool
	SuperID     string
	seen        map[string]bool
}

func parseRoadmap(rel string, data []byte) ([]rawItem, []Diagnostic, *Diagnostic) {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	if !utf8.Valid(data) {
		d := Diagnostic{Code: "SS_ENCODING_INVALID", Severity: "error", Message: "Roadmap encoding is not valid UTF-8.", Path: rel}
		return nil, nil, &d
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	var (
		items   []rawItem
		cur     *rawItem
		diags   []Diagnostic
		fence   byte
		fenceN  int
		inFence bool
	)
	flush := func() {
		if cur != nil {
			items = append(items, *cur)
			cur = nil
		}
	}
	fail := func(code, msg string, line int) *Diagnostic {
		d := Diagnostic{Code: code, Severity: "error", Message: msg, Path: rel, Line: line}
		return &d
	}
	fenceLine, commentLine, fenceIndent, commentIndent := 0, 0, 0, 0
	inComment, inTable := false, false
	prev := ""
	for i, line := range lines {
		lineNo := i + 1
		last := prev
		prev = line
		if inFence {
			if fenceIndent > 0 && strings.TrimSpace(line) != "" {
				if w, _ := leadingIndent(line); w < fenceIndent {
					// The list item holding the fence ended before the fence
					// closed. Renderers close it there; we fail visibly.
					return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Code fence inside a list item is not closed before the item ends.", fenceLine)
				}
			}
			if closesFence(line, fence, fenceN) {
				if w, _ := leadingIndent(line); w != fenceIndent {
					// CommonMark limits a closer's indentation; rather than
					// guess where a renderer ends the block, fail visibly.
					return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Code fence closer is indented differently from its opener.", lineNo)
				}
				inFence = false
			}
			continue
		}
		if inComment {
			if commentIndent > 0 && strings.TrimSpace(line) != "" {
				if w, _ := leadingIndent(line); w < commentIndent {
					return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "HTML comment inside a list item is not closed before the item ends.", commentLine)
				}
			}
			if strings.Contains(line, "-->") {
				inComment = false
			}
			continue
		}
		if n, ch, ok := openFence(line); ok {
			inFence, fence, fenceN, fenceLine = true, ch, n, lineNo
			fenceIndent, _ = leadingIndent(line)
			continue
		}
		// HTML comments are not rendered, so like fences they are not live
		// roadmap data. A comment opens only at three or fewer columns of
		// indentation; deeper "<!--" is indented code or continuation text.
		// An unclosed comment is fatal below.
		if w, b := leadingIndent(line); w <= 3 && strings.HasPrefix(line[b:], "<!--") {
			trim := strings.TrimSpace(line)
			if end := strings.Index(trim[4:], "-->"); end >= 0 {
				if cur != nil && strings.TrimSpace(trim[4+end+3:]) != "" {
					return nil, nil, fail("SS_ITEM_MALFORMED", "Text follows an HTML comment on the same line under an item.", lineNo)
				}
			} else {
				inComment, commentLine, commentIndent = true, lineNo, w
			}
			continue
		}
		if inTable {
			if strings.TrimSpace(line) != "" && strings.Contains(line, "|") {
				if tableItemRow(line) {
					return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Unsupported item-like table layout.", lineNo)
				}
				continue
			}
			inTable = false
		}
		if delimiterRow(line) {
			// A GFM table, with or without outer pipes. Its header is the
			// previous line; its rows follow until a blank or pipe-less line.
			if tableItemRow(last) {
				return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Unsupported item-like table layout.", lineNo-1)
			}
			inTable = true
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "|") && tableItemRow(line) {
			return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Unsupported item-like table layout.", lineNo)
		}
		if id, title, status, ok := matchItem(line); ok {
			flush()
			cur = &rawItem{ID: id, Title: title, Status: status, Line: lineNo, seen: map[string]bool{}}
			continue
		}
		if itemLikeLine(line) {
			return nil, nil, fail("SS_ITEM_MALFORMED", "Item heading does not match the supported shape.", lineNo)
		}
		if _, _, _, ok := heading(line); ok || (setextUnderline(line) && paragraphLine(last)) {
			// Any other heading, ATX or setext, ends the current item, so the
			// fields under it cannot attach to the item above.
			flush()
			continue
		}
		if _, marked := stripMarkers(line); marked && indentDepth(line) > maxNesting {
			return nil, nil, fail("SS_INPUT_LIMIT", "Parser nesting limit exceeded.", lineNo)
		}
		field, value, ok := matchField(line)
		if !ok {
			width, _ := leadingIndent(line)
			if near, isNear := fieldPrefix(line); isNear {
				// Metric fields are strict everywhere. A descriptive label is
				// strict only at the top level of an item; elsewhere it is prose.
				if metricFields[near] || (cur != nil && width == 0) {
					return nil, nil, fail("SS_ITEM_MALFORMED", "Field label does not match the supported spelling.", lineNo)
				}
				continue
			}
			if cur == nil {
				continue
			}
			// Closed item body (owner decision, 7 November 2026): under an item
			// heading, a line must be blank, indented continuation text, a
			// known descriptive field, or a thematic break that ends the item.
			// Anything else is fatal instead of being read as prose.
			switch {
			case strings.TrimSpace(line) == "", width > 0, descriptiveField(line):
				continue
			case thematicBreak(line):
				flush()
				continue
			}
			return nil, nil, fail("SS_ITEM_MALFORMED", "Line under an item is not a supported field or an indented continuation.", lineNo)
		}
		if cur == nil {
			if metricFields[field] {
				return nil, nil, fail("SS_ITEM_MALFORMED", "Item field appears outside an item.", lineNo)
			}
			continue
		}
		if cur.seen[field] {
			return nil, nil, fail("SS_ITEM_MALFORMED", "Item repeats a recognized field.", lineNo)
		}
		cur.seen[field] = true
		switch field {
		case "Description":
			cur.Description = value
		case "Outcome":
			cur.Outcome = value
		case "Scope":
			cur.Scope = value
		case "Depends on":
			cur.DepsLine = lineNo
			ids, bad := parseIDList(value)
			cur.Deps, cur.DepsBad = ids, bad
		case "Spec dir":
			cur.HasSpec = true
			cur.SpecLine = lineNo
			cur.SpecDir = unwrapPath(trimSentenceDot(value))
		case "superseded-by":
			cur.HasSuper = true
			cur.SuperLine = lineNo
			cur.SuperRaw = value
			id, bad := parseSingleID(value)
			cur.SuperID, cur.SuperBad = id, bad
		}
	}
	if inFence {
		return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Code fence is not closed.", fenceLine)
	}
	if inComment {
		return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "HTML comment is not closed.", commentLine)
	}
	flush()
	if len(items) > maxItems {
		return nil, nil, fail("SS_INPUT_LIMIT", "Roadmap item limit exceeded.", 0)
	}
	edges := 0
	for _, it := range items {
		edges += len(it.Deps)
		if it.HasSuper {
			edges++
		}
	}
	if edges > maxEdges {
		return nil, nil, fail("SS_INPUT_LIMIT", "Reference edge limit exceeded.", 0)
	}
	return items, diags, nil
}

func matchItem(line string) (id, title, status string, ok bool) {
	const prefix = "### "
	if !strings.HasPrefix(line, prefix) {
		return "", "", "", false
	}
	rest := strings.TrimRight(line, " \t")
	// ### NNN — Title [status: token]
	if len(rest) < len(prefix)+3 {
		return "", "", "", false
	}
	body := rest[len(prefix):]
	if len(body) < 6 || !isDigit3(body[:3]) || !strings.HasPrefix(body[3:], " — ") {
		return "", "", "", false
	}
	id = body[:3]
	after := body[3+len(" — "):]
	mark := strings.LastIndex(after, " [status: ")
	if mark <= 0 || !strings.HasSuffix(after, "]") {
		return "", "", "", false
	}
	title = strings.TrimSpace(after[:mark])
	status = strings.TrimSpace(after[mark+len(" [status: ") : len(after)-1])
	if title == "" || status == "" || strings.ContainsAny(title, "\r\n") {
		return "", "", "", false
	}
	return id, title, status, true
}

// fieldLabels maps each accepted label to the field it fills. "Scope (in)" is
// the speckit-roadmap spelling of Scope; using both in one item is a repeat.
var fieldLabels = []struct{ label, field string }{
	{"Description", "Description"},
	{"Outcome", "Outcome"},
	{"Scope (in)", "Scope"},
	{"Scope", "Scope"},
	{"Depends on", "Depends on"},
	{"Spec dir", "Spec dir"},
	{"superseded-by", "superseded-by"},
}

// matchField recognizes "- Label: value". The speckit-roadmap bold spellings
// "- **Label:** value" and "- **Label**: value" are read as the same field.
func matchField(line string) (name, value string, ok bool) {
	if !strings.HasPrefix(line, "- ") {
		return "", "", false
	}
	body := line[2:]
	for _, f := range fieldLabels {
		for _, p := range []string{f.label + ":", "**" + f.label + ":**", "**" + f.label + "**:"} {
			if strings.HasPrefix(body, p) {
				return f.field, strings.TrimSpace(body[len(p):]), true
			}
		}
	}
	return "", "", false
}

// trimSentenceDot removes one trailing period, as in "Depends on: 002."
func trimSentenceDot(value string) string {
	value = strings.TrimSpace(value)
	return strings.TrimSpace(strings.TrimSuffix(value, "."))
}

// stripNote removes one trailing parenthetical note, as in "002 (verified)".
func stripNote(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, ")") {
		if open := strings.LastIndex(s, " ("); open > 0 {
			return strings.TrimSpace(s[:open])
		}
	}
	return s
}

// parseIDList reads "001, 002". It accepts "none", one trailing period, and a
// parenthetical note after each ID. Any other element marks the list malformed.
func parseIDList(value string) ([]string, bool) {
	value = trimSentenceDot(value)
	if value == "" {
		return nil, true
	}
	if strings.EqualFold(value, "none") {
		return nil, false
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		id := stripNote(p)
		if !isDigit3(id) {
			return nil, true
		}
		out = append(out, id)
	}
	return out, false
}

func parseSingleID(value string) (string, bool) {
	id := trimSentenceDot(value)
	if !isDigit3(id) {
		return "", true
	}
	return id, false
}

func unwrapPath(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 && value[0] == '`' && value[len(value)-1] == '`' && strings.Count(value, "`") == 2 {
		return strings.TrimSpace(value[1 : len(value)-1])
	}
	return value
}

func isDigit3(s string) bool {
	return len(s) == 3 && s[0] >= '0' && s[0] <= '9' && s[1] >= '0' && s[1] <= '9' && s[2] >= '0' && s[2] <= '9'
}

func indentDepth(line string) int {
	n := 0
	for _, r := range line {
		if r == ' ' {
			n++
			continue
		}
		if r == '\t' {
			n += 2
			continue
		}
		break
	}
	return n / 2
}

type taskFile struct {
	Checked int
	Total   int
	Links   []taskLink
	Checks  []TaskCheck
}

type taskLink struct {
	Path    string
	Checked bool
	Line    int
	Kind    string // local, unsupported
}

func parseTasks(rel string, data []byte) (taskFile, *Diagnostic) {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	if !utf8.Valid(data) {
		d := Diagnostic{Code: "SS_ENCODING_INVALID", Severity: "error", Message: "Tasks file encoding is not valid UTF-8.", Path: rel}
		return taskFile{}, &d
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	var tf taskFile
	inFence := false
	var fence byte
	var fenceN, fenceLine, fenceIndent, commentIndent int
	inComment, commentLine := false, 0
	allLines := strings.Split(text, "\n")
	refs := referenceDefinitions(allLines)
	inTask, taskChecked := false, false
	// inList tracks whether the current line can belong to a list item. A
	// line indented four or more columns outside a list is an indented code
	// block, so it is neither a task, a fence, nor a comment.
	inList, prevBlank := false, true
	for i, line := range allLines {
		if inFence {
			if fenceIndent > 0 && strings.TrimSpace(line) != "" {
				if w, _ := leadingIndent(line); w < fenceIndent {
					d := Diagnostic{Code: "SS_LAYOUT_UNSUPPORTED", Severity: "error", Message: "Code fence inside a task list item is not closed before the item ends.", Path: rel, Line: fenceLine}
					return taskFile{}, &d
				}
			}
			if closesFence(line, fence, fenceN) {
				if w, _ := leadingIndent(line); w != fenceIndent {
					d := Diagnostic{Code: "SS_LAYOUT_UNSUPPORTED", Severity: "error", Message: "Code fence closer in tasks file is indented differently from its opener.", Path: rel, Line: i + 1}
					return taskFile{}, &d
				}
				inFence = false
			}
			continue
		}
		if inComment {
			if commentIndent > 0 && strings.TrimSpace(line) != "" {
				if w, _ := leadingIndent(line); w < commentIndent {
					d := Diagnostic{Code: "SS_LAYOUT_UNSUPPORTED", Severity: "error", Message: "HTML comment inside a task list item is not closed before the item ends.", Path: rel, Line: commentLine}
					return taskFile{}, &d
				}
			}
			if strings.Contains(line, "-->") {
				inComment = false
			}
			continue
		}
		blank := strings.TrimSpace(line) == ""
		width, b := leadingIndent(line)
		_, _, _, isHeading := heading(line)
		switch {
		case blank:
		case isHeading:
			inList = false
		case listMarker(line) && (width <= 3 || inList):
			inList = true
		case width <= 3 && prevBlank && !strings.HasPrefix(line[b:], ">"):
			inList = false // a paragraph after a blank line ends the list
		}
		prevBlank = blank
		if !blank && width >= 4 && !inList {
			inTask = false
			continue // indented code block: not tasks, fences, or comments
		}
		if n, ch, ok := openFenceIndent(line, 1<<30); ok {
			inFence, fence, fenceN, fenceLine = true, ch, n, i+1
			fenceIndent = width
			continue
		}
		if strings.HasPrefix(line[b:], "<!--") {
			if !strings.Contains(line[b+4:], "-->") {
				inComment, commentLine, commentIndent = true, i+1, width
			}
			continue
		}
		checked, body, ok := matchTask(line)
		if !ok {
			// A wrapped continuation of a task can carry a link. It is not
			// checked, so it is reported rather than ignored.
			if strings.TrimSpace(line) == "" {
				inTask = false
			} else if _, _, _, isHeading := heading(line); isHeading {
				inTask = false
			} else if inTask {
				for _, link := range extractLinks(line, refs) {
					tf.Links = append(tf.Links, taskLink{Kind: "unsupported", Checked: taskChecked, Line: i + 1, Path: link.Path})
				}
			}
			continue
		}
		inTask, taskChecked = true, checked
		if indentDepth(line) > maxNesting {
			d := Diagnostic{Code: "SS_INPUT_LIMIT", Severity: "error", Message: "Task nesting limit exceeded.", Path: rel, Line: i + 1}
			return taskFile{}, &d
		}
		tf.Total++
		tf.Checks = append(tf.Checks, TaskCheck{Line: i + 1, Checked: checked})
		if tf.Total > maxTasks {
			d := Diagnostic{Code: "SS_INPUT_LIMIT", Severity: "error", Message: "Task checkbox limit exceeded.", Path: rel, Line: i + 1}
			return taskFile{}, &d
		}
		if checked {
			tf.Checked++
		}
		for _, link := range extractLinks(body, refs) {
			link.Checked = checked
			link.Line = i + 1
			tf.Links = append(tf.Links, link)
		}
	}
	if inFence {
		d := Diagnostic{Code: "SS_LAYOUT_UNSUPPORTED", Severity: "error", Message: "Code fence in tasks file is not closed.", Path: rel, Line: fenceLine}
		return taskFile{}, &d
	}
	if inComment {
		d := Diagnostic{Code: "SS_LAYOUT_UNSUPPORTED", Severity: "error", Message: "HTML comment in tasks file is not closed.", Path: rel, Line: commentLine}
		return taskFile{}, &d
	}
	return tf, nil
}

// matchTask recognizes "- [ ] text" with a -, *, or + marker, or an ordered
// marker such as "1. [x] text" or "2) [ ] text".
// Like GFM, it allows block-quote markers and one or more spaces or tabs
// after the list marker.
func matchTask(line string) (checked bool, body string, ok bool) {
	trim := strings.TrimLeft(line, " \t")
	for strings.HasPrefix(trim, ">") {
		trim = strings.TrimLeft(trim[1:], " \t")
	}
	var after string
	switch {
	case len(trim) > 1 && (trim[0] == '-' || trim[0] == '*' || trim[0] == '+'):
		after = trim[1:]
	default:
		j := 0
		for j < len(trim) && j < 9 && trim[j] >= '0' && trim[j] <= '9' {
			j++
		}
		if j == 0 || j+1 >= len(trim) || (trim[j] != '.' && trim[j] != ')') {
			return false, "", false
		}
		after = trim[j+1:]
	}
	if after == "" || (after[0] != ' ' && after[0] != '\t') {
		return false, "", false
	}
	rest := strings.TrimLeft(after, " \t")
	if len(rest) < 3 || rest[0] != '[' || rest[2] != ']' {
		return false, "", false
	}
	mark := rest[1]
	if mark != ' ' && mark != 'x' && mark != 'X' {
		return false, "", false
	}
	if len(rest) < 4 || (rest[3] != ' ' && rest[3] != '\t') {
		return false, "", false
	}
	return mark == 'x' || mark == 'X', strings.TrimSpace(rest[4:]), true
}

// extractLinks finds inline links "[label](target)" in a task. A target may
// be wrapped in <...> and may carry a quoted title. Reference-style links
// "[label][ref]" and autolinks "<path>" are reported as unsupported rather
// than skipped, so they are never silently ignored.
func extractLinks(body string, refs map[string]bool) (out []taskLink) {
	inline := 0
	defer func() {
		// Every "](" should belong to an inline link found above. Any left
		// over (nested brackets, for example) is a link we could not parse.
		for extra := strings.Count(body, "](") - inline; extra > 0; extra-- {
			out = append(out, taskLink{Kind: "unsupported"})
		}
	}()
	for i := 0; i < len(body); i++ {
		if body[i] == '<' {
			if end := strings.IndexByte(body[i+1:], '>'); end > 0 {
				inner := body[i+1 : i+1+end]
				if !strings.ContainsAny(inner, " \t") && strings.ContainsAny(inner, "./:") {
					out = append(out, taskLink{Kind: "unsupported"})
					i += end + 1
				}
			}
			continue
		}
		if body[i] != '[' {
			continue
		}
		image := i > 0 && body[i-1] == '!'
		endLabel := strings.IndexByte(body[i+1:], ']')
		if endLabel < 0 {
			continue
		}
		endLabel += i + 1
		if endLabel+1 < len(body) && body[endLabel+1] == '[' {
			if endRef := strings.IndexByte(body[endLabel+2:], ']'); endRef >= 0 {
				out = append(out, taskLink{Kind: "unsupported"})
				i = endLabel + 2 + endRef
			}
			continue
		}
		if endLabel+1 >= len(body) || body[endLabel+1] != '(' {
			// A shortcut reference "[label]" whose definition exists.
			if refs[strings.ToLower(strings.TrimSpace(body[i+1:endLabel]))] {
				out = append(out, taskLink{Kind: "unsupported"})
				i = endLabel
			}
			continue
		}
		end := strings.IndexByte(body[endLabel+2:], ')')
		if end < 0 {
			continue
		}
		end += endLabel + 2
		target, okTarget := linkDestination(body[endLabel+2 : end])
		i = end
		inline++
		if image || !okTarget || unsupportedTarget(target) {
			out = append(out, taskLink{Kind: "unsupported", Path: ""})
			continue
		}
		out = append(out, taskLink{Kind: "local", Path: target})
	}
	return out
}

// linkDestination splits "path" or "path \"title\"" or "<path> 'title'" into
// the path. Anything else after the path makes the link unsupported.
func linkDestination(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "<") {
		end := strings.IndexByte(raw, '>')
		if end < 0 {
			return "", false
		}
		dest, rest := raw[1:end], strings.TrimSpace(raw[end+1:])
		return dest, rest == "" || quotedTitle(rest)
	}
	cut := strings.IndexAny(raw, " \t")
	if cut < 0 {
		return raw, true
	}
	return raw[:cut], quotedTitle(strings.TrimSpace(raw[cut:]))
}

func quotedTitle(s string) bool {
	if len(s) < 2 {
		return false
	}
	q := s[0]
	return (q == '"' || q == '\'') && s[len(s)-1] == q
}

func unsupportedTarget(target string) bool {
	t := strings.TrimSpace(target)
	if t == "" || strings.HasPrefix(t, "#") {
		return true
	}
	lower := strings.ToLower(t)
	if strings.Contains(lower, "://") || strings.HasPrefix(lower, "mailto:") {
		return true
	}
	if strings.ContainsAny(t, "?#") {
		return true
	}
	// An absolute path or a scheme such as file: or javascript: is not a
	// root-relative local reference.
	if strings.HasPrefix(t, "/") || strings.HasPrefix(t, "\\") {
		return true
	}
	if colon := strings.IndexByte(t, ':'); colon >= 0 && !strings.Contains(t[:colon], "/") {
		return true
	}
	return false
}

// referenceDefinitions collects "[label]: target" definitions, lowercased.
func referenceDefinitions(lines []string) map[string]bool {
	refs := map[string]bool{}
	for _, line := range lines {
		width, b := leadingIndent(line)
		rest := line[b:]
		if width > 3 || !strings.HasPrefix(rest, "[") {
			continue
		}
		end := strings.Index(rest, "]:")
		if end <= 1 {
			continue
		}
		if strings.TrimSpace(rest[end+2:]) != "" {
			refs[strings.ToLower(strings.TrimSpace(rest[1:end]))] = true
		}
	}
	return refs
}

// listMarker reports a line that starts a list item: "-", "*", or "+" and a
// space, or "1." / "1)" and a space, after indentation and quote markers.
func listMarker(line string) bool {
	trim := strings.TrimLeft(line, " \t")
	for strings.HasPrefix(trim, ">") {
		trim = strings.TrimLeft(trim[1:], " \t")
	}
	if len(trim) > 1 && (trim[0] == '-' || trim[0] == '*' || trim[0] == '+') && (trim[1] == ' ' || trim[1] == '\t') {
		return true
	}
	j := 0
	for j < len(trim) && j < 9 && trim[j] >= '0' && trim[j] <= '9' {
		j++
	}
	return j > 0 && j+1 < len(trim) && (trim[j] == '.' || trim[j] == ')') && (trim[j+1] == ' ' || trim[j+1] == '\t')
}
