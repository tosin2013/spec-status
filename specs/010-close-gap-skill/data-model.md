# Data model: Close-gap skill

The skill holds no persistent state; these are the concepts it reasons with, taken from `spec status --json` (schema 1.0).

## Gap item

| Field | Source |
|---|---|
| id, title, status, status_known | `items[]` |
| active | `items[].active` (only active items are classified) |
| cause | first match in research R1's precedence table |
| evidence | `items[].source` (path:line) and the matching diagnostic's `source` |
| blocking dependency | `items[].dependencies[]` with `state` ≠ `satisfied` (cause blocked) |

Causes, in classification precedence: `unknown-status`, `missing-artifacts`, `status-ahead-of-tasks`, `tasks-done-status-early`, `blocked`, `ready-for-debrief`; not part of the gap: `not-started`, `in-progress` (no recorded inconsistency).

## Gap summary

| Field | Source |
|---|---|
| recorded verified X/Y | `summary.verified_active_items` / `summary.active_items` (as `spec` prints) |
| per-cause counts and IDs | the classified gap items |
| order of work | research R2 |

## Proposal

| Field | Value |
|---|---|
| item, cause, evidence | from the gap item |
| change | one change (research R3) |
| command | the Spec Kit command that applies it |
| state | `proposed` → `approved` → `applied` → `confirmed` (re-read shows the change) or `not confirmed`; or `declined` / `skipped` |

Transitions: only the user moves `proposed` to `approved`, `declined` or `skipped`. `applied` happens only through the named command. `confirmed` requires a re-run of `spec status`. `verified` can be the change only for a `ready-for-debrief` item after a debrief with no must-address finding and the user's agreement.
