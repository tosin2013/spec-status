# Quickstart: validating the drift workflow

```sh
CGO_ENABLED=0 go test ./internal/conformance        # pins.json read; pinned copies exist
scripts/check-upstream-drift.sh                     # real upstream, all projects (needs gh auth)
```

**Seeded checks (SC-003):**

1. `scripts/check-upstream-drift.sh --project spec-kit --candidate spec-kit=testdata/upstream/spec-kit/v1.2.0/templates/tasks-template.md` gives `no unresolved change`, exit 0.
2. The same, with a copy that has one added `- [ ] T999 [P9] x` line, gives `1 unresolved change`, exit 1, and the summary lists the key.
3. `scripts/check-upstream-drift.sh --project spec-kit --path spec-kit=templates/nope.md` gives `template not found: templates/nope.md`, exit 1.
4. `scripts/check-upstream-drift.sh --project spec-kit --tag v0.0.0-nope` gives `fetch error: …`, exit 1.
5. `scripts/check-upstream-drift.sh --project spec-kit --tag v1.2.0` gives `same tag as the pin`, exit 0.

**CI:** `gh workflow run drift.yml --ref 013-upstream-drift-workflow-impl`, then read the run summary. Record the real result.
