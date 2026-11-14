# Feature Specification: Real-world field forms and small defect fixes

**Feature Branch**: `008-fields-fixes-specify`

**Created**: 2026-11-14

**Status**: Draft

**Input**: Roadmap items 015 ("Real-world field forms") and 016 ("Small defect fixes"), chosen for v0.3.0 in decision D-2026-11-14-01 (3) and (4). Both are small changes to what `spec status` accepts and reports, so they share one feature. The exact accepted forms are set here, consistent with C-08 (no silently misread `Depends on`, `Spec dir` or `superseded-by`) and C-10 (line-level label rule).

## Context

Each item comes from README Limits ("visible over-strictness" or "open non-blocking defects"). Behavior measured with the v0.2.0-era build on 2026-11-14:

| # | Input or invocation | Today | Origin |
|---|---|---|---|
| A | `- **Superseded by:** 002`, `- Superseded by: 002`, `- superseded by: 002` | fatal `SS_ITEM_MALFORMED` (the C-10 label rule), exit 2 | review D6 (`specs/002-upstream-format-conformance/reclassification.md`) |
| B | `- **Spec dir:** [specs/001-.../ — set once the spec exists; otherwise omit]` (the pinned template's own text) or `_set once the spec exists_` | value read as a path; error `SS_SPEC_DIRECTORY_MISSING`, exit 1 | review D6; template line 71 |
| C | `- **Depends on:** [NNN, NNN — or "none"]` (the pinned template's own text) | error `SS_DEPENDENCY_MALFORMED`, exit 1 | template line 68 |
| D | `- Depends on: 002, 002` | dependency listed twice; `SS_DEPENDENCY_BLOCKED` emitted twice for one fact | README open defects |
| E | Human text | never shows `status_counts` or `invalid_supersession_items`; an invalid `superseded-by` looks like a valid one; a tasks file shared by two items is not attributed | README open defects |
| F | `spec status --json=false --bogus` | prints a JSON error although JSON was turned off | README open defects |
| G | `spec status --root ""` or `--roadmap ""` | silently uses the default | README open defects |
| H | `spec status --json --fail-on warning --bogus` | invocation-error JSON reports `fail_on: "error"` | README open defects (D10) |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - A capitalized "Superseded by" field is read, not fatal (Priority: P1)

A maintainer writes `- **Superseded by:** 002`, the way the other labels are capitalized, and `spec status` reads it as the `superseded-by` relationship instead of stopping the whole report.

**Why this priority**: today one natural spelling makes the entire report fatal; it is the most disruptive of these limits.

**Independent Test**: the three spellings in row A each produce the same report as `- superseded-by: 002`.

**Acceptance Scenarios**:

1. **Given** an item with `- **Superseded by:** 002`, `- Superseded by: 002` or `- superseded by: 002` (plain or bold, the bold forms as for other labels), **When** `spec status` runs, **Then** the item's `superseded_by` is `002` and the report is byte-identical to one using `- superseded-by: 002` on the same line.
2. **Given** an item with both `superseded-by` and `Superseded by`, **When** `spec status` runs, **Then** it is the existing "repeats a recognized field" fatal (one relationship, two spellings).
3. **Given** `- SUPERSEDED BY: 002` or `- Superseded-By: 002` (other case forms), **When** `spec status` runs, **Then** it stays the existing visible fatal; only the three listed spellings are added.
4. **Given** prose containing `Superseded by:` outside a field line, **When** `spec status` runs, **Then** the C-10 fatal is unchanged.

---

### User Story 2 - A template placeholder left in a field means "not set yet" (Priority: P1)

A maintainer copies the speckit-roadmap template and leaves `- **Spec dir:** [specs/001-.../ — set once the spec exists; otherwise omit]` or `- **Depends on:** [NNN, NNN — or "none"]` in a new entry. `spec status` treats the field as not set and says so with an informational finding, instead of reporting a missing directory or a malformed list as an error.

**Why this priority**: the pinned template's own text currently produces errors, so a freshly templated roadmap fails a default CI run.

**Independent Test**: a roadmap with the template's placeholder lines exits 0 at the default threshold and reports one info finding per placeholder; a bracketed real ID is still an error.

**Acceptance Scenarios**:

1. **Given** a `Spec dir` or `Depends on` value that is entirely enclosed in `[` `]`, or entirely in single `_` or `*` emphasis, **When** `spec status` runs, **Then** the field is treated as not given and a new `SS_FIELD_PLACEHOLDER` (info) is reported at that line.
2. **Given** a placeholder `Spec dir` on a `planned` item, **When** `spec status` runs, **Then** the item behaves as having no spec link (the existing `SS_SPEC_REFERENCE_MISSING` info) and no `SS_SPEC_DIRECTORY_MISSING` error.
3. **Given** a placeholder `Spec dir` on a `specced` or later item, **When** `spec status` runs, **Then** the existing stage rule still reports the missing spec link as an error: a placeholder never satisfies a stage.
4. **Given** `- Depends on: [002]` or `- Depends on: _002, 003_` (enclosed text that is itself a valid ID list or `none`), **When** `spec status` runs, **Then** it stays `SS_DEPENDENCY_MALFORMED` (error): a real dependency is never silently dropped (C-08).
5. **Given** a Markdown link such as `- Spec dir: [specs/001-x](specs/001-x)`, **When** `spec status` runs, **Then** behavior is unchanged (it is not entirely enclosed in brackets).

---

### User Story 3 - A repeated dependency ID is reported once, with a warning (Priority: P2)

A maintainer types `- Depends on: 002, 002`. The report lists the dependency once and warns about the repetition, instead of listing it twice and doubling every dependency finding.

**Why this priority**: it is a data-quality defect that inflates findings, but it never blocked a report.

**Independent Test**: the row-D roadmap gives one dependency entry, one `SS_DEPENDENCY_DUPLICATE` warning and one `SS_DEPENDENCY_BLOCKED`.

**Acceptance Scenarios**:

1. **Given** a `Depends on` list naming the same ID more than once, **When** `spec status` runs, **Then** the item's `dependencies` contain that ID once, at its first position, and one `SS_DEPENDENCY_DUPLICATE` (warning) per repeated ID is reported at the field line.
2. **Given** that dependency is unfinished, **When** `spec status` runs, **Then** `SS_DEPENDENCY_BLOCKED` is reported once for it.
3. **Given** a list without repeats, **When** `spec status` runs, **Then** output is unchanged.

---

### User Story 4 - Human text shows what the JSON already knows (Priority: P2)

A reviewer reading the human report sees the per-status counts, how many items have an invalid supersession, which `superseded-by` values are invalid, and which item a shared tasks file is counted under, without opening the JSON.

**Why this priority**: text and JSON are two views of one report (constitution Art. II); these four gaps are the known places where the text view is thinner.

**Independent Test**: on fixtures with each case, the human text contains the new lines; the JSON is byte-identical to before.

**Acceptance Scenarios**:

1. **Given** any complete report, **When** human text is printed, **Then** the summary panel shows a `statuses` line listing each raw status with its count, in the same key order as JSON `status_counts`.
2. **Given** a report with one or more invalid supersessions, **When** human text is printed, **Then** the panel shows `invalid supersession N`, and each affected row's `superseded-by` line is marked `(invalid)`. With none, neither appears.
3. **Given** two items whose `Spec dir` resolves to the same tasks file, **When** human text is printed, **Then** the second item's row says that its tasks are shared with the first item and counted once.
4. **Given** any input, **When** JSON is printed, **Then** it is byte-identical to the previous build for this story (display only).

---

### User Story 5 - Command-line edge cases behave as written (Priority: P3)

A script author gets errors that match what they asked for: `--json=false` really means text, an empty `--root` or `--roadmap` is refused instead of silently ignored, and a JSON invocation error reports the threshold they passed.

**Why this priority**: small, but each one surprises scripts.

**Independent Test**: rows F, G and H each give the new behavior.

**Acceptance Scenarios**:

1. **Given** `spec status --json=false --bogus`, **When** it runs, **Then** the invocation error is plain text on stderr, exit 2; `--json=true` and `--json` still give one JSON document.
2. **Given** `--root ""`, `--root=`, `--roadmap ""` or `--roadmap=`, **When** `spec status` (or `spec inventory` for `--root`) runs, **Then** it is `SS_INVOCATION_INVALID` (`--root needs a non-empty value`), exit 2.
3. **Given** `spec status --json --fail-on warning --bogus`, **When** it runs, **Then** the invocation-error JSON reports `fail_on: "warning"`; an invalid or missing `--fail-on` value reports `"error"`, as inventory already does.

---

### Edge Cases

- `Superseded by` with a placeholder value (`[NNN]`): placeholder rules apply only to `Spec dir` and `Depends on`; a non-ID `superseded-by` value stays the existing `SS_SUPERSESSION_MALFORMED`.
- Emphasis forms: `_text_` and `*text*` count; `__text__`, `**text**` (bold values) and mixed `_text*` do not, and keep today's behavior.
- An empty bracket `[]` or empty emphasis is a placeholder (nothing inside to misread).
- `Depends on: 002, 003, 002` → `002, 003` with one duplicate warning.
- `Depends on: 002, 002, 002` → one entry and one warning for 002 (per repeated ID, not per occurrence).
- The `statuses` line wraps like other panel lines at narrow widths; unknown status tokens appear raw and sanitized.
- Filters and `--explain` keep their rule: they never change the panel, diagnostics, exit code or JSON.

## Requirements *(mandatory)*

### Functional Requirements

**015 — field forms**

- **FR-001**: The field labels `Superseded by` and `superseded by` (plain, or bold as other labels are) MUST be read exactly like `superseded-by`, including the C-10 line-level rule and the closed-field-body rule. Other case or punctuation variants keep today's behavior.
- **FR-002**: A `Spec dir` or `Depends on` value that is entirely enclosed in one pair of `[` `]`, or in one pair of single `_` or `*`, MUST be treated as not given, with `SS_FIELD_PLACEHOLDER` (info) at the field line — unless, for `Depends on`, the enclosed text is itself a valid ID list or `none`, in which case it stays `SS_DEPENDENCY_MALFORMED`.
- **FR-003**: A placeholder MUST never satisfy a stage requirement, a dependency, or a spec link; stage findings apply as if the field were absent.

**016 — defect fixes**

- **FR-004**: A `Depends on` list that names an ID more than once MUST yield that ID once (first position) in `dependencies` and one `SS_DEPENDENCY_DUPLICATE` (warning) per repeated ID at the field line; dependency findings for the repeat MUST NOT be emitted twice.
- **FR-005**: The human summary panel MUST show a `statuses` line with every raw status and its count, in JSON key order; MUST show `invalid supersession N` when N > 0; each row whose supersession is invalid MUST mark its `superseded-by` line `(invalid)`; a row whose tasks file is shared with an earlier item MUST say so and that it is counted once.
- **FR-006**: `--json=false` MUST mean text mode for every outcome, including invocation errors.
- **FR-007**: An empty `--root` or `--roadmap` value MUST be `SS_INVOCATION_INVALID`, exit 2, in both `spec status` and (for `--root`) `spec inventory`.
- **FR-008**: Status invocation-error JSON MUST carry the requested `--fail-on` when its value is valid, else `"error"`.

**Both**

- **FR-009**: JSON MUST stay schema 1.0; the two new codes (`SS_FIELD_PLACEHOLDER`, `SS_DEPENDENCY_DUPLICATE`) are additive. JSON MUST be byte-identical for inputs without any of rows A–D and invocations without F–H.
- **FR-010**: The diagnostics, JSON, CLI and grammar contracts, README (Limits, findings, example output) and the extension's `reference/output.md` MUST be updated; the README Limits entries describing rows A–H MUST be removed.
- **FR-011**: No runtime side effects or new dependencies (Art. II); unsupported input still fails visibly (Art. III).

### Key Entities

- **Placeholder value**: a `Spec dir` or `Depends on` value enclosed entirely in brackets or single emphasis, not itself a valid ID list.
- **Duplicate dependency**: a repeated ID within one item's `Depends on` list.
- **Display additions**: status counts line, invalid supersession count and row mark, shared-tasks attribution.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every row A–H in the Context table gives its new behavior, checked by a test per row.
- **SC-002**: A roadmap copied from the pinned speckit-roadmap template with its placeholder lines left in exits 0 at the default threshold.
- **SC-003**: JSON output for every existing fixture root is byte-identical to the previous build (none contains rows A–D); human text differs only by the new `statuses` line (and, where applicable, the other display additions), and the README example is updated to match.
- **SC-004**: No real dependency ID is ever dropped: every fixture where an enclosed value contains a valid ID list still reports `SS_DEPENDENCY_MALFORMED`.
- **SC-005**: Every README Limits entry describing rows A–H is removed, and no new Limits entry is needed for these rows.

## Assumptions

- **Placeholders get their own info code** so a template line never vanishes silently (C-01); info does not fail CI at the default threshold.
- **The `statuses` line is always shown**, so every human report gains one panel line and the README example changes; JSON is unaffected. The owner can ask for it only when there is more than one status.
- **Spellings**: only `Superseded by` / `superseded by` are added, because review D6 found that form; widening to all case forms would weaken C-10's visible-failure rule.
- **Deduplicating** a repeated dependency changes the `dependencies` array only for inputs that have a repeat (none of the fixtures), with a warning, so it is not silent.
- **Out of scope**: other README Limits entries (prose labels, indented forms, annotated dependencies Q5, the FIFO-as-absent and over-long `Spec dir` defects), and anything in item 011 (README sample test), except keeping the example correct.
