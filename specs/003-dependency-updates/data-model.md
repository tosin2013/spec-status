# Data model: Dependency update automation

This feature stores no data. The "entities" are repository configuration and one classification rule.

## Update configuration (`.github/dependabot.yml`)

| Field | Value | Source |
|---|---|---|
| version | 2 | Dependabot schema |
| updates[gomod].directory | `/` | single module at the repository root |
| updates[gomod].schedule | weekly, Monday | FR-001 |
| updates[gomod].groups.go-minor-patch | patterns `*`, update-types `minor`, `patch` | FR-001; majors stay ungrouped |
| updates[github-actions].directory | `/` | Dependabot reads `.github/workflows/` |
| updates[github-actions].schedule | weekly, Monday | FR-002 |
| updates[github-actions].groups.actions | patterns `*` | FR-002 |

Security updates are not set in this file. They are a repository setting, enabled on 2026-11-11 (FR-003), and they ignore the schedule and groups.

## Automatic check (`.github/workflows/ci.yml`)

| Field | Value | Source |
|---|---|---|
| triggers | `pull_request` to `main`; `push` to `main` | FR-004 |
| permissions | `contents: read` (workflow level) | FR-007 |
| matrix.os | `ubuntu-latest`, `macos-latest` | FR-005, R3 |
| matrix.go | `1.26.x`, `stable` | R1, R2 |
| fail-fast | false (each cell reports its own result) | Acceptance scenario 2.3 |
| steps | checkout (pinned) → setup-go (pinned) → assert platform → assert `go.mod` minimum → `go vet ./...` → `CGO_ENABLED=0 go test ./...` → govulncheck v1.9.0 | FR-005, FR-006, R2–R5 |

There are 4 jobs: 2 operating systems × 2 Go lines. One job passes only if every one of its steps passes.

## Dependency proposal → merge-rule outcome (D-2026-11-11-04)

The rule classifies a proposal by comparing the module paths in its `go.mod` against those on `main`:

| Condition on the proposal | Outcome |
|---|---|
| Module paths are unchanged (only versions differ), the `go` line is unchanged, and CI passes | **Merge after review**: the owner approves the PR |
| Any module path added (direct or `// indirect`) | **Decision needed**: record a D-entry naming each new module before merging |
| `go` line raised | **Decision needed** |
| CI fails | **Do not merge** in either case |

Removing a module or an indirect requirement does not need a decision. A workflow-component (actions) proposal is classified as merge after review once CI passes, because actions are not part of the binary.

The rule has no other states (SC-004).
