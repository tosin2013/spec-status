# Two-mode fixtures (feature 006)

Synthetic fixtures for roadmap recognition (`spec status`) and feature inventory (`spec inventory`). They are shaped like the evidence in `specs/001-spec-status/two-mode-direction.md`; no external repository content is copied.

Every expectation file is written by hand from the contracts in `specs/006-two-mode-inventory/contracts/`. Never regenerate one from `spec` output.

- `expected-status.json`: what `spec status` must report (exit code, fatal code and message, item count, whether the explicitly-empty line is shown). Read by `internal/status/recognition_test.go`.
- `expected-inventory.json`: the full `spec inventory --json` document. Read by `internal/status/inventory_test.go`.

| Fixture | Roadmap | Purpose |
|---|---|---|
| `unrecognized-empty` | 0 bytes | unrecognized, "Roadmap file is empty." |
| `unrecognized-whitespace` | spaces, tabs, blank lines | same |
| `unrecognized-narrative` | headings and prose | unrecognized |
| `planned-in-fence` | `## Planned Specs` only inside a fence | unrecognized |
| `planned-wrong-form` | `###`, lower case, indented and closed-hash forms | unrecognized |
| `explicit-empty` | template sections, `## Planned Specs`, grouping headings only | explicitly empty, exit 0 |
| `items-no-planned` | one item, no `## Planned Specs` | unchanged report |
| `narrative-with-label` | a `Depends on` label outside items | existing C-10 fatal, not unrecognized |
| `no-roadmap` | none | `SS_ROADMAP_UNAVAILABLE` suggesting inventory; main inventory fixture |
| `brownfield` | none | inventory: one feature folder 24/24, 16 checklist marks excluded |
| `no-specs` | none | inventory: no `specs/` folder |
| `with-roadmap` | one item | inventory ignores the roadmap |

Symlink, FIFO, sensitive-name and limit cases are created in temporary directories by the tests, because they cannot be committed portably.
