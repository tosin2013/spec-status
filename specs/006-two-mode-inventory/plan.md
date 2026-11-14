# Implementation Plan: Two-mode roadmap and feature inventory

**Branch**: `006-two-mode-specify` | **Date**: 2026-11-14 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/006-two-mode-inventory/spec.md` (roadmap item 008, decision D-2026-11-12-08, target v0.3.0).

## Summary

Two changes to the `spec` CLI:

1. **Roadmap recognition in `spec status`.** The parser also reports whether it saw a visible `## Planned Specs` heading. After a successful parse with zero items, that heading means "explicitly empty" (exit 0); its absence is the new fatal `SS_ROADMAP_UNRECOGNIZED` (exit 2). Reports with items are byte-identical. A probe build showed only one existing test (`TestEmptyRoadmapJSON`) depends on the old zero-item result (research R2). No JSON field is added (R3).
2. **`spec inventory`.** A new Cobra subcommand and an `internal/status` evaluator that lists the immediate children of `specs/` through a new bounded, non-following `safeio.ListDir`, classifies Spec Kit feature folders by name, stats `spec.md`/`plan.md`/`tasks.md`, counts tasks.md checkboxes with the existing `parseTasks`, and renders a separate inventory JSON document (`inventory_schema_version` 1.0) or Lip Gloss human text.

## Technical Context

**Language/Version**: Go 1.26 (go.mod `go 1.26.0`), `CGO_ENABLED=0`

**Primary Dependencies**: standard library for parsing, safeio and evaluation; Cobra/pflag in `internal/cli`; Lip Gloss/termenv/uniseg only in `internal/status/style.go`. No new module.

**Storage**: none (read-only; no caches)

**Testing**: `go test ./...`: unit tests in `internal/status`, `internal/safeio`, `internal/cli`; hand-counted fixtures under `testdata/two-mode/` and `testdata/conformance/`; the extension reference test in `internal/extension`.

**Target Platform**: linux/amd64 and darwin/arm64 (released and smoke-tested); other platforms build, and `ListDir` returns `unsupported` there like reads.

**Project Type**: CLI

**Performance Goals**: inventory of up to 1024 feature folders within the existing bounds; one directory listing plus at most three stats and one read per feature folder.

**Constraints**: constitution Art. I–III and V; offline; deterministic bytes; root-confined; sensitive names never opened; no recursive inventory; JSON schema 1.0 for `spec status` unchanged.

**Scale/Scope**: about 600–900 lines of Go plus tests; three new diagnostic codes; one new JSON document; docs in five contract files, README and the extension reference.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Article / constraint | Check | Result |
|---|---|---|
| I Truthful evidence (C-01) | Inventory reports file presence and checkbox counts only; says "checked boxes are not verification"; never infers status, priority, dependencies. Unrecognized ≠ empty. | Pass |
| II Deterministic, standalone (C-02) | One evaluation feeds text and JSON per mode; byte order; no timestamps; no network, subprocess, writes; roadmap extension not required. `ListDir` uses `syscall.Open`, not `os.OpenFile`. | Pass |
| II amendment record | Art. II says rollout is undecided; this feature decides it (D-2026-11-12-08). An amendment note is drafted and applied only with owner approval (R11). | Pass, owner action noted |
| III Narrow, safe input (C-03) | Exact `## Planned Specs` match; one listing of `specs/`, no recursion; symlinks never followed; sensitive names not stat'ed; existing limits plus `MaxDirEntries`. Unreliable tasks.md is fatal. | Pass |
| IV History (C-04) | Untouched: inventory has no supersession; roadmap mode unchanged for item-bearing files. | Pass |
| V Compatibility (C-05) | Exits 0/1/2, `--fail-on` default `error`, `none` never masks fatal; JSON mode one document including invocation errors; separate inventory schema version. Pre-1.0 behavior change listed in release notes. | Pass |
| VI No inferred quarters | No timing fields. | Pass |
| VII Distribution claims | No new platform claims. | Pass |
| VIII Mutation authority | No writes; skills changes limited to the output reference. | Pass |
| C-08 Upstream-pinned bar | Recognition heading taken from the pinned speckit-roadmap v0.1.0 template; feature-folder forms from spec-kit v1.2.0 numbering. | Pass |
| C-09 Two-mode direction | Missing, empty, unrecognized and malformed stay distinct; inventory is explicit, never a fallback; checklists excluded. | Pass |
| C-10 Line-level label rule | Recognition runs after the parse, so C-10 fatals are unchanged. | Pass |
| Dependency allowlist | No new import outside the allowed files. | Pass |

Post-design re-check (after contracts and data model): no change. No violations; Complexity Tracking is empty.

## Project Structure

### Documentation (this feature)

```text
specs/006-two-mode-inventory/
├── spec.md
├── plan.md                         # this file
├── research.md                     # R1–R12
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── roadmap-recognition.md      # spec status changes, new code, fixtures
│   ├── inventory-cli.md            # command, flags, exits, streams
│   ├── inventory-json.md           # inventory_schema_version 1.0
│   └── inventory-text.md           # human text layout
├── checklists/requirements.md
└── tasks.md                        # /speckit-tasks, not created here
```

### Source Code (repository root)

```text
internal/safeio/
├── fs.go               # + MaxDirEntries, DirEntry, (*FS).ListDir
├── open_unix.go        # + openDirNoFollow (O_DIRECTORY|O_NOFOLLOW)
├── open_other.go       # + openDirNoFollow → unsupported
└── fs_test.go          # + ListDir tests (symlink, FIFO, sensitive, limit, change, order)

internal/status/
├── parse.go            # parseRoadmap also returns plannedSpecs
├── run.go              # recognition rule; missing-roadmap message; shared root/.specify setup
├── inventory.go        # new: EvaluateInventory, folder classification, file states
├── inventory_json.go   # new: InventoryJSON view
├── style.go            # + inventory human text; explicit-empty panel line
├── render.go           # fatal text reused for inventory
├── status_test.go      # TestEmptyRoadmapJSON replaced by recognition tests
├── recognition_test.go # new
└── inventory_test.go   # new

internal/cli/
├── cli.go              # + newInventory; rejected flags; inventory-shaped invocation JSON
└── cli_test.go         # + inventory invocation, exit codes, JSON-on-error

testdata/two-mode/      # new hand-counted fixtures (recognition and inventory)
testdata/conformance/   # unchanged (re-run only)

extensions/spec-status/reference/output.md      # + 3 codes
specs/001-spec-status/{diagnostics.md,json-contract.md,cli-contract.md,roadmap-format.md}
specs/001-spec-status/inventory-json.md          # copy of the inventory contract
README.md                                        # Limits, missing-roadmap row, inventory section
```

**Structure Decision**: single Go module, existing packages. Inventory lives in `internal/status` to share `parseTasks` and the diagnostic model (R6); the only filesystem addition is in `internal/safeio` (R5).

## Phasing for tasks

1. **Recognition (US1, P1):** fixtures and tests first, then parser flag, rule, messages, explicit-empty line, docs. Independently shippable.
2. **safeio.ListDir:** tests first, then implementation. Blocks inventory.
3. **Inventory (US2, P1):** fixtures and expected files by hand, then evaluator, JSON, text, CLI.
4. **Safety and determinism (US3, P2):** temp-dir tests for symlinks, FIFO, sensitive names, limits, repeat runs.
5. **Docs and release:** contracts into `specs/001-spec-status/`, README, extension reference, tasks.md implementation record, roadmap status, release-note text, constitution amendment draft for owner approval.

## Risks

- **Over-inclusive sensitive matcher** skips some legitimately named feature folders with an error (R8). Same as roadmap mode; listed in README Limits.
- **Message change for a missing roadmap** changes JSON `message` bytes for that fatal; agents matching on message text (not code) would notice. Release notes call it out.
- **Darwin residual limit:** like reads, `ListDir` on Darwin cannot confirm the opened directory's path (`confirmOpened` is Linux-only); the before/after identity check and `O_NOFOLLOW` remain.

## Complexity Tracking

None.
