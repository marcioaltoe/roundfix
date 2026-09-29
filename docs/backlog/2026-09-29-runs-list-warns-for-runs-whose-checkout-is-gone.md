---
type: fix
status: open
created: 2026-09-29
spec: null
reason: null
---

# `runs list` warns on every call about Runs whose checkout was deleted

## Symptom

On 2026-09-29, after the extra checkouts `~/dev/roundfix-wt-0173` through `-0180` and one delivery item worktree were removed, `roundfix runs list` printed `No Runs found.` and then 21 stderr warnings. Each warning read `inspect retained terminal Runs in repository "<removed path>": inspect terminal Run: stat recorded Git root …: no such file or directory`. The warnings repeat on every call, and no command clears them.

## Where

The retained-worktree note in `internal/worktree/worktree.go`: terminal Runs are grouped by recorded Git root, and a failed `recordedGitRoot` is reported as a warning per group.

## Expected

A terminal Run whose recorded checkout no longer exists has nothing to retain. `runs list` does not warn about it. The Run is left to `gc` and `reconcile` like any other terminal Run.

## Evidence

`bin/roundfix runs list` on `a626494c`, 2026-09-29: 21 stderr lines and exit 0.
