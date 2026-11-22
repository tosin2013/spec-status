# Contract: drift workflow and checker

## `scripts/check-upstream-drift.sh`

```
scripts/check-upstream-drift.sh [--project NAME] [--tag TAG] [--candidate NAME=FILE] [--path NAME=REPO_PATH]
```

- Default: every project in `testdata/upstream/pins.json`, each at its latest published non-prerelease release.
- `--tag` requires `--project`.
- `--candidate` skips the download for that project and reads a local file. It assumes the project has one template.
- **Needs:** `go`, plus `gh` (authenticated, or with `GH_TOKEN`), `jq` and `curl` unless every project uses `--candidate`.
- **Exit:** 0 if every project is `same tag as the pin` or `no unresolved change`. 1 otherwise, or on a usage error.
- **Writes:** only under a temporary directory, removed on exit, and to `$GITHUB_STEP_SUMMARY` if that is set.

**Summary format:**

```
## Upstream drift check

| Project | Pinned | Candidate | Commit | Result |
|---|---|---|---|---|
| spec-kit | v1.2.0 | v1.4.3 | 1a2b3c4 | no unresolved change |
| speckit-roadmap | v0.1.0 | v0.1.1 | 5d6e7f8 | 2 unresolved changes |

### speckit-roadmap: unresolved changes
(added:/removed: lines from TestDrift, in a code block)
```

Result values: `same tag as the pin`, `no unresolved change`, `N unresolved change(s)`, `template not found: PATH`, `fetch error: MESSAGE`.

## `.github/workflows/drift.yml`

- **Triggers:** weekly `schedule` (Monday 06:17 UTC) and `workflow_dispatch` (inputs: `project` = all|spec-kit|speckit-roadmap, `tag` = optional).
- **Permissions:** `contents: read` only. No other secrets.
- **Steps:** checkout, then setup-go (1.26.x), then the script with inputs passed as environment variables.
- **Never:** pushes, creates issues or pull requests, uploads artifacts, or moves pins.
- Scheduled runs pause after 60 days without repository activity. Re-enable with `gh workflow enable drift.yml`.
