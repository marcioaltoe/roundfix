# QA-16 — Read-only boundary

Status: pass.

Before and after repeated built audits of the real Spec:

- `/Users/marcio/.roundfix/roundfix.db` retained hash
  `6c22a903d7ebb1b7c0038bab4d24a908ae99b572`.
- `git show-ref` retained every ref.
- `git worktree list --porcelain` retained all three registered worktrees and
  heads.
- `git status --porcelain=v1` retained only this gate's untracked QA directory.

The clean and scratch fixtures likewise retained refs, HEADs, and worktrees.
The byte-identical delivery fixture passed. A production absence assertion for
`net/http|os.(Create|WriteFile|RemoveAll)` under `internal/specaudit` exited 0.
All production Git calls were reviewed from the single runner: log, diff-tree,
ls-tree, rev-list, diff, symbolic-ref, rev-parse, worktree list, for-each-ref,
and status with `--no-optional-locks`; reclaim commands are returned strings.
