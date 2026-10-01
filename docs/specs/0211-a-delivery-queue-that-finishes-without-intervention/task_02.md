---
task: task_02
spec: 0211-a-delivery-queue-that-finishes-without-intervention
status: pending
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
