---
task: task_06
spec: 0090-a-gate-that-could-have-failed
status: pending
type: test
complexity: low
---

# Task 06: Source every wait budget from its shared constant

## Overview

On 2026-08-09 the authoritative gate returned exit 2 and then exit 0 on one
unchanged tree, because a wait in the agent harness used five seconds where every
sibling wait uses the shared ninety-second budget. The same gate then failed in
CI on a documentation-only commit, in a process-tree test carrying its own
literal deadline. Neither assertion was wrong. Both made the gate answer a
question about the machine when it was asked a question about the tree.

## Requirements

1. MUST replace every literal wait budget in the agent and store test harnesses
   with the shared constant its siblings already use.
2. MUST give the process-tree controller test a budget sourced the same way,
   sized for a loaded CI runner rather than an idle laptop.
3. MUST keep every assertion able to fail when the behaviour it guards is
   genuinely absent; this Task changes how long a test waits, never what it
   asserts.
4. MUST add a check that fails when a wait budget is restated as a literal
   beside an assertion instead of drawn from its shared constant.

## Subtasks

- [ ] Replace the literal budgets.
- [ ] Give the store harness a shared budget.
- [ ] Add the check that catches a restated literal.

## Acceptance Criteria

- [ ] No test in the agent or store harness constructs a wait budget from a
      literal duration at the call site.
- [ ] The guard fails when such a literal is reintroduced.
- [ ] Every touched test still fails when its subject is removed.

## Bounded scope

This Task may create or modify only:

- `internal/agent/acpx_runner_test.go`
- `internal/store/process_unix_test.go`
- `internal/store/process_budget_test.go`
- `docs/specs/0090-a-gate-that-could-have-failed/task_06.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/store -run '^TestWaitBudget' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestWaitBudgetIsNeverRestatedAtACallSite'` — expected: exits 0.
- `test -z "$(grep -nE 'time\.After\([0-9]+ ?\* ?time\.(Second|Millisecond)\)' internal/agent/acpx_runner_test.go)"` — expected: exits 0, proving no literal wait budget remains in the agent harness.
- `test -z "$(grep -nE 'newOwnerProcessController\([0-9]+ ?\*' internal/store/process_unix_test.go)"` — expected: exits 0, proving the controller budget is no longer built from literals at the call site.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent ./internal/store -count=1 2>&1 | grep -cq '^ok' ; test "$(GOCACHE="$PWD/.gocache" go test ./internal/agent ./internal/store -count=1 2>&1 | grep -c '^ok')" -eq 2` — expected: exits 0, proving both packages pass rather than one being skipped.

## References

- `_prd.md` → Goal 3; Core Features, a gate that is the same twice.
- `_techspec.md` → Build Order 6.
