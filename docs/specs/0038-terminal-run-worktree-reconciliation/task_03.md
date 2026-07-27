---
task: task_03
spec: 0038-terminal-run-worktree-reconciliation
status: pending
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

- [ ] Map classifier evidence into the guarded store request.
- [ ] Journal complete reconciliation evidence transactionally.
- [ ] Preserve non-Integration-Pending outcomes.
- [ ] Order durable state before Git cleanup.
- [ ] Add stale and incomplete-evidence refusal cases.
- [ ] Add repeat reconciliation coverage.

## Acceptance Criteria

- [ ] A safe Integration Pending Run becomes Clean with one evidence event.
- [ ] The event contains both outcomes, branches, heads, worktree, and action.
- [ ] Unresolved, Failed, Stopped, and other terminal outcomes remain unchanged
      after safe cleanup.
- [ ] Missing or stale evidence produces no state or Git mutation.
- [ ] Database failure starts no worktree or branch removal.
- [ ] Repeating the operation produces no duplicate transition or evidence.

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
