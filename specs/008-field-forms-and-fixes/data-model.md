# Data model: Real-world field forms and small defect fixes

## Parser (`rawItem`) additions

| Field | Type | Set when | Used for |
|---|---|---|---|
| `Placeholders` | `[]placeholder{Field string; Line int}` | a `Spec dir` or `Depends on` value is a placeholder (research R2) | one `SS_FIELD_PLACEHOLDER` each |
| `DupDeps` | `[]string` | a `Depends on` list repeats an ID; each repeated ID once, in first-repeat order | one `SS_DEPENDENCY_DUPLICATE` each |

A placeholder leaves the field unset (`HasSpec` false; `Deps` nil, `DepsBad` false). Deduplicated `Deps` keep first-occurrence order.

## Field labels

| Accepted label (plain, `**label:**`, `**label**:`) | Field |
|---|---|
| `superseded-by` (existing) | `superseded-by` |
| `Superseded by` (new) | `superseded-by` |
| `superseded by` (new) | `superseded-by` |

## Placeholder value

| Value (after trimming spaces and one trailing period) | Placeholder? |
|---|---|
| `[` … `]`, whole value, one pair | yes, unless (Depends on) the inside is a valid ID list or `none` |
| `_` … `_` or `*` … `*`, whole value, single markers | yes, same exception |
| `[]`, `__`, `**` with nothing inside | `[]` yes; empty emphasis yes |
| `**…**`, `__…__`, `[x](y)`, `_x*` | no (today's behavior) |

## Display-only `Item` addition

| Field | Type | Meaning | In JSON? |
|---|---|---|---|
| `TasksSharedWith` | string | ID of the first item that read the same tasks file; "" otherwise | no |

## Summary display (no model change)

| Panel line | Source | Shown |
|---|---|---|
| `statuses <s> <n>  …` | `Summary.StatusCounts`, keys in byte order | always (`statuses none` when empty) |
| `invalid supersession N` | `Summary.InvalidSupersessionItems` | only when N > 0 |

## CLI state

| Value | Rule |
|---|---|
| JSON mode | last `--json` occurrence parsed as a boolean; `--json=false` is text |
| `--root`, `--roadmap` | empty value is an invocation error |
| `fail_on` in status invocation JSON | `requestedFailOn(args)` or `"error"` |
