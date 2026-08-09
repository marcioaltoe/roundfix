---
task: task_01
spec: 0090-a-gate-that-could-have-failed
status: completed
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

## Result

### Implementation

- Added two daemon Task-cycle characterization tests using the existing
  `taskFakeRunner` and `taskFakeVerifier`; neither test executes a repository
  command.
- The vacuous-gate case records the current `agent>verify>commit` path and the
  Daemon-owned `completed` settlement when the fake Verification already
  passes without an Agent change.
- The unobserved-verdict case records the current failed settlement beside a
  command that returned a real non-zero verdict. It normalizes both terminal
  reasons and requires the same exact command-failure shape.
- Each test's documentation comment declares its break: `task_03` changes the
  vacuous case to refusal before the Agent turn, and `task_02` changes the
  unobserved case to an explicit unknown cause.

### Focused checks

- Before implementation, `rtk ls
  internal/daemon/verification_probe_characterization_test.go` exited `1`
  because the required characterization file did not exist.
- The first focused Go test attempt stopped before compilation because the
  sandbox denied the default macOS Go build cache under
  `/Users/marcio/Library/Caches/go-build`.
- `GOCACHE=/private/tmp/roundfix-task-01-gocache rtk go test
  ./internal/daemon -run
  '^TestVerificationProbeCharacterization(VacuousGateSettlesCompleted|UnobservedVerdictSettlesFailed)$'
  -count=1` passed both selected tests.

### Acceptance evidence

- Vacuous gate: the focused test asserts one completed Task, the persisted
  `completed` status, the fake command ledger, and the
  `agent>verify>commit` sequence.
- Unobserved verdict: the focused test asserts that both the unobserved runner
  error and a real non-zero verdict settle `failed`, produce no commit, and
  normalize to
  `Verification failed: command "<command>" exited with <cause>; diagnostics: <diagnostics>`.
- Declared breaks: source inspection finds one `Declared break` comment naming
  `task_03` and its pre-Agent refusal, and one naming `task_02` and its explicit
  unknown cause.
- Production scope: no production file changed; the only implementation path
  is the Task-authorized characterization test file.

### Follow-up

- The PRD and TechSpec attribute the authoritative repository gate decision to
  ADR-0083, but accepted ADR-0083 concerns adopted Spec sources. Correcting
  those out-of-scope Spec artifacts belongs to a separate documentation slice.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
