# Quickstart: validating the close-gap skill

## Repository checks

```
CGO_ENABLED=0 go test ./internal/extension
bash scripts/check-extension-install.sh        # needs Spec Kit 1.2.0+; expects 6 skills
```

Expected: command structure, links, roles (with `verification-gap`), manifest version 0.3.0 matching the README heading, and the `close-gap-all-causes` Signals all pass; the install check reports 6 skills.

## Agent walkthrough (manual or agent evaluation, as for feature 005)

In a scratch copy of `testdata/skills/close-gap-all-causes` with the extension installed (`specify extension add --dev extensions/spec-status`):

1. Run `/speckit.spec-status.close-gap`. Expected: `recorded verified 1/8`; causes as in contracts/close-gap-command.md "Fixture"; order 002, 003, 004, 005, 008, 006; 007 listed as not started; first proposal is the debrief for 002; no file changed.
2. Decline the proposal. Expected: it moves to 003 and asks whether the work is really done.
3. Say stop. Expected: a summary with nothing applied.
4. Run `/speckit.spec-status.review pm`. Expected: leads with the verification-gap section; first follow-up is close-gap.

## Release checks

The release workflow's archive check expects 6 command files and extension version 0.3.0 on the v0.3.0 tag.
