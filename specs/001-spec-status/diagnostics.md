# Diagnostic registry — 0.1

**PROFILE-R1 codes, implemented in the development build and approved as built.** The code registry and report severities were approved with DC-03 on 7 November 2026 (D-2026-11-07-03); fatal classification, confinement and limit handling with DC-01, DC-02 and DC-04–DC-06 on 11 November 2026 (D-2026-11-11-02). Versioned meanings, not prose matching, are the intended machine contract. All codes share SS_ prefix to distinguish roadmap findings from uncontrolled operating-system strings.

## Severity and exit invariants

Retained severity assignments: duplicate ID error; unknown status warning; missing/self/cyclic dependency errors; ordinary unfinished alone no finding; blocked unfinished warning; planned missing spec link info; checked missing explicit local target error; unchecked missing target info. The other rows' severities below were approved as built under D-2026-11-07-03. A fatal error is always exit 2 independent of threshold; error-level nonfatal finding defaults to exit 1. All findings remain visible even with threshold none.

## Registry

| Code | Severity / class | Trigger and source | Authority status |
|---|---|---|---|
| SS_INVOCATION_INVALID | error, fatal | Invalid command/option/value; no invented source line | Fatal meaning retained; Approved as built (D-2026-11-11-02) |
| SS_CONFIG_INVALID | error, fatal | Invalid effective configuration if such config surface approved | Reserved and not emitted: there is no config file (DC-04, D-2026-11-11-02) |
| SS_ROOT_UNAVAILABLE | error, fatal | No selected root or required root unusable | Approved as built (D-2026-11-11-02) |
| SS_ROADMAP_UNAVAILABLE | error, fatal | Required roadmap missing/unreadable | Required-input meaning retained |
| SS_ENCODING_INVALID | error, fatal | Unsupported/invalid encoding prevents interpretation | Approved as built (D-2026-11-11-02) |
| SS_LAYOUT_UNSUPPORTED | error, fatal if inventory unreliable | Item-like unsupported table/header; source path/line | Visible failure retained; Approved as built (D-2026-11-11-02) |
| SS_ITEM_MALFORMED | error, fatal if inventory unreliable | Broken required item heading/fields | Approved as built (D-2026-11-11-02) |
| SS_DUPLICATE_ID | error, report | Duplicate roadmap occurrence IDs; include each source | Severity retained; occurrence counting approved as built (D-2026-11-07-03) |
| SS_STATUS_UNKNOWN | warning, report | Unknown raw lifecycle retained | Retained |
| SS_DEPENDENCY_MALFORMED | error, report | Unsupported declared ID/reference | Approved as built (D-2026-11-07-03) |
| SS_DEPENDENCY_MISSING | error, report | Referenced roadmap ID absent | Retained severity |
| SS_DEPENDENCY_AMBIGUOUS | error, report | Referenced ID duplicated | Approved as built (D-2026-11-07-03) |
| SS_DEPENDENCY_SELF | error, report | Row depends on itself | Retained severity |
| SS_DEPENDENCY_CYCLE | error, report | Dependency cycle; show direct references/membership | Retained severity |
| SS_DEPENDENCY_BLOCKED | warning, report | Ordinary unfinished/deferred/abandoned target blocks selected policy | Retained blocked-unfinished rule; scope and severity approved under D-2026-11-07-03 |
| SS_DEPENDENCY_HUMAN_REVIEW | warning, report | Any dependency on valid superseded row, even successor verified | Human review retained; severity and count treatment approved under D-2026-11-07-03 |
| SS_VERIFIED_DEPENDENCY_CONTRADICTION | error, report | Recorded verified row has unsatisfied/review dependency | Finding retained; approved as built (D-2026-11-07-03) |
| SS_SUPERSESSION_MALFORMED | error, report | Non-single-ID/multiple replacement declaration | Validation retained; approved as built (D-2026-11-07-03) |
| SS_SUPERSESSION_MISSING | error, report | Replacement ID missing | Validation retained; approved as built (D-2026-11-07-03) |
| SS_SUPERSESSION_AMBIGUOUS | error, report | Replacement target duplicated | Validation retained; approved as built (D-2026-11-07-03) |
| SS_SUPERSESSION_SELF | error, report | Replacement points to same item | Validation retained; approved as built (D-2026-11-07-03) |
| SS_SUPERSESSION_CYCLE | error, report | Replacement cycle, original history retained | Validation retained; approved as built (D-2026-11-07-03) |
| SS_SUPERSESSION_CHAIN_INVALID | error, report | Direct edge leads into invalid downstream chain | Approved as built under DC-02 and D-02A (D-2026-11-11-02) |
| SS_SPEC_REFERENCE_MISSING | info for planned; error for specced/later | Stage lacks required spec link | Planned info retained; matrix and error approved under D-2026-11-07-03 |
| SS_SPEC_DIRECTORY_MISSING | error, report | Explicit safe linked spec directory absent | Detection retained; approved as built (D-2026-11-07-03) |
| SS_STAGE_FILE_MISSING | error under the stage matrix | Stage-appropriate spec/plan/tasks missing | Detection retained; stage matrix and severity approved under D-2026-11-07-03 |
| SS_TASK_TARGET_MISSING | error if checked, info if unchecked | Explicit supported local link target absent | Retained |
| SS_REFERENCE_UNSUPPORTED | warning, report unless required parse fatal | Explicit reference syntax outside approved grammar | Approved as built (D-2026-11-11-02); never turns URL/prose into path |
| SS_PATH_UNSAFE | error; fatal for required source | Escape/traversal/symlink/mount/regular-file policy violation | Confinement retained; approved as built (D-2026-11-11-02) |
| SS_REFERENCE_SENSITIVE_SKIPPED | error, report; fatal if required source | Excluded sensitive reference, no content read | Exclusion/report retained; approved as built (D-2026-11-11-02) |
| SS_INPUT_LIMIT | error, fatal on required input/work | Approved bytes/files/path/nesting/items/edge/task budget exceeded | Explicit limit failure retained; fixed limits approved under DC-06 (D-2026-11-11-02) |
| SS_SOURCE_CHANGED | error, fatal | Detectable source change during input read | No-authoritative-snapshot policy (D-05B); approved as built (D-2026-11-11-02) |

No quarter diagnostic/code is defined. No finding is a mutation instruction or semantic code verdict. Codes emitted for a single issue must not inflate summary blocker count; e.g. a human-review dependency may also support a verified-contradiction finding, but the active row counts once.

## Location, redaction and ordering

Use root-relative normalized path and 1-based physical line where obtained; no absolute host paths, secret contents or raw OS stack traces. Missing source uses null, not an invented line. Related references preserve original dependency and successor/invalid chain for review, not rewrite. Diagnostic ordering follows JSON contract; every platform must produce equivalent safe meanings even when OS-specific errors differ.

## Recovery policy

Required root/roadmap/input encoding, unreliable item inventory, required unsafe/limit failures and observed source changes are fatal outcomes with no authoritative totals (D-04A, D-2026-11-11-02). Ordinary structural references can produce a complete supported-scope report plus findings. A skipped optional linked file is visible coverage and error, not equivalent to zero tasks or a verified outcome. Unsupported layout must not be ignored merely to keep exit 0. Any different approved recovery policy must adjust JSON completeness/summary semantics and fixtures explicitly.

## Approval and testing

The registry, format, matrix, safety matcher and JSON schema are approved as built (D-2026-11-07-03, D-2026-11-11-02). Fixtures assert code/severity/source and threshold behavior, not only free-form messages. Unknown future codes must not cause consumers to reinterpret exit 1 as a crash. Registry additions/meaning changes follow documented CLI/JSON compatibility policy. Tests have been run by the implementer, not independently verified; no security review of the matcher has been done.
