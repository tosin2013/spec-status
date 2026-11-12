# Data model: Agent skills extension

No stored data. The entities are files in the extension and the states a proposal moves through.

## Extension (`extensions/spec-status/extension.yml`)

| Field | Value / rule |
|---|---|
| schema_version | `'1.0'` |
| extension.id | `spec-status` |
| extension.version | equals the release tag without `v` at release time (R7) |
| extension.license | `Apache-2.0` |
| requires.speckit_version | `'>=1.2.0'`, the version these command files were checked against |
| requires.tools | `[{name: spec, required: true, version: '>=0.2.0'}]`. Documentation only in Spec Kit 1.2.0 (R1); commands check `spec version` |
| provides.commands | exactly the five `speckit.spec-status.*` commands, each with `file: commands/<name>.md` |
| hooks.after_implement | `{command: speckit.spec-status.reconcile, optional: true}`; no other hooks |

## Command (`commands/*.md`)

| Part | Rule |
|---|---|
| frontmatter | `description` (required); `argument-hint` (review: `[as pm|developer|qa|architect] [focus]`) |
| sections | `## User Input` (with `$ARGUMENTS`), `## Goal`, `## Operating Constraints`, `## Outline`, `## Follow-ups` |
| first step | read `.specify/extensions/spec-status/reference/rules.md`; run `spec version` (FR-004) |
| links | only from the allowed set (FR-008), or from the role's partner list |

## Output reference (`reference/output.md`)

Tables:
- **Codes:** `| code | severity | fatal when | meaning | typical cause | fix (command) |`, one row per `SS_*` code (32).
- **JSON fields:** root, summary, item, dependency and diagnostic fields.
- **Dependency states** (7), **task states** (4), and **status tokens** (9).

## Role (`reference/roles.md`)

| Field | Rule |
|---|---|
| role | `pm`, `developer`, `qa`, `architect`; `default` when none is given |
| leads with | ordered section names from review's section list |
| first follow-ups | ordered; allowed set only |
| role checks | short, for example qa: "flag `implemented`/`verified` items with unchecked tasks, an unverified dependency, or error diagnostics" |
| partners | rows of the partner table for this role |

## Partner (rows in `reference/roles.md`)

| Field | Example (bug) |
|---|---|
| role | developer |
| extension id | `bug` |
| version checked | 1.0.2 |
| catalog | default |
| install | `specify extension add bug` |
| requirements | Spec Kit `>=0.9.0` |
| commands | `speckit.bug.assess`, `speckit.bug.fix`, `speckit.bug.test` |
| adds | "when a checked task's work is reported broken" |

## Mismatch (produced by reconcile)

| Field | Rule |
|---|---|
| kind | one of: `tasks-done-status-early`, `status-ahead-of-tasks`, `status-ahead-of-artifacts`, `verified-dependency-unverified`, `unknown-status`, `checked-target-missing` |
| item | ID, plus file:line from `--explain` |
| evidence | the JSON fields and diagnostic codes that show it |
| proposal | one change and the command that applies it |

## Proposal states

`proposed` → (user approves) → `applied` (through the named Spec Kit command) → `confirmed` (re-run `spec status` shows the mismatch is gone) or `unconfirmed` (still present; report and stop). Declining leaves the proposal as `proposed`, and nothing changes.
