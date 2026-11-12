# Data model: Release engineering

No stored data. The entities are build-time values, files and GitHub objects.

## Build identity

| Field | Source | Rule |
|---|---|---|
| kind | `cli.Kind` stamp; `debug.ReadBuildInfo()` | `release` \| `module` \| `development` ([contracts/build-identity.md](contracts/build-identity.md)) |
| version | `cli.Version` (release, development) or `Main.Version` without `v` (module) | release: equals the tag without `v`; development default `0.0.0-dev` |
| commit | `cli.Commit` (release, development); `unknown` (module) | release: 40-hex `$GITHUB_SHA` |
| label | derived from kind | exact strings in the contract |

## Release tag

| Field | Rule |
|---|---|
| name | `^v(0\|[1-9][0-9]*)\.(0\|[1-9][0-9]*)\.(0\|[1-9][0-9]*)(-rc\.(0\|[1-9][0-9]*))?$` |
| commit | reachable from `origin/main` |
| prerelease | true when the name contains `-rc.` |
| creator | the owner (a ruleset can enforce this; offered, research R8) |

## Archive

| Field | Rule |
|---|---|
| name | `spec_<version>_<os>_<arch>.tar.gz`, with `<os>_<arch>` ∈ {`linux_amd64`, `darwin_arm64`} |
| layout | one top-level directory with the same base name, containing `spec`, `LICENSE`, `README.md` |
| binary | `CGO_ENABLED=0`, `-trimpath`, `-s -w`, stamped kind `release` |

## Checksum file

| Field | Rule |
|---|---|
| name | `spec_<version>_checksums.txt` |
| content | `sha256sum` lines for both archives, computed once in the build job |

## GitHub Release

States follow `draft → published`. The draft is created by the workflow; publishing is done only by the owner.

| Field | Rule |
|---|---|
| tag | the pushed tag (`--verify-tag`) |
| draft | always true when created |
| prerelease | equals the tag's prerelease flag |
| assets | exactly the three files |
| notes | the template in [contracts/release-assets.md](contracts/release-assets.md), followed by generated notes |
| after publish | never edited; corrections ship as a new version (immutable releases setting offered, R8) |

## Workflow jobs

`verify → ci (reused ci.yml) → build → smoke[linux/amd64, darwin/arm64] → release`. Only `release` holds `contents: write`.
