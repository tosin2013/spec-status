# Research: Agent skills extension

Feature 005 · roadmap item 012 · 2026-11-12. The sources were the installed Spec Kit CLI (`specify` 1.2.0, `specify_cli` package), the roadmap extension in `.specify/extensions/roadmap/`, `specify extension search` and `specify extension info`, and the partner repositories via `gh api`. The commands are at the end.

## R1 — Extension package format (FR-001, FR-002)

- **Decision**: Model `extensions/spec-status/` on `.specify/extensions/roadmap/`:
  - an `extension.yml` with `schema_version: '1.0'`, `extension{id,name,version,description,author,repository,license,homepage}`, `requires`, `provides.commands[{name,file,description}]`, `hooks`, and `tags`;
  - `commands/*.md` files, each starting with frontmatter (`description`, optional `argument-hint`).
- **Rendering**: Spec Kit renders each command into a skill per agent integration (`specify_cli/agents.py`). For Claude, `argument-hint` from the command frontmatter is copied into the skill (`apply_argument_hint`), so `review` can advertise `[as pm|developer|qa|architect]`. Commands refer to each other in text, as the roadmap commands do with `/speckit.roadmap.write`.
- **Install sources**:
  - `specify extension add --dev <dir>`;
  - `--from <url>`, accepting `.zip`, `.tar.gz` or `.tgz`, with `extension.yml` at the root of the archive or inside one top-level directory (`specify_cli/extensions/__init__.py`, `install_from_archive`).
- **Installed location**: `.specify/extensions/spec-status/`. Commands read shared files from there by path, the same way roadmap commands read `.specify/extensions/roadmap/templates/…`.
- **`requires.tools`**: Assure's manifest uses it (`tools: [{name, required, version}]`), but Spec Kit 1.2.0 does not appear to enforce it: the install code has no tool lookup. We list `spec` there as documentation. FR-004 is still met by each command checking `spec version` itself.
- **Alternatives considered**: plain `SKILL.md` files, which the owner rejected (D-2026-11-12-04); publishing to a catalog (out of scope).

## R2 — Partner extensions (FR-024 to FR-027, D-2026-11-12-06)

Each manifest was read from its repository on 2026-11-12.

| Role | Partner (id, version, catalog) | Commands | Requirements | Notes |
|---|---|---|---|---|
| architect | `arch-governance` 1.2.3 (catalog lists 1.2.2), community (discovery only), MIT, `ashbrener/spec-kit-arch-governance` | `speckit.arch-governance.validate` (read-only, fail-closed citation validator), `.gate` (before_implement; advisory or blocking), `.install`, `.sync`, `.repin`, `.issues` (the last three dry-run by default); hooks `after_specify` and `after_plan` → validate, `before_implement` → gate | Spec Kit `>=0.1.0`; no tools or extensions | Compatible with Spec Kit 1.2.0. No release assets; install from the tag archive `https://github.com/ashbrener/spec-kit-arch-governance/archive/refs/tags/v1.2.3.zip` (one nested folder, which Spec Kit accepts, R1) |
| qa | `assure` 2.7.5, community (discovery only), in `samykabu/sanduq` | `speckit.assure.init`, `speckit.assure.analyze`, `speckit.assure.document`; hooks `after_tasks`/`before_implement` → analyze, `after_implement` → document | Spec Kit `>=0.2.0`; requires extension `illustrate >=2.0.0,<3.0.0`; tools `python >=3.11` (required), `git` (optional) | Has a required extension dependency |
| developer | `bug` 1.0.2, default catalog (github/spec-kit) | `speckit.bug.assess`, `speckit.bug.fix`, `speckit.bug.test` | Spec Kit `>=0.9.0` | Installable: `specify extension add bug` |

- **Decision**: The role file records these exact commands, versions and requirements. A partner command is recommended only if the project's `.specify/extensions/.registry` lists it under `registered_commands`. Install hints quote each partner's requirements:
  - arch-governance's hint gives the tag-archive install command;
  - Assure's hint says it needs the `illustrate` extension and Python 3.11+;
  - community partners also get the "not vetted by Spec Kit" note.
- **Consequence**: in this project (Spec Kit 1.2.0), Bug Triage and arch-governance can be installed and tested. Assure needs `illustrate` and Python 3.11+, so it is tested through a synthetic registry entry and its "not installed" hint.
- **adrkit dropped (D-2026-11-12-07)**: adrkit 0.1.5 was the first architect candidate. Its extension manifest on `main` still says `speckit_version: ">=0.13.0,<1.2.0"`, with a comment that widening past 1.1 requires re-verification under its ADR-0007. The project's newer CLI releases (up to v0.20.0, 2026-10-11) do not change that. It also needs the `adr` CLI. A catalog search for ADR extensions (names, descriptions and tags matching ADR, "decision record" or MADR) found only adrkit and arch-governance. Revisit adrkit when it supports Spec Kit 1.2 or later.

## R3 — Detecting installed extensions (FR-024)

- **Decision**: Read `.specify/extensions/.registry`, which is JSON. Find `extensions.<id>.enabled` and `extensions.<id>.registered_commands.<agent>[]`. A partner counts as installed when it is present and `enabled: true`, and a command counts when it appears in `registered_commands` for any agent.
- **Rationale**: Spec Kit maintains this file itself; this project's copy shows `roadmap` with its four commands. If the file is missing or unreadable, the command behaves as if no partners are installed and says so once (spec edge case).
- **Not used**: running `specify extension list` from the skill. That needs the `specify` CLI at review time, and reading the file is enough.

## R4 — Shared files and how the commands use them (FR-015, FR-021, FR-025)

- **Decision**: three files under `extensions/spec-status/reference/`:
  - `output.md`: the output vocabulary;
  - `roles.md`: roles and partners;
  - `rules.md`: the shared rules from FR-006, FR-007 and FR-009.

  Every command starts by reading `rules.md`, and reads the others as needed. All three are Markdown tables, so agents can read them and the Go test can parse them with the standard library (no YAML or JSON dependency, constitution Art. II).
- **Rationale**: one copy of the rules and vocabulary for all five commands (D-2026-11-12-05), and parseable without adding a module.
- **Alternatives considered**: copying the rules into each command, which drifts; a JSON or YAML role file, which agents read less easily and which needs a parser or a new module for the test.

## R5 — Keeping the reference in step with the contracts (FR-015, FR-016, FR-022, FR-027)

- **Decision**: a new test-only package `internal/extension`, written with only `_test.go` files and the standard library. It checks:
  1. **Codes**: every code in a `| SS_… |` row of `specs/001-spec-status/diagnostics.md` (32 today) appears in `output.md` with the same severity, and `output.md` has no extra codes.
  2. **Fields**: the first-column field names of the Root, Summary and Item tables in `json-contract.md` (comma-separated cells split) all appear in `output.md`. So do the dependency states (`satisfied, unfinished, missing, ambiguous, self, cycle, human_review`), the task states (`read, missing, skipped, not_linked`) and the nine status tokens.
  3. **Manifest**: every `file:` under `provides.commands` in `extension.yml` exists, and every `commands/*.md` is declared. The manifest is read with a line regex over the fixed format, not with a YAML parser.
  4. **Command structure**: each command has frontmatter with `description`, plus the sections `## User Input`, `## Goal`, `## Operating Constraints` and `## Outline`, and it references `reference/rules.md`.
  5. **Links**: every `/speckit…` reference in a command or in `roles.md` is either in the allowed set (FR-008: the five spec-status commands, the six core commands, the four roadmap commands) or inside a role's partner list. Partner commands must match `speckit.<id>.<name>`, and their `<id>` must be a partner declared in `roles.md`.
  6. **Roles**: every section a role names exists in review's section list, and every role follow-up is allowed.
  7. **Version**: `extension.yml`'s `extension.version` equals the version in the extension `README.md`'s heading. The release job checks it against the tag (R7).
- **Rationale**: the CLI's import allowlist test is unaffected, because test files are not product code. The repository already reads Markdown contracts this way in `internal/conformance`.

## R6 — Mismatch fixtures (FR-017, SC-001, SC-004, SC-008, SC-009)

- **Decision**:
  - **Fixtures**: `testdata/skills/<case>/` holds a small Spec Kit project (`.specify/memory/roadmap.md`, `specs/...`) and an `EXPECTED.md` written by hand. `EXPECTED.md` gives the finding, the evidence (item and file:line), the proposed correction, the applying command, and for role cases the expected first follow-ups.
  - **Cases**:
    - reconcile, 6: `tasks-done-status-planned`, `implemented-unchecked`, `verified-tasks-missing`, `verified-dep-unverified`, `unknown-status`, `checked-target-missing`;
    - unblock, 3: `dependency-chain`, `superseded-dependency`, `dependency-cycle`;
    - roles, 1: `roles-mixed`, one roadmap with signals for all four roles;
    - partners, 2: `partner-installed`, which has a `.specify/extensions/.registry` registering `speckit.bug.*`, and `partner-absent`.
  - **Signal test**: a Go test runs `status.Evaluate` on each fixture. It asserts that the raw signals the skills rely on are present: the status, tasks counts and state, the diagnostic codes, and the dependency states. This means a future CLI change can't quietly invalidate a fixture.
  - **Behavior**: agent behavior is checked by hand or by an agent run against `EXPECTED.md`. Results are recorded in the implementation record, not in CI.
- **Rationale**: agents can't run in CI deterministically. The signal test keeps the fixtures honest, and the expectations are written by people, not generated.

## R7 — Release asset and version (FR-019)

- **Decision**: In `.github/workflows/release.yml`:
  - **Build job**: zip `extensions/spec-status/` as `spec-status-extension_<version>.zip`, with one top-level directory `spec-status/`, and add it to the checksum file.
  - **Verify job**: fail unless `extensions/spec-status/extension.yml`'s version equals the tag version.
  - **Smoke job**: unzip the asset and check that `extension.yml` and five command files are present.
  - **Release notes**: add an install line: `specify extension add spec-status --from <asset URL>`.
  - `contracts/release-assets.md` in feature 004 gains the fourth asset.
- **Rationale**: `--from` accepts a zip with one nested directory (R1). Checking the version against the tag stops a release from shipping a stale extension version.
- **Install test**: `specify init` uses templates bundled in the CLI (`core_pack`), so a throwaway project can be created locally. A script `scripts/check-extension-install.sh` (development tooling) creates a temporary project with `specify init --here --force --non-interactive --integration claude`, runs `specify extension add --dev <repo>/extensions/spec-status`, and checks that five `speckit-spec-status-*` skills are rendered. It runs locally and its output is recorded. It is not a CI job in v0.2.0, so CI gains no Python or pip dependency (feature 003's supply-chain stance).

## R8 — Hooks (FR-003)

- **Decision**: one optional hook, `after_implement` → `speckit.spec-status.reconcile`, with a prompt such as "Check recorded statuses against tasks now?". No other hooks.
- **Rationale**: the roadmap's `after_implement` (debrief) and any partner hooks are listed alongside, and the user chooses. Spec Kit treats every optional hook as a suggestion.

## R9 — Sample prompts (FR-023)

- **Decision**: `roles.md` takes over the substance of `project-manager.txt` (pm) and `software-engineer.txt` (developer). The two `.txt` files become a three-line pointer to the extension, keeping their "SAMPLE TEXT … Not loaded by spec" banner. `companion.md` keeps linking to them, unchanged in meaning.

## Commands used

```sh
specify --version                                   # 1.2.0
specify extension search adr|architecture|qa|test   # adrkit, Assure, Bug Triage
specify extension info adrkit                       # Spec Kit >=0.13.0,<1.2.0; adr >=0.3.0 (dropped)
curl …/extensions/catalog{,.community}.json         # ADR filter: only adrkit and arch-governance
gh api repos/ashbrener/spec-kit-arch-governance/contents/extension.yml   # speckit.arch-governance.*; speckit >=0.1.0
gh api repos/samykabu/sanduq/contents/extensions/assure/extension.yml   # speckit.assure.{init,analyze,document}
gh api repos/github/spec-kit/contents/extensions/bug/extension.yml      # speckit.bug.{assess,fix,test}
gh api repos/mbeacom/adrkit/contents/packages/adapters/spec-kit/extension.yml  # speckit.adrkit.{context,check,draft}
grep install_from_archive / apply_argument_hint in specify_cli (uv tool install)  # archive formats, nesting, argument-hint
cat .specify/extensions/.registry                   # registered_commands per agent
```
