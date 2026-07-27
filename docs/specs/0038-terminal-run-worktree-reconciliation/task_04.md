---
task: task_04
spec: 0038-terminal-run-worktree-reconciliation
status: completed
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

- [x] Add Reconcile Command parsing and usage.
- [x] Resolve one-Run and repository-wide scopes.
- [x] Render deterministic text and versioned JSON.
- [x] Connect dry-run and explicit apply behavior.
- [x] Preserve mixed safe and unsafe scan results.
- [x] Add idempotent repeat and invalid-selector coverage.

## Acceptance Criteria

- [x] Default invocation changes neither Git nor the Run Database.
- [x] One-Run and repository-wide results are ordered newest-first.
- [x] Text and JSON expose classification, branches, heads, worktree, evidence,
      action, and refusal reason consistently.
- [x] `--apply` removes all and only fresh `safe` results.
- [x] Dirty, unintegrated, unknown, and released results are preserved and
      reported without turning a complete scan into failure.
- [x] Operational inspection or apply failure returns the Run failure exit code
      and names the next safe action.
- [x] A second apply reports `released` and performs zero mutations.

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

## Result

Implemented the Reconcile Command as a dry-run-first CLI over the shared
proof-based classifier. A positional Run ID and flags work in either order;
the no-ID form queries only terminal spec Runs for the canonical current Git
root. Text and `roundfix-reconcile/v1` JSON carry the same ordered evidence,
actions, refusal reasons, and summary counts. Apply opens the writable Run
Database only when the current scan contains a `safe` result, then delegates
each such result to stale-proof revalidation and cleanup. Expected preserved
states remain successful, while inspection, persistence, cleanup, and output
failures return the Run failure exit code with a rerun command.

Verification:

- `rtk go test ./internal/cli -run 'TestRunReconcile.*(DryRun|Apply|JSON|Repository|Invalid|Mixed|Idempotent)' -count=1`
  — passed, 12 tests.
- `rtk go test -race ./internal/cli ./internal/worktree -run 'Test.*Reconcile' -count=1`
  — passed, 12 tests across both packages.
- `rtk go build -buildvcs=false ./cmd/roundfix` — passed.
- `rtk go test ./internal/cli -count=1` — 779 tests passed; the sole remaining
  sandbox failure was the pre-existing owner-process integration test because
  `/bin/ps` was denied. Running that exact test outside the sandbox passed.
- `rtk git diff --check` — passed.

Acceptance evidence:

- `TestRunReconcileDryRunReadOnly` compares the Run Database bytes and the
  complete Git worktree/ref surface before and after the default invocation;
  both remain identical.
- `TestRunReconcileRepositoryScopeNewestFirst` pins creation timestamps,
  observes newest-first current-repository output, and excludes a newer Run
  belonging to another repository. A selected one-Run report contains exactly
  that Run.
- `TestRunReconcileDryRunReadOnly` and
  `TestRunReconcileJSONMatchesTextFields` cover classification, Run and target
  branches and heads, worktree, evidence, action, refusal reason, schema
  version, mode, repository, apply command, and summary counts.
- `TestRunReconcileApplyMixedResults` creates `safe`, `dirty`, `unintegrated`,
  `unknown`, and `released` real-Git fixtures. Apply removes only the freshly
  revalidated `safe` worktree and branch; all three preserved work surfaces
  remain and the scan exits successfully.
- `TestRunReconcileApplyFailureNamesNextSafeAction` uses a locked clean
  worktree to make real Git removal fail. The command exits with the Run
  failure code, reports one operational failure and retained evidence, and
  prints the exact safe rerun action.
- `TestRunReconcileIdempotentApply` proves the first apply reconciles
  Integration Pending to Clean and releases Git state. The second apply
  reports `released`, performs zero Git or Run Database mutations, and reports
  `applied=0`.

Follow-ups: none.
