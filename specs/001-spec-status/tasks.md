# Dependency-ordered work backlog — unstarted documentation draft

Every task is **unstarted**; drafting this backlog authorizes no code, workflow, repository change or publication. Sequence is a proposal, not a one-day estimate or an approved 6–12 day commitment. “Parallel” below means logically independent future work only after dependencies and separate authorization are satisfied.

## Gate 0 — approve the contract before coding

| ID | Depends on | Future work / completion evidence | Acceptance mapping |
|---|---|---|---|
| TASK-01 | — | Review constitution, settled boundaries and scope; record exact implementation authorization separately. Confirm no inferred quarters, strict supersession, no runtime side effects. | AC-19, AC-24 |
| TASK-02 | TASK-01 | Resolve DC-01 and DC-03 grammar/discovery/stage matrix with approved examples; resolve local link/checkbox/duplicate behavior and unsupported input boundary. Produce fixture definitions, not guessed compatibility. | AC-02, AC-04–07, AC-13 |
| TASK-03 | TASK-01 | Resolve DC-02 and related DC-03 graph/metric rules: chains, invalid references, review severity/blocking, recorded contradictions. Preserve owner history rules. | AC-03, AC-08–12 |
| TASK-04 | TASK-01 | Resolve DC-04 CLI/schema/code/rounding/date/task scope, preserving --json, superseded_by and 0/1/2. Record no-config/date-independent default or alternative explicitly. | AC-03, AC-17–19 |
| TASK-05 | TASK-01 | Resolve DC-05/06 safe-open, symlink/mount/exclusion/regular-file/mid-read/limit policy; independently review current toolchain and security evidence. Calibrate candidate caps against fixtures. | AC-13–16, AC-20 |
| TASK-06 | TASK-02–05 | Reconcile specification/contracts/traceability and fixture oracles; remove contradictions. Name approved real dogfood inputs under DC-08. Implementation gate requires separate approval, not inferred from task completion. | AC-01–20, AC-24 |

## Gate 1 — candidate core, only after implementation authorization

| ID | Depends on | Future work / completion evidence | Acceptance mapping |
|---|---|---|---|
| TASK-07 | TASK-06 | Establish minimal stdlib Go build/test baseline, development identity and fixture runner under approved repository controls. No release workflow or code is created by this backlog itself. | AC-01, AC-18–19 |
| TASK-08 | TASK-07 | Implement sole filesystem boundary per approved policy with bounded reads/metadata and safe stable failures; adversarial fixture coverage before parser accesses sources. | AC-13–16 |
| TASK-09 | TASK-08 | Implement supported roadmap/fence/encoding parser with raw statuses/source mapping, duplicate occurrence preservation and visible unsupported failures. | AC-02, AC-04–05 |
| TASK-10 | TASK-09 | Implement explicit spec-dir/stage-file/task/reference inspection only through safe boundary; no recursive repository scan or semantic verification. | AC-05–07 |
| TASK-11 | TASK-09 | Implement replacement validity/chain/cycle model preserving raw/history and no mutation; no dependency redirection. | AC-09–10 |
| TASK-12 | TASK-09, TASK-11 | Implement dependency direct-blocker/cycle/human-review/contradiction results using original references and approved scope/count rules. | AC-08, AC-11–12 |
| TASK-13 | TASK-10–12 | Implement raw/active/verified/task metrics and coverage; invalid/empty/weighted/shared-file cases explicit. | AC-03, AC-09–12 |
| TASK-14 | TASK-04, TASK-13 | Implement common result rendering, text/JSON schema, canonical order/rounding/redaction and threshold exits/fatal envelopes; no independent AI mode. | AC-17–19 |
| TASK-15 | TASK-08, TASK-14 | End-to-end safety/no-side-effect and compatibility fixtures across candidate targets; bounded fuzz/adversarial cases, document residual limits. | AC-01–19 |
| TASK-16 | TASK-15 | Run disclosed cold-start benchmark and approved two-roadmap hand-count dogfood; record observed outcomes and limitations rather than target claims. | AC-20, AC-24 |

TASK-10 and TASK-11 may be independent after parser/safety APIs stabilize; they must reconverge before metrics/rendering. All graph work remains bounded; no “quick” bypass of safe reads to unblock development.

## Gate 2 — release engineering, separately authorized

| ID | Depends on | Future work / completion evidence | Acceptance mapping |
|---|---|---|---|
| TASK-17 | TASK-06 | Resolve DC-07/08 support matrix/runner coverage, pins, public license/security/support/hosting expectations and identity policy. Public docs namespace intent is not a selected license or release approval. | AC-21–24 |
| TASK-18 | TASK-15, TASK-17 | Design and, only when authorized, implement least-privilege tag-triggered test-before-publish build/package workflow using reviewed pins. Audit version/tag mismatch and failing-gate behavior. | AC-21 |
| TASK-19 | TASK-16, TASK-18 | Produce candidate archives/checksum/version/target evidence and execute runtime smoke on every advertised target; no publishing authorization inferred. | AC-22 |
| TASK-20 | TASK-19 | Complete independently exercised no-Go/source rootless install/upgrade docs, OS limits/caveats, immutable correction policy and any approved rollback/identity claims. | AC-23–24 |
| TASK-21 | TASK-20 | Present tested release candidate and complete evidence to owner; first stable tag created only by owner after explicit RC approval. Publishing requires its own authorized action. | AC-21–24 |

## Future work — not part of v1

| ID | Depends on | Boundary |
|---|---|---|
| TASK-22 | Separate future scope approval, stable core contract | Design optional supported-extension mutation interface with each-ID/change/path approval, source reapproval and persisted reread. Not a v1 runtime dependency. AC-25 applies only if approved. |
| TASK-23 | Separate future scope/demand/maintenance approval | Evaluate deferred dashboards/trends/notifications/digests/multirepo/package-manager ideas. Do not create background stores/listeners/sends/schedules merely because listed. |

## Definition of done for a future task

Observed evidence matches its approved contract and acceptance cases; traceability updated, unsupported behavior disclosed, no unresolved decision silently encoded, no claimed success without execution. Product tasks require implementation/test evidence. Documentation-only checks cannot mark product tasks complete. License, root safety, malformed metrics and support promises cannot be accepted through implication or omission.
