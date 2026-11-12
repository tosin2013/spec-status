package status

import (
	"strings"
	"testing"

	"github.com/rivo/uniseg"
)

func TestSanitizeDropsTerminalControls(t *testing.T) {
	in := "Hello\x1b[31mX\x1b[0m \x1b]8;;http://evil.example\x1b\\link\x07tail"
	got := sanitizeTerminal(in)
	if strings.Contains(got, "\x1b") || strings.Contains(got, "http://evil") {
		t.Fatalf("controls survived: %q", got)
	}
	if !strings.Contains(got, "Hello") || !strings.Contains(got, "X") || !strings.Contains(got, "link") || !strings.Contains(got, "tail") {
		t.Fatalf("text was discarded: %q", got)
	}
}

func TestTruncateKeepsWideRunesWhole(t *testing.T) {
	got := truncate("你好世界", 5)
	if strings.Contains(got, "\uFFFD") || !utf8Valid(got) {
		t.Fatalf("invalid %q", got)
	}
	if uniseg.StringWidth(got) > 5 {
		t.Fatalf("width %d for %q", uniseg.StringWidth(got), got)
	}
	if !strings.HasPrefix(got, "你好") {
		t.Fatalf("split a wide character: %q", got)
	}
}

func TestHumanTextColorIsLabelsNotJSON(t *testing.T) {
	rep := sampleReport()
	plain := Text(rep, View{Width: 80})
	color := Text(rep, View{Width: 80, Color: true})
	if strings.Contains(plain, "\x1b") {
		t.Fatal("plain text contains ANSI")
	}
	if !strings.Contains(color, "\x1b") {
		t.Fatal("styled text has no ANSI")
	}
	for _, text := range []string{plain, color} {
		for _, phrase := range []string{"recorded verified", "historical", "active", "human review", "implemented", "verified"} {
			if !strings.Contains(stripANSI(text), phrase) {
				t.Fatalf("missing %q\n%s", phrase, text)
			}
		}
	}
	if Text(rep, View{Width: 80, Color: true}) != color {
		t.Fatal("styled text is not deterministic")
	}
	js := string(JSON(rep))
	if strings.Contains(js, "\x1b[") {
		t.Fatal("JSON gained styling")
	}
	narrow := stripANSI(Text(rep, View{Width: 24}))
	if !strings.Contains(narrow, "Preserve") || !strings.Contains(narrow, "1/2") || !strings.Contains(narrow, "╭") {
		t.Fatal(narrow)
	}
	if stripANSI(color) != plain {
		t.Fatalf("styled facts differ from plain\nplain:\n%s\nstripped:\n%s", plain, stripANSI(color))
	}
}

func TestSummaryKeepsZeroDenominatorsUnadorned(t *testing.T) {
	rep := sampleReport()
	rep.Summary.ActiveItems = 0
	rep.Summary.VerifiedActiveItems = 0
	rep.Summary.VerifiedPercent = percentOf(0, 0)
	rep.Summary.Tasks = TaskSummary{Checked: 0, Total: 0, Percent: percentOf(0, 0), Scope: taskScope}
	rep.Summary.BlockedCount = 0
	rep.Summary.HumanReviewCount = 0
	rep.Items = nil
	rep.Diagnostics = nil
	text := Text(rep, View{Width: 40})
	if strings.Contains(text, "[") || strings.Contains(text, "#") {
		t.Fatalf("empty progress decoration:\n%s", text)
	}
	if !strings.Contains(text, "0/0 (N/A)") || !strings.Contains(text, "human review") || !strings.Contains(text, "active blockers") {
		t.Fatalf("counts clipped:\n%s", text)
	}
	if !strings.Contains(text, "╭") {
		t.Fatal("missing summary border")
	}
}

func TestEscapeInTitleStaysRawInJSONOnly(t *testing.T) {
	rep := sampleReport()
	rep.Items[1].Title = "Safe\x1b[31mNO\x1b[0m"
	text := Text(rep, View{Width: 80, Color: true})
	if strings.Contains(text, "\x1b[31mNO") {
		t.Fatal("title injected color")
	}
	if !strings.Contains(stripANSI(text), "SafeNO") && !strings.Contains(stripANSI(text), "Safe NO") && !strings.Contains(stripANSI(text), "Safe") {
		t.Fatal(stripANSI(text))
	}
	js := string(JSON(rep))
	if !strings.Contains(js, `\u001b`) && !strings.Contains(js, `\u001B`) {
		t.Fatalf("JSON lost the raw escape: %s", js)
	}
}

func sampleReport() Report {
	valid := true
	return Report{
		Outcome: "report", Complete: true, Roadmap: ".specify/memory/roadmap.md", FailOn: "error",
		Summary: &Summary{
			TotalItems: 3, ActiveItems: 2, VerifiedActiveItems: 1,
			VerifiedPercent: percentOf(1, 2), SupersededItems: 1,
			BlockedCount: 1, HumanReviewCount: 1,
			StatusCounts: map[string]int{"implemented": 1, "planned": 1, "verified": 1},
			Tasks:        TaskSummary{Checked: 1, Total: 2, Percent: percentOf(1, 2), Scope: taskScope, FilesRead: 2},
		},
		Items: []Item{
			{ID: "001", Title: "Preserve original feature", Status: "implemented", StatusKnown: true, SupersededBy: "002", Active: false, Supersession: Supersession{Supplied: true, Valid: &valid, TerminalID: "002"}},
			{ID: "002", Title: "Replacement feature", Status: "verified", StatusKnown: true, Active: true},
			{ID: "003", Title: "Follow-on work", Status: "planned", StatusKnown: true, Active: true, Dependencies: []Dependency{{ID: "001", State: "human_review"}}},
		},
		Diagnostics: []Diagnostic{{Severity: "warning", Code: "SS_DEPENDENCY_HUMAN_REVIEW", Message: "Dependency on a superseded item requires human review.", ItemID: "003", Path: ".specify/memory/roadmap.md", Line: 18}},
	}
}

func utf8Valid(s string) bool {
	return strings.ToValidUTF8(s, "") == s
}

// Build label in the panel heading (feature 004, contracts/build-identity.md).
func TestPanelHeadingFollowsBuildLabel(t *testing.T) {
	rep := sampleReport()
	def := Text(rep, View{Width: 80})
	dev := Text(rep, View{Width: 80, Build: "development build"})
	rel := Text(rep, View{Width: 80, Build: "release 0.1.0"})
	if def != dev {
		t.Fatalf("empty Build must render as development build:\n%s\n%s", def, dev)
	}
	if !strings.Contains(def, "Spec Status  development build") {
		t.Fatalf("default heading changed:\n%s", def)
	}
	if !strings.Contains(rel, "Spec Status  release 0.1.0") || strings.Contains(rel, "development build") {
		t.Fatalf("release heading:\n%s", rel)
	}
	if strings.Replace(rel, "release 0.1.0    ", "development build", 1) != def {
		t.Fatalf("only the heading label may differ:\n%s\n%s", rel, def)
	}
}
