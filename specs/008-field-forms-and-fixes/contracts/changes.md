# Contract: field forms, new diagnostics, display and CLI changes

Amends, at implementation: `specs/001-spec-status/roadmap-format.md` §2, `diagnostics.md`, `json-contract.md`, `cli-contract.md`, README, and `extensions/spec-status/reference/output.md`. Decisions D-2026-11-14-01 (3), (4); rules in [../data-model.md](../data-model.md).

## 1. Grammar (`roadmap-format.md` §2)

- **Superseded by.** `Superseded by` and `superseded by` are accepted spellings of `superseded-by`, plain or bold like other labels. Other case or punctuation variants remain near misses (fatal). Two spellings on one item are a repeated field (fatal).
- **Placeholders.** A `Spec dir` or `Depends on` value that is wholly enclosed in one pair of `[` `]`, or of single `_` or `*`, is a template placeholder: the field is treated as not given and `SS_FIELD_PLACEHOLDER` (info) is reported at its line. Exception: when the enclosed text is itself a valid `Depends on` list or `none`, the value is read as today and is `SS_DEPENDENCY_MALFORMED`. A placeholder never satisfies a stage requirement, a dependency or a spec link.
- **Duplicates.** A `Depends on` list that names an ID more than once lists it once (first position); `SS_DEPENDENCY_DUPLICATE` (warning) is reported once per repeated ID at the field line.

## 2. Registry rows (`diagnostics.md`)

| Code | Severity / class | Trigger and source | Authority status |
|---|---|---|---|
| SS_FIELD_PLACEHOLDER | info, report | A `Spec dir` or `Depends on` value that is a template placeholder; source is the field line; item ID set | D-2026-11-14-01 (3), feature 008 |
| SS_DEPENDENCY_DUPLICATE | warning, report | A `Depends on` list repeats an ID; one per repeated ID; source is the field line; item ID set | D-2026-11-14-01 (4), feature 008 |

Messages (exact):
- `SS_FIELD_PLACEHOLDER`: `Field holds a template placeholder and is treated as not set.`
- `SS_DEPENDENCY_DUPLICATE`: `Dependency list names the same item more than once; it is counted once.`

## 3. JSON (`json-contract.md`, schema 1.0, no field change)

- `dependencies` contains each ID at most once, in first-declared order.
- New codes appear in the existing diagnostic shape.

## 4. Human text

Summary panel, after `items …`:

```
│ statuses implemented 1 planned 1 verified 1                                  │
```

When `invalid_supersession_items > 0`, an extra panel line `invalid supersession N`.

Row relations:

```
001  implemented  active      Preserve original feature
  superseded-by 009 (invalid)
002  verified     active      Replacement feature
  tasks shared with 001; counted once
```

The exact panel is regenerated in the README example from a real run and checked by hand.

## 5. CLI (`cli-contract.md`)

| Invocation | Result |
|---|---|
| `spec status --json=false --bogus` | plain-text `SS_INVOCATION_INVALID` on stderr, exit 2 |
| `spec status --json=true --bogus` | JSON invocation error, exit 2 |
| `spec status --root ""` / `--root=` / `--roadmap ""` / `--roadmap=` | `SS_INVOCATION_INVALID: --root needs a non-empty value` (or `--roadmap`), exit 2 |
| `spec inventory --root ""` | same |
| `spec status --json --fail-on warning --bogus` | JSON invocation error with `"fail_on":"warning"` |
| `spec status --json --fail-on bogus` | JSON invocation error with `"fail_on":"error"` |

## 6. Fixtures (`testdata/field-forms/`, hand-written `expected-status.json`)

| Fixture | Content | Expected |
|---|---|---|
| `superseded-bold` | `- **Superseded by:** 002` | report; 001 `superseded_by` 002 |
| `superseded-plain` | `- Superseded by: 002` | same |
| `superseded-lower` | `- superseded by: 002` | same |
| `superseded-upper` | `- SUPERSEDED BY: 002` | fatal `SS_ITEM_MALFORMED` |
| `superseded-twice` | both `superseded-by` and `Superseded by` | fatal `SS_ITEM_MALFORMED` |
| `template-placeholders` | the pinned template's `Spec dir` and `Depends on` placeholder lines on a `planned` item | exit 0; two `SS_FIELD_PLACEHOLDER`; `SS_SPEC_REFERENCE_MISSING` info; no error |
| `placeholder-italic` | `- **Spec dir:** _set once the spec exists_` on a `specced` item | `SS_FIELD_PLACEHOLDER` + `SS_SPEC_REFERENCE_MISSING` error |
| `bracketed-real-id` | `- Depends on: [002]` | `SS_DEPENDENCY_MALFORMED` (error), no placeholder |
| `duplicate-deps` | `- Depends on: 002, 003, 002, 002` | dependencies `002, 003`; one `SS_DEPENDENCY_DUPLICATE`; `SS_DEPENDENCY_BLOCKED` once per unfinished ID |
