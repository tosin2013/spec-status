# Feature Specification: Release engineering

**Feature Branch**: `004-release-engineering`

**Created**: 2026-11-11

**Status**: Draft

**Roadmap item**: 004 (depends on 002, 003, 009)

**Input**: Owner choices of 11 November 2026 (D-2026-11-11-06):
- **Release flow:** the owner pushes a SemVer tag. A workflow runs the full CI, builds both archives, smoke-tests each archive on its own platform, and creates a **draft** GitHub Release with checksums. Nothing is public until the owner publishes it. Release-candidate tags (`vX.Y.Z-rc.N`) are allowed.
- **First version:** `v0.1.0` (pre-1.0).
- **Channels:** GitHub Release archives, plus a documented `go install …@vX.Y.Z`. The version output reports the module version for `go install` builds.

Earlier decisions bind this feature:
- D-2026-11-11-02: linux/amd64 and darwin/arm64, each advertised only after a real runner proves it; SHA-256 checksums only, with no signing; only the latest release supported.
- D-2026-11-11-03: GitHub is the only home.
- D-2026-11-11-05: Go 1.26 minimum.
- [build-release.md](../001-spec-status/build-release.md): the retained pipeline contract.

## Context

`spec status` has no release. Users can only build from source, and every build calls itself a "development build, not a release": in `spec version`, the status panel, and `--help`. Feature 003 (roadmap 009) adds the test-and-scan CI that this feature reuses. This feature is the step from "code on `main`" to "a tagged, tested, downloadable version that a user can install without Go".

## User Scenarios & Testing *(mandatory)*

### User Story 1 - The owner cuts a release safely (Priority: P1)

As the owner, I push a version tag and get a complete draft release to review: archives for both platforms, a checksum file, and notes, all built from exactly the tagged commit, tested, and smoke-tested on each platform. I publish it only when satisfied. A failed step leaves nothing public.

**Why this priority**: Nothing else in this feature exists without it.

**Independent Test**: Push `v0.1.0-rc.1` on a commit that passes CI. A draft pre-release appears with two archives and a checksum file, and the run shows passing tests, builds and smoke tests. Then push a tag on a commit with a failing test: no release or draft is created.

**Acceptance Scenarios**:

1. **Given** the owner pushes a tag `vX.Y.Z` or `vX.Y.Z-rc.N` on a commit of `main`, **When** the release workflow runs, **Then** it checks out exactly that commit, runs the same tests and scans as CI, and stops before building if any of them fail.
2. **Given** the checks pass, **When** the workflow builds, **Then** it produces one archive per supported target and a checksum file listing the SHA-256 of each archive.
3. **Given** the archives, **When** the smoke test runs on a linux/amd64 runner and on a darwin/arm64 runner, **Then** each extracts its own archive, verifies the checksum, and runs `spec version`, `spec --help`, `spec status` (human and JSON) on the example roadmap, and an exit-code check. Any failure stops publication.
4. **Given** all steps pass, **When** the workflow finishes, **Then** a **draft** GitHub Release exists for the tag, with the archives and checksum file attached and generated notes. A `-rc.N` tag's draft is marked pre-release.
5. **Given** a draft release, **When** the owner reviews and publishes it, **Then** it becomes public. The workflow itself never publishes.
6. **Given** a tag that is not a valid SemVer `v` tag, or that points at a commit not on `main`, **When** the workflow runs, **Then** it fails before building, with a message.
7. **Given** a published release, **When** anyone tries to change its assets, **Then** the process forbids it: corrections ship as a new version (build-release.md item 8).

---

### User Story 2 - A user installs without Go (Priority: P1)

As a user on linux/amd64 or darwin/arm64, I can download the right archive, verify its checksum, put `spec` on my PATH without root, run it, and later upgrade or roll back by the same steps.

**Why this priority**: The point of a release.

**Independent Test**: In a clean environment without Go and without root, follow the README install section for the published `v0.1.0`. `spec version` shows `0.1.0` and identifies a release build. Repeat the upgrade steps with a later release candidate.

**Acceptance Scenarios**:

1. **Given** the README install section, **When** a user follows it, **Then** it names the exact archive for each supported platform, shows how to verify it against the checksum file, extract it, place `spec` in a user-owned PATH directory, and run it. It needs no root and no Go.
2. **Given** a macOS user who downloaded the archive in a browser, **When** they first run the unsigned binary, **Then** the README explains the Gatekeeper prompt and the documented way to allow it. It does not claim the binary is signed or notarized.
3. **Given** an installed version, **When** the user follows the upgrade or rollback steps, **Then** they replace the binary with another release's binary by the same verify-extract-place steps. Earlier releases stay downloadable.
4. **Given** an unsupported platform (for example Windows or linux/arm64), **When** a user reads the README, **Then** it says plainly that no release archive exists for it, and that building from source or `go install` may work but is untested.

---

### User Story 3 - A Go user installs with `go install` and sees the real version (Priority: P2)

As a Go developer, I can run `go install github.com/tosin2013/spec-status/cmd/spec@v0.1.0`, and `spec version` reports `0.1.0`. It says the binary was built from the module version by my own toolchain, not by the release workflow.

**Why this priority**: This channel exists the moment a tag is pushed. Without the change it would report `0.0.0-dev`.

**Independent Test**: After the tag exists, run `go install …@v0.1.0` in a clean module cache. `spec version` shows version `0.1.0` and a label that distinguishes it from a release archive.

**Acceptance Scenarios**:

1. **Given** a binary built by `go install` at a tagged module version, **When** the user runs `spec version`, **Then** it shows that module version, the Go version and the target, and a label stating it was built from the module version, not a release archive.
2. **Given** a build from a source checkout with no release stamping, **When** the user runs `spec version`, **Then** it still says development build.
3. **Given** a binary from a release archive, **When** the user runs `spec version`, **Then** it shows the release version, the tagged commit, the Go version, the target, and a label identifying a release build.

---

### User Story 4 - Release builds stop calling themselves development builds (Priority: P2)

As a user of a release, the status panel and help text no longer say "development build". Development builds still do.

**Why this priority**: Truthful labeling (constitution Art. I and VII). It is visible on every run.

**Independent Test**: Compare `spec status --color=never` and `spec --help` from a release-stamped build and a plain build. Only the build-identity label differs, and JSON bytes are identical.

**Acceptance Scenarios**:

1. **Given** a release-stamped binary, **When** it prints the human status panel, **Then** the panel heading shows the release version instead of "development build".
2. **Given** a release-stamped binary, **When** it prints `--help`, **Then** the help text drops "This is a development build" and keeps the read-only statement.
3. **Given** any build, **When** it prints `--json`, **Then** the bytes are identical to a development build's for the same input. Build identity is not part of the JSON 1.0 schema.

---

### User Story 5 - Each release is traceable and documented (Priority: P3)

As a user or reviewer, each release states its version, tagged commit, the Go version used to build it, the targets it supports and the evidence for them, the checksum file, the CLI/JSON compatibility policy for pre-1.0 versions, and its known limits.

**Why this priority**: Required by build-release.md "Identity and compatibility". Mostly text.

**Independent Test**: Read the release notes of `v0.1.0-rc.1`. Each item above is present, and nothing claims signing, provenance, reproducible builds, or an audit.

**Acceptance Scenarios**:

1. **Given** a draft release, **When** the owner opens it, **Then** the notes include the version, commit, Go version, targets with links to the smoke-test runs, the checksum file name, the pre-1.0 compatibility statement, and a link to README Limits.
2. **Given** the notes and README, **When** read, **Then** they make no claim of signed binaries, provenance, reproducible builds, a security audit, or support for any platform other than the two chosen.

### Edge Cases

- **Tag pushed by someone other than the owner**: only users with write access can push tags. The workflow does not publish, so a stray tag yields at most a draft that the owner deletes. The tag protection setting is an owner choice (see Assumptions).
- **Tag on a commit not reachable from `main`**: fail before building (Acceptance 1.6).
- **Re-running a release for the same tag**: if a release for the tag already exists, the workflow fails rather than replacing assets.
- **A smoke test passes on one platform and fails on the other**: no draft is created. Releasing with one platform needs a decision and is not an automatic fallback.
- **darwin/arm64 runner unavailable**: the job queues or fails visibly. It is never skipped.
- **A reachable vulnerability appears between CI and the tag**: the release workflow's own scan fails, and no draft is created.
- **Version mismatch**: if the binary's reported version differs from the tag, the smoke test fails.
- **`go install` from a pseudo-version or a fork**: `spec version` reports what Go's build information gives (a pseudo-version or `(devel)`). Development labeling applies when no module version exists.
- **macOS quarantine**: unsigned binaries trigger Gatekeeper. Document it; never instruct users to disable Gatekeeper system-wide.

## Requirements *(mandatory)*

### Functional Requirements

**Release workflow**
- **FR-001**: A release workflow MUST run only when a tag matching `v<MAJOR>.<MINOR>.<PATCH>` or `v<MAJOR>.<MINOR>.<PATCH>-rc.<N>` is pushed. It MUST fail before building if the tagged commit is not reachable from `main`.
- **FR-002**: It MUST check out exactly the tagged commit and run the same tests, `go vet`, vulnerability scan and hygiene checks as feature 003's CI, on both supported platforms, before building.
- **FR-003**: It MUST build with cgo disabled, for linux/amd64 and darwin/arm64 only, with an exact Go version that is recorded. It MUST stamp the version (the tag without the `v`), the full commit, and a release marker into the binary.
- **FR-004**: It MUST produce `spec_<version>_linux_amd64.tar.gz` and `spec_<version>_darwin_arm64.tar.gz`. Each contains the `spec` binary, `LICENSE` and `README.md`. It MUST also produce `spec_<version>_checksums.txt` in `sha256sum` format.
- **FR-005**: It MUST smoke-test each archive on a runner of its own platform: verify the checksum, extract it, run `spec version` and check that the version and commit match, run `spec --help`, run `spec status` in human and JSON form against the example roadmap with the expected exit code, and check that the JSON is valid. Any failure MUST stop the run before a release is created.
- **FR-006**: It MUST create a **draft** GitHub Release for the tag, with the archives and checksum file attached and notes meeting FR-015. An `-rc.N` tag MUST produce a draft marked pre-release. It MUST NOT publish.
- **FR-007**: Only the job that creates the draft release MAY have `contents: write`. Every other job MUST be read-only. The workflow MUST use no secrets beyond the automatic token, MUST pin every action by commit SHA, and MUST fail rather than replace assets if a release for the tag already exists.
- **FR-008**: Build outputs MUST pass from the build job to the smoke-test and release jobs unchanged, so the bytes tested are the bytes attached. The checksum file is computed once, at build.

**Build identity (product change)**
- **FR-009**: `spec version` MUST distinguish three kinds of build:
  - **Release build** (stamped by the release workflow): version, commit, Go version, target, and the label "release build".
  - **Module build** (no stamp, but Go build information carries a module version, as from `go install …@vX.Y.Z`): that version, the VCS revision if Go recorded one, Go version, target, and a label saying it was built from the module version, not a release archive.
  - **Development build** (neither): the current output, unchanged.
- **FR-010**: The human status panel heading and `--help` MUST drop "development build" wording for release builds and keep it for development builds. Module builds MUST say "module build", not "release".
- **FR-011**: JSON output MUST be byte-identical across build kinds for the same input. Exit codes and diagnostics MUST be unchanged.
- **FR-012**: The build-identity change MUST use only the Go standard library and MUST have tests for all three kinds.

**Documentation and policy**
- **FR-013**: README MUST gain an install section meeting User Story 2:
  - per-platform archive names, checksum verification, extraction, a user-owned PATH directory, no root, no Go;
  - upgrade and rollback steps;
  - the macOS Gatekeeper note;
  - unsupported platforms named as unsupported;
  - the `go install` line with its labeling.
- **FR-014**: README and release notes MUST state the pre-1.0 policy: before `v1.0.0`, a minor version may change CLI flags, exit codes or the JSON schema, and each such change is listed in the notes. Patch versions do not break compatibility.
- **FR-015**: Release notes MUST include the version, commit, Go version, targets with links to their smoke-test runs, the checksum file name, the compatibility statement and a link to README Limits. They MUST NOT claim signing, provenance, reproducible builds, an audit, or other platforms.
- **FR-016**: Once both platforms have passed the smoke test of a published release, README MAY list linux/amd64 and darwin/arm64 as supported for that version, citing the run (DC-07). Until then they stay "chosen, not yet proven".
- **FR-017**: A short owner release checklist MUST exist in the repository: tag, review the draft, check the assets and notes, publish, and post-publish checks. Every step is performed by the owner.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A release candidate (`v0.1.0-rc.1`) produces a complete draft (2 archives, 1 checksum file, notes) with no manual build steps, in one workflow run.
- **SC-002**: 100% of release runs with any failing test, scan, build or smoke step create no draft release.
- **SC-003**: A user without Go and without root installs `v0.1.0` from the README on each supported platform in under 5 minutes, and the checksum check succeeds.
- **SC-004**: `spec version` correctly identifies release, module and development builds in 3 of 3 cases.
- **SC-005**: JSON output for the example and conformance fixtures is byte-identical between a release build and a development build of the same commit.
- **SC-006**: The bytes of each downloaded archive match the checksum file and the archive that was smoke-tested (one checksum computed at build, checked again at smoke test and by the user).
- **SC-007**: The first published release, `v0.1.0`, follows an owner-approved release candidate, as build-release.md requires.

## Assumptions

- The release workflow reuses feature 003's CI checks. 003 must be merged first.
- Releases are built with the current stable Go patch at release time, and that exact version is recorded in notes and `spec version`. The Go 1.26 minimum (D-2026-11-11-05) is for source and `go install` builds.
- Asset names, archive contents and the checksum format follow the defaults in FR-004. They can be changed in planning without new owner input.
- Release notes start from GitHub's generated notes, completed in the draft by the workflow's template and edited by the owner before publishing.
- Hosted runners for linux/amd64 and darwin/arm64 count as "real runners" for DC-07, because the binary actually runs there. Cross-compiling alone never does.
- Whether to turn on GitHub's tag protection rules and immutable-releases setting for this repository is an owner settings choice, offered during implementation and not applied without approval.
- Roadmap 004's wording "first stable tag" becomes "first tagged release (v0.1.0)" (D-2026-11-11-06).
- The constitution needs no amendment. The build-identity change is a labeling change within Art. I and VII, uses the standard library (Art. II), and leaves the JSON schema unchanged (Art. V).

## Out of Scope

- Signing, notarization, attestations, provenance, SBOMs (excluded from v1, D-2026-11-11-02).
- Homebrew, Scoop, WinGet, Linux packages, and container images.
- Windows or linux/arm64 archives.
- Automatic publishing, automatic version bumps, and changelog tooling.
- Performance benchmarking as a release gate.
