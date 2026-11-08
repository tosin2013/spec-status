# Reclassification of acceptance-review findings under D-2026-11-07-09

Task T024, FR-015. Every finding from the seven independent acceptance reviews of feature 001, classified under the bar that replaces D-2026-11-07-06.

**Classes**:
- **fixed**: now read correctly, with a regression test.
- **visible**: now a fatal error or diagnostic, not a silent miscount.
- **limit**: stated in README Limits; non-blocking under the new bar.
- **over-strict**: visible rejection of input a renderer accepts; non-blocking.
- **open blocking**: none remain.

"Upstream shape" means the finding's input matches a row of the shape inventory. "Metric" means it involves a `Depends on`, `Spec dir`, or `superseded-by` value.

## Review 1 (speckit-roadmap-compat, before #8)

| ID | Finding | Upstream shape | Metric | Class | Where |
|---|---|---|---|---|---|
| D1 | Headings at other levels and table rows dropped; fields attached to the item above | no | yes | visible | #8 `strict_test.go` |
| D2 | One-line, unclosed, or indented fences hid later items and tasks | no | yes | visible | #8 |
| D3 | Near-miss field labels ignored (superseded item kept credit) | no | yes | visible | #8 |
| D4 | FIFO swapped in mid-open could hang the read | no | no | fixed | #8 `TestOpenNoFollowDoesNotBlockOnFIFO` |
| D5 | Checkbox cap per file, not in total | no | no | fixed | #8 `TestAggregateCheckboxBudget` |
| D6 | Link titles and reference-style links mishandled | no | no | fixed | #8 `TestTaskLinkForms` |
| D7 | Every edge out of a cycle member labelled `cycle` | no | yes | fixed | #9 `TestCycleLabelsOnlyCycleEdges` |
| D8 | Skipped spec dirs shown as `not_linked`; FIFO tasks.md shown as missing | no | no | limit | README Limits |
| D9 | Invalid supersession not marked in human text | no | no | limit | README Limits |
| D10 | macOS `/tmp` vs `/private/tmp` root; invocation JSON `fail_on` | no | no | limit | README Limits |

## Review 2 (main after #8)

| ID | Finding | Upstream shape | Metric | Class | Where |
|---|---|---|---|---|---|
| D1 | `- superseded-by 002` without a colon ignored | no | yes | visible | #9 prefix rule |
| D2 | GFM table without outer pipes gave 0 items | no | no | visible | #9 `TestTablesWithoutOuterPipes` |
| D3 | Colon-less, backtick-wrapped, block-quoted labels dropped | no | yes | visible | #9 |
| D4 | Bullet item rows ignored | no | no | visible | #9 |
| D5 | Items inside HTML comments counted | no | yes | fixed | #9 `TestHTMLCommentsAreNotLive` |
| D6 | Setext heading did not end an item | no | yes | fixed | #9 `TestSetextHeadingEndsItem` |
| D7 | Cycle labels hid human review | no | yes | fixed | #9 |
| D8 | Invalid supersession chain let dependents be satisfied | no | yes | fixed | #9 `TestInvalidSupersessionStillNeedsHumanReview` |
| D9 | Root discovery passed an unusable `.specify` | no | no | fixed | #9 `TestDiscoveryStopsAtUnusableSpecify` |
| D10 | Scheme and absolute links checked as local paths | no | no | fixed | #9 `TestSchemeAndAbsoluteReferences` |
| D11 | Human text omits `status_counts` and invalid supersession | no | no | limit | README Limits |
| D12 | Tasks-file fatal cited the roadmap path | no | no | fixed | #12 `TestTasksFatalCitesTasksFile` |
| D13 | FIFO or directory targets reported as absent | no | no | limit | README Limits |
| D14 | Long or empty `Spec dir` reported as unsafe | no | no | limit | README Limits |
| D15 | Deep-indented prose tripped the nesting limit | no | no | fixed | #9 |
| D16 | `--root version --json`; invocation `fail_on`; empty `--root`/`--roadmap` use defaults | no | no | limit | README Limits |
| D17 | macOS `/tmp` symlinked prefix rejected | no | no | limit | README Limits |
| D18 | Some checkbox spellings not counted | no | no | fixed | #10, #11 |

## Review 3 (main after #9)

| ID | Finding | Upstream shape | Metric | Class | Where |
|---|---|---|---|---|---|
| D1 | Unclosed fence inside a list item swallowed later items | no | yes | visible | #10 `TestFenceInsideListItemMustClose` |
| D2 | Unicode lookalikes (U+2011, U+00A0, U+200B) slipped past field rules | no | yes | visible | #10 `TestClosedItemBody` |
| D3 | GFM checkbox spellings not counted | no | no | fixed | #10 `TestCheckboxSpellingsLikeGFM` |
| D4 | Nested-bracket, shortcut, and continuation-line links unchecked | no | no | visible | #10 `TestUncheckableTaskLinksAreReported` |
| D5 | One tasks file double-counted under two paths | no | no | fixed | #10 `TestOneTasksFileUnderTwoPathsCountsOnce` |
| D6 | Non-regular files labelled missing | no | no | limit | README Limits |
| D7 | macOS `/tmp` prefix | no | no | limit | README Limits |
| D8 | Human text omits facts | no | no | limit | README Limits |
| N2 | Bare `Depends on:` under an item read as prose | no | yes | visible | #10 closed item body |

## Review 4 (main after #10)

| ID | Finding | Upstream shape | Metric | Class | Where |
|---|---|---|---|---|---|
| D1, D4 | Fence closer at a different indentation ended a block early | no | yes | visible | #11 `TestFenceCloserMustMatchOpenerIndent` |
| D2, D3, D5 | Comment boundaries in list items and at 4-space indentation | no | yes | visible | #11 `TestCommentBoundaries` |
| D6 | Checkbox in indented code counted | no | no | visible | #11 `TestIndentedCodeInTasks` |
| N1 | Text after a one-line comment hid a dependency | no | yes | visible | #11 |
| N2 | Hard link to a file outside the root is read | no | no | limit | README Limits |
| N3 | Wrong message for an indented fence | no | no | fixed | #11 |
| N4 | `--json=false` with an invalid flag still emits JSON | no | no | limit | README Limits |

## Review 5 (PR #11 at c3f450e)

| ID | Finding | Upstream shape | Metric | Class | Where |
|---|---|---|---|---|---|
| B1 | Empty comment `<!-->` hid content | no | yes | visible | #11 `TestEmptyCommentsAreFatal` |
| B2 | Block-quoted fences and checkboxes counted | no | no | visible | #11 |
| B3 | Indented code inside list items counted | no | no | visible | #11 |
| B4 | Nested task after an indented details paragraph dropped | **yes** (hand edit covered by the hand-edits fixture: `task-nested` after a details paragraph) | no | fixed | #11 `TestNestedTasksAfterIndentedParagraph` |
| B5 | Raw `<style>`-type blocks read as live | no | yes | visible | #11 `TestRawHTMLBlocksAreFatal` |
| N1 | Two `[status:` markers | no | no | visible | #11 |
| N2 | Cyrillic homoglyph in a nested label | no | yes | visible | feature 002 `TestLetterLookalikesInMetricLabels` |
| N3 | Fences or comments indented 1–3 columns with less-indented content are fatal | no | no | over-strict | visible |

## Review 6 (PR #11 at af59ab2)

| ID | Finding | Upstream shape | Metric | Class | Where |
|---|---|---|---|---|---|
| D1 | `Depends on` wrapped onto an indented line lost IDs | **yes** (`field-continuation`) | yes | fixed | #12 `TestWrappedFieldValuesAreJoined` |
| D2 | Checkboxes inside an HTML block counted | no | no | visible | #12 `TestHTMLBlocksHideStructure` |
| D3 | Column-0 constructs did not close list items | no | no | visible | #12 `TestBlockConstructsCloseListItems` |
| D4 | Other GFM list rules | no | no | visible | #12 `TestRemainingGFMListRules` |
| D5 | Item heading after a type-7 HTML line counted | no | yes | visible | #12 |
| — | Tasks-file fatal lacked a line number | no | no | fixed | #12 |

## Review 7 (PR #12 at c4824a6)

| ID | Finding | Upstream shape | Metric | Class | Where |
|---|---|---|---|---|---|
| B1 | Lone HTML tag after a heading hid content only from renderers | no | yes (roadmap) | visible | #12 `TestLoneTagAfterHeadingStartsHTMLBlock` |
| B2 | Column-0 paragraph after a fence inside a task did not close it | no | no | visible | #12 `TestParagraphAfterFenceClosesTask` |
| B3 | HTML block types 3 to 5 not recognized | no | yes | visible | #12 `TestOtherHTMLBlockTypesAreFatal` |
| B4 | Tabs counted as a fixed 4 columns | no | no | fixed | #12 `TestTabsAdvanceToTabStops` |
| B5 | Stale list-item stack | no | no | visible | #12 `TestMarkerClosesItemsItDoesNotReach` |
| obs. | Bare indented `superseded-by:` under an item ignored | no | yes | visible | #12 `TestBareMetricLabelUnderItemIsFatal` |
| N1 | `<details>` inside a task item closed by an indented `</details>` with the next task directly after is fatal | no | no | over-strict | visible; passes with a blank line |
| N2 | Reference definition after a fence that follows a task gets a spurious `SS_REFERENCE_UNSUPPORTED` warning | no | no | over-strict | visible warning |

## Review 8 (PR #13 at 9e8dff0, the first review under D-2026-11-07-09)

The reviewer recounted four fixtures by hand (all matched), compared tasks counts with pandoc on 11 fixture files, the pinned template, and 23 realistic `/speckit-tasks` variations (all matched apart from ordered-list checkboxes, which pandoc does not render), and found no wrong count on any upstream shape. Rule (b) failed:

| ID | Finding | Upstream shape | Metric | Class | Where |
|---|---|---|---|---|---|
| D1 | An item heading indented one to three spaces after a field line was joined to that field; the next item's `Spec dir` and `superseded-by` attached to the item above | no | yes | visible | `TestBlockStartsAreNotFieldContinuations` |
| D2 | An indented HTML block under an item did not hide its fields or item headings | no | yes | visible | `TestHTMLBlocksUnderItems` |
| D3 | A `Depends on` value continued after a blank line lost the continued IDs | no | yes | visible | `TestMetricValueAfterBlankLineIsFatal` |
| D4 | A struck-through item heading, which speckit.roadmap.write produces for superseded entries, was fatal | **yes** (new hand-edit shape `item-heading-struck`) | no | fixed | `TestStruckThroughItemHeading`; hand-edits fixture |
| — | `expected.json` did not assert `spec_dir` or `superseded_by` | — | yes | fixed | comparator and all five fixtures |
| — | The shape inventory did not say how paragraph text is treated | — | no | fixed | inventory wording |

The previous version of this table said "Open blocking: none". Review 8 showed that was wrong for rule (b); the corrected result is below.

## Review 9 (main at ed5d88e, narrowed to D-2026-11-07-09)

No rule (a) defect: a realistic `/speckit-tasks` file, 2/3/4-space, tab, and 8-space nesting, lazy wraps, details paragraphs, fences, comments, quotes, breaks, and `<details>` all matched GFM. Rule (b) failed in four layouts:

| ID | Finding | Upstream shape | Metric | Class | Where |
|---|---|---|---|---|---|
| D1 | An inline tag (`<br>`, `<span>`) on a continuation line cut a `Depends on` or `superseded-by` value short (a regression from the review-8 D1 fix) | no | yes | visible | `TestInlineTagInsideWrappedValue` |
| D2 | A field line containing `\|` after an indented table was swallowed by the table | no | yes | fixed | `TestTableEndsAtListMarker` |
| D3 | Text after a multi-line comment's `-->` under an item was dropped | no | yes | visible | `TestTextAfterMultiLineCommentIsFatal` |
| D4 | A comment, quote, `<details>`, or heading between value lines dropped the later line | no | yes | visible | `TestBlockSeparatesMetricValue` |
| D5 | A metric label written as a heading was ignored | no | yes | visible | `TestMetricLabelAsHeadingIsFatal` |
| D6 | Prose bullets starting "Depends on …" were fatal; `- **Superseded by:** 002` fatal; a placeholder `Spec dir` an error | no | no | over-strict (prose fixed; the other two listed in README Limits) | `TestDependsOnProseIsNotAField` |
| D7 | An empty `- [ ]` was not counted | no | no | fixed | `TestEmptyCheckboxCounts` |

## Line-level label rule (D-2026-11-08-02)

Each review from 7 to 9 found another layout that hid or split a metric field. Instead of a tenth case list, a check after parsing makes any visible line with a `Depends on`, `Spec dir`, or `superseded-by` label and a colon fatal unless it was accepted as a field line under an item (`labelrule_test.go`). Of its twelve fatal cases, eight were accepted silently on `main` at 5135231, mostly outside items (tables, text after a comment, inline HTML, prose, lookalike spacing) plus a label inside a `Notes` value. Real roadmaps and every conformance fixture give byte-identical output. Over-strictness: prose that contains such a label with a colon is fatal; it is listed in README Limits.

## Result

- Upstream-shape findings: 2 (review 5 B4, review 6 D1). Both fixed and covered by fixtures.
- Metric findings: all fixed or visible. Review 5 N2 (letter homoglyphs) was the last; while writing this table it was found to fall under D-2026-11-07-09 (b), so feature 002 extended lookalike normalization to common Cyrillic and Greek letters, with `TestLetterLookalikesInMetricLabels`.
- Every finding classed as a limit or over-strict is listed in README Limits.
- Open blocking after review 9 fixes and the line-level label rule: **none known**. A further independent review (T026) is needed to confirm it before roadmap item 001 is recorded verified.
