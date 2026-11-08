# Contract: the drift check

A development-only test. It never runs as part of `spec`, and it never fetches anything.

## Invocation

```
go test ./internal/conformance -run TestDrift -args \
  -project=spec-kit -candidate=/path/to/downloaded/tasks-template.md
```

Without `-candidate` the test is skipped. `TestPinnedInventory` always runs: it checks that every shape key extracted from the pinned templates appears in [shape-inventory.md](shape-inventory.md) and is covered by a fixture.

## Report (test output)

```
drift: spec-kit pinned v1.2.0 (4ca46a07) vs candidate /path/…/tasks-template.md
added:   task-marker:[Pn]           line 54: - [ ] T003 [P1] Configure …
removed: front-matter
changed: field label "Spec dir" -> "Spec directory"
result:  3 unresolved changes; the pin cannot move
```

- Placeholder text (`[FEATURE NAME]`, `[Title]`, `T###` numbers) is normalized before comparison, so wording-only edits report nothing.
- Exit status: the test passes when there are no unresolved changes. Changes count as resolved when each one is listed in `testdata/upstream/<project>/resolutions.md` with a fixture name and a decision ID.

## Moving a pin

1. Download the candidate upstream release by hand and record its tag and commit.
2. Run the drift check against each relevant template.
3. For each change: add or extend a fixture with a hand-counted `expected.json`, record a decision (support, or reject with visible failure or a README limit), and list both in `resolutions.md`.
4. When the drift check passes and `go test ./...` passes, copy the candidate templates into `testdata/upstream/<project>/<tag>/` with LICENSE and NOTICE, mark the old pin `previous`, and update README's supported-versions line.
