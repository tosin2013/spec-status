# Quickstart: validating the public Go package

```sh
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go test -race ./pkg/status
go doc ./pkg/status
```

1. **Byte identity:** `go test ./pkg/status -run TestMatchesCLI -v` lists every fixture root × 4 thresholds, all passing.
2. **Round trip:** `go test ./pkg/status -run TestRoundTrip` passes.
3. **Example:** `go test ./pkg/status -run Example -v` prints the example roadmap's summary as documented.
4. **CLI unchanged:** for every fixture root, text, `--explain`, `--json` and exit codes are byte-identical to the pre-change build.
5. **Guardrail:** temporarily adding an outside import (for example `github.com/spf13/pflag`) to `pkg/status` makes `TestProductCodeDoesNotImportSideEffectPackages` fail.
6. **After the release:** in a scratch module, run `go get github.com/tosin2013/spec-status/pkg/status@v0.4.0` and a 10-line program that loads `testdata/example-roadmap` from a checkout. Its output equals `spec status --json` for the same root.
