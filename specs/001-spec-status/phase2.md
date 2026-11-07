# Optional future extension-mediated write-back

**Draft future contract — out of v1.** This document preserves a future approval contract, not an implemented command, supported extension integration or authorization to mutate files. No writer CLI syntax, native executable or extension version is assumed.

## Isolation and truth source

V1 `status`, JSON and CI are strictly read-only and never silently launch an agent/extension. Future write-back must be an explicit separate mode using the roadmap extension's own supported command workflow, not direct behind-the-extension rewrites. Establish exact supported invocation/version during separately approved design. The roadmap file remains the sole source of truth; no competing status store.

AI/extension dependencies, if needed, are isolated to the optional mode. Unavailable integration reports its limitation without breaking standalone reporting. Installation of the core must not require the extension, account, credentials or AI. The future writer does not expand what current status is allowed to do.

## Per-mutation approval invariant

Before **each** mutation, present affected roadmap IDs, exact proposed changes and target roadmap path for explicit owner approval. A diagnostic, recommendation, LLM request, prior unrelated approval or an entire task backlog is not approval. A replacement reference or dependency migration is not implicitly authorized by approving a status change.

Read and bind proposal to observed source. If that source changes after approval, refresh proposal and obtain approval again; do not apply stale instructions, silently reconcile changes or reuse earlier approval. Exact binding/expiry/race-safe mechanism remains a future design decision, not solved by a timestamp mention. Required approvals cannot be delegated to an autonomous consumer implicitly.

## Observed-result contract

1. **Proposed:** safe preview only; no claimed persisted change.
2. **Approved:** specific change/path/IDs consented against unchanged observed source; approval is not a successful write.
3. **Attempted:** supported extension workflow invoked under the separate approved mode; success text is not persisted confirmation.
4. **Confirmed:** reread actual roadmap, compare intended fields and report observed persisted outcome.
5. **Unconfirmed/failed/source-changed:** disclose limitation and stop; no claimed completion or automatic extra mutation/retry. If another change is needed, propose it and obtain approval again.

If integration writes something different or partial, do not silently repair it. Reread and explain observed state without inventing semantic correctness. No rewrite/redirect of dependencies follows automatically from a successor link. Preserve original row, ID, raw status and history. Authoring replacement references requires this separately approved workflow and is never a v1 parser action.

## Future acceptance boundary

Require explicit ID/change/path approval, stale-source rejection/reapproval, integration unavailable behavior, per-change rather than batch-implied authority, persisted reread and proposed/attempted/confirmed distinctions. Ordinary status/JSON/CI must still be shown never to invoke integration. No background listener, schedule, external notification, dashboard, digest or multi-repo write is authorized here. Those optional uses each need their own scope/design/approval; listing them is not a commitment for phase two.
