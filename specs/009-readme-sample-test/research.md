# Research: README sample test

Phase 0 for `specs/009-readme-sample-test/plan.md`.

## R1 — Where the test runs

- **Decision:** `internal/cli/readme_samples_test.go`, driving the real entry point `Run(args, repoRoot, stdout, stderr)` with the label's arguments. `withBuild(t, resolveBuild("", "0.0.0-dev", "unknown", nil, false))` forces the development-build heading; `t.Setenv("COLUMNS", "80")` fixes the width (`textColumns`, `term_unix.go`); a `bytes.Buffer` is not a terminal, so `--color` auto/never both give no ANSI.
- **Rationale:** a reader copies the README command; running it through `Run` checks the same flag parsing, help-free path and output stream the reader gets. A status-package test calling `Text` directly would miss flag handling.
- **Alternatives:** keep per-sample tests in `internal/status` (rejected: duplicates the label parsing, misses CLI parsing); run the built binary as a subprocess (rejected: tests avoid subprocesses, and `go test` should not need a build step).

## R2 — Label form

- **Decision:** an HTML comment on its own line, before the fenced block (blank lines allowed between):
  `<!-- spec-sample: status --color=never --root testdata/example-roadmap -->`
  The text after `spec-sample:` is split on spaces into arguments (no quoting needed by any current sample; a quote character in a label fails the test rather than being guessed).
- **Rationale:** GitHub hides comments, so readers see no change; the README's visible "This is `spec …`" sentence stays as prose.
- **Alternatives:** parse the visible sentence (rejected: prose wording varies, e.g. "This is `--human-review --explain` on the same example"); a fenced-block info string such as ```` ```spec-sample ```` (rejected: visible and unusual in rendered Markdown).

## R3 — Safety rules for labels

- First argument must be `status` or `inventory`.
- `--root` is required and its value must start with `testdata/` (no `..`), so samples never read outside fixtures.
- Allowed flags: `--color=never`, `--root`, `--status`, `--blocked`, `--human-review`, `--active`, `--historical`, `--explain`, `--fail-on`. Anything else (including `--json`, which has no README sample yet) fails the test as unsupported until a sample needs it.
- The block is compared with stdout; if the run is fatal (exit 2 with stderr), the block is compared with stderr.

## R4 — Unlabeled output detection

- A fenced block (```` ``` ```` with no info string) whose first line starts with `╭` (the panel's top-left corner) or `spec: ` is output-looking; without a preceding label it fails. Blocks with an info string (`sh`, `yaml`, `text`) and the agent-prompt block are never output-looking.
- Measured on the current README: exactly three output-looking blocks (lines 85, 158, 218), all of which get labels.

## R5 — Stale sample

- The "Filters" block (`--human-review --explain`) lacks the `statuses implemented 1 planned 1 verified 1` panel line added in feature 008. The block is replaced from a real run and the test then guards it.

## R6 — Retiring the feature-006 test

- `TestReadmeInventorySample` (`internal/status/inventory_test.go`) checks one sample by its visible sentence. It is removed; the inventory block gets a label and is covered by the shared test.

## R7 — Constraints

Test code only: standard library, no network or subprocess, reads only `README.md` and `testdata/`. Runtime untouched (Art. II). CLAUDE.md's "if output changes, update the README sample" note gains "`TestReadmeSamples` fails until you do".
