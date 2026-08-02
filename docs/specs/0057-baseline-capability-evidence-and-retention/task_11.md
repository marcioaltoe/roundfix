---
task: task_11
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: medium
---

# Task 11: Report the result as a status matrix

## Overview

A successful apply currently reads as an update being complete, when what was
proven is that postimages were written. Retention, alignment, repository
Verification, and idempotence are separate facts and some of them may not have
run at all. This Task reports five axes separately and reserves completion
language for the case that earns it.

## Requirements

1. MUST report the final result as five separate axes: approved postimages,
   semantic retention, profile alignment, repository Verification, and
   idempotence.
2. MUST report each axis as verified or not run, so an axis that never
   executed is never read as passing.
3. MUST use completion language only when semantic retention is verified and
   the idempotence check passed.
4. MUST derive each axis from the evidence the run actually produced, not from
   the absence of an error.
5. MUST carry the same five axes in machine output, additively.
6. MUST leave the transaction, apply, and digest confirmation behavior
   unchanged.

## Subtasks

- [ ] Report the five axes separately.
- [ ] Distinguish verified from not run on each.
- [ ] Gate completion language on retention and idempotence.
- [ ] Add the axes to machine output additively.

## Acceptance Criteria

- [ ] The final result shows all five axes, each verified or not run.
- [ ] A run where the idempotence check did not execute shows it as not run,
      not as verified.
- [ ] Completion language appears only when retention is verified and
      idempotence passed.
- [ ] A run with verified postimages but unverified retention does not read as
      complete.
- [ ] Machine output carries the same five axes, with every prior field
      unchanged.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/apply.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestResultStatusMatrix -count=1` — expected:
  exit 0; five axes, each verified or not run.
- `go test ./internal/baseline -run TestCompletionLanguageRequiresRetention -count=1`
  — expected: exit 0; verified postimages alone never read as complete.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 7; Core Features 10; User Experience.
- `_techspec.md` → API Contracts; Build Order 11.
