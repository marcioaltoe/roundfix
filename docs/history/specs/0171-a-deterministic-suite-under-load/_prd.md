---
spec: 0171-a-deterministic-suite-under-load
status: archived
created: 2026-09-25
surfaces: [backend]
archived: "2026-09-25"
source_slug: 0171-a-deterministic-suite-under-load
---


# A deterministic suite under load

QA gates run the repository Verification on a machine that is also running
other Runs, and in the week to 2026-09-25 four Runs were spent on tests that
pass in isolation. Spec 0163's gate failed on
`TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks`
("timed out waiting for 2 Agent starts; got 1" after 91 s) and then on
`TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork`
("expected Clean exit, got 1"). Spec 0164's gate failed on
`TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`.
A documentation-only `make verify` failed on `TestBootstrapSerializesAcrossSiblings`.
All four pass three times out of three in isolation, on the Spec branches and on
`main`.

Reading the tests shows three shapes of the same defect:

- **A product timeout the test does not examine.** Three Implement tests
  configure `bootstrap_timeout: 1s` and the worktree tests pass one- and
  five-second `BootstrapSpec.Timeout` values, although none of those tests is
  about the timeout. A loaded machine that takes longer than a second to start `sh`
  fails the Run, and the test reports a wrong exit code.
- **A wait that watches only one channel.** `waitImplementAgentStarts` in
  `internal/cli/implement_test.go` waits 90 s for an Agent start and never
  looks at the Implement Command it started, so a command that already ended
  is reported as a slow Agent 90 s later, with its real error unread. The
  Daemon helpers in `internal/daemon/task_engine_test.go` are already bound to
  `t.Deadline()`, but they too ignore an early `TaskCycle` return, so the same
  case waits out the whole package timeout.
- **A fixed wall-clock budget.** `implementWaitBudget` (90 s),
  `detachStartupBudget` (30 s) and the 90-second waits beside them measure the
  machine, not the code; `TestBootstrapSerializesAcrossSiblings` has no bound at
  all on its result loop, so a stuck sibling hangs until the package times out.

A reproduction attempt on 2026-09-25 ran three `internal/daemon`, three
`internal/worktree` and two `internal/cli` test binaries at once beside ten CPU
burners (load average 36): every test passed. The failures are rare, so this
Spec removes the shapes that make them possible and proves the removal, rather
than claiming a reproduction it does not have.

Two defects make the same suite depend on things outside the machine. Tests of
`fetch`, `resolve`, `watch` and `implement` reach `gh api` through
`maybeReportVersionFreshness` in `internal/cli/upgrade.go`, because the checked-in
version is not a development version and only `internal/cli/upgrade_test.go`
injects a fake lookup; Spec 0165's gate hit `api.github.com` again. And
`provisionFakeAdapter` in `internal/agent/acpx_runner_test.go` symlinks
`os.Args[0]` as written, so a compiled test binary started by a relative path
creates adapters that cannot execute.

## Project Constraints

- Identifier strategy: not applicable — no identifier is created or changed;
  the Spec edits tests and adds a test helper package. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the release lookup calls GitHub through
  `gh api`; this Spec removes that call from every test process and adds none.
  No credential is read. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0089 makes code under test take its
  environment explicitly, which the release lookup must follow; ADR-0125 keeps a
  spawned fixture a compiled binary, which the adapter repair preserves;
  ADR-0126 makes each spawning package install the suite guard; ADR-0034 and
  ADR-0056 own the bootstrap and Verification Capacity behavior these tests
  examine, which stays unchanged. ADR-0093 checks Spec consistency by citation,
  ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a path
  governed once bounded, ADR-0155 makes the `qa` Task declare the matrix and
  ADR-0156 makes a declared promise name a consuming Task. This Spec's gate is
  bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `GovernedPath` is empty; `internal/cli/cli_test.go`,
  `docs/references/coverage-record.json` and the Makefile stay untouched.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A test that passes in isolation passes on a loaded machine.
- A test that still fails names its cause, not a timeout.
- No test process performs the live release lookup.
- The adapter fixture works however the test binary is started.

## Core Features

1. **Deadline-bound waits that watch the work.** A test helper package,
   `internal/testwait`, bounds every wait by the test's own deadline and fails a
   wait at once, with the work's result, when the work it watches ends first.
2. **Implement tests wait on events.** `internal/cli/implement_test.go` waits
   through that package, watches the Implement Command it started, and
   configures no product timeout shorter than its deadline unless the timeout is
   what the test examines.
3. **Daemon and worktree tests wait on events.** `internal/daemon/task_engine_test.go`
   and `internal/worktree/worktree_test.go` do the same, including every round
   of `TestBootstrapSerializesAcrossSiblings`.
4. **No live release lookup in tests.** Every `internal/cli` test process
   answers the release lookup offline, every built binary a test runs for an
   operational command finds a fresh version cache, and guard tests fail when a
   live lookup could happen.
5. **An adapter fixture that survives a relative path.** The fake adapter links
   an absolute path to the compiled test binary.

## Non-Goals / Out of Scope

- Changing any product timeout, Run behavior, command output or exit code. The
  one production change, serializing Git worktree administration per
  repository (task_07), changes none of them.
- The attach budget `attachDetachBudget`, which Spec 0124 owns.
- Other network boundaries: review sources and pull request calls are already
  injected per test through command dependencies.
- Moving every in-process test onto a test-owned HOME.
- Makefile, CI, `-timeout` or suite-budget changes.

## Success Metrics

1. None of the three named test files waits on a fixed duration: no
   `time.After`, `time.NewTimer` or named wall-clock budget remains in them, and
   no bootstrap timeout shorter than the test deadline remains where the timeout
   is not under test.
2. The four tests named in the adopted Backlog Entry, and the bootstrap tests
   beside them, pass every iteration of a `-count` and `-cpu 1,4` stress run.
3. A wait whose work ended first fails within one second and names the work's
   result; a wait that times out fails no earlier than the test deadline minus
   the helper's margin.
4. `fetch`, `resolve`, `watch` and `implement` run by a test record the offline
   lookup's answer in their HOME, and a built binary leaves a seeded fresh cache
   byte-identical.
5. The fake adapter runs when the compiled test binary is started by a relative
   path, and still runs when it is started by an absolute path.

## Decisions

- **Bound by the deadline, fail on the event.** A fixed budget either fails a
  slow machine or hides a stuck one; `t.Deadline()` is the only bound that
  cannot fail earlier than `go test` itself would, and watching the work's end
  turns a stuck wait into a named failure.
- **Raise the timeouts the tests do not examine.** A timeout that is the
  subject of a test stays, written so load can only make its outcome more
  certain; every other one follows the deadline.
- **Assign the offline lookup once, at package initialisation.** The finding
  warns that a mutable package-level switch races parallel tests. The default is
  assigned in a test file's `init`, before any test starts, and never
  reassigned, as `newOutcomeNotifier` already is in the same package; per-test
  fakes keep using command dependencies. A development version for every test
  was rejected because it would stop the freshness check from being tested at
  all.
- **Seed the cache for built binaries.** A built binary cannot take a test
  double, and ADR-0089 rejects a production-only switch; a fresh cache in the
  binary's HOME exercises the real skip path.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a wait that returns on the
happy path but still sleeps out its budget when the work ends early, or a guard
that passes because the network happened to be down, would pass the positive
case alone.

## Research basis

The four failures, their messages and the isolation reruns are recorded in the
adopted Backlog Entry, from the QA reports of Runs
`run_20260925T153433Z_8603eb3e79157622`, `run_20260925T155346Z_b801d0ee9b2a247a`
and `run_20260924T223414Z_54e6f85ae98b6550`. The two adopted Findings record the
release lookup and the adapter link. Go documents `testing.T.Deadline` as the
time at which the test binary will exceed its `-timeout`
(https://pkg.go.dev/testing#T.Deadline).

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
