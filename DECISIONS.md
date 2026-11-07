# Clarifications and decision log

**Draft.** The owner has authorized drafting and public-project documentation, not implementation or release. This document distinguishes requirements from proposed choices. PROFILE-R1 is one consistent recommended candidate; **none of its unresolved choices is ratified by drafting it**.

## Settled boundaries

- Omit inferred quarters entirely from v1. Preserve recorded timing, without inventing dates or treating question IDs as calendar data.
- Retain strict v1 history-preserving supersession and human review of its dependents; no automatic dependency redirection.
- Offline read-only Go standard-library core, `CGO_ENABLED=0`; no runtime AI/network/subprocess/Git/credentials/writes/extensions.
- Preserve raw statuses and unknowns, recorded-verified versus checkbox progress, JSON/text parity, default threshold error and exits 0/1/2.
- Future extension write-back is optional, separate and per-mutation approved; optional integrations remain deferred.
- Public project documentation targets [tosin2013/spec-status](https://github.com/tosin2013/spec-status). Repository establishment is separate from this package; no license, code implementation or release is created here.

## Proposed profile: approval register

| ID | Needs choosing / alternatives | PROFILE-R1 recommendation, not approved | Consequence / implementation gate |
|---|---|---|---|
| DC-01 | Exact input headings/fields, whitespace, nested bullets, version marker; alias/tolerant parser versus narrow dialect; root/override bases | Narrow v0.1.1-shaped headings and documented one-field extension; exact examples in ROADMAP-FORMAT. BOM/CRLF normalization, no aliases/tables. Optional advisory compatibility marker, no migration required. Explicit root then nearest ancestor discovery; roadmap override root-relative, spec dir root-relative, task links tasks-file-relative. | Approve real fixture examples, uniqueness and source-base semantics before parser implementation; malformed item-like content invalidates reliable totals. |
| DC-02 | Metric effect of invalid replacements; legal chains; mandatory provenance versus optional/future fields | Only validated references exclude active rows; invalid supplied reference retained, diagnosed and not excluded solely by its presence. Allow acyclic same-roadmap chains, preserve each direct edge and optionally show terminal ID; no redirect. Do not require new reason/date/decision fields. | Approve graph-validity propagation and ambiguity cases; supersession itself remains in v1. Optional provenance extension needs a later format decision. |
| DC-03 | Human-review severity and blocked counting; dependency satisfaction alternatives; lifecycle-file matrix and stable codes | Warning HUMAN_REVIEW, distinct review count and once-per-dependent blocked count; verified-only default, alternate implemented-or-verified policy remains uncommitted. Candidate stage matrix in format doc; candidate codes in diagnostics. | Approve stage matrix, contradiction severity, blocked scope and code registry. No auto-satisfaction through successor. |
| DC-04 | Name/config/schema/rounding/task scope; retained --as-of use; new aliases versus existing interface | Working `spec status`; no config file initially, explicit flags then defaults. No --as-of without a concrete non-quarter use; no clock dependence. Candidate JSON schema 1.0, two-decimal half-up percentages, stable snake_case. Unique linked task-file aggregation with separately labelled per-item reuse. | Command name, schema/envelope/canonical bytes and task/reference grammar require confirmation. Retained --json/exit meanings are not reopened. If date feature is chosen, revise all related contracts before implementation. |
| DC-05 | Contained symlinks versus reject all linked symlinks; mount and mid-read behavior; sensitive exclusions and patched Go baseline | Regular-file-only reads, reject symlink components in referenced paths initially; deny cross-mount referenced reads where platform can prove policy. Fail required-input/source-changed report rather than claim atomicity. Concrete deny examples below; supported patched toolchain selected only after current evidence review. | A capability gap must fail closed, not downgrade confinement. Root canonicalization, case behavior, matcher, race safety and platform-equivalent enforcement require approval. No API/toolchain version selected by this draft. |
| DC-06 | Fixed versus configurable limits; fatal versus partial totals | Candidate caps: 2 MiB/file; 32 MiB total actual bytes read; 1,024 distinct opened input files; 4,096-byte normalized path representation; 32 parser nesting levels. Fixed initial profile; explicit limit error, no authoritative partial totals. | New proposed numbers, uncalibrated. Approve/correct against fixtures, including 100 items/10,000 tasks. Include lookup/edge/time-work budgets in final review; byte caps alone are not sufficient. |
| DC-07 | Initial two versus five/six targets; checksum-only versus additional signed identity gate; OS minimums/rollback | Candidate linux/amd64 + darwin/arm64 only after real runners prove them; other targets not advertised supported. Direct GitHub Release archives, exact commit, pinned tooling, gated checks and owner RC approval retained. Truthful checksum-only integrity initially; attestation/signing a separate owner choice. | Actual runners, patched toolchain, support window, OS versions, packaging rules, permission review and identity policy require approval before release work. No binary support claimed today. |
| DC-08 | Dogfood targets; public artifact versus feature expansion; license/security/support/fork expectations | Personal utility remains baseline; public-documentation intent accepted, not a market commitment. Proposed approved Spec Kit roadmap plus one owner-named portfolio roadmap, hand-count parity. Public repository is tosin2013/spec-status; repository operations remain outside this package. | Name approved inputs. Repository home is established; license, security disclosure route, support window and mirror/Bitbucket expectations remain unresolved. No license file invented. |

## PROFILE-R1 refinements requiring explicit approval

These make contracts draftable; they are new technical recommendations, not hidden owner decisions.

- **D-01A Format:** ID is exactly three ASCII digits, unique within one roadmap. Heading punctuation is an em dash, fields at one bullet level; backtick-wrapped or bare local paths accepted only as declared in the format doc. No alias status normalization. Mixed invalid UTF-8 is fatal; optional UTF-8 BOM and CRLF accepted.
- **D-02A Graph validity:** a replacement chain is valid only if all nodes/edges to its end resolve uniquely and are acyclic. An edge into an invalid chain does not exclude its source. Duplicate ID rows are retained as occurrences; ambiguous references fail validation. No automatic terminal-successor inference of dependency satisfaction.
- **D-03A Blocker scope:** show per-item direct blockers for all rows; summary blocked/review counts cover active rows only. A row with multiple reasons counts once. A verified-but-blocked row still contributes its recorded-verified numerator and emits contradiction; do not silently recategorize its label.
- **D-03B Stage candidate:** planned and earlier do not require spec links; specced requires link/spec.md; in-progress/implemented/verified additionally require plan.md/tasks.md in candidate matrix. Missing linked artifacts do not establish code absence. This is a proposed stricter mapping, not a source-settled stage requirement; choosing a minimal mapping is a valid alternative that needs explicit acceptance changes.
- **D-04A Error model:** fatal input gives `complete:false`, null summary, no authoritative items, sorted safe diagnostics and structured fatal `error`. Normal result has `complete:true`; skipped optional evidence is visible through coverage/findings, not a claim of complete repository inspection. No OS error strings or secret input snippets on stdout.
- **D-05A Sensitive path candidate:** deny any path component `.git`, `.ssh`, `.aws`, `.gnupg`; deny `.env` / `.env.*`, case-insensitive secret/credentials/private-key/answer-key basename patterns, and `.pem`, `.key`, `.p12`, `.pfx` suffixes. This over-inclusive starting matcher is incomplete and must be reviewed; it is not a security certification or permission to open unmatched secrets. No content sniffing of excluded files. Approve normalization/case policy per supported platform.
- **D-05B Root and races:** explicit root must exist; recommend canonicalize selected root once, then enforce all later paths beneath it through a race-resistant safe-open boundary. No promise of an atomic multi-file snapshot. If observable source identity/size/mtime changes during a read, invalidate result. Same-size undetectable edits are not claimed impossible; acceptance requires documented limits or stronger snapshot design, not a false guarantee.
- **D-06A Work budget:** 1,024 decoded roadmap items, 16,384 total dependency/replacement/reference edges and 100,000 recognized task checkboxes are additional candidate work caps for approval. Algorithmic cycle checks must avoid unbounded recursion. Final path resolution and read-call budgets must be measurable. Configurable overrides, if approved later, must never disable root confinement or sensitive exclusions.

## Readiness ledger

| Gate | Current state |
|---|---|
| Documentation drafting | Authorized; this package is draft work |
| Required review returns | Received; semantic/compatibility decisions remain, not missing replies |
| Source certainty | Requirements coverage assembled; original input-format prose partly condensed, not fully verbatim recovered |
| PROFILE-R1 adoption | Pending owner/technical clarification, per DC-01–08 |
| Implementation start | Not authorized by this package; blocked on relevant decisions and separate direction |
| Public repository | tosin2013/spec-status established separately; no repository operations performed by this package |
| Release | No authorization, workflow, binary or tests; release evidence and approvals still required |

## Change procedure

Record a choice with decision ID, approver, date, chosen alternative, rejected alternatives and impacted requirements/tests. Do not replace “proposed” with “approved” merely because it appears in several documents. A scoped decision may clear one gate without clearing implementation/release. Revisions of input/schema compatibility must identify migration/fixtures. Licensing decisions must be explicit; public visibility does not imply a permissive license.
