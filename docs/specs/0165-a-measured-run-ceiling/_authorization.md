---
status: approved
granted: 2026-09-24
action: bound Active Implement Runs machine-wide by a measured User Config ceiling and make the budget test independent of the header line
consuming: 0165-a-measured-run-ceiling
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0165

Approved on 2026-09-24: the maintainer chose delivery 7 in reduced scope — the
budget-test flake, a measured ceiling on simultaneous Runs, no Makefile or CI
change.

## Governed paths

None. Measured with `GovernedPath`: `internal/config`, `internal/store`,
`internal/cli/implement.go`, their tests and `docs/user-guide/configuration.md`
are ordinary source.

## Limits

- No Makefile, CI, cache or `.roundfixrc.yml` change.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
