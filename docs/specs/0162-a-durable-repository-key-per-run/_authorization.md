---
status: approved
granted: 2026-09-24
action: record a durable repository key on each Run and make listing, reconcile, gc and Run Window lookups use it
consuming: 0162-a-durable-repository-key-per-run
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0162

Minted on 2026-09-24 as the follow-up that Spec 0157 recorded, under the
maintainer's standing rules: at the corrective ceiling, deliver with limits
recorded and carry them to a Spec of their own.

## Governed paths

None. Measured with `GovernedPath`: `internal/store`, `internal/config`,
`internal/cli/reconcile.go`, `internal/cli/gc.go`, `internal/cli/window.go` and
their tests are ordinary source.

## Limits

- No artifact directory is moved and no Run record is deleted.
- No skill, Baseline asset or generated guide changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
