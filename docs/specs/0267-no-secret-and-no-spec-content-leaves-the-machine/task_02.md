---
task: task_02
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
status: completed
type: backend
complexity: medium
---

# Task 02: No Agent Session or Verification command receives a judge key

## Overview

The QA of Spec 0259 ran `roundfix archive --plan`, whose Archive Advice sent
the Spec's files to the judge provider, because every Agent Session inherits
the judge's keys. The hazard: the environment of the Daemon comes from the
operator's shell, which holds the judge's credentials, and its readers are the
Agent and every command the Agent or a Verification runs. This Task adds
`judge.WithoutKeyVariables` and calls it where a Run's command environments
are built: the acpx environment of every Agent Session and `ExecVerifier`.
A Light Tier session keeps its OpenRouter key under the implementation Stage
Key name. It is verifiable on its own through the named tests.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [A QA Run can send Spec content to the judge provider](references/2026-10-13-a-qa-run-can-send-spec-content-to-the-judge-provider.md):
   no process a Run starts for an Agent or a Verification can reach the judge.
2. MUST add `internal/judge/environment.go` with `WithoutKeyVariables` as
   `_techspec.md` → Interfaces states, deriving the names from
   `Questions.KeyVariables` of the embedded questions (ADR-0239), never from a
   copied list, and `internal/judge/environment_test.go` with
   `TestWithoutKeyVariablesRemovesEveryJudgeVariable`, which requires the
   result to lack every name `Load().KeyVariables()` returns, to keep every
   other entry in order, and to drop a duplicated or empty-valued entry too.
3. MUST make `acpxCommandEnv` drop every Judge Key Variable from its base and
   make `lightSessionEnvironment` follow `_techspec.md` → API Contract 5 and
   Invariant 5; and MUST make `ExecVerifier.Verify` follow API Contract 6.
4. MUST change, in `internal/agent/light_session_test.go`, only the
   expectations of `TestOpenModelImplementationFallsBackToTheSharedKey` and
   `TestLightSessionReadsTheImplementKeyByName` that pin the shared variable
   in the light session's option or environment: they move to
   `{env:ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY}`, the implementation name
   holding the shared value, and the shared name absent.
   `TestOpenModelImplementationPrefersTheImplementStageKey` keeps its bytes.
5. MUST add `internal/agent/judge_keys_test.go` with
   `TestNoAgentSessionReceivesAJudgeKeyVariable`: through the existing fake
   acpx, a standard session and a light session scoped to the shared key run
   every acpx command a Task turn uses (prepare, prompt, cancel, end, list),
   and each fake invocation records the names, never the values, of the
   variables it received. The test requires no Judge Key Variable in any
   record, a control variable set by the test present in every record, and
   the light session's records to hold the implementation Stage Key name.
6. MUST add `internal/daemon/verification_environment_test.go` with
   `TestVerificationReceivesNoJudgeKeyVariable`, which sets each Judge Key
   Variable and a control variable to runtime-built values and runs
   `ExecVerifier.Verify` with a command that exits 0 only when every Judge Key
   Variable is unset and the control is set. It changes the process
   environment, so it carries a `// Sequential:` comment naming it
   (ADR-0259).
7. MUST NOT change the judge's own key resolution, the Light Spend Log, the
   generic `OPENROUTER_API_KEY` removal, Git hooks or the QA format command.

## Subtasks

- [ ] Add the judge helper and its test.
- [ ] Drop the variables in the acpx environment and rename the light key.
- [ ] Run Verification commands without the variables.
- [ ] Adjust the two light session expectations and add the class tests.

## Acceptance Criteria

- [ ] Every acpx command of a standard and a light session receives no Judge
      Key Variable and keeps every other variable.
- [ ] A light session scoped to the shared key reads it under the
      implementation Stage Key name.
- [ ] A Verification command sees none of the Judge Key Variables.

## Consumers

- `judge.WithoutKeyVariables` and the Agent Session environment — consumers: `./internal/agent`, `./internal/daemon`

## Context

- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- instruction: `docs/adr/0239-each-openrouter-stage-reads-its-own-roundfix-key-first.md`
- instruction: `docs/adr/0238-a-light-tier-runs-low-complexity-tasks-on-an-open-model.md`
- creates: `internal/judge/environment.go`
- creates: `internal/judge/environment_test.go`
- interface: `internal/agent/codex_spawn.go`
- interface: `internal/agent/light_session.go`
- interface: `internal/agent/light_session_test.go`
- creates: `internal/agent/judge_keys_test.go`
- interface: `internal/daemon/daemon.go`
- creates: `internal/daemon/verification_environment_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestWithoutKeyVariablesRemovesEveryJudgeVariable|TestNoAgentSessionReceivesAJudgeKeyVariable|TestOpenModelImplementationFallsBackToTheSharedKey|TestOpenModelImplementationPrefersTheImplementStageKey|TestLightSessionReadsTheImplementKeyByName|TestVerificationReceivesNoJudgeKeyVariable)$' ./internal/judge ./internal/agent ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestWithoutKeyVariablesRemovesEveryJudgeVariable TestNoAgentSessionReceivesAJudgeKeyVariable TestOpenModelImplementationFallsBackToTheSharedKey TestOpenModelImplementationPrefersTheImplementStageKey TestLightSessionReadsTheImplementKeyByName TestVerificationReceivesNoJudgeKeyVariable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three new tests do not exist.

## References

- [_prd.md](_prd.md) — Goals; User Story 4; Core Feature 5; Success Metric 1
- [_techspec.md](_techspec.md) — API Contract 5; API Contract 6; Invariant 5; Testing Approach; Build Order 2
- [2026-10-13-a-qa-run-can-send-spec-content-to-the-judge-provider.md](references/2026-10-13-a-qa-run-can-send-spec-content-to-the-judge-provider.md)
- ADR-0273; ADR-0201; ADR-0239; ADR-0238; ADR-0259


## Result

Implemented the Task 02 environment boundary for Daemon Verification.
`judge.WithoutKeyVariables` derives its names from the embedded
`Load().KeyVariables()`, returns an independent copy, and drops every matching
entry, including duplicates and empty values. Invalid embedded questions stop
execution rather than pass credentials through.

`acpxCommandEnv` filters its base before applying runtime overrides.
`lightSessionEnvironment` references the implementation Stage Key name and
copies a selected shared value into that name; it still removes the generic
OpenRouter variable. `ExecVerifier.Verify` explicitly uses the filtered process
environment. Judge key selection and the Light Spend Log are unchanged.

The existing fake acpx in `internal/agent/acpx_runner_test.go` has one recording
hook for the new class test. The recorder lives in `judge_keys_test.go` and
writes command names and environment variable names, never variable values.
Only the two authorized light-session tests' expectations changed;
`TestOpenModelImplementationPrefersTheImplementStageKey` retains its bytes.

### Acceptance evidence

- Every acpx command: `TestNoAgentSessionReceivesAJudgeKeyVariable` passed for
  standard and shared-key light sessions. It checks every recorded invocation
  against the invocation log, requires prepare (`sessions ensure`), prompt,
  cancel, end (`sessions close`) and list, rejects every judge variable, and
  requires the control and unrelated variable names in each record. All light
  records contain the implementation Stage Key name.
- Shared-key Light Tier: `TestOpenModelImplementationFallsBackToTheSharedKey`
  passed, proving the shared value reaches the implementation variable and the
  shared name is absent. `TestLightSessionReadsTheImplementKeyByName` passed
  through fake acpx, including inherited configuration and invalid-config
  cases. The unchanged implementation-stage preference test passed too.
- Verification environment: `TestVerificationReceivesNoJudgeKeyVariable`
  passed. Its real shell command exits zero only when every embedded judge
  variable is unset and the runtime-built control value is retained; it also
  checks that the verifier-owning child process's fixture values were not
  changed. The parent test runs in parallel; its Sequential comment names
  the child fixture's isolated process environment mutation.
- Helper contract: `TestWithoutKeyVariablesRemovesEveryJudgeVariable` passed
  for all embedded variable names, duplicate and empty entries, exact-name
  matching, unrelated entry order, and input/result independence. Its preserved
  entries include `PATH`, `HOME` and `GOCACHE`.

### Focused checks

- Before the production change,
  `rtk proxy go test -count=1 -run 'TestNoAgentSessionReceivesAJudgeKeyVariable|TestVerificationReceivesNoJudgeKeyVariable' ./internal/agent ./internal/daemon`
  exited 1: all three judge variable names reached both Agent session types,
  and the Verification command rejected the inherited judge environment.
- After implementation,
  `rtk proxy go test -count=1 -v -run 'TestWithoutKeyVariables|TestNoAgentSession|TestOpenModelImplementation|TestLightSessionReads|TestVerificationReceives' ./internal/judge ./internal/agent ./internal/daemon`
  exited 0; all six selected top-level tests passed.
- `rtk proxy go test -count=1 -run 'TestACPXCommandEnv|TestACPXRunnerCommandEnv|TestLightSessionSendsNoWarmup' ./internal/agent`
  exited 0, covering existing override, guard, inherited environment, and light
  setup behavior.
- `rtk proxy go test -count=1 -run '^TestExecVerifier' ./internal/daemon`
  exited 0, covering existing successful cleanup, failed diagnostics, temporary
  failure classification, cancellation and infrastructure error behavior.
- `rtk proxy git diff --check` exited 0.

The worktree's symlink prevents sandbox initialization; authorized reads,
edits and focused checks ran with sandbox escalation. The Daemon-owned
`status: in_progress` was already present at entry and is preserved. Declared
Task Verification, repository Verification and incremental Verification were
not run by this Agent. No commit, push, pull request, manifest edit or other Task edit was
made. No follow-up work was identified within this slice.


### Verification Feedback — attempt 1

The Daemon's configured Verification reported a parallel-test contract failure:
this Task's environment-mutating test raised `internal/daemon` to 16 Sequential
Tests, above its ceiling of 15. Diagnostic artifact inspected:
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261014T102741Z_dd0d6c13d87ea4f9/verification/batch-002-attempt-1.log`.

The repair isolates this Task's environment mutation in a child process. The
named parent test now calls `t.Parallel()` first. A test-only dispatch in
`internal/daemon/main_test.go` runs the child fixture before the child starts
its test harness. That fixture sets every embedded Judge Key Variable and the
control to runtime-built values, invokes the real `ExecVerifier.Verify` with
the same shell assertions, checks its own environment remains intact, and
removes its temporary directory. The Sequential comment names the isolated
child's state. No repository ceiling, production behavior, or unrelated test
changed in this repair.

Focused repair evidence:

- `rtk proxy go test -count=1 -v -run '^TestVerificationReceivesNoJudgeKeyVariable$' ./internal/daemon`
  exited 0 with the repaired named test passing.
- `rtk proxy go test -count=1 -v -run '^TestEveryTestInAParallelTestPackageRunsInParallel$' ./internal/testfixture`
  exited 0, including its repository check; `internal/daemon` has 15 Sequential
  Tests against the unchanged ceiling of 15.
- A temporary Go overlay omitted only `ExecVerifier`'s environment-filter
  assignment and its now-unused import, then ran the named environment test.
  The test failed at the shell environment assertion as expected. This proves
  process isolation retains the regression oracle. The overlay and probe
  files were removed without modifying repository source.
- `rtk proxy git diff --check` exited 0.

Result evidence updated for handback. Task status is preserved; the Agent did
not rerun declared Task Verification or either configured repository
Verification command. No commit, push, pull request, manifest edit or other
Task edit was made. The Daemon owns the full configured Verification rerun.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/agent/acpx_runner_test.go`
- `internal/daemon/main_test.go`
