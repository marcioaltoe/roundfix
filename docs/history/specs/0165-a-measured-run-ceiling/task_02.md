---
task: task_02
spec: 0165-a-measured-run-ceiling
status: completed
type: test
complexity: low
---

# Task 02: A budget test that does not race the header

## Overview

`TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch` configures a 500 ms Run Budget and reads the Run id only from the `Implement Run: <id>` header, which a loaded machine prints after the budget has expired, so the test fails under load, on `main` too.

## Requirements

1. MUST make the budget test identify its Run from the Run Database, or from any stderr line naming it, not from the header alone.
2. MUST add `TestBudgetExceededRunIsFoundWithoutTheHeaderLine`, proving the lookup succeeds for a stderr that lacks the header.
3. MUST NOT change the configured budget or any production code.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The budget test passes whether or not the header line was printed.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/implement_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestBudgetExceededRunIsFoundWithoutTheHeaderLine|TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestBudgetExceededRunIsFoundWithoutTheHeaderLine TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the first case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The budget test

## Result

The budget integration test now reads the Run id from the terminal
`Implement Run <id> reached BudgetExceeded.` diagnostic instead of depending
on the later `Implement Run: <id>` header. The configured 500 ms budget is
unchanged, and no production code changed. Its fixture now starts with the
Task whose preservation it checks already completed, so budget expiry during
Run Worktree setup or during the next Task produces the same assertions. The new
`TestBudgetExceededRunIsFoundWithoutTheHeaderLine` pins the non-header lookup
against the diagnostic form emitted by BudgetExceeded paths.

Acceptance evidence:

- `TestBudgetExceededRunIsFoundWithoutTheHeaderLine` passed with stderr that
  contains the BudgetExceeded terminal diagnostic and no header line.
- `TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch` passed while
  resolving the Run through the non-header diagnostic, then confirmed the
  BudgetExceeded state, the completed Task, a recoverable next Task, and the
  preserved Run Worktree and Run Branch whether the next Task started or not.

Focused checks:

- Pre-change signal: repository search found no
  `TestBudgetExceededRunIsFoundWithoutTheHeaderLine`, and the existing budget
  test called the header-only `implementRunIDFromStderr` helper.
- An interim focused rerun after only changing the parser reached
  BudgetExceeded before `task_01` settled and failed the old timing-sensitive
  completed-Task assertion; the fixture now removes that incidental race.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^TestBudgetExceededRunIsFoundWithoutTheHeaderLine$' ./internal/cli` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=5 -run '^TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch$' ./internal/cli` — passed five consecutive runs.
- `rtk make verify-incremental` — passed vet, the full Go suite, skill
  consistency checks, and the Roundfix build.

The Daemon-owned command under `## Verification` was not run in this Agent
turn.
