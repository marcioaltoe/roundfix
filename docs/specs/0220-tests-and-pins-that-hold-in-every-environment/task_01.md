---
task: task_01
spec: 0220-tests-and-pins-that-hold-in-every-environment
status: pending
type: test
complexity: medium
---

# Task 01: A fixture's process group is proved ended by its live members

## Overview

The detached-child death tests decide that a fixture ended by sending signal 0
to its process group. On macOS a group whose members have all exited, but are
not yet reaped, answers `EPERM`, and the stopgap skips the test on that answer.
This Task reads the group's live members from the process table instead,
excludes exited and unreaped processes, guards the recorded PID by its start
identity, removes the skip, and characterizes the platform answer with a
zombie group the test holds itself.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-02, "The detached-child test meets
   EPERM on its process group" (adopted at
   [references/2026-10-02-the-detached-child-test-meets-eperm-on-its-process-group.md](references/2026-10-02-the-detached-child-test-meets-eperm-on-its-process-group.md)),
   by removing the `EPERM` skip. No death test may call `t.Skip`, `t.Skipf` or
   `t.SkipNow`.
2. MUST add `detachFixtureGroupLiveMembers(pgid int) ([]int, error)` in three
   platform test files, with the shape of the TechSpec's Interfaces:
   - Darwin reads `kern.proc.pgrp` through `unix.SysctlKinfoProcSlice` and
     drops entries whose `Proc.P_stat` is `SZOMB` (5 in `<sys/proc.h>`, held in
     a named test constant);
   - Linux scans `/proc/<pid>/stat`, parses the state and process-group fields
     after the last `)`, drops state `Z`, and skips a PID that vanished during
     the scan (`ENOENT` or `ESRCH`);
   - the other Unix file keeps the signal probe and lists the group as live
     only when `kill(-pgid, 0)` succeeds.

   A failed reading MUST be returned as an error, never as an empty list.
3. MUST decide that a recorded fixture ended only when its group has no live
   member, or when its PID now carries a start identity other than the one
   `store.OwnerProcessIdentity` returned while the fixture was alive. The death
   tests MUST record that identity after they see the fixture alive and before
   they kill the inner test binary.
4. MUST make the group cleanup accept `ESRCH`, and accept `EPERM` only when the
   reading lists no live member; any other answer MUST fail the test with the
   group ID and the error.
5. MUST add `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded`. It starts
   `cat` in its own process group with a pipe as its standard input, so the
   child lives until the test closes that pipe, asserts that the reading lists the child, closes the
   input, and does not wait for the child. Once the reading lists no live
   member while `kill(pid, 0)` still succeeds, it MUST assert that
   `kill(-pgid, 0)` answers `EPERM` on Darwin and succeeds on Linux, and that
   the fixture is reported ended. It then MUST reap the child and assert that
   `kill(-pgid, 0)` answers `ESRCH`. Every wait MUST be bounded by the test
   deadline through `testwait`.
6. MUST keep the `SIGKILL` wait-status assertion and the "no cleanup signal
   before the assertion" order of both death tests, and keep
   `TestRunImplementDetachSurvivesCallerProcessGroupKill` and
   `TestDetachedChildIsTerminatedAtTeardown` passing with their survival
   assertions unchanged.
7. MUST signal only PIDs the test started or read from its own records, and
   MUST NOT change production code.

## Subtasks

- [ ] Add the three platform readings of a process group's live members.
- [ ] Decide fixture end by live members and the recorded start identity.
- [ ] Accept `EPERM` in cleanup only when no live member remains, and remove
      the skip.
- [ ] Add the zombie-group characterization test.

## Acceptance Criteria

- [ ] A group whose only member exited and is unreaped is reported ended, on
      Darwin while `kill(-pgid, 0)` answers `EPERM` and on Linux while it
      succeeds.
- [ ] A live child in its own group is listed by the reading.
- [ ] Both death tests pass 10 consecutive iterations with no skip.
- [ ] The test package still compiles for Linux and FreeBSD.

## Context

- interface: `internal/cli/implement_detach_teardown_test.go`
- interface: `internal/cli/implement_test.go`
- creates: `internal/cli/detach_fixture_group_darwin_test.go`
- creates: `internal/cli/detach_fixture_group_linux_test.go`
- creates: `internal/cli/detach_fixture_group_other_test.go`
- instruction: `internal/store/process.go`
- instruction: `internal/store/process_darwin.go`
- instruction: `internal/store/process_linux.go`
- instruction: `internal/testwait/testwait.go`
- instruction: `docs/adr/0213-a-test-fixture-process-ends-with-the-test-binary-that-started-it.md`

## Verification

- `out="$(go test -count=10 -timeout 900s -v -run '^(TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded|TestImplementDetachChildEndsWhenItsTestBinaryDies|TestDetachSurvivorEndsWhenItsTestBinaryDies)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; if printf '%s\n' "$out" | grep -qF -- '--- SKIP'; then printf '%s\n' "$out"; exit 1; fi; for name in TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded TestImplementDetachChildEndsWhenItsTestBinaryDies TestDetachSurvivorEndsWhenItsTestBinaryDies; do n="$(printf '%s\n' "$out" | grep -cF -- "--- PASS: $name (")"; test "$n" -eq 10 || { printf 'passes of %s: %s, want 10\n' "$name" "$n" >&2; exit 1; }; done` — expected: exit 0; before this Task the characterization test does not exist, so the command fails.
- `test -f internal/cli/detach_fixture_group_darwin_test.go && test -f internal/cli/detach_fixture_group_linux_test.go && test -f internal/cli/detach_fixture_group_other_test.go && dir="$(mktemp -d)" && GOOS=linux go test -buildvcs=false -c -o "$dir/cli-linux.test" ./internal/cli && GOOS=freebsd go test -buildvcs=false -c -o "$dir/cli-freebsd.test" ./internal/cli` — expected: exit 0; before this Task the platform files do not exist, so the command fails.
- `! grep -Eq 't[.]Skip' internal/cli/implement_detach_teardown_test.go` — expected: exit 0; before this Task the probe calls `t.Skipf`, so the command fails.

## References

- `_prd.md` → Goals 1-2; Core Feature 1; Success Metric 1; Acceptance evidence
- `_techspec.md` → Why an `EPERM` arrives; Interfaces; Testing Approach 1-2;
  Build Order 1
- [references/2026-10-02-the-detached-child-test-meets-eperm-on-its-process-group.md](references/2026-10-02-the-detached-child-test-meets-eperm-on-its-process-group.md)
- ADR-0213; ADR-0125; ADR-0126; ADR-0028
