# Quickstart: validating status/tasks consistency diagnostics

Run from the repository root after implementation.

## Build and test

```
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build -o spec ./cmd/spec
```

Expected: all packages pass, including the conformance test (with the hand-edited `template-filled` expectation) and the extension reference test (codes listed in `reference/output.md`).

## 1. Ahead of tasks (US1)

```
./spec status --root testdata/skills/implemented-unchecked --json          # exit 0
./spec status --root testdata/skills/implemented-unchecked --fail-on warning # exit 1
```

Expected: one `SS_STATUS_AHEAD_OF_TASKS` (warning) for 001, source at the item heading, `related` pointing at the first unchecked line of `specs/001-import/tasks.md`.

## 2. Tasks done, status early (US2)

```
./spec status --root testdata/skills/tasks-done-status-planned --json      # exit 0
./spec status --root testdata/skills/tasks-done-status-planned --fail-on info # exit 1
```

Expected: one `SS_TASKS_DONE_STATUS_EARLY` (info) for 001.

## 3. Nothing else changes (SC-002, SC-003)

Capture outputs for every fixture root and this repository with the pre-change binary, then with the new one, and diff (the 006 scratch `snap.sh` approach). Expected: only `template-filled`, `implemented-unchecked`, `roles-mixed` and `tasks-done-status-planned` differ, only by the added diagnostics; no default-threshold exit code changes; this repository is byte-identical.

## 4. Docs and agents

- `specs/001-spec-status/diagnostics.md` and `extensions/spec-status/reference/output.md` list both codes; `go test ./internal/extension` passes.
- The reconcile command cites the codes and still derives the cases from counts for older `spec` builds.
