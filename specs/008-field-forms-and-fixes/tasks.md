---
description: "Tasks for feature 008: real-world field forms and small defect fixes"
---

# Tasks: Real-world field forms and small defect fixes

**Input**: `specs/008-field-forms-and-fixes/` (spec.md, plan.md, research.md, data-model.md, contracts/changes.md, quickstart.md). Roadmap items 015 and 016, decision D-2026-11-14-01 (3), (4), target v0.3.0.

**Tests**: Required. Success criteria SC-001 (one test per row A–H), SC-002 (raw template exits 0), SC-003 (JSON byte identity) and SC-004 (no real ID dropped) are test-measured. Tests are written before the code and must fail first. Fixture expectations are written by hand from contracts/changes.md §6, never generated from `spec` output (CLAUDE.md).

**Organization**: US1 (Superseded by) and US2 (placeholders) are P1; US3 (duplicate IDs) and US4 (human text) are P2; US5 (CLI edge cases) is P3. They touch different code paths and are independently testable.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: can run in parallel (different files, no dependency on an unfinished task)
- **[Story]**: US1–US5 from spec.md

---

## Phase 1: Setup

- [X] T001 Branch from up-to-date `main` after PR #57 merges (or from its branch, opening the implementation PR against `main`); set `.specify/feature.json` to `specs/008-field-forms-and-fixes`; build the pre-change binary into the scratchpad and capture human (`COLUMNS=80 --color=never`) and `--json` outputs, with exit codes, for every fixture root and this repository (the 007 `snap.sh` root list)
- [X] T002 [P] Create `testdata/field-forms/README.md` listing each fixture from contracts/changes.md §6, stating that they are synthetic and every `expected-status.json` is hand-written; the raw-template case is built at test time from `testdata/upstream/speckit-roadmap/v0.1.0/templates/roadmap-template.md`

---

## Phase 2: Foundational

**Purpose**: the shared fixture test harness used by US1–US3.

- [X] T003 Write `internal/status/field_forms_test.go` with a table-driven `TestFieldFormsFixtures` that, for each `testdata/field-forms/*/expected-status.json`, runs `Evaluate`, `JSON` and `ExitCode` and checks the expected fields: `exit`, fatal `code`, per-item `superseded_by`, per-item `dependencies` IDs in order, and diagnostic `(item, code)` counts (exact). Reuse the `recognitionExpect` style from `recognition_test.go`, extended with `superseded_by`, `dependencies` and `diagnostics`

**Checkpoint**: harness compiles; fixtures added per story make it fail until implemented.

---

## Phase 3: User Story 1 — A capitalized "Superseded by" field is read (P1) 🎯 MVP

**Goal**: `Superseded by` / `superseded by` read as `superseded-by`; other variants stay fatal.

**Independent Test**: the three spellings give the same report as `superseded-by`; `SUPERSEDED BY` and a doubled field stay fatal.

- [X] T004 [P] [US1] Create fixtures with hand-written `expected-status.json` in `testdata/field-forms/`: `superseded-bold` (`- **Superseded by:** 002`), `superseded-plain` (`- Superseded by: 002`), `superseded-lower` (`- superseded by: 002`) each with items 001 and 002 → report, 001 `superseded_by` `002`, exit 0; `superseded-upper` (`- SUPERSEDED BY: 002`) → fatal `SS_ITEM_MALFORMED`, exit 2; `superseded-twice` (both `- superseded-by: 002` and `- Superseded by: 002` on 001) → fatal `SS_ITEM_MALFORMED`
- [X] T005 [US1] Add a test in `internal/status/field_forms_test.go` that each accepted spelling produces JSON byte-identical to the same roadmap written with `- superseded-by: 002` on the same line, and that prose `Notes: replaced, Superseded by: 002` outside a field line stays the C-10 fatal
- [X] T006 [US1] Add `{"Superseded by", "superseded-by"}` and `{"superseded by", "superseded-by"}` to `fieldLabels` in `internal/status/parse.go` (research R1); no other parser change
- [X] T007 [US1] Run `CGO_ENABLED=0 go test ./internal/status`; T004–T005 pass

---

## Phase 4: User Story 2 — A template placeholder means "not set yet" (P1)

**Goal**: bracketed or single-emphasis `Spec dir` / `Depends on` values are not set, with `SS_FIELD_PLACEHOLDER` (info); a valid enclosed ID list stays malformed.

**Independent Test**: the raw pinned template exits 0 with two placeholder findings; `Depends on: [002]` stays `SS_DEPENDENCY_MALFORMED`.

- [X] T008 [P] [US2] Create fixtures with hand-written `expected-status.json` in `testdata/field-forms/`: `template-placeholders` (a `planned` item with the template's exact lines `- **Depends on:** [NNN, NNN — or "none"]` and `- **Spec dir:** [specs/001-.../ — set once the spec exists; otherwise omit]`) → exit 0, two `SS_FIELD_PLACEHOLDER` for 001, one `SS_SPEC_REFERENCE_MISSING` (info), no `SS_DEPENDENCY_MALFORMED` or `SS_SPEC_DIRECTORY_MISSING`; `placeholder-italic` (a `specced` item with `- **Spec dir:** _set once the spec exists_`) → `SS_FIELD_PLACEHOLDER` plus `SS_SPEC_REFERENCE_MISSING` (error), exit 1; `bracketed-real-id` (`- Depends on: [002]` with item 002 present) → `SS_DEPENDENCY_MALFORMED`, no placeholder, exit 1
- [X] T009 [US2] Add tests in `internal/status/field_forms_test.go`: (a) copy `testdata/upstream/speckit-roadmap/v0.1.0/templates/roadmap-template.md` into a temp project at test time and assert a complete report, exit 0, exactly two `SS_FIELD_PLACEHOLDER` (lines 68 and 71) and no error (SC-002); (b) unit cases for `placeholderValue`: `[x]`, `[]`, `_x_`, `*x*` → placeholder; `**x**`, `__x__`, `_x*`, `[x](y)`, `x` → not; `Depends on` values `[002]`, `_002, 003_`, `[none]` → not a placeholder (valid list inside) and `SS_DEPENDENCY_MALFORMED`; (c) the message is exactly `Field holds a template placeholder and is treated as not set.` with the field line as source
- [X] T010 [US2] Implement `placeholderValue(value string) (inner string, ok bool)` in `internal/status/parse.go` per data-model "Placeholder value" ("whole value, one pair"; single `_` or `*` only; trimming spaces and one trailing period first), add `Placeholders []placeholder{Field string; Line int}` to `rawItem`, and apply it in `applyField` for `Spec dir` and `Depends on`: a placeholder records the entry and leaves the field unset (`HasSpec` false; `Deps` nil, `DepsBad` false); for `Depends on`, when `parseIDList(inner)` is valid, skip the placeholder path and parse the original value as today
- [X] T011 [US2] In `build` (`internal/status/analyze.go`), emit one `SS_FIELD_PLACEHOLDER` (info, item ID, roadmap path, the placeholder's line) per recorded placeholder
- [X] T012 [US2] Run `CGO_ENABLED=0 go test ./internal/status`; T008–T009 pass

---

## Phase 5: User Story 3 — A repeated dependency ID is reported once (P2)

**Goal**: one entry per ID, one `SS_DEPENDENCY_DUPLICATE` (warning) per repeated ID, no doubled dependency findings.

**Independent Test**: `duplicate-deps` gives dependencies `002, 003`, one duplicate warning, and `SS_DEPENDENCY_BLOCKED` once per unfinished ID.

- [X] T013 [P] [US3] Create `testdata/field-forms/duplicate-deps` with hand-written `expected-status.json`: item 001 `- Depends on: 002, 003, 002, 002` with planned items 002 and 003 → 001 dependencies `["002","003"]`, exactly one `SS_DEPENDENCY_DUPLICATE` for 001, `SS_DEPENDENCY_BLOCKED` exactly once per dependency (two), exit 0
- [X] T014 [US3] Add a test in `internal/status/field_forms_test.go` that the duplicate message is exactly `Dependency list names the same item more than once; it is counted once.` with `DepsLine` as the source, and that a list without repeats is unchanged
- [X] T015 [US3] In `applyField` (`internal/status/parse.go`), after `parseIDList` succeeds, keep the first occurrence of each ID and record each repeated ID once in a new `rawItem.DupDeps []string`; in `build` (`analyze.go`) emit one `SS_DEPENDENCY_DUPLICATE` (warning) per entry
- [X] T016 [US3] Run `CGO_ENABLED=0 go test ./internal/status ./internal/conformance`; all pass

---

## Phase 6: User Story 4 — Human text shows what the JSON already knows (P2)

**Goal**: `statuses` panel line, `invalid supersession N`, `(invalid)` row mark, shared-tasks row note; JSON unchanged.

**Independent Test**: human text on targeted roadmaps contains each addition; JSON is byte-identical to before.

- [X] T017 [US4] Add tests in `internal/status/style_test.go`: (a) the example roadmap's panel contains `statuses implemented 1  planned 1  verified 1` (keys in byte order, two spaces between pairs); an explicitly empty roadmap shows `statuses none`; an unknown status token appears raw and sanitized; (b) a roadmap with an invalid supersession (`superseded-by: 009`, no 009) shows a panel line `invalid supersession 1` and the row line `superseded-by 009 (invalid)`, and a malformed reference (`superseded-by: soon`) shows `superseded-by soon (invalid)`; with no invalid supersession neither appears; (c) two items with the same `Spec dir` show `tasks shared with 001; counted once` under the second; (d) `JSON` output for each of these is identical to what it was before (compare against an expectation captured from the T001 baseline build, not from the new build)
- [X] T018 [US4] Add display-only `TasksSharedWith string` to `Item` in `internal/status/model.go` (not added to `itemToView`); in `build`, track the first item ID per canonical `Tasks.Source` with state `read` and set it on later items
- [X] T019 [US4] In `internal/status/style.go`: add the `statuses` line after the `items …` line in `summaryLines` (sorted keys, `<status> <count>` pairs joined by two spaces, `statuses none` when empty, wrapped with `wrapText`); add `invalid supersession N` when `InvalidSupersessionItems > 0`; in `writeRelations`, print `superseded-by <id or sanitized raw reference> (invalid)` when `Supersession.Supplied` and `Valid` is false, and `tasks shared with <ID>; counted once` when `TasksSharedWith` is set
- [X] T020 [US4] Update every existing test that pins whole human panels (search `internal/status/*_test.go` and `internal/cli/*_test.go` for panel lines such as `task files read`) so each expectation gains exactly the new `statuses` line, edited by hand
- [X] T021 [US4] Run `CGO_ENABLED=0 go test ./...`; all pass

---

## Phase 7: User Story 5 — Command-line edge cases behave as written (P3)

**Goal**: rows F, G, H.

**Independent Test**: the three commands in contracts/changes.md §5 give the listed results.

- [X] T022 [P] [US5] Add tests in `internal/cli/cli_test.go` for every row of contracts/changes.md §5: `--json=false --bogus` → plain text on stderr, exit 2, empty stdout; `--json=true --bogus` → one JSON document; `--root ""`, `--root=`, `--roadmap ""`, `--roadmap=` on `status`, and `--root ""` on `inventory` → `SS_INVOCATION_INVALID` with `--root needs a non-empty value` / `--roadmap needs a non-empty value`, exit 2; `--json --fail-on warning --bogus` → `"fail_on":"warning"`; `--json --fail-on bogus` → `"fail_on":"error"`; a `--root` value spelled `--json` does not turn JSON mode on
- [X] T023 [US5] In `internal/cli/cli.go`: add `jsonRequested(args)` (last `--json` via `walkArgs`, parsed with `strconv.ParseBool`, unparsable = false) and use it at the three `countFlag(args, "--json") > 0` sites; in `singleString.Set`, reject an empty value for `root` and `roadmap` with `--<name> needs a non-empty value`; make `invocationJSON` take the threshold from `requestedFailOn(args)` (default `"error"`)
- [X] T024 [US5] Run `CGO_ENABLED=0 go test ./internal/cli`; all pass

---

## Phase 8: Polish and cross-cutting

- [X] T025 [P] Update `specs/001-spec-status/roadmap-format.md` §2 (the two spellings, the placeholder rule with its valid-list exception, duplicate handling), `specs/001-spec-status/diagnostics.md` (two registry rows from contracts/changes.md §2), `specs/001-spec-status/json-contract.md` (dependencies listed once per ID) and `specs/001-spec-status/cli-contract.md` (`--json=false`, empty values, `fail_on`)
- [X] T026 [P] Add both codes to `extensions/spec-status/reference/output.md` with severities matching the registry and fixes naming allowed commands (`/speckit.roadmap.write` for placeholders and duplicates); run `CGO_ENABLED=0 go test ./internal/extension`
- [X] T027 Update `README.md`: remove the Limits lines for rows A–H (the `Superseded by` and placeholder `Spec dir` bullets under "Field labels and values reported as errors", and the open-defect lines for human-text omissions, `--json=false`, empty `--root`/`--roadmap` and invocation `fail_on`), add findings rows for `SS_FIELD_PLACEHOLDER` and `SS_DEPENDENCY_DUPLICATE`, and replace the example panel with the new `spec status --color=never --root testdata/example-roadmap` output copied from a run and checked line by line
- [X] T028 Compare every T001 output: JSON byte-identical for every root (no fixture contains rows A–D); human text differs only by the `statuses` line and, where applicable, the other display additions; no exit code changed; record the result in the PR description
- [X] T029 Run `CGO_ENABLED=0 go vet ./...`, `gofmt -l .`, `go mod tidy -diff`, the full suite on darwin/arm64 and linux/amd64 (container), and `spec status --fail-on error` on this repository; record platforms and Go versions
- [X] T030 Add a row for feature 008 to the "Implementation record" table in `specs/001-spec-status/tasks.md` (implemented, executed by the implementer, not independently verified; platforms actually run), and draft the v0.3.0 release-note entries in the PR description
- [ ] T031 Doomer review on the source change (prepare before editing, submit, answer every Quick lead, finish) before handing off the PR
- [ ] T032 Propose roadmap 015 and 016 `specced` → `implemented` through `/speckit.roadmap.write` after the owner approves; `verified` only after the v0.3.0 release checks

---

## Dependencies and execution order

- T001 → everything; T003 before story fixtures are run.
- US1 (T004–T007), US2 (T008–T012), US3 (T013–T016) depend on T003; US2 and US3 both edit `applyField` in `parse.go`, so T010 and T015 run sequentially.
- US4 (T017–T021) is independent of US1–US3 except that T020 must run after all parser changes so full-text expectations are edited once.
- US5 (T022–T024) is independent (CLI only).
- Polish after all stories; T031 before handoff; T032 last.

## Parallel examples

- After T003: T004, T008, T013 (fixtures), T017 and T022 (tests) in parallel.
- In Polish: T025 and T026 in parallel.

## Implementation strategy

1. **MVP = US1 + US2**: they remove the fatal and the errors a freshly templated roadmap hits today.
2. US3, US4, US5 follow in the same PR; US4 is the only one that changes every human report.
3. One PR against `main`; owner merges; release through the owner when 013 and 011 are also in.
