---
task: task_04
spec: 0038-terminal-run-worktree-reconciliation
status: pending
type: backend
complexity: high
---

# Task 04: Deliver the Reconcile Command and apply contract

## Overview

Expose proof-based reconciliation through a deterministic, dry-run-first CLI
for one terminal spec Run or every terminal spec Run in the current repository.
Text and versioned JSON report the same evidence; `--apply` acts only on results
proven safe during the current invocation.

## Requirements

1. MUST implement `roundfix reconcile [run-id] [--apply] [--format text|json]`.
2. MUST reject Active Runs, review Runs, missing Runs, and cross-repository
   selectors before any Git mutation.
3. MUST scope the no-ID form to terminal spec Runs in the current Git root.
4. MUST default to a read-only report and print the exact apply command.
5. MUST render deterministic text and a versioned JSON envelope with ordered
   results and summary counts.
6. MUST apply only entries classified and revalidated `safe` in that
   invocation; no force or assertion bypass is permitted.
7. MUST keep expected preserved states successful while making operational Git
   or database failures non-zero and actionable.

## Subtasks

- [ ] Add Reconcile Command parsing and usage.
- [ ] Resolve one-Run and repository-wide scopes.
- [ ] Render deterministic text and versioned JSON.
- [ ] Connect dry-run and explicit apply behavior.
- [ ] Preserve mixed safe and unsafe scan results.
- [ ] Add idempotent repeat and invalid-selector coverage.

## Acceptance Criteria

- [ ] Default invocation changes neither Git nor the Run Database.
- [ ] One-Run and repository-wide results are ordered newest-first.
- [ ] Text and JSON expose classification, branches, heads, worktree, evidence,
      action, and refusal reason consistently.
- [ ] `--apply` removes all and only fresh `safe` results.
- [ ] Dirty, unintegrated, unknown, and released results are preserved and
      reported without turning a complete scan into failure.
- [ ] Operational inspection or apply failure returns the Run failure exit code
      and names the next safe action.
- [ ] A second apply reports `released` and performs zero mutations.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/runs.go`
- interface: `internal/worktree/worktree.go`

## Verification

- `rtk go test ./internal/cli -run 'TestRunReconcile.*(DryRun|Apply|JSON|Repository|Invalid|Mixed|Idempotent)' -count=1`
  — expected: command scope, output, mutation, refusal, and repeat contracts
  pass.
- `rtk go test -race ./internal/cli ./internal/worktree -run 'Test.*Reconcile' -count=1`
  — expected: CLI inspection and apply coordination are race-free.
- `rtk go build -buildvcs=false ./cmd/roundfix`
  — expected: the public command builds with repository build settings.

## References

- `_prd.md` → Goals 1–2; User Stories 1–2 and 4–5; Core Features 4–7; User
  Experience; Success Metrics.
- `_techspec.md` → API Contracts; Integration Points; Build Order 4.
- `../../adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md` →
  explicit dry-run/apply command.
