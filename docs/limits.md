# Known limits

The short list is in the [README](../README.md#limits). This page has every known limit in detail.


Under the acceptance bar (decision D-2026-11-07-09), a wrong result on input outside the [supported shapes](../specs/002-upstream-format-conformance/contracts/shape-inventory.md) is non-blocking if it fails visibly or is listed here. Every review finding's classification is in [reclassification.md](../specs/002-upstream-format-conformance/reclassification.md).

- Only the documented dialect is read. Under an item heading, every line must be one of these:
  - a recognized or descriptive field;
  - indented continuation;
  - a blank line;
  - a fence;
  - a comment;
  - a thematic break.

  Anything else stops the report with an error. Outside items, prose is not checked. However, anywhere outside fences, comments and code spans, a line with a `Depends on:`, `Spec dir:` or `superseded-by:` label that isn't an accepted field line under an item stops the report with an error.
- Task files are free text in a narrow dialect:
  - Only recognized checkboxes in ordinary lists count.
  - Layouts whose count would depend on block-quote, list-nesting or raw HTML rules stop the report with an error. The cases are in roadmap-format.md §3.
  - Only inline links with a simple local path are checked. Other link forms are reported as unsupported, not checked. Link-like text inside inline code is ignored.
- Reads are not an atomic snapshot. A change that keeps the same size and modification time may go undetected.
- Release archives exist for linux/amd64, linux/arm64, darwin/arm64 and darwin/amd64, each smoke-tested natively. Other platforms are untested. A Windows build compiles, but it cannot read any file: the safe-read layer has no Windows implementation yet (`CodeUnsupported`), which roadmap item 020 adds.
- On Darwin there is no `/proc/self/fd` check after a file is opened, so a parent directory swapped between the checks and the open is not detected. Reads there rely on `O_NOFOLLOW`, the symlink-component check and the device check.
- The sensitive-name matcher is incomplete by design. For example, `id_rsa` is not matched.
- A hard link inside the root that points at a file outside it can't be told apart from an ordinary file, so it is read. Symlinks are rejected.
- Labels and item headings are matched in the source text, not as rendered (D-2026-11-08-03).
  - A `Depends on`, `Spec dir` or `superseded-by` label goes unrecognized when it is written in any of these ways:
    - with an HTML entity or character reference (`Depends&nbsp;on:`, `superseded&#8209;by:`);
    - with a backslash escape (`superseded\-by:`);
    - wrapped in a link, or with inline HTML inside it;
    - with an invisible character outside the normalized set;
    - split across two lines.
  - Where such a label appears under a nested bullet, in a descriptive value, or in prose outside items, its value is ignored without a diagnostic.
  - An item heading written with entities (`### 00&#50; — B &#91;status: planned]`) is not recognized, and the item is not counted.
- A repeated ID in one `Depends on` list (`002, 002`) is listed twice without a diagnostic; the blocked count still counts the item once. A `Spec dir` written as a Markdown link is kept as raw text and reported as a missing directory, not as unsupported syntax.
- Item 001 is recorded `verified` under owner decision D-2026-11-08-05, after three review rounds by separate AI agent sessions.
  - Each round found new contrived layouts that hid a `Depends on`, `Spec dir` or `superseded-by` value, and each one was fixed.
  - No round came back with zero findings, so further layouts of that kind may exist.
  - No human review was done, and nothing was checked against GitHub's own renderer.
- In task files, a line that starts with an HTML tag is not treated as the start of a raw HTML block in these cases:
  - after a setext heading;
  - after indented code;
  - when there is a `>` inside an attribute.

  So checkboxes that renderers show as raw text may be counted.
- Per-phase task counts (`tasks.phases` in the JSON) see only ATX headings: a phase is a `## Phase …` heading. Setext headings (text underlined with `---` or `===`) neither open nor end a phase, and `## phase 1` (lower case) or `## Phases overview` is not a phase.
- Lookalike normalization covers Unicode spaces, dashes, invisible characters, and common Cyrillic and Greek letters that look like Latin ones. Other homoglyphs in a field label are not recognized.
- Open non-blocking defects from review:
  - A FIFO or directory where a tasks file is expected is reported as absent.
  - An over-long or empty `Spec dir` is reported as unsafe.
  - `--root version --json` is an invocation error.
  - An absolute `--roadmap` through a symlinked prefix such as macOS `/tmp` is rejected.
  - A spec directory rejected as unsafe shows tasks as `not_linked` and isn't counted in `files_skipped`.
- Visible over-strictness: input a renderer would accept, but `spec` rejects.
  - **Fences, comments and HTML:**
    - a fence or comment indented one to three columns whose content is less indented;
    - a `<details>` block inside a task item, closed by an indented `</details>`, with the next task directly after it (it passes with a blank line);
    - outside items, an item heading or field directly after a line that starts with an HTML tag, with no blank line between them, even where renderers would show the heading.
  - **Spurious warning:** an `SS_REFERENCE_UNSUPPORTED` warning for a reference definition after a fence that follows a task.
  - **Field labels and values reported as errors:**
    - `- **Depends on: 001**`, with the bold closing after the value;
    - a `Depends on` note containing a comma (`002 (blocked, waiting)`) or any of `<`, `[`, `"`, a backtick, `\` or `&` is reported as malformed.
  - **Prose containing a field label:** prose with one of those labels followed by a colon stops the report, for example `- Notes: rollout depends on: vendor approval` or a table cell `Depends on: 002`.
  - **Indented lines:**
    - under a metric field, any indented line that isn't plain continuation text (a list, quote, comment, HTML block or code block), even after a blank line;
    - an indented `---` under a descriptive field, which ends the item, so later fields are reported as outside an item;
    - `  2. 003` on a line under a `Depends on` field.
  - **Item headings:**
    - `###` followed by a tab before an item ID;
    - a trailing `###` closing sequence;
    - an item heading indented one to three spaces.
  - **Field bullets and continuations:**
    - the field bullet forms `-` followed by a tab, ` - ` with a leading space, `__label:__` and `1. **label:**`;
    - a lazy continuation of a field line at column 0;
    - a continuation line that starts like an item ID (`002 — `).

The original input-format prose was not fully recovered word for word. This is a working implementation of that interpretation, not a certification of upstream compatibility.
