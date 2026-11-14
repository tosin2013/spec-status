# Contract: inventory human text

Same plain-surface rules as `spec status` (CLAUDE.md "Plain surfaces"): color only on a TTY with `NO_COLOR` unset and `TERM != dumb`, or `--color=always`; width from `COLUMNS`, else 80; untrusted names sanitized of terminal escapes; Lip Gloss only in `style.go`.

## Layout (80 columns, color off)

```
╭──────────────────────────────────────────────────────────────────────────────╮
│ Spec Inventory  development build                                            │
│ mode                inventory (roadmap not read; no status inferred)         │
│ checked tasks       3/5 (60%)  [#######-----]                                │
│ folders 3 feature 2 other 1 skipped 0                                        │
│ task files read 1 missing 1 skipped 0                                        │
│ source specs                                                                 │
╰──────────────────────────────────────────────────────────────────────────────╯

Folder                 spec     plan     tasks
001-alpha              present  present  3/5
20261101-093000-beta   present  missing  missing
notes                  not a feature folder

Checked boxes are not verification.

Diagnostics
info  SS_FOLDER_NOT_FEATURE  specs/notes  Folder name is not a Spec Kit feature folder name.
```

- The exact box drawing, column widths and narrow-width (stacked) layout follow the existing `spec status` renderer; the README sample for inventory is generated from a fixture and checked by a test (the same mechanism item 011 adds).
- `specs/` missing: the panel shows `source specs (missing)` and no folder table; the diagnostic is listed.
- `roadmap_file: present` adds one line after the table: `A roadmap exists; spec status reports it.`
- Fatal: `spec: CODE: message` on stderr only.
