---
task: task_06
spec: 0154-evidence-when-the-target-branch-is-gone
status: completed
type: backend
complexity: medium
---

# Task 06: The same proof on every cleanup path

## Overview

Task 05 put the archive requirement on the branch-set fallback and left two
other cleanup paths without it. Confirmed in source: `InspectTerminalRun`
carries no archive check at all, and `PruneTerminalReport` has the same
unguarded shape.

So a deleted target whose clean Run tree is already represented on the default
branch still returns `safe`, and `--apply` removes the Worktree and Branch
without archived QA proof. Core Feature 3 forbids exactly that, and a Run
already classified safe skips the later branch-set result where Task 05's guard
lives.

Review also found that the preserved reason can lose the very clause it exists
to carry: `boundedReconciliationReason` truncates at 160 bytes, and a long Spec
slug formatted ahead of the proof clause can push `default branch` or `Spec is
not archived` past the cut.

## Requirements

1. MUST apply the archived-evidence requirement to every deleted-target cleanup
   path, including `InspectTerminalRun` and `PruneTerminalReport`, not only the
   branch-set fallback.
2. MUST preserve a Run whose target is absent and whose Spec is not archived,
   on every one of those paths.
3. MUST leave present-target behavior unchanged on all of them.
4. MUST keep the proof clause in the preserved reason when the Spec slug is long
   enough to reach the length bound, abbreviating the slug rather than the
   explanation.
5. MUST keep the reason within its existing length bound.

## Subtasks

- [ ] Apply the archive requirement on the terminal-run and prune paths.
- [ ] Reserve room for the proof clause when bounding the reason.
- [ ] Add a test per path and a test for the long-slug reason.

## Acceptance Criteria

- [ ] A deleted target with a represented tree and no archived proof returns
      preserved from `InspectTerminalRun`, not safe.
- [ ] The prune path preserves the same Run.
- [ ] Present-target classifications are unchanged on both paths.
- [ ] A slug long enough to hit the bound still yields a reason naming the
      missing proof.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestInspectTerminalRunRequiresArchivedEvidence" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestInspectTerminalRunRequiresArchivedEvidence"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestReconcileReasonKeepsProofClauseForLongSlugs" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReconcileReasonKeepsProofClauseForLongSlugs"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Core Features

## Result

### Implementation

- Deleted-target content reconciliation now requires a QA Report under the
  archived Spec path before it classifies a represented Run tree as `safe`.
  Without that proof, direct inspection returns `unintegrated`, so pruning
  preserves the Run Worktree and Run Branch.
- Present-target inspection remains ancestry-based and unchanged; the archive
  requirement applies only after the recorded target branch disappears.
- Missing-proof reasons now reserve the full evidence clause and abbreviate a
  long Spec slug with an ellipsis, while retaining the existing 160-byte bound.

### Focused checks

- Before the production change,
  `rtk env GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 -run 'Test(InspectTerminalRunRequiresArchivedEvidence|PruneTerminalReportRequiresArchivedEvidence|ReconcileReasonKeepsProofClauseForLongSlugs)$' ./internal/worktree`
  exited 1: direct inspection returned `safe`, pruning removed the Run
  surfaces, and the long-slug reason omitted `Spec is not archived`.
- After the production change,
  `rtk env GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 -run 'Test(InspectTerminalRunRequiresArchivedEvidence|PruneTerminalReportRequiresArchivedEvidence|ReconcileReasonKeepsProofClauseForLongSlugs|InspectTerminalRunSafeWhenTargetDeletedAfterSquashMerge|PruneTerminalReconciliationReachableChangedBranch)$' ./internal/worktree`
  exited 0.
- `rtk env GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 ./internal/worktree`
  exited 0.
- The first sandboxed `rtk env GOCACHE=/tmp/roundfix-task06-gocache make verify-incremental`
  run reached the full suite but failed because two process-owner integration
  tests could not read the host process table. Re-running
  `rtk make verify-incremental` with process-table permission exited 0 after
  vet, all Go package tests, skill checks, and the build.
- The Task's authored `## Verification` commands were not run; the Daemon owns
  those checks and Task settlement.

### Acceptance-criterion evidence

- **Deleted target preserves without archived proof:**
  `TestInspectTerminalRunRequiresArchivedEvidence/deleted_target_with_active_Spec_evidence_preserves`
  asserts `unintegrated` and a reason naming both the default branch and the
  missing archive state.
- **Prune preserves the same Run:**
  `TestPruneTerminalReportRequiresArchivedEvidence` asserts no pruned entry and
  verifies that both the Run Worktree and Run Branch remain.
- **Present-target behavior is unchanged:**
  `TestInspectTerminalRunRequiresArchivedEvidence/present_target_without_archived_Spec_evidence_stays_safe`
  keeps direct inspection `safe`, and
  `TestPruneTerminalReconciliationReachableChangedBranch` keeps the existing
  present-target prune behavior.
- **Long-slug reason keeps the proof:**
  `TestReconcileReasonKeepsProofClauseForLongSlugs` asserts that the bounded
  reason still names `default branch` and `Spec is not archived` and remains at
  most 160 bytes.
