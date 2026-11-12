# Contract: Release assets and notes

`<version>` is the tag without the leading `v`, for example `0.1.0` or `0.1.0-rc.1`.

## Assets (exactly three per release)

| File | Content |
|---|---|
| `spec_<version>_linux_amd64.tar.gz` | directory `spec_<version>_linux_amd64/` containing `spec` (mode 0755), `LICENSE`, `README.md` |
| `spec_<version>_darwin_arm64.tar.gz` | directory `spec_<version>_darwin_arm64/` containing `spec` (mode 0755), `LICENSE`, `README.md` |
| `spec_<version>_checksums.txt` | `sha256sum` format: `<64 hex>␣␣<archive file name>`, one line per archive, sorted by name |

There are no other files: no Windows archive, no linux/arm64 archive, no signatures and no SBOM.

## Release notes template (`notes.md`, prepended to GitHub's generated notes)

```markdown
## spec <version>

- Commit: `<full commit>`
- Built with: <go version output>, `CGO_ENABLED=0`, `-trimpath`
- Archives:
  - `spec_<version>_linux_amd64.tar.gz`: smoke-tested on linux/amd64 (<smoke job URL>)
  - `spec_<version>_darwin_arm64.tar.gz`: smoke-tested on darwin/arm64 (<smoke job URL>)
- Checksums: `spec_<version>_checksums.txt` (SHA-256). Verify before use, as described in the README install section.

**Compatibility (pre-1.0):** before v1.0.0, a minor version may change CLI flags, exit codes or the JSON schema. Any such change is listed in these notes. Patch versions do not break compatibility. JSON schema: 1.0.

**Not provided:** signed or notarized binaries, provenance attestations, reproducible-build guarantees, or archives for other platforms. Other platforms may build from source but are untested. Known limits: [README Limits](https://github.com/tosin2013/spec-status#limits).
```

For a release candidate, the heading is `## spec <version> (release candidate)`.
