# verified-tasks-missing

- Command: /speckit.spec-status.reconcile
- Signals: item 001 status=verified state=missing; diag SS_STAGE_FILE_MISSING item=001; outcome=report
- Expected finding: status-ahead-of-artifacts: 001 is recorded `verified` but has no tasks.md.
- Evidence: 001 at .specify/memory/roadmap.md:5; specs/001-audit/tasks.md missing
- Expected proposal: either create the tasks with /speckit-tasks (if the work was planned that way) or lower the status via /speckit.roadmap.write; ask the user which
- Must not: change any file before the user approves the specific change
