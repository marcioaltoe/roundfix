---
task: task_01
spec: 0038-terminal-run-worktree-reconciliation
status: completed
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

- [x] Add the five reconciliation states and result contract.
- [x] Inspect cleanliness including untracked files.
- [x] Resolve recorded Run and target branch tips.
- [x] Prove ancestry through Git porcelain.
- [x] Classify missing paths, refs, and metadata conservatively.
- [x] Add real-repository fixtures for all classifications.

## Acceptance Criteria

- [x] Every terminal fixture produces exactly one documented state.
- [x] A clean reachable branch is `safe`; a clean divergent branch is
      `unintegrated`.
- [x] Any tracked or untracked worktree change produces `dirty`.
- [x] Missing target metadata, ambiguous refs, or Git inspection uncertainty
      produces `unknown` and preserves work.
- [x] Absence of both retained worktree and Run Branch produces `released`.
- [x] Output contains the recorded paths, both heads when known, and a bounded
      deterministic reason.
- [x] No input can make inspection read outside the recorded Git root.

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

## Result

Implemented one read-only terminal Run classifier in `internal/worktree`. It
derives the five reconciliation states from the recorded Git root, Run
Worktree, deterministic Run Branch, target branch, porcelain cleanliness, and
commit ancestry. It resolves refs only after validating local-branch metadata,
rejects ambiguous short refs, and runs worktree status only against a path
registered by the recorded repository.

### Verification

- Red baseline:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task01-gocache go test ./internal/worktree -run 'TestInspectTerminalRun' -count=1`
  failed to compile because the classifier contract did not exist.
- Focused classifications:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task01-gocache go test ./internal/worktree -run 'TestInspectTerminalRun.*(Safe|Unintegrated|Dirty|Unknown|Released|UnsafePath)' -count=1`
  passed.
- Race check:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task01-gocache go test -race ./internal/worktree -run 'TestInspectTerminalRun' -count=1`
  passed, including eight concurrent inspections of one real Git fixture.
- Package regression:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task01-gocache go test ./internal/worktree -count=1`
  passed.
- Scope and whitespace:
  `rtk git -c core.fsmonitor=false diff --check` passed; only this Task file
  and its two declared `internal/worktree` interfaces are changed.

The writable `GOCACHE` wrapper was required because the sandbox denied the
default Go build cache. The Daemon remains responsible for running the two
declared Verification commands verbatim.

### Acceptance evidence

- Every classification is covered by a real temporary Git repository and the
  result contract permits only `safe`, `unintegrated`, `dirty`, `unknown`, or
  `released`.
- `TestInspectTerminalRunSafe`,
  `TestInspectTerminalRunSafeWithoutWorktree`, and
  `TestInspectTerminalRunUnintegrated` prove cleanliness plus positive
  ancestry and divergence behavior.
- `TestInspectTerminalRunDirty` proves tracked and untracked changes produce
  `dirty`, including when target metadata is missing.
- `TestInspectTerminalRunUnknownMissingTarget`,
  `TestInspectTerminalRunUnknownMissingRunBranch`, and
  `TestInspectTerminalRunUnknownAmbiguousRef` prove missing or ambiguous
  evidence remains `unknown`.
- `TestInspectTerminalRunReleased` removes both the registered worktree and
  Run Branch and proves the result is `released`.
- Every fixture checks the recorded Run ID, outcome, Run Worktree, Run Branch,
  target branch, known heads, and one deterministic single-line reason bounded
  to 160 bytes.
- `TestInspectTerminalRunUnsafePath` proves unclean paths and final or ancestor
  symlinks cannot redirect inspection. Git status runs only in the
  repository's registered Run Worktree path.

### Follow-ups

Safe cleanup, stale-head revalidation, automatic reaper migration, store
events, CLI behavior, and Run listing remain in their later Spec Tasks.
