---

description: "Task list for upstream-format conformance (feature 002)"
---

# Tasks: Upstream-format conformance

**Input**: Design documents from `/specs/002-upstream-format-conformance/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/ (shape-inventory.md, expected-result.md, drift-report.md), quickstart.md

**Tests**: The conformance fixtures and test-only checks are this feature's deliverable (FR-004 to FR-011), so they appear as core tasks, not optional TDD tasks.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

**Authorization**: Constitution governance §2–3: this list does not authorize implementation. The owner approves before `/speckit-implement` runs.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1 to US4)
- Each description names the exact file paths

## Path Conventions

- Go test-only package: `internal/conformance/` (all files end in `_test.go`, so nothing reaches the `spec` binary)
- Pinned upstream material: `testdata/upstream/<project>/<tag>/`
- Conformance fixtures: `testdata/conformance/<name>/` with `.specify/memory/roadmap.md`, `specs/...`, and `expected.json`
- Decisions and dialect: `specs/001-spec-status/decisions.md`, `specs/001-spec-status/roadmap-format.md`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Pinned upstream copies and directory layout

- [X] T001 Create directories testdata/upstream/spec-kit/v1.2.0/templates/, testdata/upstream/speckit-roadmap/v0.1.0/templates/, testdata/conformance/, and internal/conformance/ per plan.md
- [X] T002 [P] Copy tasks-template.md, spec-template.md, and plan-template.md byte-for-byte from the installed specify-cli 1.2.0 package (specify_cli/core_pack/templates/, github/spec-kit commit 4ca46a0706a662b1edc07086ef16f021820c2db0) into testdata/upstream/spec-kit/v1.2.0/templates/, add the package's MIT LICENSE as testdata/upstream/spec-kit/v1.2.0/LICENSE, and write testdata/upstream/spec-kit/v1.2.0/NOTICE.md naming the source repository, tag, commit, and copied files
- [X] T003 [P] Copy templates/roadmap-template.md from speckit-roadmap v0.1.0 (installed copy at .specify/extensions/roadmap/templates/roadmap-template.md; confirm it matches the v0.1.0 archive by checksum if the archive is downloaded) into testdata/upstream/speckit-roadmap/v0.1.0/templates/, add the Apache-2.0 LICENSE, and write testdata/upstream/speckit-roadmap/v0.1.0/NOTICE.md
- [X] T004 [P] Create testdata/upstream/spec-kit/resolutions.md and testdata/upstream/speckit-roadmap/resolutions.md, each with an empty table of columns "shape change", "fixture", "decision ID" (format from contracts/drift-report.md)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The runner and classifier that every story's fixtures and checks depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T005 Implement the expected.json loader and comparator in internal/conformance/expected_test.go per contracts/expected-result.md: required "outcome" and "exit"; "items" compared in document order by id, status, and active, with dependencies and tasks compared when present; "summary" fields compared when present; "diagnostics" compared as a multiset of (code, item) pairs where every listed pair must occur and no other pair may occur; fatal fixtures compare "fatal" code, path, and line and omit items and summary
- [X] T006 Implement the fixture runner TestConformance in internal/conformance/conformance_test.go: discover every directory under testdata/conformance/ plus testdata/speckit-roadmap-upstream/, evaluate each with status.Evaluate using the fixture directory as root, compare with its expected.json via T005, and report the fixture name and every differing fact
- [X] T007 [P] Implement the independent line classifier in internal/conformance/shapes_test.go per research.md R5: return shape keys such as heading:N, item-heading, field:LABEL, field-style:bold, field-continuation, task-line, task-checked, task-marker:[P], task-marker:[USn], task-id:T###, front-matter, fence, table, html-comment, block-quote, html-block, ordered-list, thematic-break; normalize placeholder text so wording-only differences produce no new keys; do not call internal/status parsing functions

**Checkpoint**: Runner and classifier ready; fixtures can be added story by story

---

## Phase 3: User Story 1 - Trust counts on real Spec Kit files (Priority: P1) 🎯 MVP

**Goal**: spec status matches a hand count on roadmaps and tasks files shaped by the pinned upstream templates, including real upstream output and common hand edits

**Independent Test**: `CGO_ENABLED=0 go test ./internal/conformance -run TestConformance` passes, and each fixture's expected.json was written by hand count

### Implementation for User Story 1

- [X] T008 [P] [US1] Build testdata/conformance/template-filled/ from the pinned templates with every construct filled in: roadmap with the top-level sync comment, section comments, grouping headings, thematic breaks, footer, and items using every bold field (Description, Outcome wrapped onto an indented line, Scope (in), Scope (out), Depends on, Governed by, Addresses, Spec dir, Notes) across several statuses; one linked specs/001-sample/ with spec.md, plan.md, and a tasks.md that keeps the template front matter, top comment, quoted note, fenced bash block, numbered prose list, and task lines in all four marker forms with some checked
- [X] T009 [US1] Hand-count testdata/conformance/template-filled/ and write testdata/conformance/template-filled/expected.json with counted_by set to the counter and date (depends on T008)
- [X] T010 [P] [US1] Build testdata/conformance/hand-edits/ covering FR-006: a Depends on list wrapped onto an indented line, notes in descriptive fields, nested subtasks, a details paragraph under a task followed by nested subtasks, an indented fenced command example under a task, and a details element with blank lines around its tasks; write its hand-counted expected.json
- [X] T011 [P] [US1] Write testdata/speckit-roadmap-upstream/expected.json by hand count: 7 items, 6 active, 4 recorded verified, 2 active blockers, 43 of 43 tasks in 4 files, and the diagnostics 001 SS_SPEC_REFERENCE_MISSING, 002 SS_DEPENDENCY_BLOCKED, 002 SS_VERIFIED_DEPENDENCY_CONTRADICTION, 004 SS_DEPENDENCY_MALFORMED, 007 SS_SPEC_REFERENCE_MISSING, 007 SS_DEPENDENCY_BLOCKED (recount before writing)
- [X] T012 [P] [US1] Snapshot this repository's .specify/memory/roadmap.md and the spec.md, plan.md, and tasks.md of specs/001-spec-status/ into testdata/conformance/this-repo/ and write its hand-counted expected.json, so the fixture stays stable as the live roadmap changes
- [X] T013 [US1] Copy specs/002-upstream-format-conformance/tasks.md, unedited, into testdata/conformance/upstream-generated-002-tasks/specs/002-sample/tasks.md with a minimal roadmap that links it and empty spec.md and plan.md, and write its hand-counted expected.json (research.md R7)
- [X] T014 [US1] Add focused regression tests in internal/status/upstream_shapes_test.go for the spec edge cases: tasks-template front matter creates no tasks, and the markers [P], [US1], and [Entity1] in task text produce no link diagnostics without a matching reference definition
- [X] T015 [US1] Run CGO_ENABLED=0 go test ./internal/conformance; for any mismatch on a supported shape, fix internal/status/parse.go or internal/status/analyze.go and add a failing-then-passing test in internal/status/upstream_shapes_test.go; record "no product change needed" in specs/002-upstream-format-conformance/plan.md if every fixture passes as is

**Checkpoint**: User Story 1 is independently testable; counts on upstream-shaped files are proven by hand-counted fixtures

---

## Phase 4: User Story 2 - Know exactly which upstream versions are supported (Priority: P2)

**Goal**: Published pins and a shape inventory that is checked against the pinned templates and the fixtures

**Independent Test**: `CGO_ENABLED=0 go test ./internal/conformance -run TestPinnedInventory` passes, and README names only pins that have fixtures

### Implementation for User Story 2

- [X] T016 [US2] Implement TestPinnedInventory in internal/conformance/drift_test.go: extract shape keys from every file under testdata/upstream/*/*/templates/ with the T007 classifier, assert each key is listed in specs/002-upstream-format-conformance/contracts/shape-inventory.md, and assert each supported key appears in at least one fixture under testdata/conformance/ (SC-002)
- [X] T017 [US2] Finalize specs/002-upstream-format-conformance/contracts/shape-inventory.md: remove the draft note, align every key with the T007 classifier output, and fill the fixture coverage column (depends on T016 output)
- [X] T018 [P] [US2] Add a "Supported input" section to README.md naming github/spec-kit v1.2.0 (commit 4ca46a07) and srobroek/speckit-roadmap v0.1.0, linking the shape inventory, and stating that other versions are not claimed (constitution Article VII)
- [X] T019 [P] [US2] Restructure specs/001-spec-status/roadmap-format.md so "Upstream shape inventory (supported)" links the inventory and is separate from the existing strictness rules (visible failure for out-of-dialect input), without changing any rule (FR-014, FR-017)

**Checkpoint**: Users and maintainers can see exactly what is supported and why

---

## Phase 5: User Story 3 - Handle upstream changes deliberately (Priority: P2)

**Goal**: A drift check that reports shape changes and blocks moving a pin until each change has a fixture and a decision

**Independent Test**: `CGO_ENABLED=0 go test ./internal/conformance -run TestDriftDetectsSimulatedChange` passes and reports the simulated changes

### Implementation for User Story 3

- [X] T020 [US3] Implement the comparison and TestDrift in internal/conformance/drift_test.go per contracts/drift-report.md: flags -project and -candidate, skip when -candidate is empty, report added and removed shape keys and changed field labels or task markers with line examples, and pass only when every change is listed in testdata/upstream/PROJECT/resolutions.md with a fixture and a decision ID; read local files only
- [X] T021 [P] [US3] Create testdata/conformance-drift/tasks-template-changed.md and testdata/conformance-drift/roadmap-template-changed.md (pinned copies with a new task marker [P1] and the field label "Spec dir" renamed to "Spec directory"), and add TestDriftDetectsSimulatedChange in internal/conformance/drift_test.go asserting both changes are reported and unresolved (SC-004)
- [X] T022 [P] [US3] Write testdata/upstream/README.md with the pin-move procedure from contracts/drift-report.md, including that the drift check and downloads are development-only and spec status never fetches upstream content (FR-011)

**Checkpoint**: Upstream releases can be adopted deliberately without silent shape changes

---

## Phase 6: User Story 4 - A clear, finite acceptance bar (Priority: P3)

**Goal**: The new bar is recorded, past findings are re-classified, and 001 can be judged against it

**Independent Test**: Every finding from the acceptance reviews has a class (fixed, visible failure, limit, or blocking) under D-2026-11-07-09 with no unresolved cases

### Implementation for User Story 4

- [X] T023 [US4] Add decision D-2026-11-07-09 to the Recorded decisions table in specs/001-spec-status/decisions.md with the rule from data-model.md (blocking means a wrong result on input matching a current pin's shape inventory, or a silently misread Depends on, Spec dir, or superseded-by in any layout; other exotic input fails visibly or is a README limit), and mark D-2026-11-07-06 as superseded without deleting its text
- [X] T024 [US4] Write specs/002-upstream-format-conformance/reclassification.md listing every finding from acceptance reviews 1 to 7 with its review, ID, current behavior, and class under D-2026-11-07-09
- [X] T025 [US4] Update the Limits section of README.md so every finding classed as a limit in T024 is listed and no listed limit is now fixed or fatal
- [ ] T026 [US4] Run an independent review of feature 001 against D-2026-11-07-09 and the shape inventory, record its result as a new row in the implementation record of specs/001-spec-status/tasks.md, and if it finds zero blocking defects, move roadmap item 001 to verified in .specify/memory/roadmap.md through speckit.roadmap.write with the review cited (SC-006)

**Checkpoint**: Roadmap item 001 can be closed against a finite, published bar

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Roadmap, records, and end-to-end validation

- [X] T027 [P] Add roadmap item 007 for this feature to .specify/memory/roadmap.md through speckit.roadmap.write, with Spec dir specs/002-upstream-format-conformance and Depends on 001
- [X] T028 [P] Add a row for feature 002 to the implementation record in specs/001-spec-status/tasks.md
- [X] T029 [P] Add the conformance and drift test commands to the Commands section of CLAUDE.md
- [X] T030 Run every step in specs/002-upstream-format-conformance/quickstart.md, then CGO_ENABLED=0 go test ./... on darwin/arm64 and in Docker on linux/arm64 and linux/amd64, and confirm text and JSON output on testdata/example-roadmap is unchanged

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup; blocks every user story
- **US1 (Phase 3)**: Depends on Foundational
- **US2 (Phase 4)**: Depends on Foundational; T016 coverage checks are meaningful once US1 fixtures exist
- **US3 (Phase 5)**: Depends on Foundational (T007 classifier); independent of US1 and US2
- **US4 (Phase 6)**: Depends on US1 and US2 (the bar refers to the shape inventory and needs fixtures); T026 depends on T023 to T025
- **Polish (Phase 7)**: T027 to T029 can start any time after Setup; T030 runs last

### User Story Dependencies

- **US1 (P1)**: Starts after Foundational; no dependency on other stories
- **US2 (P2)**: Starts after Foundational; uses US1 fixtures for coverage
- **US3 (P2)**: Starts after Foundational; independent
- **US4 (P3)**: Starts after US1 and US2

### Within Each User Story

- Fixtures before their expected.json, and expected.json before running the suite
- A mismatch on a supported shape is fixed in product code before the story is done (T015)

### Parallel Opportunities

- T002, T003, T004 in Setup
- T007 alongside T005 and T006
- T008, T010, T011, T012 in US1 (separate fixture directories)
- T018 and T019 in US2
- T021 and T022 in US3
- T027, T028, T029 in Polish

---

## Parallel Example: User Story 1

```bash
# Build independent fixtures together:
Task: "Build testdata/conformance/template-filled/ from the pinned templates"
Task: "Build testdata/conformance/hand-edits/ covering FR-006"
Task: "Write testdata/speckit-roadmap-upstream/expected.json by hand count"
Task: "Snapshot this repository's roadmap into testdata/conformance/this-repo/"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: TestConformance passes on hand-counted fixtures
5. This alone answers whether real Spec Kit files are counted correctly

### Incremental Delivery

1. Setup + Foundational, then the runner works
2. US1, then counts are proven on upstream-shaped files (MVP)
3. US2, then supported versions are published and checked
4. US3, then upstream changes are caught before the pin moves
5. US4, then 001 can be judged against the finite bar

### Notes

- [P] tasks touch different files and have no dependencies on incomplete tasks
- Every expected.json is a hand count; never generate one from spec status output
- No task changes runtime behavior unless a fixture exposes a misread of a supported shape
- Commit after each task or logical group
