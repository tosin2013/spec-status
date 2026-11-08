# Candidate input dialect and optional supersession extension

**Draft PROFILE-R1, pending DC-01–03 approval.** This is a proposed compatibility profile based on the supported `speckit-roadmap` v0.1.1 heading/bullet shape, not a verbatim upstream grammar or certified parser compatibility. Actual approved fixtures must validate exact field spelling/punctuation before implementation.

## 1. Source surface and path bases

Default selected root: nearest ancestor from invocation directory containing `.specify`; explicit root wins. The root is the project directory containing `.specify`, not the `.specify` directory itself. Default roadmap: selected root plus `.specify/memory/roadmap.md`. PROFILE-R1 explicit relative root is invocation-relative; explicit relative roadmap is selected-root-relative. Absolute overrides must remain confined to the selected root; `.specify` must be a directory. No fallback to another root after a required input fails.

`Spec dir` is root-relative. Only explicitly linked directories are inspected for spec.md/plan.md/tasks.md; no crawl. Explicit local Markdown links in tasks.md resolve relative to that tasks.md directory. Returned source paths are normalized selected-root-relative paths with forward slashes. Lexical resolution is not sufficient security: the approved safe-open/symlink policy applies before filesystem access. Traversal attempting to leave root is rejected, even if a later component would come back in; exact normalization cases need fixture approval.

## 2. Candidate recognition rules

- UTF-8 input only; optional leading UTF-8 BOM removed for recognition, CRLF and LF treated equivalently. Lines are 1-based physical lines after newline normalization, no leading header deletion. Invalid encoding is fatal.
- A roadmap row starts at a level-three heading `### NNN — Title [status: token]`. `NNN` is exactly three ASCII digits. At least one nonblank title character; exact one-line heading shape and spaces must be confirmed against fixtures.
- Candidate field bullets directly under a row: `- Description: ...`, `- Outcome: ...`, `- Scope: ...`, optional `- Depends on: 001, 002`, optional `- Spec dir: specs/003-example`, optional `- superseded-by: 004`. Continuation text may be retained as raw text, never executable syntax. PROFILE-R1 accepts one backtick pair around an explicitly declared local path or bare path text; no whitespace guessing/URL decoding.
- Fields occur before the next row heading. Top-level metadata/prose and non-item headings are allowed but not parsed as items. Fenced code, including matching backtick/tilde fence families, is excluded from field/item/task recognition; fence edge cases require approved fixtures.
- **Visible failure for item-like input (review fix, 7 November 2026).** Each of these is fatal rather than skipped: a heading at any level that carries `[status:`, or starts with exactly three digits (any level-three heading starting with a digit, as before), or is indented; a line outside code that carries `[status:` and is not a supported item heading (inline code such as `` `[status: x]` `` is prose); a pipe-table row whose first cell is three digits or that carries `[status:` (SS_LAYOUT_UNSUPPORTED); a list item whose label is a recognized field under another case, spacing, hyphenation, marker, or indentation, such as `- Superseded-by:` or `* Depends on:` (SS_ITEM_MALFORMED); a recognized field outside any item; an unclosed code fence. Any other heading ends the current item, so later fields cannot attach to the item above. Fences follow CommonMark openers: at most three spaces of indentation, and a backtick fence's info string has no backtick, so ```` ```x``` ```` is inline code.
- Nested non-reference bullets may be retained as descriptive text. They do not create hidden dependency/spec-dir fields. Duplicate recognized fields, malformed ID lists or unsupported item-like table/header structures are diagnosed; if inventory cannot be trusted, PROFILE-R1 fatal parse suppresses authoritative totals.
- **Owner-approved on 7 November 2026 (speckit-roadmap spelling).** The speckit-roadmap extension writes field labels in bold. `- **Label:** value` and `- **Label**: value` are read as `- Label: value` for every field above. `Scope (in)` fills `Scope`; giving both `Scope` and `Scope (in)` in one row is a repeated field. Other extension fields (`Scope (out)`, `Governed by`, `Addresses`, `Notes`) are descriptive text and create no references. `Depends on` accepts `none` (any case) as an empty list, one trailing period, and one trailing parenthetical note per ID, so `002 (verified), 003 (verified).` reads as 002 and 003. Any element that is still not three digits, such as `003-brief`, makes the whole list malformed (SS_DEPENDENCY_MALFORMED); nothing is silently dropped. `Spec dir` and `superseded-by` accept one trailing period; `superseded-by` takes no note. These rules are exercised against the extension's own v0.1.0 roadmap in `testdata/speckit-roadmap-upstream`, which is a fixture, not a certification of upstream compatibility.
- Known statuses are exact lowercase tokens: undecided, needs-info, planned, specced, in-progress, implemented, verified, deferred, abandoned. Unknown raw tokens retained without case folding/remapping; warning. Empty/malformed status is not an alias for planned.
- IDs are unique within the selected roadmap. Spec/task files do not create additional roadmap item IDs. Duplicate roadmap occurrences remain visible inventory plus error; references to duplicated IDs are ambiguous, not first-match wins.
- An advisory format declaration may be documented, but no unapproved mandatory marker or migration is required of existing v0.1.1-shaped files. Marker exact syntax/version rules remain DC-01.

Illustrative **candidate**, not a statement that upstream uses these exact strings:

```markdown
### 001 — Preserve original feature [status: implemented]
- Description: Historical original feature.
- Outcome: Original outcome remains recorded.
- Scope: Original scope is retained.
- Spec dir: specs/001-original
- superseded-by: 002

### 002 — Replacement feature [status: verified]
- Description: Replacement scope can differ.
- Outcome: Recorded as verified, not independently checked by this CLI.
- Scope: Replacement scope.
- Spec dir: specs/002-replacement

### 003 — Follow-on work [status: planned]
- Description: Still refers to the original requirement.
- Outcome: Requires a scope decision.
- Scope: Follow-on.
- Depends on: 001
```

Item 001 stays in inventory with raw implemented status and explicit successor. Item 003 must require HUMAN REVIEW even though 002 is verified. The example demonstrates relationships, not that a binary can currently parse it.

## 3. Explicit tasks and local file references

PROFILE-R1 recognized task: a bullet with `[ ]`, `[x]` or `[X]` checkbox followed by text, optionally nested within documented depth; not inside fences. Only tasks.md is counted for task progress; spec.md/plan.md contents are not used for semantic checks or checkbox percentages. No inference from words such as “done.” Empty task set is null completion.
Ordered-list checkboxes (`1. [ ] text`, `2) [x] text`) are recognized tasks. Task-file fences use the same CommonMark opener rule, and an unclosed fence in a tasks file is fatal. Checkbox totals are capped across all task files read, not per file (review fix, 7 November 2026).

An explicit ordinary Markdown link with relative nonempty path is a task-file reference. Anchors-only links, URLs/schemes and ordinary inline-code filenames are not local file candidates. Candidate parser supports simple inline links; reference-style definitions, images, escaped/nested link syntax and query/fragment combinations are unsupported until examples are approved. Never reinterpret an unsupported URL as a local path. Exact grammar is an implementation blocker, not permission to silently ignore item-like data. Declared reference fields are supported only where documented; do not invent a prose extractor.
An inline link may wrap its path in `<...>` and may carry a quoted title (`[t](file.md "Title")`); only the path is checked. Reference-style links (`[label][ref]`) and autolinks (`<path>`) are reported as SS_REFERENCE_UNSUPPORTED, never ignored (review fix, 7 November 2026).

Check reference existence safely, not contents of every referenced source file. A safe missing file from checked task is error; unchecked task is info. Sensitive/escaping references are excluded before content reads; no network lookup or path guessing. Each recognized task counts once regardless of number of links. Aggregates count distinct linked tasks files once, with per-item shared-file attribution labelled.

## 4. Optional `superseded-by` relationship

The extension is optional input data; absence means no supplied replacement. It is **not upstream-native v0.1.1**, not a new lifecycle token, and not an instruction to edit files. It contains one same-roadmap three-digit item ID, not a URL, path, ID list or free-form replacement description. Original ID/title/status/text and source location remain preserved.

Validation: supplied value well-formed; target exists uniquely; no self-edge; no directed replacement cycle. PROFILE-R1 permits forward/backward references and acyclic chains. It validates all edges to chain end before marking an origin valid; links into malformed/ambiguous/cyclic chains remain invalid. Preserve supplied direct reference on errors. Terminal successor, if shown, is derived display only and never a graph rewrite.

Valid superseded rows remain in inventory/raw-status/historical counts, shown separately, excluded from active numerator and denominator. Invalid supplied links are errors and, in PROFILE-R1, do not exclude a row solely by reference presence. Invalid labels/links must not allow artificial completion gains. Deferred/abandoned exclusions still apply independently.

Every original dependency edge pointing at a superseded item requires human review. No automatic redirect to chain end or direct successor, even when verified. This applies regardless of chain display or optional future provenance. Exact diagnostic severity/blocked-count policy is DC-03, not implied by graph mechanics.

Reason/evidence/decision-date/decision-maker fields may be useful later but are not silently required by this one-field extension. No external evidence fetching. Extension authoring is future approved write-back, never read-only parser output.

## 5. Lifecycle artifact expectations (approved 7 November 2026, D-2026-11-07-03)

| Recorded status | Spec reference / spec.md | plan.md / tasks.md |
|---|---|---|
| undecided, needs-info | Not required | Not required |
| planned | Missing link info; if explicit link exists inspect safely | Not required; missing optional files not error |
| specced | Required, missing link/file error | Not required |
| in-progress, implemented, verified | Required, missing link/file error | Required, missing error |
| deferred, abandoned, unknown | No invented stage requirement; explicit links still checked | No inferred requirement |

The owner approved this matrix on 7 November 2026 and rejected the minimal spec-only alternative (D-2026-11-07-03). Changing it later changes fixtures and diagnostics explicitly, not silently. Supersession retains recorded lifecycle; PROFILE-R1 stage findings remain visible on historical rows, without affecting active completion eligibility. A missing link or file proves only a mechanical absence, not that implementation is wrong.

## 6. Safety and compatibility limits

All accesses use the selected root and approved exclusion/regular-file/symlink policy. No recursive lookup, glob expansion, shell, Git, extension or credential access. Candidate budgets/matcher in DECISIONS require approval and fixture calibration. An accepted path string is not permission to open a secret. Quarters are not inferred; explicit unknown timing fields may be retained as raw source text but not normalized into a new date schema without approval.

Before claiming compatibility, record upstream template/version, accepted examples, unsupported constructs, line mapping and malformed behavior in a versioned fixture corpus. Original baseline prose is partly condensed, so this proposed profile must not be advertised as exact recovered upstream syntax. Format, CLI and JSON versioning are independent.
