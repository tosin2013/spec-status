# CLI contract — interface version 0.1

**PROFILE-R1, approved as built for DC-01 and DC-04 on 11 November 2026 (D-2026-11-11-02).** Commands here describe the development build's behavior, not a release or install instructions. CLI release version, interface draft version, JSON schema version and roadmap input version are separate identities.

## Command dependency

The development CLI uses Cobra for commands, help, and flags. The owner approved that dependency on 7 November 2026. Human status text may use Lip Gloss, which the owner approved separately, for headings, aligned columns, progress marks, and restrained color. Lip Gloss is not Bubble Tea and this contract adds no dashboard or watch mode. JSON, fatal text, version, and completion stay plain. `--color` is `auto`, `always`, or `never`. Auto emits no ANSI unless stdout is a terminal, `NO_COLOR` is unset, and `TERM` is not `dumb`. Color is never the only distinction. Parsing, metrics, and root-confined reads remain standard-library code. Shell completion is an explicit `spec completion` command that prints a script to stdout and does not install it. For that process, `BASH_COMP_DEBUG_FILE` is unset and then restored so Cobra's hidden `__complete` error path cannot create or append the log. Concurrent `Run` calls in one process wait on that restoration. The product binary calls it once. The generated script is shell text and is not executed by spec; if a person sources it, the shell may append to that variable. The other interface choices below are PROFILE-R1 as built, which the owner approved under DC-04 on 11 November 2026 (D-2026-11-11-02).

## Working command surface

`spec status` is the working command. The development build also exposes `spec version`, `spec completion`, and `spec help`. Status options:

| Option | Behavior / retained constraint |
|---|---|
| `--json` | One JSON document from the same report as text, including fatal errors; retained spelling |
| `--fail-on info\|warning\|error\|none` | Threshold defaults to error; all findings still emitted; retained choices |
| `--root PATH` | Explicit selected project root, relative to invocation directory; otherwise nearest ancestor with .specify |
| `--roadmap PATH` | Override input; relative to selected root, or absolute only if confined within that root |
| `--help` | Usage/scope/version guidance; no filesystem/project input required |
| `--color auto\|always\|never` | Owner-approved human-text color control. Default `auto` emits no ANSI unless stdout is a terminal, `NO_COLOR` is unset, and `TERM` is not `dumb`. Does not change JSON or exit codes |
| `--status TOKEN` | Human-text view. Exact, case-sensitive match on the item's raw status, including an unknown token. Repeat with different tokens to match any of them. Repeating the same token is an invocation error |
| `--blocked` | Human-text view. An item with any unsatisfied dependency, including a historical item. This is not the active-blockers summary count |
| `--human-review` | Human-text view. An item with any dependency in state `human_review`, including a historical item |
| `--active` | Human-text view. Active items only |
| `--historical` | Human-text view. Items outside the active set, including deferred, abandoned, and valid supersession |
| `--explain` | Human text only. Under each shown item, print mechanical reasons and file:line evidence. Does not filter items by itself |

Working identity commands `spec version` and `spec --version` report CLI version, commit, Go version, and target; development builds are clearly labelled. Version and help are metadata operations, not a roadmap report. Mixing help or version with status or JSON is an invocation error. The development build is not a release.

No quarter option exists. `--as-of` is not part of the interface; the owner rejected it under DC-04 (D-2026-11-11-02). No inferred date or inert compatibility flag. `--format=json` is not an approved alias. No write/agent/sync/Git/network option belongs to status.

No config file or environment override (DC-04, D-2026-11-11-02): explicit flags then documented defaults. No credential/environment-secret discovery. Unknown options, invalid thresholds, missing flag values, conflicting repeated scalar options and incompatible subcommands are invocation errors (exit 2). Repeated identical scalar flags are rejected uniformly as invocation errors. The low-level parser must detect an explicit --json request early enough to return structured invocation errors; edge cases such as misspelled --json are documented plain-text errors, not guessed intent.

## Output

Text report should be glanceable. It opens with a compact static summary panel grouping recorded verified, checked tasks, active blockers, and human review, then selected relative source, inventory/raw statuses, task-file coverage, superseded rows, and findings with stable code/severity/path/line. Those counts stay separate and always describe the full roadmap. A view never shrinks those denominators or hides a diagnostic. Show null as N/A. Distinguish invalid supplied replacement, historical verified status and successful active completion. Every skipped required/optional evidence source is visible without secret contents.

### Human-text views

The owner approved these views on 7 November 2026. They are not a new schema and they do not approve DC-01–DC-08.

Families are `--status` (exact raw token, case-sensitive, unknown tokens allowed), `--blocked`, `--human-review`, `--active`, and `--historical`. Different families combine with AND. Repeated different `--status` values combine with OR. `--active` together with `--historical` matches nothing. Repeating `--blocked`, `--human-review`, `--active`, `--historical`, `--explain`, or the same `--status` token is an invocation error (exit 2), including when `--json` is set.

When any filter family is set, human text prints `shown N of M items` before the item list. `M` is the full item count. If `N` is 0, the next line is `no item matches this view; totals above are the full roadmap`. That message is not a diagnostic and is not fatal. Exit 0, 1, or 2 still comes from the full result: hidden items do not suppress findings, and a no-match with no finding at the threshold stays exit 0.

`--explain` adds lines under each shown item, or under every item when no filter is set. Each line starts with `explain:`. It may state that a valid supersession excluded the row from active and recorded-verified counts, that a dependency is unsatisfied or `human_review`, and how checked tasks were counted. A file and line are printed only when that location was retained from the parse of the roadmap or a linked task file. If the location was not retained, the line says so. The text does not judge code correctness, does not crawl, and does not invent a line number. Repo-supplied terminal controls in those lines are sanitized the same way as other human text. JSON keeps raw strings and does not gain explain lines.

The summary panel, the diagnostic list, and JSON stay on the unfiltered result. `--blocked` can list a historical item that the active-blockers count does not include.

Items remain document order; diagnostics sort normalized relative path then line then code, plus stable tie-breaks specified in JSON contract. No volatile timestamp/clock context by default. Color follows the owner-approved `--color` rule above; JSON and fatal text are never colored. No special AI mode. Empty supported input does not become a fatal error simply because there are no items.

JSON stdout: exactly one UTF-8 document plus a terminal newline, no leading/trailing logs, progress spinner or human explanation. Text/log errors must not be mixed into JSON stdout. Stderr is empty in JSON mode (D-2026-11-11-02) unless a separate unavoidable process failure prevents serialization; such failures cannot be called valid JSON acceptance. Error documents and completeness are defined in [json-contract.md](json-contract.md).

## Exit contract — retained, not proposal

| Exit | Meaning |
|---|---|
| 0 | No finding reaches the configured threshold. Does not mean complete work, code correctness or independent verification. |
| 1 | Findings at/above threshold. Not a runtime crash. |
| 2 | Invocation/configuration/required-input/fatal-parse failure. Never suppressed by none. |

Severity order is info < warning < error. Threshold info fails for any finding, warning for warnings/errors, error for errors, none for no findings-based failure. Fatal conditions always exit 2 regardless of diagnostics. Ordinary unfinished work alone is not a diagnostic. Blocked unfinished is warning; unknown status warning; planned missing spec ref info; checked missing local target error and unchecked info. Other exact severities are in the diagnostic registry approved under DC-03 (D-2026-11-07-03).

## Illustrative outcomes, not executed tests

- Supported roadmap with ordinary in-progress work and no diagnostic: exit 0 at default error, not “project complete.”
- Unknown status warning only: exit 0 at default error, exit 1 at warning/info, findings remain visible.
- Checked missing explicit local target: exit 1 at default error, exit 0 at none if no fatal input failure.
- Fatal unsupported item layout, missing required roadmap, unsafe required source or exhausted required-input budget under the fixed DC-06 limits (D-2026-11-11-02): exit 2 even at none; one JSON error if --json requested.
- A superseded dependency produces human review even if its successor is verified. The warning/blocked behavior is DC-03 (D-2026-11-07-03) and does not redefine the source relationship.

## Compatibility gates

Executable/grammar/help/error behavior was approved as built on 11 November 2026 (D-2026-11-11-02); a later change needs a recorded decision. Preserve existing --json/snake_case concepts rather than inventing aliases to match earlier illustrative reviews. Document breaking command/meaning changes with CLI SemVer rules; do not conflate input v0.1.1 with CLI 0.1.1. Any decision to support an alternate satisfaction policy or date/config feature updates help, JSON effective configuration, fixtures and compatibility policy explicitly. No command examples here authorize mutations or release creation.

## Not part of the working command

[timing-evidence.md](timing-evidence.md) and [companion.md](companion.md) are draft proposals. `spec companion` is not a command in this build. No quarter option is added. Human styled status text remains for supported terminal surfaces. JSON remains the machine evidence. An external harness is not a native integration.

[two-mode-direction.md](two-mode-direction.md) is an owner decision for a later change. This build still reads `.specify/memory/roadmap.md` and still exits 2 when that file is missing. No inventory command is defined. Flag names, default selection, and how a template version is recognized are open. Filters, when a roadmap report exists, still do not hide findings.
