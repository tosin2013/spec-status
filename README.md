# Spec Status

**Planning/specification stage — documentation draft, 7 November 2026.** There is no working binary, implementation, published release or executed product test evidence in this package. Proposed behavior below is not a claim that commands work today.

Spec Status is a proposed offline Go CLI for a glanceable view of a Spec Kit roadmap: recorded statuses, task progress, dependency blockers and mechanical reference inconsistencies. Working command name: `spec status`.

## What it will and will not report

The intended v1 reads `.specify/memory/roadmap.md` and explicitly linked specification files, without requiring the roadmap extension. It will provide plain text and versioned JSON from the same deterministic result, with findings-based CI exit codes:

- `0`: no findings at the configured threshold;
- `1`: findings at or above that threshold;
- `2`: invocation, configuration, required-input or fatal-parse failure.

A checked task, a recorded `implemented` item and a recorded `verified` item are different evidence levels. This tool will **not** run Git or tests to independently prove merges, implementation correctness or verification. It reports what files record and what structural checks establish.

V1 is intended to be read-only, offline and standard-library-only Go, built with `CGO_ENABLED=0`: no AI, network, subprocess, Git, credentials, file mutations or extension dependency at runtime. Reads will be bounded and confined to a selected project root; sensitive references will be excluded. These are requirements, not completed security assurances.

## Scope

Planned v1 includes narrow `speckit-roadmap` v0.1.1 heading/bullet support, raw status preservation, recorded-verified and checkbox progress reported separately, dependency/reference diagnostics and strict history-preserving optional `superseded-by` links. Dependencies on superseded items always require human review, never automatic redirection.

**Inferred quarters are omitted from v1**, including any opt-in, calculation, denominator or inference JSON field. Explicit timing remains recorded source text where the selected input contract supports it; question IDs are not dates.

Extension-mediated write-back is a separately approved future mode requiring explicit approval before **each** mutation and a persisted reread afterward. Dashboards, historical trends, chat notifications, weekly digests, multi-repository aggregation and package managers are deferred, not promised features.

## Installation and support

Nothing is available to install from this documentation package. No download URL, installation command, released version or supported platform is asserted here. Future direct binary downloads are planned through GitHub Releases, with rootless installation and no Go toolchain required for end users. Exact platforms remain undecided; a small actually tested initial matrix is recommended rather than claiming cross-compiled targets are verified.

Unsigned SHA-256 manifests can detect inconsistency with a manifest, but do not independently authenticate its publisher. Signing/attestation policy remains open. No signatures, reproducible builds, security audit or passing tests are claimed.

## Review the documentation

| Document | Purpose |
|---|---|
| [Constitution](CONSTITUTION.md) | Non-negotiable principles and scope governance |
| [Feature specification](SPECIFICATION.md) | User outcomes, functional requirements and edge cases |
| [Decisions](DECISIONS.md) | Settled boundaries, unapproved recommended profile and implementation blockers |
| [Roadmap format](ROADMAP-FORMAT.md) | Candidate narrow grammar and optional supersession extension |
| [CLI](CLI-CONTRACT.md), [JSON](JSON-CONTRACT.md), [diagnostics](DIAGNOSTICS.md) | Candidate versioned machine/human contracts |
| [Technical plan](TECHNICAL-PLAN.md) | Proposed architecture and safety design, not implementation |
| [Acceptance strategy](ACCEPTANCE.md) | Required fixtures and evidence; no completed test claims |
| [Ordered tasks](TASKS.md) | Dependency-ordered future work, all unstarted |
| [Build/release/distribution](BUILD-RELEASE.md) | Future release gates and truthful integrity/support requirements |
| [Phase two](PHASE2.md) | Approval-isolated future write-back boundary |
| [Requirements/traceability](TRACEABILITY.md) | PRD requirement coverage into specification, tests and tasks |

All design details marked **PROFILE-R1** are coherent recommendations awaiting approval, not ratified requirements. Drafting permission does not authorize implementation, tags, releases or deployment. The original input-format prose is not completely recovered verbatim; this package is a reviewed requirements interpretation and proposed specification, not an upstream compatibility certification.

## Public project and license status

These documents are prepared for a public project under the `tosin2013` namespace. The public repository is [tosin2013/spec-status](https://github.com/tosin2013/spec-status). Repository availability is distinct from implementation or release readiness; this package contains documentation only. A license has **not** been selected; no open-source license, reuse grant or contributor agreement is invented here. Licensing and publication readiness must be resolved before any code distribution claim.

The primary objective is useful deterministic reporting. Broader usefulness is possible but demand and uniqueness are not validated. There is no committed implementation estimate or delivery date.
