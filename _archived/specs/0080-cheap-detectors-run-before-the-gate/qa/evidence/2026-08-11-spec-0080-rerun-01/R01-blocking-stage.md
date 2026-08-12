# R01 — blocking mechanical stage

Build: `c2372a9f709c9197aa5c5e89fd71da1ab46f07e6`.

`rtk env GOCACHE=/private/tmp/roundfix-spec0080-rerun-gocache go test
./internal/daemon -run
'^(TestMechanicalStageWithholdsAgentSession|TestQAMechanicalRequestSelectsTheAuthorizedTaskCommit)$'
-count=1 -v` exited 0 in 0.85 seconds for the process.

- `TestQAMechanicalRequestSelectsTheAuthorizedTaskCommit` created an isolated
  committed Git carrier, selected only the Task commit whose declared tooling
  scope applied, ran the real mechanical stage, and found the folded
  `outside.txt` path as a blocking fact in 0.13 seconds.
- `TestMechanicalStageWithholdsAgentSession` ran the Daemon TaskCycle with a
  blocking typed result and passed in 0.01 seconds. It proved zero Agent
  Session requests, zero Verification calls/events, a closed fail report with
  zero pending rows, one finding-blocked row, Daemon-owned failed Task status,
  the existing report-commit message, and a `daemon.qa` mechanical event.

The assignment forbids creating or committing a disposable public Spec Run,
so a single production CLI journey from gate start through real detector,
report persistence, event stream, and process exit remains unreachable. The
two focused supervised seams cover the exact detector and withholding halves,
and the measured process stays well inside the PRD's seconds-not-a-round goal.
Unblock the combined observation with an authorized disposable Implement Run
that may commit its fixture and QA Report.
