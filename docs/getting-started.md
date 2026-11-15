# Getting started

This walkthrough takes about ten minutes. You install `spec`, read its output on a small example, run it on your own project, and choose a next step. Nothing in it changes your files: `spec` only reads.

You need a terminal. You do not need Go unless you choose `go install`.

## 1. Install

Download the archive for your platform from the [latest release](https://github.com/tosin2013/spec-status/releases/latest), check its SHA-256 checksum, and put `spec` on your `PATH`. The [README Quick start](../README.md#quick-start) has the exact commands for linux/amd64 and darwin/arm64. With Go 1.26 or newer you can instead run:

```sh
CGO_ENABLED=0 go install github.com/tosin2013/spec-status/cmd/spec@latest
```

Check it:

```sh
spec version
```

The first line names the version. A release archive says `label: release build`; `go install` says `label: module build`.

## 2. Read the output on an example

Clone this repository (or download the source archive) and run `spec` on its small example roadmap:

```sh
spec status --color=never --root testdata/example-roadmap
```

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

- **The panel** describes the whole roadmap.
  - `recorded verified 1/2`: of the 2 active items, 1 is recorded `verified` in the roadmap. That is a label someone wrote, not proof that the code works.
  - `checked tasks 1/2`: checkboxes across the linked `tasks.md` files.
  - `active blockers 1`: one active item depends on something that is not recorded `verified`.
  - `human review 1`: one dependency points at a superseded item. `spec` never redirects it; a person decides.
  - `statuses ...` counts each raw status word. Historical and superseded items stay in history and earn no completion credit.
- **The table** has one row per item, in roadmap order, with its dependencies under it.
- **Diagnostics** are findings with a severity, a code and a `path:line`. Each code is explained in the [diagnostics registry](../specs/001-spec-status/diagnostics.md).

The exit code was 0: no finding reached the default threshold, `error`. Exit 0 does not mean the work is done.

## 3. Run it on your project

From anywhere inside your Spec Kit project:

```sh
spec status
```

`spec` finds the nearest folder that contains `.specify` and reads `.specify/memory/roadmap.md`. One of four things happens.

**You get a report.** Read it as in step 2. Exit 1 means a finding reached the threshold (`--fail-on`, default `error`); exit 0 means none did.

**There is no roadmap.** Many Spec Kit projects do not have one:

<!-- spec-sample: status --color=never --root testdata/two-mode/no-roadmap -->
```
spec: SS_ROADMAP_UNAVAILABLE: Required roadmap is unavailable. Run spec inventory for a feature list without a roadmap.
```

List your feature folders instead. Inventory never reads a roadmap and never infers a status:

```sh
spec inventory
```

To track statuses, dependencies and verification, add a roadmap with the [speckit-roadmap](https://github.com/srobroek/speckit-roadmap) extension.

**The roadmap is not in a supported layout.** `spec` reads the speckit-roadmap template, not arbitrary Markdown:

<!-- spec-sample: status --color=never --root testdata/two-mode/unrecognized-narrative -->
```
spec: SS_ROADMAP_UNRECOGNIZED: Roadmap is not a recognized roadmap: no item headings and no "## Planned Specs" heading. See the supported template, or run spec inventory for a feature list.
```

Items must be headings such as `### 001 — Title [status: planned]`. The [roadmap format](../specs/001-spec-status/roadmap-format.md) lists exactly what is read.

**The roadmap is empty.** The template's `## Planned Specs` heading is there, but no items yet. The report says so and exits 0:

<!-- spec-sample: status --color=never --root testdata/two-mode/explicit-empty -->
```
╭──────────────────────────────────────────────────────────────────────────────╮
│ Spec Status  development build                                               │
│ recorded verified   0/0 (N/A)                                                │
│ checked tasks       0/0 (N/A)                                                │
│ active blockers     0                                                        │
│ human review        0                                                        │
│ items 0 active 0 historical 0 superseded 0                                   │
│ statuses none                                                                │
│ task files read 0 missing 0 skipped 0                                        │
│ roadmap explicitly empty (Planned Specs heading, no items)                   │
│ source .specify/memory/roadmap.md                                            │
╰──────────────────────────────────────────────────────────────────────────────╯

ID   Status       Scope       Title

diagnostics  none
```

## 4. Understand a finding

Add `--explain` to see the file and line behind each row:

```sh
spec status --color=never --root testdata/skills/implemented-unchecked --explain
```

<!-- spec-sample: status --color=never --root testdata/skills/implemented-unchecked --explain -->
```
╭──────────────────────────────────────────────────────────────────────────────╮
│ Spec Status  development build                                               │
│ recorded verified   0/1 (0%)  [------------]                                 │
│ checked tasks       1/3 (33.33%)  [####--------]                             │
│ active blockers     0                                                        │
│ human review        0                                                        │
│ items 1 active 1 historical 0 superseded 0                                   │
│ statuses implemented 1                                                       │
│ task files read 1 missing 0 skipped 0                                        │
│ source .specify/memory/roadmap.md                                            │
╰──────────────────────────────────────────────────────────────────────────────╯

ID   Status       Scope       Title
001  implemented  active      Import data
  explain: checked tasks 1/3 in specs/001-import/tasks.md
  explain: checked checkbox at specs/001-import/tasks.md:3
  explain: unchecked checkbox at specs/001-import/tasks.md:4
  explain: unchecked checkbox at specs/001-import/tasks.md:5

diagnostics
warning SS_STATUS_AHEAD_OF_TASKS
  .specify/memory/roadmap.md:5 001 Recorded status is implemented or verified,
  but the linked tasks file has unchecked tasks.
```

Item 001 is recorded `implemented`, but 2 of its 3 tasks are unchecked. `SS_STATUS_AHEAD_OF_TASKS` points at the roadmap line; the `explain:` lines point at each checkbox. Either the work is unfinished or the boxes were never ticked; `spec` does not decide which. The reverse case, every box checked but the status still before `implemented`, is `SS_TASKS_DONE_STATUS_EARLY`.

## 5. Choose a next step

- **Add it to CI.** `spec status --json --fail-on warning` fails the job on warnings such as a status recorded ahead of its tasks. See [usage](usage.md#in-ci).
- **Get JSON for scripts.** `spec status --json` prints one document (schema 1.0); `spec inventory --json` has its own schema.
- **Ask an AI agent to review the roadmap.** Paste the [agent prompt](agent-prompt.md) into your coding agent.
- **Use the Spec Kit skills.** The [spec-status extension](../extensions/spec-status/README.md) adds six skills. Start with `/speckit.spec-status.review`. If few items are recorded verified, `/speckit.spec-status.close-gap` explains why and works through it one approved change at a time.
- **Know the limits.** [Known limits](limits.md) lists inputs that are rejected on purpose.
