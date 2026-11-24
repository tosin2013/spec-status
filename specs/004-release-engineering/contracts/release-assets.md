# Contract: Release assets and notes

`<version>` is the tag without the leading `v`, for example `0.1.0` or `0.1.0-rc.1`.

## Assets (exactly six per release, since v0.4.1; four from v0.2.0 to v0.4.0)

| File | Content |
|---|---|
| `spec_<version>_linux_amd64.tar.gz` | directory `spec_<version>_linux_amd64/` containing `spec` (mode 0755), `LICENSE`, `README.md` |
| `spec_<version>_linux_arm64.tar.gz` | directory `spec_<version>_linux_arm64/` containing `spec` (mode 0755), `LICENSE`, `README.md` (since v0.4.1) |
| `spec_<version>_darwin_arm64.tar.gz` | directory `spec_<version>_darwin_arm64/` containing `spec` (mode 0755), `LICENSE`, `README.md` |
| `spec_<version>_darwin_amd64.tar.gz` | directory `spec_<version>_darwin_amd64/` containing `spec` (mode 0755), `LICENSE`, `README.md` (since v0.4.1) |
| `spec-status-extension_<version>.zip` | the Spec Kit extension (feature 005): one top-level directory `spec-status/` with `extension.yml`, `README.md`, `LICENSE`, `commands/` (5 files) and `reference/` (3 files). Its `extension.version` equals `<version>`; the release `verify` job enforces this |
| `spec_<version>_checksums.txt` | `sha256sum` format: `<64 hex>␣␣<file name>`, one line per archive (the four binary archives and the extension zip), sorted by name |

There are no other files: no Windows archive (planned, roadmap 020), no signatures and no SBOM. Platform details: [014 platforms contract](../../014-more-release-platforms/contracts/platforms.md). Releases before v0.2.0 have three assets, with no extension zip.

## Release notes template (`notes.md`, prepended to GitHub's generated notes)

```markdown
## spec <version>

- Commit: `<full commit>`
- Built with: <go version output>, `CGO_ENABLED=0`, `-trimpath`
- Archives:
  - `spec_<version>_linux_amd64.tar.gz`: smoke-tested on linux/amd64 (<smoke job URL>)
  - `spec_<version>_linux_arm64.tar.gz`: smoke-tested on linux/arm64 (<smoke job URL>)
  - `spec_<version>_darwin_arm64.tar.gz`: smoke-tested on darwin/arm64 (<smoke job URL>)
  - `spec_<version>_darwin_amd64.tar.gz`: smoke-tested on darwin/amd64 (<smoke job URL>)
- Checksums: `spec_<version>_checksums.txt` (SHA-256). Verify before use, as described in the README install section.

**Compatibility (pre-1.0):** before v1.0.0, a minor version may change CLI flags, exit codes or the JSON schema. Any such change is listed in these notes. Patch versions do not break compatibility. JSON schema: 1.0.

**Not provided:** signed or notarized binaries, provenance attestations, reproducible-build guarantees, or archives for other platforms. Other platforms may build from source but are untested. Known limits: [README Limits](https://github.com/tosin2013/spec-status#limits).
```

After the "Not provided" paragraph, the notes add a section **Use it with an AI coding agent**: one sentence, then the prompt block from `docs/agent-prompt.md` (between its `agent-prompt` markers), so every release page shows a copy-and-paste install-and-review prompt. CI's hygiene job keeps the README copy identical.

For a release candidate, the heading is `## spec <version> (release candidate)`.
