# Quickstart: validating per-phase task counts

```sh
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build -o spec ./cmd/spec
```

1. **Real Spec Kit layout:**

   ```sh
   ./spec status --json --root testdata/conformance/upstream-generated-002-tasks | jq '.items[0].tasks.phases'
   ```

   Expected: seven phases in file order (lines 33, 44, 58, 79, 96, 112, 129), `story` 1–4 for the phases at lines 58–112 and `null` otherwise, with counts as hand-counted in that fixture's `expected.json`.
2. **Edge cases:**

   ```sh
   ./spec status --json --root testdata/phases | jq '.items[] | {id, state: .tasks.state, phases: .tasks.phases}'
   ```

   Expected: as listed in `testdata/phases/README.md` (`[]` for the file without phases, `null` for the missing file, the fenced `## Phase` absent).
3. **Nothing else changed:** for every fixture root, the text output is identical to the pre-change build, and the JSON with `del(.items[].tasks.phases)` is identical to the pre-change JSON.
4. **Docs:** `TestReadmeSamples` passes, so the README and `docs/` output samples still match.
