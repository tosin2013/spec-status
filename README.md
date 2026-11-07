# Spec Status

Spec Status gives a glanceable reading of a Spec Kit roadmap: recorded statuses, task checkboxes, dependency blockers, and broken references. It reads the roadmap you already have and prints those facts as text or JSON.

It helps a developer see where the work stands, a reviewer spot structural blockers, and a script fail only when a finding reaches a chosen severity. It does not decide whether the code is correct, and a recorded `verified` label is not an independent review.

## Example

This report was produced by the development build against `testdata/example-roadmap`. It is illustrative fixture output, not a released product. Exit 0 means there was no error-level finding. The warning is still there: item 003 depends on superseded item 001, and the tool does not retarget that dependency.

```
source: .specify/memory/roadmap.md
items: 3
statuses: implemented 1, planned 1, verified 1
active: 2
recorded_verified: 1/2 (50%)
tasks: 1/2 (50%) scope=unique_linked_task_files files_read=2 files_missing=0 files_skipped=0
superseded: 1
invalid_supersession: 0
blocked_active: 1
human_review_active: 1
001 implemented [historical] Preserve original feature superseded-by 002 terminal 002
002 verified [active] Replacement feature
003 planned [active] Follow-on work
  depends: 001=human_review
diagnostics:
info SS_SPEC_REFERENCE_MISSING .specify/memory/roadmap.md:14 003 Planned item has no spec link.
warning SS_DEPENDENCY_HUMAN_REVIEW .specify/memory/roadmap.md:18 003 Dependency on a superseded item requires human review.
```

## Build and run from source

These commands were run in this workspace on Linux amd64 with Go 1.22.2:

```
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build -o spec ./cmd/spec
./spec status --root testdata/example-roadmap
./spec status --root testdata/example-roadmap --json
./spec --help
./spec --version
```

Building from source requires Go. No release binary is published. An eventual installed release is intended to run without Go; that install path does not exist yet.

The Go module is `github.com/tosin2013/spec-status`. The command is `spec status`.

`--root` is the project directory that contains `.specify`. If you omit it, `spec` uses the nearest ancestor of the current directory that contains `.specify`. The roadmap is `.specify/memory/roadmap.md` unless `--roadmap` names another file inside that root.

- `0` means no finding reached the threshold. It does not mean the work is done.
- `1` means a finding reached the threshold. The default threshold is `error`.
- `2` means the command, the selected root, or the roadmap could not be evaluated.

`--fail-on` accepts `info`, `warning`, `error`, or `none`. Findings stay visible at every threshold. `none` does not hide a fatal read.

## What this build does

The development build reads a narrow heading-and-bullet roadmap shape. Unknown status words stay as written. Checkbox progress, recorded `implemented`, and recorded `verified` stay separate. Superseded items remain in the history. A dependency on a superseded item stays on that item and is marked for human review.

Text and JSON are two views of one result. Commands, help, and flags use Cobra, which the owner approved as the only dependency exception. Roadmap parsing and file reads stay in the Go standard library. The program is built with `CGO_ENABLED=0`. While it runs it does not use the network, Git, a subprocess, or write files. `spec completion bash` prints a script to stdout and does not install it.

Choices that were still open in PROFILE-R1 are used here as reversible implementation defaults, not as owner approvals. The list is in [decisions](specs/001-spec-status/decisions.md).

## Limits

Inferred quarters are not calculated. The sensitive-name filter and symlink checks are candidate safeguards, not a certification that every secret or mount escape is caught. Linux amd64 is the only platform this build was executed on. A Darwin arm64 binary was cross-compiled here and was not run. No operating system is a supported release target, and nothing is signed.

On this Linux amd64 host, nine runs of a 100-item, 10,000-task fixture had a median of about 3.9 ms. That is one development machine, not a performance promise and not an independent measurement.

There is no package-manager install, release workflow, or file-updating command. No license is selected.

## Specification

`README.md` stays at the repository root. The constitution is `.specify/memory/constitution.md`. Feature `001-spec-status` documents are under `specs/001-spec-status/`.

| Document | Purpose |
|---|---|
| [Constitution](.specify/memory/constitution.md) | Non-negotiable principles and scope governance |
| [Feature specification](specs/001-spec-status/spec.md) | User outcomes, functional requirements and edge cases |
| [Decisions](specs/001-spec-status/decisions.md) | Settled boundaries, unapproved profile, and the defaults this build uses |
| [Roadmap format](specs/001-spec-status/roadmap-format.md) | Candidate narrow grammar and optional supersession extension |
| [CLI](specs/001-spec-status/cli-contract.md), [JSON](specs/001-spec-status/json-contract.md), [diagnostics](specs/001-spec-status/diagnostics.md) | Candidate versioned machine and human contracts |
| [Technical plan](specs/001-spec-status/plan.md) | Proposed architecture and safety design |
| [Acceptance strategy](specs/001-spec-status/acceptance.md) | Required fixtures and evidence |
| [Ordered tasks](specs/001-spec-status/tasks.md) | Backlog and the current implementation record |
| [Build/release/distribution](specs/001-spec-status/build-release.md) | Future release gates |
| [Phase two](specs/001-spec-status/phase2.md) | Future write-back, not in this build |
| [Requirements/traceability](specs/001-spec-status/traceability.md) | Where each baseline requirement is recorded |

The original input-format prose is not completely recovered verbatim. This repository is a working development build of that interpretation, not an upstream compatibility certification.
