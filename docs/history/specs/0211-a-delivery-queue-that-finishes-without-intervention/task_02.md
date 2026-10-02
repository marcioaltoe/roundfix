---
task: task_02
spec: 0211-a-delivery-queue-that-finishes-without-intervention
status: completed
type: backend
complexity: medium
---

# Task 02: An archived Delivery Retry finds the start head of a Run the queue started

## Overview

The Delivery Retry of an operator-archived `qa-environment-partial` item reads
the Implement start head of the item's Run, but the lookup compares the Run's
working root with the queue's checkout root. A Run the queue owner starts runs
in the item worktree, so the lookup always refuses it. This Task matches the
Run by its repository, so the printed answer of that park, carry-forward,
operator evidence, `archive --qa-override` and `deliver retry`, returns the
item to `reviewing`.

## Requirements

1. MUST make `RunStart` accept an Implement Run with a non-empty start head
   whose repository root equals `roundconfig.RepositoryRoot` of the queue's
   checkout, as "The Run start lookup" states, deriving the Run's repository
   root from its working root when the stored value is empty.
2. MUST keep the refusal text `Run %q has no Implement start head for
   repository %q` for an absent Run, a non-Implement Run, an empty start head
   and a Run of another repository.
3. MUST NOT change `Engine.Retry`, the descent proof, the QA Archive Override
   rules or any other Delivery Retry branch.
4. MUST add to the new file `internal/cli/deliver_archived_retry_test.go`:
   `TestRunStartFindsAQueueStartedRunByItsRepository`, whose Run has a linked
   item worktree as its working root;
   `TestRunStartRefusesARunOfAnotherRepository`; and
   `TestArchivedRetryOfAQueueStartedRunReturnsToReview`, which drives
   `Engine.Retry` of a real `delivery.NewEngine` whose Workspace, Recovery and
   History are the command workflow, with a fake Pull Request boundary and a
   disposable Run Database and repository, from an item parked
   `qa-environment-partial` whose Spec the in-process archive command archived
   with `--qa-override`, and asserts stage `reviewing` and the item head
   appended to the candidate commits. Each test MUST fail against the current
   lookup.
5. MUST leave `TestOperatorArchiveHistoryReadsTheRunStartAndAncestry` and the
   operator-archive tests in `internal/delivery` unedited and passing.

## Subtasks

- [ ] Match the Run by repository in the start-head lookup.
- [ ] Add the two lookup tests over a linked item worktree.
- [ ] Add the end-to-end archived retry test.

## Acceptance Criteria

- [ ] A Run whose working root is a linked item worktree of the checkout
      yields its start head; before this Task the same fixture refuses with
      `has no Implement start head`.
- [ ] A Run of another repository is refused with the unchanged message.
- [ ] The archived retry of a queue-started Run returns the item to
      `reviewing` with the descended head as its newest candidate commit.
- [ ] The existing operator-archive tests pass unedited.

## Context

- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_archived_retry_test.go`
- instruction: `internal/cli/deliver_operator_archive_test.go`
- instruction: `internal/cli/deliver_recovery_test.go`
- instruction: `internal/delivery/operator_archive_retry_test.go`
- instruction: `internal/delivery/engine.go`
- instruction: `internal/store/store.go`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `.agents/skills/roundfix/references/deliver.md`
- instruction: `docs/specs/0211-a-delivery-queue-that-finishes-without-intervention/references/2026-10-01-an-archived-retry-cannot-find-the-runs-implement-start-head.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRunStartFindsAQueueStartedRunByItsRepository|TestRunStartRefusesARunOfAnotherRepository|TestArchivedRetryOfAQueueStartedRunReturnsToReview)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunStartFindsAQueueStartedRunByItsRepository TestRunStartRefusesARunOfAnotherRepository TestArchivedRetryOfAQueueStartedRunReturnsToReview; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestOperatorArchiveHistoryReadsTheRunStartAndAncestry|TestRunStartFindsAQueueStartedRunByItsRepository)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestOperatorArchiveHistoryReadsTheRunStartAndAncestry TestRunStartFindsAQueueStartedRunByItsRepository; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing start-head test runs unedited beside the new one, which does not exist before this Task.
- `out="$(go test -count=1 -v -run "^(TestRunStartFindsAQueueStartedRunByItsRepository|TestRunStartRefusesARunOfAnotherRepository|TestArchivedRetryOfAQueueStartedRunReturnsToReview)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunStartFindsAQueueStartedRunByItsRepository TestRunStartRefusesARunOfAnotherRepository TestArchivedRetryOfAQueueStartedRunReturnsToReview; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run "^(TestRetryResumesAnOperatorArchivedItemAtReview|TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor|TestOperatorArchiveRetryRequiresHistory)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRetryResumesAnOperatorArchivedItemAtReview TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor TestOperatorArchiveRetryRequiresHistory; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the delivery operator-archive tests run unedited after the new tests, which do not exist before this Task.

## References

- `_prd.md` → User Story 1; Core Feature 1; Success Metric 1; Acceptance evidence
- `_techspec.md` → The Run start lookup; Testing Approach 1; Build Order 2
- ADR-0154; ADR-0170

## Result

Implemented the Task 02 slice for Daemon Verification. `RunStart` now compares
`run.RepositoryRoot` with `roundconfig.RepositoryRoot` of the queue checkout.
For legacy rows with no stored repository root, it derives identity from the
Run's working root. The existing refusal text remains for absent Runs,
non-Implement Runs, empty start heads and another repository. `Engine.Retry`,
the ancestry proof and QA Archive Override rules are unchanged.

Added the three required tests in `internal/cli/deliver_archived_retry_test.go`.
All three failed against the original lookup before the production edit:
linked item Runs were refused with `has no Implement start head`, and archived
retry refused the moved head against an empty candidate. The other-repository
test first asserts the exact refusal, then checks a same-repository linked
Run as a control; that control failed against the original lookup.

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Linked item worktree yields its start head | `TestRunStartFindsAQueueStartedRunByItsRepository` passes for both a stored repository root and an empty legacy root in a disposable Run Database. |
| Another repository retains the refusal | `TestRunStartRefusesARunOfAnotherRepository` passes over two repositories with linked worktrees and asserts the exact unchanged message and an empty returned head. |
| Archived retry resumes reviewing with the descended head appended | `TestArchivedRetryOfAQueueStartedRunReturnsToReview` seeds a parked `qa-environment-partial` item with no candidate, commits operator evidence, invokes the in-process archive command with `--qa-override`, and drives a real `delivery.NewEngine` with command-workflow Workspace, Recovery and History plus a fake Pull Request boundary. It asserts the returned and persisted `reviewing` stage, candidate list containing the archived item head, cleared blocker and unchanged Git head. |
| Existing operator-archive tests pass unedited | The focused CLI check includes `TestOperatorArchiveHistoryReadsTheRunStartAndAncestry`; the focused delivery check includes the operator-archive retry tests. Neither existing test file was edited. |

Checks run:

- `rtk proxy go test -count=1 ./internal/cli -run 'Test(RunStart|ArchivedRetry)'`
  before the lookup edit: exit 1; all three required new tests failed at the
  intended lookup/retry boundary.
- `rtk proxy go test -count=1 ./internal/cli -run 'Test(RunStart|ArchivedRetry|OperatorArchiveHistory)'`
  after the lookup edit: exit 0.
- `rtk proxy go test -count=1 ./internal/delivery -run 'Test.*(OperatorArchive|ArchivedHead|EnvironmentOnlyPartial)'`:
  exit 0.
- `rtk make verify-incremental`: the sandboxed attempt exited 2 because three
  existing process-lifecycle tests could not enumerate or signal their fixture
  process groups. The permission-adjusted rerun exited 0, including Go vet,
  the full Go suite, skill checks and build. No assertions or configuration
  were changed to accommodate the sandbox.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exit 0.

The Task status was already `in_progress` on entry and remains Daemon-owned.
The authored Verification commands were not run. No other Task or Task Graph
was edited; no commit, push or Pull Request was made. No follow-up was needed
within this slice.

## Carry-forward provenance

- Source Run: `run_20261002T081705Z_a516b29ed11eed22`
- Source commit: `89953dc64cd707e59ee4bae6b92a927394142294`
