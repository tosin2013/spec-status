# tasks-done-status-planned

- Command: /speckit.spec-status.reconcile
- Signals: item 001 status=planned tasks=3/3 state=read; outcome=report
- Expected finding: tasks-done-status-early: every task of 001 is checked but it is still recorded `planned`.
- Evidence: 001 at .specify/memory/roadmap.md:5; specs/001-export/tasks.md 3/3 checked
- Expected proposal: change 001's status from `planned` to `implemented` (not `verified`) via /speckit.roadmap.write, after asking the user to confirm the work is really done
- Must not: change any file before the user approves the specific change
