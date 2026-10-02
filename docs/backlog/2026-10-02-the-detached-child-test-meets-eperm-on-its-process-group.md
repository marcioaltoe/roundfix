---
type: fix
status: open
created: 2026-10-02
spec: null
reason: null
---

# The detached-child test meets EPERM on its process group

## Symptom

`TestImplementDetachChildEndsWhenItsTestBinaryDies`, added by Spec 0213 to prove that a fixture process ends with the test binary that started it, failed in the Daemon's Verification of Spec 0212's task_01 on 2026-10-02 (Run `run_20261002T094236Z_2384f2babb94acdf`, `verification/batch-001-attempt-*.log`):

```text
implement_detach_teardown_test.go:142: probe fixture process group 95744: operation not permitted
implement_detach_teardown_test.go:75: kill fixture process group 95744: operation not permitted
```

The same test passed on `main` from an interactive shell minutes later and in Spec 0213's Linux CI. Spec 0213 was delivered to remove exactly this kind of intermittent failure.

## Where

`internal/cli/implement_detach_teardown_test.go`; the probe and cleanup signal the fixture's process group by its recorded ID.

## Expected

The test proves the fixture ended without signalling a process group it no longer owns. `EPERM` from `kill(-pgid, 0)` after the group has exited (its ID reused or owned by another session) counts as ended, or the test identifies the group by a handle it cannot lose, such as a pidfd or the recorded start time. A repeated run in the Daemon's Verification environment (`-count` high enough to have failed before) proves the fix.
