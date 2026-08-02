---
task: task_01
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: test
complexity: medium
---

# Task 01: Record how planning behaves today

## Overview

This Spec's largest slice turns a plan that completes today into one that can
stop. The only way to bound that is to record what completes today before
anything moves. This Task captures plan outcomes and diagnostics for real
repository shapes and changes no behavior; it is the gate every later slice is
measured against.

## Requirements

1. MUST record, for a corpus of repository shapes, the complete plan outcome:
   state, every diagnostic, every divergence, and every warning.
2. MUST cover at minimum a clean adoption, an idempotent re-plan after a
   verified apply, a repository with unsatisfied blocking capabilities, one
   with advisory-only divergences, and one whose Profile and catalog digests
   changed under an unchanged Baseline identifier.
3. MUST fail with a readable diff naming the affected shape and the changed
   field when a later change alters any recorded outcome.
4. MUST be regenerable through an explicit flag so an intended change is
   re-recorded deliberately, never silently.
5. MUST NOT change any production behavior, diagnostic, or exported API.

## Subtasks

- [ ] Assemble the repository-shape corpus.
- [ ] Record each shape's plan outcome deterministically.
- [ ] Add the comparison with a readable failure diff.
- [ ] Add the explicit regeneration flag.

## Acceptance Criteria

- [ ] The corpus contains a case for each shape named in Requirement 2.
- [ ] A test plans each shape and compares against its golden, passing on the
      unmodified tree.
- [ ] Deliberately altering one recorded outcome makes the test fail and name
      the affected shape.
- [ ] Two consecutive runs produce the same result and rewrite no golden.
- [ ] Goldens are re-recordable only through the explicit flag.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/profile_alignment.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0; the corpus matches the current tree.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0 on a second consecutive run, proving the comparison is stable
  and self-recording is gated.
- `grep -rq "TestBaselinePlanCharacterization" internal/baseline` — expected:
  exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 11; Success Metrics.
- `_techspec.md` → Testing Approach: characterization corpus; Build Order 1.
