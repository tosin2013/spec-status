# Shape inventory (draft, from the pinned templates)

Status key: **supported** means the shape is produced by a pinned template, and a wrong result on it is blocking under D-2026-11-07-09. **strict** means no pinned template produces it; spec status fails visibly or the case is a README limit.

This draft was taken from the pinned template files on 2026-11-07. Implementation turns it into fixtures and checks it with `TestPinnedInventory`.

## srobroek/speckit-roadmap v0.1.0: `roadmap.md` (from `templates/roadmap-template.md`)

| Shape key | Example from the template | Interpretation | Status |
|---|---|---|---|
| `html-comment` (top level, multi-line) | `<!--` … `SYNC IMPACT REPORT` … `-->` | Not live data | supported |
| `html-comment` (under a `##` section) | `<!-- THE LEDGER. … -->` | Not live data | supported |
| `heading:1` | `# [PROJECT_NAME] — Spec Roadmap` | Prose | supported |
| `heading:2` | `## Vision & End States`, `## Planned Specs`, `## Open Questions` | Ends any open item | supported |
| `heading:3` (grouping) | `### Core extension` (named in the template's comment) | Ends any open item | supported |
| `item-heading` | `### 001 — [Spec Title]  [status: planned]` (two spaces before `[status:`) | Item: ID, title, raw status | supported |
| `field:Description` / `Outcome` / `Scope (in)` (bold) | `- **Description:** …` | Item text; `Scope (in)` fills Scope | supported |
| `field:Depends on` (bold) | `- **Depends on:** [NNN, NNN — or "none"]` | Dependency list; `none`, a trailing period, and `(note)` per ID allowed | supported |
| `field:Spec dir` (bold) | `- **Spec dir:** [specs/001-.../ — …]` | Spec directory; trailing `/` and `.` allowed | supported |
| `field:Scope (out)` / `Governed by` / `Addresses` / `Notes` (bold) | `- **Notes:** …` | Descriptive text, no references | supported |
| `field-continuation` | Wrapped `Outcome` value on an indented next line | Joined to the field value | supported |
| `bullet` (outside items) | `- **C-01 — [constraint]:** …` under `## Constraints & Decisions` | Prose | supported |
| `thematic-break` | `---` between the intro and the sections, and before the footer | Ends any open item | supported |
| `footer` | `**Version**: … \| **Ratified**: … \| **Last Amended**: …` | Prose | supported |
| `field:superseded-by` | not in the upstream template | Spec Status extension (feature 001, FS-04) | supported (local extension) |

The status vocabulary is the nine tokens listed in the template legend; unknown tokens are kept raw with a warning.

## github/spec-kit v1.2.0 (commit 4ca46a07): `tasks.md` (from `core_pack/templates/tasks-template.md`)

| Shape key | Example from the template | Interpretation | Status |
|---|---|---|---|
| `front-matter` | `---` / `description: "…"` / `---` at the top | Not tasks | supported |
| `heading:1` / `heading:2` / `heading:3` | `# Tasks: …`, `## Phase 1: Setup …`, `### Tests for User Story 1 …` | Close open list items | supported |
| `task-line` | `- [ ] T001 Create project structure …` | Counted task | supported |
| `task-marker:[P]` | `- [ ] T003 [P] Configure linting …` | Part of the task text; not a link | supported |
| `task-marker:[USn]` | `- [ ] T012 [P] [US1] Create [Entity1] model in src/models/[entity1].py` | Part of the task text; not a link | supported |
| `task-text-brackets` | `[Entity1]`, `[endpoint]`, `[name]` inside descriptions | Text, not links (no reference definitions exist) | supported |
| `html-comment` (top level) | `<!--` … `-->` block explaining the sample tasks | Not tasks | supported |
| `block-quote` (prose) | `> **NOTE: Write these tests FIRST, …**` | Prose | supported |
| `fence` (top level) | ```` ```bash ```` … ```` ``` ```` | Not tasks | supported |
| `ordered-list` (prose) | `1. Complete Phase 1: Setup` | Not tasks (no checkbox) | supported |
| `bullet` (prose) | `- **[P]**: Can run in parallel …` | Not tasks | supported |
| `task-checked` | `- [x] T001 …` (produced when tasks are completed) | Counted, checked | supported |

## github/spec-kit v1.2.0: `spec.md`, `plan.md`

spec status only checks that these files exist (the stage matrix from D-2026-11-07-03). Their contents are not parsed, so no shapes are listed.

## Hand edits treated as supported (FR-006)

| Shape key | Example | Interpretation |
|---|---|---|
| `field-continuation` | A `Depends on` list wrapped onto an indented next line | Joined to the value |
| `task-nested` | `  - [ ] T001a sub-step` under a task | Counted |
| `task-details` | An indented paragraph under a task, then nested subtasks | Item stays open; nested subtasks counted |
| `task-fenced-example` | An indented fence under a task | Not tasks |
| `details-with-blank` | `<details>` / `<summary>…</summary>` / blank line / tasks / blank line / `</details>` | Tasks counted |

## Strict (not produced by any pinned template)

Everything else in `specs/001-spec-status/roadmap-format.md` §2–3 that fails visibly. Examples: item headings at other levels or without the em dash; near-miss field labels; free text under an item; tables with item rows; fences and comments that cross list-item boundaries; block-quoted checkboxes; checkboxes four or more columns past their item; raw `<pre>`/`<script>`/`<style>`/`<textarea>` blocks; empty comments. These stay fatal (FR-014).
