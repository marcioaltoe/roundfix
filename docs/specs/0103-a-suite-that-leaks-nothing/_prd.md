---
spec: 0103-a-suite-that-leaks-nothing
status: active
created: 2026-08-12
surfaces: [backend, infra]
---

# A suite that leaks nothing

The test suite writes into the repository it is reading, leaves processes running
for days, and carries a flake family whose common condition is spawn density —
and no Roundfix command can see any of it. A test that applies a baseline against
the live tree deleted tracked files while another test was copying them, so the
authoritative gate can both report a false failure and destroy the operator's
working tree. Four processes from the detach tests survived up to three days and
burned two hours and forty minutes of CPU doing nothing; `runs list` reported no
Runs, because they never registered one. A tool whose central promise is detached
execution has no command that answers what it detached that is still running.

## Project Constraints

- Identifier strategy: applicable — Run, Run Database, and Force Stop are glossary
  terms this Spec reports residue against, and a command that lists processes with
  no live Run record needs vocabulary the glossary owns. The closing node checks
  whether the work introduced or changed a term. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work is process lifetime, filesystem isolation,
  test fixtures, and a local inventory command. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0014 gives the Daemon ownership of task
  verification and status settlement, which bounds where a residue inventory may
  report without taking over settlement. The decisions extending that ownership are
  accounted and none applies: ADR-0020 ranks a parsed prompt result above the
  runtime's exit code, ADR-0038 allows one Verification repair, ADR-0057 gives the
  Daemon exclusive ownership of Implement Task status, ADR-0056 separates Task
  Capacity from Verification Capacity, and ADR-0096 with ADR-0117 place the gate's
  mechanical stage and its checks — this Spec reports residue and settles nothing,
  so it changes none of them. No accepted ADR governed test isolation or process
  residue when this Spec was written, which is why the invariants it establishes
  are new: ADR-0125 makes a spawned fixture a compiled binary, ADR-0126 places the
  isolation guard per package rather than once around the suite, and ADR-0127
  makes residue a readiness fact rather than a Run record. The closing node rests
  on two more that it follows rather than changes: ADR-0091 makes the QA gate a
  Task node of its own type, and ADR-0104 makes a Spec accept on evidence it did
  not author, which is what the upstream diagnosis and the 2026-08-06 process
  measurement supply here. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work changes Go test code, test fixtures, and production Go
  in the process-ownership and CLI packages. It edits no linter, formatter,
  test-runner, or build configuration, and it changes no test in order to make a
  failure disappear. Source: `docs/agents/agent-instructions.md`.

## Goals

1. No test writes inside the repository root, and the suite proves it.
2. No test can leave a process running after the suite ends.
3. A spawned fixture that dies is a fast, named failure rather than a long wait.
4. An operator can ask Roundfix what it detached that is still running.
5. A gate's own scratch state never reaches a Spec's committed evidence.

## User Stories

1. As a maintainer running the gate, I want the suite to fail naming the test that
   wrote inside the repository, so that a lost race does not surface as a
   stranger's error and does not edit my working tree.
2. As a maintainer who ran the suite yesterday, I want no process from it still
   running today, so that my machine is not slowly filled by tests that prove
   survival.
3. As a maintainer whose test spawned a fixture that failed to start, I want the
   failure in seconds naming the fixture, so that I am not waiting out a long
   budget for a process that is already dead.
4. As an operator of detached Runs, I want a command that reports processes with
   no live Run record, so that `ps` and personal discipline are not the only
   tools.
5. As a maintainer reviewing a Spec's evidence, I want no scratch repository or
   built binary committed there, so that the evidence is what a reader can check.

## Core Features

1. **A proven no-write invariant.** The suite proves that no test wrote inside the
   repository root, and a violation names the offending test directly rather than
   surfacing as an unrelated failure elsewhere.
2. **Detached children are terminated, not only released.** A test that proves a
   process survives its caller records that process and ends it non-cooperatively
   at teardown. The cooperative release sentinel lives outside a directory the
   test framework deletes, so a late waiter is not stranded by construction.
3. **Bounded fixture waits.** A spawned fixture's wait is bounded and observes the
   child, so a dead fixture fails in milliseconds with its own name instead of
   consuming a long budget and reporting a timeout.
4. **Fixtures that cannot race their own creation.** Spawned adapter fixtures are
   compiled test binaries rather than scripts written and immediately executed,
   removing the write-then-execute window that a dense start turns into a silent
   probe.
5. **A residue inventory.** A command reports processes Roundfix started that have
   no live Run record, so detached execution has an inventory rather than a habit.
6. **Force Stop proves the tree.** Stopping a Run proves the exit of the processes
   it spawned, not only of the registered owner.
7. **A gate's scratch state stays out of the evidence.** Scratch repositories and
   built binaries a gate produces are not committed as Spec evidence, and a
   submodule-shaped directory reference never reaches a commit.

## User Experience

The suite's failure names the test that wrote into the repository and the path it
wrote. The residue command lists each process with its age, its originating Run
if any, and what to do about it; with nothing to report it says so rather than
printing an empty table.

## Non-Goals / Out of Scope

- Changing what the detach tests prove. The survival property under test is
  correct and stays.
- Changing the concurrency or parallelism of the suite.
- Fixing any individual flake by loosening its assertion, extending a deadline, or
  retrying it.
- Killing processes the operator started outside Roundfix.
- Changing garbage collection or retention policy.

## Success Metrics

- A test that writes inside the repository root fails the suite naming itself,
  proven by exercising the check against such a test.
- After a full suite run, no process spawned by it remains.
- A spawned fixture that fails to start produces a failure in under one second.
- The residue command reports the class of processes that survived up to three
  days and consumed two hours and forty minutes of CPU on a developer machine on
  2026-08-06, measured outside this Spec by a process table nobody in the tool
  could see.
- No Spec evidence directory holds a built binary or a submodule-shaped reference.

## Decisions

- Fixtures become compiled test binaries rather than gaining a retry, because a
  retry hides a window that removing the window closes.

## Open Questions

- Whether the no-write invariant is enforced per package or once around the whole
  suite. Once around the suite is the default until answered, because a violation
  is cross-package by nature.
- Whether the residue command is its own command or a line in the existing
  readiness diagnostic. The default is the readiness diagnostic, since residue is
  a fact about the machine's current state rather than about a Run.
