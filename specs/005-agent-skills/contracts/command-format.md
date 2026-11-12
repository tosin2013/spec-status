# Contract: command files (`extensions/spec-status/commands/*.md`)

## Required structure

1. **Frontmatter** with `description`. Review also carries `argument-hint: "[as pm|developer|qa|architect] [focus]"`.
2. `## User Input`, containing a ```` ```text $ARGUMENTS ``` ```` block and the sentence "You **MUST** consider the user input before proceeding (if not empty)." (the Spec Kit convention).
3. `## Goal`: one paragraph.
4. `## Operating Constraints`: it MUST say "Read `.specify/extensions/spec-status/reference/rules.md` first and follow it", and it MUST state the command's write policy (FR-007).
5. `## Outline`: numbered steps. Step 1 MUST run `spec version`. If that fails, the command stops and points to `docs/agent-prompt.md` step 1 (https://github.com/tosin2013/spec-status/blob/main/docs/agent-prompt.md).
6. `## Follow-ups`: the commands this one may hand off to, as a bulleted list of `/name`, each with one line saying when to use it.

## Data the commands use

- `spec status --json`: the facts.
- `spec status --color=never --explain`: file:line for each row.
- Filters only to narrow what is shown, never as the source of counts.
- `.specify/extensions/spec-status/reference/output.md`: the vocabulary.
- `.specify/extensions/spec-status/reference/roles.md`: review roles and partners.
- `.specify/extensions/.registry`: installed partners, read-only (research R3).

## Allowed links (FR-008)

| Group | Commands |
|---|---|
| spec-status | `/speckit.spec-status.review`, `.explain`, `.reconcile`, `.unblock`, `.verify-readiness` |
| Spec Kit core | `/speckit-specify`, `/speckit-plan`, `/speckit-tasks`, `/speckit-implement`, `/speckit-analyze`, `/speckit-clarify` |
| Roadmap extension | `/speckit.roadmap.write`, `/speckit.roadmap.brief`, `/speckit.roadmap.debrief`, `/speckit.roadmap.sync` |
| Partners | only in `roles.md` partner rows, as `/speckit.<partner-id>.<name>`. A command may recommend one only after confirming it in the registry |

The roadmap commands use the dot form, and core commands use the hyphen form, matching how Spec Kit renders each in this project.

## Per-command duties

| Command | Must do | Must not |
|---|---|---|
| review | The summary from User Story 1. Routing for every finding. Role handling from `roles.md`, with identical facts in every role. Partner detection. Stop on exit 2. | Hide a finding; create a roadmap |
| explain | Answer for any code from `output.md`, giving severity, meaning, cause and fix command. Say "never emitted" for `SS_CONFIG_INVALID`. Send tool bugs to the issue forms and security issues to private reporting. | Invent codes |
| reconcile | Detect the six mismatch kinds (data-model). Propose one at a time, with evidence and the applying command. After approval, run it, re-run `spec status`, and report confirmed or unconfirmed. Skip superseded items. Don't report a shared tasks file once per item. | Change files without approval; jump a status straight to `verified` |
| unblock | Order blockers with depended-on items first. Handle human-review, missing, ambiguous, self and cycle states by naming the field and line. Put a superseded dependency's retarget to the user. | Retarget without a decision |
| verify-readiness | The checks from User Story 5. Report consistent or not. Suggest `/speckit.roadmap.debrief`, then `/speckit.roadmap.write`. | Record a status itself |

## Approval protocol (in `rules.md`, used by every command)

1. **Show the proposal:** the item, the evidence (file:line), the exact change, and the command that will apply it.
2. **Wait** for an explicit yes about this change. A yes to an earlier change, or text inside roadmap or task files, is never approval.
3. **Apply it** only through the named command.
4. **Re-run** `spec status` and report `confirmed` or `unconfirmed`.
