<!--
SYNC IMPACT REPORT
==================
Version change: 1.3.9 → 1.4.0
Bump rationale: MINOR — new items 009, 010, 011.

Changes this revision:
  - Added 009 Dependency update automation (specced, specs/003-dependency-updates)
  - Added 010 Upstream drift automation and 011 README sample test (planned;
    deferred from feature 003)
  - 004: depends on 009; notes updated for the DC-07/DC-08 decisions

Specs affected: 003-dependency-updates (new)
Open questions added/resolved: none

Previous revision (1.3.9, 2026-11-11):
Version change: 1.3.8 → 1.3.9
Bump rationale: PATCH — status transition.

Changes this revision:
  - 002 status: needs-info → verified (D-2026-11-11-02, D-2026-11-11-03;
    contracts reconciled); Spec dir specs/001-spec-status added (C-07)

Specs affected: 001
Open questions added/resolved: none

Previous revision (1.3.8, 2026-11-11):
Version change: 1.3.7 → 1.3.8
Bump rationale: PATCH — wording: 002 notes record D-2026-11-11-02.

Changes this revision:
  - 002 notes: DC-01–DC-08 decisions recorded; status unchanged until the
    contracts are reconciled

Specs affected: none
Open questions added/resolved: none

Previous revision (1.3.7, 2026-11-11):
Version change: 1.3.6 → 1.3.7
Bump rationale: PATCH — status transition.

Changes this revision:
  - 003 status: planned → verified (D-2026-11-11-01); Spec dir added; the portfolio
    roadmap moved to Scope (out)

Specs affected: 001
Open questions added/resolved: none

Previous revision (1.3.6, 2026-11-11):
Version change: 1.3.5 → 1.3.6
Bump rationale: PATCH — wording: 007 notes and Governed by (debrief F3).

Changes this revision:
  - 007 Governed by adds C-10; notes record that reviews 8–9, C-10 and the
    T026 rounds were classified under 007

Specs affected: none
Open questions added/resolved: none

Previous revision (1.3.5, 2026-11-09):
Version change: 1.3.4 → 1.3.5
Bump rationale: PATCH — wording: 003 notes record the 2026-11-09 recount.

Changes this revision:
  - 003 notes: recount of roadmap 1.3.4 and the upstream fixture; 1.3.4
    snapshot added as a conformance fixture; status unchanged

Specs affected: none
Open questions added/resolved: none

Previous revision (1.3.4, 2026-11-08):
Version change: 1.3.3 → 1.3.4
Bump rationale: PATCH — status transitions.

Changes this revision:
  - 001 status: implemented → verified (D-2026-11-08-05; three T026 review
    rounds and their fixes; residual risk in README Limits)
  - 007 status: implemented → verified (T026 complete)

Specs affected: 001, 007
Open questions added/resolved: none

Previous revision (1.3.3, 2026-11-08):
Version change: 1.3.2 → 1.3.3
Bump rationale: PATCH — wording: C-10 notes the closed field body
(D-2026-11-08-04).

Specs affected: 001
Open questions added/resolved: none

Previous revision (1.3.2, 2026-11-08):
Version change: 1.3.1 → 1.3.2
Bump rationale: PATCH — wording: C-08 notes its amendment D-2026-11-08-03.

Changes this revision:
  - C-08: rule (b) means any block layout; markup-encoded labels are a README
    limit (D-2026-11-08-03, after the T026 review)

Specs affected: 001, 007
Open questions added/resolved: none

Previous revision (1.3.1, 2026-11-08):
Version change: 1.3.0 → 1.3.1
Bump rationale: PATCH — status transition only.

Changes this revision:
  - 007 status: in-progress → implemented (owner, 2026-11-08, after the
    post-implementation debrief; all tasks done except review-only T026)

Specs affected: 007
Open questions added/resolved: none

Previous revision (1.3.0, 2026-11-08):
Version change: 1.2.0 → 1.3.0
Bump rationale: MINOR — new spec 008 (two-mode roadmap and inventory) and new
decisions C-09 (D-2026-11-08-01) and C-10 (D-2026-11-08-02).

Changes this revision:
  - Added spec 008 — Two-mode roadmap and feature inventory (needs-info)
  - Decisions recorded: C-09 (two-mode direction), C-10 (line-level label rule)
  - 001 now also governed by C-10
  - Open question added: Q6 (two-mode open choices)

Specs affected: 001, 008
Open questions added/resolved: Q6 added

Previous revision (1.2.0, 2026-11-07):
Version change: 1.1.0 → 1.2.0
Bump rationale: MINOR — new planned spec 007 (upstream-format conformance)
and new decision C-08 (acceptance bar D-2026-11-07-09).

Changes this revision:
  - Added spec 007 — Upstream-format conformance (in-progress)
  - Decision recorded: C-08 (D-2026-11-07-09 replaces D-2026-11-07-06)

Specs affected: 007
Open questions added/resolved: none

Previous revision (1.1.0, 2026-11-07):
Version change: 1.0.0 → 1.1.0
Bump rationale: MINOR — new decision C-07 (DC-03 approved) and DC-08 dogfood
inputs recorded; Q1 and Q2 resolved. No status transitions.

Changes this revision:
  - Decision recorded: C-07 (D-2026-11-07-03, DC-03 as built)
  - 002 notes: DC-03 approved; DC-08 dogfood inputs approved (D-2026-11-07-08)
  - 003 notes: hand-count parity observed on both approved inputs
  - Open questions resolved: Q1 (dogfood inputs), Q2 (stage matrix/blocked policy)

Specs affected: 002, 003
Open questions added/resolved: Q1, Q2 resolved

Previous revision (1.0.0, 2026-11-07): initial ratification.
Changes in 1.0.0:
  - Added spec 001 — Read-only spec status core (implemented)
  - Added spec 002 — Contract approval DC-01–DC-08 (needs-info)
  - Added spec 003 — Real-roadmap dogfood (planned)
  - Added spec 004 — Release engineering (needs-info)
  - Added spec 005 — Phase-two write-back (deferred)
  - Added spec 006 — Deferred products (deferred)
  - Decisions recorded: C-01 … C-06 from constitution Articles I–VIII

Specs affected: 001, 002, 003, 004, 005, 006
Open questions added/resolved: Q1–Q5 added

Notes: Harvested from .specify/memory/constitution.md, specs/001-spec-status/
(spec, decisions, tasks, acceptance, build-release, phase2) and README.md, then
confirmed by the owner on 2026-11-07. Statuses are recorded labels; this roadmap
does not claim independent verification. Depends-on lists use bare IDs so that
`spec status` reads them without annotation.
-->

# Spec Status — Spec Roadmap

Living, non-binding map of the specs planned for Spec Status. It is **not a
commitment to order or scope** — it captures the decisions, outcomes, and
constraints already recorded in the constitution and feature documents so they
are not lost before the spec that needs them is written. Specs are scoped and
clarified when they are actually started.
Foundations: the project [constitution](constitution.md) and the feature
documents under `specs/001-spec-status/`.

Status legend (lifecycle): **undecided** · **needs-info** · **planned** ·
**specced** · **in-progress** · **implemented** · **verified** · **deferred** ·
**abandoned**.

---

## Vision & End States

- A developer, reviewer, planner, calling LLM, or CI job gets a deterministic,
  read-only reading of a Spec Kit roadmap — recorded statuses, checkbox progress,
  dependency blockers, supersession history, and broken references — as text or
  versioned JSON from one result.
- The report never claims more than the files show: recorded `verified` is a
  label, not a review; exit 0 is "below threshold", not "done".
- An eventual released binary runs without Go, root, AI, network, or the roadmap
  extension, on platforms that real runners have proven.

## Constraints & Decisions

- **C-01 — Truthful evidence (Art. I):** report recorded state and mechanical
  inconsistencies only; keep checked tasks, recorded implemented, and recorded
  verified separate; state denominators and skipped inputs.
- **C-02 — Deterministic, standalone, side-effect-free (Art. II):** one result
  feeds text and JSON; no network, subprocess, Git, AI, credentials, or writes at
  runtime; `CGO_ENABLED=0`. Parsing and confined reads use the standard library.
  Owner-approved dependency exceptions: Cobra (commands/flags) and Lip Gloss
  (human text only).
- **C-03 — Narrow, safe input (Art. III):** a pinned heading/bullet dialect, root
  confinement, symlink and sensitive-name rejection, bounded work; unsupported
  item-like input fails visibly.
- **C-04 — History cannot become completion (Art. IV):** superseded items stay in
  history; every dependency on a superseded item requires human review and is
  never redirected.
- **C-05 — Compatibility (Art. V–VI):** `--json`, raw `superseded_by`, default
  threshold `error`, exits 0/1/2; no inferred quarters in v1.
- **C-06 — Evidence-gated distribution and mutation (Art. VII–VIII):** no support,
  signing, or provenance claims without verified evidence; any future write-back
  needs explicit per-mutation approval and is never part of the v1 core.
- **C-07 — Dependency and stage policy (DC-03, approved 2026-11-07):** a
  dependency is satisfied only by a recorded `verified` target; specced needs a
  spec link and spec.md, in-progress and later also need plan.md and tasks.md;
  human review of a superseded dependency is a warning; blocked and review
  counts cover active rows once each. Decision D-2026-11-07-03 in
  `specs/001-spec-status/decisions.md`.
- **C-08 — Upstream-pinned acceptance bar (D-2026-11-07-09, 2026-11-07):**
  supported input is what pinned github/spec-kit v1.2.0 and
  srobroek/speckit-roadmap v0.1.0 templates produce; a blocking defect is a
  wrong result on those shapes or a silently misread `Depends on`, `Spec dir`,
  or `superseded-by`; other exotic Markdown fails visibly or is a README
  limit. Replaces D-2026-11-07-06. Amended by D-2026-11-08-03 (2026-11-08):
  rule (b) covers any block layout; labels encoded with markup are a README
  limit.
- **C-09 — Two-mode direction (D-2026-11-08-01, 2026-11-08):** supported
  roadmaps are a finite set of versioned templates; missing, empty,
  unrecognized, and malformed roadmaps stay distinct; a roadmap-free feature
  inventory is a separate explicit mode and never a silent fallback;
  requirements checklists are excluded. Direction only; recorded in
  `specs/001-spec-status/two-mode-direction.md`.
- **C-10 — Line-level label rule (D-2026-11-08-02, 2026-11-08):** a visible
  roadmap line containing a `Depends on`, `Spec dir`, or `superseded-by` label
  with a colon must be an accepted field line under an item, otherwise the
  report is fatal. Closed field body (D-2026-11-08-04): an indented line under
  such a field joins its value or is fatal.

## Planned Specs

### Core

### 001 — Read-only spec status core [status: verified]

- **Description:** The development `spec status` CLI: roadmap and tasks parsing,
  supersession and dependency analysis, text and schema 1.0 JSON, exits 0/1/2,
  Lip Gloss summary panel, human-text filters, and `--explain`.
- **Outcome:** Any supported roadmap produces the same report in text and JSON,
  with no runtime side effects, and fails visibly on unsupported input.
- **Scope (in):** FS-01–FS-06 as implemented under PROFILE-R1 reversible
  defaults; speckit-roadmap bold field labels (owner-approved DC-01 refinement).
- **Scope (out):** release artifacts, write-back, inferred quarters, dashboards.
- **Depends on:** none
- **Governed by:** C-01, C-02, C-03, C-04, C-05, C-10
- **Spec dir:** specs/001-spec-status
- **Notes:** Recorded verified on 2026-11-08 by owner decision D-2026-11-08-05:
  three independent T026 review rounds by separate AI agent sessions, every
  blocking finding fixed, no round with zero findings; residual risk in README
  Limits. Not a human review or an upstream certification. PROFILE-R1 choices remain reversible until 002 lands.

### 002 — Contract approval DC-01–DC-08 [status: verified]

- **Description:** Gate 0 (TASK-01–06): owner decisions on grammar, metrics,
  CLI/JSON, safety, limits, platforms, and dogfood inputs that PROFILE-R1
  currently assumes.
- **Outcome:** Each DC item is approved, amended, or explicitly rejected, and the
  contracts, fixtures, and traceability are reconciled with those decisions.
- **Scope (in):** DC-01 through DC-08 in `specs/001-spec-status/decisions.md`.
- **Scope (out):** implementing changes the decisions require (tracked per spec).
- **Depends on:** none
- **Governed by:** C-01, C-03, C-05
- **Spec dir:** specs/001-spec-status
- **Notes:** Approved on 2026-11-07: one DC-01 refinement (speckit-roadmap field
  spelling), all of DC-03 (C-07), and the DC-08 dogfood inputs. Approved on
  2026-11-11 (D-2026-11-11-02): the rest of DC-01, DC-02, DC-04–DC-06 as built;
  DC-07 linux/amd64 and darwin/arm64 with checksum-only archives; DC-08
  Apache-2.0, GitHub private vulnerability reporting, latest release only.
  DC-08 hosting on 2026-11-11 (D-2026-11-11-03): GitHub only, no mirror;
  CONTRIBUTING.md and CODE_OF_CONDUCT.md added. Contracts reconciled with
  these decisions in PR #25. Recorded verified on 2026-11-11: every DC item
  is decided and the contracts match; the toolchain, runners and release
  workflow belong to 004.

### 008 — Two-mode roadmap and feature inventory [status: needs-info]

- **Description:** TASK-24–26: distinguish an empty supported roadmap from an
  unrecognized file, and add a separate, explicit inventory of feature folders
  for repositories without a roadmap.
- **Outcome:** An unrecognized roadmap is never reported as an empty supported
  one, and a repository without a roadmap can list its features, file
  presence, and task checkbox counts without any silent fallback.
- **Scope (in):** Contracts for the chosen command, version identification,
  folder and checkbox rules, inventory schema and exits; then the two modes
  with tests for missing, empty, unrecognized, and malformed roadmaps.
- **Scope (out):** Inferring priority, dependencies, or lifecycle from
  folders; counting requirements checklists; network or upstream lookups.
- **Depends on:** 001
- **Governed by:** C-01, C-02, C-08, C-09
- **Notes:** Needs the owner's choices in Q6 before any contract or code
  work; README Limits records today's zero-item result.

### Evidence and distribution

### 003 — Real-roadmap dogfood [status: verified]

- **Description:** TASK-16 / AC-24: run `spec status` on approved real Spec Kit
  roadmaps and compare against a hand count.
- **Outcome:** Hand-count parity on at least two owner-approved roadmaps, with any
  disagreement recorded as a finding rather than adjusted away.
- **Scope (in):** this repository's own roadmap; the speckit-roadmap upstream
  fixture.
- **Scope (out):** certifying upstream compatibility; a portfolio roadmap (may
  be added later under DC-08).
- **Depends on:** 001
- **Governed by:** C-01
- **Spec dir:** specs/001-spec-status
- **Notes:** Inputs approved under D-2026-11-07-08: this repository's roadmap and
  the speckit-roadmap v0.1.0 roadmap in `testdata/speckit-roadmap-upstream`. On
  2026-11-07 the implementer observed hand-count parity on both (recorded in
  `specs/001-spec-status/tasks.md`). Not an independent review. On 2026-11-09
  the implementer recounted this repository's roadmap at 1.3.4 (8 items, 6
  active, 2/6 verified, 1 blocker, 30/30 tasks in 2 files) and the upstream
  fixture; both match. Roadmap 1.3.4 is now a hand-counted conformance fixture
  (`testdata/conformance/this-repo-1.3.4`). The owner confirmed it done on
  2026-11-11 ({D}); counts are the implementer's, not an independent
  review. Links `specs/001-spec-status`, where TASK-16 and AC-24 live.

### 009 — Dependency update automation [status: specced]

- **Description:** Dependabot for Go modules and workflow components (security
  fixes on advisory, one weekly grouped minor/patch proposal, majors separate),
  a pull-request and main check (tests on linux/amd64 and darwin/arm64, go vet,
  govulncheck), a merge rule under the dependency allowlist, and a SECURITY.md
  dependency section.
- **Outcome:** Every pull request is tested before merge, vulnerable dependencies
  are reported and fix proposals arrive tested, and automated updates cannot
  add third-party modules without a recorded owner decision.
- **Scope (in):** update configuration, test check, merge rule (D-2026-11-11-04),
  SECURITY.md, CONTRIBUTING.md, README and build-release.md wording.
- **Scope (out):** release, tag and publication workflows; toolchain selection;
  signing; automatic merging.
- **Depends on:** 002
- **Governed by:** C-01, C-06
- **Spec dir:** specs/003-dependency-updates
- **Notes:** Owner choices of 2026-11-11. Dependabot alerts and automatic
  security fixes were enabled in repository settings the same day. No change
  to the `spec` binary.

### 004 — Release engineering [status: needs-info]

- **Description:** TASK-17–21: support matrix, pinned toolchain,
  test-before-publish tag workflow, archives and checksums, runtime smoke on each
  advertised target, install/upgrade docs, owner-approved release candidate.
- **Outcome:** A first stable tag created by the owner from a tested release
  candidate, with evidence for every advertised platform.
- **Scope (in):** build-release.md gates; DC-07 decisions; license and security
  disclosure choices.
- **Scope (out):** signing or provenance unless separately approved; package
  managers.
- **Depends on:** 002, 003, 009
- **Governed by:** C-06
- **Notes:** DC-07 and DC-08 are decided (D-2026-11-11-02, D-2026-11-11-03):
  linux/amd64 and darwin/arm64, checksum-only archives, Apache-2.0. Open: the
  toolchain version, runner evidence and the release workflow. Reuses the test
  check from 009.

### 007 — Upstream-format conformance [status: verified]

- **Description:** Define supported input as the shapes produced by pinned
  upstream templates, with hand-counted conformance fixtures, a drift check
  for upstream changes, and the acceptance bar D-2026-11-07-09.
- **Outcome:** Counts on upstream-shaped roadmaps and tasks files are proven
  by fixtures, supported versions are published, upstream changes are adopted
  deliberately, and 001 can be judged against a finite bar.
- **Scope (in):** Pins, shape inventory, conformance fixtures, drift check,
  reclassification of review findings, README support statement.
- **Scope (out):** Runtime network access or upstream lookups; release work.
- **Depends on:** 001
- **Governed by:** C-01, C-02, C-03, C-08, C-10
- **Spec dir:** specs/002-upstream-format-conformance
- **Notes:** Created with /speckit-specify, /speckit-plan, and /speckit-tasks;
  implementation approved by the owner on 2026-11-07. Recorded implemented
  on 2026-11-08; recorded verified the same day when T026 completed under
  D-2026-11-08-05. Reviews 8 and 9, the line-level label rule (C-10), and the
  three T026 rounds changed 001's parser but were classified in this
  feature's `reclassification.md`. Roadmap 1.3.4 is a hand-counted fixture
  (`testdata/conformance/this-repo-1.3.4`).

### 010 — Upstream drift automation [status: planned]

- **Description:** A scheduled workflow that fetches the latest spec-kit and
  speckit-roadmap templates and runs the conformance drift check (`TestDrift`)
  against the pins.
- **Outcome:** An upstream template change is noticed within a week, so the
  pin moves deliberately (007's intent), not by surprise.
- **Scope (in):** schedule, which releases to fetch, how a drift is reported.
- **Scope (out):** moving pins automatically; runtime network access.
- **Depends on:** 007, 009
- **Governed by:** C-02, C-08
- **Notes:** Deferred from feature 003 by the owner on 2026-11-11. Needs its own
  spec; reporting by issue would need write permission.

### 011 — README sample test [status: planned]

- **Description:** A Go test that compares the README's example output with
  `spec status` on `testdata/example-roadmap`.
- **Outcome:** The README sample cannot drift from real output unnoticed.
- **Scope (in):** the human-text and JSON samples in README.md.
- **Scope (out):** other documentation checks.
- **Depends on:** 999
- **Governed by:** C-01
- **Notes:** Deferred from feature 003 by the owner on 2026-11-11.

### Future

### 005 — Phase-two write-back [status: deferred]

- **Description:** TASK-22: optional extension-mediated mutation of roadmap
  statuses, with per-ID, per-change, per-path owner approval and persisted reread.
- **Outcome:** _to be defined_
- **Scope (in):** _to be defined_
- **Depends on:** 002
- **Governed by:** C-06
- **Notes:** Needs separate scope approval and a stable core contract;
  `specs/001-spec-status/phase2.md` records the boundary.

### 006 — Deferred products [status: deferred]

- **Description:** TASK-23: dashboards, trends, notifications, digests,
  multi-repo aggregation, package-manager channels.
- **Outcome:** _to be defined_
- **Scope (in):** _to be defined_
- **Depends on:** none
- **Governed by:** C-02
- **Notes:** Listed so the ideas are not lost; listing them is not a commitment.

## Open Questions

- ~~**Q1 —** Which roadmaps besides this one are approved dogfood inputs
  (DC-08)?~~ Resolved 2026-11-07: this repository and speckit-roadmap v0.1.0.
- ~~**Q2 —** Which lifecycle-file matrix and blocked-count policy should be
  ratified (DC-03)?~~ Resolved 2026-11-07: see C-07.
- **Q3 —** Initial targets: linux/amd64 + darwin/arm64 only, and on which runners
  (DC-07)?
- **Q4 —** License, security disclosure route, and support window (DC-08 / 004).
- **Q5 —** Should the speckit-roadmap annotated dependency style
  (`003-brief (verified)`) gain any further support, or stay a visible
  SS_DEPENDENCY_MALFORMED finding?

- **Q6 —** For 008: command names and flags, default mode, how a template
  version is identified, feature-folder and checkbox rules, inventory JSON
  schema and exit codes, and whether it ships in v1.

## Cross-Cutting Notes

- This roadmap is read by this repository's own `spec status`; keep `Depends on`
  lists as bare three-digit IDs.
- Grouping headings (`### Core`) are allowed; item headings must keep the
  `### NNN — Title [status: token]` shape.

---

**Version**: 1.4.0 | **Ratified**: 2026-11-07 | **Last Amended**: 2026-11-11
