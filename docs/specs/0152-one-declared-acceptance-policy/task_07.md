---
task: task_07
spec: 0152-one-declared-acceptance-policy
status: pending
type: backend
complexity: medium
---

# Task 07: A rendered command that needs nothing installed

## Overview

Task 03 removed the awk copy of the rule by having the rendered command invoke
`roundfix qa-report accept`. Pre-PR review raised that to a P1, and it is right:
the command then runs whichever `roundfix` is first on `PATH`.

Measured on this machine, a Homebrew 0.14.1 shadowed the current build and has
no `qa-report` command at all, so a valid report was rejected before the policy
ran. `.agents/skills/write-tasks/SKILL.md` forbids exactly this — Verification
must be satisfiable "in a fresh worktree using only repository state", and an
installed binary is ambient machine state.

Weakening the rendered command alone would open a hole: `settle` replays a
Task's Verification and does not special-case a `qa` Task, so a weaker command
would let it complete a gate whose report is ineligible.

So both settlement paths apply the decision in process, and the rendered command
stops deciding.

## Requirements

1. MUST render a command that invokes no installed binary and depends on no
   ambient machine state.
2. MUST have the rendered command still prove, from repository state alone, that
   a newest QA Report exists and that its verdict is readable, failing when
   either does not hold.
3. MUST apply the one eligibility decision in `settle` before it completes a
   `qa` Task, in process.
4. MUST keep the Daemon's gate settlement applying the decision exactly as Task
   06 left it.
5. MUST keep `settle` unchanged for every non-`qa` Task.
6. MUST NOT reintroduce a copy of the eligibility rule in awk or any other
   rendered form.

## Subtasks

- [ ] Render a hermetic command that proves presence and readability only.
- [ ] Apply the decision in `settle` for a `qa` Task.
- [ ] Add tests for the ineligible-report refusal through `settle`.

## Acceptance Criteria

- [ ] The rendered command contains no invocation of `roundfix`.
- [ ] It fails when no report exists and when the newest report's verdict cannot
      be read.
- [ ] `settle` refuses to complete a `qa` Task whose newest report is ineligible,
      and completes one whose report qualifies.
- [ ] A non-`qa` Task settles exactly as before.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/task.go`
- interface: `internal/cli/settle.go`

## Verification

- `test "$(grep -cF 'roundfix qa-report accept' internal/spec/task.go)" = "0"` — expected: exit 0; the rendered command invokes no installed binary. Before this Task it does, so the command fails.
- `out="$(go test -count=1 -v -run "^TestSettleAppliesEligibilityToAQATask" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestSettleAppliesEligibilityToAQATask"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — How the derived command reaches it
