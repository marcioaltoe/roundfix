---
spec: 0232-a-queue-that-sees-a-pull-request-merged-by-hand
prd: _prd.md
created: 2026-10-05
---

# A queue that sees a Pull Request merged by hand — Technical Spec

## Executive Summary

A Delivery Retry of a parked item now asks a Merge Observer whether the item
is already merged before it applies the Delivery Queue Limits, touches the
item workspace or runs any existing retry rule. The observer reads the item's
recorded Pull Request through `gh`, and otherwise looks for ADR-0232's merge
evidence in local Git. A merged item is recorded `merged` through a new
compare-and-set in the store, and the owner's existing pass runs the
post-merge cleanup and releases waiting items. The primary trade-off is that
the merge is recognized only when the operator retries the item, not by
`deliver status` or by the owner on its own: status stays an offline read and
the owner adds no GitHub reads per parked item, at the cost of one explicit
command after a hand merge, which is the command the Pending Question already
asks for (ADR-0237).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; queue
  items, stages, blockers and Pull Requests keep their names and numbers.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — `ObserveMerge` adds one read through
  the authenticated `gh` CLI the queue already uses, `gh pr view <n> --json`
  in the repository checkout. No credential is read, stored or sent by
  Roundfix; no test reaches GitHub, because every `gh` call in a test goes
  through a scripted `CommandRunner`. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0237: "A Delivery Retry of a
  parked item now first asks whether the item is already merged", and "A
  recorded Pull Request that was closed without merging and has no merge
  evidence keeps the item parked". ADR-0232: "A Spec now carries merge
  evidence when the default-branch head holds its archived `_prd.md`"; the
  observer reuses that proof and the cleanup keeps its tree proof. For an
  item that is not merged, ADR-0223 keeps its rule, ADR-0223: "records that head as a new candidate". ADR-0229 keeps its rule, ADR-0229: "It records that head as
  the candidate and resumes". After a merge the owner still
  applies ADR-0193: "The owner itself returns such an item to `queued` once
  every prerequisite is met". ADR-0199 acts "when a parked item is retried";
  ADR-0237 refines it for a recorded merge, which resumes no work.
  ADR-0184 binds Surface Transcripts 1 and 2. ADR-0179 bounds the Governed
  Paths. ADR-0189 and ADR-0233 bind the Roundfix Skill's version raise. The
  QA gate follows ADR-0080: "QA verdicts distinguish environment-blocked
  rows", and ADR-0091: "required to be terminal and to depend on every leaf";
  ADR-0096, ADR-0097, ADR-0104, ADR-0167, ADR-0182, ADR-0194, ADR-0195 and
  ADR-0210 bind its stage, row carry, outside evidence, Pull Request row,
  settlement, row record, re-observation and evidence snapshot. ADR-0093,
  ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check consistency by
  citation and receipt. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_02 changes the Governed Paths
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md` and
  `skills/roundfix/SKILL.md`, and no other; express maintainer
  authorization: the standing grant of 2026-09-30, "considere autorizado a
  ajustar todas as skills se necessário"; bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0232-a-queue-that-sees-a-pull-request-merged-by-hand/_authorization.md`.

## System Architecture

No new package, command or flag.

| Component | Where | Change |
| --- | --- | --- |
| Merge Observer seam | `MergeObservation`, `MergeObserver`, `EngineDependencies.Merges` in `internal/delivery/engine.go` | New optional dependency; nil keeps today's retry |
| Delivery Retry | `Engine.Retry` in `internal/delivery/engine.go` | Asks the observer first; records a merged item or refuses a closed one |
| Pull Request read | `GitHubCLI.ViewPullRequest`, `PullRequest.Merged` in `internal/delivery/github.go` | Exports the existing `gh pr view` read and merge test |
| Merged record | `Store.RecordDeliveryQueueItemMerged` in `internal/store/delivery.go` | Compare-and-set from the observed park to `merged` |
| Merge evidence | `ProveDelivery` in `internal/worktree/merged_head.go` | Exports ADR-0232's `provenDeliveryEvidence` |
| Observer adapter | `commandDeliveryWorkflow.ObserveMerge` in `internal/cli/deliver_workflow.go` | Pull Request first, then merge evidence; wired in `newCommandDeliveryEngine` |
| Retry output | `printDeliverRetryResult` in `internal/cli/deliver.go` | Prints the merge line before the retry line |
| Records | `docs/user-guide/commands/deliver.md`, the Roundfix Skill's deliver reference | Describe the rule, output and refusal |

## Implementation Design

### Interfaces

```go
// internal/delivery/engine.go
type MergeObservation struct {
	Merged         bool
	MergeCommit    string // the merge commit, or the delivery commit
	Head           string // the merged head; becomes the newest candidate
	Evidence       string // "pull request #404" or `Spec archived on default branch "main"`
	ClosedUnmerged bool   // the recorded Pull Request is closed without merging
}

type MergeObserver interface {
	ObserveMerge(ctx context.Context, gitRoot string, item store.DeliveryQueueItem) (MergeObservation, error)
}

// EngineDependencies gains: Merges MergeObserver
// RetryResult gains: Merge MergeObservation

// internal/delivery/github.go
func (client GitHubCLI) ViewPullRequest(ctx context.Context, number string) (PullRequest, error)
func (pullRequest PullRequest) Merged() bool

// internal/store/delivery.go
func (store *Store) RecordDeliveryQueueItemMerged(ctx context.Context, gitRoot string, item DeliveryQueueItem, parkedBlocker string) (ownerPID int, ownerIdentity string, err error)

// internal/worktree/merged_head.go
type DeliveryProof struct{ DefaultBranch, DefaultHead, DeliveryCommit string }
func ProveDelivery(ctx context.Context, gitRoot, specSlug, head string) (DeliveryProof, bool)
```

```text
1. Engine.Retry calls Merges.ObserveMerge only for an item whose stage is parked, after the stage check and before the Queue Token Ceiling, queue deadline and retry-limit checks, the prerequisite rule, UseItemBranch and every archived-head, conflict and corrective rule. With Merges nil, Retry is byte-for-byte today's behavior.
2. An observer error refuses the retry with "retry Delivery Queue item <q-slug>: observe merge: <error>" and leaves the item unchanged.
3. An observation with Merged set and an empty MergeCommit or Head is an error under Invariant 2.
4. A Merged observation is recorded through RecordDeliveryQueueItemMerged: the stored item must still be parked with the observed blocker; it becomes stage merged with blocker "", merge_commit set, Head appended to candidate_commits unless it is already the newest, pull_request_number, run_id, branch, worktree, warning and retry_count unchanged; the queue's owner PID and identity are returned. No retry limit, deadline or token ceiling refuses it, and no workspace, Carry-Forward, revalidation or review action runs.
5. Retry then returns RetryResult{Blocker: <park>, Stage: merged, Merge: <observation>, OwnerPID, OwnerIdentity}, and the CLI hands the queue to the live owner or starts one exactly as for any retry.
6. An observation with ClosedUnmerged set and Merged unset refuses the retry with API Contract 3 and leaves the item unchanged.
7. Any other observation continues into today's retry rules unchanged.
8. ObserveMerge, with a recorded Pull Request number, reads it with ViewPullRequest in gitRoot. It is merged evidence when PullRequest.Merged() is true, its HeadBranch equals the item's Branch, and its MergeCommit and HeadSHA are not empty; Evidence is "pull request #<n>", Head the PR head, MergeCommit its merge commit. A read failure is an observer error.
9. Otherwise ObserveMerge takes the item head as the local item branch's tip when that branch exists, else the newest candidate commit, and asks ProveDelivery(gitRoot, slug, head). A proof gives Merged with MergeCommit = DeliveryCommit, Head = that item head, Evidence `Spec archived on default branch "<branch>"`. No head, or no proof, is not merged; a Git failure in the proof is no proof.
10. ClosedUnmerged is set when the recorded Pull Request's state is CLOSED, it is not merged, and Invariant 9 found no proof.
11. ProveDelivery returns exactly what provenDeliveryEvidence decides for the local default branch: the archived _prd.md is on its head and the commit that added it is not reachable from head.
12. The owner's pass is unchanged: a merged item gets ReleaseMergedRuns and RemoveItemBranch, a failure becomes the existing cleanup warning, and releasePrerequisites runs after it. deliver status and deliver resume are unchanged.
```

### Data Models

No schema change. A recorded merge writes the existing `stage`, `blocker`,
`merge_commit` and `candidate_commits` columns of `delivery_queue_items`;
`retry_count` is untouched. `RetryResult` gains the observation it acted on.

### API Contracts

1. API Contract: the observer reads a recorded Pull Request with
   `gh pr view <n> --json number,url,state,headRefName,headRefOid,mergedAt,mergeCommit,mergeable,mergeStateStatus`,
   the existing field list, in the repository checkout.
2. API Contract: a retry that records a merge prints, on stdout, before the
   existing `Retried <slug>: <blocker> -> merged` line, exactly one line
   `Merged outside the queue: <evidence>; merge commit <sha>`, where
   `<evidence>` is `pull request #<n>` or
   `Spec archived on default branch "<branch>"` and `<sha>` is the full
   recorded merge commit.
3. API Contract: a retry refused for a closed Pull Request has the reason
   `retry Delivery Queue item "<slug>": pull request #<n> was closed without merging; reopen it or merge the Spec into the default branch, then run roundfix deliver retry <slug>`,
   printed in the existing `Retry refused` form with exit `2`.
4. API Contract: every other retry output, refusal, exit code and the owner
   hand-off lines keep their text.

### Surface Transcripts

1. Surface Transcript: a retry of an item whose recorded Pull Request was
   merged by hand records the merge and hands the queue to its owner.

   ```transcript
   $ roundfix deliver retry 0300-example
   stdout:
   Merged outside the queue: pull request #<number>; merge commit <sha>
   Retried 0300-example: checks-failed -> merged
   ...
   stderr:
   ...
   exit: 0
   ```

2. Surface Transcript: a retry of an item whose recorded Pull Request was
   closed without merging is refused and leaves the item parked.

   ```transcript
   $ roundfix deliver retry 0300-example
   stdout:
   stderr:
   Retry refused
   ...
     retry Delivery Queue item "0300-example": pull request #<number> was closed without merging; reopen it or merge the Spec into the default branch, then run roundfix deliver retry 0300-example
   ...
     stage: parked; blocker: checks-failed
   ...
   exit: 2
   ```

## Coverage Map

- Goal 1 → `Engine.Retry`, `RecordDeliveryQueueItemMerged`, `ObserveMerge` Pull Request read (Invariants 1, 4, 5, 8, 12)
- Goal 2 → `ObserveMerge` merge evidence, `ProveDelivery` (Invariants 9, 11)
- Goal 3 → `Engine.Retry` closed refusal (Invariants 6, 10; API Contract 3)
- Goal 4 → nil observer and not-merged path (Invariants 1, 7, 12; API Contract 4)
- Core Feature 1 → `MergeObserver`, `Engine.Retry`, `ViewPullRequest`, `RecordDeliveryQueueItemMerged`
- Core Feature 2 → `ProveDelivery`, `ObserveMerge`
- Core Feature 3 → owner pass (unchanged), `printDeliverRetryResult` (API Contract 2)
- Core Feature 4 → deliver guide, Roundfix Skill deliver reference
- Success Metric 1 → CLI replay of 0231 with scripted `gh` and a removed item branch
- Success Metric 2 → CLI replay of 0225 with a squash delivery commit in a disposable repository
- Success Metric 3 → engine and CLI closed-unmerged tests (Surface Transcript 2)
- Success Metric 4 → existing retry tests; nil-observer test

## Integration Points

- GitHub through `gh pr view`, the read `MergePullRequest` already makes, now
  also exported as `ViewPullRequest`.
- Git: ADR-0232's merge evidence on the local default branch; the owner's
  cleanup keeps its fetch of the delivery remote's default branch.

## Testing Approach

- The engine is tested in a new `internal/delivery/merged_outside_test.go`
  through `Engine.Retry` with a fake `MergeObserver`, the real store in a
  temporary home, and the existing fake workspace and recovery seams, which
  record that no workspace call happens for a merged item.
- `ViewPullRequest` is tested with the existing `scriptedCommandRunner`.
- The store method is tested in a new `internal/store/delivery_merged_test.go`.
- The adapter and output are tested in a new
  `internal/cli/deliver_merged_outside_test.go`: disposable Git repositories
  with a bare origin, a `commandDeliveryWorkflow` whose `GitHubCLI` uses a
  scripted runner, an engine built like `newCommandDeliveryEngine` with that
  workflow, `runDeliverRetry` through the existing dependency overrides with
  the owner start replaced, and one `Engine.Run` pass to prove the cleanup.
  No test runs the real `gh` or reads the real `~/.roundfix`.

## Build Order

1. Record a merged item from a Merge Observer in the Delivery Retry, the
   store and the Pull Request read, and state the rule and the closed refusal
   in the deliver guide.
2. Observe merges through `gh` and ADR-0232's merge evidence, print the merge
   line, and describe it in the deliver guide and the Roundfix Skill
   (depends on: 1).
3. Final QA gate (depends on: 1, 2).

## Risks & Considerations

- A retry with a recorded Pull Request now needs `gh`; when the read fails the
  retry is refused rather than republishing work that may be merged. The
  existing check stage already needs `gh` for such an item.
- Merge evidence reads the local default branch, so after a hand merge the
  operator's checkout must have pulled it; the Pull Request path does not
  need it.
- A merged head that the cleanup cannot prove against the merge commit, for
  example after a hand merge of more than the item branch, stays a cleanup
  warning and the item stays `merged`.
- The Roundfix Skill version is raised by the record command, which picks a
  free version, so a concurrent raise in another Spec resolves at merge
  (ADR-0233).

## Vocabulary Contract

- emits: `internal/delivery/engine.go`
  pattern: `was closed without merging`
  documented-in: `docs/user-guide/commands/deliver.md`
- emits: `internal/cli/deliver.go`
  pattern: `Merged outside the queue`
  documented-in: `docs/user-guide/commands/deliver.md`

No glossary term is adopted; "merged outside the queue" is used in its plain
sense beside **Delivery Retry** and **Delivery Queue**.

## Decisions

- Recognize the merge only in `deliver retry`. See ADR-0237.
- Prefer the recorded Pull Request, and fall back to ADR-0232's merge
  evidence for Pull Requests the queue never recorded.
- Refuse a retry whose recorded Pull Request is closed without merging or
  cannot be read, rather than resuming a publication.
- Leave the retry count unchanged for a recorded merge, because no work
  resumes.
