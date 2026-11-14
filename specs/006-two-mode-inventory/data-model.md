# Data model: Two-mode roadmap and feature inventory

In-memory entities for `specs/006-two-mode-inventory/plan.md`. JSON field names are in [contracts/inventory-json.md](contracts/inventory-json.md); roadmap-mode changes are in [contracts/roadmap-recognition.md](contracts/roadmap-recognition.md).

## Roadmap recognition (roadmap mode)

| Result | Condition | Outcome |
|---|---|---|
| missing | roadmap path not found | fatal `SS_ROADMAP_UNAVAILABLE`, exit 2 (message suggests `spec inventory`) |
| malformed | a parse fatal (encoding, layout, item, C-10) | that fatal, exit 2, unchanged |
| with items | parse succeeded, ≥ 1 item | today's report, unchanged |
| explicitly empty | parse succeeded, 0 items, visible `## Planned Specs` | complete report, 0 items, null denominators, exit 0 |
| unrecognized | parse succeeded, 0 items, no visible `## Planned Specs` | fatal `SS_ROADMAP_UNRECOGNIZED`, exit 2 |

Parser output gains one boolean: `plannedSpecs` (a visible `## Planned Specs` line was seen). `Report` gains nothing; explicit empty is `Complete && Summary.TotalItems == 0`.

## DirEntry (safeio)

| Field | Type | Rule |
|---|---|---|
| Name | string | one path component, as stored on disk |
| Kind | enum | `dir`, `file`, `symlink`, `other`, `sensitive` (not stat'ed), `foreign` (other device) |
| ID | string | device:inode when known (for identity, not output) |

Returned in byte order of `Name`. At most `MaxDirEntries` (1024) entries; more is a `limit` error.

## InventoryReport

| Field | Type | Notes |
|---|---|---|
| Outcome | `report` \| `fatal` | |
| Complete | bool | false on fatal |
| FailOn | `info` \| `warning` \| `error` \| `none` | default `error` |
| SpecsDir | string | always `specs` (root-relative) |
| SpecsState | `present` \| `missing` \| `skipped` | skipped = unsafe or sensitive |
| RoadmapFile | `present` \| `absent` \| `skipped` | stat of `.specify/memory/roadmap.md`, never read |
| Folders | []Folder | byte order of name |
| Summary | InventorySummary | nil on fatal |
| Diagnostics | []Diagnostic | sorted with `diagLess` |
| Fatal | *Diagnostic | set on fatal |

## Folder

| Field | Type | Notes |
|---|---|---|
| Name | string | sanitized in human text |
| Path | string | `specs/<name>` |
| Feature | bool | name matches a Spec Kit form and entry is a real directory |
| NameForm | `sequential` \| `timestamp` \| null | null when not a feature folder |
| Skipped | bool | entry was a symlink, foreign, other, or sensitive; nothing under it examined |
| Spec, Plan | FileState | null when not a feature folder |
| Tasks | TaskState | null when not a feature folder |

Validation: `Feature` true implies `Skipped` false. A skipped entry is counted only in `skipped_entries`, never as a feature or other folder, and its file states are null.

## FileState (spec.md, plan.md)

| State | Meaning |
|---|---|
| `present` | regular file exists (stat only, not read) |
| `missing` | no entry with that name |
| `skipped` | symlink, not a regular file, other device, or sensitive; `Reason` is `unsafe` or `sensitive` |

## TaskState (tasks.md)

| State | Checked/Total | Meaning |
|---|---|---|
| `read` | numbers | read and counted with the roadmap-mode rules (`parseTasks`) |
| `missing` | null | no tasks.md |
| `skipped` | null | as FileState skipped, or a non-fatal read failure; `Reason` set |

`Percent` per folder is null when Total is 0.

## InventorySummary

| Field | Rule |
|---|---|
| FeatureFolders | count of `Feature` |
| OtherFolders | real directories that are not feature folders |
| SkippedEntries | children not examined (symlink, foreign, other, sensitive) |
| Tasks.Checked / Total | sums over `read` tasks files |
| Tasks.Percent | null when Total is 0 |
| Tasks.FilesRead / FilesMissing / FilesSkipped | over feature folders' tasks.md |

## State transitions

None. Every run is a fresh, read-only snapshot; nothing is stored.
