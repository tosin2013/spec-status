### 001 — Preserve original feature [status: implemented]
- Description: Historical original feature.
- Outcome: Original outcome remains recorded.
- Scope: Original scope is retained.
- Spec dir: specs/001-original
- superseded-by: 002

### 002 — Replacement feature [status: verified]
- Description: Replacement scope can differ.
- Outcome: Recorded as verified, not independently checked by this CLI.
- Scope: Replacement scope.
- Spec dir: specs/002-replacement

### 003 — Follow-on work [status: planned]
- Description: Still refers to the original requirement.
- Outcome: Requires a scope decision.
- Scope: Follow-on.
- Depends on: 001
