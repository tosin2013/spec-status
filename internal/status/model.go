// Package status evaluates a Spec Kit roadmap.
//
// Behavior that originated in PROFILE-R1 is an implementation proposal.
// It is not an owner-approved decision. See specs/001-spec-status/decisions.md.
package status

import "github.com/tosin2013/spec-status/internal/schema"

const (
	schemaVersion = "1.0"
	inputFormat   = "speckit-roadmap-v0.1.1+spec-status-supersession-draft"
	taskScope     = "unique_linked_task_files"
	depPolicy     = "verified"

	maxItems   = 1024
	maxEdges   = 16384
	maxTasks   = 100000
	maxNesting = 32
)

// Known lifecycle tokens. Comparison is exact and case-sensitive.
var knownStatus = map[string]bool{
	"undecided": true, "needs-info": true, "planned": true, "specced": true,
	"in-progress": true, "implemented": true, "verified": true, "deferred": true, "abandoned": true,
}

// Percent is a half-up percentage to two decimal places, or null when the
// denominator is zero. It is the schema type, so text and JSON share it.
type Percent = schema.Percent

func percentOf(num, den int) Percent { return schema.PercentOf(num, den) }

// Report is the single result used for text and JSON.
type Report struct {
	Outcome     string
	Complete    bool
	Roadmap     string
	InputFormat string
	FailOn      string
	Summary     *Summary
	Items       []Item
	Diagnostics []Diagnostic
	Fatal       *Diagnostic
}

type Summary struct {
	TotalItems               int
	StatusCounts             map[string]int
	ActiveItems              int
	VerifiedActiveItems      int
	VerifiedPercent          Percent
	SupersededItems          int
	InvalidSupersessionItems int
	BlockedCount             int
	HumanReviewCount         int
	Tasks                    TaskSummary
}

type TaskSummary struct {
	Checked      int
	Total        int
	Percent      Percent
	Scope        string
	FilesRead    int
	FilesMissing int
	FilesSkipped int
}

type Item struct {
	ID           string
	Title        string
	Status       string
	StatusKnown  bool
	Path         string
	Line         int
	Active       bool
	Description  string
	Outcome      string
	Scope        string
	SpecDir      string
	SupersededBy string
	Supersession Supersession
	Dependencies []Dependency
	Tasks        ItemTasks
	// TasksSharedWith is the ID of an earlier item that read the same tasks
	// file, which the summary counts once. Human text only; not in JSON
	// (feature 008).
	TasksSharedWith string
}

type Supersession struct {
	Supplied     bool
	RawReference string
	RawSet       bool
	Parsed       string
	Valid        *bool
	TerminalID   string
	Line         int
}

type Dependency struct {
	ID            string
	Satisfied     bool
	State         string
	ReplacementID string
	Path          string
	Line          int
}

type ItemTasks struct {
	Checked int
	Total   int
	Percent Percent
	HasNums bool
	Source  string
	State   string
	Checks  []TaskCheck
	// Phases lists the tasks file's `## Phase …` sections (feature 011). Set
	// only when the file was read.
	Phases []Phase
}

// Phase is the recorded checkboxes under one `## Phase …` heading of a tasks
// file (feature 011). It says nothing about status or verification.
type Phase struct {
	Heading string
	Line    int
	Story   *int
	Checked int
	Total   int
}

// TaskCheck is text-only evidence. It is not part of the JSON schema.
type TaskCheck struct {
	Line    int
	Checked bool
}

type Diagnostic struct {
	Code     string
	Severity string
	Message  string
	Path     string
	Line     int
	ItemID   string
	Related  []Ref
}

type Ref struct {
	Path string
	Line int
}

func fatalReport(code, message, roadmap, failOn string) Report {
	d := Diagnostic{Code: code, Severity: "error", Message: message, Path: roadmap}
	if roadmap == "" {
		d.Path = ""
	}
	return Report{
		Outcome:     "fatal",
		Complete:    false,
		Roadmap:     roadmap,
		FailOn:      failOn,
		Items:       []Item{},
		Diagnostics: []Diagnostic{d},
		Fatal:       &d,
	}
}
