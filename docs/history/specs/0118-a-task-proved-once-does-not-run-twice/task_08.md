---
status: completed
type: backend
---

# Task: A carry-forward that can carry a serial graph

Corrective, from this Spec's own first Run. Six Tasks completed and passed
their Verification, and carry-forward refused every one of them. Nothing had
drifted: task_04 edits the glossary and later Tasks correct the TechSpec, so
each Task settled after them compared its declared inputs against a checkout
that had not received those edits yet. The act reported them moved and refused
the whole set, which is the deadlock this Spec exists to remove.

## Work

- Compare each candidate's declared inputs against the accumulating staged
  state — the checkout plus the carries already staged ahead of it in
  dependency order — rather than against the raw checkout. The act already
  cherry-picks into a temporary staging worktree in order, so the state to
  compare against already exists.
- This strengthens the proof, and must not weaken it. An input that moved for
  any reason other than a preceding staged carry still refuses, and the
  whole-set refusal stays exactly as it is. Prove both: a Task whose Context
  genuinely moved must still refuse when it is carried in the same set as
  Tasks that edit shared inputs.
- Change the Implement Preflight's refusal condition from "at least one Task is
  carriable" to "the set would carry". The two differ precisely when the set
  refuses as a whole, and the weaker one makes implement refuse while naming a
  command that then refuses — the deadlock Core Feature 4 forbids.
- Restore the Run Database open to its place after the profile preflight. It
  was hoisted above it, so a profile-preflight failure now creates a Run
  Database; three existing tests encode that it must not. Keep the carry-forward
  check beside the Run Window check, which is where the TechSpec placed it.
- Cover the ordering with the existing preflight tests rather than new ones:
  they already assert the invariant and currently fail.

## References

- `_prd.md` → Goal 4, User Story 3, Core Features 4 and 5
- `_techspec.md` → Build Order 7; Risks & Considerations, the declared-input
  baseline, the whole-set refusal, and the Run Database ordering
- ADR-0053 keeps the act proof-based: this changes the baseline a proof reads,
  never whether the proof runs

## Verification
- `grep -q "TestImplementCarryForwardRefusesOnlyWhenTheSetWouldCarry" internal/cli/implement_test.go && grep -q "TestCarryForwardInputsResolveAgainstStagedCarries" internal/cli/carryforward_test.go && go test -count=1 ./internal/cli`

## Result

- Task Carry-Forward now builds its proof baseline in a temporary staging
  worktree seeded from the checkout, then stages each proved candidate in Task
  Graph order before comparing the next candidate. The apply path reuses the
  same candidate-staging operation.
- The whole-set contract remains strict: a genuine shared-input change and the
  existing moved-input and mixed-set cases refuse without carrying any Task.
- The Implement Preflight now selects a prior Run only when every candidate in
  its set would carry. A mixed ready/refusing set is reported and Implement
  proceeds instead of naming a carry-forward command that would refuse.
- Agent Selection Profile Preflight runs before the Run Database opens. The Run
  Window and Task Carry-Forward checks remain adjacent immediately after the
  open.

Focused-check evidence:

- Before implementation,
  `rtk go test -count=1 ./internal/cli -run '^(TestRunImplementPreflightProbeFailureCreatesNoRun|TestImplementProfilePreflightFailureCreatesNoRunWorktreeOrAgentPrompt|TestRunImplementSelectionFailureReportsProfileRemediationWithoutCreatingRun)$'`
  reported one pass and two failures because profile-preflight refusals created
  the Run Database.
- Before the production fix,
  `rtk go test -count=1 ./internal/cli -run '^(TestCarryForwardInputsResolveAgainstStagedCarries|TestImplementCarryForwardRefusesOnlyWhenTheSetWouldCarry)$'`
  reported the serial-graph and mixed-set regressions.
- After implementation,
  `rtk go test -count=1 ./internal/cli -run '^(TestCarryForward.*|TestInspectSpecCarryForwards|TestImplementCarryForward.*|TestImplementRefusesWhenCarryForwardIsAvailable|TestImplementProceedsWhenNothingIsCarriable|TestImplementRunWindow.*|TestRunImplementPreflightProbeFailureCreatesNoRun|TestImplementProfilePreflightFailureCreatesNoRunWorktreeOrAgentPrompt|TestRunImplementSelectionFailure.*)$'`
  passed 42 tests. This covers the serial staged baseline, genuine drift,
  whole-set refusal, mixed-set Implement continuation, Run selection, Run
  Window placement, and profile-preflight database ordering.
- The Daemon-owned `## Verification` command was not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20260827T150555Z_9e1f7632d21f0c1b`
- Source commit: `5c59c27e76a5376e15ca4d480cd89c1399035cc1`
