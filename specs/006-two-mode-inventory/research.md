# Research: Two-mode roadmap and feature inventory

Phase 0 for `specs/006-two-mode-inventory/plan.md`. Each entry: decision, rationale, alternatives. Code references are to the tree at the start of planning (`756ae75` plus PRs #48/#49, docs only).

## R1 — Where recognition is decided

- **Decision:** `parseRoadmap` (`internal/status/parse.go`) also reports whether it saw a visible `## Planned Specs` heading. `Evaluate` (`internal/status/run.go`) applies the rule after a successful parse: items > 0 → today's path; items = 0 and heading seen → explicitly empty; items = 0 and no heading → fatal `SS_ROADMAP_UNRECOGNIZED`.
- **Heading match:** a line outside fences and comments (the parser's existing `hidden[]` tracking) that is exactly `## Planned Specs`, starting at column 0, optionally followed by spaces or tabs. No setext form, no leading spaces, no closing `#`s, no other case or spelling. That is the pinned template's line (`testdata/upstream/speckit-roadmap/v0.1.0/templates/roadmap-template.md:52`) and nothing wider (C-03 narrow input).
- **Order:** parse fatals (encoding, layout, malformed, C-10 label rule) are decided first and keep their codes. Recognition only applies to a parse that succeeded, so a malformed file is never called unrecognized (spec US1 scenario 4).
- **Rationale:** the parser already knows which lines are visible; a second scan would duplicate the fence/comment logic and could disagree with it.
- **Alternatives:** a separate pre-scan of raw text (rejected: would count a heading inside a fence); accepting any heading level or case (rejected: wider than the pinned template).

## R2 — Impact on existing behavior (measured)

- **Probe:** a throwaway build with the rule (items = 0 and no `## Planned Specs` → fatal) ran the full suite. Only `TestEmptyRoadmapJSON` (`internal/status/status_test.go:13`) failed; it builds a 0-byte roadmap and expects exit 0. All conformance fixtures (19 roadmaps), CLI tests and extension tests passed unchanged.
- **Decision:** that test is rewritten into the new fixture set: a 0-byte roadmap is unrecognized (exit 2), and the same assertions about null denominators move to an explicitly empty roadmap (`## Planned Specs`, no items).
- **This repository's roadmap** has items, so its output is unchanged (SC-002).

## R3 — How "explicitly empty" is shown

- **Decision:** no new JSON field. Once unrecognized files are fatal, a complete report with `total_items: 0` can only mean an explicitly empty roadmap. Human text adds one line to the summary panel, shown only when there are zero items: `roadmap explicitly empty (Planned Specs heading, no items)`.
- **Rationale:** an additive field would change the JSON bytes of every existing report and fixture (FR-004, SC-002). The meaning is already determined by existing fields.
- **Spec change:** FR-002 and FR-006 are reworded accordingly (statement in human text; JSON meaning documented in the JSON contract), see the spec's revision note.
- **Alternatives:** `summary.roadmap_state` (rejected: breaks byte identity); a new `input_format` value (rejected: same).

## R4 — Messages

- `SS_ROADMAP_UNRECOGNIZED`: `Roadmap is not a recognized roadmap: no item headings and no "## Planned Specs" heading. See the supported template, or run spec inventory for a feature list.` For a file that is empty or whitespace only: `Roadmap file is empty. …` (same code, same tail). Source: the roadmap path, no line.
- `SS_ROADMAP_UNAVAILABLE` (missing only): message gains `Run spec inventory for a feature list without a roadmap.` The JSON `message` bytes for a missing roadmap change; that is the intended FR-005 change and goes in the release notes. Other unreadable cases (permission, unavailable) keep the old message.
- Neither message calls the project invalid Spec Kit (two-mode direction).

## R5 — Directory listing in safeio

- **Fact:** `internal/safeio` can stat and read files but cannot list a directory. Inventory needs one listing: the children of `specs/`.
- **Decision:** add `(*FS).ListDir(rel string) ([]DirEntry, error)`:
  - same `prepare` path as reads (normalize, sensitive check, no symlink components, within root);
  - `Lstat` the directory: symlink or other device → `unsafe`; not a directory → `unsafe`; absent → `notfound`;
  - open with `O_RDONLY|O_NOFOLLOW|O_DIRECTORY|O_CLOEXEC` through `syscall.Open` (no `os.OpenFile`, which the side-effect test forbids), `confirmOpened` on Linux;
  - read at most `MaxDirEntries + 1` names (new constant, 1024); more → `limit`;
  - for each name: if `Sensitive(rel/name)` → kind `sensitive`, not stat'ed; else `Lstat` → kind `dir`, `file`, `symlink` or `other`, plus the device check (`other-device` reported as `symlink`-like unsafe);
  - `Lstat` the directory again; identity changed → `changed`;
  - names returned in byte order; listing does not count against `MaxFiles`/`MaxTotalBytes` (no content is read).
  - Non-Linux/Darwin builds return `unsupported`, like `openNoFollow`.
- **Rationale:** keeps safeio the only filesystem boundary (CLAUDE.md, constitution Art. III) and one listing, not recursion ("No recursive inventory").
- **Alternatives:** `os.ReadDir` in the status package (rejected: bypasses the boundary); listing each feature folder too (rejected: unnecessary — the three files are stat'ed by name).

## R6 — Package layout

- **Decision:** inventory lives in `internal/status`: `inventory.go` (evaluate), `inventory_json.go` (JSON view), and new rendering functions in `style.go` (the only file allowed to import Lip Gloss). It reuses `resolveRoot`, the `.specify` check, `parseTasks`, `Diagnostic`, `diagLess`, `Percent`, the sanitizer and the exit-code rule.
- **Rationale:** the checkbox rules must be exactly the same (FR-011); sharing the function is the only way to guarantee it. The dependency allowlist stays as is.
- **Alternatives:** a new `internal/inventory` package (rejected: would need `parseTasks` exported or duplicated).

## R7 — Feature-folder names

- **Decision:** ASCII only. Timestamp form `^[0-9]{8}-[0-9]{6}-.+$` is checked first, then sequential `^[0-9]{3,}-.+$`. The JSON reports `name_form: "timestamp" | "sequential"`. Dates and times are not validated against a calendar. `20261101-name` (a date with no time) matches the sequential form, because eight digits are "three or more digits", so it is reported as a `sequential` feature folder.
- **Correction to the spec edge case:** the spec listed `20261101-name` as "not a feature folder". Under D-2026-11-12-08 (4) ("3 or more digits") it is a sequential feature folder. The spec edge case is corrected.
- Only directories are classified; a file named like a feature is ignored (FR-009). Symlinked children are never followed: they are listed with `SS_PATH_UNSAFE` and counted as skipped.

## R8 — Inventory diagnostics

| Code | Severity, class | When |
|---|---|---|
| `SS_SPECS_DIRECTORY_MISSING` (new) | info, report | `specs/` is absent |
| `SS_FOLDER_NOT_FEATURE` (new) | info, report | an immediate child directory of `specs/` does not match either name form |
| `SS_PATH_UNSAFE` (reused) | error, report | `specs/` is a symlink, not a directory, or on another device; a child is a symlink or other non-regular entry; `spec.md`/`plan.md`/`tasks.md` is a symlink or not a regular file |
| `SS_REFERENCE_SENSITIVE_SKIPPED` (reused) | error, report | a child or file name matched the sensitive-name matcher; not opened or stat'ed |
| `SS_INPUT_LIMIT`, `SS_SOURCE_CHANGED`, `SS_ENCODING_INVALID`, `SS_LAYOUT_UNSUPPORTED` (reused) | error, fatal | as in `spec status`; a tasks.md the narrow dialect cannot count reliably is fatal, citing the tasks file and line |

- Inventory does not check task link targets (`SS_TASK_TARGET_MISSING` is roadmap-mode only): it counts checkboxes and nothing else (FR-013). Documented in the contract.
- **Risk noted:** the sensitive-name matcher is over-inclusive (`safeio` package comment), so a feature folder named like `007-api-keys` is skipped with an error. Same behavior as a `Spec dir` with that name in roadmap mode. Recorded in README Limits.
- FIFO `tasks.md`: inventory reports `skipped` with `SS_PATH_UNSAFE` ("not a regular file"), unlike roadmap mode's known "reported as absent" defect, which is left to its own item.

## R9 — Invocation errors under `spec inventory --json`

- **Decision:** when the arguments name the `inventory` command (the existing `hasCommand` helper), invocation errors print an inventory-shaped fatal document (`inventory_schema_version: "1.0"`, `outcome: "fatal"`, `error` set). Otherwise the existing status-shaped invocation JSON is unchanged.
- Roadmap-only flags (`--roadmap`, `--status`, `--blocked`, `--human-review`, `--active`, `--historical`, `--explain`) are persistent flags on the root today; `spec inventory` rejects them when set with `SS_INVOCATION_INVALID` (`--roadmap is not used by inventory`).
- **Rationale:** a consumer of `spec inventory --json` parses one document shape.

## R10 — Bounds

- One listing (≤ 1024 children), three `StatRel` calls per feature folder, and at most one read per feature folder (`tasks.md`), all under the existing `MaxFiles` (1024), `MaxTotalBytes` (32 MiB), `MaxFileBytes` (2 MiB) and the existing checkbox total `maxTasks`. Exceeding any is `SS_INPUT_LIMIT`, fatal, no partial result (US3 scenario 3).

## R11 — Constitution text

- Art. II says rollout is not decided and the current CLI fails when the roadmap is missing. The second statement stays true. The first becomes stale once this ships. Under Governance item 4, a short amendment note (affected requirements, compatibility, fixtures, risks) is drafted in this feature and **applied only with owner approval**; `/speckit-constitution` is not run (CLAUDE.md).

## R12 — Release notes and docs

- v0.3.0 notes: "A roadmap file with no item headings and no `## Planned Specs` heading now stops with `SS_ROADMAP_UNRECOGNIZED` (exit 2) instead of reporting zero items (exit 0)" and the new `spec inventory`.
- Docs touched: `specs/001-spec-status/{diagnostics.md, json-contract.md, cli-contract.md, roadmap-format.md}`, the new `contracts/inventory-json.md` copied into `specs/001-spec-status/inventory-json.md` at implementation, README (Limits, missing-roadmap row, a short inventory section), `extensions/spec-status/reference/output.md` (three new codes; the extension test enforces coverage).
- README example output: unchanged (the fixture has items).
