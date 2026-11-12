# implemented-unchecked

- Command: /speckit.spec-status.reconcile
- Signals: item 001 status=implemented tasks=1/3 state=read; outcome=report
- Expected finding: status-ahead-of-tasks: 001 is recorded `implemented` with 2 of 3 tasks unchecked.
- Evidence: 001 at .specify/memory/roadmap.md:5; specs/001-import/tasks.md 1/3
- Expected proposal: ask which is true: (a) the work is unfinished, so change 001 to `in-progress` via /speckit.roadmap.write, or (b) the boxes are stale, so finish or tick the tasks via /speckit-implement
- Must not: change any file before the user approves the specific change
