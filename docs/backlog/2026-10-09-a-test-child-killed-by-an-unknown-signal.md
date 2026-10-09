---
type: fix
status: open
created: 2026-10-09
spec: null
---

# A test's child process is killed by a signal nobody identified

## Problem

Spec 0254 shipped with one known failure, F6. The maintainer decided on
2026-10-09: "Publicar e investigar depois". In the full suite,
`TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` in `internal/cli`
sometimes sees its item-binary child, `migrate --check`, end with exit `-1`
and no output. `runRoundfix` reports `exec.ExitError.ExitCode()`, so `-1`
means a signal killed the child.

- The failure persists after the test went back to sequential, so it does not
  come from parallel tests in its own package.
- It happens only in the full multi-package suite.
- Forty isolated runs and two package runs pass.

The suspected cause is not proven. Another package may signal a stale PID
that the OS reused for this child. The suite starts more than 83,000 `git`
processes per run, and several store and CLI tests signal owner or helper
processes by recorded PID. Six `store.test TestOwnerProcessHelper` processes
were found orphaned for 45 hours on 2026-10-08, which shows those helpers can
outlive their tests.

## Direction

- Record which tests send signals, and to which PIDs: an instrumented helper,
  or `dtrace`/`dtruss` on `kill` during a full suite. Then find the sender.
- Guard every signal sent to a recorded PID. Verify that the process is still
  the expected one (start time or process group) before sending. Join or kill
  helper processes in `t.Cleanup` so none outlive their test.

## Sources

- The Archive Record of Spec 0254, the QA report `qa-report-2026-10-09-01.md`
  (in Git at the record's `source_revision`).
- Operator log, 2026-10-08: six orphaned `TestOwnerProcessHelper` processes
  were killed by exact PID.
