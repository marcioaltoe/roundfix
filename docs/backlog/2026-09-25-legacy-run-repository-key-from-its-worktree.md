---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# Legacy Runs whose original checkout vanished have no repository key and are invisible to reconcile

## Symptom

Migration 18 backfills `repository_root` only from `git_root`; Runs created from since-removed linked worktrees keep an empty key, so `runs list` and `reconcile` from the main checkout never see them although their Run Worktree is still registered there. 26 of 108 Runs on this machine are unkeyed; 9 with existing worktrees.

## Where

`internal/store/store.go` repository-key backfill; `ListRuns`; `worktree.go` prune guard.

## Expected

Derive the key from `work_dir`'s common Git directory when `git_root` is gone (re-runnable backfill or query-time fallback).

## Evidence

secondbrain `inbox/roundfix/2026-09-14-fiscus-run-sem-checkout-proprietario-perde-rota-de-limpeza.md`; live repro `reconcile run_20260924T173241Z_9612fd1f6966c280`.
