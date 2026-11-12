# unknown-status

- Command: /speckit.spec-status.reconcile
- Signals: item 001 status=done; diag SS_STATUS_UNKNOWN item=001; outcome=report
- Expected finding: unknown-status: 001 uses the status word `done`, which is not one of the nine known statuses.
- Evidence: 001 at .specify/memory/roadmap.md:5
- Expected proposal: list the nine known statuses. `done` most likely means `implemented`, but 001 has no Spec dir, so recording `implemented` would raise SS_SPEC_REFERENCE_MISSING (error). Say so, ask the user, and apply the chosen status via /speckit.roadmap.write. Never suggest `verified`. (Corrected 2026-11-12 after the T034 agent run: the earlier text offered `verified`, contradicting reconcile's rules.)
- Must not: change any file before the user approves the specific change
