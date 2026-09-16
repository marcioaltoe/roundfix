---
task: task_04
spec: 0139-a-suite-that-passes-where-it-runs
status: pending
type: backend
complexity: medium
---

# Task 04: Run the repository Verification in the QA gate step

## Overview

Today the QA Agent runs the repository Verification inside its sandbox, so
whether the gate's static check passes depends on whether the Agent chose to
leave the sandbox. This slice moves that run to the Daemon, in the QA gate step,
before the report is written and before the Agent turn:

- a pass is stated in the gate prompt;
- a failure or an unobserved outcome becomes the existing Precondition Refusal,
  which withholds the Agent with verdict `fail`.

Nothing is copied into Spec evidence and no report section is added.

## Requirements

1. MUST, in the QA gate step, run the configured repository Verification command
   once when both hold:
   - a command is configured;
   - the mechanical result does not already withhold the Agent.

   The run happens after the mechanical stage and before the QA Report is
   written, with shared Verification Capacity, attempt 1 and the QA Task as Work
   Item. It MUST take no Verification Feedback and no temporary-failure retry,
   exit status 75 included.
2. MUST leave Verification log retention as it is: kept on failure, removed on
   success.
3. MUST, on a pass, add these lines to the gate prompt after the seeded-report
   instruction, with the actual values, and continue to the Agent turn:

   ```text
   Repository Verification: already run by the Daemon outside the Agent sandbox.
   - command: <command>
   - verdict: pass (exit 0)
   - diagnostics: <path or "removed on success">
   Record this as the static gate result. Do not run the repository Verification again.
   ```

4. MUST, on a non-zero exit, set the mechanical result's Precondition Refusal and
   mark it blocking before the report is written:
   - check name: the command;
   - reason: `exited <status>; diagnostics: <path>`.

   The existing refusal report then records `verdict: fail`, the Agent is
   withheld, and the QA Task settles from that report.
5. MUST, when the runner reports an unobserved outcome, record the same refusal
   with the reason `outcome unobserved: <cause>; diagnostics: <path>`. When no
   diagnostics were retained, the reason MUST say so rather than leave the field
   empty.
6. MUST leave the Verification event publisher untouched: no new classification,
   no changed payload shape, and no claim about how the stream projects the
   attempt.
7. MUST leave any other runner error on the existing infrastructure-failure path,
   which ends the Run. Only a command verdict and an unobserved outcome become
   refusals.
8. MUST, when the Run's context is cancelled, return through the existing stop
   path without recording a refusal.
9. MUST, when no command is configured, run nothing and state in the gate prompt
   that the gate runs the repository Verification itself.
10. MUST NOT copy any log into Spec evidence, add a report section, change
    evidence-path detection, add a Mechanical Refusal Code, Verification
    classification, report frontmatter key or row status, or change the qa-gate
    skill or any governed path.
11. MUST add these tests, in the Daemon QA gate tests:
    - `TestQAGateRunsRepositoryVerificationBeforeAgentSession`
    - `TestQAGateRefusesOnFailedRepositoryVerification`
    - `TestQAGateRefusesOnUnobservedRepositoryVerification`
    - `TestQAGateSkipsRepositoryVerificationWhenMechanicalStageWithholds`
    - `TestQAGateWithoutConfiguredRepositoryVerificationLeavesItToTheGate`

## Subtasks

- [ ] Run the configured repository Verification in the QA gate step.
- [ ] State a pass in the prompt and continue to the Agent turn.
- [ ] Record a failure or unobserved outcome as a Precondition Refusal.
- [ ] Publish one payload per outcome.
- [ ] Cover the pass, failure, unobserved, withheld and unconfigured cases.

## Acceptance Criteria

- [ ] A passing repository Verification runs once before the Agent session, and
      the prompt carries the statement with the command, verdict and diagnostics.
- [ ] A failing repository Verification produces a Precondition Refusal report
      with verdict `fail` naming the command, its exit status and its diagnostics
      path. No Agent session starts.
- [ ] An unobserved outcome produces the same refusal, and its reason names the
      cause and the diagnostics state.
- [ ] Exit status 75 produces a refusal with no retry.
- [ ] When the mechanical stage withholds the Agent, no repository Verification
      runs. Without a configured command nothing runs, and the prompt states that
      the gate runs it.
- [ ] The existing verifier, mechanical-stage, QA prompt and QA report tests still
      pass, and Verification log retention is unchanged.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/spec/qa.go`

## Verification

- `grep -q 'func TestQAGateRefusesOnFailedRepositoryVerification' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^(TestQAGateRunsRepositoryVerificationBeforeAgentSession|TestQAGateRefusesOnFailedRepositoryVerification|TestQAGateRefusesOnUnobservedRepositoryVerification|TestQAGateSkipsRepositoryVerificationWhenMechanicalStageWithholds|TestQAGateWithoutConfiguredRepositoryVerificationLeavesItToTheGate)$'` — the pass, failure, unobserved, withheld and unconfigured cases hold; this fails today.
- `grep -q 'already run by the Daemon outside the Agent sandbox' internal/daemon/task_engine.go` — the gate prompt carries the statement; this fails today.
- `grep -q 'func TestQAGateRunsRepositoryVerificationBeforeAgentSession' internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/daemon -run '^(TestExecVerifier|TestMechanicalStage|TestTaskCycleQA|TestWriteMechanicalQAReport|TestQAMechanicalRequest)'` — the existing verifier and QA gate step tests still pass.

## References

- `_prd.md` → Goals 2 and 4; User Story 2; Core Features 4-7; Declared
  intentional breaks; Regression locks.
- `_techspec.md` → Implementation Design: QA gate repository Verification; Data
  Models; Testing Approach 4; Build Order 4.
- ADR-0096; ADR-0111; ADR-0135; ADR-0014; ADR-0056.
