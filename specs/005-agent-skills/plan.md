# Implementation Plan: Agent skills extension

**Branch**: `005-plan` | **Date**: 2026-11-12 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/005-agent-skills/spec.md` (decisions D-2026-11-12-04, -05, -06)

## Summary

This feature adds a Spec Kit extension, `extensions/spec-status/`, with five command files that Spec Kit renders as agent skills:
- `review`: the entry point, with roles `pm`, `developer`, `qa` and `architect`;
- `explain`;
- `reconcile`;
- `unblock`;
- `verify-readiness`.

Three shared reference files hold, in one place each:
- the output vocabulary;
- the roles and partner extensions;
- the shared rules.

A test-only Go package keeps the reference in step with the contracts, and checks the manifest, command structure, links and roles. Fixtures under `testdata/skills/` give each mismatch, role and partner case hand-written expectations, with a signal test that runs `status.Evaluate`. The release workflow ships the extension as a fourth asset. The `spec` binary does not change.

## Technical Context

**Language/Version**: Markdown command files and YAML for the extension manifest. Go 1.26+ for the test-only checks; the standard library is enough.

**Primary Dependencies**:
- Spec Kit 1.2.0 extension format, installed by the user, not by this repository.
- No new Go modules.
- The install test script uses the `specify` CLI locally.

**Storage**: N/A. Partner detection reads the project's `.specify/extensions/.registry`, which is JSON.

**Testing**:
- **CI:** `go test ./internal/extension` covers the reference, manifest, structure, links, roles, version and fixture signals. `go test ./...` is otherwise unchanged.
- **Local and recorded:** `scripts/check-extension-install.sh`, plus agent runs against the fixtures.

**Target Platform**: any agent that Spec Kit renders extension commands for. The install test uses the Claude integration.

**Project Type**: a Spec Kit extension (documentation-as-code) inside the CLI repository.

**Performance Goals**: install in under 1 minute (SC-005).

**Constraints**:
- The CLI is unchanged: its JSON, exit codes and README samples stay byte-identical.
- No write without per-change approval, and writes only through Spec Kit commands (Art. VIII, C-06).
- Partners are detected, never required.
- The test uses the standard library only.

**Scale/Scope**:
- about 10 extension files;
- 12 fixtures;
- 1 test package;
- 1 script;
- edits to `release.yml`, README, 2 sample prompts, and the feature 004 asset contract.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Article | Check | Result |
|---|---|---|
| I Truthful evidence | The rules restate that `verified` is a label, that a checkbox is not a test, and that exit 0 is not "done". Proposals report proposed, applied or confirmed. Partner limits (adrkit is incompatible with Spec Kit 1.2.0) are stated. | Pass |
| II Standalone, deterministic; allowlist | The `spec` binary is unchanged. The extension is development tooling (CLAUDE.md). The test package has only `_test.go` files and uses the standard library. | Pass |
| III Narrow, safe input | Unaffected. The skills consume output; they don't parse roadmaps themselves. | Pass |
| IV History is not false completion | reconcile never asks superseded items to match tasks, and never retargets a superseded dependency. | Pass |
| V Compatibility | Schema 1.0 is unchanged. The skills warn on a different `schema_version`. | Pass |
| VI No inferred quarters | The rules forbid inferring dates. | Pass |
| VII Distribution claims need evidence | The extension asset is smoke-checked in the release job. Partner compatibility is stated as researched. | Pass |
| VIII Explicit mutation authority | Changes go through Spec Kit commands only, after per-change approval, followed by a re-read. | Pass |
| Governance | No amendment is needed. | Pass |

**Post-design re-check:** passes. The contracts keep every write behind approval, and keep the reference test-bound to the contracts.

## Project Structure

### Documentation (this feature)

```text
specs/005-agent-skills/
├── spec.md, plan.md, research.md, data-model.md, quickstart.md
├── contracts/
│   ├── extension-manifest.md   # extension.yml shape and rules
│   ├── command-format.md       # required frontmatter/sections, links, approval protocol
│   ├── reference-files.md      # output.md, roles.md, rules.md table formats
│   └── fixtures.md             # testdata/skills layout and EXPECTED.md format
├── checklists/requirements.md
└── tasks.md                    # /speckit-tasks
```

### Source (repository root)

```text
extensions/spec-status/
├── extension.yml
├── README.md
├── LICENSE                       # copy of the repository's Apache-2.0
├── commands/
│   ├── speckit.spec-status.review.md
│   ├── speckit.spec-status.explain.md
│   ├── speckit.spec-status.reconcile.md
│   ├── speckit.spec-status.unblock.md
│   └── speckit.spec-status.verify-readiness.md
└── reference/
    ├── output.md                 # codes, severities, JSON fields, states, statuses
    ├── roles.md                  # roles and partner extensions
    └── rules.md                  # shared rules and the approval protocol
internal/extension/                # test-only package (*_test.go)
testdata/skills/<case>/            # 12 fixtures with EXPECTED.md
scripts/check-extension-install.sh
.github/workflows/release.yml      # + extension zip, version check, smoke
README.md, specs/001-spec-status/examples/prompts/*.txt, specs/004-release-engineering/contracts/release-assets.md, RELEASING.md
```

## Complexity Tracking

No constitution violations.
