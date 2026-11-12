# Tasks: Release engineering

**Input**: Design documents from `specs/004-release-engineering/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/), [quickstart.md](quickstart.md)

**Tests**: Required by FR-012 for the build-identity change: Go tests in `internal/cli` and `internal/status`, written before the code they cover. The workflow is validated by real runs (quickstart §2–§4).

**Owner-only steps**: Pushing release tags and publishing releases are the owner's actions (D-2026-11-11-06). Tasks marked **(owner)** are prepared and checked by the implementer, but performed by the owner. Throwaway negative-test tags are pushed only with the owner's approval.

**Story order note**: The release workflow's smoke test (US1) checks the release label in `spec version` and `--help`. So the labeling story US4 (P2) is built before US1 (P1). US4 is small, and the order is a dependency, not a change of priority.

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup

- [X] T001 Branch from up-to-date `main` after the 004-plan PR merges (or continue on `004-plan`), and point `.specify/feature.json` at `specs/004-release-engineering`. Build `spec-before` from `main` into the scratchpad with `CGO_ENABLED=0` for byte comparisons.
- [X] T002 Re-check, and record in the PR description, that the pins in research.md R5 and the plan are still the latest. If any moved, update research.md and use the new values.
  - `actions/checkout` v7.0.1 `3d3c42e5aac5ba805825da76410c181273ba90b1`
  - `actions/setup-go` v7.0.0 `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e`
  - `actions/upload-artifact` v7.0.2 `cf430e030ddbb5b0abf93d22962f4752f3646cd9`
  - `actions/download-artifact` v8.0.2 `9000827ccba6bdab643e8b6fd33ac0654aef8333`

---

## Phase 2: Foundational: build-identity resolver (blocks US1, US3, US4)

**Purpose**: One function that decides release, module or development ([contracts/build-identity.md](contracts/build-identity.md) table; research R3).

- [X] T003 Write `internal/cli/build_test.go` first. It tests a pure resolver `resolveBuild(kind, version, commit string, info *debug.BuildInfo, ok bool) buildIdentity` with these cases:
  - **Release:** `kind=="release"` gives release, using the stamped version and commit, even when `info` has a module sum.
  - **Module:** `info.Main.Sum=="h1:…"` and `Main.Version=="v0.1.0"` give module, version `0.1.0`, commit `unknown`.
  - **Pseudo-version module:** `Main.Sum` set and `Main.Version=="v0.0.0-20261112040751-b228f6020c7b"` give module with that version without the `v`.
  - **Checkout build:** `Main.Sum==""` with `Main.Version=="v0.0.0-…+dirty"` gives development.
  - **(devel):** `Main.Version=="(devel)"` gives development.
  - **No build info:** `ok==false` gives development.
  - **Development defaults:** version `0.0.0-dev`, commit `unknown`.
- [X] T004 Create `internal/cli/build.go` (standard library only: `runtime/debug`, `strings`):
  - **Stamp variable:** declare `var Kind = ""` beside `Version` and `Commit` in `internal/cli/cli.go`, and document it as set only by `-ldflags -X …cli.Kind=release`.
  - **Resolver:** implement `resolveBuild` and `currentBuild()`, which calls `debug.ReadBuildInfo()` once. The `buildIdentity` fields are `kind`, `version`, `commit`, `label`, `heading` and `helpSentence`, with the exact strings from the contract.
  - **Test:** make T003 pass.
- [X] T005 Confirm the `internal/status/imports_test.go` allowlist still passes with `runtime/debug` imported from `internal/cli`, and that `CGO_ENABLED=0 go test ./...` passes.

**Checkpoint**: The resolver exists and is tested. No output has changed yet.

---

## Phase 3: User Story 4 — Release builds stop calling themselves development builds (P2, built first)

**Goal**: `spec version`, `--help` and the status panel heading follow the build identity. Development output is byte-identical to today's.

**Independent test**: Build one binary with the release stamps and one without. Only the label lines and the heading differ, and `--json` output is byte-identical (quickstart §1).

- [X] T006 [P] [US4] Add tests to `internal/status/style_test.go`:
  - **Heading:** `View{Build: "release 0.1.0"}` renders `Spec Status  release 0.1.0`.
  - **Default:** `View{}` and `View{Build: "development build"}` render `Spec Status  development build` unchanged.
  - **JSON:** `JSON(rep)` bytes don't depend on `View`.
- [X] T007 [P] [US4] Add tests to `internal/cli/cli_test.go` for the release, module and development identities. Each test injects an identity through a package-level override used only in tests. Check:
  - `versionText()` gives the exact five lines per kind from the contract;
  - `--help` contains the kind's sentence and still contains "It does not change files, use the network, or infer quarters.";
  - the development output equals today's text byte for byte.
- [X] T008 [US4] In `internal/status/style.go`, add the field `Build string` to `View`. In `summaryLines`, replace the hard-coded `"  development build"` with `"  " + build`, where `build` is `"development build"` when `View.Build == ""`. Make T006 pass.
- [X] T009 [US4] In `internal/cli/cli.go`:
  - `versionText()` prints the identity's version, commit and label;
  - `rootLong` becomes a function, or a template filled at command construction, whose development sentence stays "This is a development build.";
  - the `status.View{…}` literal at the view construction sets `Build: currentBuild().heading`.
  Make T007 pass.
- [X] T010 [US4] Byte checks:
  - Build `spec-dev` (no stamps) and compare `version`, `--help`, and `status --color=never` plus `--json` on `testdata/example-roadmap` and every `testdata/conformance/*` fixture against `spec-before`. Everything must be identical.
  - Build `spec-rel` with the three `-X` stamps and compare `--json` against `spec-dev` on the same fixtures. Must be identical (FR-011, SC-005).

**Checkpoint**: Labels are truthful and the README sample is unchanged.

---

## Phase 4: User Story 1 — The owner cuts a release safely (P1) 🎯 MVP

**Goal**: A tag produces a tested draft release with three assets ([contracts/release-workflow.md](contracts/release-workflow.md), [contracts/release-assets.md](contracts/release-assets.md)).

**Independent test**: The owner pushes `v0.1.0-rc.1`. All jobs pass, and a draft pre-release with exactly three assets appears. Then the negative runs in quickstart §3 behave as specified.

- [X] T011 [US1] In `.github/workflows/ci.yml`, add `workflow_call:` under `on:`, with no other change (research R1). The existing `pull_request` and `push: branches: [main]` triggers keep working.
- [X] T012 [US1] Create `.github/workflows/release.yml` with the header from the contract:
  - **Header:** `name: release`; `on: push: tags: ['v*']`; workflow-level `permissions: contents: read`; a comment saying the workflow never publishes.
  - **Job `verify`:**
    - checkout with `fetch-depth: 0` and `persist-credentials: false`;
    - check the tag against `^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-rc\.(0|[1-9][0-9]*))?$`;
    - run `git merge-base --is-ancestor "$GITHUB_SHA" origin/main`;
    - fail if `gh release view "$GITHUB_REF_NAME"` succeeds, with `GH_TOKEN: ${{ github.token }}`;
    - outputs `version` (the tag without `v`) and `prerelease` (`true` when the tag contains `-rc.`).
- [X] T013 [US1] In `.github/workflows/release.yml`, add job `ci`: `needs: verify`, `uses: ./.github/workflows/ci.yml`.
- [X] T014 [US1] In `.github/workflows/release.yml`, add job `build` (`needs: [verify, ci]`, `ubuntu-latest`):
  - **Setup:** pinned checkout, then setup-go with `go-version: stable` and `check-latest: true`.
  - **Build each target** `linux/amd64` and `darwin/arm64`: `CGO_ENABLED=0 GOOS=… GOARCH=… go build -trimpath -ldflags "-s -w -X github.com/tosin2013/spec-status/internal/cli.Version=${version} -X github.com/tosin2013/spec-status/internal/cli.Commit=${GITHUB_SHA} -X github.com/tosin2013/spec-status/internal/cli.Kind=release"`.
  - **Package:** directory `spec_${version}_${os}_${arch}/` with `spec` (0755), `LICENSE` and `README.md`, then `tar -czf spec_${version}_${os}_${arch}.tar.gz` of it.
  - **Checksums:** `sha256sum spec_${version}_*.tar.gz | sort -k2 > spec_${version}_checksums.txt`.
  - **Outputs:** `go version` as output `go_version`.
  - **Upload:** the three files with `actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9 # v7.0.2`, artifact name `release-files`.
- [X] T015 [US1] In `.github/workflows/release.yml`, add job `smoke` (`needs: [verify, build]`, matrix `os: [ubuntu-latest, macos-latest]`, `fail-fast: false`):
  - **Platform:** assert `linux amd64` or `darwin arm64`.
  - **Download:** `actions/download-artifact@9000827ccba6bdab643e8b6fd33ac0654aef8333 # v8.0.2`.
  - **Checksum:** `sha256sum -c --ignore-missing` on Linux, `shasum -a 256 -c --ignore-missing` on macOS.
  - **Extract and run** only this platform's archive (research R6):
    - the file list is exactly `spec`, `LICENSE`, `README.md`;
    - `spec version` contains `spec ${version}`, `commit: ${GITHUB_SHA}` and `label: release build`;
    - `spec --help` contains `This is release ${version}.` and does not contain `development build`;
    - `spec status --color=never --root testdata/example-roadmap` exits with the expected code; checkout supplies the testdata;
    - `spec status --json` with the same root parses as JSON with `schema_version` `1.0`, checked with `python3 -c`;
    - `spec status --root /nonexistent` exits 2.
- [X] T016 [US1] In `.github/workflows/release.yml`, add job `release` (`needs: [verify, build, smoke]`, `ubuntu-latest`, job-level `permissions: contents: write`):
  - download the artifact;
  - re-check that `gh release view` fails;
  - render `notes.md` from the template in contracts/release-assets.md, filling in the version, commit, `needs.build.outputs.go_version` and the smoke job URLs (from `gh api repos/$GITHUB_REPOSITORY/actions/runs/$GITHUB_RUN_ID/jobs`), with the `(release candidate)` heading for `-rc.N`;
  - run `gh release create "$GITHUB_REF_NAME" spec_${version}_* --draft --verify-tag --title "$GITHUB_REF_NAME" --generate-notes --notes-file notes.md`, plus `--prerelease` when needed, with `GH_TOKEN: ${{ github.token }}`.
- [X] T017 [US1] Check `.github/workflows/release.yml` and `ci.yml` against the contracts:
  - YAML parses;
  - every `uses:` has a 40-hex SHA and a version comment;
  - no `secrets.`, `pull_request_target`, `continue-on-error`, `workflow_dispatch` or `if: always()`;
  - only job `release` has `contents: write`;
  - `gh release` appears only in `verify` (`view`) and `release` (`view`, `create`).
- [X] T018 [US1] Dry-run the build and package steps locally in the scratchpad for `0.0.0-local`. Run the same commands as T014 for both targets, then `sha256sum -c` on the result, and run the darwin/arm64 binary's `version` and `--help` checks from T015 locally.
- [ ] T019 [US1] **(owner)** After this PR merges, the owner pushes `v0.1.0-rc.1` on `main`. Watch the run, then record in the PR or implementation record:
  - the run URL;
  - that the draft pre-release has exactly three assets;
  - that the notes meet FR-015;
  - that `sha256sum -c` passes on the downloaded assets.
- [ ] T020 [US1] **(with owner approval)** Negative runs from quickstart §3, each recorded with its run URL, and each throwaway tag and draft deleted afterwards:
  - bad tag name `v0.1`: `verify` fails;
  - tag on a commit not on `main`: `verify` fails;
  - re-push of the existing `v0.1.0-rc.1` tag name: `verify` fails because the release exists.

**Checkpoint**: Release candidates work end to end, and failures create nothing.

---

## Phase 5: User Story 2 — A user installs without Go (P1)

**Goal**: The README install section meets FR-013.

**Independent test**: On clean linux/amd64 and darwin/arm64 machines, without Go and without root, follow the README for the published release in under 5 minutes (SC-003).

- [X] T021 [P] [US2] Add an "Install" section to `README.md`, above "Build and run from source":
  - **Archives:** the per-platform archive names using `<version>` and the GitHub Releases URL pattern `https://github.com/tosin2013/spec-status/releases/download/v<version>/<file>`.
  - **Download:** a `curl -fLO` example for the archive and the checksum file.
  - **Verify:** `sha256sum -c --ignore-missing` on Linux, `shasum -a 256 -c --ignore-missing` on macOS.
  - **Install:** `tar -xzf`, then move `spec` into `~/.local/bin` or another user-owned PATH directory. No root and no Go needed.
  - **Upgrade and rollback:** repeat the steps with another version. Old releases stay downloadable.
  - **macOS Gatekeeper:** browser downloads are quarantined; use `xattr -d com.apple.quarantine ./spec` for that file only. Binaries are not signed or notarized, and the README never suggests disabling Gatekeeper.
  - **Unsupported platforms:** Windows and linux/arm64 have no archive; building from source may work but is untested.
- [X] T022 [P] [US2] Add the pre-1.0 compatibility policy to the same section of `README.md` (FR-014): before v1.0.0, a minor version may change flags, exit codes or the JSON schema, and the notes list such changes; patch versions don't break compatibility.

**Checkpoint**: Users have a complete install guide.

---

## Phase 6: User Story 3 — `go install` reports the real version (P2)

**Goal**: The module path is documented, and the module-build label is verified on a real tag.

**Independent test**: After `v0.1.0-rc.1` exists, `go install …@v0.1.0-rc.1` in a clean `GOBIN` reports `spec 0.1.0-rc.1`, `commit: unknown` and the module-build label (SC-004).

- [X] T023 [P] [US3] Add a "With Go" subsection to `README.md` install: `go install github.com/tosin2013/spec-status/cmd/spec@v<version>` (needs Go 1.26+). Explain that `spec version` then reports a module build, not a release archive, with commit `unknown`.
- [ ] T024 [US3] After T019, run `GOBIN=$(mktemp -d) GOFLAGS=-modcacherw go install github.com/tosin2013/spec-status/cmd/spec@v0.1.0-rc.1` and record the `spec version` output. It must show module kind, version `0.1.0-rc.1` and `commit: unknown`. If the module proxy hasn't indexed the tag yet, retry later, and record the wait.

**Checkpoint**: The `go install` channel is truthful.

---

## Phase 7: User Story 5 — Each release is traceable and documented (P3)

**Goal**: Notes, the owner checklist and the policy docs are complete.

**Independent test**: The `v0.1.0-rc.1` draft notes contain every FR-015 item and no forbidden claim, and `RELEASING.md` covers tag to publish to post-publish.

- [X] T025 [P] [US5] Create `RELEASING.md` (FR-017), the owner checklist:
  - **Before tagging:** `main` CI green, roadmap and README updated.
  - **Tagging:** `git tag vX.Y.Z[-rc.N] <commit-on-main> && git push origin <tag>`.
  - **Watch the run.**
  - **Review the draft:** three assets, notes, `sha256sum -c`, and a manual `spec version` on a downloaded archive.
  - **Edit the notes** if needed, then publish.
  - **Post-publish:**
    - update README "supported" wording per FR-016;
    - mark roadmap 004 when v0.1.0 ships;
    - add the implementation-record row.
  - **Recovery:** delete a bad draft and its tag before publish; after publish, never replace assets, and ship a new version instead.
- [X] T026 [P] [US5] Link `RELEASING.md` from `CONTRIBUTING.md`, in one sentence under "Pull requests" or a new "Releases" line.
- [X] T027 [P] [US5] Update `specs/001-spec-status/build-release.md`:
  - the release workflow exists (`.github/workflows/release.yml`; tag → CI → build → smoke → draft; owner publishes);
  - archive names and checksum format follow `specs/004-release-engineering/contracts/release-assets.md`;
  - remove "No release workflow …" from the header while keeping "no release published" until v0.1.0 is published;
  - the release Go version is current stable, recorded in the notes;
  - D-2026-11-11-06.

---

## Phase 8: Polish and cross-cutting

- [X] T028 Run the full local verification: `CGO_ENABLED=0 go test -count=1 ./...`, `go vet ./...`, `gofmt -l .` (empty), `go mod tidy -diff`, `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...`, `spec status --fail-on error` on the repo, and the T010 byte checks again.
- [X] T029 Run the Doomer review workflow on the source changes (`internal/cli`, `internal/status`, `.github/workflows`), and record a response for every Quick finding before push.
- [ ] T030 Commit, push and open the PR. Confirm feature 003's CI passes on it: four test jobs, hygiene and dependency review.
- [ ] T031 Add an implementation-record row to `specs/001-spec-status/tasks.md`: "Release engineering (feature 004)". Cover the build identity, the release workflow, the README install section and RELEASING.md, with run URLs once T019 and T020 exist. Mark roadmap 004 `implemented` in `.specify/memory/roadmap.md` after T019 succeeds; `verified` waits for the owner-published v0.1.0 and the US2 install check.
- [ ] T032 **(owner, optional)** Offer the repository settings from research.md R8: immutable releases, and a `v*` tag ruleset limiting tag creation and deletion to the owner. Also re-offer feature 003's branch protection. Apply none without explicit approval.
- [ ] T033 **(owner)** After the release candidate is reviewed, the owner tags and publishes `v0.1.0` (SC-007). The implementer then runs the quickstart §4 checks, updates the README support wording (FR-016), and records the results.

---

## Dependencies and execution order

- **Phase 1 → Phase 2** (T003 → T004 → T005) **→ US4** (T006 and T007 in parallel → T008 → T009 → T010) **→ US1** (T011 → T012 → T013 → T014 → T015 → T016 → T017 → T018).
- **US2, US3 docs (T021–T023) and US5 (T025–T027)** edit separate files, so they can run in parallel with US1 once Phase 2 is done. T021, T022 and T023 all edit `README.md`, so run them in sequence with each other.
- **Polish:** T028 → T029 → T030 come after all file tasks. Merging the PR is the owner's step.
- **After merge:** T019 (owner tag) → T020 (approved negatives) and T024 (`go install` check) → T031 records → T033 (owner publishes v0.1.0).

### Story completion order

Foundational → US4 → US1 (MVP) → US2, US3, US5 (docs, parallel) → post-merge validation.

## Parallel example

```text
After T005:
  T006 [US4] status heading tests     T007 [US4] cli identity tests
After T010, alongside the US1 workflow lane:
  T021→T022→T023 README (one lane)    T025 RELEASING.md    T026 CONTRIBUTING.md    T027 build-release.md
```

## Implementation strategy

1. **MVP**: Phase 2, US4 and US1. A tag yields a tested draft release with truthful labels.
2. **Docs**: US2, US3 and US5, in the same PR, so the first release candidate already has install instructions and a checklist.
3. **Validation after merge**: the owner's `v0.1.0-rc.1`, then the approved negative runs and the `go install` check, then the owner publishes `v0.1.0`.
4. **Stop conditions**:
   - T010: any development-output or JSON byte difference.
   - T017: a contract violation.
   - T019: a failed release-candidate job. Fix it and re-tag with `rc.2`; never move a tag.
