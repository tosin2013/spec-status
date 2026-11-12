# Two-mode product direction — owner decision, 8 November 2026

**Approved as product direction. Not implemented. Rollout and v1 placement are not decided.** On 8 November 2026 the owner authorized recording this direction in the specs and comparing it with pull request 15. That authorization does not approve command syntax, a parser change, a merge, or a change to roadmap item 001.

The earlier product rule that a roadmap file is required for every report is amended at this direction level only. Pull requests 15 and 16 are merged. The CLI on `main` still treats a missing `.specify/memory/roadmap.md` as fatal `SS_ROADMAP_UNAVAILABLE`. That behavior stays until a later implementation is authorized. The observations below were taken at `a6a6ad04bc861608e3b0bea014ff7dc021542259`, which is an ancestor of the merged pull request 15 head `bb3af1ff8d0ea2eb074a3427714a7a11ff6718ad`.

## What is approved

Supported roadmap input is a finite set of explicitly versioned templates. It is not arbitrary Markdown. A filename, a `.specify` directory, or a Spec Kit slash command does not establish that the bytes match a supported template or version.

Five outcomes stay distinct:

| Content | Meaning |
|---|---|
| Supported | Bytes match a pinned, versioned template the product claims to read |
| Malformed or unsupported inside a supported read | The file was selected as a roadmap and breaks a supported rule, or is item-like but not the supported shape. The result is a located diagnostic |
| Unrecognized | The file is present and readable, and it is not a supported template. This is not the same as an empty supported roadmap |
| Explicitly empty | A supported roadmap with no item rows. Denominators stay null. This is not a missing file |
| Missing | The selected roadmap path is not available. Roadmap mode stays a required-input failure |

An unsupported or unrecognized roadmap is not called invalid Spec Kit. Spec Kit can be installed without the roadmap extension. The diagnostic names the path and location when a location was parsed, and it points to extension and compatibility guidance.

Roadmap reporting and roadmap-free feature inventory are separate modes.

- A supported roadmap report states raw recorded lifecycle statuses, dependencies, supersession, and linked-task checkbox counts. Recorded `verified` is not independent verification. Filters do not hide findings or shrink the full-result denominators.
- Inventory, when there is no roadmap, reports bounded feature folders and whether `spec.md`, `plan.md`, and `tasks.md` are present, missing, unreadable, or unsupported. It reports checked and total task checkboxes under a counting rule that is still open. It does not infer priority, dependency, lifecycle, or that checked boxes mean the implementation is done.
- Requirements checklists are excluded. Their checks are not task progress and not implementation evidence.
- A rejected or unrecognized roadmap does not fall through into inventory. Inventory is offered when the roadmap is absent, as its own explicit mode. It is not a silent substitute for a roadmap that was selected and rejected.
- When someone wants roadmap-level priorities or dependencies, the product may recommend the roadmap extension. Installing that extension does not by itself create priorities or dependency data, and this decision does not say every Spec Kit repository should have it.

Reads stay offline, deterministic, and confined to the selected root, with the existing symlink and sensitive-name exclusions and the existing bounds. Human text and versioned JSON state which mode ran, whether the result is complete, and the evidence. Schema 1.0 is unchanged by this document. A later inventory document needs its own explicit version. That version is open.

Two future issue classes are required as text, not as files created here:

- A supported-format bug: a named template version, the input, and a wrong or silent result on a shape that version claims to support.
- A compatibility request: a real layout that is not in the supported set, asking whether a later versioned template should include it.

## What stays open

Not selected, and not to be treated as approved: command names and flags, whether either mode is the default, how a template version is identified, the precise feature-folder naming rule, the precise checkbox grammar (indentation, marker character, empty boxes, letter case), the inventory JSON fields and exit codes, and whether either mode ships in v1 or later.

## Evidence, not a passing inventory run

These observations used pull request 15's commit. They are not an implementation of this direction.

- [anguriatech/envy](https://github.com/anguriatech/envy) `ad5eae49f28ff2cf86ce2fe2af68364385e82a7c` has a narrative `.specify/memory/roadmap.md`. The CLI returned a complete report, zero items, exit 0. That is the current silent empty result. This direction says an unrecognized file must not be reported as an empty supported roadmap. The distinction is not implemented.
- [tosin2013/animspec](https://github.com/tosin2013/animspec) `be42f7900036e8d216012f6cbcac3c601095f6bb` has no `.specify/memory/roadmap.md`. `specs/ROADMAP.md` is a pipe table. The CLI returned fatal `SS_LAYOUT_UNSUPPORTED` at line 25, exit 2. That matches the current item-like-table rule. It is not a claim that the AnimSpec project is invalid Spec Kit. Its Spec Kit metadata says `1.0.12` and it has no roadmap template.
- [mnriem/spec-kit-java-brownfield-demo](https://github.com/mnriem/spec-kit-java-brownfield-demo) `31c462513fd7b0a4ec5cef0f76258964d033cae0` has Spec Kit `0.14.3` metadata and no roadmap file. The CLI returned fatal `SS_ROADMAP_UNAVAILABLE`, exit 2. A separate manual trace, not a CLI inventory, found one immediate feature folder, `specs/001-piranha-server-management-ui`, with `spec.md`, `plan.md`, and `tasks.md` present. Under a proposed column-0 checkbox reading, `tasks.md` is 24 checked of 24. `checklists/requirements.md` has 16 checks and was excluded. `spec.md` line 5 says `Status: Draft`, which is a template field. That trace does not prove the feature files match the 0.14.3 templates, and it is not an executed inventory.

## Pull request 15

Pull request 15 merged on 8 November 2026. Its head is `bb3af1ff8d0ea2eb074a3427714a7a11ff6718ad`. Pull request 16 merged the same day. Both are part of `main`. The repairs in pull request 15 are inside the supported-roadmap grammar: inline tags inside a wrapped value, a table ending at a list marker, text after a multi-line comment, a block between a metric value and a continuation, a metric label written as a heading, prose such as “Depends on whether…” not treated as a field, and an empty checkbox counted. On Linux x86_64, Go 1.22.2, `CGO_ENABLED=0 go test ./...` passed at `a6a6ad04bc861608e3b0bea014ff7dc021542259`, and those seven tests failed on the base and passed there. Those repairs stay valid for the supported-roadmap mode. They do not implement unrecognized-versus-empty, a missing-roadmap inventory, or a new schema. Implementation of the new mode waits on the open choices above.

After that evidence was taken, the pull request gained one commit, `bb3af1ff8d0ea2eb074a3427714a7a11ff6718ad`. It changes only `README.md`: one Limits line saying a readable roadmap with no item headings reports zero items with exit 0, and that this build cannot tell an empty supported roadmap from an unrecognized file. No Go file changed, so the test evidence above applies to the parser at the new head. The Linux run was not repeated at `bb3af1f`.
