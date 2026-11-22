# Implementation Plan: Upstream drift workflow

**Branch**: `013-upstream-drift-workflow` | **Date**: 2026-11-21 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/013-upstream-drift-workflow/spec.md` (roadmap item 010, decision D-2026-11-21-01 (3), target v0.4.0).

## Summary

- **One pin list:** `testdata/upstream/pins.json` names, for each project, the GitHub repository, the pinned tag, and each template's repository path and pinned copy (R1). `internal/conformance/drift_test.go` reads it instead of the hard-coded `pinnedTemplates` map, so the test and the workflow cannot disagree (FR-009). A new test checks that every listed pinned copy exists.
- **The checker:** `scripts/check-upstream-drift.sh` holds the logic, so it runs the same way locally and in CI (R2). For each selected project it:
  - resolves the candidate tag (the latest non-prerelease release, or the given tag) and its commit with `gh api`;
  - downloads each template from `raw.githubusercontent.com/<repo>/<tag>/<path>` into a fresh temporary directory;
  - runs the existing `TestDrift` per file;
  - classifies the result as `no unresolved change`, `same tag as the pin`, `N unresolved changes`, `template not found` or `fetch error`;
  - writes a Markdown summary (to `$GITHUB_STEP_SUMMARY` in CI, otherwise to stdout).

  It exits 1 if any project failed.
- **The workflow:** `.github/workflows/drift.yml` runs on a weekly `schedule` and on `workflow_dispatch`, with inputs `project` and `tag`. It has `permissions: contents: read`, actions pinned by SHA (the same pins as ci.yml), and runs Go 1.26.x and the script (R3).
- **Docs:** `testdata/upstream/README.md`, the drift-report contract, and the CLAUDE.md commands.

## Technical Context

**Language/Version**: Go 1.26 (tests only), bash 3.2+ (portable to macOS), GitHub Actions YAML

**Primary Dependencies**: `gh` and `jq` (present on GitHub-hosted ubuntu runners and used locally), `curl`, and the Go toolchain

**Storage**: none; candidates go in `mktemp -d`, removed on exit

**Testing**:
- `go test ./internal/conformance` (the pins file is read and the copies exist; existing drift tests);
- seeded script runs: a local candidate override with a changed shape gives a drift, a bad path gives `template not found`, a bad tag gives `fetch error`, and the current pins give `same tag as the pin`;
- one real `workflow_dispatch` run on the branch.

**Target Platform**: GitHub-hosted `ubuntu-latest`; the script also runs on macOS

**Project Type**: CI tooling (no product code change)

**Performance Goals**: a few seconds of network and test time per run

**Constraints**:
- Art. II: no runtime network. Only CI and the development script fetch, and the `spec` binary and `go test ./...` never do.
- Read-only permissions; no repository writes (FR-006, FR-007).
- Untrusted candidates are only read as text (FR-008).

**Scale/Scope**: about 150 lines of bash, 40 lines of YAML, a small refactor of `drift_test.go`, and docs

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Article | Check | Result |
|---|---|---|
| I Truthful evidence | Reports the drift test's own findings and labels failure kinds distinctly; never claims drift on a fetch error | Pass |
| II Standalone, no side effects | Network only in a CI workflow and a dev script. `spec` and `go test ./...` stay offline (TestDrift is skipped without `-candidate`). No product import changes. | Pass |
| III Narrow, safe input | Candidates are untrusted text, read only by the existing test. The pin moves only through the manual procedure. | Pass |
| IV History | Not affected | Pass |
| V Compatibility | No CLI, JSON or Go API change | Pass |
| VI No quarters | Not affected | Pass |
| VII Evidence-gated distribution | No release claim; CI evidence only | Pass |
| VIII Approval | No automatic mutation; pins move only by an owner-approved PR | Pass |

Post-design re-check: unchanged, all pass.

## Project Structure

### Documentation (this feature)

```text
specs/013-upstream-drift-workflow/
├── spec.md, plan.md, research.md (R1–R6), data-model.md, quickstart.md
├── contracts/drift-workflow.md
└── checklists/requirements.md
```

### Source Code (repository root)

```text
testdata/upstream/pins.json                 # new: the single pin list
internal/conformance/drift_test.go          # reads pins.json; TestPinsFileMatchesCopies
scripts/check-upstream-drift.sh             # new: resolve, download, run TestDrift, summarize
.github/workflows/drift.yml                 # new: weekly + dispatch, read-only
testdata/upstream/README.md                 # workflow, results, 60-day pause
specs/002-upstream-format-conformance/contracts/drift-report.md   # automated invocation
CLAUDE.md; specs/001-spec-status/tasks.md   # command line; implementation record
```

**Structure Decision**: The logic goes in a script and the workflow stays thin, following `scripts/check-extension-install.sh`. Pins are data in `testdata/upstream/`, beside the pinned copies.

## Complexity Tracking

None.
