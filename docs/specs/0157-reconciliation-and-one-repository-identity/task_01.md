---
task: task_01
spec: 0157-reconciliation-and-one-repository-identity
status: pending
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
