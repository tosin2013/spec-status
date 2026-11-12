# Research: Dependency update automation

Feature 003 · roadmap item 009 · 2026-11-11. Facts were checked on this date with `gh api`, the Go module proxy, `go.dev/dl`, and local runs; the commands are listed at the end.

## R1 — Minimum Go version (was out of scope; owner decision D-2026-11-11-05)

- **Finding**: `go.mod` declares `go 1.22`, and Go 1.22 no longer receives fixes. The supported Go release lines on 2026-11-11 are **1.27** (go1.27.4) and **1.26** (go1.26.11).
  - Run under go1.24.3, `govulncheck` reports a reachable standard-library vulnerability (GO-2025-3750) and exits 3. Go 1.22 can only be worse.
  - Run under go1.26.9, it reports 0 reachable vulnerabilities.
  - The `golang.org/x/sys` fix (v0.44.0 and later) declares `go 1.25.0`, so any security update of x/sys raises the module's minimum Go version.
  - All tests pass under go1.26.11.
- **Decision**: `go.mod` declares `go 1.26.0`. CI tests both supported lines: `1.26.x` for the declared minimum and `stable` for the current release. The owner chose this on 2026-11-11.
- **Rationale**: A check that fails every pull request on stdlib advisories that only a toolchain move can fix would be ignored, and the first Dependabot security fix (x/sys) could not merge without this move anyway. 1.26 is the oldest supported line, so users on either supported line can build.
- **Alternatives considered**:
  - 1.27 minimum: drops 1.26 users for no gain.
  - Keep 1.22 with a non-blocking scan: leaves a red check on every PR, and blocks the x/sys fix.
- **Consequence for 004**: the toolchain item there narrows to which exact patch builds the release archives. The minimum is decided.

## R2 — Which Go versions CI runs, and how they are selected

- **Decision**: a matrix with `go-version: ['1.26.x', 'stable']` in `actions/setup-go`, plus one step that fails if `go.mod`'s `go` line does not start with `go 1.26`.
- **Rationale**:
  - `1.26.x` resolves to the newest 1.26 patch, so the minimum line is tested with its fixes.
  - `stable` follows the newest release with no manual edits.
  - Using `go-version-file: go.mod` instead would read `1.26.0` as an exact version and install an unpatched toolchain.
  - The consistency step keeps the hard-coded `1.26.x` honest when the minimum moves, which needs a decision under D-2026-11-11-04.
- **Note**: When Go 1.28 ships, 1.26 leaves support. The `1.26.x` job's scan will then start reporting stdlib advisories. That is the intended signal to propose raising the minimum.

## R3 — Runners and platform truthfulness

- **Decision**: `ubuntu-latest` and `macos-latest`. Each job asserts its platform (`go env GOOS GOARCH` must be `linux amd64` or `darwin arm64`) before testing.
- **Rationale**: The `-latest` labels always exist, and GitHub's `macos-latest` has been Apple-silicon arm64 for years. The assertion keeps the spec's platform claim true even if a label's meaning changes. If it does change, the job fails instead of silently testing another platform. A pinned image label such as `macos-15` would eventually be retired, and the job would fail.
- **Not claimed**: Passing jobs are CI evidence only, not DC-07's "proven on a real runner" release evidence. Feature 004 defines that evidence (runtime smoke tests of the archives).

## R4 — Vulnerability scan

- **Decision**: `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...` in every matrix job.
- **Rationale**:
  - v1.9.0 (2026-10-11) is the latest release, and it needs Go ≥1.26.0, which both matrix entries satisfy.
  - `go run pkg@version` builds the tool outside the project module and does not touch `go.mod` or `go.sum`. So the scanner never becomes a project dependency, and FR-012 and the import allowlist are unaffected.
  - govulncheck reports only advisories whose vulnerable symbols are reachable, and exits 3 when any are found. Exit code 3 fails the job (FR-006).
  - Scanning per Go version and per platform catches advisories specific to a stdlib version or to GOOS.
- **Updating it**: Dependabot cannot see a version inside a `run:` line. The pinned `@v1.9.0` is bumped by hand. The bump is recorded in the PR, and a newer scanner version is not a project dependency.
- **Alternative considered**: `golang/govulncheck-action`. It adds another third-party action and runs its own setup-go, which would duplicate the matrix's toolchain.
- **Expected first results** (local run, go1.26.9): 0 reachable. Informational only: GO-2026-5024 in `golang.org/x/sys@v0.30.0`, which the code does not call, fixed in v0.44.0. A repository alert, and then a Dependabot security PR, are likely for it, because alerts are per module, not per call.

## R5 — Workflow components and pinning

- **Decision**: Pin every action to a full commit SHA, with the release tag in a trailing comment:
  - `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`
  - `actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e # v7.0.0`
  - Both are the latest releases on 2026-11-11, resolved from their tag refs (both refs point directly to commits).
- **Rationale**: FR-007 requires exact versions, and a SHA cannot be moved after the fact the way a tag can. Dependabot's `github-actions` ecosystem updates SHA pins and keeps the version comment current.
- **Permissions**: workflow-level `permissions: contents: read`, and no secrets are referenced. Dependabot-triggered `pull_request` runs get a read-only token and no repository secrets by default.
- **Other settings**:
  - `persist-credentials: false` on checkout, so the token is not left in the git config.
  - setup-go's module cache stays on (default). It is keyed on `go.sum` and holds only downloaded modules.

## R6 — Dependabot configuration

- **Decision**: `.github/dependabot.yml`, version 2, with two `updates` entries:
  - `gomod` at `/`: weekly on Monday. One group `go-minor-patch` with patterns `["*"]` and update types `minor` and `patch`. Major updates are left out of the group, so each opens its own PR.
  - `github-actions` at `/`: weekly on Monday, one group `actions` covering all updates.
- **Rationale**:
  - Groups apply to version updates by default.
  - Security updates are configured separately: they were enabled in repository settings on 2026-11-11, arrive as soon as an advisory is published, and are not grouped. That gives the owner's "security fixes as soon as an advisory is published, routine updates weekly" split.
  - The default limit of 5 open pull requests per ecosystem is enough for 6 direct modules.
- **Not chosen**:
  - Auto-merge: out of scope.
  - `ignore` rules: none needed yet.
  - Labels and assignees: not requested.

## R7 — `go.mod` change mechanics

- **Decision**: `go mod edit -go=1.26.0`, then `go mod tidy` under go1.26.11. Review the diff: if tidy adds any module not already in `go.mod`, stop and ask for a decision (D-2026-11-11-04).
- **Rationale**: Raising the `go` line can make tidy rewrite the requirement blocks, for example by pruning or adding indirect requirements. The merge rule applies to this feature's own change too.

## Commands used

```sh
gh api repos/actions/checkout/releases/latest --jq .tag_name           # v7.0.1
gh api repos/actions/checkout/git/ref/tags/v7.0.1                      # 3d3c42e5…, type commit
gh api repos/actions/setup-go/releases/latest --jq .tag_name           # v7.0.0
gh api repos/actions/setup-go/git/ref/tags/v7.0.0                      # b7ad1dad…, type commit
curl -fsSL https://proxy.golang.org/golang.org/x/vuln/@latest          # v1.9.0
curl -fsSL https://proxy.golang.org/golang.org/x/vuln/@v/v1.9.0.mod    # go 1.26.0
curl -fsSL https://proxy.golang.org/golang.org/x/sys/@v/v0.44.0.mod    # go 1.25.0
curl -fsSL 'https://go.dev/dl/?mode=json'                              # go1.27.4, go1.26.11
GOTOOLCHAIN=local    govulncheck ./...   # go1.24.3: 1 reachable stdlib vuln, exit 3
GOTOOLCHAIN=go1.26.9 govulncheck ./...   # 0 reachable; x/sys GO-2026-5024 not called
GOTOOLCHAIN=go1.26.11 CGO_ENABLED=0 go test -count=1 ./...   # all packages ok
```
