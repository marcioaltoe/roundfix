---
status: approved
granted: 2026-09-15
created: 2026-09-15
action: replace the stale historical audit expectation with a controlled fixture and make its skip guard require reachability
consuming: 0139-a-suite-that-passes-where-it-runs
paths:
  - internal/speccheck/mechanical_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0139

On 2026-09-15 the maintainer approved the proposed record for this Spec with
"Aprovar como proposto". That proposal named three paths. The pre-PR review then
showed that only one of them is in the governed set this repository declares.
The other two, `internal/daemon/task_engine.go` and
`internal/daemon/task_engine_test.go`, are ordinary source, so this record now
bounds only the governed path. Narrowing the record removes authority and grants
none.

The same day, the maintainer decided that the repository's own suite is fixed
before Spec 0138's QA gate is retried, and that the Daemon, not the QA Agent,
runs the gate's repository Verification. This record lands in `main` before the
Task commit that consumes it.

## Why a governed path is unavoidable

`internal/speccheck/mechanical_test.go` is in the historically bounded set
(ADR-0130). It holds two things this Spec must change:

- the historical audit subtest, whose expectation went stale once the audit
  began reading a grant at its authorizing ancestor;
- the skip guard, which checks whether a Git object exists instead of whether it
  is reachable.

The repair belongs where the audit's test contract already lives.

No authorization has ever bounded these files, so they are ordinary source and
need no grant: `internal/agent/acpx_runner_test.go`,
`internal/daemon/daemon.go`, `internal/daemon/daemon_test.go`,
`internal/daemon/task_engine.go` and `internal/daemon/task_engine_test.go`.

## Approved bounded mutation

- Remove the historical subtest and its unreachable commit constants.
- Add a test that builds its own repository. It proves that a grant widened after
  its consuming commit is refused, and that a grant already bounded at the
  authorizing ancestor authorizes the change.
- Make `firstMissingMechanicalCommit` require ancestry of the audited `HEAD`,
  keeping its signature.

## Limits

- No action, operation or path beyond those above.
- No other existing assertion in the file is weakened or deleted.
- No other governed path is changed. The ordinary source this Spec changes
  (named above) needs no grant and is outside this record's limits.
- No Baseline module, linter, Verification configuration or `.roundfixrc.yml`
  edit; no paid API use, release, tag or deployment.
- Verification remains Daemon-owned; Task status remains Daemon-written.
