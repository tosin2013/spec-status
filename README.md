# Spec Status

[![ci](https://github.com/tosin2013/spec-status/actions/workflows/ci.yml/badge.svg)](https://github.com/tosin2013/spec-status/actions/workflows/ci.yml)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

**See where a Spec Kit roadmap stands, in one command.**

`spec status` is for a developer, a reviewer, or a CI job that already has a [speckit-roadmap](https://github.com/srobroek/speckit-roadmap) v0.1.0 file at `.specify/memory/roadmap.md`. It prints the recorded statuses, checkbox counts, dependency blockers, and broken links those files contain. The same facts are available as one JSON document (schema 1.0).

It does **not**:

- use AI, the network, Git, or a subprocess, and it does not write files;
- read arbitrary Markdown or a product plan as a roadmap (for a Spec Kit tree without a roadmap, [`spec inventory`](#without-a-roadmap-spec-inventory) lists the feature folders instead);
- treat a checked task or a recorded `verified` label as proof the code is correct.

A file name does not prove the bytes match that template. A missing roadmap stops with `SS_ROADMAP_UNAVAILABLE`. A layout the reader will not guess, such as an item-like table, stops with `SS_LAYOUT_UNSUPPORTED`. A file with neither an item heading nor the template's `## Planned Specs` heading stops with `SS_ROADMAP_UNRECOGNIZED`; a file with that heading and no items is reported as explicitly empty.

**New here?** Start with the ten-minute [getting started guide](docs/getting-started.md).

## Contents

- [Quick start](#quick-start)
- [Example](#example)
- [Using it](#using-it)
- [Without a roadmap: spec inventory](#without-a-roadmap-spec-inventory)
- [Use it with an AI coding agent](#use-it-with-an-ai-coding-agent)
- [Install](#install)
- [What it reads](#what-it-reads)
- [What it never does](#what-it-never-does)
- [Build from source and project documents](#build-from-source-and-project-documents)
- [Contributing, security and license](#contributing-security-and-license)
- [Limits](#limits)

## Quick start

Install the latest release binary. You do not need Go, a git checkout, or root. The first line looks up the latest version on [Releases](https://github.com/tosin2013/spec-status/releases); to pin one instead, set it yourself (for example `V=0.2.0`). This block is Linux x86-64; macOS on Apple silicon uses `darwin_arm64` and `shasum -a 256 -c`. See [Install](#install) for details.

```sh
V=$(curl -fsSLI -o /dev/null -w '%{url_effective}' https://github.com/tosin2013/spec-status/releases/latest | sed 's#.*/tag/v##')   # latest release, without the v
printf '%s\n' "$V" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' || echo "Could not look up the latest version (got '$V'); set it yourself, for example V=0.2.0" >&2
OS_ARCH=linux_amd64          # or darwin_arm64
base=https://github.com/tosin2013/spec-status/releases/download/v$V
curl -fLO "$base/spec_${V}_${OS_ARCH}.tar.gz"
curl -fLO "$base/spec_${V}_checksums.txt"
sha256sum -c --ignore-missing "spec_${V}_checksums.txt"
tar -xzf "spec_${V}_${OS_ARCH}.tar.gz"
mkdir -p "$HOME/.local/bin"
mv "spec_${V}_${OS_ARCH}/spec" "$HOME/.local/bin/"
```

Put `$HOME/.local/bin` on your `PATH`. From a project whose roadmap is the supported template:

```sh
spec status
spec status --json
spec status --color=never --explain
spec status --color=never --human-review --explain
```

`--json` is one document on stdout. Human text is the panel. Exit `0` means no finding reached the default `error` threshold. It does not mean the work is finished. Exit `1` means a finding did. Exit `2` means the command, the root, or the roadmap could not be evaluated.

If you have Go 1.26 or newer and would rather build it yourself:

```sh
CGO_ENABLED=0 go install github.com/tosin2013/spec-status/cmd/spec@latest
```

That binary reports `label: module build`, not `label: release build`. A checkout build is in [Build from source](#build-from-source).

### When the first run looks wrong

| What you see | What it means now |
|---|---|
| `SS_ROADMAP_UNAVAILABLE`, exit 2 | `.specify/memory/roadmap.md` is missing, or `--roadmap` names a file that is not there. Spec Kit without the roadmap extension looks like this. Run [`spec inventory`](#without-a-roadmap-spec-inventory) for a feature list instead. |
| `SS_ROADMAP_UNRECOGNIZED`, exit 2 | The file is there, but it has no item headings and no `## Planned Specs` heading, so it is not the supported template (for example a hand-written narrative roadmap or an empty file). Before v0.3.0 this printed zero items with exit 0. |
| `SS_LAYOUT_UNSUPPORTED`, exit 2 | The file has an item-like layout this reader will not guess, such as a table whose first cell is a three-digit id. The JSON names the file and line. The project can still be a valid Spec Kit tree. |
| A panel with `items 0` and `roadmap explicitly empty` | The file has the template's `## Planned Specs` heading and no items yet. |
| `SS_STATUS_AHEAD_OF_TASKS`, warning | An item is recorded `implemented` or `verified`, but its `tasks.md` still has unchecked tasks. Either the work is unfinished or the boxes were never ticked. It fails a run only with `--fail-on warning` or `info`. |
| `SS_TASKS_DONE_STATUS_EARLY`, info | Every task in an item's `tasks.md` is checked, but the roadmap still says `planned`, `specced` or another status before `implemented`. The roadmap is probably out of date; `spec` does not decide that the work is verified. |
| `SS_FIELD_PLACEHOLDER`, info | A `Spec dir` or `Depends on` field still holds a template placeholder, such as the speckit-roadmap template's `[specs/001-.../ — set once the spec exists; otherwise omit]`. The field is treated as not set. A bracketed real ID such as `[002]` is still reported as malformed, never dropped. |
| `SS_DEPENDENCY_DUPLICATE`, warning | A `Depends on` list names the same item twice. It is counted once. |

## Example

This is `spec status --color=never --root testdata/example-roadmap`, run on the small example roadmap in this repository. In a terminal the same words appear with restrained color; color is never the only signal.

<!-- spec-sample: status --color=never --root testdata/example-roadmap -->
```
╭──────────────────────────────────────────────────────────────────────────────╮
│ Spec Status  development build                                               │
│ recorded verified   1/2 (50%)  [######------]                                │
│ checked tasks       1/2 (50%)  [######------]                                │
│ active blockers     1                                                        │
│ human review        1  warning                                               │
│ items 3 active 2 historical 1 superseded 1                                   │
│ statuses implemented 1 planned 1 verified 1                                  │
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

How to read it:

- **recorded verified 1/2:** of the two active items, one is recorded `verified`. Item 001 was replaced by 002, so it is history and doesn't count.
- **checked tasks 1/2:** checkboxes ticked across the linked `tasks.md` files.
- **active blockers 1:** item 003 depends on 001, which was superseded. `spec` does not quietly redirect that dependency to 002. It asks for **human review** instead.
- **diagnostics:** each finding has a severity, a stable code, and the file and line it came from.
- The exit code is 0, because no finding reached the default `error` threshold. The warning is still shown.

## Using it

### Commands and common flags

| Command | What it does |
|---|---|
| `spec status` | Report on the nearest project root (the closest parent directory with `.specify`) |
| `spec status --root DIR` | Report on a specific project |
| `spec status --roadmap FILE` | Use another roadmap file inside the root |
| `spec status --json` | Print one JSON document (schema 1.0) instead of text |
| `spec status --fail-on LEVEL` | Exit 1 when a finding reaches `info`, `warning`, `error` (default) or `none` |
| `spec status --explain` | Add the mechanical reason and file:line behind each row |
| `spec inventory` | List feature folders under `specs/` without a roadmap (see [below](#without-a-roadmap-spec-inventory)) |
| `spec version` | Show the version, commit, Go version, target, and how the binary was built |
| `spec completion bash` | Print a shell completion script (zsh, fish and powershell too) |

### Exit codes

| Code | Meaning |
|---|---|
| `0` | No finding reached the `--fail-on` threshold. This does **not** mean the work is done. |
| `1` | A finding reached the threshold. |
| `2` | The command, the project root, or the roadmap could not be read or evaluated. `--fail-on none` never hides this. |

More: [filters, CI, scripts, and color and width](docs/usage.md).

## Without a roadmap: spec inventory

Many Spec Kit projects have no roadmap file. `spec inventory` lists the feature folders directly under `specs/` instead: whether `spec.md`, `plan.md` and `tasks.md` are there, and the checkboxes in `tasks.md`. It never reads the roadmap, and it does not infer a status, priority or dependency from folder names or checkboxes. Checked boxes are not verification. `spec status` never switches to inventory on its own.

A folder counts as a feature folder when its name has one of the two forms Spec Kit creates: `NNN-name` (three or more digits) or `YYYYMMDD-HHMMSS-name`. Other folders are listed as not feature folders. Nothing deeper is opened, so requirements checklists under `checklists/` are never counted.

This is `spec inventory --color=never --root testdata/two-mode/no-roadmap`:

<!-- spec-sample: inventory --color=never --root testdata/two-mode/no-roadmap -->
```
╭──────────────────────────────────────────────────────────────────────────────╮
│ Spec Inventory  development build                                            │
│ mode                inventory (roadmap not read; no status inferred)         │
│ checked tasks       3/5 (60%)  [#######-----]                                │
│ folders 6 feature 3 other 3 skipped 0                                        │
│ task files read 2 missing 1 skipped 0                                        │
│ source specs                                                                 │
╰──────────────────────────────────────────────────────────────────────────────╯

Folder                 spec     plan     tasks
001-alpha              present  present  3/5
003_name               not a feature folder
20261101-093000-beta   present  missing  missing
20261101-gamma         missing  missing  0/0
3-short                not a feature folder
notes                  not a feature folder

Checked boxes are not verification.

diagnostics
info    SS_FOLDER_NOT_FEATURE
  specs/003_name Folder name is not a Spec Kit feature folder name.
info    SS_FOLDER_NOT_FEATURE
  specs/3-short Folder name is not a Spec Kit feature folder name.
info    SS_FOLDER_NOT_FEATURE
  specs/notes Folder name is not a Spec Kit feature folder name.
```

`spec inventory --json` prints one document with its own version, `inventory_schema_version` 1.0 ([contract](specs/001-spec-status/inventory-json.md)). It takes `--root`, `--color` and `--fail-on` like `spec status`; the roadmap options are rejected. Exit codes mean the same as for `spec status`.

## Use it with an AI coding agent

[docs/agent-prompt.md](docs/agent-prompt.md) has a prompt you can paste into an AI coding agent (for example Claude Code, Cursor or Codex) that runs in the repository you want reviewed. The agent installs the latest release, asking you before it installs anything and checking the checksum. It then runs `spec status` and summarizes progress, blockers, items needing human review, diagnostics, and what is ready to start, without changing any files. If `spec` itself misbehaves, it drafts a GitHub issue from the [bug report form](https://github.com/tosin2013/spec-status/issues/new?template=bug_report.yml) and files it only with your approval. Security problems go to [private reporting](SECURITY.md). The same prompt is on every release page.

### As Spec Kit skills

If the project uses [Spec Kit](https://github.com/github/spec-kit) 1.2.0 or newer, the **spec-status extension** adds six skills to your AI coding agent: review (with pm, developer, qa and architect roles), explain, reconcile, unblock, verify-readiness and close-gap. They never change a file on their own: each change is proposed, applied through a Spec Kit command only after you approve it, and confirmed by running `spec status` again.

```sh
specify extension add spec-status --from https://github.com/tosin2013/spec-status/releases/download/v<version>/spec-status-extension_<version>.zip
```

The [extension guide](extensions/spec-status/README.md) explains each skill, the roles, the approval protocol and partner extensions.

## Install

Releases are published on [GitHub Releases](https://github.com/tosin2013/spec-status/releases) by the [release workflow](.github/workflows/release.yml). Each release has one archive per supported platform and a SHA-256 checksum file. Every archive is smoke-tested on its own platform before the owner publishes the release.

| Platform | Archive | Supported since |
|---|---|---|
| Linux, x86-64 (linux/amd64) | `spec_<version>_linux_amd64.tar.gz` | v0.1.2 ([smoke test](https://github.com/tosin2013/spec-status/actions/runs/43160883770/job/130992118717)) |
| macOS, Apple silicon (darwin/arm64) | `spec_<version>_darwin_arm64.tar.gz` | v0.1.2 ([smoke test](https://github.com/tosin2013/spec-status/actions/runs/43160883770/job/130992118760)) |

The first release is `v0.1.2`. Version `v0.1.0` was published directly rather than through the release workflow, with untested archives that included unsupported platforms. It was withdrawn and should not be used (decision D-2026-11-12-02). The `v0.1.1` tag was built by the workflow but never published, because the agent prompt and the `v0.1.0` retraction landed right after it (D-2026-11-12-03).

You need neither root nor Go. The first line looks up the latest release; to install or pin a specific version, set `V` yourself (for example `V=0.2.0`):

```sh
V=$(curl -fsSLI -o /dev/null -w '%{url_effective}' https://github.com/tosin2013/spec-status/releases/latest | sed 's#.*/tag/v##')   # latest release, without the v
printf '%s\n' "$V" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' || echo "Could not look up the latest version (got '$V'); set it yourself, for example V=0.2.0" >&2
OS_ARCH=linux_amd64          # or darwin_arm64
base=https://github.com/tosin2013/spec-status/releases/download/v$V
curl -fLO "$base/spec_${V}_${OS_ARCH}.tar.gz"
curl -fLO "$base/spec_${V}_checksums.txt"
sha256sum -c --ignore-missing "spec_${V}_checksums.txt"        # macOS: shasum -a 256 -c --ignore-missing ...
tar -xzf "spec_${V}_${OS_ARCH}.tar.gz"
mkdir -p ~/.local/bin && mv "spec_${V}_${OS_ARCH}/spec" ~/.local/bin/
spec version                                                   # label: release build
```

Make sure `~/.local/bin`, or whichever directory you chose, is on your `PATH`. The checksum shows the archive matches what the release workflow built. It does not prove who published it, and the binaries are not signed.

- **Upgrade or roll back:** repeat the same steps with another version. Published releases stay downloadable, and their files are never replaced; a fix ships as a new version.
- **macOS:** the binaries are not signed or notarized. Files downloaded with `curl` run directly. If you download with a browser, macOS quarantines the file and blocks the first run. Clear it for that one file with `xattr -d com.apple.quarantine ~/.local/bin/spec`. Don't disable Gatekeeper.
- **Other platforms** (for example Windows, linux/arm64 or Intel Macs) have no release archive. Building from source or using `go install` may work there, but it is untested and unsupported.
- **With Go 1.26 or newer:** `go install github.com/tosin2013/spec-status/cmd/spec@latest` (or a specific version, such as `@v0.2.0`). Your own toolchain builds this binary from the published module, so `spec version` reports `label: module build (built from the module version, not a release archive)` and `commit: unknown`.
- **Compatibility before 1.0:** a minor version (0.x.0) may change command flags, exit codes or the JSON schema, and the release notes list every such change. Patch versions (0.x.y) don't break compatibility. The JSON schema has its own version, currently 1.0.

Maintainers: the release steps are in [RELEASING.md](RELEASING.md).

## What it reads

- **The roadmap** at `.specify/memory/roadmap.md`, in the shape written by the [speckit-roadmap](https://github.com/srobroek/speckit-roadmap) extension (v0.1.0). Items are headings like `### 001 — Title [status: planned]`, followed by fields such as `- **Depends on:** 002`, `- **Spec dir:** specs/001-name`, and `- superseded-by: 003`.
- **Task files**: the `tasks.md` in each item's spec directory, in the shape produced by [github/spec-kit](https://github.com/github/spec-kit) v1.2.0. `spec.md` and `plan.md` are only checked for existence.
- **With `spec inventory`**: the names of the folders directly under `specs/`, and in each feature folder the same three files, read the same way. The roadmap is not read.

The project root is the directory with `.specify`. `--root` sets it; otherwise `spec` uses the nearest ancestor of the current directory.

`spec` reads a deliberately narrow format. It doesn't guess at layouts it can't read reliably: they stop the report with an error, or are listed under [Limits](#limits). The exact shapes are in the [shape inventory](specs/002-upstream-format-conformance/contracts/shape-inventory.md), checked by hand-counted fixtures in `testdata/conformance/`. A newer upstream version is only adopted after the drift check in `internal/conformance` shows no unresolved change (see `testdata/upstream/README.md`).

Unknown status words are kept exactly as written, with a warning. Checked tasks, recorded `implemented`, and recorded `verified` are always kept separate.

## What it never does

While it runs, `spec`:

- never writes, creates or deletes files, and keeps no cache;
- never uses the network, Git, subprocesses, telemetry or AI;
- never reads outside the project root, follows symlinks, or opens sensitive files such as `.env`, `.ssh/`, `*.pem` or `*.key`;
- never infers dates, estimates or quarters;
- never treats a superseded item's successor as satisfying a dependency.

It is a single static binary built with `CGO_ENABLED=0`. Roadmap parsing and file access use only the Go standard library. The only other dependencies are [Cobra](https://github.com/spf13/cobra), for commands and flags, and [Lip Gloss](https://github.com/charmbracelet/lipgloss), for terminal layout and color, both approved by the owner.

## Build from source and project documents

You need Go 1.26 or newer: `CGO_ENABLED=0 go build -o spec ./cmd/spec`. A checkout build reports itself as a `development build`. [docs/development.md](docs/development.md) covers the tests, what CI checks, and the project's design documents; [docs/architecture.md](docs/architecture.md) shows how the pieces fit together.


## Contributing, security and license

- **Contributing:** see [CONTRIBUTING.md](CONTRIBUTING.md). Bug reports and feature requests use the [issue forms](https://github.com/tosin2013/spec-status/issues/new/choose).
- **Security:** report vulnerabilities privately, as described in [SECURITY.md](SECURITY.md).
- **Conduct:** [Code of Conduct](CODE_OF_CONDUCT.md).
- **License:** [Apache-2.0](LICENSE).

## Limits

The short version:

- `spec` reads a narrow, pinned format, not general Markdown. Anything it can't read reliably stops the report with an error or is listed in [docs/limits.md](docs/limits.md).
- `verified` means the roadmap says so. It is not an independent review, and exit 0 does not mean the work is done.
- Only linux/amd64 and darwin/arm64 are supported, each since v0.1.2 after its release smoke test. Releases carry SHA-256 checksums and are not signed or notarized (decision D-2026-11-11-02).
- The sensitive-file filter and symlink checks are safeguards, not a guarantee that every secret or mount escape is caught.
- Performance: on one Linux amd64 development machine, nine runs over a 100-item, 10,000-task fixture had a median of about 3.9 ms. That is not a performance promise.
- The roadmap is recognized only by an item heading or the template's exact `## Planned Specs` line. A roadmap written in another layout stops with `SS_ROADMAP_UNRECOGNIZED`, even if a person would call it a roadmap.
- `spec inventory` skips, with an error finding, any feature folder or file whose name the sensitive-name filter matches (for example `007-client-secret`), because that filter is deliberately over-inclusive.

The full list of known limits, including inputs that are rejected on purpose, is in [docs/limits.md](docs/limits.md).
