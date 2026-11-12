# Acceptance and test strategy — documentation-stage plan

**Product tests have been run by the implementer; they are not independent verification.** The cases below are PROFILE-R1 acceptance cases; the decisions they waited on were approved as built (DC-03: D-2026-11-07-03; DC-01, DC-02, DC-04–DC-08: D-2026-11-11-02), except DC-08's mirror expectations and release-engineering choices, which remain open. Documentation checks alone cannot establish parser correctness, safe reads, runtime behavior or release readiness. “Pass” is not recorded anywhere in this test plan.

## Test levels and evidence

1. Approved fixture corpus with source input, effective configuration, expected text/JSON/counts/codes/exit and known unsupported variants. Synthetic examples first; approved real roadmap snapshots later. No private fixtures may be published without separate permission/sanitization.
2. Parser/graph/progress fixtures and Go standard-library testing/fuzzing where suitable, with bounded resource expectations. No AI oracle or semantic code inference.
3. End-to-end binary cases for both output modes, exits 0/1/2 and external observation of no runtime writes/network/subprocess. Test harness activities do not become product capabilities.
4. Platform-specific filesystem/race cases on every chosen v1 target (linux/amd64, darwin/arm64; D-2026-11-11-02), not just one Unix development host.
5. Release archive/startup/checksum/tag/identity-policy/install evidence on each advertised OS/architecture. Cross-compilation alone is never evidence of execution.

Each evidence record should identify fixture revision, selected profile decisions, product commit/build/Go version, platform/hardware where relevant, expected versus observed outcome, limitations and reviewer. No fabricated evidence log, golden file or execution result is included here.

## Acceptance cases

| ID | Inputs / action | Expected outcome and limits |
|---|---|---|
| AC-01 Standalone happy path | Valid supported roadmap with explicitly linked files, extension absent, one command in text/JSON | Same supported-scope facts/counts, source lines, no extension/Git/agent/network/subprocess. Recorded status not independently verified. |
| AC-02 Status inventory | All nine statuses, unknown raw token including case variation, duplicates | Known raw labels retained; unknown warning, active/not verified; duplicates error, occurrences retained and references ambiguous under PROFILE-R1 (D-2026-11-11-02). No alias remapping. |
| AC-03 Empty and weighted progress | Empty supported roadmap; all inactive; task files with 1/1 and 0/9 checked; shared tasks file | Active/task empty fractions null/N/A; weighted unique-file aggregate 1/10, not 50%; shared file counted once aggregate with per-item scope disclosed. Missing input not equated to empty file. |
| AC-04 Parser/layout | Heading/bullet examples, nested content, BOM/CRLF, valid UTF-8, invalid encoding, tables and malformed item-like headings, ordinary prose | Approved grammar recognized, physical line mapping stable; unsupported item-like input visibly fails rather than undercounts. Fatal inventory failure has null summary. Non-item prose does not become a row/path. |
| AC-05 Fences/references | Backtick/tilde fences containing fake rows/tasks; prose/inline filenames, URL/anchor, supported relative link, unsupported link form | Fenced content ignored for recognition; only explicit supported local references considered. No network or guessed path from URL/prose. Unsupported reference policy visible. |
| AC-06 Stage/spec files | Each chosen lifecycle status with/without link, directory, spec.md/plan.md/tasks.md | Exact approved matrix applied. Planned absent spec link info, no future-stage false errors. DC-03 settled the matrix (D-2026-11-07-03). |
| AC-07 Task targets | Identical missing explicit local target on checked and unchecked task; safe existing target | Checked missing error, unchecked info; task count once irrespective link count; existence check not semantic correctness. Unsafe target handled under safety cases, not blindly read. |
| AC-08 Dependencies | Verified/implemented/unfinished/deferred/abandoned/unknown targets, missing/ambiguous/self/cycles, multiple blockers | Default verified-only satisfaction; missing/self/cycles errors; blocked unfinished warning; cycle membership/direct blockers shown. Optional alternate policy only if explicitly adopted and tested. No date forecast. |
| AC-09 Valid history | Historical item recorded verified/implemented, optional valid successor forward/backward, no replacement field | Old ID/status/source retained, direct replacement distinct in text/JSON; valid superseded rows in inventory but not active completion numerator/denominator; absence alone leaves normal lifecycle behavior. |
| AC-10 Invalid/chain history | Malformed/duplicate replacement fields, nonexistent/ambiguous target, self/cycle, valid chain and chain into invalid node | Approved DC-02 behavior; PROFILE-R1 valid-only exclusion and downstream validity propagation; invalid supplied data retained and diagnosed. Typos cannot hide rows. Chain end never rewrites dependencies. |
| AC-11 Human-review dependency | Dependent points to superseded original; direct and terminal successor verified; original itself verified | Always HUMAN REVIEW, original edge and displayed successor preserved, never satisfied/redirected by successor. Warning plus distinct active review count (D-2026-11-07-03), once-per-row blocked count. Inactive rows still show per-item reasons. |
| AC-12 Contradiction/counts | Verified active dependent has unsatisfied/review edge; multiple error reasons; historical verified superseded row | Contradiction finding; raw status not overwritten; active verified numerator follows recorded label, not inferred readiness. Blocked summary scope and dedup match chosen policy. Historical verified does not earn active credit. |
| AC-13 Root/path safety | Relative/absolute overrides, nearest ancestor root, sibling-prefix confusion, ../ escape, drive/case boundaries, invalid root | All accesses remain confined; failed required input exit 2 (D-2026-11-11-02). No “fallback root” that silently widens scope. Path bases match approved examples. |
| AC-14 Symlink/regular-file safety | Internal/external/broken symlink, trailing-slash symlink escape, FIFO/device/socket, mount boundary, unsupported enforcement | Chosen conservative policy enforced before content read; no hang or safe-open downgrade. Each supported platform demonstrates actual protection. Reviewer-reported advisory fixtures require technical validation, not assumption. |
| AC-15 Sensitive exclusions | D-05A secret/credential/private-key/answer-key names, nested denied components, case variants, encoded-looking references | Excluded contents never opened, not merely hidden after read. Safe diagnostic without secret content/absolute path. Over- and under-match cases reviewed; no claim matcher recognizes every possible secret. |
| AC-16 Limits/change safety | At/over file-byte/total-byte/file-count/path/nesting/item/edge/task caps; changing files and repeated edits | Explicit bounded failure, no unbounded recursion/retry or authoritative truncated totals. Required budget/source-change failure exit 2 (D-2026-11-11-02). Document residual non-atomic snapshot limits; do not claim detecting all same-metadata edits. |
| AC-17 JSON/text parity | Identical inputs/config repeatedly, map ordering/ties, rounding 1/3, 2/3, zero denominator, Unicode | Same facts/rounding; byte-stable JSON, stable source/diagnostic order, no volatile dates/host paths/OS error text. Cross-platform equality tested where supported, not assumed. |
| AC-18 Exit/schema compatibility | Every severity × threshold, none plus fatal, invalid invocation with --json, required parse failure | Exact exits 0/1/2. One valid JSON document on report and fatal; stdout uncontaminated, fatal summary null. None does not suppress fatal. Retained --json/superseded_by meanings; version policy tested. |
| AC-19 Omitted timing/runtime | Search interface/result behavior for inferred quarter fields/flags/denominators; observed runtime operations | No quarter inference or opt-in, question IDs not dates; explicit timing retained only as supported raw data. No writes/cache/network/subprocess/Git/credentials/AI/extensions. Observation is necessary; documentation assertions insufficient. |
| AC-26 Human-text views | `--status`, `--blocked`, `--human-review`, `--active`, `--historical`, alone and combined, including an unknown raw token, a repeated token, a no-match, and a historical row separate from an active row; `--explain` on rows with and without a retained line; the same flags with `--json` and with `--fail-on` | Summary denominators, diagnostics, and exit 0/1/2 stay on the full result. Human text labels `shown N of M`. No-match is a message, not a fatal, unless a separate finding already fails the threshold. Explain lines cite only parsed file:line locations and say when a line was not retained. JSON bytes match the unfiltered schema 1.0 document. Repeated identical flags are exit 2. |
| AC-20 Cold-start target | Documented fixture with 100 roadmap items/10,000 recognized tasks, extension absent; repeated cold invocations on disclosed hardware/OS | Median <100 ms proposed target. Record startup/cache methodology, sample count and distribution, fixture bytes/caps and target configuration. No present measured claim or one-day delivery estimate. |
| AC-21 Release identity/gates | Approved candidate tag/commit, metadata mismatch, intentionally failed required check, least-privilege review | Failed checks/mismatched version block publication. Exact commit, Go target/version and dev distinction verified; only publisher job write permission. A real publication attempt requires separate approval. |
| AC-22 Asset/platform integrity | Every advertised archive, checksum manifest, actual target runner; unadvertised targets | Extraction/layout/version/help/sample text/JSON/startup tested, checksum matches exact asset/tag/commit. Compile success not runtime proof; signed identity/provenance/reproducibility claimed only with matching verified evidence. |
| AC-23 Distribution/immutability | No-Go/source-checkout clean user environment; user-owned PATH install/upgrade; correction version policy | Rootless install/upgrade docs independently exercised. Existing published bytes never silently replaced; fixes new version. Rollback and OS/minimum/quarantine cases added if DC-07 adopts them. No package-manager promises. |
| AC-24 Dogfood/docs/governance | Approved Spec Kit input plus second owner-selected roadmap, manual inventory/findings; README audit | Hand-count parity proposed, approved sources only. README truthful planning-stage until evidence changes; no invented license, support or releases. Public docs leak no private coordination/paths. All choices/requirements traceable before implementation approval. |
| AC-25 Future mutation boundary | Status/JSON/CI with integration unavailable; separately proposed writer attempts unchanged/changed source approvals | Core never invokes integration. Future mode approval covers each ID/change/path; source changes require reapproval; proposed/attempted/confirmed distinct, persisted reread required, no LLM-as-approval. Out of v1 implementation scope. |

## Metrics and oracle discipline

Hand counts distinguish raw inventory from active eligibility and recorded verification. Tests cannot “prove” arbitrary code satisfies a spec; fixtures compare declared evidence and mechanical findings only. Supersession evidence should include completed historical rows to catch inflated percentages. Task-file aggregation must state inclusion of inactive rows and deduplication. Invalid data that makes inventory unknown must not receive a falsely precise denominator.

Performance is an acceptance target, not an implementation estimate. Source format compatibility depends on approved representative fixtures, not the fact a draft contains examples. Toolchain advisory references and platform capabilities need current independent validation before acceptance. A release job's own checksum is insufficient independent identity proof.

## Acceptance readiness

Block final implementation acceptance if grammar/stage matrix, invalid replacement metrics, HUMAN_REVIEW count policy, schema/code meanings, safe-read limits/matcher or toolchain policies remain ambiguous. Block release acceptance if support targets, runner evidence, pins, archive manifest or RC approval missing. As of 11 November 2026 the grammar/stage matrix, invalid replacement metrics, HUMAN_REVIEW count policy, schema/code meanings and safe-read limits/matcher are approved as built (D-2026-11-07-03, D-2026-11-11-02) and the v1 support targets are chosen (D-2026-11-11-02); toolchain policy, runner evidence, pins, archive manifest and RC approval are still missing. A spec draft may be complete as documentation while these implementation/release gates remain unresolved. No [x] completion state or passing-test badge is appropriate today.

## Two-mode sketches — not acceptance of the current build

These rows record D-2026-11-08-01. They are not executed tests and not a claim that pull request 15 implements them. Details and the open choices are in [two-mode-direction.md](two-mode-direction.md).

| Sketch | Later outcome, once syntax and schema are chosen |
|---|---|
| Supported template | Raw statuses, dependencies, supersession, and linked checkbox counts. A wrong result on a pinned shape is a supported-format bug. |
| Empty supported roadmap | Zero items, null denominators, not a missing-file fatal. |
| Unrecognized file | Located diagnostic and compatibility guidance. Not called invalid Spec Kit. Not reported as an empty supported roadmap. |
| Malformed supported read | Located fatal or diagnostic. No invented rows. No fall-through to inventory. |
| Missing roadmap | Roadmap mode stays a required-input failure. A separate explicit inventory may then list bounded feature folders and file presence, plus checked and total tasks, excluding requirements checklists. |
| Inventory limits | Missing, unreadable, and unsupported files stay distinct. Checkbox completion does not verify implementation or invent a lifecycle. |

## Draft proposal sketches — not product acceptance

These rows record the 7 November 2026 drafting request. They are not AC cases, not executed tests, and not implementation approval. Details are in [timing-evidence.md](timing-evidence.md) and [companion.md](companion.md). Illustrative IDs and date shapes in those drafts are not facts about this repository.

| Sketch | Proposed outcome if later adopted |
|---|---|
| Timing, unknown status | Raw token kept, warning, not recorded verified; otherwise the active rules still apply. No inferred date. |
| Timing, valid supersession | History, ID, and raw status kept; excluded from the active denominator; not reported as success. Zero active is not 100. |
| Timing, superseded dependency | Human review, original edge kept, no redirect. |
| Timing, bad graph | Self, missing, and cycle references are diagnostics. A missing reference is not proof that code failed. |
| Timing, filters | Full totals, findings, and exit stay on the unfiltered result. |
| Timing, malformed explicit value | Reported as malformed. Not repaired into a quarter or an ETA. |
| Companion, opt-in | Ordinary `spec status` JSON unchanged. |
| Companion, bad prompt file | Exit 2, no partial envelope, no default prompt, no write, no model call. |
| Companion, envelope | Same common status model. No computed quarter or ETA. |
| Companion, malicious prompt | Does not elevate harness permissions. A harness boundary test would be required. The CLI alone would not certify that. |
