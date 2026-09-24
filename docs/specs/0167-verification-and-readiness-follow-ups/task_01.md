---
task: task_01
spec: 0167-verification-and-readiness-follow-ups
status: pending
type: backend
complexity: medium
---

# Task 01: The retry's verdict wins; completed repairs are history

## Overview

In independent mode a command failing in both the first run and the exclusive retry reaches the repair turn twice, and one passing on retry is still carried as failed; and `ValidatePreconditionRepairs` inspects completed repair Tasks, so a completed repair whose Task file lost the verbatim command blocks the Spec's next `implement`.

## Requirements

1. MUST let the exclusive retry's verdict replace the first run's for every command the retry reached, carrying first-run failures only for commands it never reached.
2. MUST list each current failure exactly once in the repair request and the reason.
3. MUST skip Tasks whose status is completed in the verbatim planning check, while still refusing a pending named Task without the command.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A command failing in both runs appears once; a command passing on retry is absent.
- [ ] A completed repair Task without the verbatim command does not refuse planning; a pending one does.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/task_engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRetryVerdictReplacesTheFirstRun|TestAFailureRepeatedOnRetryReachesRepairOnce|TestCompletedRepairTaskDoesNotBlockPlanning|TestPendingRepairTaskStillNeedsTheVerbatimCommand)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRetryVerdictReplacesTheFirstRun TestAFailureRepeatedOnRetryReachesRepairOnce TestCompletedRepairTaskDoesNotBlockPlanning TestPendingRepairTaskStillNeedsTheVerbatimCommand; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Retry verdicts
