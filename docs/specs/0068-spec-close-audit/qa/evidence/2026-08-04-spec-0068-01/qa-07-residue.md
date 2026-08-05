# QA-07 — Local and remote residue branches

Status: pass.

The public harness's pushed-and-merged fixture classified both local
`ma/spec-close-scratch-rerun` and remote
`origin/ma/spec-close-scratch-rerun` as `residue`. Their evidence says the
content is fully represented on `main`.

The CLI printed the ordered worktree/local-branch reclaim command and the
remote command `git push --delete 'origin' 'ma/spec-close-scratch-rerun'`.
Text and JSON exited 1. Repeated audits left refs and worktrees identical;
only the deliberate execution of the exact local emitted command inside the
disposable fixture changed them.

`TestAuditClassifiesResidueBranch` and the remote-backup replay also passed in
the fresh 17-test selection.
