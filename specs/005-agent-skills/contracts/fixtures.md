# Contract: skill fixtures (`testdata/skills/`)

Each fixture is a minimal Spec Kit project:

```text
testdata/skills/<case>/
├── .specify/memory/roadmap.md
├── .specify/extensions/.registry        # partner cases only
├── specs/<NNN-name>/{spec.md,plan.md,tasks.md}   # as the case needs
└── EXPECTED.md
```

## `EXPECTED.md` format (written by hand, never generated from skill output)

```markdown
# <case>
- Command: /speckit.spec-status.<command> [as <role>]
- Signals (checked by the Go signal test): <item> status=<s>; tasks <c>/<t> state=<state>; diagnostics=<codes>; dependency <id> state=<state>
- Expected finding: <kind and one sentence>
- Evidence: <item> at <path>:<line>
- Expected proposal: <exact change> via <command>   (or "none")
- Expected first follow-ups: <ordered list>        (role and partner cases)
- Must not: <for example "change any file before approval">
```

## Cases (12)

| Case | Command | Key signal |
|---|---|---|
| tasks-done-status-planned | reconcile | status `planned`, tasks 3/3 `read` |
| implemented-unchecked | reconcile | status `implemented`, tasks 1/3 |
| verified-tasks-missing | reconcile | status `verified`, tasks state `missing`, `SS_STAGE_FILE_MISSING` |
| verified-dep-unverified | reconcile | `SS_VERIFIED_DEPENDENCY_CONTRADICTION` |
| unknown-status | reconcile | `SS_STATUS_UNKNOWN` (`status: done`) |
| checked-target-missing | reconcile | `SS_TASK_TARGET_MISSING` error |
| dependency-chain | unblock | 003→002→001 unfinished |
| superseded-dependency | unblock | `SS_DEPENDENCY_HUMAN_REVIEW` |
| dependency-cycle | unblock | `SS_DEPENDENCY_CYCLE` |
| roles-mixed | review as each role | blockers, unchecked tasks, an evidence gap and a supersession together |
| partner-installed | review as developer | `.registry` registers `speckit.bug.{assess,fix,test}` |
| partner-absent | review as developer | no registry, so exactly one install hint |

**Signal test:** for each case, `status.Evaluate` must produce the signals listed in `EXPECTED.md`. Fixtures must stay valid inputs; a fatal outcome fails the test unless the case expects it.
