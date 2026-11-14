# Data model: Status/tasks consistency diagnostics

No new types, fields or JSON keys. The feature reads existing values of one `Item` and may append one existing-shape `Diagnostic`.

## Inputs (per item, after `inspectSpec`)

| Value | Source | Used for |
|---|---|---|
| `Active` | `Item.Active` (excludes deferred, abandoned, validly superseded) | only active items are compared |
| `Status`, `StatusKnown` | `Item` | status class |
| `Tasks.State` | `ItemTasks.State` (`read`, `missing`, `skipped`, `not_linked`) | only `read` is compared |
| `Tasks.Checked`, `Tasks.Total` | `ItemTasks` | checkbox fact |
| `Tasks.Source` | `ItemTasks.Source` (canonical path for shared files) | related path |
| `Tasks.Checks[].Line`, `.Checked` | `ItemTasks.Checks` (text-only evidence) | related line |

## Status class

| Class | Statuses |
|---|---|
| before implemented | `undecided`, `needs-info`, `planned`, `specced`, `in-progress` |
| implemented or verified | `implemented`, `verified` |
| not compared | `deferred`, `abandoned` (inactive anyway), unknown tokens |

## Checkbox fact

| Fact | Condition |
|---|---|
| some unchecked | `Checked < Total` |
| all checked | `Total > 0 && Checked == Total` |
| no checkboxes | `Total == 0` |

## Decision table

| Status class | some unchecked | all checked | no checkboxes |
|---|---|---|---|
| before implemented | — | `SS_TASKS_DONE_STATUS_EARLY` (info) | — |
| implemented or verified | `SS_STATUS_AHEAD_OF_TASKS` (warning) | — | — |
| not compared | — | — | — |

## Output diagnostic

| Field | Value |
|---|---|
| Code, Severity | from the decision table |
| Message | fixed text ([contracts/diagnostics.md](contracts/diagnostics.md)) |
| ItemID | the item ID |
| Path, Line | roadmap path, item heading line |
| Related | `[{Path: Tasks.Source, Line: first unchecked checkbox}]` for ahead; `[{Path: Tasks.Source, Line: first checkbox}]` for early |

At most one of the two codes per item (the classes are disjoint). Summary counts, item fields and dependency states are unchanged.
