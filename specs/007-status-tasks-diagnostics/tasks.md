---
description: "Tasks for feature 007: status/tasks consistency diagnostics"
---

# Tasks: Status/tasks consistency diagnostics

**Input**: `specs/007-status-tasks-diagnostics/` (spec.md, plan.md, research.md, data-model.md, contracts/diagnostics.md, quickstart.md). Roadmap item 014, decision D-2026-11-14-01 (2), target v0.3.0.

**Tests**: Required. Success criteria SC-001 (every scanned case gives its expected finding), SC-002 (byte identity elsewhere) and SC-003 (no default-threshold exit change) are test-measured. Tests are written before the code and must fail first. Hand-counted expectations are edited by hand, never generated from `spec` output (CLAUDE.md).

**Organization**: US1 (ahead of tasks, P1) and US2 (done but early, P2) share one helper but are tested independently; US3 (P3) is docs and agents.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: can run in parallel (different files, no dependency on an unfinished task)
- **[Story]**: US1, US2, US3 from spec.md

---

## Phase 1: Setup

- [ ] T001 Branch from up-to-date `main` after PR #54 merges; set `.specify/feature.json` to `specs/007-status-tasks-diagnostics`; build the pre-change binary into the scratchpad and capture human (`COLUMNS=80 --color=never`) and `--json` outputs, with exit codes, for every fixture root (`testdata/example-roadmap`, `testdata/speckit-roadmap-upstream`, every `testdata/conformance/*`, `testdata/skills/*`, `testdata/conformance-drift/*`, `testdata/two-mode/*` with a roadmap) and this repository

---

## Phase 2: Foundational

**Purpose**: the shared decision helper (research R1–R4), with its full decision table tested first.

- [ ] T002 Write `internal/status/status_tasks_test.go` (must fail first) covering the data-model decision table through `Evaluate` on `project(...)` temp roadmaps: implemented 1/3 → ahead; verified 1/2 → ahead; implemented 3/3 → none; implemented 0/0 → none; in-progress 2/2 → early; planned 3/3 → early; undecided and needs-info 1/1 → early; specced 1/2 → none; planned 0/0 → none; deferred 2/2 → none; abandoned 1/3 → none; validly superseded implemented 1/3 → none; unknown status `done` 2/2 → only `SS_STATUS_UNKNOWN`; implemented with missing tasks.md → only `SS_STAGE_FILE_MISSING`; two items sharing one tasks file → one finding each. For each finding assert code, severity, exact message from contracts/diagnostics.md, `ItemID`, `Path` = roadmap, `Line` = item heading line, and `Related` = `[{tasks file, first unchecked line}]` (ahead) or `[{tasks file, first checkbox line}]` (early)
- [ ] T003 Implement `statusTasksDiag(item Item, rel string) *Diagnostic` in `internal/status/analyze.go` per data-model "Decision table" and call it in `build` right after `inspectSpec` sets `item.Tasks` (only when `item.Active`, `item.Tasks.State == "read"`); append the diagnostic to `diags`; no change to items, summary or dependency states
- [ ] T004 Run `CGO_ENABLED=0 go test ./internal/status` (T002 passes) and the full suite; expect only the conformance `template-filled` failure until T006

**Checkpoint**: both codes are emitted; nothing else changed.

---

## Phase 3: User Story 1 — CI notices a status recorded ahead of its tasks (P1) 🎯 MVP

**Goal**: `SS_STATUS_AHEAD_OF_TASKS` (warning) on the fixtures that have the mismatch; `--fail-on warning` exits 1; default exit unchanged.

**Independent Test**: `spec status --root testdata/skills/implemented-unchecked` exits 0 with one warning; with `--fail-on warning` exits 1.

- [ ] T005 [P] [US1] Add a CLI-level test in `internal/status/status_tasks_test.go` (or `internal/cli/cli_test.go`): `testdata/skills/implemented-unchecked` gives exactly one `SS_STATUS_AHEAD_OF_TASKS` for 001, `ExitCode` 0 at `error`, 1 at `warning` and `info`, 0 at `none`; the JSON diagnostic matches the contract example shape (`related` path `specs/001-import/tasks.md` with the first unchecked line)
- [ ] T006 [US1] Hand-edit `testdata/conformance/template-filled/expected.json`: add `{"code": "SS_STATUS_AHEAD_OF_TASKS", "item": "001"}` to `diagnostics` (001 is recorded `implemented` with 6 of 10 tasks checked, counted from `specs/001-sample/tasks.md`), keep `exit` 0, and extend `counted_by` with "; 2026-11-14: + SS_STATUS_AHEAD_OF_TASKS for 001 (feature 007, D-2026-11-14-01)"
- [ ] T007 [P] [US1] Add `diag SS_STATUS_AHEAD_OF_TASKS item=001` to the Signals line of `testdata/skills/implemented-unchecked/EXPECTED.md`
- [ ] T008 [US1] Run `CGO_ENABLED=0 go test ./internal/conformance ./internal/extension`; all pass

**Checkpoint**: US1 complete.

---

## Phase 4: User Story 2 — The report notices finished tasks behind a stale status (P2)

**Goal**: `SS_TASKS_DONE_STATUS_EARLY` (info) on the fixtures that have it; exit 0 at default, 1 at `--fail-on info`.

**Independent Test**: `spec status --root testdata/skills/tasks-done-status-planned --fail-on info` exits 1 with one info finding.

- [ ] T009 [P] [US2] Add a test in `internal/status/status_tasks_test.go`: `testdata/skills/tasks-done-status-planned` gives exactly one `SS_TASKS_DONE_STATUS_EARLY` for 001, `ExitCode` 0 at `error` and `warning`, 1 at `info`
- [ ] T010 [P] [US2] Add `diag SS_TASKS_DONE_STATUS_EARLY item=001` to the Signals line of `testdata/skills/tasks-done-status-planned/EXPECTED.md`
- [ ] T011 [US2] Add `diag SS_STATUS_AHEAD_OF_TASKS item=002` and `diag SS_TASKS_DONE_STATUS_EARLY item=003` to the Signals line of `testdata/skills/roles-mixed/EXPECTED.md`, and check its "diagnostics" and qa/developer view prose still states only true facts with the two added findings (edit only where it enumerates diagnostics)
- [ ] T012 [US2] Run `CGO_ENABLED=0 go test ./internal/extension`; all pass

**Checkpoint**: US1 and US2 complete.

---

## Phase 5: User Story 3 — Agents and docs use the same codes (P3)

**Goal**: contract, extension reference, reconcile skill and README describe both codes identically.

**Independent Test**: `go test ./internal/extension` passes with the codes in `reference/output.md`; the docs match contracts/diagnostics.md.

- [ ] T013 [P] [US3] Add the two registry rows from contracts/diagnostics.md to `specs/001-spec-status/diagnostics.md` (Registry table, after `SS_TASK_TARGET_MISSING`), and add "status ahead of tasks warning; tasks done with an early status info" to the severity-invariants line
- [ ] T014 [P] [US3] Add both codes to `extensions/spec-status/reference/output.md` "Diagnostic codes" with severities matching the registry, Fatal when "never", meaning, typical cause, and a fix naming `/speckit.spec-status.reconcile`, `/speckit.roadmap.write` and `/speckit-implement`; run `CGO_ENABLED=0 go test ./internal/extension`
- [ ] T015 [P] [US3] In `extensions/spec-status/commands/speckit.spec-status.reconcile.md`, add to the **status-ahead-of-tasks** and **tasks-done-status-early** bullets: "reported by `spec` as `SS_STATUS_AHEAD_OF_TASKS` / `SS_TASKS_DONE_STATUS_EARLY` (v0.3.0 and later); with an older `spec`, derive it from the counts as described"; keep the propose-only, per-change-approval wording unchanged; run `CGO_ENABLED=0 go test ./internal/extension`
- [ ] T016 [US3] Update `README.md`: in "When the first run looks wrong" (or the findings guidance), one row each for `SS_STATUS_AHEAD_OF_TASKS` (recorded implemented/verified while tasks are unchecked; warning, fails only with `--fail-on warning` or `info`) and `SS_TASKS_DONE_STATUS_EARLY` (every task checked, status still before implemented; info); keep the README tone

---

## Phase 6: Polish and cross-cutting

- [ ] T017 Rebuild and compare every T001 output byte for byte: only `testdata/conformance/template-filled`, `testdata/skills/implemented-unchecked`, `testdata/skills/roles-mixed` and `testdata/skills/tasks-done-status-planned` differ, only by the added diagnostics (and the diagnostics block in human text); every exit code is unchanged; this repository is byte-identical apart from 007's own checked tasks. Record the result in the PR description
- [ ] T018 Run `CGO_ENABLED=0 go vet ./...`, `gofmt -l .`, `go mod tidy -diff`, the full suite on darwin/arm64 and on linux/amd64 (container), and `spec status --fail-on error` on this repository; record platforms and Go versions in the PR description
- [ ] T019 Add a row for feature 007 to the "Implementation record" table in `specs/001-spec-status/tasks.md` (implemented, executed by the implementer, not independently verified; platforms actually run)
- [ ] T020 Draft the v0.3.0 release-note entry in the PR description: two new codes; projects on `--fail-on warning`/`info` may start exiting 1 on real mismatches
- [ ] T021 Doomer review on the source change (prepare before editing, submit, answer every Quick lead, finish) before handing off the PR
- [ ] T022 Propose roadmap 014 `specced` → `implemented` through `/speckit.roadmap.write` after the owner approves; `verified` only after the v0.3.0 release checks

---

## Dependencies and execution order

- T001 → everything.
- T002 → T003 → T004 (Phase 2 blocks both stories).
- US1 (T005–T008) and US2 (T009–T012) depend on Phase 2 and are otherwise independent; T011 edits the same file family as T007/T010 only in a different fixture.
- US3 (T013–T016) can start after T003 fixes the messages; T014 must agree with T013's severities.
- T017–T020 after all stories; T021 before handoff; T022 last.

## Parallel examples

- After T004: T005, T007, T009, T010 in parallel (different files).
- In US3: T013, T014, T015 in parallel.

## Implementation strategy

1. **MVP = Phase 2 + US1**: the warning is the value for CI; it ships even if US2/US3 slip.
2. US2 and US3 are small and land in the same PR.
3. One PR against `main` (no stacking); owner merges; release through the owner when 013, 015, 016 and 011 are in.
