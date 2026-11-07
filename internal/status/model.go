// Package status evaluates a Spec Kit roadmap.
//
// Behavior that originated in PROFILE-R1 is an implementation proposal.
// It is not an owner-approved decision. See specs/001-spec-status/decisions.md.
package status

import "strconv"

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

// Percent is a half-up percentage to two decimal places, or null when the denominator is zero.
type Percent struct {
	ok         bool
	hundredths int
}

func percentOf(num, den int) Percent {
	if den <= 0 || num < 0 {
		return Percent{}
	}
	scaled := num * 10000
	q := scaled / den
	r := scaled % den
	if r*2 >= den {
		q++
	}
	return Percent{ok: true, hundredths: q}
}

func (p Percent) Text() string {
	if !p.ok {
		return "N/A"
	}
	whole := p.hundredths / 100
	frac := p.hundredths % 100
	if frac == 0 {
		return strconv.Itoa(whole)
	}
	if frac%10 == 0 {
		return strconv.Itoa(whole) + "." + strconv.Itoa(frac/10)
	}
	if frac < 10 {
		return strconv.Itoa(whole) + ".0" + strconv.Itoa(frac)
	}
	return strconv.Itoa(whole) + "." + strconv.Itoa(frac/10) + strconv.Itoa(frac%10)
}

func (p Percent) MarshalJSON() ([]byte, error) {
	if !p.ok {
		return []byte("null"), nil
	}
	return []byte(p.Text()), nil
}

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
}

type Supersession struct {
	Supplied     bool
	RawReference string
	RawSet       bool
	Parsed       string
	Valid        *bool
	TerminalID   string
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
