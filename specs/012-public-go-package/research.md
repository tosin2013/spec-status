# Research: Public Go package

## R1 — Package layout without an import cycle

- **Decision:** A leaf package `internal/schema` defines the exported report types, `Marshal` and `ExitCode`. `internal/status` imports it and converts its internal `Report` with `toSchema`. `pkg/status` re-exports the types with type aliases and implements `Load` on top of `internal/status`.
- **Rationale:** `pkg/status.Load` needs `status.Evaluate`, and `status.JSON` needs the types, so the types must sit below both packages. Aliases give outside callers the full field set: pkg.go.dev shows the alias and links to the documented fields. The CLI keeps calling `status.JSON`, now a two-line wrapper, so it needs no change.
- **Alternatives considered:**
  - Define the types in `pkg/status` and move JSON rendering there: the CLI would then depend on the public package, and every internal test that calls `status.JSON` would need a new import path, for no gain.
  - A second public package `pkg/status/schema`: it doubles the public surface to maintain under the experimental promise.
  - Generating the types from a schema file: there is no machine-readable schema today, and it would add a build step.

## R2 — Go types for nullable JSON fields

- **Decision:**
  - JSON fields that may be `null` become pointers: `*string` for strings, `*int` for counts and story, `*bool` for `supersession.valid`, and `*Location` for diagnostic and error sources.
  - Slices that are always arrays (`items`, `diagnostics`, `dependencies`, `related`) are non-nil after `Load`.
  - `tasks.phases` is a `[]Phase` whose nil value means `null` (unread file) and whose empty non-nil value means `[]`.
  - `summary` is `*Summary` (null on fatal).
  - `Location` is `{Path *string; Line *int}`, because the current `srcAny` emits a null path or a null line independently.
- **Rationale:** This reproduces today's bytes exactly with typed fields instead of `any`, which gives callers compile-time access. The round-trip test (R5) proves that decode-then-encode is lossless.
- **Alternatives considered:**
  - Keeping `any`: callers would need type assertions, and the bytes give no compile-time guarantee.
  - A separate `PhasesKnown bool`: it would diverge from the JSON shape.

## R3 — Percent

- **Decision:** Move `Percent` to `internal/schema`, and keep `internal/status` working through `type Percent = schema.Percent` and `schema.PercentOf`. Add `Valid() bool`, `Hundredths() int` and `UnmarshalJSON`, and keep `Text()` and `MarshalJSON` unchanged.
- **Rationale:** The value is emitted as a bare JSON number with at most two decimals, or `null`. Callers need the exact value without float rounding, so hundredths are exposed as an integer.
- **Alternatives considered:** `*float64`, rejected because a float cannot round-trip `41.67` to identical bytes in every case, and it would change the CLI's formatting path.

## R4 — Load input validation and options

- **Decision:**
  - `Options{Roadmap, FailOn string}`.
  - `Load` returns an `error` when `root == ""` or when `FailOn` is not in {"", error, warning, info, none}. Otherwise it calls `status.Evaluate(status.Option{Root: root, Roadmap: opts.Roadmap, FailOn: failOn})` with an empty `CWD`, so relative paths use the process working directory, as the CLI does.
  - Fatal problems stay in the report.
- **Rationale:** This mirrors the CLI's split. Invocation errors exit 2 before any evaluation; everything else is a report.
- **Alternatives considered:**
  - Returning an error for fatal reports: the consumer would lose the `SS_*` code and the identical decision.
  - Accepting a `CWD` option: it adds surface without a requester (spec assumption).

## R5 — Proving identity with the CLI

- **Decision:** `pkg/status/status_test.go` (external package) walks every fixture root under `testdata/` (each directory that directly contains `.specify/`) and runs four thresholds for each:
  - it runs the CLI in-process through `internal/cli.Run` with `status --json --root R --fail-on T`, capturing stdout and the exit code;
  - it calls `Load` + `Marshal` and compares the bytes, then compares `ExitCode(report)` with the CLI code;
  - it decodes the CLI bytes into `pkg/status.Report` with `encoding/json`, re-encodes them with `Marshal` and compares (round trip).
- **Rationale:** The comparison goes through the real command path, including flag handling, so a future CLI-only change to JSON would fail the test.
- **Alternatives considered:** Running a built binary as a subprocess. That is slower and needs `os/exec` in tests (allowed in tests, but unnecessary).

## R6 — ExitCode as one implementation

- **Decision:** `schema.ExitCode(r Report)`: 2 when `Outcome == "fatal"` or `!Complete`; 0 when `Configuration.FailOn == "none"`; otherwise 1 if any diagnostic's severity rank is at least the threshold rank, else 0. `status.ExitCode(rep)` returns `schema.ExitCode(toSchema(rep))`.
- **Rationale:** One rule set for the CLI and the library. The existing exit-code tests keep covering it through `status.ExitCode`.
- **Alternatives considered:** Two copies with a parity test. Two copies can drift between test runs, and the cost of delegating is negligible.

## R7 — Guardrails and docs

- **Decision:**
  - Extend `imports_test.go`: files under `internal/schema/` and `pkg/status/` may import only the standard library and `github.com/tosin2013/spec-status/internal/...`. Any dotted import path outside the module fails, with no exceptions, unlike the CLI and style surfaces.
  - `pkg/status/doc.go` states the experimental promise.
  - `docs/library.md` covers `go get`, an example, the promise, the safety properties and `SchemaVersion`.
  - Update the README (one line and a link), `docs/architecture.md` (container and component tables, a pipeline note) and CLAUDE.md (pipeline).
- **Rationale:** FR-007 and FR-010.
