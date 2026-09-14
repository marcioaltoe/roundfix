---
task: task_01
spec: 0134-a-governed-deletion-the-gate-can-see
status: completed
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

## Result

- Added `TestGovernedMutationClassificationCharacterization` at the existing
  Daemon classifier seam. Its table records addition as `true`, removal as
  `false`, and rename from a governed path to an ordinary path as `false`.
  The removal and rename rows state that Task 02 changes their answers to
  `true`; production classifier behavior is unchanged.
- Pre-change signal: `rtk rg -n '^func TestGovernedMutationClassificationCharacterization' internal/daemon/task_engine_test.go`
  exited 1 because the characterization did not yet exist.
- Focused check: `rtk go test -count=1 ./internal/daemon -run '^(TestGovernedMutationClassificationCharacterization|TestGovernedMutationDetectionUsesTheUnfilteredSnapshot|TestTaskCycleFixtureSeedIsCreatedOnce)$'`
  was initially blocked before test execution by `operation not permitted` in
  the host Go build cache. One unchanged retry with cache access passed all 6
  selected tests and subtests.
- Acceptance criterion 1: the passing `addition is a governed mutation`
  subtest asserts the current `true` answer.
- Acceptance criterion 2: the passing
  `removal is not yet a governed mutation` subtest asserts the current `false`
  answer and marks Task 02 as the owner of the change to `true`.
- Acceptance criterion 3: the passing
  `rename to an ungoverned path is not yet a governed mutation` subtest asserts
  the current `false` answer and marks Task 02 as the owner of the change to
  `true`.
- Acceptance criterion 4: the edit appends the characterization without
  changing an existing assertion. `rtk go test -count=1 ./internal/daemon`
  passed all 303 package tests after the edit.
- The Daemon-owned `## Verification` commands were not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20260913T221055Z_c70caa90b7bbad59`
- Source commit: `e57f096376ce6ad04b1fc873136154df53cdf577`
