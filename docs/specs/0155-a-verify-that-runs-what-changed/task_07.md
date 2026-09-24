---
task: task_07
spec: 0155-a-verify-that-runs-what-changed
status: completed
type: test
complexity: medium
---

# Task 07: The partition follows what the recipes run

## Overview

Corrective Task from the pre-PR review of 2026-09-24. The `internal/cli` half of the partition contract sorts tests with its own copy of the pattern, so it cannot fail when the Makefile's `-run` and `-skip` recipes drift apart or one of them is dropped.

## Requirements

1. MUST derive each set's `internal/cli` membership from what the `verify-changed` recipes actually run, reading their arguments from `make -n` or listing with `go test -list` under the same `-run` and `-skip` patterns.
2. MUST fail when the two lists do not partition every top-level test of the package.
3. MUST carry a negative control that changes a recipe pattern and observes the failure.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The contract passes on the tree and fails when a recipe pattern drifts.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/verifyselect/verifyselect_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestPartitionFollowsTheMakefileRecipes)$" ./internal/verifyselect 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestPartitionFollowsTheMakefileRecipes; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

- Implementation: the `internal/cli` contract now lists every top-level test with `go test -list`, derives core and Baseline membership from each `verify-changed-*` recipe's dry-run `-skip` or `-run` argument, and rejects missing or duplicate membership. The dry-run copy replaces only recursive `$(MAKE)` calls so GNU Make cannot execute the Baseline recipe while the test inspects it.
- Acceptance evidence: `TestPartitionFollowsTheMakefileRecipes` passes against the tree and carries a negative control that changes the copied Baseline recipe to `-run "TestNotDeclared"`; the same contract then reports a test selected by zero sets.
- Focused check: `GOCACHE=/tmp/roundfix-task07-gocache go test -count=1 ./internal/verifyselect -run '^TestPartitionFollowsTheMakefileRecipes$'` exited 0 (`ok`, 1.924s).
- Follow-up: the broader `GOCACHE=/tmp/roundfix-task07-gocache go test -count=1 ./internal/verifyselect` check reached the pre-existing package partition contract and failed because `./cmd/roundfix` is selected by both package sets. Task 07 changes only the CLI recipe contract and does not alter package membership.
- Daemon Verification was not run; the Daemon owns the declared command and Task settlement.

## Carry-forward provenance

- Source Run: `run_20260924T174549Z_5664c74b2c5ad5f5`
- Source commit: `a57e0fd96ab528f2918a57799245500b40332bf7`
