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
	for i, line := range lines {
		lineNo := i + 1
		if n, ch, ok := fenceMark(line); ok {
			if !inFence {
				inFence, fence, fenceN = true, ch, n
				continue
			}
			if ch == fence && n >= fenceN && strings.TrimSpace(line) == strings.Repeat(string(ch), n) {
				inFence = false
				continue
			}
		}
		if inFence {
			continue
		}
		if strings.Contains(line, "[status:") && strings.Contains(line, "|") {
			return nil, nil, fail("SS_LAYOUT_UNSUPPORTED", "Unsupported item-like table layout.", lineNo)
		}
		if id, title, status, ok := matchItem(line); ok {
			flush()
			cur = &rawItem{ID: id, Title: title, Status: status, Line: lineNo, seen: map[string]bool{}}
			continue
		}
		if itemLikeHeading(line) {
			return nil, nil, fail("SS_ITEM_MALFORMED", "Item heading does not match the supported shape.", lineNo)
		}
		if cur == nil {
			continue
		}
		if depth := indentDepth(line); depth > maxNesting {
			return nil, nil, fail("SS_INPUT_LIMIT", "Parser nesting limit exceeded.", lineNo)
		}
		field, value, ok := matchField(line)
		if !ok {
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
			cur.SpecDir = unwrapPath(value)
		case "superseded-by":
			cur.HasSuper = true
			cur.SuperLine = lineNo
			cur.SuperRaw = value
			id, bad := parseSingleID(value)
			cur.SuperID, cur.SuperBad = id, bad
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

func fenceMark(line string) (n int, ch byte, ok bool) {
	trim := strings.TrimSpace(line)
	if len(trim) < 3 || (trim[0] != '`' && trim[0] != '~') {
		return 0, 0, false
	}
	ch = trim[0]
	i := 0
	for i < len(trim) && trim[i] == ch {
		i++
	}
	if i < 3 {
		return 0, 0, false
	}
	return i, ch, true
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

func itemLikeHeading(line string) bool {
	if !strings.HasPrefix(line, "### ") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(line, "### "))
	if rest == "" {
		return false
	}
	if strings.Contains(line, "[status:") {
		return true
	}
	r, _ := utf8.DecodeRuneInString(rest)
	return r >= '0' && r <= '9'
}

func matchField(line string) (name, value string, ok bool) {
	if !strings.HasPrefix(line, "- ") {
		return "", "", false
	}
	body := line[2:]
	for _, field := range []string{"Description", "Outcome", "Scope", "Depends on", "Spec dir", "superseded-by"} {
		p := field + ":"
		if strings.HasPrefix(body, p) {
			return field, strings.TrimSpace(body[len(p):]), true
		}
	}
	return "", "", false
}

func parseIDList(value string) ([]string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, true
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		id := strings.TrimSpace(p)
		if !isDigit3(id) {
			return nil, true
		}
		out = append(out, id)
	}
	return out, false
}

func parseSingleID(value string) (string, bool) {
	id := strings.TrimSpace(value)
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
	var fenceN int
	for i, line := range strings.Split(text, "\n") {
		if n, ch, ok := fenceMark(line); ok {
			if !inFence {
				inFence, fence, fenceN = true, ch, n
				continue
			}
			if ch == fence && n >= fenceN && strings.TrimSpace(line) == strings.Repeat(string(ch), n) {
				inFence = false
			}
			continue
		}
		if inFence {
			continue
		}
		checked, body, ok := matchTask(line)
		if !ok {
			continue
		}
		if indentDepth(line) > maxNesting {
			d := Diagnostic{Code: "SS_INPUT_LIMIT", Severity: "error", Message: "Task nesting limit exceeded.", Path: rel, Line: i + 1}
			return taskFile{}, &d
		}
		tf.Total++
		if tf.Total > maxTasks {
			d := Diagnostic{Code: "SS_INPUT_LIMIT", Severity: "error", Message: "Task checkbox limit exceeded.", Path: rel, Line: i + 1}
			return taskFile{}, &d
		}
		if checked {
			tf.Checked++
		}
		for _, link := range extractLinks(body) {
			link.Checked = checked
			link.Line = i + 1
			tf.Links = append(tf.Links, link)
		}
	}
	return tf, nil
}

func matchTask(line string) (checked bool, body string, ok bool) {
	trim := strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(trim, "- [") && !strings.HasPrefix(trim, "* [") && !strings.HasPrefix(trim, "+ [") {
		return false, "", false
	}
	// bullet marker then space? accept "- [ ]" with the marker and a space before bracket, or "-[" rejected.
	var rest string
	switch {
	case strings.HasPrefix(trim, "- [") || strings.HasPrefix(trim, "* [") || strings.HasPrefix(trim, "+ ["):
		rest = trim[2:]
	default:
		return false, "", false
	}
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

func extractLinks(body string) []taskLink {
	var out []taskLink
	for i := 0; i < len(body); i++ {
		if body[i] != '[' {
			continue
		}
		image := i > 0 && body[i-1] == '!'
		endLabel := strings.IndexByte(body[i+1:], ']')
		if endLabel < 0 {
			continue
		}
		endLabel += i + 1
		if endLabel+1 >= len(body) || body[endLabel+1] != '(' {
			continue
		}
		end := strings.IndexByte(body[endLabel+2:], ')')
		if end < 0 {
			continue
		}
		end += endLabel + 2
		target := body[endLabel+2 : end]
		i = end
		if image || unsupportedTarget(target) {
			out = append(out, taskLink{Kind: "unsupported", Path: ""})
			continue
		}
		out = append(out, taskLink{Kind: "local", Path: target})
	}
	return out
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
	return false
}
