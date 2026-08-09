---
task: task_03
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
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
   update it in the same commit.

## Subtasks

- [ ] Re-read each issue before marking.
- [ ] Preserve Terminal outcomes.
- [ ] Leave the non-Terminal path unchanged.

## Acceptance Criteria

- [ ] An issue at `resolved`, `invalid` or `duplicated` survives a failed Batch.
- [ ] An issue at `pending` or `valid` is still marked `failed` with the reason.
- [ ] The in-memory Batch status is never trusted.

## Bounded scope

This Task may create or modify only:

- `internal/agent/agent.go`
- `internal/agent/agent_test.go`
- `internal/daemon/run_disposition_characterization_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_03.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestMarkBatchFailedKeeps' -count=1 -v 2>&1 | grep -q '^--- PASS: TestMarkBatchFailedKeepsAlreadySettledIssues/keeps_resolved_issue_untouched'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestMarkBatchFailedKeeps' -count=1 -v 2>&1 | grep -q '^--- PASS: TestMarkBatchFailedKeepsAlreadySettledIssues/marks_pending_issue_failed'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -count=1 2>&1 | grep -q '^ok'` — expected: exits 0, proving the declared break was updated.

## References

- `_prd.md` → Goal 2.
- `_techspec.md` → Build Order 3.
- ADR-0113.
- `docs/backlog/2026-08-09-a-failed-batch-erases-the-issues-it-resolved.md`
