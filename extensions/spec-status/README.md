# spec-status extension v0.2.0

A [Spec Kit](https://github.com/github/spec-kit) extension that turns [`spec status`](https://github.com/tosin2013/spec-status) output into next actions. Once installed, its five commands appear as skills in your AI coding agent (Claude Code, Cursor, Codex, Copilot and others that Spec Kit supports).

## Requirements

- Spec Kit 1.2.0 or newer.
- The `spec` CLI, v0.1.2 or newer, on your `PATH`. See the install steps: https://github.com/tosin2013/spec-status#install
- Optional: the [speckit-roadmap](https://github.com/srobroek/speckit-roadmap) extension, so that approved roadmap changes go through `/speckit.roadmap.write`.

## Install

From a release (replace `<version>`):

```sh
specify extension add spec-status --from https://github.com/tosin2013/spec-status/releases/download/v<version>/spec-status-extension_<version>.zip
```

From a clone of the repository:

```sh
specify extension add --dev /path/to/spec-status/extensions/spec-status
```

## Commands

| Command | What it does |
|---|---|
| `/speckit.spec-status.review [as pm\|developer\|qa\|architect]` | Start here. Summarizes progress, blockers, human-review dependencies, diagnostics and what is ready to start, then sends each finding to the right command |
| `/speckit.spec-status.reconcile` | Finds recorded statuses that disagree with tasks or evidence, and proposes one fix at a time |
| `/speckit.spec-status.explain [code]` | Explains any `SS_*` diagnostic and the command that fixes it |
| `/speckit.spec-status.unblock` | Lists blockers and superseded dependencies in the order to clear them |
| `/speckit.spec-status.verify-readiness [item]` | Lists the evidence still missing before an item can honestly be recorded `verified` |

**Roles** change what `review` shows first, never the facts. **Partner extensions** are recommended only when installed:
- architect: arch-governance
- qa: Assure
- developer: Bug Triage

See [reference/roles.md](reference/roles.md).

## What it will not do

It never changes a file on its own. Every change is proposed, applied only through a Spec Kit command after you approve that specific change, and then confirmed by running `spec status` again ([reference/rules.md](reference/rules.md)). The `spec` CLI stays read-only.

## License

Apache-2.0 ([LICENSE](LICENSE)).
