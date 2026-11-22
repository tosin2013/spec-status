# Contract: `github.com/tosin2013/spec-status/pkg/status`

**Stability:** experimental until 1.0 (D-2026-11-21-01, constitution Art. V).
- Report fields follow JSON schema 1.0 and change only additively within it.
- Function and type signatures may change in a minor release, and release notes call out any change.
- `SchemaVersion` changes only when the JSON schema does.

```go
package status

const SchemaVersion = "1.0"

type Options struct {
	Roadmap string // as --roadmap; "" = .specify/memory/roadmap.md
	FailOn  string // "" = "error"; or "warning", "info", "none"
}

// Load evaluates the project at root exactly as `spec status --json --root root` does.
// Fatal problems (for example SS_ROADMAP_UNAVAILABLE) are returned in the Report
// (Outcome "fatal", Error set). err is non-nil only for invalid input:
// an empty root or an unknown FailOn.
func Load(root string, opts Options) (Report, error)

// Marshal returns the exact bytes `spec status --json` prints for the report,
// including the trailing newline.
func Marshal(r Report) []byte

// ExitCode returns the exit code the binary uses for the report: 0, 1 or 2.
func ExitCode(r Report) int

// Report, Source, Configuration, Summary, TaskSummary, Item, Ref, Supersession,
// Dependency, Tasks, Phase, Diagnostic, Error, Location, Percent:
// aliases of the internal schema types (see data-model.md).
```

**Guarantees:**
- For every root and options, `Marshal(Load(...))` equals the CLI's `--json` bytes, and `ExitCode` equals the CLI's exit code. This is tested on every fixture root at every threshold.
- Decoding CLI JSON into `Report` with `encoding/json` and calling `Marshal` reproduces the bytes.
- No writes, network, subprocesses or environment reads. Root-confined reads with the core's limits. Imports are the standard library and this module only.
- Safe for concurrent use.
- Report text (titles, headings, messages) is raw, untrusted project content. Sanitize it before printing it to a terminal.
