---
spec: 0254-a-faster-make-test
status: active
created: 2026-10-08
surfaces: [backend]
---

# A faster make test

This is an engineering-framed minimal PRD. The Spec changes how the
repository's own tests run, not what Roundfix does, so it enters through the
TechSpec as a refactor (`docs/agents/spec-routing.md`).

A full `go test -count=1 -parallel 16 ./...`, the work of `make test` with the
result cache off, took 375 s to 384 s on the maintainer's 10-core machine on
2026-10-08, and `internal/cli` alone took 342 s to 383 s of it. The cause is
not slow code. Of the package's 1,409 top-level tests, 592 never call
`t.Parallel()`. Go runs those one after another before any parallel test of
the package starts, so their durations add up to 279 s to 324 s. Six more
packages repeat the pattern with 450 sequential tests: `internal/daemon`,
`internal/spec`, `internal/speccheck`, `internal/worktree`, `internal/store`
and `internal/baseline`. On GitHub's runners the suite measured 185 s to 223 s
against its 240 s budget. The Daemon pays the same chain at every Task
settlement through `make verify-changed`.

The maintainer asked for this Spec and, on 2026-10-08, placed it last in the
cycle: "A spec de make test é importante, mas acho que pode ser feita, se
complexa, após o final das demais specs e releases". The project memory's
target is a repository suite under 60 s. This Spec does not reach it, and
says so below.

## Prerequisites

None. Spec 0253 (`0253-authoring-rules-that-stop-qa-reruns`) is delivering
while this Spec is authored and adds test files under `internal/baseline`. No
Verification of this Spec reads an artifact 0253 creates. task_02 converts
every top-level test it finds in `internal/baseline` when it runs, including
any 0253 added, and its ceiling leaves room for them.

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or file format
  changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added; every test keeps its temporary directories and fake runners.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0259 (this Spec) decides the rule,
  its ceilings, the race and shuffle proof and the merged contract run.
  ADR-0089: "where a test exists precisely to verify that the process-level
  default is read correctly, it keeps the global and states in one line why it
  stays sequential", which is the `// Sequential:` form this Spec enforces.
  ADR-0125: "Fixtures are therefore compiled once", and no converted test
  starts executing a file it wrote. ADR-0126: "The guard therefore runs per
  package", and every converted package keeps its suite guard. ADR-0213: "A
  process a test starts, directly or through a detached Run, must end when the
  test binary that started it ends", which a conversion keeps. ADR-0090:
  "Caching facts across mutation boundaries was rejected"; pinned history is
  a fixed commit, so reading less of it caches nothing. ADR-0252: "Each
  Repository Contract Test declares when the selective gate runs it", and
  ADR-0253: "CI runs the Full Contract Run on every push to main"; both stay,
  and task_04 changes only how the selected contracts are invoked. ADR-0257:
  "the Agent runs focused tests of the packages it changed", which the Tasks
  follow. ADR-0244: "A Spec declares the domain terms it introduces"; task_01
  adds **Sequential Test** and **Parallel Test Package**. ADR-0170 rereads a
  Task Graph after each carried Run and ADR-0258 governs how a Spec is
  authored and how retirements write history; this Spec changes neither
  Delivery Retry nor history, and was authored against ADR-0258's rules.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No Makefile, CI workflow, lint, formatter or
  test-runner configuration changes. Four of the changed test files and the
  repository-owned guide are historically bounded Governed Paths, and the
  maintainer granted the governed paths this Spec declares on 2026-10-08.
  Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.
  Spec-contained authorization record:
  `docs/specs/0254-a-faster-make-test/_authorization.md`. Bounded files:
  `docs/agents/specific-repository.md`,
  `internal/baseline/derived_ownership_test.go`,
  `internal/cli/cli_test.go`, `internal/spec/coverage_test.go`,
  `internal/speccheck/mechanical_test.go`.

## Goals

- Every top-level test in the seven slow packages runs in parallel, or it says
  why it cannot, and a repository test keeps it that way.
- Converting a test never weakens an assertion and is proved free of shared
  state by the race detector and a shuffled order.
- The sequential tests left in `internal/cli` are the ones that truly change
  process-wide state.
- A test that materializes pinned history writes only the paths it reads.
- One `go test` runs every selected Repository Contract Test.
- Each claimed saving is measured against the starting commit on the same
  machine and load.

## User Stories

1. As the maintainer, I want `make test` with a cold cache to stop waiting on
   tests that never needed to run alone, so that a full run fits a short
   break.
2. As the Daemon, I want `make verify-changed` to finish sooner at each Task
   settlement, so that a Run spends its time on work.
3. As a Task agent, I want a test I add to be parallel by default or to say
   why not, so that the suite does not drift back to a serial chain.
4. As a reviewer, I want each conversion proved by the race detector and a
   shuffled order, so that faster never means flakier.
5. As CI, I want one `go test` for both contract tags, so that the Full
   Contract Run and the selective gate overlap their packages.

## Core Features

1. **Parallel Test Package rule.** A repository test requires every top-level
   test of a Parallel Test Package to call `t.Parallel()` first or to be a
   Sequential Test with a `// Sequential:` reason, under a per-package ceiling
   (ADR-0259).
2. **`internal/cli` runs in parallel.** Every eligible top-level test is
   converted; the race detector and `-shuffle=on` prove the package.
3. **Six more packages run in parallel.** `internal/daemon`, `internal/spec`,
   `internal/speccheck`, `internal/worktree`, `internal/store` and
   `internal/baseline` join the rule the same way.
4. **The sequential residue shrinks.** Tests in `internal/cli` that are
   sequential only because they change the working directory, the
   environment or the build identity take a per-test value instead (ADR-0089),
   and two pinned-history tests copy only the files they read.
5. **One contract invocation.** `cmd/verify-select` prints one invocation for
   the selected contracts of both tags (ADR-0259).
6. **Glossary and repository rule.** `CONTEXT.md` adds **Sequential Test**
   and **Parallel Test Package**, and `docs/agents/specific-repository.md`
   states the rule for whoever writes the next test.

## Non-Goals / Out of Scope

- The repository suite under 60 s. The converted suite is bound by the work
  it does: at least 83,660 Git processes per run and their files. Batching Git
  in production code is the next Spec.
- Changing the Makefile, `GO_TEST_PARALLEL`, `SUITE_BUDGET_SECONDS`, the CI
  workflows or the race detector's place in `make test`.
- Converting packages other than the seven; they join the rule later.
- Splitting `internal/cli` into packages, or adding a dependency.
- Changing any assertion, golden or fixture meaning, or any product behavior.

## Success Metrics

1. Success Metric: on the same machine, back to back, a full
   `go test -count=1 -parallel 16 ./...` at the audited head takes at most
   70 % of the wall time it takes at the starting commit `2ce5abe8`.
   Prototype: 227 s against 384 s, 59 %, before the residue work.
2. Success Metric: in the same pair of runs, the sequential phase of
   `internal/cli` (the summed durations of its top-level tests that never
   paused) is at most 25 % of the starting commit's. Prototype: 124 s against
   324 s, 38 %, with 43 sequential tests left; task_03 removes the slowest.
3. Success Metric: the repository test passes with every Parallel Test
   Package within its ceiling, and fails, naming the file and line, on a seeded
   top-level test that neither calls `t.Parallel()` first nor carries a
   `// Sequential:` reason.
4. Success Metric: each converted package passes once with `-race -short`
   and once with `-shuffle=on`.
5. Success Metric: `cmd/verify-select -contracts -all` prints one invocation,
   and the Full Contract Run's wall time at the audited head is below the
   starting commit's on the same machine.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The Go toolchain's own documentation, read offline with
  `go doc testing.T.Parallel` and `go help testflag` from the installed
  toolchain: "Parallel signals that this test is to be run in parallel with
  (and only with) other parallel tests", and `-parallel` "Allow parallel
  execution of test functions that call t.Parallel". Together they say a test
  that never calls `t.Parallel()` runs alone, which is the cause this Spec
  removes.
- GitHub Actions runs of `ci-verify.yml` on main, measured by GitHub's
  runners: run 37835405144 (`suite-time budget: 203s of 240s`,
  `internal/cli` 171.690 s), run 37834066256 (185 s, `internal/cli`
  163.133 s) and run 37821753645 (223 s, `internal/cli` 202.356 s). The QA
  gate reads them with `gh run view <id> --log` when the network is reachable,
  and otherwise records the row as blocked with its reason.
- Candido, Melo and d'Amorim, "Test Suite Parallelization in Open-Source
  Projects" (ASE 2017,
  https://damorim.github.io/publications/candido-etal-ase17.pdf, found and read
  through Exa during authoring): of the failures parallel runs caused, "approximately
  97.5% of the failures were caused by a null dereference and 1.6% were caused
  by concurrent access on unsynchronized data structures", and "Cases of
  likely broken test dependencies were not as prevalent as race conditions
  (0.8% of the total)". This is why each conversion runs the race detector.
  Mondal et al., "Soundy Automated Parallelization of Test Execution" (ICSME
  2021, https://damorim.github.io/publications/MondalETAL_ICSME21.pdf), adds
  that tests "can fail because of data races or broken test dependencies",
  which is why each conversion also runs once in a shuffled order.
- The Secondbrain mirrors of Specs 0071 (`projects/roundfix/mirror/docs/history/specs/0071-verification-cost.md`)
  and 0074 (`projects/roundfix/mirror/docs/history/specs/0074-git-spawn-economy.md`),
  read during authoring: earlier campaigns cut Git spawns and moved the
  regeneration gates out of `make test`, which is why this Spec targets the
  sequential chain instead.

## Unreachable Acceptance

- criterion: the CI suite time on GitHub's runners drops after merge, and
  `SUITE_BUDGET_SECONDS` can be lowered to keep it a regression guard
  reason: CI runs only after the Pull Request exists, on hardware the QA gate
  cannot reach
  satisfied-by: the operator reads the `Verify` step of the first `ci-verify.yml`
  run on main after merge, compares its `suite-time budget` line with the
  185 s to 223 s above, and proposes a lower budget in a later change

## Glossary

- adds: **Sequential Test**
- adds: **Parallel Test Package**

## Skills

- unchanged: command deliver plan — task_06 only routes process-global reads through commandDependencies, defaulting to today's read; the command's behavior, flags, output and exit codes are unchanged.
- unchanged: command deliver resume — task_06 only routes process-global reads through commandDependencies, defaulting to today's read; the command's behavior, flags, output and exit codes are unchanged.
- unchanged: command deliver retry — task_06 only routes process-global reads through commandDependencies, defaulting to today's read; the command's behavior, flags, output and exit codes are unchanged.
- unchanged: command deliver start — task_06 only routes process-global reads through commandDependencies, defaulting to today's read; the command's behavior, flags, output and exit codes are unchanged.
- unchanged: command deliver status — task_06 only routes process-global reads through commandDependencies, defaulting to today's read; the command's behavior, flags, output and exit codes are unchanged.
- unchanged: command deliver stop — task_06 only routes process-global reads through commandDependencies, defaulting to today's read; the command's behavior, flags, output and exit codes are unchanged.
- unchanged: command doctor — task_06 only routes process-global reads through commandDependencies, defaulting to today's read; the command's behavior, flags, output and exit codes are unchanged.
- unchanged: command reconcile — task_06 only routes process-global reads through commandDependencies, defaulting to today's read; the command's behavior, flags, output and exit codes are unchanged.

## Decisions

- Every top-level test in a Parallel Test Package runs in parallel or names
  why not, under a ceiling that only falls; see ADR-0259.
- Each conversion is proved by one race-detector run and one shuffled run,
  and the race detector stays out of `make test`; see ADR-0259.
- One `go test` runs the selected contracts of both tags; see ADR-0259.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the measurements, the guard, the
conversions, the contract invocation and the build order.
