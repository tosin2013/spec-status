# Build, release and distribution requirements

**Release requirements remain a draft.** A development CLI and implementer-run tests now exist in the repository. No release workflow, tag, release, published binary, artifact manifest, or supported-platform claim exists. A test-and-scan CI workflow exists (feature 003, `.github/workflows/ci.yml`): tests, `go vet` and govulncheck on linux/amd64 and darwin/arm64 with Go 1.26.x and stable, plus hygiene checks that include compile-only builds for windows/amd64 and linux/arm64, which are not support. The release workflow may reuse it. CI passes are not release evidence. On 11 November 2026 the owner chose the v1 targets, checksum-only integrity, the Apache-2.0 license, the security reporting route and the support window (DC-07, DC-08, D-2026-11-11-02). GitHub is the only home and release channel (D-2026-11-11-03). The minimum Go version is 1.26 (D-2026-11-11-05). The exact Go patch used for release builds, release-runner evidence and the release workflow remain open (feature 004).

## Identity and compatibility

Retained requirements: CLI SemVer tags `vX.Y.Z`; supported roadmap input v0.1.1 is independent, not lockstep; JSON schema is versioned independently (1.0, D-2026-11-11-02). Version output identifies CLI version, tagged commit, exact Go version and target platform; untagged/development builds plainly dev. Describe CLI/JSON breaking-change policy, roadmap dialect/extension/limitations and platform coverage for each release.

Candidate metadata policy should record tag, commit, Go/GOOS/GOARCH and build settings; controlled trimpath/VCS stamping and reproducibility limits are proposed refinements to approve. An installed version label alone is not proof of origin. Signed provenance/attestations require separate implementation and verification, not aspirational README language.

## Future tag-triggered pipeline contract

1. Owner-approved candidate commit and explicit owner-created SemVer tag are inputs. No PRD/spec/task author creates/pushes a tag by implication.
2. Check out the exact tagged commit, not moving default branch. Pin supported patched Go and reviewed action commit SHAs. Independently review current advisories and platform-safe filesystem behavior; historic patch numbers are not a substitute.
3. Validate tag/version/source consistency; run required unit/fixture/integration/security-contract tests before publication. Build with `CGO_ENABLED=0`. The analysis core stays in the standard library. The owner-approved Cobra and Lip Gloss dependencies and their pinned transitive modules are part of that build; they are not a license to add further runtime dependencies. Workflow tooling is distinct from runtime requirements.
4. Build selected target binaries, create version/OS/architecture archives and SHA-256 manifest. Linux/macOS archives `.tar.gz`; Windows `.zip` only if that target becomes supported. No source-toolchain install required for end users.
5. Execute version/help/startup/sample status/JSON/archive/checksum checks on suitable real target runners. Check tag/commit correspondence. Failed required checks block publication. Cross-compiling is not executing; do not label untested targets supported.
6. Only a narrowly scoped publish job receives `contents:write`; other jobs read-only. Review action pins, untrusted pull-request contexts and artifact promotion before later workflow implementation. No secrets are a runtime dependency.
7. Approved publication attaches archives and checksum manifest as downloadable **GitHub Release assets**, not solely temporary Actions artifacts. Exact asset promotion/build-once policy needs the selected workflow design and verification; no claim of reproducibility merely from pinned inputs.
8. Treat published versions/assets as immutable. Corrections get a new version; do not silently replace published bytes or move source identity. Validate platform enforcement where chosen rather than confusing policy with a technical guarantee.

This is a desired workflow contract, not executable CI configuration or authorization to run publication. Installation and status never invoke GitHub Actions, Git, the extension or an agent.

## Platform targets — chosen 11 November 2026 (D-2026-11-11-02)

| Option | Status / decision tradeoff |
|---|---|
| Linux/macOS/Windows × amd64/arm64 | Original six-target aspiration; not chosen for v1 |
| linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64 | Earlier five-target suggestion; not chosen for v1 (linux/arm64 and Windows rejected as v1 targets, D-2026-11-11-02) |
| linux/amd64 + darwin/arm64 | Chosen v1 targets (D-2026-11-11-02). Each is advertised only after a real runner proves it; neither is proven yet. Windows not advertised supported |

Support window: only the latest release is supported, best effort, with no response-time promise (DC-08, D-2026-11-11-02). The owner still chooses supported OS minimums, test coverage, static-linking exceptions and macOS signing/quarantine caveats. A compatibility note must clearly say untested/unavailable, not imply “works on Windows.” Requirements for all advertised targets remain the same regardless of chosen count. No OS/architecture is claimed supported in the current README.

## Integrity versus identity

- SHA-256 verifies asset consistency relative to the manifest received. If manifest and binary are maliciously replaced together, this alone does not authenticate the publisher.
- Tag/commit/version checks improve traceability but do not automatically establish independently verified publisher identity or code correctness.
- The owner chose truthful unsigned checksum-only integrity for v1, with limitations prominently stated; signing and attestations are not part of v1 (DC-07, D-2026-11-11-02). Adding them later is a separate owner decision.
- Do not claim signed binaries, trusted provenance, bit-reproducible builds, supply-chain certification or a security audit unless the exact mechanism is implemented and validated. Repeated-build byte comparison and external verification, if adopted, need their own evidence.

## End-user distribution requirements (future documentation, not current install commands)

Initial channel: direct GitHub Releases binary download. Eventual guide must identify exact released version, supported OS/architecture selection, correct archive, checksum manifest, extraction, no-Go runtime prerequisites, user-owned `PATH` directory and rootless install/upgrade. No invented URL, curl pipeline, package-manager formula or fabricated release command appears in this draft.

Exercise guide in clean environments without Go/source checkout and without root. State OS minimums, unsupported combinations and static-link exceptions. Candidate additional manual rollback guidance must preserve immutable previous version identity and not execute untrusted downloaded scripts. Package managers (Homebrew/Scoop/WinGet/Linux packages/others) remain deferred pending demand and maintenance ownership.

First stable tag follows explicitly owner-approved tested RC. Archive contents, asset names, metadata, checksum manifest and supported runner evidence must be reviewable before that sign-off. No test failure is waived merely to meet a budget aspiration.

## Public availability gates

Public documentation targets [tosin2013/spec-status](https://github.com/tosin2013/spec-status); repository establishment occurred separately from this package. The license is Apache-2.0 (`LICENSE`); security reports go through GitHub private vulnerability reporting (`SECURITY.md`); only the latest release is supported, best effort (DC-08, D-2026-11-11-02). Release readiness still requires confirmed release configuration, fixture privacy review and owner publication approval. GitHub is the only repository home and release channel; there are no mirrors, Bitbucket copies or extra release channels (D-2026-11-11-03).

Performance (<100 ms cold median on disclosed fixture) and build effort are not measured today. Neither a one-day aspiration nor an earlier unvalidated planning range is a release schedule. See acceptance AC-20–24 and future TASK-17–21; all are unstarted.
