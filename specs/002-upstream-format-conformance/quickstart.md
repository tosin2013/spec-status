# Quickstart: validating upstream-format conformance

## Prerequisites

- Go 1.22 or newer, `CGO_ENABLED=0`
- Repository checkout. No network is needed for any step below.

## 1. Conformance suite (every change)

```
CGO_ENABLED=0 go test ./internal/conformance
```

Expected: every fixture under `testdata/conformance/`, plus `testdata/speckit-roadmap-upstream`, matches its hand-counted `expected.json` ([contract](contracts/expected-result.md)). A mismatch names the fixture and the differing fact.

## 2. Pinned inventory check

```
CGO_ENABLED=0 go test ./internal/conformance -run TestPinnedInventory -v
```

Expected: every shape extracted from `testdata/upstream/*/*/` appears in [shape-inventory.md](contracts/shape-inventory.md) and is covered by at least one fixture.

## 3. Drift check against a new upstream release

```
CGO_ENABLED=0 go test ./internal/conformance -run TestDrift -v -args \
  -project=spec-kit -candidate=$HOME/Downloads/spec-kit-vX/templates/tasks-template.md
```

Expected: a report of added, removed, and changed shapes ([contract](contracts/drift-report.md)). It passes only when each change has a fixture and a decision in `resolutions.md`.

## 4. Real-project hand count (SC-003)

```
CGO_ENABLED=0 go build -o spec ./cmd/spec
./spec status --color=never --root .                               # this repository
./spec status --color=never --root testdata/speckit-roadmap-upstream
```

Expected: the counts match the hand counts recorded in `specs/001-spec-status/tasks.md` (AC-24) and in each fixture's `expected.json`.

## 5. Acceptance under the new bar

An independent reviewer is given decision D-2026-11-07-09 and the shape inventory. They classify findings as blocking only when the input matches a supported shape, or a metric field (`Depends on`, `Spec dir`, `superseded-by`) is silently misread. Expected for SC-006: zero blocking findings.
