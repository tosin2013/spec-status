# Contract: drift workflow and checker

## `scripts/check-upstream-drift.sh`

```
scripts/check-upstream-drift.sh [--project NAME] [--tag TAG] [--candidate NAME=FILE] [--path NAME=REPO_PATH]
```

- Default: every project in `testdata/upstream/pins.json`, each at its latest published non-prerelease release.
- `--tag` requires `--project`. A tag may use letters, digits, `.`, `_`, `-` and `/`, but not `..`, `//`, or a leading or trailing `/`.
- `--candidate` skips the download for that project and reads a local file. `--candidate` and `--path` must name a selected project that has exactly one template, and `--candidate` cannot be combined with `--tag`; anything else is a usage error rather than a silently ignored option.
- **Needs:** `go` and `jq` (pins.json is read with jq), plus `gh` (authenticated, or with `GH_TOKEN`) and `curl` unless every selected project uses `--candidate`.
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

Result values: `same tag as the pin`, `no unresolved change`, `N unresolved change(s)`, `template not found: PATH`, `fetch error: MESSAGE`, and `check error: MESSAGE` (the drift test itself did not run, added during implementation so that a broken test is never reported as a fetch error or a drift). Downloads use the resolved commit, not the tag name, so the file matches the commit shown in the summary.

## `.github/workflows/drift.yml`

- **Triggers:** weekly `schedule` (Monday 06:17 UTC) and `workflow_dispatch` (inputs: `project` = all|spec-kit|speckit-roadmap, `tag` = optional).
- **Permissions:** `contents: read` only. No other secrets.
- **Steps:** checkout, then setup-go (1.26.x), then the script with inputs passed as environment variables.
- **Never:** pushes, creates issues or pull requests, uploads artifacts, or moves pins.
- Scheduled runs pause after 60 days without repository activity. Re-enable with `gh workflow enable drift.yml`.

## Review fixes (2026-11-21)

- Each template is compared with its own pinned copy: the script passes `-template=N`, the template's index in `pins.json`, to `TestDrift`.
- A run in which `TestDrift` did not actually execute (renamed, filtered out or skipped) is a `check error`. The log must contain `--- PASS: TestDrift` or `--- FAIL: TestDrift`.
- Unresolved lines in the summary are fenced with more backticks than any run inside them, so upstream text cannot close the code block.
- The dispatch input's project list in `drift.yml` is written out by hand, because GitHub choice inputs cannot read a file. `TestPinsFile` fails unless it is `all` followed by the `pins.json` names in order.
- The tag returned by `releases/latest` is upstream-controlled, so it passes the same checks as `--tag`, and the resolved commit must be 40 hex characters; otherwise the result is a fetch error and nothing is downloaded.
