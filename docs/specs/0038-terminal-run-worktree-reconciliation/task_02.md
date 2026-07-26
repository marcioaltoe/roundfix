---
task: task_02
spec: 0038-terminal-run-worktree-reconciliation
status: pending
type: backend
complexity: high
---

# Task 02: Apply stale-proof cleanup and migrate terminal reaping

## Overview

Apply reconciliation only after revalidating the evidence immediately before
mutation, then move automatic terminal reaping onto the same classifier. This
slice removes safe or already released residue while preserving any work whose
cleanliness, heads, or ancestry changed after inspection.

## Requirements

1. MUST accept cleanup only for a result previously classified `safe`.
2. MUST re-resolve cleanliness, Run head, and target head immediately before
   the first mutation.
3. MUST refuse apply when either head, worktree state, or recorded metadata
   changed after inspection.
4. MUST remove the Run Worktree without force before deleting its Run Branch.
5. MUST preserve every remaining path or ref when a mutation step fails.
6. MUST make repeated cleanup report `released` without another mutation.
7. MUST replace creation-base terminal reaping with the shared classifier and
   permit automatic cleanup only for `safe` or `released`.

## Subtasks

- [ ] Add stale-proof pre-mutation revalidation.
- [ ] Apply worktree removal before branch deletion.
- [ ] Preserve partial-failure evidence and remaining resources.
- [ ] Make repeated apply idempotent.
- [ ] Route automatic terminal reaping through the classifier.
- [ ] Reproduce reachable-later-merge and unique-branch regressions.

## Acceptance Criteria

- [ ] Only a still-clean result with unchanged heads reaches Git mutation.
- [ ] A stale head or newly dirty worktree is refused without deleting a path
      or ref.
- [ ] Worktree removal uses no force option.
- [ ] Branch deletion occurs only after successful worktree removal.
- [ ] A failed removal or deletion reports the remaining recoverable surface.
- [ ] A second apply performs zero mutations and reports `released`.
- [ ] Automatic reaping removes a changed branch already reachable from target
      and preserves a unique changed branch.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/worktree/worktree.go`
- interface: `internal/worktree/worktree_test.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/implement_test.go`

## Verification

- `rtk go test ./internal/worktree -run 'TestApplyTerminalRun.*(Safe|Stale|Dirty|Failure|Released)|TestPruneTerminal.*Reconciliation' -count=1`
  — expected: only fresh safe evidence mutates Git and repeat apply is
  idempotent.
- `rtk go test ./internal/cli -run 'TestRunImplementPreflight.*Terminal.*(Safe|Unique|Reachable)' -count=1`
  — expected: automatic reaping shares the proof-based classifier.
- `rtk go test -race ./internal/worktree -run 'Test(ApplyTerminalRun|PruneTerminal)' -count=1`
  — expected: inspection-to-apply boundaries remain race-free.

## References

- `_prd.md` → Goals 1–2; User Stories 1 and 4; Core Features 4, 7, and 9;
  Success Metrics.
- `_techspec.md` → API Contracts; Testing Approach; Build Order 2.
- `../../adr/0023-runs-execute-in-per-run-worktrees.md` → Run Worktree
  ownership.
- `../../adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md` →
  apply safety.
