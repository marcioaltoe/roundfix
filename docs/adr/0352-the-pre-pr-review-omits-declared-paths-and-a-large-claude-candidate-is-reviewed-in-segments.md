---
status: accepted
created_at: 2026-11-03T12:00:00Z
updated_at: 2026-11-03T12:00:00Z
deprecated_at: null
superseded_by: null
---

# The pre-PR review omits declared paths, and a large Claude candidate is reviewed in segments

On 2026-11-03 Vortex reported a Delivery Queue item parked `review-blocked`
by the Claude review: the candidate's prompt was estimated at 541,455 tokens,
above the 500,000-token budget ADR-0227 set, although 669 KB of its
1,049,480-byte diff was one generated ORM snapshot and nothing else was large.
The only way to leave the file out was a `delivery.derived_paths` declaration
with a `regenerate` command that exits 1, which mixes conflict regeneration
with review scope. ADR-0282 splits a large candidate into Review Segments only
for `codex`, and the refusal named no file.

## Decision

**A declared path leaves the review.** Project Config gains
`pre_pr_review.omit_paths`, a list of clean repository-relative paths,
directories ending in `/` and whole-path Go `path.Match` patterns, validated
like `delivery.derived_paths` entries. The review reads it from the Project
Config at the candidate's base commit, as ADR-0282 reads derived declarations,
never from the head. A changed file it matches is left out of the diff with the
new reason `omit-path`, after every existing reason, so a file another rule
already omits keeps its reason. The prompt's omission summary names the count
only when it is not zero, so a review without the key sends the same prompt
bytes. Conflict regeneration still reads only `delivery.derived_paths`.

**Claude splits round 1 into Review Segments too.** A `claude` round-1 prompt
whose estimate exceeds the budget is reviewed as at most four consecutive
first-parent commit ranges, each a prompt as sent — its segment header, the
instructions, the segment diff, its omissions and, for the first segment, the
Spec contexts — estimated within the budget, sent in order to one review
session. A single commit whose segment exceeds the budget, a fifth segment,
and a round-2 or Correction Review prompt above the budget still block. The
`codex` segments keep the byte bound and their bytes.

**A size refusal names what to omit.** Every refusal by the byte bound, the
token budget or the segment limit keeps its text and appends the three largest
files of the diff it measured, by bytes of each file's section, with the key
that would leave them out. The Delivery Queue parks a blocked review as
`review-blocked: <first line of the record's reason>`.

## Rejected

- **Reusing `delivery.derived_paths` with an optional `regenerate`.** A derived
  path is one the queue regenerates on a conflict; a review omission is not.
- **Omitting the largest files automatically above the budget.** The reviewer
  would skip authored code nobody chose to skip.
- **Reading `linguist-generated` from `.gitattributes`.** The candidate's own
  attributes would decide what its review skips, which ADR-0282 rejects for
  Project Config.
- **Bounding Claude segments in bytes.** The Claude budget counts the whole
  prompt in tokens; a byte bound would split candidates the budget admits.
- **A fresh session per Claude segment, or a bound on the segments' sum.** One
  session keeps the Reviewer Lineage of ADR-0282; the runtime compacts earlier
  turns, as it does for `codex`, and a sum bound would admit nothing a single
  prompt does not.

## Consequences

- An adopter removes the `derived_paths` workaround and lists the file in
  `omit_paths` on its default branch; an item branched before that change
  sees the key only after its merge base includes it.
- A Claude segmented review costs up to four prompts in one session; a later
  segment that the runtime still answers with `Prompt is too long` blocks with
  ADR-0227's quoted reason.
- Every `review-blocked` park carries the record's reason after a colon; its
  Park Class is unchanged.
