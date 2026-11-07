# Requirements traceability — draft coverage, not test results

This public catalogue maps the reconciled PRD baseline identifiers to the specification, acceptance oracle and dependency-ordered future work. IDs are requirement labels, not execution receipts or claims of approval. Source-only coordination and private review links are intentionally excluded.

**97 baseline entries are accounted for: 91 requirement/user/acceptance entries and six reconciliation notes.** Presence in this table establishes documented coverage, not semantic correctness, implemented behavior, passing tests, independent verification or approval of PROFILE-R1. Exact original input-format prose was partly reconstructed from summaries; approved fixtures must resolve compatibility before implementation. No claim of complete verbatim source recovery is made.

V-03 has a scoped authority update: documentation drafting is permitted now; earlier drafting hold is historical. Implementation/release still require separate authorization and relevant decision closure. DC-08 public-documentation intent is established; repository home is now tosin2013/spec-status, established separately; license and release readiness remain unresolved. T-01–03 retain complete v1 omission of inferred quarters. Retired algorithms are not tasks.

## Individual coverage

Each row lists specification location plus a detailed contract when applicable, acceptance IDs defined in acceptance.md and task dependencies defined in tasks.md. Code/schema/grammar defaults remain proposals, as declared in decisions.md.

| Baseline ID | Preserved outcome | Specification / detailed contract | Acceptance | Future work |
|---|---|---|---|---|
| P-01 | Glanceable recorded roadmap position | FS-01–09; Constitution I–II; [spec.md](spec.md) | AC-01, AC-19, AC-24 | TASK-01, TASK-06, TASK-16 |
| P-02 | Files truth source, deterministic reporting independent of AI | FS-01–09; Constitution I–II; [spec.md](spec.md) | AC-01, AC-19, AC-24 | TASK-01, TASK-06, TASK-16 |
| P-03 | No unsupported uniqueness or market claim | FS-01–09; Constitution I–II; [spec.md](spec.md) | AC-01, AC-19, AC-24 | TASK-01, TASK-06, TASK-16 |
| P-04 | Checked, implemented and verified are distinct recorded evidence | FS-01–09; Constitution I–II; [spec.md](spec.md) | AC-01, AC-19, AC-24 | TASK-01, TASK-06, TASK-16 |
| G-01 | Standalone offline fast compiled binary, extension absent | FS-01, FS-03, FS-07–09; [spec.md](spec.md) | AC-01, AC-17, AC-20, AC-24 | TASK-01, TASK-14–16 |
| G-02 | Human/JSON/CI common facts and explicit bases | FS-01, FS-03, FS-07–09; [spec.md](spec.md) | AC-01, AC-17, AC-20, AC-24 | TASK-01, TASK-14–16 |
| G-03 | Facts versus inference and recorded completion | FS-01, FS-03, FS-07–09; [spec.md](spec.md) | AC-01, AC-17, AC-20, AC-24 | TASK-01, TASK-14–16 |
| G-04 | Personal utility baseline, broader evidence-led possibility | FS-01, FS-03, FS-07–09; [spec.md](spec.md) | AC-01, AC-17, AC-20, AC-24 | TASK-01, TASK-14–16 |
| N-01 | No runtime AI, semantic checks, tests/builds, Git, network, accounts, writes | Scope; Constitution II–VI; [.specify/memory/constitution.md](../../.specify/memory/constitution.md) | AC-04–05, AC-19, AC-24–25 | TASK-01, TASK-06, TASK-15 |
| N-02 | No arbitrary layouts, crawl, background state, dashboard or notifier | Scope; Constitution II–VI; [.specify/memory/constitution.md](../../.specify/memory/constitution.md) | AC-04–05, AC-19, AC-24–25 | TASK-01, TASK-06, TASK-15 |
| N-03 | No effort/value/readiness/date inference | Scope; Constitution II–VI; [.specify/memory/constitution.md](../../.specify/memory/constitution.md) | AC-04–05, AC-19, AC-24–25 | TASK-01, TASK-06, TASK-15 |
| N-04 | Build tooling distinct from runtime | Scope; Constitution II–VI; [.specify/memory/constitution.md](../../.specify/memory/constitution.md) | AC-04–05, AC-19, AC-24–25 | TASK-01, TASK-06, TASK-15 |
| I-01 | Default/override roadmap and root discovery | FS-01–02; candidate grammar; [roadmap-format.md](roadmap-format.md) | AC-02, AC-04–06, AC-13 | TASK-02, TASK-08–10 |
| I-02 | Narrow v0.1.1-shaped heading/bullet input | FS-01–02; candidate grammar; [roadmap-format.md](roadmap-format.md) | AC-02, AC-04–06, AC-13 | TASK-02, TASK-08–10 |
| I-03 | Nine known raw statuses, unknowns preserved/warned | FS-01–02; candidate grammar; [roadmap-format.md](roadmap-format.md) | AC-02, AC-04–06, AC-13 | TASK-02, TASK-08–10 |
| I-04 | Fence exclusion and visible unsupported layout failure | FS-01–02; candidate grammar; [roadmap-format.md](roadmap-format.md) | AC-02, AC-04–06, AC-13 | TASK-02, TASK-08–10 |
| I-05 | Explicit spec-dir and supported file reads only | FS-01–02; candidate grammar; [roadmap-format.md](roadmap-format.md) | AC-02, AC-04–06, AC-13 | TASK-02, TASK-08–10 |
| H-01 | Optional separate replacement relationship preserves history | FS-04–05; history contract; [roadmap-format.md](roadmap-format.md) | AC-09–12, AC-25 | TASK-03, TASK-11–13, TASK-22 |
| H-02 | Raw status and replacement separately in text/JSON | FS-04–05; history contract; [roadmap-format.md](roadmap-format.md) | AC-09–12, AC-25 | TASK-03, TASK-11–13, TASK-22 |
| H-03 | Targets/self/cycles/malformed/ambiguous replacement validation | FS-04–05; history contract; [roadmap-format.md](roadmap-format.md) | AC-09–12, AC-25 | TASK-03, TASK-11–13, TASK-22 |
| H-04 | Historical inventory retained, active completion exclusion | FS-04–05; history contract; [roadmap-format.md](roadmap-format.md) | AC-09–12, AC-25 | TASK-03, TASK-11–13, TASK-22 |
| H-05 | Every superseded dependency human review, never redirect | FS-04–05; history contract; [roadmap-format.md](roadmap-format.md) | AC-09–12, AC-25 | TASK-03, TASK-11–13, TASK-22 |
| H-06 | No v1 replacement authoring; future approved writer only | FS-04–05; history contract; [roadmap-format.md](roadmap-format.md) | AC-09–12, AC-25 | TASK-03, TASK-11–13, TASK-22 |
| F-01 | Summary counts, progress, blockers/findings without inference | FS-03, FS-06; result metrics; [json-contract.md](json-contract.md) | AC-03, AC-06, AC-09–12 | TASK-03–04, TASK-10, TASK-13–14 |
| F-02 | Active recorded-verified numerator/denominator and null empty | FS-03, FS-06; result metrics; [json-contract.md](json-contract.md) | AC-03, AC-06, AC-09–12 | TASK-03–04, TASK-10, TASK-13–14 |
| F-03 | Separate weighted checkbox progress, explicit scope | FS-03, FS-06; result metrics; [json-contract.md](json-contract.md) | AC-03, AC-06, AC-09–12 | TASK-03–04, TASK-10, TASK-13–14 |
| F-04 | Missing/empty evidence not proven completion | FS-03, FS-06; result metrics; [json-contract.md](json-contract.md) | AC-03, AC-06, AC-09–12 | TASK-03–04, TASK-10, TASK-13–14 |
| T-01 | Owner omission of all inferred quarters v1 | FS-02/Scope; T omission boundary; [cli-contract.md](cli-contract.md) | AC-19 | TASK-01, TASK-04, TASK-14–15 |
| T-02 | No inert quarter interface/schema/test obligation | FS-02/Scope; T omission boundary; [cli-contract.md](cli-contract.md) | AC-19 | TASK-01, TASK-04, TASK-14–15 |
| T-03 | Date flag only if justified separately, decision pending | FS-02/Scope; T omission boundary; [cli-contract.md](cli-contract.md) | AC-19 | TASK-01, TASK-04, TASK-14–15 |
| D-01 | Default verified satisfaction, alternative policy uncommitted | FS-05; dependency states; [diagnostics.md](diagnostics.md) | AC-08, AC-11–12 | TASK-03, TASK-12–13 |
| D-02 | Missing/self/cycles/direct blockers and contradictions | FS-05; dependency states; [diagnostics.md](diagnostics.md) | AC-08, AC-11–12 | TASK-03, TASK-12–13 |
| D-03 | Human-review dependency even verified successor | FS-05; dependency states; [diagnostics.md](diagnostics.md) | AC-08, AC-11–12 | TASK-03, TASK-12–13 |
| D-04 | No unblock dates; ordinary unfinished not finding | FS-05; dependency states; [diagnostics.md](diagnostics.md) | AC-08, AC-11–12 | TASK-03, TASK-12–13 |
| X-01 | Structural/stage findings and planned-link information | FS-06; structural registry; [diagnostics.md](diagnostics.md) | AC-02, AC-05–08, AC-11–12 | TASK-02–04, TASK-10, TASK-12, TASK-14 |
| X-02 | Only explicit local reference grammar, not prose/URL/anchor | FS-06; structural registry; [diagnostics.md](diagnostics.md) | AC-02, AC-05–08, AC-11–12 | TASK-02–04, TASK-10, TASK-12, TASK-14 |
| X-03 | Checked missing target error, unchecked info | FS-06; structural registry; [diagnostics.md](diagnostics.md) | AC-02, AC-05–08, AC-11–12 | TASK-02–04, TASK-10, TASK-12, TASK-14 |
| X-04 | Stable safe diagnostic source/severity/code | FS-06; structural registry; [diagnostics.md](diagnostics.md) | AC-02, AC-05–08, AC-11–12 | TASK-02–04, TASK-10, TASK-12, TASK-14 |
| C-01 | Working CLI and explicit root/roadmap/threshold surface | FS-07; invocation/exit contract; [cli-contract.md](cli-contract.md) | AC-01, AC-18–19 | TASK-04, TASK-14–15 |
| C-02 | Same human/JSON facts, no AI-specific mode | FS-07; invocation/exit contract; [cli-contract.md](cli-contract.md) | AC-01, AC-18–19 | TASK-04, TASK-14–15 |
| C-03 | Error-default thresholds and exits 0/1/2 | FS-07; invocation/exit contract; [cli-contract.md](cli-contract.md) | AC-01, AC-18–19 | TASK-04, TASK-14–15 |
| J-01 | Versioned JSON fields and raw relationships, no quarters | FS-07; proposed schema 1.0; [json-contract.md](json-contract.md) | AC-03, AC-17–19 | TASK-04, TASK-13–15 |
| J-02 | Deterministic document order and sorted findings/rounding | FS-07; proposed schema 1.0; [json-contract.md](json-contract.md) | AC-03, AC-17–19 | TASK-04, TASK-13–15 |
| J-03 | One JSON fatal document, no stdout logs or secret text | FS-07; proposed schema 1.0; [json-contract.md](json-contract.md) | AC-03, AC-17–19 | TASK-04, TASK-13–15 |
| J-04 | Compatibility rules and parity | FS-07; proposed schema 1.0; [json-contract.md](json-contract.md) | AC-03, AC-17–19 | TASK-04, TASK-13–15 |
| S-01 | Stdlib Go, CGO_ENABLED=0, no Go on user runtime | FS-08; sole I/O boundary; [plan.md](plan.md) | AC-13–16, AC-19 | TASK-05, TASK-08, TASK-15 |
| S-02 | No runtime side effects or hidden store | FS-08; sole I/O boundary; [plan.md](plan.md) | AC-13–16, AC-19 | TASK-05, TASK-08, TASK-15 |
| S-03 | Root confinement, traversal and symlink escape rejection | FS-08; sole I/O boundary; [plan.md](plan.md) | AC-13–16, AC-19 | TASK-05, TASK-08, TASK-15 |
| S-04 | Sensitive exclusions with safe skipped-reference finding | FS-08; sole I/O boundary; [plan.md](plan.md) | AC-13–16, AC-19 | TASK-05, TASK-08, TASK-15 |
| S-05 | Bounded file/parse work, explicit failures, no crawl | FS-08; sole I/O boundary; [plan.md](plan.md) | AC-13–16, AC-19 | TASK-05, TASK-08, TASK-15 |
| V-01 | Retained v1 reporting, strict history and release requirements | FS-09; readiness gates; [decisions.md](decisions.md) | AC-20–24 | TASK-01–06, TASK-17–21 |
| V-02 | Owner-approved RC before owner-created first stable tag | FS-09; readiness gates; [decisions.md](decisions.md) | AC-20–24 | TASK-01–06, TASK-17–21 |
| V-03 | Decision gates; updated drafting authority, implementation separate | FS-09; readiness gates; [decisions.md](decisions.md) | AC-20–24 | TASK-01–06, TASK-17–21 |
| V-04 | No validated effort estimate or delivery commitment | FS-09; readiness gates; [decisions.md](decisions.md) | AC-20–24 | TASK-01–06, TASK-17–21 |
| W-01 | Separate supported-extension writer, roadmap sole truth | Separate future write-back boundary; [phase2.md](phase2.md) | AC-19, AC-25 | TASK-22 |
| W-02 | Each mutation ID/change/path approval and source reapproval | Separate future write-back boundary; [phase2.md](phase2.md) | AC-19, AC-25 | TASK-22 |
| W-03 | Invocation/version designed separately; integration absence isolated | Separate future write-back boundary; [phase2.md](phase2.md) | AC-19, AC-25 | TASK-22 |
| W-04 | Persisted reread, proposed/attempted/confirmed distinctions | Separate future write-back boundary; [phase2.md](phase2.md) | AC-19, AC-25 | TASK-22 |
| W-05 | AI/extension isolated; replacement authoring not v1 | Separate future write-back boundary; [phase2.md](phase2.md) | AC-19, AC-25 | TASK-22 |
| O-01 | Deferred IDE dashboard, separate scope approval | Deferred scope only; [phase2.md](phase2.md) | AC-24–25 | TASK-23 |
| O-02 | Deferred historical trends/store | Deferred scope only; [phase2.md](phase2.md) | AC-24–25 | TASK-23 |
| O-03 | Deferred chat notifications, no implicit send authority | Deferred scope only; [phase2.md](phase2.md) | AC-24–25 | TASK-23 |
| O-04 | Deferred scheduled digest, no implicit schedule | Deferred scope only; [phase2.md](phase2.md) | AC-24–25 | TASK-23 |
| O-05 | Deferred multirepo aggregation | Deferred scope only; [phase2.md](phase2.md) | AC-24–25 | TASK-23 |
| B-01 | CLI SemVer separate input/JSON and honest version metadata | FS-09; future gated build/release; [build-release.md](build-release.md) | AC-21–23 | TASK-17–21 |
| B-02 | Exact owner tag, pinned reviewed tools, tests then CGO0 build | FS-09; future gated build/release; [build-release.md](build-release.md) | AC-21–23 | TASK-17–21 |
| B-03 | Platform choices unresolved, crosscompile not runtime proof | FS-09; future gated build/release; [build-release.md](build-release.md) | AC-21–23 | TASK-17–21 |
| B-04 | Versioned archives/checksums and gated Release assets | FS-09; future gated build/release; [build-release.md](build-release.md) | AC-21–23 | TASK-17–21 |
| B-05 | Least-privilege publish job, immutable asset policy | FS-09; future gated build/release; [build-release.md](build-release.md) | AC-21–23 | TASK-17–21 |
| B-06 | Unsigned integrity is not identity; no unverified provenance claims | FS-09; future gated build/release; [build-release.md](build-release.md) | AC-21–23 | TASK-17–21 |
| R-01 | Future direct GitHub Release channel | FS-09; rootless direct distribution; [build-release.md](build-release.md) | AC-22–24 | TASK-17, TASK-20–21 |
| R-02 | No-Go source-free rootless install/upgrade guide | FS-09; rootless direct distribution; [build-release.md](build-release.md) | AC-22–24 | TASK-17, TASK-20–21 |
| R-03 | Package-manager channels deferred | FS-09; rootless direct distribution; [build-release.md](build-release.md) | AC-22–24 | TASK-17, TASK-20–21 |
| RR-01 | Honest declared-state reporting, not independent merge verification | Reconciled claim/decision boundaries; [decisions.md](decisions.md) | AC-09–19, AC-21–24 | TASK-01–06, TASK-17 |
| RR-02 | No silent extra mandatory supersession provenance fields | Reconciled claim/decision boundaries; [decisions.md](decisions.md) | AC-09–19, AC-21–24 | TASK-01–06, TASK-17 |
| RR-03 | Bounded evidence not proof all roadmaps lack timing | Reconciled claim/decision boundaries; [decisions.md](decisions.md) | AC-09–19, AC-21–24 | TASK-01–06, TASK-17 |
| RR-04 | Current safety/toolchain evidence required, not historic pin assumption | Reconciled claim/decision boundaries; [decisions.md](decisions.md) | AC-09–19, AC-21–24 | TASK-01–06, TASK-17 |
| RR-05 | Support and reproducibility evidence distinct from promises | Reconciled claim/decision boundaries; [decisions.md](decisions.md) | AC-09–19, AC-21–24 | TASK-01–06, TASK-17 |
| RR-06 | Scope/effort limited, optional integration not core | Reconciled claim/decision boundaries; [decisions.md](decisions.md) | AC-09–19, AC-21–24 | TASK-01–06, TASK-17 |
| U-01 | Developer mid-session status use case | Scenarios/users; [spec.md](spec.md) | AC-01, AC-17–19, AC-22–25 | TASK-01, TASK-06, TASK-14–16, TASK-20, TASK-22 |
| U-02 | Planner recorded progress/timing, no inferred quarter | Scenarios/users; [spec.md](spec.md) | AC-01, AC-17–19, AC-22–25 | TASK-01, TASK-06, TASK-14–16, TASK-20, TASK-22 |
| U-03 | Reviewer mechanical blocker/reference findings | Scenarios/users; [spec.md](spec.md) | AC-01, AC-17–19, AC-22–25 | TASK-01, TASK-06, TASK-14–16, TASK-20, TASK-22 |
| U-04 | Calling LLM receives same structured facts | Scenarios/users; [spec.md](spec.md) | AC-01, AC-17–19, AC-22–25 | TASK-01, TASK-06, TASK-14–16, TASK-20, TASK-22 |
| U-05 | CI threshold/exit consumer | Scenarios/users; [spec.md](spec.md) | AC-01, AC-17–19, AC-22–25 | TASK-01, TASK-06, TASK-14–16, TASK-20, TASK-22 |
| U-06 | Binary user without Go/root/source checkout | Scenarios/users; [spec.md](spec.md) | AC-01, AC-17–19, AC-22–25 | TASK-01, TASK-06, TASK-14–16, TASK-20, TASK-22 |
| U-07 | Future owner-approved persisted mutation | Scenarios/users; [spec.md](spec.md) | AC-01, AC-17–19, AC-22–25 | TASK-01, TASK-06, TASK-14–16, TASK-20, TASK-22 |
| A-01 | Cold target <100ms on disclosed 100-item/10,000-task fixture | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-20 | TASK-06, TASK-15–16, TASK-19–21 |
| A-02 | JSON bytes, text parity, rounding/order/fatal document | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-17–18 | TASK-06, TASK-15–16, TASK-19–21 |
| A-03 | Statuses/unknowns/empty/tasks/fences/layout cases | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-02–05 | TASK-06, TASK-15–16, TASK-19–21 |
| A-04 | Dependency/cycle/contradiction/satisfaction cases | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-08, AC-11–12 | TASK-06, TASK-15–16, TASK-19–21 |
| A-05 | Lifecycle/explicit path/future-output cases | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-05–07 | TASK-06, TASK-15–16, TASK-19–21 |
| A-06 | Safety/limit/no-effects/extension-absent cases | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-01, AC-13–16, AC-19 | TASK-06, TASK-15–16, TASK-19–21 |
| A-07 | History/exclusions/invalid replacement/human review cases | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-09–12 | TASK-06, TASK-15–16, TASK-19–21 |
| A-08 | No inferred-quarter feature, explicit timing retained | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-19 | TASK-06, TASK-15–16, TASK-19–21 |
| A-09 | Every threshold and 0/1/2 case | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-18 | TASK-06, TASK-15–16, TASK-19–21 |
| A-10 | Future downloadable no-Go assets with runtime evidence | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-21–23 | TASK-06, TASK-15–16, TASK-19–21 |
| A-11 | Required-check failure gates, immutable assets, RC/install evidence | Acceptance oracle, not results; [acceptance.md](acceptance.md) | AC-21–24 | TASK-06, TASK-15–16, TASK-19–21 |

## Additional candidate refinements (not source-settled requirements)

DC-01–08 and D-01A–D-06A in decisions.md cover exact grammar/path precedence, valid-only chain exclusions, HUMAN_REVIEW counting/stage matrix, candidate JSON 1.0/rounding/no-config/no-clock, conservative symlink/mount/secret matcher and numeric budgets, small tested matrix/identity policy, and dogfood/public-readiness choices. They are mapped through AC-02–24 and TASK-02–06/17–21; their repetition across documents does not ratify them.

## Coverage and contradiction checklist

- Human/JSON one model, raw unknowns, null denominators, weighted tasks and deterministic ordering are shared across specification and contracts.
- Strict supersession is retained, not optionally cut; invalid-reference metrics are explicitly proposed, not smuggled into source authority. No silent dependency redirect, even verified successor.
- Inferred quarters omitted from interface, schema, algorithms, denominators and future implementation backlog. Explicit timing is not transformed into a new forecast/evidence source.
- Core stdlib/offline/read-only, safety-constrained and extension-independent. Future writer approval/reread isolated and no deferred dashboard/notification feature promoted.
- Compile versus runtime verification and unsigned integrity versus identity remain distinct. A development CLI and implementer-run tests now exist. They are not independent verification, a release, a supported platform list, a signature, or a chosen license.
- README/docs stage does not remain under an obsolete drafting hold; implementation and release gates remain. Public-facing text includes no private coordination/source handles/paths.

## Implementation slice

The development CLI covers the read-only reporting path: narrow roadmap parsing, raw statuses, separate checkbox and recorded-verified counts, history-safe supersession, human-review dependencies, shared text/JSON results, exits 0/1/2, and root-confined reads. In-repo tests exercise malformed input, deterministic output, and adversarial filesystem cases. Those tests were run by the implementer.

Still open as product work, not silently dropped: owner approval of DC-01–08, real dogfood roadmaps, release identity and multi-platform execution, signing, license, and phase-two mutation. Acceptance rows AC-21–AC-25 are not met. AC-20 has one host measurement and is not a passed certification.

Documentation review still must assess details, not only identifier counts. Cases where ambiguity remains (grammar/source provenance, metadata races, unsafe skips versus fatal, duplicate counts and stage matrix) are labelled as implementation blockers or candidate policies rather than falsely resolved requirements.
