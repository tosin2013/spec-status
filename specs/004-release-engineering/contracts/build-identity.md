# Contract: Build identity (`spec version`, `--help`, status panel)

The build kind is resolved once per process by `internal/cli/build.go`, in this order:

| Kind | Condition | Version shown | Commit shown |
|---|---|---|---|
| release | stamped variable `cli.Kind == "release"` (set only by `release.yml`) | `cli.Version` | `cli.Commit` |
| module | not release; `debug.ReadBuildInfo()` ok; `Main.Sum != ""`; `Main.Version` not `""` or `(devel)` | `Main.Version` with leading `v` removed | `unknown` |
| development | otherwise | `cli.Version` (default `0.0.0-dev`) | `cli.Commit` (default `unknown`) |

## `spec version` (exact lines; `<go>` = `runtime.Version()`, `<os>/<arch>` = target)

Development, unchanged from today:

```text
spec 0.0.0-dev
commit: unknown
go: <go>
target: <os>/<arch>
label: development build, not a release
```

Release:

```text
spec <version>
commit: <full 40-hex commit>
go: <go>
target: <os>/<arch>
label: release build
```

Module:

```text
spec <version>
commit: unknown
go: <go>
target: <os>/<arch>
label: module build (built from the module version, not a release archive)
```

## `--help`, second paragraph's first sentence

| Kind | Sentence |
|---|---|
| development | `This is a development build.` (unchanged) |
| release | `This is release <version>.` |
| module | `This is a module build of <version>, not a release archive.` |

The rest of the help text is unchanged in every kind, including "It does not change files, use the network, or infer quarters." and the exit-code sentence.

## Status panel heading (human text)

`status.View.Build` holds the text after `Spec Status` and two spaces.

| Kind | `View.Build` | Heading |
|---|---|---|
| development | `""` or `development build` | `Spec Status  development build` (unchanged) |
| release | `release <version>` | `Spec Status  release <version>` |
| module | `module build <version>` | `Spec Status  module build <version>` |

## Invariants (tests)

- `--json` output is byte-identical for all three kinds on the same input. `View` is never consulted by JSON rendering.
- Exit codes, diagnostics, filters and `--explain` are unchanged.
- Development-kind output of `spec version`, `--help` and `spec status` is byte-identical to the build before this feature, so the README sample needs no change.
- Resolution uses only the standard library (`runtime/debug`) and never reads files, the environment, or the clock.
