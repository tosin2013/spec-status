# Implementation Plan: Upstream-format conformance

**Branch**: `002-upstream-format-conformance` | **Date**: 2026-11-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/002-upstream-format-conformance/spec.md`

## Summary

Make pinned upstream versions the definition of supported input: github/spec-kit at commit `4ca46a07` (tag v1.2.0) and srobroek/speckit-roadmap v0.1.0. Copy their templates with licenses into `testdata/upstream/`, publish a shape inventory of what they generate, and add hand-counted conformance fixtures. Add a development-only drift check that reports shape changes between a pinned and a candidate template. Record a new acceptance bar (D-2026-11-07-09) that replaces D-2026-11-07-06, then re-classify past review findings under it. The product binary changes only if a conformance fixture exposes a misread of an upstream shape. A first probe found none: both pinned templates already parse, giving 2 items and 34 of 34 task lines.

## Technical Context

**Language/Version**: Go 1.22 (module minimum); tested here with Go 1.24.3 on darwin/arm64 and Go 1.22.2 in Docker on linux/arm64 and linux/amd64

**Primary Dependencies**: Standard library only for all new code. New code lives in `_test.go` files, so it never reaches the binary. The existing Cobra and Lip Gloss exceptions are untouched.

**Storage**: Files only: pinned templates, fixtures, and hand-written expected results under `testdata/`

**Testing**: `go test`. The conformance runner is a table-driven test over fixture directories; the drift check is a test skipped unless `-candidate` is passed.

**Target Platform**: Development machines where `go test` runs. No runtime change.

**Project Type**: CLI (feature 001) plus test-only tooling

**Performance Goals**: The conformance suite adds under 1 second to `go test ./...` on a development machine.

**Constraints**: Constitution Article II: no network, subprocess, or writes at runtime; the drift check never fetches upstream content and works on files a maintainer has downloaded. Upstream content keeps its license and attribution.

**Scale/Scope**: Two upstream projects, one pin each. About 6–10 conformance fixtures, each holding one roadmap and up to a few tasks files.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Article | Requirement | This plan | Status |
|---|---|---|---|
| I. Truthful evidence | No claims beyond measurable behavior; state limits | Support is claimed only for pinned versions that have fixtures (FR-016). Expected results are hand counts, not product output (research R6). Exotic-input limits stay listed in README. | Pass |
| II. Deterministic, standalone | Stdlib parsing; no runtime network, subprocess, or writes; dev tooling is not a runtime dependency | All new code is in test files. The drift check reads local files only. No new module dependency. | Pass |
| III. Narrow, safe input | Pinned dialect; unsupported item-like input fails visibly | The dialect becomes *more* precisely pinned. Existing visible failures are kept (FR-014, R9). | Pass |
| IV. History cannot become completion | Superseded items keep history | Not affected. D-2026-11-07-06 is kept and marked superseded rather than rewritten (R8). | Pass |
| V. Compatibility and findings | Version the input dialect, schema, and CLI separately | Pins act as the input-dialect version. No CLI or JSON schema change. | Pass |
| VI. No inferred quarters | — | Not affected. | Pass |
| VII. Distribution claims need evidence | No support claims without evidence | README names the supported upstream versions only. No platform or release claims are added. | Pass |
| VIII. Future mutation authority | No writes or agent launches | The drift check and tests never write into the repository or launch upstream commands. Upstream-generated fixtures are stored once from runs a maintainer performed. | Pass |
| Governance §2–3 | Implementation needs separate authorization | This plan and the later `tasks.md` do not authorize implementation. **Owner approval is required before `/speckit-implement`.** | Gate (owner) |

**Post-design re-check (after Phase 1)**: Still passes. The design adds test-only Go files, fixtures, and documentation, and changes product code only if a conformance fixture fails. No violations, so Complexity Tracking is empty.

## Project Structure

### Documentation (this feature)

```text
specs/002-upstream-format-conformance/
├── plan.md                     # This file
├── research.md                 # Phase 0: pins, licenses, drift-check form, oracle
├── data-model.md               # Phase 1: pin, shape, fixture, expected result, drift report, bar decision
├── quickstart.md               # Phase 1: how to run and validate
├── contracts/
│   ├── shape-inventory.md      # Supported shapes per pin, and how spec status reads each
│   ├── expected-result.md      # Format of a fixture's expected.json
│   └── drift-report.md         # What the drift check reports and how the pin moves
├── checklists/requirements.md  # From /speckit-specify
└── tasks.md                    # Phase 2: /speckit-tasks (not created here)
```

### Source Code (repository root)

```text
testdata/
├── upstream/
│   ├── spec-kit/v1.2.0/                  # pinned templates (tasks, spec, plan), LICENSE (MIT), NOTICE
│   └── speckit-roadmap/v0.1.0/           # pinned roadmap template, LICENSE (Apache-2.0), NOTICE
├── conformance/
│   ├── template-filled/                  # roadmap and tasks from the pinned templates, every construct filled
│   ├── hand-edits/                       # wrapped values, notes, nested subtasks, details paragraphs, fenced examples
│   ├── upstream-generated-002-tasks/     # tasks.md produced by /speckit-tasks for this feature, unedited
│   └── …                                 # each fixture directory has .specify/memory/roadmap.md, specs/…, expected.json
├── speckit-roadmap-upstream/             # existing upstream-generated fixture; gains an expected.json
└── example-roadmap/                      # existing; unchanged

internal/conformance/
├── conformance_test.go                   # runs status.Evaluate on each fixture and compares it with expected.json
├── shapes_test.go                        # independent line classifier producing shape keys
└── drift_test.go                         # TestDrift (needs -candidate) and TestPinnedInventory

specs/001-spec-status/decisions.md        # adds D-2026-11-07-09; marks D-2026-11-07-06 superseded
specs/001-spec-status/roadmap-format.md   # separates the upstream shape inventory from the strictness rules
README.md                                 # supported upstream versions; Limits re-checked
.specify/memory/roadmap.md                # adds item 007 for this feature
```

**Structure Decision**: Reuse the existing single Go module. All new Go code goes into a test-only package, `internal/conformance`, so it can't affect the binary. It imports `internal/status` to evaluate fixtures, the same way the existing tests do. Pinned upstream material and fixtures sit beside the existing `testdata/` fixtures, and documentation changes stay in feature 001's documents, which own the dialect and decisions.

## Complexity Tracking

No constitution violations to justify.
