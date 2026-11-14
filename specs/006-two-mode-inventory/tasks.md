---
description: "Tasks for feature 006: two-mode roadmap and feature inventory"
---

# Tasks: Two-mode roadmap and feature inventory

**Input**: `specs/006-two-mode-inventory/` (spec.md, plan.md, research.md, data-model.md, contracts/, quickstart.md). Roadmap item 008, decision D-2026-11-12-08, target v0.3.0.

**Tests**: Required. The spec's success criteria are hand-counted fixtures (SC-001, SC-003, SC-006), byte identity (SC-002) and repeat-run identity (SC-005). Tests are written first and must fail before the code they cover. Every `expected.json` is a hand count, never generated from `spec` output (CLAUDE.md).

**Organization**: by user story. US1 (recognition) and US2 (inventory) are both P1 and independent; US3 (safety) hardens US2.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: can run in parallel (different files, no dependency on an unfinished task)
- **[Story]**: US1, US2, US3 from spec.md

---

## Phase 1: Setup

- [X] T001 Branch from up-to-date `main` after PRs #47, #48 and #49 merge; set `.specify/feature.json` to `specs/006-two-mode-inventory`; build a baseline binary from `main` into the scratchpad and save its outputs for `testdata/example-roadmap` (human and `--json`), this repository's roadmap, and every `testdata/conformance/*` fixture, for the SC-002 byte comparison
- [X] T002 [P] Create the fixture root `testdata/two-mode/README.md` describing each fixture, that all are synthetic (shaped like the evidence in `specs/001-spec-status/two-mode-direction.md`, no copied content), and that every `expected.json` is hand-counted

---

## Phase 2: Foundational (blocks US2 and US3)

**Purpose**: the directory listing in the only filesystem boundary (research R5).

- [X] T003 Write `ListDir` tests in `internal/safeio/fs_test.go` that fail first: byte-ordered names; kinds `dir`, `file`, `symlink`, `other` (FIFO), `sensitive` (not stat'ed; assert with a sensitive-named FIFO that would block if opened); a symlinked or non-directory target → `unsafe`; absent → `notfound`; `..` and symlink components rejected via `prepare`; 1025 entries → `limit` with `MaxDirEntries = 1024`; directory replaced during listing (test seam) → `changed`; listing does not change the `MaxFiles`/`MaxTotalBytes` counters
- [X] T004 Implement `MaxDirEntries`, `DirEntry{Name, Kind, ID}` and `(*FS).ListDir(rel string) ([]DirEntry, error)` in `internal/safeio/fs.go` per data-model "DirEntry": `prepare`, `Lstat` (symlink, other device or non-directory → `unsafe`), open, read at most `MaxDirEntries+1` names, per-entry `Lstat` (other device → kind `foreign`), sensitive names not stat'ed, identity re-check → `changed`, sort by byte order
- [X] T005 [P] Add `openDirNoFollow` in `internal/safeio/open_unix.go` (`syscall.Open` with `O_RDONLY|O_NOFOLLOW|O_DIRECTORY|O_CLOEXEC`, `confirmOpened` on Linux) and an `unsupported` stub in `internal/safeio/open_other.go`; confirm `internal/status/imports_test.go` and the side-effect string checks still pass (no `os.OpenFile`)
- [X] T006 Run `CGO_ENABLED=0 go test ./internal/safeio` and the full suite; all pass

**Checkpoint**: `ListDir` exists and is tested; nothing calls it yet.

---

## Phase 3: User Story 1 — An unrecognized roadmap is never reported as empty (P1) 🎯 MVP

**Goal**: zero-item roadmaps are either explicitly empty (`## Planned Specs`, exit 0) or fatal `SS_ROADMAP_UNRECOGNIZED` (exit 2); item-bearing reports are byte-identical (contracts/roadmap-recognition.md).

**Independent Test**: the recognition fixtures give the expected code and exit; the T001 baseline outputs compare equal.

### Tests for User Story 1

- [X] T007 [P] [US1] Create recognition fixtures under `testdata/two-mode/` with hand-written `expected-status.json` (code, severity, exit, message prefix, `total_items`, percents): `unrecognized-empty` (0-byte roadmap), `unrecognized-whitespace`, `unrecognized-narrative` (headings and prose, no items), `planned-in-fence` (heading only inside a fence), `planned-wrong-form` (`### Planned Specs`, `## planned specs`, `  ## Planned Specs` in one file), `explicit-empty` (pinned template sections, `## Planned Specs`, only `###` grouping headings), `items-no-planned` (items, no heading), `narrative-with-label` (a `Depends on:` label outside items → existing C-10 fatal), `no-roadmap` (no roadmap file; also used by US2)
- [X] T008 [US1] Write `internal/status/recognition_test.go`: each T007 fixture through `Evaluate`, `JSON`, `Text` and `ExitCode` matches its `expected-status.json`; unrecognized with `FailOn: "none"` exits 2; `--roadmap` pointing at an unrecognized file gives the same fatal; the explicit-empty human panel contains `roadmap explicitly empty (Planned Specs heading, no items)` and a roadmap with items does not; the missing-roadmap message contains `spec inventory`. Must fail before T010
- [X] T009 [US1] Replace `TestEmptyRoadmapJSON` in `internal/status/status_test.go`: a 0-byte roadmap is now unrecognized (exit 2); move its null-denominator, no-`100`, no-timing-field and determinism assertions onto an explicitly empty roadmap

### Implementation for User Story 1

- [X] T010 [US1] In `internal/status/parse.go`, make `parseRoadmap` also return `plannedSpecs bool`: true when a line not marked `hidden` (fence or comment) is exactly `## Planned Specs` at column 0 with optional trailing spaces or tabs (research R1); update every caller, including tests
- [X] T011 [US1] In `internal/status/run.go`, after a successful parse with zero items: `plannedSpecs` → continue to `build` (explicitly empty); otherwise return `fatalReport("SS_ROADMAP_UNRECOGNIZED", …)` with the exact messages from contracts/roadmap-recognition.md (empty or whitespace-only file uses the "Roadmap file is empty." message); for safeio `notfound` only, change the `SS_ROADMAP_UNAVAILABLE` message to `Required roadmap is unavailable. Run spec inventory for a feature list without a roadmap.`
- [X] T012 [US1] In `internal/status/style.go`, add the summary-panel line `roadmap explicitly empty (Planned Specs heading, no items)` only when the report is complete with zero items; no other rendering change
- [X] T013 [US1] Run T008/T009 and the full suite; compare every T001 baseline output byte for byte (only the intended missing-roadmap message may differ, and no baseline fixture is missing a roadmap); record the comparison in the PR description

**Checkpoint**: US1 is complete and shippable on its own.

---

## Phase 4: User Story 2 — List features in a repository without a roadmap (P1)

**Goal**: `spec inventory` lists `specs/` children, file states and tasks.md checkbox counts, as human text or inventory JSON 1.0 (contracts/inventory-cli.md, inventory-json.md, inventory-text.md).

**Independent Test**: `spec inventory` on the inventory fixtures matches the hand-counted `expected.json`; `spec status` never runs inventory.

### Tests for User Story 2

- [X] T014 [P] [US2] Create inventory fixtures with hand-counted `expected-inventory.json` (the full inventory JSON document) under `testdata/two-mode/`: `no-roadmap` (from T007) gets `specs/001-alpha/` (spec, plan, tasks with 3 of 5 checked), `specs/20261101-093000-beta/` (spec only), `specs/001-alpha/checklists/requirements.md` with checked boxes, `specs/notes/` (not a feature), `specs/README.md` (a file, ignored), `specs/3-short/`, `specs/003_name/`, `specs/20261101-gamma/` (sequential form, tasks.md with zero checkboxes → 0/0); `brownfield` (one feature folder, tasks 24/24, a 16-check `checklists/requirements.md`, `spec.md` with `Status: Draft`); `no-specs` (`.specify` only); `with-roadmap` (a roadmap with items plus one feature folder, to show the roadmap is not read)
- [X] T015 [US2] Write `internal/status/inventory_test.go`: each fixture through `EvaluateInventory`, `InventoryJSON` and `ExitCode` matches `expected-inventory.json` byte for byte; `name_form` values `sequential`/`timestamp`; folder order is byte order; requirements checklist marks never counted; `brownfield` reports 24/24 (SC-006); `no-specs` gives `SS_SPECS_DIRECTORY_MISSING` (info), `specs_state: "missing"`, `percent: null`, exit 0, and exit 1 with `FailOn: "info"`; `with-roadmap` reports `roadmap_file: "present"` and the roadmap file is never read (assert via the safeio after-read hook); a tasks.md the narrow dialect rejects is fatal citing the tasks file and line. Must fail before T017
- [X] T016 [P] [US2] Write CLI tests in `internal/cli/cli_test.go`: `spec inventory` human and `--json` exit codes; each rejected flag (`--roadmap`, `--status`, `--blocked`, `--human-review`, `--active`, `--historical`, `--explain`) → `SS_INVOCATION_INVALID` `--<name> is not used by inventory`, exit 2; with `--json`, invocation errors are one inventory-shaped document carrying the requested `fail_on`; `--version` with `inventory` rejected; repeated flags rejected as for `status`; `spec status` with no roadmap does not print inventory output; `spec --help` lists `inventory`; `spec completion bash` mentions `inventory`

### Implementation for User Story 2

- [X] T017 [US2] Extract the shared setup in `internal/status/run.go` (root resolution, `safeio.New`, `.specify` check) into one helper used by `Evaluate` and the new evaluator, with no change to `spec status` results
- [X] T018 [US2] Create `internal/status/inventory.go` with the `InventoryReport`, `Folder`, `FileState`, `TaskState` and `InventorySummary` types from data-model.md and `EvaluateInventory(opt InventoryOption) InventoryReport`: list `specs` with `ListDir`; `notfound` → `SS_SPECS_DIRECTORY_MISSING` (info), `specs_state: "missing"`; `unsafe`/`sensitive` → `SS_PATH_UNSAFE`/`SS_REFERENCE_SENSITIVE_SKIPPED` (error, report), `specs_state: "skipped"`; `limit`/`changed` → fatal `SS_INPUT_LIMIT`/`SS_SOURCE_CHANGED`
- [X] T019 [US2] In `internal/status/inventory.go`, classify children (research R7): only kind `dir` can be a feature; timestamp form `^[0-9]{8}-[0-9]{6}-.+$` first, then sequential `^[0-9]{3,}-.+$`, ASCII only; other real directories → `feature: false` with `SS_FOLDER_NOT_FEATURE` (info, source `specs/<name>`); kind `file` ignored; kinds `symlink`, `foreign`, `other`, `sensitive` → `skipped: true`, counted in `skipped_entries`, with `SS_PATH_UNSAFE` or `SS_REFERENCE_SENSITIVE_SKIPPED`
- [X] T020 [US2] In `internal/status/inventory.go`, for each feature folder `StatRel` `spec.md`, `plan.md`, `tasks.md`: regular → `present`; absent → `missing`; symlink, non-regular or other device → `skipped` reason `unsafe` with `SS_PATH_UNSAFE`; sensitive → `skipped` reason `sensitive`. Read a regular `tasks.md` with `ReadRel` and count with `parseTasks` (no task-link checks); a parse fatal or `SS_ENCODING_INVALID` makes the whole run fatal citing the tasks file and line; enforce the existing `maxTasks` total; read failures map as in roadmap mode (`limit`/`changed` fatal, others `skipped`)
- [X] T021 [US2] In `internal/status/inventory.go`, compute the summary per data-model "InventorySummary" (percent null when total is 0), stat `.specify/memory/roadmap.md` for `roadmap_file` (`present`/`absent`/`skipped`, never read), sort diagnostics with `diagLess`, and add an exit-code function sharing the roadmap-mode threshold rule (fatal or incomplete → 2)
- [X] T022 [P] [US2] Create `internal/status/inventory_json.go`: `InventoryJSON(rep InventoryReport) []byte` producing exactly the key order, types and null rules in contracts/inventory-json.md (report and fatal documents; HTML escaping off; trailing newline; diagnostics reuse the schema 1.0 diagnostic view with `item_id: null`), plus a serialization-failure fallback document like `JSON`'s
- [X] T023 [US2] In `internal/status/style.go`, render inventory human text per contracts/inventory-text.md: panel (`Spec Inventory` + build, mode line, checked tasks with bar, folder counts, task file counts, source), wide table and the existing narrow stacked layout, `not a feature folder` rows, `Checked boxes are not verification.`, the `A roadmap exists; spec status reports it.` line when `roadmap_file` is present, diagnostics list; names sanitized of terminal escapes; fatal text via the existing plain `spec: CODE: message`
- [X] T024 [US2] In `internal/cli/cli.go`, add `newInventory`: `Use: "inventory"`, help text stating the roadmap is not read and checked boxes are not verification, `cobra.NoArgs`; reject roadmap-only flags when set; call `EvaluateInventory`; `--json` prints `InventoryJSON` with its exit code; human text to stdout (fatal to stderr, exit 2) with `wantColor` and `textColumns`; register it in `newRoot`
- [X] T025 [US2] In `internal/cli/cli.go`, make `fail`/`invocationJSON` emit the inventory-shaped fatal document (with the requested `fail_on` when it parsed) when `hasCommand(args, "inventory")`; status-shaped invocation JSON unchanged otherwise
- [X] T026 [US2] Run T015/T016 and the full suite; all pass

**Checkpoint**: US1 and US2 both work independently.

---

## Phase 5: User Story 3 — Inventory is as safe and plain as the roadmap report (P2)

**Goal**: the guarantees in spec US3 hold for inventory.

**Independent Test**: the temp-directory safety tests and repeat-run comparisons pass.

- [X] T027 [US3] Add temp-directory tests to `internal/status/inventory_test.go` (created at test time, not committed): a symlinked feature folder and a symlinked `tasks.md` → `skipped`, `SS_PATH_UNSAFE`, contents never read (after-read hook); a FIFO `tasks.md` → `skipped`, run returns promptly; a sensitive-named child folder and a sensitive-named file inside a feature folder → `SS_REFERENCE_SENSITIVE_SKIPPED`, not opened; `specs` as a symlink and as a regular file → `SS_PATH_UNSAFE`, `specs_state: "skipped"`, zero folders; 1025 children → fatal `SS_INPUT_LIMIT`, `folders: []`, exit 2; tasks totals over `maxTasks` → fatal
- [X] T028 [P] [US3] Add determinism tests: two `EvaluateInventory` runs on every inventory fixture give identical `InventoryJSON` bytes and identical human text at fixed width with color off; output contains no absolute path and no timestamp
- [X] T029 [P] [US3] Add a human-text test in `internal/status/style_test.go`: a folder name with terminal escape sequences is sanitized; color is off unless requested; width follows `COLUMNS` (fixed in the test) and the narrow layout is used below the existing threshold
- [X] T030 [US3] Run the full suite, including `internal/status/imports_test.go` and the dependency allowlist checks; all pass

**Checkpoint**: all three stories done.

---

## Phase 6: Polish and cross-cutting

- [X] T031 [P] Update `specs/001-spec-status/diagnostics.md`: add `SS_ROADMAP_UNRECOGNIZED` (error, fatal), `SS_SPECS_DIRECTORY_MISSING` (info, report), `SS_FOLDER_NOT_FEATURE` (info, report) with authority "D-2026-11-12-08, feature 006"; note inventory's use of existing codes
- [X] T032 [P] Update `specs/001-spec-status/json-contract.md` (a complete report with zero items means an explicitly empty roadmap; schema 1.0 unchanged) and add `specs/001-spec-status/inventory-json.md` from `specs/006-two-mode-inventory/contracts/inventory-json.md`
- [X] T033 [P] Update `specs/001-spec-status/cli-contract.md` (the `inventory` command, rejected flags, missing-roadmap message) and `specs/001-spec-status/roadmap-format.md` (recognition rule and the exact heading)
- [X] T034 [P] Update `extensions/spec-status/reference/output.md` with the three codes (meaning, severity, fatal, fix: for `SS_ROADMAP_UNRECOGNIZED` "use the supported template via `/speckit.roadmap.write`, or run `spec inventory`"); run `CGO_ENABLED=0 go test ./internal/extension`
- [X] T035 Update `README.md`: remove the zero-items Limits line and the "does not tell those two apart" sentence; the `SS_ROADMAP_UNAVAILABLE` row mentions `spec inventory`; add `SS_ROADMAP_UNRECOGNIZED` to the findings table; add a short "Without a roadmap: spec inventory" section with a sample taken from a fixture run and kept identical by a test in `internal/status/inventory_test.go`; add the over-inclusive sensitive-name skip of feature folders to Limits; keep the tone ("executed by the implementer, not independently verified")
- [X] T036 Update `docs/agent-prompt.md` and the README copy identically, only if a step mentions the missing-roadmap behavior; the CI hygiene check must still pass
- [X] T037 Draft the Art. II amendment note (affected requirements, compatibility, fixtures, risks per Governance item 4) in the PR description for owner approval; do not edit `.specify/memory/constitution.md` until the owner approves; do not run `/speckit-constitution`
- [X] T038 Draft the v0.3.0 release-note entries in the PR description: `SS_ROADMAP_UNRECOGNIZED` replaces "zero items, exit 0" (pre-1.0 behavior change); the missing-roadmap message text changed; new `spec inventory` with `inventory_schema_version` 1.0
- [X] T039 Add a row for feature 006 to the "Implementation record" table in `specs/001-spec-status/tasks.md` (implemented, executed by the implementer, not independently verified; platforms actually run)
- [X] T040 Run every scenario in `specs/006-two-mode-inventory/quickstart.md` on darwin/arm64, and the test suite on linux/amd64 (container), and record the results in the PR description; run `spec status` on this repository and confirm it is unchanged apart from the 008 entry
- [ ] T041 Propose the roadmap update for 008 (`specced` → `implemented`, notes) through `/speckit.roadmap.write` after the owner approves; `verified` only after the v0.3.0 release checks

---

## Dependencies and execution order

- **Phase 1** → everything.
- **Phase 2 (ListDir)** blocks US2 and US3, not US1.
- **US1 (Phase 3)** depends only on Phase 1. It can be done, reviewed and merged before or in parallel with Phase 2.
- **US2 (Phase 4)** depends on Phase 2. T017 before T018; T018 → T019 → T020 → T021 in one file; T022 and T023 after T021; T024 → T025.
- **US3 (Phase 5)** depends on US2.
- **Polish (Phase 6)**: T031–T034 any time after the codes are fixed (after T011 and T019); T035 after T023; T040 last but one; T041 last.

## Parallel examples

- After T001: T002, T003 and T007 in parallel (different files).
- In Phase 2: T005 alongside T004 once the function signature is fixed.
- In US2: T014 and T016 in parallel; T022 alongside T023 after T021.
- In Polish: T031, T032, T033, T034 in parallel.

## Implementation strategy

1. **MVP = US1.** It removes the silent "zero items, exit 0" result, the only behavior change in the release. Ship it as its own PR if review time is short.
2. **Then Phase 2 + US2** as the second PR: the new command.
3. **US3 and Polish** complete the feature; the README, contracts and extension reference must land with the code they describe.
4. Doomer review on every source-changing PR before delivery; owner merges; release only through the owner (tag v0.3.0 after 013–016 and 011 are also in).
