---
type: fix
status: open
created: 2026-11-09
spec: null
reason: null
---

# The disk-full tests read a real host path

## Symptom

`TestTheItemBuildOnAFullDiskIsADiskFullFailure/path_in_the_log` and `TestTheRepositoryGateOnAFullDiskParksDiskFull/full_disk` fail on a machine where `/tmp/x` exists as a regular file. Spec 0365 lost its first Run this way on 2026-11-09 (Run `run_20261109T120541Z_75e62679b1abfda6`, task_01, `make verify-changed`): an empty `/tmp/x` file created on the maintainer's MacBook earlier that morning made the reported free space `0 MiB` instead of the stubbed value, and the item parked `run-unresolved` for a defect the Spec did not touch.

## Where

`internal/cli/deliver_disk_full_test.go` writes the literal path `/tmp/x/_pkg_.a` into the fake build output. `internal/cli/deliver_disk_full.go` `diskFullFailure` calls `doctorDiskSpaceAncestor` on that path before the stubbed `availableBytes`; `os.Stat` on `/tmp/x/_pkg_.a` returns `ENOTDIR` when `/tmp/x` is a file, which is not `os.ErrNotExist`, so the ancestor lookup fails and `AvailableBytes` stays 0.

## Expected

The tests use a path under their own temporary directory (absent child of `t.TempDir()`), so no file on the host can change their outcome. Whether `doctorDiskSpaceAncestor` should also walk past `ENOTDIR` (a file where a directory is expected) is a design choice for the Spec.

## Evidence

Diagnostics `/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261109T120541Z_75e62679b1abfda6/verification/batch-001-attempt-2.log`; removing the empty `/tmp/x` (created 2026-11-09 09:05 BRT, origin unknown) made both tests pass 3 of 3 on main. Recorded by the operator on 2026-11-09.
