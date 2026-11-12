# Quickstart: validate the agent skills extension

## 1. Repository checks (CI)

```sh
CGO_ENABLED=0 go test ./internal/extension -v   # reference vs contracts, manifest, structure, links, roles, version, fixture signals
CGO_ENABLED=0 go test ./...                      # everything else unchanged
```

Expected: all pass. To prove the checks catch drift, temporarily:
- delete one `SS_*` row from `reference/output.md`: the test fails, naming the code;
- add a `/speckit-nonexistent` link: the test fails;
- give a role the section `foo`: the test fails.

Revert each change afterwards.

## 2. CLI unchanged (SC-006)

Build `spec` from `main` and from the branch, then `cmp` the `--json` and `--color=never` output on `testdata/example-roadmap` and every `testdata/conformance/*` fixture. All must be identical.

## 3. Install test (local, recorded; FR-018, SC-005)

```sh
scripts/check-extension-install.sh
```

The script:
1. Creates a temporary directory and runs `specify init --here --force --non-interactive --integration claude` there.
2. Runs `specify extension add --dev <repo>/extensions/spec-status`.
3. Checks that `.claude/skills/` has exactly the five `speckit-spec-status-*` skills, and that review's skill carries the `argument-hint`.

It prints the elapsed time, which should be under 60 s.

## 4. Behavior on fixtures (manual or agent-run; SC-001, SC-004, SC-008, SC-009)

For each `testdata/skills/<case>`:
1. Copy the case into a temporary directory.
2. Install the extension there.
3. Run the command from its `EXPECTED.md` in an agent session.
4. Compare the result with `EXPECTED.md`: the finding, evidence, proposal, follow-ups, and that no file changed before approval (`git status` clean).

Record pass or fail per case in the implementation record.

Role check (SC-008): on `roles-mixed`, run review with no role and with each of the four roles. The counts and findings must be identical; only the order and the first follow-ups differ, as listed in `roles.md`.

Partner check (SC-009):
- `partner-installed` lists `speckit.bug.*` among the first developer follow-ups.
- `partner-absent` gives exactly one hint, `specify extension add bug`.
- With Bug Triage really installed (`specify extension add bug`), the same output appears.

## 5. Release asset (after tagging v0.2.0)

The draft release has a fourth asset, `spec-status-extension_0.2.0.zip`, listed in the checksums. `specify extension add spec-status --from <asset URL>` installs it into a fresh project and renders the five skills.
