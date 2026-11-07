# Spec Status constitution — draft 0.1

**Documentation only.** Core principles below preserve the requirements baseline; technical recommendations in linked documents remain unapproved PROFILE-R1 choices. No implementation or release is authorized by this text.

## Article I — Truthful evidence

Report recorded state and mechanically established inconsistencies, never claim semantic correctness. Preserve the distinction between checked tasks, recorded implemented work and recorded verified outcomes. Do not infer a merge, passing tests or independent review from a label, prose or file presence. State denominators, skipped inputs and incomplete coverage. Marketing and README claims must not outrun measurable behavior.

## Article II — Deterministic, standalone reporting

The roadmap and explicitly linked supported files are the truth source. Human output and versioned JSON must derive from one result. Identical supported inputs and effective configuration must give identical results. Preserve document order for items, stable diagnostic ordering and no volatile default timestamps. Do not require the roadmap extension, AI, credentials or an agent to report status.

Roadmap parsing and root-confined reads stay in the Go standard library. The owner approved one dependency exception on 7 November 2026: command, help, and flag handling may use Cobra. That exception does not authorize styling libraries, color, network clients, telemetry, or other frameworks. The approved module and its required transitive modules are pinned in go.mod and go.sum. Adopting Cobra is not an independent security review of those modules.

The build target remains `CGO_ENABLED=0`, with fast cold start and honest platform exceptions. Runtime must perform no network access, subprocess/Git invocation, AI calls, account/credential access or writes, including telemetry, caches or a parallel status store. An explicit completion command may print a script to stdout and must not install it. The process must not create or append `BASH_COMP_DEBUG_FILE`, even when that variable is already set. A sourced completion script is shell behavior, not a write by spec. Development/release tooling is not a runtime dependency.

## Article III — Narrow, safe input

Support a documented pinned heading/bullet dialect, not arbitrary Markdown interpretation. Preserve nine known status tokens and raw unknown labels with warnings. Fenced examples are not live roadmap or task data. Unsupported item-like structures fail visibly; do not manufacture plausible counts.

All input and reference resolution must stay within the selected root. Reject traversal and symlink escape. Bound bytes, file counts, path processing and parsing work. Exclude credential, secret, private-key and answer-key references without reading contents. No recursive inventory. Precise matcher, limits, safe-open and symlink policies are approval gates, not solved merely by these principles.

## Article IV — History cannot become false completion

Retain superseded item IDs, original rows and lifecycle statuses. `superseded-by` is a separate optional relationship, not a new lifecycle status and not upstream-native v0.1.1 behavior. Superseded items stay in raw/historical inventory but do not earn active completion credit. Validate targets, ambiguity, self-links and cycles. Invalid-reference metric behavior must be chosen explicitly, not silently inherited from field presence.

Every dependency on a superseded item requires human review even if its successor is recorded verified. Never redirect, rewrite, auto-satisfy or infer scope equivalence. Changes to this strict v1 history requirement require explicit scope approval; deferring it is not an implementation shortcut.

## Article V — Compatibility and findings

Preserve `--json`, raw `superseded_by`, default failure threshold `error`, and exits 0/1/2. Exit 0 is below the chosen threshold, not proof of correctness; exit 1 is a finding outcome, not a runtime crash; exit 2 is fatal/invocation/configuration/required-input failure. `none` never hides fatal failure. JSON mode emits one JSON document, including fatal errors, with no log contamination on stdout.

Version input dialect, JSON schema and CLI release identity separately. Document CLI/JSON breaking changes and test parity, ordering, null denominators and rounding. Draft field names or codes are not public compatibility promises before approval.

## Article VI — No inferred quarters

Inferred quarters are excluded from v1, not experimental opt-in. No quarter buckets, percentages, denominators, flags or inference JSON fields. Explicit timing, if supported, stays recorded rather than estimated. Retired heuristic requirements are not implementation tasks. Any future timing experiment needs separate evidence and scope approval.

## Article VII — Distribution claims need evidence

Plan downloadable no-Go-needed binaries, exact tag/commit metadata, required test-before-publish gates, archives/checksums, owner-approved release candidate and rootless install/upgrade. Cross-compilation is not runtime testing. No advertised supported OS/architecture without real evidence. Published assets are immutable in policy; fixes get a new version, not silent replacement. Unsigned checksums do not prove identity. Do not claim signing, provenance, reproducibility or security assurance without verified implementation.

## Article VIII — Explicit future mutation authority

V1 status/JSON/CI never mutate or launch an agent. Future extension write-back is separate and optional. Before each mutation obtain explicit approval covering IDs, changes and target path; source change invalidates approval. A calling LLM request is not approval. Use supported extension commands, reread persisted state and distinguish proposed, attempted and confirmed outcomes. Integration absence never breaks the core.

## Governance and amendment

1. Settled requirements are catalogued in [Traceability](../../specs/001-spec-status/traceability.md); proposal IDs and alternatives are in [Decisions](../../specs/001-spec-status/decisions.md).
2. No implementation starts until relevant grammar, metric, safety and interface decisions are recorded and implementation itself is separately authorized.
3. Draft tasks and plan do not constitute that authorization. Release/publication requires its own approval and evidence.
4. Amendments identify affected requirements, compatibility, fixture changes and risks; no silent scope cuts.
5. Dashboards, trends, digests, notifications, multi-repo aggregation and package managers remain deferred; public-project intent is not a commercial commitment or proven market.
6. Original input-format prose has partial provenance limitations; acceptance must use approved fixtures rather than claiming fully recovered source or upstream certification.
