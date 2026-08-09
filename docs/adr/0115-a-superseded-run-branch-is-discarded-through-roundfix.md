---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A superseded Run Branch is discarded through Roundfix, with its evidence first

`roundfix reconcile` inspects terminal Run Worktrees and never integrates. That
read-only contract is deliberate and stays. Its consequence was not intended: a
stopped Run's branch is classified `unintegrated` and preserved, and Branch
Integrity Preflight then refuses to create any Run while such a branch exists for
the same head branch.

Measured on 2026-08-09, three stopped Run Branches refused two consecutive
`roundfix resolve` invocations. The integration command the refusal suggests,
`git merge --ff-only`, could not apply: the branches had diverged behind the work
that superseded them. Clearing them meant `git worktree remove --force` and
`git branch -D` by hand, which is the operation the Daemon exists to keep a
maintainer out of, performed at the moment the maintainer is already blocked.

Roundfix therefore owns a disposition for a Run Branch that has been superseded:
one whose commits are already present in the target, or whose Run was replaced by
a later one covering the same Tasks. The act writes what the branch held —
its commits, its files, and why it was classified superseded — before removing
anything, so the decision is auditable after the branch is gone.

The disposition is a separate, named act rather than a widening of `reconcile`.
Reconciliation reports; disposal changes the repository. Folding the second into
the first would make a read-only command destructive, and the read-only property
is what makes it safe to run while diagnosing.

Automatically discarding on refusal was rejected. The Preflight that refuses is
protecting work it cannot evaluate, and a command that silently deletes branches
to unblock itself is the wrong answer to being blocked.
