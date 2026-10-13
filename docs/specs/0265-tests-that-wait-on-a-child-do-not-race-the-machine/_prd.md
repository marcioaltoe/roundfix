---
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
status: active
created: 2026-10-13
surfaces: [backend, docs]
---

# Tests that wait on a child do not race the machine

Spec 0259 accepts its fix only after ten consecutive clean full suites
(`go test -count=1 -parallel 16 ./...`, exit 0, no `signal: killed`,
`exited -1` or Helper Residue line). Its QA gate ran on a machine shared with
other Runs, at load averages from 16 to 28, and failed three times on tests whose
subject was sound and whose real-time budgets did not survive the load:

- Run `run_20261013T004148Z_5319d413bd50ad25`: `TestBudgetExceededRunRecordsItsReason`
  (`internal/daemon`), a real 500 ms Run Budget spent by Git work.
- Run `run_20261013T013741Z_3a6796ce5e63cbc3`: the two worktree admin-lock
  disposal tests (`internal/worktree`), a 250 ms context that also covers a Git
  fork.
- Run `run_20261013T014653Z_1fc0f49baff7f737` (QA report
  `qa-report-2026-10-12-03.md`, suite 3 of 10, load average 16.6):
  `TestRunDetachedCommandTimesOutWaitingForRunCreation` and
  `TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline`
  (`internal/cli/detach_test.go`). Both inject a 500 ms liveness budget that is
  not what they assert. The re-executed child wrote no liveness byte in time,
  the parent killed it (`signal: killed`), and both reported
  `Detached Run child produced no liveness signal within 500ms`. Production
  allows ten seconds for that step.

0259's two Corrective Tasks are spent: they fixed the first two cases on its
item branch. Under the operating rules a blocking Spec makes the suite
load-independent instead, and 0259 then merges `main` and runs its QA again.
This Spec fixes the class, not the case: it fixes the two detach tests, every
other test the authoring sweep found that gives a child process, a goroutine or
Git a short real-time budget it does not assert, and adds a repository rule so
the class cannot grow back. This is a test-only change. The only production
edits add a dependency with the production default where a test could not
reach a budget otherwise; no production budget value changes.

## Prerequisites

None. This Spec is delivered before Spec 0259 merges. It keeps out of its write
scope every file 0259's item branch changes, except `internal/cli/detach_test.go`,
`internal/cli/implement_test.go`, `internal/cli/orphan_unix_test.go` and
`internal/daemon/task_engine_test.go`. The TechSpec names the functions changed
in each and shows that none of them overlaps a 0259 hunk. Specs 0260 to 0264
are authored or delivered at the same time; none of them changes a file this
Spec declares, and no Verification here reads their artifacts.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or record format
  changes; the Spec edits tests and adds two unexported dependencies with
  production defaults. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; every test and Verification uses temporary homes, temporary
  repositories and fake runners. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0271 (this Spec) decides the rule:
  "A test bounds by real time only the budget it asserts", and every fix and the
  repository rule follow it. ADR-0259: "a ceiling only falls", for the Real-Time
  Budget ceilings as for Sequential Tests; every test touched stays parallel or
  keeps its named reason, and each package with a converted test runs once under
  the race detector. ADR-0213: a fixture process "must end when the test binary
  that started it ends", so a child that now blocks instead of sleeping keeps
  watching its owner. ADR-0266: "A Proven Signal is the only signal Roundfix
  code or tests send"; no Task adds a signal. ADR-0125: "A spawned fixture is a
  compiled binary, never a written script"; no Task writes an executable.
  ADR-0126: "The guard therefore runs per package"; no package's suite guard
  changes. ADR-0089: "Code under test takes its environment explicitly", which
  is why the two new dependencies are fields beside the existing clock and
  timeout fields. ADR-0244: "A Spec declares the domain terms it introduces".
  ADR-0257: "Inside a Run the Daemon is the only full gate", so Task agents run
  focused tests. ADR-0258: "A Spec is authored against the QA rerun classes".
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AF grant of 2026-10-12 ("Sim, como no AE", the same "Concedo" for the Governed
  Paths each Spec declares); bounded files: `docs/agents/specific-repository.md`.
  No Makefile, CI workflow, lint, formatter or module file change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0265-tests-that-wait-on-a-child-do-not-race-the-machine/_authorization.md`.

## Goals

- The two detach tests assert the phase they name whatever the machine's load,
  with every assertion they make today.
- No test in the repository bounds a child process, a goroutine or Git by a
  real-time budget it does not assert, and every asserted budget races work
  that cannot finish on its own.
- A repository rule fails, naming file and line, on a new short real-time budget
  in a test file whose ceiling does not cover it.
- Production behavior does not change.

## Core Features

1. Detach tests that fire the phase they assert: a budget a detach test does not
   assert takes the test deadline or the production constant, and a child that
   races an asserted budget blocks until it is killed or its owner is gone.
2. The other `internal/cli` tests of the class wait on the event their
   assertion reads or on the test deadline.
3. The `internal/agent` tests of the class stop racing the real cancel grace and
   the real disposable-session cleanup timeout; the cleanup timeout becomes a
   runner dependency with the production default.
4. The `internal/daemon`, `internal/worktree`, `internal/store` and
   `internal/spec` tests of the class let their budget expire only after the
   event they assert, or stop bounding work they do not assert.
5. The Real-Time Budget rule in `internal/testfixture`, with a ceiling per test
   file that only falls, ADR-0271, the **Real-Time Budget** glossary term and a
   repository rule line.

## Non-Goals / Out of Scope

- The three tests 0259's item branch already fixed
  (`TestBudgetExceededRunRecordsItsReason` and the two admin-lock disposal
  tests); Spec 0262 owns the open Backlog Entry about the admin-lock test.
- The detach survivor teardown in `detach_test.go` (`terminate`), which 0259
  rewrites; its 2 s reap after `SIGKILL` does not wait for a child to start.
- Production budgets: the 10 s detach liveness budget, the 2 s detach reap wait,
  Force Stop's grace and window, and every other production value stay as they
  are.
- Raising `-parallel`, changing the race detector's place, the Makefile or the
  CI workflows.

## Success Metrics

1. Success Metric: each test fixed by this Spec passes 50 consecutive times with
   `-race -parallel 16` while synthetic CPU load runs.
2. Success Metric: three consecutive full suites
   (`go test -count=1 -parallel 16 ./...`) exit 0 under synthetic CPU load, with
   no `signal: killed` and no `exited -1` line, stopping at the first failure
   without a rerun.
3. Success Metric: the Real-Time Budget rule passes on the repository and fails,
   naming file and line, on each seeded violation: an unlisted file with a short
   context timeout, a file over its ceiling, a short budget passed to a
   same-package parameter, and a short budget held in a named constant.
4. Success Metric: no production budget value changes, which the QA gate proves
   from the diff of every non-test Go file.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not produce,
found and read through Exa during authoring:

- Silva, Bell et al., "The Effects of Computational Resources on Flaky Tests"
  (TSE 2024, preprint https://www.jonbell.net/preprint/tse24-raft.pdf): "We find
  that 46.5% of the flaky tests in our dataset are RAFTs, indicating that a
  substantial proportion of flaky-test failures happen depending on the
  resources available when running tests." This is why the acceptance runs the
  suites under synthetic load rather than on an idle machine.
- Pei, Habchi, Rwemalika, Sohn and Papadakis, "An empirical study of async wait
  flakiness in front-end testing" (https://ceur-ws.org/Vol-3245/paper7.pdf):
  "tests relying on an explicit time to wait in order to synchronize the tests
  tend to create more flakiness (38 instances) than synchronizing on the status
  of DOM elements (24 instances)", and "where time-based issues are typically
  addressed by increasing the time to wait, DOM-based issues are resolved by
  actually introducing a missing synchronization point." This is why the fixes
  wait on events instead of raising budgets.
- The Secondbrain, read through `wiki/index.md` and three `qmd query`
  searches: the mirror of the 2026-08-03 finding
  `projects/roundfix/mirror/docs/history/findings/2026-08-03-a-200ms-attach-budget-fails-under-ci-load.md`,
  the archived Specs 0171 and 0213 under
  `projects/roundfix/mirror/docs/history/specs/`, and
  `wiki/concepts/testes-paralelos-em-go.md`, which records ADR-0259's
  conversion of the suite to parallel tests. They show the class recurring
  after per-test fixes, which is why this Spec adds a rule. The testing-boss
  anti-pattern list (`projects/skills/mirror/skills/05-implementation-loop/testing-boss/references/antipatterns.md`)
  names static sleeps and fixed-timeout waits as the top empirical flakiness
  cause. No wiki page covered resource-affected flaky tests; the Exa sources
  above were captured as a pending research Inbox Entry in the Secondbrain's
  own namespace.

## Glossary

- adds: **Real-Time Budget**

## Decisions

- A test bounds by real time only the budget it asserts; see ADR-0271.
- An asserted budget races only work that blocks until it is cancelled or its
  test binary is gone; see ADR-0271.
- Where a test cannot reach a production budget, the code takes it as a
  dependency with the production default, as `detachTimeouts` and the agent's
  cancellation clock already do; see ADR-0271.
- A ceiling table in the rule's own file, not a comment at each site, holds the
  remaining budgets; see ADR-0271.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the inventory, the fixes, the rule and
the build order.
