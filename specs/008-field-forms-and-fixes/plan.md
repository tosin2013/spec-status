# Implementation Plan: Real-world field forms and small defect fixes

**Branch**: `008-fields-fixes-specify` | **Date**: 2026-11-14 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/008-field-forms-and-fixes/spec.md` (roadmap items 015 and 016, decision D-2026-11-14-01 (3), (4), target v0.3.0).

## Summary

Eight small changes, each measured first (spec Context):

- **Parser** (`internal/status/parse.go`): two more `fieldLabels` entries for `Superseded by` / `superseded by` (R1); a placeholder check for `Spec dir` and `Depends on` values that leaves the field unset and records a finding, except when the enclosed text is a valid ID list (R2); deduplication of repeated `Depends on` IDs with a record of the repeats (R4).
- **Analysis** (`analyze.go`): emit `SS_FIELD_PLACEHOLDER` (info) and `SS_DEPENDENCY_DUPLICATE` (warning); record display-only shared-tasks attribution (R3, R5).
- **Human text** (`style.go`): a `statuses` panel line, `invalid supersession N`, an `(invalid)` row mark, and a shared-tasks row note (R5).
- **CLI** (`internal/cli/cli.go`): `--json=false` honored, empty `--root`/`--roadmap` refused, requested `fail_on` in status invocation JSON (R6–R8).

JSON stays schema 1.0 and byte-identical for every existing fixture; human text gains the `statuses` line everywhere, so the README example is regenerated.

## Technical Context

**Language/Version**: Go 1.26, `CGO_ENABLED=0`

**Primary Dependencies**: standard library (parser, analysis); Cobra/pflag (CLI); Lip Gloss only in `style.go`

**Storage**: none

**Testing**: `go test ./...`; new fixtures `testdata/field-forms/` with hand-written `expected-status.json`; unit tests for display and CLI; JSON byte comparison against the pre-change build

**Target Platform**: linux/amd64, darwin/arm64

**Project Type**: CLI

**Performance Goals**: no measurable change (constant work per field)

**Constraints**: C-03 narrow input, C-08 no silently misread metric field, C-10 line-level label rule, C-05 compatibility, Art. II plain surfaces

**Scale/Scope**: roughly 150 lines of Go plus tests; two registry rows; README and contract edits

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Article / constraint | Check | Result |
|---|---|---|
| I Truthful evidence (C-01) | Placeholders are reported, not hidden; display additions show existing facts only | Pass |
| II Deterministic, standalone (C-02) | Status keys sorted; fixed messages; no I/O added; text and JSON still one report | Pass |
| III Narrow input (C-03) | Two exact spellings added; other variants stay visible near misses; placeholder rule needs a whole-value enclosure | Pass |
| IV History (C-04) | Invalid supersession now marked in text; no redirect | Pass |
| V Compatibility (C-05) | Schema 1.0; additive codes; default-threshold exits unchanged for existing fixtures; `dependencies` dedup only affects inputs with repeats; README example update noted | Pass |
| C-08 metric fields | A bracketed or emphasized valid ID list stays `SS_DEPENDENCY_MALFORMED`; a placeholder never satisfies a stage or dependency | Pass |
| C-10 label rule | Accepted spellings become field lines; prose with the label still fatal | Pass |
| VI–VIII | No timing, distribution or mutation change | Pass |
| Dependency allowlist | No new import | Pass |

Post-design re-check: no change; no violations.

## Project Structure

### Documentation (this feature)

```text
specs/008-field-forms-and-fixes/
├── spec.md
├── plan.md              # this file
├── research.md          # R1–R10
├── data-model.md
├── quickstart.md
├── contracts/
│   └── changes.md       # grammar, registry, JSON, human text, CLI, fixtures
├── checklists/requirements.md
└── tasks.md             # /speckit-tasks
```

### Source Code (repository root)

```text
internal/status/
├── parse.go                 # fieldLabels +2; placeholderValue; dedup + DupDeps; rawItem.Placeholders
├── analyze.go               # emit 2 codes; TasksSharedWith
├── model.go                 # Item.TasksSharedWith (display only)
├── style.go                 # statuses line, invalid supersession line, row marks
├── field_forms_test.go      # new: rows A–D, raw template, display (row E)
└── (existing tests updated where they compare full human text)

internal/cli/
├── cli.go                   # jsonRequested, empty root/roadmap, fail_on in invocationJSON
└── cli_test.go              # rows F–H

testdata/field-forms/        # new fixtures with hand-written expected-status.json

specs/001-spec-status/{roadmap-format.md,diagnostics.md,json-contract.md,cli-contract.md}
extensions/spec-status/reference/output.md
README.md                    # Limits, findings rows, example panel
specs/001-spec-status/tasks.md   # implementation record
```

**Structure Decision**: existing packages only.

## Phasing for tasks

1. Baseline outputs (JSON and text) for all fixture roots from the pre-change build.
2. US1 + US2 (P1): fixtures and tests first, then labels and placeholders.
3. US3 (P2): duplicate IDs.
4. US4 (P2): display additions; update full-text test expectations and the README example by hand.
5. US5 (P3): CLI edge cases.
6. Docs, byte comparison, platforms, Doomer review, roadmap proposals for 015 and 016.

## Risks

- **Full-text tests and README**: the always-on `statuses` line changes every human report; tests that pin whole panels need updating by hand, and the README example with them.
- **Placeholder false positives**: a real path written as `_specs/x_` becomes a placeholder; it is reported (info) and, for specced items, the stage rule still errors, so it is visible.

## Complexity Tracking

None.
