---
task: task_01
spec: 0134-a-governed-deletion-the-gate-can-see
status: pending
type: test
complexity: low
---

# Task 01: Characterize classification for addition, removal and rename

## Overview

Record what the governed-mutation classifier answers today for a Governed Path
that appears, one that is removed, and one that is renamed. The slice is
verifiable alone: it states the present contract, including the two answers the
next Task moves.

## Requirements

1. MUST record that a Governed Path appearing after the turn is classified as a
   mutation today.
2. MUST record that a Governed Path removed by the turn is classified as **no**
   mutation today, which is the answer Task 02 moves.
3. MUST record that renaming a Governed Path to an ungoverned name is
   classified as no mutation today, which is the same answer seen from the
   other side.
4. MUST NOT change classifier behavior, and MUST NOT weaken or delete an
   existing assertion.

## Subtasks

- [ ] Record the addition answer.
- [ ] Record the removal answer.
- [ ] Record the rename answer.

## Acceptance Criteria

- [ ] A case asserts the addition is classified as a mutation.
- [ ] A case asserts the removal is classified as no mutation today, marked as
      the behavior Task 02 changes.
- [ ] A case asserts the rename is classified as no mutation today, marked the
      same way.
- [ ] Every assertion that existed in the touched file before this Task still
      runs and passes.

## Context

- interface: `internal/daemon/task_engine.go`

## Verification

- `grep -q 'func TestGovernedMutationClassificationCharacterization' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestGovernedMutationClassificationCharacterization$'` — the three recorded answers execute; absence cannot pass.
- `grep -q 'func TestGovernedMutationClassificationCharacterization' internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestTaskCycleFixtureSeedIsCreatedOnce$'` — the pre-existing seed-reuse lock still passes beside the new characterization.

## References

- `_prd.md` → Regression locks.
- `_techspec.md` → Testing Approach observation 1; Build Order 1.
