# Diagnostic registry — draft 0.1

**PROFILE-R1 candidate codes; not implemented or ratified.** Exact vocabulary, severity where not retained, fatal classification and source-span shape require DC-01–06 approval. Versioned meanings, not prose matching, are the intended machine contract. All codes share SS_ prefix to distinguish roadmap findings from uncontrolled operating-system strings.

## Severity and exit invariants

Retained severity assignments: duplicate ID error; unknown status warning; missing/self/cyclic dependency errors; ordinary unfinished alone no finding; blocked unfinished warning; planned missing spec link info; checked missing explicit local target error; unchecked missing target info. Other rows' proposed severities below must not be called owner approvals. A fatal error is always exit 2 independent of threshold; error-level nonfatal finding defaults to exit 1. All findings remain visible even with threshold none.

## Candidate registry

| Code | Severity / class | Trigger and source | Authority status |
|---|---|---|---|
| SS_INVOCATION_INVALID | error, fatal | Invalid command/option/value; no invented source line | Fatal meaning retained; code proposed |
| SS_CONFIG_INVALID | error, fatal | Invalid effective configuration if such config surface approved | Meaning retained; no config file added by this code |
| SS_ROOT_UNAVAILABLE | error, fatal | No selected root or required root unusable | Candidate classification |
| SS_ROADMAP_UNAVAILABLE | error, fatal | Required roadmap missing/unreadable | Required-input meaning retained |
| SS_ENCODING_INVALID | error, fatal | Unsupported/invalid encoding prevents interpretation | Proposed |
| SS_LAYOUT_UNSUPPORTED | error, fatal if inventory unreliable | Item-like unsupported table/header; source path/line | Visible failure retained, exact fatal boundary proposed |
| SS_ITEM_MALFORMED | error, fatal if inventory unreliable | Broken required item heading/fields | Proposed boundary |
| SS_DUPLICATE_ID | error, report | Duplicate roadmap occurrence IDs; include each source | Severity retained; occurrence counting proposed |
| SS_STATUS_UNKNOWN | warning, report | Unknown raw lifecycle retained | Retained |
| SS_DEPENDENCY_MALFORMED | error, report | Unsupported declared ID/reference | Proposed |
| SS_DEPENDENCY_MISSING | error, report | Referenced roadmap ID absent | Retained severity |
| SS_DEPENDENCY_AMBIGUOUS | error, report | Referenced ID duplicated | Proposed explicit case |
| SS_DEPENDENCY_SELF | error, report | Row depends on itself | Retained severity |
| SS_DEPENDENCY_CYCLE | error, report | Dependency cycle; show direct references/membership | Retained severity |
| SS_DEPENDENCY_BLOCKED | warning, report | Ordinary unfinished/deferred/abandoned target blocks selected policy | Retained blocked- unfinished rule; exact scope proposed Severity approved under D-2026-11-07-03. |
| SS_DEPENDENCY_HUMAN_REVIEW | warning, report | Any dependency on valid superseded row, even successor verified | Human review retained; severity/count treatment proposed Severity approved under D-2026-11-07-03. |
| SS_VERIFIED_DEPENDENCY_CONTRADICTION | error, report | Recorded verified row has unsatisfied/review dependency | Finding retained; severity proposed |
| SS_SUPERSESSION_MALFORMED | error, report | Non-single-ID/multiple replacement declaration | Validation retained; severity proposed |
| SS_SUPERSESSION_MISSING | error, report | Replacement ID missing | Validation retained; severity proposed |
| SS_SUPERSESSION_AMBIGUOUS | error, report | Replacement target duplicated | Validation retained; severity proposed |
| SS_SUPERSESSION_SELF | error, report | Replacement points to same item | Validation retained; severity proposed |
| SS_SUPERSESSION_CYCLE | error, report | Replacement cycle, original history retained | Validation retained; severity proposed |
| SS_SUPERSESSION_CHAIN_INVALID | error, report | Direct edge leads into invalid downstream chain | PROFILE-R1 proposal |
| SS_SPEC_REFERENCE_MISSING | info for planned; error candidate for specced/later | Stage lacks required spec link | Planned info retained; exact matrix/error choice proposed |
| SS_SPEC_DIRECTORY_MISSING | error, report | Explicit safe linked spec directory absent | Detection retained; severity proposed |
| SS_STAGE_FILE_MISSING | error under candidate matrix | Stage-appropriate spec/plan/tasks missing | Detection retained; stage matrix/severity proposed Severity approved under D-2026-11-07-03. |
| SS_TASK_TARGET_MISSING | error if checked, info if unchecked | Explicit supported local link target absent | Retained |
| SS_REFERENCE_UNSUPPORTED | warning, report unless required parse fatal | Explicit reference syntax outside approved grammar | Proposed, never turns URL/prose into path |
| SS_PATH_UNSAFE | error; fatal for required source | Escape/traversal/symlink/mount/regular-file policy violation | Confinement retained; classification proposed |
| SS_REFERENCE_SENSITIVE_SKIPPED | error, report; fatal if required source | Excluded sensitive reference, no content read | Exclusion/report retained; severity/class proposed |
| SS_INPUT_LIMIT | error, fatal on required input/work | Approved bytes/files/path/nesting/items/edge/task budget exceeded | Explicit limit failure retained; exact limit/class proposed |
| SS_SOURCE_CHANGED | error, fatal | Detectable source change during input read | Proposed no-authoritative-snapshot policy |

No quarter diagnostic/code is defined. No finding is a mutation instruction or semantic code verdict. Codes emitted for a single issue must not inflate summary blocker count; e.g. a human-review dependency may also support a verified-contradiction finding, but the active row counts once.

## Location, redaction and ordering

Use root-relative normalized path and 1-based physical line where obtained; no absolute host paths, secret contents or raw OS stack traces. Missing source uses null, not an invented line. Related references preserve original dependency and successor/invalid chain for review, not rewrite. Diagnostic ordering follows JSON contract; every platform must produce equivalent safe meanings even when OS-specific errors differ.

## Recovery policy proposal

Required root/roadmap/input encoding, unreliable item inventory, required unsafe/limit failures and observed source changes are candidate fatal outcomes with no authoritative totals. Ordinary structural references can produce a complete supported-scope report plus findings. A skipped optional linked file is visible coverage and error, not equivalent to zero tasks or a verified outcome. Unsupported layout must not be ignored merely to keep exit 0. Any different approved recovery policy must adjust JSON completeness/summary semantics and fixtures explicitly.

## Approval and testing

Approve registry alongside format, matrix, safety matcher and JSON schema. Fixtures assert code/severity/source and threshold behavior, not only free-form messages. Unknown future codes must not cause consumers to reinterpret exit 1 as a crash. Registry additions/meaning changes follow documented CLI/JSON compatibility policy. No tests or security checks have been run by drafting this registry.
