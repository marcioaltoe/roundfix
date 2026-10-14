---
type: fix
status: open
created: 2026-10-13
spec: null
reason: null
---

# Trace the sender of the F6 signal

## Symptom

F6, `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate/stdout`, fails
when its `migrate --check` child is killed by a signal. It occurred in 0254,
in the starting verification of 0259, in suite 7 of 0259's ten full
suites, and in the first attempt of 0268's task_02 on 2026-10-14
(`migrate --check exited -1` instead of 2, before Signal Exit shipped): four
recorded occurrences.

## Where

`internal/cli` delivery tests that run the item binary.

## Expected

Find and fix the sender so ten full suites pass.

## Evidence

Since this Spec, the failure names the signal and reports a Signal Exit. Every
recorded occurrence was the stdout subtest, the second subtest of
`TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate`, whose build copies
the fixture over an existing `bin/roundfix`; this is an observation, not a
cause. Tracing the sender with `eslogger` needs `sudo`, which only the
maintainer can run.
