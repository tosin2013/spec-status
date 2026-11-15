# Implementation Plan: Close-gap skill

**Branch**: `011-013-specify` | **Date**: 2026-11-15 | **Spec**: [spec.md](spec.md)

**Input**: `specs/010-close-gap-skill/spec.md` (roadmap item 013, D-2026-11-14-01 (1), v0.3.0; depends on 012 and 014).

## Summary

Add a sixth extension command, `speckit.spec-status.close-gap`, that classifies every active non-verified item into exactly one cause (precedence in research R1), reports the breakdown and an order of work (R2), then proposes one change at a time and applies it only after approval through the existing Spec Kit commands, re-reading `spec status` after each (R3, R4). Extend `review`'s `pm` role with a `verification-gap` section and close-gap as its first follow-up (R5). Bump the extension to 0.3.0 with six skills everywhere they are counted (R6), and add an all-causes fixture (R7). No change to the `spec` binary.

## Technical Context

**Language/Version**: Markdown command files (Spec Kit 1.2.0 extension format); Go 1.26 tests in `internal/extension` (standard library)
**Primary Dependencies**: Spec Kit `>=1.2.0` (extension host), `spec >=0.1.2` (tool), optional partners unchanged
**Storage**: none (the skill keeps no state)
**Testing**: `go test ./internal/extension`; `scripts/check-extension-install.sh`; fixture Signals checked against real `spec`
**Target Platform**: any agent that renders Spec Kit extension commands as skills
**Project Type**: Spec Kit extension (in-repo, tested by Go)
**Performance Goals**: n/a
**Constraints**: C-01 (verified is a label), C-06 / Art. VIII (per-change approval, Spec Kit commands only), C-04 (no redirect of superseded dependencies), C-07 (blocked until dependency verified)
**Scale/Scope**: one new command file, edits to roles reference, review command, manifest, two READMEs, install script, release workflow, tests, one fixture

## Constitution Check

| Article / constraint | Check | Result |
|---|---|---|
| I Truthful evidence (C-01) | causes come only from `spec` output; `verified` proposed only after debrief and agreement; "not confirmed" reported honestly | Pass |
| II Standalone | `spec` binary unchanged; skill is development tooling (CLAUDE.md), not a runtime dependency | Pass |
| IV History (C-04) | superseded items excluded; superseded dependencies routed to human review, never redirected | Pass |
| VIII Mutation authority (C-06) | one proposal per approval; applied only through Spec Kit commands; re-read after each | Pass |
| C-07 | blocked precedes ready-for-debrief; 013 itself stays unverified until 014 is | Pass |
| Inventory (C-09) | inventory never used as status evidence | Pass |

Post-design re-check: no change.

## Project Structure

```text
specs/010-close-gap-skill/
├── spec.md, plan.md, research.md, data-model.md, quickstart.md
├── contracts/close-gap-command.md
└── tasks.md                                            # /speckit-tasks

extensions/spec-status/
├── extension.yml                                       # 0.3.0, six commands
├── commands/speckit.spec-status.close-gap.md           # new
├── commands/speckit.spec-status.review.md              # follow-up
├── reference/roles.md                                  # verification-gap; pm row
└── README.md                                           # v0.3.0 heading; six skills; close-gap section

internal/extension/{helpers_test.go,roles_test.go}      # command list, review sections
testdata/skills/close-gap-all-causes/                   # new fixture + EXPECTED.md
scripts/check-extension-install.sh                      # six skills
.github/workflows/release.yml                           # n = 6; notes text
README.md                                               # six skills
specs/001-spec-status/tasks.md                          # implementation record
```

**Structure Decision**: extension files only, plus their tests and release plumbing.

## Phasing for tasks

1. Tests first: add close-gap to `extensionCommands`, `verification-gap` to `reviewSections`, the fixture with hand-written `EXPECTED.md` (fails until the command and roles exist).
2. Command file; roles reference; review follow-up; manifest 0.3.0.
3. Six-skill surfaces: install check, release workflow, both READMEs.
4. Install check run (Spec Kit 1.2.0), agent walkthrough on the fixture, implementation record, roadmap proposal (013 `implemented`; `verified` only after 014 is verified and the v0.3.0 checks).

## Risks

- **Behavioral correctness is agent-dependent**: the Go tests check structure and fixture signals, not an agent's reasoning; the walkthrough (as in feature 005) is the behavioral check and its result is recorded as executed by the implementer.
- **Release coupling**: the release workflow's version check requires extension version = tag; 0.3.0 must not be tagged before this lands.

## Complexity Tracking

None.
