# Data model: Upstream drift workflow

| Entity | Fields | Source |
|---|---|---|
| Pin (`pins.json` project) | name, repo (owner/name), tag, templates[] {path (repository path at the tag), pinned (copy under testdata/upstream/)} | `testdata/upstream/pins.json` |
| Candidate | project, tag, commit (40-hex), files[] (downloaded, temporary) | GitHub API and raw download |
| Result (one per project) | kind: `same-tag`, `no-change`, `drift` (count + lines), `template-not-found` (path), `fetch-error` (message); failed = drift, not-found or fetch-error | script |

Run outcome: success if no project failed; otherwise exit 1. The summary has a table row per project, then details.
