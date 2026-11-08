# Timing and evidence extension — draft proposal

**Draft proposal. Not implemented. Not an approved design choice.** On 7 November 2026 the owner authorized drafting this extension in a voice request. This documentation commit and pull request record that draft. Neither the voice request nor this pull request approves implementation, a date syntax, a severity, a version number, or a change to denominators.

The settled v1 rule is unchanged: inferred quarters stay omitted. Four ordinal buckets cannot answer a delivery date. This draft does not restore quarter inference, quarter percentages, quarter flags, or quarter fields. Explicit timing, if a later decision adopts it, is recorded intent. It is not a forecast and not a guarantee.

Illustrative tokens in this document, including an ID such as `901` or a date shape `YYYY-MM-DD`, are syntax sketches. They are not items, dates, or facts about this repository.

## What would be recorded

Two optional fields are proposed and not adopted:

| Proposed field | Sketch | Meaning if later adopted |
|---|---|---|
| `target` | `YYYY-MM-DD` | A recorded intent that a completion was aimed at that calendar date |
| `timing-updated` | `YYYY-MM-DD` | A recorded intent of when that timing note was updated |

Absent either field means timing is unknown and no completion target is recorded. The CLI must not fill a date from a quarter label, a question ID, a filename, the clock, or a neighboring item.

These choices stay open: the exact syntax, how an invalid or ambiguous date is diagnosed and at what severity, the input-extension version number, and whether an invalid supersession still changes an active denominator. This draft does not pick them.

## Citable evidence model

Human text and JSON would cite one model. A later schema must add it explicitly. Schema 1.0 does not gain these fields in this draft. The input-extension version is a different identity from CLI SemVer and from the JSON schema version. The exact next version is open. No silent breaking change to the versioned JSON payload is proposed.

For each item the citable fields are:

- ID and title
- root-relative source path
- retained line, or an explicit statement that the line is unavailable
- raw status, including an unknown token kept as written
- each dependency: the original ID, the recorded edge state, and the edge's own source path and retained line when the parse kept them
- supersession: the raw relationship, whether it is included in or excluded from active and recorded-verified counts, and the scope or reason for that inclusion or exclusion when the parse supports it
- timing: the recorded target and timing-updated values, or unknown when absent
- source identity for each of those facts

A missing location is reported as unavailable. A line number is not invented. A missing reference stays a diagnostic under the diagnostic policy. It does not prove that code failed, and it does not prove that code is absent.

Human text sanitizes terminal controls before display. JSON keeps the original strings. That split already applies to status text and is part of this proposal for any future timing or evidence lines.

## What this extension would not do

No runtime network, AI, Git, subprocess, write, or crawl. No computed quarter and no ETA. A zero active denominator stays N/A, not 100. A valid superseded row keeps its ID and raw status in history, leaves the active denominator, and is not a success. A dependency on a superseded item stays on that item for human review and is not redirected. Self, missing, and cycle references stay diagnostics. Human-text filters stay views: they do not shrink totals, hide findings, or change the exit.

## Proposed acceptance sketches

These sketches are not acceptance cases of the current build. See the draft block in [acceptance.md](acceptance.md).

- An unknown raw status is kept, warned, not counted as recorded verified, and otherwise remains subject to the active rules.
- A valid superseded item keeps history, ID, and raw status, is excluded from the active denominator, and is not reported as success. Zero active items are not 100 percent.
- A dependency on a superseded item is human review and is not redirected.
- Self, missing, and cycle references are diagnostics.
- Filters preserve full totals, findings, and the exit.
- No date is inferred. A malformed explicit timing value is reported as malformed. It is not repaired into a quarter or an ETA.
