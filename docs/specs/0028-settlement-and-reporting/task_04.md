---
task: task_04
spec: 0028-settlement-and-reporting
status: pending
type: backend
complexity: low
---

# Task 04: Warn on no-op Task commits

## Overview

Make silent no-op completions visible while the Run is still observable: when a Task settles completed and its Task commit contains no change outside the Spec Root — including the currently-silent case where nothing is stageable at all — the Daemon still settles the Task but publishes a warning Run Event and one stderr warning line. Warn, never block (PRD decision).

## Requirements

1. MUST classify the Task commit's stageable paths against the Spec Root at the engine's commit step and detect the two no-op shapes: all paths inside the Spec Root, and an empty stageable set.
2. MUST publish a warning Run Event carrying the task id and the no-op shape, and emit one stderr warning line, in both shapes.
3. MUST NOT change settlement or commit behavior — the Task still settles completed and the commit (when non-empty) is still created.
4. MUST keep the existing per-path dropped-stage reporting unchanged.

## Subtasks

- [ ] Path classification against the Spec Root at the commit step
- [ ] Warning Run Event + stderr line for both no-op shapes
- [ ] Engine tests: a spec-tree-only Task commit warns and still commits; an empty stageable set warns and still settles; a Task with code changes does not warn

## Acceptance Criteria

- [ ] A Task whose only change is its own task file settles completed with the warning event and stderr line
- [ ] A Task with an empty stageable set settles completed with the warning event and stderr line
- [ ] A Task with changes outside the Spec Root produces no warning
- [ ] The full test suite passes

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/runevent/event.go`

## Verification

- `go test ./internal/daemon/...` — expected: all tests pass, including the no-op warning coverage
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 4, User Story 3, Core Feature 3, Decisions (warn, don't block); `_techspec.md` → Build Order 4, System Architecture (internal/daemon), API Contracts (no-op Task commit).
