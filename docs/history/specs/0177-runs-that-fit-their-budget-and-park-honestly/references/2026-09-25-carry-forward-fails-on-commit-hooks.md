---
type: fix
status: promoted
created: 2026-09-25
spec: 0177-runs-that-fit-their-budget-and-park-honestly
reason: null
---

# Carry-forward fails on repository commit hooks in its dependency-less staging worktree

## Symptom

Carry-forward cherry-picks and amends in a fresh staging worktree without installed dependencies, so Node-based hooks fail (`ERR_MODULE_NOT_FOUND`); reconcile prints the raw hook output and suggests a dry run. Fluxus lost 13 agent-minutes redoing completed Tasks.

## Where

`internal/cli/carryforward.go`, `internal/cli/reconcile.go` (git config for staging commits).

## Expected

Decide: bypass hooks for Roundfix-internal staging commits, or classify the refusal like the Daemon's `HookRefusalError` and name a recovery.

## Evidence

secondbrain `inbox/roundfix/2026-09-10-events-aborta-e-carry-forward-quebra-no-worktree.md` (defect 2).
