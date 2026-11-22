# Data model: Per-phase task counts

## Phase (new)

| Field | Type | Rule |
|---|---|---|
| Heading | string | ATX heading text, closing `#`s removed, trimmed (R2) |
| Line | int | 1-based line of the heading in the tasks file |
| Story | *int | N from the first `User Story N` (1–9 digits, then a non-digit or the end), else nil |
| Checked | int | checked checkboxes inside the phase |
| Total | int | checkboxes inside the phase |

A Phase is created at a qualifying level-2 heading and closed at the next level-1 or level-2 ATX heading, or at the end of the file. Invariants:
- Phases are in file order.
- `Checked ≤ Total`.
- `Σ Total ≤ taskFile.Total`, and `Σ Checked ≤ taskFile.Checked`.

## Changed types

- `taskFile` (`parse.go`): adds `Phases []Phase`. Empty when the file has no qualifying headings.
- `ItemTasks` (`model.go`): adds `Phases []Phase`, copied from the `taskFile` when the file is read (next to the `Checks` copy in `analyze.go`). Items that share one file each get a copy.
- `itemTaskView` (`render.go`): adds `Phases *[]phaseView` (`json:"phases"`): `[]`/list for state `read`, `null` otherwise.

Text rendering (`style.go`) does not read `Phases`.
