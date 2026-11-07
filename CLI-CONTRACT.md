# CLI contract — candidate interface version 0.1

**Draft PROFILE-R1, pending DC-01/DC-04.** Commands here describe intended behavior, not currently runnable software or install instructions. CLI release version, interface draft version, JSON schema version and roadmap input version are separate identities.

## Working command surface

`spec status` is the working name, not an established executable. Candidate status options:

| Option | Candidate behavior / retained constraint |
|---|---|
| `--json` | One JSON document from the same report as text, including fatal errors; retained spelling |
| `--fail-on info\|warning\|error\|none` | Threshold defaults to error; all findings still emitted; retained choices |
| `--root PATH` | Explicit selected project root, relative to invocation directory; otherwise nearest ancestor with .specify |
| `--roadmap PATH` | Override input; relative to selected root, or absolute only if confined within that root |
| `--help` | Candidate usage/scope/version guidance; no filesystem/project input required |

Working identity command `spec --version` will report CLI version, tagged commit, Go version and target; development builds clearly labelled. Version/help are metadata operations, not a roadmap report. PROFILE-R1 recommends a separate error for mixing help/version/status flags; exact parser behavior needs approval. Do not invent a working binary today.

No quarter option exists. `--as-of` is not part of the advertised candidate profile unless a concrete non-quarter deterministic use is chosen under DC-04. No inferred date or inert compatibility flag. `--format=json` is not an approved alias. No write/agent/sync/Git/network option belongs to status.

PROFILE-R1 no config file/environment override: explicit flags then documented defaults. No credential/environment-secret discovery. Unknown options, invalid thresholds, missing flag values, conflicting repeated scalar options and incompatible subcommands are invocation errors (exit 2). Repeated identical scalar flags may either reject or deduplicate; recommend reject uniformly, pending approval. The low-level parser must detect an explicit --json request early enough to return structured invocation errors; edge cases such as misspelled --json are documented plain-text errors, not guessed intent.

## Output

Text report should be glanceable: selected relative source, inventory/raw statuses, active/recorded-verified numerator+denominator+percentage, task counts with unique-file scope, superseded row count and successors, active blocked/review counts, then findings with stable code/severity/path/line. Show null as N/A. Distinguish invalid supplied replacement, historical verified status and successful active completion. Every skipped required/optional evidence source is visible without secret contents.

Items remain document order; diagnostics sort normalized relative path then line then code, plus stable tie-breaks specified in JSON contract. No volatile timestamp/clock context by default. Color/TTY heuristics are not a committed feature; candidate profile uses plain text on all terminals for stability. No special AI mode. Empty supported input does not become a fatal error simply because there are no items.

JSON stdout: exactly one UTF-8 document plus candidate terminal newline, no leading/trailing logs, progress spinner or human explanation. Text/log errors must not be mixed into JSON stdout. PROFILE-R1 recommends stderr empty in JSON mode unless a separate unavoidable process failure prevents serialization; such failures cannot be called valid JSON acceptance. Error documents and completeness are defined in [JSON-CONTRACT.md](JSON-CONTRACT.md).

## Exit contract — retained, not proposal

| Exit | Meaning |
|---|---|
| 0 | No finding reaches the configured threshold. Does not mean complete work, code correctness or independent verification. |
| 1 | Findings at/above threshold. Not a runtime crash. |
| 2 | Invocation/configuration/required-input/fatal-parse failure. Never suppressed by none. |

Severity order is info < warning < error. Threshold info fails for any finding, warning for warnings/errors, error for errors, none for no findings-based failure. Fatal conditions always exit 2 regardless of diagnostics. Ordinary unfinished work alone is not a diagnostic. Blocked unfinished is warning; unknown status warning; planned missing spec ref info; checked missing local target error and unchecked info. Other exact severities remain DC-03 proposals in diagnostic registry.

## Illustrative outcomes, not executed tests

- Supported roadmap with ordinary in-progress work and no diagnostic: exit 0 at default error, not “project complete.”
- Unknown status warning only: exit 0 at default error, exit 1 at warning/info, findings remain visible.
- Checked missing explicit local target: exit 1 at default error, exit 0 at none if no fatal input failure.
- Fatal unsupported item layout, missing required roadmap, unsafe required source or exhausted required-input budget under candidate policy: exit 2 even at none; one JSON error if --json requested.
- A superseded dependency produces human review even if its successor is verified. Candidate warning/blocked behavior is in DC-03 and does not redefine the source relationship.

## Compatibility gates

Approve executable/grammar/help/error behavior before coding. Preserve existing --json/snake_case concepts rather than inventing aliases to match earlier illustrative reviews. Document breaking command/meaning changes with CLI SemVer rules; do not conflate input v0.1.1 with CLI 0.1.1. Any decision to support an alternate satisfaction policy or date/config feature updates help, JSON effective configuration, fixtures and compatibility policy explicitly. No command examples here authorize mutations or release creation.
