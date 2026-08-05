# QA-03 — Real Spec operator journey

Status: pass.

Fresh built-command results:

- `bin/roundfix --help` exited 0 and listed
  `roundfix spec audit <slug> [--format <text|json>]` plus the Spec command's
  audit purpose.
- `bin/roundfix spec audit --help` exited 0 and documented the read-only
  promise, two formats, classifications, reclaim output, and exits 0/1/2.
- `bin/roundfix spec audit 0068-spec-close-audit` required only the slug and
  exited 1 with an attention report for the current real repository. It named
  the pending Run Branch, preserved Run Worktree, and three undelivered
  artifacts held by `ma/spec-0068-implementation`.
- A fresh JSON process returned the same state in one
  `roundfix-specaudit/v1` object.

Pre/post ref, worktree, status, and Run Database snapshots were identical.
