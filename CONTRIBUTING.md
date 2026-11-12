# Contributing to spec-status

Thanks for your interest. `spec status` is a small, read-only Go CLI that reports
on a Spec Kit roadmap. This page explains how to propose a change and which rules
a change must keep.

By participating you agree to the [Code of Conduct](CODE_OF_CONDUCT.md). Report
security issues privately as described in [SECURITY.md](SECURITY.md), not in a
public issue.

## Before you start

- **Open an issue first** for anything beyond a small fix, so the direction can
  be agreed before you write code.
- Behavior is defined by the documents in `specs/001-spec-status/` (spec, CLI and
  JSON contracts, roadmap grammar, diagnostics) and by recorded owner decisions in
  `specs/001-spec-status/decisions.md`. A change that alters a contract, a
  diagnostic, an exit code, the JSON schema, or a limit needs a recorded owner
  decision. Say so in the issue.
- Supported input is pinned to upstream templates (see
  `specs/002-upstream-format-conformance/`). Support for a new input shape is
  added deliberately, with a fixture, not by loosening the parser.

## Build and test

Building needs Go 1.26 or newer (D-2026-11-11-05). Always build with `CGO_ENABLED=0`.

```sh
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build -o spec ./cmd/spec
./spec status --color=never --root testdata/example-roadmap
CGO_ENABLED=0 go test ./internal/conformance   # hand-counted conformance fixtures
```

If output for `testdata/example-roadmap` changes, update the README sample.

## Rules a change must keep

These come from `.specify/memory/constitution.md` and are enforced by tests.

- **No runtime side effects.** No network, subprocesses, Git, file writes,
  caches, or telemetry. `internal/status/imports_test.go` fails otherwise.
- **Dependency allowlist.** Parsing, file access, and analysis use the Go
  standard library only. Cobra and pflag are allowed for the command layer and
  Lip Gloss for human text styling; adding any other dependency needs owner
  approval.
- **All file reads go through `internal/safeio`**, which keeps reads inside the
  project root and refuses symlinks and sensitive files.
- **Truthful evidence.** Checked tasks, recorded `implemented`, and recorded
  `verified` stay separate. Never infer verification or a date.
- **Visible failure.** Input the parser does not support fails with a diagnostic;
  it is never guessed at.
- **Determinism.** Same input, same bytes: no timestamps, stable ordering.
- **Text and JSON are two views of one report.** Human-text filters must not
  change the summary, diagnostics, exit code, or JSON.

## Test fixtures

- Every `expected.json` under `testdata/conformance/` is a **hand count**. Never
  regenerate one from `spec` output; count the input yourself and say so in its
  `counted_by` field.
- Add a test that fails before your change and passes after it.

## Dependency updates

Dependabot proposes Go module and GitHub Actions updates: security fixes as soon as an advisory is published, and one grouped minor/patch proposal per ecosystem each week. Major updates arrive as separate proposals. Every proposal runs CI.

The merge rule (decision D-2026-11-11-04) compares a proposal's `go.mod` with `main`:

| Proposal | Outcome |
|---|---|
| Only versions of modules already in `go.mod` change, the `go` line is unchanged, and CI passes | **Merge after review:** the owner approves the pull request |
| Any module path is added, direct or `// indirect` | **Decision needed:** record a decision naming each new module in `specs/001-spec-status/decisions.md` before merging |
| The `go` line is raised | **Decision needed** |
| CI fails | **Do not merge** |

Removing a module needs no decision. A GitHub Actions update is merge after review once CI passes, because actions are not part of the binary. The `dependency-review` check on each pull request lists the modules it adds or changes.

The minimum Go version is 1.26 (D-2026-11-11-05). CI tests the newest 1.26 patch and the current stable release.

## Pull requests

- Keep each pull request focused on one change, with tests.
- Update the README, the relevant contract in `specs/`, and the implementation
  record in `specs/001-spec-status/tasks.md` when behavior changes.
- Describe evidence plainly: say what you ran and on which platform, and do not
  claim more than that.
- `CGO_ENABLED=0 go test ./...` must pass.

## License

The project is licensed under [Apache-2.0](LICENSE). Unless you state otherwise,
a contribution you submit is licensed under the same terms (section 5 of the
license).
