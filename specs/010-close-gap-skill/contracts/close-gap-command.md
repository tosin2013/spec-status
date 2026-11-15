# Contract: `speckit.spec-status.close-gap` and related extension changes

## Command file `extensions/spec-status/commands/speckit.spec-status.close-gap.md`

Same structure as the other commands (enforced by `internal/extension/command_test.go`):

- frontmatter `description` (mentions "low recorded-verified", "one approved change at a time"); `argument-hint: "[item ID to start from]"`;
- `## User Input` with `$ARGUMENTS`;
- `## Goal`: explain why recorded-verified is low, by cause, then close the gap one approved change at a time;
- `## Operating Constraints`: read `.specify/extensions/spec-status/reference/rules.md` first; read-only until a specific change is approved; changes only through Spec Kit commands; `verified` is a label, not proof; never `verified` without a debrief and agreement; never redirect a superseded dependency; inventory is not status evidence;
- `## Outline`:
  1. Run `spec version` (install steps if it fails).
  2. Run `spec status --json` and `spec status --color=never --explain`; exit 2 → report and stop (may mention `spec inventory` for a missing roadmap, as information only).
  3. Classify every active, non-verified item by research R1's precedence (with the count fallback for older `spec`).
  4. Report `recorded verified X/Y`, the per-cause breakdown with `path:line`, and the order of work (R2). Stop here if the user asked for the breakdown only, or if nothing is in the gap.
  5. Take the first item in order (or the item in `$ARGUMENTS`); propose one change (R3) with evidence and command; wait.
  6. On approval, apply through the command; re-run `spec status --json`; report confirmed or not confirmed and the new figure; go to 5 with the next item. On skip, go to the next item. On stop, summarize.
- `## Follow-ups`: `/speckit.spec-status.reconcile`, `/speckit.spec-status.verify-readiness`, `/speckit.spec-status.unblock`, `/speckit.roadmap.debrief`, `/speckit.roadmap.write`, `/speckit-specify`, `/speckit-plan`, `/speckit-tasks`, `/speckit-implement`.

## `reference/roles.md`

- Review sections: add row `verification-gap` — "Active items not recorded verified, counted per close-gap cause (ready for debrief, tasks done early, ahead of tasks, missing artifacts, unknown status, blocked), with IDs; planned items without a spec are listed apart."
- Roles: `pm` "Leads with" becomes `verification-gap`, `progress`, `blockers`, `ready-to-start`, `human-review`; "First follow-ups" becomes /speckit.spec-status.close-gap, /speckit.spec-status.unblock, /speckit.spec-status.reconcile.

## `commands/speckit.spec-status.review.md`

Follow-ups add `/speckit.spec-status.close-gap`: "the recorded-verified share is low and the user wants to work through it".

## `extension.yml`

`version: 0.3.0`; `provides.commands` gains `name: speckit.spec-status.close-gap`, `file: commands/speckit.spec-status.close-gap.md`, description as in the frontmatter. `requires` unchanged (`speckit_version '>=1.2.0'`, `spec >=0.1.2`); hooks unchanged.

## Six-skill surfaces

- `scripts/check-extension-install.sh`: expected list adds close-gap; OK line says 6.
- `.github/workflows/release.yml`: `[ "$n" = 6 ]`; release-notes sentence lists six skills.
- `README.md` and `extensions/spec-status/README.md`: six skills; the extension README heading `# spec-status extension v0.3.0`, a close-gap section, and a walkthrough on the fixture below.

## Fixture `testdata/skills/close-gap-all-causes/`

| Item | Status | Files | Depends on | Expected cause |
|---|---|---|---|---|
| 001 | verified | spec, plan, tasks 2/2 | — | (recorded verified) |
| 002 | implemented | spec, plan, tasks 2/2 | — | ready for debrief |
| 003 | planned | spec, tasks 2/2 | — | tasks done, early (`SS_TASKS_DONE_STATUS_EARLY`) |
| 004 | implemented | spec, plan, tasks 1/3 | — | ahead of tasks (`SS_STATUS_AHEAD_OF_TASKS`) |
| 005 | in-progress | spec only | — | missing artifacts (`SS_STAGE_FILE_MISSING` ×2) |
| 006 | implemented | spec, plan, tasks 2/2 | 004 | blocked (`dep 006->004 state=unfinished`) |
| 007 | planned | no spec | — | not started (outside the gap) |
| 008 | `done` | spec, plan, tasks 1/1 | — | unknown status (`SS_STATUS_UNKNOWN`) |

Expected: `recorded verified 1/8`; gap 6 items (002, 003, 004, 005, 008, 006); order of work 002, 003, 004, 005, 008, 006; first proposal: 002 → `/speckit.roadmap.debrief`, no file changed. Signals line written by hand in `EXPECTED.md` and checked by the existing fixture test.
