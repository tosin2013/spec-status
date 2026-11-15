# Quickstart: validating the README sample test

```
CGO_ENABLED=0 go test ./internal/cli -run TestReadmeSamples -v
```

Expected: three samples checked (Example, Filters, Without a roadmap), all pass.

Drift check (SC-002): change one character inside any labeled block in `README.md` and rerun: the test fails, naming the label's command and the README line, with expected and actual text. Restore the file.

Coverage check (SC-003): add a fenced block whose first line is `╭───` with no label above it; rerun: the test fails naming that line. Remove it.

Full suite: `CGO_ENABLED=0 go test ./...` passes; `TestReadmeInventorySample` no longer exists.
