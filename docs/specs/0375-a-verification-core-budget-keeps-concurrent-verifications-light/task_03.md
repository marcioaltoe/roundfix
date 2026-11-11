---
task: task_03
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
status: pending
type: backend
complexity: high
---

# Task 03: The repository gate and the Linux Lane receive the budget, and doctor reports it

## Overview

The Delivery Queue's repository gate computes the Verification Core Budget
once and gives it to its repository Verification and to the Linux Lane
container, as `roundfix deliver lane` does for its own container. `roundfix
doctor` prints one `verification-cores:` line that names the budget, its
source and the active Implement Runs, and warns when the bounds let the
Verifications take more cores than the machine has. The doctor and deliver
guides and the Roundfix `deliver` and `setup` references say the same.

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   `commandDeliveryWorkflow.Gate` calls `daemon.ExecVerifier` without a
   budget, `runLinuxLane` builds the container environment only from its
   fixed `-e` arguments, and `roundfix doctor` prints no line about cores;
   name every existing test that pins the Linux Lane's `docker run`
   arguments and every test that pins the whole standard output of `roundfix
   doctor` with the default doctor dependencies, on the default branch when
   this Task runs.
2. MUST change `commandDeliveryWorkflow.Gate` (`internal/cli/deliver_workflow.go`)
   to compute `verificationCoreBudget(ctx, workflow.loaded.Config).Cores` once
   and pass it as `VerifyRequest.Cores` to the repository Verification and as
   `Cores` to its `linuxLaneRequest`; `roundfix deliver lane`
   (`internal/cli/deliver_lane.go`) passes its own budget the same way. The
   item binary build and the derived-file regeneration pass no budget.
3. MUST add `Cores int` to `linuxLaneRequest` and make `runLinuxLane`
   (`internal/cli/deliver_linux_lane.go`) insert, when `Cores > 0`, the four
   arguments `-e`, `GOMAXPROCS=<n>`, `-e`, `ROUNDFIX_VERIFICATION_CORES=<n>`
   right after `-e`, `GOSUMDB=off` and before the first `-v`, as API Contract
   4 says. Update the pinned argument list in
   `TestLinuxLaneRunsTheChangedPackagesInAContainer`
   (`internal/cli/deliver_linux_lane_test.go`) and every other test the
   characterization names: the fixture's command context gives eight logical
   CPUs, so the list gains `-e GOMAXPROCS=4 -e ROUNDFIX_VERIFICATION_CORES=4`.
4. MUST add `HealthCheckVerificationCores = "verification-cores"` to
   `internal/cli/health.go`, a `verificationCores` dependency to
   `doctorDependencies` (`internal/cli/doctor.go`), default
   `defaultDoctorVerificationCoresResult` in a new
   `internal/cli/doctor_verification_cores.go`, printed right after the
   `verification-capacity:` result and only when the dependency is not nil;
   the default counts `store.ActiveImplementRuns` through a read-only open of
   the Run Database (none is 0) and formats exactly as API Contract 5 says,
   including the two next actions word for word. Doctor writes nothing and
   the line never changes its exit code. Every test the characterization
   names that pins the whole doctor output gains the line or sets the
   dependency to nil the way it treats `diskSpace`; none is known at
   authoring time.
5. MUST create `internal/cli/deliver_core_budget_test.go` with
   `TestTheRepositoryGateCommandReceivesTheCoreBudget` (the delivery item
   fixture with `runs.max_verifying` 2 and twelve logical CPUs; the gate
   command writes `$GOMAXPROCS|$ROUNDFIX_VERIFICATION_CORES` to a file: `3|3`,
   and `5|5` when the row adds `GOMAXPROCS=5`), using task_02's hermetic
   budget-input helper, and `internal/cli/doctor_verification_cores_test.go`
   with `TestDoctorReportsTheVerificationCoreBudget` (rows `default`, `capacity
   off`, `undivided with three Runs`, `explicit within the machine`, `explicit
   beyond the machine`, `unreadable Run Database`, each comparing the whole
   printed line with the six lines of `_techspec.md` → Testing Approach, and
   row `default dependency without a Run Database`, the default dependency in
   a temporary home, which prints the first line with `0 active Implement
   Runs`). Every top-level test calls `t.Parallel()` first.
6. MUST keep every other existing Doctor, Deliver, Gate and Linux Lane test
   of `./internal/cli` and every test of `./internal/delivery` passing.
7. MUST document: in `docs/user-guide/commands/doctor.md`, a
   `verification-cores:` bullet after the `verification-capacity:` bullet
   with the three forms of API Contract 5; in
   `docs/user-guide/commands/deliver.md`, that the repository gate's
   Verification and the Linux Lane container receive `GOMAXPROCS` and
   `ROUNDFIX_VERIFICATION_CORES`; in the canonical
   `.agents/skills/roundfix/references/setup.md`, the doctor line; and in the
   canonical `.agents/skills/roundfix/references/deliver.md`, the gate and
   lane variables. Then run `make skills-sync` and raise the Roundfix skill
   version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Never write inside a `### QA settlement` section, never hand-edit a
   mirror, version field or record.
8. MUST, as the Task's last step, after its last guide and skill edit,
   re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Of the Governed Paths the Task may change only
   `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/deliver.md`,
   `.agents/skills/roundfix/references/setup.md` and
   `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Characterize the gate, the lane and the doctor output.
- [ ] Give the gate's Verification and the Linux Lane the budget.
- [ ] Add the doctor line.
- [ ] Add the gate and doctor tests and update the pinned lane arguments.
- [ ] Write the doctor and deliver guides and the `setup` and `deliver` references, sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The repository gate's command and the Linux Lane container see both variables equal to the budget.
- [ ] `roundfix doctor` prints the prescribed `verification-cores:` line for each source and warns exactly on oversubscription or an unreadable Run Database.
- [ ] The doctor and deliver guides and both copies of the `setup` and `deliver` references describe the budget.

## Context

- instruction: `docs/adr/0385-each-verification-receives-a-verification-core-budget-in-its-environment.md`
- instruction: `internal/cli/deliver_item_binary_test.go`
- instruction: `internal/cli/doctor_test.go`
- instruction: `internal/cli/doctor_storage.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_lane.go`
- interface: `internal/cli/deliver_linux_lane.go`
- interface: `internal/cli/deliver_linux_lane_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- creates: `internal/cli/doctor_verification_cores.go`
- creates: `internal/cli/doctor_verification_cores_test.go`
- creates: `internal/cli/deliver_core_budget_test.go`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The repository gate's budget and the Linux Lane arguments — consumers: `./internal/cli`, `./internal/delivery`
- The `verification-cores:` doctor line — consumers: `./internal/cli`
- The doctor and deliver guides and the `setup` and `deliver` references — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run '^(TestTheRepositoryGateCommandReceivesTheCoreBudget|TestDoctorReportsTheVerificationCoreBudget|TestLinuxLaneRunsTheChangedPackagesInAContainer)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheRepositoryGateCommandReceivesTheCoreBudget TestDoctorReportsTheVerificationCoreBudget TestLinuxLaneRunsTheChangedPackagesInAContainer; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 40m -run 'Doctor|Deliver|Gate|Lane' ./internal/cli && go test -count=1 -timeout 30m ./internal/delivery` — expected: exit 0. Before this Task the two new tests do not exist and the lane test does not expect the two variables; afterwards the three pass, and the Doctor, Deliver, Gate and Lane tests of the CLI and the Delivery Queue package pass.

- `grep -qF 'verification-cores:' docs/user-guide/commands/doctor.md || { echo 'missing: doctor guide line' >&2; exit 1; }; grep -qF 'ROUNDFIX_VERIFICATION_CORES' docs/user-guide/commands/deliver.md || { echo 'missing: deliver guide variable' >&2; exit 1; }; for f in .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do grep -qF 'verification-cores:' "$f" || { echo "missing doctor line: $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do grep -qF 'ROUNDFIX_VERIFICATION_CORES' "$f" || { echo "missing variable: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither guide nor reference names the line or the variable; afterwards they do, and the skills suite and the documentation contracts pass.

## References

- `_prd.md` → Goal 5; User Story 3; User Story 5; Core Feature 5; Core Feature 6; Core Feature 7; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior 4 and 6; Interfaces; API Contract 4; API Contract 5; API Contract 6; Invariant 4; Declared breaks; Testing Approach; Build Order 3
- ADR-0385
