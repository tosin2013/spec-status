package status

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/rivo/uniseg"
)

// View controls human text only. JSON does not use it.
// Color is an explicit choice made by the caller. Width is a character budget.
type View struct {
	Color       bool
	Width       int
	Statuses    []string
	Blocked     bool
	HumanReview bool
	Active      bool
	Historical  bool
	Explain     bool
	// Build is the panel heading label after "Spec Status"; empty means
	// "development build". It never affects JSON (feature 004).
	Build string
}

func normalizeWidth(n int) int {
	if n < 20 {
		return 20
	}
	if n > 240 {
		return 240
	}
	return n
}

// sanitizeTerminal removes terminal controls from repository-supplied text
// before it is placed in human output. JSON keeps the original strings.
func sanitizeTerminal(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			i++
			continue
		}
		if r == 0x1b {
			i += size
			i = skipEscape(s, i)
			continue
		}
		if r < 0x20 || r == 0x7f || (r >= 0x80 && r <= 0x9f) {
			if r == '\n' || r == '\r' || r == '\t' {
				b.WriteByte(' ')
			}
			i += size
			continue
		}
		b.WriteRune(r)
		i += size
	}
	return b.String()
}

func skipEscape(s string, i int) int {
	if i >= len(s) {
		return i
	}
	switch s[i] {
	case '[':
		i++
		for i < len(s) {
			c := s[i]
			i++
			if c >= 0x40 && c <= 0x7e {
				break
			}
		}
	case ']':
		i++
		for i < len(s) {
			if s[i] == 0x07 {
				i++
				break
			}
			if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '\\' {
				i += 2
				break
			}
			i++
		}
	default:
		i++
	}
	return i
}

func displayWidth(s string) int {
	return uniseg.StringWidth(s)
}

func truncate(s string, width int) string {
	s = sanitizeTerminal(s)
	if width <= 0 {
		return ""
	}
	if displayWidth(s) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	var b strings.Builder
	used := 0
	limit := width - 1
	for _, r := range s {
		w := uniseg.StringWidth(string(r))
		if used+w > limit {
			break
		}
		b.WriteRune(r)
		used += w
	}
	b.WriteString("…")
	return b.String()
}

func pad(s string, width int) string {
	gap := width - displayWidth(s)
	if gap <= 0 {
		return s
	}
	return s + strings.Repeat(" ", gap)
}

type painter struct {
	on bool
	r  *lipgloss.Renderer
}

func newPainter(on bool) *painter {
	// The renderer writes nowhere. Profile is explicit, so Lip Gloss does not
	// query the terminal or open /dev/tty. Ascii keeps the plain fallback free of ANSI.
	profile := termenv.Ascii
	if on {
		profile = termenv.ANSI
	}
	r := lipgloss.NewRenderer(io.Discard, termenv.WithProfile(profile))
	r.SetColorProfile(profile)
	r.SetHasDarkBackground(true)
	return &painter{on: on, r: r}
}

func (p *painter) paint(s, kind string) string {
	if !p.on || s == "" {
		return s
	}
	st := p.r.NewStyle()
	switch kind {
	case "heading":
		st = st.Bold(true)
	case "verified":
		st = st.Foreground(lipgloss.Color("2"))
	case "implemented":
		st = st.Foreground(lipgloss.Color("6"))
	case "historical":
		st = st.Faint(true)
	case "warning":
		st = st.Foreground(lipgloss.Color("3"))
	case "error":
		st = st.Foreground(lipgloss.Color("1"))
	case "info":
		st = st.Foreground(lipgloss.Color("4"))
	case "bar":
		st = st.Foreground(lipgloss.Color("2"))
	case "border":
		st = st.Foreground(lipgloss.Color("8"))
	default:
		return s
	}
	return st.Render(s)
}

func renderHuman(rep Report, view View) string {
	if rep.Outcome == "fatal" {
		return fatalText(rep)
	}
	width := normalizeWidth(view.Width)
	p := newPainter(view.Color)
	var b strings.Builder
	if rep.Summary == nil {
		b.WriteString("summary  N/A\n")
		return b.String()
	}
	b.WriteString(summaryPanel(p, rep, width, view.Build))
	shown := filterItems(rep.Items, view)
	if view.filtering() {
		b.WriteString("shown " + itoa(len(shown)) + " of " + itoa(len(rep.Items)) + " items\n")
		if len(shown) == 0 {
			b.WriteString("no item matches this view; totals above are the full roadmap\n")
		}
	}
	b.WriteString("\n")
	if width >= 60 {
		writeTable(&b, p, shown, view, width)
	} else {
		writeStacked(&b, p, shown, view, width)
	}
	b.WriteString("\n")
	if len(rep.Diagnostics) == 0 {
		b.WriteString("diagnostics  none\n")
		return b.String()
	}
	b.WriteString(p.paint("diagnostics", "heading") + "\n")
	for _, d := range rep.Diagnostics {
		sev := pad(d.Severity, 7)
		b.WriteString(p.paint(sev, d.Severity) + " " + d.Code + "\n")
		loc := d.Path
		if d.Line > 0 {
			loc += ":" + itoa(d.Line)
		}
		if loc == "" {
			loc = "-"
		}
		item := d.ItemID
		if item == "" {
			item = "-"
		}
		msg := sanitizeTerminal(loc + " " + item + " " + d.Message)
		for _, line := range wrapText(msg, width-2) {
			b.WriteString("  " + line + "\n")
		}
	}
	return b.String()
}

func wrapText(s string, width int) []string {
	if width < 8 {
		width = 8
	}
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return []string{""}
	}
	var lines []string
	var cur strings.Builder
	used := 0
	flush := func() {
		if cur.Len() > 0 {
			lines = append(lines, cur.String())
			cur.Reset()
			used = 0
		}
	}
	for _, word := range fields {
		word = truncate(word, width)
		w := displayWidth(word)
		if used == 0 {
			cur.WriteString(word)
			used = w
			continue
		}
		if used+1+w > width {
			flush()
			cur.WriteString(word)
			used = w
			continue
		}
		cur.WriteByte(' ')
		cur.WriteString(word)
		used += 1 + w
	}
	flush()
	return lines
}

const summaryLabelWidth = 18

func summaryPanel(p *painter, rep Report, width int, build string) string {
	s := rep.Summary
	padX := 1
	if width < 48 {
		padX = 0
	}
	inner := width - 2 - 2*padX
	if inner < 16 {
		inner = 16
	}
	lines := summaryLines(p, s, rep.Roadmap, inner, build)
	return frame(p, lines, inner, padX)
}

func summaryLines(p *painter, s *Summary, roadmap string, inner int, build string) []string {
	if build == "" {
		build = "development build"
	}
	var lines []string
	lines = append(lines, fitSegments(inner, p.paint("Spec Status", "heading"), "  "+sanitizeTerminal(build)))
	lines = append(lines, metricLines(p, "recorded verified", s.VerifiedActiveItems, s.ActiveItems, s.VerifiedPercent.Text(), "verified", inner)...)
	lines = append(lines, metricLines(p, "checked tasks", s.Tasks.Checked, s.Tasks.Total, s.Tasks.Percent.Text(), "", inner)...)
	lines = append(lines, countLine(p, "active blockers", itoa(s.BlockedCount), "", inner)...)
	reviewValue := itoa(s.HumanReviewCount)
	reviewKind := ""
	if s.HumanReviewCount > 0 {
		reviewValue += "  warning"
		reviewKind = "warning"
	}
	lines = append(lines, countLine(p, "human review", reviewValue, reviewKind, inner)...)
	ctx := "items " + itoa(s.TotalItems) + "  active " + itoa(s.ActiveItems) + "  historical " + itoa(s.TotalItems-s.ActiveItems) + "  superseded " + itoa(s.SupersededItems)
	lines = append(lines, wrapText(ctx, inner)...)
	scope := "task files read " + itoa(s.Tasks.FilesRead) + "  missing " + itoa(s.Tasks.FilesMissing) + "  skipped " + itoa(s.Tasks.FilesSkipped)
	lines = append(lines, wrapText(scope, inner)...)
	lines = append(lines, wrapText("source "+sanitizeTerminal(roadmap), inner)...)
	return lines
}

func metricLines(p *painter, label string, num, den int, percent, kind string, inner int) []string {
	value := itoa(num) + "/" + itoa(den) + " (N/A)"
	bar := ""
	if den > 0 {
		value = itoa(num) + "/" + itoa(den) + " (" + percent + "%)"
		bar = progressMarks(num, den, 12)
	}
	return countLine(p, label, value, kind, inner, bar)
}

func countLine(p *painter, label, value, kind string, inner int, bar ...string) []string {
	mark := ""
	if len(bar) > 0 {
		mark = bar[0]
	}
	plainLabel := pad(label, summaryLabelWidth)
	plain := plainLabel + "  " + value
	if mark != "" {
		plain = plain + "  [" + mark + "]"
	}
	if displayWidth(plain) <= inner && mark != "" {
		head := p.paint(label, kind) + strings.Repeat(" ", summaryLabelWidth-displayWidth(label))
		return []string{head + "  " + value + "  [" + p.paint(mark, "bar") + "]"}
	}
	if displayWidth(plainLabel+"  "+value) <= inner {
		head := p.paint(label, kind) + strings.Repeat(" ", summaryLabelWidth-displayWidth(label))
		return []string{head + "  " + value}
	}
	var out []string
	if displayWidth(label) <= inner {
		out = append(out, p.paint(label, kind))
	} else {
		out = append(out, wrapText(label, inner)...)
	}
	out = append(out, "  "+value)
	return out
}

func progressMarks(num, den, width int) string {
	if den <= 0 || width <= 0 {
		return ""
	}
	filled := (num*width + den/2) / den
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return strings.Repeat("#", filled) + strings.Repeat("-", width-filled)
}

func fitSegments(inner int, head, tail string) string {
	line := head + tail
	if visibleWidth(line) <= inner {
		return line
	}
	return line
}

func frame(p *painter, lines []string, inner, padX int) string {
	border := lipgloss.RoundedBorder()
	span := inner + 2*padX
	pad := strings.Repeat(" ", padX)
	var b strings.Builder
	b.WriteString(p.paint(border.TopLeft+strings.Repeat(border.Top, span)+border.TopRight, "border"))
	b.WriteByte('\n')
	for _, line := range lines {
		line = padVisible(line, inner)
		b.WriteString(p.paint(border.Left, "border"))
		b.WriteString(pad)
		b.WriteString(line)
		b.WriteString(pad)
		b.WriteString(p.paint(border.Right, "border"))
		b.WriteByte('\n')
	}
	b.WriteString(p.paint(border.BottomLeft+strings.Repeat(border.Bottom, span)+border.BottomRight, "border"))
	b.WriteByte('\n')
	return b.String()
}

func visibleWidth(s string) int {
	return displayWidth(stripANSI(s))
}

func padVisible(s string, width int) string {
	gap := width - visibleWidth(s)
	if gap <= 0 {
		return s
	}
	return s + strings.Repeat(" ", gap)
}

func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == 0x1b {
			i = skipEscape(s, i+1)
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

func writeTable(b *strings.Builder, p *painter, items []Item, view View, width int) {
	idW, statusW, scopeW := 4, 12, 11
	titleW := width - (idW + statusW + scopeW + 3)
	if titleW < 8 {
		writeStacked(b, p, items, view, width)
		return
	}
	b.WriteString(p.paint(pad("ID", idW)+" "+pad("Status", statusW)+" "+pad("Scope", scopeW)+" Title", "heading") + "\n")
	for _, it := range items {
		writeItemRow(b, p, it, idW, statusW, scopeW, titleW)
		if view.Explain {
			writeExplain(b, it, width)
		}
	}
}

func writeStacked(b *strings.Builder, p *painter, items []Item, view View, width int) {
	for _, it := range items {
		scope := "active"
		scopeKind := ""
		if !it.Active {
			scope = "historical"
			scopeKind = "historical"
		}
		id := truncate(it.ID, min(5, width))
		rest := width - displayWidth(id) - 1
		if rest < 1 {
			rest = 1
		}
		b.WriteString(id + " " + p.paint(truncate(it.Status, rest), statusKind(it)) + "\n")
		b.WriteString("  " + p.paint(truncate(scope, width-2), scopeKind) + "\n")
		b.WriteString("  " + truncate(it.Title, width-2) + "\n")
		writeRelations(b, p, it, width)
		if view.Explain {
			writeExplain(b, it, width)
		}
	}
}

func (v View) filtering() bool {
	return len(v.Statuses) > 0 || v.Blocked || v.HumanReview || v.Active || v.Historical
}

func filterItems(items []Item, view View) []Item {
	if !view.filtering() {
		return items
	}
	out := make([]Item, 0, len(items))
	for _, it := range items {
		if itemMatches(it, view) {
			out = append(out, it)
		}
	}
	return out
}

func itemMatches(it Item, view View) bool {
	if len(view.Statuses) > 0 {
		ok := false
		for _, status := range view.Statuses {
			if it.Status == status {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	if view.Blocked && !itemBlocked(it) {
		return false
	}
	if view.HumanReview && !itemReview(it) {
		return false
	}
	if view.Active && !it.Active {
		return false
	}
	if view.Historical && it.Active {
		return false
	}
	return true
}

func itemBlocked(it Item) bool {
	for _, d := range it.Dependencies {
		if !d.Satisfied {
			return true
		}
	}
	return false
}

func itemReview(it Item) bool {
	for _, d := range it.Dependencies {
		if d.State == "human_review" {
			return true
		}
	}
	return false
}

func writeExplain(b *strings.Builder, it Item, width int) {
	for _, line := range explainLines(it) {
		for _, wrapped := range wrapText(line, width-2) {
			b.WriteString("  " + wrapped + "\n")
		}
	}
}

func explainLines(it Item) []string {
	var lines []string
	excluded := it.Supersession.Supplied && it.Supersession.Valid != nil && *it.Supersession.Valid
	if excluded {
		if it.Supersession.Line > 0 && it.Path != "" {
			lines = append(lines, "explain: excluded from active and recorded-verified counts; superseded-by "+sanitizeTerminal(it.SupersededBy)+" at "+it.Path+":"+itoa(it.Supersession.Line))
		} else {
			lines = append(lines, "explain: excluded from active and recorded-verified counts; superseded-by field line unavailable")
		}
	}
	for _, d := range it.Dependencies {
		if d.Satisfied && d.State != "human_review" {
			continue
		}
		loc := "source line unavailable"
		if d.Line > 0 && d.Path != "" {
			loc = d.Path + ":" + itoa(d.Line)
		}
		line := "explain: dependency " + sanitizeTerminal(d.ID) + " is " + sanitizeTerminal(d.State) + " at " + loc
		if d.State == "human_review" {
			line += "; the successor is not used and this is not independent verification"
		}
		lines = append(lines, line)
	}
	switch {
	case it.Tasks.HasNums && it.Tasks.Source != "":
		lines = append(lines, "explain: checked tasks "+itoa(it.Tasks.Checked)+"/"+itoa(it.Tasks.Total)+" in "+it.Tasks.Source)
		if len(it.Tasks.Checks) == 0 {
			lines = append(lines, "explain: checkbox line locations were not retained")
		}
		for _, c := range it.Tasks.Checks {
			mark := "unchecked"
			if c.Checked {
				mark = "checked"
			}
			if c.Line > 0 {
				lines = append(lines, "explain: "+mark+" checkbox at "+it.Tasks.Source+":"+itoa(c.Line))
			} else {
				lines = append(lines, "explain: "+mark+" checkbox line unavailable in "+it.Tasks.Source)
			}
		}
	case it.Tasks.State == "not_linked" || it.Tasks.State == "":
		lines = append(lines, "explain: no task file linked; checkbox line evidence unavailable")
	default:
		lines = append(lines, "explain: task file was "+it.Tasks.State+"; checkbox line evidence unavailable")
	}
	if len(lines) == 0 {
		lines = append(lines, "explain: no dependency, supersession-exclusion, or task-file evidence for this item")
	}
	return lines
}

func writeItemRow(b *strings.Builder, p *painter, it Item, idW, statusW, scopeW, titleW int) {
	scope := "active"
	scopeKind := ""
	if !it.Active {
		scope = "historical"
		scopeKind = "historical"
	}
	id := pad(truncate(it.ID, idW), idW)
	status := pad(truncate(it.Status, statusW), statusW)
	scopeText := pad(truncate(scope, scopeW), scopeW)
	title := truncate(it.Title, titleW)
	b.WriteString(id + " " + p.paint(status, statusKind(it)) + " " + p.paint(scopeText, scopeKind) + " " + title + "\n")
	writeRelations(b, p, it, idW+statusW+scopeW+titleW+3)
}

func writeRelations(b *strings.Builder, p *painter, it Item, width int) {
	if it.SupersededBy != "" {
		rel := "superseded-by " + sanitizeTerminal(it.SupersededBy)
		if it.Supersession.Valid != nil && *it.Supersession.Valid && it.Supersession.TerminalID != "" && it.Supersession.TerminalID != it.SupersededBy {
			rel += " terminal " + sanitizeTerminal(it.Supersession.TerminalID)
		}
		b.WriteString("  " + p.paint(truncate(rel, width-2), "historical") + "\n")
	}
	if len(it.Dependencies) == 0 {
		return
	}
	var parts []string
	review := false
	for _, d := range it.Dependencies {
		label := sanitizeTerminal(d.ID) + " " + sanitizeTerminal(d.State)
		parts = append(parts, label)
		if d.State == "human_review" {
			review = true
		}
	}
	line := "depends " + strings.Join(parts, ", ")
	kind := ""
	if review {
		kind = "warning"
		line = "human review: " + line
	}
	b.WriteString("  " + p.paint(truncate(line, width-2), kind) + "\n")
}

func statusKind(it Item) string {
	if !it.StatusKnown {
		return ""
	}
	switch it.Status {
	case "verified":
		return "verified"
	case "implemented":
		return "implemented"
	default:
		return ""
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [16]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
