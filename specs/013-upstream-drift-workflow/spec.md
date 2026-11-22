# Feature Specification: Upstream drift workflow

**Feature Branch**: `013-upstream-drift-workflow`

**Created**: 2026-11-21

**Status**: Draft

**Input**: Roadmap item 010, "Upstream drift automation", chosen for v0.4.0 by the owner in decision D-2026-11-21-01 (3): "a scheduled, read-only workflow; scope set in its feature spec." Roadmap scope: schedule, which releases to fetch, how a drift is reported. Out of scope: moving pins automatically, and runtime network access. Notes: "reporting by issue would need write permission."

## Context

spec-status supports exactly what pinned upstream templates produce: github/spec-kit v1.2.0 (tasks template) and srobroek/speckit-roadmap v0.1.0 (roadmap template). This is the acceptance bar D-2026-11-07-09. Feature 002 built a drift check that compares a downloaded candidate template with the pin and lists added or removed shapes. Today someone has to remember to download a new release and run it by hand.

Upstream moves on regardless. On 2026-11-21, spec-kit's latest release was v1.4.3 and speckit-roadmap's was v0.1.1, and nobody had checked either. The outcome promised by roadmap item 010 is that "an upstream template change is noticed within a week, so the pin moves deliberately, not by surprise."

The check is development tooling. The `spec` binary never fetches anything (Art. II). Only this repository's CI reaches the network, to download public template files.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - The owner learns about upstream template changes within a week (Priority: P1)

Once a week, without anyone starting it, the repository checks the latest release of each upstream project against its pin. If a supported template's shapes changed in a way that has not been resolved, the run fails and its summary says which project, which release, and which shapes were added or removed. The owner gets GitHub's normal notice of a failed run.

**Why this priority**: This is the whole outcome of roadmap item 010.

**Independent Test**: Run the workflow against a candidate that differs from the pin (a seeded test release). The run fails and its summary names the changed shape keys. Run it against the current upstream releases: it passes or fails according to the real drift, and the summary lists both projects with their tags.

**Acceptance Scenarios**:

1. **Given** an upstream release whose template has an unresolved shape change, **When** the scheduled run happens, **Then** the run fails. Its summary lists the project, the candidate tag and commit, the pinned tag, and each added or removed shape key with an example line.
2. **Given** upstream releases whose templates have no unresolved change (identical, wording-only, or every change listed in `resolutions.md`), **When** the run happens, **Then** it succeeds. The summary lists each project with "no unresolved change" and the tags compared.
3. **Given** the latest upstream release equals the pin, **When** the run happens, **Then** it still compares and reports "same tag as the pin", and succeeds unless the files differ.

---

### User Story 2 - A maintainer checks a specific release on demand (Priority: P2)

Before moving a pin, a maintainer starts the workflow by hand and names a tag for one project, or leaves it empty for the latest release. They get the same report without downloading anything locally.

**Why this priority**: Moving a pin starts with exactly this check (drift-report contract, "Moving a pin", step 2).

**Independent Test**: Start the workflow by hand for spec-kit with tag `v1.2.0`. It reports "same tag as the pin" and passes. With a non-existent tag, it fails with "release not found", not with a drift.

**Acceptance Scenarios**:

1. **Given** a manual run naming project `spec-kit` and tag `v1.4.3`, **When** it runs, **Then** only that project is checked, at that tag.
2. **Given** a tag that does not exist, **When** it runs, **Then** it fails with a fetch error that names the tag. This is reported differently from a drift.

---

### User Story 3 - Failures say what went wrong (Priority: P2)

A run can fail for three different reasons: a drift, a template that is missing at the expected path in the new release (renamed or removed upstream), or a network or API problem. The summary must make clear which one happened, so the owner does not chase a phantom drift.

**Why this priority**: A check that cries wolf gets ignored, and then the outcome is lost.

**Independent Test**: Point one project's template path at a file that does not exist. The summary says "template not found at path …", and no shape comparison is reported for it.

**Acceptance Scenarios**:

1. **Given** the template file is missing at the new tag, **When** the run happens, **Then** it fails with "template not found" for that project and still reports the other project.
2. **Given** the upstream API cannot be reached, **When** the run happens, **Then** it fails with a fetch error, and no drift is claimed.

### Edge Cases

- **Prereleases and drafts** are ignored for the scheduled run. "Latest" means the latest published, non-prerelease release.
- **A scheduled run stops after 60 days of repository inactivity** (GitHub behavior). The docs say so, and say how to re-enable it.
- **Candidate files are untrusted.** They are downloaded into a fresh directory and only read as text by the existing drift check. They are never executed, sourced or committed.
- **Pins never move automatically.** The workflow never writes to the repository, opens issues or pull requests, or changes `resolutions.md`.
- **Two projects in one run:** a failure in one project does not hide the other project's result.
- **Large or odd files:** the candidate is read only by the drift test, which already treats unknown shapes as reportable changes.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST have a workflow that runs on a weekly schedule and on manual dispatch.
- **FR-002**: For each pinned project (spec-kit, speckit-roadmap), a scheduled run MUST resolve the latest published, non-prerelease release. It MUST download each pinned template file at that tag from the project's repository, and run the existing drift check against the pin.
- **FR-003**: A manual run MUST accept an optional project and an optional tag, and check only what was named. Empty means all projects at their latest release.
- **FR-004**: A run MUST fail when any project has an unresolved shape change, a missing template, or a fetch error. It MUST succeed otherwise.
- **FR-005**: The run summary MUST list, for each project:
  - the pinned tag, the candidate tag and the candidate commit;
  - one result: no unresolved change, same tag as the pin, N unresolved changes (each added or removed key with an example line), template not found, or fetch error.
- **FR-006**: The workflow MUST have read-only repository permissions and use no secrets beyond the default token. It MUST pin third-party actions by commit, as the existing workflows do.
- **FR-007**: The workflow MUST NOT write to the repository, create issues or pull requests, or move pins.
- **FR-008**: Candidate files MUST be treated as untrusted text: downloaded into a fresh directory, only read by the drift check, and never executed or committed.
- **FR-009**: The pinned projects, their pinned tags and their template paths MUST come from one place that the drift test also uses, so the workflow cannot check a different file than the test pins.
- **FR-010**: The documentation MUST describe the workflow, how to start it by hand, what each result means, and the 60-day schedule pause. That is `testdata/upstream/README.md` and the drift-report contract.
- **FR-011**: The `spec` binary and its tests MUST NOT gain any network access (Art. II).

### Key Entities

- **Pinned project:** the upstream repository, the pinned tag, and the template paths within the repository.
- **Candidate:** a release tag and commit of a pinned project, with its downloaded template files.
- **Drift result:** one per project: no change, same tag, unresolved changes (list), template not found, or fetch error.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An unresolved template change in a new upstream release is reported by a failed run within 7 days of the release, with no manual step.
- **SC-002**: The summary names the project, tag and changed shape keys, so a maintainer can start the "Moving a pin" procedure without rerunning anything locally.
- **SC-003**: In seeded tests, each failure kind (drift, missing template, fetch error) produces its own distinct result line, with no false drift.
- **SC-004**: The workflow has read-only permissions and creates no repository changes in 100% of runs (checked from the workflow definition).

## Assumptions

- GitHub's failed-run notification is enough reporting for a single-owner project. Opening an issue is out of scope, because it needs write permission (roadmap 010 notes).
- Template paths are the repository paths at the tag: `templates/tasks-template.md` for spec-kit and `templates/roadmap-template.md` for speckit-roadmap. On 2026-11-21, the spec-kit file at `v1.2.0` matched the pinned SHA-256, and its git blob was identical at `v1.4.3`.
- Weekly is often enough. The roadmap outcome is "within a week".
- The existing `TestDrift` behavior is unchanged: normalization, and the resolution list in `resolutions.md`. This feature automates invoking it and summarizing its result, and may refactor where the pins are listed (FR-009).
- Moving a pin stays a manual, decision-recorded procedure.
