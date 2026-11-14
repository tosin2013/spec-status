# Feature Specification: Two-mode roadmap and feature inventory

**Feature Branch**: `006-two-mode-specify`

**Created**: 2026-11-14

**Status**: Draft

**Revised**: 2026-11-14 during planning: FR-002 and FR-006 (explicit empty needs no new JSON field, research R3) and one edge case (folder names, research R7).

**Input**: Roadmap item 008, "Two-mode roadmap and feature inventory", with the owner's Q6 choices in decision D-2026-11-12-08 (`specs/001-spec-status/decisions.md`) and the direction in `specs/001-spec-status/two-mode-direction.md` (C-09). Targets v0.3.0 (D-2026-11-14-01).

## Context

Today `spec status` has two silent gaps:

- A roadmap file that is present but is not the supported template (for example a narrative document) reports **zero items, exit 0**. That looks exactly like a genuinely empty supported roadmap. README Limits records this.
- A Spec Kit project without the roadmap extension has no roadmap file. `spec status` stops with `SS_ROADMAP_UNAVAILABLE`, exit 2, and the user gets nothing about the features that do exist under `specs/`.

The owner chose (D-2026-11-12-08):

1. A separate command, `spec inventory`. `spec status` stays the roadmap report.
2. A missing roadmap still exits 2 with `SS_ROADMAP_UNAVAILABLE`; the message suggests `spec inventory`. There is no automatic switch.
3. A roadmap is recognized if it has at least one supported item heading, or the pinned speckit-roadmap template's `## Planned Specs` heading. Heading and no items: explicitly empty, exit 0. Neither: unrecognized, new fatal `SS_ROADMAP_UNRECOGNIZED`, exit 2.
4. Feature folders are immediate children of `specs/` named in either of Spec Kit's two forms.
5. Task checkboxes are counted with the existing tasks.md rules. Requirements checklists are excluded.
6. Inventory has its own JSON document, `inventory_schema_version: "1.0"`. Exit codes keep their meaning. `spec status` schema 1.0 is unchanged.
7. Ships in v0.3.0.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - An unrecognized roadmap is never reported as empty (Priority: P1)

A maintainer runs `spec status` in a repository whose `.specify/memory/roadmap.md` is a hand-written narrative, not the speckit-roadmap template. Instead of a clean, empty report, they get a fatal finding that names the file and says it is not a recognized roadmap, with a pointer to the supported template and to `spec inventory`.

**Why this priority**: A silent "zero items, exit 0" is a wrong result that CI would treat as success. This is the truthfulness defect the direction (C-09) exists to remove, and it changes existing behavior, so it must land first and be in the release notes.

**Independent Test**: Run `spec status` and `spec status --json` on three fixtures: a narrative file, a template-shaped file with `## Planned Specs` and no items, and a file with one supported item. Check the exit codes (2, 0, 0) and that only the first reports `SS_ROADMAP_UNRECOGNIZED`.

**Acceptance Scenarios**:

1. **Given** a readable roadmap with neither a supported item heading nor a `## Planned Specs` heading, **When** `spec status` runs, **Then** it exits 2 with fatal `SS_ROADMAP_UNRECOGNIZED`, naming the roadmap path, in both human text and JSON.
2. **Given** a roadmap with the `## Planned Specs` heading and no item headings, **When** `spec status` runs, **Then** it exits 0, reports zero items, says the roadmap is explicitly empty, and the completion denominators are null.
3. **Given** a roadmap with at least one supported item heading and no `## Planned Specs` heading (as hand-written roadmaps that follow the item grammar have), **When** `spec status` runs, **Then** the report is the same as today.
4. **Given** a roadmap that is recognized but contains malformed or unsupported item-like structure, **When** `spec status` runs, **Then** the existing malformed/unsupported diagnostics are reported exactly as today; it is never reported as unrecognized.
5. **Given** an unrecognized roadmap, **When** `spec status --fail-on none` runs, **Then** it still exits 2 (a fatal is never masked).

---

### User Story 2 - List features in a repository without a roadmap (Priority: P1)

A developer on a Spec Kit project that never installed the roadmap extension runs `spec status`, sees `SS_ROADMAP_UNAVAILABLE` with a suggestion, then runs `spec inventory`. They get each feature folder under `specs/`, whether its `spec.md`, `plan.md` and `tasks.md` are present, and checked/total task checkboxes, with a summary, in human text or JSON.

**Why this priority**: This is the new value for the large group of Spec Kit users without a roadmap; without it the only answer they get is a fatal error.

**Independent Test**: Run `spec inventory` and `spec inventory --json` on a fixture with two Spec Kit feature folders (one with all three files and a partially checked tasks.md, one with only spec.md), a `checklists/requirements.md` full of checked boxes, and a non-feature folder. Compare with a hand-counted expected result.

**Acceptance Scenarios**:

1. **Given** a repository with `specs/001-alpha/` (spec, plan, tasks with 3 of 5 boxes checked) and `specs/20261101-093000-beta/` (spec only), **When** `spec inventory` runs, **Then** both are listed in folder-name order with their file states, alpha shows 3/5 and beta shows tasks missing, and the summary shows 2 feature folders and 3/5 checked tasks.
2. **Given** a feature folder with `checklists/requirements.md` containing checked boxes, **When** inventory runs, **Then** those boxes are not counted anywhere.
3. **Given** `specs/notes/` and `specs/README.md`, **When** inventory runs, **Then** the folder is listed as not a feature folder (its contents not read) and the file is ignored as not a folder; neither is counted as a feature.
4. **Given** no `.specify/memory/roadmap.md`, **When** `spec status` runs, **Then** it exits 2 with `SS_ROADMAP_UNAVAILABLE` and its message suggests `spec inventory`; it does not run inventory itself.
5. **Given** a repository that also has a roadmap, **When** `spec inventory` runs, **Then** it reports the feature folders without reading the roadmap and does not infer status, priority or dependencies from it.

---

### User Story 3 - Inventory is as safe and plain as the roadmap report (Priority: P2)

A CI job and an AI agent consume `spec inventory --json`. They rely on the same guarantees as `spec status`: offline, read-only, deterministic bytes, root-confined, sensitive names never opened, located findings, exit codes 0/1/2 with `--fail-on`.

**Why this priority**: The new mode must not weaken any constitution guarantee (Art. I–III); it builds on P1 and P2 but is what makes the output trustworthy for automation.

**Independent Test**: Run inventory twice on the same fixture and compare bytes; run it on fixtures containing a symlinked feature folder, a symlinked tasks.md, a FIFO named tasks.md, and a sensitive-named file; check findings, file states and exit codes.

**Acceptance Scenarios**:

1. **Given** the same input, **When** inventory runs twice, **Then** the JSON and the human text (with color off and a fixed width) are byte-identical, with no timestamps.
2. **Given** a feature folder that is a symlink, or a `tasks.md` that is a symlink or not a regular file, **When** inventory runs, **Then** the item is reported as skipped with a located finding, its contents are not read, and it is counted in the skipped files total.
3. **Given** more feature folders or bytes than the existing input limits, **When** inventory runs, **Then** it stops with the existing input-limit fatal, exit 2, and reports no partial result as complete.
4. **Given** inventory found a warning, **When** run with the default threshold, **Then** it exits 0; with `--fail-on warning` it exits 1.
5. **Given** `--json`, **When** any invocation error occurs, **Then** the error is a single JSON document on stdout, as for `spec status`.

---

### Edge Cases

- A roadmap that is empty (0 bytes) or whitespace only: unrecognized (neither heading), exit 2. Its message says the file is empty.
- `## Planned Specs` inside a fenced code block or HTML comment does not count; only a visible level-2 heading does. Spelling and case must match the pinned template.
- A roadmap with `## Planned Specs` and only `###` grouping headings (no items): explicitly empty, exit 0.
- `--roadmap PATH` naming a readable, unrecognized file: same fatal as the default path.
- `specs/` missing: inventory exits 0, reports that the specs folder is missing, zero feature folders, null denominators. `specs/` present but not a directory, a symlink, or unreadable: a located finding and the same zero-folder result; never a guess.
- Two folders with the same number (`003-a`, `003-b`): both listed, no inference about which is current.
- A folder named `3-short` (fewer than three digits) or `003_name` (underscore): not a feature folder. `20261101-name` (a date without a time) has eight digits, so it is a sequential feature folder.
- A feature folder with no files at all: listed, all three files missing.
- A tasks.md with zero checkboxes: read, 0/0; the item contributes nothing to the denominator; overall checked-task percentage is null if no tasks file has any checkbox.
- Nested folders under a feature folder (`contracts/`, `checklists/`) are not read, and nothing inside them is counted.
- Untrusted folder names and file contents containing terminal escape sequences are sanitized in human text.
- Root discovery: inventory uses the same root selection as `spec status` (nearest ancestor with a non-symlink `.specify`, or `--root`).

## Requirements *(mandatory)*

### Functional Requirements

**Roadmap recognition (`spec status`)**

- **FR-001**: `spec status` MUST treat a readable roadmap as recognized when it contains at least one supported item heading, or a visible `## Planned Specs` level-2 heading (outside fences and comments, matching the pinned speckit-roadmap v0.1.0 template).
- **FR-002**: A recognized roadmap with no item headings MUST report as explicitly empty: exit 0, zero items, null completion denominators, and an explicit "empty roadmap" statement in human text. In JSON, a complete report with zero items MUST mean exactly this (unrecognized files no longer produce one), and the JSON contract MUST say so.
- **FR-003**: A readable roadmap that is not recognized MUST produce a new fatal diagnostic `SS_ROADMAP_UNRECOGNIZED` (error, fatal, exit 2) that names the roadmap path and points to the supported template and to `spec inventory`. It MUST NOT say the project is invalid Spec Kit.
- **FR-004**: Recognition MUST NOT change the result for any roadmap that has a supported item heading: existing reports, diagnostics, JSON bytes and exit codes stay identical, including all existing conformance fixtures and this repository's own roadmap.
- **FR-005**: When the roadmap is missing, `spec status` MUST keep exiting 2 with `SS_ROADMAP_UNAVAILABLE`, and its message MUST suggest `spec inventory`. It MUST NOT fall back to inventory, and an unrecognized or rejected roadmap MUST NOT fall back either.
- **FR-006**: `spec status` JSON MUST stay schema 1.0 with no field added or removed. The new code is additive; existing fields keep their meaning.

**Inventory command**

- **FR-007**: A new command `spec inventory` MUST exist with `--json`, `--root`, `--color` and `--fail-on`, with the same meanings, validation and defaults as on `spec status`. It MUST NOT accept roadmap-only options (`--roadmap`, the human-text filters, `--explain`).
- **FR-008**: Inventory MUST use the same root discovery as `spec status` and read only through the existing confined filesystem boundary.
- **FR-009**: Inventory MUST examine only immediate children of `<root>/specs/`. A child directory is a feature folder when its name is `NNN-name` (three or more ASCII digits, a hyphen, a non-empty name) or `YYYYMMDD-HHMMSS-name` (a date and time, a hyphen, a non-empty name). Other child directories MUST be listed as not feature folders, without reading their contents. Child files MUST be ignored.
- **FR-010**: For each feature folder, inventory MUST report the state of `spec.md`, `plan.md` and `tasks.md` as present, missing or skipped (with the reason: unsafe, sensitive, not a regular file, or limit), using the same states as task files in `spec status`.
- **FR-011**: Inventory MUST count `tasks.md` checkboxes with exactly the same rules as `spec status` counts linked tasks files, and MUST NOT count boxes in any other file, including `checklists/` content.
- **FR-012**: Inventory MUST report a summary: feature folders found, non-feature folders, files read/missing/skipped, and total checked/total task checkboxes with a percentage that is null when the total is zero.
- **FR-013**: Inventory MUST NOT infer or report lifecycle status, priority, dependencies, or completion from folder names, file presence, checkbox counts or `spec.md` metadata, and MUST say in human text that checked boxes are not verification.
- **FR-014**: Inventory MUST NOT read the roadmap. It MAY state whether a roadmap file exists at the default path, as a fact, without opening it.
- **FR-015**: Inventory JSON MUST be a separate document with `inventory_schema_version: "1.0"`, containing the mode, the root-relative specs path and its state, the folders (name, whether a feature folder, file states, checked and total), the summary, diagnostics, and the effective threshold. Field names and types are fixed in a contract before implementation.
- **FR-016**: Inventory findings (skipped or unsafe entries, missing or unusable `specs/`, non-feature folders) MUST be located diagnostics with the existing ordering and severity model; their codes and severities are fixed in the diagnostics contract. Exit codes MUST be 0 below `--fail-on`, 1 at or above it, 2 for fatal or invocation errors; `--fail-on none` never masks a fatal.
- **FR-017**: Folder order MUST be byte order of folder names; output MUST be deterministic with no timestamps.
- **FR-018**: Human text for inventory MUST follow the existing plain-surface rules (color only under the existing conditions, width from `COLUMNS` else 80, escapes sanitized), and MUST state which mode ran.

**Documentation and compatibility**

- **FR-019**: The diagnostics contract, JSON contract (with a new inventory contract), CLI contract, roadmap grammar and README MUST describe the new behavior, and README Limits MUST drop the "zero items, exit 0" limit. The README example output, if affected, MUST be updated.
- **FR-020**: The release notes for v0.3.0 MUST list the unrecognized-roadmap fatal as a pre-1.0 behavior change.
- **FR-021**: The agent skills extension's `reference/output.md` MUST list `SS_ROADMAP_UNRECOGNIZED` and any new inventory codes, so the existing extension test (every code in the diagnostics contract appears there) keeps passing. Broader skill changes for inventory are out of scope (item 013).
- **FR-022**: No new runtime dependency, network access, subprocess, write, cache or telemetry is added (Art. II); the import and side-effect tests MUST keep passing.

### Key Entities

- **Roadmap recognition result**: one of recognized-with-items, explicitly empty, unrecognized, or missing (malformed stays a diagnostic inside a recognized read).
- **Feature folder**: an immediate child of `specs/` matching a Spec Kit name form; has three file states and a checked/total count.
- **Non-feature folder**: any other immediate child directory; listed by name only.
- **Inventory report**: mode, specs folder state, folders, summary, diagnostics, threshold; rendered as human text or inventory JSON 1.0.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of hand-counted fixtures for missing, empty-supported, unrecognized (narrative, empty file, whitespace-only) and malformed roadmaps give the expected outcome and exit code; none of the unrecognized ones exits 0.
- **SC-002**: All existing tests, conformance fixtures and this repository's own `spec status` output are unchanged byte for byte.
- **SC-003**: On hand-counted inventory fixtures (both name forms, a non-feature folder, a requirements checklist, missing files, unsafe entries), every folder, file state and checkbox count matches the expected file exactly.
- **SC-004**: A user without a roadmap goes from the `SS_ROADMAP_UNAVAILABLE` message to a feature list in one command, with no option other than `spec inventory`.
- **SC-005**: Repeated inventory runs on the same input produce identical output bytes in 100% of runs.
- **SC-006**: A fixture shaped like the brownfield example in `two-mode-direction.md` (one feature folder, tasks 24/24, a 16-check requirements checklist) reports 24/24 and excludes all 16 checklist marks.

## Assumptions

- `--color` and `--fail-on` on `spec inventory` follow from "exit codes keep their meaning" and the existing plain-surface rules; `--fail-on` defaults to `error`.
- Inventory works whether or not a roadmap exists; it is explicit, so running it in a roadmap repository is allowed and simply ignores the roadmap.
- A missing `specs/` folder is not fatal for inventory: it is a reported state with zero folders and exit 0, because a fresh Spec Kit project legitimately has none.
- Duplicate feature numbers are listed without a finding; resolving them is a human decision.
- Only `spec.md`, `plan.md` and `tasks.md` are inspected; other files (research, contracts, quickstart) are out of scope.
- Fixtures are synthetic, shaped like the evidence repositories in `two-mode-direction.md`; no external repository content is copied, and nothing is fetched at runtime.
- The existing input limits (files, bytes, path length) apply unchanged to inventory.
- Out of scope: inferring priority, dependencies or lifecycle from folders; counting requirements checklists; network or upstream lookups; other roadmap template versions; changes to agent skills beyond the output reference.
