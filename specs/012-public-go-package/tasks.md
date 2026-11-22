# Tasks: Public Go package

**Input**: Design documents from `specs/012-public-go-package/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/go-api.md, quickstart.md

**Tests**: Required. The spec calls for byte identity on every fixture (FR-009), a race test (FR-008), an import check (FR-007), a compiled example (FR-011) and an unchanged CLI (FR-012).

**Organization**: Tasks are grouped by user story. US1 and US2 are both P1 and share the schema refactor in Phase 2.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story the task belongs to (US1, US2, US3)

## Phase 1: Setup

- [ ] T001 Build the pre-change binary `CGO_ENABLED=0 go build -o $TMPDIR/spec-base ./cmd/spec` from `main`. Save text, `--explain`, `--json` and exit codes for every fixture root under `testdata/` into `$TMPDIR/pkg-baseline/`, using the snapshot script from feature 011 (`$TMPDIR/snap.sh`).

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: One set of schema types, with the CLI rewired onto them and its output unchanged.

- [ ] T002 Create `internal/schema/schema.go` (package doc: "schema 1.0 report types shared by the CLI and pkg/status; standard library only"). Add these types, with JSON tags in the current field order and the Go types from `specs/012-public-go-package/data-model.md`:
  - `Report`, `Source`, `Configuration`, `Summary`, `TaskSummary`, `Item`, `Ref`, `Supersession`, `Dependency`, `Tasks`, `Phase`, `Diagnostic`, `Error`, `Location`;
  - "pointers where the JSON allows null; `Phases []Phase` nil = null, empty = []".
- [ ] T003 Move `Percent` from `internal/status/model.go` into `internal/schema/percent.go`:
  - keep the unexported fields `ok` and `hundredths`;
  - add `PercentOf(num, den int) Percent`, keeping the half-up rounding of `percentOf`;
  - add `Valid() bool` and `Hundredths() int`;
  - keep `Text()` and `MarshalJSON` unchanged;
  - add `UnmarshalJSON`, which accepts `null` or a non-negative number with at most two decimals and rejects anything else.

  In `internal/status/model.go`, replace the type with `type Percent = schema.Percent` and make `percentOf` call `schema.PercentOf`.
- [ ] T004 Create `internal/schema/marshal.go` with `Marshal(r Report) []byte`: a `json.Encoder` with `SetEscapeHTML(false)` and a trailing newline. On encode failure, return the existing fallback document from `internal/status/render.go` byte for byte.
- [ ] T005 Create `internal/schema/exit.go` with `ExitCode(r Report) int` (research R6):
  - 2 if `Outcome == "fatal"` or `!Complete`;
  - 0 if `Configuration.FailOn == "none"`;
  - otherwise 1 if any diagnostic's severity rank (info 1, warning 2, error 3) is at least the threshold rank, else 0.
- [ ] T006 Rewrite the JSON path in `internal/status/render.go`:
  - replace `docView` and the other `*View` structs, `srcAny`, `sourceRef` and `nullString` with `toSchema(rep Report) schema.Report`, keeping every rule of the current `JSON` (fatal handling, empty `status_counts` as `{}`, non-nil slices, `phases` null unless read);
  - `JSON(rep)` becomes `schema.Marshal(toSchema(rep))`;
  - export `Schema(rep Report) schema.Report` as the bridge for `pkg/status`.
- [ ] T007 Make `ExitCode` in `internal/status/run.go` return `schema.ExitCode(toSchema(rep))`.
- [ ] T008 Run `CGO_ENABLED=0 go test ./...` and compare every fixture root against `$TMPDIR/pkg-baseline/`: text, `--explain`, `--json` bytes and exit codes must be byte-identical (FR-012). Fix any difference in `internal/status/render.go` before continuing.

**Checkpoint**: The CLI is byte-identical and all existing tests pass on the new schema types.

---

## Phase 3: User Story 1 - A Go gate evaluates a project in-process (Priority: P1) 🎯 MVP

**Goal**: An outside Go program loads a project with one call and gets typed schema 1.0 values.

**Independent Test**: `go test ./pkg/status -run 'TestLoad'` (external test package) passes.

### Tests for User Story 1

- [ ] T009 [P] [US1] In `pkg/status/status_test.go` (`package status_test`), add `TestLoadExampleRoadmap`: load `../../testdata/example-roadmap` and check the following against the CLI's documented example:
  - `SchemaVersion == status.SchemaVersion == "1.0"`;
  - `Outcome == "report"`;
  - the item IDs and statuses;
  - the summary counts;
  - `Items[i].Tasks.Phases` nil or non-nil according to the state.
- [ ] T010 [P] [US1] In `pkg/status/status_test.go`, add `TestLoadFatalAndInvalid`:
  - a root without a roadmap gives `Outcome == "fatal"`, `Error.Code == "SS_ROADMAP_UNAVAILABLE"`, `Summary == nil` and a nil `error`;
  - `Load("", …)` returns an error mentioning the root;
  - `Load(root, Options{FailOn: "loud"})` returns an error mentioning fail-on;
  - a roadmap outside the root gives a fatal `SS_PATH_UNSAFE` in the report.

### Implementation for User Story 1

- [ ] T011 [US1] Create `pkg/status/status.go`:
  - type aliases for every schema type (`type Report = schema.Report`, … `type Percent = schema.Percent`);
  - `const SchemaVersion = "1.0"`;
  - `type Options struct{ Roadmap, FailOn string }`;
  - `Load(root string, opts Options) (Report, error)`: validate "root must not be empty" and "FailOn must be one of "", error, warning, info, none", then return `status.Schema(status.Evaluate(status.Option{Root: root, Roadmap: opts.Roadmap, FailOn: opts.FailOn}))`;
  - `Marshal(r Report) []byte` returning `schema.Marshal(r)`;
  - `ExitCode(r Report) int` returning `schema.ExitCode(r)`.
- [ ] T012 [US1] Create `pkg/status/doc.go` with the package doc: purpose, the "exactly as `spec status --json --root`" guarantee, fatal-in-report versus Go error, the safety properties, "report text is untrusted; sanitize before printing to a terminal", and the experimental-until-1.0 promise (contracts/go-api.md).

**Checkpoint**: US1 tests pass.

---

## Phase 4: User Story 2 - The library and the binary cannot disagree (Priority: P1)

**Goal**: Proven byte and exit-code identity with the CLI, a lossless round trip, and concurrency safety.

**Independent Test**: `go test -race ./pkg/status -run 'TestMatchesCLI|TestRoundTrip|TestConcurrent'` passes.

### Tests for User Story 2

- [ ] T013 [US2] In `pkg/status/status_test.go`, add `TestMatchesCLI`. For every directory under `../../testdata` that directly contains `.specify/` (non-symlink), and each threshold in {error, warning, info, none}:
  - run `cli.Run([]string{"status", "--json", "--root", abs, "--fail-on", t}, cwd, &out, &errb)` from `internal/cli`;
  - assert `Marshal(rep)` equals `out` byte for byte;
  - assert `ExitCode(rep)` equals the returned code.

  Fail the test if fewer than 40 roots are found, so that a broken walk cannot pass silently.
- [ ] T014 [US2] In `pkg/status/status_test.go`, add `TestRoundTrip`: for each root's CLI output, `json.Unmarshal` into `status.Report`, then `Marshal`, and assert the bytes are equal. Add Percent cases: `null`, `0`, `50`, `41.67`, `0.05`, and rejection of `"41"` and `1.234`.
- [ ] T015 [US2] In `pkg/status/status_test.go`, add `TestConcurrentLoads`: 8 goroutines × 3 roots (the example roadmap, `testdata/phases` and a fatal root), each result compared with a sequential `Marshal`. Run under `-race` in T021.

### Implementation for User Story 2

- [ ] T016 [US2] Fix any mismatch found by T013–T015 in `internal/schema/*.go` or `internal/status/render.go` (never in the tests or fixtures).

---

## Phase 5: User Story 3 - Consumers know what is promised (Priority: P2)

**Goal**: The docs and a compiled example state how to use the package and what is promised.

**Independent Test**: `go test ./pkg/status -run Example` passes. `docs/library.md` exists and is linked from the README.

- [ ] T017 [P] [US3] Create `pkg/status/example_test.go` with `Example()`. It loads `../../testdata/example-roadmap` and prints `SchemaVersion`, `Outcome`, the item count, verified/active and the exit code, with an `// Output:` block taken from a real run.
- [ ] T018 [P] [US3] Create `docs/library.md`:
  - `go get github.com/tosin2013/spec-status/pkg/status@latest`;
  - a short program;
  - fatal outcomes versus errors;
  - `Marshal` and `ExitCode`;
  - decoding the CLI JSON into `Report`;
  - the safety properties;
  - untrusted text;
  - the experimental-until-1.0 promise and `SchemaVersion`.

  No output-looking block without a `spec-sample` label. Code blocks are Go source, which `TestReadmeSamples` ignores; confirm that with the test.
- [ ] T019 [P] [US3] Add one line and a link to `docs/library.md` in `README.md` (near "Use it with an AI coding agent"). In `docs/architecture.md`, add `pkg/status` and `internal/schema` to the containers and components tables and the components diagram, and update the version line. In `CLAUDE.md`, describe the pipeline (`status.JSON` → `schema.Marshal(toSchema)`) and the new packages.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T020 Extend `internal/status/imports_test.go`: non-test files under `internal/schema/` and `pkg/status/` may import only the standard library and `github.com/tosin2013/spec-status/internal/...`. Seed an outside import once to see the test fail, then revert it.
- [ ] T021 Run `CGO_ENABLED=0 go test ./...`, `go test -race ./pkg/status ./internal/status`, `go vet ./...`, gofmt and `go doc ./pkg/status`. Repeat the T008 baseline comparison.
- [ ] T022 [P] Add a "Go API" note (pkg/status, experimental until 1.0, same bytes) in `specs/001-spec-status/json-contract.md`. Add a row "Public Go package (feature 012, roadmap 018)" to the implementation record in `specs/001-spec-status/tasks.md`, with the T013 root count, the T008/T021 results and the seeded-fault checks, marked "executed by the implementer, not independently verified".

---

## Dependencies & Execution Order

### Phase Dependencies

- T001 comes before any code change.
- Phase 2 (T002–T008) blocks every story. T002–T005 create `internal/schema`. T006–T007 rewire `internal/status`. T008 verifies the CLI.
- US1 depends on Phase 2. US2 depends on US1's `Load`. US3 can start after US1.
- Polish comes last.

### Parallel Opportunities

- T009 and T010 run together, before T011 (the same test file, but independent functions).
- T017, T018 and T019 are different files.
- T022 is independent documentation.

## Parallel Example: User Story 3

```text
T017 pkg/status/example_test.go
T018 docs/library.md
T019 README.md, docs/architecture.md, CLAUDE.md
```

## Implementation Strategy

1. Baseline, then the schema refactor with a byte-identical CLI (Phase 2). This is the riskiest step, so it is done and verified first.
2. US1 `Load` is the MVP for issue #65.
3. US2 proves identity, the round trip and concurrency.
4. US3 docs, then guardrails and the record.
5. Doomer review on the source changes, then one PR against `main`. The owner merges.
