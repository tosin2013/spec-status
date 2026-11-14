# Feature Specification: Status/tasks consistency diagnostics

**Feature Branch**: `007-status-tasks-specify`

**Created**: 2026-11-14

**Status**: Draft

**Input**: Roadmap item 014, "Status/tasks consistency diagnostics", chosen for v0.3.0 by the owner in decision D-2026-11-14-01 (2): "new mechanical findings when a recorded status disagrees with the item's tasks … Codes, severities and exact rules are set in the feature spec. Reported facts only; never infers verification (C-01)." Item 013 (the close-gap skill) depends on it.

## Context

Today `spec status` shows a recorded status and a checkbox count side by side, but says nothing when they disagree. An item recorded `implemented` with 1 of 3 tasks checked reports "diagnostics none" and exit 0 at every threshold. Only the agent skills' reconcile command notices, so CI cannot act on it.

A scan of the existing fixtures (2026-11-14) shows where the disagreement appears and where existing findings already cover it:

| Case | Example | Covered today? |
|---|---|---|
| `implemented`/`verified`, tasks read, some unchecked | `testdata/conformance/template-filled` 001 (6/10), `testdata/skills/implemented-unchecked` 001 (1/3), `testdata/skills/roles-mixed` 002 (1/2) | No |
| Status before `implemented`, tasks read, all checked | `testdata/skills/tasks-done-status-planned` 001 (3/3), `testdata/skills/roles-mixed` 003 (2/2) | No |
| `implemented`/`verified`, no readable tasks file | `testdata/skills/verified-tasks-missing` 001, `testdata/speckit-roadmap-upstream` 001 | Yes: `SS_STAGE_FILE_MISSING`, `SS_SPEC_REFERENCE_MISSING`, `SS_SPEC_DIRECTORY_MISSING`, `SS_PATH_UNSAFE` |
| `implemented`/`verified`, tasks file read with no checkboxes (0/0) | this repository's 001–003 (task files that are tables) | No disagreement to report |

The two uncovered cases get new codes. The agent skills already use these names for them: "status-ahead-of-tasks" and "tasks-done-status-early".

## User Scenarios & Testing *(mandatory)*

### User Story 1 - CI notices a status recorded ahead of its tasks (Priority: P1)

A maintainer runs `spec status --fail-on warning` in CI. An item is recorded `implemented`, but its linked tasks file still has unchecked tasks. The run reports a located warning naming the item and the tasks file, and exits 1, so the mismatch is caught before anyone treats the item as done.

**Why this priority**: This is the case that most often hides unfinished work behind a confident label, and the outcome of roadmap item 014 is that CI and the human report see what the reconcile skill sees.

**Independent Test**: Run `spec status --json` and `--fail-on warning` on a fixture with one `implemented` item at 1/3 tasks; check for exactly one new warning with the item ID, the roadmap line and the tasks file, and exit 1. With the default threshold, exit stays 0.

**Acceptance Scenarios**:

1. **Given** an active item recorded `implemented` or `verified` whose linked tasks file was read with checked < total, **When** `spec status` runs, **Then** it reports `SS_STATUS_AHEAD_OF_TASKS` (warning) for that item, located at the item's heading line in the roadmap, with the tasks file as related evidence.
2. **Given** that finding, **When** run with the default threshold (`error`), **Then** the exit code is unchanged from today (0 unless something else is an error); **When** run with `--fail-on warning` or `info`, **Then** it exits 1.
3. **Given** an `implemented` item whose tasks are all checked, **When** `spec status` runs, **Then** no new finding is reported.
4. **Given** an `implemented` item whose tasks file was read with zero checkboxes, **When** `spec status` runs, **Then** no new finding is reported (there is no checkbox evidence either way).
5. **Given** an `implemented` or `verified` item with no readable tasks file, **When** `spec status` runs, **Then** only the existing findings for the missing or skipped file are reported; the new code is not added on top.

---

### User Story 2 - The report notices finished tasks behind a stale status (Priority: P2)

A developer checks off every task of an item but forgets to update the roadmap, which still says `planned`. `spec status` reports an informational finding so the stale status is visible, without failing CI at the default threshold.

**Why this priority**: The reverse case understates progress rather than overstating it, so it matters less for trust, but it is the most common reason a recorded-verified percentage looks low, which is what the close-gap skill (013) works on.

**Independent Test**: Run `spec status --json` on a fixture with one `planned` item at 3/3 tasks; check for one `SS_TASKS_DONE_STATUS_EARLY` info finding and exit 0 at the default threshold, 1 at `--fail-on info`.

**Acceptance Scenarios**:

1. **Given** an active item whose status is `undecided`, `needs-info`, `planned`, `specced` or `in-progress`, and whose linked tasks file was read with total > 0 and every task checked, **When** `spec status` runs, **Then** it reports `SS_TASKS_DONE_STATUS_EARLY` (info) at the item's heading line, with the tasks file as related evidence.
2. **Given** the same item with one task unchecked, **When** `spec status` runs, **Then** no new finding is reported.
3. **Given** an item recorded `implemented` with all tasks checked, **When** `spec status` runs, **Then** no new finding is reported (implemented is not "early"; whether it is verified is never inferred).

---

### User Story 3 - Agents and docs use the same codes (Priority: P3)

The agent skills' explain and reconcile commands, and the README findings table, describe the two codes, so a person or agent reading `spec status` output gets the same meaning and fix everywhere.

**Why this priority**: It completes the loop for 013, but the CLI behavior in US1 and US2 stands on its own.

**Independent Test**: The extension's reference test passes with the new codes listed; the skills fixtures that contain these mismatches name the new codes in their expected signals.

**Acceptance Scenarios**:

1. **Given** the diagnostics contract lists the two codes, **When** the extension tests run, **Then** `reference/output.md` lists both with the same severities and a fix that names an allowed command.
2. **Given** the skills fixtures `implemented-unchecked`, `tasks-done-status-planned` and `roles-mixed`, **When** `spec status` runs on them, **Then** their expected signals include the new findings.

---

### Edge Cases

- **Historical items** (deferred, abandoned, validly superseded) are not compared: they earn no active completion credit and their recorded status is history (C-04).
- **Unknown status tokens** are not compared; they already get `SS_STATUS_UNKNOWN`.
- **Two items sharing one tasks file**: each item is compared with that file's counts and may get its own finding.
- **Skipped or unsafe tasks files**: not compared; the existing skip finding stands.
- **An item recorded `verified` with unchecked tasks** gets `SS_STATUS_AHEAD_OF_TASKS` in addition to any `SS_VERIFIED_DEPENDENCY_CONTRADICTION`; the two describe different facts.
- **Filters and `--explain`** change only which item rows are shown; the new findings are always in the diagnostics list, the JSON and the exit code.
- **Determinism**: findings sort with the existing diagnostic order; no counts or timestamps go into message text.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `spec status` MUST report `SS_STATUS_AHEAD_OF_TASKS`, severity warning, for each active item recorded `implemented` or `verified` whose linked tasks file was read and has at least one unchecked task (checked < total).
- **FR-002**: `spec status` MUST report `SS_TASKS_DONE_STATUS_EARLY`, severity info, for each active item recorded `undecided`, `needs-info`, `planned`, `specced` or `in-progress` whose linked tasks file was read with total > 0 and checked = total.
- **FR-003**: Both findings MUST name the item ID, be located at the item's heading line in the roadmap, and carry the linked tasks file as related evidence. Messages MUST be fixed text without counts, stating the recorded status class and the checkbox fact, and MUST NOT claim the item is or is not done or verified.
- **FR-004**: The two findings MUST NOT be reported for historical items, unknown statuses, items whose tasks file is missing, not linked or skipped, or tasks files with zero checkboxes.
- **FR-005**: The findings MUST NOT change any summary count, denominator, item field, or dependency state. They only add diagnostics (and therefore can change the exit code at thresholds at or below their severity).
- **FR-006**: JSON MUST stay schema 1.0; the codes are additive diagnostics in the existing shape.
- **FR-007**: Output for inputs without either mismatch MUST be byte-identical to the previous build, in human text and JSON.
- **FR-008**: The diagnostics contract, README findings table and the extension's `reference/output.md` MUST describe both codes (meaning, severity, typical cause, fix). The hand-counted conformance fixture `template-filled` MUST be updated by hand with the new finding and a note of why its expectation changed.
- **FR-009**: The skills fixtures that contain these mismatches MUST list the new findings in their expected signals; the reconcile command MAY cite the codes, but MUST keep proposing changes only through Spec Kit commands with per-change approval.
- **FR-010**: No runtime side effects, new dependencies or inferred verification are added (Art. I–II).

### Key Entities

- **Status/tasks comparison**: for one active item with a known status and a read tasks file: the recorded status class (before implemented, implemented or verified) and the checkbox fact (some unchecked, all checked, no checkboxes).
- **Finding**: code, severity, item ID, roadmap path and heading line, related tasks file, fixed message.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every fixture case in the scan table above produces exactly the finding (or absence of one) listed for it.
- **SC-002**: Outputs for all fixture roots without a mismatch, and this repository's own report, are byte-identical to the previous build.
- **SC-003**: At the default threshold, no existing fixture's exit code changes.
- **SC-004**: A CI job using `--fail-on warning` fails on a roadmap with an `implemented` item that has unchecked tasks, and passes once the tasks are checked or the status is lowered.
- **SC-005**: Repeated runs on the same input produce identical output in 100% of runs.

## Assumptions

- **Severities**: warning for status ahead of tasks (overstated progress, worth failing a strict CI), info for tasks done with an early status (understated progress). Neither is an error, so default-threshold CI does not start failing on upgrade. The owner can raise either in `/speckit-clarify`.
- **0/0 is not a mismatch**: a tasks file with no checkboxes gives no checkbox evidence. This keeps this repository's verified items 001–003, whose task files are tables, free of findings.
- **"Before implemented"** means the five lifecycle statuses before `implemented`; `deferred` and `abandoned` are historical and excluded.
- **No double reporting**: a missing, unlinked or skipped tasks file is already an error or skip finding under the stage matrix (C-07), so the new codes apply only when the tasks file was read.
- **Release notes**: v0.3.0 lists the two codes; a project already running `--fail-on warning` or `info` may start exiting 1 on real mismatches, which is the intended outcome.
- **Out of scope**: proposing or applying status changes (the skills' job, with approval); comparing `spec.md`/`plan.md` content; any inference about verification; changing the close-gap skill (item 013).
