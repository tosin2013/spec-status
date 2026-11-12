# dependency-cycle

- Command: /speckit.spec-status.unblock
- Signals: diag SS_DEPENDENCY_CYCLE; outcome=report
- Expected finding: 001 and 002 depend on each other.
- Evidence: 001 Depends on at .specify/memory/roadmap.md:7; 002 Depends on at .specify/memory/roadmap.md:11
- Expected proposal: propose removing one of the two Depends on edges (the user chooses which) via /speckit.roadmap.write
- Must not: change any file before the user approves the specific change
