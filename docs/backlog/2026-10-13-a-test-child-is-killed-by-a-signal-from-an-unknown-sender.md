---
type: fix
status: open
created: 2026-10-13
spec: null
reason: null
---

# A test child is killed by a signal from an unknown sender (F6)

## Symptom

`TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate/stdout` fails when the
`migrate --check` child it runs exits with `-1` (killed by a signal). It
happened in 0254, again in the starting verification of 0259 and in suite 7 of
the ten full suites of its QA. It passes 20 of 20 alone, no other Roundfix
process was running, and the copy-over-a-running-binary hypothesis was ruled
out twice.

## Where

`internal/cli` delivery tests that run the item binary; the sender is unknown.

## Expected

First, the failure names the signal (`syscall.WaitStatus.Signal`) in the test
report; then the sender is found and the cause fixed, so ten full suites pass.
Tracing the sender with `eslogger` needs `sudo`, which only the maintainer can
run.

## Evidence

Operator log entries 305 and 306 (2026-10-13); v0.72.0 shipped with a QA
override the maintainer approved and lists F6 as a known issue.
