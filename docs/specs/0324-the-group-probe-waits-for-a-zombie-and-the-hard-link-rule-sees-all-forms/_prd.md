---
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
status: active
created: 2026-10-26
surfaces: [backend]
---

# The group probe waits for a zombie, and the Hard Link rule sees all forms

Minimal PRD for a test-suite fix, minted by `write-techspec`; no product
behavior changes. On macOS `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded`
fails under load because it probes its process group while the member is
still exiting, which Darwin counts as a live member, instead of waiting for
the zombie the test names. After this Spec the test waits until the member's
own state is a zombie before it asserts `EPERM`, a new test proves that the
zombie reading never takes an exiting member for a zombie, and the Hard Link
rule also reports a dot-imported `Link` and a `Link` taken as a function
value, in every Go test file of the module.

## Sources

Two Backlog Entries, written and adopted at birth by this Spec's planning
change (see `references/_index.md`):

1. [`references/2026-10-26-the-unreaped-group-test-fails-under-the-full-suite-on-macos.md`](references/2026-10-26-the-unreaped-group-test-fails-under-the-full-suite-on-macos.md),
   `type: fix`: the Spec 0317 QA gate of 2026-10-25 failed on this test
   ("got <nil>, want EPERM") while the whole `internal/cli` package ran.
2. [`references/2026-10-26-the-hard-link-rule-misses-a-dot-import-and-a-function-value.md`](references/2026-10-26-the-hard-link-rule-misses-a-dot-import-and-a-function-value.md),
   `type: fix`: the Spec 0311 pre-PR review finding F1 of 2026-10-24,
   dismissed then for scope.

Both are test-suite defects found by the cycle AP triage and share the
suite-audit context of `docs/references/2026-10-24-test-suite-audit.md`. At
the authoring base `docs/backlog/` and `docs/findings/` hold no entry.

The operator's third ask, a per-test count of Git starts for `internal/cli`,
is not measured anywhere; the planning change records it as a follow-up in
the audit reference and no Task measures it.

## Prerequisites

None in `requires`. Authored against `aa8b3db7`, with Specs 0319, 0320, 0321
and 0322 active under `docs/specs/`; this Spec delivers after them and shares
no file with them.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; only test code changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  network access; the tests start local processes and parse local files.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0324 (a test never gives an
  executable a second hard link; a repository rule in `internal/testfixture`
  fails on any hard link a test creates) is the rule task_02 widens. The
  platform-split files task_01 changes keep compiling for every release
  platform under the repository's Windows compile contract, a rule task_01
  proves with `GOOS` vets rather than a decision it changes. Source:
  `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized; the record carries `paths: []`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms/_authorization.md`;
  bounded files: none governed; the record's `paths` is `[]`.

## Goals

- The unreaped-group test fails only when Darwin's group probe of a zombie
  member returns something other than `EPERM`, never because the probe ran
  while the member was still exiting, whatever the load.
- The test keeps its four assertions: `EPERM` on Darwin and success on
  Linux for the unreaped group, the group read as ended, and `ESRCH` after
  the reap.
- A test proves that the zombie reading the wait uses answers "not a zombie"
  for a process observed while it exits.
- The Hard Link rule reports every reference to a hard-link function of
  `os`, `syscall` or `golang.org/x/sys/unix` in any Go test file of the
  module, called or taken as a value, through a named, aliased or dot import.

## Core Features

1. The zombie wait. Each platform file of `detachFixtureGroupLiveMembers`
   gains `detachFixtureMemberIsZombie(pid)`: on Darwin the process's own
   `P_stat` is `SZOMB`, on Linux its `/proc/<pid>/stat` state is `Z`, and on
   other Unix systems it returns an error. The unreaped-group test's poll
   ends only when the group has no live member and that reading says zombie.
2. The exiting-member proof. A new Darwin test launches short-lived children
   until it observes one exiting but not yet a zombie, asserts the zombie
   reading says no, then waits until it says yes and asserts the group probe
   returns `EPERM`.
3. The widened Hard Link rule. `TestNoTestCreatesAHardLink` walks every
   `_test.go` file under the repository root that the `go` command would
   build, reports selector references to `Link`, and to `Linkat` of
   `syscall` and `unix`, whether or not they are called, and reports an
   unshadowed bare `Link` or `Linkat` identifier in a file that dot-imports
   one of those packages; it ignores fields, methods, struct literal keys,
   labels and selected members named `Link`.

## Non-Goals / Out of Scope

- Changing `detachFixtureGroupLiveMembers`, `detachFixtureEnded`,
  `TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded`, or any production
  process-control code: `store.SignalProvenProcessGroup` already treats a
  successful group signal as signalled and `EPERM` with no live member as
  ended.
- A raw `waitid` system call: Darwin has no wrapper in the Go modules the
  repository requires, and `wait4` with `WNOWAIT` reaps on Darwin (measured).
- Measuring per-test Git starts in `internal/cli` or refactoring Git-bound
  tests (recorded as follow-up in the audit reference).
- The Makefile, the CI workflows, `go.mod`, `go.sum`, guides, skills,
  `CONTEXT.md`.

## Success Metrics

1. Success Metric: the load harness of the audit reference (eight concurrent
   copies of the compiled `internal/cli` test binary, each running the
   unreaped-group test 2,000 times, beside eight CPU-bound loops) passes
   three rounds, 48,000 runs, with no failure at the audited head on macOS.
   At authoring the Delivery Base failed 87 of 48,000 at load averages from
   49 to 55, and a prototype of the fix failed none.
2. Success Metric: the new exiting-member test passes, and a zombie reading
   that also counts an exiting member as a zombie makes it fail naming the
   exiting state.
3. Success Metric: the unreaped-group test still asserts `EPERM` on Darwin,
   success on Linux, the group read as ended and `ESRCH` after the reap:
   removing any one of those assertions is a visible diff in the QA scope
   audit, and replacing the Darwin `EPERM` expectation by success makes the
   test fail.
4. Success Metric: `TestNoTestCreatesAHardLink` passes on the repository;
   its ten existing seeded cases pass unchanged; new seeds for a
   dot-imported `os` `Link` call, a dot-imported `unix` `Linkat` call,
   `os.Link` taken as a value twice, a dot-imported `syscall` `Link` value,
   and a test under `cmd/` and under `skills/` are each reported with file
   and line; a local `Link`, a field, method, key and label named `Link`, a
   dot-imported `os` that has no `Linkat`, and tests under `testdata`,
   hidden and underscore directories are not.
5. Success Metric: `GOOS=windows go vet` passes for `./internal/cli` and
   `./internal/testfixture`, and `GOOS=linux go vet ./internal/cli` passes.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- Apple's published XNU source, `bsd/kern/kern_sig.c` `killpg1`
  (https://github.com/apple-oss-distributions/xnu/blob/main/bsd/kern/kern_sig.c),
  read through Exa on 2026-10-26: the group iteration keeps a member when
  "`return !(p->p_flag & P_SYSTEM) && p->p_stat != SZOMB;`", and the call
  ends with "`error = (karg.nfound > 0 ? 0 : (posix ? EPERM : ESRCH));`". An
  exiting member that is not yet `SZOMB` is counted and the probe succeeds;
  a group whose only member is a zombie has none counted and fails with
  `EPERM`. The same passage is copied in
  `docs/references/2026-10-24-test-suite-audit.md`, committed by the planning
  change before any Task runs.
- The Spec 0317 QA record in the adopted Backlog Entry, written by another
  Spec's gate before this Spec existed: "probe unreaped fixture group 6892:
  got <nil>, want EPERM" under the full `internal/cli` package, passing alone.

## Unreachable Acceptance

- criterion: no Run delivered after this Spec fails
  `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` in its
  Verification or its QA gate.
  reason: it needs later Runs' Verification logs, which no Verification or
  QA row may read.
  satisfied-by: the operator's reading of the queue records and Run logs of
  the Specs delivered after this one, recorded session evidence.

## Validation

The decisions were checked on 2026-10-26 against these sources.

- Reproduction and mechanism, measured at `aa8b3db7` and recorded in the
  addendum of `docs/references/2026-10-24-test-suite-audit.md`: the load
  harness failed the test 1 time in 16,000 at load average 15 and 87 times
  in 48,000 at 49 to 55; a scratch program showed success only from probes
  of an exiting member and `EPERM` from every one of 70,991 probes of a
  zombie.
- A prototype of both Tasks in the authoring clone: the harness passed
  48,000 runs with no failure; the new exiting-member test passed and failed
  when the zombie reading also accepted an exiting member; the widened rule
  passed on the repository with every seed, and failed its dot-import seeds
  when dot imports were ignored and its function-value seed when only calls
  were inspected; `go test ./internal/testfixture` (every repository rule)
  passed, and `GOOS=windows go vet` passed for both packages.
- Secondbrain, read only: `qmd query "macOS kill process group EPERM zombie
  unreaped test flake" --all --files --min-score 0.3` returned the decision
  records on proven signals and on rebuilt test executables, the digest of
  the macOS kill of a signed binary, the audit reference and ADR-0324; `qmd query "roundfix hard link rule test dot import
  os.Link" --all --files --min-score 0.3` returned a digest on worktrees
  under a symbolic link, the README, `docs/agents/specific-repository.md`,
  the Spec 0311 record and the decision record on rebuilt test
  executables. None records the exiting window or a
  dot-import decision.
- Exa: the XNU `killpg1` source quoted above, and a published report of the
  same behavior in another project (NousResearch hermes-agent pull request
  104707: "XNU's `killpg1` skips `SZOMB` members and returns EPERM when no
  eligible member remains").

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AQ batch 2 brief for Spec 0324 of 2026-10-26, the maintainer's cycle AQ answers quoted in `_authorization.md`, and the characterization above.

Taken from those sources, with no further interview:

- Fix the test's precondition, not its assertion: the poll waits for the
  member's own zombie state; the `EPERM` expectation stays exact.
- Read the state through the sysctl the file already uses, polled under the
  test deadline, instead of a raw `waitid` system call number.
- Widen the Hard Link rule with `go/ast` only: references, not only calls;
  dot imports; every test file the `go` command builds. No `go/types`.
- No new ADR: ADR-0324's decision already covers every hard link a test
  creates, and the test fix decides nothing beyond its own precondition.

## Open Questions

None.
