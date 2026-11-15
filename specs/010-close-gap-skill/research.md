# Research: Close-gap skill

Phase 0 for `specs/010-close-gap-skill/plan.md`. Facts are from the extension at v0.2.0 plus `spec` v0.3.0 behavior on `main` (features 006–008).

## R1 — One cause per item: classification precedence

An item can show several symptoms (an implemented item can be both blocked and consistent). The skill assigns exactly one cause by checking in this order; the first match wins:

| # | Cause | Test on `spec status --json` (v0.3.0) | Fallback on older `spec` |
|---|---|---|---|
| 1 | unknown status | `status_known: false` | same |
| 2 | missing artifacts | an error `SS_STAGE_FILE_MISSING`, `SS_SPEC_REFERENCE_MISSING` or `SS_SPEC_DIRECTORY_MISSING` for the item | same codes |
| 3 | status ahead of tasks | `SS_STATUS_AHEAD_OF_TASKS` for the item | status `implemented`/`verified`, `tasks.state` `read`, `checked` < `total` |
| 4 | tasks done, status early | `SS_TASKS_DONE_STATUS_EARLY` for the item | status before `implemented`, `tasks.state` `read`, `total` > 0, `checked` = `total` |
| 5 | blocked | any dependency `state` other than `satisfied` | same |
| 6 | ready for debrief | status `implemented`, `tasks.state` `read`, `checked` = `total` > 0, every dependency satisfied, no error diagnostic for the item | same |
| 7 | not started | status `planned`/`undecided`/`needs-info` without the above (reported, not part of the gap) | same |
| — | other | anything left (for example `specced` or `in-progress` with partial tasks): "in progress, no recorded inconsistency" — reported, not part of the gap | same |

Only active items are classified (`active: true`); recorded `verified` items are the numerator, not gap items. Rationale for blocked (5) before ready (6): C-07 says a dependency is satisfied only by a recorded `verified` target, so an item waiting on an unverified dependency cannot honestly be recorded `verified` yet, whatever its own evidence.

## R2 — Order of work

1. ready for debrief (raises recorded-verified without new work);
2. tasks done, status early (one confirmed status change each);
3. status ahead of tasks (needs a user decision);
4. missing artifacts (new spec work);
5. unknown status (fix the status word);
6. blocked, in dependency order: an item is worked only after its blocking dependency has been handled (or the user skips it); human-review dependencies are reported, never redirected.

Within a lane, document order (roadmap order) — deterministic and matches `spec`'s item order. Unknown status sits after missing artifacts because it never raises recorded-verified by itself; the spec (US1 scenario 2) lists the other lanes and this placement is recorded as an assumption for the owner.

## R3 — Proposals per cause

| Cause | Proposal (one change, after the user approves it) | Command |
|---|---|---|
| ready for debrief | run the debrief; if no 🎯 must-address finding and the user agrees, record `verified` | `/speckit.roadmap.debrief`, then `/speckit.roadmap.write` |
| tasks done, early | after the user confirms the work is done, advance to `implemented` | `/speckit.roadmap.write` |
| ahead of tasks | ask: unfinished work (lower the status, usually to `in-progress`) or stale boxes (finish or tick tasks) | `/speckit.roadmap.write` or `/speckit-implement` |
| missing artifacts | create the first missing artifact, or correct the status; ask which | `/speckit-specify`, `/speckit-plan`, `/speckit-tasks`, or `/speckit.roadmap.write` |
| unknown status | replace the word with one of the nine known statuses the user chooses | `/speckit.roadmap.write` |
| blocked | handle the dependency's own cause first; route human-review dependencies to `/speckit.spec-status.unblock` | as for the dependency |

These reuse the wording and rules already in reconcile, verify-readiness and unblock, so close-gap introduces no new judgment; it sequences them.

## R4 — Loop and approval protocol

- Each proposal states: item ID, cause, evidence (`path:line`), the single change, the command, and what will be re-checked.
- Nothing runs until the user approves that proposal (the approval protocol in `reference/rules.md`; Art. VIII). "Skip" moves on; "stop" ends with a summary.
- After an applied change, the skill re-runs `spec status --json` and reports the new `recorded verified X/Y` before the next proposal. A change that did not take effect (re-read shows the old state) is reported as not confirmed, never as done.
- The skill keeps no state between runs; re-running starts from the current roadmap.

## R5 — Review integration

- `reference/roles.md` "Review sections" gains `verification-gap`: counts per cause for active non-verified items, with the item IDs.
- The `pm` row leads with `verification-gap` and lists `/speckit.spec-status.close-gap` first in "First follow-ups".
- `commands/speckit.spec-status.review.md` Follow-ups names close-gap (the existing test requires every command to be reachable from review).
- Tests: `internal/extension/roles_test.go` `reviewSections` gains `verification-gap`; `helpers_test.go` `extensionCommands` gains close-gap (so `allowedLinks` and the reachability test include it).

## R6 — Six-skill touch points (measured with grep)

| File | Change |
|---|---|
| `extensions/spec-status/extension.yml` | version `0.3.0`; sixth `provides.commands` entry |
| `extensions/spec-status/README.md` | heading `v0.3.0` (manifest test checks it); "The five skills" → six, with a close-gap section and walkthrough |
| `scripts/check-extension-install.sh` | expect six skills including close-gap |
| `.github/workflows/release.yml` | archive check `n = 6`; release-notes skill list |
| `README.md` (root) | "adds five skills" → six |
| `internal/extension/*_test.go` | command list, review sections |

## R7 — Fixture

`testdata/skills/close-gap-all-causes/`: one active item per cause (R1) plus one verified item, with spec/plan/tasks files as each case needs; hand-written `EXPECTED.md` using the existing Signals grammar (`item`, `diag`, `dep`, `summary`, `outcome`), plus the expected breakdown, order of work and first proposal. The existing fixture test checks the Signals against real `spec` output.

**Measured (2026-11-15)** by building the planned fixture in a scratch directory and running the v0.3.0 development build: `recorded verified 1/8`, `blocked_count 1`, and exactly the expected signals — 003 `SS_TASKS_DONE_STATUS_EARLY` (info), 004 `SS_STATUS_AHEAD_OF_TASKS` (warning), 005 `SS_STAGE_FILE_MISSING` ×2 (error), 006 `SS_DEPENDENCY_BLOCKED` (warning), 007 `SS_SPEC_REFERENCE_MISSING` (info, so not "missing artifacts", which needs an error), 008 `SS_STATUS_UNKNOWN` (warning).

## R8 — Constraints

The `spec` binary does not change (FR-010). The skill is read-only until an approved change, and applies changes only through Spec Kit commands (C-06). `spec inventory` output is never status evidence. Tool requirement stays `spec >=0.1.2` thanks to the count fallbacks in R1.
