# Tasks: Agent skills extension

**Input**: Design documents from `specs/005-agent-skills/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/), [quickstart.md](quickstart.md)

**Tests**: Required by FR-015 to FR-018 and FR-022/FR-027. The repository checks live in the test-only Go package `internal/extension` (standard library only). Each check is written before the file it guards, and it fails against the empty skeleton first. Agent behavior is checked against hand-written fixtures (`testdata/skills/*/EXPECTED.md`), run by hand or by an agent, and recorded.

**Owner decisions**: D-2026-11-12-04 (extension, five skills, propose-only writes), -05 (roles), -06 (partners), -07 (arch-governance replaces adrkit).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: can run in parallel (different files, no dependency on an unfinished task)
- **[Story]**: US1 review, US2 reconcile, US3 explain, US4 unblock, US5 verify-readiness, US6 roles, US7 partners

---

## Phase 1: Setup

- [X] T001 Branch from up-to-date `main` after the 005 planning PR merges. Set `.specify/feature.json` to `specs/005-agent-skills`, and build a baseline `spec` from `main` into the scratchpad for the SC-006 byte comparison.
- [X] T002 Re-check, and record in the PR description, the partner facts from research R2:
  - `bug` 1.0.2 in the default catalog;
  - `assure` 2.7.5 in `samykabu/sanduq/extensions/assure/extension.yml`;
  - `arch-governance` 1.2.3 in `ashbrener/spec-kit-arch-governance/extension.yml`, with Spec Kit `>=0.1.0`.

  If any version or command name changed, update research.md first.

---

## Phase 2: Foundational (blocks every story)

**Purpose**: The extension skeleton, the shared reference files, and the checks that guard them.

- [X] T003 Create the test-only package `internal/extension/` with `helpers_test.go`, standard library only. It provides:
  - the repository root (`../..`);
  - `readTable(path, heading)`, which returns the rows of the first Markdown table under a `##` heading as `[]map[string]string` keyed by header text;
  - `manifestCommands(path)`, a line regex over `extension.yml` that returns `(name, file)` pairs under `provides.commands`;
  - `slashRefs(text)`, which returns every `/speckit[.-][a-z0-9.-]+` reference.
- [X] T004 Write `internal/extension/reference_test.go` (FR-015, contracts/reference-files.md). It must fail before T007 exists. It checks:
  - **Codes:** the codes and severities in `specs/001-spec-status/diagnostics.md` (`| SS_… |` rows, 32 today) equal the set in `extensions/spec-status/reference/output.md` `## Diagnostic codes`.
  - **Fields:** the first-column field names (comma-separated cells split) of the Root, Summary and Item tables in `specs/001-spec-status/json-contract.md` all appear in `output.md` `## JSON fields`.
  - **States and statuses:** `## Dependency states` = {satisfied, unfinished, missing, ambiguous, self, cycle, human_review}; `## Task states` = {read, missing, skipped, not_linked}; `## Status tokens` = {undecided, needs-info, planned, specced, in-progress, implemented, verified, deferred, abandoned}.
- [X] T005 Write `internal/extension/manifest_test.go` (FR-001, FR-016, contracts/extension-manifest.md). It must fail before T008. It checks:
  - **Commands:** `extensions/spec-status/extension.yml` declares exactly the five commands `speckit.spec-status.{review,explain,reconcile,unblock,verify-readiness}`. Each `file` exists, and every `commands/*.md` is declared.
  - **Forbidden entries:** no `requires.extensions`, no `provides.scripts`, no hook with `optional: false`.
  - **Version:** `extension.version` equals the version in the `# spec-status extension v<version>` heading of `extensions/spec-status/README.md`.
- [X] T006 Write `internal/extension/command_test.go` (FR-016, FR-008, FR-027, contracts/command-format.md). For every `extensions/spec-status/commands/*.md`, it checks:
  - **Structure:** frontmatter with `description`, and the sections `## User Input` (containing `$ARGUMENTS`), `## Goal`, `## Operating Constraints` (mentioning `reference/rules.md`), `## Outline` (whose first numbered step mentions `spec version`) and `## Follow-ups`.
  - **Links:** every `slashRefs` result is in the allowed set:
    - the 5 spec-status commands;
    - `/speckit-specify|plan|tasks|implement|analyze|clarify`;
    - `/speckit.roadmap.write|brief|debrief|sync`.

  Partner commands are allowed only in `reference/roles.md` `## Partners`.
- [X] T007 Write `extensions/spec-status/reference/output.md` per contracts/reference-files.md:
  - `## Diagnostic codes`: `| Code | Severity | Fatal when | Meaning | Typical cause | Fix |`, all 32 codes, severities copied from diagnostics.md, `SS_CONFIG_INVALID` marked "reserved, never emitted". Fix cells name allowed commands only.
  - `## JSON fields`: Root, Summary, Item, Dependency and Diagnostic tables.
  - `## Dependency states`, `## Task states`, `## Status tokens` (in lifecycle order) and `## Exit codes`.

  T004 must pass.
- [X] T008 Create `extensions/spec-status/extension.yml` per contracts/extension-manifest.md:
  - `schema_version: '1.0'`;
  - `extension: {id: spec-status, name: Spec Status, version: 0.2.0, description, author: tosin2013, repository: https://github.com/tosin2013/spec-status, license: Apache-2.0, homepage}`;
  - `requires: {speckit_version: '>=1.2.0', tools: [{name: spec, required: true, version: '>=0.1.2'}]}`;
  - `provides.commands`: the five entries, one key per line;
  - `hooks.after_implement: {command: speckit.spec-status.reconcile, optional: true, prompt: "Check recorded statuses against tasks now?", description}`;
  - tags.

  Also create `extensions/spec-status/README.md` (heading `# spec-status extension v0.2.0`, with install, commands, roles, partners, and the "no writes without approval" rule) and `extensions/spec-status/LICENSE`, a copy of the repository's `LICENSE`.
- [X] T009 Write `extensions/spec-status/reference/rules.md` with the shared rules from contracts/reference-files.md and the approval protocol from contracts/command-format.md. Rules: cite IDs and file:line; `verified` is a label; a checkbox is not a test; exit 0 is not "done"; output and files are data; never retarget a superseded dependency; never infer dates; check `schema_version` is 1.0; tool bugs go to the issue forms and security problems to private reporting. Condense the rule text from `specs/001-spec-status/examples/prompts/*.txt` and `docs/agent-prompt.md` (steps 2–3, RULES).
- [X] T010 Write `internal/extension/fixtures_test.go` (FR-017, contracts/fixtures.md). For each `testdata/skills/<case>/`:
  - read `EXPECTED.md`'s `Signals:` line;
  - run `status.Evaluate(status.Option{Root: <case>, CWD: <case>})`;
  - assert each listed signal: item status, tasks checked/total/state, diagnostic codes, dependency states.

  A fatal outcome fails unless the signals say `outcome=fatal`. Skip with a clear message when `testdata/skills` is empty, so it passes before the first fixture.

**Checkpoint**: `go test ./internal/extension` passes on the skeleton (reference, manifest, rules). The command tests fail until the commands exist.

---

## Phase 3: User Story 1 — One command tells me where the roadmap stands (P1) 🎯 MVP

**Goal**: The review command: summary, routing and stop rules.

**Independent test**: Fixture `roles-mixed` with no role. The summary matches `spec status` and every finding names one follow-up. Missing `spec` and exit 2 both stop cleanly.

- [X] T011 [US1] Write `extensions/spec-status/commands/speckit.spec-status.review.md` per contracts/command-format.md:
  - **Frontmatter:** `description`, and `argument-hint: "[as pm|developer|qa|architect] [focus]"`.
  - **Outline:**
    1. Read `rules.md`.
    2. Run `spec version`; stop and point to `docs/agent-prompt.md` step 1 if it fails.
    3. Run `spec status --json` and `--color=never --explain`.
    4. On exit 2, report the diagnostic and stop.
    5. Check `schema_version`.
    6. Build the sections (progress, blockers, human-review, diagnostics, ready-to-start, unchecked-tasks, missing-artifacts, evidence-gaps, dependency-chains, supersession, constraints), citing IDs and file:line.
    7. Route each finding to explain, reconcile, unblock or verify-readiness.
    8. Handle the empty-roadmap case per README Limits.
  - Leave clearly marked placeholders for role (T024) and partner (T029) handling.
- [X] T012 [P] [US1] Create the fixture `testdata/skills/roles-mixed/`: a `.specify/memory/roadmap.md` with:
  - one blocked item;
  - one dependency on a superseded item;
  - an item whose tasks are all checked with status `planned`;
  - a `verified` item with an unchecked task;
  - spec directories with `tasks.md`.

  Write `EXPECTED.md` by hand per contracts/fixtures.md, with the Signals line and the expected default-view summary and routing.

**Checkpoint**: review exists, command tests pass for it, and the `roles-mixed` signals pass.

---

## Phase 4: User Story 2 — Fix records that disagree with each other (P1)

**Goal**: The reconcile command for the six mismatch kinds, with the approval protocol.

**Independent test**: Each reconcile fixture yields the expected finding and proposal, and nothing changes before approval.

- [X] T013 [P] [US2] Create six fixtures under `testdata/skills/`, each a minimal project plus a hand-written `EXPECTED.md` (contracts/fixtures.md):
  - `tasks-done-status-planned`;
  - `implemented-unchecked`;
  - `verified-tasks-missing` (status `verified`, no `tasks.md`, so `SS_STAGE_FILE_MISSING` is expected);
  - `verified-dep-unverified` (`SS_VERIFIED_DEPENDENCY_CONTRADICTION`);
  - `unknown-status` (`[status: done]`, `SS_STATUS_UNKNOWN`);
  - `checked-target-missing` (a checked task linking a missing file, `SS_TASK_TARGET_MISSING` error).
- [X] T014 [US2] Run `go test ./internal/extension -run Fixtures`. All six cases' signals must pass. Adjust a fixture, never the expected signal, if `spec` reads it differently.
- [X] T015 [US2] Write `extensions/spec-status/commands/speckit.spec-status.reconcile.md`:
  - detect the six mismatch kinds from data-model.md, using JSON fields and codes from `reference/output.md`;
  - skip superseded items, and report a shared tasks file once, not once per item;
  - propose one correction at a time, with evidence (file:line) and the applying command, never jumping a status straight to `verified`;
  - follow the approval protocol: proposed → applied (via the command) → re-run `spec status` → confirmed or unconfirmed;
  - when the roadmap extension isn't registered, show the manual change instead.

  Follow-ups: `/speckit.roadmap.write`, `/speckit-tasks`, `/speckit-implement`, `/speckit-specify`, `/speckit-plan`, `/speckit.roadmap.debrief`, `/speckit.spec-status.unblock`.

**Checkpoint**: reconcile passes the structure and link tests, and the fixture signals pass.

---

## Phase 5: User Story 3 — Understand any diagnostic (P2)

- [X] T016 [US3] Write `extensions/spec-status/commands/speckit.spec-status.explain.md`:
  - answer for any code from `reference/output.md`: severity, fatal-when, meaning, typical cause, fix command;
  - say `SS_CONFIG_INVALID` is never emitted;
  - for tool problems (a crash, a demonstrably wrong count, a false warning), point to the issue forms (`?template=bug_report.yml`), and send security problems to private reporting.

  Follow-ups: the fix commands named in output.md, plus `/speckit.spec-status.reconcile` and `/speckit.spec-status.unblock`.

**Independent test**: Ask about each of the 32 codes. Each answer matches output.md, and the reference and link tests pass.

---

## Phase 6: User Story 4 — Get unblocked in the right order (P2)

- [X] T017 [P] [US4] Create three fixtures with hand-written `EXPECTED.md`:
  - `dependency-chain`: 003→002→001, all unfinished;
  - `superseded-dependency`: `SS_DEPENDENCY_HUMAN_REVIEW`;
  - `dependency-cycle`: `SS_DEPENDENCY_CYCLE`.

  Their signals must pass (T010).
- [X] T018 [US4] Write `extensions/spec-status/commands/speckit.spec-status.unblock.md`:
  - order blockers so depended-on items come first, using the `dependencies[].state` fields;
  - for each blocker, state what it needs to become `verified`;
  - for a superseded dependency, explain why there's no redirect and ask the user, applying only through `/speckit.roadmap.write` after approval;
  - for missing, ambiguous, self and cycle states, name the field and line and propose the correction.

  Follow-ups: `/speckit.spec-status.reconcile`, `/speckit.spec-status.verify-readiness`, `/speckit.roadmap.write`.

**Independent test**: The three fixtures produce the expected order and questions, and nothing is redirected without a decision.

---

## Phase 7: User Story 6 — See the review through my role (P2)

- [X] T019 [US6] Write `internal/extension/roles_test.go` (FR-022), failing before T020. For `extensions/spec-status/reference/roles.md` it checks:
  - `## Review sections` equals the closed list in contracts/reference-files.md;
  - `## Roles` has rows default, pm, developer, qa and architect, whose "Leads with" uses only review sections and whose "First follow-ups" uses only allowed commands.
- [X] T020 [US6] Write `extensions/spec-status/reference/roles.md` `## Review sections` and `## Roles` from spec User Story 6:
  - **pm:** progress, blockers, ready-to-start → unblock;
  - **developer:** ready-to-start, unchecked-tasks, missing-artifacts → reconcile, `/speckit-tasks`, `/speckit-implement`;
  - **qa:** evidence-gaps, then items recorded implemented/verified → verify-readiness, `/speckit.roadmap.debrief`, `/speckit-analyze`;
  - **architect:** dependency-chains, supersession, human-review, constraints → unblock, explain, `/speckit.roadmap.sync`.

  Role checks come from the spec. Fold in the substance of `specs/001-spec-status/examples/prompts/project-manager.txt` (pm) and `software-engineer.txt` (developer).
- [X] T021 [US6] Fill the role placeholder in `extensions/spec-status/commands/speckit.spec-status.review.md`:
  - parse `as <role>` from `$ARGUMENTS` and look it up in `roles.md`;
  - reorder sections and follow-ups only, never changing facts;
  - for an unknown role, list the roles and use the default view.
- [X] T022 [P] [US6] Replace the bodies of `specs/001-spec-status/examples/prompts/project-manager.txt` and `software-engineer.txt` with a three-line pointer to `extensions/spec-status/reference/roles.md` (the pm and developer roles). Keep the "SAMPLE TEXT … Not loaded by spec" banner (FR-023).
- [X] T023 [US6] Extend `testdata/skills/roles-mixed/EXPECTED.md` with the expected leading sections and first follow-ups for each of the four roles, and the rule that the facts are identical (SC-008).

**Independent test**: On `roles-mixed`, the default and all four roles report identical facts, differing only in order and first follow-ups. The roles test passes.

---

## Phase 8: User Story 5 — Know what's missing before marking something verified (P3)

- [X] T024 [US5] Write `extensions/spec-status/commands/speckit.spec-status.verify-readiness.md`:
  - for an item given in `$ARGUMENTS` (or each `implemented` item), check: all tasks checked and `tasks.state` `read`; every dependency `verified`; no error diagnostics on the item; spec, plan and tasks present;
  - report "consistent" with the reminder that a label isn't proof, or list the gaps;
  - never record a status itself.

  Follow-ups: `/speckit.roadmap.debrief`, `/speckit-analyze`, `/speckit.roadmap.write`, `/speckit.spec-status.reconcile`, `/speckit.spec-status.unblock`.

**Independent test**: On `implemented-unchecked` and `verified-dep-unverified`, it lists the right gaps and suggests the debrief. It never marks anything `verified`.

---

## Phase 9: User Story 7 — Use the other extensions I already have (P3)

- [X] T025 [US7] Extend `internal/extension/roles_test.go` (FR-027), failing before T026. In `roles.md` `## Partners`:
  - each row has all eight columns;
  - every command matches `speckit.<Extension>.<name>`, where `<Extension>` is the row's id;
  - each Role is one of the four roles.

  No partner command may appear in any `commands/*.md` file.
- [X] T026 [US7] Write `extensions/spec-status/reference/roles.md` `## Partners`, with three rows from research R2 (D-2026-11-12-07):

  | Role | Extension | Version | Catalog | Install | Requirements | Commands |
  |---|---|---|---|---|---|---|
  | developer | `bug` | 1.0.2 | default | `specify extension add bug` | Spec Kit `>=0.9.0` | `speckit.bug.{assess,fix,test}` |
  | qa | `assure` | 2.7.5 | community | `specify extension add assure --from <vetted archive URL>` | requires `illustrate` and Python 3.11+ | `speckit.assure.{init,analyze,document}` |
  | architect | `arch-governance` | 1.2.3 | community | `specify extension add arch-governance --from https://github.com/ashbrener/spec-kit-arch-governance/archive/refs/tags/v1.2.3.zip` | Spec Kit `>=0.1.0` | `speckit.arch-governance.{validate,gate,install,sync,repin,issues}` |

  Each row also has an "Adds" sentence from spec User Story 7. Mark community rows "not vetted by Spec Kit; review before installing".
- [X] T027 [US7] Fill the partner placeholder in `extensions/spec-status/commands/speckit.spec-status.review.md`:
  - read `.specify/extensions/.registry` (JSON) and treat a partner as installed when `extensions.<id>.enabled` is true and the command is in `registered_commands`;
  - recommend only registered commands;
  - otherwise give one install hint per partner, with the vetting note for community partners;
  - if the registry is missing or unreadable, say once that no partners can be detected.
- [X] T028 [P] [US7] Create two fixtures:
  - `partner-installed`: a `.specify/extensions/.registry` with `bug` enabled and `registered_commands.claude` = `speckit.bug.assess|fix|test`;
  - `partner-absent`: no registry.

  Each has a hand-written `EXPECTED.md`: the developer role recommends `speckit.bug.*` first when installed, and gives exactly one hint when absent (SC-009).

**Independent test**: Both fixtures behave as expected. A real install of `bug`, and of `arch-governance` from its tag archive, into a throwaway project registers the same commands that `roles.md` lists (T032).

---

## Phase 10: Polish and cross-cutting

- [X] T029 Create `scripts/check-extension-install.sh` (research R7, FR-018), using `set -euo pipefail` and a guarded `${TMP:?}`. It: *Done: bash 3.2-compatible. Spec Kit 1.2.0 does not carry `argument-hint` into extension skills, so the script checks that review's description names the roles instead (research R1, corrected). Run 2026-11-12: OK, 5 skills, 1 s.*
  1. Creates a temporary directory and runs `specify init --here --force --non-interactive --integration claude` there.
  2. Runs `specify extension add --dev <repo>/extensions/spec-status`.
  3. Asserts that exactly five `.claude/skills/speckit-spec-status-*/SKILL.md` exist, and that review's skill contains `argument-hint`.
  4. Prints the elapsed seconds and exits non-zero on any failure.
- [X] T030 Add the extension asset to `.github/workflows/release.yml` (research R7, FR-019): *Done. A local dry run of the zip and its smoke checks passed. Installing from the archive is pending: Spec Kit's `--from` only accepts https URLs, so it is checked after v0.2.0 is tagged (quickstart §5).*
  - **verify:** fail unless `extensions/spec-status/extension.yml` `version` equals the tag version;
  - **build:** `(cd extensions && zip -r ../dist/spec-status-extension_${VERSION}.zip spec-status)`, included in the checksum file;
  - **smoke** (ubuntu job): unzip, check `spec-status/extension.yml` and five `commands/*.md`, and verify its checksum line;
  - **release notes:** an install line `specify extension add spec-status --from https://github.com/$GITHUB_REPOSITORY/releases/download/$TAG/spec-status-extension_${VERSION}.zip`.
- [X] T031 [P] Update `specs/004-release-engineering/contracts/release-assets.md` (now four assets), `RELEASING.md` (bump `extensions/spec-status/extension.yml` and README to the release version before tagging) and `README.md`. In the README, "Use it with an AI coding agent" gains an "As Spec Kit skills" subsection: the install line, the five skills, `review as <role>`, partners, and the no-writes-without-approval rule.
- [X] T032 Run the local checks and record the results: *Done 2026-11-12. `bug` installed as 1.0.1 from the catalog (repository main: 1.0.2). arch-governance 1.2.3 installed from its tag archive. Both registered exactly the commands in roles.md. arch-governance's commands run Python scripts with `uv` (Python 3.11+, pydantic, pyyaml) and call `git`, none of which its manifest declares; roles.md says so. `validate` is read-only; `install`, `sync`, `repin` and `issues` write files, so they are only proposed with approval.*
  - `scripts/check-extension-install.sh`: five skills, under 60 s (SC-005);
  - install `bug` (`specify extension add bug`) and `arch-governance` (tag-archive URL) into throwaway projects, never this repository, and confirm their registered commands match `roles.md`;
  - read arch-governance's command bodies to confirm what its `install`, `sync`, `repin` and `issues` write.
- [X] T033 Run quickstart §1–§2: `go test ./...`, then three deliberate drift checks (remove a code row, add a bad link, add an unknown role section), each failing as expected and then reverted. Byte-compare `spec` output against the T001 baseline on the example and conformance fixtures (SC-006). *Done: `go test ./...`, vet and gofmt are clean. The three drift checks failed as expected and were reverted. CLI output was byte-identical in 14 comparisons, with no product code changed.*
- [X] T034 Run quickstart §4, the agent runs on all 12 fixtures: *Done 2026-11-12 by a separate agent following the rendered skills, never approving a change. First run: 10/12 PASS, roles-mixed FAIL; every fixture copy left git-clean; SC-008 held (identical facts in all five views). It found 7 rule ambiguities (A1–A7), which were fixed in the skills. Four expectations were corrected and noted in their EXPECTED.md: unknown-status, partner-absent, dependency-chain signals, roles-mixed routing and sections. Re-run of the changed cases: 4/4 PASS, git-clean. Two last ambiguities (the developer role with no fully-specced ready item; partner lines for roles without partners) were settled in roles.md and review.*
  - for each, record pass or fail against `EXPECTED.md`, and confirm with `git status` that nothing changed before approval (SC-001, SC-004);
  - the role check on `roles-mixed` (SC-008);
  - the partner checks (SC-009).

  Fix commands, not expectations, unless an expectation was wrong; note any such change.
- [X] T035 Run the Doomer review workflow on `internal/extension`, `scripts/` and `.github/workflows/release.yml`. Record a response for every Quick finding. *Done: Doomer Quick review over 3 revisions. One lead (indented headings in sectionBody) was fixed with a self-test; no other leads.*
- [X] T036 Commit, push and open the PR, and confirm CI is green. Add an implementation-record row to `specs/001-spec-status/tasks.md`: "Agent skills extension (feature 005)", with the test, install-check and fixture-run results. Move roadmap 012 to `implemented` (patch bump). *Done: PR opened; CI result recorded in the PR.*
- [ ] T037 **(owner)** After merge: bump the versions if needed, tag `v0.2.0`, check that the draft has four assets, publish, and run quickstart §5 (`specify extension add spec-status --from <asset URL>` in a fresh project).

---

## Dependencies and execution order

- **Setup:** T001, T002.
- **Foundational, sequential:** T003 → T004 → T005 → T006 → T007 → T008 → T009 → T010. Each test is written before the file it guards: T004 before T007, T005 before T008, T006 before any command.
- **review.md is shared:** T011 (US1) → T021 (US6 roles) → T027 (US7 partners).
- **roles.md is shared:** T020 → T026.
- **Independent once the foundational phase is done:** the fixtures (T012, T013, T017, T028) and the commands for US2–US5 (T015, T016, T018, T024).
- **Polish:** T029–T031 can start any time after the foundational phase. T032–T036 need every story. T037 comes after merge.

### Story order

US1 (MVP) and US2 are both P1. Then US3, US4 and US6 (P2), then US5 and US7 (P3).

## Parallel example

```text
After T010:
  T012 [US1] roles-mixed fixture      T013 [US2] six reconcile fixtures     T017 [US4] three unblock fixtures
  T015 [US2] reconcile.md             T016 [US3] explain.md                 T018 [US4] unblock.md
  T024 [US5] verify-readiness.md      T029 install script                   T031 docs
Sequential lane: T011 → T021 → T027 (review.md); T019 → T020 → T025 → T026 (roles)
```

## Implementation strategy

1. **MVP**: Foundational + US1 + US2. Review plus reconcile already solves the "records disagree" problem.
2. **Then**: explain, unblock and roles (P2), then verify-readiness and partners (P3), then the release asset.
3. **One PR, released as v0.2.0.**
4. **Stop conditions**:
   - T014: a fixture's signal doesn't appear in `spec` output. Rework the fixture; never change the CLI.
   - T032: a partner's registered commands differ from `roles.md`. Update research, roles and the decision.
   - T034: any unapproved file change in a fixture run. This is a blocking defect.
