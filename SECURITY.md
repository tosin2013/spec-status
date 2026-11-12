# Security policy

## Reporting a vulnerability

Report vulnerabilities privately through GitHub's private vulnerability reporting:
open the repository's **Security** tab and choose **Report a vulnerability**.
Do not open a public issue for a suspected vulnerability.

## Supported versions

Only the latest release is supported. Fixes are made on a best-effort basis by
the maintainer; there is no promised response time. There is no release yet, so
reports currently apply to the `main` branch.

## Scope

`spec status` is a read-only command that reads files beneath a project root. It
makes no network calls, runs no subprocesses, and writes nothing. Relevant reports
include reads outside the root, reads of excluded sensitive files, and terminal
escape sequences from roadmap text reaching human output. The sensitive-name
list in `internal/safeio` is a starting point, not a security certification
(decision D-2026-11-11-02).
