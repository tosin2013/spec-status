# Field-form fixtures (feature 008)

Synthetic roadmaps for roadmap items 015 and 016 (`specs/008-field-forms-and-fixes/contracts/changes.md` §6). Each `expected-status.json` is written by hand from the contract and the existing rules; never regenerate one from `spec` output. It lists the exit code, a fatal code if any, selected item fields, and the complete diagnostic set as (item, code) pairs. Read by `internal/status/field_forms_test.go`.

| Fixture | Case |
|---|---|
| `superseded-bold`, `superseded-plain`, `superseded-lower` | the two added spellings of `superseded-by` |
| `superseded-upper` | another case form: still fatal |
| `superseded-twice` | both spellings on one item: repeated field, fatal |
| `template-placeholders` | the pinned template's `Depends on` and `Spec dir` placeholder lines on a planned item |
| `placeholder-italic` | an italic `Spec dir` placeholder on a specced item: still needs a real spec link |
| `bracketed-real-id` | `Depends on: [002]` stays malformed; a real ID is never dropped |
| `duplicate-deps` | a repeated dependency ID listed once, with one warning |

The raw pinned template case is built at test time from `testdata/upstream/speckit-roadmap/v0.1.0/templates/roadmap-template.md`, so no second copy is committed.
