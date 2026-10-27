---
task: task_01
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
status: completed
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


## Result

Implemented API Contracts 1 and 2 within the four declared test files.
The initial checkout had no zombie reader or exiting-member proof; the only
pre-existing changed path was this Task's Daemon-owned `in_progress` status.

- Acceptance criterion 1: the unreaped-member poll now requires both no live
  group members and a positive per-PID zombie reading, fails on a reading
  error, and observes both values. A byte comparison against HEAD, removing
  only the authorized wait edits, confirmed the teardown file is otherwise
  identical. The existing Darwin group reader and exiting-member test are
  also byte-identical. Darwin checks only `P_stat == SZOMB`; Linux checks
  the state after the last `)` with the existing parsing guards; other Unix
  systems return a PID-specific unsupported error.
- Acceptance criterion 2: the new parallel Darwin proof brackets its reader
  with raw PID reads, retries launches when the window closes, then waits
  for the zombie and asserts group-probe `EPERM`. Every started child is
  reaped by a callback-local defer before another launch or on failure;
  `/usr/bin/true` exits naturally, so no cleanup signal is needed. All waits
  use the test deadline. Focused command:
  `go test -count=10 -timeout 3m -v -run '^(TestDetachFixtureZombieReadingRejectsAnExitingMember|TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded)$' ./internal/cli`
  exited 0: both tests passed ten times. The proof observed state `2`, flags
  `0x6004`, in each run and reached the zombie `EPERM` assertion.
- Negative focused check: a temporary Go overlay under
  `/tmp/roundfix-zombie-probe-h2ya_f5i` made the Darwin reader also accept
  `P_WEXIT`. `go test -overlay /tmp/roundfix-zombie-probe-h2ya_f5i/overlay.json -count=1 -timeout 3m -v -run '^TestDetachFixtureZombieReadingRejectsAnExitingMember$' ./internal/cli`
  exited 1 as intended: `exiting fixture 92562 read as a zombie: state=2
  flags=0x6004; before state=2 flags=0x6004`. The repository reader was
  unchanged by this probe.
- Acceptance criterion 3: focused cross-compilation exited 0 for each of
  Windows, Linux and FreeBSD, using `GOOS=<os> go test -c -o
  /tmp/roundfix-zombie-probe-h2ya_f5i/cli-<os>.test ./internal/cli` for
  `windows`, `linux` and `freebsd`. These compile the package and its tests;
  they do not establish runtime behavior on those systems. The declared
  platform vets remain Daemon-owned.
- `git diff --check` exited 0. Final changed-path inspection contains only
  the four declared test files and this Task's Result record.

The Task status and authored checkboxes are unchanged. Declared Verification,
repository Verification and the full suite were not run. No follow-up work
was added to this slice.
