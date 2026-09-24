---
status: approved
granted: 2026-09-24
action: record a durable repository key on each Run and make listing, reconcile, gc and Run Window lookups use it
consuming: 0162-a-durable-repository-key-per-run
paths:
  - internal/cli/cli_test.go
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

`internal/cli/cli_test.go` rides the standing grant of 2026-09-21 for governed
paths a slice genuinely needs. It seeds outdated Run Databases by dropping the
columns later migrations add; a new `runs.repository_root` column must be
dropped there too, or the Branch Integrity Preflight migration test seeds a
database that already has it. The first Run of this Spec proved the need: its
Task edited that seed and was refused for lack of a bound path.

Everything else — `internal/store`, `internal/config`,
`internal/cli/reconcile.go`, `internal/cli/gc.go`, `internal/cli/window.go` and
their other tests — is ordinary source, measured with `GovernedPath`.

## Limits

- No artifact directory is moved and no Run record is deleted.
- No skill, Baseline asset or generated guide changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
