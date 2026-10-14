---
task: task_02
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
status: pending
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
