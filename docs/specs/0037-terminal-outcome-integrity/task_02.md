---
task: task_02
spec: 0037-terminal-outcome-integrity
status: pending
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

- [ ] Add the latest-active-scope store query.
- [ ] Route cleanup through registered Agent Selection scopes.
- [ ] Make registered-session absence idempotent.
- [ ] Persist closed lifecycle after successful cleanup.
- [ ] Add deterministic ordering and lifecycle-state tests.
- [ ] Remove unconditional Run-wide session-name cleanup.

## Acceptance Criteria

- [ ] Active Task, QA, and review scopes are returned once in stable order.
- [ ] Failed, closed, and superseded lifecycle attempts are not targeted.
- [ ] A Run with no active lifecycle record performs zero Agent Session calls.
- [ ] An already absent registered session closes without a warning.
- [ ] Other cleanup failures remain visible and do not invent lifecycle state.
- [ ] Existing Agent Selection history and sensitive-field protections pass.

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
