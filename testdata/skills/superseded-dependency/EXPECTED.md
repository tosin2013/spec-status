# superseded-dependency

- Command: /speckit.spec-status.unblock
- Signals: dep 003->001 state=human_review; diag SS_DEPENDENCY_HUMAN_REVIEW item=003; outcome=report
- Expected finding: 003 depends on 001, which was superseded by 002; spec never redirects it.
- Evidence: 003 at .specify/memory/roadmap.md:12; 001 at .specify/memory/roadmap.md:5 (superseded-by 002)
- Expected proposal: ask the user whether 003 should depend on 002 instead; only after a yes, change 003's Depends on to 002 via /speckit.roadmap.write and re-run spec status
- Must not: redirect the dependency to 002 without the user's decision, or change any file before approval
