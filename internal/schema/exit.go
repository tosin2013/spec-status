package schema

// ExitCode is the exit code `spec status` uses for a report: 2 for a fatal
// or incomplete report, 1 when a diagnostic reaches the fail_on threshold,
// otherwise 0. fail_on "none" never hides a fatal report.
func ExitCode(r Report) int {
	if r.Outcome == "fatal" || !r.Complete {
		return 2
	}
	if r.Configuration.FailOn == "none" {
		return 0
	}
	rank := map[string]int{"info": 1, "warning": 2, "error": 3}
	need := rank[r.Configuration.FailOn]
	if need == 0 {
		return 0
	}
	for _, d := range r.Diagnostics {
		if rank[d.Severity] >= need {
			return 1
		}
	}
	return 0
}
