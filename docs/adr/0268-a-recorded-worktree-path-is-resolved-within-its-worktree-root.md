---
status: accepted
created_at: 2026-10-12T00:00:00Z
updated_at: 2026-10-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A recorded worktree path is resolved within its Worktree Root, and a lock test waits on a synchronization point

On 2026-10-10 the maintainer made `~/.roundfix/worktrees` a symbolic link to
another volume. Roundfix records every Run, Task and item worktree in the
configured spelling, `<worktree.location>/<repo-slug>/<name>`, while Git
writes the resolved path when it adds a worktree (`real_path` in
`add_worktree`, `builtin/worktree.c`) and lists that path. Reconciliation
refused every recorded path that crossed any symbolic link, so each terminal
Run became `unknown` with the reason `recorded Run Worktree path is unsafe`:
550 reasons in one reconcile and 24 GB that nothing reclaimed by 2026-10-12.
In the same cycle a test that proved the worktree administration lock by
letting a 250 ms deadline expire failed a Delivery Queue gate under load,
because the deadline also had to cover a real `git rev-parse`.

**A recorded worktree path is resolved within its Worktree Root.** The
Worktree Root of a recorded path is the directory two levels above it, the
`worktree.location` it was created under. Symbolic links are followed only
inside the Worktree Root. The repository directory and the worktree directory
below it must be real directories, never links, so the resolved path is always
the resolved Worktree Root joined with the same two names. A recorded path
whose repository or worktree directory is a link resolves outside its Worktree
Root and stays unsafe, as today. When a recorded path passes, it and a
registered path that resolve to the same directory are the same worktree.
Path comparison resolves the longest existing prefix of a path and keeps the
missing remainder verbatim, so a recorded worktree whose directory is already
gone still compares equal to the path Git keeps for it. A symbolic link in a
recorded Git root stays refused; this decision is about worktree paths only.

**A lock test proves waiting through a synchronization point.** A test that
claims an operation waits for the worktree administration lock drives it
through the package's runner seam, which already passes `execGitRunner{}` from
each exported function to an unexported one. The test's runner delegates to
real Git and reports when the Git common directory is resolved, the last step
before the lock. The test then cancels the operation while the lock is held and
requires `context.Canceled`, no worktree command, and the worktree unchanged;
after the release the same operation succeeds. No deadline, sleep or elapsed
time decides the result.

**Rejected.** Passing the configured `worktree.location` into every inspection
would change the signature of every reconcile, release, prune and disposition
caller and their tests, the daemon's disposition has no configuration to pass,
and a later change of `worktree.location` would strand every Run created under
the old one. Accepting any symbolic link would let a repository or worktree
directory replaced by a link redirect cleanup to another tree. Recording the
resolved path at creation would leave every existing record unsafe and needs a
Run Database migration. A longer deadline only moves the flake. `testing/synctest`
cannot see a Git subprocess or a file lock, which are not durably blocking.

## Consequences

- A Run, Task or item worktree under a symbolic-linked `worktree.location` is
  reconciled, released and pruned like one under a real directory, and
  `roundfix reconcile --apply` reclaims those an earlier version kept.
- A recorded path whose repository or worktree directory is a link, and a
  recorded Git root that crosses a link, remain refused.
- The disposal tests no longer depend on how fast the host runs Git.
