---
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
prd: _prd.md
created: 2026-11-13
---

# A parked item always has a way back through a Delivery Retry — Technical Spec

## Executive Summary

`Engine.Retry` returns a parked item whose recorded Item Branch is empty to
`queued`, on the path a prerequisite park already takes, and
`store.RetryDeliveryQueueItem` clears a worktree recorded without its branch
when it requeues. `commandDeliveryWorkflow.Authorization` resolves the Spec
Root with `ResolveSpecsRootOrAbsent`, and `spec.ReadSpecAuthorization` reads
an absent Spec Root through its nearest existing ancestor, so the grant is
read from Git history whatever the working tree holds. `checkCandidate`
returns an item whose Pull Request holds an earlier candidate to `gating`,
where the repository gate runs again and publication pushes the current
candidate; the `review-stale` Park Class gains a next action. The trade-off
is one more gate run for a candidate the queue already gated, in exchange for
never pushing an ungated candidate (ADR-0398).

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0398 (new), ADR-0193 and ADR-0223
  apply; ADR-0229, ADR-0237, ADR-0262 and ADR-0312 do not (see ADR
  Alternatives). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BE answers of 2026-11-13 quoted in the authorization record; the
  Governed Paths are the Roundfix skill files task_04 changes. No Makefile,
  workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured at `f953f77f`:

1. `Engine.Retry` (`internal/delivery/engine.go`) returns an item to `queued`
   only for `archived-item-unproven` and `prerequisite-unmerged` parks; every
   other park calls `ItemWorkspace.UseItemBranch`, and
   `commandDeliveryWorkflow.UseItemBranch` refuses an empty branch with `use
   item branch: branch is required`.
2. `advanceItem` assigns `item.Branch` only after `CreateItemBranch` returns.
   `CreateItemBranch` records the branch and worktree pair with
   `RecordDeliveryQueueItemWorktree` before it creates the worktree, so a
   failure after that point parks through `Engine.park`, whose
   `UpdateDeliveryQueueItem` writes the in-memory empty branch over the
   recorded one and leaves the recorded worktree.
   `RecordDeliveryQueueItemWorktree` only fills a row whose branch and
   worktree are both empty, and refuses a row it reads back without a branch.
3. `RetryDeliveryQueueItem` (`internal/store/delivery.go`) counts a retry for
   every blocker except `runtime-infrastructure` and `disk-full`, and changes
   neither branch nor worktree.
4. `commandDeliveryWorkflow.Authorization` (`internal/cli/deliver_workflow.go`)
   resolves the Spec Root with the strict `roundconfig.ResolveSpecsRoot`, which
   refuses an absent directory. `spec.ReadSpecAuthorization`'s
   `resolveAuthorizationLocation` runs `git -C <spec root> rev-parse
   --show-toplevel` through `authorizationRepositoryLocation`, which fails on a
   missing directory, so the resolution is `unresolved` and `engine.authorized`
   would park `unauthorized`.
5. `checkCandidate` parks `review-stale` whenever the current-head check
   report's head differs from the item's current candidate. An archived item
   with an unchanged candidate and a recorded Pull Request re-enters at
   `checking`, so a candidate recorded by an earlier retry but never pushed
   parks there.
6. `ClassifyPark` gives `review-stale` the generic answer `resolve the
   blocker, then run roundfix deliver retry <slug>`.

## System Architecture

No new package, command, flag or file format. The changes stay in
`internal/delivery` (retry, checking, park answer), `internal/store` (the
requeue write), `internal/cli` (the authorization read) and `internal/spec`
(the authorization location).

## Implementation Design

### Interfaces

```go
// internal/delivery/engine.go
func earlierCandidate(item store.DeliveryQueueItem, head string) bool

// internal/spec/authorization.go
func authorizationSpecsRootLocation(ctx context.Context, specsRoot string) (repoRoot, prefix, commonDir string, err error)
```

### Data Models

None: no column, record or config key changes. A requeued item's `worktree`
and `worktree_provisioned` are cleared only when its stored branch is empty.

### API Contracts

1. API Contract 1 — the branchless retry. In `Engine.Retry`, after the merge
   observation, the token, deadline and retry-limit checks, the condition
   that returns an item to `queued` becomes
   `blockerMatches(item.Blocker, BlockerArchivedItemUnproven) ||
   prerequisitePark(item.Blocker) || strings.TrimSpace(item.Branch) == ""`,
   with a comment that an item parked before it has an Item Branch never
   started and the owner creates its branch as on a first start. Such a retry
   reads no item evidence and calls no workspace method. The retry is counted
   as `RetryDeliveryQueueItem` already counts it.
2. API Contract 2 — the requeue write. In `RetryDeliveryQueueItem`, inside its
   transaction and before the stage update, when the target stage is
   `DeliveryStageQueued` it runs `UPDATE delivery_queue_items SET worktree =
   '', worktree_provisioned = 0 WHERE git_root = ? AND position = ? AND
   branch = ''`, with a comment that a worktree recorded without its branch is
   half of a pair the owner never finished. A row with a branch is unchanged.
3. API Contract 3 — the authorization read. `Authorization` resolves the Spec
   Root with `roundconfig.ResolveSpecsRootOrAbsent(workflow.loaded, gitRoot)`,
   with a comment that the grant is read from Git history.
   `resolveAuthorizationLocation` calls `authorizationSpecsRootLocation` for
   the Spec Root, which walks up with `os.Lstat` to the nearest existing
   ancestor, calls `authorizationRepositoryLocation` there and returns its
   prefix joined with the missing components in order; any `Lstat` error other
   than not-exist, or reaching the filesystem root, falls back to
   `authorizationRepositoryLocation(ctx, specsRoot)`. An existing Spec Root
   resolves exactly as today. Nothing is created on disk.
4. API Contract 4 — the checking stage. In `checkCandidate`, when the report's
   head differs from the current candidate and `earlierCandidate(*item,
   report.HeadSHA)` holds, the owner writes to its log
   `roundfix: Delivery Queue item <slug>: pull request #<n> holds earlier candidate <pull request head>, not candidate <candidate>; returning to gating to push it`
   and sets the stage to `gating`; otherwise it parks `review-stale` as today.
   `earlierCandidate` is true only for a non-empty head equal to an entry of
   `CandidateCommits` other than the last.
5. API Contract 5 — the park answer. `ClassifyPark` keeps `review-stale` in
   the `review` Park Class and sets its next action to
   `in <worktree>, merge into the item branch any commit its Pull Request's head branch holds that the item branch lacks, then run roundfix deliver retry <slug>`.
   `review-blocked` keeps the generic answer.
6. API Contract 6 — the `deliver` guide. In
   `docs/user-guide/commands/deliver.md`: (a) in the retry table, before the
   row `| Active Spec with any unfinished Task | `running` |`, the row
   `| Item parked before it has an Item Branch | `queued`, counting a retry |`,
   and right after the table, before "After the retry, a live owner", the
   paragraph Text A; (b) right after "not hide that pre-archive grant.", Text
   B; (c) right after Text A, its own paragraph, Text C; (d) after the
   paragraph that begins "A `secret-found` park's next action is", its own
   paragraph, Text D. The texts, word for word:
   - Text A: "An item parked before it has an Item Branch re-enters at
     `queued`, counting a retry, and the owner creates its branch as on a
     first start."
   - Text B: "The grant is read from Git history, so it needs no Spec Root in
     the working tree: when the last file under the Spec Root leaves the item
     branch, as after a conflict resolution that keeps the archive, Git
     removes the empty directory and the read still finds the grant."
   - Text C: "At `checking`, a Pull Request whose head is an earlier candidate
     of the item never received the current one: the owner logs `roundfix:
     Delivery Queue item <slug>: pull request #<n> holds earlier candidate
     <sha>, not candidate <sha>; returning to gating to push it`, and the item
     returns to `gating`, whose repository gate runs again before the push.
     Any other Pull Request head parks the item `review-stale`."
   - Text D: "A `review-stale` park's next action is `in <worktree>, merge
     into the item branch any commit its Pull Request's head branch holds that
     the item branch lacks, then run roundfix deliver retry <slug>`."
7. API Contract 7 — the Roundfix skill. Both copies of `references/deliver.md`
   carry Text A followed by a blank line and Text C right after the sentence
   that ends "or at `checking` with one.", Text B right after the sentence
   that ends "`push`, `pull_request`, and `merge`.", and Text D as its own
   paragraph after the paragraph that begins "A `secret-found` park's next
   action is". Text C and Text D stay on one line each in the reference.
8. API Contract 8 — the glossary. The **Delivery Retry** entry of
   `CONTEXT.md` gains, at the end of its definition paragraph: "An item parked
   before it has an Item Branch returns to `queued`, and its owner creates the
   branch as on a first start (ADR-0398).", where `ADR-0398` is a Markdown
   link to
   `docs/adr/0398-a-parked-item-always-has-a-way-back-through-a-delivery-retry.md`,
   as the entry's existing ADR link is written.

### Surface Transcripts

None. No command, flag or exit code is added. The one changed output line,
the `review-stale` `Park:` line and its Pending Question answer, is already
asserted whole, stdout, stderr and exit code, by
`TestDeliverStatusPrintsTheItemWorktree` and
`TestDeliverStatusPrintsAParkLinePerParkedItem`, which task_03 updates; the
`Retried <slug>: <blocker> -> queued` line is the existing retry line with a
stage it already prints for a prerequisite park.

### Invariants

1. A retry of an item with a recorded Item Branch takes exactly the paths it
   takes today; only an item with no recorded branch is requeued.
2. A requeue clears a worktree only from a row without a branch; a row with a
   branch keeps its workspace.
3. The delivery authorization is still read at the parent of the archive
   commit; no read creates a directory, and an absent Spec Root outside the
   repository is still refused.
4. The owner pushes only a candidate the repository gate passed in the same
   pass; `deliver retry` never pushes.
5. A Pull Request head that is not a candidate of the item still parks
   `review-stale` with no push.

### Declared breaks

- `TestExistingGenericParkAnswersAreUnchanged`
  (`internal/delivery/park_class_test.go`) drops `BlockerReviewStale` from its
  list; every other blocker keeps the generic answer.
- `TestPendingQuestionAnswersEachBlockerClass`
  (`internal/delivery/question_test.go`): its `other blocker` row uses
  `BlockerReviewBlocked` instead of `BlockerReviewStale`, with the same
  answer.
- `TestDeliverStatusPrintsTheItemWorktree` (`internal/cli/deliver_test.go`):
  its `Park:` line and `Answer:` line carry Text D's action for
  `/worktrees/delivery-item`.
- `TestDeliverStatusPrintsAParkLinePerParkedItem`
  (`internal/cli/deliver_park_status_test.go`): its `Park: last review:` line
  carries Text D's action with the item's empty worktree, `in , merge into
  ...`.
- A Task that finds another test pinning the generic `review-stale` answer on
  the default branch when it runs changes it the same way and says so.
- `retry Delivery Queue item "<slug>": use item branch "": use item branch:
  branch is required` is no longer produced by a retry; no test or document
  pins it.

## Coverage Map

- Goal 1, User Story 1, Core Feature 1, Success Metric 1 → API Contracts 1
  and 2, Invariants 1 and 2.
- Goal 2, User Story 2, Core Feature 2, Success Metric 2 → API Contract 3,
  Invariant 3.
- Goal 3, User Story 3, Core Feature 3, Success Metric 3 → API Contract 4,
  Invariants 4 and 5.
- Goal 4, User Story 4, Core Feature 4, Success Metric 4 → API Contract 5.
- Core Feature 5, Success Metric 5 → API Contracts 6 to 8.

Deviations from the PRD: None.

## Absence Evidence

1. `earlierCandidate` — command: `grep -rn 'func earlierCandidate\|returning to gating' internal/delivery`; result: no match at `f953f77f`.
2. `authorizationSpecsRootLocation` — command: `grep -rn 'authorizationSpecsRootLocation\|nearest existing ancestor' internal/spec`; result: no match at `f953f77f`.

## ADR Alternatives

- ADR-0398: the design is the decision; it takes none of the six options the
  record rejects.
- ADR-0193: it rejected treating the owner's release of a prerequisite park
  as a Delivery Retry; that release stays uncounted, and the explicit retry of
  a branchless item is counted like the explicit retry of a prerequisite park.
- ADR-0223: it rejected rewriting other Specs' paths in the archive commit;
  the archive stays an exact move, and a post-archive correction still
  returns to review before publication.
- ADR-0229: not applicable — it resumes an operator-archived item from its
  Run start, a retry path this Spec does not change.
- ADR-0237: not applicable — it records a merge made outside the queue, which
  the retry still observes first and this Spec does not change.
- ADR-0262: not applicable — it governs `roundfix cancel` and the retry's
  refusal of a cancelled Spec, which this Spec does not change.
- ADR-0312: not applicable — it resumes a new queue's archived item branch at
  review, which this Spec does not change.

## Integration Points

- `deliver retry` builds the production engine with its terminal's standard
  error; the detached owner's log is its console log, where the line of API
  Contract 4 appears.
- The owner's next pass runs `advanceItem` on the requeued item: the
  prerequisite check, `CreateItemBranch`, the archived-item proof and
  revalidation, as for a new item.

## Testing Approach

New tests, each calling `t.Parallel()` first:

- `internal/delivery/branchless_retry_test.go`:
  `TestRetryRequeuesAnItemParkedBeforeItHasABranch` (an item parked
  `delivery-error: read prerequisites: refresh prerequisite default branch
  "main": git fetch origin main: Operation timed out` with no branch, on
  `newRetryDeliveryEngine` with `fakeItemRecovery`: the result and the stored
  item are `queued`, blocker empty, retry count 1, no branch; the recovery was
  never asked and no workspace branch was used).
- `internal/cli/deliver_branchless_retry_test.go`:
  `TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch` (subtests `no
  workspace` and `worktree without branch`, on `newItemBranchStartRepository`
  with `newItemRecoveryWorkflowForRepository` and `worktree.location` under
  the test root; the second first records a pair with
  `RecordDeliveryQueueItemWorktree` and then parks with an empty branch; the
  production engine from `commandDeliveryEngineDependencies` retries; the
  stored item is `queued` with retry count 1 and no branch, worktree or
  provisioning; `CreateItemBranch` then returns a new Item Branch whose
  worktree head is `origin/main`).
- `internal/spec/authorization_absent_root_test.go`:
  `TestASpecAuthorizationIsReadAtARevisionWhenTheSpecsRootLeftTheWorkingTree`
  (a committed record granting `implement` and `push`, then `docs` removed
  from the working tree: granted, `push` permitted, the record path and
  revision unchanged).
- `internal/cli/deliver_operator_archive_test.go`:
  `TestTheDeliveryAuthorizationOfAnArchivedSpecNeedsNoSpecsRoot` (the
  archive fixture of `TestTheDeliveryAuthorizationIsReadBeforeTheArchiveCommit`,
  `git ls-files -- docs/specs` empty, `docs/specs` removed: `push`,
  `pull_request`, `merge`, and `docs/specs` still absent).
- `internal/delivery/unpushed_candidate_test.go`:
  `TestCheckingReturnsACandidateThePullRequestLacksToTheGate` (subtests
  `earlier candidate` and `foreign head`: an item at `checking` with
  candidates `pushed-head` and `resolved-head` and Pull Request `318`, on an
  engine built like `newTestDeliveryEngine` with `Log` set to a buffer and a
  `findPullRequest` hook whose Pull Request head follows the fake remote; the
  first merges after one push, the remote head is `resolved-head`, the
  workflow recorded `gate`, and the log holds the line of API Contract 4 for
  `#318`, `pushed-head` and `resolved-head`; the second parks `review-stale`
  with no push and no such line).
- `internal/delivery/park_class_test.go`:
  `TestAReviewStaleParkNamesItsNextAction` (class `review`, next action and
  Pending Question answer for worktree `/worktrees/example`, word for word).

Sabotages, each on a temporary copy of the audited head made with `git
archive` and committed in a fresh Git repository, each observed by the named
test on the named case (all run by the author on 2026-11-13 against the
prototype):

1. Drop `|| strings.TrimSpace(item.Branch) == ""` from the requeue condition →
   `TestRetryRequeuesAnItemParkedBeforeItHasABranch` and
   `TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch/no_workspace`.
2. Skip the worktree clearing of API Contract 2 →
   `TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch/worktree_without_branch`.
3. Resolve the Spec Root with `ResolveSpecsRoot` in `Authorization` →
   `TestTheDeliveryAuthorizationOfAnArchivedSpecNeedsNoSpecsRoot`.
4. Call `authorizationRepositoryLocation` instead of
   `authorizationSpecsRootLocation` for the Spec Root →
   `TestASpecAuthorizationIsReadAtARevisionWhenTheSpecsRootLeftTheWorkingTree`
   and `TestTheDeliveryAuthorizationOfAnArchivedSpecNeedsNoSpecsRoot`.
5. Never call `earlierCandidate` in `checkCandidate` →
   `TestCheckingReturnsACandidateThePullRequestLacksToTheGate/earlier_candidate`.
6. Make `earlierCandidate` true for any non-empty head →
   `TestCheckingReturnsACandidateThePullRequestLacksToTheGate/foreign_head`.
7. Return to `publishing` instead of `gating` →
   `TestCheckingReturnsACandidateThePullRequestLacksToTheGate/earlier_candidate`.
8. Restore the generic answer for `review-stale` →
   `TestAReviewStaleParkNamesItsNextAction`.

## Build Order

1. The branchless retry and the requeue write, their tests, the guide's table
   row and Text A, and the **Delivery Retry** glossary entry (task_01).
2. The authorization read, its tests and the guide's Text B (task_02)
   (depends on: 1, for the shared `deliver` guide and Behavior Surface Record).
3. The checking stage and the `review-stale` answer, their tests, the declared
   breaks and the guide's Texts C and D (task_03) (depends on: 2, for the
   shared guide and record, and on 1 for `internal/delivery/engine.go`).
4. The Roundfix `deliver` reference and the skill version (task_04) (depends
   on: 3).
5. The QA gate (task_05) (depends on: 4).

No other Spec is required. Spec 0388 delivers first and changes another
section of the Roundfix skill; the shared version fields and records are named
in `_prd.md` → Prerequisites. Each Task works on the tree it runs on, cites
functions rather than lines, and task_04 raises the skill version with the
record command.

## Risks & Considerations

- A requeued item whose failed worktree creation left a directory on disk
  gets a new Item Branch and worktree beside it; the leftover stays until the
  operator removes it (Non-Goals).
- A requeued item counts a retry, so a queue with a low retry limit can refuse
  a later retry of the same item; the limit already governs prerequisite
  parks the same way.
- The extra gate run costs one repository Verification when a Pull Request
  lags its candidate.
- A network failure still parks the item; the operator retries it.

## Glossary

- changes: **Delivery Retry**

## Decisions

- Requeue an item without an Item Branch, counted, and clear a half-recorded
  pair. See ADR-0398 and ADR-0193.
- Read the grant through the nearest existing ancestor of an absent Spec Root.
  See ADR-0398.
- Return a Pull Request behind its candidate to the gate, never push from the
  retry or from `checking`. See ADR-0398 and ADR-0223.
- Name the `review-stale` next action. See ADR-0398.
