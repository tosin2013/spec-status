# Contract: CI workflow (`.github/workflows/ci.yml`)

A conforming workflow meets every **MUST** below. Each one maps to a spec requirement.

## Triggers

- MUST run on `pull_request` targeting `main` and on `push` to `main` (FR-004).
- MUST NOT run on tags or releases, and MUST NOT use `pull_request_target` (FR-008, FR-007).

## Permissions and secrets

- MUST declare `permissions: contents: read` at workflow level, with no job-level widening (FR-007).
- MUST NOT reference `secrets.*` (FR-007).
- Checkout MUST set `persist-credentials: false`.

## Pinning

- Every `uses:` MUST reference a full 40-character commit SHA, with the release tag in a trailing comment (FR-007).
- The scanner MUST be invoked with an exact version (`@vX.Y.Z`), never `@latest`.

## Matrix

- `os`: `ubuntu-latest`, `macos-latest`. `go`: `1.26.x`, `stable`. `fail-fast: false` (FR-005, Acceptance 2.3).

## Steps, in order, per job

1. Checkout.
2. setup-go with `go-version: ${{ matrix.go }}`.
3. **Platform assertion**: `go env GOOS GOARCH` MUST equal `linux amd64` on ubuntu and `darwin arm64` on macOS, otherwise fail (R3).
4. **Minimum assertion**: the `go` line in `go.mod` MUST start with `go 1.26`, otherwise fail (R2).
5. `go vet ./...` MUST pass (FR-006).
6. `CGO_ENABLED=0 go test ./...` MUST pass (FR-005).
7. `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...` MUST exit 0. Any reachable advisory, or a database fetch failure, fails the job (FR-006, edge cases).

## Must not

- Build archives, compute release checksums, create tags or releases, upload artifacts, comment on PRs, or push (FR-008).
- Skip or `continue-on-error` any of steps 3–7.

## Observable result

- Four checks appear on each pull request, one per `os` × `go` cell, each pass or fail.
- Whether passing checks are *required* for merge is a branch-protection setting. It is outside this contract; quickstart notes it as an optional owner step.
