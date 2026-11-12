# Clarifications and decision log

**Draft.** The owner has authorized drafting and public-project documentation, not implementation or release. This document distinguishes requirements from proposed choices. PROFILE-R1 is one consistent recommended candidate; **none of its unresolved choices is ratified by drafting it**.

## Settled boundaries

- Omit inferred quarters entirely from v1. Preserve recorded timing, without inventing dates or treating question IDs as calendar data.
- Retain strict v1 history-preserving supersession and human review of its dependents; no automatic dependency redirection.
- Offline read-only core, `CGO_ENABLED=0`; no runtime AI/network/subprocess/Git/credentials/writes/extensions. Parsing and confined reads stay in the standard library. On 7 November 2026 the owner approved Cobra for command, help, and flag handling only. The owner later approved Lip Gloss for human text layout and restrained color only, then a compact static summary panel, then human-text filters and `--explain`. Filters are views: they do not change summary denominators, diagnostics, exits, or the schema 1.0 JSON payload. Those exceptions do not approve Bubble Tea, dashboards, telemetry, inferred quarters, or a security review of the modules. DC-01–DC-08 remain unapproved. The command process unsets `BASH_COMP_DEBUG_FILE` for the call so Cobra cannot append that log, including on hidden `__complete` errors, then restores the previous value. Concurrent calls in one process are serialized for that reason. The product binary makes one call. Generated completion scripts remain shell text and are not executed by spec.
- Preserve raw statuses and unknowns, recorded-verified versus checkbox progress, JSON/text parity, default threshold error and exits 0/1/2.
- Future extension write-back is optional, separate and per-mutation approved; optional integrations remain deferred.
- Public project documentation targets [tosin2013/spec-status](https://github.com/tosin2013/spec-status). Repository establishment is separate from this package; no license, code implementation or release is created here.

## Proposed profile: approval register

| ID | Needs choosing / alternatives | PROFILE-R1 recommendation, not approved | Consequence / implementation gate |
|---|---|---|---|
| DC-01 | Exact input headings/fields, whitespace, nested bullets, version marker; alias/tolerant parser versus narrow dialect; root/override bases | Narrow v0.1.1-shaped headings and documented one-field extension; exact examples in roadmap-format.md. BOM/CRLF normalization, no aliases/tables. Optional advisory compatibility marker, no migration required. Explicit root then nearest ancestor discovery; roadmap override root-relative, spec dir root-relative, task links tasks-file-relative. | Approve real fixture examples, uniqueness and source-base semantics before parser implementation; malformed item-like content invalidates reliable totals. Owner approved on 7 November 2026 one refinement only: speckit-roadmap bold field labels, the `Scope (in)` alias, and annotated `Depends on` lists, as written in roadmap-format.md §2. The rest of DC-01 remains open. **Approved 11 November 2026 (D-2026-11-11-02).** |
| DC-02 | Metric effect of invalid replacements; legal chains; mandatory provenance versus optional/future fields | Only validated references exclude active rows; invalid supplied reference retained, diagnosed and not excluded solely by its presence. Allow acyclic same-roadmap chains, preserve each direct edge and optionally show terminal ID; no redirect. Do not require new reason/date/decision fields. | Approve graph-validity propagation and ambiguity cases; supersession itself remains in v1. Optional provenance extension needs a later format decision. **Approved 11 November 2026 (D-2026-11-11-02).** |
| DC-03 | Human-review severity and blocked counting; dependency satisfaction alternatives; lifecycle-file matrix and stable codes | Warning HUMAN_REVIEW, distinct review count and once-per-dependent blocked count; verified-only default, alternate implemented-or-verified policy remains uncommitted. Candidate stage matrix in roadmap-format.md; candidate codes in diagnostics.md. | Approve stage matrix, contradiction severity, blocked scope and code registry. No auto-satisfaction through successor. **Approved 7 November 2026 (D-2026-11-07-03, see Recorded decisions).** |
| DC-04 | Name/config/schema/rounding/task scope; retained --as-of use; new aliases versus existing interface | Working `spec status`; no config file initially, explicit flags then defaults. No --as-of without a concrete non-quarter use; no clock dependence. Candidate JSON schema 1.0, two-decimal half-up percentages, stable snake_case. Unique linked task-file aggregation with separately labelled per-item reuse. | Command name, schema/envelope/canonical bytes and task/reference grammar require confirmation. Retained --json/exit meanings are not reopened. If date feature is chosen, revise all related contracts before implementation. **Approved 11 November 2026 (D-2026-11-11-02).** |
| DC-05 | Contained symlinks versus reject all linked symlinks; mount and mid-read behavior; sensitive exclusions and patched Go baseline | Regular-file-only reads, reject symlink components in referenced paths initially; deny cross-mount referenced reads where platform can prove policy. Fail required-input/source-changed report rather than claim atomicity. Concrete deny examples below; supported patched toolchain selected only after current evidence review. | A capability gap must fail closed, not downgrade confinement. Root canonicalization, case behavior, matcher, race safety and platform-equivalent enforcement require approval. No API/toolchain version selected by this draft. **Approved 11 November 2026 (D-2026-11-11-02).** |
| DC-06 | Fixed versus configurable limits; fatal versus partial totals | Candidate caps: 2 MiB/file; 32 MiB total actual bytes read; 1,024 distinct opened input files; 4,096-byte normalized path representation; 32 parser nesting levels. Fixed initial profile; explicit limit error, no authoritative partial totals. | New proposed numbers, uncalibrated. Approve/correct against fixtures, including 100 items/10,000 tasks. Include lookup/edge/time-work budgets in final review; byte caps alone are not sufficient. **Approved 11 November 2026 (D-2026-11-11-02).** |
| DC-07 | Initial two versus five/six targets; checksum-only versus additional signed identity gate; OS minimums/rollback | Candidate linux/amd64 + darwin/arm64 only after real runners prove them; other targets not advertised supported. Direct GitHub Release archives, exact commit, pinned tooling, gated checks and owner RC approval retained. Truthful checksum-only integrity initially; attestation/signing a separate owner choice. | Actual runners, patched toolchain, support window, OS versions, packaging rules, permission review and identity policy require approval before release work. No binary support claimed today. **Approved 11 November 2026 (D-2026-11-11-02).** |
| DC-08 | Dogfood targets; public artifact versus feature expansion; license/security/support/fork expectations | Personal utility remains baseline; public-documentation intent accepted, not a market commitment. Proposed approved Spec Kit roadmap plus one owner-named portfolio roadmap, hand-count parity. Public repository is tosin2013/spec-status; repository operations remain outside this package. | Name approved inputs. Repository home is established; license, security disclosure route, support window and mirror/Bitbucket expectations remain unresolved. No license file invented. **Dogfood inputs approved 7 November 2026 (D-2026-11-07-08).** License, security route, support window and mirror expectations remain open. **Approved 11 November 2026 (D-2026-11-11-02); hosting: GitHub only, no mirror (D-2026-11-11-03).** |

## PROFILE-R1 refinements requiring explicit approval

These make contracts draftable; they are new technical recommendations, not hidden owner decisions.

- **D-01A Format:** ID is exactly three ASCII digits, unique within one roadmap. Heading punctuation is an em dash, fields at one bullet level; backtick-wrapped or bare local paths accepted only as declared in roadmap-format.md. No alias status normalization. Mixed invalid UTF-8 is fatal; optional UTF-8 BOM and CRLF accepted.
- **D-02A Graph validity:** a replacement chain is valid only if all nodes/edges to its end resolve uniquely and are acyclic. An edge into an invalid chain does not exclude its source. Duplicate ID rows are retained as occurrences; ambiguous references fail validation. No automatic terminal-successor inference of dependency satisfaction.
- **D-03A Blocker scope:** show per-item direct blockers for all rows; summary blocked/review counts cover active rows only. A row with multiple reasons counts once. A verified-but-blocked row still contributes its recorded-verified numerator and emits contradiction; do not silently recategorize its label.
- **D-03B Stage candidate:** planned and earlier do not require spec links; specced requires link/spec.md; in-progress/implemented/verified additionally require plan.md/tasks.md in candidate matrix. Missing linked artifacts do not establish code absence. This is a proposed stricter mapping, not a source-settled stage requirement; choosing a minimal mapping is a valid alternative that needs explicit acceptance changes.
- **D-04A Error model:** fatal input gives `complete:false`, null summary, no authoritative items, sorted safe diagnostics and structured fatal `error`. Normal result has `complete:true`; skipped optional evidence is visible through coverage/findings, not a claim of complete repository inspection. No OS error strings or secret input snippets on stdout.
- **D-05A Sensitive path candidate:** deny any path component `.git`, `.ssh`, `.aws`, `.gnupg`; deny `.env` / `.env.*`, case-insensitive secret/credentials/private-key/answer-key basename patterns, and `.pem`, `.key`, `.p12`, `.pfx` suffixes. This over-inclusive starting matcher is incomplete and must be reviewed; it is not a security certification or permission to open unmatched secrets. No content sniffing of excluded files. Approve normalization/case policy per supported platform.
- **D-05B Root and races:** explicit root must exist; recommend canonicalize selected root once, then enforce all later paths beneath it through a race-resistant safe-open boundary. No promise of an atomic multi-file snapshot. If observable source identity/size/mtime changes during a read, invalidate result. Same-size undetectable edits are not claimed impossible; acceptance requires documented limits or stronger snapshot design, not a false guarantee.
- **D-06A Work budget:** 1,024 decoded roadmap items, 16,384 total dependency/replacement/reference edges and 100,000 recognized task checkboxes are additional candidate work caps for approval. Algorithmic cycle checks must avoid unbounded recursion. Final path resolution and read-call budgets must be measurable. Configurable overrides, if approved later, must never disable root confinement or sensitive exclusions.

## Recorded decisions

| ID | Approver, date | Chosen | Rejected | Impacted requirements and tests |
|---|---|---|---|---|
| D-2026-11-07-01 | Owner, 7 November 2026 | DC-01 refinement only: speckit-roadmap bold field labels, `Scope (in)`, annotated `Depends on` lists (roadmap-format.md §2) | Strict bare lists only | FS-01; AC-04; `internal/status/speckit_test.go`; fixture `testdata/speckit-roadmap-upstream` |
| D-2026-11-07-03 | Owner, 7 November 2026 | DC-03 as built: satisfaction is recorded `verified` only; D-03B strict stage matrix; human review of a superseded dependency is a warning; per-item blocker reasons on every row; active-only blocked and review counts, once per row; verified-but-blocked keeps its numerator and emits a contradiction | Implemented-or-verified satisfaction; minimal spec-only matrix; human review as error | FS-05, FS-06; AC-06, AC-08, AC-11, AC-12; no code change; existing tests stand |
| D-2026-11-07-08 | Owner, 7 November 2026 | DC-08 dogfood inputs: this repository's `.specify/memory/roadmap.md` and the speckit-roadmap v0.1.0 roadmap in `testdata/speckit-roadmap-upstream` | A private portfolio roadmap (may be added later) | AC-24; roadmap item 003 |
| D-2026-11-07-04 | Owner, 7 November 2026 | After the second acceptance review: (a) the prefix rule for field labels in roadmap-format.md §2; (b) item-like headings narrowed to `[status:` or three digits, a dash, and a space; (c) a dependency on any item with a supplied `superseded-by` line, valid or not, requires human review, while an invalid chain still leaves the item active (D-02A unchanged for metrics) | Closed item body (every bullet must carry a known label); keeping the stricter numeric-heading rule; judging dependents of an invalid chain on the target's own status | FS-01, FS-04, FS-05; AC-04, AC-10, AC-11; `internal/status/review2_test.go` |
| D-2026-11-07-05 | Owner, 7 November 2026 | After the third acceptance review: a closed item body (roadmap-format.md §2), Unicode lookalike normalization for recognition, and task-file fixes for checkbox spelling, indented fences, unparsed links, and file identity | Another targeted round of variant fixes; adopting a CommonMark library (would need a constitution amendment to Art. II) | FS-01, FS-03; AC-03, AC-04, AC-05, AC-07; `internal/status/review2_test.go`, `review3_test.go` |
| D-2026-11-07-06 | Owner, 7 November 2026 | Acceptance bar for 001: `verified` may be recorded when an independent review finds no blocking defect (a wrong count, status, dependency state, or exit code without a diagnostic) inside the documented dialect, input outside it fails visibly, and the remaining limits are stated in README.md "Limits" | Requiring that no reviewer can find any defect of any kind | Roadmap item 001; README.md Limits **Superseded on 7 November 2026 by D-2026-11-07-09**; kept as the bar reviews 1 to 7 were judged against. |
| D-2026-11-07-09 | Owner, 7 November 2026 | Acceptance bar from feature 002. A blocking defect is (a) a wrong item count, task count, status, dependency state, or exit code, without a diagnostic, on input that matches the shape inventory of a current upstream pin (`specs/002-upstream-format-conformance/contracts/shape-inventory.md`), or (b) a silently misread `Depends on`, `Spec dir`, or `superseded-by` value in any layout. A wrong result on other exotic input fails visibly or is listed in README Limits, and is then non-blocking. Existing visible failures stay (FR-014). | Keeping D-2026-11-07-06, which treated every divergence from CommonMark renderers as blocking and proved unbounded across seven reviews | Roadmap item 001; feature 002 FR-012 to FR-015; `specs/002-upstream-format-conformance/reclassification.md` |
| D-2026-11-07-07 | Owner, 7 November 2026 | After the fifth acceptance review: a narrow tasks dialect (roadmap-format.md §3) that measures indentation from the enclosing list item and fails visibly on block-quoted structure, checkboxes four or more columns past their item, empty comments, and raw HTML blocks | Modeling CommonMark containers in full; stating the miscounts as README limits | FS-03; AC-03, AC-05; `internal/status/review5_test.go` |

## Readiness ledger

| Gate | Current state |
|---|---|
| Documentation drafting | Authorized; this package is draft work |
| Required review returns | Received; semantic/compatibility decisions remain, not missing replies |
| Source certainty | Requirements coverage assembled; original input-format prose partly condensed, not fully verbatim recovered |
| PROFILE-R1 adoption | DC-03, one DC-01 refinement and the DC-08 dogfood inputs approved on 7 November 2026; the rest of DC-01, DC-02, DC-04–DC-07 and DC-08's license, security route and support window approved on 11 November 2026 (D-2026-11-11-02); DC-08 hosting approved the same day (D-2026-11-11-03). Complete. |
| Implementation start | Owner authorized implementation on 7 November 2026. That authorization does not approve DC-01–08. |
| Public repository | tosin2013/spec-status. This decision log does not perform repository operations. |
| Release | Not authorized. In-repo tests have been run by the implementer; they are not independent verification or release evidence. No workflow or published binary. |

## Implementation proposals in the development build

These are reversible defaults chosen so a useful read-only v1 could be built. They are not owner approvals, and repeating them here does not close DC-01–08. Cobra is separate: the owner approved it for command handling, and it is not one of these unresolved proposals.

- Narrow grammar from DC-01 and D-01A: `### NNN — Title [status: token]`, the six field labels in roadmap-format.md, UTF-8 with optional BOM and CRLF, and no table items. `..` in a reference is rejected. Spec directories are root-relative. Task links are relative to the tasks file. Task bullets may use `-`, `*`, or `+`.
- History metrics from DC-02 and D-02A: only a fully valid acyclic chain is removed from the active counts. An invalid `superseded-by` value is kept and diagnosed, and the row stays active. A displayed terminal ID does not rewrite dependencies.
- Dependencies from DC-03 and D-03A: satisfaction means recorded `verified` only. Human review of a superseded dependency is a warning. An active row counts once in the blocked summary. The stage file matrix is the D-03B candidate.
- Interface from DC-04: command `spec status`, no config file, no `--as-of`, candidate JSON schema 1.0, half-up percentages to two decimals, and one count per distinct tasks file.
- Confinement from DC-05 and D-05A/D-05B: regular files only, symlink components rejected, Linux `O_NOFOLLOW` plus `/proc/self/fd` and a device-id check, and the candidate sensitive-name matcher. A detected size or mtime change invalidates the read.
- Budgets from DC-06 and D-06A: 2 MiB per file, 32 MiB total, 1,024 opened files, 4,096-byte paths, 32 nesting levels, 1,024 items, 16,384 edges, and 100,000 checkboxes.

Not chosen, and not implemented: the implemented-or-verified dependency alternative, a date flag, signing, a support matrix, a license, and phase-two writes.

## Change procedure

Record a choice with decision ID, approver, date, chosen alternative, rejected alternatives and impacted requirements/tests. Do not replace “proposed” with “approved” merely because it appears in several documents. A scoped decision may clear one gate without clearing implementation/release. Revisions of input/schema compatibility must identify migration/fixtures. Licensing decisions must be explicit; public visibility does not imply a permissive license.

## Draft proposals recorded, not chosen

On 7 November 2026 the owner authorized drafting, not implementation and not a decision on the open choices. This pull request records two drafts: [timing-evidence.md](timing-evidence.md) and [companion.md](companion.md).

Retained baseline: inferred quarters stay omitted from v1. Four ordinal buckets cannot answer a delivery date. Approved exceptions that stay approved are Cobra for commands, help, and flags; Lip Gloss for human text; the static summary panel; and human-text filters with `--explain`. Those approvals do not make the rest of the roadmap profile approved. DC-01–DC-08 stay open in this document. A separate open pull request may record other owner decisions; this section does not adopt or revert that pull request.

Assistant drafting choices, not owner decisions: where to file the drafts, the tentative command spelling `spec companion --prompt-file`, the 64 KiB sketch, and the citable field list. Syntax, date severity, input-extension version, invalid-supersession denominator policy, envelope version, loader placement, and consent rules remain open.

A clarification in the same documentation authorization says a harness means an external tool, with Claude Code, Grok Bot, Muse, and OpenCode as illustrative consumers only. That is not verified compatibility and not permission to embed a model in the core.

## Recorded decision — two modes, 8 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-08-01 | Owner, 8 November 2026 | Product direction in [two-mode-direction.md](two-mode-direction.md): versioned templates only; separate roadmap and inventory modes; missing, empty, unrecognized, and malformed stay distinct; no silent inventory fallback; checklists excluded; extension recommended but not required | Command syntax, default mode, version identification, folder and checkbox grammar, inventory schema and exits, v1 placement. Those remain open |

## Recorded decision — line-level label rule, 8 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-08-02 | Owner, 8 November 2026 (approved by "proceed" on the proposed rule) | Any visible roadmap line containing a `Depends on`, `Spec dir`, or `superseded-by` label with a colon must be an accepted field line under an item, otherwise fatal. Rule (b) of D-2026-11-07-09 then holds by construction rather than by listing layouts. Prose containing those labels with a colon is rejected and listed in README Limits | Listing further layouts one review at a time; marking item 001 verified without a review |

## Recorded decision — markup-encoded labels, 8 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-08-03 | Owner, 8 November 2026, after the T026 review | Amends rule (b) of D-2026-11-07-09: "any layout" means any block layout. A `Depends on`, `Spec dir`, or `superseded-by` label or item heading encoded with markup (HTML entities or character references, backslash escapes, a link around the label, inline HTML inside it, invisible characters outside the normalized set, or a label split across lines) is outside rule (b) and is listed in README Limits. A list or quote nested under a metric field is fixed to fail visibly (review B1) | Decoding entities, escapes, and inline HTML before matching (rendering Markdown); recording 001 verified without a further review |

## Recorded decision — closed field body, 8 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-08-04 | Owner, 8 November 2026, after the final T026 review | Under a `Depends on`, `Spec dir`, or `superseded-by` field, an indented line either joins the value as paragraph text or is fatal, so no block nested under a metric field can be silently ignored. A code block nested under a metric field is fatal | Listing nested block types one review at a time; recording 001 verified without a further review |

## Recorded decision — recording 001 verified, 8 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-08-05 | Owner, 8 November 2026, after the third T026 review round | Amends SC-006 of feature 002. Roadmap items 001 and 007 are recorded `verified` on the evidence of three independent T026 review rounds (at ac62214, dadd8fd, and 4574c13) and the fixes for every blocking finding they reported, although no single round returned zero blocking defects. Rule (a) was clean in all three rounds. The residual risk, that further contrived rule (b) layouts exist, is stated in README Limits. The reviews were separate AI agent sessions, not human reviewers, and none checked GitHub's own renderer | Running a fourth review before recording verified; leaving 001 at implemented |

## Recorded decision — dogfood complete, 11 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-11-01 | Owner, 11 November 2026 | Roadmap item 003 (TASK-16 dogfood part / AC-24) is done and recorded `verified` on the implementer's hand-count parity with both approved inputs: this repository's roadmap (2026-11-07 and again at 1.3.4 on 2026-11-09) and the speckit-roadmap v0.1.0 roadmap. 003 links `specs/001-spec-status`, where TASK-16 and AC-24 are defined. The owner-named portfolio roadmap is out of 003's scope; it can be added later under DC-08 | Waiting for a portfolio roadmap; an independent recount |

## Recorded decision — contract approval, 11 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-11-02 | Owner, 11 November 2026 | Closes DC-01–DC-08. **DC-01** (rest), **DC-02**, **DC-04**, **DC-05**, **DC-06**: PROFILE-R1 approved as built, including D-01A, D-02A, D-04A, D-05A, D-05B and D-06A: the narrow speckit-roadmap-shaped dialect with fatal failure on unsupported structure; only valid acyclic supersession chains exclude rows, with no redirection and no new provenance fields; `spec status` with no config file and no `--as-of`, JSON schema 1.0, two-decimal half-up percentages, unique linked task-file aggregation; regular-file-only, symlink-component rejection, the D-05A sensitive-name list as a starting point rather than a certification, invalidation on mid-read change; fixed limits of 2 MiB per file, 32 MiB total, 1,024 files, 4,096-byte paths, 1,024 items, 16,384 edges, 100,000 checkboxes and 32 nesting levels, fatal when exceeded. **DC-07**: v1 targets linux/amd64 and darwin/arm64, each advertised only after a real runner proves it; GitHub Release archives with SHA-256 checksums only. **DC-08**: Apache-2.0 license; security reports through GitHub private vulnerability reporting; only the latest release is supported, best effort, no response-time promise | Configurable limits; a config file; `--as-of`; linux/arm64 or Windows as v1 targets; signing or attestations in v1; MIT; an email security contact. Not covered: DC-08's mirror or Bitbucket expectations, which remain open |

## Recorded decision — hosting and community files, 11 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-11-03 | Owner, 11 November 2026 | Closes the last open DC-08 point. GitHub (`tosin2013/spec-status`) is the only repository home and release channel; no mirror or Bitbucket copy. The repository adds `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md` (Contributor Covenant 2.1, enforcement contact tosin.akinosho@gmail.com) | A Bitbucket or other mirror; extra release channels |

## Recorded decision — dependency update policy, 11 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-11-04 | Owner, 11 November 2026 | Dependency proposals (feature 003, roadmap item 009): a proposal that only changes versions of modules already in `go.mod` may merge after the CI check passes and the owner approves the pull request. A proposal that adds any module not already in `go.mod` (direct or indirect) or changes the minimum Go version must not merge until an owner decision naming the change is recorded here. Updates: Dependabot security fixes on advisory plus one weekly grouped minor/patch proposal for Go modules and one for workflow components, majors separate. CI: tests on linux/amd64 and darwin/arm64, go vet, govulncheck, on pull requests and `main`. Dependabot alerts and automatic security fixes enabled 11 November 2026. No constitution amendment: Art. II limits which modules are imported, and new modules stay behind this rule | A decision entry for every bump; security-only or monthly updates; tests-only CI; automatic merging |

## Recorded decision — minimum Go version, 11 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-11-05 | Owner, 11 November 2026, during planning of feature 003 | `go.mod` declares `go 1.26.0`, the oldest supported Go release line (1.27 and 1.26 are supported). CI tests the newest 1.26 patch and the current stable release. Reasons: Go 1.22 no longer receives fixes, govulncheck reports reachable standard-library advisories on old toolchains, and the `golang.org/x/sys` security fix needs Go ≥1.25. Feature 004 still chooses the exact toolchain patch used for release builds. A later change of the minimum needs a recorded decision (D-2026-11-11-04) | Go 1.27 minimum; keeping Go 1.22 with a non-blocking scan |

## Recorded decision — release flow, first version and channels, 11 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-11-06 | Owner, 11 November 2026 (feature 004, roadmap item 004) | **Flow:** the owner pushes a SemVer tag (`vX.Y.Z`, or `vX.Y.Z-rc.N` for candidates). The release workflow runs the full CI, builds linux/amd64 and darwin/arm64 archives, smoke-tests each on its own platform, and creates a **draft** GitHub Release with a SHA-256 checksum file. Only the owner publishes. **First version:** `v0.1.0` (pre-1.0; a minor version may break compatibility, and the notes list it); roadmap 004's "first stable tag" becomes "first tagged release". **Channels:** GitHub Release archives plus documented `go install …@vX.Y.Z`. `spec version` reports the module version for `go install` builds and labels release, module and development builds distinctly. JSON is unchanged | Manual-dispatch releases; publishing directly from the tag; starting at v1.0.0; releases-only (go install reporting 0.0.0-dev); a Homebrew tap |

## Recorded decision — roadmap sync, 12 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-12-01 | Owner, 12 November 2026, after the roadmap sync (`.specify/memory/roadmap-sync-2026-11-12.md`) | (1) Roadmap item 009 (dependency automation, feature 003) is recorded `verified`. The evidence: all CI checks passed on PR #29; three seeded-fault runs each failed exactly the expected jobs; Dependabot's first proposal (PR #31) passed CI and the merge rule with only version changes. Task T026 (offering branch protection) is closed as offered; the setting stays an owner choice. (2) Q5 (further support for annotated `Depends on` lists) is deferred past v0.1.0; the visible `SS_DEPENDENCY_MALFORMED` finding stays. (3) The false `SS_REFERENCE_UNSUPPORTED` warning for link-like text inside inline code in task files (sync finding F5) is fixed before the first release candidate. (4) The repository adds a pull request template, a bug report form, a feature request form, and an issue chooser that routes security reports to SECURITY.md | Keeping 009 `implemented` until branch protection is decided; deciding Q5 now; deferring the F5 fix |

## Recorded decision — v0.1.0 withdrawn, release protections, 12 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-12-02 | Owner, 12 November 2026 | **Facts:** at the owner's request, the Cursor integration merged PR #38 and published `v0.1.0`. The tag is annotated and points at `7112c3e`; the release was authored by `cursor[bot]`. It was built outside `.github/workflows/release.yml` with Go 1.26.0. It had untested archives for five platforms (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64) and an extra `SHA256SUMS`, and none went through the smoke tests. The workflow's own run for the tag stopped at `verify` because the release already existed. **Chosen:** (1) `v0.1.0` is withdrawn: its release was turned back into a draft. The tag stays at `7112c3e`, so `go install …@v0.1.0` checksums never change. The version number is never reused for other bytes. (2) The first release is `v0.1.1`, tagged by the owner and built, smoke-tested and drafted by the workflow, then published by the owner. (3) Repository protections are enabled: a tag ruleset "release tags: owner only" (creating, updating and deleting `v*` tags only by repository admins); branch protection on `main` requiring the four test jobs and `hygiene`, with no force-push or deletion and admins able to bypass; immutable releases. (4) The `v0.1.0-rc.1` draft is deleted and its tag kept | Keeping `v0.1.0` and editing its assets; leaving it published; moving the `v0.1.0` tag to a workflow-built commit |

## Recorded decision — first release is v0.1.2, 12 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-12-03 | Owner, 12 November 2026 | The `v0.1.1` tag (`a2c5dd2`, the merge of PR #39) was built, smoke-tested and drafted by the release workflow (run 43159089291), but it came before PR #40, which added the AI-agent prompt (README, release notes) and `retract v0.1.0` in `go.mod`. The `spec` code is identical. The unpublished `v0.1.1` draft is deleted, and its tag stays and is never reused. The first published release is `v0.1.2`, tagged by the owner on `main` after this record merges | Publishing v0.1.1 with the prompt added to its notes by hand; publishing v0.1.1 as is |

## Recorded decision — agent skills extension, 12 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-12-04 | Owner, 12 November 2026 | Roadmap item 012, targeting v0.2.0. **Packaging:** a Spec Kit extension (id `spec-status`, tracked in this repository), installed with Spec Kit's extension mechanism. Its commands render as agent skills. **First set:** review (entry point), explain, reconcile (status vs tasks and evidence mismatches), unblock, verify-readiness. They link to each other and to the Spec Kit and roadmap commands. **Writes:** read-only by default. A change is proposed and applied only through Spec Kit commands (for example `/speckit.roadmap.write`) after the user approves each change (Art. VIII, C-06). **Timing:** specify now; implement after v0.1.2 is closed out. The `spec` binary does not change | Plain skill files outside Spec Kit; strictly read-only skills; including install, CI-setup and issue-reporting skills in the first set; deferring the spec |

## Recorded decision — role views in the review skill, 12 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-12-05 | Owner, 12 November 2026 | Personas for feature 005 are views of `speckit.spec-status.review`, selected with an argument: `pm`, `developer`, `qa`, `architect`. They are defined in one shared, tested file. A role reorders sections and prioritizes follow-ups, and never changes or hides facts. The two sample prompts become the `pm` and `developer` definitions. Still five skills | Four full persona skills (duplicated rules and routing); thin per-role wrapper skills in v0.2.0; leaving personas out of v0.2.0 |

## Recorded decision — partner extensions behind roles, 12 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-12-06 | Owner, 12 November 2026 | Feature 005 roles may recommend partner Spec Kit extensions when installed: architect → adrkit (community catalog, v0.1.5 when checked), qa → Assure (community, v2.7.5), developer → Bug Triage Workflow (`bug`, default catalog, v1.0.2). They are detected from the project's extension registry, recommended only by registered command names, and never required, installed or run without approval. Community-catalog partners carry a "not vetted by Spec Kit" note. Exact commands are confirmed from each repository during planning | Defining the mechanism without naming partners; deferring partner integration past v0.2.0 |

## Recorded decision — ADR partner for the architect role, 12 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-12-07 | Owner, 12 November 2026 | Amends D-2026-11-12-06. The architect role's partner is **Architecture Governance** (`arch-governance` 1.2.3, community catalog, MIT, `ashbrener/spec-kit-arch-governance`). It supports Spec Kit `>=0.1.0` and needs no extra tools. Its commands are `speckit.arch-governance.validate` (read-only), `.gate`, `.install`, `.sync`, `.repin` and `.issues`. **adrkit is dropped**: its 0.1.5 extension requires Spec Kit `<1.2.0` and the `adr` CLI, and no newer extension release supports 1.2. Revisit adrkit when it does | Keeping adrkit alongside arch-governance; vetting arch-governance with a throwaway install before choosing |

## Recorded decision — two-mode choices (Q6), 12 November 2026

| ID | Approver, date | Chosen | Not chosen |
|---|---|---|---|
| D-2026-11-12-08 | Owner, 12 November 2026 | Resolves Q6 for roadmap item 008 (two-mode roadmap and feature inventory), building on the 8 November direction (D-2026-11-08-01). - (1) **Command:** a separate subcommand, `spec inventory [--json] [--root]`. `spec status` stays the roadmap report. - (2) **Missing roadmap:** `spec status` with no roadmap file keeps exiting 2 with `SS_ROADMAP_UNAVAILABLE`; its message suggests `spec inventory`. There is no automatic switch. - (3) **Recognition:** a roadmap is supported if it has at least one supported item heading, or the pinned speckit-roadmap template's `## Planned Specs` heading. A file with the heading and no items is **explicitly empty**: exit 0, null denominators. A file with neither is **unrecognized**: a new fatal `SS_ROADMAP_UNRECOGNIZED`, exit 2. This replaces today's "zero items, exit 0" for narrative files and is listed in the release notes as a pre-1.0 behavior change. - (4) **Feature folders:** immediate children of `specs/` named `NNN-name` (3 or more digits) or `YYYYMMDD-HHMMSS-name`, the two forms Spec Kit creates. Other folders are reported as not feature folders, never guessed. - (5) **Checkboxes:** counted with exactly the same tasks.md rules as roadmap mode. Requirements checklists are excluded. - (6) **Output:** a separate inventory JSON document with `inventory_schema_version: "1.0"`, listing folders (spec, plan and tasks state; checked and total) and a summary. Exit codes keep their meaning: 0 below the threshold, 1 finding, 2 fatal. `spec status` schema 1.0 is unchanged. - (7) **Release:** v0.3.0. | `spec status --inventory` as a flag; automatically running inventory when the roadmap is missing; an explicit version-marker comment; treating an unrecognized roadmap as a warning with exit 0; requiring two template headings; counting every folder under specs/; a simpler column-0 checkbox rule; deferring 008 past v0.3.0 |

Where an earlier bullet in this file says every DC-01–DC-08 choice remains unapproved, the Recorded decisions table is the later record for the IDs it actually lists. D-2026-11-08-01 does not close the rest of DC-01–DC-08. It also does not change roadmap item 001.
