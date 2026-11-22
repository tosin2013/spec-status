# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`spec status` is a read-only Go CLI that reports on a Spec Kit roadmap (`.specify/memory/roadmap.md`): recorded statuses, task checkboxes, dependency blockers, supersession history, and broken references, as human text or JSON. Module: `github.com/tosin2013/spec-status`.

## Commands

```
CGO_ENABLED=0 go test ./...                                   # all tests
CGO_ENABLED=0 go test ./internal/status -run TestHistoryHumanReviewDoesNotRedirect   # single test
CGO_ENABLED=0 go build -o spec ./cmd/spec
./spec status --color=never --root testdata/example-roadmap   # human text against the fixture
./spec status --root testdata/example-roadmap --json
./spec status --color=never --root testdata/example-roadmap --human-review --explain
CGO_ENABLED=0 go test ./internal/conformance                  # hand-counted conformance fixtures + shape inventory
CGO_ENABLED=0 go test ./internal/conformance -run TestDrift$ -v -args -project=spec-kit -candidate=PATH   # upstream drift check
```

Supported input is pinned to the upstream templates in `testdata/upstream/` (github/spec-kit v1.2.0, srobroek/speckit-roadmap v0.1.0); see `specs/002-upstream-format-conformance/`. Every `expected.json` under `testdata/conformance/` is a hand count — never regenerate one from `spec` output. The acceptance bar is decision D-2026-11-07-09 in `specs/001-spec-status/decisions.md`.

Build identity is stamped with `-ldflags "-X github.com/tosin2013/spec-status/internal/cli.Version=... -X ...cli.Commit=..."` (defaults `0.0.0-dev` / `unknown`). Always build with `CGO_ENABLED=0`.

## Architecture

Pipeline: `cmd/spec/main.go` → `cli.Run` → `status.Evaluate` → `status.JSON` or `human.Text` → `status.ExitCode`. `status.JSON` is `schema.Marshal(toSchema(rep))` and `status.ExitCode` delegates to `schema.ExitCode`: `internal/schema` (standard library only) holds the schema 1.0 report types, shared with the public Go package `pkg/status` (`Load`, `Marshal`, `ExitCode`; type aliases; experimental until 1.0, see `docs/library.md`). `pkg/status/status_test.go` fails if the library's bytes or exit code differ from the CLI on any fixture root.

- `internal/cli` — the Cobra command tree (`status`, `version`, `completion`), flag validation (`precheck`, repeated-flag checks, custom `pflag.Value` types), color-mode decision, and invocation-error output. In `--json` mode even invocation errors are emitted as a single JSON document on stdout. It unsets `BASH_COMP_DEBUG_FILE` around each call so Cobra cannot write a log.
- `internal/status` — the core:
  - `parse.go`: narrow line-based roadmap and tasks.md parser (items are `### 001 — Title [status: planned]` headings followed by `- Field: value` bullets such as `Spec dir`, `Depends on`, `superseded-by`; fenced blocks are skipped). The grammar is in `specs/001-spec-status/roadmap-format.md`. Unsupported item-like structure is a fatal parse error, not a guess.
  - `analyze.go` (`build`): turns raw items into a `Report` — supersession chains/cycles, dependency classification (a dependency on a superseded item is always `human_review`, never redirected to the successor), spec/plan/tasks file inspection via `safeio`, metrics, and sorted diagnostics.
  - `model.go`: the `Report` model; `render.go`: the versioned JSON view (schema 1.0, via `toSchema`).
- `internal/human` — `style.go`: Lip Gloss human text (`human.Text`, `human.InventoryText`, `human.View`) — summary panel, table/stacked rows, filters, `--explain`, and plain fatal text. It is a separate package so `pkg/status` links no third-party module; its tests live beside it.
  - `run.go`: root discovery (nearest ancestor with a non-symlink `.specify`), roadmap resolution, mapping `safeio` error codes to `SS_*` diagnostics, and exit codes.
- `internal/safeio` — the only filesystem boundary. Root-confined, bounded reads; rejects `..`, symlink components, non-regular files (FIFOs), sensitive names (secrets/keys, not opened at all), and files that change mid-read. All roadmap/task reads must go through it.

Text and JSON are two views of one `Report`. Human-text filters (`--status`, `--blocked`, `--human-review`, `--active`, `--historical`) and `--explain` only select/annotate item rows; they must not change the summary panel, diagnostics, exit code, or JSON bytes.

Exit codes: `0` no finding at/above `--fail-on` (default `error`), `1` a finding reached it, `2` fatal/invocation failure. `--fail-on none` never masks a fatal.

## Hard constraints (enforced by tests and the constitution)

These come from `.specify/memory/constitution.md`; don't relax them without explicit owner approval.

- **No runtime side effects**: no network, subprocess/Git, writes, caches, or telemetry. `internal/status/imports_test.go` fails on imports of `os/exec`, `net/*`, `plugin`, and on strings like `os.WriteFile`, `os.Create`, `os.OpenFile`, `os.Remove`, `exec.Command` in non-test code.
- **Dependency allowlist**: only Cobra/pflag/`golang.org/x/sys/unix` may be imported from `internal/cli` and `cmd/`; only Lip Gloss/termenv/uniseg from `internal/human/style.go`; `pkg/status` and `internal/schema` are standard library only, and `TestLibraryDependencyClosureIsStandardLibrary` checks everything `pkg/status` links. Everything else (parsing, safeio, analysis) is standard library only. Bubble Tea, Glamour, `lipgloss/tree`, Huh, dashboards, and watch modes are explicitly banned.
- **Truthful evidence**: keep checked tasks, recorded `implemented`, and recorded `verified` separate; never infer verification. Unknown status tokens are preserved raw with a warning. Superseded items stay in history and get no active completion credit.
- **Determinism**: document order for items, stable diagnostic ordering (`diagLess`), no timestamps.
- **Plain surfaces**: JSON, fatal text, version, and completion are never colored. Color only when stdout is a TTY, `NO_COLOR` unset, and `TERM != dumb` (or `--color=always`). Width comes from `COLUMNS`, else 80; never query the terminal or open `/dev/tty`. Untrusted roadmap text is sanitized of terminal escapes in human output.
- **No inferred quarters** — no quarter flags, fields, or estimation.

## Specs and docs

Feature docs live in `specs/001-spec-status/` (spec, CLI/JSON/diagnostics contracts, roadmap grammar, decisions, tasks). `tasks.md` has an "Implementation record" table tracking what's implemented vs. owner-approved vs. not started; keep it and the README accurate when behavior changes. README and docs deliberately avoid overstating claims (e.g. tests are "executed by the implementer, not independently verified"; only Linux amd64 and Darwin arm64 were run) — match that tone.

Spec Kit scaffolding (`.specify/templates`, `.specify/scripts`, `.specify/extensions/roadmap`, and the `/speckit-*` / `speckit.roadmap.*` skills under `.claude/`) is development tooling only. The `spec` CLI must never depend on it at runtime (constitution Art. II). `/speckit-constitution` would overwrite the hand-written constitution, so don't run it without owner approval.

`testdata/example-roadmap` is the synthetic fixture used in README examples; if output changes, update the README sample. Every output block in README.md and `docs/*.md` is labeled `<!-- spec-sample: <command> -->`, and `TestReadmeSamples` (`internal/cli/readme_samples_test.go`) fails until each labeled block matches real output byte for byte; an output-looking block without a label also fails. The README stays short: usage detail is in `docs/usage.md`, the full limits in `docs/limits.md`, build and project documents in `docs/development.md`, the system overview in `docs/architecture.md`, and the agent prompt only in `docs/agent-prompt.md`, and the Go library in `docs/library.md`.
