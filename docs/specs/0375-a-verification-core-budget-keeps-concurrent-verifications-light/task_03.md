---
task: task_03
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
status: completed
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

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, golang-error-handling, golang-context, systematic-debugging, no-workarounds.

### Starting-tree characterization

- Starting HEAD: `4a313873276b5776400cd948b42297435cbc86e8`; the only pre-existing change was this Task's Daemon-owned status. Read-only inspection confirms `commandDeliveryWorkflow.Gate` calls `daemon.ExecVerifier` without `Cores`, `runLinuxLane` builds its container environment solely from the fixed `-e` arguments through `GOSUMDB=off`, and Doctor has no cores dependency, result or output line.
- Inspected both the starting HEAD and default branch `main` at `996769bbe4eb8e4187fbde68e90ec1ebb456e4a9` with `git ls-tree` and `git show`. The sole whole `docker run` argument pin is `TestLinuxLaneRunsTheChangedPackagesInAContainer`. `TestDeliverLaneSelectsPackagesAndImage` checks selected run arguments; cancellation tests pin `docker rm -f`, not `docker run`.
- No existing whole-output Doctor assertion uses the entire unmodified default Doctor dependency set. Literal whole-output assertions are `TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts`, `TestRunDoctorMissingRepositoryRoot`, and `TestRunDoctorRealRepositoryCheckDoesNotMutateState`; all use `withDoctorFakeLoadedAndReadiness`, whose optional dependencies are nil. The last explicitly restores only the default capacity dependency. Baseline-derived whole-output comparisons are `TestDoctorWarnsWhenALocationIsLowOnDisk`, `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn`, and `TestDoctorReportsTheMachineVerificationCapacity`; each also uses that injected fixture and selectively restores its subject dependency. Their nil cores dependency remains isolated, like diskSpace. The new tests enable the default cores dependency explicitly.

### Implementation and acceptance evidence

- Repository gate and Linux Lane: `Gate` computes the budget once and shares it between the repository `VerifyRequest` and `linuxLaneRequest`; the on-demand lane computes its own budget. Container variables are inserted immediately after `GOSUMDB=off` and before the mounts only for positive budgets. The item build and derived regeneration verifier requests remain unchanged with zero `Cores`.
  `TestTheRepositoryGateCommandReceivesTheCoreBudget` uses task_02's hermetic inputs and the delivery-item fixture: the real shell writes `3|3` on twelve CPUs at capacity 2, and `5|5` with explicit `GOMAXPROCS=5`. `TestLinuxLaneRunsTheChangedPackagesInAContainer` pins the entire argument list with both variables at 4 on eight CPUs. `TestDeliverLanePassesItsVerificationCoreBudget` additionally checks the on-demand lane's adjacent argument sequence at budgets 3 and 5.
- Doctor: the optional default dependency prints immediately after capacity, or after disk space without that dependency. It reads all active Implement Runs through `openGCStore(ctx, true, home)`, reusing the read-only snapshot reader, and treats an absent database as zero. It reports each source, oversubscription and database errors with the exact next actions. The concurrent-core multiplication uses `int64`.
  `TestDoctorReportsTheVerificationCoreBudget` compares the whole stdout against the six prescribed complete lines and the absent-database line. Its real temporary database includes active Implement Runs in different repositories, an excluded terminal Implement Run and an excluded active review-kind Run. Every row also preserves both healthy and failed readiness exit codes and compares the entire temporary home before and after; missing databases stay absent. The baseline run with a nil cores dependency proves the opt-out and isolates existing transcript fixtures.
- Documentation: the doctor guide and canonical setup reference describe all three forms, the source, active count and warning actions; the deliver guide and canonical deliver reference describe both variables on the gate and lane, explicit overrides and the unchanged build/regeneration scope. `make skills-sync` generated the mirrors. The version recorder raised both Roundfix skill version fields to `0.1.216` and appended its digest; no mirror, version or record was hand-edited. `TestAuthorialSkillSync` and the owned-version check passed on the generated files.

### Focused checks

- Before wiring, `go test -count=1 -timeout 5m ./internal/cli -run '^(TestTheRepositoryGateCommandReceivesTheCoreBudget|TestDeliverLanePassesItsVerificationCoreBudget)$'` failed as intended: the shell saw only its inherited Go value and Docker had neither budget variable.
- `go test -count=1 -timeout 5m ./internal/cli -run '^(TestTheRepositoryGateCommandReceivesTheCoreBudget|TestDeliverLanePassesItsVerificationCoreBudget|TestDoctorReportsTheVerificationCoreBudget|TestLinuxLaneRunsTheChangedPackagesInAContainer|TestRunDoctorRealRepositoryCheckDoesNotMutateState)$'` — exit 0.
- `go test -count=1 -timeout 8m ./internal/cli -run '^(TestRunDoctor|TestDoctor|TestLinuxLane|TestDeliverLane|TestDeliveryStep|TestDeliveryGate|TestTheRepositoryGate)'` — exit 0.
- `go test -count=1 -timeout 5m ./internal/cli -run '^(TestGate|TestRepositoryGate|TestTheRepositoryGateHolds|TestTheRepositoryGateOnAFullDisk|TestTheItemBuildOnAFullDisk|TestAConflictConfinedToDeclaredLines|TestADerivedConflict|TestAnArchivedItemsDerivedRecord)'` — exit 0 outside the sandbox. The sandbox attempt failed in the archived item's nested coverage regeneration; the same check passed without a source change when rerun with the required filesystem access. Earlier Go-cache and canonical skill-version write denials were likewise resolved by rerunning the same commands with filesystem access. Initial test-fixture compilation/required-field mistakes were corrected before these green runs.
- `make skills-sync` — exit 0; `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — exit 0.
- `make baseline-digests` — exit 0, `changed: false`; no additional derived paths changed.
- `go test -count=1 ./skills -run '^(TestAuthorialSkillSync|TestEveryOwnedSkillVersionIsRecorded)$'` — exit 0.
- `git diff --check` and `gofmt -l` over every changed Go file — exit 0, no formatting diagnostics.

### Final regeneration and settlement boundary

After this Result write, the last implementation step is the required command
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
It records the Behavior Surface fingerprints after all guide, skill, version and
code edits; its outcome accompanies the Agent handback. Read-only postflight
checks the resulting path set and confirms canonical/mirror equality.

Declared Verification, the full CLI/delivery and documentation suites, and Task
settlement remain for the Daemon. No Task status, sibling Task or Task Graph was
edited by this Agent, and no commit, push, pull request or message was made.

## Carry-forward provenance

- Source Run: `run_20261111T064558Z_6d985452394229fa`
- Source commit: `a77ea3a6e532307aa3fa27eca66e30eb305c321a`
