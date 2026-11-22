package status_test

import (
	"fmt"

	"github.com/tosin2013/spec-status/pkg/status"
)

// Load a project and read its report, as `spec status --json` would print it.
func Example() {
	rep, err := status.Load("../../testdata/example-roadmap", status.Options{FailOn: "warning"})
	if err != nil {
		fmt.Println("invalid input:", err)
		return
	}
	if rep.Outcome == "fatal" {
		fmt.Println("fatal:", rep.Error.Code)
		return
	}
	s := rep.Summary
	fmt.Println("schema", rep.SchemaVersion, "outcome", rep.Outcome)
	fmt.Printf("recorded verified %d/%d (%s%%)\n", s.VerifiedActiveItems, s.ActiveItems, s.VerifiedPercent.Text())
	for _, it := range rep.Items {
		phases := "unknown"
		if it.Tasks.Phases != nil {
			phases = fmt.Sprint(len(it.Tasks.Phases))
		}
		fmt.Println(it.ID, it.Status, "tasks:", it.Tasks.State, "phases:", phases)
	}
	fmt.Println("exit code", status.ExitCode(rep))
	// Output:
	// schema 1.0 outcome report
	// recorded verified 1/2 (50%)
	// 001 implemented tasks: read phases: 0
	// 002 verified tasks: read phases: 0
	// 003 planned tasks: not_linked phases: unknown
	// exit code 1
}
