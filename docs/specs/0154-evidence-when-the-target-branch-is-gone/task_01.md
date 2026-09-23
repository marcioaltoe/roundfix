---
task: task_01
spec: 0154-evidence-when-the-target-branch-is-gone
status: pending
type: backend
complexity: medium
---

# Task 01: Ask the default branch when the target is gone

## Overview

`internal/worktree/worktree.go:422` returns as soon as the target branch is
absent, so nothing is checked. That is the state every squash-merged pull
request leaves behind: measured on this repository, 22 retained Run Worktrees,
all preserved for an absent target, and `reconcile --apply` released none.

Ancestry would not have helped either. One Run Branch in twenty-two is an
ancestor of `main`, because a squash merge writes a new commit and leaves the
originals outside the default branch's lineage.

The content-evidence rule that already handles a missed ancestry —
`supersedingQAReport` — is what this Task reaches for.

## Requirements

1. MUST resolve the default branch when the target branch is absent, and apply
   the accepted content-evidence rule against it.
2. MUST classify the Run as it would have been classified with the target
   present, when that evidence proves delivery.
3. MUST reuse `supersedingQAReport` rather than adding a second way to prove
   integration.
4. MUST preserve the Run when the default branch cannot be resolved.
5. MUST NOT release a Run whose Worktree holds uncommitted changes.
6. MUST leave every classification unchanged when the target branch is present.

## Subtasks

- [ ] Resolve the default branch on the absent-target path.
- [ ] Apply the accepted evidence rule against it.
- [ ] Add tests for evidence present, evidence absent and default unresolvable.

## Acceptance Criteria

- [ ] A fixture squash-merged into the default branch with its target deleted is
      classified releasable.
- [ ] The same fixture without that evidence stays preserved.
- [ ] An unresolvable default branch preserves.
- [ ] A dirty Worktree preserves regardless of evidence.
- [ ] Present-target classifications are asserted unchanged.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReconcileFallsBackToTheDefaultBranch" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReconcileFallsBackToTheDefaultBranch"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestReconcilePreservesWithoutDefaultBranchEvidence" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReconcilePreservesWithoutDefaultBranchEvidence"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The fallback
