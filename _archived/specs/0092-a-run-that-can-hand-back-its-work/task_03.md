---
task: task_03
spec: 0092-a-run-that-can-hand-back-its-work
status: completed
type: backend
complexity: medium
---

# Task 03: Let a failed Batch keep what its Agent achieved

## Overview

`MarkBatchFailed` overwrites every Review Issue in its Batch, including the ones
the Agent resolved. Its sibling `SettleAssignedIssues` already re-reads each issue
and skips settled ones, which is why the asymmetry reads as an oversight. It is
not: preserving outcomes alone changes what a Run reports, so this Task preserves
them and Task 04 fixes the outcome. Applied alone, this Task leaves a Run able to
report `Clean` on a crashed Agent — which is why the two are ordered.

## Requirements

1. MUST re-read each Review Issue from disk before marking it, because the
   in-memory Batch carries the status the issues held when it was assembled.
2. MUST preserve any issue already in a Terminal status; a Batch failing on one
   issue is not evidence against the others.
3. MUST mark every non-Terminal issue `failed` with the Batch's terminal reason,
   exactly as today.
4. MUST NOT change the Run outcome in this Task; the six outcome-contract tests
   Task 01 enumerated stay as they are until Task 04.
5. MUST break the characterization case Task 01 declared for the overwrite, and
   update it in the same commit. Task 01 recorded it as
   `TestRunDispositionCharacterizationFailedBatchOverwritesSettledIssues`, whose
   name states the behaviour this Task removes; the updated case is
   `TestRunDispositionCharacterizationFailedBatchKeepsSettledIssues` and asserts
   the settled issue keeps its status and reason. Renaming is what makes the
   update visible: a case that keeps its old name cannot be told apart from one
   that was never touched.

## Subtasks

- [x] Re-read each issue before marking.
- [x] Preserve Terminal outcomes.
- [x] Leave the non-Terminal path unchanged.

## Acceptance Criteria

- [x] An issue at `resolved`, `invalid` or `duplicated` survives a failed Batch.
- [x] An issue at `pending` or `valid` is still marked `failed` with the reason.
- [x] The in-memory Batch status is never trusted.

## Bounded scope

This Task may create or modify only:

- `internal/agent/agent.go`
- `internal/agent/agent_test.go`
- `internal/daemon/run_disposition_characterization_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_03.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestMarkBatchFailedKeeps' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^[[:space:]]*--- PASS: TestMarkBatchFailedKeepsAlreadySettledIssues/keeps_resolved_issue_untouched'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestMarkBatchFailedKeeps' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^[[:space:]]*--- PASS: TestMarkBatchFailedKeepsAlreadySettledIssues/marks_pending_issue_failed'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunDispositionCharacterizationFailedBatch' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunDispositionCharacterizationFailedBatchKeepsSettledIssues'` — expected: exits 0, proving the declared break was updated to the new behaviour rather than left failing or deleted. Naming the post-update case is what makes this command able to fail: the pre-work name asserts today's overwrite and passes against the unchanged tree, so asserting it would approve the Task before any work happened.

## References

- `_prd.md` → Goal 2.
- `_techspec.md` → Build Order 3.
- ADR-0113.
- `docs/backlog/2026-08-09-a-failed-batch-erases-the-issues-it-resolved.md`

## Result

`MarkBatchFailed` now parses each assigned Review Issue from disk, preserves an
on-disk Terminal status, and applies the existing `failed` write with the Batch
reason to every non-Terminal status. The regression table deliberately gives
the in-memory Batch the opposite status from disk, so its result proves the
persisted issue is the source of truth. The characterization case now has the
required `KeepsSettledIssues` name and preserves both an `invalid` status and
its recorded reason.

Focused evidence:

- Before the production edit,
  `GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run '^TestMarkBatchFailedKeepsAlreadySettledIssues$' -count=1`
  exited 1: the `resolved`, `invalid`, and `duplicated` cases were overwritten
  as `failed`, while the two non-Terminal cases passed.
- After the final code and test edits,
  `GOCACHE="$PWD/.gocache" rtk go test ./internal/agent ./internal/daemon -run '^(TestSettleAssignedIssues|TestMarkBatchFailedKeepsAlreadySettledIssues|TestRunDispositionCharacterizationFailedBatchKeepsSettledIssues)$' -count=1`
  exited 0 with 14 tests passing across both packages.
- `rtk git diff --name-only` listed only this Task's four bounded paths. The six
  Task 04 outcome-contract tests remain outside the diff.

Acceptance evidence:

- Terminal survival: named table cases preserve `resolved`, `invalid`, and
  `duplicated`; the latter two also preserve their reasons, and `duplicated`
  preserves its canonical issue path.
- Non-Terminal failure: named table cases change `pending` and `valid` to
  `failed` with the Batch reason; an already-`failed` issue also receives the
  current Batch reason.
- Disk authority: every table row seeds a conflicting in-memory status, and all
  14 focused test executions pass only when `MarkBatchFailed` follows disk.

The Agent did not run any declared `## Verification` command; the Daemon owns
those commands and Task settlement.

Verification Feedback attempt 1 showed that the first declared command's Go
test run passed the parent and every named subtest, but its column-one `grep`
could not match Go's indented subtest result line. The first two matchers now
accept standard leading whitespace while retaining the exact `PASS` marker and
exact test/subtest names. A focused `awk` inspection of the Daemon diagnostic
confirmed the original expression had zero matches and the whitespace-aware
expression had one; production code and tests did not change during this
repair.
