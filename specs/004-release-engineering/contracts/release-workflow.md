# Contract: Release workflow (`.github/workflows/release.yml`)

## Trigger and permissions

- `on: push: tags: ['v*']` only. No `workflow_dispatch`, `pull_request` or `pull_request_target`.
- Workflow level: `permissions: contents: read`. Only job `release` declares `contents: write`.
- No `secrets.*`. The `release` job uses `GH_TOKEN: ${{ github.token }}`.
- Every `uses:` is pinned to a 40-hex SHA with a version comment. No `continue-on-error`.

## Jobs (each `needs` the previous)

1. **verify** (`ubuntu-latest`): checkout with `fetch-depth: 0` and `persist-credentials: false`.
   - Fail unless `$GITHUB_REF_NAME` matches `^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-rc\.(0|[1-9][0-9]*))?$`.
   - Fail unless `git merge-base --is-ancestor "$GITHUB_SHA" origin/main`.
   - Fail if `gh release view "$GITHUB_REF_NAME"` succeeds, because a release already exists. This job's read-only token sees **published** releases only. An existing **draft** for the tag is caught by the `release` job's re-check, which has write access, so nothing is created or replaced either way. Observed in T020 (2026-11-12).
   - Outputs: `version` (the tag without `v`) and `prerelease` (`true` for `-rc.N`).
2. **ci**: `uses: ./.github/workflows/ci.yml`. These are the same test, hygiene and scan jobs as on pull requests. `dependency-review` is skipped because the event is not a pull request.
3. **build** (`ubuntu-latest`): setup-go with `stable` and `check-latest: true`.
   - For each of `linux/amd64`, `linux/arm64`, `darwin/amd64` and `darwin/arm64` (four since v0.4.1, feature 014): `CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X github.com/tosin2013/spec-status/internal/cli.Version=<version> -X …cli.Commit=$GITHUB_SHA -X …cli.Kind=release"`.
   - Package per [release-assets.md](release-assets.md), then write the checksum file once.
   - Print `go version` and save it as output `go_version`.
   - Upload all release files as one artifact.
4. **smoke** (matrix `ubuntu-latest`, `ubuntu-24.04-arm`, `macos-latest`, `macos-15-intel`, `fail-fast: false`):
   - Assert the platform from `uname -s -m`: `Linux x86_64`, `Linux aarch64`, `Darwin arm64` or `Darwin x86_64`; anything else fails.
   - Download the artifact and verify the checksums. Extract only this platform's archive and check its file list.
   - Run the checks in [research.md](../research.md) R6. Any mismatch fails the job.
5. **release** (`ubuntu-latest`, `permissions: contents: write`):
   - Download the artifact. Re-check that no release exists.
   - Render `notes.md` from the template in [release-assets.md](release-assets.md), one archive line per smoke job. Fail if any expected smoke job URL or archive is missing.
   - Run `gh release create "$GITHUB_REF_NAME" <4 archives, extension zip, checksums> --draft --verify-tag --title "$GITHUB_REF_NAME" --generate-notes --notes-file notes.md`, adding `--prerelease` when needed.

## Must not

- Publish a release, edit or delete an existing release or asset, push commits or tags, or move a tag.
- Rebuild in `smoke` or `release`. Only the build job's artifact bytes are used.
- Run any job after a failed predecessor. All jobs use `needs`, with no `if: always()`.

## Observable outcomes

| Situation | Result |
|---|---|
| Valid tag on `main`, all checks pass | One draft release with 3 assets; pre-release if `-rc.N` |
| Tag not SemVer, or not on `main`, or release exists | `verify` fails; nothing else runs |
| Any CI, build or smoke failure | No draft is created |
