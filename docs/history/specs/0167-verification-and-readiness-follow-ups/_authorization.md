---
status: approved
granted: 2026-09-24
action: let the exclusive retry's verdicts replace the first run's, skip completed repair Tasks in the verbatim planning check, and print a degraded access policy
consuming: 0167-verification-and-readiness-follow-ups
paths:
  - internal/cli/cli_test.go
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

`internal/cli/cli_test.go` rides the standing grant of 2026-09-21 for governed
paths a slice genuinely needs: it pins the command surface's text output, which
the degraded-access line changes. The first Run of this Spec proved the need:
Task 02 edited it and was refused for lack of a bound path.

Everything else — `internal/daemon`, `internal/cli/profiles_validate.go`,
`internal/cli/doctor.go`, `internal/config`, `internal/cli/implement.go` and
their other tests — is ordinary source, measured with `GovernedPath`.

## Limits

- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
