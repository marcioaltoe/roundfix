---
status: approved
granted: 2026-09-24
action: add a selective verification gate and point Runs and pull request CI at it
consuming: 0155-a-verify-that-runs-what-changed
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0155

The maintainer approved this delivery on 2026-09-24, choosing one repository
with selective verification gates.

## Governed paths

None. Measured as the intersection of this Spec's changed paths with the literal
set in `internal/speccheck/governed.go`: empty. `Makefile`, `.roundfixrc.yml`,
`.github/workflows/ci-verify.yml` and the new selector package are ordinary
source.

## Limits

- `make verify` keeps its meaning: the complete gate.
- No change to the generated guides, Baseline assets or skill files.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
