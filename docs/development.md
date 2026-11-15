# Development

How to build and test `spec`, what CI checks, and where the project's design documents are. The system overview is in [architecture.md](architecture.md).

## Build from source

You need Go 1.26 or newer.

```sh
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build -o spec ./cmd/spec
./spec status --color=never --root testdata/example-roadmap
./spec status --root testdata/example-roadmap --json
./spec status --color=never --root testdata/example-roadmap --human-review --explain
./spec --help
```

A binary built from a source checkout reports itself as a `development build`.

Every pull request and push to `main` runs [CI](../.github/workflows/ci.yml):
- `go vet`, the tests and `govulncheck`, on linux/amd64 and darwin/arm64, with the newest Go 1.26 patch and the current stable Go;
- a hygiene job: formatting, `go.mod`/`go.sum` tidiness and integrity, `spec status --fail-on error` on this repository's own roadmap, and compile-only builds for windows/amd64 and linux/arm64;
- on pull requests, a dependency review.

CI passes are not release evidence. Releases add their own smoke tests (see [RELEASING.md](../RELEASING.md)). Dependabot proposes dependency updates, which merge under the rule in [CONTRIBUTING.md](../CONTRIBUTING.md).

## Project documents

This project is built with [Spec Kit](https://github.com/github/spec-kit) and tracks its own work in [its roadmap](../.specify/memory/roadmap.md), which `spec status` reads in CI.

| Document | Purpose |
|---|---|
| [Constitution](../.specify/memory/constitution.md) | The project's non-negotiable principles |
| [Roadmap](../.specify/memory/roadmap.md) | Planned, current and finished work |
| [Feature specification](../specs/001-spec-status/spec.md) | User outcomes, requirements and edge cases |
| [Decisions](../specs/001-spec-status/decisions.md) | Every owner decision, with its date and alternatives |
| [Roadmap format](../specs/001-spec-status/roadmap-format.md) | The exact grammar that is read |
| [CLI](../specs/001-spec-status/cli-contract.md), [JSON](../specs/001-spec-status/json-contract.md), [diagnostics](../specs/001-spec-status/diagnostics.md) | The versioned command, output and finding contracts |
| [Upstream conformance](../specs/002-upstream-format-conformance/) | Pinned upstream formats, shape inventory, drift check |
| [Dependency automation](../specs/003-dependency-updates/), [Release engineering](../specs/004-release-engineering/) | CI, Dependabot, and the release workflow |
| [Agent skills](../specs/005-agent-skills/), [close-gap](../specs/010-close-gap-skill/) | The extension and its skills |
| [Two modes](../specs/006-two-mode-inventory/), [status/tasks diagnostics](../specs/007-status-tasks-diagnostics/), [field forms](../specs/008-field-forms-and-fixes/), [README samples](../specs/009-readme-sample-test/) | The v0.3.0 features |
| [Build and release requirements](../specs/001-spec-status/build-release.md) | Release gates and distribution rules |
| [Tasks and implementation record](../specs/001-spec-status/tasks.md) | What was built, with evidence |
| [Two-mode direction](../specs/001-spec-status/two-mode-direction.md) | The two modes: roadmap report and roadmap-free inventory. Implemented in [feature 006](../specs/006-two-mode-inventory/) |
| [Timing](../specs/001-spec-status/timing-evidence.md), [companion](../specs/001-spec-status/companion.md), [phase two](../specs/001-spec-status/phase2.md) | Draft proposals only. Not implemented |
