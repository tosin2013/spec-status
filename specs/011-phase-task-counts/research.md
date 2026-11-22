# Research: Per-phase task counts

## R1 — Where to recognize phase headings

- **Decision:** Inside `parseTasks`, in the existing branch that treats an ATX heading at width ≤ 3 as a heading (`heading(line)` in `internal/status/strict.go`, branch around `parse.go:714`).
- **Rationale:** By that point, lines inside fences, HTML comments and raw HTML blocks have already been skipped with `continue`. A `>` line starts with `>`, so `heading()` never sees it as a heading, and indented code (width > 3) is rejected by `heading()` itself. Phase recognition therefore inherits exactly the rules that decide which checkboxes count (FR-005), without a second scanner.
- **Alternatives considered:** A separate pass over the file with its own fence tracking, rejected because the two scanners could disagree. Deriving phases in `analyze.go` from line numbers, rejected because the heading lines are not kept after parsing.

## R2 — Heading text and qualification

- **Decision:** Take `text` from `heading()`. Remove a closing sequence of `#` characters when it is preceded by a space or tab, or when the whole text is `#`s (CommonMark closing sequence), then trim. A level-2 heading qualifies when the text equals `Phase` or starts with `Phase` followed by an ASCII space, a tab or `:`. The match is case-sensitive.
- **Rationale:** This matches every Spec Kit heading (`Phase 1: Setup (Shared Infrastructure)`, `Phase N: Polish & Cross-Cutting Concerns`) and rejects prose (`Phases overview`, `phase 1`). Stored text is never rewritten beyond the CommonMark marker removal, so the `🎯 MVP` suffix is kept.
- **Alternatives considered:** Matching through the `lookalikes` replacer, rejected because phases do not change counts or relationships, so a near-miss heading only loses a phase entry, never a task. Matching case-insensitively, rejected because it is not upstream's form.

## R3 — Where a phase ends

- **Decision:** Any recognized ATX heading of level 1 or 2 closes the open phase. A qualifying level-2 heading then opens the next one. Level 3–6 headings do not close it.
- **Rationale:** Spec Kit's `### Tests for User Story N` and `### Implementation for User Story N` sit inside a phase, and `## Dependencies & Execution Order` ends the last one. A level-1 heading is a higher section, so it also ends the phase (spec edge case). Setext headings are not recognized by `heading()` and are a documented limit.
- **Alternatives considered:** Ending phases only at level 2 (the wording of #66), extended to level 1 because a phase cannot sensibly continue past a top-level heading.

## R4 — Counting

- **Decision:** At the point where `parseTasks` appends a `TaskCheck` (`parse.go:847`), also increment the open phase's `Total` and, if checked, its `Checked`.
- **Rationale:** A checkbox is in a phase exactly when it is counted in the item totals (FR-004). The sums are consistent by construction.
- **Alternatives considered:** Bucketing `Checks` by line after parsing works too, but needs the heading lines kept separately. Counting inline keeps one source of truth.

## R5 — JSON rendering and null rule

- **Decision:** `itemTaskView` gets `Phases *[]phaseView` with `json:"phases"`, placed after `state`. `itemToView` sets it to a non-nil (possibly empty) slice only when `State == "read"`, otherwise `nil`, which renders as `null`. `phaseView` fields are in this order: `heading`, `line`, `story` (`*int`), `checked`, `total`.
- **Rationale:** Field order is deterministic (Art. V determinism). `null` versus `[]` follows the null-denominator rule (spec US2).
- **Alternatives considered:** Omitting the field when the file is unread (`omitempty`), rejected because consumers would need two code paths and the "field present" check would be ambiguous.

## R6 — Conformance and fixtures

- **Decision:**
  - `expectedTasks` in `internal/conformance/expected_test.go` gains an optional `phases` (pointer to a slice, so absent means "not checked" and `[]` is checked). `compareItem` compares heading, line, story, checked and total in order.
  - Phases are hand-counted for every conformance fixture with a read tasks file. The primary evidence is `upstream-generated-002-tasks`, the unedited output of `/speckit-tasks`: seven phases, phase 3–6 stories 1–4, plus a fenced `#` line after the phases.
  - A new synthetic fixture `testdata/phases/` covers the edge cases: no phases, a fenced `## Phase`, tasks before the first phase and after `## Dependencies`, an empty phase, level-1 ending, closing `#`s, `Phases overview`, a 10-digit story number, duplicate headings, a missing tasks file, and a tasks file shared by two items.
- **Rationale:** Hand counts are the acceptance bar (D-2026-11-07-09). `expected.json` is never regenerated from `spec` output.
- **Alternatives considered:** Adding phases only to unit tests, rejected because FR-009 requires the conformance record.
