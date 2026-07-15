---
task: task_03
spec: 0029-launch-and-recovery-fixes
status: pending
type: backend
complexity: medium
---

# Task 03: Settle Batches with the model-rejection reason instead of protocol error

## Overview

Wire the not-advertised classification through the Daemon: a Batch that dies because the Agent Session rejected the Agent Model settles its Work Items with a terminal reason naming the rejected model and the advertised list, journals the same, and the final report shows the actionable message — replacing the generic `Agent Batch failed after acpx exited with code 1: agent/protocol error` observed in the 2026-07-14 field failure.

## Requirements

1. MUST detect the typed model-rejection error on Batch failure paths (review Batches and spec Task Batches) and use its rendered message as the terminal reason for the affected Work Items, shaped like `Agent Model "<model>" not advertised by runtime "<runtime>"; advertised: <list>`.
2. MUST carry the same message into the Run Event Journal and the end-of-run report's reason lines (the report plumbing from specs 0027/0028 is already in place — this task supplies the better reason string).
3. MUST leave every other Batch failure classification and reason unchanged.
4. MUST keep Verification and settlement mechanics untouched — only the reason content changes.

## Subtasks

- [ ] Detect the typed error at the Batch-failure settle points for both review and spec engines
- [ ] Thread the rendered reason into issue/task terminal reasons and Run Events
- [ ] Engine tests: a fake runner failing with the rejection stderr settles Work Items with the model-rejection reason; an unrelated failure keeps its current reason
- [ ] Report-level test: the final report's reason line shows the model-rejection message

## Acceptance Criteria

- [ ] A spec Task whose Batch fails on the rejection carries the model-rejection terminal reason in its journal payload and report line
- [ ] A review Batch failing the same way settles its Review Issues with that reason in their artifacts
- [ ] Unrelated agent failures produce byte-identical reasons to today
- [ ] The full test suite passes

## Context

- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/agent/agent.go`

## Verification

- `grep -rq "not advertised" internal/daemon` — expected: exit 0 (reason wiring exists)
- `go test ./internal/daemon/... ./internal/agent/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 3, Core Feature 2, Problem 2; `_techspec.md` → Build Order 3, API Contracts (Batch failure reason), Coverage Map (Goal 3).
