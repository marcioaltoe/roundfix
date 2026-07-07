---
spec: 0022-cleanup-robustness
status: active
created: 2026-07-07
surfaces: [cli, docs]
---

# Cleanup Robustness

Engineering-framed minimal PRD (bug fix; no product behavior change beyond
removing a failure mode). Field report: tax-poc, 2026-07-07.

## Problem statement

Worktree cleanup uses `git worktree remove` without `--force`, so bootstrap
debris (node_modules, env files, nested repos) makes git refuse — and a
cleanup failure after a fully integrated Clean Run converts the outcome to
Failed with exit 1. A Run whose work is complete, verified, and integrated
reports failure because deleting a directory failed.

## Goals

- Worktree and branch cleanup succeeds in the presence of untracked
  bootstrap debris.
- A cleanup failure after successful integration never changes the Run's
  outcome, report, or exit code: it warns, journals, and keeps the worktree
  for manual removal.
- Operators know the sanctioned recovery for an unsatisfiable task
  Verification: fix the task file's Verification, re-run Settle — never a
  skip-verification flag.

## Core features

1. Every Roundfix-owned worktree removal (Clean cleanup, Task cleanup,
   force-stop reap, preflight sweep) uses `git worktree remove --force`.
2. Cleanup failure after integration degrades: one stderr warning naming the
   kept path, one journaled Daemon event, outcome and exit code unchanged.
3. Docs: Settle recovery via editing the task's Verification; a
   skip-verification flag is explicitly rejected (verification is the only
   gate). The write-tasks skill gains the two authoring rules from the field
   report: task Verification must be hermetic and satisfiable, and
   commit/push never appear in acceptance criteria.

## Non-goals

- Budget default changes (config already covers long Runs).
- Any change to verification, settle, or integration semantics.
