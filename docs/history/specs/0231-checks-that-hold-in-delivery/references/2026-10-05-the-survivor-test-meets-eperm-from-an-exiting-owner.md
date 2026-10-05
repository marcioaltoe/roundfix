---
type: fix
status: promoted
created: 2026-10-05
spec: 0231-checks-that-hold-in-delivery
reason: null
---

# The detach survivor test meets EPERM from an owner that is already exiting

## Symptom

On 2026-10-05 the Daemon's `make verify-changed` for Run
`run_20261005T101255Z_f599a5cc8ea7f4c2` (batch 3, attempt 1) failed one test:

```text
--- FAIL: TestRunImplementDetachSurvivesCallerProcessGroupKill (5.07s)
    implement_test.go:1862: kill fixture process group 52154: operation not permitted; live members=[52154]; reading error=<nil>
```

Every other package passed. Spec 0220 had already replaced an earlier skip of
this `EPERM` with a reading of the group's live members, so the test now fails
instead of skipping whenever it meets the window. The operator carried the
Run forward (intervention log entry 161).

## Where

`internal/cli` tests: the cleanup `reapOwner` of
`TestRunImplementDetachSurvivesCallerProcessGroupKill` calls
`killDetachFixtureGroup`, which accepts `EPERM` only when
`detachFixtureGroupLiveMembers` (Darwin: `kern.proc.pgrp`, excluding `SZOMB`)
reports no live member.

## Expected

The test passes in every environment where it ran before, without a skip,
while a fixture that is still running is never counted as ended.

## Evidence

Measured during authoring on 2026-10-05 (macOS, Darwin 27, no sandbox):

- The detached owner in this test does not hold after its Run reaches Clean
  (the hold applies only without the detach handshake), so it exits on its
  own while the test's cleanup sends `SIGKILL` to its group.
- XNU's `proc_exit` marks the process `P_REF_DEAD` ("will not be visible via
  proc_find") long before it sets `p_stat = SZOMB`. `killpg1` skips a member
  that `proc_find` cannot reference, counts nothing, and returns `EPERM`,
  while `kern.proc.pgrp` still reports the member with `p_stat` 2 (`SRUN`)
  and `P_WEXIT` set in `p_flag`.
- A probe that starts a child in its own group and polls `kill(-pgid, 0)`
  caught that window in 150 of 150 launches (`/usr/bin/true`,
  `sh -c 'exit 0'`, `cat`) and in 80 of 80 launches of a Go child. In the
  210 launches whose probe also read `p_flag`, every `EPERM` came with the
  only non-zombie member's `P_WEXIT` set, and none with a member lacking it.
- The Daemon runs Verification outside the Agent sandbox
  (`internal/daemon/task_engine.go`), and the probe reproduced the window in
  an ordinary shell, so the sandbox is not the cause.
