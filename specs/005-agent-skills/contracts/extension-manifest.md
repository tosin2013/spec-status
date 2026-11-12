# Contract: `extensions/spec-status/extension.yml`

The shape follows `.specify/extensions/roadmap/extension.yml` (Spec Kit 1.2.0). A conforming manifest:

- MUST have `schema_version: '1.0'` and `extension.id: spec-status`.
- MUST set `extension.version` to the release version. The release `verify` job fails if it differs from the tag (research R7).
- MUST set `requires.speckit_version: '>=1.2.0'`.
- SHOULD list `requires.tools: [{name: spec, required: true, version: '>=0.2.0'}]` for documentation. Commands still check `spec version` themselves (FR-004).
- MUST list exactly five `provides.commands`, each with `name`, `file` and `description`:
  `speckit.spec-status.review`, `.explain`, `.reconcile`, `.unblock`, `.verify-readiness`. Each `file` is `commands/<name>.md`.
- MUST NOT declare `requires.extensions` (partners are optional, FR-026), `provides.scripts`, or any mandatory hook.
- MAY declare exactly one hook: `after_implement: {command: speckit.spec-status.reconcile, optional: true, prompt: …, description: …}`.
- MUST be readable by the repository test with a line regex. That means one key per line, and list items as `- name:` / `file:` pairs as in the roadmap manifest; no flow-style mappings for `provides.commands`.

**Checked by:** `internal/extension` (manifest test, version test), and `scripts/check-extension-install.sh`, which confirms that Spec Kit accepts and renders it.
