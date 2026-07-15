---
task: task_05
spec: 0028-settlement-and-reporting
status: pending
type: backend
complexity: medium
---

# Task 05: Return per-Task outcomes and report failure reasons

## Overview

Make the implement report diagnostic on its own: the Daemon's cycle result carries each Task's outcome (status plus the same one-line reason it journals), and the implement report renders an indented `  reason: <one line>` after every failed and skipped Task line. The existing per-Task line shape stays byte-stable, so current parsers keep working.

## Requirements

1. MUST extend the Daemon's task-cycle result with per-Task outcomes — task id, final status, and a one-line reason — populated at every settle and skip point with the same reason strings already published to the Run Event Journal.
2. MUST render one indented `  reason: <one line>` line after each failed and skipped Task line in the implement report; completed Tasks gain no extra line and the existing `task_NN <status> — <title>` line shape is unchanged.
3. MUST name the failed command and exit status, with a pointer to the diagnostics, when the failure came from Verification.
4. MUST source reasons from the in-memory result, not a journal read-back at report time (TechSpec decision).

## Subtasks

- [ ] Accumulate per-Task outcomes in the engine at settle and skip points
- [ ] Extend the cycle result and thread it to the implement report renderer
- [ ] Render the indented reason lines for failed and skipped Tasks
- [ ] Tests: engine test asserting outcomes match journaled reasons for failed, skipped, and completed Tasks; report rendering fixtures with reasons present and absent

## Acceptance Criteria

- [ ] A Run with a failed Task prints the failed line followed by an indented reason naming the failed step (command and exit status for Verification failures)
- [ ] A skipped Task's reason names the unmet needs, matching the journal payload
- [ ] Completed Task lines are byte-identical to today
- [ ] The full test suite passes

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/implement.go`

## Verification

- `grep -q "TaskOutcome" internal/daemon/task_engine.go` — expected: exit 0
- `go test ./internal/daemon/... ./internal/cli/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 3, User Story 5, Core Feature 5; `_techspec.md` → Build Order 5, Interfaces (TaskOutcome, TaskCycleResult), API Contracts (implement report), Decisions (in-memory result, additive line).
