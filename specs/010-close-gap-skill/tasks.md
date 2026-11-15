---
description: "Tasks for feature 010: close-gap skill"
---

# Tasks: Close-gap skill

**Input**: `specs/010-close-gap-skill/` (spec.md, plan.md, research.md, data-model.md, contracts/close-gap-command.md, quickstart.md). Roadmap item 013, D-2026-11-14-01 (1), v0.3.0; depends on 012 and 014.

**Tests**: Required. The extension tests (structure, links, roles, manifest, fixture signals) are updated first and must fail until the command and reference changes exist. The all-causes fixture's `EXPECTED.md` is written by hand from contracts/close-gap-command.md "Fixture", never generated from `spec` output. Agent behavior is checked by a recorded walkthrough, as in feature 005.

**Organization**: US1 (breakdown) and US2 (one change at a time) are both in the one command file; US3 (pm review) and US4 (six skills, 0.3.0) are reference and plumbing changes.

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup

- [ ] T001 Branch from up-to-date `main` after PR #60 merges (or from its branch, opening the implementation PR against `main`); set `.specify/feature.json` to `specs/010-close-gap-skill`; confirm `go test ./internal/extension` passes before any change

---

## Phase 2: Foundational (tests and fixture first)

- [ ] T002 Add `"speckit.spec-status.close-gap"` to `extensionCommands` in `internal/extension/helpers_test.go` (after verify-readiness) and `"verification-gap"` to `reviewSections` in `internal/extension/roles_test.go`; run `go test ./internal/extension` and confirm it fails (missing command file, manifest entry, roles row)
- [ ] T003 [P] Create `testdata/skills/close-gap-all-causes/` exactly as in contracts/close-gap-command.md "Fixture" (items 001–008, their `Spec dir` folders with `spec.md`, `plan.md` and `tasks.md` checkboxes as listed, 006 `Depends on: 004`, 008 status `done`), and write `EXPECTED.md` by hand: Command `/speckit.spec-status.close-gap`; Signals `item 002 status=implemented tasks=2/2 state=read; item 003 status=planned tasks=2/2 state=read; item 004 status=implemented tasks=1/3 state=read; diag SS_TASKS_DONE_STATUS_EARLY item=003; diag SS_STATUS_AHEAD_OF_TASKS item=004; diag SS_STAGE_FILE_MISSING item=005; diag SS_STATUS_UNKNOWN item=008; dep 006->004 state=unfinished; summary active=8 verified=1 blocked=1 human_review=0; outcome=report`; Expected finding: `recorded verified 1/8`, the cause per item, 007 not started; Expected order of work `002, 003, 004, 005, 008, 006`; Expected proposal: 002 → `/speckit.roadmap.debrief`; Must not: change any file before the user approves the specific change

**Checkpoint**: tests fail only for the missing command, manifest entry and roles row; the fixture's Signals pass against real `spec`.

---

## Phase 3: User Story 1 — See why the recorded percentage is low (P1) 🎯 MVP

**Goal**: the read-only breakdown by cause with evidence and order of work.

**Independent Test**: on the fixture, an agent following the command lists every active non-verified item under the expected cause and changes nothing.

- [ ] T004 [US1] Create `extensions/spec-status/commands/speckit.spec-status.close-gap.md` with frontmatter (`description` naming low recorded-verified and one approved change at a time; `argument-hint: "[item ID to start from]"`), `## User Input` (`$ARGUMENTS`), `## Goal`, `## Operating Constraints` (read `.specify/extensions/spec-status/reference/rules.md` first; read-only until a specific change is approved; changes only through Spec Kit commands; `verified` is a label, not proof; never `verified` without a debrief with no must-address finding and the user's agreement; never redirect a superseded dependency; `spec inventory` is never status evidence), and Outline steps 1–4 from contracts/close-gap-command.md: `spec version`; `spec status --json` and `--color=never --explain` (exit 2: report and stop); classify active non-verified items by the precedence table in research R1 including the older-`spec` fallbacks; report `recorded verified X/Y`, the per-cause breakdown with `path:line`, items outside the gap, and the order of work (research R2)
- [ ] T005 [US1] Add `extensions/spec-status/commands/speckit.spec-status.close-gap.md` to `provides.commands` in `extensions/spec-status/extension.yml` (name, file, description matching the frontmatter); run `go test ./internal/extension`

---

## Phase 4: User Story 2 — Close the gap one approved change at a time (P1)

**Goal**: the proposal loop with approval, apply, re-read.

**Independent Test**: on the fixture, the first proposal is the debrief for 002; nothing is written before approval; declining moves to 003.

- [ ] T006 [US2] Add Outline steps 5–6 to `extensions/spec-status/commands/speckit.spec-status.close-gap.md`: take the first item in the order (or the item in `$ARGUMENTS`); propose one change per research R3's table with item, cause, evidence, the change, the command and what will be re-checked; wait for approval; on approval apply only through that command, re-run `spec status --json`, report `confirmed` or `not confirmed` and the new `recorded verified X/Y`; on skip move to the next item; on stop summarize what changed and what is left; blocked items only after their blocking dependency has been handled or skipped; human-review dependencies routed to `/speckit.spec-status.unblock`
- [ ] T007 [US2] Add `## Follow-ups` to the close-gap command listing `/speckit.spec-status.reconcile`, `/speckit.spec-status.verify-readiness`, `/speckit.spec-status.unblock`, `/speckit.roadmap.debrief`, `/speckit.roadmap.write`, `/speckit-specify`, `/speckit-plan`, `/speckit-tasks`, `/speckit-implement` (all in `allowedLinks`); run `go test ./internal/extension`

---

## Phase 5: User Story 3 — The pm review shows the gap and points to close-gap (P2)

**Goal**: `verification-gap` section; pm leads with it; close-gap first pm follow-up; review links to close-gap.

**Independent Test**: the roles test passes with the new section; the review reachability test passes.

- [ ] T008 [P] [US3] In `extensions/spec-status/reference/roles.md`, add the `verification-gap` row to "Review sections" (text from contracts/close-gap-command.md), change the `pm` row's "Leads with" to `verification-gap`, `progress`, `blockers`, `ready-to-start`, `human-review` and its "First follow-ups" to /speckit.spec-status.close-gap, /speckit.spec-status.unblock, /speckit.spec-status.reconcile
- [ ] T009 [P] [US3] In `extensions/spec-status/commands/speckit.spec-status.review.md`, add `/speckit.spec-status.close-gap` to `## Follow-ups` ("the recorded-verified share is low and the user wants to work through it"), and mention the `verification-gap` section where the outline lists sections
- [ ] T010 [US3] Run `CGO_ENABLED=0 go test ./internal/extension`; all pass

---

## Phase 6: User Story 4 — The extension ships six skills as version 0.3.0 (P2)

**Goal**: version 0.3.0 and "six skills" everywhere they are counted.

**Independent Test**: manifest test (version matches README heading) passes; install check expects six; release archive check expects six.

- [ ] T011 [US4] Set `version: 0.3.0` in `extensions/spec-status/extension.yml` and the heading `# spec-status extension v0.3.0` in `extensions/spec-status/README.md`
- [ ] T012 [P] [US4] In `scripts/check-extension-install.sh`, add `close-gap` to the `want=(review explain reconcile unblock verify-readiness)` list (the script checks `speckit-spec-status-<name>/SKILL.md` for each) and change the OK message to six skills including close-gap
- [ ] T013 [P] [US4] In `.github/workflows/release.yml`, change the archive check to `[ "$n" = 6 ]` with message "expected 6 command files", and the release-notes sentence to list six skills including close-gap
- [ ] T014 [US4] In `extensions/spec-status/README.md`: rename "The five skills" to "The six skills" (and its contents link), add a close-gap section (what it does, the causes and order, the approval loop, what it never does), add a walkthrough on `testdata/skills/close-gap-all-causes`, and update the pm role description; keep the guide's tone
- [ ] T015 [P] [US4] In the root `README.md`, change "adds five skills" to six and name close-gap in the list
- [ ] T016 [US4] Run `CGO_ENABLED=0 go test ./...`; all pass

---

## Phase 7: Polish

- [ ] T017 Run `bash scripts/check-extension-install.sh` with Spec Kit 1.2.0 into a throwaway project; record the result (six skills rendered) in the PR description
- [ ] T018 Agent walkthrough per quickstart.md on a scratch copy of `testdata/skills/close-gap-all-causes` with the extension installed: breakdown and order as expected, first proposal 002 debrief, decline moves to 003, stop summarizes, nothing written; then `review pm` leads with verification-gap. Record what matched and any differences in the PR description (executed by the implementer, not independently verified)
- [ ] T019 Add a row for feature 010 to the "Implementation record" table in `specs/001-spec-status/tasks.md`, and update the feature-005 row's skill count if it states five
- [ ] T020 Draft the v0.3.0 release-note entry in the PR description (close-gap skill; extension 0.3.0 with six skills; pm review's verification-gap section)
- [ ] T021 Doomer review on the change (prepare before editing, submit, answer every Quick lead, finish) before handing off the PR
- [ ] T022 Propose roadmap 013 `specced` → `implemented` through `/speckit.roadmap.write` after the owner approves; `verified` only after 014 is recorded verified and the v0.3.0 release checks pass

---

## Dependencies and execution order

- T001 → T002, T003 → US1 (T004 → T005) → US2 (T006 → T007, same file) → US3 (T008, T009 parallel → T010) → US4 (T011 → T012, T013, T015 parallel → T014 → T016) → Polish.
- T017 needs Spec Kit 1.2.0 installed locally.

## Parallel examples

- T002 and T003 together.
- T008 and T009 together.
- T012, T013 and T015 together.

## Implementation strategy

1. **MVP = US1**: the read-only breakdown answers the owner's original question ("why is recorded verified 2/12?") on its own.
2. US2 adds the guided loop; US3 makes it discoverable; US4 ships it.
3. One PR against `main`; release v0.3.0 only after this lands (the release workflow requires extension version = tag).
