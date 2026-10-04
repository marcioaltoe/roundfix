---
task: task_01
spec: 0220-tests-and-pins-that-hold-in-every-environment
status: completed
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

## Result

Implemented the Task's test-only slice. Darwin reads `kern.proc.pgrp` and
excludes the named `SZOMB` state; Linux reads the state and group after the
last `)` in procfs stat records, excludes `Z`, and tolerates only vanished
PIDs (`ENOENT`/`ESRCH`). Other Unix systems retain the signal probe and return
permission failures as errors. Failed readings cannot prove an empty group.

Both death tests record the live fixture's start identity before killing the
inner binary. The exit assertion now requires no live group member or a
changed leader identity; an absent leader with live children does not prove
exit. Cleanup accepts `ESRCH`, and accepts `EPERM` only after a successful
reading with no live members. Removed the skip, preserving the `SIGKILL`
wait-status assertion and the absence of fixture cleanup signals before the
exit assertion. The existing survival assertions remain unchanged.

The new characterization test holds a real `cat` child alive on an input
pipe, verifies live membership and identity, closes the pipe without reaping,
checks the platform's zombie-group signal answer and the ended predicate,
then reaps through `testwait` and checks `ESRCH`. It also checks that a
different recorded identity reports the old fixture ended. The zombie
characterization applies to Darwin and Linux; other Unix platforms retain
the specified fallback and compile the package.

Focused evidence from this Agent turn:

| Acceptance criterion | Evidence |
| --- | --- |
| An exited, unreaped group is reported ended | `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` passed on Darwin, asserting `kill(pid, 0)` succeeds, `kill(-pgid, 0)` answers `EPERM`, the ended predicate is true, and the reaped group answers `ESRCH`. Linux runtime behavior remains for a Linux executor. |
| A live child is listed | The same test passed its explicit live-child membership assertion and rejected ended for the matching start identity. |
| Both death tests pass ten times without skipping | Each death test passed once in the final focused run with no skip. Source inspection found no `t.Skip` call in the teardown file. The authored ten-iteration gate remains exclusively for Daemon Verification. |
| Linux and FreeBSD test package compilation | Both final foreign-platform compile checks exited 0, using `-exec /usr/bin/true` to avoid executing their binaries. |

Commands and outcomes:

- `rtk proxy go test -count=1 -timeout 180s -v -run '^(TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded|TestImplementDetachChildEndsWhenItsTestBinaryDies|TestDetachSurvivorEndsWhenItsTestBinaryDies|TestRunImplementDetachSurvivesCallerProcessGroupKill|TestDetachedChildIsTerminatedAtTeardown)$' ./internal/cli`
  — exit 0; all five tests passed on Darwin, no skip (5.519s package time).
- `GOOS=linux rtk proxy go test -buildvcs=false -exec /usr/bin/true -run '^$' ./internal/cli`
  — exit 0; compiled the Linux test package, no foreign tests executed.
- `GOOS=freebsd rtk proxy go test -buildvcs=false -exec /usr/bin/true -run '^$' ./internal/cli`
  — exit 0; compiled the FreeBSD test package, no foreign tests executed.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

The first sandboxed characterization attempt could not open the host Go
build cache (`operation not permitted`). Focused Go checks then ran with
approved host cache/process-table access. Before editing, the task file was
the sole pre-existing changed path; the old teardown source contained the
`EPERM` skip and lacked the characterization test and platform readers.

No production code, other Task, Task Graph, tooling configuration or skill
content changed. Task status and authored Verification remain Daemon-owned;
no authored Verification command was run, and no commit, push or PR was made.

## Carry-forward provenance

- Source Run: `run_20261004T010416Z_bdbf1cc82a3473ae`
- Source commit: `39cd3d980a0c1eb78aa741e228100b2e241f05b7`
