# Implementation Plan: Dependency update automation

**Branch**: `003-dependency-updates` | **Date**: 2026-11-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/003-dependency-updates/spec.md`

## Summary

This feature adds repository tooling only:
- a Dependabot configuration for Go modules and GitHub Actions;
- a CI workflow that runs tests, `go vet` and `govulncheck` on pull requests and `main`, on linux/amd64 and darwin/arm64 and on both supported Go release lines;
- a hygiene job (gofmt, `go mod tidy -diff`, `go mod verify`, `spec status --fail-on error` on this repository, build-only checks for windows/amd64 and linux/arm64) and a dependency-review job on pull requests;
- CodeQL code scanning (default setup, a repository setting);
- the merge rule (D-2026-11-11-04), written into CONTRIBUTING.md;
- a dependency-vulnerability section in SECURITY.md;
- truthful README and build-release wording.

Planning found that the declared `go 1.22` makes the vulnerability scan fail on every pull request. The owner therefore moved the minimum Go version into this feature: `go.mod` becomes `go 1.26.0` (D-2026-11-11-05, research R1). The `spec` binary's behavior and dependency set are unchanged.

## Technical Context

**Language/Version**:
- Go. `go.mod` moves from `go 1.22` to `go 1.26.0`.
- CI tests the latest `1.26.x` and `stable` (1.27.x today).

**Primary Dependencies**: None added to the module. Development tooling:
- `actions/checkout` v7.0.1, `actions/setup-go` v7.0.0 and `actions/dependency-review-action` v5.0.0, all pinned by commit SHA;
- CodeQL default setup (repository setting);
- `golang.org/x/vuln/cmd/govulncheck` v1.9.0, run with `go run pkg@version` so it stays outside `go.mod`;
- GitHub Dependabot.

**Storage**: N/A

**Testing**:
- The existing `CGO_ENABLED=0 go test ./...` and `go vet ./...` run in CI.
- The workflow and Dependabot files are checked by the validation steps in [quickstart.md](quickstart.md) and the contracts in [contracts/](contracts/).

**Target Platform**: GitHub-hosted runners `ubuntu-latest` (asserted linux/amd64) and `macos-latest` (asserted darwin/arm64).

**Project Type**: CLI (existing). This feature touches only repository configuration and documentation, plus the `go` line of `go.mod`.

**Performance Goals**: A full CI run should finish in about 5 minutes on each job. This is not a requirement; the check only needs to report before merge.

**Constraints**:
- No runtime change.
- Read-only workflow token.
- No secrets.
- SHA-pinned actions.
- No release or publish steps.
- No new module in `go.mod` without an owner decision.

**Scale/Scope**:
- New files: 2 (`.github/dependabot.yml`, `.github/workflows/ci.yml`).
- Edited files: about 6 (`go.mod`, possibly `go.sum`, SECURITY.md, CONTRIBUTING.md, README.md, build-release.md), plus records (tasks.md implementation record, decisions, roadmap).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Article | Check | Result |
|---|---|---|
| I Truthful evidence | Docs state that CI exists only for tests and scans, that its runs are not release evidence, and that fixes reach users only through releases (none exist yet). | Pass |
| II Deterministic, standalone reporting; dependency allowlist | No product code or import changes. govulncheck stays outside `go.mod`. The merge rule keeps new modules behind owner decisions. Raising the `go` line is an owner decision (D-2026-11-11-05), not a new dependency. | Pass |
| III Narrow, safe input | Unaffected. | Pass |
| IV History cannot become false completion | Unaffected. | Pass |
| V Compatibility and findings | CLI, exit codes and JSON are unchanged. Building now needs Go ≥1.26. That is a build-environment change, recorded in the decision log and README. There is no released compatibility promise yet. | Pass |
| VI No inferred quarters | Unaffected. | Pass |
| VII Distribution claims need evidence | No supported-platform, signing, audit or certification claim. CI passes are described as CI evidence only. | Pass |
| VIII Explicit future mutation authority | The workflow cannot write: `contents: read`, no auto-merge. | Pass |
| Governance | No amendment needed (research R1, spec Assumptions). | Pass |

**Post-design re-check (after Phase 1):** still passes. The contracts add no product behavior. The only module-file change is the `go` line, plus whatever `go mod tidy` does. That is reviewed under R7, and it stops for a decision if a new module appears.

## Project Structure

### Documentation (this feature)

```text
specs/003-dependency-updates/
├── spec.md
├── plan.md              # this file
├── research.md          # R1–R11
├── data-model.md        # configuration entities and the merge-rule classification
├── quickstart.md        # validation steps
├── contracts/
│   ├── ci-workflow.md   # required jobs, steps, permissions, failure behavior
│   └── dependabot.md    # required ecosystems, schedule, grouping
├── checklists/requirements.md
└── tasks.md             # /speckit-tasks
```

### Source Code (repository root)

```text
.github/
├── dependabot.yml           # new
└── workflows/
    └── ci.yml               # new
go.mod                       # `go 1.22` → `go 1.26.0`
go.sum                       # only if `go mod tidy` changes it (reviewed, R7)
SECURITY.md                  # + dependency vulnerabilities section
CONTRIBUTING.md              # + merge rule, Go version
README.md                    # CI and Go-version wording
specs/001-spec-status/build-release.md   # CI exists; not a release workflow; minimum Go decided
specs/001-spec-status/decisions.md       # D-2026-11-11-05
specs/001-spec-status/tasks.md           # implementation-record row
.specify/memory/roadmap.md               # 009 status, 004 notes
```

Product code is not touched: `cmd/`, `internal/` and `testdata/` stay as they are.

## Complexity Tracking

No constitution violations.
