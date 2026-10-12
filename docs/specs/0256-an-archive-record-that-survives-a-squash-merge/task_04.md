---
task: task_04
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: pending
type: backend
complexity: high
---

# Task 04: The Delivery Queue records the pull request in a Pull Request Record Commit

## Overview

The Delivery Queue archives a Spec before it opens the pull request, so its
archive commit cannot carry the number. This Task adds the Pull Request
Record Commit. At the start of the Checking stage, the queue stores the
opened pull request's number in the record through
`roundfix archive <slug> --pull-request <n>`. It proves that the commit
changes nothing else, commits it and returns to Publishing to push it. It is
verifiable on its own: the engine with a fake recorder, the command
workflow against real Git, and an end-to-end run against the fake forge whose
merged head carries the number.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12
   [an Archive Record points at a branch commit a squash merge drops](references/2026-10-12-an-archive-record-points-at-a-branch-commit-a-squash-merge-drops.md)
   on the delivery side, implementing `_techspec.md` → API Contract 4 and API
   Contract 5 exactly.
2. MUST add `delivery.PullRequestRecorder` and `delivery.PullRequestRecord`
   as in `_techspec.md` → Interfaces, as an `EngineDependencies` field. The
   command engine MUST wire it. An engine built without one skips the step,
   as an engine without a conflict resolver skips conflict resolution.
3. MUST run the step at the start of `checkCandidate`, before the first check
   read. On a recorded commit it appends the head to the item's candidate
   commits, writes the log line of API Contract 4, and returns to Publishing
   without the repository gate. An inexact commit, or a returned head equal
   to the candidate head, parks `review-stale`. An error parks through the
   existing error path as `delivery-error`. Nothing merges while the record
   at the candidate head lacks the number.
4. MUST make the command workflow's `RecordPullRequest` idempotent and safe
   to resume, as API Contract 4 states. It returns an empty result for a
   Legacy Archive Folder and for a record that already names the number. It
   adopts an exact child commit of the candidate head found at the checkout
   `HEAD`, and commits an uncommitted record that is the only working-tree
   change without rewriting it. Otherwise it runs the Archive Command through
   the step executable and proves the result. The commit subject MUST be
   `docs: record pull request #<n> in the <slug> Archive Record`.
5. MUST add these tests, each top-level test calling `t.Parallel()` first or
   carrying a `// Sequential: <reason>` comment, with no network connection:
   - in `internal/delivery`: `TestDeliveryRecordsThePullRequestBeforeChecks`,
     `TestDeliveryRetryAdoptsAPullRequestRecordCommit`,
     `TestDeliveryParksAnInexactPullRequestRecordCommit` and
     `TestDeliveryParksWhenThePullRequestCannotBeRecorded`;
   - in `internal/cli`: `TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord`,
     `TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath`,
     `TestDeliveryWorkflowAdoptsAnUncommittedRecord` and
     `TestCommandDeliveryEngineRecordsThePullRequest`. The last runs a queue
     item to merge against the fake forge and reads `pull_request` from the
     record at the merged head.
6. MUST keep `TestDeliveryArchiveStageCommitsTheArchiveRecord`,
   `TestDeliveryAcceptsAnExactRetirement` and
   `TestDeliveryEngineMergesAQueueWithoutAnOperator` passing. An existing test
   whose expected stages or candidate commits move because of the new commit
   MAY change only to that new contract. The Result MUST name each one and
   what moved.
7. MUST prove each new gate can fail. The Result MUST record two sabotages,
   restored afterwards: the proof accepting a commit that also changes another
   path, and the step skipped at Checking. For each it names the test that
   failed.
8. MUST NOT change the archive stage, its exact-archive proof, the Baseline
   delivery order, the pull request boundary or any blocker string. It MUST
   NOT push any ref other than the item branch.

## Subtasks

- [ ] Add the recorder dependency and the Checking-stage step.
- [ ] Implement the command workflow's record, adopt and commit paths.
- [ ] Implement the exactness proof.
- [ ] Add the engine, workflow and end-to-end tests.
- [ ] Record the sabotages and any existing test that moved.

## Acceptance Criteria

- [ ] A queue run against the fake forge merges a head whose record names the opened pull request.
- [ ] An interruption between opening the pull request and recording it produces exactly one Pull Request Record Commit after retry.
- [ ] A commit that changes anything but the record's number parks `review-stale`.
- [ ] A failed recording parks the item, and nothing merges without the number.

## Context

- instruction: `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- instruction: `docs/adr/0223-a-delivery-keeps-its-work-across-archive-requeue-and-review.md`
- instruction: `docs/adr/0237-a-delivery-retry-records-a-merge-made-outside-the-queue.md`
- instruction: `internal/cli/deliver_archive_record_test.go`
- instruction: `internal/delivery/engine_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_test.go`
- creates: `internal/delivery/record_pull_request_test.go`
- creates: `internal/cli/deliver_record_pull_request_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliveryRecordsThePullRequestBeforeChecks|TestDeliveryRetryAdoptsAPullRequestRecordCommit|TestDeliveryParksAnInexactPullRequestRecordCommit|TestDeliveryParksWhenThePullRequestCannotBeRecorded|TestDeliveryEngineMergesAQueueWithoutAnOperator)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryRecordsThePullRequestBeforeChecks TestDeliveryRetryAdoptsAPullRequestRecordCommit TestDeliveryParksAnInexactPullRequestRecordCommit TestDeliveryParksWhenThePullRequestCannotBeRecorded TestDeliveryEngineMergesAQueueWithoutAnOperator; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four engine tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord|TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath|TestDeliveryWorkflowAdoptsAnUncommittedRecord|TestCommandDeliveryEngineRecordsThePullRequest|TestDeliveryArchiveStageCommitsTheArchiveRecord|TestDeliveryAcceptsAnExactRetirement)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath TestDeliveryWorkflowAdoptsAnUncommittedRecord TestCommandDeliveryEngineRecordsThePullRequest TestDeliveryArchiveStageCommitsTheArchiveRecord TestDeliveryAcceptsAnExactRetirement; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four workflow tests do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; User Stories 1 and 3; Core Feature 2; Success Metric 1
- [_techspec.md](_techspec.md) — Interfaces; API Contract 4; API Contract 5; Testing Approach; Risks & Considerations; Build Order 4
- ADR-0261; ADR-0223; ADR-0229; ADR-0237; ADR-0232
