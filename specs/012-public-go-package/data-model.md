# Data model: Public Go package

The types live in `internal/schema` and are re-exported as aliases from `pkg/status`. Each field maps 1:1 to schema 1.0 (`specs/001-spec-status/json-contract.md`), in the same order as the JSON.

| Type | Fields (Go type → JSON) |
|---|---|
| Report | SchemaVersion string; Outcome string; Complete bool; Source Source; Configuration Configuration; Summary *Summary (null on fatal); Items []Item (always array); Diagnostics []Diagnostic (always array); Error *Error |
| Source | Roadmap *string; InputFormat *string |
| Configuration | FailOn, DependencyPolicy, TaskScope string |
| Summary | TotalItems int; StatusCounts map[string]int (always object); ActiveItems, VerifiedActiveItems int; VerifiedPercent Percent; SupersededItems, InvalidSupersessionItems, BlockedCount, HumanReviewCount int; Tasks TaskSummary |
| TaskSummary | Checked, Total int; Percent Percent; Scope string; FilesRead, FilesMissing, FilesSkipped int |
| Item | ID, Title, Status string; StatusKnown bool; Source Ref; Active bool; Description, Outcome, Scope, SpecDir, SupersededBy *string; Supersession Supersession; Dependencies []Dependency; Tasks Tasks |
| Ref | Path string; Line int |
| Supersession | Supplied bool; RawReference *string; Valid *bool; TerminalID *string |
| Dependency | ID string; Satisfied bool; State string; ReplacementID *string; Source Ref |
| Tasks | Checked, Total *int; Percent Percent; Source *string; State string; Phases []Phase (nil = null: file not read; empty = no phases) |
| Phase | Heading string; Line int; Story *int; Checked, Total int |
| Diagnostic | Code, Severity, Message string; Source *Location; ItemID *string; Related []Ref (always array) |
| Error | Code, Message string; Source *Location |
| Location | Path *string; Line *int (either may be null) |
| Percent | unexported (ok, hundredths); Valid() bool; Hundredths() int; Text() string; MarshalJSON / UnmarshalJSON (bare number or null) |

`Options` (pkg/status only): Roadmap string (as `--roadmap`), FailOn string ("" = "error"; else error, warning, info or none).

Validation: `Load` rejects `root == ""` and unknown `FailOn` with an `error`. Everything else, fatal outcomes included, is a `Report`.
