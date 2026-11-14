# Contract: two new `spec status` diagnostics

Amends `specs/001-spec-status/diagnostics.md` at implementation. Decision D-2026-11-14-01 (2); rules in [../data-model.md](../data-model.md).

## Registry rows

| Code | Severity / class | Trigger and source | Authority status |
|---|---|---|---|
| SS_STATUS_AHEAD_OF_TASKS | warning, report | Active item recorded `implemented` or `verified` whose linked tasks file was read with checked < total; source is the item heading line; related is the first unchecked checkbox | D-2026-11-14-01 (2), feature 007 |
| SS_TASKS_DONE_STATUS_EARLY | info, report | Active item recorded `undecided`, `needs-info`, `planned`, `specced` or `in-progress` whose linked tasks file was read with total > 0 and every task checked; source is the item heading line; related is the first checkbox | D-2026-11-14-01 (2), feature 007 |

## Messages (exact)

- `SS_STATUS_AHEAD_OF_TASKS`: `Recorded status is implemented or verified, but the linked tasks file has unchecked tasks.`
- `SS_TASKS_DONE_STATUS_EARLY`: `Every task in the linked tasks file is checked, but the recorded status is before implemented.`

## Not reported

Historical items (deferred, abandoned, validly superseded); unknown status tokens; tasks files that are missing, not linked or skipped (the stage-matrix and confinement findings already cover them); tasks files with zero checkboxes.

## JSON example (one diagnostic, schema 1.0, unchanged shape)

```json
{"code":"SS_STATUS_AHEAD_OF_TASKS","severity":"warning","message":"Recorded status is implemented or verified, but the linked tasks file has unchecked tasks.","source":{"path":".specify/memory/roadmap.md","line":5},"item_id":"001","related":[{"path":"specs/001-import/tasks.md","line":4}]}
```

## Exit behavior

| `--fail-on` | ahead (warning) only | early (info) only |
|---|---|---|
| `error` (default) | 0 | 0 |
| `warning` | 1 | 0 |
| `info` | 1 | 1 |
| `none` | 0 | 0 |

## Fixtures

| Fixture | Expected new finding |
|---|---|
| `testdata/conformance/template-filled` | 001 ahead (hand-edited `expected.json`) |
| `testdata/skills/implemented-unchecked` | 001 ahead |
| `testdata/skills/roles-mixed` | 002 ahead, 003 early |
| `testdata/skills/tasks-done-status-planned` | 001 early |
| unit tests (new) | implemented 3/3 → none; implemented 0/0 → none; verified 1/2 → ahead; in-progress 2/2 → early; planned 1/2 → none; deferred 2/2 → none; validly superseded implemented 1/3 → none; unknown status 2/2 → none (only `SS_STATUS_UNKNOWN`); implemented with missing tasks.md → only `SS_STAGE_FILE_MISSING`; two items sharing one tasks file → one finding each |
