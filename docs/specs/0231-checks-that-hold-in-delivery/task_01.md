---
task: task_01
spec: 0231-checks-that-hold-in-delivery
status: pending
type: test
complexity: low
---

# Task 01: The detach fixture tests count a process-group member that is already exiting as ended

## Overview

The finding of 2026-10-05 ([the survivor test meets EPERM from an exiting owner](references/2026-10-05-the-survivor-test-meets-eperm-from-an-exiting-owner.md))
is that the darwin reading of a fixture's process group reports a member that
is already exiting as live, while the kernel already refuses to count it as a
signal target, so the survivor test's cleanup fails on `EPERM`. This Task
makes the reading exclude exiting members and adds a test that catches the
window and proves the reading. No production code changes.

## Requirements

1. MUST answer the finding of 2026-10-05 named in the Overview: the darwin
   `detachFixtureGroupLiveMembers` MUST exclude a member whose `kinfo_proc`
   `p_flag` carries `P_WEXIT` (0x2000 in Darwin's `<sys/proc.h>`), as it
   already excludes `p_stat` `SZOMB`, named by a constant like the zombie
   state (TechSpec Invariant 9).
2. MUST add the darwin test `TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded`
   to the darwin file: it starts short-lived children in their own process
   groups and, reading `kern.proc.pgrp` directly, polls under `testwait`
   until it observes a group whose only member is not a zombie and carries
   `P_WEXIT` while `kill(-pgid, 0)` returns `EPERM`; at that moment the
   reading MUST report no live member. It MUST reap every child it starts,
   MUST fail at its deadline if it never observes the window, and MUST NOT
   skip.
3. MUST keep `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` proving
   that a running member is live and that a different start identity ends the
   fixture, unchanged.
4. MUST NOT change the Linux or other Unix readings, `killDetachFixtureGroup`,
   `detachFixtureEnded`, any detach production code, `internal/store`, or the
   survivor and death tests themselves.

## Subtasks

- [ ] Exclude exiting members in the darwin reading.
- [ ] Add the exiting-window test.

## Acceptance Criteria

- [ ] A group whose only member is exiting reads as ended; a running member
      still reads as live.
- [ ] The detach death, survivor, unreaped-member and exiting-member tests
      pass ten consecutive iterations each, with no skip.

## Context

- instruction: `docs/adr/0213-a-test-fixture-process-ends-with-the-test-binary-that-started-it.md`
- instruction: `internal/cli/implement_detach_teardown_test.go`
- interface: `internal/cli/detach_fixture_group_darwin_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded '` — expected: exit 0; before this Task the test does not exist, so no pass line is printed and the command fails.
- `out="$(go test -count=10 -v -run '^(TestRunImplementDetachSurvivesCallerProcessGroupKill|TestImplementDetachChildEndsWhenItsTestBinaryDies|TestDetachSurvivorEndsWhenItsTestBinaryDies|TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded|TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunImplementDetachSurvivesCallerProcessGroupKill TestImplementDetachChildEndsWhenItsTestBinaryDies TestDetachSurvivorEndsWhenItsTestBinaryDies TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded; do test "$(printf '%s\n' "$out" | grep -c -- "--- PASS: $name ")" = 10 || { printf 'want 10 passes of %s\n' "$name" >&2; exit 1; }; done; ! printf '%s\n' "$out" | grep -q -- '--- SKIP'` — expected: exit 0; before this Task the exiting-member test does not exist, so its ten passes are missing and the command fails.

## References

- `_prd.md` → Goals; Core Feature 1; Success Metric 1; Acceptance evidence
- `_techspec.md` → Interfaces; Invariant 9; Testing Approach; Build Order 1
- ADR-0213
