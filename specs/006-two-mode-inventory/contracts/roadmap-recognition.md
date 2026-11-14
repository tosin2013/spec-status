# Contract: roadmap recognition (`spec status`)

Amends `specs/001-spec-status/roadmap-format.md`, `diagnostics.md` and `json-contract.md` at implementation. Decision D-2026-11-12-08 (3); research R1–R4.

## Rule

After a successful parse of the selected roadmap:

1. One or more supported item headings → the report is exactly as before this feature (bytes, diagnostics, exit code).
2. Zero items and a **visible** line that is exactly `## Planned Specs` (column 0, optional trailing spaces or tabs, not inside a fence or HTML comment) → **explicitly empty**: `outcome: "report"`, `complete: true`, `total_items: 0`, `verified_percent: null`, `tasks.percent: null`, exit 0 (no findings).
3. Otherwise → fatal `SS_ROADMAP_UNRECOGNIZED`, exit 2, under every `--fail-on` value.

Parse fatals are decided before this rule and keep their codes. The rule applies to the default path and to `--roadmap PATH` alike.

## New diagnostic

| Code | Severity / class | Trigger and source |
|---|---|---|
| `SS_ROADMAP_UNRECOGNIZED` | error, fatal | Readable roadmap with no supported item heading and no visible `## Planned Specs` heading; source is the roadmap path, line null |

Messages (exact):

- non-empty file: `Roadmap is not a recognized roadmap: no item headings and no "## Planned Specs" heading. See the supported template, or run spec inventory for a feature list.`
- empty or whitespace-only file: `Roadmap file is empty. See the supported template, or run spec inventory for a feature list.`

## Changed message

`SS_ROADMAP_UNAVAILABLE` when the roadmap path does not exist (`safeio` `notfound`) only:
`Required roadmap is unavailable. Run spec inventory for a feature list without a roadmap.`
Other unavailable causes keep `Required roadmap is unavailable.`

## JSON (schema 1.0, unchanged shape)

- No field is added or removed. A fatal uses the existing fatal document (`summary: null`, `items: []`, `error` set).
- Documented meaning: a complete report with `summary.total_items == 0` is an explicitly empty supported roadmap.

## Human text

- Fatal: `spec: SS_ROADMAP_UNRECOGNIZED: <message>` on stderr, plain, exit 2 (existing fatal text).
- Explicitly empty: the summary panel adds the line `roadmap explicitly empty (Planned Specs heading, no items)`. Nothing else changes. The line is not shown when there are items.

## Fixtures (hand-counted, `testdata/conformance/` or unit tests)

| Fixture | Expected |
|---|---|
| 0-byte roadmap | fatal `SS_ROADMAP_UNRECOGNIZED` (empty message), exit 2 |
| whitespace-only roadmap | same |
| narrative roadmap (headings, prose, no items) | fatal `SS_ROADMAP_UNRECOGNIZED`, exit 2 |
| `## Planned Specs` only inside a fence | fatal `SS_ROADMAP_UNRECOGNIZED` |
| `### Planned Specs` / `## planned specs` / `  ## Planned Specs` | fatal `SS_ROADMAP_UNRECOGNIZED` |
| pinned template sections with `## Planned Specs` and only `###` grouping headings | explicitly empty, exit 0 |
| items without `## Planned Specs` | unchanged report |
| narrative file with a field label outside items | existing parse fatal (C-10), not unrecognized |
| missing roadmap | `SS_ROADMAP_UNAVAILABLE` with the new message, exit 2 |
| unrecognized with `--fail-on none` | exit 2 |
