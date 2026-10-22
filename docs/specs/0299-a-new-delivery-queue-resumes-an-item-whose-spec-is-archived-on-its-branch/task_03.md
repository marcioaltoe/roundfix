---
task: task_03
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
status: completed
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
7. MUST keep every child process these tests start bounded. Without
   `delivery.item_binary`, `stepExecutable` returns `os.Executable()`, which
   in a test is the `internal/cli` test binary: started without
   `ROUNDFIX_CLI_TEST_HELPER=1` it runs the whole suite again instead of the
   command, and the suite starts the same test again, without end (on
   2026-10-22 an authoring characterization test of this shape left 405
   `cli.test archive 0001-widget-flow --pull-request 1` processes running
   for two hours). Each new test either reaches no step that starts
   `roundfix` (assert it) or starts it only as the CLI helper the way the
   existing deliver tests do (`cliTestHelperEnv` in the child's
   environment, never through `t.Setenv` in a parallel test). No child
   outlives its test, and the Verification bounds the run with
   `-timeout 10m`.

## Subtasks

- [ ] Implement the archived-item proof and wire it into the queue.
- [ ] Resolve archive paths without the Specs Root directory.
- [ ] Add the deliver start tests for the resume, the three refusals and the unarchived control.
- [ ] Run the unchanged archive, record, publication and retry tests.
- [ ] Keep every child the new tests start bounded: no step starts `roundfix`, or it runs as the CLI helper.

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

- `out="$(go test -count=1 -timeout 10m -v -run '^(TestDeliverStartResumesAnArchivedItemBranchAtReview|TestDeliverStartParksAnUnprovenArchivedItemBranch|TestDeliverStartRevalidatesAContinuedBranchWithoutAnArchive|TestResumeAcceptsARealArchiveCommit|TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord|TestInspectItemReportsAnArchivedSpec|TestDeliverStartContinuesTheItemBranchThatHoldsWork|TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview|TestDeliveryPublicationRequestsNoCodeRabbitReview|TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverStartResumesAnArchivedItemBranchAtReview/another_active_Spec TestDeliverStartResumesAnArchivedItemBranchAtReview/only_active_Spec TestDeliverStartParksAnUnprovenArchivedItemBranch/inexact_archive_commit TestDeliverStartParksAnUnprovenArchivedItemBranch/folder_removed_without_archive TestDeliverStartParksAnUnprovenArchivedItemBranch/no_removing_commit TestDeliverStartRevalidatesAContinuedBranchWithoutAnArchive TestResumeAcceptsARealArchiveCommit TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord TestInspectItemReportsAnArchivedSpec TestDeliverStartContinuesTheItemBranchThatHoldsWork TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview TestDeliveryPublicationRequestsNoCodeRabbitReview TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check. Afterwards the resume, refusal and control tests pass through `roundfix deliver start` with the unchanged archive, record, recovery, continuation, retry and publication tests, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 5; Acceptance evidence
- `_techspec.md` → Current behavior; API Contracts 1, 5 and 6; Invariants 1-5; Testing Approach; Build Order 3
- ADR-0223
- ADR-0312

## Result

Implemented the Task 03 slice for Daemon Verification:

- `commandDeliveryWorkflow.ProveArchivedItem` reads the recorded item
  worktree's Git objects, returns the zero proof for an active PRD, emits
  each specified refusal, and proves the removing commit against its full
  first-parent hash with the existing `archiveCommitIsExact`. It reads no
  commit subject and performs no writes. Its comment records that refusal
  paths and commit prefixes reach the Run Database blocker and status,
  without environment values or file contents.
- The command engine wires `ArchivedItems: workflow`, with the interface
  assertion. `archivePaths` now uses `resolveHistoricalSpecsRoot`.
  Byte comparisons against HEAD confirmed `archiveCommitIsExact`,
  `archiveRetirementIsExact`, `reconcileArchiveCommit`, `InspectItem` and
  the `deliver_revalidate.go` file are unchanged.
- Added only `internal/cli/deliver_archived_resume_test.go`; no existing
  test, golden or fixture was changed. Every top-level test calls
  `t.Parallel()` first. Disposable homes, a bare origin, a hardened clone,
  real `spec.Archive` and `spec.RecordArchivePullRequest` supply the Git
  history. Only readiness, owner launch and implementation/review/forge
  boundaries are replaced. The owner uses real workspace, recovery,
  history, revalidator, archiver, recorder and archive prover behavior.

Acceptance evidence from focused checks:

1. Both `TestDeliverStartResumesAnArchivedItemBranchAtReview` subtests
   (`another_active_Spec`, `only_active_Spec`) passed through `deliver
   start`: merged with existing PR #619, zero Implement calls, exactly one
   review at the record commit head, zero PR creations, and unchanged
   candidate/item head. The head and absent Specs Root are observed before
   the real workspace's normal merged-branch cleanup.
2. All three `TestDeliverStartParksAnUnprovenArchivedItemBranch` subtests
   passed: the persisted blocker equals `archived-item-unproven: ` plus
   the exact cause, with zero Implement/review calls, and `deliver status`
   prints the specified `finding` Park Class and next action.
3. `TestDeliverStartRevalidatesAContinuedBranchWithoutAnArchive` passed
   with real revalidation and one Implement call, then the fake runner
   deliberately returns unresolved before any archive/publication step.

Commands and outcomes:

- `go test -count=1 -timeout 3m -v -run '^TestDeliverStart(ResumesAnArchivedItemBranchAtReview|ParksAnUnprovenArchivedItemBranch|RevalidatesAContinuedBranchWithoutAnArchive)$' ./internal/cli`
  — exit 0; all three tests and five named subtests passed.
- `go test -count=1 -timeout 4m -run '^(TestResumeAcceptsARealArchiveCommit|TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord|TestInspectItemReportsAnArchivedSpec|TestDeliverStartContinuesTheItemBranchThatHoldsWork|TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview|TestDeliveryPublicationRequestsNoCodeRabbitReview|TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug)$' ./internal/cli`
  — exit 0; unchanged consumer tests passed.
- `GOOS=windows go test -c -o /tmp/roundfix-task03-cli.test.exe ./internal/cli`
  — exit 0 for the final test source; Windows compilation checked, not
  Windows execution.
- A temporary Go overlay restoring strict Specs Root resolution ran only
  `TestDeliverStartResumesAnArchivedItemBranchAtReview/only_active_Spec`:
  expected exit 1 with `delivery-error` and the missing Specs Root cause.
  This confirms the regression test detects the old behavior without
  changing repository files.
- `git diff --check` — exit 0.

Child-process bound: tests use the test context for Git processes and set
`cliTestHelperEnv=1` through per-test command environment overrides, never
`t.Setenv`. Any regression that reaches a roundfix step therefore dispatches
as the CLI helper instead of recursively running the suite. The resume
assertions prove no new commit or Implement Run; refusals stop before Run
or review. No real-time wait is introduced.

Environment: the sandbox could not initialize its symlinked writable root;
commands used approved escalation in the resolved Run checkout. No project
configuration changed. Task status, the declared Verification command and
repository Verification remain Daemon-owned and were not run or settled by
this turn. No commit, push or PR was made.

## Carry-forward provenance

- Source Run: `run_20261022T213854Z_c379322ae483e8bf`
- Source commit: `69680008faaf0de2fa27cf8744f6a1a2b596298f`
