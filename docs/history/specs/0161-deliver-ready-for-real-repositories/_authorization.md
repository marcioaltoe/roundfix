---
status: approved
granted: 2026-09-24
action: make roundfix deliver wait for unreported checks, never publish to the default branch, and survive parks, crashes and reused owner PIDs
consuming: 0161-deliver-ready-for-real-repositories
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0161

Minted on 2026-09-24 as the blocking follow-up that Spec 0156 recorded, under
the maintainer's standing rules: at the corrective ceiling, deliver with limits
recorded and carry them to a Spec of their own.

## Governed paths

None. Measured with `GovernedPath` over every path this Spec changes:
`internal/delivery`, `internal/cli/deliver.go`,
`internal/cli/deliver_workflow.go`, `internal/store/delivery.go` and their
tests are ordinary source.

## Limits

- No new stage and no new external action.
- No skill, Baseline asset or generated guide changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
