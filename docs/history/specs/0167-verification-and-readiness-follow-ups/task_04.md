---
task: task_04
spec: 0167-verification-and-readiness-follow-ups
status: completed
type: backend
complexity: medium
---

# Task 04: A quiet project init and a budget test without a race

## Overview

Carried from Spec 0165: `roundfix init` writes the default template, which contains `runs.max_active`, to Project Config, so every later command warns that the key is ignored there; and `TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch` needs a Task to settle inside a real 500 ms budget and flakes without load.

## Requirements

1. MUST leave `runs.max_active` out of the template written to Project Config, keeping it in the User Config template.
2. MUST let a project-scope `init` followed by any command print no ignored-key warning.
3. MUST inject the Implement budget clock through the CLI command dependencies, and make the budget test cross the budget only after the fake committer records the settled Task, keeping the test's assertions.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Project-scope `init` then `Load` warns nothing.
- [ ] The budget test no longer depends on wall-clock timing.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/config/config.go`
- interface: `internal/cli/implement.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestProjectInitWritesNoRunCeiling|TestProjectInitThenLoadWarnsNothing|TestBudgetClockCrossesAfterTheSettledTask)$" ./internal/config ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProjectInitWritesNoRunCeiling TestProjectInitThenLoadWarnsNothing TestBudgetClockCrossesAfterTheSettledTask; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

### Implementation

- Project-scope `init` now renders a scope-specific template without the
  User Config-only `runs.max_active` key. `DefaultConfigYAML` and user-scope
  `init` retain the configured Run ceiling.
- The Implement command's budget clock now comes from `commandDependencies`.
  CLI budget checks, elapsed-reason rendering and the Daemon Task cycle share
  that clock; production uses `time.Now`.
- `TestBudgetClockCrossesAfterTheSettledTask` replaces the real 500 ms race
  with two explicit milestones: the fake committer records `task_01`, then
  `task_02` starts and advances the injected clock beyond the one-hour test
  budget. The original outcome, commit, Verification, Task-status, Run
  Worktree and Run Branch assertions remain.

### Acceptance evidence

- Project init then Load warns nothing: `TestProjectInitWritesNoRunCeiling`
  observed no Run ceiling in Project Config and the ceiling still present in
  User Config; `TestProjectInitThenLoadWarnsNothing` loaded the generated
  Project Config with empty stderr. The full `internal/config` package passed.
- The budget test has no elapsed-time dependency:
  `TestBudgetClockCrossesAfterTheSettledTask` passed with one Verification and
  one commit, `task_01` completed, `task_02` in progress, and the recoverable
  BudgetExceeded surfaces preserved. The same test passed under `-race`.

### Focused checks

- Pre-change: the two project-init tests failed because Project Config
  contained `max_active` and Load printed the ignored-key warning; the budget
  test did not compile because `commandDependencies` had no budget clock.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test -count=1
  ./internal/config` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -run
  '^(TestBudgetClockCrossesAfterTheSettledTask|TestBudgetExceededKeepsWorkSettledBeforeTheBudget|TestBudgetExceededRunIsFoundWithoutTheHeaderLine)$'
  ./internal/cli` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test -race -count=1 -run
  '^TestBudgetClockCrossesAfterTheSettledTask$' ./internal/cli` — passed.
- `rtk git diff --check` — passed.

The Daemon-owned command under `## Verification` was not run in this Agent
turn.
