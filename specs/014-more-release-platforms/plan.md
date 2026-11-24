# Implementation Plan: Release archives for linux/arm64 and darwin/amd64

**Branch**: `014-more-release-platforms` | **Date**: 2026-11-24 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/014-more-release-platforms/spec.md` (roadmap 019, issue #82, D-2026-11-24-01 (1) and (3); target v0.4.1).

## Summary

This is a workflow and docs change only; no Go code changes.

- **`release.yml`:**
  - The build loop gains `linux/arm64` and `darwin/amd64`. The existing checksum glob covers them.
  - The smoke matrix becomes four native runners, and its platform `case` gains `Linux aarch64` and `Darwin x86_64` (R2).
  - The draft job maps each smoke job by name to its archive line, and attaches all four archives. It fails if any expected smoke job URL is missing (R3).
- **`ci.yml`:**
  - The test matrix gains `ubuntu-24.04-arm` and `macos-15-intel`.
  - The "Assert platform" step checks both the Go toolchain (`go env GOOS GOARCH`) and the hardware (`uname -s -m`), so a toolchain running under emulation cannot pass for a native run (R4).
- **Docs and contracts:** they list four platforms. The extension version becomes 0.4.1.

## Technical Context

**Language/Version**: GitHub Actions YAML, bash. Go 1.26 for the unchanged build.

**Primary Dependencies**: GitHub-hosted runners `ubuntu-latest`, `ubuntu-24.04-arm`, `macos-latest`, `macos-15-intel`

**Storage**: none

**Testing**:
- actionlint;
- `go test ./...` (TestReadmeSamples, TestDocsInstallCommandsAreNotPinned, extension version tests);
- the PR's own CI, where the new test legs run;
- the v0.4.1 release run, where the four smoke legs run.

**Target Platform**: release assets for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64

**Project Type**: CI and release tooling, plus docs

**Performance Goals**: two extra CI jobs per PR and two extra smoke jobs per release

**Constraints**:
- Art. VII and C-06: advertise a platform only after a native smoke test (D-2026-11-11-02).
- No CLI, JSON or library change (FR-009).
- Install lines stay unpinned (TestDocsInstallCommandsAreNotPinned).

**Scale/Scope**: about 40 changed YAML lines and about 15 doc and contract edits

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Article | Check | Result |
|---|---|---|
| I Truthful evidence | Each archive line links its own native smoke job. The platform asserts fail on a mismatch. | Pass |
| II Standalone | No runtime change. The build stays `CGO_ENABLED=0`. | Pass |
| III Narrow input | Unchanged | Pass |
| IV History | Unchanged | Pass |
| V Compatibility | No flag, JSON or exit-code change. The assets are additive (patch release). | Pass |
| VI No quarters | Not affected | Pass |
| VII Evidence-gated distribution | No archive is drafted without its native smoke test (FR-004). No signing or provenance claims. | Pass |
| VIII Approval | The owner tags and publishes | Pass |

Post-design re-check: unchanged, all pass.

## Project Structure

```text
specs/014-more-release-platforms/   spec, plan, research (R1–R5), data-model, quickstart, contracts/platforms.md
.github/workflows/release.yml       build loop, smoke matrix and case, draft job
.github/workflows/ci.yml            test matrix, platform assert
README.md, docs/{agent-prompt,getting-started,architecture,development,limits}.md, RELEASING.md
specs/004-release-engineering/contracts/{release-assets,release-workflow}.md
extensions/spec-status/{extension.yml,README.md}   version 0.4.1
specs/001-spec-status/tasks.md      implementation record
```

## Complexity Tracking

None.
