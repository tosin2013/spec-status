# Data model: Upstream-format conformance

These entities are documentation and test data. None of them is added to the product's runtime model or JSON schema.

## Upstream pin

| Field | Description | Validation |
|---|---|---|
| project | `github/spec-kit` or `srobroek/speckit-roadmap` | One current pin per project |
| version | Tag, for example `v1.2.0` | Required |
| commit | Full commit SHA when known (Spec Kit: `4ca46a0706a662b1edc07086ef16f021820c2db0`) | Required if the source records it |
| pinned_on | Date the owner approved the pin | ISO date |
| templates | Copied template files with paths | Each file is byte-identical to the upstream package copy |
| license | License file and NOTICE | Present for every pin |

**State**: `proposed` → `current` → `previous`. A pin becomes `current` only when its drift report is fully resolved (see Drift report) and the conformance suite passes. The old pin becomes `previous`; its shapes stay supported unless a recorded decision drops them.

## Shape

| Field | Description |
|---|---|
| key | Shape key from the classifier, for example `item-heading`, `field:Depends on`, `task-marker:[P]` |
| example | One literal line from a pinned template |
| interpretation | What spec status does with it: counted, linked, ignored as prose, and so on |
| pins | The pins whose templates produce it |
| status | `supported` (upstream-produced) or `strict` (not produced upstream; fails visibly or is a stated limit) |

The inventory itself is [contracts/shape-inventory.md](contracts/shape-inventory.md).

## Conformance fixture

| Field | Description | Validation |
|---|---|---|
| name | Directory name under `testdata/conformance/` (or an existing fixture directory) | Unique |
| origin | `template-filled`, `upstream-generated`, or `hand-edit` | Required |
| pins | Pins it covers | At least one current pin |
| files | `.specify/memory/roadmap.md` plus any linked `specs/*/` files | Root-relative, inside the fixture |
| expected | `expected.json`, a hand count | Format in [contracts/expected-result.md](contracts/expected-result.md) |

Rule: every `supported` shape is covered by at least one fixture (SC-002).

## Expected result

A hand-counted subset of facts: items with statuses and dependency states, supersession, the summary counts, task counts, diagnostic codes per item, and exit code. Format: [contracts/expected-result.md](contracts/expected-result.md).

## Drift report

| Field | Description |
|---|---|
| pin | The current pin compared against |
| candidate | Path and declared version of the candidate template |
| added / removed | Shape keys present in only one of the two |
| changed_labels | Field labels or task markers whose spelling changed |
| resolution | For each entry, a fixture name and a decision ID (support or reject) |

Rule: the pin can't move while any entry lacks both a fixture and a decision (FR-009).

## Bar decision

| Field | Value |
|---|---|
| id | `D-2026-11-07-09` |
| supersedes | `D-2026-11-07-06` (kept, marked superseded) |
| approver, date | Owner, 2026-11-07 |
| rule | Blocking means a wrong result on input that matches a current pin's shape inventory, or a silently misread `Depends on`, `Spec dir`, or `superseded-by` in any layout. Other exotic input fails visibly or is a README limit. |
| reclassified | The findings from acceptance reviews 1–6 (and 7 when it completes), each marked fixed, visible failure, limit, or blocking |
