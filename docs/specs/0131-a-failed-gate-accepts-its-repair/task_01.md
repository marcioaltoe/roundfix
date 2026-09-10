---
task: task_01
spec: 0131-a-failed-gate-accepts-its-repair
status: pending
type: test
complexity: low
---

# Task 01: Characterize the loader's answer for each gate verdict

## Overview

Record what the Task Graph loader answers today for a settled QA gate standing
above an incomplete dependency, separately for a completed gate and a failed
one. The slice is verifiable alone: it states the present contract, including
the answer the next Task changes.

## Requirements

1. MUST record that a completed gate above one incomplete dependency refuses,
   naming the error identity and the dependency the error reports.
2. MUST record that a failed gate above one incomplete dependency refuses today
   with the same error identity, which is the answer the next Task changes.
3. MUST record that a gate covering every dependency loads for both verdicts,
   so the characterization separates staleness from coverage.
4. MUST NOT change loader behavior, and MUST NOT weaken or delete an existing
   assertion.

## Subtasks

- [ ] Build Spec fixtures for a completed and a failed gate above an incomplete dependency.
- [ ] Record the refusal each one produces today, with its error identity.
- [ ] Record the loading case for both verdicts.

## Acceptance Criteria

- [ ] A case asserts the completed gate's refusal and the dependency it names.
- [ ] A case asserts the failed gate's refusal today, marked as the behavior the
      next Task changes.
- [ ] A case asserts that both verdicts load when every dependency is completed.
- [ ] Every assertion that existed in the file before this Task still runs.

## Context

- interface: `internal/spec/spec.go`

## Verification

- `grep -q 'func TestGateStalenessCharacterizesEachVerdict' internal/spec/spec_test.go && go test -count=1 ./internal/spec -run '^TestGateStalenessCharacterizesEachVerdict$'` — the per-verdict characterization executes; its absence cannot pass.
- `grep -q 'func TestGateStalenessCharacterizesEachVerdict' internal/spec/spec_test.go || exit 1; go test -count=1 ./internal/spec -run '^TestLoadRejectsAppendedTaskUnderSettledGate'` — the pre-existing stale-gate assertion still runs and passes beside the new characterization.

## References

- `_prd.md` → Success Metrics; Decisions: Regression locks.
- `_techspec.md` → Testing Approach observations 1-3; Build Order 1.
