package status

import (
	"bytes"
	"encoding/json"

	"github.com/tosin2013/spec-status/internal/schema"
)

// Inventory JSON document, inventory_schema_version 1.0
// (specs/006-two-mode-inventory/contracts/inventory-json.md). Struct field
// order is the key order. spec status JSON (schema 1.0) is separate.

const inventorySchemaVersion = "1.0"

type invDocView struct {
	SchemaVersion string              `json:"inventory_schema_version"`
	Outcome       string              `json:"outcome"`
	Complete      bool                `json:"complete"`
	Mode          string              `json:"mode"`
	Source        invSourceView       `json:"source"`
	Configuration invConfigView       `json:"configuration"`
	Summary       *invSummaryView     `json:"summary"`
	Folders       []invFolderView     `json:"folders"`
	Diagnostics   []schema.Diagnostic `json:"diagnostics"`
	Error         *schema.Error       `json:"error"`
}

type invSourceView struct {
	SpecsDir    string  `json:"specs_dir"`
	SpecsState  *string `json:"specs_state"`
	RoadmapFile *string `json:"roadmap_file"`
}

type invConfigView struct {
	FailOn    string `json:"fail_on"`
	TaskScope string `json:"task_scope"`
}

type invSummaryView struct {
	FeatureFolders int            `json:"feature_folders"`
	OtherFolders   int            `json:"other_folders"`
	SkippedEntries int            `json:"skipped_entries"`
	Tasks          invTaskSumView `json:"tasks"`
}

type invTaskSumView struct {
	Checked      int     `json:"checked"`
	Total        int     `json:"total"`
	Percent      Percent `json:"percent"`
	FilesRead    int     `json:"files_read"`
	FilesMissing int     `json:"files_missing"`
	FilesSkipped int     `json:"files_skipped"`
}

type invFolderView struct {
	Name     string        `json:"name"`
	Path     string        `json:"path"`
	Feature  bool          `json:"feature"`
	NameForm *string       `json:"name_form"`
	Skipped  bool          `json:"skipped"`
	Spec     *invFileView  `json:"spec"`
	Plan     *invFileView  `json:"plan"`
	Tasks    *invTasksView `json:"tasks"`
}

type invFileView struct {
	State  string  `json:"state"`
	Reason *string `json:"reason"`
}

type invTasksView struct {
	State   string  `json:"state"`
	Reason  *string `json:"reason"`
	Checked any     `json:"checked"`
	Total   any     `json:"total"`
	Percent any     `json:"percent"`
}

// InventoryJSON renders one inventory document, newline-terminated.
func InventoryJSON(rep InventoryReport) []byte {
	view := invDocView{
		SchemaVersion: inventorySchemaVersion,
		Outcome:       rep.Outcome,
		Complete:      rep.Complete,
		Mode:          "inventory",
		Source:        invSourceView{SpecsDir: inventorySpecsDir, SpecsState: strPtr(rep.SpecsState), RoadmapFile: strPtr(rep.RoadmapFile)},
		Configuration: invConfigView{FailOn: or(rep.FailOn, "error"), TaskScope: inventoryTaskScope},
		Folders:       []invFolderView{},
		Diagnostics:   []schema.Diagnostic{},
	}
	if rep.Outcome != "fatal" && rep.Summary != nil {
		s := rep.Summary
		view.Summary = &invSummaryView{
			FeatureFolders: s.FeatureFolders, OtherFolders: s.OtherFolders, SkippedEntries: s.SkippedEntries,
			Tasks: invTaskSumView{
				Checked: s.Tasks.Checked, Total: s.Tasks.Total, Percent: s.Tasks.Percent,
				FilesRead: s.Tasks.FilesRead, FilesMissing: s.Tasks.FilesMissing, FilesSkipped: s.Tasks.FilesSkipped,
			},
		}
		for _, f := range rep.Folders {
			view.Folders = append(view.Folders, folderToView(f))
		}
	}
	for _, d := range rep.Diagnostics {
		view.Diagnostics = append(view.Diagnostics, diagToSchema(d))
	}
	if rep.Outcome == "fatal" {
		d := Diagnostic{Code: "SS_ROOT_UNAVAILABLE", Severity: "error", Message: "Required input is unavailable."}
		if rep.Fatal != nil {
			d = *rep.Fatal
		} else if len(rep.Diagnostics) > 0 {
			d = rep.Diagnostics[0]
		}
		view.Error = &schema.Error{Code: d.Code, Message: d.Message, Source: locationOf(d.Path, d.Line)}
		view.Complete = false
		view.Summary = nil
		view.Folders = []invFolderView{}
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(view); err != nil {
		return []byte("{\"inventory_schema_version\":\"1.0\",\"outcome\":\"fatal\",\"complete\":false,\"mode\":\"inventory\",\"source\":{\"specs_dir\":\"specs\",\"specs_state\":null,\"roadmap_file\":null},\"configuration\":{\"fail_on\":\"error\",\"task_scope\":\"feature_folder_tasks_md\"},\"summary\":null,\"folders\":[],\"diagnostics\":[{\"code\":\"SS_INVOCATION_INVALID\",\"severity\":\"error\",\"message\":\"Report could not be serialized.\",\"source\":null,\"item_id\":null,\"related\":[]}],\"error\":{\"code\":\"SS_INVOCATION_INVALID\",\"message\":\"Report could not be serialized.\",\"source\":null}}\n")
	}
	return buf.Bytes()
}

// InventoryInvocationJSON is the inventory-shaped document for an
// invocation error under spec inventory --json (research R9).
func InventoryInvocationJSON(failOn, message string) []byte {
	return InventoryJSON(inventoryFatal("SS_INVOCATION_INVALID", message, "", or(failOn, "error")))
}

func folderToView(f Folder) invFolderView {
	v := invFolderView{Name: f.Name, Path: f.Path, Feature: f.Feature, NameForm: strPtr(f.NameForm), Skipped: f.Skipped}
	if f.Spec != nil {
		v.Spec = &invFileView{State: f.Spec.State, Reason: strPtr(f.Spec.Reason)}
	}
	if f.Plan != nil {
		v.Plan = &invFileView{State: f.Plan.State, Reason: strPtr(f.Plan.Reason)}
	}
	if f.Tasks != nil {
		t := &invTasksView{State: f.Tasks.State, Reason: strPtr(f.Tasks.Reason)}
		if f.Tasks.State == "read" {
			t.Checked, t.Total, t.Percent = f.Tasks.Checked, f.Tasks.Total, f.Tasks.Percent
		}
		v.Tasks = t
	}
	return v
}
