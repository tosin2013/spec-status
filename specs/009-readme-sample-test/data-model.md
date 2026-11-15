# Data model: README sample test

Test-only values; nothing at runtime.

| Entity | Fields | Rules |
|---|---|---|
| Label | `Line int`, `Args []string` | parsed from `<!-- spec-sample: … -->`; `Args[0]` in {status, inventory}; `--root` present, value under `testdata/`; only allowed flags; no quote characters |
| Sample | `Line int` (opening fence), `Body string` (between fences, with trailing newline) | the first fenced block after a label with no info string |
| Output-looking block | `Line int` | a fenced block with no info string whose first line starts with `╭` or `spec: ` |

Relationships: each Label has exactly one Sample; each output-looking block must be the Sample of some Label.
