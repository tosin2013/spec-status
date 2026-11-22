package status

import (
	"sort"
	"strconv"
	"strings"

	"github.com/tosin2013/spec-status/internal/schema"
)

// JSON renders one UTF-8 document and a trailing newline. The CLI and
// pkg/status both produce it through toSchema and schema.Marshal.
func JSON(rep Report) []byte {
	return schema.Marshal(toSchema(rep))
}

// Schema returns the schema 1.0 view of a report, as `--json` prints it. It
// is the bridge used by pkg/status.
func Schema(rep Report) schema.Report {
	return toSchema(rep)
}

func toSchema(rep Report) schema.Report {
	out := schema.Report{
		SchemaVersion: schemaVersion,
		Outcome:       rep.Outcome,
		Complete:      rep.Complete,
		Source:        schema.Source{Roadmap: strPtr(rep.Roadmap), InputFormat: strPtr(rep.InputFormat)},
		Configuration: schema.Configuration{FailOn: or(rep.FailOn, "error"), DependencyPolicy: depPolicy, TaskScope: taskScope},
		Items:         []schema.Item{},
		Diagnostics:   []schema.Diagnostic{},
	}
	if rep.Complete && rep.Summary != nil {
		s := rep.Summary
		counts := s.StatusCounts
		if counts == nil {
			counts = map[string]int{}
		}
		out.Summary = &schema.Summary{
			TotalItems: s.TotalItems, StatusCounts: counts, ActiveItems: s.ActiveItems,
			VerifiedActiveItems: s.VerifiedActiveItems, VerifiedPercent: s.VerifiedPercent,
			SupersededItems: s.SupersededItems, InvalidSupersessionItems: s.InvalidSupersessionItems,
			BlockedCount: s.BlockedCount, HumanReviewCount: s.HumanReviewCount,
			Tasks: schema.TaskSummary{
				Checked: s.Tasks.Checked, Total: s.Tasks.Total, Percent: s.Tasks.Percent, Scope: s.Tasks.Scope,
				FilesRead: s.Tasks.FilesRead, FilesMissing: s.Tasks.FilesMissing, FilesSkipped: s.Tasks.FilesSkipped,
			},
		}
		for _, it := range rep.Items {
			out.Items = append(out.Items, itemToSchema(it))
		}
	}
	for _, d := range rep.Diagnostics {
		out.Diagnostics = append(out.Diagnostics, diagToSchema(d))
	}
	if rep.Outcome == "fatal" {
		d := Diagnostic{Code: "SS_ROADMAP_UNAVAILABLE", Severity: "error", Message: "Required input is unavailable."}
		if rep.Fatal != nil {
			d = *rep.Fatal
		} else if len(rep.Diagnostics) > 0 {
			d = rep.Diagnostics[0]
		}
		out.Error = &schema.Error{Code: d.Code, Message: d.Message, Source: locationOf(d.Path, d.Line)}
		out.Summary = nil
		out.Items = []schema.Item{}
		out.Complete = false
		out.Outcome = "fatal"
		out.Source.InputFormat = nil
	}
	return out
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

func itemToSchema(it Item) schema.Item {
	deps := make([]schema.Dependency, 0, len(it.Dependencies))
	for _, d := range it.Dependencies {
		deps = append(deps, schema.Dependency{
			ID: d.ID, Satisfied: d.Satisfied, State: d.State,
			ReplacementID: strPtr(d.ReplacementID), Source: schema.Ref{Path: d.Path, Line: d.Line},
		})
	}
	var checked, total *int
	if it.Tasks.HasNums {
		c, t := it.Tasks.Checked, it.Tasks.Total
		checked, total = &c, &t
	}
	var valid *bool
	if it.Supersession.Valid != nil {
		v := *it.Supersession.Valid
		valid = &v
	}
	var raw *string
	if it.Supersession.RawSet {
		r := it.Supersession.RawReference
		raw = &r
	}
	return schema.Item{
		ID: it.ID, Title: it.Title, Status: it.Status, StatusKnown: it.StatusKnown,
		Source:      schema.Ref{Path: it.Path, Line: it.Line},
		Active:      it.Active,
		Description: strPtr(it.Description), Outcome: strPtr(it.Outcome), Scope: strPtr(it.Scope),
		SpecDir: strPtr(it.SpecDir), SupersededBy: strPtr(it.SupersededBy),
		Supersession: schema.Supersession{Supplied: it.Supersession.Supplied, RawReference: raw, Valid: valid, TerminalID: strPtr(it.Supersession.TerminalID)},
		Dependencies: deps,
		Tasks: schema.Tasks{Checked: checked, Total: total, Percent: it.Tasks.Percent, Source: strPtr(it.Tasks.Source),
			State: or(it.Tasks.State, "not_linked"), Phases: phasesToSchema(it.Tasks)},
	}
}

// phasesToSchema returns nil (JSON null) unless the tasks file was read.
func phasesToSchema(t ItemTasks) []schema.Phase {
	if t.State != "read" {
		return nil
	}
	out := make([]schema.Phase, 0, len(t.Phases))
	for _, ph := range t.Phases {
		p := schema.Phase{Heading: ph.Heading, Line: ph.Line, Checked: ph.Checked, Total: ph.Total}
		if ph.Story != nil {
			v := *ph.Story
			p.Story = &v
		}
		out = append(out, p)
	}
	return out
}

func diagToSchema(d Diagnostic) schema.Diagnostic {
	rel := make([]schema.Ref, 0, len(d.Related))
	for _, r := range d.Related {
		rel = append(rel, schema.Ref{Path: r.Path, Line: r.Line})
	}
	return schema.Diagnostic{
		Code: d.Code, Severity: d.Severity, Message: d.Message,
		Source: locationOf(d.Path, d.Line), ItemID: strPtr(d.ItemID), Related: rel,
	}
}

// locationOf is the schema form of sourceRef: nil when both parts are
// absent, otherwise each absent part is null.
func locationOf(path string, line int) *schema.Location {
	if path == "" && line == 0 {
		return nil
	}
	loc := &schema.Location{}
	if path != "" {
		p := path
		loc.Path = &p
	}
	if line > 0 {
		l := line
		loc.Line = &l
	}
	return loc
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
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

type srcView struct {
	Path string `json:"path"`
	Line int    `json:"line"`
}

type errView struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Source  any    `json:"source"`
}
