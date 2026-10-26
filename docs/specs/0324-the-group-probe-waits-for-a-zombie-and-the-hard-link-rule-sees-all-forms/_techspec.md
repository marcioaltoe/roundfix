---
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
prd: _prd.md
created: 2026-10-26
---

# The group probe waits for a zombie, and the Hard Link rule sees all forms — Technical Spec

## Executive Summary

Two test-suite defects, both in tests only. First,
`TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` ends its wait when
the darwin live-member reading reports none, a reading that also leaves out
a member still exiting; Darwin's `killpg1` counts an exiting member, so the
group probe that follows can succeed where the test expects `EPERM`. The wait
gains a second condition, the member's own zombie state, read through the
same sysctl, and a new test proves that reading never takes an exiting
member for a zombie. Second, the Hard Link rule inspects only calls through
named imports under `internal/`; it is widened to every reference, dot
imports, and every Go test file of the module.

The trade-off accepted: the unreaped-group test still polls, every 10 ms
under the test deadline, instead of blocking on a kernel event. Darwin's only
non-reaping wait, `waitid` with `WNOWAIT`, has no wrapper in the modules the
repository requires, and a raw system call number in a test costs more than a
poll whose end condition cannot race: once a held child is a zombie it stays
one until the test reaps it.

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local processes and local files
  only. Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0324, as `_prd.md` → Project
  Constraints states. Source: `docs/agents/domain.md`,
  `docs/agents/specific-repository.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms/_authorization.md`;
  bounded files: none governed; the record's `paths` is `[]`.

## Current behavior

Read on 2026-10-26 at `aa8b3db7`.

1. `internal/cli/implement_detach_teardown_test.go`
   `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` (unix build
   constraint; returns early outside Darwin and Linux) starts `cat` with
   `Setpgid`, closes its input, and polls through `testwait.Poll` ("fixture
   cat exited but unreaped") until `detachFixtureGroupLiveMembers(pid)` is
   empty, asserting `kill(pid, 0)` succeeds on every poll. It then asserts
   `kill(-pid, 0)` is `EPERM` on Darwin and success on Linux,
   `detachFixtureEnded` is true, reaps, and asserts `kill(-pid, 0)` is
   `ESRCH`.
2. `internal/cli/detach_fixture_group_darwin_test.go`
   `detachFixtureGroupLiveMembers` reads `kern.proc.pgrp` and keeps a member
   whose `P_stat` is not `SZOMB` (5) and whose `P_flag` lacks `P_WEXIT`
   (0x2000). `TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded` already
   documents the exiting window: while `P_WEXIT` is set and the state is not
   `SZOMB`, it accepts a nil group probe and keeps launching until the probe
   returns `EPERM`.
3. `internal/cli/detach_fixture_group_linux_test.go` keeps every `/proc`
   process of the group whose state is not `Z`; a Linux process becomes `Z`
   only when its exit completes, so the Linux poll already ends on a zombie.
   `detach_fixture_group_other_test.go` probes the group with a signal.
4. Darwin `killpg1` (see `_prd.md` → Acceptance evidence) skips `SZOMB`
   members and returns success when it counted any member, `EPERM` (POSIX
   `kill`) when none. The failure of 2026-10-25 and the 87 failures of
   2026-10-26 are probes that ran in the exiting window
   (`docs/references/2026-10-24-test-suite-audit.md`, addendum of
   2026-10-26).
5. Other readers of the same state, swept with
   `grep -rn 'P_stat' internal`: `internal/store/process_darwin.go`
   `liveProcessGroupMembers` (production; treats exiting and zombie alike as
   not live, and `SignalProvenProcessGroup` accepts either probe result
   correctly) and `internal/suiteguard/residue_darwin.go` (counts an exiting
   process as gone, which is right for residue). Neither asserts a probe
   result in the exiting window; neither changes.
6. `internal/testfixture/hard_link_rule_test.go` `hardLinkViolations` walks
   `internal/` only, maps each import of `os`, `syscall` and
   `golang.org/x/sys/unix` by name or alias, ignores a `.` import as a name
   that never appears as a receiver, and inspects only `*ast.CallExpr` whose
   `Fun` is a selector on an unshadowed identifier. `TestNoTestCreatesAHardLink`
   has a `repository` subtest and ten seeded subtests. Go test files outside
   `internal/`: 17 under `skills/` and one vendored skill asset under
   `.agents/skills/` (a hidden directory the `go` command never builds).

## System Architecture

No production package changes. task_01 changes the `internal/cli` detach
fixture test files; task_02 changes one `internal/testfixture` test file.
The two Tasks share no file.

### The zombie wait in the unreaped-group test

The unreaped-group test keeps holding `cat` unreaped, but its poll now ends
only when the live-member reading is empty and the member's own state is a
zombie, read by `detachFixtureMemberIsZombie` in the platform file of the
host. Once a held child is a zombie it stays one until the test reaps it, so
the probe that follows always sees the zombie-only group Darwin answers with
`EPERM`, whatever the load (API Contract 1).

### The exiting-member proof

A Darwin test catches a short-lived child inside its exiting window, the
state the old wait mistook for the end, and checks that the zombie reading
says no there and yes once the child becomes a zombie, when the group probe
is `EPERM`. A reading that counts the exiting flag as a zombie fails it
(API Contract 2).

### The widened Hard Link rule

The rule keeps its message and its `repository` subtest, but inspects every
reference instead of every call, follows dot imports of `os`, `syscall` and
`golang.org/x/sys/unix`, and walks every test file the `go` command builds,
so `cmd/` and `skills/` are covered and `testdata`, hidden and underscore
directories are not (API Contract 3).

## Implementation Design

### Interfaces

```go
// detach_fixture_group_{darwin,linux,other}_test.go
// detachFixtureMemberIsZombie reports whether pid is an exited, unreaped
// process. Darwin: kern.proc.pid P_stat == SZOMB (an exiting P_WEXIT
// process is not). Linux: /proc/<pid>/stat state "Z". Other Unix: error.
// An error is a failed reading, never "not a zombie".
func detachFixtureMemberIsZombie(pid int) (bool, error)
```

```go
// hard_link_rule_test.go
// hardLinkViolations walks root, skipping directories named testdata or
// starting with "." or "_", parses each *_test.go, and returns
// "<rel>:<line>: <pkg>.<Name> creates a hard link; copy or use a symbolic link"
// for each reference, in source order per file.
func hardLinkViolations(root string) ([]string, error)
```

### Data Models

None. No persisted format changes.

### API Contracts

1. API Contract: the zombie wait. In
   `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` the poll
   "fixture cat exited but unreaped" returns ready only when
   `detachFixtureGroupLiveMembers(pid)` is empty and
   `detachFixtureMemberIsZombie(pid)` is true; a reading error fails the
   test; its observation names the live members and the zombie answer. The
   per-poll `kill(pid, 0)` success check, the Darwin `EPERM` and Linux success
   expectations, the `detachFixtureEnded` assertion, the reap and the `ESRCH`
   expectation keep their bytes. A comment states why an empty live reading
   is not enough.
2. API Contract: the exiting-member proof.
   `TestDetachFixtureZombieReadingRejectsAnExitingMember` in
   `detach_fixture_group_darwin_test.go`, `t.Parallel()` first, launches
   `/usr/bin/true` with `Setpgid` under a `testwait.Poll` bounded by
   `testwait.Bound(t)`, reaps each child before the next launch, and spins on
   `kern.proc.pid` until it reads `P_WEXIT` without `SZOMB`. It then calls
   `detachFixtureMemberIsZombie` and reads the raw state again: when the
   second read is `SZOMB` the window closed and the launch proves nothing;
   otherwise a `true` answer fails the test naming the state and flags. After
   a `false` answer it polls until the reading is `true` and asserts
   `kill(-pid, 0)` is `EPERM`.
3. API Contract: the widened Hard Link rule. `hardLinkViolations` walks the
   repository root instead of `internal/`, skipping `testdata` and
   directories whose name starts with `.` or `_`, as the `go` command does.
   It reports: every `*ast.SelectorExpr` `X.Link` (and `X.Linkat` unless `X`
   is `os`) whose `X` is an unshadowed identifier naming such an import,
   called or not; and, in a file with a `.` import of such a package, every
   identifier `Link` or `Linkat` with no resolved object that is not a
   selected member, a function or method name, a field or parameter name, a
   struct literal key, or a label, attributed to the first dot-imported
   package that exports it (`os` exports no `Linkat`). The message format
   and the `repository` subtest are unchanged; `TestNoTestCreatesAHardLink`
   keeps its ten seeded subtests byte for byte and gains the seeds of
   `_prd.md` → Success Metric 4.

### Surface Transcripts

None. No command surface changes.

### Invariants

1. The unreaped-group test never asserts a group probe result before
   `detachFixtureMemberIsZombie` returned true for its member.
2. `detachFixtureMemberIsZombie` returns false for a process whose state is
   not `SZOMB` (Darwin) or `Z` (Linux), including an exiting one.
3. Every expectation of the unreaped-group test other than its wait keeps
   its bytes.
4. The Hard Link rule reports each seeded form with its file and line and
   never a seeded non-reference.

## Coverage Map

- Goal 1 (no failure from the exiting window) → The zombie wait in the unreaped-group test (API Contract 1)
- Goal 2 (four assertions kept) → The zombie wait in the unreaped-group test (API Contract 1, Invariant 3)
- Goal 3 (zombie reading proved) → The exiting-member proof (API Contract 2)
- Goal 4 (every hard-link form) → The widened Hard Link rule (API Contract 3)
- Core Feature 1 → API Contract 1
- Core Feature 2 → API Contract 2
- Core Feature 3 → API Contract 3
- Success Metric 1 → API Contract 1, the QA load harness
- Success Metric 2 → API Contract 2
- Success Metric 3 → API Contract 1, Invariant 3
- Success Metric 4 → API Contract 3, Invariant 4
- Success Metric 5 → API Contracts 1-3, task_01 and task_02 Verification

Deviations from the PRD: None.

## Absence Evidence

1. `detachFixtureMemberIsZombie` — command: `grep -rn 'MemberIsZombie' internal cmd skills`;
   result: no match at `aa8b3db7`. `grep -rln 'P_stat' internal` lists only
   the darwin fixture file, `internal/store/process_darwin.go` and
   `internal/suiteguard/residue_darwin.go`, none of which reads one process's
   zombie state apart from its exiting flag.
2. `TestDetachFixtureZombieReadingRejectsAnExitingMember` — command:
   `grep -rn 'ZombieReading' internal`; result: no match at `aa8b3db7`.

## ADR Alternatives

- ADR-0324: the decision ("a repository rule in `internal/testfixture` fails
  on any hard link created by a test") stands; the widening only removes
  forms the first rule could not see. None of its rejected alternatives
  (removing links only in the top-level process, retrying a killed fake,
  disabling the host check) is taken.

## Integration Points

None outside the repository. The Linux Lane runs the changed `internal/cli`
and `internal/testfixture` tests in `golang:1.26`; the Linux reading of
API Contract 1 runs there.

## Testing Approach

The seams are the existing test files; no production seam.

- task_01: the unreaped-group test with its new wait, and the new
  exiting-member test, run by name with `-v`; `GOOS=windows`, `GOOS=linux`
  and `GOOS=freebsd` vet of `./internal/cli`.
- task_02: `TestNoTestCreatesAHardLink` with `-v`, every subtest named.
- QA, on temporary copies made with `git archive`:
  - the load harness of the audit addendum, three rounds at the audited head;
  - sabotage: make the darwin `detachFixtureMemberIsZombie` also return true
    when `P_WEXIT` is set →
    `TestDetachFixtureZombieReadingRejectsAnExitingMember` fails "read as a
    zombie" (its input is the observed exiting window; measured at
    authoring, it fails within the first launches);
  - sabotage: change the Darwin expectation of the unreaped-group test from
    `EPERM` to success → that test fails, proving the zombie probe returns
    `EPERM`;
  - sabotage: drop the recording of dot imports in `hardLinkViolations` →
    the seeds "dot import of os", "dot import of unix Linkat" and
    "dot-imported syscall Link as a value" fail;
  - sabotage: inspect only selectors that are a call's `Fun` → the seed
    "os Link as a function value" fails;
  - sabotage: walk `internal/` only → the seeds "a test outside internal"
    and "a test under skills" fail;
  - sabotage: stop skipping struct literal keys → the seed "dot import and a
    field, method, key and label named Link" fails.

## Build Order

1. task_01: the zombie wait and the exiting-member proof (API Contracts 1
   and 2).
2. task_02: the widened Hard Link rule (API Contract 3); independent of
   step 1.
3. task_03: the QA gate (depends on: 1, 2).

## Risks & Considerations

- The exiting window is short; the new test may launch several children
  before it observes one. The existing exiting-member test has the same
  shape and passed every authoring run within three launches; both are
  bounded by the test deadline, never by a short window.
- A bare `Link` identifier in a file with a dot import of `os`, `syscall`
  or `unix` and no local declaration cannot refer to anything else: a
  package-level `Link` would conflict with the dot import and fail to
  compile.
- The rule now also reads `cmd/` and `skills/` test files; none creates a
  hard link at `aa8b3db7` (prototype `repository` subtest passed).

## Glossary

None.

## Decisions

- Wait for the zombie by reading `kern.proc.pid`, not through a raw
  `waitid` call (`_prd.md` → Decisions).
- Leave `detachFixtureGroupLiveMembers` unchanged: its exclusion of exiting
  members is what `detachFixtureEnded` and the exiting-member test rely on.
- `go/ast` with parser object resolution, as the rule already uses; no
  `go/types`.
- No new ADR.
