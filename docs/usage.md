# Using spec status

The [README](../README.md) covers installation, the commands and the exit codes. This page has the details: filters, CI, scripts, and color and width.

## Filters

Filters narrow the item rows in the text report. The summary panel always describes the whole roadmap. Diagnostics stay visible, and `--json` output never changes.

| Flag | Shows |
|---|---|
| `--status TOKEN` | Items whose raw status matches exactly (repeat the flag to allow several) |
| `--blocked` | Items with an unsatisfied dependency, including historical ones |
| `--human-review` | Items with a dependency on a superseded item |
| `--active` / `--historical` | Current items, or superseded and deferred history |

Flags combine with AND. This is `--human-review --explain` on the same example:

<!-- spec-sample: status --color=never --root testdata/example-roadmap --human-review --explain -->
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

## In CI

Fail the build when the roadmap has an error-level finding, such as a dependency on a missing item or a broken required reference:

```yaml
# GitHub Actions steps, after actions/setup-go with Go 1.26+
- run: go install github.com/tosin2013/spec-status/cmd/spec@latest   # pin a version you tested, such as @v0.2.0, for reproducible builds
- run: spec status --color=never --fail-on error
```

This repository runs exactly that check on its own roadmap in every pull request. Use `--fail-on warning` to be stricter.

## In scripts

`--json` prints one document with a stable, versioned schema (`schema_version` `"1.0"`). The JSON is identical whichever filters or color settings you use. The fields are documented in the [JSON contract](../specs/001-spec-status/json-contract.md).

```sh
spec status --json | jq '.summary | {active_items, verified_active_items, blocked_count}'
```

## Color and width

`--color=auto` (the default) uses color only when output is a terminal, `NO_COLOR` is unset, and `TERM` isn't `dumb`. `--color=always` and `--color=never` force it either way. JSON, errors, `spec version` and completion scripts are never colored. Width comes from `COLUMNS`, otherwise 80; `spec` never queries the terminal.
