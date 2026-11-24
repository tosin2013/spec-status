# Contract: supported release platforms (from v0.4.1)

This amends `specs/004-release-engineering/contracts/release-assets.md` and `release-workflow.md`.

- **Assets per release:** four `spec_<v>_<os>_<arch>.tar.gz` archives (the platforms in data-model.md), `spec-status-extension_<v>.zip` and `spec_<v>_checksums.txt` (SHA-256 of all archives and the zip, sorted by name).
- **Smoke tests:** each archive is smoke-tested on its own native runner with the same checks:
  - checksum OK;
  - exact listing;
  - `spec version` shows the version, commit, `target: <os>/<arch>` and `label: release build`;
  - help says the release version;
  - text and JSON exits agree;
  - a missing root exits 2;
  - the JSON has schema_version 1.0.

  The extension zip is checked once, on `ubuntu-latest`.
- **Draft notes:** one "smoke-tested on <os>/<arch> (<smoke job URL>)" line per archive. The draft job fails if any smoke job URL is missing, and the release is not drafted unless every smoke job passed.
- **CI:** test jobs run on all four runners. Each asserts both the Go toolchain platform and the machine (`uname`).
- **Not provided:** Windows (roadmap 020), signed or notarized binaries, provenance.
