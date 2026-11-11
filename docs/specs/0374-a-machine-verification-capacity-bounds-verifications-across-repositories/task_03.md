---
task: task_03
spec: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
status: pending
type: backend
complexity: high
---

# Task 03: The repository gate takes the capacity, and runs list, deliver status and doctor name its holders

## Overview

The Delivery Queue's repository gate takes one machine unit for its
repository Verification and its Linux Lane together. `roundfix runs list`,
`roundfix deliver status` and `roundfix doctor` read the capacity's entries,
write nothing, and name each holder and waiter with its step, place in line
and start time. Their command guides and the Roundfix `runs`, `deliver` and
`setup` references show the new block and line.

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   `commandDeliveryWorkflow.Gate` runs the repository Verification and the
   Linux Lane without any machine-wide bound, and `runs list`, `deliver
   status` and `doctor` print no line about another repository's
   Verification.
2. MUST change `commandDeliveryWorkflow.Gate` in
   `internal/cli/deliver_workflow.go` as `_techspec.md` → API Contract 7 says:
   after the shared cache trim and before the repository Verification,
   acquire with `store.AcquireMachineVerification` holder kind `delivery`
   with the owner's checkout root (`workflow.loaded.GitRoot`), the Spec slug
   and step `repository gate`, with `Runs.MaxVerifying`; release after the
   Linux Lane on every return path, a release failure written as
   `roundfix: warning: <release error>` to the queue log; write on the first
   wait, word for word,
   `roundfix: Delivery Queue item <spec>: waiting for the Machine Verification Capacity; <held> of <capacity> held by <holders>`
   (`<holders>`: each holding entry as
   `<DescribeMachineVerificationHolder>, <step>`, joined by `; `); return an
   acquisition error wrapped with `wait for the Machine Verification Capacity`.
   No other queue step takes a unit.
3. MUST create `internal/cli/machine_verification_status.go` with one
   renderer that prints the block of API Contract 8 word for word:
   `Machine Verification Capacity: <held> of <capacity> held, <waiting> waiting`
   (or `Machine Verification Capacity: off, <held> held, <waiting> waiting`
   when the reader's `runs.max_verifying` is `0`), then
   `  holding: <holder description>, <step>, since <acquired_at>` per holder,
   `  waiting <n>: <holder description>, <step>, since <since>` per waiter
   numbered from 1, and `  unreadable: <path>` per unreadable entry, times in
   RFC3339 UTC, holders and waiters in ticket order.
4. MUST print the block from `store.ReadMachineVerification` with the
   command's loaded Roundfix Home and config: `runs list` (`internal/cli/runs.go`)
   on standard output after its listing in both formats, including after
   `No Runs found.` when no Run Database exists; `deliver status`
   (`internal/cli/deliver.go`) right before its `Limits:` line; each only
   when the view has an entry or an unreadable path. A read error is one
   stderr line `roundfix: warning: read Machine Verification Capacity: <error>`
   and does not change the exit code.
5. MUST add the `verificationCapacity` dependency to `doctorDependencies` in
   `internal/cli/doctor.go`, set by `defaultDoctorDependencies`, and print its
   result after the `disk-space:` results and before `linux-lane:` as API
   Contract 8 says: name `verification-capacity`
   (`HealthCheckVerificationCapacity` in `internal/cli/health.go`), status
   `ok` with detail `runs.max_verifying <n>; <held> held, <waiting> waiting`
   (`off; <held> held, <waiting> waiting` at `0`), followed by
   `; holding: <holder description>, <step>` and
   `; waiting <n>: <holder description>, <step>` per entry; status `warn`
   with `; unreadable entry <path>` and next action
   `inspect <path> and remove it if no Roundfix process wrote it` when an
   entry is unreadable, and `warn` naming the error when the directory cannot
   be read. Doctor writes nothing.
6. MUST create `internal/cli/machine_verification_status_test.go` with the
   tests `_techspec.md` → Testing Approach names for task_03:
   `TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns` (rows `passing
   gate` and `failing gate` through the delivery item fixture of
   `deliver_item_binary_test.go` with `Runs.MaxVerifying` `1` and
   `delivery.linux_lane` off; the configured `defaults.verification` is a
   shell command that copies the entry files' `state` and `kind` lines from
   the temporary Roundfix Home's `.roundfix/verification-capacity` directory,
   named by its absolute path in the command, to a file outside the
   repository; the copy shows exactly one `holding` `delivery` entry, and no
   entry remains after the gate), `TestRunsListNamesTheMachineCapacityHolderAndWaiters`,
   `TestDeliverStatusNamesTheMachineCapacityHolderAndWaiters`,
   `TestDoctorReportsTheMachineVerificationCapacity` (rows `empty`,
   `holder and waiter`, `off`, `unreadable entry`) and
   `TestRunsListMachineCapacityTranscript`, which asserts Surface Transcript 1
   whole: every line of standard output and standard error, in order, and the
   exit code, with only the transcript's two repository placeholders
   replaced. Entries are written by the test as files in the format of
   `_techspec.md` → Data Models with the test's own PID and identity, under a
   temporary `HOME`. Every top-level test calls `t.Parallel()` first.
7. MUST update every existing test that pins the whole standard output of
   `roundfix doctor` with the default doctor dependencies on the default
   branch when the Task runs, adding the `verification-capacity:` line, and
   name each in the Task's `## Result`; every other existing test of `runs
   list`, `deliver status`, `deliver` and `doctor` passes unchanged.
8. MUST document the surfaces with the exact strings of API Contracts 7 and
   8: `docs/user-guide/commands/runs.md` and `docs/user-guide/commands/deliver.md`
   each gain a `#### Machine Verification Capacity` subsection with the block
   as a `text` fence (`deliver.md` also the queue log line), and
   `docs/user-guide/commands/doctor.md` gains a `verification-capacity:`
   bullet after the `disk-space:` bullet giving its `ok` and `warn` forms.
   Change the canonical references `.agents/skills/roundfix/references/runs.md`
   (the block), `references/deliver.md` (the gate's unit and its log line,
   beside the Run slot wait) and `references/setup.md` (the doctor line and
   the user template key), each naming the Machine Verification Capacity;
   then run `make skills-sync` and raise the Roundfix skill version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Never write inside a `### QA settlement` section, never hand-edit a
   mirror, version field or record.
9. MUST, as the Task's last step, after its last guide and skill edit,
   re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Of the Governed Paths the Task may change only
   `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/deliver.md`,
   `.agents/skills/roundfix/references/runs.md`,
   `.agents/skills/roundfix/references/setup.md` and
   `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Characterize the gate and the three surfaces.
- [ ] Take one machine unit around the repository Verification and the Linux Lane.
- [ ] Add the renderer and print it from runs list and deliver status.
- [ ] Add the doctor check.
- [ ] Add the tests, the transcript included, and update the pinned doctor outputs.
- [ ] Document the three surfaces in their guides and references, sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The repository gate's command runs while its `delivery` entry is held, and no entry remains after it.
- [ ] `runs list`, `deliver status` and `doctor` print the prescribed lines for holders and waiters, and nothing new when the capacity is empty.
- [ ] None of the three surfaces writes under the Roundfix Home, and their guides and both reference copies show the same strings.

## Context

- instruction: `docs/adr/0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md`
- instruction: `internal/cli/deliver_test.go`
- instruction: `internal/cli/deliver_disk_full_test.go`
- instruction: `internal/cli/deliver_item_binary_test.go`
- instruction: `internal/cli/runs_list_table_test.go`
- instruction: `internal/cli/doctor_test.go`
- instruction: `internal/cli/doctor_disk_space_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/runs.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- creates: `internal/cli/machine_verification_status.go`
- creates: `internal/cli/machine_verification_status_test.go`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The capacity block of `runs list` and `deliver status` and the `verification-capacity:` doctor line — consumers: `./internal/cli`
- The repository gate's machine unit — consumers: `./internal/delivery`, `./internal/cli`
- The three command guides and the `runs`, `deliver` and `setup` references — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run '^(TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns|TestRunsListNamesTheMachineCapacityHolderAndWaiters|TestDeliverStatusNamesTheMachineCapacityHolderAndWaiters|TestDoctorReportsTheMachineVerificationCapacity|TestRunsListMachineCapacityTranscript)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns TestRunsListNamesTheMachineCapacityHolderAndWaiters TestDeliverStatusNamesTheMachineCapacityHolderAndWaiters TestDoctorReportsTheMachineVerificationCapacity TestRunsListMachineCapacityTranscript; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 40m -run 'Doctor|RunsList|DeliverStatus|Deliver|Gate' ./internal/cli && go test -count=1 -timeout 20m ./internal/delivery` — expected: exit 0. Before this Task the five tests do not exist; afterwards they pass, and every existing doctor, runs list and Delivery Queue test passes.

- `for f in docs/user-guide/commands/runs.md docs/user-guide/commands/deliver.md; do grep -q '^#### Machine Verification Capacity$' "$f" || { echo "missing section: $f" >&2; exit 1; }; done; grep -qF 'verification-capacity:' docs/user-guide/commands/doctor.md || { echo 'missing: doctor line' >&2; exit 1; }; for f in .agents/skills/roundfix/references/runs.md .agents/skills/roundfix/references/deliver.md .agents/skills/roundfix/references/setup.md skills/roundfix/references/runs.md skills/roundfix/references/deliver.md skills/roundfix/references/setup.md; do grep -qF 'Machine Verification Capacity' "$f" || { echo "missing phrase: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no command guide or reference names the capacity; afterwards they document it, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record pass.

## References

- `_prd.md` → Goal 4; User Story 3; Core Feature 3; Core Feature 8; Success Metric 5; Success Metric 6; Success Metric 7
- `_techspec.md` → Current behavior 3 and 7; API Contracts 7-9; Surface Transcript 1; Invariants 1 and 5; Declared breaks; Testing Approach; Build Order 3
- ADR-0384
