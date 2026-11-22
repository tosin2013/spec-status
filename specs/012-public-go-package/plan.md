# Implementation Plan: Public Go package

**Branch**: `012-public-go-package` | **Date**: 2026-11-21 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/012-public-go-package/spec.md` (roadmap item 018, issue #65, decision D-2026-11-21-01 (2), target v0.4.0).

## Summary

The work introduces one set of schema 1.0 report types, used by both the CLI and the new public package.

- **`internal/schema` (new):** a leaf package, standard library only. It defines the exported report types (`Report`, `Item`, `Tasks`, `Phase`, …) with JSON tags and precise Go types: pointers where the JSON allows `null`, no `any`. It moves `Percent` here and gives it accessors and `UnmarshalJSON`. It also defines `Marshal` (the exact encoder settings and trailing newline that `--json` uses) and `ExitCode`.
- **`internal/status`:** builds a `schema.Report` from its internal `Report` (`toSchema`, which replaces the private `docView` structs). `status.JSON` becomes `schema.Marshal(toSchema(rep))`, and `status.ExitCode` delegates to `schema.ExitCode`. The CLI is unchanged.
- **`pkg/status` (new):** re-exports the types as aliases (`type Report = schema.Report`, …). It adds `SchemaVersion`, `Options`, `Load(root, opts)` (validates, then calls `status.Evaluate` and `toSchema` through an exported bridge `status.Schema(rep)`), `Marshal` and `ExitCode`. The package depends only on this module and the standard library.

Byte identity holds by construction, because both paths call `toSchema` and `schema.Marshal`. It is also checked by a test over every fixture root and all four thresholds, plus a JSON round trip (decode the CLI bytes into `pkg/status.Report`, encode again, compare bytes).

## Technical Context

**Language/Version**: Go 1.26, `CGO_ENABLED=0`

**Primary Dependencies**: standard library only for `internal/schema` and `pkg/status`

**Storage**: none

**Testing**:
- `go test ./...` and `go test -race ./pkg/status`;
- `pkg/status` external tests (`package status_test`): the byte identity with the CLI over every fixture root × 4 thresholds, run through `internal/cli.Run`; the round trip; concurrency (8+ goroutines); invalid input;
- a compiled `Example`;
- the imports test extended to `internal/schema` and `pkg/status`;
- a byte comparison of the CLI against the pre-change build on every fixture root.

**Target Platform**: linux/amd64, darwin/arm64 (any platform the module builds on, for library users)

**Project Type**: CLI plus library

**Performance Goals**: the same work as `spec status --json`. The conversion is linear in report size.

**Constraints**:
- Art. II: no side effects; standard library only in the new packages.
- Art. V: schema 1.0 unchanged; the Go API is versioned separately and experimental (D-2026-11-21-01).
- Art. I: no new inference.
- FR-012: the CLI's bytes are unchanged.

**Scale/Scope**: about 250 lines of types and conversion (mostly moved out of `render.go`), about 80 lines of public API, tests, and docs.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Article | Check | Result |
|---|---|---|
| I Truthful evidence | The library returns the same report as the JSON. Nothing is added or inferred. | Pass |
| II Standalone, no side effects | New packages use only the standard library and this module. No writes, network, exec or environment reads. The imports test covers the new paths. | Pass |
| III Narrow, safe input | Same core, same root confinement, limits and fatal codes | Pass |
| IV History | Unchanged | Pass |
| V Compatibility | Schema 1.0 bytes unchanged (byte tests). The Go API is a separately versioned surface, experimental until 1.0 (amended Art. V). | Pass |
| VI No quarters | Not affected | Pass |
| VII Evidence-gated distribution | No support claim beyond what the post-release scratch-module check shows | Pass |
| VIII Approval | Spec, plan, tasks and implementation each go through owner-merged PRs | Pass |

Post-design re-check: unchanged, all pass.

## Project Structure

### Documentation (this feature)

```text
specs/012-public-go-package/
├── spec.md
├── plan.md              # this file
├── research.md          # R1–R7
├── data-model.md        # schema types and Go types per JSON field
├── quickstart.md
├── contracts/
│   └── go-api.md        # public API contract
└── checklists/requirements.md
```

### Source Code (repository root)

```text
internal/schema/          # new, stdlib only
├── schema.go             # Report … Phase types, Percent, Location, Ref
├── marshal.go            # Marshal (SetEscapeHTML(false), trailing newline, serialization fallback)
└── exit.go               # ExitCode
internal/status/
├── render.go             # toSchema replaces docView/itemView…; JSON = schema.Marshal(toSchema(rep))
├── model.go              # type Percent = schema.Percent; percentOf → schema.PercentOf
├── run.go                # ExitCode delegates to schema.ExitCode(toSchema(rep))
└── imports_test.go       # stdlib-only rule for internal/schema and pkg/status
pkg/status/               # new public package
├── doc.go                # package doc, API promise
├── status.go             # aliases, SchemaVersion, Options, Load, Marshal, ExitCode
├── status_test.go        # byte identity, round trip, invalid input, concurrency
└── example_test.go       # Example on testdata/example-roadmap
docs/library.md; README.md link; docs/architecture.md; CLAUDE.md
specs/001-spec-status/{json-contract.md,tasks.md}
```

**Structure Decision**: The layout is a leaf `internal/schema` package with type aliases in `pkg/status` (research R1). It has no import cycle: `pkg/status` imports `internal/status`, which imports `internal/schema`.

## Complexity Tracking

| Addition | Why needed | Simpler alternative rejected because |
|---|---|---|
| New `internal/schema` package | Breaks the cycle between `pkg/status` (needs `Evaluate`) and the shared types (needed by `internal/status`) | Types in `pkg/status` would force `internal/status` to import a package that imports it. Moving `Evaluate` into `pkg/status` would put the whole core in the public API. |
