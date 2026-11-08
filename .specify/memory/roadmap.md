<!--
SYNC IMPACT REPORT
==================
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

## Planned Specs

### Core

### 001 — Read-only spec status core [status: implemented]

- **Description:** The development `spec status` CLI: roadmap and tasks parsing,
  supersession and dependency analysis, text and schema 1.0 JSON, exits 0/1/2,
  Lip Gloss summary panel, human-text filters, and `--explain`.
- **Outcome:** Any supported roadmap produces the same report in text and JSON,
  with no runtime side effects, and fails visibly on unsupported input.
- **Scope (in):** FS-01–FS-06 as implemented under PROFILE-R1 reversible
  defaults; speckit-roadmap bold field labels (owner-approved DC-01 refinement).
- **Scope (out):** release artifacts, write-back, inferred quarters, dashboards.
- **Depends on:** none
- **Governed by:** C-01, C-02, C-03, C-04, C-05
- **Spec dir:** specs/001-spec-status
- **Notes:** Recorded implemented, not verified: tests were run by the
  implementer on Linux amd64 and Darwin arm64 and have not been independently
  reviewed. PROFILE-R1 choices remain reversible until 002 lands.

### 002 — Contract approval DC-01–DC-08 [status: needs-info]

- **Description:** Gate 0 (TASK-01–06): owner decisions on grammar, metrics,
  CLI/JSON, safety, limits, platforms, and dogfood inputs that PROFILE-R1
  currently assumes.
- **Outcome:** Each DC item is approved, amended, or explicitly rejected, and the
  contracts, fixtures, and traceability are reconciled with those decisions.
- **Scope (in):** DC-01 through DC-08 in `specs/001-spec-status/decisions.md`.
- **Scope (out):** implementing changes the decisions require (tracked per spec).
- **Depends on:** none
- **Governed by:** C-01, C-03, C-05
- **Notes:** Approved on 2026-11-07: one DC-01 refinement (speckit-roadmap field
  spelling), all of DC-03 (C-07), and the DC-08 dogfood inputs. Still open: the
  rest of DC-01, DC-02, DC-04–DC-07, and DC-08's license, security route and
  support window.

### Evidence and distribution

### 003 — Real-roadmap dogfood [status: planned]

- **Description:** TASK-16 / AC-24: run `spec status` on approved real Spec Kit
  roadmaps and compare against a hand count.
- **Outcome:** Hand-count parity on at least two owner-approved roadmaps, with any
  disagreement recorded as a finding rather than adjusted away.
- **Scope (in):** this repository's own roadmap; the speckit-roadmap upstream
  fixture; one owner-named portfolio roadmap (DC-08).
- **Scope (out):** certifying upstream compatibility.
- **Depends on:** 001
- **Governed by:** C-01
- **Notes:** Inputs approved under D-2026-11-07-08: this repository's roadmap and
  the speckit-roadmap v0.1.0 roadmap in `testdata/speckit-roadmap-upstream`. On
  2026-11-07 the implementer observed hand-count parity on both (recorded in
  `specs/001-spec-status/tasks.md`). Not an independent review. The outcome looks
  met; a status change awaits owner confirmation and a spec link, which C-07
  requires from specced onward.

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
- **Depends on:** 002, 003
- **Governed by:** C-06
- **Notes:** Blocked on DC-07 (targets, runners, identity policy) and on a license
  choice; no license is selected today.

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

## Cross-Cutting Notes

- This roadmap is read by this repository's own `spec status`; keep `Depends on`
  lists as bare three-digit IDs.
- Grouping headings (`### Core`) are allowed; item headings must keep the
  `### NNN — Title [status: token]` shape.

---

**Version**: 1.1.0 | **Ratified**: 2026-11-07 | **Last Amended**: 2026-11-07
