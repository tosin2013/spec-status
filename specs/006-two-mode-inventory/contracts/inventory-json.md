# Contract: inventory JSON document, `inventory_schema_version` 1.0

New document for `spec inventory --json`. `spec status` JSON schema 1.0 is unchanged. Decision D-2026-11-12-08 (6); data model in [../data-model.md](../data-model.md).

## Rules

- One UTF-8 JSON object per run, newline-terminated, HTML escaping off, keys in the order below.
- Deterministic: no timestamps, no absolute paths, folders in byte order of name, diagnostics in the roadmap-mode order (`diagLess`).
- Percentages use the roadmap-mode number format (up to two decimals) or `null` when the denominator is 0.
- Additive changes may bump to 1.x; removing or changing a field's meaning needs 2.0.

## Report document

```json
{
  "inventory_schema_version": "1.0",
  "outcome": "report",
  "complete": true,
  "mode": "inventory",
  "source": {
    "specs_dir": "specs",
    "specs_state": "present",
    "roadmap_file": "absent"
  },
  "configuration": {
    "fail_on": "error",
    "task_scope": "feature_folder_tasks_md"
  },
  "summary": {
    "feature_folders": 2,
    "other_folders": 1,
    "skipped_entries": 0,
    "tasks": {
      "checked": 3,
      "total": 5,
      "percent": 60,
      "files_read": 1,
      "files_missing": 1,
      "files_skipped": 0
    }
  },
  "folders": [
    {
      "name": "001-alpha",
      "path": "specs/001-alpha",
      "feature": true,
      "name_form": "sequential",
      "skipped": false,
      "spec": {"state": "present", "reason": null},
      "plan": {"state": "present", "reason": null},
      "tasks": {"state": "read", "reason": null, "checked": 3, "total": 5, "percent": 60}
    },
    {
      "name": "20261101-093000-beta",
      "path": "specs/20261101-093000-beta",
      "feature": true,
      "name_form": "timestamp",
      "skipped": false,
      "spec": {"state": "present", "reason": null},
      "plan": {"state": "missing", "reason": null},
      "tasks": {"state": "missing", "reason": null, "checked": null, "total": null, "percent": null}
    },
    {
      "name": "notes",
      "path": "specs/notes",
      "feature": false,
      "name_form": null,
      "skipped": false,
      "spec": null,
      "plan": null,
      "tasks": null
    }
  ],
  "diagnostics": [
    {"code": "SS_FOLDER_NOT_FEATURE", "severity": "info", "message": "Folder name is not a Spec Kit feature folder name.", "source": {"path": "specs/notes", "line": null}, "item_id": null, "related": []}
  ],
  "error": null
}
```

## Fields

| Field | Type | Values / meaning |
|---|---|---|
| `inventory_schema_version` | string | `"1.0"` |
| `outcome` | string | `report`, `fatal` |
| `complete` | bool | `false` on fatal |
| `mode` | string | `"inventory"` |
| `source.specs_dir` | string | `"specs"` |
| `source.specs_state` | string | `present`, `missing`, `skipped` |
| `source.roadmap_file` | string \| null | `present`, `absent`, `skipped`; `null` on fatal before the root was usable |
| `configuration.fail_on` | string | effective threshold |
| `configuration.task_scope` | string | `"feature_folder_tasks_md"` |
| `summary` | object \| null | `null` on fatal |
| `summary.feature_folders` / `other_folders` / `skipped_entries` | int | see data model |
| `summary.tasks.checked` / `total` | int | sums over read tasks files |
| `summary.tasks.percent` | number \| null | null when `total` is 0 |
| `summary.tasks.files_read` / `files_missing` / `files_skipped` | int | over feature folders |
| `folders[]` | array | `[]` on fatal |
| `folders[].name` / `path` | string | raw name; root-relative slash path |
| `folders[].feature` | bool | |
| `folders[].name_form` | string \| null | `sequential`, `timestamp`; null when not a feature |
| `folders[].skipped` | bool | entry not examined |
| `folders[].spec` / `plan` | object \| null | `{state: present\|missing\|skipped, reason: null\|unsafe\|sensitive}`; null when not a feature folder |
| `folders[].tasks` | object \| null | `{state: read\|missing\|skipped, reason, checked, total, percent}`; numbers null unless `read` |
| `diagnostics[]` | array | same object shape as `spec status` schema 1.0 diagnostics; `item_id` is always null |
| `error` | object \| null | `{code, message, source}` on fatal, as in `spec status` |

## Fatal document

```json
{"inventory_schema_version":"1.0","outcome":"fatal","complete":false,"mode":"inventory","source":{"specs_dir":"specs","specs_state":null,"roadmap_file":null},"configuration":{"fail_on":"error","task_scope":"feature_folder_tasks_md"},"summary":null,"folders":[],"diagnostics":[{"code":"SS_INPUT_LIMIT","severity":"error","message":"Input budget exceeded.","source":{"path":"specs","line":null},"item_id":null,"related":[]}],"error":{"code":"SS_INPUT_LIMIT","message":"Input budget exceeded.","source":{"path":"specs","line":null}}}
```

On fatal, `source.specs_state` may be `null`. Invocation errors use this shape with `SS_INVOCATION_INVALID` and the requested `fail_on` when it parsed, else `"error"`.

## Diagnostics used by inventory

| Code | Severity / class | Source |
|---|---|---|
| `SS_SPECS_DIRECTORY_MISSING` (new) | info, report | `specs` |
| `SS_FOLDER_NOT_FEATURE` (new) | info, report | `specs/<name>` |
| `SS_PATH_UNSAFE` | error, report | the entry or file path |
| `SS_REFERENCE_SENSITIVE_SKIPPED` | error, report | the entry or file path |
| `SS_ROOT_UNAVAILABLE`, `SS_INPUT_LIMIT`, `SS_SOURCE_CHANGED`, `SS_ENCODING_INVALID`, `SS_LAYOUT_UNSUPPORTED`, `SS_INVOCATION_INVALID` | error, fatal | as in `spec status`; tasks.md fatals cite the tasks file and line |

Not emitted by inventory: roadmap, dependency, supersession, stage and task-link codes.
