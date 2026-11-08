# Spec Status

Spec Status gives a glanceable reading of a Spec Kit roadmap: recorded statuses, task checkboxes, dependency blockers, and broken references. It reads the roadmap you already have and prints those facts as text or JSON.

It helps a developer see where the work stands, a reviewer spot structural blockers, and a script fail only when a finding reaches a chosen severity. It does not decide whether the code is correct, and a recorded `verified` label is not an independent review.

## Example

This report was produced by the development build against `testdata/example-roadmap`. It is illustrative fixture output, not a released product. Exit 0 means there was no error-level finding. The warning is still there: item 003 depends on superseded item 001, and the tool does not retarget that dependency.

This is the plain report from `spec status --color=never --root testdata/example-roadmap` on the development build. A terminal, or `--color=always`, uses the same words with restrained color. Color is not the only distinction: `historical`, `active`, `recorded verified`, and `human review` are written out. `verified` here means the roadmap recorded that word. It is not an independent review.

```
╭──────────────────────────────────────────────────────────────────────────────╮
│ Spec Status  development build                                               │
│ recorded verified   1/2 (50%)  [######------]                                │
│ checked tasks       1/2 (50%)  [######------]                                │
│ active blockers     1                                                        │
│ human review        1  warning                                               │
│ items 3 active 2 historical 1 superseded 1                                   │
│ task files read 2 missing 0 skipped 0                                        │
│ source .specify/memory/roadmap.md                                            │
╰──────────────────────────────────────────────────────────────────────────────╯

ID   Status       Scope       Title
001  implemented  historical  Preserve original feature
  superseded-by 002
002  verified     active      Replacement feature
003  planned      active      Follow-on work
  human review: depends 001 human_review

diagnostics
info    SS_SPEC_REFERENCE_MISSING
  .specify/memory/roadmap.md:14 003 Planned item has no spec link.
warning SS_DEPENDENCY_HUMAN_REVIEW
  .specify/memory/roadmap.md:18 003 Dependency on a superseded item requires
  human review.
```

## Filters

Filters are a human-text view. The summary panel above stays the full roadmap, diagnostics stay visible, and `--json` still prints the unfiltered schema 1.0 document. A view that matches nothing says so and is not fatal by itself.

This is `spec status --color=never --root testdata/example-roadmap --human-review --explain` on the same fixture:

```
╭──────────────────────────────────────────────────────────────────────────────╮
│ Spec Status  development build                                               │
│ recorded verified   1/2 (50%)  [######------]                                │
│ checked tasks       1/2 (50%)  [######------]                                │
│ active blockers     1                                                        │
│ human review        1  warning                                               │
│ items 3 active 2 historical 1 superseded 1                                   │
│ task files read 2 missing 0 skipped 0                                        │
│ source .specify/memory/roadmap.md                                            │
╰──────────────────────────────────────────────────────────────────────────────╯
shown 1 of 3 items

ID   Status       Scope       Title
003  planned      active      Follow-on work
  human review: depends 001 human_review
  explain: dependency 001 is human_review at .specify/memory/roadmap.md:18; the
  successor is not used and this is not independent verification
  explain: no task file linked; checkbox line evidence unavailable

diagnostics
info    SS_SPEC_REFERENCE_MISSING
  .specify/memory/roadmap.md:14 003 Planned item has no spec link.
warning SS_DEPENDENCY_HUMAN_REVIEW
  .specify/memory/roadmap.md:18 003 Dependency on a superseded item requires
  human review.
```

`--status` matches the raw status exactly, including an unknown token. Repeat it with a different token to match any of those tokens. `--blocked`, `--human-review`, `--active`, and `--historical` combine with AND. `--blocked` includes a historical item with an unsatisfied dependency; that is separate from the active-blockers count. `--explain` adds a mechanical reason and a file:line only when the parse retained that location.

## Build and run from source

These commands were run in this workspace on Linux amd64 with Go 1.22.2. The tests, build, and the plain and colored status reports were also run on Darwin arm64 with Go 1.24.3:

```
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build -o spec ./cmd/spec
./spec status --color=never --root testdata/example-roadmap
./spec status --color=always --root testdata/example-roadmap
./spec status --root testdata/example-roadmap --json
./spec status --color=never --root testdata/example-roadmap --human-review --explain
./spec status --color=never --root testdata/example-roadmap --status planned --blocked
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

The development build reads a narrow heading-and-bullet roadmap shape, including the bold field labels (`- **Depends on:** 002.`) that the speckit-roadmap extension writes. Unknown status words stay as written. Checkbox progress, recorded `implemented`, and recorded `verified` stay separate. Superseded items remain in the history. A dependency on a superseded item stays on that item and is marked for human review.

Text and JSON are two views of one result. Human-text filters and `--explain` select or annotate item rows only; they do not change the summary, the findings, or the JSON bytes. Commands, help, and flags use Cobra. Human text uses Lip Gloss for headings, columns, the summary panel, and restrained status color. The owner approved those two dependencies and no others. Roadmap parsing and file reads stay in the Go standard library. The program is built with `CGO_ENABLED=0`. While it runs it does not use the network, Git, a subprocess, or write files.

`--color=auto` is the default: color only when stdout is a terminal, `NO_COLOR` is unset, and `TERM` is not `dumb`. `--color=always` and `--color=never` are explicit. JSON, fatal text, version, and completion stay plain. Width comes from `COLUMNS` when that value is usable, otherwise 80. The program does not query the terminal and does not open `/dev/tty`.

`spec completion bash` prints a script to stdout and does not install it. The process does not create or append `BASH_COMP_DEBUG_FILE`. The script is shell text; if you source it later, the shell may append to that variable. Spec itself does not.

Choices that were still open in PROFILE-R1 are used here as reversible implementation defaults, not as owner approvals. The list is in [decisions](specs/001-spec-status/decisions.md).

## Limits

Inferred quarters are not calculated. The sensitive-name filter and symlink checks are candidate safeguards, not a certification that every secret or mount escape is caught. This build was executed on Linux amd64 and Darwin arm64, by the implementer, on one machine each. On Darwin the post-open path check that Linux performs through `/proc/self/fd` is skipped; reads there rely on `O_NOFOLLOW`, the symlink-component check, and the device check. No operating system is a supported release target, and nothing is signed.

On this Linux amd64 host, nine runs of a 100-item, 10,000-task fixture had a median of about 3.9 ms. That is one development machine, not a performance promise and not an independent measurement.

There is no package-manager install, release workflow, or file-updating command. No license is selected.

## Draft proposals

These are not commands in this build. On 7 November 2026 the owner authorized writing them down. That did not approve implementation or the open choices.

- [Timing and evidence](specs/001-spec-status/timing-evidence.md) keeps inferred quarters omitted. Optional recorded dates are a proposal, not a forecast.
- [Companion prompts](specs/001-spec-status/companion.md) sketches an opt-in prompt file and a local evidence envelope. The core would still not embed a model.
- [Project manager sample](specs/001-spec-status/examples/prompts/project-manager.txt) and [software engineer sample](specs/001-spec-status/examples/prompts/software-engineer.txt) are editable proposed texts, not embedded prompts and not write-back.

A harness, in that companion draft, means an external tool. Claude Code, Grok Bot, Muse, and OpenCode are illustrative consumers only. They are not verified integrations. If one of them is allowed to run the CLI, it can show the human status text on a supported terminal and can read the JSON. That is not an automatic graphic and not a model inside `spec`.

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
| [Timing draft](specs/001-spec-status/timing-evidence.md), [companion draft](specs/001-spec-status/companion.md) | Draft proposals only, not implemented |
| [Requirements/traceability](specs/001-spec-status/traceability.md) | Where each baseline requirement is recorded |

The original input-format prose is not completely recovered verbatim. This repository is a working development build of that interpretation, not an upstream compatibility certification.
