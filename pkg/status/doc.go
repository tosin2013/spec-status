// Package status evaluates a Spec Kit project the same way
// `spec status --json --root <root>` does, in-process, and returns the
// schema 1.0 report as Go values.
//
//	rep, err := status.Load("path/to/project", status.Options{})
//	if err != nil {
//		// invalid input only: an empty root or an unknown FailOn
//	}
//	if rep.Outcome == "fatal" {
//		// for example SS_ROADMAP_UNAVAILABLE, in rep.Error.Code
//	}
//	os.Exit(status.ExitCode(rep)) // the binary's 0, 1 or 2
//
// Load returns fatal problems, such as a missing or unrecognized roadmap,
// inside the report (Outcome "fatal", Error set), exactly as the JSON does.
// Marshal returns the exact bytes the CLI prints, and ExitCode returns its
// exit code; the CLI produces both through the same code. Decoding the CLI's
// JSON into Report with encoding/json and calling Marshal reproduces it.
//
// Safety: Load only reads files under the root, with the CLI's confinement
// and limits (no symlinks, no sensitive names, bounded sizes). It never
// writes, opens network connections or starts processes, no environment
// variable changes its result, and it is safe for concurrent use. Report text such as titles,
// headings and messages is raw project content: sanitize it before printing
// it to a terminal.
//
// Stability: this API is experimental until 1.0. Report fields follow JSON
// schema 1.0 (SchemaVersion) and change only additively within it. Function
// and type signatures may change in a minor release; release notes call out
// any change.
package status
