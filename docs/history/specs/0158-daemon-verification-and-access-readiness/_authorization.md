---
status: approved
granted: 2026-09-24
action: let a Task declare independent Verification, let the frozen authorization name Tasks that repair a red repository gate, and prove access policy in readiness
consuming: 0158-daemon-verification-and-access-readiness
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0158

Approved on 2026-09-24 as delivery 3b of the restructured queue, under the
maintainer's standing authorization for the queue order and its operations.

## Governed paths

None. Measured as the intersection of this Spec's changed paths with the literal
set in `internal/speccheck/governed.go`: empty. `internal/spec/task.go`,
`internal/authorization`, `internal/daemon`, `internal/agent`,
`internal/cli/profile_preflight.go`, `docs/adr` and
`docs/user-guide/configuration.md` are ordinary.

## Limits

- No skill, Baseline asset or generated guide changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
