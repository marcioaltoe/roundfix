---
task: task_09
spec: 0155-a-verify-that-runs-what-changed
status: completed
type: test
complexity: low
---

# Task 09: Duplicate membership is a named failure

## Overview

QA finding F-001 of 2026-09-24: Task 07 moved `internal/cli` membership to Makefile recipe inspection and removed the negative control that selects one `internal/cli` test in both recipes, so the partition contract no longer proves that duplicate membership is caught and named.

## Requirements

1. MUST add to `TestPartitionFollowsTheMakefileRecipes` a subtest named `overlapping recipes name the duplicated test` that makes both recipes select one `internal/cli` test and asserts the failure names that test.
2. MUST add to `TestPartitionCoversEveryTestExactlyOnce` a subtest named `package selected by both sets is named` that asserts the failure names the duplicated package.
3. MUST NOT change production code.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Both negative controls exist and pass on the tree.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/verifyselect/verifyselect_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestPartitionFollowsTheMakefileRecipes|TestPartitionCoversEveryTestExactlyOnce)$" ./internal/verifyselect 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestPartitionFollowsTheMakefileRecipes/overlapping_recipes_name_the_duplicated_test TestPartitionCoversEveryTestExactlyOnce/package_selected_by_both_sets_is_named; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither subtest exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

- Added `overlapping recipes name the duplicated test` to exercise the real
  Makefile recipe parser with a copied recipe whose Baseline pattern also
  selects one core `internal/cli` test. The assertion requires the resulting
  two-set diagnostic to name that test exactly.
- Added `package selected by both sets is named` to put one repository package
  in both computed package sets. The assertion requires the resulting two-set
  diagnostic to name that package exactly.
- Focused acceptance evidence:
  - `GOCACHE=/tmp/roundfix-task09-gocache go test -count=1 -run '^TestPartitionCoversEveryTestExactlyOnce/package_selected_by_both_sets_is_named$' ./internal/verifyselect` exited 0.
  - `GOCACHE=/tmp/roundfix-task09-gocache go test -count=1 -run '^TestPartitionFollowsTheMakefileRecipes/overlapping_recipes_name_the_duplicated_test$' ./internal/verifyselect` exited 0.
- The first focused invocation with the default Go cache did not reach
  compilation because the sandbox denied access to `~/Library/Caches/go-build`;
  rerunning with the task-scoped cache above resolved the environment issue.
- The Daemon-owned `## Verification` command was not run.
