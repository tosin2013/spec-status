# Quickstart: validate release engineering

## 1. Local: build identity ([contracts/build-identity.md](contracts/build-identity.md))

```sh
CGO_ENABLED=0 go test ./...                                   # includes the new internal/cli build-kind tests
# development (unchanged bytes)
CGO_ENABLED=0 go build -o spec-dev ./cmd/spec && ./spec-dev version
# release-stamped
CGO_ENABLED=0 go build -trimpath -ldflags "-X github.com/tosin2013/spec-status/internal/cli.Version=0.1.0-rc.1 -X github.com/tosin2013/spec-status/internal/cli.Commit=$(git rev-parse HEAD) -X github.com/tosin2013/spec-status/internal/cli.Kind=release" -o spec-rel ./cmd/spec
./spec-rel version            # label: release build
./spec-rel --help | head -3   # "This is release 0.1.0-rc.1."
./spec-rel status --color=never --root testdata/example-roadmap | head -3   # "Spec Status  release 0.1.0-rc.1"
# JSON identical across kinds
for r in testdata/example-roadmap testdata/conformance/*; do cmp <(./spec-dev status --root $r --json) <(./spec-rel status --root $r --json); done
# development text identical to main's build (README sample unchanged)
```

A module build can only be observed after a tag exists. See §4.

## 2. Release candidate run on GitHub (owner pushes the tag)

1. The owner tags a commit on `main`: `git tag v0.1.0-rc.1 && git push origin v0.1.0-rc.1`.
2. Actions → release: `verify`, `ci` (4 test jobs plus hygiene), `build`, `smoke` (2 platforms) and `release` all pass.
3. Releases shows a **draft pre-release** `v0.1.0-rc.1` with exactly three assets ([contracts/release-assets.md](contracts/release-assets.md)). The notes contain the commit, Go version, smoke-job URLs, the compatibility statement and the Limits link.
4. Download both archives and run `sha256sum -c spec_0.1.0-rc.1_checksums.txt`. Both lines say OK.

## 3. Negative runs (SC-002)

These use throwaway tags, which are deleted afterwards along with any draft:
- **Bad name:** `v0.1` or `v0.1.0-beta`. `verify` fails, nothing else runs.
- **Not on main:** a tag on a commit that is only on a scratch branch. `verify` fails.
- **Release exists:** re-push the tag of an existing draft. `verify` fails, and the draft is untouched.

A failing-commit run isn't needed: CI, build and smoke all gate `release` through `needs`, and CI failures were already shown by feature 003's seeded faults.

## 4. After the owner publishes `v0.1.0`

- On a clean linux/amd64 machine and a clean darwin/arm64 machine, without Go and without root, follow the README install section (SC-003, under 5 minutes). `spec version` shows `0.1.0` and `label: release build`.
- `GOBIN=$(mktemp -d) go install github.com/tosin2013/spec-status/cmd/spec@v0.1.0`, then `spec version` shows `0.1.0`, `commit: unknown`, and the module-build label (SC-004).
- README lists linux/amd64 and darwin/arm64 as supported for v0.1.0, citing the smoke runs (FR-016).

## 5. Optional owner settings ([research.md](research.md) R8)

- Enable immutable releases.
- Add a tag ruleset on `v*` restricting tag creation and deletion to the owner.

Neither is applied without the owner's approval.
