# Quickstart: validating two-mode roadmap and feature inventory

Run from the repository root after implementation. Fixture paths are planned names; the hand-counted `expected.json` files are written before the code (CLAUDE.md: never regenerate an expected file from `spec` output).

## Build and test

```
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build -o spec ./cmd/spec
```

Expected: all packages pass, including the import/side-effect test, the extension reference test (new codes listed in `extensions/spec-status/reference/output.md`), and the conformance fixtures.

## 1. Unchanged reports (SC-002)

```
./spec status --root testdata/example-roadmap --json > /tmp/after.json
git worktree add /tmp/before main && (cd /tmp/before && CGO_ENABLED=0 go build -o /tmp/spec-before ./cmd/spec)
/tmp/spec-before status --root testdata/example-roadmap --json | cmp - /tmp/after.json
git worktree remove /tmp/before
./spec status --color=never        # this repository: same output as before
```

Expected: `cmp` prints nothing; this repository's report has its 16 items and the same summary.

## 2. Roadmap recognition (US1)

| Command | Expected |
|---|---|
| `./spec status --root testdata/two-mode/unrecognized-narrative` | stderr `spec: SS_ROADMAP_UNRECOGNIZED: …`, exit 2 |
| `./spec status --root testdata/two-mode/unrecognized-empty --json` | fatal document, `error.code` `SS_ROADMAP_UNRECOGNIZED`, message starts `Roadmap file is empty.`, exit 2 |
| `./spec status --root testdata/two-mode/explicit-empty --json` | `outcome` `report`, `total_items` 0, both percents `null`, exit 0 |
| `./spec status --root testdata/two-mode/explicit-empty --color=never` | panel contains `roadmap explicitly empty (Planned Specs heading, no items)` |
| `./spec status --root testdata/two-mode/unrecognized-narrative --fail-on none` | exit 2 |
| `./spec status --root testdata/two-mode/no-roadmap` | `SS_ROADMAP_UNAVAILABLE`, message mentions `spec inventory`, exit 2 |

Full list: [contracts/roadmap-recognition.md](contracts/roadmap-recognition.md) "Fixtures".

## 3. Inventory (US2)

```
./spec inventory --root testdata/two-mode/no-roadmap --color=never
./spec inventory --root testdata/two-mode/no-roadmap --json
```

Expected (hand count in `testdata/two-mode/no-roadmap/expected.json`): `001-alpha` 3/5 with spec, plan, tasks present; `20261101-093000-beta` spec only; `notes` not a feature folder with an info finding; `checklists/requirements.md` marks not counted; summary 2 feature folders, 3/5; exit 0; `--fail-on info` exits 1.

Brownfield-shaped fixture (SC-006): `./spec inventory --root testdata/two-mode/brownfield --json` → one feature folder, tasks 24/24, the 16 checklist marks excluded.

`./spec inventory --root testdata/two-mode/no-specs` → `SS_SPECS_DIRECTORY_MISSING` (info), zero folders, null percent, exit 0.

`./spec inventory --roadmap x` → `SS_INVOCATION_INVALID`, exit 2; with `--json`, an inventory-shaped fatal document.

## 4. Safety and determinism (US3)

Symlink and FIFO cases are created at test time in a temporary directory (they cannot be committed portably). The Go tests check:

- a symlinked feature folder and a symlinked `tasks.md` → `skipped`, `SS_PATH_UNSAFE`, contents not read;
- a FIFO `tasks.md` → `skipped`, the run does not block;
- a sensitive-named child → `SS_REFERENCE_SENSITIVE_SKIPPED`, never stat'ed;
- 1025 children of `specs/` → `SS_INPUT_LIMIT`, exit 2, `folders: []`;
- two runs on the same input produce identical JSON and human bytes.

## 5. Docs

- README Limits no longer contains the zero-items line; the missing-roadmap row mentions `spec inventory`.
- `specs/001-spec-status/diagnostics.md` lists the three new codes; `json-contract.md` documents the explicit-empty meaning; `cli-contract.md` has the inventory command.
- `go test ./internal/extension` passes (output reference covers the new codes).
