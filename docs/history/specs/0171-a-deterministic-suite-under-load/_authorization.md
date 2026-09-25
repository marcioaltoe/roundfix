---
status: approved
granted: 2026-09-25
action: make the repository test suite deterministic on a loaded machine, keep every test process off the live release lookup, and resolve the adapter fixture's link target
consuming: 0171-a-deterministic-suite-under-load
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0171

The maintainer approved the 2026-09-25 efficiency sequence (Ondas 0 and 1) in
chat on 2026-09-25, and this Spec is one of its deliveries. The Go sources it
touches would ride the standing grant of 2026-09-21 for governed source, but
none of them is governed.

## Governed paths

None. Measured with `GovernedPath` on 2026-09-25: `internal/testwait/testwait.go`,
`internal/testwait/testwait_test.go`, `internal/cli/implement_test.go`,
`internal/cli/version_freshness_isolation_test.go`,
`internal/daemon/task_engine_test.go`,
`internal/daemon/run_disposition_characterization_test.go`,
`internal/worktree/worktree_test.go`, `internal/agent/acpx_runner_test.go` and
this Spec's own files are ordinary.

## What is not governed and must stay untouched

`internal/cli/cli_test.go`, `docs/references/coverage-record.json` and the
Makefile are governed. No Task edits them: the harness in `cli_test.go` is not
needed, every recorded test keeps its name, and no suite flag changes. A Task
that finds it cannot avoid one of them stops and reports instead of editing it.

## Sanctioned regeneration

None.

## Limits

- No production source, command output, flag, exit code or schema change.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
