# Roadmap Sync Report: spec-status (whole ledger)

**Date**: 2026-11-12
**Feature**: all specs (`specs/001`–`004`)
**Roadmap**: [.specify/memory/roadmap.md](roadmap.md) v1.4.2, 11 entries
**Verdict**: ⚠️ PROCEED WITH UPDATES

---

## Summary

All 11 ledger entries were reconciled against the 4 spec directories on disk. `docs/adr/` is absent, so the C-NN constraints are resolved from the roadmap itself.

The ledger and disk agree structurally:
- every lifecycle entry points to a filled spec directory;
- every spec directory is claimed by an entry;
- no dependency points to a missing or abandoned item.

What remains is status and staleness drift:
- **009** is recorded `implemented`, but it has the evidence needed for `verified`. Because a dependency is satisfied only by a `verified` target, this currently blocks **004**.
- Two open questions (Q3, Q4) were answered on 2026-11-11 but are still listed as open.
- 004's implementation is in an unmerged PR, so its ledger state (`specced`, 0/33 tasks on `main`) lags behind the branch (27/33).

`spec status` on `main` also reports four false `SS_REFERENCE_UNSUPPORTED` warnings, caused by link-like text inside inline code in a tasks file. That is a small product finding, not ledger drift.

## Divergence counts

| Type | Count |
|---|---|
| orphan-spec | 0 |
| phantom-entry | 0 |
| status-lagging (dir exists, pre-commitment status) | 0 |
| abandoned-but-active | 0 |
| dependency-contradiction (missing or abandoned target) | 0 |
| superseded-ADR | n/a (no ADR directory) |
| status drift / stale ledger text (beyond the standard types) | 4 |
| product finding surfaced by the audit | 1 |

## Ledger vs disk

| Entry | Status | Spec dir | On disk | Note |
|---|---|---|---|---|
| 001 Core | verified | specs/001-spec-status | filled | ok |
| 002 Contract approval | verified | specs/001-spec-status | filled (shared) | ok; decisions live in 001 |
| 003 Dogfood | verified | specs/001-spec-status | filled (shared) | ok; TASK-16/AC-24 |
| 007 Upstream conformance | verified | specs/002-upstream-format-conformance | filled, 30/30 tasks | ok |
| 009 Dependency automation | implemented | specs/003-dependency-updates | filled, 25/26 tasks | F1 |
| 004 Release engineering | specced | specs/004-release-engineering | filled, 0/33 on main (27/33 in PR #33) | F2 |
| 008 Two-mode | needs-info | none | n/a | waits on Q6 |
| 010 Drift automation | planned | none | n/a | ok; needs a spec |
| 011 README sample test | planned | none | n/a | ok; needs a spec |
| 005 Phase-two write-back | deferred | none | n/a | ok |
| 006 Deferred products | deferred | none | n/a | ok |

## Findings

| ID | Severity | Category | Finding | Suggestion |
|----|----------|----------|---------|------------|
| F1 | 🎯 | status drift (blocking) | 009 is `implemented`, and 004 depends on it. A dependency counts only when its target is `verified` (C-07), so `spec status` reports 004 blocked. 009 has the evidence for verification: all CI checks pass on PR #29; seeded faults failed exactly the expected jobs (three runs recorded); Dependabot's first real PR (#31, pflag and x/sys) came through the merge rule with CI green and no new modules. The only open task is T026, offering branch protection, which is an owner choice and not part of 009's outcome. | Owner decides: record 009 `verified` (via `/speckit.roadmap.write`), closing T026 as "offered". This unblocks 004. |
| F2 | 💡 | status lagging (expected) | 004 is `specced` on `main`, while its implementation (27/33 tasks) is in PR #33. | After #33 merges, move 004 to `in-progress`. Move it to `implemented` after the owner's `v0.1.0-rc.1` run passes (T019), and to `verified` after `v0.1.0` is published and the install check passes on both platforms. |
| F3 | 💡 | stale ledger text | Q3 (initial targets) and Q4 (license, security route, support window) are still listed as open. D-2026-11-11-02 and D-2026-11-11-03 answered them. | Strike both through with "Resolved 2026-11-11", as Q1 and Q2 already are. |
| F4 | 🤔 | open question | Q5 (whether the speckit-roadmap annotated dependency style should gain more support) has had no decision since 2026-11-07. It does not block v0.1.0, because the current behavior is a visible `SS_DEPENDENCY_MALFORMED`. | Owner: decide now, or mark it deferred past v0.1.0. |
| F5 | 💡 | product finding | `spec status` on `main` reports 4× `SS_REFERENCE_UNSUPPORTED` at `specs/004-release-engineering/tasks.md:89`. The line contains a regex inside backticks (`` `^v(0|[1-9][0-9]*)…` ``), and the tasks-file reference scan treats `[…](…)`-like text inside inline code as a link. It is a visible warning, not a miscount, so it is non-blocking under D-2026-11-07-09. | Add a roadmap item, or include it in 011-style maintenance: skip code spans when scanning task-file references, with a regression test. Until then it can be listed in README Limits. |

## Completion view

**v0.1.0 (first tagged release), roadmap scope:**

| Area | State |
|---|---|
| Core reporting (001), contract decisions (002), dogfood (003), conformance (007) | recorded `verified` |
| Dependency automation and CI (009) | `implemented`; evidence in hand (F1) |
| Release engineering (004) | code in PR #33; release candidate, publish and install check pending (owner) |
| Remaining owner steps | merge #33 → tag `v0.1.0-rc.1` → review the draft → publish `v0.1.0` |

That is 4 of the 6 v0.1.0 items `verified`, one ready to verify, and one in final validation.

**After v0.1.0, already on the ledger:**

| Entry | Needs |
|---|---|
| 008 Two-mode roadmap and feature inventory | Owner choices in Q6 (command syntax, default mode, template-version identification, folder and checkbox rules, inventory JSON schema and exits, v1 placement), then a spec |
| 010 Upstream drift automation | A spec: which releases to fetch, how drift is reported, network access as development tooling only |
| 011 README sample test | A small spec, or a direct task under 001 maintenance |
| F5 code-span reference fix | A small fix with a test, possibly folded into 011 |
| 005 Phase-two write-back, 006 Deferred products | Deferred; no action until reopened |
| Q5 | Decide or defer |

## Recommended Actions *(proposed only, never auto-applied)*

1. **F1**: Owner decision on recording 009 `verified`. Apply with `/speckit.roadmap.write`. This is the one finding that blocks the release path.
2. **F3**: Mark Q3 and Q4 resolved (wording, patch bump).
3. **F2**: Update 004's status as milestones land (after merge, after the release candidate, after publish).
4. **F4**: Owner: decide or defer Q5.
5. **F5**: Track the code-span false warning, either as a new planned item or under 011.

---

*Generated by `/speckit.roadmap.sync`. Read-only: this report proposes changes and applies none.*
