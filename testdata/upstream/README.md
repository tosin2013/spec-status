# Pinned upstream templates

These directories hold the upstream templates that define the input spec status supports (feature 002):

| Project | Current pin | Files |
|---|---|---|
| github/spec-kit | `v1.2.0` (commit `4ca46a0706a662b1edc07086ef16f021820c2db0`) | `spec-kit/v1.2.0/` |
| srobroek/speckit-roadmap | `v0.1.0` | `speckit-roadmap/v0.1.0/` |

Each pin keeps its upstream LICENSE and a NOTICE with the source, tag, commit, and file checksums. The supported shapes are listed in `specs/002-upstream-format-conformance/contracts/shape-inventory.md`.

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
4. When the drift check passes and `CGO_ENABLED=0 go test ./...` passes, copy the candidate templates into `PROJECT/TAG/templates/` with LICENSE and NOTICE, update `pinnedTemplates` in `internal/conformance/drift_test.go`, update the shape inventory, and update the "Supported input" section of README.md.
5. Shapes supported under the previous pin stay supported unless a recorded decision removes them (FR-010). Note any removal for users.
