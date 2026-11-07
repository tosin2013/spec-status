# JSON report contract — proposed schema 1.0

**Draft PROFILE-R1, not a ratified or implemented JSON Schema.** This Markdown contract proposes one versioned interoperable result. Field shapes, code registry, rounding and byte format require DC-02–04 approval before implementation. Examples are illustrative fixture expectations, not real runtime output.

## Version and compatibility

Propose top-level `schema_version: "1.0"`; not CLI SemVer or input v0.1.1. Within a stable schema major, additive optional fields may be non-breaking only if documented consumer rules permit ignoring them. Reinterpreting status, denominators, exit meanings or existing field types is breaking. Adding a diagnostic code requires documented unknown-code tolerance and does not silently change threshold meaning. Input compatibility matrix stays independently documented.

Human-text filters and `--explain` do not change this document. Passing `--status`, `--blocked`, `--human-review`, `--active`, `--historical`, or `--explain` with `--json` still emits the full schema 1.0 payload for the unfiltered result. Those flags add no field, remove no item, and do not reinterpret a denominator. A later schema that exposed a view would need an explicit additive field or a new schema version; this contract does not add one.

Retain `superseded_by` separate from raw lifecycle status and `--json` invocation. No camelCase replacement or --format alias. No quarter/inference fields, even null placeholders. Explicit timing is preserved only as supported source text; no new normalized date fields are adopted by this draft. No private absolute paths, secret contents, volatile timestamps or process-dependent IDs.

## Root document (candidate required fields)

| Field | Type / semantics |
|---|---|
| schema_version | string, candidate "1.0" |
| outcome | "report" or "fatal" |
| complete | boolean: usable report under declared supported read scope, not proof of full project truth |
| source | object: `roadmap` relative path or null; `input_format` string or null; root is implicit selected root, not absolute path |
| configuration | object: `fail_on`, `dependency_policy`, `task_scope`; declared candidate effective settings, no environment dumps |
| summary | summary object on report; null on fatal, not a misleading partial percentage |
| items | array, document order on report; [] on candidate fatal envelope |
| diagnostics | array of safe findings including fatal diagnostics, sorted as below |
| error | null for report; `{code,message,source}` for fatal; stable generic message, no raw OS/secret error text |

`complete:true` allows explicitly reported skipped optional evidence; it means the supported inventory is structurally interpretable. Coverage fields enumerate skipped/missing files. If a required roadmap cannot be interpreted or a required limit is exhausted, candidate fatal `complete:false` invalidates all authoritative totals rather than publishing undercounted progress. Exact partial-result policy is DC-05/DC-06; no silent truncation.

## Summary object

| Field | Candidate type / meaning |
|---|---|
| total_items | nonnegative integer, roadmap occurrences including historical and duplicate occurrences |
| status_counts | ordered object raw status -> nonnegative count; includes all occurrences, unknown raw keys retained |
| active_items | nonnegative integer after approved lifecycle/supersession exclusions |
| verified_active_items | nonnegative integer active rows recorded verified; not independent verification |
| verified_percent | number 0..100 or null when active_items zero |
| superseded_items | nonnegative integer count valid superseded rows under PROFILE-R1; no completion credit |
| invalid_supersession_items | nonnegative integer; supplied invalid reference, not excluded solely by it in candidate profile |
| blocked_count | count of distinct active dependent rows blocked/review-required under chosen DC-03 policy |
| human_review_count | distinct active dependent rows requiring supersession review, subset of blocked_count in PROFILE-R1 |
| tasks | `{checked,total,percent,scope,files_read,files_missing,files_skipped}` with nonnegative integer counts, nullable percent, explicit unique-linked-file scope |

Raw status_counts include historical verified rows; verified_active_items excludes valid superseded rows. A verified active row with blockers retains its recorded label/numerator while showing a contradiction—this does not assert readiness. Unknown labels remain active/not verified. Null denominator never maps to 100%. Scope all unique explicitly linked tasks files is a proposal, not silently implied active-task progress. Per-item task counts can reuse files; aggregate deduplicates source files rather than summing duplicated attribution.

## Item object

| Field | Candidate type / meaning |
|---|---|
| id, title, status | strings preserving source ID/title/raw status |
| status_known | boolean, no remapping |
| source | `{path,line}`; path root-relative; line positive 1-based physical line |
| active | boolean under chosen metric policy |
| description, outcome, scope | raw supported text, absent candidate fields represented null, not inferred prose |
| spec_dir | root-relative normalized path or null; unsafe supplied value must not be opened or leaked verbatim if sensitive |
| superseded_by | parsed direct target ID or null; never substituted for lifecycle |
| supersession | `{supplied,raw_reference,valid,terminal_id}`; valid boolean or null if absent; raw_reference sanitized/null on unsafe data; terminal_id only when chain valid and chosen to display |
| dependencies | array of direct dependency result objects in supplied order |
| tasks | `{checked,total,percent,source,state}`; counts/percent nullable when missing/skipped, state read/missing/skipped/not_linked |

Dependency result candidate: `id`, `satisfied`, `state`, `replacement_id`, `source`. State is satisfied/unfinished/missing/ambiguous/self/cycle/human_review. `replacement_id` is display only; supplied original `id` is never rewritten. A dependency on a valid superseded row is `satisfied:false`, state human_review even if original/successor verified. Invalid replacement elsewhere does not invent a successor; malformed-reference diagnostic remains visible. If multiple issues overlap, preserve individual diagnostics and choose primary state deterministically (self/cycle/ambiguous/missing before satisfaction); priority is a DC-03 refinement to approve.

Supersession raw_reference, duplicate fields and malformed values must be preserved sufficiently for mechanical diagnostics without printing credentials or arbitrary uncontrolled text. Never include secret snippets in JSON. Unknown lifecycle tokens are intentionally raw but require normal JSON escaping; length limits prevent unbounded output.

## Diagnostic and fatal source shape

Diagnostic candidate fields: `code`, `severity` (info/warning/error), `message`, `source` (`{path,line}` or null), `item_id` nullable, `related` array of safe source references. Codes carry semantics; messages are stable generic human explanations, not the only machine contract. `related` may show original dependency/replacement locations, never infer scope equivalence. A malformed source with no physical line uses null source or null line under the final approved schema; do not fabricate line 1.

Sorting candidate: path (null first, forward-slash normalized), line (null before positive), code, item_id (null first), safe related-reference ordering, then message. Status map keys lexicographic by raw UTF-8 bytes; items remain occurrence document order; dependency arrays preserve declared order. Exact tie-break behavior and cross-platform path-case treatment require fixtures. Do not rely on unspecified Go map iteration.

## Canonical serialization and rounding proposal

UTF-8 JSON, fixed field order as contract tables, no indentation, standard JSON escaping, one terminal LF; no logs stdout. Numeric percentages use integer rational calculation rounded half-up to two decimal places, emit no unnecessary trailing zeros; count integers have no float conversion. For example 1/3 becomes 33.33; 2/3 66.67; zero denominator null. Text uses the same rounded result and N/A for null. Two runs with identical files/configuration must be byte-identical under candidate clock-independent profile; platform canonicalization still needs actual evidence.

## Minimal normal example (candidate fixture)

```json
{"schema_version":"1.0","outcome":"report","complete":true,"source":{"roadmap":".specify/memory/roadmap.md","input_format":"speckit-roadmap-v0.1.1+spec-status-supersession-draft"},"configuration":{"fail_on":"error","dependency_policy":"verified","task_scope":"unique_linked_task_files"},"summary":{"total_items":0,"status_counts":{},"active_items":0,"verified_active_items":0,"verified_percent":null,"superseded_items":0,"invalid_supersession_items":0,"blocked_count":0,"human_review_count":0,"tasks":{"checked":0,"total":0,"percent":null,"scope":"unique_linked_task_files","files_read":0,"files_missing":0,"files_skipped":0}},"items":[],"diagnostics":[],"error":null}
```

Empty recognized roadmap is not a missing required input. Future nonempty fixtures must include every candidate item field; this minimal example makes no parser implementation claim.

## Fatal example (candidate fixture, exit 2)

```json
{"schema_version":"1.0","outcome":"fatal","complete":false,"source":{"roadmap":null,"input_format":null},"configuration":{"fail_on":"error","dependency_policy":"verified","task_scope":"unique_linked_task_files"},"summary":null,"items":[],"diagnostics":[{"code":"SS_ROOT_UNAVAILABLE","severity":"error","message":"Selected root is unavailable.","source":null,"item_id":null,"related":[]}],"error":{"code":"SS_ROOT_UNAVAILABLE","message":"Selected root is unavailable.","source":null}}
```

Fatal outcome and exit 2 survive fail_on none. Output state and exit policy are related, not duplicates: report outcome can exit 0 or 1 depending on diagnostics/threshold, never quietly convert fatal failure into a report. A failure to serialize at all is an acceptance failure, not a valid example of the promised JSON contract.
