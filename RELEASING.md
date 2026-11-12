# Releasing spec-status

Only the owner creates release tags and publishes releases (decision D-2026-11-11-06). The workflow `.github/workflows/release.yml` builds, tests and drafts a release; it never publishes. The design is in `specs/004-release-engineering/`.

## Versions

- `v0.1.0` was published directly, outside the workflow, and then withdrawn (D-2026-11-12-02). It is never reused. `v0.1.1` was tagged and drafted but never published (D-2026-11-12-03). The first release is `v0.1.2`.
- Repository protections: only admins can create, move or delete `v*` tags (tag ruleset); `main` requires the four CI test jobs and `hygiene`; published releases are immutable.

- Tags are `vX.Y.Z`, or `vX.Y.Z-rc.N` for release candidates.
- Before v1.0.0, a minor version may break compatibility (flags, exit codes, JSON schema). The notes must list every such change.
- A stable version follows a release candidate that the owner approved. The first stable version is `v0.1.0`.

## 1. Before tagging

- [ ] `main` CI is green: four test jobs, `hygiene`, and CodeQL.
- [ ] `spec status --fail-on error` passes on this repository, and the roadmap and README reflect what ships.
- [ ] No open Dependabot security PR or alert is waiting that should be in this release.

## 2. Tag

```sh
git fetch origin && git checkout origin/main
git tag vX.Y.Z-rc.N          # or vX.Y.Z
git push origin vX.Y.Z-rc.N
```

The tag must point at a commit on `main`, or the workflow stops at `verify`.

## 3. Watch the run

Actions → **release** runs these jobs in order:

1. `verify`
2. `ci`: the same checks as on pull requests
3. `build archives`
4. `smoke (ubuntu-latest)` and `smoke (macos-latest)`
5. `create draft release`

If any job fails, no draft is created.

## 4. Review the draft

Releases → the draft (a release candidate is marked pre-release):

- [ ] Exactly three assets: `spec_<v>_linux_amd64.tar.gz`, `spec_<v>_darwin_arm64.tar.gz`, `spec_<v>_checksums.txt`.
- [ ] Download them, then run `sha256sum -c spec_<v>_checksums.txt` (on macOS, `shasum -a 256 -c`). Both lines should say OK.
- [ ] Extract your platform's archive and run `spec version`. It should show `spec <v>`, the tagged commit, and `label: release build`.
- [ ] Read the notes. Check the commit, Go version, smoke-test links, the compatibility statement, and the "Not provided" line. Add any breaking changes, and edit wording if needed.

## 5. Publish

Click **Publish release**. After publishing, never edit or replace its assets. A mistake ships as a new version.

## 6. After publishing

- [ ] Check `go install github.com/tosin2013/spec-status/cmd/spec@v<v>`. `spec version` should show the module-build label. The module proxy can take a few minutes to see a new tag.
- [ ] For the first stable release, on each platform whose smoke test passed: update the README "Install" section from "chosen but not yet proven" to supported for that version, citing the smoke-run links (FR-016).
- [ ] Add an implementation-record row in `specs/001-spec-status/tasks.md` with the run and release URLs, and update the roadmap.

## Recovery

- **Bad draft, not published:** delete the draft release and the tag (`git push origin :refs/tags/vX.Y.Z-rc.N`). Then fix `main` and tag the next candidate (`-rc.N+1`). Never move a tag to a different commit.
- **Re-pushing a tag that already has a release:** `verify` fails on purpose.
- **Problem found after publishing:** ship a new patch version. Do not replace assets.
