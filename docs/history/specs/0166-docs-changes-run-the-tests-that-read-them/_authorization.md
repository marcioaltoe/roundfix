---
status: approved
granted: 2026-09-24
action: make documentation and root Markdown changes select both test sets in the selective gate
consuming: 0166-docs-changes-run-the-tests-that-read-them
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0166

Minted on 2026-09-24 as the follow-up Spec 0155 recorded, under the
maintainer's standing rules for carried limits.

## Governed paths

None. Measured with `GovernedPath`: `internal/verifyselect` and its tests are
ordinary source.

## Limits

- No Makefile, CI or `.roundfixrc.yml` change.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
