---
task: task_04
spec: 0156-a-delivery-loop-that-outlives-the-session
status: pending
type: backend
complexity: medium
---

# Task 04: The command family

## Overview

Expose the queue as `roundfix deliver start <slug>...`, `status`, `resume` and `stop`, starting the owner detached with the mechanism `implement --detach` already uses.

## Requirements

1. MUST record the queue and start a detached owner on `deliver start`.
2. MUST print each item's stage and blocker on `deliver status`.
3. MUST end the owner on `deliver stop` and restart it from the persisted queue on `deliver resume`.
4. MUST refuse unknown flags and unknown Spec slugs, matching the surrounding commands.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The command family appears on the public surface.
- [ ] `deliver status` prints each item's stage and blocker.
- [ ] Unknown slugs are refused before any queue is recorded.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/detach.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestDeliverCommandRecordsAQueueAndReportsIt$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestDeliverCommandRecordsAQueueAndReportsIt"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `go run -buildvcs=false ./cmd/roundfix --help 2>&1 | grep -q "roundfix deliver"` — expected: exit 0; the command family appears on the public surface. Before this Task the usage has no deliver line, so the command fails.

## References

- [_techspec.md](_techspec.md) — Components
