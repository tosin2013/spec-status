# Research: more release platforms

## R1 — Runners

- **Decision:** `ubuntu-24.04-arm` for linux/arm64 and `macos-15-intel` for darwin/amd64. Both are on the actions/runner-images list (checked 2026-11-24), the arm64 Linux runner is available to public repositories, and the labels are pinned to an OS version.
- **Rationale:** native hardware for each archive (D-2026-11-11-02).
- **Alternatives considered:**
  - `macos-26-intel`: newer, but chosen later if `macos-15-intel` is retired.
  - QEMU emulation on x64: not native evidence.
  - `macos-latest-large`: a paid runner.

## R2 — Smoke `case` values

- **Decision:** `uname -s -m` gives `Linux x86_64`, `Linux aarch64`, `Darwin arm64` or `Darwin x86_64`, which map to linux/amd64, linux/arm64, darwin/arm64 and darwin/amd64. Linux uses `sha256sum` and macOS `shasum -a 256`. Any other value fails the job.
- **Rationale:** the existing `target:` check in `spec version` then proves the binary matches the machine.

## R3 — Draft notes and assets

- **Decision:** a small table in bash maps each smoke job name to its platform and archive:

  | Smoke job | Platform | Archive |
  |---|---|---|
  | `smoke (ubuntu-latest)` | linux/amd64 | `spec_<v>_linux_amd64.tar.gz` |
  | `smoke (ubuntu-24.04-arm)` | linux/arm64 | `spec_<v>_linux_arm64.tar.gz` |
  | `smoke (macos-latest)` | darwin/arm64 | `spec_<v>_darwin_arm64.tar.gz` |
  | `smoke (macos-15-intel)` | darwin/amd64 | `spec_<v>_darwin_amd64.tar.gz` |

  For each row, the job prints "smoke-tested on <platform> (<url>)", and fails if the URL is empty. `gh release create` takes the four archives explicitly, in the same order.
- **Rationale:** FR-004 and FR-005. A missing smoke job must never produce a release line without a link.
- **Alternatives considered:** a glob over `dist/*.tar.gz`. That would attach an archive whose smoke job did not run.

## R4 — CI platform assert

- **Decision:** keep the `go env GOOS GOARCH` check and add a `uname -s -m` check. Expected values:

  | Runner | `go env` | `uname` |
  |---|---|---|
  | ubuntu-latest | linux amd64 | Linux x86_64 |
  | ubuntu-24.04-arm | linux arm64 | Linux aarch64 |
  | macos-latest | darwin arm64 | Darwin arm64 |
  | macos-15-intel | darwin amd64 | Darwin x86_64 |

- **Rationale:** a Go toolchain running under Rosetta would report amd64 on an arm64 machine; `uname -m` reports the real machine there.
- **Alternatives considered:** the `go env` check alone, which is what exists today.

## R5 — Version and docs

- **Decision:**
  - Extension 0.4.1.
  - README platform table: two rows marked "v0.4.1"; smoke-test links are added after the release.
  - Agent prompt: Linux x86-64 → linux_amd64, Linux arm64/aarch64 → linux_arm64, macOS Apple silicon → darwin_arm64, macOS Intel → darwin_amd64. Anything else, including Windows, has no archive yet.
  - RELEASING.md asset checklist: 4 archives, checksums, extension zip (the current text says "exactly three", which is stale).
- **Rationale:** FR-007 and FR-008.
