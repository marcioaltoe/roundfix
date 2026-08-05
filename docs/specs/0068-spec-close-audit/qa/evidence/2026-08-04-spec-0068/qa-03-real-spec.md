# QA-03 — Real Spec operator journey

Status: pass.

The built `bin/roundfix --help` and `bin/roundfix spec audit --help` both list
`roundfix spec audit <slug> [--format <text|json>]` and state the read-only
contract. Running only the real slug exited 1 with an attention report:

- The current Run Branch classified `pending` with six commits, two
  branch-only files, and 22 differing shared files not on `origin/HEAD`.
- The current Run Worktree classified `preserved` with explicit
  no-matching-Run evidence in the Agent-visible Run Database.
- `internal/specaudit/audit.go` and `audit_test.go` were undelivered and held
  by the Run Branch.

A fresh JSON process returned the same survivors and artifacts in one
`roundfix-specaudit/v1` object. This Agent-visible Run Database does not bind
the current worktree path to a Run, so QA-09 uses an injected Active Run
fixture for that independent safety proof.
