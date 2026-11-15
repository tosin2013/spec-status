# Contract: README output samples

## Label

```
<!-- spec-sample: <subcommand> <flags…> -->
```

- On its own line, immediately before the sample's fenced block (blank lines allowed).
- `<subcommand>` is `status` or `inventory`; flags are those listed in research R3; `--root testdata/<fixture>` is required; no quotes.

## Sample

A fenced block with no info string, directly after the label. Its contents must equal, byte for byte, what `spec <subcommand> <flags…>` prints from the repository root with color off, `COLUMNS=80` and the development-build heading (stdout, or stderr for a fatal run).

## Rules enforced by `TestReadmeSamples`

| Situation | Result |
|---|---|
| Label and block match real output | pass |
| Block differs | fail: label command, README line, expected and actual text |
| Output-looking block (first line starts with `╭` or `spec: `) without a label | fail: README line |
| Label without a following fenced block | fail |
| Unsupported subcommand or flag, missing `--root`, root outside `testdata/`, quote in label | fail |
| Zero labels found | fail (the README is expected to have samples) |

## Current labels (after this feature)

| README section | Label |
|---|---|
| Example | `<!-- spec-sample: status --color=never --root testdata/example-roadmap -->` |
| Filters | `<!-- spec-sample: status --color=never --root testdata/example-roadmap --human-review --explain -->` |
| Without a roadmap | `<!-- spec-sample: inventory --color=never --root testdata/two-mode/no-roadmap -->` |
