# Tasks: Dependency update automation

**Input**: Design documents from `specs/003-dependency-updates/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/), [quickstart.md](quickstart.md)

**Tests**:
- The spec asks for no new Go tests. Validation is the existing suite, the contract checks in [quickstart.md](quickstart.md), and live CI runs (seeded-fault runs included).
- No product code changes, so there is no test-first step for product behavior.

**Organization**: by user story from spec.md. US1 and US2 share `.github/workflows/ci.yml`, so their workflow tasks run in sequence. US2 creates the file, and US1 and US6 add jobs and steps to it.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: can run in parallel (different files, no dependency on an unfinished task)
- **[Story]**: US1–US6 from spec.md

## Path conventions

All paths are relative to the repository root. New files go under `.github/`. Product code (`cmd/`, `internal/`, `testdata/`) is not touched.

---

## Phase 1: Setup

**Purpose**: Confirm the starting state and record facts that the later steps check against.

- [X] T001 Confirm the branch `003-dependency-updates` is up to date with `main` (after #26 and #27 merge, rebase or recreate from `main`), then build a baseline binary from `main` to scratchpad `spec-before` under go1.26.11 for the SC-005 comparison (quickstart §2)
- [X] T002 Re-check, and record in the PR description, that the pinned versions in research.md R4/R5/R9 are still the latest: `actions/checkout` v7.0.1 `3d3c42e5aac5ba805825da76410c181273ba90b1`, `actions/setup-go` v7.0.0 `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e`, `actions/dependency-review-action` v5.0.0 `a1d282b36b6f3519aa1f3fc636f609c47dddb294`, `golang.org/x/vuln` v1.9.0. If any moved, update research.md and use the new value

---

## Phase 2: Foundational (blocks every story)

**Purpose**: The minimum Go version (FR-015, D-2026-11-11-05). Every CI job asserts it.

- [X] T003 Set the minimum Go version in `go.mod` with `GOTOOLCHAIN=go1.26.11 go mod edit -go=1.26.0`, then run `GOTOOLCHAIN=go1.26.11 go mod tidy`
- [X] T004 Review `git diff -- go.mod go.sum` (research R7). Only the `go` line and tidy's rewrites are allowed. If any module path not already in `go.mod` appears, STOP: under D-2026-11-11-04 that needs an owner decision before continuing
- [X] T005 Under go1.26.11, run `CGO_ENABLED=0 go test -count=1 ./...`, `go vet ./...` and `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...`. Expect all packages ok, no vet output, and "No vulnerabilities found." If govulncheck reports a reachable advisory, stop and report it

**Checkpoint**: The module builds and passes on Go 1.26. The stories can start.

---

## Phase 3: User Story 2 — Every pull request is tested before merge (P1) 🎯 MVP

**Goal**: Four test jobs (linux/amd64 and darwin/arm64 × Go `1.26.x` and `stable`) on every pull request to `main` and every push to `main`.

**Independent test**: Open a PR. Four test checks appear and pass. Push a commit that breaks one test: all four fail. Revert it: all four pass (quickstart §3).

- [X] T006 [US2] Create `.github/workflows/ci.yml` per [contracts/ci-workflow.md](contracts/ci-workflow.md):
  - **Name and triggers:** `name: ci`; `on: pull_request: branches: [main]` and `push: branches: [main]`.
  - **Permissions:** workflow-level `permissions: contents: read`.
  - **Job `test`:** `strategy.fail-fast: false`; matrix `os: [ubuntu-latest, macos-latest]` and `go: ['1.26.x', 'stable']`.
  - **Checkout step:** `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1` with `persist-credentials: false`.
  - **Go setup step:** `actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e # v7.0.0` with `go-version: ${{ matrix.go }}`.
  - **Rules:** no `secrets.*`, no `pull_request_target`, no `continue-on-error`.
- [X] T007 [US2] In `.github/workflows/ci.yml` job `test`, add a platform-assertion step. It must fail unless `go env GOOS GOARCH` prints `linux amd64` on `ubuntu-latest` and `darwin arm64` on `macos-latest` (research R3).
- [X] T008 [US2] In `.github/workflows/ci.yml` job `test`, add a minimum-assertion step that fails unless `grep -E '^go 1\.26(\.|$)' go.mod` matches (research R2, FR-015).
- [X] T009 [US2] In `.github/workflows/ci.yml` job `test`, add the steps `go vet ./...` and `CGO_ENABLED=0 go test ./...`, in that order, after the assertions.

**Checkpoint**: Pull requests get four tested results.

---

## Phase 4: User Story 1 — Learn about a vulnerable dependency without looking (P1)

**Goal**: Reachable advisories fail CI. Repository alerts and security-fix PRs are on.

**Independent test**: Repository settings return `204` for alerts and `enabled:true` for automated security fixes. govulncheck runs in every test job. If an x/sys alert or PR appears, it's classified under the merge rule (quickstart §3).

- [X] T010 [US1] In `.github/workflows/ci.yml` job `test`, add a final step `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...`. Use the exact version, never `@latest`, and no `continue-on-error` (FR-006, research R4).
- [X] T011 [P] [US1] Verify the repository settings (FR-003, [contracts/dependabot.md](contracts/dependabot.md)):
  - `gh api repos/tosin2013/spec-status/vulnerability-alerts -i` returns `204`.
  - `gh api repos/tosin2013/spec-status/automated-security-fixes` returns `"enabled":true`.
  - Record the open alerts with `gh api repos/tosin2013/spec-status/dependabot/alerts --jq '.[].security_advisory.ghsa_id'` in the PR description.

**Checkpoint**: Vulnerabilities are detected in CI and by repository alerts.

---

## Phase 5: User Story 3 — Routine updates arrive in a manageable batch (P2)

**Goal**: Weekly grouped update proposals for Go modules and actions, with majors kept separate.

**Independent test**: The file meets [contracts/dependabot.md](contracts/dependabot.md). After merge, Insights → Dependency graph → Dependabot lists both ecosystems.

- [X] T012 [P] [US3] Create `.github/dependabot.yml` with `version: 2` and exactly two `updates` entries:
  - **Go modules:** `package-ecosystem: gomod`, `directory: /`, `schedule: {interval: weekly, day: monday}`, `groups: {go-minor-patch: {patterns: ["*"], update-types: [minor, patch]}}`.
  - **Actions:** `package-ecosystem: github-actions`, `directory: /`, `schedule: {interval: weekly, day: monday}`, `groups: {actions: {patterns: ["*"]}}`.
  - **Leave out:** auto-merge, `ignore`, registries, `insecure-external-code-execution`, and any `applies-to: security-updates` group.

**Checkpoint**: Scheduled updates are configured.

---

## Phase 6: User Story 4 — Merges keep the dependency allowlist honest (P2)

**Goal**: The merge rule (D-2026-11-11-04) is written where contributors will read it.

**Independent test**: Apply the rule in CONTRIBUTING.md to "bump an approved module, no new module" (merge after review) and "adds an indirect module" (decision needed). The [data-model.md](data-model.md) table gives the same answers.

- [X] T013 [P] [US4] Add a "Dependency updates" section to `CONTRIBUTING.md`:
  - **Merge rule** (copy the outcomes table from data-model.md):
    - Only versions of modules already in `go.mod` change, the `go` line is unchanged, and CI passes: merge after owner approval.
    - Any module path added (direct or `// indirect`), or the `go` line raised: a decision recorded in `specs/001-spec-status/decisions.md` first.
    - CI fails: do not merge.
    - Removing a module needs no decision. Action updates are merge after review once CI passes.
  - **Go version:** the minimum is Go 1.26 (D-2026-11-11-05).
  - **Build and test section:** update it to say Go 1.26 or newer is required.

**Checkpoint**: The rule is published.

---

## Phase 7: User Story 5 — Reporters and users know how dependency vulnerabilities are handled (P3)

**Goal**: SECURITY.md explains how dependency vulnerabilities are detected and fixed, and when a fix reaches users, without overstating it.

**Independent test**: SECURITY.md answers the three SC-006 questions and makes none of the claims FR-013 forbids.

- [X] T014 [P] [US5] Add a "Dependency vulnerabilities" section to `SECURITY.md`:
  - **Detection:** Dependabot alerts, govulncheck in CI (reachable code only), dependency review on pull requests, and CodeQL code scanning.
  - **Handling:** fixes arrive as pull requests, are reviewed by the maintainer under the merge rule, and land on `main`.
  - **Delivery:** a fix reaches users only in a release. No release exists yet, so fixes are on `main` until the first release.
  - **Limits:** no audit, certification, signing or response-time claim.
  - **Scope section:** extend it to include vulnerable dependencies.

**Checkpoint**: Disclosure is truthful.

---

## Phase 8: User Story 6 — Hygiene and security checks catch what tests miss (P2)

**Goal**: A hygiene job, a dependency-review job, and CodeQL code scanning.

**Independent test**: SC-007. On a scratch branch, a misformatted file, then an untidy `go.mod`, then a broken roadmap field each fail only `hygiene`, naming the cause. A PR that changes modules shows the dependency-review summary. The code-scanning check appears.

- [X] T015 [US6] Add job `hygiene` to `.github/workflows/ci.yml` (contract "Job hygiene", FR-016 to FR-018):
  - **Setup:** `runs-on: ubuntu-latest`, with the same pinned checkout (`persist-credentials: false`) and setup-go using `go-version: '1.26.x'`.
  - **Steps, each its own named step:**
    1. `test -z "$(gofmt -l .)" || { gofmt -l .; exit 1; }`.
    2. `go mod tidy -diff`.
    3. `go mod verify`.
    4. `CGO_ENABLED=0 go build -o spec ./cmd/spec && ./spec status --color=never --fail-on error`.
    5. `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./...`.
    6. `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ./...`.
- [X] T016 [US6] Add job `dependency-review` to `.github/workflows/ci.yml` (FR-019):
  - `if: github.event_name == 'pull_request'`, `runs-on: ubuntu-latest`.
  - Pinned checkout (`persist-credentials: false`), then `actions/dependency-review-action@a1d282b36b6f3519aa1f3fc636f609c47dddb294 # v5.0.0` with `fail-on-severity: low`.
  - No `comment-summary-in-pr`, and no write permissions.
- [X] T017 [P] [US6] Enable CodeQL default setup for Go (FR-020, research R10). The owner chose this on 2026-11-11.
  - Run `gh api -X PATCH repos/tosin2013/spec-status/code-scanning/default-setup -f state=configured`.
  - Verify that `gh api repos/tosin2013/spec-status/code-scanning/default-setup --jq .state` returns `configured`.
  - Record the result in the PR.

**Checkpoint**: All six workflow checks and code scanning exist.

---

## Phase 9: Polish and cross-cutting

**Purpose**: Truthful docs, records, and end-to-end validation.

- [X] T018 [P] Update `README.md` (FR-014). Say that CI runs tests, `go vet` and govulncheck on linux/amd64 and darwin/arm64 with Go 1.26.x and stable. Say it also runs hygiene checks, including compile-only builds for windows/amd64 and linux/arm64 that make no support claim, and dependency review. State that CI is not a release workflow and its passes are not release evidence, and that building needs Go 1.26 or newer.
- [X] T019 [P] Update `specs/001-spec-status/build-release.md` (FR-014):
  - A test-and-scan workflow exists, which the release workflow (004) may reuse.
  - The minimum Go version is decided (D-2026-11-11-05), and only the release-build patch stays open.
  - Compile-only builds are not support.
- [X] T020 [P] Add an implementation-record row to `specs/001-spec-status/tasks.md`: "Dependency automation and CI (feature 003)". State the CI checks, the Dependabot configuration, CodeQL, Go 1.26, and that it was run by the implementer, with run URLs.
- [X] T021 Run quickstart §1 locally: the go line, gofmt, tidy -diff and verify, own-roadmap `spec status --fail-on error`, both cross-builds, then tests, vet and govulncheck under go1.26.11. Then run the workflow-file greps: every `uses:` has a 40-hex SHA, and there is no `secrets.`, `pull_request_target` or `continue-on-error`.
- [X] T022 Run quickstart §2 (SC-005). Build `spec-after` under go1.26.11. Then `cmp` the `--json` and `--color=never` text output against `spec-before` for `testdata/example-roadmap`, `testdata/speckit-roadmap-upstream` and every `testdata/conformance/*` fixture. Every result must be byte-identical.
- [X] T023 Commit, push and open or update the PR. Confirm six checks pass: four test cells, `hygiene`, and `dependency-review`. Record the run URLs in the PR and in the T020 row.
- [X] T024 Seeded-fault runs (SC-007 and US2's independent test), on throwaway commits on a scratch branch, reverted afterwards. Record each run URL:
  1. Break one test: all four test cells fail.
  2. Misformat a file: only `hygiene` fails.
  3. Break a roadmap field: only `hygiene` fails, showing the diagnostic.
- [X] T025 Update `.specify/memory/roadmap.md`: 009 becomes `implemented`, and its notes cite the PR and run URLs. Item 004's notes say that the CI check exists and the minimum Go version is decided. Bump the patch version with a sync note.
- [ ] T026 Offer the owner the optional branch-protection step (quickstart §4), which requires the CI checks before merging to `main`. Do not apply it without explicit approval.

---

## Dependencies and execution order

- **Phase 1 → Phase 2 → stories.** T003–T005 must pass first, because every job asserts Go 1.26.
- **Shared file `.github/workflows/ci.yml`:** run T006 → T007 → T008 → T009 → T010 → T015 → T016 in that order (US2, then US1, then US6).
- **Independent of the workflow file, so parallel:** T011 (settings check), T012 (dependabot.yml), T013 (CONTRIBUTING.md), T014 (SECURITY.md), T017 (CodeQL setting), and T018–T020 (docs).
- **Polish:** T021 and T022 come after all file tasks. T023 needs T021 and T022. T024 needs T023. T025 needs T023 and T024. T026 comes last.

### Story completion order

US2 (MVP) → US1 → US6, along the workflow-file chain. US3, US4 and US5 can land at any time after Phase 2.

## Parallel example

After T005, run these in parallel:

```text
T011 [US1] verify repository security settings
T012 [US3] create .github/dependabot.yml
T013 [US4] CONTRIBUTING.md dependency-updates section
T014 [US5] SECURITY.md dependency-vulnerabilities section
T017 [US6] enable CodeQL default setup
```

Meanwhile, one sequential lane builds `.github/workflows/ci.yml`: T006 → T010, then T015 and T016.

## Implementation strategy

1. **MVP**: Phase 2 plus US2 (T003–T009). Every PR is tested on both platforms and both Go lines.
2. **Then**: US1 (scan in CI), US3 (updates), US4 (rule), US5 (disclosure), US6 (hygiene and security).
3. **One pull request**: all of it lands in a single PR, because the pieces depend on each other. For example, the first Dependabot PRs need the CI check and the merge rule to already exist.
4. **Stop conditions**: T004 (a new module after tidy) and T005 (a reachable advisory on Go 1.26) stop the work and go back to the owner.
