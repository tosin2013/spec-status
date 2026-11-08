# Pinned upstream: github/spec-kit v1.2.0

- Source: https://github.com/github/spec-kit
- Tag: `v1.2.0`
- Commit: `4ca46a0706a662b1edc07086ef16f021820c2db0` (recorded by the installed specify-cli 1.2.0 in its `direct_url.json`)
- License: MIT, copyright GitHub, Inc. (`LICENSE`, copied from the package's `licenses/LICENSE`)

Copied byte-for-byte from the package directory `specify_cli/core_pack/templates/`:

| File | SHA-256 |
|---|---|
| `templates/tasks-template.md` | `1e448e9153462e8c5b5a55231a54c756d870a82950d4dd9425a2bf5b7008a0d9` |
| `templates/spec-template.md` | `3945437fc35cd30a5b2bf7beea680337c3516826d3efa5a6b92c4a7eca1ba28e` |
| `templates/plan-template.md` | `dfba43da2b2a207ff8c177a0fba03af074fd45e261696b4b64172041fad2e992` |

These are the package copies, not the copies `specify init` writes into a project, which substitute agent-specific command names (research R2). They define the Spec Kit shapes spec status supports (feature 002). spec status reads `tasks.md` content and only checks that `spec.md` and `plan.md` exist.
