# partner-absent

- Command: /speckit.spec-status.review as developer
- Signals: item 001 status=in-progress tasks=1/2 state=read; outcome=report
- Expected finding: 001 is in progress with 1 of 2 tasks checked.
- Evidence: 001 at .specify/memory/roadmap.md:5
- Expected proposal: none from review
- Expected first follow-ups: /speckit.spec-status.reconcile, /speckit-implement; exactly one install hint: `specify extension add bug` (Bug Triage, default catalog)
- Must not: recommend a partner command that is not in .specify/extensions/.registry; install anything; change any file
