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
	inComment, inTable, inHTML := false, false, false
	prev := ""
	// pending holds the last field line until its indented continuation
	// lines, if any, have been joined to it, the way renderers show a
	// wrapped value. It is applied before any other line is processed.
	var pending *pendingField
	// blockEnd is true on the line right after a fence or comment closed.
	// Like a blank line or heading, that lets a type-7 HTML block start: in
	// CommonMark only a paragraph prevents it.
	blockEnd := false
	metricOpen, metricBlank := false, false
	// tagRun is true from a line outside items that starts with an HTML tag
	// up to the next blank line. Renderers may show that run as a raw HTML
	// block, so an item heading or field inside it is fatal even when the
	// tag does not start a block by the stricter checks below (T026 round 3
	// B1: after a setext heading, after indented code, or with ">" inside
	// an attribute).
	tagRun := false
	// hidden marks lines inside fences and comments; tails holds the visible
	// text after a comment closes on a line; metricLine marks lines accepted
	// as a Depends on, Spec dir, or superseded-by field. They feed the
	// line-level label check after the loop.
	hidden := make([]bool, len(lines)+1)
	tails := map[int]string{}
	metricLine := map[int]bool{}
	apply := func() {
		if pending != nil && cur != nil {
			applyField(cur, pending.field, pending.value, pending.line)
		}
		pending = nil
	}
	for i, line := range lines {
		lineNo := i + 1
		last := prev
		prev = line
		prevBlockEnd := blockEnd
		blockEnd = false
		if inFence {
			hidden[lineNo] = true
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
				inFence, blockEnd = false, true
			}
			continue
		}
		if inComment {
			hidden[lineNo] = true
			if commentIndent > 0 && strings.TrimSpace(line) != "" {
				if w, _ := leadingIndent(line); w < commentIndent {
					return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "HTML comment inside a list item is not closed before the item ends.", commentLine)
				}
			}
			if end := strings.Index(line, "-->"); end >= 0 {
				if cur != nil && strings.TrimSpace(line[end+3:]) != "" {
					// Text after the closing --> renders as visible text; it
					// is not a field (review 9 D3). Fail visibly.
					return nil, nil, fail("SS_ITEM_MALFORMED", "Text follows the end of an HTML comment on the same line under an item.", lineNo)
				}
				tails[lineNo] = line[end+3:]
				inComment, blockEnd = false, true
			}
			continue
		}
		if strings.TrimSpace(line) == "" {
			tagRun = false
		} else if w, b := leadingIndent(line); cur == nil && w <= 3 && htmlTagLine(line[b:]) {
			tagRun = true
		}
		// D3 (review 8): a Depends on, Spec dir, or superseded-by value
		// continued on an indented line after a blank line renders as part
		// of the same bullet, but is not joined to the value: fail visibly.
		if cur != nil && metricOpen {
			switch w, b := leadingIndent(line); {
			case strings.TrimSpace(line) == "":
				metricBlank = true
			case w == 0:
				metricOpen, metricBlank = false, false
			case listMarker(line) || strings.HasPrefix(line[b:], ">"):
				// A list or quote nested under the field renders as part of
				// the value, with or without a blank line, comment, or other
				// block before it. It is never joined to the value, so fail
				// visibly (T026 review B1).
				return nil, nil, fail("SS_ITEM_MALFORMED", "A list or quote nested under a Depends on, Spec dir, or superseded-by field is not supported.", lineNo)
			case metricBlank:
				return nil, nil, fail("SS_ITEM_MALFORMED", "A Depends on, Spec dir, or superseded-by value continued after a blank line or another block is not supported.", lineNo)
			case startsBlock(line, b):
				// Closed field body (owner decision D-2026-11-08-04): under a
				// metric field, an indented line either joins the value as
				// paragraph text or is fatal. A comment, HTML block, fence,
				// break, or heading renders inside the field's bullet but is
				// never joined, at any indentation (final T026 review B-1 to B-3).
				return nil, nil, fail("SS_ITEM_MALFORMED", "Under a Depends on, Spec dir, or superseded-by field, only indented continuation text is supported.", lineNo)
			}
		}
		if cur != nil && bareMetricLabel(line) {
			// "superseded-by: 002" or "Depends on: 001" without a list marker
			// renders as text, but its author meant a field: fail visibly.
			return nil, nil, fail("SS_ITEM_MALFORMED", "A Depends on, Spec dir, or superseded-by label without a list marker is not supported.", lineNo)
		}
		if pending != nil {
			w, b := leadingIndent(line)
			if w > 0 && strings.TrimSpace(line) != "" {
				if listMarker(line) {
					if metricFields[pending.field] {
						return nil, nil, fail("SS_ITEM_MALFORMED", "A list under a Depends on, Spec dir, or superseded-by field is not supported.", lineNo)
					}
				} else if !startsBlock(line, b) {
					// Only paragraph text continues a value. A line that starts
					// a block (heading, break, HTML, fence, comment, quote, or
					// an item-like line) is handled on its own below.
					pending.value = strings.TrimSpace(pending.value + " " + strings.TrimSpace(line))
					continue
				}
			}
			apply()
		}
		if inHTML {
			if strings.TrimSpace(line) == "" {
				inHTML = false
				continue
			}
			// A raw HTML block runs to the next blank line; renderers show its
			// lines as raw HTML, not as roadmap structure.
			if _, _, _, ok := matchItem(line); ok || itemLikeLine(line) {
				return nil, nil, fail("SS_ITEM_MALFORMED", "Roadmap item inside a raw HTML block; renderers show it as text.", lineNo)
			}
			if _, _, ok := matchField(line); ok {
				return nil, nil, fail("SS_ITEM_MALFORMED", "Item field inside a raw HTML block; renderers show it as text.", lineNo)
			}
			if _, ok := fieldPrefix(line); ok {
				return nil, nil, fail("SS_ITEM_MALFORMED", "Item field inside a raw HTML block; renderers show it as text.", lineNo)
			}
			continue
		}
		if n, ch, ok := openFence(line); ok {
			inFence, fence, fenceN, fenceLine = true, ch, n, lineNo
			fenceIndent, _ = leadingIndent(line)
			hidden[lineNo] = true
			continue
		}
		// HTML comments are not rendered, so like fences they are not live
		// roadmap data. A comment opens only at three or fewer columns of
		// indentation; deeper "<!--" is indented code or continuation text.
		// An unclosed comment is fatal below.
		if w, b := leadingIndent(line); w <= 3 && (strings.HasPrefix(line[b:], "<!-->") || strings.HasPrefix(line[b:], "<!--->")) {
			return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Empty HTML comment is read differently by renderers.", lineNo)
		}
		if w, b := leadingIndent(line); w <= 3 && rawHTMLBlock(line[b:]) {
			return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Raw HTML block (pre, script, style, textarea) is not supported.", lineNo)
		}
		if w, b := leadingIndent(line); w <= 3 && otherHTMLBlock(line[b:]) {
			return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "HTML processing instruction, CDATA, or declaration block is not supported.", lineNo)
		}
		if w, b := leadingIndent(line); ((cur == nil && w <= 3) || (cur != nil && w >= 1 && w <= 5)) && htmlBlockStart(line[b:], notParagraph(last) || prevBlockEnd) {
			// Under an item, an indented HTML block (for example <details>)
			// hides its content from renderers, so a field or item heading
			// inside it is fatal (review 8 D2). A column-0 HTML line under an
			// item stays fatal under the closed item body.
			inHTML = true
			continue
		}
		if w, b := leadingIndent(line); w <= 3 && strings.HasPrefix(line[b:], "<!--") {
			trim := strings.TrimSpace(line)
			hidden[lineNo] = true
			if end := strings.Index(trim[4:], "-->"); end >= 0 {
				blockEnd = true
				tails[lineNo] = trim[4+end+3:]
				if cur != nil && strings.TrimSpace(trim[4+end+3:]) != "" {
					return nil, nil, fail("SS_ITEM_MALFORMED", "Text follows an HTML comment on the same line under an item.", lineNo)
				}
			} else {
				inComment, commentLine, commentIndent = true, lineNo, w
			}
			continue
		}
		if inTable {
			_, _, _, isHeading := heading(line)
			if strings.TrimSpace(line) != "" && strings.Contains(line, "|") && !listMarker(line) && !isHeading {
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
		if tagRun {
			_, _, isField := matchField(line)
			_, isNear := fieldPrefix(line)
			if _, _, _, ok := matchItem(line); ok || itemLikeLine(line) || isField || isNear {
				return nil, nil, fail("SS_ITEM_MALFORMED", "Roadmap item or field directly after an HTML tag line; renderers may show it as raw HTML.", lineNo)
			}
		}
		if id, title, status, ok := matchItem(line); ok {
			flush()
			cur = &rawItem{ID: id, Title: title, Status: status, Line: lineNo, seen: map[string]bool{}}
			continue
		}
		if itemLikeLine(line) {
			return nil, nil, fail("SS_ITEM_MALFORMED", "Item heading does not match the supported shape.", lineNo)
		}
		if _, _, text, ok := heading(line); ok && bareMetricLabel(text) {
			return nil, nil, fail("SS_ITEM_MALFORMED", "A Depends on, Spec dir, or superseded-by label written as a heading is not supported.", lineNo)
		}
		if _, indent, _, ok := heading(line); ok && cur != nil && indent > 0 {
			// A heading indented under an item sits inside the item's list
			// content for renderers; reading it as the item's end could drop
			// what follows (review 9 D4).
			return nil, nil, fail("SS_ITEM_MALFORMED", "An indented heading under an item is not supported.", lineNo)
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
			case thematicBreak(line):
				// An indented break (one to three spaces) also ends the item.
				flush()
				continue
			case strings.TrimSpace(line) == "", width > 0, descriptiveField(line):
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
		pending = &pendingField{field: field, value: value, line: lineNo}
		metricOpen, metricBlank = metricFields[field], false
		metricLine[lineNo] = metricFields[field]
	}
	apply()
	if inFence {
		return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Code fence is not closed.", fenceLine)
	}
	if inComment {
		return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "HTML comment is not closed.", commentLine)
	}
	// Line-level label rule (owner direction, 8 November 2026): a visible
	// line that contains a Depends on, Spec dir, or superseded-by label with
	// a colon must be an accepted field line under an item. This holds by
	// construction, whatever Markdown layout hides or splits the field.
	for i, line := range lines {
		lineNo := i + 1
		if metricLine[lineNo] {
			continue
		}
		visible := line
		if hidden[lineNo] {
			visible = tails[lineNo]
		}
		if metricLabelIn(visible) {
			return nil, nil, fail("SS_ITEM_MALFORMED", "A Depends on, Spec dir, or superseded-by label appears outside a supported field line.", lineNo)
		}
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
	// speckit.roadmap.write strikes through superseded entries:
	// "### ~~006 — Title~~  [status: abandoned]" is read as item 006.
	struck := strings.HasPrefix(body, "~~")
	if struck {
		body = body[2:]
	}
	if len(body) < 6 || !isDigit3(body[:3]) || !strings.HasPrefix(body[3:], " — ") {
		return "", "", "", false
	}
	id = body[:3]
	after := body[3+len(" — "):]
	if strings.Count(after, "[status:") > 1 {
		return "", "", "", false // two status markers: ambiguous, fatal upstream
	}
	mark := strings.LastIndex(after, " [status: ")
	if mark <= 0 || !strings.HasSuffix(after, "]") {
		return "", "", "", false
	}
	title = strings.TrimSpace(after[:mark])
	if struck {
		if !strings.HasSuffix(title, "~~") || strings.Contains(strings.TrimSuffix(title, "~~"), "~~") {
			return "", "", "", false
		}
		title = strings.TrimSpace(strings.TrimSuffix(title, "~~"))
	}
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
	if strings.ContainsAny(value, "<[\"`\\&") {
		// Markup in a dependency list can hide part of the raw text from
		// readers (a link title, comment, or attribute spanning a comma), so
		// the list is malformed rather than split (T026 round 3 B2).
		return nil, true
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
	// Narrow tasks dialect (owner decision, 7 November 2026). items holds the
	// open list items; indentation is measured from the innermost item's
	// content column. Where the answer would depend on CommonMark container
	// rules, the file is fatal instead of guessed.
	var items []listItem
	prevBlank, inHTML := true, false
	// lastPara is true when the previous non-blank line was paragraph text
	// (including a list item's first line). Only a paragraph allows lazy
	// continuation and stops a type-7 HTML block from starting.
	lastPara := false
	fatalAt := func(msg string, line int) (taskFile, *Diagnostic) {
		d := Diagnostic{Code: "SS_LAYOUT_UNSUPPORTED", Severity: "error", Message: msg, Path: rel, Line: line}
		return taskFile{}, &d
	}
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
				inFence, lastPara = false, false
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
				inComment, lastPara = false, false
			}
			continue
		}
		blank := strings.TrimSpace(line) == ""
		width, b := leadingIndent(line)
		wasBlank := prevBlank
		prevBlank = blank
		if blank {
			inTask, inHTML = false, false
			continue
		}
		prevPara := lastPara
		lastPara = true
		if inHTML {
			lastPara = false
			// A raw HTML block runs to the next blank line; renderers show a
			// checkbox in it as raw text.
			if _, _, ok := matchTask(line); ok {
				return fatalAt("Checkbox inside a raw HTML block; renderers show it as text.", i+1)
			}
			continue
		}
		if _, _, _, isHeading := heading(line); isHeading && width <= 3 {
			items, inTask, lastPara = nil, false, false
			continue
		}
		closesItems := strings.HasPrefix(line[b:], ">") || thematicBreak(line) || strings.HasPrefix(line[b:], "<") || func() bool { _, _, ok := openFence(line); return ok }()
		if closesItems && !listMarker(line) {
			lastPara = false
			// A quote, break, fence, comment, or HTML line less indented than
			// an item's content ends that item; lazy continuation applies to
			// paragraph text only.
			for len(items) > 0 && width < items[len(items)-1].content {
				items = items[:len(items)-1]
			}
		}
		if strings.HasPrefix(line[b:], ">") {
			// Block quotes: a fence, comment, checkbox, or indented code inside
			// a quote depends on container rules, so it is fatal.
			inner := line[b:]
			for strings.HasPrefix(inner, ">") {
				inner = strings.TrimPrefix(strings.TrimPrefix(inner[1:], " "), "\t")
			}
			innerW, innerB := leadingIndent(inner)
			if _, _, ok := openFence(inner); ok || innerW >= 4 || strings.HasPrefix(inner[innerB:], "<!--") {
				return fatalAt("Fence, comment, or indented code inside a block quote in a tasks file is not supported.", i+1)
			}
			if _, _, ok := matchTask(line); ok {
				return fatalAt("Checkbox inside a block quote in a tasks file is not supported.", i+1)
			}
			inTask = false
			continue
		}
		isMarker := listMarker(line)
		sibling := false
		if isMarker {
			ordered, number := orderedMarker(line[b:])
			// A marker closes items at its indentation or deeper, and items
			// whose content column it does not reach.
			for len(items) > 0 && (items[len(items)-1].marker >= width || width < items[len(items)-1].content) {
				if items[len(items)-1].marker == width && items[len(items)-1].ordered == ordered {
					sibling = true
				}
				items = items[:len(items)-1]
			}
			_, _, isTask := matchTask(line)
			if ordered && number != 1 && !wasBlank && prevPara && !sibling {
				// An ordered list not starting at 1 cannot interrupt a
				// paragraph, so this line is paragraph text.
				if isTask {
					return fatalAt("Checkbox in an ordered list that does not start at 1 directly after text; renderers read it as text.", i+1)
				}
				continue
			}
			if spacesAfterMarker(line[b:]) >= 5 && isTask {
				return fatalAt("Checkbox five or more spaces after its list marker is indented code to renderers.", i+1)
			}
			if inner := afterMarker(line[b:]); listMarker(inner) {
				if _, _, ok := matchTask(inner); ok {
					return fatalAt("Nested list marker on the same line as its parent is not supported.", i+1)
				}
			}
		} else if wasBlank || !prevPara {
			// A line after a blank line, or after a block that is not a
			// paragraph (a closed fence or comment), is not a lazy
			// continuation: it closes the items it is not indented into.
			for len(items) > 0 && width < items[len(items)-1].content {
				items = items[:len(items)-1]
			}
		}
		base := 0
		for j := len(items) - 1; j >= 0; j-- {
			if items[j].content <= width {
				base = items[j].content
				break
			}
		}
		if width-base >= 4 && (isMarker || wasBlank) {
			// Four or more columns past the enclosing content: renderers read
			// this as code or continuation text, not as a list item.
			if _, _, ok := matchTask(line); ok {
				return fatalAt("Checkbox is indented four or more columns past its list item; renderers read it as code or text.", i+1)
			}
			inTask = false
			continue
		}
		if isMarker {
			ordered, _ := orderedMarker(line[b:])
			items = append(items, listItem{marker: width, content: contentColumn(line), ordered: ordered})
		}
		if strings.HasPrefix(line[b:], "<!-->") || strings.HasPrefix(line[b:], "<!--->") {
			return fatalAt("Empty HTML comment is read differently by renderers.", i+1)
		}
		if rawHTMLBlock(line[b:]) {
			return fatalAt("Raw HTML block (pre, script, style, textarea) in a tasks file is not supported.", i+1)
		}
		if otherHTMLBlock(line[b:]) {
			return fatalAt("HTML processing instruction, CDATA, or declaration block in a tasks file is not supported.", i+1)
		}
		if width-base <= 3 && !isMarker && htmlBlockStart(line[b:], wasBlank || !prevPara) {
			inHTML, inTask = true, false
			continue
		}
		if n, ch, ok := openFenceIndent(line, base+3); ok && !isMarker {
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
	if len(rest) == 3 {
		return mark == 'x' || mark == 'X', "", true // "- [ ]" with no text is still a checkbox on GitHub
	}
	if rest[3] != ' ' && rest[3] != '\t' {
		return false, "", false
	}
	return mark == 'x' || mark == 'X', strings.TrimSpace(rest[4:]), true
}

// extractLinks finds inline links "[label](target)" in a task. A target may
// be wrapped in <...> and may carry a quoted title. Reference-style links
// "[label][ref]" and autolinks "<path>" are reported as unsupported rather
// than skipped, so they are never silently ignored.
func extractLinks(body string, refs map[string]bool) (out []taskLink) {
	body = stripCodeSpans(body)
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

// stripCodeSpans removes inline code spans from a task line before links are
// looked for: text inside a code span is literal, never a link (CommonMark
// 6.1). A span opens with a run of N backticks and closes at the next run of
// exactly N backticks. A run with no matching close is literal text and is
// kept, so a stray backtick cannot hide a real link after it. Roadmap sync
// finding F5, decision D-2026-11-12-01.
func stripCodeSpans(line string) string {
	var b strings.Builder
	for i := 0; i < len(line); {
		if line[i] != '`' {
			b.WriteByte(line[i])
			i++
			continue
		}
		n := 0
		for i+n < len(line) && line[i+n] == '`' {
			n++
		}
		closeAt := -1
		for j := i + n; j < len(line); {
			if line[j] != '`' {
				j++
				continue
			}
			m := 0
			for j+m < len(line) && line[j+m] == '`' {
				m++
			}
			if m == n {
				closeAt = j
				break
			}
			j += m
		}
		if closeAt < 0 {
			b.WriteString(line[i : i+n])
			i += n
			continue
		}
		b.WriteByte(' ')
		i = closeAt + n
	}
	return b.String()
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

type listItem struct {
	marker, content int
	ordered         bool
}

// orderedMarker reports whether rest starts with an ordered marker, and its number.
func orderedMarker(rest string) (bool, int) {
	n, j := 0, 0
	for j < len(rest) && j < 9 && rest[j] >= '0' && rest[j] <= '9' {
		n = n*10 + int(rest[j]-'0')
		j++
	}
	return j > 0 && j < len(rest) && (rest[j] == '.' || rest[j] == ')'), n
}

// markerWidth is the width of the list marker at the start of rest.
func markerWidth(rest string) int {
	if ok, _ := orderedMarker(rest); ok {
		j := 0
		for rest[j] >= '0' && rest[j] <= '9' {
			j++
		}
		return j + 1
	}
	return 1
}

// spacesAfterMarker counts the spaces between a list marker and its content.
func spacesAfterMarker(rest string) int {
	m, n := markerWidth(rest), 0
	for m+n < len(rest) && rest[m+n] == ' ' {
		n++
	}
	return n
}

// afterMarker returns the content after a list marker and its spaces.
func afterMarker(rest string) string {
	m := markerWidth(rest)
	if m >= len(rest) {
		return ""
	}
	return strings.TrimLeft(rest[m:], " \t")
}

// contentColumn returns the column where a list item's content starts: the
// marker width plus one to four following spaces (one if more than four).
func contentColumn(line string) int {
	width, b := leadingIndent(line)
	rest := line[b:]
	m := 1
	if rest[0] >= '0' && rest[0] <= '9' {
		for m < len(rest) && rest[m] >= '0' && rest[m] <= '9' {
			m++
		}
		m++ // the "." or ")"
	}
	spaces := 0
	for m+spaces < len(rest) && rest[m+spaces] == ' ' {
		spaces++
	}
	if spaces == 0 || spaces > 4 || m+spaces == len(rest) {
		spaces = 1
	}
	return width + m + spaces
}

// rawHTMLBlock reports a CommonMark type-1 HTML block opener (<pre>,
// <script>, <style>, <textarea>), whose content is hidden or raw until the
// closing tag.
func rawHTMLBlock(rest string) bool {
	lower := strings.ToLower(rest)
	for _, tag := range []string{"<pre", "<script", "<style", "<textarea"} {
		if strings.HasPrefix(lower, tag) {
			after := lower[len(tag):]
			return after == "" || after[0] == ' ' || after[0] == '>' || after[0] == '\t'
		}
	}
	return false
}

type pendingField struct {
	field, value string
	line         int
}

// applyField stores a field value, after any continuation lines were joined.
func applyField(cur *rawItem, field, value string, lineNo int) {
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

// htmlBlockTags are the CommonMark type-6 block tags.
var htmlBlockTags = map[string]bool{}

func init() {
	for _, t := range strings.Fields("address article aside base basefont blockquote body caption center col colgroup dd details dialog dir div dl dt fieldset figcaption figure footer form frame frameset h1 h2 h3 h4 h5 h6 head header hr html iframe legend li link main menu menuitem nav noframes ol optgroup option p param search section summary table tbody td tfoot th thead title tr track ul") {
		htmlBlockTags[t] = true
	}
}

// htmlBlockStart reports a CommonMark type-6 or type-7 HTML block opener.
// Such a block runs to the next blank line. Type 7 (any single complete tag
// alone on its line) cannot interrupt a paragraph, so it needs a blank line
// or the start of the file before it.
func htmlBlockStart(rest string, afterBlank bool) bool {
	if !strings.HasPrefix(rest, "<") || strings.HasPrefix(rest, "<!--") {
		return false
	}
	t := rest[1:]
	t = strings.TrimPrefix(t, "/")
	n := 0
	for n < len(t) && (t[n] >= 'a' && t[n] <= 'z' || t[n] >= 'A' && t[n] <= 'Z' || n > 0 && (t[n] >= '0' && t[n] <= '9' || t[n] == '-')) {
		n++
	}
	if n == 0 {
		return false
	}
	name, after := strings.ToLower(t[:n]), t[n:]
	if htmlBlockTags[name] && (after == "" || after[0] == ' ' || after[0] == '\t' || after[0] == '>' || strings.HasPrefix(after, "/>")) {
		return true
	}
	if !afterBlank {
		return false
	}
	trim := strings.TrimRight(rest, " \t")
	return strings.HasSuffix(trim, ">") && strings.Count(trim, "<") == 1 && strings.Count(trim, ">") == 1
}

// notParagraph reports a previous line after which a type-7 HTML block may
// start: a blank line, a heading, or a thematic break.
func notParagraph(last string) bool {
	if strings.TrimSpace(last) == "" || thematicBreak(last) {
		return true
	}
	_, _, _, ok := heading(last)
	return ok
}

// otherHTMLBlock reports CommonMark HTML block types 3 to 5 (processing
// instructions, declarations, CDATA), which span blank lines.
func otherHTMLBlock(rest string) bool {
	if strings.HasPrefix(rest, "<?") || strings.HasPrefix(rest, "<![CDATA[") {
		return true
	}
	return len(rest) > 2 && strings.HasPrefix(rest, "<!") && (rest[2] >= 'A' && rest[2] <= 'Z' || rest[2] >= 'a' && rest[2] <= 'z')
}

// bareMetricLabel reports a line without a list marker that starts with a
// Depends on, Spec dir, or superseded-by label and a colon, at any
// indentation and in any lookalike spelling.
func bareMetricLabel(line string) bool {
	rest := strings.TrimLeft(normalizeLookalikes(line), " \t")
	if listMarker(rest) || strings.HasPrefix(rest, ">") {
		return false
	}
	rest = strings.ToLower(strings.NewReplacer("\\:", ":", "*", "", "`", "", "_", " ", "-", " ").Replace(rest))
	rest = strings.Join(strings.Fields(rest), " ")
	for _, label := range []string{"depends on", "spec dir", "superseded by", "dependson", "specdir", "supersededby"} {
		if strings.HasPrefix(rest, label) {
			if after := strings.TrimLeft(rest[len(label):], " "); strings.HasPrefix(after, ":") {
				return true
			}
		}
	}
	return false
}

// htmlTagLine reports text starting with an HTML open or close tag name,
// such as "<span", "</div", or "<img", whether or not the tag is complete.
func htmlTagLine(rest string) bool {
	if len(rest) < 2 || rest[0] != '<' {
		return false
	}
	i := 1
	if rest[i] == '/' {
		i++
	}
	return i < len(rest) && (rest[i] >= 'A' && rest[i] <= 'Z' || rest[i] >= 'a' && rest[i] <= 'z')
}

// metricLabelIn reports visible text containing a Depends on, Spec dir, or
// superseded-by label followed by a colon, after normalizing lookalikes,
// case, emphasis, hyphens, and spacing. Code spans and inline comments are
// not visible labels and are removed first.
func metricLabelIn(text string) bool {
	text = normalizeLookalikes(text)
	for {
		start := strings.Index(text, "<!--")
		if start < 0 {
			break
		}
		end := strings.Index(text[start+4:], "-->")
		if end < 0 {
			break
		}
		text = text[:start] + " " + text[start+4+end+3:]
	}
	text = strings.ToLower(withoutCodeSpans(text))
	text = strings.NewReplacer("\\:", ":", "*", "", "_", " ", "-", " ").Replace(text)
	text = strings.Join(strings.Fields(text), " ")
	for _, label := range []string{"depends on", "spec dir", "superseded by", "dependson", "specdir", "supersededby"} {
		for off := 0; ; {
			k := strings.Index(text[off:], label)
			if k < 0 {
				break
			}
			k += off
			off = k + 1
			if k > 0 {
				if c := text[k-1]; c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
					continue
				}
			}
			if strings.HasPrefix(strings.TrimLeft(text[k+len(label):], " "), ":") {
				return true
			}
		}
	}
	return false
}

// startsBlock reports a line that begins a new block rather than continuing
// a field's paragraph: a heading, a thematic break, an HTML block, a fence,
// a comment, a block quote, or an item-like line. rest starts at byte b.
func startsBlock(line string, b int) bool {
	rest := line[b:]
	if _, _, _, ok := heading(line); ok || thematicBreak(line) || itemLikeLine(line) {
		return true
	}
	if _, _, ok := openFence(line); ok {
		return true
	}
	return strings.HasPrefix(rest, "<!--") || strings.HasPrefix(rest, ">") ||
		rawHTMLBlock(rest) || otherHTMLBlock(rest) || htmlBlockStart(rest, false)
}
