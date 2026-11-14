package status

import (
	"sort"

	"github.com/tosin2013/spec-status/internal/safeio"
)

// Feature 006: roadmap-free feature inventory (D-2026-11-12-08,
// specs/006-two-mode-inventory). Inventory lists the immediate children of
// specs/, stats spec.md, plan.md and tasks.md in feature folders, and counts
// tasks.md checkboxes with the roadmap-mode rules. It never reads the
// roadmap and never infers status, priority or dependencies.

const (
	inventorySpecsDir  = "specs"
	inventoryTaskScope = "feature_folder_tasks_md"
	inventoryRoadmap   = ".specify/memory/roadmap.md"
)

// inventoryAfterRead is a test seam: when set, it observes every file
// inventory reads. Production leaves it nil.
var inventoryAfterRead func(path string)

// InventoryOption selects one read-only inventory.
type InventoryOption struct {
	Root   string
	FailOn string
	CWD    string
}

// InventoryReport is the single result used for inventory text and JSON.
type InventoryReport struct {
	Outcome     string
	Complete    bool
	FailOn      string
	SpecsState  string // present, missing, skipped; "" on an early fatal
	RoadmapFile string // present, absent, skipped; "" on an early fatal
	Folders     []Folder
	Summary     *InventorySummary
	Diagnostics []Diagnostic
	Fatal       *Diagnostic
}

// Folder is one immediate child directory of specs/ (or a skipped entry).
type Folder struct {
	Name     string
	Path     string
	Feature  bool
	NameForm string // sequential, timestamp; "" when not a feature folder
	Skipped  bool
	Spec     *FileState
	Plan     *FileState
	Tasks    *TaskState
}

// FileState is the presence of spec.md or plan.md. Contents are not read.
type FileState struct {
	State  string // present, missing, skipped
	Reason string // unsafe, sensitive; "" otherwise
}

// TaskState is tasks.md: read and counted, missing, or skipped.
type TaskState struct {
	State   string // read, missing, skipped
	Reason  string
	Checked int
	Total   int
	Percent Percent
}

// InventorySummary totals one inventory.
type InventorySummary struct {
	FeatureFolders int
	OtherFolders   int
	SkippedEntries int
	Tasks          TaskSummary
}

func inventoryFatal(code, message, path, failOn string) InventoryReport {
	d := Diagnostic{Code: code, Severity: "error", Message: message, Path: path}
	return InventoryReport{
		Outcome:     "fatal",
		Complete:    false,
		FailOn:      failOn,
		Folders:     []Folder{},
		Diagnostics: []Diagnostic{d},
		Fatal:       &d,
	}
}

// EvaluateInventory lists feature folders under the selected root. It does not write.
func EvaluateInventory(opt InventoryOption) InventoryReport {
	failOn := opt.FailOn
	if failOn == "" {
		failOn = "error"
	}
	fs, ferr := openProject(opt.CWD, opt.Root)
	if ferr != nil {
		return inventoryFatal(ferr.Code, ferr.Message, "", failOn)
	}
	if inventoryAfterRead != nil {
		fs.SetAfterReadHook(inventoryAfterRead)
	}
	rep := InventoryReport{Outcome: "report", Complete: true, FailOn: failOn, Folders: []Folder{}}
	rep.RoadmapFile = roadmapFileState(fs)
	sum := &InventorySummary{}

	entries, err := fs.ListDir(inventorySpecsDir)
	switch code := safeCode(err); {
	case err == nil:
		rep.SpecsState = "present"
	case code == safeio.CodeNotFound:
		rep.SpecsState = "missing"
		rep.Diagnostics = append(rep.Diagnostics, Diagnostic{Code: "SS_SPECS_DIRECTORY_MISSING", Severity: "info", Message: "Specs directory is absent.", Path: inventorySpecsDir})
	case code == safeio.CodeLimit:
		return inventoryFatal("SS_INPUT_LIMIT", "Input budget exceeded.", inventorySpecsDir, failOn)
	case code == safeio.CodeChanged:
		return inventoryFatal("SS_SOURCE_CHANGED", "Input changed while it was being read.", inventorySpecsDir, failOn)
	case code == safeio.CodeSensitive:
		rep.SpecsState = "skipped"
		rep.Diagnostics = append(rep.Diagnostics, sensitiveDiag(inventorySpecsDir))
	default:
		rep.SpecsState = "skipped"
		rep.Diagnostics = append(rep.Diagnostics, Diagnostic{Code: "SS_PATH_UNSAFE", Severity: "error", Message: "Specs path is a symlink, not a directory, or outside the confinement rules.", Path: inventorySpecsDir})
	}

	for _, e := range entries {
		path := inventorySpecsDir + "/" + e.Name
		switch e.Kind {
		case safeio.KindFile:
			continue // files directly under specs/ are not folders
		case safeio.KindDir:
			form := featureNameForm(e.Name)
			if form == "" {
				sum.OtherFolders++
				rep.Folders = append(rep.Folders, Folder{Name: e.Name, Path: path})
				rep.Diagnostics = append(rep.Diagnostics, Diagnostic{Code: "SS_FOLDER_NOT_FEATURE", Severity: "info", Message: "Folder name is not a Spec Kit feature folder name.", Path: path})
				continue
			}
			sum.FeatureFolders++
			f := Folder{Name: e.Name, Path: path, Feature: true, NameForm: form}
			if fatal := inspectFeature(fs, &f, &rep, &sum.Tasks, failOn); fatal != nil {
				return *fatal
			}
			rep.Folders = append(rep.Folders, f)
		default:
			sum.SkippedEntries++
			rep.Folders = append(rep.Folders, Folder{Name: e.Name, Path: path, Skipped: true})
			if e.Kind == safeio.KindSensitive {
				rep.Diagnostics = append(rep.Diagnostics, sensitiveDiag(path))
			} else {
				rep.Diagnostics = append(rep.Diagnostics, Diagnostic{Code: "SS_PATH_UNSAFE", Severity: "error", Message: "Entry is a symlink, special file, or on another device and was not examined.", Path: path})
			}
		}
	}
	sum.Tasks.Percent = percentOf(sum.Tasks.Checked, sum.Tasks.Total)
	sum.Tasks.Scope = inventoryTaskScope
	rep.Summary = sum
	sort.SliceStable(rep.Diagnostics, func(i, j int) bool { return diagLess(rep.Diagnostics[i], rep.Diagnostics[j]) })
	return rep
}

// featureNameForm classifies a folder name by the two forms Spec Kit
// creates (research R7): YYYYMMDD-HHMMSS-name, checked first, then
// NNN-name with three or more ASCII digits. Dates are not validated.
func featureNameForm(name string) string {
	digits := func(s string) int {
		n := 0
		for n < len(s) && s[n] >= '0' && s[n] <= '9' {
			n++
		}
		return n
	}
	if len(name) > 16 && digits(name) == 8 && name[8] == '-' && digits(name[9:]) == 6 && name[15] == '-' {
		return "timestamp"
	}
	if n := digits(name); n >= 3 && n+1 < len(name) && name[n] == '-' {
		return "sequential"
	}
	return ""
}

// inspectFeature fills the file states of one feature folder and adds its
// tasks.md counts to the totals. It returns a fatal report when the input
// cannot be counted reliably.
func inspectFeature(fs *safeio.FS, f *Folder, rep *InventoryReport, totals *TaskSummary, failOn string) *InventoryReport {
	for _, name := range []string{"spec.md", "plan.md"} {
		st, fatal := fileState(fs, f.Path+"/"+name, rep, failOn)
		if fatal != nil {
			return fatal
		}
		if name == "spec.md" {
			f.Spec = st
		} else {
			f.Plan = st
		}
	}
	p := f.Path + "/tasks.md"
	st, fatal := fileState(fs, p, rep, failOn)
	if fatal != nil {
		return fatal
	}
	f.Tasks = &TaskState{State: st.State, Reason: st.Reason}
	switch st.State {
	case "missing":
		totals.FilesMissing++
		return nil
	case "skipped":
		totals.FilesSkipped++
		return nil
	}
	data, err := fs.ReadRel(p)
	if err != nil {
		switch code := safeCode(err); code {
		case safeio.CodeLimit:
			r := inventoryFatal("SS_INPUT_LIMIT", "Input budget exceeded.", p, failOn)
			return &r
		case safeio.CodeChanged:
			r := inventoryFatal("SS_SOURCE_CHANGED", "Input changed while it was being read.", p, failOn)
			return &r
		case safeio.CodeSensitive:
			rep.Diagnostics = append(rep.Diagnostics, sensitiveDiag(p))
			f.Tasks = &TaskState{State: "skipped", Reason: "sensitive"}
		default:
			rep.Diagnostics = append(rep.Diagnostics, unsafeFileDiag(p))
			f.Tasks = &TaskState{State: "skipped", Reason: "unsafe"}
		}
		totals.FilesSkipped++
		return nil
	}
	tf, perr := parseTasks(p, data)
	if perr != nil {
		r := inventoryFatal(perr.Code, perr.Message, perr.Path, failOn)
		r.Fatal.Line, r.Diagnostics[0].Line = perr.Line, perr.Line
		return &r
	}
	totals.FilesRead++
	totals.Checked += tf.Checked
	totals.Total += tf.Total
	if totals.Total > maxTasks {
		r := inventoryFatal("SS_INPUT_LIMIT", "Task checkbox limit exceeded across task files.", p, failOn)
		return &r
	}
	f.Tasks = &TaskState{State: "read", Checked: tf.Checked, Total: tf.Total, Percent: percentOf(tf.Checked, tf.Total)}
	return nil
}

// fileState stats one file without reading it.
func fileState(fs *safeio.FS, p string, rep *InventoryReport, failOn string) (*FileState, *InventoryReport) {
	ent, err := fs.StatRel(p)
	if err != nil {
		switch safeCode(err) {
		case safeio.CodeLimit:
			r := inventoryFatal("SS_INPUT_LIMIT", "Input budget exceeded.", p, failOn)
			return nil, &r
		case safeio.CodeChanged:
			r := inventoryFatal("SS_SOURCE_CHANGED", "Input changed while it was being read.", p, failOn)
			return nil, &r
		case safeio.CodeSensitive:
			rep.Diagnostics = append(rep.Diagnostics, sensitiveDiag(p))
			return &FileState{State: "skipped", Reason: "sensitive"}, nil
		default:
			rep.Diagnostics = append(rep.Diagnostics, unsafeFileDiag(p))
			return &FileState{State: "skipped", Reason: "unsafe"}, nil
		}
	}
	switch {
	case !ent.Exists:
		return &FileState{State: "missing"}, nil
	case !ent.Regular:
		rep.Diagnostics = append(rep.Diagnostics, unsafeFileDiag(p))
		return &FileState{State: "skipped", Reason: "unsafe"}, nil
	}
	return &FileState{State: "present"}, nil
}

// roadmapFileState reports whether the default roadmap path exists. The
// file is never opened (FR-014).
func roadmapFileState(fs *safeio.FS) string {
	ent, err := fs.StatRel(inventoryRoadmap)
	switch {
	case err != nil:
		return "skipped"
	case ent.Exists:
		return "present"
	default:
		return "absent"
	}
}

func safeCode(err error) string {
	if err == nil {
		return ""
	}
	if se, ok := err.(*safeio.Error); ok {
		return se.Code
	}
	return safeio.CodeUnavailable
}

func sensitiveDiag(path string) Diagnostic {
	return Diagnostic{Code: "SS_REFERENCE_SENSITIVE_SKIPPED", Severity: "error", Message: "Entry was skipped because it matched the sensitive-name matcher; it was not opened.", Path: path}
}

func unsafeFileDiag(path string) Diagnostic {
	return Diagnostic{Code: "SS_PATH_UNSAFE", Severity: "error", Message: "File is a symlink, not a regular file, or could not be read under the confinement rules; it was not read.", Path: path}
}

// InventoryExitCode applies the threshold with the roadmap-mode rule:
// fatal or incomplete is 2, a finding at or above --fail-on is 1.
func InventoryExitCode(rep InventoryReport) int {
	return ExitCode(Report{Outcome: rep.Outcome, Complete: rep.Complete, FailOn: rep.FailOn, Diagnostics: rep.Diagnostics})
}
