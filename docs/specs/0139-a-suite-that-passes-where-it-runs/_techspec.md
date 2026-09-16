---
spec: 0139-a-suite-that-passes-where-it-runs
prd: _prd.md
created: 2026-09-15
---

# A suite that passes where it runs — Technical Spec

## Executive Summary

Three test repairs and one Daemon change make the repository Verification mean
the same thing on this machine as in CI:

- **ACPX fixtures** use symlinks instead of per-test hard links, plus one-time
  links for the package-directory case.
- **Task-cycle waits** end at the running test's deadline.
- **The historical audit subtest** gives way to a fixture the test builds, and
  its skip guard requires reachability.
- **The QA gate step** runs the configured repository Verification before the
  Agent turn and outside the sandbox, and retains its log:
  - a pass is recorded in the seeded report and stated in the prompt;
  - a failure or an unobserved outcome is recorded as a Precondition Refusal
    that withholds the Agent, as ADR-0096 requires for a blocking machine fact.

The design accepts three trade-offs:

- **Bootstrap.** This Spec's own Run is driven by the binary that started it,
  and that binary does not have the new step. Its first gate still runs
  `make verify` inside the Agent sandbox. The gate that proves the new step has
  to run on a binary rebuilt from the assembled tree.
- **Fewer observations.** A red repository Verification stops the gate before
  any flow row is observed. The gate learns less in that round, but its verdict
  no longer depends on the Agent.
- **Stale skill text.** The qa-gate skill still tells the Agent to run
  `make verify`. The prompt overrides that, because Spec 0138 owns the skill's
  text and aligns it afterwards.

## Project Constraints

- Identifier strategy: applicable — Verification event classifications and reasons, QA Report keys and Mechanical Refusal Codes keep their spelling and meaning; the QA gate's repository Verification publishes ordinary Verification events and a failed one reuses the existing Precondition Refusal record, so no token is coined. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, request or transport is created or changed; the network reach of the version freshness check in tests is recorded as a separate finding and left unchanged. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the repair touches Daemon-owned Verification, the QA gate's mechanical facts, Spec artifact placement and the test contract of a governed file. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0014 applies: the Daemon runs task verification and settles task status, so the QA gate's repository Verification is Daemon-run and the Agent only reads its result.
  ADR-0038 applies: the Daemon allows one Verification repair, and the QA gate's repository Verification adds no repair turn of its own.
  ADR-0056 applies: Spec Runs separate Task Capacity and Verification Capacity, so the QA gate's repository Verification acquires Verification Capacity like any other attempt.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and a repository Verification failure the Daemon observed stays a failure, never an environment block.
  ADR-0089 applies: code under test takes its environment explicitly, so no test may depend on local Git objects, machine load or churn of the running test binary.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, the repository verification gate is one of the facts it names, and a failed or unobserved result is a blocking fact that withholds the Agent Session.
  ADR-0117 applies: a defect is checked by the stage that can produce it, and the Daemon outside the Agent sandbox is the stage that can produce the repository Verification result.
  ADR-0130 applies: the audit judges governed paths, and history keeps the set honest, so only the one governed test file this Spec repairs needs its grant.
  ADR-0035 is not applicable to this change: Spec Root is configurable and external Spec artifacts stay uncommitted, and this Spec writes no Spec artifact beyond the QA Report the Daemon already writes.
  ADR-0036 is not applicable to this change: review artifacts are committed in a separate docs commit, and this Spec changes no review Run or its artifact commit.
  ADR-0029 is not applicable to this change: review artifacts live with the Spec, and this Spec changes no review artifact location.
  ADR-0142 is not applicable to this change: head-bound Review Source Evidence decides the watch outcome, and this Spec changes no watch Run or Review Source evidence.
  ADR-0057 applies: the Daemon exclusively owns Implement Task status, and the QA gate's repository Verification settles the QA Task only through the existing refusal report.
  ADR-0091 applies: the QA gate is a Task node of its own type, and the repository Verification runs inside that node's step rather than as a separate command.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row replays the repository's own gate against Spec 0138's recorded failure.
  ADR-0111 applies: an unobserved Verification is unknown, not a verdict, so an outcome the runner could not observe becomes a refusal whose Run Event carries the unknown classification rather than a pass or an ordinary failure.
  ADR-0135 applies: an absent diagnostic is a reported state, not an empty message, so a refusal whose outcome left no log says so instead of carrying an empty link.
  ADR-0020 is not applicable to this change: a parsed prompt result outranks the acpx exit code, and this Spec changes no Agent result handling; the ACPX repair touches only test fixtures.
  ADR-0093 is not applicable to this change: Spec consistency is checked by citation, never by inference, and this Spec changes no Spec Consistency Check rule.
  ADR-0097 is not applicable to this change: a QA row carries forward only on declared, unmoved evidence, and this Spec changes no carry-forward rule.
  ADR-0127 is not applicable to this change: process residue is a readiness fact, not a Run record, and this Spec changes no residue inventory.
- Tooling authority: applicable — express maintainer authorization: "Aprovar como proposto", 2026-09-15, recorded in `docs/specs/0139-a-suite-that-passes-where-it-runs/_authorization.md`; bounded files: `internal/speccheck/mechanical_test.go`. No sanctioned regeneration applies. The Daemon and ACPX files this Spec changes are ordinary source outside the governed set. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Responsibility | Change |
| --- | --- | --- |
| ACPX adapter test fixtures | Stand in for acpx and runtime adapters by re-executing the test binary | Symlink per test; one-time links for the package-directory adapter |
| Task-cycle test helpers | Wait for scheduler starts and TaskCycle results | Bound every wait by the test deadline |
| Historical audit subtest and skip guard | Prove how the audit judges a grant against history | Controlled fixture; reachability guard |
| QA gate step | Seeds the report, withholds or starts the QA Agent | Runs the configured repository Verification before writing the report; refuses or records |
| Precondition Refusal report | Records a gate that stopped before its matrix | Reused unchanged, with the Verification command as the named check |

## Implementation Design

### ACPX fixtures

- **Per-test adapter fixture.** Create a symlink to the test binary instead of a
  hard link. Exec keeps the symlink path as the first argument, so the fixture's
  sidecar lookup keeps working.
- **Package-directory adapter.** Its resolver follows symlinks past the package
  directory, so it keeps real links. They are created once, by a helper called
  from the package's test main before the tests run, and removed after.
- **Existing guarantees.** The test proving the fixture binary survives
  concurrent exec keeps passing. No test removes `t.Parallel`, adds a retry or
  skips.

### Task-cycle waits

- **One helper.** It returns the time remaining before the running test's
  deadline minus a small margin, or a generous fallback when the test has no
  deadline. Every wait for a scheduler start or a TaskCycle result uses it
  instead of a fixed two-second timer.
- **Every waiting helper.** Scheduler starts, TaskCycle results, Verification
  starts, integrated Tasks and published events all route through it.
- **Hangs still fail.** A wait that reaches its bound fails the test. The failure
  names what the test waited for and includes every goroutine's stack, because the
  helper fails before the test binary's own timeout would report them.

### Historical audit fixture and guard

- **Controlled fixture.** The test builds its own repository. One Task commit
  changes a governed path, and a later commit widens the authorization record to
  bound that path.
  - When the audit reads the record at the authorizing ancestor, it returns
    exactly one `QA-AUTH-PATHS` finding for that path.
  - When the record already bounded the path at the ancestor, it returns no
    finding.
- **Reachability guard.** The skip guard changes from "the object exists" to
  "the commit is an ancestor of the audited `HEAD`". Its signature stays, so the
  project-root characterization test keeps its meaning.

### QA gate repository Verification

This step runs after the mechanical stage computes its result and before the QA
Report is written. It applies only when that result does not already withhold the
Agent and a repository Verification command is configured.

1. **Run the command.** Move the Run to the Verifying state and run the configured
   command once, with shared Verification Capacity, attempt 1 and the QA Task as
   the Work Item. It takes no Verification Feedback and no temporary-failure
   retry, exit status 75 included; that exit becomes an ordinary refusal here, and
   its event carries whatever metadata the existing machinery assigns. Its log
   keeps today's retention: kept on failure, removed on success.
2. **Pass.** Write the seeded report exactly as today and add the prompt statement
   below after the seeded-report instruction, then continue to the Agent turn.
3. **Command failure.** Before the report is written, set the mechanical result's
   Precondition Refusal and mark it blocking. The check name is the command, and
   the reason is `exited <status>; diagnostics: <path>`. The existing refusal
   contract writes `verdict: fail` and `rows_blocked_precondition: 1`, the Agent
   is withheld, and the QA Task settles from that report like any withheld gate.
4. **Unobserved outcome.** When the runner reports that the command could not
   start or its diagnostics could not be prepared, record the same refusal with
   the reason `outcome unobserved: <cause>; diagnostics: <path>`. When no
   diagnostics were retained, the reason says so rather than leaving the field
   empty.
5. **Stop.** When the Run's context is cancelled, return through the existing stop
   path without recording a refusal.
6. **No command configured.** Run nothing, and state in the prompt that the gate
   runs the repository Verification itself.

```text
Repository Verification: already run by the Daemon outside the Agent sandbox.
- command: <command>
- verdict: pass (exit 0)
- diagnostics: <path or "removed on success">
Record this as the static gate result. Do not run the repository Verification again.
```

**Events.** The step publishes the ordinary Verification events for its attempt
through the existing publisher. It adds no classification and changes no payload
shape. Today that publisher emits an unobserved outcome without its
classification, so the Run Event Stream projects it as an ordinary failure; a
finding records that gap, and this Spec does not close it. What the gate's record
does carry is the refusal itself, which names the command, the outcome and the
diagnostics state.

**Runner errors.** `VerificationCommandError` and `VerificationUnknownError` are
the two outcomes this step turns into a refusal. Any other runner error, such as
a diagnostics file that cannot be closed or renamed, stays an infrastructure
failure on the existing path and ends the Run, exactly as it does for a Task
Verification today.

**What this step does not do.** It copies no log into Spec evidence, adds no
report section and changes no evidence-path detection. A pass is recorded by the
Agent from the prompt statement; a failure is recorded by the refusal the Daemon
writes. Machine-validated evidence for this command is left to a later Spec.

### Data Models

- **No schema change.** The QA Report keeps its frontmatter keys, row statuses and
  typed counts. A refused gate uses the existing Precondition Refusal contract.
- **No new vocabulary.** No frontmatter key, row status, refusal code or
  Verification classification is added.

### API Contracts

No command, flag or exit code changes. Run Events are the existing Verification
events, carrying the QA Work Item.

## Coverage Map

- Goal 1 → ACPX fixtures, Task-cycle waits, historical audit fixture and guard.
- Goal 2 → QA gate repository Verification, Verification command runner.
- Goal 3 → Task-cycle waits.
- Goal 4 → Precondition Refusal report reuse, runner retention default,
  regression locks across all components.
- User Story 1 → ACPX fixtures, Task-cycle waits, historical audit fixture and
  guard.
- User Story 2 → QA gate repository Verification.
- User Story 3 → Task-cycle waits.
- User Story 4 → historical audit fixture and guard.
- Core Features 1-3 → the three test components.
- Core Features 4-7 → QA gate repository Verification, Precondition Refusal
  report.

## Integration Points

- **Configured command.** The Implement Command already passes the configured
  repository Verification command to the Daemon for Task preconditions, and the
  QA gate step reads that same value.
- **Spec 0138.** Its skill text still tells the Agent to run `make verify`. For a
  gate that reaches the Agent, the prompt statement takes precedence. A refused
  gate never reaches the Agent.

## Testing Approach

1. **ACPX.** The concurrent-exec fixture test keeps passing. A stressed package
   run at `-count=5 -parallel 16` passes on this machine, where it reliably
   failed before. A structural check proves no per-test hard link to the test
   binary remains.
2. **Task-cycle.** A structural check proves no fixed two-second wait remains. A
   stressed package run under concurrent load passes, and the helper's own test
   passes.
3. **Audit.** The fixture test and the remaining audit tests pass. Structural
   checks prove the unreachable commits are gone and the guard uses ancestry.
4. **QA gate step.** Five tests cover the five outcomes:
   - a pass runs once before the Agent session and is stated in the prompt;
   - a command failure becomes a refusal that withholds the Agent with verdict
     `fail` and publishes an unclassified failure;
   - an unobserved outcome becomes a refusal naming its cause and diagnostics
     state, with the event the existing publisher already emits;
   - a mechanical withholding runs no Verification;
   - a missing command runs nothing and says so in the prompt.
5. **Outside evidence.** The terminal QA Task replays Spec 0138's recorded
   failure and the base revision, measures the assembled tree, and runs the audit
   fixture in a fresh clone.

## Build Order

1. Stop ACPX fixtures churning links to the test binary (depends on: none).
2. Bind Task-cycle waits to the test deadline (depends on: none).
3. Replace the historical audit subtest with a controlled fixture and a
   reachability guard (depends on: none).
4. Run the repository Verification in the QA gate step, with runner retention
   and refusal (depends on: 2).
5. Terminal QA (depends on: 1, 3, 4).

Steps 1, 2 and 3 touch disjoint files. Steps 2 and 4 share the Daemon test file,
so step 4 follows step 2.

## Risks & Considerations

- **Bootstrap.** The Run's gate is driven by the binary that started the Run, and
  that binary lacks step 4. If that gate fails only because `make verify` ran
  inside the sandbox, the Supervisor:
  1. rebuilds the binary from the assembled tree;
  2. carries the settled Tasks forward;
  3. reruns the gate with the rebuilt binary.

  Record that first failure as expected, not as a defect of steps 1 to 4.
- **Unobserved flow rows.** A red repository Verification now stops the gate
  before flow rows run. This is ADR-0096's blocking-fact rule. The refusal names
  the log, so the next round starts from the failing command.
- **Unvalidated diagnostics.** The refusal names a diagnostics path under the Run
  artifacts, which no checker resolves. A later Spec owns machine-validated
  evidence for this command.
- **Unclassified unobserved events.** The Run Event Stream cannot tell an
  unobserved QA Verification from an ordinary failure, because the publisher omits
  the classification the projection requires. The refusal still records the
  outcome; the stream gap is left to its own Spec.
- **macOS-only reproduction.** The ACPX kill and the load timeouts reproduce
  here, not in CI. The structural checks make Verification fail on any platform
  until the work exists.
- **Fleet behavior.** Any repository with a configured repository Verification
  now has it run by the Daemon during QA, which refuses on failure. The same
  command ran before, inside the Agent.

## Decisions

- **Refusal reuse.** A failed or unobserved repository Verification reuses the
  existing Precondition Refusal instead of adding a Mechanical Refusal Code. No
  vocabulary is added, and settlement stays mechanical.
- **Opt-in retention.** Retention on success is an opt-in request field, so Task
  Verification keeps removing successful logs.
- **No retry.** The gate step takes no temporary-failure retry, because a
  repository Verification before QA is a single observation.
- **No evidence contract here.** Three review rounds kept finding new gaps in
  copying the log, linking it and validating that link across Spec Roots. This
  Spec therefore records the Verification in the prompt and in the refusal, and
  leaves machine-validated evidence to a later Spec.
- **No event-shape promise.** The step reuses the existing publisher untouched.
  Making an unobserved Verification legible in the Run Event Stream is a
  pre-existing gap, recorded as a finding rather than folded in here.
- **Prompt over skill.** A passing result is stated in the gate prompt. That
  keeps the qa-gate skill out of this Spec and avoids a conflict with Spec 0138's
  pending rewrite.

## Vocabulary Contract

No token is coined.
