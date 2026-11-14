# Spec Status constitution — draft 0.1

**Documentation only.** Core principles below preserve the requirements baseline; technical recommendations in linked documents remain unapproved PROFILE-R1 choices. No implementation or release is authorized by this text.

## Article I — Truthful evidence

Report recorded state and mechanically established inconsistencies, never claim semantic correctness. Preserve the distinction between checked tasks, recorded implemented work and recorded verified outcomes. Do not infer a merge, passing tests or independent review from a label, prose or file presence. State denominators, skipped inputs and incomplete coverage. Marketing and README claims must not outrun measurable behavior.

## Article II — Deterministic, standalone reporting

The roadmap and explicitly linked supported files are the truth source for roadmap reporting. Human output and versioned JSON must derive from one result. Identical supported inputs and effective configuration must give identical results. Preserve document order for items, stable diagnostic ordering and no volatile default timestamps. Do not require the roadmap extension, AI, credentials or an agent to report status.

On 8 November 2026 the owner amended the earlier roadmap-required rule at product-direction level only. A supported roadmap report and a roadmap-free feature inventory are distinct. Rollout is not decided, and the current CLI still fails when the roadmap file is missing. The record is [two-mode-direction.md](../../specs/001-spec-status/two-mode-direction.md). It does not approve command syntax or an implementation. Under Governance item 4, this amendment records:

- Affected requirements: this article's truth-source wording (now scoped to roadmap reporting), [spec.md](../../specs/001-spec-status/spec.md) §7, the supported-input paragraph in [roadmap-format.md](../../specs/001-spec-status/roadmap-format.md), the missing-roadmap note in [cli-contract.md](../../specs/001-spec-status/cli-contract.md), and the schema note in [json-contract.md](../../specs/001-spec-status/json-contract.md). The Article III rule that unsupported item-like input fails visibly is unchanged.
- Compatibility: none today. CLI behavior, exit codes, and JSON schema 1.0 are unchanged, and a missing roadmap still exits 2. A later inventory output needs its own schema version, which is open.
- Fixture changes: none in this amendment. Later implementation (TASK-25) needs fixtures for missing, empty, unrecognized, and malformed roadmaps, and for inventory.
- Risks: readers may assume an inventory mode exists; it does not. Until implemented, an empty or unrecognized roadmap still gives a zero-item report with exit 0, which README Limits states. How a template version is identified is undecided.
- Scope cuts: none.

On 14 November 2026 the owner approved this amendment recording the rollout (D-2026-11-14-02). The two modes were chosen in D-2026-11-12-08 and implemented in feature 006 ([specs/006-two-mode-inventory](../../specs/006-two-mode-inventory/spec.md)) for v0.3.0. Where the 8 November text says rollout is not decided or that an inventory mode does not exist, this paragraph is the later record. `spec status` remains the roadmap report. It still exits 2 when the roadmap is missing, and that message now suggests `spec inventory`. A readable file that has neither an item heading nor the pinned template's `## Planned Specs` heading is the fatal `SS_ROADMAP_UNRECOGNIZED`, so an unrecognized file is never reported as an empty roadmap. `spec inventory` is a separate, explicit command with its own JSON version (`inventory_schema_version` 1.0). It reads one directory listing of `specs/` and three named files per feature folder. It never reads recursively and never reads the roadmap, and it infers no status, priority or dependency. Under Governance item 4:

- Affected requirements: this article's two-mode paragraph; [roadmap-format.md](../../specs/001-spec-status/roadmap-format.md) §2 (recognition rule); [cli-contract.md](../../specs/001-spec-status/cli-contract.md) (inventory command, missing-roadmap message); [json-contract.md](../../specs/001-spec-status/json-contract.md) (meaning of zero items) and [inventory-json.md](../../specs/001-spec-status/inventory-json.md); the [diagnostics](../../specs/001-spec-status/diagnostics.md) registry (three new codes).
- Compatibility: one pre-1.0 behavior change, listed in the v0.3.0 release notes. A file that is not a recognized roadmap now exits 2 instead of 0. The missing-roadmap message text changed. Schema 1.0 fields are unchanged and the new codes are additive.
- Fixture changes: `testdata/two-mode/` holds hand-counted recognition and inventory fixtures. Symlink, FIFO, sensitive-name and limit cases are built at test time.
- Risks: a roadmap written in another layout now stops instead of reporting zero items. The over-inclusive sensitive-name matcher skips some legitimately named feature folders with an error finding (README Limits). Inventory counts are checkbox counts, not verification.
- Scope cuts: none.

Roadmap parsing and root-confined reads stay in the Go standard library. The owner approved two dependency exceptions: Cobra on 7 November 2026 for command, help, and flag handling, and Lip Gloss for human text layout and restrained color. Those exceptions do not authorize Bubble Tea, dashboards, watch modes, network clients, telemetry, or other frameworks. The approved modules and their required transitive modules are pinned in go.mod and go.sum. Adopting them is not an independent security review.

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

Inferred quarters are excluded from v1, not experimental opt-in. No quarter buckets, percentages, denominators, flags or inference JSON fields. Explicit timing, if supported, stays recorded rather than estimated. Retired heuristic requirements are not implementation tasks. Any future timing experiment needs separate evidence and scope approval. A draft proposal for recorded dates and a separate companion-prompt draft are [timing-evidence.md](../../specs/001-spec-status/timing-evidence.md) and [companion.md](../../specs/001-spec-status/companion.md). They are not amendments to this article and they do not approve implementation.

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
