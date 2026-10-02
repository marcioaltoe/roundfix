---
task: task_06
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: pending
type: backend
complexity: medium
---

# Task 06: Apply releases a merged Run whose target branch is gone

## Overview

Corrective Task for QA finding F1 of the 2026-10-02 QA Report (row 5, Requirement 5). After a squash merge deletes the target branch, `roundfix reconcile` classifies the Run superseded through the merged-head proof, but `reconcile --apply` does not release it. `applyTerminalRunWithStore` in `internal/worktree/worktree.go` passes the empty `TargetHead` of the absent branch into the Store's `IntegrationReconciliation`, and `validateIntegrationReconciliation` in `internal/store/store.go` refuses it with `target head is required`. The proof head that justified the classification never reaches the Store.

## Requirements

1. MUST, when a terminal Run is classified superseded by a merged-head proof while its target branch is absent, send the Store the proof it was classified on: the merged head as `TargetHead`, and as `TargetBranch` the branch or ref the merged head was read from. A Run whose target branch is present keeps its current request unchanged.
2. MUST NOT relax `validateIntegrationReconciliation` for any other request: an empty target head without a merged-head proof is still refused.
3. MUST add the tests named in Verification to the new file `internal/worktree/merged_spec_absent_target_test.go`, over real Git repositories and the real Store, as `internal/worktree/worktree_test.go` already does:
   - `TestApplyReleasesAMergedRunWhoseTargetBranchIsGone`: dry-run reports superseded, and apply removes the Run Worktree and the Run Branch and records the reconciliation with the merged head.
   - `TestApplyKeepsAnUndeclaredLeftoverWhenTheTargetIsGone`: a Run with an undeclared leftover stays preserved, and nothing is removed.
   - `TestApplyStillRefusesAnEmptyTargetHeadWithoutProof`: an empty target head without a merged-head proof is still refused.
4. MUST leave every existing `internal/worktree` and `internal/store` test unedited and passing.
5. MUST prove each new test fails before the fix and passes after, and record that in the Result.

## Subtasks

- [ ] Carry the merged-head proof into the Store request for an absent target.
- [ ] Add the three tests over real repositories and the real Store.

## Acceptance Criteria

- [ ] `reconcile --apply` releases a superseded Run whose target branch is gone.
- [ ] An undeclared leftover is still preserved, and an empty target head without proof is still refused.

## Context

- interface: `internal/worktree/worktree.go`
- interface: `internal/worktree/merged_head.go`
- interface: `internal/store/store.go`
- creates: `internal/worktree/merged_spec_absent_target_test.go`
- instruction: `internal/worktree/worktree_test.go`
- instruction: `internal/worktree/merged_spec_leftovers_test.go`
- instruction: `docs/adr/0212-a-merged-spec-supersedes-its-runs-spec-directory-work-and-declared-leftovers.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestApplyReleasesAMergedRunWhoseTargetBranchIsGone|TestApplyKeepsAnUndeclaredLeftoverWhenTheTargetIsGone|TestApplyStillRefusesAnEmptyTargetHeadWithoutProof)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestApplyReleasesAMergedRunWhoseTargetBranchIsGone TestApplyKeepsAnUndeclaredLeftoverWhenTheTargetIsGone TestApplyStillRefusesAnEmptyTargetHeadWithoutProof; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0
- `out="$(go test -count=1 ./internal/worktree ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; test -f internal/worktree/merged_spec_absent_target_test.go` — expected: exit 0

## References

- QA Report 2026-10-02 → F1; row 5
- task_04 → Requirement 5
- ADR-0212; ADR-0053; ADR-0161
