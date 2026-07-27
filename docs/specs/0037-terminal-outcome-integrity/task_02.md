---
task: task_02
spec: 0037-terminal-outcome-integrity
status: completed
type: backend
complexity: medium
---

# Task 02: Target only registered Agent Sessions during cleanup

## Overview

Use persisted Agent Selection lifecycle evidence as the exclusive Agent Session
cleanup registry. Cleanup becomes deterministic and idempotent: only scopes
whose latest lifecycle is active are targeted, while absent, failed, and
already closed sessions produce no misleading cleanup action.

## Requirements

1. MUST query the latest Agent Selection lifecycle per Run scope.
2. MUST return only scopes whose latest persisted lifecycle is active.
3. MUST cancel and close eligible Agent Sessions once in deterministic scope
   order.
4. MUST record closed lifecycle after an idempotent close.
5. MUST treat a registered but already absent Agent Session as an idempotent
   close.
6. MUST not derive or target an Agent Session without active lifecycle evidence.
7. MUST retain non-absence cleanup failures as secondary diagnostics.

## Subtasks

- [x] Add the latest-active-scope store query.
- [x] Route cleanup through registered Agent Selection scopes.
- [x] Make registered-session absence idempotent.
- [x] Persist closed lifecycle after successful cleanup.
- [x] Add deterministic ordering and lifecycle-state tests.
- [x] Remove unconditional Run-wide session-name cleanup.

## Acceptance Criteria

- [x] Active Task, QA, and review scopes are returned once in stable order.
- [x] Failed, closed, and superseded lifecycle attempts are not targeted.
- [x] A Run with no active lifecycle record performs zero Agent Session calls.
- [x] An already absent registered session closes without a warning.
- [x] Other cleanup failures remain visible and do not invent lifecycle state.
- [x] Existing Agent Selection history and sensitive-field protections pass.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/store/agent_selection.go`
- interface: `internal/store/agent_selection_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/agent/sessions.go`

## Verification

- `rtk go test ./internal/store -run 'TestAgentSelection.*(Active|Lifecycle|Cleanup|Scope)' -count=1`
  — expected: only latest active scopes are returned in deterministic order.
- `rtk go test ./internal/cli -run 'Test.*(RegisteredAgentSession|AgentSessionCleanup|PrimaryFailure)' -count=1`
  — expected: cleanup targets registered scopes only and registered absence is
  idempotent.

## References

- `_prd.md` → Goal 4; User Story 4; Core Features 5–6; Success Metrics.
- `_techspec.md` → Data Models; API Contracts: Agent Session cleanup rules;
  Build Order 2.
- `../../adr/0051-tasks-and-qa-own-agent-sessions.md` → Work Item-scoped Agent
  Sessions.

## Result

Implemented persisted Agent Selection lifecycle as the exclusive Force Stop
Agent Session cleanup registry. Cleanup now reads each scope's latest
lifecycle, targets only active Task, QA, and review scopes in deterministic
order, reconstructs the exact fallback session name and runtime, and records
`closed` only after close succeeds or acpx proves the registered session is
already absent. Non-absence cancel and close errors remain visible; a failed
close leaves the lifecycle active for a later retry.

Verification:

- `rtk go test ./internal/store -run 'TestAgentSelection.*(Active|Lifecycle|Cleanup|Scope)' -count=1`
  — passed, 3 tests.
- `rtk go test ./internal/cli -run 'Test.*(RegisteredAgentSession|AgentSessionCleanup|PrimaryFailure)' -count=1`
  — passed, 4 tests.
- `rtk go test ./internal/agent ./internal/store ./internal/cli -count=1`
  — passed, 1,108 tests.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

Acceptance evidence:

- `TestAgentSelectionActiveScopesReturnsLatestLifecycleInStableOrder` proves
  Task, QA, and review scope ordering and excludes failed, closed, and
  superseded lifecycles.
- `TestRunStopForceAgentSessionCleanupSkipsRunWithoutActiveLifecycle` proves
  zero cancel and close calls without active persisted evidence.
- `TestRunStopForceRegisteredAgentSessionCleanupTargetsActiveScopesInOrder`
  proves one ordered cancel/close pair per registered scope, including a
  fallback Agent Session, followed by persisted `closed` lifecycles.
- `TestRunStopForceRegisteredAgentSessionAbsenceIsIdempotent` proves wrapped
  acpx missing-session responses are silent and still persist `closed`.
- `TestRunStopForceAgentSessionCleanupFailureRemainsVisibleWithoutClosedLifecycle`
  proves other cleanup failures remain diagnostic and do not fabricate a
  closed lifecycle.
- The affected-package run includes the existing Agent Selection history,
  lifecycle transition, event payload, schema privacy, and sensitive-field
  protection tests.

Follow-up: TechSpec Build Order 5 still owns primary-before-secondary
publication ordering and winner-only terminal outcome wiring.
