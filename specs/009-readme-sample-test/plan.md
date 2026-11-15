# Implementation Plan: README sample test

**Branch**: `011-013-specify` | **Date**: 2026-11-15 | **Spec**: [spec.md](spec.md)

**Input**: `specs/009-readme-sample-test/spec.md` (roadmap item 011, D-2026-11-14-01 (5), v0.3.0).

## Summary

Label each README output block with an HTML comment naming its `spec` command; one test in `internal/cli` finds every label, runs the command through `Run` from the repository root (color off, 80 columns, development heading) and compares bytes; output-looking blocks without labels fail. Fix the stale "Filters" sample and retire the feature-006 inventory-only test.

## Technical Context

**Language/Version**: Go 1.26, `CGO_ENABLED=0`
**Primary Dependencies**: standard library (test only)
**Storage**: none
**Testing**: `go test ./internal/cli`
**Target Platform**: linux/amd64, darwin/arm64 (CI matrix)
**Project Type**: CLI (documentation guard)
**Performance Goals**: three in-process runs, well under a second
**Constraints**: no runtime change; test reads only `README.md` and `testdata/`
**Scale/Scope**: one test file (~150 lines), three README labels, one corrected block, one removed test

## Constitution Check

| Article | Check | Result |
|---|---|---|
| I Truthful evidence | README claims about output become tested facts | Pass |
| II Standalone | test only; no subprocess, network or writes; runtime untouched | Pass |
| III–VIII | no input, schema, distribution or mutation change | Pass |
| Dependency allowlist | standard library in a test file | Pass |

Post-design re-check: no change.

## Project Structure

```text
specs/009-readme-sample-test/
├── spec.md, plan.md, research.md, data-model.md, quickstart.md
├── contracts/readme-samples.md
└── tasks.md                     # /speckit-tasks

internal/cli/readme_samples_test.go   # new: TestReadmeSamples
internal/status/inventory_test.go     # remove TestReadmeInventorySample
README.md                             # 3 labels; corrected Filters block
CLAUDE.md                             # point the README-sample note at the test
```

**Structure Decision**: test lives with the CLI because it exercises real flag parsing (research R1).

## Phasing for tasks

1. Write `TestReadmeSamples` (fails: no labels yet, and the Filters block is stale).
2. Add the three labels; correct the Filters block from a real run; test passes.
3. Remove the old inventory test; update CLAUDE.md; drift and coverage checks from quickstart.

## Complexity Tracking

None.
