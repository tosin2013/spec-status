# Quickstart: validating field forms and small defect fixes

Run from the repository root after implementation.

```
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build -o spec ./cmd/spec
```

## 1. Field forms (US1, US2)

```
for f in testdata/field-forms/*/; do ./spec status --root "$f" --json; echo "exit=$?"; done
```

Expected: each matches its `expected-status.json` (contracts/changes.md §6). In particular `template-placeholders` exits 0, and `bracketed-real-id` still reports `SS_DEPENDENCY_MALFORMED`.

Template check (SC-002): copy `testdata/upstream/speckit-roadmap/v0.1.0/templates/roadmap-template.md` into a scratch project's `.specify/memory/roadmap.md` and run `spec status`: exit 0, placeholder info findings, no errors from the template text.

## 2. Duplicates (US3)

`./spec status --root testdata/field-forms/duplicate-deps --json` → one entry per ID, one `SS_DEPENDENCY_DUPLICATE`.

## 3. Human text (US4)

```
./spec status --color=never --root testdata/example-roadmap
```

Expected: the panel has a `statuses` line; the README example matches this output exactly.

## 4. CLI (US5)

| Command | Expected |
|---|---|
| `./spec status --json=false --bogus` | plain text on stderr, exit 2 |
| `./spec status --root ""` | `SS_INVOCATION_INVALID: --root needs a non-empty value`, exit 2 |
| `./spec status --json --fail-on warning --bogus` | JSON with `"fail_on":"warning"`, exit 2 |

## 5. Nothing else changes (SC-003)

Compare JSON for every fixture root with the pre-change build (the 006/007 snapshot approach): identical. Human text differs only by the added panel lines and row notes.
