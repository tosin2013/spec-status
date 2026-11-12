package status

import "testing"

// Task-file links inside inline code spans are not links (CommonMark §6.1).
// Roadmap sync 2026-11-12, finding F5; decision D-2026-11-12-01.

func taskDiagCounts(t *testing.T, tasks string) (unsupported, missing int) {
	t.Helper()
	body := "### 001 — A [status: planned]\n- Spec dir: specs/a\n"
	dir := project(t, body, map[string]string{"specs/a/tasks.md": tasks})
	rep := Evaluate(Option{Root: dir, CWD: dir})
	if rep.Outcome == "fatal" {
		t.Fatalf("fatal: %+v", rep.Fatal)
	}
	for _, d := range rep.Diagnostics {
		switch d.Code {
		case "SS_REFERENCE_UNSUPPORTED":
			unsupported++
		case "SS_TASK_TARGET_MISSING":
			missing++
		}
	}
	return unsupported, missing
}

func TestLinkLikeTextInCodeSpanIsNotALink(t *testing.T) {
	for _, tasks := range []string{
		"- [ ] Check the tag against `^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)$`\n",
		"- [ ] Write `[label](missing.md)` in the notes\n",
		"- [ ] Write ``a ` [x](missing.md) `` literally\n",
		"- [ ] Autolink in code `<docs/missing.md>`\n",
		"- [ ] Reference in code `[a][b]`\n",
		"- [ ] First line\n  continued with `[x](missing.md)` in code\n",
	} {
		if u, m := taskDiagCounts(t, tasks); u != 0 || m != 0 {
			t.Errorf("%q: unsupported=%d missing=%d, want 0 0", tasks, u, m)
		}
	}
}

func TestLinksOutsideCodeSpansStillFound(t *testing.T) {
	for _, c := range []struct {
		tasks       string
		unsupported int
		missing     int
	}{
		// A real link next to a code span is still checked.
		{"- [ ] See `code` and [file](missing.md)\n", 0, 1},
		// An unmatched backtick is literal text, so the link after it counts.
		{"- [ ] Stray ` then [file](missing.md)\n", 0, 1},
		// Runs of different length do not close each other.
		{"- [ ] Open `` then ` and [file](missing.md)\n", 0, 1},
		// Unsupported shapes outside code stay visible.
		{"- [ ] Ref [a][b] and `code`\n\n[b]: x.md\n", 1, 0},
		{"- [ ] Web [w](https://example.test/a) and `code`\n", 1, 0},
	} {
		u, m := taskDiagCounts(t, c.tasks)
		if u != c.unsupported || m != c.missing {
			t.Errorf("%q: unsupported=%d missing=%d, want %d %d", c.tasks, u, m, c.unsupported, c.missing)
		}
	}
}
