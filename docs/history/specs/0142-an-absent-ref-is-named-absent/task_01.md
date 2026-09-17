---
task: task_01
spec: 0142-an-absent-ref-is-named-absent
status: completed
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

- [x] Replace the single count refusal with an absent answer and an ambiguous
      answer.
- [x] Use the existing existence check for absence.
- [x] Cover absent, ambiguous and resolvable names at the resolver's fixture
      seam.

## Acceptance Criteria

- [x] A fixture repository with no such ref answers absent and names the branch.
- [x] A fixture with a tag and a branch of the same short name answers ambiguous.
- [x] A plain branch resolves to its head, as today.
- [x] A caller can distinguish the two refusals without reading their text.

## Context

- interface: `internal/worktree/worktree.go`

## Verification

- `grep -q "errBranchAbsent" internal/worktree/worktree.go` — expected: exit 0; the resolver has a distinct absent answer. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestResolveLocalBranchTellsAbsentFromAmbiguous$" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Feature 1; User Story 1; Goal 1; Success Metrics 1-2;
`_techspec.md` → Implementation Design: Absent and ambiguous are separate
answers; API Contract 2; Build Order 1; ADR-0127.

## Result

Implemented distinct sentinel errors for absent and ambiguous local branches.
The resolver now uses `localBranchExists` before ambiguity detection, wraps
`errBranchAbsent` when the local branch does not exist, and preserves the
existing ambiguity sentence while wrapping a separate `errBranchAmbiguous`.
Candidate counting now identifies ambiguity only when more than one matching
ref exists.

Added `TestResolveLocalBranchTellsAbsentFromAmbiguous` at the real Git fixture
seam. Its absent case matches `errBranchAbsent`, rejects
`errBranchAmbiguous`, and checks that the error names `missing-branch`. Its
ambiguous case creates both a local branch and tag, matches only
`errBranchAmbiguous`, and checks the existing sentence exactly. Its resolvable
case compares the returned head with `main`'s commit.

Focused evidence:

- Before the implementation,
  `rtk env GOCACHE=/private/tmp/roundfix-go-cache go test ./internal/worktree -run TestResolveLocalBranchTellsAbsentFromAmbiguous -count=1`
  failed to compile because `errBranchAbsent` and `errBranchAmbiguous` did not
  exist.
- After the implementation, the same focused command passed:
  `ok roundfix/internal/worktree 0.689s`.
- A focused regression run covering the new resolver test plus the existing
  ambiguous and missing Run Branch reconciliation tests passed:
  `ok roundfix/internal/worktree 0.544s`.
- `rtk env GOCACHE=/private/tmp/roundfix-go-cache make verify-incremental`
  reached and passed `internal/worktree` in 20.212s. The overall incremental
  check did not finish successfully because two unrelated `internal/cli`
  process ownership integration tests could not read the process table in the
  sandbox (`operation not permitted`).

The Task's declared `## Verification` commands were not run; the Daemon owns
that verification and settlement.
