---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# A killed carry-forward leaves a locked staging worktree nothing sweeps

## Symptom

Carry-forward stages in `$TMPDIR/roundfix-carry-forward-*`; a process killed during `git worktree add` leaves a `locked initializing` admin entry and directory that the single `--force` cleanup cannot remove and nothing sweeps. Tests also write staging dirs into the global TMPDIR.

## Where

`internal/cli/carryforward.go` staging creation and cleanup.

## Expected

Sweep stale staging worktrees on the next carry-forward or reconcile (`git worktree prune`, `remove --force --force`), and keep test staging under `t.TempDir()`.

## Evidence

Session 2026-09-25: `$TMPDIR/roundfix-carry-forward-2309636395` locked since 2026-09-24 14:57.
