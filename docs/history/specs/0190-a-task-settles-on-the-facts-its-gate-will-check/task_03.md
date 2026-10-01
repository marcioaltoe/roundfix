---
task: task_03
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: completed
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

Implemented the Task-start Spec Consistency baseline and the in-process check
`settlement check: spec consistency`. The worker reads refusing findings after
the pre-work probe and before creating the Agent Session owner, only for a
plan with Settlement Checks. It stores their `RefusalReason` values in the
worker's local Task plan; repair and temporary retry attempts use that same
baseline. A baseline read error settles the Task failed with the error in its
reason before Agent work.

The check reads refusing findings inside each Verification attempt, reports
new reasons with every location and fix line, and uses the existing diagnostic
artifact and Verification Feedback path. Checker errors fail the check.
`PreconditionRefusal` now calls the exported `RefusalReason`; the existing
formatting and deduplication behavior remain unchanged.

Focused evidence for each acceptance criterion:

| Criterion | Evidence |
| --- | --- |
| New refusing finding fails and Feedback names its code | `TestSettlementRefusesASpecConsistencyFindingTheTaskIntroduced` passed: both attempts refuse the finding, the repair prompt includes `SC-TEST`, the sentence, both locations and fix, final diagnostics contain the code, the fixture Task settles failed and no commit is requested. |
| Finding present at Task start does not fail | `TestSettlementIgnoresASpecConsistencyFindingPresentAtStart` passed: changing locations and fix text preserves the baseline identity, with one Agent request and one commit request. |
| Baseline error fails before Agent work | `TestSpecConsistencyBaselineErrorFailsTheTaskBeforeAgentWork` passed: one checker read, no Agent request or commit request, and the fixture's failed status/reason records `cannot read baseline`. |
| Checker error inside attempt fails | `TestSpecConsistencyCheckerErrorFailsTheCheck` passed: both attempts fail, Feedback names the check and `cannot read attempt findings`, and no commit is requested. |
| RefusalReason is the gate's recorded code and sentence | `TestRefusalReasonIsTheReasonTheGateRecords` passed for code and sentence, whitespace normalization, code only, sentence only and unnamed findings, including duplicate gate findings. |

`TestSpecConsistencyRepairKeepsTheTaskStartBaseline` also passed: a finding on
the first attempt triggers repair; removing it allows the second attempt to
settle, with exactly three checker reads across the cycle.

Focused checks:

- Before implementation, `GOCACHE=/tmp/roundfix-task03-gocache go test
  ./internal/daemon ./internal/speccheck -run
  'Test(SpecConsistency|SettlementRefusesA|SettlementIgnoresA|RefusalReason)'
  -count=1` exited 1: `RefusalReason` was undefined and the Daemon tests
  observed no refusing-finding reads.
- After implementation, `GOCACHE=/tmp/roundfix-task03-gocache go test
  ./internal/daemon ./internal/speccheck -run
  'Test(SpecConsistency|SettlementRefusesA|SettlementIgnoresA|RefusalReason|PreconditionRefusal|SettlementChecks|GatelessGraph)'
  -count=1` exited 0 for both packages.
- A fresh focused rerun after import formatting, `GOCACHE=/tmp/roundfix-task03-gocache
  go test ./internal/daemon ./internal/speccheck -run
  'Test(SpecConsistency|SettlementRefusesA|SettlementIgnoresA|RefusalReason|GateRefusal|MechanicalStageStoresThePrecondition)'
  -count=1`, exited 0 for both packages, including the existing gate-refusal
  and mechanical-report regression tests.
- The first `rtk make verify-incremental` exited 2. Process-owner integration
  tests could not read the sandboxed process table (`operation not permitted`).
  Suiteguard also detected this session's Result/import edits made while the
  suite was running. The retry uses the required process access and keeps the
  worktree unchanged throughout the command. The retry of `rtk make
  verify-incremental` exited 0: formatting, vet, package tests, skill sync/check
  and CLI build passed.

The authored Verification commands were not run. Task status, the Task Graph,
other Task files, `mechanical_test.go`, and files created by task_02 were not
edited. No commit, push or Pull Request was performed. No follow-up was found.

## Carry-forward provenance

- Source Run: `run_20260930T230750Z_197aa9aff15490a9`
- Source commit: `3a2ddff9ffd9ea55afb9ba31d9dcb0bdf5a87598`
