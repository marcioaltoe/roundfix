---
status: pending
created_at: 2026-09-25
updated_at: 2026-09-25
---

# Delivery — A squash merge leaves the Spec's Runs, branches and worktrees behind (2026-09-25)

After a day of parallel deliveries (Specs 0155–0168, PRs #238–#248), this
repository held 16 `roundfix/run-*` branches, 16 Run Worktrees under
`~/.roundfix/worktrees/` and a locked carry-forward staging worktree, most of
them for Specs already squash-merged and archived. The maintainer asked that
every finished squash merge be followed by a cleanup of completed Runs and
branches. A read-only `roundfix reconcile --format json` on main 7a9b6ec6
explained each retained item.

## 1. Nothing releases a Spec's Runs after its merge

- Symptom / evidence: `roundfix deliver` records the `merged` stage and stops
  (`internal/delivery/engine.go`, `mergeCandidate`). Automatic pruning runs only
  at `implement` preflight and at `stop`, and releases only Runs classified
  Safe or Released. `roundfix gc` prunes journals and artifacts, never worktrees
  or branches. Manual merges have no release step at all.
- Root cause: no merge-triggered release exists; cleanup depends on a later,
  unrelated command happening to run, and on proofs that fail (below).
- Action / suggestion: a release step after every merge receipt, and a manual
  equivalent for merges done outside the loop. Routed to Backlog Entry
  `docs/backlog/2026-09-25-post-merge-cleanup-of-runs-and-branches.md`.

## 2. The deleted-target proof cannot pass for an archived Spec

- Symptom / evidence: `run_20260924T223414Z_54e6f85ae98b6550` (Spec 0164,
  merged in #246, archived, local branch deleted) is retained with
  `11 Run-only files, 35 differing shared files against default branch "main"`.
  The 11 files are exactly the Spec folder archiving moved; the 35 include every
  later change on `main`. `inspectDeletedTargetRunByContent` compares the whole
  tree with the current `main` yet requires the Spec to be archived.
- Root cause: whole-tree equality against a moving `main`.
- Action / suggestion: prove against the merge commit, or only the Run's
  changed paths with the archive rename mapped. Backlog Entry
  `docs/backlog/2026-09-25-reconcile-proves-archived-deleted-target-runs.md`.

## 3. Nine Runs are invisible to reconcile

- Symptom / evidence: the Runs of Specs 0155 (×4), 0156, 0157, 0160 (×2) and
  0162 have `repository_root = ''` and a `git_root` under removed scratchpad
  worktrees; `ListRuns` never returns them, and the prune guard skips them,
  although their Run Worktrees are registered in this repository.
- Root cause: the schema v18 backfill only resolves `git_root`, never
  `work_dir`.
- Action / suggestion: Backlog Entry
  `docs/backlog/2026-09-25-legacy-run-repository-key-from-its-worktree.md`.

## 4. Superseded Runs survive `--apply`

- Symptom / evidence: `run_20260924T205403Z_4627184d18e75ee8` (Spec 0159,
  merged in #247) is classified superseded, but only `--discard-superseded`
  acts on that classification.
- Root cause: `--apply` keeps its legacy classification by design.
- Action / suggestion: Backlog Entry
  `docs/backlog/2026-09-25-reconcile-apply-ignores-superseded-runs.md`.

## 5. A killed carry-forward leaks a locked staging worktree

- Symptom / evidence: `$TMPDIR/roundfix-carry-forward-2309636395` is
  `locked initializing` since 2026-09-24 14:57; a single `--force` cannot remove
  a locked worktree and nothing sweeps `roundfix-carry-forward-*`.
- Root cause: interrupted `git worktree add` with no later sweep.
- Action / suggestion: Backlog Entry
  `docs/backlog/2026-09-25-carry-forward-staging-worktree-leaks.md`.

## What worked — keep

Runs still Active, an Unresolved Run with uncommitted changes, and a Run
carrying the only QA evidence of an unmerged Spec were retained correctly:
reconcile refused to release anything it could not prove.
