---
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
prd: _prd.md
created: 2026-09-30
---

# A Task settles on the facts its gate will check — Technical Spec

## Executive Summary

The Daemon's Task engine gains one stage, Settlement Checks, for every non-QA
Task of a Task Graph that has a QA gate Task. The repository Verification is
appended to the Task's Verification as its last command, through the path a
named repair Task already uses. Two in-process checks run inside the same
Verification attempt: the refusing Spec Consistency findings the Task
introduced, and the gate's authorization audit applied to the commit the
Daemon is about to create. All three fail as Verification failures, so the
existing single repair, Temporary Verification Failure and Repeated Failure
handling serve them unchanged. The trade-off this design accepts is cost for
earliness: each settlement runs the repository Verification once more (206 to
223 seconds on this repository), against a gate refusal that costs a rerun of
15 to 60 minutes. It also accepts that the stage sees one Task's tree, not the
integrated tree of a parallel Wave.

## Project Constraints

- Identifier strategy: not applicable — no new entity identifier; the config
  key, the check labels and the glossary term follow existing naming. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local processes, Git, the Spec
  tree and the Run Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014, ADR-0038, ADR-0056, ADR-0057,
  ADR-0093, ADR-0094, ADR-0096, ADR-0111, ADR-0117, ADR-0130, ADR-0135,
  ADR-0148, ADR-0159, ADR-0160, ADR-0166, ADR-0176, ADR-0178 and ADR-0179
  hold; this Spec adds ADR-0182; its gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104,
  ADR-0155, ADR-0156 and ADR-0167. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer's authorization of
  2026-09-30 for every skill, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| Gate scope | `TaskCycle` in `internal/daemon/task_engine.go` | Marks the non-QA plan as settling under Settlement Checks when the graph has a QA gate Task |
| Repository Verification at settlement | `executeTaskWorker` in `internal/daemon/task_engine.go` | Requires the configured repository command for every Task of a gated graph, unless the switch is off |
| In-process checks | New `internal/daemon/settlement_checks.go` and `internal/daemon/settlement_spec_consistency.go`; `runVerificationAttempt` in `internal/daemon/engine.go` | A Verification attempt runs its checks after its commands and before its verdict |
| Spec Consistency findings | `internal/speccheck/mechanical.go` | Exports the reason the gate records for one refusing finding |
| Prospective commit audit | `internal/speccheck/mechanical.go` | One per-commit audit function serves existing commits and a commit that does not exist yet |
| Switch | `internal/config/config.go`, `internal/cli/implement.go` | `verification.repository_at_settlement`, carried into the Task plan |

Nothing in the QA gate changes. `runQAGate` keeps its precondition, its
repository Verification and its mechanical stage.

## Implementation Design

### Interfaces

```go
// internal/daemon — a new Engine dependency, defaulted like MechanicalStage.
type SettlementChecker interface {
	// RefusingFindings returns the findings that would make the Spec's own
	// gate precondition refuse.
	RefusingFindings(specsRoot, workDir, slug string) ([]speccheck.Finding, error)
	// AuditCommit applies the gate's authorization audit to a commit that
	// does not exist yet.
	AuditCommit(ctx context.Context, request speccheck.MechanicalRequest,
		commit speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error)
}

// internal/daemon/engine.go — one in-process check of a Verification attempt.
type verificationCheck struct {
	Label string // "settlement check: spec consistency"
	Run   func(ctx context.Context, diagnosticPath string) (failure string, err error)
}

// internal/speccheck
type ProspectiveTaskCommit struct {
	TaskID   string
	TaskFile string
	Parent   string   // the revision the commit will be created on
	Changed  []string // the paths the commit will stage
}

func AuditProspectiveTaskCommit(ctx context.Context, request MechanicalRequest,
	commit ProspectiveTaskCommit) ([]MechanicalFinding, error)

// RefusalReason is the code and sentence the gate records for one finding.
func RefusalReason(finding Finding) string
```

`TaskPlan` gains `RepositoryVerificationAtSettlement bool` and an unexported
`settlementChecks bool`. `verificationAttemptRequest` gains
`Checks []verificationCheck`.

### Data Models

`config.Verification` gains `RepositoryAtSettlement bool`, `true` by default.
The Run Database schema does not change.

### When the stage applies

`TaskCycle` already separates the QA gate Task with `taskPlanWithoutQAGate`.
It sets `settlementChecks` on the non-QA plan when that Task exists. Every
decision below reads that field, so a Task Graph without a QA gate Task takes
today's path exactly.

### The repository Verification at settlement

`executeTaskWorker` passes the configured repository command as
`requiredRepositoryVerification` in two cases. The first is today's: a named
repair Task that entered on a red repository (ADR-0160). The second is new:
`settlementChecks` and `RepositoryVerificationAtSettlement` are both true.

- `taskWithRequiredRepositoryVerification` already appends the command as the
  last one, and already leaves a Task that declares the same command
  unchanged. Both attempts use it, so nothing else moves.
- `verifyTaskPreWork` keeps receiving the Task as authored. The pre-work
  prober never runs the appended command, which would pass before the work and
  read as vacuous.
- The entry precondition (`verifyRepositoryPrecondition`) keeps applying only
  to a Task that declares the command. The Daemon adds no entry run.

### Checks inside a Verification attempt

`runVerificationAttempt` runs `req.Checks` after the commands it reached,
whatever their result, and before it publishes the attempt's verdict. It skips
them only when the attempt ends without a verdict: a Stop Request, an
infrastructure error or an unobserved command.

- A check that returns a failure, or an error, is appended to
  `CommandFailures` as a `VerificationCommandError` whose `Command` is the
  check's label and whose `OutputPath` is its diagnostics file. An error never
  passes: its text becomes the diagnostics.
- The attempt's verdict is `failed` when a command or a check failed.
- Diagnostics go to
  `runs/<run-id>/verification/batch-NNN-attempt-M-settlement-<name>.log` in
  the Artifact Directory.
- The retry of a Temporary Verification Failure runs the checks again. They
  take under a second, so no state is carried between the two.

`verifyTask` fills `Checks` only when `settlementChecks` is true, from one
builder in `settlement_checks.go` and one in `settlement_spec_consistency.go`. `repairTaskVerification` needs no change:
it already turns every entry of `CommandFailures` into Verification Feedback.

### Spec Consistency findings the Task introduced

- **Baseline.** `executeTaskWorker` calls `RefusingFindings` after the
  pre-work probe and before the Agent starts. It keeps the set of
  `speccheck.RefusalReason` values. An error there settles the Task `failed`
  before Agent work, with the error in its reason.
- **Check.** The check calls `RefusingFindings` again and fails when a reason
  is absent from the baseline. Its diagnostics list each new finding's code,
  sentence, locations and fix line.
- **One definition.** The default `SettlementChecker` runs `speccheck.Check`,
  `PromoteGaps` and `GatePrecondition`, the calls `qaGatePrecondition` makes.
  `qaGatePrecondition` and the default checker share one helper.
  `PreconditionRefusal` builds its reasons with `RefusalReason`, so the gate
  and the settlement identify a finding the same way.

### The prospective Task commit audit

- **One audit.** `detectMechanicalAuthPaths` evaluates each Task commit
  through one per-commit function that takes the commit's parent revision and
  its changed paths. For an existing commit those are `<sha>^1` and the
  commit's diff. `AuditProspectiveTaskCommit` calls the same function with
  `ProspectiveTaskCommit.Parent` and `.Changed`. It returns the findings and
  adds no skip records.
- **Request.** `qaMechanicalRequest` resolves the authorization record and
  the delivery target. That resolution moves into one helper, which the
  settlement check also calls, so both read the same record at the same
  revisions.
- **Paths.** The check resolves the prospective commit with
  `prepareTaskCommit`, the function the settlement itself uses, and takes its
  `stageable` list. The parent is `git rev-parse HEAD` in the Task's work
  directory.
- **When.** The audit runs only when a staged path other than the Task file
  is a Governed Path. Otherwise the check passes without reading a grant,
  which is what the gate does for a commit with no Governed Path.
- **Wording.** A finding for a prospective commit reads
  `Task <id>'s prospective commit changes <path> outside authorization grant <record>'s exact bounded files`.
  Every other field of the finding is unchanged.

The existing operation check in `executeTask`
(`spec.RequireGovernedOperation`) stays where it is, after the final attempt.

### API Contracts

1. API Contract: config key `verification.repository_at_settlement` — a
   boolean in User Config or Project Config, `true` when absent. `false` stops
   the Daemon from appending the repository Verification at settlement and
   changes nothing else. A non-boolean value is refused at config load, as the
   other boolean keys are. The generated config template carries the key and
   one comment line.
2. API Contract: Run Event Stream — each in-process Settlement Check
   publishes `verification` events with phase `started` and then
   `command-passed` or `failed`. They carry the Task, the attempt and a
   `command` that is the label `settlement check: spec consistency` or
   `settlement check: authorization`. Only these two Roundfix-defined labels
   are projected into the public stream; the appended repository Verification
   publishes the events any Verification command publishes, whose command
   stays redacted as today.
3. API Contract: Task failure reason — a Task failed by an in-process
   Settlement Check settles with
   `Settlement check failed: <label>: <first diagnostic line>; diagnostics: <path>`.
   A Task failed by the appended repository Verification settles with today's
   `Verification failed: command "<command>" exited with …` reason.

## Vocabulary Contract

This Spec coins one term, **Settlement Check**: a gate fact the Daemon checks
for a Task before it settles. It is distinct from the QA settlement table,
which maps a QA verdict to a Spec's closing. The two check labels and the
`Settlement check failed:` reason prefix are the only new emitted words.
task_04 documents them in `docs/user-guide/commands.md`, and its Verification
asserts each one there.

## Coverage Map

- Goal 1 → Checks inside a Verification attempt; The repository Verification
  at settlement.
- Goal 2 → Checks inside a Verification attempt; API Contract 3.
- Goal 3 → API Contract 1; API Contract 2.
- Goal 4 → When the stage applies.
- User Story 1 → The repository Verification at settlement.
- User Story 2 → The prospective Task commit audit.
- User Story 3 → API Contract 1.
- User Story 4 → Build Order 4.
- Core Feature 1 → When the stage applies.
- Core Feature 2 → The repository Verification at settlement; API Contract 1.
- Core Feature 3 → Spec Consistency findings the Task introduced.
- Core Feature 4 → The prospective Task commit audit.
- Core Feature 5 → Checks inside a Verification attempt; API Contract 3.
- Core Feature 6 → API Contract 2; Build Order 4.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 3.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 1; Testing Approach 4.
- Success Metric 5 → Testing Approach 1.

## Integration Points

- **The QA gate.** Unchanged. It reruns all three facts on the integrated
  tree.
- **Task Worktrees.** In a parallel Wave each Task is checked in its own Task
  Worktree, whose `HEAD` is the Run Branch tip the Task started from.
  Integration cherry-picks the commit, which keeps its merge base with the
  delivery target, so the gate's audit reads the same grant.
- **Verification Capacity.** The appended command is one more command of the
  Task's attempt and holds the same capacity slot.
- **The Settle Command.** Out of scope. It keeps settling on declared
  Verification.

## Testing Approach

1. **Repository Verification at settlement.** New
   `internal/daemon/settlement_repository_verification_test.go`, over the
   existing Task-cycle fixture (`newTaskCycleFixture`, `taskFakeRunner`,
   `taskFakeVerifier`) with a graph that has a QA gate Task:
   - the verifier sees the declared command, then the repository command;
   - a repository failure returns Verification Feedback, and the Task settles
     `completed` after a repair;
   - a repository failure on the final attempt settles the Task `failed` and
     creates no commit;
   - a Task that declares the command runs it once;
   - the switch set to `false` runs only the declared command;
   - the pre-work probe sees only the declared command.

   `TestTaskCycleWithoutRepositoryGateSkipsPrecondition` stays unedited and
   pins the gate-less graph.
2. **One audit.** New `internal/speccheck/prospective_commit_audit_test.go`,
   over `gittest` repositories:
   - a Governed Path outside the grant is refused;
   - a bounded path and a sanctioned regeneration output are accepted;
   - a change to the authorization record is refused as self-approval;
   - the prospective audit, and the gate's audit of the commit created from
     the same parent and paths, report the same findings.
3. **In-process checks through the engine.** New
   `internal/daemon/settlement_checks_test.go`, with a fake
   `SettlementChecker`:
   - a finding absent at Task start fails the check and returns Feedback;
   - a finding present at Task start does not;
   - an audit finding returns Feedback, and a repaired tree settles
     `completed`;
   - a failing check on the final attempt settles the Task `failed` with API
     Contract 3's reason and no commit;
   - a checker error fails the check;
   - a gate-less graph never calls the checker.

   The `SettlementChecker` seam is new. It is needed because the real checks
   read a Spec tree and Git history that the Task-cycle fixture fakes. Two
   tests run the default checker, one over a `gittest` repository and one
   over a temporary Spec tree.
4. **Switch.** New `internal/config/settlement_config_test.go` covers the
   default, `false`, a non-boolean refusal and the generated template. New
   `internal/cli/implement_settlement_checks_test.go` runs the Implement
   Command over a gated graph with the switch on and off.
5. **Guides.** Phrase checks on the user guide, the two skills, their mirrors
   and `CONTEXT.md`, plus `make skills-sync-check`.
6. **Real history.** The QA gate reads the recorded 0187 Run and the
   published presubmit documentation as evidence this Spec did not author.

## Build Order

1. The repository Verification at settlement, the gate scope and the switch,
   with the configuration guide, task_01 (depends on: none).
2. Checks inside a Verification attempt, the `SettlementChecker` seam with
   its default, and the prospective Task commit audit, task_02 (depends on: 1).
3. Spec Consistency findings the Task introduced: the baseline, the check and
   `RefusalReason`, task_03 (depends on: 2).
4. The Roundfix skill, the `write-tasks` rule, the commands guide and the
   glossary term, task_04 (depends on: 1, 2, 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

Tasks 1 to 3 all change `internal/daemon/task_engine.go`, so they run in
series.

## Risks & Considerations

- **The one declared break.** In a gated graph every Task attempt now ends
  with the repository Verification. A graph whose Tasks left the repository
  red between slices now fails at the first such Task. The `write-tasks` rule
  and the switch are the two answers.
- **A red repository on entry.** No entry run is added, so a Task can fail for
  a failure it did not cause. Its reason names the command and its
  diagnostics, and the named repair Task path of ADR-0160 is the way to fix a
  red repository inside a Run.
- **A parallel Wave.** Two Tasks can each pass and integrate into a red tree.
  The next settlement or the gate reports it.
- **The hidden repository check.** The appended command runs only after the
  declared commands pass. After a declared failure on the first attempt it
  first runs on the final attempt, where a failure has no repair left.
- **Parity.** The audit and the finding identity each have one definition,
  used by the gate and by the settlement. The parity test in Testing Approach
  2 fails if they diverge.
- **Event consumers.** The check labels are new `command` values in existing
  events. No field is added or removed.
- **New config key.** The binary must know the key before any config file
  carries it, so this Spec does not write it to this repository's Project
  Config.

## Decisions

- Settlement Checks apply only in a graph with a QA gate Task. See ADR-0182.
- The repository Verification is appended to the Task's Verification instead
  of getting its own loop.
- The two in-process checks run inside the Verification attempt, so one
  Feedback turn carries every failure that attempt found.
- Spec Consistency findings are compared with the Task's start.
- The audit of a prospective commit and of an existing commit is one function.
- One switch, for the repository Verification only.
