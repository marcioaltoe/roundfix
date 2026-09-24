---
status: approved
granted: 2026-09-24
action: make reconciliation refuse unproven candidates, recognise archived report copies, and give a repository one identity across its worktrees
consuming: 0157-reconciliation-and-one-repository-identity
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0157

Approved on 2026-09-24 as part of delivery 3 of the restructured queue.

## Governed paths

None. Measured as the intersection of this Spec's changed paths with the literal
set in `internal/speccheck/governed.go`: empty. `internal/config`,
`internal/worktree`, `internal/store` and their tests are ordinary source.

## Limits

- The main checkout's repository identity must not change.
- No skill, Baseline asset or generated guide changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
