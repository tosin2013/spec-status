# Implementation Plan: Release engineering

**Branch**: `004-release-engineering` (work branch `004-plan`) | **Date**: 2026-11-12 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/004-release-engineering/spec.md`

## Summary

This feature adds two things:
- **A tag-triggered `release.yml`.** It runs these jobs in order:
  1. Verify the tag: SemVer, on `main`, and no existing release.
  2. Reuse CI, through `workflow_call` on `ci.yml`.
  3. Build linux/amd64 and darwin/arm64 archives and one checksum file, with release stamping.
  4. Smoke-test each downloaded archive on its own platform.
  5. Create a **draft** GitHub Release with `gh`, in the only job with write permission.
- **A small product change.** It distinguishes release, module (`go install`) and development builds in `spec version`, `--help` and the status panel heading, using `runtime/debug`. JSON is unchanged.

The docs gain a README install section, the pre-1.0 policy and an owner release checklist.

## Technical Context

**Language/Version**:
- Go. The minimum is 1.26 (D-2026-11-11-05).
- Release binaries are built with the current stable Go (`check-latest`), and the exact version is recorded.

**Primary Dependencies**:
- No new modules. `runtime/debug` is standard library.
- Workflow components, all pinned by commit SHA:
  - `actions/checkout` v7.0.1, `actions/setup-go` v7.0.0;
  - `actions/upload-artifact` v7.0.2 and `actions/download-artifact` v8.0.2, both new;
  - the `gh` CLI preinstalled on runners.

**Storage**: N/A

**Testing**:
- New Go tests in `internal/cli` cover the three build kinds (FR-012). An existing status test is extended for the `View.Build` heading.
- The release workflow is validated by a release candidate run (quickstart) and by negative runs: a bad tag, and a failing commit.

**Target Platform**: linux/amd64 and darwin/arm64 archives (D-2026-11-11-02), smoke-tested on `ubuntu-latest` and `macos-latest`.

**Project Type**: CLI (existing).

**Performance Goals**: One release run should produce a draft in roughly 10 minutes. This is not a gate.

**Constraints**:
- **Permissions:** only the release job may write; everything else is read-only. No secrets beyond `github.token`.
- **Publishing:** the workflow never publishes or replaces assets.
- **Bytes:** the tested bytes are the attached bytes.
- **Output:** JSON is byte-identical across build kinds. Development-build text output is byte-identical to today's, so the README sample and existing tests are unchanged.
- **Claims:** no signing, provenance, reproducibility or audit claims.

**Scale/Scope**:
- New: `.github/workflows/release.yml`, `internal/cli/build.go`, `internal/cli/build_test.go`, `RELEASING.md`.
- Edited: `ci.yml` (adds `workflow_call`), `internal/cli/cli.go`, `internal/status/style.go`, `README.md`, `CONTRIBUTING.md`, the build-release.md and tasks.md records, and the roadmap.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Article | Check | Result |
|---|---|---|
| I Truthful evidence | Labels state exactly how a binary was built. Development builds still say so. Release notes link smoke-run evidence. No claim beyond it. | Pass |
| II Deterministic, standalone; dependency allowlist | `runtime/debug` only, no new module. Output stays deterministic: the label depends only on build-time inputs, not on time or the environment. | Pass |
| III Narrow, safe input | Unaffected. | Pass |
| IV History cannot become false completion | Unaffected. | Pass |
| V Compatibility and findings | JSON 1.0 bytes, exit codes and diagnostics are unchanged. Human text changes only in the build label, and only for release and module builds. The pre-1.0 policy is documented. | Pass |
| VI No inferred quarters | No dates in output. Release notes carry no schedule. | Pass |
| VII Distribution claims need evidence | A platform is listed as supported only after the published release's smoke test passes on it (FR-016). Unsupported platforms are named. No signing claim. | Pass |
| VIII Explicit future mutation authority | The workflow cannot publish. Only the owner publishes. | Pass |
| Governance | No amendment needed (spec Assumptions). | Pass |

**Post-design re-check (after Phase 1):** still passes. [contracts/build-identity.md](contracts/build-identity.md) keeps the development-build bytes unchanged. [contracts/release-workflow.md](contracts/release-workflow.md) confines write permission to one job and leaves publishing to the owner.

## Project Structure

### Documentation (this feature)

```text
specs/004-release-engineering/
├── spec.md
├── plan.md                 # this file
├── research.md             # R1–R9
├── data-model.md           # build identity, archive, checksum file, release, workflow jobs
├── quickstart.md           # local identity checks, RC run, negative runs, publish
├── contracts/
│   ├── build-identity.md   # exact version/help/panel text per build kind
│   ├── release-workflow.md # jobs, triggers, permissions, steps, failure behavior
│   └── release-assets.md   # archive names, layout, checksum format, notes template
├── checklists/requirements.md
└── tasks.md                # /speckit-tasks
```

### Source Code (repository root)

```text
.github/workflows/
├── ci.yml          # + workflow_call trigger (no other change)
└── release.yml     # new
internal/cli/
├── build.go        # new: identity resolution (release | module | development)
├── build_test.go   # new
└── cli.go          # versionText, rootLong and View.Build use the identity; adds `Kind` stamp variable
internal/status/
└── style.go        # View.Build field; heading uses it, empty → "development build"
README.md           # install section, upgrade/rollback, Gatekeeper, go install, pre-1.0 policy
RELEASING.md        # owner release checklist (FR-017)
CONTRIBUTING.md     # link to RELEASING.md
specs/001-spec-status/build-release.md, specs/001-spec-status/tasks.md, .specify/memory/roadmap.md
```

## Complexity Tracking

No constitution violations.
