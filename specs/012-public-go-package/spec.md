# Feature Specification: Public Go package

**Feature Branch**: `012-public-go-package`

**Created**: 2026-11-21

**Status**: Draft

**Input**: Roadmap item 018, "Public Go package", issue #65, chosen for v0.4.0 by the owner in decision D-2026-11-21-01 (2): "a public package at `github.com/tosin2013/spec-status/pkg/status`: `Load`, the schema 1.0 report types, and the exact `--json` bytes from the same code the CLI uses. The Go API is experimental until 1.0: report fields follow schema 1.0 and only gain additive fields, but Go signatures may change in a minor release, and release notes call out any change." Constitution Art. V (amended 2026-11-21) versions the Go API separately.

## Context

specfoundry is moving its gate to Go. Today the gate runs `spec status --json --root <dir>` as a subprocess and parses schema 1.0. It would rather call spec-status in-process: there would be no binary to install, and it would use the same logic. All of spec-status's logic is under `internal/`, which other modules cannot import.

The library must give the same answer as the binary. specfoundry plans a test that runs both and compares the decisions, and spec-status must make a divergence impossible by construction, not just unlikely.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - A Go gate evaluates a project in-process (Priority: P1)

A Go program imports the public package, points it at a project directory, and gets the report back as Go values: items, statuses, dependencies, task counts including phases, the summary and diagnostics. The field names and meanings are the ones documented for schema 1.0.

**Why this priority**: This is the request in issue #65.

**Independent Test**: A test in a separate package (as an outside module would be) loads `testdata/example-roadmap` and checks the item IDs, statuses, the summary counts and one item's phases against the documented example.

**Acceptance Scenarios**:

1. **Given** a project with a supported roadmap, **When** a program loads it through the package, **Then** it receives a report whose outcome is `report`, whose schema version is `"1.0"`, and whose items, summary and diagnostics match the CLI's JSON for the same project.
2. **Given** a project whose roadmap is missing or unrecognized, **When** it is loaded, **Then** the program receives a report with outcome `fatal` and the error code (for example `SS_ROADMAP_UNAVAILABLE`), exactly as the JSON would show. This is not reported as a Go error.
3. **Given** a root of `""` or an unknown fail-on value, **When** it is loaded, **Then** the call returns a Go error that names the invalid input, and no report.

---

### User Story 2 - The library and the binary cannot disagree (Priority: P1)

A maintainer of a consuming tool needs to know that switching from the subprocess to the library changes nothing. Serializing the library's report gives the exact bytes `spec status --json` prints, and the library's exit decision equals the binary's exit code.

**Why this priority**: The consumer's correctness depends on it. A silent drift would let a gate allow work the CLI would refuse.

**Independent Test**: For every fixture root under `testdata/`, compare the library's serialized report with the CLI's `--json` output byte for byte, and the library's exit decision with the CLI's exit code, for each fail-on threshold.

**Acceptance Scenarios**:

1. **Given** any fixture root and any threshold (`error`, `warning`, `info`, `none`), **When** the library's report is serialized, **Then** the bytes equal `spec status --json --root <root> --fail-on <threshold>` output, including the trailing newline.
2. **Given** the same inputs, **When** the library's exit decision is requested, **Then** it equals the binary's exit code (0, 1 or 2).
3. **Given** a future change to the JSON (for example a new additive field), **When** it is made, **Then** the library reflects it with no separate edit, because both use one implementation. The byte-identity test enforces this.

---

### User Story 3 - Consumers know what is promised (Priority: P2)

A developer evaluating the package reads its documentation and knows how to call it, which guarantees hold (the report fields follow schema 1.0, additively) and which do not yet (Go signatures may change before 1.0, called out in release notes).

**Why this priority**: Adoption without a clear promise leads to surprise breakage (Art. V).

**Independent Test**: The library documentation page and the package documentation exist. A compiled example runs against the example roadmap and prints the expected output.

**Acceptance Scenarios**:

1. **Given** the repository, **When** a developer opens `docs/library.md`, **Then** it shows installation (`go get`), a usage example, the experimental promise, the relation between the package version and `SchemaVersion`, and the safety properties (read-only, root-confined, no network).
2. **Given** the package's example, **When** the test suite runs, **Then** the example compiles, runs and matches its documented output.

### Edge Cases

- **Relative root:** a relative root is resolved against the process's working directory, exactly as `--root` is. The root is discovered the same way as with `--root`: the nearest ancestor with a non-symlink `.specify`.
- **Roadmap option:** a roadmap path behaves exactly like `--roadmap`. A path outside the root gives a fatal `SS_PATH_UNSAFE` inside the report, not a Go error.
- **Fail-on default:** an empty fail-on value means `error`, the CLI default. `none` never turns a fatal outcome into exit 0.
- **Concurrency:** concurrent loads of the same or different roots from several goroutines are safe and give the same results as sequential loads. The package has no mutable package-level state.
- **Side effects:** the package never writes files, opens network connections, starts processes, reads environment variables for behavior, or changes the process's working directory or environment.
- **Unreadable and unsafe inputs** (symlinks, sensitive names, oversize files, files that change during a read) give the same fatal codes or diagnostics as the CLI.
- **Untrusted text** in the report (titles, headings, messages) is returned as raw text, as in the JSON. Sanitizing terminal escapes applies only to human text, which this package does not produce.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The module MUST provide a public, importable package at `github.com/tosin2013/spec-status/pkg/status`.
- **FR-002**: The package MUST provide a load operation that takes a project root and options (roadmap path and fail-on threshold, with empty meaning the CLI defaults). It returns either a report or an error.
  - It returns an error only for invalid input: an empty root, or a fail-on value other than `error`, `warning`, `info` or `none`.
- **FR-003**: The report MUST expose every schema 1.0 field: schema version, outcome, complete, source, configuration, summary, items (including `tasks.phases`), diagnostics and error. Each MUST have the same name, meaning and null rules as `json-contract.md`.
  - Fields that are `null` in the JSON MUST be representable as absent.
- **FR-004**: The package MUST provide the schema version as a constant equal to `"1.0"`. The `schema_version` value in every report MUST equal it.
- **FR-005**: Serializing a report through the package MUST produce exactly the bytes that `spec status --json` prints for the same root and options.
  - The CLI's JSON output MUST be produced by the same code path, so that the two cannot diverge.
- **FR-006**: The package MUST provide an exit decision (0, 1 or 2) for a report that equals the binary's exit code for the same inputs.
- **FR-007**: The package MUST keep the core's runtime guarantees (Art. II, III):
  - no writes, network, subprocesses or environment-dependent behavior;
  - root confinement and the existing input limits;
  - only the standard library and this module's own internal packages as imports. An automated import check MUST enforce this.
- **FR-008**: Loads MUST be safe for concurrent use. A test MUST run concurrent loads under the race detector.
- **FR-009**: A test MUST check, for every fixture root under `testdata/` and every fail-on threshold, that the library's bytes equal the CLI's `--json` output and that its exit decision equals the CLI's exit code.
- **FR-010**: The package documentation MUST state the API promise: experimental until 1.0. Report fields follow schema 1.0 and change only additively within it. Go signatures may change in a minor release, and release notes call out any change.
  - `docs/library.md` MUST give installation, an example, that promise, and the safety properties.
  - The README MUST link to it.
  - `docs/architecture.md` MUST show the package.
- **FR-011**: A compiled example MUST load `testdata/example-roadmap` and print a short summary, checked by the test suite.
- **FR-012**: The CLI's text output, JSON bytes, exit codes and flags MUST be unchanged by this feature (byte comparison with the pre-change build on every fixture root).

### Key Entities

- **Report:** the schema 1.0 document as Go values. Its fields are the outcome, completeness, source, configuration, summary, items, diagnostics and fatal error.
- **Options:** the roadmap path and the fail-on threshold. The zero value means the CLI defaults.
- **Exit decision:** 0 (no finding reached the threshold), 1 (a finding did) or 2 (fatal).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For 100% of fixture roots and all four thresholds, the library's serialized report equals the CLI's `--json` bytes, and the exit decisions match.
- **SC-002**: A program outside this module can import the package and load a project with one call and no binary installed. This is shown after the release by a scratch module built against the tagged version.
- **SC-003**: The CLI's text, JSON and exit codes are byte-identical to the pre-change build on every fixture root.
- **SC-004**: Concurrent loads (at least 8 goroutines) pass under the race detector, with results identical to sequential loads.
- **SC-005**: The import check fails if the package (or anything it imports from this module) gains a dependency outside the standard library and this module.

## Assumptions

- Root discovery, path handling and every finding come from the existing core, unchanged. The package adds a public surface, not new behavior.
- The package version is the module version (v0.4.0 first). `SchemaVersion` changes only when the JSON schema changes, as today.
- `Options` has no `CWD` field: relative paths use the process working directory, as the CLI does. Callers that need another base pass an absolute root.
- Human text, filters and `--explain` are not part of the library. They are presentation, and they depend on the human-text styling library, which the package must not import.
- Out of scope: an inventory API (`spec inventory`), a stable 1.0 promise, streaming or incremental evaluation, and any write path.
