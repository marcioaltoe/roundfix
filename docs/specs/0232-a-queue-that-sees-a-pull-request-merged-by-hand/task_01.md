---
task: task_01
spec: 0232-a-queue-that-sees-a-pull-request-merged-by-hand
status: completed
type: backend
complexity: medium
---

# Task 01: A Delivery Retry records a parked item as merged when a Merge Observer reports its merge

## Overview

The Backlog Entry of 2026-10-05
([a Pull Request merged by hand stays parked](references/2026-10-05-a-pull-request-merged-by-hand-stays-parked.md))
records that a Delivery Retry never asks whether a parked item's work is
already merged. This Task adds the Merge Observer seam to the Delivery
Engine, records a merged item through a new store compare-and-set, refuses a
retry whose recorded Pull Request was closed without merging, exports the
existing Pull Request read, and states the rule in the deliver command
reference (ADR-0237). The adapter that reads GitHub and Git is task_02's; here
the observer is a test fake.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-05 named in the Overview by adding
   `MergeObservation`, `MergeObserver`, `EngineDependencies.Merges` and
   `RetryResult.Merge` with the shapes in TechSpec → Interfaces, and by making
   `Engine.Retry` follow Invariants 1 to 7: the observer is asked only for a
   parked item, after the stage check and before the Queue Token Ceiling,
   queue deadline and retry-limit checks, the prerequisite rule, the item
   workspace and every archived-head, conflict and corrective rule; with
   `Merges` nil, `Retry` behaves exactly as before.
2. MUST add `Store.RecordDeliveryQueueItemMerged` per Invariant 4: in one
   write transaction it refuses unless the stored item at the item's position
   has the same slug, stage `parked` and the observed blocker; it sets stage
   `merged`, blocker empty, the merge commit and the candidate commits (the
   merged head appended unless it is already the newest), leaves
   `retry_count`, `run_id`, `branch`, `worktree`, `warning` and
   `pull_request_number` unchanged, applies no retry limit, and returns the
   queue's owner PID and identity as `RetryDeliveryQueueItem` does.
3. MUST refuse a merged observation with an empty merge commit or head, and an
   observer error, with `retry Delivery Queue item "<slug>": observe merge: <error>`
   (Invariants 2 and 3), and a closed-unmerged observation with the reason of
   API Contract 3, each leaving the stored item unchanged.
4. MUST export `GitHubCLI.ViewPullRequest`, the existing `gh pr view <n> --json`
   read of API Contract 1, and `PullRequest.Merged`, the existing merge test,
   without changing `PullRequestBoundary` or any existing call.
5. MUST describe in `docs/user-guide/commands/deliver.md`, in the retry
   section, that a retry first records a parked item `merged` when its
   recorded Pull Request was merged from the item branch or when the default
   branch already archives its Spec through a delivery commit outside the
   item branch, without counting a retry and regardless of the retry limit,
   deadline or token ceiling, and that the owner then runs the post-merge
   cleanup; add that case as the first row of the re-entry table; and quote
   the refusal of API Contract 3 verbatim with `<slug>` and `<n>`
   placeholders, also adding it to the list of retry refusals.
6. MUST add, in the new test files named in Verification, tests through
   `Engine.Retry` with a fake observer and the real store in a temporary home:
   a merged Pull Request is recorded `merged` with its merge commit and head,
   the retry count unchanged and no workspace, recovery or revalidation call;
   a merge is recorded at the retry limit, past the queue deadline and at the
   token ceiling; a closed-unmerged observation and an observer error are
   refused with the item unchanged; a not-merged observation and a nil
   observer follow today's rules; the store method refuses a changed blocker;
   and `ViewPullRequest` parses a merged and a closed Pull Request through the
   scripted command runner.
7. MUST NOT reach GitHub from any test, change `deliver status`, `deliver
   resume`, the owner's pass, the post-merge cleanup, the Run Database schema,
   `PullRequestBoundary`, or any skill.

## Subtasks

- [ ] Add the Merge Observer seam and the merged branch of `Engine.Retry`.
- [ ] Add the store compare-and-set for a recorded merge.
- [ ] Export the Pull Request read and merge test.
- [ ] Add the engine, store and GitHub read tests.
- [ ] Describe the rule and the refusal in the deliver guide.

## Acceptance Criteria

- [ ] A retry of a parked item whose observer reports a merge records it
      `merged` with the observed merge commit and head, leaves the retry
      count unchanged and touches no item workspace, whatever its Queue
      Limits.
- [ ] A closed-unmerged observation or an observer error refuses the retry
      and leaves the item unchanged.
- [ ] Every other retry, and every retry with no observer, behaves as before.
- [ ] The deliver guide states the rule, the table row and the refusal.

## Context

- instruction: `docs/adr/0237-a-delivery-retry-records-a-merge-made-outside-the-queue.md`
- instruction: `docs/adr/0223-a-delivery-keeps-its-work-across-archive-requeue-and-review.md`
- instruction: `docs/adr/0199-a-queue-token-ceiling-stops-new-starts-and-never-a-running-task.md`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/github.go`
- interface: `internal/store/delivery.go`
- interface: `internal/delivery/retry_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- creates: `internal/delivery/merged_outside_test.go`
- creates: `internal/store/delivery_merged_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestARetryRecordsAMergedPullRequestWithoutTouchingTheWorkspace|TestARetryRecordsAMergeAtTheRetryLimitDeadlineAndTokenCeiling|TestARetryOfAClosedUnmergedPullRequestIsRefusedAndLeavesTheItemParked|TestARetryWhoseMergeObservationFailsIsRefused|TestARetryOfAnUnmergedItemFollowsTheExistingRules|TestGitHubCLIViewPullRequestReportsMergedAndClosedPullRequests|TestRecordDeliveryQueueItemMergedRequiresTheObservedPark)$' ./internal/delivery ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestARetryRecordsAMergedPullRequestWithoutTouchingTheWorkspace TestARetryRecordsAMergeAtTheRetryLimitDeadlineAndTokenCeiling TestARetryOfAClosedUnmergedPullRequestIsRefusedAndLeavesTheItemParked TestARetryWhoseMergeObservationFailsIsRefused TestARetryOfAnUnmergedItemFollowsTheExistingRules TestGitHubCLIViewPullRequestReportsMergedAndClosedPullRequests TestRecordDeliveryQueueItemMergedRequiresTheObservedPark; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the seven tests exists, so their pass lines are missing and the command fails.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'pull request #<n> was closed without merging; reopen it or merge the Spec into the default branch, then run roundfix deliver retry <slug>' || { printf 'missing refusal in %s\n' docs/user-guide/commands/deliver.md >&2; exit 1; }; grep -qF -- 'was closed without merging' internal/delivery/engine.go && grep -qF -- 'func (store *Store) RecordDeliveryQueueItemMerged(' internal/store/delivery.go` — expected: exit 0; before this Task neither the guide nor the code names the refusal or the store method, so the command fails.

## References

- `_prd.md` → Goals; Core Feature 1; Success Metric 3; Success Metric 4
- `_techspec.md` → Interfaces; Invariants 1 to 7; Data Models; API Contract 1; API Contract 3; API Contract 4; Vocabulary Contract; Testing Approach; Build Order 1
- ADR-0237; ADR-0223; ADR-0199

## Result

Implemented task_01's Merge Observer seam and merge-recording path. A parked
item's merge observation runs before Queue Limits, prerequisite re-entry and
workspace recovery. A valid merge records `merged`, the merge commit and the
newest candidate head without counting a retry; it returns the observed park,
merge evidence and owner identity. The store compare-and-set requires the same
position, slug, parked stage and blocker, and updates only stage, blocker,
merge commit and candidate commits in one write transaction. Existing
`PullRequestBoundary` methods and calls remain unchanged; `ViewPullRequest`
and `Merged` expose the existing read and merge predicate.

Acceptance evidence:

- Merged retry and Queue Limits: tests in `merged_outside_test.go` cover
  checks, missing-workspace, prerequisite, corrective and conflict parks;
  retry allowance exhausted with an expired deadline; and recorded token
  usage at the ceiling. They compare the full persisted item, preserve its
  retry count and metadata, check owner PID/identity and the returned merge,
  and assert no workspace, recovery, revalidation or workflow action. A
  separate case preserves an already-newest candidate without appending it.
- Refusal with unchanged item: closed-unmerged observations assert the exact
  API Contract 3 reason. Observer errors retain `errors.Is` identity; missing
  or blank merge commits and heads receive the `observe merge` error prefix.
  Every refusal compares the full stored item and checks no recovery/workspace
  action. Store tests reject changed blockers, stages, slugs and positions.
- Existing retry behavior: nil and not-merged observers both re-enter
  `checking` with the existing workspace/inspection actions and one counted
  retry. The focused run also exercises existing retry and GitHub CLI tests.
  A non-parked item refuses before calling the observer. Scripted `gh pr view`
  cases parse merged and closed Pull Requests without reaching GitHub.
- Guide: the retry section states both merge evidence paths, unchanged retry
  count, Queue Limit bypass and owner cleanup; its first re-entry row is
  `merged`. The refusal is quoted verbatim with `<slug>` and `<n>` and included
  in the refusal list. A Python inspection confirmed those contract elements.

Focused checks:

- Baseline inspection of committed `HEAD` confirmed that `MergeObserver`,
  `RecordDeliveryQueueItemMerged` and `ViewPullRequest` were absent.
- `GOCACHE=/tmp/roundfix-task01-go-cache rtk proxy go test ./internal/delivery ./internal/store -run 'Test(ARetry|AMergeObserver|GitHubCLI|RecordDeliveryQueueItemMerged|Retry)' -count=1`
  exited 0; delivery and store both reported `ok` (12.069s and 5.702s).
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.
- `GOCACHE=/tmp/roundfix-task01-go-cache rtk make verify-incremental`
  exited 0 on the final unchanged code/test tree with process-table access:
  formatting, vet, the full Go suite, skill sync/check and build passed.
  The first attempt exited 2: my concurrent edits triggered the repository
  fingerprint guard, and two existing CLI force-stop tests lacked sandbox
  process-table access. Holding the worktree unchanged and rerunning with
  that access resolved both causes without code or configuration changes.
- The initial check using the shared Go build cache was blocked by sandbox
  cache access; the task-scoped `/tmp` cache resolved that environment issue.

The declared Verification commands were not run, status remains Daemon-owned,
and no commit, push or Pull Request was made. The real GitHub/Git observer,
CLI wiring/output and skill updates remain task_02's slice.

## Carry-forward provenance

- Source Run: `run_20261005T214431Z_77dc874bf3a7c619`
- Source commit: `37300f673412959996cac17c2e4ffe1c6478fc61`
