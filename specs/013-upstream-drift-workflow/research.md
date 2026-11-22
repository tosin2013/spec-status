# Research: Upstream drift workflow

## R1 — One pin list: `testdata/upstream/pins.json`

- **Decision:** Use a JSON file, read by the Go test (`encoding/json`) and by the script (`jq`):

  ```json
  {"projects":[
    {"name":"spec-kit","repo":"github/spec-kit","tag":"v1.2.0",
     "templates":[{"path":"templates/tasks-template.md","pinned":"spec-kit/v1.2.0/templates/tasks-template.md"}]},
    {"name":"speckit-roadmap","repo":"srobroek/speckit-roadmap","tag":"v0.1.0",
     "templates":[{"path":"templates/roadmap-template.md","pinned":"speckit-roadmap/v0.1.0/templates/roadmap-template.md"}]}]}
  ```

  `pinned` is relative to `testdata/upstream/`. `drift_test.go` builds its project → files map from it. `TestPinsFile` checks that every `pinned` file exists, that the project names match the `resolutions.md` directories, and that the pinned tag appears in the `pinned` path.
- **Rationale:** FR-009. Moving a pin becomes one edit, which the test checks.
- **Alternatives considered:**
  - Keeping the Go map and parsing it from bash: brittle.
  - YAML: no standard-library parser in Go.
  - Deriving the pins from directory names: the repository paths would not be recorded anywhere.

## R2 — The checker script

- **Decision:** `scripts/check-upstream-drift.sh [--project NAME] [--tag TAG] [--candidate NAME=FILE] [--path NAME=REPO_PATH]`.
  - `--candidate` uses a local file instead of downloading. It is for seeded tests and for manual pin moves from a local download.
  - `--path` overrides the repository path. It is for the missing-template test.
  - Resolution: `gh api repos/R/releases/latest -q .tag_name` (GitHub's "latest" excludes drafts and prereleases), then the commit from `gh api repos/R/commits/TAG -q .sha`. A failure in either is a `fetch error`.
  - Download: `curl -fsS --max-time 30 -o FILE https://raw.githubusercontent.com/R/TAG/PATH`. HTTP 404 means `template not found`; any other failure is a `fetch error`.
  - Check: `go test ./internal/conformance -run 'TestDrift$' -count=1 -v -args -project=P -candidate=F`. The changed keys are the `added:` and `removed:` lines marked `(unresolved)`. Test exit 0 means no unresolved change.
  - `same tag as the pin` is reported when the candidate tag equals the pinned tag and there is no unresolved change.
  - Summary: a Markdown table (project, pinned, candidate tag and short commit, result) plus the unresolved lines per project. It goes to `$GITHUB_STEP_SUMMARY` when that is set, and always to stdout. The exit code is 1 if any project's result is a failure.
  - Candidate files are written only under `mktemp -d`, with a trap that removes them. Their content is never interpolated into a shell command. In the summary, lines from the test output are listed inside a code block.
- **Rationale:** The same code runs locally and in CI, and it is testable without CI (US3 seeded runs). It is portable to bash 3.2, like `check-extension-install.sh`.
- **Alternatives considered:**
  - All logic in workflow YAML: not testable locally.
  - A Go program: it would need net/http in this module, against the spirit of the import guard. Even under `scripts/`, this is avoidable.

## R3 — Workflow shape

- **Decision:**
  - Triggers: `on: schedule: - cron: '17 6 * * 1'` (Mondays 06:17 UTC, off the hour to avoid peak load) and `workflow_dispatch` with inputs `project` (choice: all, spec-kit, speckit-roadmap; default all) and `tag` (string, default empty; used only with a single project).
  - `permissions: contents: read`.
  - One job on `ubuntu-latest`: checkout and setup-go at the ci.yml SHAs, with Go 1.26.x. It runs the script with `GH_TOKEN: ${{ github.token }}`, passing inputs through `env:`, not by interpolation into `run:`.
  - `timeout-minutes: 10`.
- **Rationale:**
  - FR-001 to FR-006.
  - Passing inputs through env avoids script injection from dispatch inputs.
  - The token raises API rate limits and is read-only.
- **Alternatives considered:**
  - Daily runs: more noise than the "within a week" outcome needs.
  - A matrix per project: the summary would be split, and the shared setup duplicated.

## R4 — Reporting

- **Decision:** A failed run is the notification, and the step summary is the report. No issues, comments or artifacts with write scope.
- **Rationale:** The roadmap 010 notes say issue creation needs write permission, and the owner is the only maintainer.
- **Alternatives considered:**
  - Opening or updating an issue: needs `issues: write`. It was deferred, and can be a later option.

## R5 — The 60-day pause

- **Decision:** Document that GitHub disables scheduled workflows in a repository with no activity for 60 days, and that re-enabling is one click in the Actions tab (or `gh workflow enable drift.yml`).
- **Rationale:** FR-010. The outcome depends on the schedule actually running.

## R6 — Verifying it

- **Decision:**
  - Seeded local runs:
    - `--candidate spec-kit=<pinned copy>` gives `same tag`, or no change;
    - a copy with an added `[Pn]` marker line gives `1 unresolved change`;
    - `--path spec-kit=templates/nope.md` gives `template not found`;
    - `--tag v0.0.0-nope` gives `fetch error`.
  - One real dispatch run on the feature branch with the current upstream releases. Record its result as found, whatever it is.
- **Rationale:** SC-003. There must be evidence for each result kind, and no false drift.
