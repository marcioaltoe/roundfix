---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A large authored data file is reviewed apart, in Data Review Segments

A repository can declare in `pre_pr_review.data_paths` the files it writes by
hand as data, such as rule catalogs. The pre-PR review takes those files out of
the code diff and reviews their diff apart, in up to four prompts of their
own in the same review session, so the reviewer still reads them and they no
longer push the code diff past its bound.

## Context

On 2026-11-05 Fiscus reported that an authored catalog of fiscal rules,
`icms-es-v2.json` at about 550 KB, passed the review bound on its own. The
only way through was to declare it in `delivery.derived_paths` (Fiscus pull
request #263), which leaves it out of the review as `generated`, so the
reviewer no longer reads the fiscal rules. `pre_pr_review.omit_paths`
(ADR-0352) leaves a file out the same way. Neither is right for a file that is
authored, not generated. The maintainer answered on 2026-11-05:

> Uma categoria nova para arquivos grandes de dados escritos pelo projeto,
> revisados à parte em vez de omitidos.

## Decision

**A declared Review Data Path is reviewed apart.** Project Config gains
`pre_pr_review.data_paths`, a list of clean repository-relative paths,
directories ending in `/` and whole-path Go `path.Match` patterns, validated
like `pre_pr_review.omit_paths`. The review reads it from the Project Config
at the candidate's base commit, never from the head. A changed file it matches
gets the reason `data-path` after every existing reason, so a file another
rule omits stays omitted. It is left out of every code diff the review sends,
in every round, and the code prompt names it as reviewed apart.

**Its diff goes into Data Review Segments.** The diff of the data files from
the start of the round's code diff to the head is cut into at most four
segments, each sent after the code prompts, in the same session, with its own
two-line header. Whole files are packed in path order while a segment fits
the provider's bound, the `codex` byte bound or the `claude` token budget; a
file that does not fit alone is cut at its hunks, and a hunk that does not fit
alone at its lines, each piece keeping the file header and a recomputed hunk
header so every anchor still names a line of the file. A fifth segment blocks
the review with the size refusal ADR-0352 made actionable.

**The record keeps them apart.** The record lists the segments as
`dataSegments`, a finding from one carries `dataSegment`, its anchor is checked
against that segment's diff, and the answer file joins the answers under their
own headers. `delivery.derived_paths` and `pre_pr_review.omit_paths` keep their
meaning; a review without the key sends the same bytes as before.

## Considered Options

- **A structured summary of the data file.** Counting added, removed and
  changed keys needs a parser per format and hides the values, which are what
  the reviewer must read.
- **Raising the bound for declared files.** ADR-0282 measured why the bound
  exists; a larger prompt fails in the provider, not in Roundfix.
- **Reusing Review Segments.** A Review Segment is a commit range; a data file
  added in one commit is still one range above the bound.
- **Reading the category from `.gitattributes`.** The candidate's own
  attributes would choose how its review treats a file, which ADR-0282 and
  ADR-0352 reject.

## Consequences

- An adopter that declared an authored catalog in `delivery.derived_paths`
  moves it to `pre_pr_review.data_paths` on its default branch; an item sees
  the key once its merge base includes it.
- A review with data files costs one prompt per Data Review Segment, up to
  four, in the session it already holds.
- ADR-0282 rejected splitting the diff by file or by bytes for Review
  Segments; this decision splits by file and by hunk only for declared data
  files, whose findings anchor to lines of the file, and leaves Review
  Segments as commit ranges. It amends
  [ADR-0282](0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md)
  and refines
  [ADR-0352](0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md).
