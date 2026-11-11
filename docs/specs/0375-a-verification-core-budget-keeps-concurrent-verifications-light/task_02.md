---
task: task_02
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
status: completed
type: backend
complexity: high
---

# Task 02: Every Verification and Agent Session of an Implement Run receives the Verification Core Budget

## Overview

The Implement Command computes the Verification Core Budget once from its
loaded User Config and the process environment and hands it to the Daemon.
Every Verification command of the Run — the pre-work probe and every attempt
request — and every Agent Session of the Run receive `GOMAXPROCS` and
`ROUNDFIX_VERIFICATION_CORES`. The configuration guide and the Roundfix
`implement` reference describe the budget and how a non-Go test runner opts
in.

## Requirements

1. MUST characterize first, in the Task's `## Result`: name the two call
   sites of `Verifier.Verify` in `internal/daemon` (`verificationAttempt` and
   `probeCommandsForRun`) and every caller of `probeCommandsForRun`, and the
   two places in `internal/agent/acpx_runner.go` that append the Run Marker to
   the session overrides.
2. MUST add `VerificationCores int` to `daemon.Dependencies`
   (`internal/daemon/engine.go`) and pass it as `VerifyRequest.Cores` from
   `verificationAttempt`; add a `cores int` parameter to `probeCommandsForRun`
   (`internal/daemon/verification_probe.go`) after `verifier`, passed
   `engine.deps.VerificationCores` by `verifyTaskPreWork`
   (`internal/daemon/task_engine.go`) and `0` by `ProbeCommands`; and set
   `req.Cores = engine.deps.VerificationCores` at the start of
   `runAgentSession` (`internal/daemon/agent_session_owner.go`), as API
   Contract 3 says. With `VerificationCores` 0 every request is what it is
   today.
3. MUST add `Cores int` to `agent.ExecuteRequest` (`internal/agent/agent.go`)
   and, in `internal/agent/acpx_runner.go`, append
   `coreBudgetOverrides(req.Cores)` right after each Run Marker append:
   nothing for `Cores < 1`, otherwise `GOMAXPROCS=<n>` and
   `ROUNDFIX_VERIFICATION_CORES=<n>`, so `acpxCommandEnv` replaces the
   inherited copies.
4. MUST create `internal/cli/verification_cores.go` with
   `verificationCoreBudget(ctx context.Context, config roundconfig.Config) corebudget.Budget`
   as API Contract 3 says; add `logicalCPUs func() int` to
   `commandDependencies` in `internal/cli/cli.go`, `runtime.NumCPU` in
   `defaultCommandDependencies`; add `cores int` to `implementCapacities` in
   `internal/cli/implement.go`, set from
   `verificationCoreBudget(ctx, loadedConfig.Config).Cores` where
   `runImplementCommand` builds it, and pass it as
   `Dependencies.VerificationCores` in `executeImplementCycle`. No other
   engine construction sets it.
5. MUST create `internal/daemon/run_core_budget_test.go` with
   `TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget`: a Task
   cycle over the existing fixtures with a verifier that records each
   request and delegates to the existing fake, the plan's repository
   Verification set and run at settlement; row `budget five` asserts every
   recorded Verification request, at least one pre-work probe request and the
   settlement command among them, and every Agent Session request recorded by
   the fake runner carry `Cores` 5; row `no budget` asserts 0.
6. MUST create `internal/agent/core_budget_test.go` with a recorder that the
   fake acpx of `internal/agent/acpx_runner_test.go` calls right after its
   judge-environment recorder (one added call there), writing per invocation
   the command key and the values of both variables, and
   `TestAnAgentSessionReceivesTheCoreBudget`: with an inherited
   `GOMAXPROCS=12`, a `Run` of a fake ACP runtime with `Cores` 3 shows
   `sessions ensure` and `prompt` with `3` and `3`; with `Cores` 0 they show
   `12` and no `ROUNDFIX_VERIFICATION_CORES`.
7. MUST create `internal/cli/verification_cores_test.go` with a helper that
   gives the command twelve logical CPUs and an environment without either
   variable unless the row adds `GOMAXPROCS`, and
   `TestImplementGivesTheEngineTheVerificationCoreBudget`: `roundfix
   implement` over the existing implement workspace with User Config
   `runs:\n  max_verifying: 2\n` gives every Verification request and every
   Agent Session request 3, and 7 when the row adds `GOMAXPROCS=7`. Every
   top-level test calls `t.Parallel()` first.
8. MUST keep every existing test of `./internal/daemon`, `./internal/agent`
   and every Implement test of `./internal/cli` passing unchanged.
9. MUST add `## Verification Core Budget` to `docs/user-guide/configuration.md`
   right after `## Machine Verification Capacity`, covering API Contract 1,
   who receives the budget and who does not (PRD Core Features 3 to 5 and
   Non-Goals), and containing these sentences word for word: "Roundfix gives
   every Verification and every Agent Session of an Implement Run, and the
   repository gate of a Delivery Queue, a Verification Core Budget: half of
   the logical CPUs, at least 1, divided by `runs.max_verifying`, or by
   `runs.max_active` while `runs.max_verifying` is `0`, and never below 1.";
   "It sets `GOMAXPROCS` and `ROUNDFIX_VERIFICATION_CORES` to that number.";
   "A `GOMAXPROCS` that holds a positive whole number in the environment of
   the process that starts the Run or the queue takes its place."; "A test
   runner that reads neither variable keeps its own parallelism; changing its
   configuration is a tooling change your repository authorizes." It shows
   one Vitest configuration example whose `maxWorkers` reads
   `ROUNDFIX_VERIFICATION_CORES`.
10. MUST add one paragraph to item 9 of the canonical
    `.agents/skills/roundfix/references/implement.md` that names the
    Verification Core Budget, `GOMAXPROCS`, `ROUNDFIX_VERIFICATION_CORES`,
    the share and divisor, the explicit `GOMAXPROCS` rule and the Agent
    Sessions; then run `make skills-sync` and raise the Roundfix skill version
    with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
    Never write inside a `### QA settlement` section, never hand-edit a
    mirror, version field or record.
11. MUST, as the Task's last step, after its last guide and skill edit,
    re-record the Behavior Surface Record with
    `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
    Of the Governed Paths the Task may change only
    `.agents/skills/roundfix/SKILL.md`,
    `.agents/skills/roundfix/references/implement.md` and
    `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Characterize the verifier and Agent Session call sites.
- [ ] Carry the budget through the Daemon and the acpx overrides.
- [ ] Compute it in the Implement Command.
- [ ] Add the Daemon, Agent and CLI tests.
- [ ] Write the configuration guide section and the `implement` reference, sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Every Verification request and Agent Session request of an Implement Run carries the budget its process computed, and none carries one when the engine has none.
- [ ] An Agent Session's acpx invocations see both variables set to the budget, replacing an inherited `GOMAXPROCS`.
- [ ] The configuration guide and both copies of the `implement` reference describe the budget.

## Context

- instruction: `docs/adr/0385-each-verification-receives-a-verification-core-budget-in-its-environment.md`
- instruction: `internal/daemon/task_engine_test.go`
- instruction: `internal/agent/judge_keys_test.go`
- instruction: `internal/cli/implement_settlement_checks_test.go`
- instruction: `internal/cli/implement_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/verification_probe.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/agent_session_owner.go`
- creates: `internal/daemon/run_core_budget_test.go`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`
- creates: `internal/agent/core_budget_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- creates: `internal/cli/verification_cores.go`
- creates: `internal/cli/verification_cores_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `daemon.Dependencies.VerificationCores` and `agent.ExecuteRequest.Cores` — consumers: `./internal/daemon`, `./internal/agent`, `./internal/cli`
- The configuration guide and the `implement` reference — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run '^TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget ' || { echo 'missing pass: TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget' >&2; exit 1; }; go test -count=1 -timeout 30m -run 'Probe|PreWork|AgentSession|VerificationCapacity' ./internal/daemon` — expected: exit 0. Before this Task the test does not exist; afterwards it passes, and the existing probe, pre-work, Agent Session and capacity tests pass unchanged.

- `out="$(go test -count=1 -timeout 30m -v -run '^TestAnAgentSessionReceivesTheCoreBudget$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestAnAgentSessionReceivesTheCoreBudget ' || { echo 'missing pass: TestAnAgentSessionReceivesTheCoreBudget' >&2; exit 1; }; go test -count=1 -timeout 30m -run 'JudgeKey|ACPX|Session' ./internal/agent` — expected: exit 0. Before this Task the test does not exist; afterwards it passes, and the existing acpx and judge-key tests pass unchanged.

- `out="$(go test -count=1 -timeout 30m -v -run '^TestImplementGivesTheEngineTheVerificationCoreBudget$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestImplementGivesTheEngineTheVerificationCoreBudget ' || { echo 'missing pass: TestImplementGivesTheEngineTheVerificationCoreBudget' >&2; exit 1; }; go test -count=1 -timeout 40m -run 'Implement' ./internal/cli` — expected: exit 0. Before this Task the test does not exist; afterwards it passes, and every existing Implement test of the CLI passes.

- `tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qE -- 'It sets .GOMAXPROCS. and .ROUNDFIX_VERIFICATION_CORES. to that number.' || { echo 'missing: configuration sentence' >&2; exit 1; }; grep -q '^## Verification Core Budget$' docs/user-guide/configuration.md || { echo 'missing: configuration section' >&2; exit 1; }; for f in .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md; do grep -qF 'ROUNDFIX_VERIFICATION_CORES' "$f" || { echo "missing variable: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither the guide nor the reference names the variable; afterwards both copies of the reference and the guide describe it, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record pass.

## References

- `_prd.md` → Goal 1; Goal 3; User Story 1; User Story 2; User Story 3; Core Feature 3; Core Feature 4; Core Feature 7; Success Metric 2; Success Metric 5
- `_techspec.md` → Current behavior 2 and 3; Interfaces; API Contract 3; API Contract 6; Invariants 1 and 2; Declared breaks; Testing Approach; Build Order 2
- ADR-0385

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, context7-cli, knowledge-workspace.

Characterization before implementation:
- `internal/daemon/engine.go`: the `verificationAttempt` path (the method is named `runVerificationAttempt` in this checkout) calls `engine.deps.Verifier.Verify` for attempts, including repository settlement.
- `internal/daemon/verification_probe.go`: `probeCommandsForRun` calls `verifier.Verify`. Its only callers are `ProbeCommands` in that file and `verifyTaskPreWork` in `internal/daemon/task_engine.go`.
- `internal/agent/acpx_runner.go`: `ensureSession` (called by `PrepareSession`) and `RunPrompt` append the Run Marker to `codexEnv` before session selection and prompt execution respectively.
- Initial inspection: the three required budget propagation tests and the CLI budget helper do not exist; requests have no Run budget wiring. The only pre-existing working-tree change is this Task's Daemon-owned status.

Implementation and acceptance evidence:
- Every Run request: the Implement Command computes the budget once from its loaded config and command environment, with injectable logical CPUs and the `runtime.NumCPU` fallback. Only `executeImplementCycle` supplies `Dependencies.VerificationCores`; attempts, pre-work probes and both branches of `runAgentSession` receive it. `ProbeCommands` explicitly supplies zero. The new Daemon test observes all Verification requests, requires a pre-work probe and the repository command in a Task-attempt settlement artifact, and checks every recorded Agent request with budgets five and zero. The new CLI test exercises `implement` with twelve logical CPUs and `runs.max_verifying: 2`, asserting three on every Verification and Agent request, or seven with explicit `GOMAXPROCS=7`.
- acpx environment: both Run Marker sites append `coreBudgetOverrides`, using the existing environment replacement path. The new fake-acpx recorder observes every invocation without credential values. The Agent test requires both `sessions ensure` and `prompt`: inherited `GOMAXPROCS=12` and `ROUNDFIX_VERIFICATION_CORES=19` become `3`/`3`; at zero, `GOMAXPROCS` stays `12` and the Roundfix variable is absent (checked with `LookupEnv`).
- Documentation: the configuration section follows Machine Verification Capacity, preserves all four required sentences, explains recipients, exclusions, rounding and explicit values, and gives the authorized Vitest opt-in example. The example's numeric `maxWorkers` API was checked against Vitest's official documentation through Context7. Item 9 of the canonical implement reference describes the same budget and Agent inheritance; `make skills-sync` generated its byte-identical mirror. The version recorder raised both Roundfix skill versions to `0.1.215`.

Focused checks (all exit 0 unless stated):
- Before implementation, `go test ./internal/agent ./internal/daemon ./internal/cli -run 'CoreBudget/' -count=1 -timeout 5m` exposed the missing `Cores`, `VerificationCores` and `logicalCPUs` fields. The first sandboxed attempt also could not access parts of the Go build cache; the approved rerun outside the sandbox confirmed the compile failures.
- `go test ./internal/agent ./internal/daemon ./internal/cli -run 'CoreBudget/' -count=1 -timeout 8m`: passed.
- `go test -race ./internal/agent ./internal/daemon ./internal/cli -run 'CoreBudget/' -count=1 -v -timeout 8m`: passed every new test and the existing ExecVerifier budget test, including all budget/zero rows.
- After strengthening the Agent recorder's absence and replacement assertions: `go test ./internal/agent -race -run 'CoreBudget/' -count=1 -timeout 5m`: passed. After identifying settlement by its artifact: `go test ./internal/daemon -race -run 'CoreBudget/' -count=1 -timeout 5m`: passed.
- `go test ./internal/agent ./internal/daemon ./internal/cli -run '^(TestProbeCommands|TestFallbackEligibilitySurvivesASelectionFailure|TestNoAgentSessionReceivesAJudgeKeyVariable|TestApplySessionSelectionDisposableAndLiveOrder|TestImplementGatedGraphRunsRepositoryVerificationAtSettlement|TestImplementRepositoryAtSettlementOffRunsOnlyDeclaredCommands|TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement)$' -count=1 -timeout 8m`: passed existing tests unchanged.
- `make skills-sync`: passed. `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed. `go test ./skills -run '^(TestAuthorialSkillSync|TestCheckValidatesRoundfixSkillArtifacts)$' -count=1`: passed.
- `make baseline-digests`: passed, reporting `changed: false`; no Baseline artifacts changed.
- Final implementation step, after all guide and skill edits: `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed; the generator refreshed only the configuration guide's digest in `docs/references/behavior-surfaces.json`.
- `git diff --check`: passed. Scope postflight accounted for tracked and untracked files, confirmed skill mirror equality, preserved authored Task text apart from Result and the pre-existing Daemon status, and verified `_tasks.md` and every sibling Task byte-identical to HEAD.

Declared Verification and the complete package/Implement suites remain for the Daemon. No status settlement, commit, push, PR or cross-project message was performed. Delivery Queue/doctor/Linux Lane implementation and the glossary remain the later Tasks' slices.

## Carry-forward provenance

- Source Run: `run_20261111T064558Z_6d985452394229fa`
- Source commit: `4a313873276b5776400cd948b42297435cbc86e8`
