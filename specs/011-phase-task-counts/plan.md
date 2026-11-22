# Implementation Plan: Per-phase task counts

**Branch**: `011-phase-task-counts` | **Date**: 2026-11-21 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/011-phase-task-counts/spec.md` (roadmap item 017, issue #66, decision D-2026-11-21-01 (1), target v0.4.0).

## Summary

Each item's `tasks` object in `spec status --json` gets an additive `phases` array.

- **Parsing:** `parseTasks` (`internal/status/parse.go`) already walks every line. It already recognizes ATX headings at the point where fences, comments, raw HTML blocks, blank lines and block quotes have been handled (around line 714, the `heading(line)` branch). At that point it will also record qualifying `## Phase …` headings (R1–R3) and close the open phase on any level-1 or level-2 heading.
- **Counting:** each recorded checkbox is counted in the phase open at that moment (R4).
- **Model:** `taskFile` and `ItemTasks` carry the phases. `itemToView` renders them as `[]` for a read file and `null` for any other state (R5).
- **Unchanged:** human text, the summary, diagnostics and exit codes.
- **Conformance:** the comparison in `internal/conformance` learns an optional `phases` expectation. The real Spec Kit-generated fixture `upstream-generated-002-tasks` and the other fixtures with read tasks files get hand-counted phases (R6).

## Technical Context

**Language/Version**: Go 1.26, `CGO_ENABLED=0`

**Primary Dependencies**: standard library only (parse, model, render)

**Storage**: none

**Testing**: `go test ./...`. New unit tests in `internal/status/phases_test.go`, a synthetic fixture `testdata/phases/`, the conformance comparison extended for `phases`, `TestReadmeSamples` re-run, and a byte comparison of every fixture root against the pre-change build (text identical; JSON identical after deleting `phases`).

**Target Platform**: linux/amd64, darwin/arm64

**Project Type**: CLI

**Performance Goals**: no new reads. Constant work per line and per checkbox. The phase count is bounded by the line count, which the existing size limits already bound.

**Constraints**:
- Art. I: counts only, no completion inference.
- Art. II: no side effects; standard library.
- Art. III: no new fatal errors, and the same hidden-content rules as checkboxes.
- Art. V: schema 1.0, additive field, text unchanged.
- C-08: setext headings are a documented limit.

**Scale/Scope**: about 60 lines of Go plus tests, one new fixture, hand edits to five `expected.json` files, and contract/grammar/limits docs.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Article | Check | Result |
|---|---|---|
| I Truthful evidence | Phases report checked/total under a heading. No "done" flag, no status inference. `null` for unread files keeps "unknown" separate from "none". | Pass |
| II Standalone, no side effects | No new reads, imports or dependencies; standard library only | Pass |
| III Narrow, safe input | Reuses the existing heading recognition after fence, comment, HTML and quote handling. No new fatal errors; setext is a listed limit. | Pass |
| IV History | Phases are per item and do not touch supersession or active credit | Pass |
| V Compatibility | Additive field under schema 1.0; text, exit codes and other bytes unchanged (SC-002) | Pass |
| VI No quarters | Not affected | Pass |
| VII Evidence-gated distribution | No release claim. Shipping in v0.4.0 follows the release checks. | Pass |
| VIII Approval | Spec, plan, tasks and implementation each go through owner-merged PRs | Pass |

Post-design re-check: unchanged, all pass.

## Project Structure

### Documentation (this feature)

```text
specs/011-phase-task-counts/
├── spec.md
├── plan.md              # this file
├── research.md          # R1–R6
├── data-model.md        # Phase, changes to taskFile / ItemTasks / itemTaskView
├── quickstart.md        # validation commands
├── contracts/
│   └── phases-json.md   # JSON field contract and grammar
└── checklists/requirements.md
```

### Source Code (repository root)

```text
internal/status/
├── parse.go          # taskFile.Phases; phase recognition in parseTasks
├── model.go          # ItemTasks.Phases, Phase type
├── analyze.go        # copy phases into ItemTasks (next to the Checks copy, ~line 604)
├── render.go         # itemTaskView.Phases; itemToView
└── phases_test.go    # new
internal/conformance/expected_test.go   # optional phases expectation
testdata/phases/                         # new synthetic fixture (roadmap + tasks files)
testdata/conformance/*/expected.json     # hand-added phases
specs/001-spec-status/{json-contract.md,roadmap-format.md,tasks.md}
docs/limits.md                           # setext headings do not delimit phases
```

**Structure Decision**: Single Go module, existing packages. No new package.

## Complexity Tracking

None. No constitution violations.
