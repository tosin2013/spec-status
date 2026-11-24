# Data model: more release platforms

| Platform | Archive | Smoke runner | CI test runner | `uname -s -m` | First release |
|---|---|---|---|---|---|
| linux/amd64 | spec_<v>_linux_amd64.tar.gz | ubuntu-latest | ubuntu-latest | Linux x86_64 | v0.1.2 |
| linux/arm64 | spec_<v>_linux_arm64.tar.gz | ubuntu-24.04-arm | ubuntu-24.04-arm | Linux aarch64 | v0.4.1 |
| darwin/arm64 | spec_<v>_darwin_arm64.tar.gz | macos-latest | macos-latest | Darwin arm64 | v0.1.2 |
| darwin/amd64 | spec_<v>_darwin_amd64.tar.gz | macos-15-intel | macos-15-intel | Darwin x86_64 | v0.4.1 |

Each archive contains `spec_<v>_<os>_<arch>/{spec,LICENSE,README.md}` and has a line in `spec_<v>_checksums.txt`.
