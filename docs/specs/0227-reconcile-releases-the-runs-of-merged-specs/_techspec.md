---
spec: 0227-reconcile-releases-the-runs-of-merged-specs
prd: _prd.md
created: 2026-10-04
---

# Reconcile releases the Runs of merged Specs — Technical Spec

## Executive Summary

Three places decide the defects the PRD measures. `inspectRunAtMergedHead` in
`internal/worktree/merged_head.go` compares the files of every non-Task,
non-QA commit with the merged head, so a Run of a merged Spec turns
`unintegrated` as soon as a later Spec edits one of those files;
`resolveMergedReleaseEvidence` in `internal/cli/deliver_workflow.go` accepts a
candidate only by ancestry or an equal tree, which a squash merge onto a moved
default branch never has; and `roundfix reconcile` never lists
`roundfix/deliver-<slug>-<16 hex>` branches. This Spec adds one proof, merge
evidence, used as a fallback after today's proof fails; a tree comparison
against `git merge-tree --write-tree` for the cleanup; and an item-branch
inspection in reconcile that reuses the merge evidence. The trade-off accepted
is that a committed change in a merged Spec's Run or item branch that never
reached the delivery is released with it: the archive on the default branch,
not a file comparison, decides that the Spec's work is done (ADR-0232). Task
commits still need their Task completed, and dirty work keeps ADR-0212's rule.

## Project Constraints

- Identifier strategy: not applicable — no identifier changes; Runs, item
  branches and the `roundfix-reconcile/v1` schema name keep their names, and
  the report only gains a field. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — every proof reads local Git and
  the Run Database; no credential, forge read or network call is added, and
  no test reaches GitHub. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0232 (this Spec) decides merge
  evidence, the cleanup proof and item branches, ADR-0232: "Any other
  difference, or a Git without that mode, still refuses with today's text".
  ADR-0161 releases a merged Spec's Runs on the merged head and stays the
  first proof, ADR-0161: "The Delivery Queue merge record is the first source
  of that head". ADR-0212 keeps its Spec-directory and dirty-path rules,
  ADR-0212: "Any other path, committed or not, keeps the Run preserved",
  which this Spec narrows to Specs without merge evidence. ADR-0053 keeps
  reconciliation read-only unless applied, ADR-0053: "Ambiguous, dirty, and
  unintegrated work stays intact", and ADR-0115 discards a superseded Run
  Branch through Roundfix with its evidence first, ADR-0115: "Reconciliation reports; disposal changes the
  repository". ADR-0187 and ADR-0189 govern the Roundfix Skill edit and its
  version. ADR-0184 has this TechSpec state its changed command surfaces,
  ADR-0184: "A TechSpec now declares numbered Surface Transcripts". The gate
  is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155 and ADR-0156,
  and ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check this Spec's
  consistency. ADR-0178, ADR-0179 and ADR-0166 decide each Task commit's
  grant and record undeclared paths. ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry, ADR-0167 cites ADR-0080 but decides the pre-PR row of a
  qualifying partial, ADR-0170 cites ADR-0053 but decides Task Carry-Forward,
  which this Spec leaves unchanged, ADR-0182 cites ADR-0117 but settles a Task on
  the facts its gate will check, and ADR-0192 cites ADR-0178 but decides
  derived-path conflicts. ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097
  but decide what a QA row records, when it is observed again and its
  evidence snapshot, and ADR-0229 cites ADR-0167 but decides the Delivery
  Retry of an operator archive; this Spec changes none of them, so none
  applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and
  its `SKILL.md` mirror are Governed Paths under the maintainer's skill
  authorization and this Spec's grant of 2026-10-04. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0227-reconcile-releases-the-runs-of-merged-specs/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `skills/roundfix/SKILL.md`.

## System Architecture

`InspectTerminalRunMerged` (`internal/worktree/worktree.go`) classifies a
terminal Run. On an ancestry miss, an absent target branch or a dirty Run
Worktree it asks `chooseMergedHead` for a source, either the Delivery Queue
merge record's head or the default branch holding the archived QA Report, and
proves the Run there with `inspectRunAtMergedHead`. Both `roundfix reconcile`
(`internal/cli/reconcile.go`) and the delivery owner's
`releaseMergedSpecRuns` (`internal/cli/deliver_release.go`) call it, and the
Run Branch set classification of `runBranch` candidates reaches
`inspectRunAtMergedHead` too. Apply re-inspects through
`revalidateTerminalRunApply` and refuses stale evidence.

This Spec adds, without new packages:

- **Merge evidence**, in `merged_head.go`, a proof any caller of
  `inspectRunAtMergedHead` and `inspectTerminalRunMerged` falls back to when
  today's proof leaves the Run `unintegrated`.
- **A moved-main cleanup proof**, a third check in
  `resolveMergedReleaseEvidence`.
- **Item branch reconciliation**, a new file
  `internal/worktree/item_branch_reconcile.go` beside the existing item
  worktree code, called from a new `inspectReconcileItemBranches` step in
  `reconcile.go`.

## Implementation Design

### Interfaces

```go
// internal/worktree/merged_head.go
type deliveryEvidence struct {
	defaultBranch  string // e.g. "main"
	defaultHead    string
	deliveryCommit string // commit that added <archive root>/<slug>/_prd.md
}
func provenDeliveryEvidence(ctx context.Context, runner gitRunner, gitRoot, slug, head string) (deliveryEvidence, bool)
// supersededByDelivery proves every commit of head that the default head lacks.
func supersededByDelivery(ctx context.Context, runner gitRunner, gitRoot, slug, head string, evidence deliveryEvidence) (reason string, refusal string, proven bool)

// internal/worktree/item_branch_reconcile.go
type ItemBranchReconciliation struct {
	Branch, SpecSlug, Head, Worktree string
	Releasable                       bool
	Proof, RefusalReason             string
	evidence                         *itemBranchEvidence
}
func InspectItemBranches(ctx context.Context, gitRoot, location string, live []string) ([]ItemBranchReconciliation, error)
func ApplyItemBranch(ctx context.Context, result ItemBranchReconciliation) error
```

### Merge evidence

1. `provenDeliveryEvidence` resolves the default branch and its head with
   `resolveDefaultBranchHead`, requires `<archive root>/<slug>/_prd.md` at that
   head (`specArchivedAtMergedHead`), finds the delivery commit with
   `git log -1 --diff-filter=A --format=%H <head> -- <archive root>/<slug>/_prd.md`,
   and fails when that commit is reachable from the branch head being proven
   (`git merge-base --is-ancestor`), because such a branch started after the
   delivery.
2. `supersededByDelivery` lists `git rev-list <branch head> ^<default head>`.
   A commit with a `Roundfix-Task` trailer must carry `Roundfix-Spec: <slug>`
   and its Task must be completed in the archived Spec at the default head
   (`taskCompletedAtMergedHead`); otherwise it refuses with today's
   `unrepresentedMergedHeadCommit` cause. Every other commit, QA Report and
   Spec-directory commits included, is counted as superseded.
3. A Run proven this way is `superseded` with the reason
   `Run work is superseded by the delivery of Spec <slug>: delivery commit <12 hex> archived it on default branch "<branch>"; <n> Task commit(s) completed, <m> other commit(s) superseded`,
   bounded by `boundedReconciliationReason`. Its evidence records the default
   head and branch as `proofHead` and `proofRef`, so apply re-proves it and a
   moved default branch refuses as stale.
4. The fallback runs only where today's proof would return `unintegrated`:
   at the end of `inspectRunAtMergedHead`, and in `inspectTerminalRunMerged`
   when no merged-head source qualifies on an ancestry miss or an absent
   target branch. A `safe` or `superseded` result from today's proof keeps its
   reason, so every existing evidence string is unchanged.
5. A dirty Run Worktree stays `dirty` unless the proof is `safe` or
   `superseded` and `dirtyPathsInArchivedSpec` accepts every dirty path at the
   head that carried the archived Spec (ADR-0212's rule, unchanged). The
   dirty branch of `inspectTerminalRunMerged` tries merge evidence when no
   merged-head source holds the archived Spec.

### The cleanup proof for a moved default branch

1. `resolveMergedReleaseEvidence` keeps its ancestry and equal-tree checks.
2. When both miss and the merge commit has a first parent, it runs
   `git merge-tree --write-tree <merge>^1 <candidate>`; a zero exit whose one
   output line equals `<merge>^{tree}` proves the candidate.
3. A conflict, a non-zero exit, a Git without `--write-tree`, a merge commit
   without a parent, or a different tree refuses with today's
   `candidate head "<head>" is not represented by merge commit "<merge>"`.

### Item branch reconciliation

1. Only a full scan (`roundfix reconcile` without a Run ID) inspects item
   branches: every local `refs/heads/roundfix/deliver-*` whose name matches
   `^roundfix/deliver-(<slug>)-[0-9a-f]{16}$`.
2. `live` holds each branch a Delivery Queue item of this repository records
   while its stage is not `merged`, read through the same reader as
   `loadDeliveryMergedHeads`. A live branch is preserved with
   `item branch belongs to live Delivery Queue item "<slug>"`.
3. Other branches are releasable only with merge evidence for their slug and
   a proven `supersededByDelivery`; otherwise preserved with the proof's
   refusal or `Spec "<slug>" has no merge evidence on the default branch`.
4. A worktree registered on the branch must sit at
   `ItemRefFor(gitRoot, location, branch).Path` and report an empty
   `git status --porcelain=v1 -z --untracked-files=all`; otherwise the branch
   is preserved with the reason naming the worktree.
5. `ApplyItemBranch` re-inspects the branch, refuses when its head, delivery
   commit, default head or worktree state changed, removes the clean worktree
   with `git worktree remove` (no force) and deletes the branch with
   `git branch -D`.

### Data Models

No schema change. The Run Database records a released Run as today through
`store.IntegrationReconciliation`. Item branches have no record; the report
prints each head before `--apply` deletes it (ADR-0232).

### API Contracts

1. API Contract: merge-evidence release. Input: a terminal spec Run whose
   Spec is archived on the default branch by a delivery commit the Run Branch
   lacks. Output: classification `superseded`, the reason of Merge evidence
   step 3, action `would release with --apply`, and after `--apply` the Run
   Worktree and Run Branch removed with an Integration Reconciliation
   recorded. Failure: a Task commit of another Spec or with a Task not
   completed keeps `unintegrated` naming the commit; a dirty path outside the
   archived scope keeps `dirty`.
2. API Contract: post-merge cleanup. `ReleaseMergedRuns` accepts a merge
   commit equal to the `git merge-tree --write-tree` merge of the candidate
   into the merge commit's first parent; every other unproven candidate
   returns `release merged Spec "<slug>" Runs: candidate head "<head>" is not represented by merge commit "<merge>"`.
3. API Contract: item branch candidates. The `roundfix-reconcile/v1` report
   gains `itemBranchCandidates`, a list of debris entries with `kind`
   `itemBranch`, `specSlug`, `itemBranch`, `head`, `worktree`, `proof`,
   `action` and `refusalReason`; `debrisSummary` gains
   `itemBranchCandidates` and `itemBranchesApplied`. The list is present,
   possibly empty, on a full scan and absent when a Run ID selects one Run,
   because only a full scan inspects item branches. A releasable branch's
   action is `would reclaim with --apply`, then `released`; a preserved one
   appears in `preservedCandidates` with `kind` `itemBranch`. Existing keys
   keep their meaning, and the other debris entries omit the two new fields.

### Surface Transcripts

1. Surface Transcript: a dry-run in a repository holding an item branch of a
   merged Spec lists it and counts it.

   ```transcript
   $ roundfix reconcile
   stdout:
   Repository: <repository>
   Mode: dry-run
   ...
   Item branch candidate: roundfix/deliver-<slug> (Spec <spec>)
     head: <head>
     worktree: -
     proof: item branch is superseded by the delivery of Spec <proof>
     action: would reclaim with --apply
     refusal-reason: -
   ...
   Debris summary: <counts> item-branch-candidates=1 item-branches-applied=0
   Apply with: roundfix reconcile --apply
   stderr:
   exit: 0
   ```

## Vocabulary Contract

No glossary edit in this Spec (`CONTEXT.md` is not touched). The Spec uses
**Run Worktree Reconciliation**, **Reconcile Command**, **Delivery Queue** and
**Branch Disposition** as the glossary defines them; "merge evidence",
"delivery commit" and "item branch" are descriptive phrases. The QA gate
records whether the glossary's Run Worktree Reconciliation entry needs merge
evidence named, as a candidate term for the maintainer.

## Coverage Map

- Goal 1 → Merge evidence; API Contract 1.
- Goal 2 → The cleanup proof for a moved default branch; API Contract 2.
- Goal 3 → Item branch reconciliation; API Contract 3; Surface Transcript 1.
- Goal 4 → Merge evidence step 5; Item branch reconciliation step 4.
- Core Feature 1 → Merge evidence.
- Core Feature 2 → The cleanup proof for a moved default branch.
- Core Feature 3 → Item branch reconciliation.
- Core Feature 4 → Build Order 4.
- Success Metric 1 → Testing Approach 1 and 2.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 1 to 4.

## Integration Points

- Local Git only: `rev-list`, `log --diff-filter=A`, `merge-base
  --is-ancestor`, `merge-tree --write-tree`, `for-each-ref`, `worktree`
  and `branch -D`. The Delivery Queue merge record and live items come from
  the Run Database through the existing reader.
- No forge, provider or network boundary.

## Testing Approach

1. **Merge evidence in the worktree package**, new file
   `internal/worktree/merge_evidence_test.go`, with the fixtures of
   `merged_head_test.go` and `merged_spec_leftovers_test.go`: a Run whose
   non-Task commit touched a file the default branch later changed and
   another it renamed is `superseded` with the delivery-commit reason, with
   and without a merge record; it stays `unintegrated` when a Task commit's
   Task is pending, when the Spec has no archived `_prd.md`, and when the Run
   Branch holds the delivery commit; a dirty path outside the scope stays
   `dirty`; apply refuses after the default branch moves. The archived-Spec
   case `TestArchivedSpecDirectoryProofStillComparesOtherCommitPaths` changes
   its expectation to `superseded`, the one declared break here; the
   non-archived refusals of `merged_head_test.go` pass unedited.
2. **Reconcile end to end**, new file
   `internal/cli/reconcile_merge_evidence_test.go`, through the public
   command runner on a disposable repository and Run Database: the PRD's
   measured fixture (an item-branch operator commit inherited by the Run, a
   later Spec editing one file and renaming another, no merge record) reports
   `superseded` in dry-run and `released` after `--apply`.
3. **Cleanup proof**, new file `internal/cli/deliver_release_moved_main_test.go`:
   a squash merge of the candidate after another default-branch commit lets
   `ReleaseMergedRuns` release the Spec's Run; a conflict resolved
   differently in the merge commit still refuses. The tests of
   `deliver_release_evidence_test.go` pass unedited.
4. **Item branches**, new files
   `internal/worktree/item_branch_reconcile_test.go` and
   `internal/cli/reconcile_item_branch_test.go`: an item branch of a merged
   Spec is listed and then deleted by `--apply`, with its clean worktree; the
   branch of a live item, a dirty item worktree, an item branch of an
   unmerged Spec and a branch whose head moved after inspection stay; a full
   scan's JSON carries `itemBranchCandidates` and a Run-ID report does not.
   The JSON key list asserted in `internal/cli/cli_test.go`, a Governed Path
   this Spec does not touch, reads a Run-ID report and passes unedited.

## Build Order

1. Merge evidence for terminal Runs, with Testing Approach 1 and 2, task_01
   (depends on: 4, which writes the guide of the reconcile surface it
   changes).
2. The cleanup proof for a moved default branch, with Testing Approach 3,
   task_02 (depends on: 4, which writes the deliver guide sentence it
   changes).
3. Item branch reconciliation, with Testing Approach 4, task_03 (depends on:
   1, whose merge evidence it reuses, and 4).
4. The `reconcile` and `deliver` guides, the Roundfix Skill's `reconcile`
   reference, the skill version and its mirrors, written from this TechSpec,
   task_04 (depends on: none).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **Committed work released with a merged Spec.** A hand commit in a Run of a
  merged Spec that never reached the delivery is released. The Run head stays
  in the Run Database's reconciliation record and the item head in the
  report, so the commits can be restored until Git collects them; dirty work
  is never released outside the declared scope.
- **A Spec archived without its delivery.** An archive commit pushed to the
  default branch outside a delivery would count as delivery. Archiving moves
  the Spec into history only as the last commit of its implementation
  delivery, so the default branch holding it is that delivery.
- **Old Git.** Git before `--write-tree` keeps today's cleanup warning; the
  Runs are released later by reconcile through merge evidence.
- **Concurrent skill edits.** Specs 0228 and 0229 may raise the same skill
  version; task_04 raises it by one patch level from the tree it starts on.

## Decisions

- Merge evidence is a fallback after today's proof, so existing reasons and
  tests keep their bytes. See ADR-0232.
- The delivery commit is found by the path it added, not by its message:
  0181's and 0184's squash merges do not name their slug.
- Dirty work keeps ADR-0212's scope rule, and a dirty item worktree keeps its
  branch. See ADR-0232.
- Item branches are inspected only by a full scan and released only by
  `--apply`, with no Branch Disposition record. See ADR-0232.
