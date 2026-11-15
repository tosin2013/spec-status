# Feature Specification: README sample test

**Feature Branch**: `011-013-specify`

**Created**: 2026-11-15

**Status**: Draft

**Input**: Roadmap item 011, "README sample test": "A Go test that compares the README's example output with `spec status` on `testdata/example-roadmap`. Outcome: the README sample cannot drift from real output unnoticed." Deferred from feature 003 on 2026-11-11; scheduled for v0.3.0 by D-2026-11-14-01 (5).

## Context

The README shows three blocks of real `spec` output (measured 2026-11-15):

| README block | Command shown | Checked by a test today? | Current state |
|---|---|---|---|
| "Example" | `spec status --color=never --root testdata/example-roadmap` | no | matches real output (hand-updated in feature 008) |
| "Filters" | `--human-review --explain` on the same example | no | **stale**: missing the `statuses` panel line added in feature 008 |
| "Without a roadmap" | `spec inventory --color=never --root testdata/two-mode/no-roadmap` | yes (feature 006, `TestReadmeInventorySample`) | matches |

The stale "Filters" block is exactly the drift this item exists to prevent: a feature changed the output, and the one sample nobody checked was missed. The README has no JSON output sample today.

CLAUDE.md already asks that the README sample be updated when output changes; this feature makes that a failing test instead of a reminder.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - A changed output fails until the README is updated (Priority: P1)

A contributor changes how `spec` prints something. The test suite fails, naming the README sample that no longer matches and showing the difference, so the README is updated in the same change.

**Why this priority**: it is the whole outcome of item 011.

**Independent Test**: change one character in a README sample; the suite fails naming that sample; restore it; the suite passes.

**Acceptance Scenarios**:

1. **Given** every labeled README sample matches real output, **When** the test suite runs, **Then** it passes.
2. **Given** a sample differs from real output by one character, **When** the suite runs, **Then** it fails, naming the sample's command and README line, and showing the expected and actual text.
3. **Given** the stale "Filters" sample today, **When** this feature lands, **Then** the sample is corrected and the test passes.

---

### User Story 2 - Every output sample is covered, not just one (Priority: P1)

A maintainer adds a new output sample to the README. They label it with the exact command it shows, and the test checks it automatically; an unlabeled block that looks like `spec` output fails the test, so no sample can silently escape checking.

**Why this priority**: the stale "Filters" block shows that checking a single named sample is not enough.

**Independent Test**: add an unlabeled block starting with the panel's top border; the suite fails asking for a label.

**Acceptance Scenarios**:

1. **Given** a README output block, **When** it is preceded by a label naming the `spec` command (subcommand and flags) that produces it, **Then** the test runs that command against the repository's fixtures and compares the block byte for byte.
2. **Given** a fenced block that looks like `spec` human output (it starts with the panel's top border line or with `spec: `) but has no label, **When** the suite runs, **Then** it fails naming the README line.
3. **Given** a label whose command uses an option the test does not support, or a root outside `testdata/`, **When** the suite runs, **Then** it fails clearly instead of skipping the sample.
4. **Given** the existing inventory sample, **When** this feature lands, **Then** it is checked by the same mechanism, and the separate feature-006 test is retired.

---

### Edge Cases

- Samples are rendered as in the README: color off, 80 columns, development-build heading. The label carries the command; width and color are fixed by the test, not by the label.
- A label must sit immediately before its fenced block (blank lines allowed); a label with no following block fails.
- The label is invisible in rendered Markdown (it does not change what readers see).
- Exit-code text is not part of a sample; only the printed bytes are compared.
- Fatal samples (`spec: CODE: message`) are compared from the error stream, if one is ever added.
- The agent-prompt block, shell snippets and YAML snippets are not output samples and are not checked.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Each README block of `spec` output MUST be preceded by a machine-readable, render-invisible label giving the exact `spec` subcommand and flags that produce it.
- **FR-002**: A test MUST run each labeled command in-process against the repository's fixtures, with color off, 80 columns and the development-build heading, and compare the output to the block byte for byte.
- **FR-003**: On mismatch, the test MUST name the command and README line and show both texts.
- **FR-004**: A fenced block that looks like `spec` human output without a label MUST fail the test; so MUST a label with no block, an unsupported option, or a root outside `testdata/`.
- **FR-005**: The three current samples MUST be labeled, and the stale "Filters" sample corrected from a real run.
- **FR-006**: The feature-006 inventory sample test MUST be replaced by the shared mechanism, keeping its coverage.
- **FR-007**: The test MUST be standard library only, run within `go test ./...` and need no network, subprocess or files outside the repository (Art. II; it is test code, not runtime).
- **FR-008**: CLAUDE.md's note about updating the README sample MUST point to the test.

### Key Entities

- **Sample label**: the command a README block shows.
- **Sample**: the fenced block that follows a label.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of README `spec` output blocks are labeled and checked (three today).
- **SC-002**: Changing any single character of any sample makes `go test ./...` fail with a message naming that sample.
- **SC-003**: Adding an unlabeled output-looking block makes the suite fail.
- **SC-004**: After this feature, the README samples match real output byte for byte, including the corrected "Filters" sample.

## Assumptions

- The label is an HTML comment (invisible when rendered on GitHub), for example `<!-- spec-sample: status --color=never --root testdata/example-roadmap -->`; the exact form is fixed in the plan.
- The `--color=never` flag stays in the label because the README tells readers to run exactly that command; the test also forces color off.
- Only human-text samples exist; if a JSON sample is added later, the same label and comparison apply.
- Out of scope: checking prose claims, `docs/agent-prompt.md` (already compared by CI), and other documents.
