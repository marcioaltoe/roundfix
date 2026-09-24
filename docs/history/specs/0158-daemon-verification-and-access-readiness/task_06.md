---
task: task_06
spec: 0158-daemon-verification-and-access-readiness
status: completed
type: backend
complexity: medium
---

# Task 06: A repair entry that the gate still closes

## Overview

Corrective Task from the pre-PR review of 2026-09-24. A named repair Task that enters a red repository gate can settle completed without the gate passing again, because Verification is reloaded from the Task file and the Agent can delete the repository command. Entry is also granted when the precondition had an unknown cause rather than a red verdict; the precondition's diagnostic is overwritten by the Task's first attempt; and exact matching of the configured command lets a padded bullet skip the entry check for an unnamed Task.

## Requirements

1. MUST remember that a Task entered on a red gate, and at settlement require the configured repository command, from the plan frozen at Run start, to run and pass regardless of the Task file's current Verification; otherwise settle failed.
2. MUST let a named Task enter only when the precondition produced a failed command verdict, never on an unknown cause.
3. MUST write the precondition check's diagnostics to their own path so later attempts cannot overwrite them.
4. MUST keep the trimmed match that decides whether a Task's Verification carries the repository command, and apply the exact-match rule only when validating `precondition_repairs` at planning.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] An Agent that deletes the repository command cannot settle the repair Task completed.
- [ ] An unknown-cause precondition does not open the entry.
- [ ] The red-entry diagnostic survives the Task's attempts.
- [ ] A padded repository command still triggers the entry check.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/task_engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRepairTaskCannotDropTheRepositoryCommand|TestRepairEntryNeedsAnObservedRedGate|TestPreconditionDiagnosticsAreNotOverwritten|TestPaddedRepositoryCommandStillChecksEntry|TestNamedTaskRepairsKnownRedPrecondition)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRepairTaskCannotDropTheRepositoryCommand TestRepairEntryNeedsAnObservedRedGate TestPreconditionDiagnosticsAreNotOverwritten TestPaddedRepositoryCommandStillChecksEntry TestNamedTaskRepairsKnownRedPrecondition; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task four of the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Precondition repair

## Result

### Implementation

- The Task worker now records entry only for an authorized precondition with
  an observed failed-command verdict. An unknown precondition settles the Task
  failed before the pre-work probe or Agent Session.
- A red-entry repair carries the repository command frozen in the Run plan
  into every post-Agent Verification attempt. Reloading the Task file cannot
  delete, pad, or replace that settlement obligation.
- Repository precondition diagnostics use a dedicated
  `batch-NNN-precondition.log` artifact, separate from Task attempt output.
- Runtime repository-command detection compares trimmed commands. Planning
  validation for `precondition_repairs` retains the exact verbatim match.

### Focused checks

- Before the production change,
  `rtk env GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 -run '^(TestRepairTaskCannotDropTheRepositoryCommand|TestRepairEntryNeedsAnObservedRedGate|TestPreconditionDiagnosticsAreNotOverwritten|TestPaddedRepositoryCommandStillChecksEntry|TestPreconditionRepairPlanningRequiresExactRepositoryCommand)$' ./internal/daemon`
  exited 1: the four corrective behavior tests failed and the planning
  exact-match control passed.
- After the production change, the same focused command exited 0.
- `rtk env GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 ./internal/daemon`
  exited 0.
- `rtk git diff --check` exited 0.
- The Task's declared `## Verification` command was not run; the Daemon owns
  that check.

### Acceptance evidence

- `TestRepairTaskCannotDropTheRepositoryCommand` passed: after the Agent
  removed the repository command, both settlement attempts restored the
  frozen command and the still-red command settled the Task failed.
- `TestRepairEntryNeedsAnObservedRedGate` passed: an unknown repository verdict
  ran only the entry check, started no Agent turn, and settled with the unknown
  cause retained.
- `TestPreconditionDiagnosticsAreNotOverwritten` passed: the red-entry marker
  remained at its dedicated path after Task Verification wrote its own
  diagnostic.
- `TestPaddedRepositoryCommandStillChecksEntry` passed: an unnamed Task with a
  padded repository-command bullet ran the entry check and started no Agent
  turn.
- `TestPreconditionRepairPlanningRequiresExactRepositoryCommand` passed: the
  runtime's trimmed match did not weaken the planning-time verbatim authority
  check.
