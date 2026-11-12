# Feature Specification: Agent skills extension

**Feature Branch**: `005-agent-skills`

**Created**: 2026-11-12

**Status**: Draft

**Roadmap item**: 012 (target release v0.2.0)

**Input**: Owner decision D-2026-11-12-04:
- **Packaging:** a Spec Kit extension, id `spec-status`, tracked in this repository. Its commands render as agent skills.
- **First set:** review (the entry point), explain, reconcile, unblock, verify-readiness.
- **Links:** the skills link to each other and to the Spec Kit and roadmap commands.
- **Writes:** read-only by default. A change is proposed and applied only through Spec Kit commands, after the user approves each change.
- **The `spec` binary does not change.**

Owner decision D-2026-11-12-05 (same day): personas (project manager, developer, QA, architect) are **views of the review skill**, chosen with an argument. They are not separate skills, so the rules and routing exist in one place.

Owner decision D-2026-11-12-06 (same day): roles may recommend **partner Spec Kit extensions** when the project has them installed. Candidates are adrkit (ADRs) for architect, Assure (QA) for qa, and Bug Triage Workflow for developer. Partners are detected, never required or installed by this extension.

## Context

`spec status` reports what a roadmap and its task files record. People and agents then have to work out what to do next. That is hardest when the record disagrees with itself, for example:
- every task is checked but the item is still `planned`;
- an item is `verified` while its tasks are unchecked, or its task file is missing;
- an item is `verified` while one of its dependencies isn't.

Today the only guidance is two sample prompts (`specs/001-spec-status/examples/prompts/`) and the install-and-review prompt (`docs/agent-prompt.md`).

This feature turns that guidance into five agent skills. Each skill knows the `spec status` output vocabulary: every diagnostic code, the JSON schema 1.0 fields, dependency and task states, and status words. The skills hand off to each other and to the Spec Kit commands that actually change the roadmap, specs or tasks. They ship as a Spec Kit extension, like the roadmap extension this project already uses, so any project that uses Spec Kit can install them.

The `spec` binary stays read-only and gains no AI, network or write behavior (constitution Art. II and VIII). The skills are development tooling around it.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - One command tells me where the roadmap stands and what to do next (Priority: P1)

As a developer using an AI coding agent in a Spec Kit project, I run the review skill. I get a short, accurate summary of the roadmap, and each finding is assigned to the skill that can deal with it.

**Why this priority**: It's the entry point. Every other skill is reached from it.

**Independent Test**: In a project with a roadmap that has one blocker, one human-review dependency and one status/task mismatch, run the review skill. The summary matches `spec status` and cites item IDs and file:line. Each of the three findings names the right follow-up skill.

**Acceptance Scenarios**:

1. **Given** a Spec Kit project with a supported roadmap and `spec` installed, **When** the user runs the review skill, **Then** it summarizes these, each with item IDs and file:line taken from the output:
   - recorded-verified progress;
   - checked tasks;
   - active blockers;
   - human-review dependencies;
   - diagnostics grouped by severity;
   - items ready to start.
2. **Given** findings of different kinds, **When** the summary is shown, **Then** each finding names one follow-up: explain (any diagnostic), reconcile (a status/task or evidence mismatch), unblock (a blocker or human-review dependency), or verify-readiness (an item the user wants to mark `verified`).
3. **Given** `spec` is not installed, **When** any skill runs, **Then** it stops and points to the install steps instead of guessing at the roadmap.
4. **Given** `spec status` exits 2 (for example, the roadmap is missing), **When** the review skill runs, **Then** it reports the diagnostic and stops, and it does not create a roadmap.

---

### User Story 2 - Fix records that disagree with each other (Priority: P1)

As a developer, when the recorded status of an item disagrees with its tasks or evidence, the reconcile skill shows me each mismatch, says which side is likely wrong, and proposes one correction at a time. A correction is applied only through the right Spec Kit command and only after I approve it.

**Why this priority**: This is the problem that prompted the feature. Mismatches are the most common reason the report and reality diverge.

**Independent Test**: For each mismatch fixture under `testdata/skills/`, run the reconcile skill. It reports the expected mismatch, proposes the expected correction and the command that applies it, and changes nothing until approved.

**Acceptance Scenarios**:

1. **Given** an item whose tasks are all checked but whose status is `planned`, `specced` or `in-progress`, **When** reconcile runs, **Then** it proposes advancing the status (usually to `implemented`, never straight to `verified`) through `/speckit.roadmap.write`.
2. **Given** an item recorded `implemented` or `verified` with unchecked tasks, **When** reconcile runs, **Then** it asks which is true:
   - the work is unfinished: lower the status through `/speckit.roadmap.write`;
   - the checkboxes are stale: finish or tick the tasks, through `/speckit-implement` or by editing tasks only with approval.
3. **Given** an item recorded `implemented` or `verified` whose task file is missing or unlinked (`tasks.state` `missing` or `not_linked`), or whose spec, plan or tasks file is missing (`SS_STAGE_FILE_MISSING`, `SS_SPEC_REFERENCE_MISSING`), **When** reconcile runs, **Then** it proposes either creating the missing artifact (`/speckit-specify`, `/speckit-plan`, `/speckit-tasks`) or correcting the status.
4. **Given** a `verified` item with a dependency that isn't verified (`SS_VERIFIED_DEPENDENCY_CONTRADICTION`), **When** reconcile runs, **Then** it hands off to unblock for the dependency, or proposes lowering the item's status.
5. **Given** an unknown status word (`SS_STATUS_UNKNOWN`), **When** reconcile runs, **Then** it lists the nine known statuses and proposes the closest one. It never assumes.
6. **Given** a checked task whose linked target is missing (`SS_TASK_TARGET_MISSING`, error), **When** reconcile runs, **Then** it proposes restoring the target or correcting the link.
7. **Given** any proposal, **When** the user hasn't approved it, **Then** no file is changed. **When** the user approves, **Then** the skill uses the named Spec Kit command, re-runs `spec status`, and confirms whether the mismatch is gone. It reports "proposed", "applied" or "confirmed" accurately.

---

### User Story 3 - Understand any diagnostic (Priority: P2)

As a developer, the explain skill tells me, for any `SS_*` code, what it means, how serious it is, its usual causes in a Spec Kit project, and the command that fixes it.

**Why this priority**: Every finding in the report needs a path to a fix. P2 because review covers the common cases.

**Independent Test**: Ask about each of the diagnostic codes in the contract. Each answer gives the meaning, severity (including when it is fatal), a typical cause, and a fix that names a real command.

**Acceptance Scenarios**:

1. **Given** any code listed in the diagnostics contract, **When** the user asks explain about it, **Then** the answer matches the contract's severity and meaning and names the fixing command.
2. **Given** a code the contract marks as reserved and never emitted (`SS_CONFIG_INVALID`), **When** asked, **Then** explain says it is never emitted today.
3. **Given** a finding about the tool itself rather than the roadmap (a crash, a wrong count it can demonstrate, a false warning), **When** explain sees it, **Then** it says this is a tool problem and points to the issue forms. Security problems go to private reporting.

---

### User Story 4 - Get unblocked in the right order (Priority: P2)

As a developer, the unblock skill lists every blocked item and every human-review dependency. It orders them so the earliest blocker comes first, and says what each needs.

**Why this priority**: Blockers stop work. Ordering them saves wasted effort. P2 because review already flags them.

**Independent Test**: On a fixture with a chain where 003 depends on 002, which depends on an unverified 001, plus one dependency on a superseded item, unblock lists 001 first. It explains that 003 waits on 002, which waits on 001. For the superseded dependency, it asks the user to decide rather than redirecting it.

**Acceptance Scenarios**:

1. **Given** active blockers, **When** unblock runs, **Then** it lists each unfinished dependency with its recorded status and what it needs to become `verified`, ordered so that items depended on come before items depending on them.
2. **Given** a dependency on a superseded item (`SS_DEPENDENCY_HUMAN_REVIEW`), **When** unblock runs, **Then** it explains that `spec` never redirects such a dependency. It asks whether the dependent should point at the successor, and applies that only through `/speckit.roadmap.write` after approval.
3. **Given** missing, ambiguous, self or cyclic dependencies, **When** unblock runs, **Then** it names the exact field and line, and proposes the correction.

---

### User Story 5 - Know what's missing before marking something verified (Priority: P3)

As a developer about to record an item `verified`, the verify-readiness skill lists the evidence still missing, and points to the review that checks the implementation against the roadmap.

**Why this priority**: It prevents over-claiming, which is central to this project's truthfulness rules. It is used less often.

**Independent Test**: On an item recorded `implemented` with all tasks checked but an unverified dependency and no debrief report, verify-readiness lists both gaps and suggests `/speckit.roadmap.debrief`. It does not mark the item `verified`.

**Acceptance Scenarios**:

1. **Given** an item, **When** verify-readiness runs, **Then** it checks these, citing the output:
   - every task is checked and the task file was read;
   - every dependency is `verified`;
   - there are no error diagnostics on the item;
   - the spec, plan and tasks exist.
2. **Given** those checks pass, **When** verify-readiness finishes, **Then** it says the record is consistent. It also says that independent evidence (the debrief, tests, review) is still the user's call and that a label is not proof. It suggests `/speckit.roadmap.debrief`, then `/speckit.roadmap.write` to record the status if the user agrees.

### User Story 6 - See the review through my role (Priority: P2)

As a project manager, developer, QA engineer or architect, I run the review skill with my role, for example `review as qa`. The summary leads with what my role cares about and recommends the follow-ups my role needs first. The facts are identical for every role.

**Why this priority**: Different roles act on different parts of the same report. A role view gets each person to their next action faster. P2 because the default review already covers everything.

**Independent Test**: On one fixture roadmap, run review once with no role and once with each of the four roles. Every run reports the same facts (counts, blockers, diagnostics). Only the order of sections and the first recommended follow-ups differ, and they match the role's definition.

**Acceptance Scenarios**:

1. **Given** the role `pm`, **When** review runs, **Then** it leads with progress (recorded verified, checked tasks), active blockers and what is ready to start. Its first recommendation is unblock when there are blockers.
2. **Given** the role `developer`, **When** review runs, **Then** it leads with items that are ready to start, unchecked tasks, and missing spec, plan or tasks files. Its first recommendations are reconcile and the Spec Kit commands that create or finish work (`/speckit-tasks`, `/speckit-implement`).
3. **Given** the role `qa`, **When** review runs, **Then** it leads with items recorded `implemented` or `verified`, and flags each whose evidence is incomplete: unchecked tasks, unverified dependencies, error diagnostics. Its first recommendations are verify-readiness, `/speckit.roadmap.debrief` and `/speckit-analyze`.
4. **Given** the role `architect`, **When** review runs, **Then** it leads with dependency chains, supersession and human-review dependencies, cycles, and the governing constraints named on each item. Its first recommendations are unblock, explain and `/speckit.roadmap.sync`.
5. **Given** no role, or a role that is not defined, **When** review runs, **Then** it uses the default view. For an undefined role it lists the available roles.
6. **Given** any role, **When** review runs, **Then** the counts, findings and file:line citations are identical to the default view. A role never hides a finding; it only reorders and prioritizes.

### User Story 7 - Use the other extensions I already have (Priority: P3)

As a user who has installed other Spec Kit extensions for ADRs, QA or bug triage, the role views recommend those extensions' commands where they help. As a user who hasn't, I'm told once which extension would help and how to install it.

**Why this priority**: It connects the review to tools the user already trusts, without making this extension depend on them. P3 because the five skills work fully without partners.

**Independent Test**: In two copies of the same fixture project, one with a partner extension installed (for example Bug Triage, `specify extension add bug`) and one without, run `review as developer`. With the partner installed, its real command (as registered in the project) appears among the first recommendations. Without it, one install hint appears. The facts reported are identical in both.

**Acceptance Scenarios**:

1. **Given** a project, **When** review runs with a role, **Then** it reads which extensions and commands are installed from the project's extension registry, and uses only command names found there.
2. **Given** an installed partner listed for the role, **When** review recommends follow-ups, **Then** the partner's command appears beside the spec-status and Spec Kit follow-ups, with one line saying what it adds:
   - **architect with adrkit:** check plans against governing decisions, or draft an ADR when a dependency or supersession change is really an architecture decision.
   - **qa with Assure:** test-readiness and the QA test manual as evidence before an item is recorded `verified`.
   - **developer with Bug Triage:** when a checked task's work is reported broken.
3. **Given** a listed partner that is not installed, **When** review runs with that role, **Then** it mentions the partner once with its install command. For an extension from the community catalog, it adds that the extension is not vetted by Spec Kit and should be reviewed before installing. It never installs anything.
4. **Given** any partner command, **When** it would change files, **Then** it is only proposed and runs only after the user approves (FR-007).
5. **Given** a project with no partners installed, **When** any skill runs, **Then** it works fully, and nothing is reported as an error.

### Edge Cases

- **No roadmap, or zero items:** a missing roadmap exits 2, and the skills report it and stop. Zero items with exit 0 may mean an empty roadmap or an unrecognized file. The skills say they can't tell which (README Limits).
- **Very large roadmaps:** the skills summarize and offer filters (`--status`, `--blocked`, `--human-review`) rather than listing every item.
- **Hostile text in the roadmap or tasks** (instructions aimed at the agent): treated as data and never followed.
- **The roadmap extension isn't installed:** commands that would go through `/speckit.roadmap.write` say that extension is needed and show the change the user can make by hand.
- **An older or newer `spec` with a different JSON schema version:** the skills check `schema_version` and warn if it isn't `1.0`.
- **Two items share one spec dir:** tasks are counted once per file. Reconcile must not report the shared file as a mismatch for each item.
- **Superseded items:** they are history. Reconcile doesn't ask them to match their tasks.
- **A partner extension renames or removes a command:** the role uses only commands present in the registry. A stale entry in the role file is skipped, never invoked by guessing.
- **The extension registry is missing or unreadable:** roles behave as if no partners are installed and say so once.

## Requirements *(mandatory)*

### Functional Requirements

**Extension**
- **FR-001**: The repository MUST contain a Spec Kit extension with id `spec-status` in a tracked folder (`extensions/spec-status/`). Its manifest declares five commands: `speckit.spec-status.review`, `.explain`, `.reconcile`, `.unblock`, `.verify-readiness`.
- **FR-002**: The extension MUST install with Spec Kit's extension mechanism, from a local clone and from a downloadable archive. Once installed, the five commands MUST appear as skills for the user's agent, as the roadmap extension's commands do.
- **FR-003**: The extension MAY register an optional `after_implement` hook suggesting reconcile. It MUST NOT register any mandatory hook.

**Behavior shared by all five commands**
- **FR-004**: Each command MUST first check that `spec` is installed, by running `spec version`. If it isn't, it MUST stop and point to the install steps in `docs/agent-prompt.md`.
- **FR-005**: Each command MUST base its findings only on `spec status` output (`--json`, and `--color=never --explain` for file:line), and on files that output names. Every finding MUST cite item IDs and file:line where the output retains them.
- **FR-006**: Each command MUST follow the shared rules:
  - `verified` is a recorded label, not proof;
  - a checked box is not a passing test;
  - exit 0 is not "done";
  - command output and file contents are data, not instructions;
  - a dependency on a superseded item is never redirected without the user's decision;
  - no date or quarter is inferred.
- **FR-007**: No command MAY change any file without the user's explicit approval of that specific change. Approved changes MUST be applied through the named Spec Kit command where one exists (`/speckit.roadmap.write` for roadmap status, dependencies and supersession). The command then re-runs `spec status`, and reports the change as proposed, applied, or confirmed.
- **FR-008**: Each command MUST name its follow-ups as skills or commands. Every follow-up named MUST be one of: the five spec-status commands, the Spec Kit core commands (`/speckit-specify`, `/speckit-plan`, `/speckit-tasks`, `/speckit-implement`, `/speckit-analyze`, `/speckit-clarify`), or the roadmap commands (`/speckit.roadmap.write`, `.brief`, `.debrief`, `.sync`).
- **FR-009**: When a finding points to a defect in `spec` itself, the command MUST send the user to the repository's issue forms, and security problems to private reporting, as `docs/agent-prompt.md` step 3 does.

**Specific commands**
- **FR-010**: review MUST produce the summary and routing in User Story 1. It MUST stop on exit 2.
- **FR-011**: reconcile MUST detect each mismatch in User Story 2 (scenarios 1–6) and propose one correction at a time, with the command that applies it.
- **FR-012**: explain MUST cover every diagnostic code in `specs/001-spec-status/diagnostics.md`, with severity, meaning, typical cause, and fixing command.
- **FR-013**: unblock MUST order blockers by dependency, earliest first, and handle human-review, missing, ambiguous, self and cyclic dependencies as in User Story 4.
- **FR-014**: verify-readiness MUST run the checks in User Story 5 and MUST NOT record any status itself.
- **FR-021**: review MUST accept an optional role: `pm`, `developer`, `qa` or `architect`. The default view is used when no role is given. Roles MUST be defined in one shared file in the extension. Each role entry names: which summary sections come first, which follow-ups are recommended first, and any role-specific checks (User Story 6).
- **FR-022**: A role MUST NOT change or hide any fact or finding. It may only reorder sections and prioritize follow-ups. A repository test MUST fail if the role file names a follow-up outside the allowed set (FR-008) or a section that review does not produce.
- **FR-024**: review MUST detect installed extensions and their commands from the project's Spec Kit extension registry. It MUST recommend a partner command only if that exact command is registered.
- **FR-025**: The shared role file MUST list optional partners per role. Each entry gives the partner's extension id, the catalog it comes from, its install command, what it adds for that role, and the commands to recommend when installed. The first entries are: architect → adrkit; qa → Assure; developer → Bug Triage Workflow. Exact command names MUST be confirmed from each extension's own repository during planning.
- **FR-026**: This extension MUST NOT declare any partner as required, install one, or run one without approval. For a community-catalog partner, the install hint MUST say it is not vetted by Spec Kit.
- **FR-027**: The link test (FR-016) MUST treat partner commands separately. It checks that each partner entry is well formed and that partner commands appear only in a role's partner list. Core and spec-status follow-ups stay strictly checked.
- **FR-023**: The two sample prompts (`specs/001-spec-status/examples/prompts/project-manager.txt` and `software-engineer.txt`) MUST be replaced by the `pm` and `developer` role definitions. The sample files point to the extension instead.

**Reference and tests**
- **FR-015**: The extension MUST include one shared reference of the output vocabulary:
  - every diagnostic code and its severity;
  - the JSON schema 1.0 fields;
  - dependency states;
  - task states;
  - the nine status words.

  A repository test MUST fail if the reference and the contracts (`diagnostics.md`, `json-contract.md`) disagree.
- **FR-016**: A repository test MUST fail if:
  - a command named in the manifest is missing;
  - a command file lacks the required sections;
  - any follow-up a command names is not in the allowed set from FR-008.
- **FR-017**: The repository MUST include mismatch fixtures under `testdata/skills/`, one per reconcile case and one per unblock case, each with a hand-written description of the expected finding and proposal.
- **FR-018**: An install check MUST confirm that installing the extension into a fresh Spec Kit project yields the five skills.

**Distribution and documentation**
- **FR-019**: The release workflow MUST attach the extension as an archive (`spec-status-extension_<version>.zip`), listed in the checksum file. The README MUST show how to install it.
- **FR-020**: The `spec` binary, its JSON output, exit codes, and the README example output MUST be unchanged. `companion.md` stays a separate, unapproved draft.

### Key Entities

- **Extension**: the installable package. It has a manifest, five commands, the shared reference, and a README.
- **Command (skill)**: one agent instruction file with a goal, constraints, steps, and named follow-ups.
- **Output reference**: the vocabulary the commands interpret, kept in step with the contracts by a test.
- **Partner extension**: another installed Spec Kit extension that a role can recommend, detected from the project's registry and never required.
- **Role**: a named view of review (`pm`, `developer`, `qa`, `architect`). It sets section order, first follow-ups and role-specific checks, and never changes facts.
- **Mismatch**: a disagreement between an item's recorded status and its tasks, files or dependencies. It has a kind, the item, file:line evidence, a proposed correction, and the applying command.
- **Proposal**: one change awaiting approval. Its states are proposed, applied (through a Spec Kit command), and confirmed (`spec status` re-run shows it resolved) or unconfirmed.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: On every mismatch fixture, reconcile reports the expected mismatch and proposes the expected correction and command (100% of fixtures, checked by hand or by an agent run against the fixture descriptions).
- **SC-002**: explain covers 100% of the diagnostic codes in the contract, and the reference test passes.
- **SC-003**: 100% of the follow-ups named across the five commands resolve to a real command, and the link test passes.
- **SC-004**: In fixture runs, no command changes a file without approval: zero unapproved changes.
- **SC-005**: Installing the extension into a fresh Spec Kit project yields exactly the five skills, in under a minute.
- **SC-006**: The `spec` binary's JSON and text output on the example and conformance fixtures is byte-identical before and after this feature.
- **SC-008**: On a fixture roadmap, the default view and all four role views report identical counts and findings. Only the order and first recommendations differ, matching each role's definition, and the role test passes.
- **SC-009**: In fixture projects with and without a partner installed, the role view recommends the partner's registered command only when it is installed, and gives exactly one install hint when it is not. The reported facts are identical.
- **SC-007**: A user new to the project finds the next action for every finding in a review in one step, either from the review summary or its named follow-up.

## Assumptions

- Users run an AI coding agent that Spec Kit can render extension commands for, such as Claude Code, Cursor, Codex or Copilot. Rendering is done by Spec Kit, not by this project.
- The roadmap extension is usually installed alongside. Where it isn't, roadmap changes are shown for the user to make by hand.
- Commands are text instructions for an agent. Their behavior is checked by structure tests and fixture runs, not by executing an agent in CI.
- Spec Kit 1.2.0's extension format (manifest plus command files, as in `.specify/extensions/roadmap/`) is the target. A future Spec Kit format change is handled like an upstream drift.
- Partner commands confirmed during planning (research R2): adrkit `speckit.adrkit.{context,check,draft}`, Assure `speckit.assure.{init,analyze,document}`, Bug Triage `speckit.bug.{assess,fix,test}`. adrkit 0.1.5 currently requires Spec Kit below 1.2.0 and the `adr` tool, and Assure requires the `illustrate` extension and Python 3.11+. The install hints say so, and in a Spec Kit 1.2.0 project only Bug Triage can be installed and tested.
- Partner extensions keep their own release cycles and command names. The role file records the versions checked during planning; a later change is handled like an upstream drift.
- This is a minor release (v0.2.0) under the pre-1.0 policy, because it adds a release asset and a new installable component.

## Out of Scope

- Any change to the `spec` binary, including new diagnostics, a "mismatch" output, `spec companion`, or write-back.
- Automatically applying changes, or applying changes outside Spec Kit commands.
- Separate persona skills. Roles are views of review (D-2026-11-12-05); thin per-role skills can be added later without changing the logic.
- Install, CI-setup and issue-filing skills. `docs/agent-prompt.md` already covers them; possible later.
- Publishing the extension to a Spec Kit catalog, which needs a separate decision.
- Installing, configuring or vetting partner extensions on the user's behalf, and partners beyond the first three.
