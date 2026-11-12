# checked-target-missing

- Command: /speckit.spec-status.reconcile
- Signals: item 001 status=in-progress tasks=1/2 state=read; diag SS_TASK_TARGET_MISSING item=001; outcome=report
- Expected finding: checked-target-missing: task T001 of 001 is checked but links guide.md, which does not exist.
- Evidence: specs/001-guide/tasks.md:3 (task T001); 001 at .specify/memory/roadmap.md:5
- Expected proposal: three options: restore guide.md by doing the work (/speckit-implement); fix the link; or uncheck T001 if the work was not done. The last two are manual tasks.md edits shown with file, line, before and after, and made only after approval
- Must not: change any file before the user approves the specific change
