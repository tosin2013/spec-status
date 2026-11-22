package human

import (
	"github.com/rivo/uniseg"
	"github.com/tosin2013/spec-status/internal/status"
	"path/filepath"
	"strings"
	"testing"
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
	js := string(status.JSON(rep))
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
	rep.Summary.Tasks = status.TaskSummary{Checked: 0, Total: 0, Percent: percentOf(0, 0), Scope: taskScope}
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
	js := string(status.JSON(rep))
	if !strings.Contains(js, `\u001b`) && !strings.Contains(js, `\u001B`) {
		t.Fatalf("JSON lost the raw escape: %s", js)
	}
}

func sampleReport() status.Report {
	valid := true
	return status.Report{
		Outcome: "report", Complete: true, Roadmap: ".specify/memory/roadmap.md", FailOn: "error",
		Summary: &status.Summary{
			TotalItems: 3, ActiveItems: 2, VerifiedActiveItems: 1,
			VerifiedPercent: percentOf(1, 2), SupersededItems: 1,
			BlockedCount: 1, HumanReviewCount: 1,
			StatusCounts: map[string]int{"implemented": 1, "planned": 1, "verified": 1},
			Tasks:        status.TaskSummary{Checked: 1, Total: 2, Percent: percentOf(1, 2), Scope: taskScope, FilesRead: 2},
		},
		Items: []status.Item{
			{ID: "001", Title: "Preserve original feature", Status: "implemented", StatusKnown: true, SupersededBy: "002", Active: false, Supersession: status.Supersession{Supplied: true, Valid: &valid, TerminalID: "002"}},
			{ID: "002", Title: "Replacement feature", Status: "verified", StatusKnown: true, Active: true},
			{ID: "003", Title: "Follow-on work", Status: "planned", StatusKnown: true, Active: true, Dependencies: []status.Dependency{{ID: "001", State: "human_review"}}},
		},
		Diagnostics: []status.Diagnostic{{Severity: "warning", Code: "SS_DEPENDENCY_HUMAN_REVIEW", Message: "Dependency on a superseded item requires human review.", ItemID: "003", Path: ".specify/memory/roadmap.md", Line: 18}},
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

func TestInventoryTextSanitizesAndStaysPlain(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".specify", ".gitkeep"), "")
	evil := "001-x\x1b[31mred\x1b]0;title\x07"
	writeFile(t, filepath.Join(dir, "specs", evil, "spec.md"), "# x\n")
	rep := status.EvaluateInventory(status.InventoryOption{Root: dir})
	plain := InventoryText(rep, View{Width: 80})
	if strings.ContainsRune(plain, 0x1b) || strings.ContainsRune(plain, 0x07) {
		t.Fatalf("terminal escapes leaked: %q", plain)
	}
	if !strings.Contains(plain, "001-xred") {
		t.Fatalf("sanitized name missing: %q", plain)
	}
	narrow := InventoryText(rep, View{Width: 30})
	if !strings.Contains(narrow, "\n  spec present  plan missing\n") {
		t.Fatalf("narrow layout not used at 30 columns:\n%s", narrow)
	}
}

// Feature 008 US4: human text shows status counts, invalid supersession and
// shared tasks; status.JSON is unchanged.
func TestHumanTextStatusesInvalidAndShared(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "example-roadmap")
	text := Text(status.Evaluate(status.Option{Root: root}), View{Width: 80})
	if !strings.Contains(text, "statuses implemented 1 planned 1 verified 1") {
		t.Errorf("statuses line missing:\n%s", text)
	}
	if strings.Contains(text, "invalid supersession") || strings.Contains(text, "(invalid)") {
		t.Errorf("invalid supersession shown without one:\n%s", text)
	}
	empty := Text(status.Evaluate(status.Option{Root: filepath.Join("..", "..", "testdata", "two-mode", "explicit-empty")}), View{Width: 80})
	if !strings.Contains(empty, "statuses none") {
		t.Errorf("empty roadmap statuses line:\n%s", empty)
	}
	base := "- Description: d\n- Outcome: o\n- Scope: s\n"
	dir := project(t, "### 001 — A [status: planned]\n"+base+"- superseded-by: 009\n\n### 002 — B [status: planned]\n"+base+"- superseded-by: soon\n\n### 003 — C [status: wip\x1b[31m]\n"+base, nil)
	rep := status.Evaluate(status.Option{Root: dir, CWD: dir})
	text = Text(rep, View{Width: 80})
	for _, want := range []string{"invalid supersession 2", "superseded-by 009 (invalid)", "superseded-by soon (invalid)", "statuses planned 2 wip 1"} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q:\n%s", want, text)
		}
	}
	if strings.ContainsRune(text, 0x1b) {
		t.Error("terminal escape leaked from a status token")
	}
	shared := project(t, "### 001 — A [status: in-progress]\n"+base+"- Spec dir: specs/001-x\n\n### 002 — B [status: in-progress]\n"+base+"- Spec dir: specs/001-x\n", map[string]string{
		"specs/001-x/spec.md": "# s\n", "specs/001-x/plan.md": "# p\n", "specs/001-x/tasks.md": "- [x] T001 a\n- [ ] T002 b\n",
	})
	srep := status.Evaluate(status.Option{Root: shared, CWD: shared})
	stext := Text(srep, View{Width: 80})
	if strings.Count(stext, "tasks shared with 001; counted once") != 1 {
		t.Errorf("shared tasks note:\n%s", stext)
	}
	if strings.Contains(string(status.JSON(srep)), "shared") || strings.Contains(string(status.JSON(rep)), "(invalid)") {
		t.Error("display additions leaked into status.JSON")
	}
}
