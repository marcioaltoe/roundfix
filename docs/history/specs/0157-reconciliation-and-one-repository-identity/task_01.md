---
task: task_01
spec: 0157-reconciliation-and-one-repository-identity
status: completed
type: backend
complexity: medium
---

# Task 01: Refuse a candidate without evidence

## Overview

Spec 0154's absent-target fallback can make a clean Run Branch a cleanup candidate that then revalidates without evidence; `cleanupTerminalRun` dereferences that nil evidence and panics inside `reconcile --apply`.

## Requirements

1. MUST refuse, with a reason naming the failed revalidation, a candidate whose revalidation returns no evidence.
2. MUST make the refusal before any Git command that removes a worktree or branch.
3. MUST keep every releasable candidate with positive evidence released as today.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A candidate revalidating without evidence is refused and nothing is removed.
- [ ] A candidate with positive evidence is still cleaned up.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestApplyRunBranchCandidateRefusesACandidateWithoutEvidence$" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestApplyRunBranchCandidateRefusesACandidateWithoutEvidence"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md)

## Result

### Implementation

- `ApplyRunBranchCandidate` now refuses a freshly `unintegrated` candidate when
  worktree revalidation supplies no evidence, before `cleanupTerminalRun` can
  issue a worktree-removal or branch-deletion command.
- Existing state-specific refusals remain unchanged, and candidates with fresh
  positive evidence continue through the existing cleanup path.

### Focused checks

- Before the production change,
  `rtk env GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 -run TestApplyRunBranchCandidateRefusesACandidateWithoutEvidence ./internal/worktree`
  exited 1 with a nil-pointer panic in `cleanupTerminalRun`.
- After the production change,
  `rtk env GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 -run 'TestApplyRunBranchCandidate(RefusesACandidateWithoutEvidence|RevalidatesProofAndCleanWorktree)$' ./internal/worktree`
  exited 0.
- After preserving the established dirty-candidate diagnostic,
  `rtk env GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 -run 'Test(ReconcileFallsBackToTheDefaultBranch|ApplyRunBranchCandidate(RefusesACandidateWithoutEvidence|RevalidatesProofAndCleanWorktree|PreservesNewlyDirtyWorktree))$' ./internal/worktree`
  exited 0.
- The first `make verify-incremental` run reached the test suite but the sandbox
  blocked a GitHub-dependent check. The approved rerun then exposed two existing
  dirty-candidate tests whose diagnostic had been shadowed; after narrowing the
  guard to the cleanup-bound path, the final
  `rtk env GOCACHE=/tmp/roundfix-task01-gocache make verify-incremental` exited
  0 after formatting, vet, all Go package tests, skill checks, and the build.
- The Task's authored `## Verification` command was not run; the Daemon owns
  that check and Task settlement.

### Acceptance-criterion evidence

- **Candidate without evidence is refused and nothing is removed:**
  `TestApplyRunBranchCandidateRefusesACandidateWithoutEvidence` reproduces the
  absent-target candidate, asserts a refusal naming worktree revalidation and
  missing evidence, and verifies that both the Run Worktree and Run Branch
  remain.
- **Candidate with positive evidence is cleaned up:**
  `TestApplyRunBranchCandidateRevalidatesProofAndCleanWorktree` supplies fresh
  positive evidence, applies the candidate, and verifies that its Run Worktree
  and Run Branch are removed while the current Run's surfaces remain.
