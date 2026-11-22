# Specification Quality Checklist: Per-phase task counts

**Purpose**: Validate specification completeness and quality before proceeding to planning
**Created**: 2026-11-21
**Feature**: [spec.md](../spec.md)

## Content Quality

- [x] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

## Requirement Completeness

- [x] No [NEEDS CLARIFICATION] markers remain
- [x] Requirements are testable and unambiguous
- [x] Success criteria are measurable
- [x] Success criteria are technology-agnostic (no implementation details)
- [x] All acceptance scenarios are defined
- [x] Edge cases are identified
- [x] Scope is clearly bounded
- [x] Dependencies and assumptions identified

## Feature Readiness

- [x] All functional requirements have clear acceptance criteria
- [x] User scenarios cover primary flows
- [x] Feature meets measurable outcomes defined in Success Criteria
- [x] No implementation details leak into specification

## Notes

- The JSON field names (`tasks.phases`, `heading`, `line`, `story`, `checked`, `total`) and Markdown heading forms are the user-facing contract of a CLI, not implementation details. Earlier specs (006–008) name JSON fields and diagnostic codes the same way.
- No clarification markers: the owner settled the scope in D-2026-11-21-01, and issue #66 fixes the shape. Choices made here: level-1 headings also end a phase; setext headings are a listed limit; `phases` is `null` for unread files; a story number has 1–9 digits.
