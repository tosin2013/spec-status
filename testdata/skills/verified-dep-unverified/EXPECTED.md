# verified-dep-unverified

- Command: /speckit.spec-status.reconcile
- Signals: diag SS_VERIFIED_DEPENDENCY_CONTRADICTION item=001; dep 001->002 state=unfinished; outcome=report
- Expected finding: verified-dependency-unverified: 001 is recorded `verified` while its dependency 002 is only `planned`.
- Evidence: 001 at .specify/memory/roadmap.md:5; 002 at .specify/memory/roadmap.md:10
- Expected proposal: hand off to /speckit.spec-status.unblock for 002, or lower 001 to `implemented` via /speckit.roadmap.write; ask the user which
- Must not: change any file before the user approves the specific change
