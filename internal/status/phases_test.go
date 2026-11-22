package status

import (
	"fmt"
	"strings"
	"testing"
)

// Feature 011: per-phase task counts.

func phaseString(ps []Phase) string {
	var b strings.Builder
	for _, p := range ps {
		story := "null"
		if p.Story != nil {
			story = fmt.Sprint(*p.Story)
		}
		fmt.Fprintf(&b, "%d|%s|%s|%d/%d\n", p.Line, p.Heading, story, p.Checked, p.Total)
	}
	return b.String()
}

func mustParseTasks(t *testing.T, src string) taskFile {
	t.Helper()
	tf, d := parseTasks("specs/x/tasks.md", []byte(src))
	if d != nil {
		t.Fatalf("unexpected diagnostic: %+v", *d)
	}
	return tf
}

func TestPhasesGeneratedLayout(t *testing.T) {
	src := strings.Join([]string{
		"# Tasks: Sample", // 1
		"",
		"- [x] T000 before any phase", // 3
		"",
		"## Phase 1: Setup (Shared Infrastructure)", // 5
		"",
		"- [x] T001 a",
		"- [ ] T002 b",
		"",
		"## Phase 2: Foundational", // 10
		"",
		"- [x] T003 c",
		"",
		"## Phase 3: User Story 1 - First (Priority: P1) 🎯 MVP", // 14
		"",
		"### Tests for User Story 1",
		"",
		"- [x] T004 d",
		"",
		"### Implementation for User Story 1",
		"",
		"- [ ] T005 e",
		"  - [x] T005a nested",
		"",
		"## Phase N: Polish & Cross-Cutting Concerns ##", // 25
		"",
		"- [ ] T006 f",
		"",
		"## Dependencies & Execution Order", // 29
		"",
		"- [x] T007 after phases",
		"",
	}, "\n")
	tf := mustParseTasks(t, src)
	want := "5|Phase 1: Setup (Shared Infrastructure)|null|1/2\n" +
		"10|Phase 2: Foundational|null|1/1\n" +
		"14|Phase 3: User Story 1 - First (Priority: P1) 🎯 MVP|1|2/3\n" +
		"25|Phase N: Polish & Cross-Cutting Concerns|null|0/1\n"
	if got := phaseString(tf.Phases); got != want {
		t.Fatalf("phases:\n%s\nwant:\n%s", got, want)
	}
	inChecked, inTotal := 0, 0
	for _, p := range tf.Phases {
		inChecked += p.Checked
		inTotal += p.Total
	}
	// Item totals include T000 (checked) and T007 (checked), outside phases.
	if tf.Total != inTotal+2 || tf.Checked != inChecked+2 {
		t.Fatalf("totals %d/%d, in-phase %d/%d", tf.Checked, tf.Total, inChecked, inTotal)
	}
}

func TestPhasesHeadingRules(t *testing.T) {
	src := strings.Join([]string{
		"## Phase", // 1: bare word qualifies
		"- [x] a",
		"# Top level", // 3: ends the phase
		"- [x] b",
		"## Phases overview", // 5: not a phase
		"- [ ] c",
		"## phase 2", // 7: case-sensitive
		"- [ ] d",
		"##Phase 3", // 9: not an ATX heading
		"- [ ] e",
		"   ## Phase:colon", // 11: three spaces of indentation is a heading
		"- [x] f",
		"## Phase 4 #",                     // 13: closing sequence removed
		"### Phase 5 is a level-3 heading", // 14: stays inside
		"- [x] g",
		"## Phase 6 C#", // 16: '#' not preceded by a space is text
		"## Phase 4 #",  // 17: duplicate heading text
		"- [ ] h",
		"",
	}, "\n")
	tf := mustParseTasks(t, src)
	want := "1|Phase|null|1/1\n" +
		"11|Phase:colon|null|1/1\n" +
		"13|Phase 4|null|1/1\n" +
		"16|Phase 6 C#|null|0/0\n" +
		"17|Phase 4|null|0/1\n"
	if got := phaseString(tf.Phases); got != want {
		t.Fatalf("phases:\n%s\nwant:\n%s", got, want)
	}
}

func TestPhaseStory(t *testing.T) {
	cases := map[string]string{
		"Phase 3: User Story 1 - x":                      "1",
		"Phase 9: User Story 12 (Priority)":              "12",
		"Phase 9: User Story 123456789":                  "123456789",
		"Phase 9: User Story 1234567890":                 "null",
		"Phase 9: User Story x":                          "null",
		"Phase 9: user story 3":                          "null",
		"Phase 2: Foundational":                          "null",
		"Phase 5: User Story 2 and User Story 3":         "2",
		"Phase 5: User Story cleanup, then User Story 2": "2",
		"Phase 5: User Story 1234567890 / User Story 3":  "3",
		"Phase 5: User Story x and User Story y":         "null",
	}
	for h, want := range cases {
		got := "null"
		if p := phaseStory(h); p != nil {
			got = fmt.Sprint(*p)
		}
		if got != want {
			t.Errorf("phaseStory(%q) = %s, want %s", h, got, want)
		}
	}
}

func TestPhasesHiddenHeadings(t *testing.T) {
	src := strings.Join([]string{
		"```",
		"## Phase 1: fenced",
		"- [x] fenced, not a task",
		"```",
		"<!--",
		"## Phase 2: commented",
		"-->",
		"    ## Phase 3: indented code",
		"",
		"- [x] T001 real",
		"",
	}, "\n")
	tf := mustParseTasks(t, src)
	if len(tf.Phases) != 0 {
		t.Fatalf("hidden headings opened phases:\n%s", phaseString(tf.Phases))
	}
	if tf.Total != 1 || tf.Checked != 1 {
		t.Fatalf("totals %d/%d, want 1/1", tf.Checked, tf.Total)
	}
}

func TestPhasesSetextIsALimit(t *testing.T) {
	src := "## Phase 1: Setup\n\n- [x] a\n\nNot a phase end\n---\n\n- [ ] b\n"
	tf := mustParseTasks(t, src)
	if got, want := phaseString(tf.Phases), "1|Phase 1: Setup|null|1/2\n"; got != want {
		t.Fatalf("phases:\n%s\nwant:\n%s", got, want)
	}
}

// TestPhasesFixture checks testdata/phases against the hand counts in its
// README, including the JSON null/[] rule.
func TestPhasesFixture(t *testing.T) {
	rep := Evaluate(Option{Root: "../../testdata/phases", CWD: "."})
	if rep.Outcome != "report" {
		t.Fatalf("outcome %s: %+v", rep.Outcome, rep.Fatal)
	}
	p001 := "5|Phase 1: Setup|null|1/2\n" +
		"19|Phase 2: Foundational|null|0/0\n" +
		"21|Phase 3: User Story 1 - Read one phase (Priority: P1)|1|2/3\n" +
		"32|Phase 4: User Story 1234567890 - too long a number|null|0/1\n"
	want := map[string]struct{ state, phases string }{
		"001": {"read", p001},
		"002": {"read", ""},
		"003": {"missing", ""},
		"004": {"not_linked", ""},
		"005": {"read", p001},
	}
	for _, it := range rep.Items {
		w, ok := want[it.ID]
		if !ok {
			t.Fatalf("unexpected item %s", it.ID)
		}
		if it.Tasks.State != w.state {
			t.Errorf("%s: state %s, want %s", it.ID, it.Tasks.State, w.state)
		}
		if got := phaseString(it.Tasks.Phases); got != w.phases {
			t.Errorf("%s: phases:\n%s\nwant:\n%s", it.ID, got, w.phases)
		}
	}
	// Shared files must not alias: changing 001's copy leaves 005's alone.
	rep.Items[0].Tasks.Phases[0].Total = 99
	if rep.Items[4].Tasks.Phases[0].Total != 2 {
		t.Fatal("items sharing a tasks file alias the same phases")
	}
	js := string(JSON(Evaluate(Option{Root: "../../testdata/phases", CWD: "."})))
	for _, frag := range []string{
		`"state":"read","phases":[]}`,
		`"state":"missing","phases":null}`,
		`"state":"not_linked","phases":null}`,
		`"phases":[{"heading":"Phase 1: Setup","line":5,"story":null,"checked":1,"total":2},`,
		`{"heading":"Phase 3: User Story 1 - Read one phase (Priority: P1)","line":21,"story":1,"checked":2,"total":3}`,
	} {
		if !strings.Contains(js, frag) {
			t.Errorf("JSON lacks %s", frag)
		}
	}
	if strings.Contains(js, "Phase 9") || strings.Contains(js, "Phase 8") {
		t.Error("hidden heading reached the JSON")
	}
}
