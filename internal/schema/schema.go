// Package schema holds the `spec status --json` schema 1.0 report types,
// shared by the CLI (internal/status) and the public package pkg/status, so
// the two cannot drift. Standard library only.
//
// Field order and JSON tags are the wire format. A pointer field is null in
// the JSON when nil. Tasks.Phases is null when nil (the tasks file was not
// read) and [] when empty.
package schema

// Version is the JSON schema version of Report.
const Version = "1.0"

// Report is one `spec status --json` document.
type Report struct {
	SchemaVersion string        `json:"schema_version"`
	Outcome       string        `json:"outcome"`
	Complete      bool          `json:"complete"`
	Source        Source        `json:"source"`
	Configuration Configuration `json:"configuration"`
	Summary       *Summary      `json:"summary"`
	Items         []Item        `json:"items"`
	Diagnostics   []Diagnostic  `json:"diagnostics"`
	Error         *Error        `json:"error"`
}

// Source names the roadmap that was read and its recognized input format.
type Source struct {
	Roadmap     *string `json:"roadmap"`
	InputFormat *string `json:"input_format"`
}

// Configuration is the threshold and the fixed policies of the run.
type Configuration struct {
	FailOn           string `json:"fail_on"`
	DependencyPolicy string `json:"dependency_policy"`
	TaskScope        string `json:"task_scope"`
}

// Summary describes the whole roadmap. It is nil for a fatal report.
type Summary struct {
	TotalItems               int            `json:"total_items"`
	StatusCounts             map[string]int `json:"status_counts"`
	ActiveItems              int            `json:"active_items"`
	VerifiedActiveItems      int            `json:"verified_active_items"`
	VerifiedPercent          Percent        `json:"verified_percent"`
	SupersededItems          int            `json:"superseded_items"`
	InvalidSupersessionItems int            `json:"invalid_supersession_items"`
	BlockedCount             int            `json:"blocked_count"`
	HumanReviewCount         int            `json:"human_review_count"`
	Tasks                    TaskSummary    `json:"tasks"`
}

// TaskSummary counts checkboxes across unique linked task files.
type TaskSummary struct {
	Checked      int     `json:"checked"`
	Total        int     `json:"total"`
	Percent      Percent `json:"percent"`
	Scope        string  `json:"scope"`
	FilesRead    int     `json:"files_read"`
	FilesMissing int     `json:"files_missing"`
	FilesSkipped int     `json:"files_skipped"`
}

// Item is one roadmap item, in document order.
type Item struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Status       string       `json:"status"`
	StatusKnown  bool         `json:"status_known"`
	Source       Ref          `json:"source"`
	Active       bool         `json:"active"`
	Description  *string      `json:"description"`
	Outcome      *string      `json:"outcome"`
	Scope        *string      `json:"scope"`
	SpecDir      *string      `json:"spec_dir"`
	SupersededBy *string      `json:"superseded_by"`
	Supersession Supersession `json:"supersession"`
	Dependencies []Dependency `json:"dependencies"`
	Tasks        Tasks        `json:"tasks"`
}

// Ref is a file and 1-based line.
type Ref struct {
	Path string `json:"path"`
	Line int    `json:"line"`
}

// Location is a file and line where either part may be absent (null).
type Location struct {
	Path *string `json:"path"`
	Line *int    `json:"line"`
}

// Supersession is an item's raw superseded-by reference and its evaluation.
type Supersession struct {
	Supplied     bool    `json:"supplied"`
	RawReference *string `json:"raw_reference"`
	Valid        *bool   `json:"valid"`
	TerminalID   *string `json:"terminal_id"`
}

// Dependency is one Depends on entry and its state.
type Dependency struct {
	ID            string  `json:"id"`
	Satisfied     bool    `json:"satisfied"`
	State         string  `json:"state"`
	ReplacementID *string `json:"replacement_id"`
	Source        Ref     `json:"source"`
}

// Tasks is an item's linked tasks file. Checked and Total are nil unless the
// file was read. Phases is nil (null) unless the file was read.
type Tasks struct {
	Checked *int    `json:"checked"`
	Total   *int    `json:"total"`
	Percent Percent `json:"percent"`
	Source  *string `json:"source"`
	State   string  `json:"state"`
	Phases  []Phase `json:"phases"`
}

// Phase is the recorded checkboxes under one `## Phase …` heading. It says
// nothing about status or verification.
type Phase struct {
	Heading string `json:"heading"`
	Line    int    `json:"line"`
	Story   *int   `json:"story"`
	Checked int    `json:"checked"`
	Total   int    `json:"total"`
}

// Diagnostic is one finding.
type Diagnostic struct {
	Code     string    `json:"code"`
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	Source   *Location `json:"source"`
	ItemID   *string   `json:"item_id"`
	Related  []Ref     `json:"related"`
}

// Error is the fatal error of a report whose Outcome is "fatal".
type Error struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Source  *Location `json:"source"`
}
