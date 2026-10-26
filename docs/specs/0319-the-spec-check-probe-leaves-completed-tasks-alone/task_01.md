---
task: task_01
spec: 0319-the-spec-check-probe-leaves-completed-tasks-alone
status: pending
type: backend
complexity: medium
---

# Task 01: Spec check probes only the Tasks a Run would still execute

## Overview

`roundfix spec check --run-verification` probes every Task in the graph, so a
resumed Spec's completed Task is reported `vacuous` or `unknown` and the check
exits `1`, while the Daemon never probes it. This Task exports the Daemon's
selection rule as `TaskAwaitsWork`, makes the Daemon's scheduler and QA
routing and the check read it, and reports each command the check leaves out
as `not probed` with its reason, in text, JSON and the verdict line.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-25-the-spec-check-probe-re-runs-completed-tasks.md`](references/2026-10-25-the-spec-check-probe-re-runs-completed-tasks.md)
   of 2026-10-25 (`_prd.md` → Sources 1) by implementing API Contracts 1-5 of
   `_techspec.md`.
2. MUST add `TaskAwaitsWork` to `internal/daemon/verification_probe.go` with
   the signature and comment of `_techspec.md` → Interfaces, and make
   `initialTaskRunStatuses` and both QA gate routing conditions in
   `internal/daemon/task_engine.go` read it. No other status comparison in the
   Daemon changes, and which Tasks the Daemon runs, probes and gates is
   unchanged.
3. MUST, in `internal/cli/spec_check.go` `probeSpecVerifications`, call
   `daemon.ProbeCommands` only for a Task `daemon.TaskAwaitsWork` selects —
   never a second status comparison — and for every other Task run no
   command and add the `notProbed` entries and ordered text lines of
   `_techspec.md` → Data Models and Interfaces. The exit code keeps reading
   only `verification.commands`; `uncommittedVerificationSources` keeps
   receiving every Task file.
4. MUST add `NotProbed` to `speccheck.VerificationCoverage` and the verdict
   line suffix of Interfaces to `RenderText`, keeping today's lines when it is
   zero, and print `No authored Verification commands.` only when the graph
   has no command at all.
5. MUST add `internal/cli/spec_check_completed_tasks_test.go`,
   `internal/daemon/task_awaits_work_test.go` and the two
   `TestVerdictLineStatesProbeCoverage` cases exactly as `_techspec.md` →
   Testing Approach describes, with those test and subtest names. Each new
   top-level test calls `t.Parallel()` first; the tests never read this
   repository's Specs and never reach the network.
6. MUST NOT change the prober's classification, its parser step, the
   disposable checkout, the JSON schema string, any existing assertion, any
   document or skill, the Makefile or the CI workflows.

## Subtasks

- [ ] Export the selection rule and make the Daemon read it.
- [ ] Select, report and count not-probed commands in the check.
- [ ] Add the CLI, Daemon and verdict-line tests.

## Acceptance Criteria

- [ ] A completed Task's commands never run, are named `not probed` in text and JSON, and leave the exit code alone.
- [ ] A vacuous, unknown or malformed command of a Task still to run, or of a pending QA gate, still fails the check.
- [ ] The Daemon's Task-cycle tests pass unchanged, and a rule change in `TaskAwaitsWork` fails both callers' tests.

## Consumers

- `daemon.TaskAwaitsWork` — consumers: `internal/cli`
- `verification.notProbed` (JSON) and `speccheck.VerificationCoverage.NotProbed` — consumers: `internal/cli`

## Context

- instruction: `internal/daemon/verification_probe_test.go`
- instruction: `internal/daemon/task_engine_test.go`
- instruction: `internal/cli/spec_check_test.go`
- instruction: `internal/cli/spec_check_provenance_test.go`
- creates: `internal/cli/spec_check_completed_tasks_test.go`
- creates: `internal/daemon/task_awaits_work_test.go`
- interface: `internal/daemon/verification_probe.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/speccheck/report.go`
- interface: `internal/speccheck/report_test.go`

## Verification

- `out="$(go test -timeout 15m -count=1 -v -run '^(TestSpecCheckLeavesCompletedTasksAlone|TestSpecCheckRunVerification|TestSpecCheckReportsAMalformedCommand|TestSpecCheckNamesAnUncommittedVerificationSource|TestSpecCheckNamesNoSourceWhenCommitted)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSpecCheckLeavesCompletedTasksAlone/a_resumed_Spec_probes_only_the_Tasks_still_to_run TestSpecCheckLeavesCompletedTasksAlone/a_pending_vacuous_command_still_fails TestSpecCheckLeavesCompletedTasksAlone/a_pending_unknown_command_still_fails TestSpecCheckLeavesCompletedTasksAlone/a_pending_malformed_command_still_fails TestSpecCheckLeavesCompletedTasksAlone/a_pending_QA_gate_is_still_probed TestSpecCheckLeavesCompletedTasksAlone/every_Task_completed TestSpecCheckLeavesCompletedTasksAlone/an_uncommitted_completed_Task_source_is_still_named TestSpecCheckLeavesCompletedTasksAlone/without_a_completed_Task_the_JSON_has_no_notProbed_key TestSpecCheckRunVerification TestSpecCheckReportsAMalformedCommand TestSpecCheckNamesAnUncommittedVerificationSource TestSpecCheckNamesNoSourceWhenCommitted; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -timeout 15m -count=1 -v -run '^(TestTaskAwaitsWorkSelectsEveryStatusButCompleted|TestTaskCycleRerunsStaleTasksAndSkipsCompletedTasks|TestTaskCycleQAStepSkippedUnlessEveryTaskCompleted|TestTaskCycleQAOnlyRunWhenEveryTaskAlreadyCompleted|TestTaskCycleQAStepRequiresEveryGraphTaskCompleted)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestTaskAwaitsWorkSelectsEveryStatusButCompleted " || { printf 'missing pass: TestTaskAwaitsWorkSelectsEveryStatusButCompleted\n' >&2; exit 1; }; out="$(go test -timeout 10m -count=1 -v -run '^TestVerdictLineStatesProbeCoverage$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in probe_left_completed_Tasks_unprobed probe_left_completed_Tasks_unprobed_and_could_not_execute_a_command probe_ran probe_did_not_run probe_ran_but_a_command_could_not_be_executed; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestVerdictLineStatesProbeCoverage/$name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the new test files and verdict-line cases do not exist, so the first pass check fails. Afterwards a completed Task is left unprobed and named, Tasks still to run keep their verdicts and exit code, the Daemon's selection is unchanged, and the verdict line counts not-probed commands.

## References

- `_prd.md` → Sources; Goals 1-4; Core Features 1 and 2; Success Metrics 1-5; Validation
- `_techspec.md` → Current behavior; Interfaces; Data Models; API Contracts 1-5; Invariants 1-5; Testing Approach; Build Order 2
