# Pinned upstream templates

These directories hold the upstream templates that define the input spec status supports (feature 002):

| Project | Current pin | Files |
|---|---|---|
| github/spec-kit | `v1.2.0` (commit `4ca46a0706a662b1edc07086ef16f021820c2db0`) | `spec-kit/v1.2.0/` |
| srobroek/speckit-roadmap | `v0.1.0` | `speckit-roadmap/v0.1.0/` |

Each pin keeps its upstream LICENSE and a NOTICE with the source, tag, commit, and file checksums. The supported shapes are listed in `specs/002-upstream-format-conformance/contracts/shape-inventory.md`.

The pins are listed once, in `pins.json`: project name, GitHub repository, pinned tag, and each template's repository path and pinned copy. The drift test and the weekly workflow both read it, and `TestPinsFile` checks it against these directories.

## Weekly drift check (feature 013)

`.github/workflows/drift.yml` runs every Monday at 06:17 UTC. For each project it takes the latest published, non-prerelease release, downloads the pinned templates at that release into a temporary directory, and runs the drift check below. It has read-only permissions, never writes to the repository, opens no issues or pull requests, and never moves a pin.

- **Start it by hand:** `gh workflow run drift.yml` (all projects), or `gh workflow run drift.yml -f project=spec-kit -f tag=v1.4.3`.
- **Run the same check locally:** `scripts/check-upstream-drift.sh [--project NAME] [--tag TAG]`. It needs `go`, `jq`, `curl` and a logged-in `gh`. `--candidate NAME=FILE` checks a file you already downloaded.
- **Reporting:** a failed run is the alert. GitHub notifies the owner, and the run summary lists each project's pinned tag, candidate tag and commit, and result.

| Result | Meaning | Run |
|---|---|---|
| same tag as the pin | The latest release is the pinned one | passes |
| no unresolved change | A newer release, and its template shapes match the pin (or every change is in `resolutions.md`) | passes |
| N unresolved change(s) | Shapes were added or removed; the summary lists each key with an example line. Start "Moving a pin" below | fails |
| template not found: PATH | The template is not at that path in the release (renamed or removed upstream) | fails |
| fetch error: … | The release, commit or file could not be fetched. Nothing was compared, so no drift is claimed | fails |
| check error: … | The drift test itself did not run, or did not report a result (for example, it was renamed) | fails |

GitHub pauses scheduled workflows in a repository with no activity for 60 days. Re-enable it in the Actions tab or with `gh workflow enable drift.yml`.

## Moving a pin

The drift check is development tooling. It reads only files you have downloaded; spec status itself never fetches or reads upstream content (constitution Article II, FR-011).

1. Download the candidate upstream release by hand. Record its tag and commit.
2. Run the drift check against each template that changed:

   ```
   CGO_ENABLED=0 go test ./internal/conformance -run TestDrift$ -v -args \
     -project=spec-kit -candidate=/path/to/tasks-template.md
   ```

3. For each reported change (`added: KEY` or `removed: KEY`):
   - add or extend a fixture under `testdata/conformance/` with a hand-counted `expected.json`;
   - record a decision (support it, or reject it with a visible failure or a README limit) in `specs/001-spec-status/decisions.md`;
   - add a row to `PROJECT/resolutions.md` with the change exactly as reported, the fixture name, and the decision ID.
4. When the drift check passes and `CGO_ENABLED=0 go test ./...` passes, copy the candidate templates into `PROJECT/TAG/templates/` with LICENSE and NOTICE, update the project's `tag` and `pinned` paths in `pins.json` (for a new project, also add its name to the `project` options in `.github/workflows/drift.yml`), update the shape inventory, and update the "Supported input" section of README.md.
5. Shapes supported under the previous pin stay supported unless a recorded decision removes them (FR-010). Note any removal for users.
