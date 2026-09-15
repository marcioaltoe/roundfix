---
spec: 0139-a-suite-that-passes-where-it-runs
prd: _prd.md
created: 2026-09-15
---

# A suite that passes where it runs — Technical Spec

## Executive Summary

Three test repairs and one Daemon change make the repository Verification mean
the same thing on this machine as it does in CI.

- **ACPX fixtures** switch from per-test hard links to symlinks, plus links created
  once for the package-directory case.
- **Task-cycle waits** end at the running test's deadline instead of after a fixed
  two seconds.
- **The historical audit subtest** gives way to a fixture the test builds, and its
  skip guard now requires reachability.
- **The QA gate step** runs the configured repository Verification itself, before
  the Agent turn and outside the Agent sandbox. It records the result in the
  seeded QA Report and tells the Agent not to run it again.

The main trade-off is the bootstrap. This Spec's own Run is driven by the binary
that existed when it started, which does not run Verification in the QA step. Its
first gate therefore still runs `make verify` inside the Agent sandbox. The gate
that proves the Daemon change has to be driven by a binary rebuilt from the
assembled tree.

The second trade-off is scope. The qa-gate skill still tells the Agent to run
`make verify`. This Spec overrides that through the gate prompt rather than
editing the skill, because Spec 0138 owns the skill's static-gate text and aligns
it afterwards.

## Project Constraints

- Identifier strategy: applicable — Verification event classifications and reasons, QA Report keys and Mechanical Refusal Codes keep their spelling and meaning; the QA gate's repository Verification publishes ordinary Verification events and coins no token. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, request or transport is created or changed; the network reach of the version freshness check in tests is recorded as a separate finding and left unchanged. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the repair touches Daemon-owned Verification, the QA gate's mechanical facts and the test contract of governed files. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0014 applies: the Daemon runs task verification and settles task status, so the QA gate's repository Verification is Daemon-run and the Agent only reads its result.
  ADR-0038 applies: the Daemon allows one Verification repair, and the QA gate's repository Verification adds no repair turn of its own.
  ADR-0056 applies: Spec Runs separate Task Capacity and Verification Capacity, so the QA gate's repository Verification acquires Verification Capacity like any other attempt.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and a repository Verification failure the Daemon observed stays a failure, never an environment block.
  ADR-0089 applies: code under test takes its environment explicitly, so no test may depend on local Git objects, machine load or churn of the running test binary.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, and the repository verification gate is one of the facts it names; this Spec runs it there without making any verdict more permissive.
  ADR-0117 applies: a defect is checked by the stage that can produce it, and the Daemon outside the Agent sandbox is the stage that can produce the repository Verification result.
  ADR-0130 applies: the audit judges governed paths, and history keeps the set honest, so the three governed paths this Spec repairs need its own grant.
  ADR-0057 applies: the Daemon exclusively owns Implement Task status, and the QA gate's repository Verification settles no Task status of its own.
  ADR-0091 applies: the QA gate is a Task node of its own type, and the repository Verification runs inside that node's step rather than as a separate command.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row compares the repository's own gate before and after this Spec on the machine where it failed.
  ADR-0020 is not applicable to this change: a parsed prompt result outranks the acpx exit code, and this Spec changes no Agent result handling; the ACPX repair touches only test fixtures.
  ADR-0093 is not applicable to this change: Spec consistency is checked by citation, never by inference, and this Spec changes no Spec Consistency Check rule.
  ADR-0097 is not applicable to this change: a QA row carries forward only on declared, unmoved evidence, and this Spec changes no carry-forward rule.
  ADR-0127 is not applicable to this change: process residue is a readiness fact, not a Run record, and this Spec changes no residue inventory.
- Tooling authority: applicable — express maintainer authorization: "Aprovar como proposto", 2026-09-15, recorded in `docs/specs/0139-a-suite-that-passes-where-it-runs/_authorization.md`; bounded files: `internal/daemon/task_engine.go`, `internal/daemon/task_engine_test.go`, `internal/speccheck/mechanical_test.go`. No sanctioned regeneration applies. The ACPX fixture file carries no historical grant and needs none. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Responsibility | Change |
| --- | --- | --- |
| ACPX adapter test fixtures | Stand in for acpx and runtime adapters by re-executing the test binary | Symlink per test; one-time links for the package-directory adapter |
| Task-cycle test helpers | Wait for scheduler starts and TaskCycle results | Bound every wait by the test deadline |
| Historical audit subtest and skip guard | Prove how the audit judges a grant against history | Controlled fixture; reachability guard |
| QA gate step | Seeds the report, withholds or starts the QA Agent | Runs the configured repository Verification before the Agent turn; records it; tells the Agent |
| Verification machinery | Runs commands with capacity, captures logs, publishes events | Reused unchanged |
| Mechanical stage | Computes machine facts from declarations and Git | None |

## Implementation Design

### ACPX fixtures

- **Adapter fixture.** Create a symlink to the test binary instead of a hard link.
  Exec keeps the symlink path as the first argument, so the fixture's sidecar
  lookup keeps working.
- **Package-directory adapter.** Its resolver follows symlinks past the package
  directory, so it still needs real links. Create those links once, in the
  package's test main, before the tests run, and remove them after.
- **Existing fixture test.** The test that proves the fixture binary survives
  concurrent exec keeps passing. No test removes `t.Parallel`, adds a retry or
  skips.

### Task-cycle waits

- **One helper.** Add a helper that returns the time left before the running
  test's deadline, minus a small margin, or a generous fallback when the test has
  no deadline. Every Task-cycle wait for a scheduler start or a TaskCycle result
  uses it instead of a fixed two-second timer.
- **Hangs still fail.** A real hang reaches the deadline, and the test binary's
  timeout reports every goroutine.
- **Assertions unchanged.** Assertions that nothing extra started are checks, not
  waits, and stay as they are.

### Historical audit fixture and guard

- **Controlled fixture.** The subtest that reads three unreachable commits from
  the local object store is replaced by a test that builds its own repository. In
  that repository, one Task commit changes a governed path, and the authorization
  record is widened in a later commit.
  - When the audit reads the record at the authorizing ancestor, it returns one
    `QA-AUTH-PATHS` finding for that path.
  - When the record already bounded the path at the ancestor, the audit returns no
    finding.
- **Reachability guard.** The skip guard for historical commits changes from
  "the object exists" to "the commit is an ancestor of the audited history". Its
  signature stays the same, so the characterization test that proves a commit
  from another repository does not resolve keeps its meaning.

### QA gate repository Verification

This runs in the QA gate step, after the mechanical stage has seeded the report,
and only when the stage did not withhold the Agent:

1. If no repository Verification command is configured, skip to the prompt and
   state that the gate runs the repository Verification itself.
2. Otherwise, move the Run to the Verifying state and run the configured command
   once. It runs through the same attempt request the Task precondition uses,
   with shared Verification Capacity, attempt 1, and the QA Task as the Work Item.
   It carries no failure classification, has no Verification Feedback, and
   publishes ordinary Verification events.
3. Copy the captured log into the Spec's QA evidence directory for this Run.
   Append a "Repository Verification" section to the seeded report with the
   command, verdict, exit status and a relative link to that log.
4. Add to the gate prompt, after the seeded-report instruction, a statement that
   the Daemon already ran the repository Verification outside the Agent sandbox.
   The statement gives the command, verdict, exit status and evidence path. It
   tells the Agent to record that observation as the static gate result and not
   to run the repository Verification again.
5. Continue to the Agent turn unchanged. The Daemon settles nothing from the
   verdict.

```text
Repository Verification: already run by the Daemon outside the Agent sandbox.
- command: <command>
- verdict: <pass|fail> (exit <status>)
- evidence: <qa evidence path>
Record this as the static gate result. Do not run the repository Verification again.
```

### Data Models

There is no schema change. The seeded QA Report gains one Markdown section, and
its frontmatter keys, row statuses and typed blocked counts are unchanged. The QA
evidence directory gains one log file.

### API Contracts

No command, flag or exit code changes. Run Events are the existing Verification
events, carrying the QA Work Item.

## Coverage Map

- Goal 1 → ACPX fixtures, Task-cycle waits, historical audit fixture and guard.
- Goal 2 → QA gate repository Verification.
- Goal 3 → Task-cycle waits.
- Goal 4 → QA gate repository Verification (verdict unchanged), regression locks
  across all four components.
- User Story 1 → ACPX fixtures, Task-cycle waits, historical audit fixture and
  guard.
- User Story 2 → QA gate repository Verification.
- User Story 3 → Task-cycle waits.
- User Story 4 → historical audit fixture and guard.
- Core Features 1-3 → the three test components.
- Core Features 4-6 → QA gate repository Verification.

## Integration Points

- **Configured command.** The Implement Command already passes the configured
  repository Verification command to the Daemon for Task preconditions, and the
  QA gate step reads the same value.
- **QA Report commit.** The copied log lives under the Spec's QA evidence
  directory, so the QA Report commit carries it next to the report. Evidence-path
  resolution finds it inside the repository.
- **Spec 0138.** Its skill text still tells the Agent to run `make verify`. The
  prompt statement takes precedence for the gate. Spec 0138 aligns the skill
  wording when its gate reruns.

## Testing Approach

1. **ACPX.** The fixture test proving concurrent exec survival keeps passing. A
   stressed package run, `-count=5 -parallel 16`, passes on this machine, where it
   reliably failed before. A structural check proves no per-test hard link to the
   test binary remains outside the package's test main.
2. **Task-cycle.** A structural check proves no fixed two-second wait remains in
   the Task-cycle helpers. A stressed package run under concurrent load passes,
   and focused Task-cycle tests still pass.
3. **Audit.** The new fixture test passes. The guard characterization test still
   passes. A structural check proves the unreachable historical commits are no
   longer read.
4. **QA gate step.** The existing seeding and withholding tests extend with three
   tests:
   - With a configured command, the Verification runs once before the Agent
     session, the report section and evidence exist, and the prompt carries the
     statement.
   - When the mechanical stage withholds the Agent, the Verification does not run.
   - Without a configured command, nothing runs, and the prompt says the gate runs
     it.
5. **Outside evidence.** The terminal QA Task runs `make verify` at the base and at
   the assembled tree, and runs the new audit test in a fresh clone.

## Build Order

1. Stop ACPX fixtures churning links to the test binary (depends on: none).
2. Bind Task-cycle waits to the test deadline (depends on: none).
3. Replace the historical audit subtest with a controlled fixture and a
   reachability guard (depends on: none).
4. Run the repository Verification in the QA gate step (depends on: 2).
5. Terminal QA (depends on: 1, 3, 4).

Steps 1, 2 and 3 touch disjoint files. Steps 2 and 4 share the Daemon test file,
so step 4 follows step 2.

## Risks & Considerations

- **Bootstrap.** The Run's gate is driven by the binary that started the Run,
  which lacks step 4.
  - If that gate fails only because `make verify` ran inside the sandbox, the
    Supervisor rebuilds the binary from the assembled tree, carries the settled
    Tasks forward, and reruns the gate with the rebuilt binary.
  - Record that first failure as expected. Do not treat it as a defect of steps
    1 to 4.
- **Log size.** A full `make verify` log can reach hundreds of kilobytes. QA
  evidence has carried logs of that size before; the copy is one file per gate
  run.
- **macOS-specific reproduction.** The ACPX kill and the load timeouts reproduce
  here, not in CI. The structural checks make Verification fail on any platform
  until the work is done.
- **Fleet behavior.** Any repository with a configured repository Verification
  now has it run by the Daemon during QA. The same command ran before, inside the
  Agent. The cost moves; it does not grow.

## Decisions

- **Prompt override instead of skill edit.** The Daemon states the repository
  Verification in the gate prompt, which keeps the governed skill out of this
  Spec and avoids a conflict with Spec 0138's pending skill rewrite.
- **Mechanical-stage-free.** The QA gate step owns running the command. The
  mechanical stage stays a pure function of declarations and Git, so its governed
  file is not touched.
- **No new classification.** The Verification publishes ordinary events, so no
  vocabulary is added.

## Vocabulary Contract

No token is coined.
