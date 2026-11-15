# Feature Specification: Close-gap skill

**Feature Branch**: `011-013-specify`

**Created**: 2026-11-15

**Status**: Draft

**Input**: Roadmap item 013, "Close-gap skill", chosen for v0.3.0 in decision D-2026-11-14-01 (1): "a sixth extension command, `speckit.spec-status.close-gap`, that works through items whose recorded status lags their evidence (for example a low recorded-verified percentage), one item and one proposed change at a time, plus a verification-gap section in `review` for the `pm` role. Propose-only, per-change approval (C-06)." It came from the owner's question on 2026-11-14 about repositories reporting "recorded verified 2/12 (16.67%)".

## Context

A low recorded-verified percentage has several different causes, and each needs a different next step:

| Cause | How `spec status` shows it (v0.3.0) | Right next step |
|---|---|---|
| Work finished, roadmap not updated | `SS_TASKS_DONE_STATUS_EARLY` | confirm, then advance the status to `implemented` |
| Status claims more than the tasks show | `SS_STATUS_AHEAD_OF_TASKS` | decide: unfinished work or stale boxes |
| Missing spec, plan or tasks for the recorded stage | `SS_STAGE_FILE_MISSING`, `SS_SPEC_REFERENCE_MISSING`, `SS_SPEC_DIRECTORY_MISSING` | create the artifact, or correct the status |
| Implemented and consistent, never checked against the roadmap | status `implemented`, tasks all checked, dependencies satisfied, no error | debrief, then record `verified` only if the user agrees |
| Waiting on an unverified dependency | dependency `state` `unfinished`, `SS_DEPENDENCY_BLOCKED` | close the gap on the dependency first |
| Planned work not started | `planned`, no spec | not a gap to close; report only |

The five existing skills each handle one of these. Today a person has to know which to run in which order. Close-gap is the guided path through them for the whole roadmap, in an order that makes the recorded percentage move honestly.

The skill's two rules mirror the constitution: verified is a label, not proof (Art. I), and nothing changes without explicit approval of each change (Art. VIII, C-06). Inventory output (feature 006) is not roadmap evidence and is never used to change a status.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See why the recorded percentage is low (Priority: P1)

A project manager runs close-gap on a roadmap with 2 of 12 active items recorded verified. Before anything changes, they get a breakdown of the other 10 items by cause (the table above), with item IDs and file:line evidence, and the order in which close-gap proposes to work through them.

**Why this priority**: understanding comes before any change, and this part is read-only.

**Independent Test**: on a fixture with one item per cause, the breakdown lists each item under the right cause with its evidence, and changes nothing.

**Acceptance Scenarios**:

1. **Given** a roadmap with active items in every cause, **When** close-gap runs, **Then** it reports `recorded verified X/Y` exactly as `spec status` does, then each non-verified active item under one cause, citing item ID, `path:line` and the code or field that shows the cause.
2. **Given** the breakdown, **When** it is shown, **Then** it states the order of work: items ready for debrief first, then tasks-done-early items, then ahead-of-tasks items, then missing artifacts, then blocked items in dependency order (upstream first), and notes planned items as not part of the gap.
3. **Given** a project with no roadmap (`SS_ROADMAP_UNAVAILABLE`) or an unrecognized one, **When** close-gap runs, **Then** it reports the diagnostic and stops; it may mention `spec inventory`, but never treats inventory counts as evidence for a status.

---

### User Story 2 - Close the gap one approved change at a time (Priority: P1)

The project manager agrees to start. Close-gap takes the first item, proposes exactly one change with its evidence and the Spec Kit command that would apply it, and waits. Only after the user approves that change does it run the command, re-run `spec status` to confirm the result, and move to the next item.

**Why this priority**: this is the "loop" the owner asked for, and its safety rules are the core of the design.

**Independent Test**: on the fixture, the first proposal is for the first item in the stated order; nothing is written before approval; after a simulated approval the next proposal reflects the re-read status.

**Acceptance Scenarios**:

1. **Given** a tasks-done-early item, **When** close-gap proposes, **Then** it asks the user to confirm the work is really done and proposes advancing to `implemented` (never directly to `verified`) via `/speckit.roadmap.write`.
2. **Given** an implemented, consistent item, **When** close-gap proposes, **Then** it proposes `/speckit.roadmap.debrief` first, and proposes recording `verified` via `/speckit.roadmap.write` only after a debrief with no must-address findings and the user's explicit agreement.
3. **Given** an ahead-of-tasks item, **When** close-gap proposes, **Then** it asks which is true (unfinished work, or stale boxes) and proposes the matching command, without choosing for the user.
4. **Given** a blocked item, **When** close-gap reaches it, **Then** it works on the blocking dependency first, and never redirects a dependency on a superseded item.
5. **Given** any proposal, **When** the user declines or says stop, **Then** close-gap records nothing, moves on or stops as asked, and summarizes what changed and what is left.
6. **Given** an approved change was applied, **When** close-gap continues, **Then** it re-runs `spec status` and reports the new recorded-verified figure before the next proposal.

---

### User Story 3 - The pm review shows the gap and points to close-gap (Priority: P2)

A project manager runs the review skill in the `pm` role. It leads with a short verification-gap section (the counts per cause) and offers close-gap as the first follow-up.

**Why this priority**: it makes the new skill discoverable from the entry point people already use.

**Independent Test**: the roles reference lists the new section for `pm` and close-gap among its first follow-ups; the review skill names close-gap in its follow-ups.

**Acceptance Scenarios**:

1. **Given** `review` with role `pm`, **When** the section order is decided, **Then** `verification-gap` comes first, followed by the existing pm sections.
2. **Given** the pm view, **When** follow-ups are listed, **Then** `/speckit.spec-status.close-gap` is first.
3. **Given** other roles, **When** they run, **Then** their views are unchanged except that review's follow-up list names close-gap.

---

### User Story 4 - The extension ships six skills as version 0.3.0 (Priority: P2)

A user installs the extension from the v0.3.0 release archive and gets six skills, with the guide describing close-gap.

**Why this priority**: release plumbing; required for anyone to use the skill.

**Independent Test**: the install check renders six skills; the release workflow's archive check expects six command files and extension version 0.3.0.

**Acceptance Scenarios**:

1. **Given** the extension manifest, **When** it is read, **Then** it declares six commands and version `0.3.0`, still requires Spec Kit `>=1.2.0`, and still marks partners optional.
2. **Given** the install check, **When** it runs, **Then** it expects six skills including close-gap.
3. **Given** the release workflow's archive check, **When** a v0.3.0 tag is built, **Then** it expects six command files.

---

### Edge Cases

- **Older `spec`** (before v0.3.0, no status/tasks codes): close-gap derives the same causes from the counts and states, as reconcile does.
- **Nothing to close** (every active item verified, or only planned items left): close-gap says so and stops.
- **Supersession**: superseded items are history and never part of the gap; a dependency on one stays a human-review decision.
- **Unknown status tokens**: listed as a cause of their own (fix the status word via `/speckit.roadmap.write`), never mapped to a known status.
- **Roles**: close-gap takes an optional item ID to start from; it has no role argument.
- **Partners**: close-gap may mention an installed partner only through the existing roles reference, never as a requirement.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The extension MUST add a sixth command, `speckit.spec-status.close-gap`, with the same structure and checks as the other five (frontmatter, User Input, Goal, Operating Constraints pointing to `reference/rules.md`, Outline starting with `spec version`, Follow-ups).
- **FR-002**: Close-gap MUST first report, read-only, `recorded verified X/Y` and every active non-verified item under exactly one cause from the Context table, with item ID and `path:line` evidence, and the order of work.
- **FR-003**: Close-gap MUST then propose at most one change at a time, each with evidence and the exact Spec Kit command, and apply nothing before the user approves that specific change (Art. VIII, C-06).
- **FR-004**: After an applied change, close-gap MUST re-run `spec status` and report the result before the next proposal.
- **FR-005**: Close-gap MUST never propose `verified` without a prior `/speckit.roadmap.debrief` with no must-address findings and the user's explicit agreement, and only for an item already recorded `implemented`. Advancing an item to `implemented` (for example from `planned` when every task is checked) MUST follow the user's confirmation that the work is really done, as reconcile does; no proposal ever moves an item straight to `verified`.
- **FR-006**: Close-gap MUST never redirect a dependency on a superseded item, never treat `spec inventory` output as status evidence, and never edit files except through Spec Kit commands.
- **FR-007**: The roles reference MUST add a `verification-gap` review section, lead the `pm` role with it, and list close-gap first among pm follow-ups; the review command MUST name close-gap in its follow-ups.
- **FR-008**: The manifest MUST declare six commands at version `0.3.0`; the install check, release archive check, extension guide and README skill list MUST say six skills and describe close-gap.
- **FR-009**: Skills fixtures MUST include at least one roadmap with an item in every cause, with a hand-written expected breakdown, order of work and first proposal.
- **FR-010**: The `spec` binary MUST NOT change for this feature.

### Key Entities

- **Gap item**: an active, non-verified item with exactly one cause and its evidence.
- **Proposal**: one change, its evidence, the command that applies it, and its approval state (proposed, approved, applied and confirmed, declined).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: On the all-causes fixture, an agent following close-gap lists 100% of active non-verified items under the expected cause and proposes first the expected item and change.
- **SC-002**: In every fixture run, zero files change before an explicit approval.
- **SC-003**: Close-gap never proposes `verified` without a debrief step first.
- **SC-004**: A user can see why recorded-verified is low, by cause, in one command, without reading JSON.
- **SC-005**: The extension tests, the install check (six skills) and the release archive check all pass.

## Assumptions

- **Order of work** puts "ready for debrief" first because it is the only cause that can raise recorded-verified without new work; the owner can reorder in `/speckit-clarify`.
- **Planned items without a spec** are not part of the gap: there is no evidence to reconcile.
- **Version**: extension 0.3.0 ships with `spec` v0.3.0; the tool requirement stays `spec >=0.1.2` because close-gap falls back to counts on older builds.
- **Blocked state**: the roadmap shows 013 blocked until 014 is recorded verified (C-07). That does not block building it; 013 is not recorded verified before 014.
- **Out of scope**: any change to the `spec` binary; automated status changes; changes to partner extensions; a role argument for close-gap.
