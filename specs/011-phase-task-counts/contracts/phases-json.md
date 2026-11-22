# Contract: `tasks.phases` in `spec status --json`

Schema: `schema_version` stays `"1.0"`. The field is additive.

## Shape

Each item's `tasks` object:

```json
"tasks": {"checked": 5, "total": 12, "percent": 41.67, "source": "specs/014-x/tasks.md", "state": "read",
  "phases": [
    {"heading": "Phase 1: Setup", "line": 30, "story": null, "checked": 3, "total": 3},
    {"heading": "Phase 3: User Story 2 - Refuse an unprepared item (Priority: P1)", "line": 110, "story": 2, "checked": 2, "total": 4}
  ]}
```

| `state` | `phases` |
|---|---|
| `read` | array in file order; `[]` when there is no phase heading |
| `missing`, `skipped`, `not_linked` | `null` |

## Grammar

1. **Phase heading:** a line that the tasks dialect treats as an ATX heading. That means up to three spaces of indentation, not inside a fence, HTML comment, raw HTML block or indented code, and not a block-quote line. The heading must be level 2 (`##` followed by a space, a tab or the end of the line). Its text, after a CommonMark closing `#` sequence is removed and the text is trimmed, must be `Phase` or start with `Phase` followed by a space, a tab or `:`. The match is case-sensitive.
2. **End:** a phase ends at the next ATX heading of level 1 or 2 (qualifying or not), or at the end of the file. Level 3–6 headings do not end it.
3. **Membership:** a checkbox counted in the item's totals is counted in the phase that is open at its line. A checkbox outside every phase counts only in the totals.
4. **Story:** the decimal number in the first `User Story ` followed by 1–9 digits and then a non-digit or the end of the text; occurrences of `User Story ` without such a number are skipped (`User Story cleanup, then User Story 2` gives 2). Otherwise `null`.
5. **Not recognized:** setext headings neither open nor end a phase (a documented limit).

## Guarantees

- `checked ≤ total` for each phase. The sum over phases is at most the item's `checked`/`total`, with equality when every counted checkbox lies inside a phase.
- No new diagnostic and no new fatal error. Human text and exit codes are unchanged.
- Phases describe recorded checkboxes only. They do not mean done, ready or verified.
