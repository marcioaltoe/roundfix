---
task: task_03
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
status: pending
type: backend
complexity: medium
---

# Task 03: The queue proves an archive on a continued item branch and delivers it through deliver start

## Overview

Task 02 gave the Delivery Engine a prover of an archived item branch; the
command workflow does not implement it, and `archivePaths` fails in an item
worktree whose only Spec was archived. This Task implements the proof with the
queue's existing archive-commit proof, wires it into the queue, lets
`archivePaths` resolve an absent Specs Root, and proves the whole path through
`roundfix deliver start` in temporary repositories: the shape that parked
Spec 0292's item `delivery-error` on 2026-10-21 now reaches `merged`, and an
unproven branch parks with its cause.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-22,
   [`references/2026-10-22-a-new-delivery-queue-parks-an-item-whose-spec-is-archived-on-its-branch.md`](references/2026-10-22-a-new-delivery-queue-parks-an-item-whose-spec-is-archived-on-its-branch.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1, 5 and 6 of
   `_techspec.md` in `internal/cli/deliver_workflow.go`.
2. MUST add `commandDeliveryWorkflow.ProveArchivedItem` with the four
   outcomes and the three refusal texts of API Contract 5, judging the commit
   that removed `<source>/_prd.md` with the existing `archiveCommitIsExact`
   against the full hash of its first parent, reading no commit subject and
   writing nothing; add `var _ delivery.ArchivedItemProver =
   (*commandDeliveryWorkflow)(nil)` and wire `ArchivedItems: workflow` in
   `newCommandDeliveryEngine`.
3. MUST make `archivePaths` resolve the Specs Root with
   `resolveHistoricalSpecsRoot` (API Contract 6), keeping its result for a
   present root; `archiveCommitIsExact`, `archiveRetirementIsExact`,
   `reconcileArchiveCommit`, `Revalidate` and `InspectItem` keep their bytes.
4. MUST state the hazard in the code's behavior: the proof reads only Git
   objects of the item worktree the queue recorded, and its refusal text,
   which names repository paths and commit prefixes, reaches the Run
   Database blocker and `roundfix deliver status`; no environment value or
   file content enters it.
5. MUST add, in the new file `internal/cli/deliver_archived_resume_test.go`,
   the three tests `_techspec.md` → Testing Approach lists for that file,
   each top-level test calling `t.Parallel()` first, compiling for Windows,
   using a temporary home, a bare origin and a clone built like
   `newItemBranchStartRepository`, the Archive Command's real archive
   (`spec.Archive`) and `spec.RecordArchivePullRequest` for the Pull Request
   Record Commit, replacing `deliveryReadiness` and `startDeliveryOwner`, and
   waiting on no real time:
   - `TestDeliverStartResumesAnArchivedItemBranchAtReview` with subtests
     `another_active_Spec` and `only_active_Spec`; the fake forge holds open
     Pull Request #619 for the item branch at the record commit's head, finds
     it rather than creating one, and merges it; the test asserts stage
     `merged`, Pull Request number `619`, no Implement call, exactly one
     review whose head is the record commit's head, no Pull Request created
     and an unchanged item head. The `only_active_Spec` subtest seeds no other
     Spec, so the item worktree has no `docs/specs` directory, which a strict
     `archivePaths` fails.
   - `TestDeliverStartParksAnUnprovenArchivedItemBranch` with subtests
     `inexact_archive_commit` (the archive commit also adds `unrelated.txt`),
     `folder_removed_without_archive` (the Spec folder removed with `git rm`
     and no archive) and `no_removing_commit` (an orphan item branch whose
     only commit holds the archived tree); each asserts the blocker
     `archived-item-unproven: ` followed by its exact cause, no Implement
     call and no review, and `roundfix deliver status` printing `Park:
     <slug> finding: ` followed by the next action of API Contract 3.
   - `TestDeliverStartRevalidatesAContinuedBranchWithoutAnArchive`: a
     continued branch whose head still holds the Spec folder makes one
     Implement call through the real `Revalidate`.
6. MUST run the existing readers of `archivePaths` and of an archived item by
   name, unchanged: `TestResumeAcceptsARealArchiveCommit`,
   `TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord`,
   `TestInspectItemReportsAnArchivedSpec`,
   `TestDeliverStartContinuesTheItemBranchThatHoldsWork`,
   `TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview`,
   `TestDeliveryPublicationRequestsNoCodeRabbitReview` and
   `TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug`. No existing
   test, golden or fixture changes.

## Subtasks

- [ ] Implement the archived-item proof and wire it into the queue.
- [ ] Resolve archive paths without the Specs Root directory.
- [ ] Add the deliver start tests for the resume, the three refusals and the unarchived control.
- [ ] Run the unchanged archive, record, publication and retry tests.

## Acceptance Criteria

- [ ] Through `roundfix deliver start`, a continued branch with an exact archive commit and a Pull Request Record Commit reaches `merged` with its open Pull Request, no Implement call, one review of its head and no new commit, with and without another active Spec.
- [ ] An inexact archive commit, a folder removed without an archive and a branch without a removing commit each park with their cause, and status prints the `finding` Park Class with the next action.
- [ ] A continued branch that still holds the Spec folder is revalidated and runs as before.

## Consumers

- `ProveArchivedItem`, its wiring in `newCommandDeliveryEngine` and `archivePaths` — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0312-a-new-delivery-queue-resumes-an-archived-item-branch-at-review.md`
- instruction: `docs/adr/0223-a-delivery-keeps-its-work-across-archive-requeue-and-review.md`
- instruction: `internal/delivery/engine.go`
- instruction: `internal/cli/archive_reader.go`
- instruction: `internal/cli/deliver_revalidate.go`
- instruction: `internal/cli/deliver_item_branch_test.go`
- instruction: `internal/cli/deliver_test.go`
- instruction: `internal/cli/deliver_record_pull_request_test.go`
- instruction: `internal/cli/deliver_recovery_test.go`
- instruction: `internal/cli/deliver_archived_retry_test.go`
- instruction: `internal/cli/deliver_publication_scope_test.go`
- instruction: `internal/cli/deliver_queue_subject_test.go`
- instruction: `internal/cli/implement_test.go`
- instruction: `internal/cli/archive_test.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_archived_resume_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliverStartResumesAnArchivedItemBranchAtReview|TestDeliverStartParksAnUnprovenArchivedItemBranch|TestDeliverStartRevalidatesAContinuedBranchWithoutAnArchive|TestResumeAcceptsARealArchiveCommit|TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord|TestInspectItemReportsAnArchivedSpec|TestDeliverStartContinuesTheItemBranchThatHoldsWork|TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview|TestDeliveryPublicationRequestsNoCodeRabbitReview|TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverStartResumesAnArchivedItemBranchAtReview/another_active_Spec TestDeliverStartResumesAnArchivedItemBranchAtReview/only_active_Spec TestDeliverStartParksAnUnprovenArchivedItemBranch/inexact_archive_commit TestDeliverStartParksAnUnprovenArchivedItemBranch/folder_removed_without_archive TestDeliverStartParksAnUnprovenArchivedItemBranch/no_removing_commit TestDeliverStartRevalidatesAContinuedBranchWithoutAnArchive TestResumeAcceptsARealArchiveCommit TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord TestInspectItemReportsAnArchivedSpec TestDeliverStartContinuesTheItemBranchThatHoldsWork TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview TestDeliveryPublicationRequestsNoCodeRabbitReview TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check. Afterwards the resume, refusal and control tests pass through `roundfix deliver start` with the unchanged archive, record, recovery, continuation, retry and publication tests, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 5; Acceptance evidence
- `_techspec.md` → Current behavior; API Contracts 1, 5 and 6; Invariants 1-5; Testing Approach; Build Order 3
- ADR-0223
- ADR-0312
