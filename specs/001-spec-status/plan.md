# Technical plan — proposed architecture, no implementation

**Draft PROFILE-R1.** This file remains the design record. A development implementation now exists beside it and uses these proposals as reversible defaults, not as owner approval. This file contains no executable code.

## Architecture and processing boundaries

Propose a small Go standard-library application with logically separate responsibilities, not necessarily a mandated directory tree:

1. **Invocation/configuration:** parse commands and flags with the owner-approved Cobra dependency; establish whether JSON is requested early enough to format invocation errors without mixing Cobra usage text into the JSON envelope; explicit flags over documented defaults, no configuration file in candidate profile. Validate root and threshold without executing anything. Parsing and filesystem code stay outside Cobra.
2. **Root-confined input boundary:** sole component permitted to inspect filesystem state. Normalize proposed path bases, enforce root/symlink/mount/regular-file and sensitive-name policy before reads, account every byte/file/path/work budget. Return stable categorized errors and source metadata, not raw secret contents or platform exception text.
3. **Narrow parser:** stream bounded UTF-8 lines/fence state into roadmap occurrences and explicitly linked supported documents. Separate raw strings/source spans from parsed values; no arbitrary Markdown library, AI or semantic parser. Reject unsupported item-like structures that make inventory unreliable. Preserve duplicate occurrences and unknown statuses.
4. **Graph and structural validation:** build unique-ID lookup plus ambiguity index. Validate dependency and supersession graphs independently with bounded, non-recursive cycle analysis. Record direct links/cycle membership and supersession validity without rewriting input or relationships.
5. **Progress/result model:** compute active eligibility, raw counts, recorded-verified fraction and unique-file checkbox aggregation. Build blockers/human-review states, file-stage/reference findings and coverage. Use integer counts and deterministic rounding; do not derive dates/effort from order.
6. **Renderers and exit policy:** one immutable in-memory result to text and candidate JSON 1.0, diagnostics sorted path/line/code plus documented tie-break. Human text may use the owner-approved Lip Gloss dependency for layout, restrained color, and the static summary panel. Human-text filters and `--explain` select or annotate item rows only; the summary, diagnostics, exits, and JSON payload stay on the full result. JSON stays unstyled and keeps raw field text. Threshold only controls exit 0/1, never finding visibility; fatal exit 2 precedes threshold. No output file, cache, spinner, or hidden runtime process.

This division makes no API/module-path promises and does not choose a repository layout or license. Documentation for this draft is placed at `.specify/memory/constitution.md` and under `specs/001-spec-status/`. Those directories record the draft and leave the repository root clear for future program code. Go package layout, module path, and license remain unchosen.

## Core data model — design only

Roadmap occurrences retain ID/title/raw status, recognized flag, description/outcome/scope, source path/line, direct dependency fields and optional spec-dir/replacement field. Duplicate IDs remain separate occurrence objects; lookup ambiguity is explicit. Avoid maps as a presentation order source.

A replacement record separates supplied raw value, parsed ID, validity, invalid reasons and optionally terminal successor. Lifecycle never mutates because a relationship is supplied. Dependency results retain original target and any displayed successor; satisfaction and human review are independent facts. Direct blockers do not become transitive forecasts.

Task data links counts to source file and coverage; per-item reuse is visible, aggregate unique-file scope explicit. Source absence, excluded unsafe evidence and zero tasks are not interchangeable. Safe source paths are selected-root-relative; unresolved source may be null, not guessed from a private machine path.

Candidate `complete` means structurally interpretable report under declared read scope, not complete project truth or complete evidence verification. Fatal required input makes summary null; optional skips remain visible as findings and coverage. No claim of an atomic repository snapshot.

## Safety design choices that must be settled

- Select a supported patched Go baseline using current evidence. A historic patch number is not by itself an acceptable current baseline. Review root-confined API capabilities and known advisories before choosing mechanisms; no unsafe fallback on another platform.
- Conservative PROFILE-R1 rejects symlink components for linked paths; canonical root treatment and mount handling need approved platform semantics. Lexical prefix comparisons alone cannot enforce confinement, case/volume boundaries or races.
- Do not read excluded secrets to classify them. Apply approved normalized path/name matcher before open and again at the safe-open boundary as appropriate. An unmatched name is not evidence that arbitrary repository content is harmless; input surface stays narrow.
- Regular files only; FIFOs/devices/sockets cannot block indefinitely. Opened descriptors must be validated within the safe boundary. Cap actual read bytes, file count, parsed items/edges/tasks, paths and nesting. Bound work even when files lie about size or change.
- Detect observable input changes during read; PROFILE-R1 invalidates required-input report with exit 2. No retry loop that can hang under continuous edits. Same-size/same-metadata races are a documented residual limitation unless a stronger snapshot design is approved.
- Existence checks for explicit task targets need confined metadata access, not opening every linked file's contents. Skipped sensitive references appear as stable safe diagnostics; avoid echoing credential snippets/absolute secret paths.
- No persistent output, temp file, telemetry, network, subprocess, Git, account/credential retrieval, runtime extension or hidden shell.

## Determinism and complexity

Propose linear parser work in permitted bytes, indexed graph work in item+edge counts, stable sorts for diagnostics/status keys and integer progress arithmetic. Do not recursively walk a repository or recursively traverse unbounded graphs. Deterministic key order and slash-normalized relative paths avoid host-map iteration variance. Platform-specific filesystem behavior still needs fixtures; cross-platform byte parity is a target, not assumed.

PROFILE-R1 is date-independent, no inert `--as-of`; if a real non-quarter use is approved, record its explicit date/input context and update determinism fixtures before coding. Text and JSON use one set of counts and facts; JSON errors cannot include uncontrolled OS wording that makes bytes volatile.

## Validation strategy and design tradeoffs

Use documented synthetic fixtures before real dogfood inputs. Unit-like fixture expectations cover parser/fences/raw statuses; graph fixtures cover duplicates, valid/invalid chains and human review; safety fixtures cover traversal, symlink/trailing slash, mount policy, regular-file and bounded-work cases. Integration harness may observe runtime behavior externally without adding test execution to product runtime. No external harness dependency is silently introduced into stdlib core.

The conservative profile trades tolerance for trustworthy reporting: malformed item-like input can invalidate inventory; invalid replacement cannot hide a row; no root/sensitive bypass. A less strict chosen policy must explain partial coverage and update acceptance explicitly. The file-stage matrix, identity evidence and two-target release profile remain decisions, not root-cause findings or validated APIs.

## Delivery boundary

[Tasks](tasks.md) describe future dependency order, all unstarted. [Acceptance](acceptance.md) defines observable gates. [Build/release](build-release.md) describes later controlled automation; this draft creates no workflow. First implementation needs explicit authorization after contract decisions; first stable release needs an owner-approved tested RC and owner-created tag. Package-manager channels and optional phase-two write-back cannot become accidental prerequisites for the standalone core.
