# Research: Status/tasks consistency diagnostics

Phase 0 for `specs/007-status-tasks-diagnostics/plan.md`. Code references are to `main` at `bda37da` plus this branch (docs only).

## R1 — Where the comparison runs

- **Decision:** in `build` (`internal/status/analyze.go`), inside the per-item loop, immediately after `inspectSpec` returns the item's `ItemTasks` and the `Active` flag is set. One small helper, `statusTasksDiag(item Item, rel string) *Diagnostic`, decides.
- **Rationale:** the inputs it needs (active, known status, `Tasks.State`, `Checked`, `Total`, `Source`, `Checks`) all exist at that point; no new pass and no parser change. Diagnostics are sorted at the end of `build` with `diagLess`, so insertion order does not matter.
- **Alternatives:** a post-pass over `rep.Items` (rejected: same result, an extra loop); inside `inspectSpec` (rejected: it does not know `Active`, which depends on supersession).

## R2 — Exact rules

| Condition (all required) | Finding |
|---|---|
| `Active`; status `implemented` or `verified`; `Tasks.State == "read"`; `Checked < Total` | `SS_STATUS_AHEAD_OF_TASKS`, warning |
| `Active`; status in `undecided`, `needs-info`, `planned`, `specced`, `in-progress`; `Tasks.State == "read"`; `Total > 0`; `Checked == Total` | `SS_TASKS_DONE_STATUS_EARLY`, info |

`Active` already excludes `deferred`, `abandoned` and validly superseded items. Unknown statuses fail both status tests. `Checked < Total` is false for 0/0, so empty task files never trigger the first rule; the second requires `Total > 0`.

## R3 — Location and related evidence

- **Source:** roadmap path and the item heading line (`it.Line`), like `SS_STATUS_UNKNOWN`, because the status being compared is recorded on that line.
- **Related:** one reference to the tasks file the counts came from (`Tasks.Source`, which already resolves shared files to the first path), at a real checkbox line from `Tasks.Checks`: the **first unchecked** checkbox for "ahead", the **first** checkbox for "early". Checkbox lines are always recorded by `parseTasks` (`parse.go:781`), so no line is fabricated (json-contract: "do not fabricate line 1").
- `related` is not shown in human text today (only the JSON); the human diagnostics list shows path:line, item and message, which is enough to find both files. No human-text change is planned.

## R4 — Messages (fixed text)

- `SS_STATUS_AHEAD_OF_TASKS`: `Recorded status is implemented or verified, but the linked tasks file has unchecked tasks.`
- `SS_TASKS_DONE_STATUS_EARLY`: `Every task in the linked tasks file is checked, but the recorded status is before implemented.`
- No counts (determinism, and counts are already in the item row and JSON). No claim about whether work is done or verified (C-01).

## R5 — Impact, measured

A scan with the current build over all 20 fixture roots (spec Context table):

| Root | Gains | Expectation change |
|---|---|---|
| `testdata/conformance/template-filled` | 001 `SS_STATUS_AHEAD_OF_TASKS` | Hand-edit `expected.json`: add `{"code":"SS_STATUS_AHEAD_OF_TASKS","item":"001"}` and a `counted_by` note dated 2026-11-14 citing feature 007. The conformance test compares (item, code) counts exactly, so it fails until then. `exit` stays 0 (warning < error). |
| `testdata/skills/implemented-unchecked` | 001 ahead | Add `diag SS_STATUS_AHEAD_OF_TASKS item=001` to Signals |
| `testdata/skills/roles-mixed` | 002 ahead, 003 early | Add both `diag` clauses |
| `testdata/skills/tasks-done-status-planned` | 001 early | Add `diag SS_TASKS_DONE_STATUS_EARLY item=001` |
| all other roots, this repository | nothing | byte-identical |

Skills `diag` clauses are "must contain" checks (`internal/extension/fixtures_test.go`), so the fixtures would pass without the edits; the edits make the expected signals complete (FR-009). No default-threshold exit code changes anywhere (SC-003).

## R6 — Extension reference and skills

- `extensions/spec-status/reference/output.md` must list both codes with the contract severities (the reference test enforces it) and a fix naming allowed commands: `/speckit.spec-status.reconcile`, `/speckit.roadmap.write`, `/speckit-implement`.
- `commands/speckit.spec-status.reconcile.md`: the two mismatch bullets gain "reported by `spec` as `SS_…` since v0.3.0; older builds: derive it from the counts as before", so the skill works with both binaries (the extension requires `spec >= 0.1.2`). Proposals still go through Spec Kit commands with per-change approval.
- The extension version bump to 0.3.0 happens with item 013, not here.

## R7 — Docs and release

- `specs/001-spec-status/diagnostics.md`: two registry rows, authority "D-2026-11-14-01 (2), feature 007"; severity invariants line updated.
- README findings guidance: a short row for each code; Limits unchanged.
- v0.3.0 release notes: the two codes; projects on `--fail-on warning`/`info` may start exiting 1 on real mismatches.

## R8 — Constitution and constraints

No new dependency, I/O, or side effect; the comparison is arithmetic on already-read values. C-01 holds (mechanical inconsistency, no inference), C-04 (historical items excluded), C-05 (schema 1.0, default threshold unchanged), C-07 (no duplicate of stage-matrix findings).
