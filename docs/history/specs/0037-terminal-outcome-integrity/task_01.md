---
task: task_01
spec: 0037-terminal-outcome-integrity
status: completed
type: data
complexity: high
---

# Task 01: Guard terminal completion and reconciliation

## Overview

Make terminal Run completion an atomic first-writer-wins operation while
preserving identical replay and the one evidence-backed Integration Pending
reconciliation. This slice is complete when competing completion attempts
cannot overwrite the stored outcome, release its lock twice, or publish
conflicting durable evidence.

## Requirements

1. MUST make ordinary Run completion update only a non-terminal Run.
2. MUST return an explicit transitioned result for the winning completion and
   an idempotent non-transition result for an identical replay.
3. MUST return a typed conflict for a different requested terminal outcome
   without changing the Run, completion timestamp, Active Run lock, or journal.
4. MUST release the Active Run lock only in the transaction that wins ordinary
   terminal completion.
5. MUST provide the sole guarded Integration Pending to Clean reconciliation,
   requiring the recorded Run and target heads and journaling both outcomes.
6. MUST reject every other terminal-to-terminal transition.
7. MUST wrap database and transaction failures with the operation and Run ID.

## Subtasks

- [x] Add explicit completion and conflict result contracts.
- [x] Implement compare-and-set terminal persistence and winner-only lock release.
- [x] Implement evidence-backed Integration Pending reconciliation.
- [x] Add idempotent replay and conflicting-outcome coverage.
- [x] Add deterministic concurrent-completion coverage.
- [x] Prove reconciliation rejects stale or invalid source outcomes.

## Acceptance Criteria

- [x] One non-terminal completion stores its outcome and reports transitioned.
- [x] Repeating the same outcome changes no persisted field or journal entry.
- [x] A competing outcome returns the typed conflict and leaves the winner,
      timestamp, lock state, and event history unchanged.
- [x] A deterministic race produces exactly one winner and one stable terminal
      row.
- [x] Integration Pending becomes Clean only with complete reconciliation
      evidence recorded transactionally.
- [x] No other terminal outcome can be rewritten.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/store/store_test.go`
- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`

## Verification

- `rtk go test ./internal/store -run 'Test(CompleteRun|TerminalOutcome|ReconcileIntegration)' -count=1`
  — expected: winner, replay, conflict, lock, and guarded reconciliation cases
  pass.
- `rtk go test -race ./internal/store -run 'Test(CompleteRun|TerminalOutcome)' -count=1`
  — expected: concurrent competing completion has one race-free winner.

## References

- `_prd.md` → Goals 1; User Stories 1, 3, and 5; Core Features 1–2 and 7;
  Success Metrics.
- `_techspec.md` → Interfaces; Data Models; Build Order 1.
- `../../adr/0052-run-completion-is-compare-and-set.md` → terminal
  compare-and-set and guarded reconciliation.

## Result

The store slice meets every acceptance criterion, and the Settle Command was
adapted to the Spec's no-rewrite decision: Settle never rewrites a settled
Run's terminal outcome, so a recovered `Unresolved` Run keeps `Unresolved` as
history and recovery is expressed only through the Settle report and its
commits. The repository regression suite is green.

### What changed

- Ordinary completion now returns an explicit transition result, atomically
  changes only a non-terminal Run, and releases the Active Run lock only for
  the winner.
- Identical replay returns the stored Run without changing timestamps, locks,
  or Run Events. A competing terminal outcome returns
  `TerminalOutcomeConflictError` with the stored and requested outcomes.
- Integration Pending reconciliation validates complete head and target
  evidence, compare-and-sets the Run to Clean, and records both outcomes and
  the evidence in one transaction.
- Existing callers were adapted only to unwrap the explicit completion result;
  winner-only event and notification publication remains Task 05's slice.
- `integrateSettledRun` in `internal/cli/settle.go` no longer calls
  `CompleteRun`: it still returns the integration command on
  `runworktree.ModePending` and still cleans up the Run Worktree on clean
  integration, but leaves the Run's stored terminal outcome untouched.
- `TestRunSettleRetargetsKeptRunWorktreeAndCleansUpAfterIntegration` and
  `TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration` now
  assert the settled Run's stored state remains `Unresolved`; all worktree
  cleanup, checkout, and commit assertions are unchanged.

### Verification

- Focused store check:
  `GOFLAGS=-buildvcs=false go test ./internal/store -run
  'Test(CompleteRun|TerminalOutcome|ReconcileIntegration)' -count=1` passed
  (`ok roundfix/internal/store 0.434s`).
- Race check:
  `GOCACHE=/private/tmp/roundfix-task01-gocache go test -race ./internal/store
  -run 'Test(CompleteRun|TerminalOutcome)' -count=1` passed.
- Settle recovery check:
  `GOFLAGS=-buildvcs=false go test ./internal/cli -run
  '^TestRunSettleRetargetsKept(Run|Task)WorktreeAndCleansUpAfterIntegration$'
  -count=1` passed (`ok roundfix/internal/cli 1.004s`).
- Repository regression check:
  `GOFLAGS=-buildvcs=false go test ./...` passed for every package; no
  failures.
- `gofmt -l internal/cli/settle.go internal/cli/settle_test.go` reported no
  files.

### Acceptance evidence

- `TestCompleteRunWinnerAndIdenticalReplay` proves the first completion reports
  `Transitioned: true`, releases the lock once, and an identical replay reports
  no transition without changing the row or journal.
- `TestTerminalOutcomeConflictPreservesWinner` proves the typed conflict keeps
  the winning outcome, completion timestamp, lock count, and event count.
- `TestCompleteRunConcurrentTerminalOutcomesHaveOneWinner` plus the race check
  proves exactly one competing completion wins and the stored terminal row is
  stable.
- `TestReconcileIntegrationPendingRecordsEvidence` and
  `TestReconcileIntegrationRollsBackWhenJournalFails` prove the guarded
  reconciliation and its Run Event commit or roll back together.
- `TestReconcileIntegrationRejectsIncompleteEvidence`,
  `TestReconcileIntegrationRejectsStaleTargetBranch`, and
  `TestReconcileIntegrationRejectsEveryOtherSourceOutcome` prove incomplete,
  stale, and invalid source evidence changes neither the Run nor its journal.
- `TestTerminalOutcomeEveryStoredTerminalStateIsImmutable` proves ordinary
  completion cannot rewrite any terminal outcome, including Integration
  Pending.
- The two updated Settle tests prove a recovered Run keeps `Unresolved` as its
  stored outcome while Settle still integrates the work, cleans up the kept
  worktrees, and releases the Active Run lock.
