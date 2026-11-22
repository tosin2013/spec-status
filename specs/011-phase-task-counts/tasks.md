# Tasks: Per-phase task counts

**Input**: Design documents from `specs/011-phase-task-counts/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/phases-json.md, quickstart.md

**Tests**: Required. FR-009 and FR-010 call for hand-counted fixtures, and SC-002 calls for a byte comparison against the pre-change build.

**Organization**: Tasks are grouped by user story. US1 and US2 are both P1 and share the parser change in Phase 2.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story the task belongs to (US1, US2, US3)

## Phase 1: Setup

- [ ] T001 Build the pre-change baseline: `CGO_ENABLED=0 go build -o $TMPDIR/spec-base ./cmd/spec` from `main`. Save the text output (`--color=never`, plus `--explain`) and the `--json` output for every fixture root under `testdata/` (each directory with `.specify/`) into `$TMPDIR/phases-baseline/`, for SC-002.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The model and parser change that both P1 stories use.

- [ ] T002 Add the `Phase` type (`Heading string; Line int; Story *int; Checked, Total int`) and `ItemTasks.Phases []Phase` in `internal/status/model.go`. Doc comment: "recorded checkboxes under a `## Phase` heading; says nothing about status".
- [ ] T003 Add `Phases []Phase` to `taskFile`, and recognize phases in `parseTasks` in `internal/status/parse.go`:
  - In the existing branch `if _, _, _, isHeading := heading(line); isHeading && width <= 3` (around line 714), any level 1 or 2 heading closes the open phase.
  - A level-2 heading whose text, after a CommonMark closing `#` sequence is removed and the text is trimmed, is `Phase` or starts with `Phase` followed by a space, a tab or `:` (case-sensitive), opens a new phase at line `i+1` (research R2, R3).
  - Add a helper `phaseHeading(level int, text string) (heading string, ok bool)`.
  - Add a helper `phaseStory(text string) *int`: the first `User Story ` followed by 1–9 digits and then a non-digit or the end of the text, otherwise nil (FR-006).
- [ ] T004 In `parseTasks`, where `tf.Checks` gets a `TaskCheck` (around line 847), increment the open phase's `Total`, and its `Checked` if the box is checked (R4), in `internal/status/parse.go`.
- [ ] T005 Copy `tf.Phases` into `ItemTasks.Phases` for read files, next to the `Checks` copy (around line 604), in `internal/status/analyze.go`. Use a fresh slice per item so shared files never alias.
- [ ] T006 Add `phaseView` (fields in this order: `heading`, `line`, `story` as `*int`, `checked`, `total`) and `itemTaskView.Phases *[]phaseView` (`json:"phases"`, after `state`) in `internal/status/render.go`. `itemToView` sets a non-nil slice (`[]` when empty) only when the state is `read`, and leaves it nil (`null`) otherwise (R5).

**Checkpoint**: `go build ./...` passes. Every existing test passes except the JSON byte expectations that now include `phases`.

---

## Phase 3: User Story 1 - A gate reads one phase's progress from the JSON (Priority: P1) 🎯 MVP

**Goal**: The phases of a Spec Kit-generated tasks.md are reported with the right headings, lines, stories and counts.

**Independent Test**: `spec status --json --root testdata/conformance/upstream-generated-002-tasks | jq '.items[0].tasks.phases'` matches the hand count.

### Tests for User Story 1

- [ ] T007 [P] [US1] Extend `expectedTasks` with an optional `Phases *[]expectedPhase` (`heading`, `line`, `story`, `checked`, `total`), and make `compareItem` compare them in order (count, then each field), in `internal/conformance/expected_test.go`.
- [ ] T008 [US1] Hand-count the seven phases of `testdata/conformance/upstream-generated-002-tasks/specs/002-sample/tasks.md` (headings at lines 33, 44, 58, 79, 96, 112, 129; stories 1–4 at 58–112) by reading the file, not `spec` output. Add them to `testdata/conformance/upstream-generated-002-tasks/expected.json`, and note the count method in `counted_by`.
- [ ] T009 [P] [US1] Unit tests in `internal/status/phases_test.go`:
  - generated layout with level-3 subsections inside a phase;
  - `## Dependencies` ends the last phase;
  - a level-1 heading ends a phase;
  - closing `#`s are removed;
  - `Phase N: Polish` qualifies, `Phases overview` and `phase 1` do not;
  - story extraction (`User Story 12` gives 12, a 10-digit number gives nil, no story gives nil);
  - duplicate headings give two entries;
  - the per-phase sums equal the in-phase share of the totals.

### Implementation for User Story 1

- [ ] T010 [US1] Run `go test ./internal/status ./internal/conformance`, and fix `internal/status/parse.go` until T007–T009 pass.

**Checkpoint**: US1 is independently testable on the real generated fixture.

---

## Phase 4: User Story 2 - "No phases" and "unknown" stay different (Priority: P1)

**Goal**: `[]` for a read file without phases, `null` for unread files. Hidden headings are ignored.

**Independent Test**: `spec status --json --root testdata/phases | jq '.items[] | {id, s: .tasks.state, p: .tasks.phases}'` matches `testdata/phases/README.md`.

### Tests for User Story 2

- [ ] T011 [P] [US2] Create the fixture `testdata/phases/` with `.specify/memory/roadmap.md` and `specs/*/tasks.md`. It covers:
  - a read file with no phase headings;
  - a fenced ```` ``` ```` block containing `## Phase 9` and a checkbox;
  - a `## Phase` inside an HTML comment;
  - tasks before the first phase;
  - tasks after `## Dependencies`;
  - an empty phase (0/0);
  - an item with a missing tasks file;
  - an item with no linked tasks file;
  - two items sharing one tasks file.

  Write a `testdata/phases/README.md` that lists the hand-counted expected `phases` per item.
- [ ] T012 [US2] Add tests in `internal/status/phases_test.go` that evaluate `testdata/phases` and assert the README's hand counts, including JSON `"phases":[]` and `"phases":null` byte fragments from `status.JSON`.

### Implementation for User Story 2

- [ ] T013 [US2] Run the tests, and fix `internal/status/parse.go` and `internal/status/render.go` until T011–T012 pass.

**Checkpoint**: US1 and US2 both pass.

---

## Phase 5: User Story 3 - Nothing else changes (Priority: P2)

**Goal**: Text, exit codes and the other JSON bytes are identical to v0.3.0 behavior.

**Independent Test**: Compare against the T001 baseline.

- [ ] T014 [US3] Rebuild `spec` and compare with `$TMPDIR/phases-baseline/` for every fixture root:
  - the text output and exit codes are byte-identical;
  - the JSON is identical after `jq -S 'del(.items[].tasks.phases)'`.

  Record the root count and result for the implementation record.
- [ ] T015 [P] [US3] Hand-add `phases` to the other conformance fixtures with read tasks files, counted from each tasks.md:
  - `testdata/conformance/hand-edits/expected.json`
  - `testdata/conformance/template-filled/expected.json`
  - `testdata/conformance/this-repo/expected.json`
  - `testdata/conformance/this-repo-1.3.4/expected.json`
- [ ] T016 [US3] Run `go test ./internal/cli -run TestReadmeSamples`, and update any labeled JSON sample in `README.md` or `docs/*.md` that now shows `phases`. Run the command for the real output; do not hand-edit samples.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T017 [P] Document the `phases` field and its null/empty rule in the item `tasks` row of `specs/001-spec-status/json-contract.md`, linking `specs/011-phase-task-counts/contracts/phases-json.md`.
- [ ] T018 [P] Add the phase grammar to the tasks section of `specs/001-spec-status/roadmap-format.md`.
- [ ] T019 [P] Add to `docs/limits.md`: setext headings neither open nor end a phase, and only `## Phase…` ATX headings are phases.
- [ ] T020 [P] Add a row "Per-phase task counts (feature 011, roadmap 017)" to the implementation record in `specs/001-spec-status/tasks.md`, with the T014 result. Mark it "implemented by the implementer, not independently verified".
- [ ] T021 Run `CGO_ENABLED=0 go test ./...`, gofmt and `go vet ./...`, and the quickstart steps in `specs/011-phase-task-counts/quickstart.md`.

---

## Dependencies & Execution Order

### Phase Dependencies

- Setup (T001) comes first, because the baseline must be built before any code changes.
- Foundational (T002–T006) blocks every story.
- US1 (Phase 3) and US2 (Phase 4) can run in either order once Phase 2 is done.
- US3 (Phase 5) runs after US1 and US2, because it compares the final behavior.
- Polish runs after US3.

### Parallel Opportunities

- T007 and T009 are in different files.
- T011 is independent of T007–T009.
- T015 and T017–T020 are documentation and fixture files.

## Parallel Example: User Story 1

```text
T007 internal/conformance/expected_test.go
T009 internal/status/phases_test.go
```

## Implementation Strategy

1. T001 baseline, then T002–T006, then US1 (the real generated fixture). That is the MVP for issue #66.
2. US2 locks the null/empty rule and the hidden-heading cases.
3. US3 proves that nothing else moved. Then the docs.
4. Doomer review on the source changes, then one PR against `main`. The owner merges.
