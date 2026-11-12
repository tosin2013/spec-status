# Security policy

## Reporting a vulnerability

Report vulnerabilities privately through GitHub's private vulnerability reporting:
open the repository's **Security** tab and choose **Report a vulnerability**.
Do not open a public issue for a suspected vulnerability.

## Supported versions

Only the latest release is supported. Fixes are made on a best-effort basis by
the maintainer; there is no promised response time. There is no release yet, so
reports currently apply to the `main` branch.

## Dependency vulnerabilities

Vulnerabilities in dependencies are found in four ways:

- **Dependabot alerts:** GitHub reports published advisories for any module in `go.mod` and `go.sum`.
- **govulncheck in CI:** every pull request and every push to `main` runs Go's vulnerability scanner. It fails when the code can reach a vulnerable function, in a module or in the Go standard library.
- **Dependency review:** each pull request lists the modules it adds or changes, with known advisories and licenses. It fails if an added or changed version has a known vulnerability.
- **CodeQL code scanning:** GitHub's static analysis for Go runs on pushes and pull requests.

Fixes arrive as pull requests, often opened by Dependabot. The maintainer reviews them under the merge rule in [CONTRIBUTING.md](CONTRIBUTING.md) and merges them to `main`. **A fix reaches users only in a release.** No release exists yet, so fixes are on `main` until the first release.

These are automated checks. They are not a security audit or a supply-chain certification. Release binaries are not signed.

## Scope

`spec status` is a read-only command that reads files beneath a project root. It
makes no network calls, runs no subprocesses, and writes nothing. Relevant reports
include vulnerable dependencies, reads outside the root, reads of excluded sensitive files, and terminal
escape sequences from roadmap text reaching human output. The sensitive-name
list in `internal/safeio` is a starting point, not a security certification
(decision D-2026-11-11-02).
