## What and why

<!-- One or two sentences: what changes, and the problem it solves. -->

- Roadmap item / spec: <!-- e.g. 004, specs/004-release-engineering; or "none" for a small fix -->
- Issue: <!-- Fixes #NNN, if any -->

## Kind of change

- [ ] Bug fix (wrong count, status, dependency, exit code, or crash)
- [ ] New behavior or flag (needs a recorded owner decision for contract changes)
- [ ] Docs only
- [ ] Tooling / CI / dependencies

## Rules this change must keep

See [CONTRIBUTING.md](../CONTRIBUTING.md) and `.specify/memory/constitution.md`.

- [ ] No runtime side effects: no network, subprocess, Git, file writes, caches, or telemetry
- [ ] No new dependency outside the allowlist (Cobra/pflag in `internal/cli`, Lip Gloss in `style.go`, standard library elsewhere)
- [ ] All file reads go through `internal/safeio`
- [ ] JSON output, exit codes, and diagnostics are unchanged, or the change is listed below with its decision ID
- [ ] Unsupported input still fails visibly. Nothing is guessed.
- [ ] Hand-counted fixtures (`testdata/conformance/*/expected.json`) were counted by hand, not regenerated from `spec` output

## Dependency updates (Dependabot or manual)

- [ ] Only versions of modules already in `go.mod` changed, and the `go` line is unchanged: merge after review
- [ ] A module was added or the `go` line was raised: a decision naming it is recorded in `specs/001-spec-status/decisions.md` (rule D-2026-11-11-04)

## Evidence

<!-- Say exactly what you ran and where. Do not claim more than that. -->

- [ ] `CGO_ENABLED=0 go test ./...` on <!-- OS/arch, Go version -->
- [ ] `go vet ./...`, `gofmt -l .` (empty), `go mod tidy -diff`
- [ ] Output compared before and after (text and `--json`) where behavior could change:
- [ ] CI is green on this PR

## Docs

- [ ] README (including Limits and the example output, if it changed)
- [ ] The relevant contract in `specs/`, and the implementation record in `specs/001-spec-status/tasks.md`
- [ ] Not needed

## Behavior changes

<!-- List every user-visible change: flags, exit codes, JSON fields, diagnostics, human text. Before v1.0.0 a minor release may break compatibility, and each break must be listed in the release notes. Write "none" if none. -->
