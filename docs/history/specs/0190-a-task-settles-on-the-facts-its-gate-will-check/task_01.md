---
task: task_01
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: completed
type: backend
complexity: high
---

# Task 01: A Task of a gated graph runs the repository Verification when it settles

## Overview

Today the Daemon appends the configured repository Verification to a Task's
Verification only for a named repair Task that entered on a red repository
(`taskWithRequiredRepositoryVerification` in `internal/daemon/task_engine.go`).
Every other Task settles on its declared commands, and the QA gate is the first
to run the repository command. This Task makes every non-QA Task of a Task
Graph that has a QA gate Task end each Verification attempt with that command,
and adds the config key that turns it off. A Task Graph without a QA gate Task
keeps today's path exactly.

## Requirements

1. MUST set an unexported `settlementChecks` field on the non-QA plan in
   `TaskCycle` when `taskPlanWithoutQAGate` finds a QA gate Task, and leave it
   unset for a graph without one.
2. MUST pass the configured repository command as the required repository
   Verification in `executeTaskWorker` when `settlementChecks` and
   `TaskPlan.RepositoryVerificationAtSettlement` are both true, as well as in
   today's named-repair case. Both attempts use it. An empty configured command
   appends nothing.
3. MUST keep `verifyTaskPreWork` receiving the Task as authored, so the
   pre-work prober never runs the appended command. MUST keep
   `verifyRepositoryPrecondition` applying only to a Task that declares the
   command, and add no entry run.
4. MUST add API Contract 1:
   - `config.Verification.RepositoryAtSettlement`, `true` when the key is
     absent from both config files;
   - a non-boolean value refused at config load, naming
     `verification.repository_at_settlement`;
   - the key and one comment line in the generated config template, under
     `verification:`;
   - the value carried into `TaskPlan.RepositoryVerificationAtSettlement` where
     the Implement Command builds the Task plan.
5. MUST describe the key in `docs/user-guide/configuration.md`: a row in the
   settings table, the guide's copy of the template, and a sentence that it
   turns off only the repository Verification at settlement.
6. MUST keep the template's existing `defaults.verification` comment line
   byte-identical, because `internal/config/config_test.go` pins it.
7. MUST keep `TestTaskCycleWithoutRepositoryGateSkipsPrecondition` unedited and
   passing.
8. MUST NOT edit `internal/cli/cli_test.go`, rename or remove a top-level test,
   or write the new key to `.roundfixrc.yml`.
9. MUST put the new tests in the three files this Task creates. The engine
   tests use the existing Task-cycle fixture (`newTaskCycleFixture`,
   `taskFakeRunner`, `taskFakeVerifier`) with a graph that has a QA gate Task.
   The Implement Command tests run in a temporary repository with a fake
   runner and assert on the non-QA Task's Verification commands.

## Subtasks

- [ ] Mark the gated plan and require the repository command at settlement.
- [ ] Add the config key, its template line and its wiring into the Task plan.
- [ ] Describe the key in the configuration guide.
- [ ] Add one test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] In a gated graph the verifier sees the Task's declared command, then the
      repository command, on each attempt.
- [ ] A repository command that fails on the first attempt returns
      Verification Feedback naming that command, and the Task settles
      `completed` after the repair.
- [ ] A repository command that fails on the final attempt settles the Task
      `failed`, its reason names the command, and no Task commit exists.
- [ ] A Task that declares the repository command runs it once per attempt.
- [ ] With `RepositoryVerificationAtSettlement` false, a gated graph runs only
      the declared commands.
- [ ] The pre-work probe of a gated graph sees only the declared commands.
- [ ] A graph without a QA gate Task runs only the declared commands, proven by
      the existing unedited test.
- [ ] The key defaults to `true`, reads `false`, refuses a non-boolean, and is
      in the generated template.
- [ ] The Implement Command runs the repository command for a non-QA Task of a
      gated graph with the key absent, and does not with the key `false`.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/config/config.go`
- interface: `internal/cli/implement.go`
- interface: `docs/user-guide/configuration.md`
- creates: `internal/daemon/settlement_repository_verification_test.go`
- creates: `internal/config/settlement_config_test.go`
- creates: `internal/cli/implement_settlement_checks_test.go`
- instruction: `docs/adr/0182-a-task-settles-on-the-facts-its-gate-will-check.md`
- instruction: `docs/adr/0160-only-the-frozen-authorization-opens-a-red-repository-gate.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands|TestRepositoryVerificationFailureAtSettlementReturnsFeedback|TestRepositoryVerificationFailureOnFinalAttemptFailsTheTask|TestTaskDeclaringRepositoryVerificationRunsItOnce|TestRepositoryAtSettlementOffRunsOnlyDeclaredCommands|TestPreWorkProbeNeverRunsTheSettlementRepositoryVerification|TestTaskCycleWithoutRepositoryGateSkipsPrecondition)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands TestRepositoryVerificationFailureAtSettlementReturnsFeedback TestRepositoryVerificationFailureOnFinalAttemptFailsTheTask TestTaskDeclaringRepositoryVerificationRunsItOnce TestRepositoryAtSettlementOffRunsOnlyDeclaredCommands TestPreWorkProbeNeverRunsTheSettlementRepositoryVerification TestTaskCycleWithoutRepositoryGateSkipsPrecondition; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task six of the seven named tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestRepositoryAtSettlementDefaultsToTrue|TestRepositoryAtSettlementReadsFalse|TestRepositoryAtSettlementRefusesANonBoolean|TestConfigTemplateCarriesRepositoryAtSettlement)$" ./internal/config 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRepositoryAtSettlementDefaultsToTrue TestRepositoryAtSettlementReadsFalse TestRepositoryAtSettlementRefusesANonBoolean TestConfigTemplateCarriesRepositoryAtSettlement; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the four named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestImplementGatedGraphRunsRepositoryVerificationAtSettlement|TestImplementRepositoryAtSettlementOffRunsOnlyDeclaredCommands)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestImplementGatedGraphRunsRepositoryVerificationAtSettlement TestImplementRepositoryAtSettlementOffRunsOnlyDeclaredCommands; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither named test exists, so the command fails.
- `for pair in "docs/user-guide/configuration.md|verification.repository_at_settlement" "docs/user-guide/configuration.md|repository_at_settlement: true"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task the guide does not name the key, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Goal 3; Goal 4; User Story 1; User Story 3; Core Feature 1; Core Feature 2; Success Metric 1; Success Metric 4; Success Metric 5; Recorded limits
- [_techspec.md](_techspec.md) — When the stage applies; The repository Verification at settlement; API Contract 1; Testing Approach 1; Testing Approach 4; Build Order 1
- ADR-0182; ADR-0014; ADR-0038; ADR-0056; ADR-0111; ADR-0148; ADR-0159; ADR-0160

## Result

Implemented the Task 01 slice. `TaskCycle` marks the non-QA plan only when
`taskPlanWithoutQAGate` finds a QA gate Task. The worker requires the configured
repository command at settlement when both that marker and
`RepositoryVerificationAtSettlement` are true, while preserving the named
red-repository repair path. Both attempts reuse the existing append/deduplication
helper. Pre-work probing and repository entry preconditions still receive the
Task as authored; no repository entry run was added.

`verification.repository_at_settlement` defaults to `true`, accepts layered
boolean overrides, and rejects non-booleans with the full key name. Both
config templates include its default and one comment line. The Implement
Command carries the loaded value into its Task plan. The configuration guide
includes the setting, template line, and explanation of the narrow opt-out.

### Focused implementation evidence

The initial focused build exposed the missing `RepositoryAtSettlement` and
`RepositoryVerificationAtSettlement` fields before implementation. After the
implementation and test-fixture corrections, this focused command exited 0:

```sh
GOCACHE=/tmp/roundfix-task01-gocache rtk proxy go test -count=1 ./internal/config ./internal/daemon ./internal/cli -run 'Test(RepositoryAtSettlement|ConfigTemplateCarriesRepository|GatedGraphRuns|RepositoryVerificationFailure|TaskDeclaringRepository|PreWorkProbeNever|EmptySettlement|UngatedGraphIgnores|TaskCycleWithoutRepositoryGate|ImplementGatedGraph|ImplementRepositoryAtSettlement)'
```

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Declared command followed by repository command on each attempt | `TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands` passed; a first-attempt repository failure exposes both complete command sequences. |
| First repository failure produces Feedback and repair settles completed | `TestRepositoryVerificationFailureAtSettlementReturnsFeedback` passed; asserts the command in Feedback, the same Session, completed status, and one Task commit through the existing fake committer. |
| Final repository failure settles failed, names the command, and creates no Task commit | `TestRepositoryVerificationFailureOnFinalAttemptFailsTheTask` passed; asserts failed disk status, reason, two Agent turns, and zero committer calls. |
| Declared repository command runs once per attempt | `TestTaskDeclaringRepositoryVerificationRunsItOnce` passed; accounts separately for the existing entry precondition and asserts one repository command in each attempt. |
| Switch false leaves only declared settlement commands | `TestRepositoryAtSettlementOffRunsOnlyDeclaredCommands` passed. |
| Pre-work probe sees only authored commands | `TestPreWorkProbeNeverRunsTheSettlementRepositoryVerification` passed; records probe requests, confirms the appended settlement command, and checks that authored Verification was not mutated. |
| Graph without QA retains declared-only behavior | Existing `TestTaskCycleWithoutRepositoryGateSkipsPrecondition` passed unchanged. New `TestUngatedGraphIgnoresEnabledRepositorySettlementVerification` also passed with the switch explicitly enabled. |
| Config default, false, invalid type, and generated template | `TestRepositoryAtSettlementDefaultsToTrue`, `TestRepositoryAtSettlementReadsFalse`, `TestRepositoryAtSettlementRefusesANonBoolean`, and `TestConfigTemplateCarriesRepositoryAtSettlement` passed. Cases include User Config, Project Config precedence, quoted booleans, integers, collections, arbitrary text, and null. |
| Implement Command honors absent and false values | `TestImplementGatedGraphRunsRepositoryVerificationAtSettlement` and `TestImplementRepositoryAtSettlementOffRunsOnlyDeclaredCommands` passed in temporary Git repositories with the existing fake runner. Both assert non-QA settlement commands and confirm QA still runs its repository Verification. |

`TestEmptySettlementRepositoryCommandAppendsNothing` also passed for an empty
configured command. All new tests are in the three files this Task declares.

### Incremental check and scope

- `GOCACHE=/tmp/roundfix-task01-gocache rtk make verify-incremental`: initial
  sandbox run exited 2 because two existing force-stop integration tests could
  not read the host process table (`operation not permitted`). The rerun with
  host process access exited 0: formatting, vet, package tests, skill checks,
  and build passed. Cached package results were reused by this incremental
  check; the focused command above was uncached.
- `git -c core.fsmonitor=false diff --check`: exited 0.
- Python byte comparisons confirmed `internal/daemon/task_engine_test.go`,
  `internal/cli/cli_test.go`, `_tasks.md`, and `.roundfixrc.yml` are unchanged
  from HEAD. The template's existing `defaults.verification` comment remains
  byte-identical. Local guide inspection confirmed its new key, default line,
  and opt-out explanation.
- The only pre-existing tracked difference was this Task's Daemon-written
  `status: in_progress`; it was preserved. Authored Task content and checkboxes
  remain unchanged. Only this Result was added by the Agent.

Declared `## Verification` commands were not run. Task status and terminal
settlement remain Daemon-owned. No commit, push, or Pull Request was created.
No follow-up implementation was added from another Task's slice.

## Carry-forward provenance

- Source Run: `run_20260930T230750Z_197aa9aff15490a9`
- Source commit: `affd73fc161d64662ed02ab8da753260e46f7762`
