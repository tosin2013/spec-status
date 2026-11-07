# Feature specification — deterministic roadmap status

**Draft feature SS-001, proposed specification version 0.1.** This file is the specification, not a release. A development CLI now implements a read-only slice; that code does not ratify the open choices below. It preserves the PRD baseline. Every exact grammar, field/code choice or policy labelled PROFILE-R1 is pending [DC-01–08](decisions.md); drafting does not ratify it.

## 1. Product outcome and users

A developer mid-session gets a glanceable roadmap position without opening files. A reviewer gets structural blockers and explicit broken references. A planner gets recorded progress, not effort or dates. A calling LLM receives the same facts through JSON. CI can fail on a configured severity without misclassifying ordinary unfinished work as a defect. A future binary user should not need Go, root or the roadmap extension.

The value is deterministic roadmap reporting, reference integrity and declared-state inconsistencies. The tool does not independently reconcile actual merges, test results or semantic implementation evidence. No uniqueness, proven market or shipping-date claim is part of acceptance.

## 2. Scope and normative language

- **Retained requirement:** source PRD requirement carried forward.
- **PROFILE-R1:** recommended, internally coherent draft behavior, not approved.
- **Future/deferred:** no v1 implementation obligation.

In v1 the runtime reads supported files only. No AI/network/accounts/credentials/subprocess/Git/tests/builds/extensions/writes, no cache or second state store. Development tools may test/build the product separately. Parsing and root-confined reads stay in the Go standard library. The owner-approved exceptions are Cobra for command, help, and flag handling, and Lip Gloss for human text layout and restrained color. They are not a network or telemetry stack, not Bubble Tea or another application framework, and not an independent security review. Static build target `CGO_ENABLED=0`, with disclosed exceptions rather than blanket portability claims.

Inferred quarters are omitted entirely. No quarter calculations, percentages, denominators, flags or JSON inference fields. Explicit supported timing is recorded source text, not date interpretation. No v1 dashboard, trend store, notifications, digest, multi-repo aggregation or package-manager channel. Optional future extension write-back is specified separately in [phase2.md](phase2.md).

## 3. Functional requirements

### FS-01 — Resolve a supported, bounded source

Read `.specify/memory/roadmap.md` by default; permit explicit `--root` and `--roadmap`. Root discovery is nearest ancestor containing `.specify`, unless overridden. Exact path-base/precedence behavior is PROFILE-R1 in the format/CLI contracts. An inaccessible required root/roadmap or fatal parse uses exit 2; do not call it an empty roadmap.

Support only the pinned heading/bullet dialect and explicitly linked spec directories. Inspect `spec.md`, `plan.md`, `tasks.md` as applicable; no recursive crawl. Only explicit supported reference fields and local Markdown task links are path candidates. Prose, inline code, URLs and anchors are not local paths to open. Fences are ignored for item/task recognition. Unsupported item-like content must diagnose rather than yield guessed totals.

### FS-02 — Preserve recorded items and source locations

Retain item document order, ID/title/raw lifecycle status, description/outcome/scope, direct dependencies and explicit spec link. Recognize `undecided`, `needs-info`, `planned`, `specced`, `in-progress`, `implemented`, `verified`, `deferred`, `abandoned`. Unknown labels remain exact raw strings, active/not verified, with warning; never map them to known tokens. A status label does not prove implementation/merge/verification.

Every item/reference/finding has source path/line where obtainable. Duplicate IDs are errors. PROFILE-R1 preserves duplicate occurrences in inventory but prohibits unique-ID graph resolution through ambiguity. Do not delete a row to make counts look consistent. Malformed item headers that make inventory unreliable fail parsing rather than disappearing from totals.

### FS-03 — Summary and progress with explicit denominators

Report total occurrence inventory, raw per-status counts, active count, recorded-verified count/percentage, blocked count and actionable diagnostics. Valid superseded rows remain historical inventory with replacement reference and separate display.

Active excludes deferred, abandoned and superseded rows; unknowns stay active/not verified. Recorded-verified percentage is active recorded-verified rows divided by active rows. Empty denominator is null/N/A, never 100%. Invalid replacement exclusions depend on DC-02; PROFILE-R1 leaves invalid rows active unless independently deferred/abandoned.

Checkbox percentage is checked recognized tasks divided by recognized tasks outside fences, separate from recorded-verified percentage. Zero tasks is null/N/A. Aggregate is weighted by tasks, never mean item percentages. PROFILE-R1 aggregates each distinct explicitly linked tasks file once and labels that scope, including files associated with inactive rows; per-item counts may reuse a shared file and must not be naively summed. Missing tasks input is visible as coverage, not “all done.” Counts imply neither effort/value/readiness nor completeness of code.

### FS-04 — History-safe supersession

Optional `superseded-by` is a documented Spec Status extension, not upstream-native and not a lifecycle enum. Preserve old row/ID/raw status and direct successor reference; supersession is never completion. Text and JSON distinguish it from completed and abandoned. Validate target existence, ambiguous ID resolution, self-link, cycles and malformed/multiple references.

PROFILE-R1 permits same-roadmap acyclic chains, requires the entire chain valid before an edge excludes a row and keeps every direct edge. Optional terminal-successor display never changes relationships. An invalid supplied reference is preserved and diagnosed, not accepted as a metric shortcut. Exact syntax, chains and invalid-reference effects are approval blockers, not historical claims that the source decided them.

### FS-05 — Dependencies and human review

Default satisfaction is a recorded `verified` dependency. An implemented-or-verified alternative is proposed, not automatically included. Missing IDs, self-dependencies and dependency cycles are errors. Unfinished/deferred/abandoned references block under the selected rule. A recorded verified item with unsatisfied dependencies is a contradiction finding; do not rewrite its status or force its percentage down.

Every dependency on a superseded row requires HUMAN REVIEW even with a verified successor. Keep original dependency and replacement visible. Never redirect, auto-satisfy or infer replacement scope equivalence. PROFILE-R1 warning severity, separate review count and once-per-dependent active blocked count await DC-03 approval. Per-item reasons are shown even for inactive rows; active summary scope is labelled. Cycle membership and direct blockers are reported, not dates or transitive delivery predictions.

### FS-06 — Structural findings, not semantic mismatch claims

Report duplicates, unknown statuses, missing linked directories, missing stage-appropriate files/spec links, invalid dependency/replacement edges and broken explicit local task references. Planned missing spec link is informational; specced-or-later missing spec reference is a finding. Exact stage/file requirements are an open approved-matrix gate; PROFILE-R1 candidate is not a source-proven lifecycle law.

Checked task link to nonexistent file: error. Unchecked task link to nonexistent file: info, because it may name future output. Inline-code filenames and prose are not links. Never infer that code fails a spec merely because a structural artifact is missing. Findings do not grant edit permission.

### FS-07 — Output and CI compatibility

Use one underlying result model for text and versioned JSON. Text highlights totals, progress bases, blockers, supersession and path/line diagnostics. JSON includes raw statuses, relationships, tasks, coverage, stable codes and errors with deterministic ordering. No special AI mode; a consumer's request is not action authority.

Default failure threshold error; selectable info/warning/error/none. Exit 0 below threshold; exit 1 at/above; exit 2 invocation/config/required-input/fatal parse. None never suppresses fatal. Ordinary unfinished work is not a diagnostic; blocked unfinished is warning. JSON produces exactly one document including fatal failure, logs excluded from stdout. No volatile default timestamps. PROFILE-R1 clock-independent bytes, schema and rounding are draft contract choices.

### FS-08 — Safe read boundary

Constrain all reads and references to selected root, reject traversal/symlink escapes and exclude sensitive credential/secret/private-key/answer-key files. Narrow supported inputs only. Bound file size/count, total bytes, paths, nesting and graph/task work; explicit limit failures, not silent truncation. PROFILE-R1 numbers and conservative symlink/mount matcher policy are unapproved candidates. A root API alone is neither secret classification nor complete mount isolation. No unsafe fallback when platform capability differs.

Do not open excluded file contents or emit them through diagnostics. Root-required failures are fatal; optional skipped evidence must be visible. PROFILE-R1 cannot promise an atomic multi-file snapshot; observable changes invalidate report, and residual limitations must be documented/tested.

### FS-09 — Release and user delivery requirements

CLI SemVer `vX.Y.Z` differs from input v0.1.1 and JSON version. Planned version output identifies CLI version, tagged commit, Go version and target, with dev builds labelled. Owner-created tag on approved commit eventually triggers exact-source, pinned-reviewed-tooling tests/build/package checks. Required failures block publication. Release assets are tagged-commit-matching archives plus SHA-256, not temporary workflow outputs only. No-Go user installs/upgrades rootlessly. No tag, workflow or binary is created by this package.

Owner-approved RC precedes first stable tag. Platform matrix remains undecided; cross-compilation is not runtime proof. Only publisher job gets narrowly scoped write permission; immutable asset policy, fixes through new versions. No identity/provenance/reproducibility claims unless implemented and verified; unsigned checksums authenticate neither publisher nor its manifest. See [build-release.md](build-release.md).

## 4. Scenario acceptance outline

| Scenario | Expected retained outcome; exact candidate details in acceptance contract |
|---|---|
| Extension absent | Standalone supported-file report, no agent/extension invocation |
| Empty supported roadmap | Zero inventory; active/task percentages null/N/A |
| Unknown status | Raw value retained, warning, active/not verified |
| Checked task references future file | Missing path error; unchecked counterpart info |
| Old verified item superseded | Old verified label retained historically, not active completion credit |
| Dependency on old item, replacement verified | Human review, original edge retained, no redirect |
| Invalid/missing/self/cyclic replacement | Finding; PROFILE-R1 no metric exclusion from invalid field |
| Table or malformed item-like layout | Visible parse/unsupported finding, no fabricated inventory |
| Escaping/sensitive/oversized input | Safe explicit failure/skip according to approved policy, no secret read or authoritative truncated totals |
| JSON fatal invocation/input | One parseable JSON envelope and exit 2 |
| Threshold none with structural findings | Findings still present, exit 0 if no fatal failure |
| No product exists | README states planning stage, no installation/release/test claims |

## 5. Completion, dependencies and open issues

Functional acceptance requires the detailed cases in [acceptance.md](acceptance.md); coverage mapping is [traceability.md](traceability.md). No test has been executed or implementation accepted. The <100 ms median cold-start target applies to a documented 100-item/10,000-task fixture on disclosed hardware/OS and is not a measured fact.

Before implementation: resolve DC-01–06 and affected PROFILE-R1 details, select approved real fixture inputs and separately authorize implementation. Release-specific decisions DC-07–08 must be resolved before release/support claims; license selection remains open; the repository is established separately and no repository operation is performed here. All pending choices and alternative effects must be recorded rather than turning this draft into implicit approval.
