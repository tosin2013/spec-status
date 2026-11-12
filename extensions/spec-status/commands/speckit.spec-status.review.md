---
description: "Summarize where the Spec Kit roadmap stands from spec status output, and send every finding to the command that can resolve it. Optional role view: pm, developer, qa or architect."
argument-hint: "[as pm|developer|qa|architect] [focus]"
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Goal

Give the user an accurate, short picture of the roadmap: what is recorded verified, how many tasks are checked, what is blocked and why, which dependencies need a human decision, which findings `spec status` reported, and what is ready to start. Every finding is cited and sent to exactly one follow-up command. With a role (`as qa`, for example) the same facts are reordered for that role.

## Operating Constraints

- Read `.specify/extensions/spec-status/reference/rules.md` first and follow it. The vocabulary is in `.specify/extensions/spec-status/reference/output.md`; roles and partner extensions are in `.specify/extensions/spec-status/reference/roles.md`.
- **Read-only.** This command never changes a file and never proposes edits itself; it routes findings to the commands that do (they follow the approval protocol in `rules.md`).
- A role never changes or hides a fact. Counts, findings and citations are identical in every view.

## Outline

1. Run `spec version`. If it fails, stop and tell the user to install `spec` (step 1 of https://github.com/tosin2013/spec-status/blob/main/docs/agent-prompt.md).
2. Parse `$ARGUMENTS`: `as <role>` selects a role from `roles.md` (`pm`, `developer`, `qa`, `architect`); anything else is a focus (for example an item ID). No role means `default`. An unknown role: list the roles and use `default`.
3. From the repository root, run `spec status --json` and `spec status --color=never --explain`.
   - Exit `2`: report the diagnostic (`error.code`, message, source) and stop. Do not create a roadmap.
   - Check `schema_version` is `1.0` (`rules.md`).
   - `summary.total_items` is `0`: say the roadmap is empty or not in a supported shape, and that `spec` cannot tell which; stop.
4. Build each section in `roles.md` "Review sections" that has content, from the JSON (facts) and the `--explain` output (`path:line`):
   - `progress`: `summary.verified_active_items` / `summary.active_items`, and `summary.tasks.checked` / `summary.tasks.total`. Say "recorded verified".
   - `blockers`: active items with a dependency whose `state` is not `satisfied`.
   - `human-review`: dependencies with `state` `human_review`.
   - `diagnostics`: group by severity (error, warning, info); one line each with code, item and `path:line`.
   - `ready-to-start`: active `planned` or `specced` items whose dependencies are all `satisfied`.
   - `unchecked-tasks`: active items with `tasks.state` `read` and `checked` < `total`.
   - `missing-artifacts`: `SS_SPEC_REFERENCE_MISSING`, `SS_SPEC_DIRECTORY_MISSING`, `SS_STAGE_FILE_MISSING`, and items with `tasks.state` `missing`.
   - `evidence-gaps`: items recorded `implemented` or `verified` with unchecked tasks, `tasks.state` other than `read`, a dependency that is not `satisfied`, or an error diagnostic.
   - `dependency-chains`: for each blocked item, follow `dependencies` back to the first item that is not satisfied.
   - `supersession`: items with `superseded_by`, their successors, and `supersession.valid` false.
   - `constraints`: each affected item's `Governed by` line from the roadmap (cite it).
5. Order the sections by the selected role's "Leads with" list, then the remaining sections in the order above. Keep every fact identical across roles.
6. Route every finding to one follow-up (see Follow-ups): a diagnostic the user may not understand → explain; a status that disagrees with tasks, files or dependencies → reconcile; a blocker or human-review dependency → unblock; an item the user wants to record `verified`, or any evidence gap → verify-readiness.
7. Partner extensions: read `.specify/extensions/.registry` (JSON). For each partner row in `roles.md` for the selected role, if `extensions.<id>.enabled` is true, recommend only those of its commands listed under `registered_commands`, with the row's "Adds" sentence. If a role's partner is not installed, mention it once with its install command and requirements (and the "not vetted" note for community partners). If the registry is missing or unreadable, say once that installed partners could not be detected. Never install or run a partner.
8. End with the role's first follow-ups, then at most three next actions, each citing the item and command.

## Follow-ups

- `/speckit.spec-status.reconcile`: a recorded status disagrees with tasks, files or dependencies.
- `/speckit.spec-status.unblock`: blockers and dependencies on superseded items.
- `/speckit.spec-status.explain`: what a diagnostic code means and how to fix it.
- `/speckit.spec-status.verify-readiness`: before recording an item `verified`, or for any evidence gap.
- `/speckit-specify`, `/speckit-plan`, `/speckit-tasks`, `/speckit-implement`: create or finish the work for an item.
- `/speckit.roadmap.sync`: check the roadmap against the spec folders on disk.
