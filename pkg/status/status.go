package status

import (
	"errors"
	"fmt"

	"github.com/tosin2013/spec-status/internal/schema"
	core "github.com/tosin2013/spec-status/internal/status"
)

// SchemaVersion is the JSON schema version of every Report.
const SchemaVersion = schema.Version

// The report types. Field names, order and null rules match the JSON
// contract (specs/001-spec-status/json-contract.md).
type (
	Report        = schema.Report
	Source        = schema.Source
	Configuration = schema.Configuration
	Summary       = schema.Summary
	TaskSummary   = schema.TaskSummary
	Item          = schema.Item
	Ref           = schema.Ref
	Location      = schema.Location
	Supersession  = schema.Supersession
	Dependency    = schema.Dependency
	Tasks         = schema.Tasks
	Phase         = schema.Phase
	Diagnostic    = schema.Diagnostic
	Error         = schema.Error
	Percent       = schema.Percent
)

// Options are the `spec status` flags that affect the report. The zero
// value means the CLI defaults.
type Options struct {
	// Roadmap is the roadmap path, as --roadmap. Empty means
	// .specify/memory/roadmap.md under the root.
	Roadmap string
	// FailOn is the threshold for ExitCode: "error" (default when empty),
	// "warning", "info" or "none".
	FailOn string
}

// ErrInvalidInput is wrapped by every error Load returns.
var ErrInvalidInput = errors.New("spec status: invalid input")

// Load evaluates the project at root exactly as
// `spec status --json --root root` does. A relative root or roadmap is
// resolved against the process working directory, as the CLI does. The
// error is non-nil only for invalid input: an empty root or an unknown
// FailOn. Fatal problems are returned in the Report.
func Load(root string, opts Options) (Report, error) {
	if root == "" {
		return Report{}, fmt.Errorf("%w: root must not be empty", ErrInvalidInput)
	}
	switch opts.FailOn {
	case "", "error", "warning", "info", "none":
	default:
		return Report{}, fmt.Errorf("%w: FailOn %q must be one of error, warning, info, none", ErrInvalidInput, opts.FailOn)
	}
	rep := core.Evaluate(core.Option{Root: root, Roadmap: opts.Roadmap, FailOn: opts.FailOn})
	return core.Schema(rep), nil
}

// Marshal returns the exact bytes `spec status --json` prints for the
// report, including the trailing newline.
func Marshal(r Report) []byte { return schema.Marshal(r) }

// ExitCode returns the exit code `spec status` uses for the report: 2 for a
// fatal report, 1 when a diagnostic reaches Configuration.FailOn, else 0.
func ExitCode(r Report) int { return schema.ExitCode(r) }
