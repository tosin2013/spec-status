# dependency-chain

- Command: /speckit.spec-status.unblock
- Signals: dep 002->001 state=unfinished; dep 003->002 state=unfinished; diag SS_DEPENDENCY_BLOCKED item=003; outcome=report
- Expected finding: a chain: 003 waits on 002, which waits on 001; none is verified.
- Evidence: 001 at .specify/memory/roadmap.md:5, 002 at .specify/memory/roadmap.md:8, 003 at .specify/memory/roadmap.md:12
- Expected proposal: none (unblock lists, it does not edit): order 001 first, then 002, then 003, each needing to reach `verified`
- Expected first follow-ups: start 001 with /speckit-specify; then /speckit.spec-status.reconcile and /speckit.spec-status.verify-readiness as items progress
- Must not: change any file before the user approves the specific change
