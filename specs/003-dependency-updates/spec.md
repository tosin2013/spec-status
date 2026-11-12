# Feature Specification: Dependency update automation

**Feature Branch**: `003-dependency-updates`

**Created**: 2026-11-11

**Status**: Draft

**Roadmap item**: 009 (ahead of 004, release engineering)

**Input**: Owner choices of 11 November 2026: automated dependency update proposals (security fixes as soon as an advisory exists, plus one weekly grouped proposal for minor and patch bumps; major bumps separately; workflow components tracked too); an automatic check on every pull request and on `main` (tests on the two chosen v1 platforms, static analysis, and a known-vulnerability scan); a merge policy under the constitution's dependency allowlist; repository vulnerability alerts and automatic security fixes (enabled on 11 November 2026); and a SECURITY.md section on dependency vulnerabilities.

## Context

`spec status` is a read-only CLI with six direct third-party modules (the command framework, its flag library, the text-styling library and two of its helpers, and the operating-system interface package) and about ten indirect ones. The constitution (Art. II) allows only owner-approved modules in each package; `internal/status/imports_test.go` enforces that list. It does not say whether a version change of an approved module needs approval.

Today nothing tells the owner when a dependency has a published vulnerability or falls behind, and nothing runs the tests on a proposed change before merge: the repository has no continuous integration. Release engineering (roadmap item 004) needs the same test check, so this feature builds it first and leaves the release workflow to 004.

This feature is development tooling only. The `spec` binary gains no behavior, network access, dependency, or output change.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Learn about a vulnerable dependency without looking (Priority: P1)

As the owner, when a vulnerability is published for a module the project uses, I am told, and a proposed fix arrives as a pull request whose tests have already run, so I can decide quickly whether to merge.

**Why this priority**: A known vulnerability in a dependency is the most serious dependency risk, and today nothing reports it.

**Independent Test**: With a module pinned to a version that has a published advisory (in a scratch fork or by inspecting an existing alert), confirm an alert appears in the repository's security view and a fix proposal is opened with the automatic check attached.

**Acceptance Scenarios**:

1. **Given** a published advisory for a module version in `go.sum`, **When** the repository's alerting runs, **Then** an alert lists the module, the affected version, and the fixed version.
2. **Given** such an alert with a fixed version available, **When** automatic security fixes run, **Then** a pull request proposes the smallest version change that fixes it, and the automatic check runs on it.
3. **Given** an advisory whose vulnerable code the project actually calls, **When** the automatic check runs on any pull request or on `main`, **Then** the known-vulnerability scan fails and names the advisory.

---

### User Story 2 - Every pull request is tested before merge (Priority: P1)

As the owner, every pull request, including automated dependency proposals, shows whether the full test suite passes on both chosen v1 platforms, so I never merge an untested change.

**Why this priority**: Without it, automated proposals are unsafe to merge, and 004 needs the same check.

**Independent Test**: Open a pull request that breaks one test; the check fails on both platforms. Revert; the check passes.

**Acceptance Scenarios**:

1. **Given** a pull request to `main`, **When** it is opened or updated, **Then** the check runs the full test suite with cgo disabled on linux/amd64 and darwin/arm64, runs static analysis, and runs the known-vulnerability scan, and reports one pass or fail result per job.
2. **Given** a push to `main`, **When** it lands, **Then** the same check runs.
3. **Given** a failing test on only one platform, **When** the check runs, **Then** that platform's job fails and the other still reports its own result.
4. **Given** the check's own workflow components (checkout, language setup), **When** they are referenced, **Then** each is pinned to an exact version so a dependency proposal can update it deliberately.

---

### User Story 3 - Routine updates arrive in a manageable batch (Priority: P2)

As the owner, routine minor and patch updates arrive as at most one grouped proposal per week for modules and one for workflow components, while each major version update arrives on its own, so I can review breaking changes separately.

**Why this priority**: Keeps dependencies current without a stream of single-module proposals; less urgent than security fixes.

**Independent Test**: Inspect the update configuration: weekly schedule, one group for minor and patch module updates, one for workflow components, and major updates excluded from the groups.

**Acceptance Scenarios**:

1. **Given** several modules with new minor or patch versions in one week, **When** the weekly run happens, **Then** one pull request updates all of them together.
2. **Given** a module with a new major version, **When** the weekly run happens, **Then** it arrives as its own pull request, not in the group.
3. **Given** workflow components with new versions, **When** the weekly run happens, **Then** one grouped pull request updates them.

---

### User Story 4 - Merges keep the dependency allowlist honest (Priority: P2)

As the owner, I have a written rule for which dependency proposals I may merge on a passing check alone and which need a recorded decision first, so automated proposals cannot quietly widen the set of third-party code the binary contains.

**Why this priority**: The constitution restricts dependencies; automation must not erode that.

**Independent Test**: Read CONTRIBUTING.md and the decision log; apply the rule to two example proposals, one bumping an approved module with no new indirect modules and one adding a new indirect module, and get "merge after review" and "decision needed" respectively.

**Acceptance Scenarios**:

1. **Given** a proposal that only changes versions of modules already in `go.mod`, **When** the check passes and the owner approves the pull request, **Then** it may merge without a decision entry.
2. **Given** a proposal that adds a module not previously in `go.mod` (direct or indirect), **When** it is reviewed, **Then** it is not merged until an owner decision naming that module is recorded in the decision log.
3. **Given** a proposal of either kind, **When** the existing dependency-allowlist test runs, **Then** it still passes, and a proposal that makes it fail is not merged.

---

### User Story 5 - Reporters and users know how dependency vulnerabilities are handled (Priority: P3)

As a security reporter or user, SECURITY.md tells me how dependency vulnerabilities are found and handled and when a fix reaches me, without overstating it.

**Why this priority**: Required for truthful disclosure; low effort.

**Independent Test**: Read SECURITY.md: it names the detection methods, says a fix is delivered only through a release, says no release exists yet, and claims no audit or certification.

**Acceptance Scenarios**:

1. **Given** SECURITY.md, **When** read, **Then** it states that dependency vulnerabilities are detected by repository alerts and by a known-vulnerability scan in the automatic check, that fixes are proposed as pull requests and reviewed by the maintainer, and that a fix reaches users only in a release.
2. **Given** no release exists, **When** SECURITY.md is read, **Then** it says so and that fixes land on `main` until the first release.
3. **Given** SECURITY.md, **When** read, **Then** it makes no claim of a security audit, supply-chain certification, signed artifacts, or a response time.

### Edge Cases

- **A vulnerability with no fixed version yet**: the alert stays open; no proposal is opened. The owner decides whether a workaround or a README limit is needed. The scan fails only if the vulnerable code is reachable.
- **A vulnerability in the language toolchain or standard library**: the scan reports it, and no module proposal can fix it. If a patch release of the tested line fixes it, the next run picks that patch up automatically. If only a newer line fixes it, or the minimum line has left support, the owner decides on raising the minimum (D-2026-11-11-04). The result is recorded, never suppressed.
- **A proposal that also changes the minimum language version in `go.mod`**: it needs an owner decision, like a new dependency (D-2026-11-11-04).
- **A grouped proposal where one module breaks the tests**: the check fails; the owner can ask for the group to be split or close it. Nothing merges on a failing check.
- **Automated proposals from forks or with workflow edits**: the check runs with read-only repository permissions and no secrets, so a proposal cannot publish or write anything.
- **Advisory database or scan service unreachable**: the scan job fails visibly, it does not pass silently.
- **darwin/arm64 runner unavailable**: the darwin job fails or queues visibly; it is not skipped or marked as passed.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST carry an update configuration that proposes Go module updates on a weekly schedule, grouping all minor and patch updates into one proposal and opening each major update separately.
- **FR-002**: The same configuration MUST track the versions of workflow components used by the repository's automated checks, grouped into one weekly proposal.
- **FR-003**: Repository vulnerability alerts and automatic security-fix proposals MUST be enabled (done 11 November 2026), and security-fix proposals MUST arrive regardless of the weekly schedule.
- **FR-004**: An automatic check MUST run on every pull request to `main` and every push to `main`.
- **FR-005**: The check MUST run the full test suite with cgo disabled on linux/amd64 and on darwin/arm64, as separate jobs with separate results.
- **FR-006**: The check MUST run the language's static analysis and a known-vulnerability scan that reports only advisories affecting code the project can reach, and MUST fail when either finds a problem.
- **FR-007**: The check MUST use read-only repository permissions, MUST use no secrets, and MUST pin every workflow component to an exact version.
- **FR-008**: The check MUST NOT build release archives, create tags, or publish anything; that is roadmap item 004.
- **FR-009**: CONTRIBUTING.md MUST state the merge rule: a proposal that only changes versions of modules already in `go.mod` may merge after the check passes and the owner approves the pull request; a proposal that adds any module not already in `go.mod`, or changes the minimum language version, MUST NOT merge until an owner decision naming the change is recorded in `specs/001-spec-status/decisions.md`.
- **FR-010**: The merge rule MUST be recorded as an owner decision in the decision log.
- **FR-011**: SECURITY.md MUST gain a dependency-vulnerability section meeting User Story 5's acceptance scenarios.
- **FR-012**: The `spec` binary's behavior, output, exit codes, JSON bytes, imports, and dependency set MUST be unchanged by this feature; all existing tests, including the dependency-allowlist test, MUST pass.
- **FR-013**: No document MAY claim a security audit, supply-chain certification, signed or attested artifacts, reproducible builds, or a supported platform as a result of this feature (constitution Art. VII).
- **FR-015**: `go.mod` MUST declare `go 1.26.0`, the oldest supported Go release line (D-2026-11-11-05, added during planning). The check MUST test the newest patch of that line and the current stable release, and MUST fail if `go.mod`'s minimum and the tested minimum line differ.
- **FR-014**: The README and the release document (`specs/001-spec-status/build-release.md`) MUST state accurately that a test check exists on pull requests and `main`, that it is not a release workflow, and which platforms it runs on.

### Key Entities

- **Update configuration**: the repository file that tells the update service which ecosystems to watch, on what schedule, and how to group proposals.
- **Automatic check**: the workflow run on pull requests and `main`; one job per platform for tests, plus static analysis and the vulnerability scan.
- **Dependency proposal**: a pull request opened by the update service, either a security fix or a scheduled update.
- **Merge rule**: the owner decision classifying proposals as "merge after review" or "decision needed".

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of pull requests opened after this feature merges show a pass or fail result for the test suite on both chosen platforms, static analysis, and the vulnerability scan before merge.
- **SC-002**: On the day of merge, the vulnerability scan on `main` either passes or each reported advisory is listed in the pull request with a decision (fix, accept with README limit, or wait for a fixed version).
- **SC-003**: In any week, routine updates produce at most one grouped proposal for modules and one for workflow components, plus one per major update.
- **SC-004**: Applying the merge rule to any dependency proposal yields exactly one outcome, "merge after review" or "decision needed", with no case left undecided.
- **SC-005**: The `spec` binary built from `main` after this feature produces byte-identical output to the build before it on the example and conformance fixtures.
- **SC-006**: A reader of SECURITY.md can answer three questions without asking: how dependency vulnerabilities are found, how they are fixed, and when a fix reaches them.

## Assumptions

- GitHub is the only repository home (D-2026-11-11-03), so the update service and check are GitHub's own (Dependabot and GitHub Actions); the owner chose them on 11 November 2026.
- Hosted runners exist for linux/amd64 and darwin/arm64. Running the tests there is evidence for this check only; it is not the "proven on a real runner" release evidence that DC-07 requires for advertising a platform, which 004 defines.
- During planning, the scan showed that the declared `go 1.22` (no longer supported) fails on reachable standard-library advisories. The x/sys security fix also needs Go ≥1.25. The owner therefore set the minimum to Go 1.26 in this feature (D-2026-11-11-05). Feature 004 still chooses the exact toolchain patch used to build release archives.
- The known-vulnerability scan is Go's official scanner (`govulncheck`), chosen by the owner; it runs only in the check, never in the binary.
- The constitution needs no amendment: Art. II restricts which modules may be imported, and the merge rule keeps new modules behind an owner decision. Version changes of approved modules are not new dependencies.
- Code review by Doomer or other tools is separate from this feature and not required by the merge rule.

## Out of Scope

- Release, tag, archive, checksum, and publication workflows (roadmap item 004).
- The exact toolchain patch used for release builds (004). The minimum Go version is in scope (FR-015).
- Signing, attestations, or provenance (excluded from v1 by D-2026-11-11-02).
- Automatic merging of any proposal.
- Coverage thresholds, linters beyond the language's own static analysis, and fuzzing in the check.
