---
task: task_01
spec: 0038-terminal-run-worktree-reconciliation
status: pending
type: backend
complexity: high
---

# Task 01: Classify terminal Run Worktrees with positive proof

## Overview

Add one conservative classifier for terminal spec Run Worktrees and Run
Branches. It derives current reconciliation state from recorded Run metadata
and real Git evidence so every consumer receives the same safe, unintegrated,
dirty, unknown, or released result.

## Requirements

1. MUST inspect terminal spec Runs from their recorded Git root, worktree,
   Run Branch, and target branch metadata.
2. MUST classify `safe` only when the retained worktree is clean and the Run
   Branch tip is an ancestor of the current target tip.
3. MUST classify tracked or untracked worktree changes as `dirty`.
4. MUST distinguish `unintegrated`, `unknown`, and `released` without treating
   missing proof as safe.
5. MUST report both resolved heads and one bounded reason.
6. MUST reject symlink, path, or metadata input that could escape the recorded
   repository boundary.
7. MUST use real Git fixtures for every state and ambiguous-ref failure.

## Subtasks

- [ ] Add the five reconciliation states and result contract.
- [ ] Inspect cleanliness including untracked files.
- [ ] Resolve recorded Run and target branch tips.
- [ ] Prove ancestry through Git porcelain.
- [ ] Classify missing paths, refs, and metadata conservatively.
- [ ] Add real-repository fixtures for all classifications.

## Acceptance Criteria

- [ ] Every terminal fixture produces exactly one documented state.
- [ ] A clean reachable branch is `safe`; a clean divergent branch is
      `unintegrated`.
- [ ] Any tracked or untracked worktree change produces `dirty`.
- [ ] Missing target metadata, ambiguous refs, or Git inspection uncertainty
      produces `unknown` and preserves work.
- [ ] Absence of both retained worktree and Run Branch produces `released`.
- [ ] Output contains the recorded paths, both heads when known, and a bounded
      deterministic reason.
- [ ] No input can make inspection read outside the recorded Git root.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/worktree/worktree.go`
- interface: `internal/worktree/worktree_test.go`
- interface: `internal/store/store.go`

## Verification

- `rtk go test ./internal/worktree -run 'TestInspectTerminalRun.*(Safe|Unintegrated|Dirty|Unknown|Released|UnsafePath)' -count=1`
  — expected: real Git fixtures prove all five states and conservative path
  handling.
- `rtk go test -race ./internal/worktree -run 'TestInspectTerminalRun' -count=1`
  — expected: concurrent inspection fixtures are race-free.

## References

- `_prd.md` → Goals 1; User Stories 1 and 4; Core Features 1–3; Success
  Metrics.
- `_techspec.md` → Interfaces; Data Models; Testing Approach; Build Order 1.
- `../../adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md` →
  positive cleanliness and ancestry proof.
