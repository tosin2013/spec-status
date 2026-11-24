# Quickstart: validating more release platforms

**Before the merge:**
- `actionlint .github/workflows/ci.yml .github/workflows/release.yml` is clean.
- `CGO_ENABLED=0 go test ./...` passes, including the README samples, unpinned install lines and extension version checks.
- For each target, `GOOS=… GOARCH=… CGO_ENABLED=0 go build ./cmd/spec` succeeds (local cross-compile).
- PR CI shows `test (ubuntu-24.04-arm, …)` and `test (macos-15-intel, …)` passing, each printing the platform assert.

**Release v0.4.1:**
- The run shows `smoke (ubuntu-latest)`, `smoke (ubuntu-24.04-arm)`, `smoke (macos-latest)` and `smoke (macos-15-intel)` all green, each with the expected `target:`.
- The draft has 4 archives, the zip and the checksums, and its notes have 4 smoke-tested lines with links.

**After publishing:**
- The latest lookup gives 0.4.1, and all 4 checksums are OK.
- darwin/amd64 runs locally with `arch -x86_64 ./spec version`.
- linux/arm64 runs with `docker run --rm --platform linux/arm64 -v "$PWD":/w alpine /w/spec version` if Docker is available.
