---
description: "Tasks for feature 009: README sample test"
---

# Tasks: README sample test

**Input**: `specs/009-readme-sample-test/` (spec.md, plan.md, research.md, data-model.md, contracts/readme-samples.md, quickstart.md). Roadmap item 011, D-2026-11-14-01 (5), v0.3.0.

**Tests**: The feature is a test. It is written first and must fail before the README is labeled and corrected.

**Organization**: US1 (drift fails, P1) and US2 (every sample covered, P1) share one test file; US2's guards are added after US1's comparison works.

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup

- [X] T001 Branch from up-to-date `main` after PR #60 merges (or from its branch, opening the implementation PR against `main`); set `.specify/feature.json` to `specs/009-readme-sample-test`

---

## Phase 2: User Story 1 — A changed output fails until the README is updated (P1) 🎯 MVP

**Goal**: each labeled README sample is compared byte for byte with real output.

**Independent Test**: one changed character in a sample makes `TestReadmeSamples` fail naming the sample; restoring it passes.

- [X] T002 [US1] Write `internal/cli/readme_samples_test.go` with `TestReadmeSamples`: read `../../README.md`; find each line matching `<!-- spec-sample: … -->` (research R2); split the text after `spec-sample:` on spaces into arguments; take the next fenced block with no info string as the sample (body between fences, ending in a newline); run `Run(args, repoRoot, &stdout, &stderr)` with `withBuild(t, resolveBuild("", "0.0.0-dev", "unknown", nil, false))` and `t.Setenv("COLUMNS", "80")`; compare the body with stdout (or stderr when exit is 2); on mismatch fail with the label, README line, and both texts. It must fail now (no labels yet)
- [X] T003 [US1] Add the three labels to `README.md` exactly as listed in contracts/readme-samples.md "Current labels", each on its own line immediately before its fenced block (Example, Filters, Without a roadmap)
- [X] T004 [US1] Replace the stale "Filters" block in `README.md` with the output of `spec status --color=never --root testdata/example-roadmap --human-review --explain` (development build, `COLUMNS=80`), checking line by line that the only difference from the old block is the `statuses implemented 1 planned 1 verified 1` panel line
- [X] T005 [US1] Run `CGO_ENABLED=0 go test ./internal/cli -run TestReadmeSamples -v`; three samples pass

---

## Phase 3: User Story 2 — Every output sample is covered, not just one (P1)

**Goal**: unlabeled output-looking blocks and unsafe or unsupported labels fail.

**Independent Test**: an unlabeled block starting with `╭` fails the test; an unsupported label fails.

- [X] T006 [US2] Extend `TestReadmeSamples` in `internal/cli/readme_samples_test.go` per contracts/readme-samples.md "Rules": fail on a fenced block with no info string whose first line starts with `╭` or `spec: ` and no preceding label; on a label with no following fenced block; on a first argument other than `status` or `inventory`; on a missing `--root` or a root not starting with `testdata/` or containing `..`; on a flag outside research R3's allowed list; on a quote character in a label; and when zero labels are found
- [X] T007 [US2] Add `TestReadmeSampleGuards` in `internal/cli/readme_samples_test.go` that runs the same parser on small in-memory Markdown strings (not the README) for each failure rule in T006 and asserts each is reported, plus one valid case that passes
- [X] T008 [US2] Remove `TestReadmeInventorySample` from `internal/status/inventory_test.go` (its sample is now labeled and covered by `TestReadmeSamples`); keep the rest of that file unchanged
- [X] T009 [US2] Run `CGO_ENABLED=0 go test ./...`; all pass

---

## Phase 4: Polish

- [X] T010 [P] Update the README-sample note in `CLAUDE.md` ("if output changes, update the README sample") to say `TestReadmeSamples` (`internal/cli/readme_samples_test.go`) fails until the labeled samples match
- [X] T011 Run the quickstart drift and coverage checks by hand (change one character in a sample; add an unlabeled `╭` block) and confirm each fails with a clear message, then restore `README.md`; record the result in the PR description
- [X] T012 Add a row for feature 009 to the "Implementation record" table in `specs/001-spec-status/tasks.md` (implemented, executed by the implementer, not independently verified)
- [X] T013 Doomer review on the change (prepare before editing, submit, answer every Quick lead, finish) before handing off the PR
- [X] T014 Propose roadmap 011 `specced` → `implemented` through `/speckit.roadmap.write` after the owner approves; `verified` after the v0.3.0 release checks

---

## Dependencies and execution order

- T001 → T002 → T003 → T004 → T005 (US1); T006 → T007 → T008 → T009 (US2, same test file, after US1).
- T010 any time; T011–T012 after T009; T013 before handoff; T014 last.

## Implementation strategy

US1 alone already guards the three samples; US2 closes the "unlabeled sample" gap that let the Filters block go stale. Both land in one PR.
