# Feature Specification: Per-phase task counts

**Feature Branch**: `011-phase-task-counts`

**Created**: 2026-11-21

**Status**: Draft

**Input**: Roadmap item 017, "Per-phase task counts", issue #66, chosen for v0.4.0 by the owner in decision D-2026-11-21-01 (1): "each item's `tasks` object in `spec status --json` gains an additive `phases` array (heading, line, story, checked, total) for level-2 `## Phase …` sections of tasks.md; JSON schema stays 1.0."

## Context

Spec Kit's generated `tasks.md` is divided into phases: Setup, Foundational, one phase per user story, and Polish. Each phase is a level-2 heading that starts with `Phase`, and its checkbox tasks sit under it. Today `spec status --json` reports only one checked/total count for the whole file.

specfoundry (issue #66) builds one *slice* of a roadmap item per run, usually one user-story phase. Its gate refuses a slice whose tasks are already ticked. Its grader checks that only that slice's tasks were ticked. Its rule is that only the roadmap and `spec status` JSON decide what may run, so it must not parse tasks.md itself. It needs per-phase counts from `spec status`.

This feature reports recorded facts only (C-01). A phase's counts say how many boxes are checked under a heading. They do not say that the phase is done, verified or ready.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - A gate reads one phase's progress from the JSON (Priority: P1)

An automated gate runs `spec status --json` on a project whose item has a Spec Kit tasks.md. It finds the item and then the phase for "User Story 2". It reads that phase's checked and total counts and refuses to start the slice if any task in it is already checked.

**Why this priority**: This is the request in issue #66. Without it, the consumer must parse tasks.md itself, which its rules forbid.

**Independent Test**: Run `spec status --json` on a fixture whose tasks.md has the generated layout. Each phase's heading, line, story number and counts match a hand count of the file.

**Acceptance Scenarios**:

1. **Given** a tasks.md with `## Phase 1: Setup`, `## Phase 2: Foundational`, `## Phase 3: User Story 1 …`, `## Phase 4: User Story 2 …`, `## Phase 5: Polish …` and then `## Dependencies & Execution Order`, **When** the report is produced, **Then** the item's `tasks.phases` lists exactly five phases in file order, each with its heading text, its 1-based heading line, its counts, and `story` 1 and 2 for the two user-story phases and `null` for the others.
2. **Given** the same file, **When** the counts are compared, **Then** the per-phase checked and total counts add up to the item's checked and total counts for the tasks that sit inside phases.
3. **Given** a phase with a level-3 heading inside it (for example `### Tests for User Story 1`), **When** the report is produced, **Then** the tasks under the level-3 heading count toward the enclosing phase, and no new phase starts.

---

### User Story 2 - "No phases" and "unknown" stay different (Priority: P1)

A consumer must tell apart a tasks file that has no phase headings from an item whose tasks file was not read at all. Treating "not read" as "no phases" would let a gate believe there is nothing ticked.

**Why this priority**: It is the same truthfulness rule as null denominators (Art. I). A wrong reading here leads the consumer to a wrong decision.

**Independent Test**: Run the report on three items: one whose tasks.md has no phase headings, one whose tasks file is missing, and one with no linked tasks file. Check that `phases` is `[]`, `null` and `null`.

**Acceptance Scenarios**:

1. **Given** a read tasks.md with checkboxes and no phase headings, **When** the report is produced, **Then** `tasks.phases` is `[]`, and the item's counts are unchanged from today.
2. **Given** an item whose tasks file state is `missing`, `skipped` or `not_linked`, **When** the report is produced, **Then** `tasks.phases` is `null`.
3. **Given** a read tasks.md with a phase heading and no checkboxes under it, **When** the report is produced, **Then** that phase is listed with checked 0 and total 0.

---

### User Story 3 - Nothing else changes (Priority: P2)

Existing users of `spec status` (human text, CI exit codes, and scripts that read schema 1.0) see no change apart from the new field.

**Why this priority**: The JSON schema stays 1.0 and the change must be additive (C-05).

**Independent Test**: Compare the output for every existing fixture root with the previous build. Human text and exit codes are identical, and the JSON differs only by the new `phases` field.

**Acceptance Scenarios**:

1. **Given** any existing fixture, **When** human text is produced, with or without filters and `--explain`, **Then** it is byte-identical to the previous build.
2. **Given** any existing fixture, **When** JSON is produced, **Then** removing every `tasks.phases` field gives the previous build's bytes, and `schema_version` is still `"1.0"`.

### Edge Cases

- **Tasks before the first phase heading**, or after a later non-phase level-2 heading such as `## Dependencies`: they count toward the item totals only and belong to no phase.
- **Heading forms:**
  - An ATX level-2 heading (`## `, with up to three spaces of indentation) whose text starts with the word `Phase` opens a phase.
  - The word must be exactly `Phase` and be followed by whitespace, a colon or the end of the text. `## Phase 1: Setup`, `## Phase N: Polish` and `## Phase` count. `## Phases overview` and `## phase 1` do not.
  - Optional closing `#` characters are removed from the reported heading text, as are leading and trailing spaces.
- **Ending a phase:** a phase runs until the next ATX heading of level 1 or 2, or the end of the file. Level 3 to 6 headings stay inside it.
- **Hidden headings:** an ATX heading inside a fence, an HTML comment, indented code, a raw HTML block or a block quote is not a heading for this purpose. This follows the same rules that already decide which checkboxes count.
- **Setext headings** (text underlined with `---` or `===`) neither open nor end a phase. Spec Kit never generates them; this is listed as a limit.
- **Story number:**
  - `story` is the number from the first `User Story N` in the heading text (case-sensitive, N one or more decimal digits), otherwise `null`.
  - N with more than nine digits gives `null`, not a truncated number.
- **Duplicate headings:** two phases with the same heading text are both listed, each with its own line.
- **Shared tasks file:** two items that link the same tasks file each get the same phases, as they already get the same per-item counts.
- **Fatal tasks files:** when a tasks file is fatal under the existing rules, the report is fatal as today, and there are no phases to show.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: For each item, `spec status --json` MUST include `tasks.phases`.
  - It is an array when the tasks file state is `read`, and `null` for every other state.
- **FR-002**: Each phase entry MUST contain exactly these fields, in this order:
  - `heading`: the heading text, without the `##` marker, closing `#`s or surrounding spaces;
  - `line`: the 1-based line of the heading in the tasks file;
  - `story`: an integer, or `null`;
  - `checked` and `total`: integers.
- **FR-003**: Phases MUST be listed in file order. A phase MUST start at a qualifying level-2 ATX heading and run to the next level-1 or level-2 ATX heading, or to the end of the file.
- **FR-004**: A checkbox MUST be counted in a phase exactly when it is counted in the item totals today and sits between that phase's heading and the phase's end.
  - No checkbox is counted in two phases.
  - Checkboxes outside every phase are counted only in the item totals.
- **FR-005**: Heading recognition MUST ignore headings inside fences, HTML comments, indented code, raw HTML blocks and block quotes, consistent with the existing tasks dialect.
  - This feature introduces no new fatal error.
- **FR-006**: `story` MUST be taken from the first `User Story N` in the heading (N: 1 to 9 decimal digits), otherwise `null`.
- **FR-007**: Human text, the summary, diagnostics, exit codes and all other JSON fields MUST be unchanged.
  - `schema_version` stays `"1.0"`.
- **FR-008**: The JSON contract (`specs/001-spec-status/json-contract.md`) and the tasks grammar (`specs/001-spec-status/roadmap-format.md`) MUST document the field, its null/empty rule and the phase grammar.
  - The README limits MUST list setext headings.
- **FR-009**: Hand-counted conformance fixtures (`testdata/conformance/*/expected.json`) MUST be updated by hand with the expected `phases` values, never regenerated from `spec` output.
- **FR-010**: A fixture MUST cover Spec Kit's generated tasks layout (Setup, Foundational, two user stories with level-3 subsections, Polish, then a non-phase section with checkboxes). Another MUST cover a file without phases, a fenced `## Phase` line, and tasks before the first phase.

### Key Entities

- **Phase:** a section of one tasks file, with heading text, heading line, an optional user-story number, and checked/total counts of the tasks inside it. It belongs to the item's tasks object and says nothing about status or verification.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For the generated-layout fixture, every phase's heading, line, story, checked and total match a hand count. The sum over phases equals the in-phase share of the item totals.
- **SC-002**: For every existing fixture root, human text and exit codes are byte-identical to the previous build. JSON is identical once the `phases` fields are removed.
- **SC-003**: For items with an unread tasks file, `phases` is `null` in 100% of cases. For read files without phase headings, it is `[]` in 100% of cases.
- **SC-004**: A consumer can decide whether a user-story phase has any checked task from one `spec status --json` run, without reading tasks.md.

## Assumptions

- Spec Kit v1.2.0's tasks template (`testdata/upstream/spec-kit/v1.2.0/templates/tasks-template.md`) is the reference layout. Its phases are ATX level-2 headings, and its later sections (`## Dependencies & Execution Order`, `## Parallel Example …`, `## Implementation Strategy`, `## Notes`) end the last phase.
- Phase names other than "User Story N" (Setup, Foundational, Polish) carry no story number, and consumers identify them by heading text.
- The field is additive under the schema 1.0 rules, so no schema version change is needed. Consumers that ignore unknown fields are unaffected.
- The public Go package (roadmap 018) will expose this field later. This feature changes only the CLI's JSON.
- Not in scope: phase output in human text, per-phase diagnostics, and any inference of status from phases.
