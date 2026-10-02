---
spec: 0213-a-test-suite-that-does-not-flake
status: archived
created: 2026-10-01
surfaces: [backend]
archived: "2026-10-02"
source_slug: 0213-a-test-suite-that-does-not-flake
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; the only blocked rows (10, 13) need an open Pull Request, and row 10 proves the ETXTBSY fixture change on Linux, which the PR's CI Verification gate runs
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 4bd6116ffa23b780f7faa5b93aaa81f4d24f5e62
---


# A test suite that does not flake

Six tests fail now and then in `make verify` or in the CI Verification gate,
then pass on rerun. Each failure parks a delivery that was otherwise ready, and
the operator retries it by hand (interventions 59 and 60 on 2026-10-01). The
adopted Backlog Entry
([references/2026-09-30-time-bound-tests-and-a-leaked-detached-child.md](references/2026-09-30-time-bound-tests-and-a-leaked-detached-child.md))
names three of the tests, and its 2026-10-01 addendum names three more. The
authoring session reproduced five of the six, and reading the code gives the
cause of each:

- **A wait that ends at the wrong event.** The linger subtest of
  `TestBatchClosesOnCountLingerAndImmediate` (`internal/store`) waits until
  the pending batch is empty. The linger flush empties the batch before its
  commit lands, so the test can read zero committed events. Under load, 74 of
  8,000 iterations failed with `expected 1 committed event after linger,
  got 0`.
- **A real-time allowance shorter than the work before it.**
  `TestTaskBudgetReasonNamesTheSettlementThatRenewedIt` (200 ms) and
  `TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement`
  (250 ms) in `internal/daemon` give the first Task a real allowance. Under
  load, the allowance can expire before that Task settles. 9 and 5 of 1,920
  iterations failed with `watched work ended: ... BudgetExceeded`.
- **A script executed the moment it is written.**
  `TestDoctorAcceptsAnAdapterAtTheFloor/codex` (`internal/cli`) writes a shell
  script and executes it at once, while parallel tests fork. On Linux a child
  forked in that window keeps the script open for writing, and `execve`
  refuses with `ETXTBSY` (golang/go#22315). The Doctor Command reported
  `did not prove required package lineage` in CI run 36903640400. ADR-0125
  already forbids the shape. Eight more written executables in non-governed
  `internal/cli` test files share it.
- **A shared fixture that outlives the tests that read it.**
  `TestAssetsSyncProvenanceAndPreMutationRefusals/*` (`internal/baseline`)
  copies one on-disk template that lives in the system temporary directory for
  the whole package run. No test or cleanup in the process removes it. On
  2026-10-01 the maintainer cleared system temporaries during the 0203 gate
  run, and every copy then failed with `open .git/objects/05: no such file or
  directory`. Removing the template during a run reproduces that message in
  four of four consumer tests. Without that removal, 200 stressed iterations
  of each Assets Sync test produced no failure.
- **A detached child whose wait outlives its test.**
  `TestRunImplementDetachSurvivesCallerProcessGroupKill` releases its detached
  child only from a `t.Cleanup`. When the test binary dies (timeout, `Ctrl-C`,
  cancelled job), cleanup never runs. The child and its fake ACPX then poll for
  a release file forever. Three of three interrupted runs left both processes
  alive. Spec 0103 reaped the survivors of `detach_test.go` at teardown but did
  not bound their wait on the test binary, and it never reached this test.

This is a test-only fix with no product behavior change. No production code
changes, because no test exposes a production race.

## Project Constraints

- Identifier strategy: not applicable — no identifier is created or changed;
  the Spec edits tests and test helpers. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network call or
  HTTP surface is touched; every fixture is local. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — under ADR-0125, "Fixtures are
  therefore compiled once", never written as scripts, and Core Feature 3
  brings the `internal/cli` adapters, fake ACPX and hooks into line with it.
  ADR-0213 (this Spec) makes every fixture process end with the test binary
  that started it. Under ADR-0126, "the guard runs per package", and the death
  tests keep their inner test binary inside the outer test's temporary
  directory. ADR-0028 owns detach, ADR-0098 owns batch durability, and
  ADR-0164 with ADR-0158 own the Run Budget whose renewal these tests examine;
  their behavior stays unchanged. ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check this
  Spec's consistency by citation and receipt. The gate is bound by ADR-0080,
  ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167. ADR-0182 does
  not apply, because no Task settlement fact changes, and ADR-0184 does not
  apply, because no command surface changes. ADR-0030 (opt-in agent run logs),
  ADR-0096, ADR-0097, ADR-0194, ADR-0195 and ADR-0210 (the gate's
  mechanical stage, row carry and evidence snapshots) and ADR-0170 (carry of
  already completed Tasks) do not apply, because no log,
  gate stage, carry rule or Task Carry-Forward behavior changes. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `GovernedPath` is empty, measured with a `go test -overlay`
  probe that wrote nothing. The governed `internal/cli/cli_test.go`, the
  Makefile, CI workflows, lint configuration and `go.mod` stay untouched.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- Each of the six named tests passes every iteration of a repeated run under
  load, at counts where it failed before.
- No `internal/cli` test executes a file the test process wrote. A guard fails
  on any new executable a test writes anywhere under `internal/`, outside a
  named residual inventory.
- No process a test starts outlives the test binary that started it.
- Removing the system temporary directory's contents during the package run
  no longer breaks the Assets Sync tests that have already started.

## Core Features

1. **Waits that end at the committed event.** The linger subtest waits until
   the linger commit is readable, bounded by the test deadline.
2. **A Run Budget proved on a fake clock.** Both short-budget tests take every
   expiry decision from the injected clock. The renewing settlement is ordered
   after the stalled Task starts, so no real allowance races the work.
3. **Script fixtures that are the compiled test binary.** An `internal/cli`
   script fixture is a symlink to the compiled test binary, with its script
   body in a non-executable sidecar, which `TestMain` runs through `/bin/sh`.
   The adapter, ACPX and hook fixtures of non-governed `internal/cli` test
   files use it. A repository guard names every test that writes an
   executable outside the residual inventory (ADR-0125).
4. **Fixture processes that end with their test binary.** The fake ACPX
   release wait and the `detach_test.go` survivor children watch the test
   binary's process and exit when it is gone. The Implement detach test reaps
   its detached child by the Owner PID its Run records (ADR-0213).
5. **An Assets Sync template held in memory.** `TestMain` builds the template
   once, without a `testing.T`, and reads it into memory. It then removes the
   on-disk build, so each test writes its own copy into its own temporary
   directory.

## Non-Goals / Out of Scope

- Any production change: the Run Budget watchdog, the journal writer, detach
  and the Doctor Command keep their behavior.
- The governed `internal/cli/cli_test.go` (`fakeACPXVersionCommand`), and the
  written executables outside `internal/cli` (`internal/baseline`,
  `internal/daemon`, `internal/preflight`, `internal/store`). They are named in
  the guard's residual inventory and left for a follow-up.
- A repository-contract test under the `repocontract` tag, which would need a
  Makefile change.
- Makefile, CI, `-timeout` or suite-budget changes, and retries of any kind.

## Success Metrics

1. The linger subtest passes a run of 8 parallel workers × 10 rounds ×
   `-count=50 -cpu 1,4` with no failure (before: 74 of 8,000 iterations
   failed).
2. Both Run Budget tests pass 8 workers × 6 rounds × `-count=20 -cpu 1,4`
   beside ten CPU burners with no failure (before: 9 and 5 of 1,920
   iterations failed). They also pass when the cycle's real elapsed time
   exceeds the configured maximum.
3. The guard reports zero executable writes in non-governed `internal/cli`
   test files, fails on a fixture tree that writes one, and holds the residual
   inventory at its recorded count.
4. After the test binary of the Implement detach test is killed with
   `SIGKILL`, its detached child and fake ACPX are gone within the test
   deadline (before: 3 of 3 interrupted runs leaked both processes).
5. No `assets-sync-template` directory exists once `TestMain` has built the
   template. The Assets Sync tests pass repeated runs, and a test's copy is a
   valid Git repository.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The CI log of run 36903640400 attempt 1 (PR #317), read with
  `gh run view 36903640400 --attempt 1 --log-failed`. It shows the
  `adapter_floor_test.go:50` failure quoted above; the rerun passed.
- Go issue golang/go#22315 (<https://github.com/golang/go/issues/22315>). It
  states that "the fd being written by one thread can leak into the forked
  child of a second thread", which leaves the file open for writing until
  that child execs, so another exec of it fails with `ETXTBSY`. The local
  Darwin probe could not reproduce this (0 of 2,400 write-then-exec
  attempts), because Darwin allows such an exec. The Linux rate rests on the
  CI evidence.
- The operator's intervention log (entries 59 and 60) and the session memory
  of 2026-10-01, which record the cleared system temporaries during the 0203
  gate run and the CI rerun of the Doctor flake.
- The v0.22.0 queue Finding,
  `docs/history/findings/2026-09-30-the-v0-22-0-queue-needed-twelve-manual-interventions.md`,
  section 4, which measured the 266 ms Run Budget, the linger failure and the
  five-hour orphan.

## Decisions

- Wait on the event the assertion reads, bounded by the test deadline, never
  on a proxy that changes earlier.
- Prove budget expiry on the injected clock, and order the renewal after the
  stalled Task starts. Injecting a timer into production for this was
  rejected, because the watchdog has no race.
- Make script fixtures compiled-binary re-executions, as ADR-0125 requires.
  Holding `syscall.ForkLock` around the write was rejected, because ADR-0125
  already decides the class and the re-execution costs about 10 ms.
- Bound every fixture wait on its test binary's life, and reap by recorded
  PID. See ADR-0213.
- Keep the Assets Sync template's speed, and move it from disk into memory,
  rather than rebuilding it per test.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
