# Using spec status from Go

`github.com/tosin2013/spec-status/pkg/status` runs the same evaluation as `spec status --json --root <dir>` inside your Go program, with no binary to install. You get the schema 1.0 report as Go values, the exact bytes the CLI would print, and the exit code it would use.

**Stability: experimental until 1.0.**
- Report fields follow [JSON schema 1.0](../specs/001-spec-status/json-contract.md) and change only additively within it.
- Go function and type signatures may change in a minor release, and the release notes call out any change.
- `status.SchemaVersion` changes only when the JSON schema does.

## Install

```sh
go get github.com/tosin2013/spec-status/pkg/status@latest
```

It needs Go 1.26 or newer, and it pulls in no third-party modules: the package imports only the standard library and this module.

## Load a project

```go
package main

import (
	"fmt"
	"os"

	"github.com/tosin2013/spec-status/pkg/status"
)

func main() {
	rep, err := status.Load("path/to/project", status.Options{FailOn: "warning"})
	if err != nil {
		// Invalid input only: an empty root or an unknown FailOn.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if rep.Outcome == "fatal" {
		// The roadmap could not be read, for example SS_ROADMAP_UNAVAILABLE.
		fmt.Fprintln(os.Stderr, rep.Error.Code, rep.Error.Message)
		os.Exit(status.ExitCode(rep)) // always 2
	}
	for _, it := range rep.Items {
		fmt.Println(it.ID, it.Status, it.Tasks.State)
		for _, ph := range it.Tasks.Phases { // nil when the tasks file was not read
			fmt.Println("  ", ph.Heading, ph.Checked, "/", ph.Total)
		}
	}
	os.Exit(status.ExitCode(rep))
}
```

| Function | Behaves like |
|---|---|
| `Load(root, Options{Roadmap, FailOn})` | `spec status --json --root root [--roadmap …] [--fail-on …]`. A relative path resolves against the process working directory |
| `Marshal(report)` | The exact bytes that command prints, including the trailing newline. For a report you built or filtered yourself, nil `Items`, `Diagnostics`, `Dependencies` and `Related` are written as `[]`, as the schema requires |
| `ExitCode(report)` | Its exit code: 0 below the threshold, 1 when a finding reaches it, 2 when fatal. An empty or unrecognized `Configuration.FailOn` is treated as `error`, the CLI default |

**Errors versus fatal reports:**
- `Load` returns a Go `error` (wrapping `status.ErrInvalidInput`) only when the CLI would reject the flags themselves.
- Every problem with the project, such as a missing, unrecognized or unsafe roadmap, comes back as a report with `Outcome == "fatal"` and the same `SS_*` code as the JSON.

**Null fields:**
- A field that can be `null` in the JSON is a pointer in Go (`*string`, `*int`, `*bool`, `*Location`).
- `Tasks.Phases` is nil when the tasks file was not read, and an empty slice when it has no `## Phase` headings.
- Percentages are `status.Percent`. Use `Valid()`, `Hundredths()` (4167 for 41.67) or `Text()`, never a float.

## If you already run the binary

The CLI's `--json` output decodes into `status.Report` with `encoding/json`, and `status.Marshal` reproduces it byte for byte. A tool can therefore start by running `spec status --json`, decode the output into these types, and later switch to `Load` with no other change. The test suite checks this round trip, and checks that `Load` equals the CLI on every fixture at every `--fail-on` threshold.

## Safety

`Load` has the CLI's guarantees:
- It reads only under the root, rejecting `..`, symlinks, sensitive file names and oversize files, with bounded work.
- It never writes, opens network connections or starts processes, and no environment variable changes its result.
- It is safe to call from several goroutines.

Report text such as titles, headings and messages is copied from the project's files. Treat it as untrusted, and sanitize it before printing it to a terminal or rendering it as HTML.

## Limits

- Only `spec status` has a library form. `spec inventory`, human text, filters and `--explain` are CLI-only for now.
- Like the CLI, the package reads the narrow, pinned roadmap format described in [limits.md](limits.md).
