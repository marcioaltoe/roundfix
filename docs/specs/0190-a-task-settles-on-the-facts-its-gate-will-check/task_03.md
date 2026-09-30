---
task: task_03
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: pending
type: backend
complexity: medium
---

# Task 03: A Task is refused the Spec Consistency findings it introduced

## Overview

The QA gate's precondition runs the strict Spec Consistency Check and refuses
on its findings (`qaGatePrecondition` in `internal/daemon/task_engine.go`). A
Task that introduces such a finding settles `completed` today. This Task adds
the second in-process Settlement Check: the Daemon takes the refusing findings
when the Task starts and again inside each Verification attempt, and fails the
check for a finding that is new. A finding that was already present is not
attributed to the Task.

## Requirements

1. MUST take the baseline through `SettlementChecker.RefusingFindings`, in
   `executeTaskWorker`, after the pre-work probe and before the Agent starts,
   only when the plan's `settlementChecks` is true. A baseline error MUST
   settle the Task `failed` before Agent work, with the error in its reason.
2. MUST export `speccheck.RefusalReason`, and make `PreconditionRefusal` build
   its reasons with it, so the gate and the settlement identify a finding the
   same way. No reason the gate records today changes.
3. MUST add the check labeled `settlement check: spec consistency`, built in
   the new file `internal/daemon/settlement_spec_consistency.go`. It fails
   when a `RefusalReason` value is absent from the baseline. Its diagnostics
   list each new finding's code, sentence, locations and fix line.
4. MUST keep the QA gate's precondition result unchanged for the same tree.
5. MUST NOT rename or remove a top-level test, or edit
   `internal/speccheck/mechanical_test.go` or a file task_02 created.
6. MUST put the new tests in the two test files this Task creates, with any
   fake checker they need.

## Subtasks

- [ ] Export the refusal reason the gate records.
- [ ] Take the baseline at Task start and add the check.
- [ ] Add one test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A refusing finding absent at Task start fails the check, and the
      Verification Feedback names the finding's code.
- [ ] A refusing finding present at Task start does not fail the check.
- [ ] A baseline error settles the Task `failed` and the Agent never starts.
- [ ] A checker error inside the attempt fails the check.
- [ ] `RefusalReason` returns the code and sentence the gate records.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/speccheck/mechanical.go`
- creates: `internal/daemon/settlement_spec_consistency.go`
- creates: `internal/daemon/settlement_spec_consistency_test.go`
- creates: `internal/speccheck/refusal_reason_test.go`
- instruction: `docs/adr/0182-a-task-settles-on-the-facts-its-gate-will-check.md`
- instruction: `docs/adr/0096-the-qa-gate-proves-machine-facts-before-it-spends-an-agent-turn.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSettlementRefusesASpecConsistencyFindingTheTaskIntroduced|TestSettlementIgnoresASpecConsistencyFindingPresentAtStart|TestSpecConsistencyBaselineErrorFailsTheTaskBeforeAgentWork|TestSpecConsistencyCheckerErrorFailsTheCheck)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSettlementRefusesASpecConsistencyFindingTheTaskIntroduced TestSettlementIgnoresASpecConsistencyFindingPresentAtStart TestSpecConsistencyBaselineErrorFailsTheTaskBeforeAgentWork TestSpecConsistencyCheckerErrorFailsTheCheck; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the four named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestRefusalReasonIsTheReasonTheGateRecords)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRefusalReasonIsTheReasonTheGateRecords; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the named test does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Core Feature 3; Core Feature 5; Success Metric 2; Recorded limits
- [_techspec.md](_techspec.md) — Spec Consistency findings the Task introduced; Testing Approach 3; Build Order 3
- ADR-0182; ADR-0093; ADR-0094; ADR-0096; ADR-0176

## Result
