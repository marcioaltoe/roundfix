# QA-18 — Pushed-and-merged Supervisor scratch worktree

Status: pass on build
`30ec663cf7ae65b3f03fcd696a576dc8fa578359`.

The public fixture harness created a disposable repository, bare `origin`, and
registered scratch worktree at
`/private/tmp/roundfix-qa0068-rerun.5LvE3b/scratch-worktree`. Its branch
`ma/spec-close-scratch-rerun` was pushed, squash-merged into `main`, and left
checked out in the worktree.

Fresh text and JSON invocations of
`bin/roundfix spec audit 0068-spec-close-audit` both exited `1`. They classified
the local branch, remote branch, and worktree as `residue`. The worktree
evidence states both required proofs: the branch is pushed to
`origin/ma/spec-close-scratch-rerun`, and its content is fully represented on
`main`.

The worktree and checked-out local branch carried the same command:

```text
git worktree remove -- '/private/tmp/roundfix-qa0068-rerun.5LvE3b/scratch-worktree' && git branch -D -- 'ma/spec-close-scratch-rerun'
```

Independent confirmation:

- Refs and `git worktree list --porcelain` were identical before and after two
  audits, proving the audit did not execute the command.
- The JSON document parsed as `roundfix-specaudit/v1` and exposed the same
  ordered command on the worktree and local branch.
- Executing that exact emitted command only inside the disposable fixture
  succeeded; the worktree path and local branch were absent afterward.

Harness: `qa-public-fixtures.sh` in this evidence directory. The earlier
F-001 reproduction is fixed on the current build.
