package schema

// ExitCode is the exit code `spec status` uses for a report: 2 for a fatal
// or incomplete report, 1 when a diagnostic reaches the fail_on threshold,
// otherwise 0. fail_on "none" never hides a fatal report; an empty fail_on
// means "error", the CLI default.
func ExitCode(r Report) int {
	severities := make([]string, len(r.Diagnostics))
	for i, d := range r.Diagnostics {
		severities[i] = d.Severity
	}
	return ExitCodeFor(r.Outcome, r.Complete, r.Configuration.FailOn, severities)
}

// ExitCodeFor applies the same rule to its four inputs, so callers need not
// build a Report. An unrecognized failOn compares as "error".
func ExitCodeFor(outcome string, complete bool, failOn string, severities []string) int {
	if outcome == "fatal" || !complete {
		return 2
	}
	if failOn == "none" {
		return 0
	}
	rank := map[string]int{"info": 1, "warning": 2, "error": 3}
	need, ok := rank[failOn]
	if !ok {
		need = rank["error"]
	}
	for _, s := range severities {
		if rank[s] >= need {
			return 1
		}
	}
	return 0
}
