---
task: task_01
spec: 0154-evidence-when-the-target-branch-is-gone
status: completed
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

## Result

### Implementation

- The absent-target branch-set classifier now resolves the repository default
  branch and applies `supersedingQAReport` against its head for each terminal
  Run Branch.
- Positive evidence records the same releasable proof used by the
  present-target path. An unresolved default branch, unresolved Run Branch,
  Active Run, or missing evidence keeps the prior preserved classification.
- Default-branch resolution is shared with terminal Run inspection so both
  reconciliation paths use the same detection and ambiguity rules.

### Acceptance evidence

- `TestReconcileFallsBackToTheDefaultBranch` passed after target-branch work was
  squash-merged into the default branch, the target was deleted, and a newer
  default-branch QA Report was added; the classification named that report as
  its releasable proof.
- `TestReconcilePreservesWithoutDefaultBranchEvidence/default_branch_has_no_superseding_evidence`
  passed and produced no releasable candidate.
- `TestReconcilePreservesWithoutDefaultBranchEvidence/default_branch_cannot_be_resolved`
  passed and produced no releasable candidate.
- `TestReconcileFallsBackToTheDefaultBranch` also dirtied the registered
  Worktree after classification; apply refused the candidate as `dirty` and
  left both the Worktree and Run Branch present.
- The focused present-target compatibility run passed 13 safe, superseded,
  unintegrated, dirty, ambiguous-target and existing branch-set assertions.

### Focused checks

- Initial regression command:
  `rtk go test -count=1 ./internal/worktree -run 'TestReconcile(FallsBackToTheDefaultBranch|PreservesWithoutDefaultBranchEvidence)$'`
  — the positive case failed on the old absent-target early return while the
  three negative/subtest cases passed.
- After implementation:
  `rtk go test -count=1 ./internal/worktree -run 'Test(Reconcile(FallsBackToTheDefaultBranch|PreservesWithoutDefaultBranchEvidence)|ClassifyRunBranchSetPreservesAbsentTarget|ApplyRunBranchCandidatePreservesNewlyDirtyWorktree|InspectTerminalRunUnknownWhenDeletedTargetDefaultBranchCannotBeResolved)$'`
  — 7 tests passed.
- Present-target compatibility:
  `rtk go test -count=1 ./internal/worktree -run 'Test(InspectTerminalRun(Safe|Unintegrated|ClassifiesSupersededQAReport|Dirty)|ClassifyRunBranchSetPreservesAbsentTarget)$'`
  — 13 tests passed.
- Repository incremental gate: `rtk make verify-incremental` passed, including
  formatting, vet, all Go tests, skill checks and the build.
- `rtk git diff --check` passed.
- The two commands under `## Verification` were not run; the Daemon owns them.

### Follow-up

- Task 02 owns preserved reason text that names the missing default-branch
  proof. This slice retains the existing absent-target reason for every
  unproven case.
