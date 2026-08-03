---
task: task_02
spec: 0072-qa-is-a-task-not-a-flag
status: completed
type: backend
complexity: high
---

# Task 02: Route the gate from the graph, not the request

## Overview

The Daemon stops asking the request whether to run the gate and asks the
graph. `plan.QA` derives from the loaded graph's gate declaration; a node of
type `qa`, once every dependency settles `completed`, is executed by the
existing `runQAGate` — same session kind, same before-snapshot, same report,
verdict, typed blocked-row counts, and `daemon.qa` events, per ADR-0080. A
`qa` node must never reach an Agent session: routing happens on Task Type,
and the engine refuses to schedule one as Agent work.

## Requirements

1. MUST derive the gate decision from the graph: a declared gate node means
   the cycle ends with the gate; a declined or legacy graph means it does
   not.
2. MUST execute the `qa` node through the existing gate step
   (`runQAGate`), recording its settlement in the node's task-file status
   like any Task, while report, verdict, and events stay byte-compatible
   with today's `--qa` output.
3. MUST NOT start the gate while any of its dependencies is unsettled — the
   existing withholding, now expressed by dependency semantics.
4. MUST refuse to schedule a `qa`-typed node as Agent work, with an error
   naming the node, even on a hand-edited graph.
5. MUST migrate the daemon tests that pass a QA request flag to graphs that
   declare a gate node, changing setup only — every assertion about report,
   verdict, counts, and events stays as it is.
6. MUST keep a failed Task suppressing the gate exactly as today: the gate
   node stays pending, resumable, and unreported.

## Subtasks

- [ ] Derive `plan.QA` (or its replacement) from the graph.
- [ ] Route `qa` nodes to `runQAGate`; settle their status from the verdict.
- [ ] Add the Agent-scheduling refusal for `qa` nodes.
- [ ] Migrate the gate tests from flag-driven to graph-driven setup.
- [ ] Prove the failed-Task path leaves the gate pending and resumable.

## Acceptance Criteria

- [ ] A graph ending in a gate node runs the gate after the last Task
      settles, never before, and the report and verdict match today's
      output for the same inputs.
- [ ] A graph with `qa: declined` or no declaration finishes its cycle with
      no gate step and no gate events.
- [ ] A cycle with one failed Task ends with the gate node still pending
      and no gate report.
- [ ] A `qa` node forced toward Agent scheduling fails with the node named.
- [ ] `git status --porcelain` shows no path outside `internal/daemon/`,
      `internal/spec/`, and this task file.

## Verification

- `go test ./internal/daemon -count=1 -run 'QA|Gate' -v | grep -q -- "--- PASS"`
  — expected: exit 0.
- `go test ./internal/daemon -count=1` — expected: exit 0.
- `go build -buildvcs=false ./...` — expected: exit 0.

## References

- `_prd.md` → Core Features 2, 5; Decisions (withholding stays).
- `_techspec.md` → Integration Points; Risks (a gate node must never reach
  an Agent); ADR-0080.

## Result

### Implementation

- `TaskCycle` now removes the graph's `qa`-typed node from ordinary Task
  scheduling and invokes the existing `runQAGate` only after that node's
  declared dependencies settle completed. The legacy request QA value is
  inert.
- The gate verdict settles the QA task file `completed` only for `pass` and
  `failed` for every other verdict. When a report exists, its existing commit
  also carries the settled QA task file; the report path, verdict, QA Agent
  session, and `daemon.qa` payload retain their existing contracts.
- A withheld gate never enters the scheduler's skipped path, so a failed
  dependency leaves the QA node pending, resumable, and absent from Task
  outcomes and gate events.
- The Task worker refuses a `qa`-typed node before worktree creation or Agent
  work and names the node in the error.
- Daemon QA fixtures now write validated graph declarations and terminal QA
  nodes instead of enabling a request flag. Declined and legacy fixtures prove
  that leftover request state cannot start a gate.

### Focused checks

- Red signal before implementation:
  `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/daemon -count=1 -run '^TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport/pass$'`
  failed because the QA node was counted as ordinary completed Agent work and
  produced no QA verdict.
- `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/daemon -run '^(TestTaskCycleQA|TestPerWorkAgentSessionMixedTaskTypesAndQA|TestTaskSchedulerRefusesQATaskAsAgentWork|TestQACommitDropsExecutableFileAndCommitsRemainingPaths|TestTaskCycleStopRequestBeforeQAStepSkipsQA|TestTaskCycleStopRequestMidWaveSkipsQAWithEveryTaskCompleted|TestTaskNeedsCompletedCoversEveryGateDependency)'`
  passed.
- `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/daemon -run '^(TestTaskCycleGatesDependenciesAndSkipsFailedDependencyChainsUnderConcurrency|TestTaskCycleValidatesPlan|TestInitialTaskRunStatusesSeedsEarlierRunCompletions)$'`
  passed.
- Final post-edit focused check:
  `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/daemon -run '^(TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport|TestTaskCycleQAStepSkippedUnlessEveryTaskCompleted|TestTaskCycleDeclinedAndLegacyGraphsIgnoreQARequestState|TestTaskSchedulerRefusesQATaskAsAgentWork|TestTaskCycleStopRequestBeforeQAStepSkipsQA|TestTaskCycleStopRequestMidWaveSkipsQAWithEveryTaskCompleted|TestTaskNeedsCompletedCoversEveryGateDependency|TestTaskCycleGatesDependenciesAndSkipsFailedDependencyChainsUnderConcurrency)$'`
  passed.
- `git diff --check` passed.
- The commands under `## Verification` were not run; Daemon Verification owns
  them.

### Acceptance evidence

- Gate ordering and byte-compatible gate output: the existing QA verdict
  matrix, per-work Agent-session test, QA-only Run test, prompt tests, and Stop
  Request tests pass with graph-declared QA setup. The matrix still asserts the
  same report paths, verdicts, QA Batch ordinal, commit message, and
  `daemon.qa` payload, and now also proves task-file settlement.
- Declined and legacy graphs: `TestTaskCycleDeclinedAndLegacyGraphsIgnoreQARequestState`
  passes with no QA prompt, report, verdict, commit, or `daemon.qa` event even
  when legacy request state is true.
- Failed dependency: `TestTaskCycleQAStepSkippedUnlessEveryTaskCompleted`
  passes and proves the gate task remains `pending`, has no Task outcome, and
  produces no QA report or event.
- Agent refusal: `TestTaskSchedulerRefusesQATaskAsAgentWork` passes and proves
  the refusal names a hand-built QA node before any Agent call.
- Changed-path scope: postflight shows changes only under `internal/daemon/`
  and this task file; no path outside the Task's allowed scope was introduced.

### Follow-up

- Task 03 owns removal of the now-inert `TaskPlan.QA` field and the Implement
  Command's request parameter.
