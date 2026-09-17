---
task: task_01
spec: 0142-an-absent-ref-is-named-absent
status: pending
type: backend
complexity: medium
---

# Task 01: Tell an absent ref from an ambiguous one

## Overview

The branch resolver counts candidate refs and refuses when the count is not one,
so a name that matches nothing is reported with the sentence written for a name
that matches several. This slice gives each case its own answer, using the
existence check already defined beside the resolver.

## Requirements

1. MUST report a short name that matches no ref as absent, naming the branch.
2. MUST keep reporting a short name that matches more than one ref as ambiguous,
   with today's sentence and branch name.
3. MUST let a caller tell the two cases apart from the returned error, without
   parsing its text.
4. MUST leave a name that resolves to exactly one ref resolving as it does today.
5. MUST name the absent case `errBranchAbsent`, so the Task's own check and later
   readers agree.

## Subtasks

- [ ] Replace the single count refusal with an absent answer and an ambiguous
      answer.
- [ ] Use the existing existence check for absence.
- [ ] Cover absent, ambiguous and resolvable names at the resolver's fixture
      seam.

## Acceptance Criteria

- [ ] A fixture repository with no such ref answers absent and names the branch.
- [ ] A fixture with a tag and a branch of the same short name answers ambiguous.
- [ ] A plain branch resolves to its head, as today.
- [ ] A caller can distinguish the two refusals without reading their text.

## Context

- interface: `internal/worktree/worktree.go`

## Verification

- `grep -q "errBranchAbsent" internal/worktree/worktree.go` — expected: exit 0; the resolver has a distinct absent answer. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestResolveLocalBranchTellsAbsentFromAmbiguous$" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Feature 1; User Story 1; Goal 1; Success Metrics 1-2;
`_techspec.md` → Implementation Design: Absent and ambiguous are separate
answers; API Contract 2; Build Order 1; ADR-0127.
