---
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
prd: _prd.md
created: 2026-11-05
---

# A main merge is a review delta, and a segment is measured from its Effective Base — Technical Spec

## Executive Summary

The pre-PR review measures every diff inside the candidate between two
commits as they stand. When the item merged the default branch, two things go
wrong (`_prd.md` → Sources, points 2 and 4). `decideReviewLineage`
(`internal/cli/review_lineage.go`) starts round 1 whenever the merge base
moved, and `planReviewSegments` (`internal/cli/review_segments.go`) diffs an
item commit cut from an old default branch against the new merge base, so the
default branch's work reads as reverted and blocks the review by size.

The fix adds one measurement, `reviewMeasuredFrom` in a new
`internal/cli/review_effective_base.go`: a diff from commit `S` to a later
commit `C` starts from `S` when both share one Effective Base, from `C`'s
Effective Base when `S` is on the default branch, and otherwise from the tree
`git merge-tree --write-tree` writes for `S` carried onto `C`'s Effective Base
(ADR-0366). Segments use it for every diff they measure. The lineage survives
a moved merge base when the move is only conflict-free merges of the default
branch, and its round-2 or Correction Review delta starts from the carried
tree. When the lineage restarts, the decision names why, and `deliver status`
and `deliver retry` print that reason before the review runs.

The trade-off accepted: each measured diff costs two `git merge-base` calls
and, after a merge, one `git merge-tree`; a carry that conflicts keeps today's
measurement instead of resolving anything.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or schema version
  is minted; the Pre-PR Review Record gains two optional fields. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0368 (a large authored data file reviewed apart, a sibling Spec in the same queue) touches the same review code but carries no obligation for this Spec; ADR-0366 (new) records the decision;
  it refines ADR-0169 ("The candidate's own merge of the default branch still
  moves the merge base forward") and amends ADR-0197 ("A rebase moves the
  merge base, so it starts a new lineage at round 1"), ADR-0282 and ADR-0352
  (how a Review Segment is measured) and ADR-0286 (where the Correction
  Review's delta starts); task_04 writes the four notices.
  Accounted for and unchanged, each one citation edge from a listed ADR on a part this Spec leaves alone: ADR-0092, ADR-0120, ADR-0165, ADR-0192, ADR-0196, ADR-0227, ADR-0233, ADR-0247, ADR-0248, ADR-0250, ADR-0251, ADR-0254, ADR-0256, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0269, ADR-0272, ADR-0277, ADR-0278, ADR-0279, ADR-0281, ADR-0290, ADR-0295, ADR-0299, ADR-0306, ADR-0308, ADR-0317, ADR-0345, ADR-0348, ADR-0349. They keep the queue's own conflict merge without a repeated review (ADR-0192), finding validation (ADR-0196), the Claude budget (ADR-0227), and the history, archive, record, module and binary decisions the design never touches. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ batch 2 answers of 2026-11-05 quoted in the authorization record; the
  Governed Paths are the Roundfix skill files task_04 changes; no Makefile,
  workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-05 at `ad77a797`; `_prd.md` → Validation records the runs.

1. `runReviewCommand` (`internal/cli/review.go`) resolves the base tip,
   then the merge base `B` of the tip and `HEAD`, reads the checkout's prior
   record and calls `decideReviewLineage`. That function returns round 1 when
   the prior record is absent, belongs to another checkout, has another
   provider or another `BaseCommit`, or when the head does not descend from
   the prior head (`git merge-base --is-ancestor`, exit 1). Otherwise it
   reuses, repeats a blocked round, reaches the ceiling, or plans round 2 with
   `PreviousHead` and the prior findings.
2. `runConfiguredReviewSession` builds the round-2 and Correction Review delta
   with `reviewScopedDiff(ctx, gitRoot, plan.Lineage.PreviousHead, headCommit,
   specRoots, gitRunner, baseCommit)`; `buildRoundTwoPrompt` and
   `buildCorrectionReviewPrompt` (`review_lineage.go`) frame it.
3. `planReviewSegments` lists `git rev-list --first-parent --reverse
   <base>..<head>`, and for each commit diffs `previous..commit` (its own diff,
   which must fit alone) and `segmentBase..commit` (the growing segment).
   With `previous` starting at `B`, a first item commit cut from an older
   default branch diffs against `B` and shows the default branch's later work
   as removed; the merge commit diffs against the item's earlier commit and
   shows it as added. `reviewSegment` holds `base`, `head`, `diffBytes` and,
   for Claude, `estimatedPromptTokens`; `reviewSegmentHeader` writes the one
   header line.
4. `reviewScopedDiff` accepts any tree-ish as its base: it runs `git ls-tree`,
   `git show <base>:<path>` and `git diff <base> <head>`, all of which accept a
   tree id, and reads Project Config declarations from the separate
   declaration base.
5. `runDeliverStatus` (`internal/cli/deliver.go`) prints one line per item,
   the `Warning:` lines, the `Owner:` line, the `Park:` lines, limits, usage
   and the pending question. `runDeliverRetry` prints `Retried <slug>:
   <blocker> -> <stage>` through `printDeliverRetryResult`, then hands off to
   a live owner or starts one. Neither reads a review record.
6. `git merge-tree --write-tree` is already used by
   `internal/cli/deliver_workflow.go`.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/cli/review_effective_base.go` (new): `reviewEffectiveBase` and
  `reviewMeasuredFrom`.
- `internal/cli/review_segments.go`: every measured diff starts from
  `reviewMeasuredFrom`; `reviewSegment.MeasuredFrom`; the header line.
- `internal/cli/review_lineage.go`: the lineage across a moved merge base,
  `reviewLineage.DeltaBase`, `reviewLineagePlan.Restart`, the prompt sentence.
- `internal/cli/review.go`: `runReviewCommand` computes `DeltaBase` for a
  round-2 plan; `runConfiguredReviewSession` starts the delta from it.
- `internal/cli/review_lineage_notice.go` (new): `reviewLineageNotice`.
- `internal/cli/deliver.go`: the two calls of the notice.
- Tests: `internal/cli/review_effective_base_test.go`,
  `internal/cli/review_main_merge_lineage_test.go` and
  `internal/cli/deliver_review_lineage_notice_test.go` (new).
- Documentation: the `review` and `deliver` guides, the Roundfix `review`
  and `deliver` references, `CONTEXT.md` and four ADR amendment notices.

## Implementation Design

### Interfaces

```go
// internal/cli/review_effective_base.go

// reviewEffectiveBase returns git merge-base <base> <commit>.
func reviewEffectiveBase(ctx context.Context, git preflight.GitRunner, root, base, commit string) (string, error)

// reviewMeasuredFrom returns what a diff from `from` to `to` starts from so
// that it holds only the candidate's own change (API Contract 1), and whether
// `from` carried onto `to`'s Effective Base without conflict.
func reviewMeasuredFrom(ctx context.Context, git preflight.GitRunner, root, base, from, to string) (string, bool, error)

// internal/cli/review_lineage_notice.go

// reviewLineageNotice returns the Review: line of API Contract 5, or "".
func reviewLineageNotice(ctx context.Context, git preflight.GitRunner, loaded roundconfig.Loaded, specSlug, worktree string) string
```

### Data Models

- `reviewSegment` gains `MeasuredFrom string` with JSON name `measuredFrom`,
  `omitempty`: the commit or tree the segment's diff starts from when it is
  not the segment's `base`.
- `reviewLineage` gains `DeltaBase string` with JSON name `deltaBase`,
  `omitempty`: the tree the round-2 or Correction Review delta starts from
  when it is not `previousHead`.
- `reviewLineagePlan` gains `Restart string`, not persisted.

A record without the new fields reads as before.

### API Contracts

1. API Contract: the measurement. `reviewEffectiveBase(base, commit)` is the
   trimmed output of `git merge-base <base> <commit>`; `base` is the review's
   merge base `B`, so the answer equals the commit's merge base with the
   current default branch. `reviewMeasuredFrom(base, from, to)` computes `E1
   = reviewEffectiveBase(base, from)` and `E2 = reviewEffectiveBase(base,
   to)` and returns:
   - `from, true` when `E1 == E2`;
   - `E2, true` when `E1 == from` (`from` is on the default branch);
   - otherwise the first line of `git merge-tree --write-tree <from> <E2>`
     and `true` on exit 0, or `from, false` on exit 1 (a conflict).
   Any other Git failure is returned as an error naming both commits.
2. API Contract: segments. In `planReviewSegments`, the own diff of each
   commit and the diff of each growing segment start from
   `reviewMeasuredFrom(base, <start>, <commit>)` instead of `<start>`. A
   segment whose diff starts from something other than its `base` records it
   as `MeasuredFrom`, and `reviewSegmentHeader` then appends this second line
   to the segment's header, in the prompt sent and in every fit estimate:

   ```text
   This segment's diff starts from <measuredFrom> instead of <segment base>, so it leaves out what the default branch changed.
   ```

   `codex` and `claude` share the plan. A conflicting carry measures from
   `<start>` as at `ad77a797`. The size refusals, the bound, the budget, the
   four-segment limit and the first header line keep their bytes.
3. API Contract: the lineage across a moved merge base. `decideReviewLineage`
   keeps the same checkout rule and then decides in this order, returning
   round 1 with `Restart` set by API Contract 4 at the first failure:
   1. the provider is unchanged;
   2. when the head moved, it descends from the prior head;
   3. when the merge base moved: the head moved; the prior merge base is an
      ancestor of the new one; every merge listed by `git rev-list
      --first-parent --merges --parents --reverse <prior head>..<head>` has
      every parent after the first contained in the new merge base; and
      `reviewMeasuredFrom(<new base>, <prior head>, <head>)` carries.
   When all hold, the plan continues exactly as when the merge base did not
   move. For every round-2 plan with a `PreviousHead`, `runReviewCommand`
   then sets `DeltaBase` to `reviewMeasuredFrom(B, PreviousHead, HEAD)` when
   that differs from `PreviousHead`, else to `""`, before the reuse check, and
   `runConfiguredReviewSession` starts the round-2 or Correction Review delta
   from `DeltaBase` when it is set, and its `measuredBase`, which the delta's
   size refusal measures from, is that same object. `buildRoundTwoPrompt` writes, after the
   line ending `The diff below is the change since round 1 reviewed <previous
   head>.` and before `Round 1 findings:`, and `buildCorrectionReviewPrompt`
   writes after its own first paragraph, this line when `DeltaBase` is set:

   ```text
   The candidate merged the default branch since then; the diff starts from tree <deltaBase>, the reviewed head with the default branch's changes applied, so it holds only the candidate's own change.
   ```

   Round 2 stays unsegmented; its bound, anchors, session reuse, ceiling and
   Correction Review conditions are unchanged.
4. API Contract: the restart reasons, word for word, with every commit as Git
   prints it in full:
   1. `the review provider changed from <prior provider> to <provider>`
   2. `the head does not descend from the reviewed head <prior head>`
   3. `the merge base moved from <prior base> to <base> at the reviewed head`
   4. `the merge base moved from <prior base> to <base>, which does not descend from it`
   5. `the merge base moved from <prior base> to <base> through merge <merge> of <parent>, which the default branch does not contain`
   6. `the merge base moved from <prior base> to <base> and the reviewed head <prior head> conflicts with it`
   Reason 5 names the oldest such merge and its first parent the new merge
   base lacks. A prior record of another checkout, or no prior record, sets
   no reason.
5. API Contract: the notice. `reviewLineageNotice` returns `""` unless the
   configured provider is `codex` or `claude`, the worktree is a Git work
   tree, its prior record exists at
   `reviewCheckoutDir(<artifact directory>, <worktree root>)` with another
   head, and `decideReviewLineage` returns a `Restart`. It resolves the base
   as `roundfix review` without `--base` does and the artifact directory with
   `roundconfig.ResolveArtifactDirectory` from the loaded Project Config,
   writes nothing, and returns `""` on any error. Otherwise it returns:

   ```text
   Review: <slug> head <head> opens a new Reviewer Lineage at round 1, a full review of the candidate: <reason>
   ```

   `deliver status` prints the line for each `parked` or `reviewing` item,
   after the `Warning:` lines and before the `Owner:` line. `deliver retry`
   prints it right after `Retried <slug>: <blocker> -> <stage>` when the
   retried stage is `reviewing`, reading the item's worktree from the queue
   before the store closes. Exit codes and every other line are unchanged.
6. API Contract: the rule in writing (the Effective Base paragraph of the
   `review` guide in task_01, everything else in task_04), word for word, line
   wrapping free:
   - In both copies of the Roundfix `review` reference and in
     `docs/user-guide/commands/review.md`, after the paragraph ending
     "all segments are one round.", a paragraph:

     ```markdown
     Each diff the segment plan measures starts from the Effective Base of the commit it ends at, the merge base of that commit with the review's merge base. A commit cut from an older default branch is measured from that older commit, and a merge of the default branch from the earlier commit's change carried onto the merged base with `git merge-tree --write-tree`, so no segment shows the default branch's own work. A segment whose diff starts elsewhere than its base records `measuredFrom`, the commit or tree it starts from, and its prompt adds a second header line: `This segment's diff starts from <measuredFrom> instead of <segment base>, so it leaves out what the default branch changed.` A carry that conflicts measures the segment from its base.
     ```

   - In the same three files, the sentence "Reviews in one checkout form a
     **Reviewer Lineage** when the provider and merge base stay the same and
     each reviewed head descends from the last." becomes "Reviews in one
     checkout form a **Reviewer Lineage** when the provider stays the same,
     each reviewed head descends from the last, and the merge base stays the
     same or moves only by a conflict-free merge of the default branch.", and
     the sentence "A rebase, changed provider or merge base, or a head that
     does not descend starts a new lineage at round 1." becomes:

     ```markdown
     A merge of the default branch keeps the lineage when every merge the head gained since the reviewed head brings in only commits the new merge base contains and the reviewed head merges onto the new merge base without conflict. The round-2 or Correction Review delta then starts from the reviewed head carried onto the new merge base, recorded as `lineage.deltaBase`, and the prompt adds: `The candidate merged the default branch since then; the diff starts from tree <deltaBase>, the reviewed head with the default branch's changes applied, so it holds only the candidate's own change.` A rebase, a changed provider, a head that does not descend, or any other move of the merge base starts a new lineage at round 1, and `roundfix deliver status` and `roundfix deliver retry` name the reason before the review runs.
     ```

   - In `docs/user-guide/commands/deliver.md` and both copies of the
     Roundfix `deliver` reference, after the paragraph ending "for `deliver
     resume` to recover.":

     ```markdown
     When an item's worktree holds a Pre-PR Review Record of an earlier head and the next review of its current head will open a new Reviewer Lineage at round 1, `deliver status` prints, for a parked or reviewing item, after any `Warning:` lines and before the `Owner:` line: `Review: <slug> head <head> opens a new Reviewer Lineage at round 1, a full review of the candidate: <reason>`. The reason is one of the six the `review` guide lists; a conflict-free merge of the default branch keeps the lineage and prints nothing.
     ```

     and, in the guide after the paragraph ending "An item-level refusal
     leaves the stored item unchanged." and in both copies of the reference
     right after the paragraph above:

     ```markdown
     When the retried item returns to `reviewing` and the next review of its current head will open a new Reviewer Lineage at round 1, `deliver retry` prints the same `Review:` line right after `Retried <slug>: <blocker> -> reviewing`, before the owner runs the review.
     ```

   - The `review` guide and both copies of the reference also list the six
     reasons of API Contract 4 after the lineage paragraph, one bullet each,
     introduced by "A restart names one of these reasons:".
   - `CONTEXT.md`, after the **Review Segment** entry:

     ```markdown
     **Effective Base**:
     The merge base of one candidate commit with the Delivery Base: the commit of the default branch it was cut from or last merged. The pre-PR review measures each Review Segment and each round-2 delta from it, carrying the earlier commit's own change onto it, so neither shows the default branch's work as the candidate's (ADR-0366).
     _Avoid_: Segment base, merge base of the segment
     ```

     In `CONTEXT.md` the closing `(ADR-0366)` is the Markdown link
     `docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`,
     as the neighbouring entries write theirs.
     The **Review Segment** entry's "each within the applicable bound" becomes
     "each measured from its Effective Base and within the applicable bound",
     and its ADR list gains ADR-0366. The **Reviewer Lineage** entry gains,
     before "It may include one Correction Review", the sentence "A
     conflict-free merge of the default branch keeps it, and the delta then
     starts from the round-1 head carried onto the new merge base.", and its
     ADR list gains ADR-0366.
   - The four amendment notices, each a new last paragraph, with `ADR-0366`
     the Markdown link to
     `0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`
     and `updated_at: 2026-11-05T12:00:00Z`:
     - ADR-0197: "**Amended by ADR-0366 (2026-11-05).** A conflict-free merge
       of the default branch into a descending head keeps the Reviewer
       Lineage although it moves the merge base; round 2 reviews the delta
       from the reviewed head carried onto the new merge base. A rebase, or a
       merge that brings in another branch or conflicts, still starts round
       1."
     - ADR-0282: "**Amended by ADR-0366 (2026-11-05).** Each Review Segment,
       and each commit within it, is measured from its Effective Base, so a
       commit cut from an older default branch or a conflict-free merge of
       the default branch no longer shows the default branch's work; a commit
       still above the bound measured that way still blocks."
     - ADR-0286: "**Amended by ADR-0366 (2026-11-05).** When the head merged
       the default branch without conflict, the Correction Review's delta
       starts from the round-2 head carried onto the new merge base; its
       conditions and its single run are unchanged."
     - ADR-0352: "**Amended by ADR-0366 (2026-11-05).** A Claude Review
       Segment's prompt is estimated on its diff measured from its Effective
       Base."

### Surface Transcripts

Both run in a temporary repository that is the queue's checkout and the
item's worktree, with a codex review record of its round-1 findings at
`<reviewed head>` and the head rewritten with `git commit --amend` to
`<head>`; the item is parked `review-findings`.

1. Surface Transcript: `deliver status` names the restart.

   ```transcript
   $ roundfix deliver status
   stdout:
   0356-example	parked	review-findings	<worktree>
   Review: 0356-example head <head> opens a new Reviewer Lineage at round 1, a full review of the candidate: the head does not descend from the reviewed head <reviewed head>
   Park: 0356-example finding: resolve the blocker, then run roundfix deliver retry 0356-example
   ...
   stderr:
   exit: 0
   ```

2. Surface Transcript: `deliver retry` names it before the owner runs the
   review (the test's owner seam prints the last line).

   ```transcript
   $ roundfix deliver retry 0356-example
   stdout:
   Retried 0356-example: review-findings -> reviewing
   Review: 0356-example head <head> opens a new Reviewer Lineage at round 1, a full review of the candidate: the head does not descend from the reviewed head <reviewed head>
   Delivery Owner: test-owner
   stderr:
   exit: 0
   ```

### Invariants

```text
1. No segment prompt and no round-2 or Correction Review prompt of a candidate that merged the default branch without conflict holds a line the default branch added after the candidate was cut.
2. A candidate with no merge of the default branch and a single Effective Base is measured, segmented and prompted byte for byte as at ad77a797.
3. Every row of TestReviewLineageDecidesEachRound keeps its round, reuse and ceiling, the rebase row among them.
4. Every size refusal, the review bound, the Claude budget and the four-segment limit keep their text and values.
5. deliver status and deliver retry print no Review: line for an item without a prior record, with the same head, or whose lineage continues, and keep every other line and exit code.
```

### Declared breaks

- A conflict-free merge of the default branch no longer starts round 1: the
  next review is round 2, the Correction Review, or the ceiling, and a round-2
  record may carry `lineage.deltaBase`.
- A segmented record may carry `segments[].measuredFrom`, its prompt the
  second header line, and its `diffBytes` drop by the default branch's work.
- `deliver status` and `deliver retry` may print one `Review:` line per item.
- No existing test changes: the prototype ran `go test -count=1
  ./internal/cli` unchanged. No test is renamed or removed, so the coverage
  record is not re-recorded.

## Coverage Map

- Goal 1, Core Features 1 and 3, Success Metric 2 → API Contracts 1 and 3,
  `TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange`,
  `TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch`,
  `TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix`,
  Invariant 1.
- Goal 2, Core Feature 3, Success Metric 3 → API Contract 4,
  `TestReviewLineageRestartNamesItsReason`,
  `TestReviewLineageDecidesEachRound`, Invariant 3.
- Goal 3, Core Feature 4, Success Metric 4 → API Contract 5, Surface
  Transcripts 1 and 2, `TestDeliverStatusNamesAReviewThatOpensANewLineage`,
  `TestDeliverRetryNamesAReviewThatOpensANewLineage`, Invariant 5.
- Goal 4, Core Features 1 and 2, Success Metric 1 → API Contracts 1 and 2,
  `TestReviewSegmentsMeasureEachCommitFromItsEffectiveBase`, Invariants 1, 2
  and 4.
- Goal 5, Core Feature 5, Success Metric 5 → API Contract 6.

Deviations from the PRD: None.

## Absence Evidence

1. `reviewEffectiveBase`, `reviewMeasuredFrom` and an Effective Base
   measurement in the review — command: `grep -rn 'EffectiveBase\|MeasuredFrom\|measuredFrom\|merge-tree' internal/cli/review*.go`;
   result: no match; `merge-tree` appears only in `internal/cli/deliver_workflow.go`.
2. `DeltaBase`, `Restart` and a lineage across a moved merge base — command:
   `grep -n 'DeltaBase\|deltaBase\|Restart' internal/cli/review_lineage.go internal/cli/review.go`;
   result: no match; `decideReviewLineage` returns round 1 on any
   `BaseCommit` difference.
3. `reviewLineageNotice` and a `Review:` status line — command:
   `grep -rn 'reviewLineageNotice\|"Review: ' internal/cli/*.go`; result: no
   match.
4. The three new test files — command: `ls internal/cli/review_effective_base_test.go internal/cli/review_main_merge_lineage_test.go internal/cli/deliver_review_lineage_notice_test.go`;
   result: no such files.

## ADR Alternatives

- ADR-0368: it rejects a structured summary, a larger bound, reusing Review Segments and `.gitattributes`; this design takes none of them.
- ADR-0366: the decision this Spec records; it rejects a tree-equality test
  of the merge, `--remerge-diff`, a larger bound, history rewriting and a
  delta across a merge of another branch, and the design takes none of them.
- ADR-0169: unchanged; the merge base still moves with the candidate's own
  merge, and the full candidate diff still starts from it. The design takes
  none of its rejected alternatives.
- ADR-0197: amended; its rule that a moved merge base always starts round 1
  narrows to moves other than a conflict-free merge of the default branch.
  It rejected no option this design takes; its two-round ceiling stays.
- ADR-0282: amended; segments stay consecutive first-parent commit ranges,
  as it chose over splitting by file or by bytes; only where each diff starts
  changes, and it still rejects raising the bound.
- ADR-0286: amended; its Correction Review keeps its conditions and single
  run, and it still rejects a third round and a dismissal-opened review.
- ADR-0352: amended; Claude segments keep one session and the token budget it
  chose over a byte bound or a fresh session per segment.
- ADR-0192: unchanged; the queue's own merge on a conflicting Pull Request still repeats no review, and its rejected rebase stays rejected.
- ADR-0196: unchanged; findings are still validated before they park, and none of its rejected options is taken.
- ADR-0227: unchanged; the Claude budget and its quoted refusal keep their values, and a larger bound stays rejected.
- ADR-0233: unchanged; a review-only correction still returns to review, and none of its rejected options is taken.
- ADR-0272: unchanged; the record still survives a fix and the declared-path bound of a retry still holds; none of its rejected options is taken.
- ADR-0278: unchanged; the design touches no History Root rule.
- ADR-0281: unchanged; the design rebuilds no test executable.
- ADR-0290: unchanged; the design reclaims no Item Binary or Run Worktree.
- ADR-0299: unchanged; the design writes no history citation.
- ADR-0092: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0120: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0247: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0248: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0250: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0254: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0256: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0258: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0260: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0261: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0277: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0279: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0295: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0308: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0317: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0345: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0165: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0251: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0262: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0269: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0348: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0349: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.
- ADR-0306: unchanged; accounted for in Project Constraints, and the design takes none of its rejected alternatives.

## Integration Points

None external. The `git` binary is a Base Test Tool; `git merge-tree
--write-tree` needs Git 2.38, which the Linux Lane image and the existing
`deliver_workflow.go` use already.

## Testing Approach

All tests are hermetic: `t.TempDir()` repositories through
`newReviewCommandFixture`, `commitScopeFile`, `commitLineageFile` and
`gittest.Run`, the fake reviewers `segmentReviewRunner` and
`newLineageRunner`, the real queue store under the fixture's home, and
`runCLI`. Every new top-level test calls `t.Parallel()` first, and each new
file opens with the `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT`
comment lines. Reviews run as `review --base main` so the merge base follows
the default branch.

- `internal/cli/review_effective_base_test.go` (task_01):
  - `TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange`, after the
    default branch gains `upstream.txt` and is merged into the item: `same
    Effective Base` (the merge to a later commit returns the merge), `start on
    the base` (the new tip to the fixture's own commit returns the old base),
    `merge of the default branch` (the fixture's commit to the merge returns a
    tree with no diff to the merge), and `conflicting carry` (a commit that
    wrote `upstream.txt` differently returns itself and `false`).
  - `TestReviewSegmentsMeasureEachCommitFromItsEffectiveBase`, subtests
    `codex` and `claude` with `unit` = `reviewDiffBound` or
    `2*claudeReviewPromptBudget` bytes: the item commits `item.txt` and
    `item2.txt` (two thirds of `unit` each) after the fixture's commit, the
    default branch gains `mainwork.txt` (`unit + unit/5`) and is merged, then
    `part1.txt` (two thirds of `unit`). Exit `0`, outcome `reviewed`, three
    segments: the first ends at `item.txt`'s commit with `measuredFrom` the
    old base, the second ends at the merge with a tree `measuredFrom`, the
    third has none; each prompt holds the second header line exactly when its
    segment has `measuredFrom`, and none holds 1,000 `m` in a row.
- `internal/cli/review_main_merge_lineage_test.go` (task_02):
  - `TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch`:
    round 1 findings, a fix disposed `--fixed-by`, the default branch gains
    `upstream.txt` (`main moved on`) and is merged; the next review exits `0`,
    round 2, `previousHead` the round-1 head, `deltaBase` a tree, base the new
    tip; its prompt holds the round-2 delta markers, the fix and the sentence,
    and neither `main moved on` nor `BEGIN CANDIDATE DIFF`.
  - `TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix`: round
    1 clean, round 2 findings on `delta.txt`, a fix disposed `--fixed-by`, the
    default branch merged; the review is the Correction Review with
    `deltaBase` set, its prompt holds the correction delta markers, the fix
    and the sentence, not `main moved on`.
  - `TestReviewLineageRestartNamesItsReason`: `decideReviewLineage` rows
    `clean merge` (no reason, round 2) and one row per reason of API
    Contract 4: `provider`, `rewritten head`, `same head`, `base not
    descending`, `merge of another branch` and `conflict` (a default branch
    commit rewriting `review.txt`'s second line, merged with `-X ours`).
- `internal/cli/deliver_review_lineage_notice_test.go` (task_03):
  - `TestDeliverStatusNamesAReviewThatOpensANewLineage`, subtests `rewritten
    head` (Surface Transcript 1's `Review:` line) and `conflict-free merge of
    the default branch` (no `Review:` line).
  - `TestDeliverRetryNamesAReviewThatOpensANewLineage`: Surface Transcript 2
    whole, through the `newDeliveryEngine` and `startDeliveryOwner` seams of
    `deliver_retry_test.go`.
- Existing suites: `./internal/cli` whole, `./skills`, `-tags docscontract
  ./internal/docscontract`.
- Sabotages, each run by the author on the prototype, with the case whose
  input makes it observable:
  - segments measured from `<start>` again:
    `TestReviewSegmentsMeasureEachCommitFromItsEffectiveBase` fails on both
    subtests with exit `2`;
  - `reviewMeasuredFrom` returning `from` when `from` is on the base: row
    `start on the base` fails;
  - the merge-parent check skipped: row `merge of another branch` of
    `TestReviewLineageRestartNamesItsReason` returns no reason;
  - the conflict check skipped: row `conflict` fails;
  - the delta started from `PreviousHead` again: both lineage tests fail on
    `main moved on`;
  - the notice never printed: both notice tests fail.

## Build Order

1. `reviewMeasuredFrom` and segments measured from it, with
   `review_effective_base_test.go`, and the Effective Base paragraph of the
   `review` guide with the Behavior Surface Record re-recorded last
   (task_01).
2. The lineage across a moved merge base, `DeltaBase`, `Restart`, the prompt
   sentence, with `review_main_merge_lineage_test.go` (task_02) (depends on:
   1).
3. `reviewLineageNotice` and its two calls in `deliver.go`, with
   `deliver_review_lineage_notice_test.go` (task_03) (depends on: 2).
4. The guides, the Roundfix `review` and `deliver` references, `CONTEXT.md`, the four ADR
   amendment notices, the skill version and the Behavior Surface Record
   recorded last (task_04) (depends on: 3).
5. QA gate (task_05) (depends on: 4).

task_01 writes the guide paragraph so the CLI files it, task_02 and task_03
change have a guide in their dependency closure; task_02 calls
`reviewMeasuredFrom` from task_01; task_03 reads the `Restart` task_02
sets. task_04 documents all three. No Verification depends on another Spec.
Spec 0355 may land first and change `runDeliverStatus` and `runDeliverRetry`;
the two calls of task_03 are placed by the lines they follow, not by line
number, and Surface Transcript 1 ends in `...` so a changed `Limits:` line
still matches.

## Risks & Considerations

- A merge that the default branch's work and the item's fix resolve
  differently from the carried tree shows the difference in the delta; that
  is the resolution, which a reviewer should read.
- A carry that conflicts keeps the old measurement, so a candidate whose
  early commit conflicts with the default branch can still block by size; the
  operator's linear-commit workaround remains for that case.
- The notice reads the Project Config of the checkout running `deliver
  status`; an item whose worktree configures another artifact directory gets
  no notice, never a wrong one.
- Two `git merge-base` calls per measured commit add milliseconds per
  segment candidate; segmentation runs only above the bound.

## Glossary

- adds: **Effective Base**
- changes: **Reviewer Lineage**
- changes: **Review Segment**

## Decisions

- One measurement function for segments and deltas; no new package.
- The carry uses `git merge-tree --write-tree` without `--merge-base` (that
  option needs Git 2.40); Git's own merge base of the two commits is the
  earlier Effective Base in every case the plan reaches.
- The merge-parent rule checks every merge since the reviewed head, so a
  merge of another branch restarts the lineage as the maintainer's answer
  requires.
- The notice is advisory: any read failure prints nothing.
- ADR-0366 records the decision, refines ADR-0169 and amends ADR-0197,
  ADR-0282, ADR-0286 and ADR-0352.
- Secondbrain and Exa as `_prd.md` → Validation records.
- Four implementation Tasks in series and a QA gate.
