package extension

import "testing"

func TestSectionBodyStopsAtIndentedHeading(t *testing.T) {
	text := "## A\n| x |\n|---|\n| 1 |\n  ## B\n| y |\n|---|\n| 2 |\n"
	if got := readTable(sectionBody(text, "A")); len(got) != 1 || got[0]["x"] != "1" {
		t.Fatalf("section A bled into B: %v", got)
	}
	if got := readTable(sectionBody(text, "B")); len(got) != 1 || got[0]["y"] != "2" {
		t.Fatalf("indented heading B not found: %v", got)
	}
	if _, ok := level2Heading("    ## code"); ok {
		t.Fatal("four-space indent is code, not a heading")
	}
}
