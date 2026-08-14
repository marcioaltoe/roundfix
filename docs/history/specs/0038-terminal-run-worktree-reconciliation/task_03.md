---
task: task_03
spec: 0038-terminal-run-worktree-reconciliation
status: completed
type: data
complexity: high
---

# Task 03: Reconcile Integration Pending with durable evidence

## Overview

Connect fresh `safe` Git proof to the guarded Integration Pending transition
without weakening ordinary terminal immutability. The Run becomes Clean and
records its reconciliation evidence transactionally before any retained Git
surface is released; other outcomes retain their stored state.

## Requirements

1. MUST build reconciliation input from the freshly revalidated Run and target
   heads.
2. MUST invoke the guarded Integration Pending to Clean store operation before
   Git cleanup.
3. MUST record prior outcome, current outcome, classification, both branches,
   both heads, worktree, and action in one durable Run Event.
4. MUST leave every non-Integration-Pending terminal outcome unchanged while
   recording safe cleanup evidence.
5. MUST refuse stale, incomplete, or mismatched reconciliation input.
6. MUST perform no Git cleanup when durable reconciliation persistence fails.
7. MUST make repeated reconciliation and cleanup idempotent.

## Subtasks

- [x] Map classifier evidence into the guarded store request.
- [x] Journal complete reconciliation evidence transactionally.
- [x] Preserve non-Integration-Pending outcomes.
- [x] Order durable state before Git cleanup.
- [x] Add stale and incomplete-evidence refusal cases.
- [x] Add repeat reconciliation coverage.

## Acceptance Criteria

- [x] A safe Integration Pending Run becomes Clean with one evidence event.
- [x] The event contains both outcomes, branches, heads, worktree, and action.
- [x] Unresolved, Failed, Stopped, and other terminal outcomes remain unchanged
      after safe cleanup.
- [x] Missing or stale evidence produces no state or Git mutation.
- [x] Database failure starts no worktree or branch removal.
- [x] Repeating the operation produces no duplicate transition or evidence.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/store/store_test.go`
- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/worktree/worktree_test.go`

## Verification

- `rtk go test ./internal/store -run 'TestReconcileIntegration.*(Safe|Stale|Invalid|Idempotent|Outcome)' -count=1`
  — expected: only Integration Pending with complete fresh evidence becomes
  Clean transactionally.
- `rtk go test ./internal/worktree -run 'TestApplyTerminalRun.*(IntegrationPending|StoreFailure|Outcome)' -count=1`
  — expected: durable reconciliation precedes Git cleanup and other outcomes
  remain unchanged.

## References

- `_prd.md` → Goal 4; User Story 5; Core Feature 6; Success Metrics.
- `_techspec.md` → Data Models; Integration Points: Run Database; Build Order
  3.
- `../0037-terminal-outcome-integrity/_techspec.md` → guarded terminal
  reconciliation boundary.
- `../../adr/0052-run-completion-is-compare-and-set.md` → sole terminal
  transition exception.

## Result

Implemented the durable reconciliation barrier shared by explicit terminal Run
apply and automatic terminal reaping. Freshly revalidated `safe` evidence is
mapped into one guarded Store request before Git mutation. The Store validates
the recorded Run outcome, deterministic Run Branch, target branch, and
worktree; promotes only Integration Pending to Clean; preserves every other
terminal outcome; and appends the complete cleanup proof in the same
transaction. Exact evidence replay is a no-op.

### Verification

- Red baseline:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/store -run 'TestReconcileIntegration.*(Safe|Stale|Invalid|Idempotent|Outcome)' -count=1`
  failed to compile because the Store request lacked prior outcome,
  classification, Run Branch, worktree, and action evidence.
- Store reconciliation family:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/store -run 'TestReconcileIntegration' -count=1`
  passed.
- Task-focused worktree reconciliation:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/worktree -run 'TestApplyTerminalRun.*(IntegrationPending|StoreFailure|Outcome)' -count=1`
  passed.
- Worktree package regression:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/worktree -count=1`
  passed.
- Automatic-reaper integration:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/cli -run 'TestRunImplementPreflight.*Terminal.*(Safe|Unique|Reachable)' -count=1`
  passed.
- CLI compile check:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/cli -run '^$' -count=1`
  passed.
- The broad `internal/store` package run reached unrelated sandbox failures
  because `/bin/ps` is not permitted in owner-process tests; the complete
  reconciliation test family above passed. The Daemon remains responsible for
  running the two declared Verification commands verbatim.

### Acceptance evidence

- `TestApplyTerminalRunIntegrationPendingPersistsBeforeCleanup` proves a fresh
  safe Integration Pending Run becomes Clean, records exactly one event before
  cleanup, removes both Git surfaces, and produces no duplicate event on
  repeat.
- `TestReconcileIntegrationSafeCompleteEvidence` decodes the durable payload
  and checks prior/current outcomes, classification, Run Branch and head,
  target branch and head, worktree, and action.
- `TestReconcileIntegrationOutcomePreservedWithEvidence` and
  `TestApplyTerminalRunOutcomeRemainsUnchanged` cover Unresolved, Failed,
  Stopped, and Timed Out outcomes: cleanup evidence is recorded while the Run
  row remains unchanged.
- `TestReconcileIntegrationInvalidCompleteEvidence`,
  `TestReconcileIntegrationRejectsStaleTargetBranch`, and the existing
  stale-head, changed-metadata, and newly-dirty worktree tests prove incomplete,
  mismatched, or stale proof changes neither durable outcome nor Git surface.
- `TestApplyTerminalRunStoreFailureStartsNoCleanup` injects persistence failure
  after fresh revalidation and observes zero worktree/branch mutation calls.
- `TestReconcileIntegrationIdempotent` and the repeated real-Git apply in
  `TestApplyTerminalRunIntegrationPendingPersistsBeforeCleanup` prove replay
  performs no duplicate transition, event, worktree removal, or branch
  deletion.

### Follow-ups

The public Reconcile Command, result rendering, and selector behavior remain
Task 04.
