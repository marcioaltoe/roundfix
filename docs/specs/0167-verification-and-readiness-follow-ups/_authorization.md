---
status: approved
granted: 2026-09-24
action: let the exclusive retry's verdicts replace the first run's, skip completed repair Tasks in the verbatim planning check, and print a degraded access policy
consuming: 0167-verification-and-readiness-follow-ups
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0167

Minted on 2026-09-24 as the follow-up Spec 0158 recorded, under the
maintainer's standing rules for carried limits.

## Governed paths

None. Measured with `GovernedPath`: `internal/daemon`, `internal/cli/profiles_validate.go`,
`internal/cli/doctor.go` and their tests are ordinary source.

## Limits

- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
