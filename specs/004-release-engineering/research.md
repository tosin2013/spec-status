# Research: Release engineering

Feature 004 · roadmap item 004 · 2026-11-12. Facts were checked with `gh api`, `gh release create --help` (gh 2.103.3), the Go module proxy, and the source; the commands are at the end.

## R1 — Reusing the CI checks (FR-002)

- **Decision**: Add `workflow_call:` to the triggers of `.github/workflows/ci.yml`. `release.yml` calls it with `uses: ./.github/workflows/ci.yml` after verifying the tag.
- **Rationale**: One definition of the tests, vet, govulncheck and hygiene checks. A called workflow checks out the caller's ref, which here is the tag. Inside a called workflow `github.event_name` is the caller's event (`push`), so the `dependency-review` job (`if: github.event_name == 'pull_request'`) is skipped. That's correct, because a tag has no base to compare against.
- **Permissions**: A called workflow cannot widen the caller's permissions. The release workflow's top-level `contents: read` therefore applies to all CI jobs.
- **Alternative considered**: Copying the steps into `release.yml`, which would let the two drift apart.

## R2 — Tag validation and "on main" (FR-001, Acceptance 1.6)

- **Decision**: `on: push: tags: ['v*']`. The first job, `verify`:
  1. Checks the tag against `^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-rc\.(0|[1-9][0-9]*))?$` and fails otherwise.
  2. Checks out with `fetch-depth: 0`, then runs `git merge-base --is-ancestor "$GITHUB_SHA" origin/main`. Fails otherwise.
  3. Runs `gh release view "$TAG"`, and fails if a release, draft or published, already exists (FR-007).
- **Rationale**: GitHub's `tags:` filter is a glob and cannot express SemVer, so the exact pattern is checked in a step. The `v*` filter only limits which pushes start the workflow.

## R3 — Build-identity mechanism (FR-009 to FR-012)

- **Current code**: `internal/cli/cli.go`:
  - `Version = "0.0.0-dev"` and `Commit = "unknown"` are set by `-ldflags -X`.
  - `versionText()` always prints `label: development build, not a release`.
  - `rootLong` says "This is a development build."
  - `internal/status/style.go:295` hard-codes `"  development build"` in the panel heading.
- **Decision**:
  - Add a third stamped variable `cli.Kind`. The release workflow sets it to `release`; its default is empty.
  - A new standard-library helper, `internal/cli/build.go`, resolves the identity, in order:
    - (a) `Kind == "release"` → **release build**: `Version`, `Commit`, label `release build`.
    - (b) else, `debug.ReadBuildInfo()` with `Main.Sum` not empty and `Main.Version` not empty or `(devel)` → **module build**: version = `Main.Version` without the leading `v`; commit = `unknown`, because module-mode builds carry no VCS stamp; label `module build (built from the module version, not a release archive)`.
    - (c) else → **development build**: the current output, byte for byte.
  - `status.View` gains a `Build string` field: the panel heading text after "Spec Status". Empty means `development build`, so existing callers and tests keep their output.
  - The CLI fills it with `development build`, `release <version>`, or `module build <version>`.
  - `--help` uses the label: "This is a development build." for development builds, "This is release <v>." for release builds, and "This is a module build of <v>, not a release archive." for module builds.
- **Rationale**:
  - `runtime/debug` is standard library (constitution Art. II).
  - The JSON path never reads `View`, so JSON stays byte-identical (FR-011).
  - Keeping the empty-means-development default leaves the README sample and the existing tests unchanged.
- **Why `Main.Sum`, not `Main.Version`** (checked 2026-11-12 with `go version -m`):
  - Since Go 1.24, `go build` inside a git checkout stamps `Main.Version` from git. This checkout recorded `v0.0.0-20261112040751-b228f6020c7b+dirty` plus `vcs.revision`, `vcs.time` and `vcs.modified`, with an empty `Main.Sum`. At a clean tagged commit it would record the tag itself.
  - A real `go install pkg@version` from the module proxy records `Main.Version` together with a module checksum (`h1:…`), and no `vcs.*` settings. The comparison case was `golang.org/x/vuln@v1.9.0`.
  - So a non-empty `Main.Sum` is the signal for "built from a downloaded, checksum-verified module version". Testing only `Main.Version` would call every source build, even a clean checkout of the `v0.1.0` tag, a module build.
  - With `-buildvcs=false`, `Main.Version` is `(devel)`.
  - Development-build output stays exactly as it is today. Showing the VCS revision for checkout builds would be a separate change.

## R4 — Release toolchain and build flags (FR-003)

- **Decision**:
  - Setup: `actions/setup-go` with `go-version: stable` and `check-latest: true`.
  - Build: `CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X …cli.Version=$V -X …cli.Commit=$GITHUB_SHA -X …cli.Kind=release" -o spec ./cmd/spec`, per target.
  - Record `go version` in the build log and in the notes.
- **Rationale**:
  - stable is the most-patched release line, and the binaries are static, so users need no Go.
  - `-trimpath` keeps runner paths out of the binary. `-s -w` shrinks it. Neither changes behavior.
  - No reproducibility claim is made (build-release.md), only that the flags are recorded.
- **Darwin**: Go's linker ad-hoc signs darwin/arm64 binaries, which Apple silicon requires to execute. That is not Developer ID signing or notarization, and no such claim is made.

## R5 — Archives and checksums (FR-004, FR-008)

- **Decision**:
  - Archive: `spec_<v>_<os>_<arch>.tar.gz`, containing one top-level directory `spec_<v>_<os>_<arch>/` with `spec`, `LICENSE`, `README.md`. The directory avoids scattering files into the user's current directory.
  - Checksums: `sha256sum spec_<v>_*.tar.gz > spec_<v>_checksums.txt`, run once in the build job.
  - Hand-off: `actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9 # v7.0.2` and `actions/download-artifact@9000827ccba6bdab643e8b6fd33ac0654aef8333 # v8.0.2`, latest on 2026-11-12.
- **Rationale**: The smoke and release jobs download the exact build bytes. The smoke test re-verifies them against the checksum file, so the tested bytes are the attached bytes (SC-006).

## R6 — Smoke test (FR-005)

- **Decision**: A matrix job on `ubuntu-latest` (asserted linux/amd64) and `macos-latest` (asserted darwin/arm64). Each job:
  1. Downloads the artifact.
  2. Verifies the checksum: `sha256sum -c --ignore-missing` on Linux, `shasum -a 256 -c --ignore-missing` on macOS.
  3. Extracts its own archive and checks the directory contains exactly `spec`, `LICENSE`, `README.md`.
  4. Runs `spec version`: requires `spec <v>`, `commit: $GITHUB_SHA` and `label: release build`.
  5. Runs `spec --help`: exit 0, contains `release <v>`, does not contain "development build".
  6. Runs `spec status --color=never --root testdata/example-roadmap` with the expected exit code.
  7. Runs `spec status --json` with the same root: valid JSON with `schema_version` 1.0.
  8. Runs `spec status --root /nonexistent`: exit 2.
- **Rationale**: This covers build-release.md item 5 (version, help, startup, sample status, JSON, archive, checksum checks) on real runners. It runs the downloaded archive, never a fresh build.

## R7 — Creating the draft release (FR-006, FR-007, FR-015)

- **Decision**: A `release` job, the only one with `permissions: contents: write`. It runs `gh release create "$TAG" spec_<v>_* --draft --verify-tag --title "$TAG" --generate-notes --notes-file notes.md`, plus `--prerelease` for `-rc.N` tags, with `GH_TOKEN: ${{ github.token }}`.
  - `notes.md` is rendered from a template in the workflow. It holds the version, commit, Go version, the targets with the run URL of each smoke job, the checksum file name, the pre-1.0 compatibility statement, and the README Limits link.
  - gh prepends `--notes`/`--notes-file` text to the generated notes.
  - The job re-checks `gh release view` right before creating, to close the race window.
- **Rationale**: gh is preinstalled on hosted runners, so no third-party release action is needed. `--draft` keeps it private. `--verify-tag` aborts if the tag disappeared.

## R8 — Repository settings to offer (not applied without approval)

- **Immutable releases**: `GET repos/tosin2013/spec-status/immutable-releases` returned `{"enabled":false}`. When enabled, a published release's assets and tag cannot change, which enforces build-release.md item 8 instead of relying on policy. Drafts stay editable until publish, so the draft flow still works.
- **Tag ruleset**: no rulesets exist. A tag ruleset on `v*` can restrict tag creation and deletion to the owner. Without one, anyone with write access can push a tag. That yields at most a draft, which only the owner can publish.

## R9 — macOS Gatekeeper (FR-013, US2)

- **Finding**: Files downloaded by a browser get the `com.apple.quarantine` attribute. An unsigned, un-notarized binary is then blocked on first run. Files fetched with `curl` don't get the attribute.
- **Decision**: The README shows the `curl` path first. For browser downloads it documents `xattr -d com.apple.quarantine ./spec` for that one file, and never suggests disabling Gatekeeper.

## Commands used

```sh
gh api repos/actions/upload-artifact/releases/latest     # v7.0.2 → cf430e03…
gh api repos/actions/download-artifact/releases/latest   # v8.0.2 → 9000827c…
gh release create --help                                 # --draft, --prerelease, --verify-tag, --generate-notes, -F; notes are prepended to generated notes
gh api repos/tosin2013/spec-status/immutable-releases    # {"enabled":false}
gh api repos/tosin2013/spec-status/rulesets              # []
grep -n 'development build' internal/cli/cli.go internal/status/style.go
```
