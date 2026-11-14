# Contract: `spec inventory` command

Amends `specs/001-spec-status/cli-contract.md` at implementation. Decision D-2026-11-12-08 (1), (2), (6); research R9.

## Synopsis

```
spec inventory [--json] [--root DIR] [--color auto|always|never] [--fail-on info|warning|error|none]
```

- Lists immediate children of `<root>/specs/` and, for feature folders, the state of `spec.md`, `plan.md`, `tasks.md` and the tasks.md checkbox counts.
- Does not read the roadmap, run Git or tests, write files, or use the network.
- Takes no positional arguments.

## Options

| Option | Meaning | Default |
|---|---|---|
| `--json` | one inventory JSON document on stdout ([inventory-json.md](inventory-json.md)) | human text |
| `--root DIR` | project directory containing `.specify`; same discovery and validation as `spec status` | nearest ancestor with a non-symlink `.specify` |
| `--color` | human text color, same rules as `spec status` | `auto` |
| `--fail-on` | lowest severity that exits 1 | `error` |

Repeated-flag rules, value validation and `--version` exclusivity are the same as `spec status`.

## Rejected options

`--roadmap`, `--status`, `--blocked`, `--human-review`, `--active`, `--historical`, `--explain` → `SS_INVOCATION_INVALID`, exit 2, message `--<name> is not used by inventory`.

## Exit codes

| Code | When |
|---|---|
| 0 | complete inventory, no finding at or above `--fail-on` |
| 1 | complete inventory, a finding at or above `--fail-on` |
| 2 | fatal (root unavailable, limit, source changed, unreliable tasks.md) or invocation error; `--fail-on none` never masks it |

## Output streams

- Human text on stdout; fatal text (`spec: CODE: message`) on stderr, plain, exit 2.
- `--json`: exactly one document on stdout for every outcome, including invocation errors, which use the inventory fatal shape when the `inventory` command is named.

## Help and completion

- `spec --help` lists `inventory`. `spec inventory --help` states that inventory does not read the roadmap and that checked boxes are not verification.
- `spec completion <shell>` gains the subcommand; it still prints to stdout only.

## Relationship to `spec status`

- `spec status` never runs inventory. A missing roadmap's message suggests it.
- `spec inventory` works whether or not a roadmap exists; it reports `roadmap_file` as a stat fact.
