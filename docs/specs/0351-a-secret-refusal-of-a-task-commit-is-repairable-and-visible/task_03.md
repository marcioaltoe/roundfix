---
task: task_03
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
status: pending
type: backend
complexity: medium
---

# Task 03: A run-unresolved park names its failed Tasks, their reasons and the skipped Tasks

## Overview

When an Implement Run ends Unresolved, the Delivery Queue parks the item with
the bare blocker `run-unresolved` and the answer "resolve the blocker, then
run roundfix deliver retry <slug>" (`runCandidate`,
`internal/delivery/engine.go`; `ClassifyPark`,
`internal/delivery/park_class.go`), so Fiscus could not see that the Task
commit's Secret Scan failed `task_01` or why its dependents were skipped.
After this Task the Run result carries the failed Tasks with their reasons
and the skipped Tasks from the Run's own events, the park records them after
`run-unresolved: `, and a Task commit that carried a secret gets a next
action that names the repair. This answers the second and third problems of
the report reported by fiscus on 2026-11-04
([references](references/2026-11-04-fiscus-a-task-commit-secret-refusal-is-invisible.md)).

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 10m -v -run '^(TestAnUnresolvedRunParksWithItsRunID|TestExistingGenericParkAnswersAreUnchanged)$' ./internal/delivery`
   and record in the Result that it passes with the bare blocker and the
   generic answer.
2. MUST add `TaskFailure`, `RunResult.FailedTasks`, `RunResult.SkippedTasks`
   and `RunUnresolvedBlocker` to `internal/delivery/engine.go` and use the
   function in `runCandidate`, exactly as `_techspec.md` → API Contract 8
   states: parts joined by `; `, each failed Task as `<task> failed:
   <reason>` with the reason whitespace-normalized and, past 300 characters,
   cut to its first 300 characters followed by `…`, then `skipped: <ids
   joined by ", ">`; the bare `run-unresolved` with no parts. Runtime
   infrastructure and the QA environment partial keep precedence.
3. MUST make `commandDeliveryWorkflow.runResult`
   (`internal/cli/deliver_workflow.go`) fill those fields from the
   `daemon.task` events it already reads (API Contract 7): `phase: settled`
   with `status: failed` gives a `TaskFailure` with the event's `reason`,
   `phase: skipped` gives a skipped Task, in event order, the Task taken from
   the payload `task` or else the event's Work Item. The exhausted runtime
   recovery reading keeps its behavior.
4. MUST make `ClassifyPark` answer a `run-unresolved` blocker whose detail
   contains ` failed: Task commit carries a secret: ` with exactly (API
   Contract 9):
   `if the failed Task's text prescribes the value the blocker locates, amend it in <worktree> to one no Secret Pattern matches (letters only, or fewer than 20 characters); then run roundfix deliver retry <slug>, whose rerun returns any Secret Match to the Agent as Verification Feedback`,
   with `<worktree>` the item's worktree and `<slug>` its Spec; every
   `run-unresolved` blocker keeps class `finding`, and every other one keeps
   "resolve the blocker, then run roundfix deliver retry <slug>".
5. MUST create `internal/delivery/run_unresolved_park_test.go` with
   `TestARunUnresolvedParkNamesItsFailedAndSkippedTasks`, calling
   `t.Parallel()` first, through the engine with the fake workflow and the
   real store (as `TestAnUnresolvedRunParksWithItsRunID` does), with cases:
   - `secret refusal`: `task_01` failed with `Task commit carries a secret:
     fixture.txt:2: credential-assignment`, `task_02` and `task_03` skipped;
     the blocker is
     `run-unresolved: task_01 failed: Task commit carries a secret: fixture.txt:2: credential-assignment; skipped: task_02, task_03`,
     and `ClassifyPark` and `PendingQuestionFor` give class `finding` and the
     answer of requirement 4;
   - `two failures without skips`: two failed Tasks give
     `run-unresolved: task_01 failed: <r1>; task_02 failed: <r2>` and the
     generic answer;
   - `long reason`: a 400-character reason is cut to 300 characters and `…`;
   - `nothing failed or skipped`: the bare `run-unresolved` and the generic
     answer.
6. MUST create `internal/cli/deliver_run_unresolved_park_test.go`, each test
   calling `t.Parallel()` first:
   - `TestDeliveryRunResultReadsFailedAndSkippedTasks`: a store Run, as in
     `TestDeliveryRunResultReadsAnExhaustedLostRollout`, holding a `task_01`
     `settled` `failed` event with a reason, `task_02` and `task_03`
     `skipped` events, a `task_04` `settled` `completed` event and a 0.130.1
     `secret_scan` event; `runResult` returns exactly the one failure and the
     two skipped Tasks in that order.
   - `TestDeliverStatusNamesTheSecretRefusalOfARunUnresolvedPark`: as
     `TestDeliverStatusPrintsAParkLinePerParkedItem` does, one item
     `0077-example` parked with the blocker of requirement 5's `secret
     refusal` case and worktree `/worktrees/0077-example`; assert the whole
     standard output of `roundfix deliver status`: the item row with that
     blocker and worktree, `Park: 0077-example finding: <answer>`, the
     `Limits:` and `Usage:` lines, `Pending question: 0077-example parked
     <blocker>` and `Answer: <answer>`, with the answer of requirement 4.
7. MUST add after the `secret-found` next-action paragraph of
   `docs/user-guide/commands/deliver.md`, verbatim apart from line wrapping:
   "A `run-unresolved` park names each Task the Run settled `failed`, with its
   reason cut to 300 characters, and the Tasks the Run skipped:
   `run-unresolved: <task> failed: <reason>; skipped: <task>, <task>`. A Run
   with neither keeps the bare `run-unresolved` and its answer. When a failed
   Task's commit carried a secret, the next action is `if the failed Task's
   text prescribes the value the blocker locates, amend it in <worktree> to
   one no Secret Pattern matches (letters only, or fewer than 20 characters);
   then run roundfix deliver retry <slug>, whose rerun returns any Secret
   Match to the Agent as Verification Feedback`."
   Then, as the LAST step after the last guide edit, re-record the Behavior
   Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   and commit `docs/references/behavior-surfaces.json`. task_04 writes the
   same paragraph into the Roundfix skill's deliver reference.
8. MUST keep every existing test of `internal/delivery` and every
   `TestDelivery*`, `TestDeliverStatus*` and `TestDeliverRetry*` test of
   `internal/cli` byte-identical and passing, and run the two sabotages of
   `_techspec.md` → Testing Approach for this Task on a temporary copy,
   recording which named test fails.
9. MUST NOT add a Park Class or a blocker name, change the Run Database
   schema, or touch `internal/daemon`, `internal/runcause`, `internal/store`,
   any skill, any other document, the Makefile, CI workflows, `go.mod` or
   `go.sum`.

## Subtasks

- [ ] Characterize the bare park.
- [ ] Carry failed and skipped Tasks in the Run result and build the park detail.
- [ ] Add the secret next action.
- [ ] Add the engine, Run result and status tests.
- [ ] Run the sabotages on a temporary copy.
- [ ] Document the park detail in the deliver guide and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] A run-unresolved park names each failed Task with its bounded reason and every skipped Task, in status and the Pending Question, and a secret refusal gets its repair as the next action.
- [ ] A Run with no failed or skipped Task keeps the bare blocker and today's answer, and every run-unresolved blocker stays class `finding`.

## Context

- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- creates: `internal/delivery/run_unresolved_park_test.go`
- interface: `internal/delivery/retry_test.go`
- interface: `internal/delivery/park_class_test.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_run_unresolved_park_test.go`
- interface: `internal/cli/deliver_runtime_infrastructure_test.go`
- interface: `internal/cli/deliver_park_status_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`
- instruction: `.agents/skills/roundfix/references/deliver.md`

## Consumers

- The Run result and the `run-unresolved` blocker printed by status and the Pending Question — consumers: `./internal/delivery`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestARunUnresolvedParkNamesItsFailedAndSkippedTasks|TestAnUnresolvedRunParksWithItsRunID|TestParkClass.*|TestExistingGenericParkAnswersAreUnchanged|TestEveryBlockerHasAParkClass|TestRetry.*)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestARunUnresolvedParkNamesItsFailedAndSkippedTasks ' || { echo 'missing pass: TestARunUnresolvedParkNamesItsFailedAndSkippedTasks' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -v -run '^(TestDeliveryRunResultReadsFailedAndSkippedTasks|TestDeliverStatusNamesTheSecretRefusalOfARunUnresolvedPark|TestDeliveryRunResult.*|TestDeliverStatus.*)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryRunResultReadsFailedAndSkippedTasks TestDeliverStatusNamesTheSecretRefusalOfARunUnresolvedPark; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'park names each Task the Run settled' || { echo 'missing park text in docs/user-guide/commands/deliver.md' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task none of the three new tests exists, so the first pass line is missing; afterwards the park names its failed and skipped Tasks, status prints the secret repair, the deliver guide states the detail with a current Behavior Surface Record, and the existing park, retry and status tests pass.

## References

- `_prd.md` → Goals; Core Feature 5; Core Feature 6; Success Metric 5
- `_techspec.md` → Current behavior 6; API Contract 7; API Contract 8; API Contract 9; Invariant 6; Declared breaks; Testing Approach; Build Order 3
