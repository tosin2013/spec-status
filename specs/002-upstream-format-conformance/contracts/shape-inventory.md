# Shape inventory

Each row is one shape key produced by the classifier in `internal/conformance/shapes_test.go`. `TestPinnedInventory` checks this table:

- every key the pinned templates use must have a row here;
- every key the conformance fixtures use must have a row here;
- every `supported` key must be used by at least one fixture.

**Status**:
- `supported`: produced by a pinned template. A wrong result on it is blocking under D-2026-11-07-09.
- `hand-edit`: not in a template, but a common edit to generated files. Also supported, and covered by fixtures.
- `extension`: a Spec Status field that upstream doesn't use (FS-04).

Anything not listed is out of dialect. It fails visibly or is a README limit (`specs/001-spec-status/roadmap-format.md` §2–3).

Pins: github/spec-kit v1.2.0 (commit `4ca46a07`) and srobroek/speckit-roadmap v0.1.0. The template copies are in `testdata/upstream/`.

## Roadmap shapes (srobroek/speckit-roadmap v0.1.0, `templates/roadmap-template.md`)

| Key | Example | Interpretation | Status | Fixtures |
|---|---|---|---|---|
| `html-comment` | `<!--` SYNC IMPACT REPORT … `-->`, and section comments | Not live data | supported | template-filled, this-repo, speckit-roadmap-upstream |
| `heading:1` | `# [PROJECT_NAME] — Spec Roadmap` | Prose | supported | template-filled, hand-edits |
| `heading:2` | `## Planned Specs` | Ends any open item | supported | template-filled, hand-edits, this-repo |
| `heading:3` | `### Core` (grouping heading) | Ends any open item | supported | template-filled, this-repo, speckit-roadmap-upstream |
| `item-heading` | `### 001 — [Spec Title]  [status: planned]` | Item: ID, title, raw status | supported | template-filled, hand-edits, this-repo, upstream-generated-002-tasks, speckit-roadmap-upstream |
| `field-style:bold` | `- **Description:** …` | Bold labels read like plain ones | supported | template-filled, hand-edits, this-repo |
| `field:Description` | `- **Description:** …` | Item text | supported | template-filled, hand-edits |
| `field:Outcome` | `- **Outcome:** …` | Item text | supported | template-filled, hand-edits |
| `field:Scope (in)` | `- **Scope (in):** …` | Fills Scope | supported | template-filled, hand-edits |
| `field:Scope (out)` | `- **Scope (out):** …` | Descriptive text | supported | template-filled |
| `field:Depends on` | `- **Depends on:** 004 (verified).` | Dependency list; `none`, a trailing period, and one `(note)` per ID allowed | supported | template-filled, hand-edits, this-repo, speckit-roadmap-upstream |
| `field:Governed by` | `- **Governed by:** C-01` | Descriptive text, no references | supported | template-filled, this-repo |
| `field:Addresses` | `- **Addresses:** PRD goal 1` | Descriptive text | supported | template-filled |
| `field:Spec dir` | `- **Spec dir:** specs/001-sample/` | Spec directory; trailing `/` and `.` allowed | supported | template-filled, hand-edits, this-repo, speckit-roadmap-upstream |
| `field:Notes` | `- **Notes:** …` | Descriptive text | supported | template-filled, hand-edits |
| `field-continuation` | The template's `Outcome` value wrapped onto an indented line | Joined to the field value | supported | template-filled, hand-edits |
| `bullet` | `- **C-01 — …:** …` under `## Constraints & Decisions` | Prose outside items | supported | template-filled, this-repo |
| `thematic-break` | `---` before the sections and before the footer | Ends any open item | supported | template-filled |
| `field-style:plain` | `- Spec dir: specs/002-sample` | Plain labels (feature 001 dialect) | hand-edit | upstream-generated-002-tasks |
| `field:superseded-by` | `- superseded-by: 002` | Successor reference (FS-04) | extension | covered by `internal/status` tests |
| `block-quote` | `> …` prose between sections | Prose outside items (under an item it is fatal) | hand-edit | speckit-roadmap-upstream |

The status vocabulary is the nine tokens in the template legend. Unknown tokens are kept raw with a warning.

## Tasks shapes (github/spec-kit v1.2.0, `core_pack/templates/tasks-template.md`)

| Key | Example | Interpretation | Status | Fixtures |
|---|---|---|---|---|
| `front-matter` | `---` / `description: "…"` / `---` | Not tasks | supported | template-filled, upstream-generated-002-tasks |
| `heading:1` | `# Tasks: [FEATURE NAME]` | Closes open list items | supported | template-filled, hand-edits |
| `heading:2` | `## Phase 1: Setup (Shared Infrastructure)` | Closes open list items | supported | template-filled, upstream-generated-002-tasks |
| `heading:3` | `### Tests for User Story 1 …` | Closes open list items | supported | template-filled, upstream-generated-002-tasks |
| `task-line` | `- [ ] T001 Create project structure …` | Counted, unchecked | supported | template-filled, hand-edits, upstream-generated-002-tasks |
| `task-id:T###` | `T001` at the start of the task text | Part of the text | supported | template-filled, hand-edits, upstream-generated-002-tasks, speckit-roadmap-upstream |
| `task-marker:[P]` | `- [ ] T003 [P] Configure linting …` (bracket groups right after the task ID are markers) | Part of the text, not a link | supported | template-filled, upstream-generated-002-tasks |
| `task-marker:[USn]` | `- [ ] T012 [P] [US1] Create …` | Part of the text, not a link | supported | template-filled, upstream-generated-002-tasks |
| `task-text-brackets` | `[Entity1]`, `[endpoint]` in task text | Text, not links (no reference definitions) | supported | upstream-generated-002-tasks |
| `html-comment` | The `<!--` … `-->` block explaining the sample tasks | Not tasks | supported | template-filled |
| `block-quote` | `> **NOTE: Write these tests FIRST, …**` | Prose | supported | template-filled |
| `fence` | ```` ```bash ```` … ```` ``` ```` | Not tasks | supported | template-filled, hand-edits |
| `ordered-list` | `1. Complete Phase 1: Setup` | Prose (no checkbox) | supported | template-filled |
| `bullet` | `- **[P]**: Can run in parallel …` | Prose | supported | template-filled |
| `thematic-break` | `---` between phases | Closes open list items | supported | template-filled |
| `task-checked` | `- [x] T001 …` (a completed task) | Counted, checked | hand-edit | template-filled, hand-edits, speckit-roadmap-upstream |
| `task-nested` | `  - [x] T001a sub-step` | Counted | hand-edit | hand-edits |
| `html-block` | `<details>` … blank line … tasks … blank line … `</details>` | Tasks inside after a blank line are counted | hand-edit | hand-edits |
| `table` | Tables in hand-written tasks files (for example `specs/001-spec-status/tasks.md`) | Prose; rows carrying `[status:` or an item ID are fatal in roadmaps | hand-edit | this-repo |

## spec.md and plan.md (github/spec-kit v1.2.0)

spec status only checks that these files exist (stage matrix, D-2026-11-07-03). It doesn't parse their contents, so they have no shapes.
