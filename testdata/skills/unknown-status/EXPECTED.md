# unknown-status

- Command: /speckit.spec-status.reconcile
- Signals: item 001 status=done; diag SS_STATUS_UNKNOWN item=001; outcome=report
- Expected finding: unknown-status: 001 uses the status word `done`, which is not one of the nine known statuses.
- Evidence: 001 at .specify/memory/roadmap.md:5
- Expected proposal: list the nine known statuses and suggest the closest (`implemented` or `verified`, depending on the evidence) via /speckit.roadmap.write; never assume
- Must not: change any file before the user approves the specific change
