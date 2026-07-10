---
task: task_03
spec: 0023-agent-model-catalog-and-isolation
status: pending
type: infra
complexity: high
---

# Task 03: Reject unavailable selections during Preflight Validation

## Overview

Prove the effective selection against the installed ACP adapter before an
operational command creates durable Run state. The slice is verifiable by
accepting supported values and rejecting unsupported model metadata or
reasoning values with zero Run rows and zero durable Agent Sessions.

## Requirements

1. MUST validate the exact model and reasoning pair through a uniquely named disposable Agent Session in the Git root.
2. MUST run the existing acpx/runtime readiness checks before selection validation and send no Agent prompt.
3. MUST close the disposable session on success, rejection, cancellation, and partial setup failure with bounded cleanup.
4. MUST complete selection validation before `resolve`, `watch`, or `implement` creates a Run or durable Agent Session.
5. MUST fail with exit `2` and an actionable stderr diagnostic naming runtime, model, reasoning, and both recovery paths.
6. MUST never substitute a model or reasoning value when the adapter rejects the selection.

## Subtasks

- [ ] Add the workdir-aware probe request and disposable-session lifecycle.
- [ ] Distinguish selection rejection from runtime infrastructure failure.
- [ ] Join cleanup failures without losing the original validation error.
- [ ] Wire operational command preflight ahead of Run creation.
- [ ] Add adapter and CLI regression coverage for every terminal path.

## Acceptance Criteria

- [ ] A supported selection completes preflight, closes its disposable session, and permits normal Run creation.
- [ ] A missing-model-metadata rejection creates no Run row and leaves no disposable or durable Agent Session.
- [ ] A rejected reasoning value reports the exact runtime/model/reasoning tuple and no fallback attempt occurs.
- [ ] The disposable session receives model then reasoning and receives no prompt.
- [ ] Cancellation and cleanup failures preserve context/error identity and do not leak a session silently.
- [ ] Existing Doctor/setup runtime readiness checks remain functional with the evolved probe contract.

## Verification

- `rtk go test ./internal/agent ./internal/cli` - expected: disposable preflight, cleanup, rejection, zero-side-effect, and compatibility tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- interface: `internal/store/store.go`

## References

`_prd.md` -> User Story 6; Core Feature 4; User Experience; Success Metrics. `_techspec.md` -> System Architecture; Interfaces: Agent Runner preflight; API Contracts; Build Order 3. ADR-0037; ADR-0039.
