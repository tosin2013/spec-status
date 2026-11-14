# Research: Real-world field forms and small defect fixes

Phase 0 for `specs/008-field-forms-and-fixes/plan.md`. Code references are to `main` at the 1.6.4 roadmap merge plus this branch (docs only). Every "today" fact was measured on 2026-11-14 with the current build (spec Context table).

## R1 — "Superseded by" spellings (row A)

- **Decision:** add `{"Superseded by", "superseded-by"}` and `{"superseded by", "superseded-by"}` to `fieldLabels` (`internal/status/parse.go:446`). `matchField` then accepts the plain and both bold forms (`label:`, `**label:**`, `**label**:`) exactly as it does for other labels.
- **Why this is enough:** the C-10 line check skips lines that `matchField` accepted (`metricLine`), and the near-miss detector in `strict.go` (`fieldPrefix`) runs only after `matchField` rejects a line, so accepted spellings never reach either fatal. Repeating the relationship under two spellings hits the existing `cur.seen[field]` "repeats a recognized field" fatal, because both map to field `superseded-by`.
- **Kept fatal:** `SUPERSEDED BY:`, `Superseded-By:` and other forms are not in the table, so `fieldPrefix` still reports them as near misses (visible failure, C-03/C-10).
- **Alternatives:** case-insensitive label matching for all fields (rejected: widens every label and weakens C-10); a lookalike-normalized alias (rejected: lookalikes are deliberately not field spellings).

## R2 — Placeholder detection (rows B, C)

- **Decision:** a helper `placeholderValue(value string) (inner string, ok bool)` applied in `applyField` to `Spec dir` and `Depends on` values, after trimming spaces and one trailing period:
  - enclosed in exactly one pair of `[` `]` with no other `]` or `[` at the ends (the whole value), or
  - enclosed in one pair of single `_` or single `*` where the second and second-to-last characters are not the same marker (so `__x__` and `**x**` are not placeholders).
- **Depends on exception (C-08):** if `parseIDList(inner)` returns a valid list (not bad), the value is **not** a placeholder; it goes through `parseIDList(value)` as today and stays `SS_DEPENDENCY_MALFORMED` (the `[` is already in the markup reject set; `_002_` is not three digits).
- **Effect:** a placeholder records `{field, line}` on the raw item (new `Placeholders []placeholder` in `rawItem`) and leaves the field unset: `HasSpec` stays false for `Spec dir`; `Deps` stays nil and `DepsBad` false for `Depends on`.
- **Markdown links** (`[x](y)`) are not wholly enclosed in brackets, so unaffected.
- **Alternatives:** matching only the template's exact text (rejected: brittle, and the italic form from review D6 would not match); treating placeholders as absent with no finding (rejected: C-01, a template line would vanish silently).

## R3 — New diagnostics

| Code | Severity | Source | Emitted in |
|---|---|---|---|
| `SS_FIELD_PLACEHOLDER` | info, report | roadmap path, the placeholder field's line; item ID | `build`, one per recorded placeholder |
| `SS_DEPENDENCY_DUPLICATE` | warning, report | roadmap path, `DepsLine`; item ID | `build`, one per repeated ID |

Messages (fixed):
- `SS_FIELD_PLACEHOLDER`: `Field holds a template placeholder and is treated as not set.`
- `SS_DEPENDENCY_DUPLICATE`: `Dependency list names the same item more than once; it is counted once.`

Placeholder findings for a `Spec dir` coexist with the stage rule's `SS_SPEC_REFERENCE_MISSING` (info for planned, error for specced and later): they state different facts (template text present; required link absent).

## R4 — Duplicate dependency IDs (row D)

- **Decision:** in `applyField` for `Depends on`, after `parseIDList`, keep the first occurrence of each ID in order and record each repeated ID once in a new `rawItem.DupDeps []string`. `build` emits one `SS_DEPENDENCY_DUPLICATE` per entry. Dependency classification, blockers and cycle detection then see each ID once, so `SS_DEPENDENCY_BLOCKED` is not doubled.
- **Measured:** no fixture has a repeated ID, so JSON is unchanged for all fixtures (spec SC-003).

## R5 — Human text additions (row E)

- **`statuses` line**: in `summaryLines` (`style.go`), after the `items …` line: `statuses ` followed by `<status> <count>` pairs joined by two spaces, keys in byte order (the order Go's JSON encoder uses for `status_counts`), raw tokens sanitized, wrapped with `wrapText`. Shown on every complete report, including an explicitly empty one (`statuses none`).
- **`invalid supersession N`**: a panel line only when `InvalidSupersessionItems > 0`.
- **Row mark**: in `writeRelations`, when `Supersession.Supplied` and `Valid` is false, print `superseded-by <id or raw reference> (invalid)`. Today a malformed reference (no parsed ID) prints no line at all; it now prints the sanitized raw reference.
- **Shared tasks**: `build` records, per canonical tasks path (`Tasks.Source`), the first item that read it; a later item with the same source gets a new display-only `Item.TasksSharedWith` (not in `itemToView`, so JSON is unchanged). `writeRelations` prints `tasks shared with <ID>; counted once`.
- **Impact:** every human report gains the `statuses` line, so the README example and any test that compares full human text change; JSON is unchanged (FR-009). The 006 README inventory sample is unaffected (inventory has its own renderer).

## R6 — `--json=false` (row F)

- **Decision:** replace `countFlag(args, "--json") > 0` (three call sites in `internal/cli/cli.go`) with `jsonRequested(args)`: the last `--json` occurrence (repeats are already an invocation error) is true for `--json`, `--json=true` and other `strconv.ParseBool` true values, false for `--json=false` and false values; an unparsable value is treated as not requested (Cobra then reports the bad value as plain text).
- The `walkArgs` helper from 006 is reused so a `--root` value spelled `--json` is not mistaken for the flag.

## R7 — Empty `--root` / `--roadmap` (row G)

- **Decision:** in `singleString.Set`, for `root` and `roadmap`, an empty value returns `--<name> needs a non-empty value`. Both `--root ""` and `--root=` reach `Set` with `""`. Applies to `spec inventory --root` too (shared flag).

## R8 — `fail_on` in status invocation JSON (row H)

- **Decision:** `fail()`'s status-shaped `invocationJSON` takes the threshold from `requestedFailOn(args)` (added in 006), defaulting to `"error"`.

## R9 — Docs and fixtures

- `specs/001-spec-status/roadmap-format.md` §2: the two spellings and the placeholder rule; `diagnostics.md`: two rows; `cli-contract.md`: `--json=false`, empty values, `fail_on`; `json-contract.md`: dependencies deduplicated (one entry per ID).
- README: remove the Limits lines for rows A–H (the two "Field labels and values" bullets, the human-text, `--json=false`, empty `--root`/`--roadmap` and `fail_on` defect lines), add findings rows for the two codes, regenerate the example panel by hand from the new output.
- `extensions/spec-status/reference/output.md`: two codes.
- New fixtures: `testdata/field-forms/` with hand-written expectations for rows A–D; unit tests for E–H.

## R9a — The raw pinned template (SC-002), measured

Running the current build on `testdata/upstream/speckit-roadmap/v0.1.0/templates/roadmap-template.md` copied into a scratch project gives a complete report (exit 1 at the default threshold) with exactly two findings, both from template placeholders: `SS_DEPENDENCY_MALFORMED` (001, line 68) and `SS_SPEC_DIRECTORY_MISSING` (001, line 71). After this feature both become `SS_FIELD_PLACEHOLDER` (info) and the run exits 0. A unit test copies the pinned template into a temp project at test time, so no second copy of the upstream file is committed.

## R10 — Constitution and constraints

No new dependency or side effect. C-03/C-10 stay strict except for two exact spellings; C-08 holds (a placeholder is never a dependency, and a bracketed valid ID list stays an error); C-01 (placeholders visible as findings); C-05 (schema 1.0, additive codes, default-threshold exit unchanged for every existing fixture).
