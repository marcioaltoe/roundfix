---
task: task_01
spec: 0037-terminal-outcome-integrity
status: pending
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

- [ ] Add explicit completion and conflict result contracts.
- [ ] Implement compare-and-set terminal persistence and winner-only lock release.
- [ ] Implement evidence-backed Integration Pending reconciliation.
- [ ] Add idempotent replay and conflicting-outcome coverage.
- [ ] Add deterministic concurrent-completion coverage.
- [ ] Prove reconciliation rejects stale or invalid source outcomes.

## Acceptance Criteria

- [ ] One non-terminal completion stores its outcome and reports transitioned.
- [ ] Repeating the same outcome changes no persisted field or journal entry.
- [ ] A competing outcome returns the typed conflict and leaves the winner,
      timestamp, lock state, and event history unchanged.
- [ ] A deterministic race produces exactly one winner and one stable terminal
      row.
- [ ] Integration Pending becomes Clean only with complete reconciliation
      evidence recorded transactionally.
- [ ] No other terminal outcome can be rewritten.

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
