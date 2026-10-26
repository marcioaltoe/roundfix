---
task: task_01
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
status: pending
type: test
complexity: medium
---

# Task 01: The unreaped-group test waits for its member to be a zombie

## Overview

On macOS `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` probes its
process group as soon as the live-member reading is empty, which also
happens while the member is still exiting, and Darwin counts an exiting
member, so the probe can succeed where the test expects `EPERM`. This Task
makes the wait also require the member's own zombie state, through a new
reading in each platform file, and adds a Darwin test proving the reading
never takes an exiting member for a zombie. Only test files change.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-the-unreaped-group-test-fails-under-the-full-suite-on-macos.md`](references/2026-10-26-the-unreaped-group-test-fails-under-the-full-suite-on-macos.md)
   of 2026-10-26 (the Spec 0317 QA failure of 2026-10-25), whose mechanism
   the addendum of 2026-10-26 in
   `docs/references/2026-10-24-test-suite-audit.md` records, by implementing
   API Contracts 1 and 2 of `_techspec.md`.
2. MUST add `detachFixtureMemberIsZombie(pid int) (bool, error)` to each of
   the three platform files:
   - `internal/cli/detach_fixture_group_darwin_test.go`: read
     `unix.SysctlKinfoProc("kern.proc.pid", pid)`; an error, or a result
     whose `P_pid` is not `pid`, is an error; return
     `P_stat == detachFixtureZombieState` and nothing else, so a process with
     `P_WEXIT` set and another state is not a zombie;
   - `internal/cli/detach_fixture_group_linux_test.go`: read
     `/proc/<pid>/stat`, take the state field after the last `)` with the
     same parsing guards `detachFixtureGroupLiveMembers` uses (a missing `)`
     or an incomplete field is an error, checked before slicing), and return
     whether it is `Z`;
   - `internal/cli/detach_fixture_group_other_test.go`: return an error
     naming the PID and that the reading is unsupported on this system.
3. MUST change only the wait of
   `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` in
   `internal/cli/implement_detach_teardown_test.go`: the poll "fixture cat
   exited but unreaped" returns ready only when the live members are empty
   and `detachFixtureMemberIsZombie(pid)` is true, fails the test on a
   reading error, reports both values in its observation, and is preceded by
   a comment saying that an empty live reading may still be an exiting
   member. Every other line of the test keeps its bytes: the per-poll
   `kill(pid, 0)` check, the Darwin `EPERM` and Linux success expectations,
   the `detachFixtureEnded` assertion, the reap and the `ESRCH` expectation.
4. MUST add `TestDetachFixtureZombieReadingRejectsAnExitingMember` to
   `internal/cli/detach_fixture_group_darwin_test.go` as API Contract 2
   states: `t.Parallel()` as its first statement; launches of
   `/usr/bin/true` with `Setpgid` under one `testwait.Poll` whose deadline is
   `testwait.Bound(t)`, each child reaped before the next launch and on
   every failure path; a raw `kern.proc.pid` read before and after the
   reading under test, so a window that closed during the reading proves
   nothing and the launch continues; a failure that names the PID, state and
   flags when an exiting member reads as a zombie; then a wait until the
   reading says zombie and an assertion that `kill(-pid, 0)` is `EPERM`.
5. MUST NOT bound any wait by a fixed real-time window (the Real-Time
   Budget rule), signal any process other than through the existing
   cleanup helpers, or change `detachFixtureGroupLiveMembers`,
   `detachFixtureEnded`, `TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded`,
   any other test, any production file, the Makefile or the CI workflows.

## Subtasks

- [ ] Add the zombie reading to the Darwin, Linux and other Unix files.
- [ ] Make the unreaped-group wait require the zombie reading.
- [ ] Add the exiting-member proof test.

## Acceptance Criteria

- [ ] The unreaped-group test probes its group only after its member reads as a zombie, and keeps every assertion.
- [ ] The exiting-member proof passes and names the state when an exiting member reads as a zombie.
- [ ] `internal/cli` vets for Windows, Linux and FreeBSD.

## Consumers

None. Only test files of `internal/cli` change; no production contract or shared helper outside the package moves.

## Context

- instruction: `docs/references/2026-10-24-test-suite-audit.md`
- interface: `internal/cli/implement_detach_teardown_test.go`
- interface: `internal/cli/detach_fixture_group_darwin_test.go`
- interface: `internal/cli/detach_fixture_group_linux_test.go`
- interface: `internal/cli/detach_fixture_group_other_test.go`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded|TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded|TestDetachFixtureZombieReadingRejectsAnExitingMember)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded TestDetachFixtureZombieReadingRejectsAnExitingMember; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -q 'detachFixtureMemberIsZombie(pid)' internal/cli/implement_detach_teardown_test.go || { echo 'the unreaped-group wait does not read the zombie state' >&2; exit 1; }; GOOS=windows go vet ./internal/cli && GOOS=linux go vet ./internal/cli && GOOS=freebsd go vet ./internal/cli` — expected: exit 0 on the macOS Run host. Before this Task the exiting-member proof does not exist, so the pass check fails on its name. Afterwards the three detach tests pass, the unreaped-group wait reads the zombie state, and every platform split compiles.

## References

- `_prd.md` → Goals 1-3; Core Features 1-2; Success Metrics 2, 3 and 5
- `_techspec.md` → Current behavior 1-5; Interfaces; API Contracts 1-2; Invariants 1-3; Testing Approach; Build Order 1
