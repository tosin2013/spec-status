# Tasks: Release archives for linux/arm64 and darwin/amd64

**Input**: Design documents from `specs/014-more-release-platforms/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/platforms.md, quickstart.md

**Tests**: The evidence is the workflow runs themselves: PR CI legs and release smoke jobs (FR-003, FR-006). Local checks are actionlint, cross-compiles and `go test ./...`.

**Organization**: Tasks are grouped by user story. US1 and US2 (the two archives) share the release workflow edits in Phase 2.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story the task belongs to (US1–US4)

## Phase 1: Setup

- [ ] T001 Cross-compile all four targets locally (`GOOS`/`GOARCH` = linux/amd64, linux/arm64, darwin/amd64, darwin/arm64; `CGO_ENABLED=0 go build -o /dev/null ./cmd/spec`). Run `actionlint` on `.github/workflows/release.yml` and `.github/workflows/ci.yml` to get a clean starting point.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The release workflow builds, smokes and drafts four archives.

- [ ] T002 In `.github/workflows/release.yml`, set the build loop in "Build, package and checksum" to `for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64`. Keep the checksum line's `spec_"${VERSION}"_*.tar.gz` glob.
- [ ] T003 In `.github/workflows/release.yml`, change the smoke matrix to `os: [ubuntu-latest, ubuntu-24.04-arm, macos-latest, macos-15-intel]`, and extend the platform `case` (research R2):
  - `"Linux aarch64"` → `os=linux; arch=arm64; sum="sha256sum"`;
  - `"Darwin x86_64"` → `os=darwin; arch=amd64; sum="shasum -a 256"`.

  Keep the `*)` error. Keep "Smoke-test the extension archive" at `if: matrix.os == 'ubuntu-latest'`.
- [ ] T004 In `.github/workflows/release.yml`, change the draft job "Create draft release" (research R3):
  - for each row in `smoke (ubuntu-latest)|linux/amd64`, `smoke (ubuntu-24.04-arm)|linux/arm64`, `smoke (macos-latest)|darwin/arm64` and `smoke (macos-15-intel)|darwin/amd64`, look up the job URL; if it is empty, fail with `::error::`;
  - otherwise print the archive line ``  - `spec_${VERSION}_<os>_<arch>.tar.gz`: smoke-tested on <os>/<arch> (<url>)``;
  - pass the four archives explicitly to `gh release create`, in the same order.

**Checkpoint**: `actionlint` is clean, and a read-through confirms that every archive name appears in build, smoke, notes and assets.

---

## Phase 3: User Story 1 - An arm64 Linux image installs a verified binary (Priority: P1) 🎯 MVP

**Goal**: linux/arm64 is shipped, checksummed and natively smoke-tested.

**Independent Test**: The release run's `smoke (ubuntu-24.04-arm)` passes with `target: linux/arm64`.

- [ ] T005 [US1] Check `.github/workflows/release.yml` for the linux/arm64 path end to end: build target, the case row `Linux aarch64`, the notes row and the asset path. Run a local dry run of the smoke `case` logic with `uname -s -m` values `Linux aarch64` and `Darwin x86_64` in a scratch bash script.

---

## Phase 4: User Story 2 - An Intel Mac user installs a verified binary (Priority: P1)

**Goal**: darwin/amd64 is shipped, checksummed and natively smoke-tested.

**Independent Test**: `smoke (macos-15-intel)` passes with `target: darwin/amd64`.

- [ ] T006 [US2] Locally, build darwin/amd64 with the release flags (`-trimpath -ldflags "-s -w -X …Version=0.4.1-test -X …Kind=release"`) and run it under Rosetta: `arch -x86_64 ./spec version` must print `target: darwin/amd64` and `label: release build`. Then delete the binary.

---

## Phase 5: User Story 3 - Every change is tested where the binaries run (Priority: P2)

**Goal**: PR CI runs the test suite on four platforms with a two-way platform assert.

**Independent Test**: This PR's CI shows passing `test (ubuntu-24.04-arm, …)` and `test (macos-15-intel, …)` jobs.

- [ ] T007 [US3] In `.github/workflows/ci.yml`, change the test matrix to `os: [ubuntu-latest, ubuntu-24.04-arm, macos-latest, macos-15-intel]` (keeping the go versions). In "Assert platform", map each runner to its `go env` pair and its `uname -s -m` value (research R4: ubuntu-latest `linux amd64`/`Linux x86_64`, ubuntu-24.04-arm `linux arm64`/`Linux aarch64`, macos-latest `darwin arm64`/`Darwin arm64`, macos-15-intel `darwin amd64`/`Darwin x86_64`), and fail on either mismatch.
- [ ] T008 [US3] Check every other `ci.yml` step for per-runner assumptions, such as a step or `if:` keyed to `macos-latest` or `ubuntu-latest`, and extend them where needed. Note in the PR that branch protection still requires only the original test jobs; changing that is the owner's choice.

---

## Phase 6: User Story 4 - Docs and the agent prompt know the new platforms (Priority: P2)

**Goal**: Every install path maps all four platforms.

**Independent Test**: `go test ./internal/cli` passes (README samples and unpinned installs), and no doc says only two platforms are supported.

- [ ] T009 [P] [US4] `docs/agent-prompt.md` STEP 1 platform line: Linux x86-64 → linux_amd64, Linux arm64 (aarch64) → linux_arm64, macOS Apple silicon → darwin_arm64, macOS Intel → darwin_amd64. Anything else (including Windows) has no release archive, so stop. Keep the `<!-- agent-prompt:* -->` markers and the line widths of the text block.
- [ ] T010 [P] [US4] `README.md`:
  - the platform table gets rows for linux/arm64 and darwin/amd64, marked "v0.4.1" (links added after the release);
  - both `OS_ARCH=` comments list all four values;
  - the Limits line names the four platforms and says Windows is planned (roadmap 020).
- [ ] T011 [P] [US4] `docs/getting-started.md`, `docs/architecture.md` (CI row, release diagram build node, supported-platform sentence), `docs/development.md` (CI line) and `docs/limits.md`: four platforms.
- [ ] T012 [P] [US4] `RELEASING.md`: the asset checklist becomes four archives, the checksums file and the extension zip (it says "Exactly three assets" today). The smoke list becomes four jobs.
- [ ] T013 [P] [US4] `specs/004-release-engineering/contracts/release-assets.md` and `release-workflow.md`: four archives, four smoke runners and the per-archive notes lines, linking `specs/014-more-release-platforms/contracts/platforms.md`.

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T014 `extensions/spec-status/extension.yml` `version: 0.4.1` and the `extensions/spec-status/README.md` heading `v0.4.1` (the release `verify` job and `go test ./internal/extension`).
- [ ] T015 Run `actionlint` (both workflows), `CGO_ENABLED=0 go test ./...` and `go vet ./...`, and check that `spec status --fail-on error` on this repository exits 0. `git diff --stat` must show no `.go` file changed.
- [ ] T016 [P] Add a row "Release platforms (feature 014, roadmap 019)" to the implementation record in `specs/001-spec-status/tasks.md`, with the local results and "native CI and smoke evidence comes from the PR CI and the v0.4.1 release run". Mark it "executed by the implementer, not independently verified".

---

## Dependencies & Execution Order

### Phase Dependencies

- T001, then Phase 2 (T002–T004, one file, in sequence).
- US1 and US2 then check Phase 2's work. US3 (ci.yml) is independent of Phase 2. US4 docs are independent.
- Polish comes last.

### Parallel Opportunities

- T009–T013 are separate doc files.
- T007 (ci.yml) can run alongside Phase 2 (release.yml).

## Parallel Example: User Story 4

```text
T009 docs/agent-prompt.md
T010 README.md
T011 docs/*.md
T012 RELEASING.md
T013 specs/004-release-engineering/contracts/*.md
```

## Implementation Strategy

1. The release workflow (Phase 2), proven by US1 and US2 locally, then by the v0.4.1 run. This is the MVP for issue #82.
2. CI legs (US3), proven by this PR's own CI.
3. Docs (US4), version, record. Then the PR, merge, tag v0.4.1, smoke evidence, post-release links and 019 `verified`.
