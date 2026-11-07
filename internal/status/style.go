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
	Color bool
	Width int
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
	if !on {
		return &painter{}
	}
	// The renderer writes nowhere and is given an explicit ANSI profile.
	// It does not query the terminal or open /dev/tty.
	r := lipgloss.NewRenderer(io.Discard, termenv.WithProfile(termenv.ANSI))
	r.SetColorProfile(termenv.ANSI)
	r.SetHasDarkBackground(true)
	return &painter{on: true, r: r}
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
	head := p.paint("Spec Status", "heading")
	b.WriteString(head + "  development build\n")
	b.WriteString("source  " + truncate(rep.Roadmap, width-8) + "\n")
	if rep.Summary == nil {
		b.WriteString("summary  N/A\n")
		return b.String()
	}
	s := rep.Summary
	b.WriteString("items " + itoa(s.TotalItems) + "   active " + itoa(s.ActiveItems) + "   historical " + itoa(s.TotalItems-s.ActiveItems) + "\n")
	b.WriteString("statuses  " + truncate(formatCounts(s.StatusCounts), width-10) + "\n")
	b.WriteString(metricLine(p, "recorded verified", s.VerifiedActiveItems, s.ActiveItems, s.VerifiedPercent.Text(), width))
	b.WriteString(metricLine(p, "tasks", s.Tasks.Checked, s.Tasks.Total, s.Tasks.Percent.Text(), width))
	b.WriteString(truncate("  task scope "+sanitizeTerminal(s.Tasks.Scope)+"  read "+itoa(s.Tasks.FilesRead)+"  missing "+itoa(s.Tasks.FilesMissing)+"  skipped "+itoa(s.Tasks.FilesSkipped), width) + "\n")
	b.WriteString("superseded " + itoa(s.SupersededItems) + "   invalid supersession " + itoa(s.InvalidSupersessionItems) + "\n")
	review := "human review " + itoa(s.HumanReviewCount)
	if s.HumanReviewCount > 0 {
		review = p.paint(review, "warning")
	}
	b.WriteString("blocked active " + itoa(s.BlockedCount) + "   " + review + "\n")
	b.WriteString("\n")
	if width >= 60 {
		writeTable(&b, p, rep, width)
	} else {
		writeStacked(&b, p, rep, width)
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

func metricLine(p *painter, label string, num, den int, percent string, width int) string {
	text := label + "  " + itoa(num) + "/" + itoa(den) + " (" + percent + "%)"
	barWidth := width - displayWidth(text) - 3
	if barWidth < 8 {
		return truncate(text, width) + "\n"
	}
	filled := 0
	if den > 0 {
		filled = (num*barWidth + den/2) / den
		if filled > barWidth {
			filled = barWidth
		}
	}
	bar := strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled)
	return text + "  [" + p.paint(bar, "bar") + "]\n"
}

func writeTable(b *strings.Builder, p *painter, rep Report, width int) {
	idW, statusW, scopeW := 4, 12, 11
	titleW := width - (idW + statusW + scopeW + 3)
	if titleW < 8 {
		writeStacked(b, p, rep, width)
		return
	}
	b.WriteString(p.paint(pad("ID", idW)+" "+pad("Status", statusW)+" "+pad("Scope", scopeW)+" Title", "heading") + "\n")
	for _, it := range rep.Items {
		writeItemRow(b, p, it, idW, statusW, scopeW, titleW)
	}
}

func writeStacked(b *strings.Builder, p *painter, rep Report, width int) {
	for _, it := range rep.Items {
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
	}
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
