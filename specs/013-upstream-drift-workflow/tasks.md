# Tasks: Upstream drift workflow

**Input**: Design documents from `specs/013-upstream-drift-workflow/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/drift-workflow.md, quickstart.md

**Tests**: Required. SC-003 asks for seeded runs for each failure kind, and FR-009 for a check of the shared pin list.

**Organization**: Tasks are grouped by user story. The pin list and the checker script (Phase 2) serve all three stories.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story the task belongs to (US1, US2, US3)

## Phase 1: Setup

- [X] T001 Confirm the tools the script needs, locally: `gh auth status`, `jq --version`, `curl --version`, and bash 3.2 compatibility (`/bin/bash --version`). Record the upstream latest tags (`gh release view --repo github/spec-kit`, `--repo srobroek/speckit-roadmap`) for the record in T020.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: One pin list and a checker that the workflow and local runs share.

- [X] T002 Create `testdata/upstream/pins.json` with both projects, as in research R1: name, repo, tag, and templates[] {path, pinned}. `pinned` is relative to `testdata/upstream/`.
- [X] T003 In `internal/conformance/drift_test.go`, replace the hard-coded `pinnedTemplates` map with a loader for `testdata/upstream/pins.json` (`encoding/json`). It builds the same project → pinned-file list, and every existing use (`pinnedShapes`, the inventory test loop) keeps working.
- [X] T004 Add `TestPinsFile` in `internal/conformance/drift_test.go`. It checks that:
  - every `pinned` file exists;
  - the project names match `testdata/upstream/<name>/resolutions.md`;
  - each `pinned` path contains `<name>/<tag>/`;
  - each `repo` matches `owner/name`.
- [X] T005 Create `scripts/check-upstream-drift.sh` (bash 3.2-portable, `set -euo pipefail`) with the options from contracts/drift-workflow.md: `--project`, `--tag`, `--candidate NAME=FILE`, `--path NAME=REPO_PATH`. For each selected project from `pins.json`:
  - resolve the tag (`gh api repos/R/releases/latest -q .tag_name`, or the given tag) and its commit (`gh api repos/R/commits/TAG -q .sha`);
  - download each template with `curl -fsS --max-time 30` into `mktemp -d` (with a trap that removes it); HTTP 404 means "template not found: PATH", any other failure means "fetch error: MESSAGE";
  - run `go test ./internal/conformance -run 'TestDrift$' -count=1 -v -args -project=P -candidate=F`, and collect the `(unresolved)` `added:`/`removed:` lines;
  - classify the result as `same tag as the pin`, `no unresolved change`, `N unresolved change(s)`, `template not found: PATH` or `fetch error: MESSAGE`.

  Write the Markdown summary to stdout and to `$GITHUB_STEP_SUMMARY` when it is set. Exit 1 if any project failed. Never put candidate content into a command line.

**Checkpoint**: `go test ./internal/conformance` passes, and the script runs locally.

---

## Phase 3: User Story 1 - The owner learns about upstream template changes within a week (Priority: P1) 🎯 MVP

**Goal**: A weekly, read-only run that fails on unresolved drift, with a clear summary.

**Independent Test**: The seeded drift run fails and names the key. A run against the current upstream reports both projects.

### Tests for User Story 1

- [X] T006 [US1] Seeded run: copy `testdata/upstream/spec-kit/v1.2.0/templates/tasks-template.md` to a temporary file and append `- [ ] T999 [P9] seeded`. `scripts/check-upstream-drift.sh --project spec-kit --candidate spec-kit=<copy>` must exit 1 and list an unresolved `added:` key. With the unchanged copy, it must exit 0 with `no unresolved change`.
- [X] T007 [US1] Real run: `scripts/check-upstream-drift.sh` (all projects, latest releases). Record each project's tag, commit and result as found, for the implementation record and the PR. Do not change pins or resolutions in this feature.

### Implementation for User Story 1

- [X] T008 [US1] Create `.github/workflows/drift.yml`:
  - `name: drift`;
  - `on: schedule: [{cron: '17 6 * * 1'}]` and `workflow_dispatch` with inputs `project` (choice: all, spec-kit, speckit-roadmap; default all) and `tag` (string, optional);
  - `permissions: contents: read`;
  - one job on `ubuntu-latest` with `timeout-minutes: 10`: checkout and setup-go at the SHAs ci.yml uses, with `go-version: '1.26.x'`; then run `scripts/check-upstream-drift.sh`, with inputs mapped through `env:` (`DRIFT_PROJECT`, `DRIFT_TAG`) and `GH_TOKEN: ${{ github.token }}`.

  The script reads `DRIFT_PROJECT` and `DRIFT_TAG` when no flags are given (treating `all` and empty as unset). No `${{ inputs.* }}` inside `run:`.
- [X] T009 [US1] Add the env-var handling (`DRIFT_PROJECT`, `DRIFT_TAG`) to `scripts/check-upstream-drift.sh`, and validate the project name against `pins.json` before use.

**Checkpoint**: The workflow YAML parses (actionlint if available; otherwise a `python3 -c "import yaml"` load or `gh workflow view` after merge), and the script handles both flags and env vars.

---

## Phase 4: User Story 2 - A maintainer checks a specific release on demand (Priority: P2)

**Goal**: Check a named project and tag.

**Independent Test**: `--project spec-kit --tag v1.2.0` gives `same tag as the pin`.

- [X] T010 [US2] Seeded run: `scripts/check-upstream-drift.sh --project spec-kit --tag v1.2.0` must exit 0 with `same tag as the pin`. `--tag` without `--project` must exit 1 with a usage error.

---

## Phase 5: User Story 3 - Failures say what went wrong (Priority: P2)

**Goal**: Distinct results for drift, a missing template and a fetch error, without one project hiding another.

**Independent Test**: The seeded runs give the right result line, and no drift is claimed.

- [X] T011 [P] [US3] Seeded run: `--project spec-kit --path spec-kit=templates/nope.md` must exit 1 with `template not found: templates/nope.md`, and no drift lines.
- [X] T012 [P] [US3] Seeded run: `--project spec-kit --tag v0.0.0-nope` must exit 1 with `fetch error: …`, and no drift lines.
- [X] T013 [US3] Seeded run with both projects, where one fails (`--path spec-kit=templates/nope.md`): the summary still reports speckit-roadmap's own result.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T014 [P] Update `testdata/upstream/README.md`:
  - the weekly workflow and how to start it (`gh workflow run drift.yml -f project=spec-kit -f tag=v1.4.3`);
  - the local script;
  - the meaning of each result;
  - that a failed run is the alert;
  - the 60-day pause and `gh workflow enable drift.yml`;
  - that pins move only through the manual procedure, with `pins.json` as the single list to edit (step 4).
- [X] T015 [P] Update `specs/002-upstream-format-conformance/contracts/drift-report.md` with a section "Automated weekly check (feature 013)", linking `specs/013-upstream-drift-workflow/contracts/drift-workflow.md`.
- [X] T016 [P] Add the script line to the Commands block in `CLAUDE.md`, and note that `pins.json` is the pin list.
- [X] T017 [P] Check `.github/workflows/drift.yml` against the repository's workflow conventions: the header comment states purpose and that it publishes nothing, and the action pins match ci.yml. If ci.yml or any test enumerates workflows (for example an action-pin or permissions check), extend it to cover drift.yml.
- [X] T018 Run `CGO_ENABLED=0 go test ./...` and `go vet ./...`; `shellcheck scripts/check-upstream-drift.sh` if installed; and `bash -n` on the script. `spec status` on this roadmap must still exit 0.
- [X] T019 Re-run T006 and T010–T013 after any change, and confirm `git status` shows no files written by the script outside its temporary directory.
- [X] T020 [P] Add a row "Upstream drift workflow (feature 013, roadmap 010)" to the implementation record in `specs/001-spec-status/tasks.md`. Include the seeded-run results, the T007 real results (tags, commits, outcome) and "first scheduled/dispatched CI run happens after merge". Mark it "executed by the implementer, not independently verified".

---

## Dependencies & Execution Order

### Phase Dependencies

- T001, then T002–T005 (T003 and T004 depend on T002; T005 depends on T002).
- US1 needs Phase 2. US2 and US3 need the script (T005) and can run in any order after it.
- Polish comes last. T019 repeats the seeded runs after the final edits.

### Parallel Opportunities

- T011 and T012 run in parallel (independent invocations).
- T014–T017 and T020 are different files.

## Parallel Example: Polish

```text
T014 testdata/upstream/README.md
T015 specs/002-upstream-format-conformance/contracts/drift-report.md
T016 CLAUDE.md
T020 specs/001-spec-status/tasks.md
```

## Implementation Strategy

1. Pin list and script (Phase 2), proven by the seeded drift and the real run (US1 tests).
2. The workflow (T008–T009): the MVP for roadmap item 010.
3. Manual tags (US2) and distinct failure kinds (US3).
4. Docs and record. Doomer review if available, then one PR against `main`. The first CI dispatch runs after the merge.
