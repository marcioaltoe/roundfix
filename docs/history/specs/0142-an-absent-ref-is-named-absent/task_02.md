---
task: task_02
spec: 0142-an-absent-ref-is-named-absent
status: completed
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

- [x] Branch the classifier on the resolver's absent answer.
- [x] Preserve each candidate with the reason that names the branch.
- [x] Cover an absent-target set, an ambiguous-target set and a present-target
      set at the classifier's fixture seam.

## Acceptance Criteria

- [x] A set whose recorded target branch was deleted preserves every candidate
      with the absent-branch reason.
- [x] A set whose recorded target branch is ambiguous still fails.
- [x] A set whose target branch exists reaches the same outcome as before this
      Task, in the same run as the absent one.
- [x] Nothing new is released.

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

## Result

The Run Branch set classifier now consumes only the resolver's typed absent
answer. It enumerates the recorded set's existing Run Branches, preserves each
one with `reconciliationReasonTargetBranchAbsent`, and returns no current or
releasable branch. Other resolver errors, including the typed ambiguous answer,
still fail the set through the existing error path. A resolved target continues
through the existing QA Report and supersession logic unchanged.

`TestClassifyRunBranchSetPreservesAbsentTarget` covers the three required
states at the real Git fixture seam in one test. The present-target fixture
retains its current branch, current report, releasable branch, proof, and empty
preserved summary. The ambiguous fixture still returns `errBranchAmbiguous`.
The deleted-target fixture preserves both candidates with the exact absent
reason and asserts that no current or releasable branch is reported.

Focused evidence:

- Before the production change,
  `GOCACHE=/private/tmp/roundfix-task02-go-cache rtk go test -count=1 -run '^TestClassifyRunBranchSet' ./internal/worktree`
  failed to compile because `reconciliationReasonTargetBranchAbsent` did not
  exist.
- After the production change, the same classifier-focused command passed all
  six matching tests.
- `GOCACHE=/private/tmp/roundfix-task02-go-cache rtk go test -count=1 ./internal/worktree`
  passed all 109 package tests.
- `rtk git diff --check` reported no whitespace errors.
- `GOCACHE=/private/tmp/roundfix-task02-go-cache rtk make verify-incremental`
  reached and passed `internal/worktree`. The overall gate exited 2 because two
  unrelated `internal/cli` force-stop integration tests could not read the
  process table in the sandbox (`operation not permitted`).

Acceptance evidence:

- Deleted target: both recorded existing candidates are present only in
  `Preserved`, each with `target branch "main" is absent`.
- Ambiguous target: the classifier returns an error matching
  `errBranchAmbiguous`.
- Present target: the same test observes the prior current branch, QA Report,
  releasable branch, release proof, and summary counts.
- Release safety: the absent-target result has no current branch, releasable
  branch, releasable proof, or release path.

The Task's declared `## Verification` commands were not run; the Daemon owns
that verification and settlement.
