package status

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

// JSON renders one UTF-8 document and a trailing newline.
func JSON(rep Report) []byte {
	view := docView{
		SchemaVersion: schemaVersion,
		Outcome:       rep.Outcome,
		Complete:      rep.Complete,
		Source: sourceView{
			Roadmap:     nullString(rep.Roadmap),
			InputFormat: nullString(rep.InputFormat),
		},
		Configuration: configView{FailOn: or(rep.FailOn, "error"), DependencyPolicy: depPolicy, TaskScope: taskScope},
		Items:         []itemView{},
		Diagnostics:   []diagView{},
		Error:         nil,
	}
	if rep.Complete && rep.Summary != nil {
		s := rep.Summary
		view.Summary = &summaryView{
			TotalItems: s.TotalItems, StatusCounts: s.StatusCounts, ActiveItems: s.ActiveItems,
			VerifiedActiveItems: s.VerifiedActiveItems, VerifiedPercent: s.VerifiedPercent,
			SupersededItems: s.SupersededItems, InvalidSupersessionItems: s.InvalidSupersessionItems,
			BlockedCount: s.BlockedCount, HumanReviewCount: s.HumanReviewCount,
			Tasks: taskSummaryView{
				Checked: s.Tasks.Checked, Total: s.Tasks.Total, Percent: s.Tasks.Percent, Scope: s.Tasks.Scope,
				FilesRead: s.Tasks.FilesRead, FilesMissing: s.Tasks.FilesMissing, FilesSkipped: s.Tasks.FilesSkipped,
			},
		}
		if view.Summary.StatusCounts == nil {
			view.Summary.StatusCounts = map[string]int{}
		}
		for _, it := range rep.Items {
			view.Items = append(view.Items, itemToView(it))
		}
	}
	for _, d := range rep.Diagnostics {
		view.Diagnostics = append(view.Diagnostics, diagToView(d))
	}
	if rep.Outcome == "fatal" {
		d := Diagnostic{Code: "SS_ROADMAP_UNAVAILABLE", Severity: "error", Message: "Required input is unavailable."}
		if rep.Fatal != nil {
			d = *rep.Fatal
		} else if len(rep.Diagnostics) > 0 {
			d = rep.Diagnostics[0]
		}
		view.Error = &errView{Code: d.Code, Message: d.Message, Source: sourceRef(d.Path, d.Line)}
		view.Summary = nil
		view.Items = []itemView{}
		view.Complete = false
		view.Outcome = "fatal"
		view.Source.InputFormat = nil
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(view); err != nil {
		return []byte("{\"schema_version\":\"1.0\",\"outcome\":\"fatal\",\"complete\":false,\"source\":{\"roadmap\":null,\"input_format\":null},\"configuration\":{\"fail_on\":\"error\",\"dependency_policy\":\"verified\",\"task_scope\":\"unique_linked_task_files\"},\"summary\":null,\"items\":[],\"diagnostics\":[{\"code\":\"SS_INVOCATION_INVALID\",\"severity\":\"error\",\"message\":\"Report could not be serialized.\",\"source\":null,\"item_id\":null,\"related\":[]}],\"error\":{\"code\":\"SS_INVOCATION_INVALID\",\"message\":\"Report could not be serialized.\",\"source\":null}}\n")
	}
	return buf.Bytes()
}

// Text renders the human report. Null quantities are N/A.
// view.Color adds restrained ANSI labels. It does not change the facts.
// Fatal text stays plain so error output is stable.
func Text(rep Report, view View) string {
	if rep.Outcome == "fatal" {
		return fatalText(rep)
	}
	return renderHuman(rep, view)
}

func fatalText(rep Report) string {
	d := Diagnostic{Code: "SS_ROADMAP_UNAVAILABLE", Message: "Required input is unavailable."}
	if rep.Fatal != nil {
		d = *rep.Fatal
	} else if len(rep.Diagnostics) > 0 {
		d = rep.Diagnostics[0]
	}
	return "spec: " + d.Code + ": " + d.Message + "\n"
}

func formatCounts(m map[string]int) string {
	if len(m) == 0 {
		return "none"
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+" "+strconv.Itoa(m[k]))
	}
	return strings.Join(parts, ", ")
}

func itemToView(it Item) itemView {
	deps := make([]depView, 0, len(it.Dependencies))
	for _, d := range it.Dependencies {
		deps = append(deps, depView{
			ID: d.ID, Satisfied: d.Satisfied, State: d.State,
			ReplacementID: nullString(d.ReplacementID), Source: srcView{Path: d.Path, Line: d.Line},
		})
	}
	var checked, total *int
	if it.Tasks.HasNums {
		c, t := it.Tasks.Checked, it.Tasks.Total
		checked, total = &c, &t
	}
	valid := any(nil)
	if it.Supersession.Valid != nil {
		v := *it.Supersession.Valid
		valid = v
	}
	raw := any(nil)
	if it.Supersession.RawSet {
		raw = it.Supersession.RawReference
	}
	return itemView{
		ID: it.ID, Title: it.Title, Status: it.Status, StatusKnown: it.StatusKnown,
		Source:      srcView{Path: it.Path, Line: it.Line},
		Active:      it.Active,
		Description: nullString(it.Description), Outcome: nullString(it.Outcome), Scope: nullString(it.Scope),
		SpecDir: nullString(it.SpecDir), SupersededBy: nullString(it.SupersededBy),
		Supersession: superView{Supplied: it.Supersession.Supplied, RawReference: raw, Valid: valid, TerminalID: nullString(it.Supersession.TerminalID)},
		Dependencies: deps,
		Tasks:        itemTaskView{Checked: checked, Total: total, Percent: it.Tasks.Percent, Source: nullString(it.Tasks.Source), State: or(it.Tasks.State, "not_linked")},
	}
}

func diagToView(d Diagnostic) diagView {
	rel := make([]srcView, 0, len(d.Related))
	for _, r := range d.Related {
		rel = append(rel, srcView{Path: r.Path, Line: r.Line})
	}
	return diagView{
		Code: d.Code, Severity: d.Severity, Message: d.Message,
		Source: sourceRef(d.Path, d.Line), ItemID: nullString(d.ItemID), Related: rel,
	}
}

func sourceRef(path string, line int) any {
	if path == "" && line == 0 {
		return nil
	}
	var lineVal any
	if line > 0 {
		lineVal = line
	}
	var pathVal any
	if path != "" {
		pathVal = path
	}
	return srcAny{Path: pathVal, Line: lineVal}
}

type srcAny struct {
	Path any `json:"path"`
	Line any `json:"line"`
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func or(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

type docView struct {
	SchemaVersion string       `json:"schema_version"`
	Outcome       string       `json:"outcome"`
	Complete      bool         `json:"complete"`
	Source        sourceView   `json:"source"`
	Configuration configView   `json:"configuration"`
	Summary       *summaryView `json:"summary"`
	Items         []itemView   `json:"items"`
	Diagnostics   []diagView   `json:"diagnostics"`
	Error         *errView     `json:"error"`
}

type sourceView struct {
	Roadmap     any `json:"roadmap"`
	InputFormat any `json:"input_format"`
}

type configView struct {
	FailOn           string `json:"fail_on"`
	DependencyPolicy string `json:"dependency_policy"`
	TaskScope        string `json:"task_scope"`
}

type summaryView struct {
	TotalItems               int             `json:"total_items"`
	StatusCounts             map[string]int  `json:"status_counts"`
	ActiveItems              int             `json:"active_items"`
	VerifiedActiveItems      int             `json:"verified_active_items"`
	VerifiedPercent          Percent         `json:"verified_percent"`
	SupersededItems          int             `json:"superseded_items"`
	InvalidSupersessionItems int             `json:"invalid_supersession_items"`
	BlockedCount             int             `json:"blocked_count"`
	HumanReviewCount         int             `json:"human_review_count"`
	Tasks                    taskSummaryView `json:"tasks"`
}

type taskSummaryView struct {
	Checked      int     `json:"checked"`
	Total        int     `json:"total"`
	Percent      Percent `json:"percent"`
	Scope        string  `json:"scope"`
	FilesRead    int     `json:"files_read"`
	FilesMissing int     `json:"files_missing"`
	FilesSkipped int     `json:"files_skipped"`
}

type itemView struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Status       string       `json:"status"`
	StatusKnown  bool         `json:"status_known"`
	Source       srcView      `json:"source"`
	Active       bool         `json:"active"`
	Description  any          `json:"description"`
	Outcome      any          `json:"outcome"`
	Scope        any          `json:"scope"`
	SpecDir      any          `json:"spec_dir"`
	SupersededBy any          `json:"superseded_by"`
	Supersession superView    `json:"supersession"`
	Dependencies []depView    `json:"dependencies"`
	Tasks        itemTaskView `json:"tasks"`
}

type srcView struct {
	Path string `json:"path"`
	Line int    `json:"line"`
}

type superView struct {
	Supplied     bool `json:"supplied"`
	RawReference any  `json:"raw_reference"`
	Valid        any  `json:"valid"`
	TerminalID   any  `json:"terminal_id"`
}

type depView struct {
	ID            string  `json:"id"`
	Satisfied     bool    `json:"satisfied"`
	State         string  `json:"state"`
	ReplacementID any     `json:"replacement_id"`
	Source        srcView `json:"source"`
}

type itemTaskView struct {
	Checked *int    `json:"checked"`
	Total   *int    `json:"total"`
	Percent Percent `json:"percent"`
	Source  any     `json:"source"`
	State   string  `json:"state"`
}

type diagView struct {
	Code     string    `json:"code"`
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	Source   any       `json:"source"`
	ItemID   any       `json:"item_id"`
	Related  []srcView `json:"related"`
}

type errView struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Source  any    `json:"source"`
}
