---
task: task_02
spec: 0038-terminal-run-worktree-reconciliation
status: completed
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

- [x] Add stale-proof pre-mutation revalidation.
- [x] Apply worktree removal before branch deletion.
- [x] Preserve partial-failure evidence and remaining resources.
- [x] Make repeated apply idempotent.
- [x] Route automatic terminal reaping through the classifier.
- [x] Reproduce reachable-later-merge and unique-branch regressions.

## Acceptance Criteria

- [x] Only a still-clean result with unchanged heads reaches Git mutation.
- [x] A stale head or newly dirty worktree is refused without deleting a path
      or ref.
- [x] Worktree removal uses no force option.
- [x] Branch deletion occurs only after successful worktree removal.
- [x] A failed removal or deletion reports the remaining recoverable surface.
- [x] A second apply performs zero mutations and reports `released`.
- [x] Automatic reaping removes a changed branch already reachable from target
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

## Result

Implemented stale-proof terminal Run cleanup and moved automatic terminal
reaping onto the shared reconciliation classifier. Inspection now seals the
recorded Run metadata and Git proof used by apply. Apply repeats inspection
immediately before mutation, accepts only inspected `safe` evidence or an
already `released` no-op, removes the Run Worktree without force, and deletes
the Run Branch only after worktree removal succeeds. Typed apply failures
report which path or ref remains.

Automatic reaping now loads each recorded terminal Run and removes a Run
Worktree or Run Branch only after the classifier returns `safe`; `released`
permits idempotent task-residue handling without another Run mutation.
Implement and force-stop call sites pass the full Run record to that lookup.

### Verification

- Red baseline:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/worktree -run 'TestApplyTerminalRun.*(Safe|Stale|Dirty|Failure|Released)|TestPruneTerminal.*Reconciliation' -count=1`
  failed to compile because `ApplyTerminalRun`, its revalidation seam, and the
  recorded-Run reaper lookup did not exist.
- Focused apply and reaper:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/worktree -run 'TestApplyTerminalRun.*(Safe|Stale|Dirty|Failure|Released)|TestPruneTerminal.*Reconciliation' -count=1`
  passed.
- Implement Preflight:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/cli -run 'TestRunImplementPreflight.*Terminal.*(Safe|Unique|Reachable)' -count=1`
  passed.
- Race check:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test -race ./internal/worktree -run 'Test(ApplyTerminalRun|PruneTerminal)' -count=1`
  passed.
- Worktree package regression:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/worktree -count=1`
  passed.
- Additional CLI package run:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/cli -count=1`
  reached the unrelated host restriction
  `fork/exec /bin/ps: operation not permitted` in
  `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`;
  the Task's focused CLI command passed.

The writable `GOCACHE` wrapper was required because the sandbox denied the
default Go build cache. The Daemon remains responsible for running the three
declared Verification commands verbatim.

### Acceptance evidence

- `TestApplyTerminalRunSafeEvidenceOnly`, `TestApplyTerminalRunStaleHeads`,
  `TestApplyTerminalRunStaleMetadata`, and `TestApplyTerminalRunDirty` prove
  that only sealed, still-clean `safe` evidence with unchanged heads reaches
  mutation; every stale case preserves the Run Worktree and Run Branch.
- `TestApplyTerminalRunRemovalFailure` records the exact
  `git worktree remove <path>` arguments, proves `--force` is absent, and
  proves no branch deletion follows a failed removal.
- `TestApplyTerminalRunDeletionFailure` records worktree removal before
  `git branch -D` and proves a deletion failure leaves the Run Branch
  recoverable.
- Both failure tests inspect `TerminalRunApplyError` and prove its remaining
  worktree and branch fields match the resources left after the failed step.
- `TestApplyTerminalRunReleased` repeats apply against the original safe
  evidence, observes zero mutation calls, and confirms a fresh inspection
  reports `released`.
- `TestPruneTerminalReconciliationReachableChangedBranch` removes a changed
  Run Branch already reachable from the recorded target;
  `TestPruneTerminalReconciliationPreservesUniqueChangedBranch` preserves a
  divergent changed branch.
- `TestRunImplementPreflightTerminalReachableChangedBranch` and
  `TestRunImplementPreflightTerminalUniqueChangedBranch` reproduce those two
  outcomes through the real Implement Preflight and Run Database lookup.

### Follow-ups

Durable reconciliation events, the explicit Reconcile Command, and terminal
Run listing remain in their later Spec Tasks.
