# Contract: Dependabot configuration (`.github/dependabot.yml`)

A conforming configuration meets every **MUST** below.

- MUST be `version: 2`.
- MUST contain exactly two `updates` entries (FR-001, FR-002):
  - `package-ecosystem: gomod`, `directory: /`, `schedule.interval: weekly`, `schedule.day: monday`. It MUST have one group that matches all dependencies (`patterns: ["*"]`) and is limited to `update-types: [minor, patch]`, so that major updates open separate pull requests.
  - `package-ecosystem: github-actions`, `directory: /`, `schedule.interval: weekly`, `schedule.day: monday`. It MUST have one group matching all actions (`patterns: ["*"]`), limited to `update-types: [minor, patch]`, so that major action updates also open separate pull requests (User Story 3).
- MUST NOT configure auto-merge, `insecure-external-code-execution`, private registries, or `ignore` rules (spec Out of Scope; none are needed).
- Security updates come from the repository setting enabled on 2026-11-11 (FR-003). This file MUST NOT group them: no `applies-to: security-updates` group.

## Repository settings (not in the file)

| Setting | Required state | Checked with |
|---|---|---|
| Dependabot alerts | enabled | `gh api repos/tosin2013/spec-status/vulnerability-alerts -i` → `204` |
| Code scanning (CodeQL default setup, Go) | configured (FR-020) | `gh api repos/tosin2013/spec-status/code-scanning/default-setup --jq .state` → `configured` |
| Dependabot security updates | enabled | `gh api repos/tosin2013/spec-status/automated-security-fixes` → `{"enabled":true,...}` |

Alerts and security updates were enabled and checked on 2026-11-11. Code scanning is applied during implementation.
