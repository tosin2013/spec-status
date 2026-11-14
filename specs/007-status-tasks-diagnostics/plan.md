# Implementation Plan: Status/tasks consistency diagnostics

**Branch**: `007-status-tasks-specify` | **Date**: 2026-11-14 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/007-status-tasks-diagnostics/spec.md` (roadmap item 014, decision D-2026-11-14-01 (2), target v0.3.0).

## Summary

Add two report diagnostics to `spec status`: `SS_STATUS_AHEAD_OF_TASKS` (warning) when an active `implemented`/`verified` item's read tasks file has unchecked tasks, and `SS_TASKS_DONE_STATUS_EARLY` (info) when an active item before `implemented` has a read tasks file with every task checked. One helper in `build` (`internal/status/analyze.go`) decides from values already computed; the finding sits at the item heading with the tasks file's first relevant checkbox as related evidence (research R1–R4). No schema, summary or human-text layout change. Four fixtures gain findings; one hand-counted conformance expectation is edited by hand (R5).

## Technical Context

**Language/Version**: Go 1.26, `CGO_ENABLED=0`

**Primary Dependencies**: standard library only for this change (status package analysis)

**Storage**: none

**Testing**: `go test ./...`; new unit tests in `internal/status`; conformance (`internal/conformance`) and extension (`internal/extension`) suites re-run; byte comparison of all fixture roots against the pre-change build

**Target Platform**: linux/amd64, darwin/arm64

**Project Type**: CLI

**Performance Goals**: constant work per item; no new reads

**Constraints**: Art. I (mechanical facts only), Art. II (no side effects), Art. IV (historical items excluded), Art. V (schema 1.0, default threshold unchanged); C-07 (no duplicate of stage-matrix findings)

**Scale/Scope**: about 40 lines of Go plus tests; two registry rows; doc and fixture updates

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Article / constraint | Check | Result |
|---|---|---|
| I Truthful evidence (C-01) | Reports a mechanical disagreement between two recorded facts; messages claim nothing about done or verified; 0/0 is not called a mismatch | Pass |
| II Deterministic, standalone (C-02) | Arithmetic on values already read; fixed messages; sorted with `diagLess`; no I/O added | Pass |
| III Narrow input | No grammar change | Pass |
| IV History (C-04) | Only active items compared; superseded/deferred/abandoned excluded | Pass |
| V Compatibility (C-05) | Additive codes in schema 1.0; default `--fail-on error` exit unchanged for every fixture; release notes list the change for warning/info thresholds | Pass |
| VI–VIII | No timing, distribution or mutation change; skills still propose only | Pass |
| C-07 stage matrix | Missing/unlinked/skipped tasks stay with existing findings; new codes need a read tasks file | Pass |
| Dependency allowlist | No new import | Pass |

Post-design re-check: no change; no violations.

## Project Structure

### Documentation (this feature)

```text
specs/007-status-tasks-diagnostics/
├── spec.md
├── plan.md              # this file
├── research.md          # R1–R8
├── data-model.md        # decision table
├── quickstart.md
├── contracts/
│   └── diagnostics.md   # registry rows, messages, exit table, fixtures
├── checklists/requirements.md
└── tasks.md             # /speckit-tasks
```

### Source Code (repository root)

```text
internal/status/
├── analyze.go                 # + statusTasksDiag, called after inspectSpec in build
└── status_tasks_test.go       # new: decision-table unit tests

testdata/conformance/template-filled/expected.json   # hand-edited: + 001 SS_STATUS_AHEAD_OF_TASKS
testdata/skills/{implemented-unchecked,roles-mixed,tasks-done-status-planned}/EXPECTED.md   # + diag signals

specs/001-spec-status/diagnostics.md                  # + 2 registry rows
extensions/spec-status/reference/output.md            # + 2 codes
extensions/spec-status/commands/speckit.spec-status.reconcile.md   # cite codes, keep count fallback
README.md                                             # findings guidance rows
specs/001-spec-status/tasks.md                        # implementation record row
```

**Structure Decision**: existing packages only; the change is one helper in the analysis step (R1).

## Phasing for tasks

1. Baseline: capture all fixture-root outputs from the pre-change build.
2. Tests first: decision-table unit tests; hand-edit the conformance expectation (fails until implemented).
3. Implement `statusTasksDiag` and the call in `build`.
4. Byte comparison: only the four expected roots differ, only by the new diagnostics.
5. Docs and agents: registry, output reference, reconcile wording, skills fixture signals, README, implementation record, roadmap status proposal.

## Risks

- A project already using `--fail-on warning` will start failing on genuine ahead-of-tasks mismatches after upgrading. Intended; called out in release notes.
- Agents pinned to an older `spec` will not see the codes; the reconcile skill keeps deriving the cases from counts (R6).

## Complexity Tracking

None.
