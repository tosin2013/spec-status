<!--
SYNC IMPACT REPORT
==================
Version change: (none) → 1.0.0
Bump rationale: initial ratification.

Changes this revision:
  - Added spec 001 — Sample feature
  - Added spec 004 — Verified base

Specs affected: 001, 002, 003, 004
Open questions added/resolved: Q1 added

Notes: Conformance fixture built from speckit-roadmap v0.1.0 templates/roadmap-template.md.
-->

# Sample Project — Spec Roadmap

Living, non-binding map of the specs planned for Sample Project. It is **not a
commitment to order or scope**. Foundations: the project [constitution](constitution.md).

Status legend (lifecycle): **undecided** · **needs-info** · **planned** ·
**specced** · **in-progress** · **implemented** · **verified** · **deferred** ·
**abandoned**.

---

## Vision & End States

<!-- The big-picture goals and outcomes this project is driving toward. -->

- A working sample that exercises every roadmap construct.

## Constraints & Decisions

<!-- Cross-cutting constraints and the durable "why". -->

- **C-01 — Read-only:** the sample never writes files.
- **C-02 — Narrow input:** only documented shapes are used.

## Planned Specs

<!-- THE LEDGER. The spine of the roadmap. Optional grouping headers (###) may
     organize entries by theme/phase. -->

### Core

### 001 — Sample feature  [status: implemented]

- **Description:** A sample feature with every field filled in.
- **Outcome:** The sample runs end to end and its tasks are tracked in a
  tasks file generated from the Spec Kit template.
- **Scope (in):** Parsing, reporting.
- **Scope (out):** Writing files.
- **Depends on:** 004 (verified).
- **Governed by:** C-01, C-02
- **Addresses:** PRD goal 1
- **Spec dir:** specs/001-sample/
- **Notes:** Implemented, not yet verified.

### 002 — Follow-up  [status: planned]

- **Description:** Follow-up work on the sample.
- **Outcome:** _to be defined_
- **Scope (in):** Reporting improvements.
- **Depends on:** 001.
- **Notes:** Waits for 001 to be verified.

### Later

### 003 — Deferred idea  [status: deferred]

- **Description:** An idea kept for later.
- **Outcome:** _to be defined_
- **Scope (in):** _to be defined_
- **Notes:** Deferred until demand appears.

### 004 — Verified base  [status: verified]

- **Description:** The base the sample builds on.
- **Outcome:** The base is in place and reviewed.
- **Scope (in):** Base structure.
- **Depends on:** none
- **Spec dir:** specs/004-base

## Open Questions

<!-- Unresolved items not yet attached to a specific spec. -->

- Q1 — Should the sample cover more statuses?

## Cross-Cutting Notes

<!-- Architecture-level notes that span multiple specs. -->

- Every construct comes from the pinned template.

---

**Version**: 1.0.0 | **Ratified**: 2026-11-07 | **Last Amended**: 2026-11-07
