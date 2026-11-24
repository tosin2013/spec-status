# Feature Specification: Release archives for linux/arm64 and darwin/amd64

**Feature Branch**: `014-more-release-platforms`

**Created**: 2026-11-24

**Status**: Draft

**Input**: Roadmap item 019, issue #82, decision D-2026-11-24-01 (1): "add `linux/arm64` and `darwin/amd64` release archives, covered by the checksums file and each smoke-tested natively (`ubuntu-24.04-arm`, `macos-15-intel`); no CLI, JSON or library change." Pull-request CI runs the full test suite on every supported platform (D-2026-11-24-01 (3)).

## Context

Every spec-status release ships two binary archives, linux/amd64 and darwin/arm64. Each is advertised only because a release smoke test ran on that platform (D-2026-11-11-02). Two groups have no verified binary:

- **Arm64 Linux:** for example, Chalkeus's tools image for Apple Silicon Docker and arm64 clusters (issue #82).
- **Intel Macs.**

The binary is built without C dependencies, so producing more targets is easy. What has to be added is the evidence: each archive smoke-tested on its own kind of machine, and the test suite run there on every change. The release notes, install docs and agent prompt must also stop telling these users their platform is unsupported.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - An arm64 Linux image installs a verified binary (Priority: P1)

A team building a `linux/arm64` container image downloads `spec_<v>_linux_arm64.tar.gz` and the checksums file from the latest release, verifies the checksum, and runs `spec`. The release notes link a smoke test that ran this archive on a native arm64 Linux machine.

**Why this priority**: This is the request in issue #82, and it blocks Chalkeus roadmap item 031.

**Independent Test**: In the release run, the smoke job on a native arm64 Linux runner verifies the archive's checksum and contents, and runs the binary. It reports target `linux/arm64` and label `release build`, and produces text, JSON and fatal exits as on the other platforms.

**Acceptance Scenarios**:

1. **Given** a release, **When** a user lists its assets, **Then** `spec_<v>_linux_arm64.tar.gz` is present and appears in `spec_<v>_checksums.txt`.
2. **Given** that archive on an arm64 Linux machine, **When** the user verifies and runs it, **Then** `spec version` shows the release version, the tagged commit, `target: linux/arm64` and `label: release build`.
3. **Given** the release notes, **When** a user reads the archive list, **Then** the linux/arm64 line links the smoke job that tested it.

---

### User Story 2 - An Intel Mac user installs a verified binary (Priority: P1)

An Intel Mac user follows the README install steps with `darwin_amd64` and gets a working, checksum-verified binary. The notes link its native smoke test.

**Why this priority**: This is the other half of issue #82.

**Independent Test**: The release smoke job on a native Intel macOS runner passes the same checks, with `target: darwin/amd64`.

**Acceptance Scenarios**:

1. **Given** a release, **When** a user lists its assets, **Then** `spec_<v>_darwin_amd64.tar.gz` is present and appears in the checksums file.
2. **Given** that archive on an Intel Mac, **When** the user verifies and runs it, **Then** `spec version` shows `target: darwin/amd64` and `label: release build`.

---

### User Story 3 - Every change is tested where the binaries run (Priority: P2)

A maintainer opening a pull request sees the full test suite run on linux/amd64, linux/arm64, darwin/arm64 and darwin/amd64. A platform-specific regression is therefore caught before a release, not by a smoke test at tag time.

**Why this priority**: Owner decision D-2026-11-24-01 (3). Smoke tests prove the archive runs; the test suite proves its behavior.

**Independent Test**: A pull request's CI shows passing test jobs for the two new runners, each confirming it ran on the expected operating system and architecture.

**Acceptance Scenarios**:

1. **Given** a pull request, **When** CI runs, **Then** test jobs run on the native arm64 Linux and Intel macOS runners, and each job checks that it really ran on that platform.
2. **Given** a runner that resolves to an unexpected platform, **When** the job starts, **Then** it fails instead of passing silently on the wrong platform.

---

### User Story 4 - Docs and the agent prompt know the new platforms (Priority: P2)

A user, or an AI agent following the published prompt, on arm64 Linux or an Intel Mac is told which archive to download, instead of "no release archive: stop".

**Why this priority**: Without this, the new archives are invisible to the install path most users and agents follow.

**Independent Test**: The README install section, getting-started guide, agent prompt and release notes map all four platforms to their archive names. No doc still says only two platforms are supported.

**Acceptance Scenarios**:

1. **Given** the agent prompt, **When** an agent runs on Linux arm64 or macOS Intel, **Then** it maps the machine to `linux_arm64` or `darwin_amd64`.
2. **Given** the README's supported-platform list, **When** a user reads it, **Then** all four platforms are listed, each with the version that first shipped it, and a smoke-test link once released.

### Edge Cases

- **An Intel macOS runner is retired:** the release must not silently skip the platform. If a smoke job cannot run on a native machine, the release workflow fails, and the platform is not advertised. A later decision is needed to drop it or move it to another runner.
- **A runner label resolves to a different architecture:** the smoke and test jobs check the actual operating system and architecture and fail if they differ.
- **Archive layout:** the new archives contain exactly the same files as the existing ones: a top-level directory with `spec`, LICENSE and README.md.
- **The extension archive** is still checked once, not once per platform.
- **Draft release notes** list exactly the archives that were smoke-tested, each with its own smoke-job link.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Each release MUST include `spec_<v>_linux_arm64.tar.gz` and `spec_<v>_darwin_amd64.tar.gz`, built the same way as the existing archives (same build flags and build identity), alongside the existing two.
- **FR-002**: The checksums file MUST list all four archives and the extension archive.
- **FR-003**: Each archive MUST be smoke-tested on a native machine of its own platform with the existing checks:
  - checksum verified;
  - exact archive contents;
  - `spec version` shows the version, commit, `target: <os>/<arch>` and `release build`;
  - help text says the release version;
  - status text and JSON exit codes agree on the example roadmap;
  - a missing root exits 2;
  - the JSON reports schema 1.0.
- **FR-004**: The release MUST NOT be drafted unless every archive's smoke test passed.
- **FR-005**: The draft release notes MUST list each archive as "smoke-tested on <os>/<arch>" with the link to its smoke job. All four archives, the checksums file and the extension archive MUST be attached.
- **FR-006**: Pull-request CI MUST run the full test suite on native linux/arm64 and darwin/amd64 runners, in addition to the existing linux/amd64 and darwin/arm64 runners. Each job MUST fail if the runner's actual platform differs from the expected one.
- **FR-007**: The documentation MUST list and map all four platforms:
  - README install and supported platforms (with first version and smoke-test link);
  - getting started;
  - the agent prompt (which also feeds the release notes);
  - architecture, development and limits;
  - RELEASING.md checklist;
  - the release-engineering contracts.
- **FR-008**: The extension version MUST equal the release version (0.4.1) so the release workflow accepts the tag.
- **FR-009**: The CLI's behavior, JSON, exit codes and the Go library MUST be unchanged.
- **FR-010**: Windows MUST remain out of scope; docs keep saying it is not provided yet (roadmap 020).

### Key Entities

- **Supported platform:** an operating system and architecture pair with a release archive, a native smoke test in each release, and native CI test jobs.
- **Release archive:** `spec_<v>_<os>_<arch>.tar.gz` with a checksum entry and a smoke-job link in the notes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The v0.4.1 release has 4 binary archives, and all 4 are listed in the checksums file.
- **SC-002**: 4 of 4 archives have a passing native smoke test linked from the release notes, each showing the expected `target:`.
- **SC-003**: Every pull request after this change runs the test suite on 4 platforms, and a job running on the wrong platform fails.
- **SC-004**: A user on arm64 Linux or an Intel Mac can install a checksum-verified binary by following the README alone, with no source build.
- **SC-005**: CLI output on the existing fixtures is unchanged; there are no code changes outside CI, release tooling and docs.

## Assumptions

- GitHub-hosted runners available on 2026-11-24 (actions/runner-images README): `ubuntu-24.04-arm` (arm64 Linux, available to public repositories) and `macos-15-intel` (x64 macOS). `macos-26-intel` also exists; `macos-15-intel` is chosen as the established label.
- `uname -m` reports `aarch64` on arm64 Linux and `x86_64` on Intel macOS.
- Branch protection requires the existing four test jobs. Whether to add the new jobs to the required set is the owner's repository setting; this feature does not change it.
- v0.4.1 is a patch release under the pre-1.0 policy: new assets, no compatibility change.
