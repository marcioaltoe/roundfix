---
task: task_01
spec: 0090-a-gate-that-could-have-failed
status: pending
type: test
complexity: medium
---

# Task 01: Record what a vacuous and an unobserved gate do today

## Overview

The corpus that every later Task is measured against. It captures two current
behaviours before anything changes: a Task whose `## Verification` already
passes before the Agent runs settles `completed`, and a Verification the runner
could not observe settles `failed` with the same shape as a real refutation.
Both are the behaviours this Spec breaks, so both are declared here rather than
discovered later.

## Requirements

1. MUST add characterization tests that assert today's behaviour for a Task
   whose Verification command exits zero against the unchanged tree.
2. MUST add characterization tests that assert today's behaviour for a
   Verification command whose verdict the runner could not observe.
3. MUST name each declared break in the test's own documentation comment,
   stating which Task changes it and to what.
4. MUST use the fake command runner already used by the daemon tests rather than
   executing real repository commands.
5. MUST NOT change any production behaviour.

## Subtasks

- [ ] Capture the vacuous-gate behaviour.
- [ ] Capture the unobserved-verdict behaviour.
- [ ] Declare both breaks in comments naming their Task.

## Acceptance Criteria

- [ ] A test proves that a Task settles `completed` today when its Verification
      passed before the Agent turn.
- [ ] A test proves that an unobservable Verification is reported today with the
      same shape as a command that ran and failed.
- [ ] Each test's comment names the Task that will break it.

## Bounded scope

This Task may create or modify only:

- `internal/daemon/verification_probe_characterization_test.go`
- `docs/specs/0090-a-gate-that-could-have-failed/task_01.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestVerificationProbeCharacterization' -count=1 -v 2>&1 | grep -q '^--- PASS: TestVerificationProbeCharacterizationVacuousGateSettlesCompleted'` — expected: exits 0. A `-run` pattern that selects no cases still exits 0, so this asserts the named case actually ran and passed.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestVerificationProbeCharacterization' -count=1 -v 2>&1 | grep -q '^--- PASS: TestVerificationProbeCharacterizationUnobservedVerdictSettlesFailed'` — expected: exits 0.
- `grep -c 'Declared break: task_0' internal/daemon/verification_probe_characterization_test.go | grep -qE '^[2-9]'` — expected: exits 0, proving both breaks are declared and each names its Task.

## References

- `_prd.md` → Goals 1 and 4.
- `_techspec.md` → Build Order 1; Testing Approach.

