# roles-mixed

- Command: /speckit.spec-status.review (no role, then as pm, developer, qa, architect)
- Signals: item 001 status=verified tasks=2/2 state=read; item 002 status=verified tasks=1/2 state=read; item 003 status=planned tasks=2/2 state=read; dep 004->005 state=unfinished; dep 008->006 state=human_review; summary active=7 verified=2 blocked=2 human_review=1; outcome=report
- Expected finding (default view): progress (verified 2 of 7 active items: 001 and 002; 006 is superseded history); blockers 004 (waits on 005) and 008 (human review of superseded 006); evidence gap on 002 (`verified` with 1 of 2 tasks checked); 003 has all tasks checked but is `planned`; ready to start: 005 and 007.
- Evidence: 001 .specify/memory/roadmap.md:5, 002 .specify/memory/roadmap.md:9, 003 .specify/memory/roadmap.md:13, 004 .specify/memory/roadmap.md:17, 005 .specify/memory/roadmap.md:21, 006 .specify/memory/roadmap.md:24, 007 .specify/memory/roadmap.md:28, 008 .specify/memory/roadmap.md:31
- Expected proposal: none from review itself. Primary routing: 002 and 003 → /speckit.spec-status.reconcile (002 has /speckit.spec-status.verify-readiness as secondary, because it is recorded `verified`); 004 and 008 → /speckit.spec-status.unblock. ready-to-start is 005 and 007 (003 is left out because its tasks are all checked). missing-artifacts is empty: the info-level spec-link findings on planned items are left out, and 006 is superseded history.
- Expected first follow-ups and leading sections, per view (facts identical in every view):
  - default: leads with progress, blockers, human-review, diagnostics, ready-to-start; first follow-ups /speckit.spec-status.reconcile, /speckit.spec-status.unblock.
  - pm: leads with progress (2 of 7 recorded verified), blockers (004 waits on 005; 008 needs a decision about superseded 006), ready-to-start (005, 007), human-review; first follow-ups /speckit.spec-status.unblock, /speckit.spec-status.reconcile; no dates.
  - developer: leads with ready-to-start (005, 007), unchecked-tasks (002: 1 of 2), missing-artifacts (none), blockers; first follow-ups /speckit.spec-status.reconcile, /speckit-tasks, /speckit-implement; mentions 003 (all tasks checked, still planned).
  - qa: leads with evidence-gaps (002: recorded verified with 1 of 2 tasks checked), diagnostics, unchecked-tasks; first follow-ups /speckit.spec-status.verify-readiness, /speckit.roadmap.debrief, /speckit-analyze.
  - architect: leads with dependency-chains (004 → 005), supersession (006 superseded by 007), human-review (008 → 006), constraints (none recorded); first follow-ups /speckit.spec-status.unblock, /speckit.spec-status.explain, /speckit.roadmap.sync.
  - Partners: none registered, so qa mentions Assure once, architect mentions arch-governance once, and developer mentions Bug Triage once, each with its install command.
- Must not: show different counts or findings in any role view; change any file.
