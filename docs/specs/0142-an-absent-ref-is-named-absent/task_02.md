---
task: task_02
spec: 0142-an-absent-ref-is-named-absent
status: pending
type: backend
complexity: medium
---

# Task 02: Preserve a set whose target branch is absent

## Overview

The Run Branch set classifier resolves the recorded target branch before it
walks the candidates, and any failure there refuses the whole set. This slice
preserves each candidate with a reason naming the absent branch instead, and
keeps every present-target outcome exactly as it is.

## Requirements

1. MUST classify a Run Branch set whose recorded target branch is absent by
   preserving each candidate, with a reason naming that branch as absent.
2. MUST keep failing the set when the recorded target branch is ambiguous,
   because the classifier cannot know which ref the Run meant.
3. MUST leave every outcome, reason and summary count unchanged when the target
   branch exists.
4. MUST release nothing that is preserved today.
5. MUST name the reason builder `reconciliationReasonTargetBranchAbsent`, so the
   Task's own check and later readers agree.

## Subtasks

- [ ] Branch the classifier on the resolver's absent answer.
- [ ] Preserve each candidate with the reason that names the branch.
- [ ] Cover an absent-target set, an ambiguous-target set and a present-target
      set at the classifier's fixture seam.

## Acceptance Criteria

- [ ] A set whose recorded target branch was deleted preserves every candidate
      with the absent-branch reason.
- [ ] A set whose recorded target branch is ambiguous still fails.
- [ ] A set whose target branch exists reaches the same outcome as before this
      Task, in the same run as the absent one.
- [ ] Nothing new is released.

## Context

- interface: `internal/worktree/worktree.go`

## Verification

- `grep -q "reconciliationReasonTargetBranchAbsent" internal/worktree/worktree.go` — expected: exit 0; the classifier names the absent branch. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestClassifyRunBranchSetPreservesAbsentTarget$" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 2-3; User Stories 2-3; Goals 2-3; Success Metric 3;
Regression locks; Acceptance evidence;
`_techspec.md` → Implementation Design: An absent target branch preserves its
set; API Contracts 1 and 3; Build Order 2; ADR-0057.
