# Phases fixture (feature 011)

Hand-counted from the files on 2026-11-21, not from `spec` output. It checks `tasks.phases` in `spec status --json`.

| Item | Tasks file | `tasks.state` | Expected `tasks.phases` |
|---|---|---|---|
| 001 | `specs/001-phases/tasks.md`: 8 tasks, 4 checked | read | `Phase 1: Setup` line 5, story null, 1/2. `Phase 2: Foundational` line 19, null, 0/0. `Phase 3: User Story 1 - Read one phase (Priority: P1)` line 21, story 1, 2/3. `Phase 4: User Story 1234567890 - too long a number` line 32, null (ten digits), 0/1 |
| 002 | `specs/002-no-phases/tasks.md`: 2 tasks, 1 checked | read | `[]` |
| 003 | `specs/003-missing-tasks/tasks.md` does not exist | missing | `null` |
| 004 | none (no Spec dir) | not_linked | `null` |
| 005 | shares 001's file | read | same as 001 |

**What 001 contains:**
- **Not in any phase:** T000, before the first phase, and T007, under `## Dependencies & Execution Order`. Both count only in the item totals.
- **Ignored:** the fenced `## Phase 9` (lines 10–13) and the commented `## Phase 8` (lines 15–17). They open no phase, and the fenced checkbox is not a task.
- **Level-3 headings** (lines 23 and 27) stay inside Phase 3.

The diagnostics (`SS_STAGE_FILE_MISSING` for 003, `SS_SPEC_REFERENCE_MISSING` for 004) are expected and unrelated to phases.
