# R01 — blocking mechanical stage

The real nested Implement Run is environment-blocked: this QA assignment
forbids commits, while a disposable public Run requires a committed Spec graph
and creates Task/QA commits. The current child prompt also lacks the seeded
report and changed-context fields, so it cannot serve as observation of the
new stage that was meant to invoke it.

Equivalent fresh evidence:

- `rtk ... go test ./internal/daemon -run '^TestMechanicalStageWithholdsAgentSession$' -count=1 -v`
  passed one focused case. The whole focused process took 12.66 seconds,
  including a cold package compile; this is an upper-bound process timing, not
  the required gate-start-to-event timing.
- The case observes the QA Task still pending during the stage, a closed fail
  report, zero Agent requests, zero Verification calls/events, one mechanical
  event, and the existing Daemon settlement afterward.
- The stage publishes its actual `duration_ms` at
  `internal/daemon/task_engine.go:1806-1815`.

The required production wall-clock observation is not available from this
report-only Run and is not treated as passed.
