---
task: task_02
spec: 0023-agent-model-catalog-and-isolation
status: completed
type: backend
complexity: high
---

# Task 02: Apply explicit selection to Agent Sessions

## Overview

Carry the resolved Agent Model and Default Reasoning Effort through the Agent
boundary and apply both to every real acpx session. The slice is verifiable
from the adapter command stream for all supported runtimes, including resumed
sessions and custom Agent Commands.

## Requirements

1. MUST require a concrete model and reasoning value in the runtime contract used for Agent work.
2. MUST map Codex reasoning to `reasoning_effort` and Claude/OpenCode reasoning to `effort`.
3. MUST assign model before reasoning and before the existing permission and sandbox configuration.
4. MUST repeat the desired selection when acpx reconnects or resumes the Agent Session.
5. MUST never inspect, mutate, or rely on runtime-owned model configuration.
6. MUST preserve cancellation, wrapped adapter errors, full-access behavior, and custom Agent Command support.

## Subtasks

- [x] Extend the runtime and execution contracts with reasoning effort.
- [x] Add the runtime-specific acpx config-option mapping.
- [x] Apply model and reasoning in deterministic session preparation order.
- [x] Preserve desired options across reconnect and resume paths.
- [x] Cover all runtimes and custom Agent Commands in the acpx harness.

## Acceptance Criteria

- [x] Codex session preparation sends the selected model followed by `reasoning_effort` and its selected value.
- [x] Claude and OpenCode session preparation send the selected model followed by `effort` and its selected value.
- [x] Permission mode and Codex sandbox configuration still run after the selection is accepted.
- [x] A resumed session receives the same concrete selection instead of inheriting current runtime defaults.
- [x] An adapter rejection identifies the failed selection operation and preserves the original error for inspection.
- [x] A custom Agent Command that lacks the required ACP option contract fails instead of using hidden defaults.

## Verification

- `rtk go test ./internal/agent` - expected: runtime mapping, acpx argument order, reconnect, custom-command, and error-path tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`

## References

`_prd.md` -> User Story 1; Core Feature 3; Non-Goals. `_techspec.md` -> Interfaces: Agent Runner preflight; Integration Points; Build Order 2. ADR-0037.

## Result

Status: completed.

Acceptance evidence:

- Codex preparation: `TestACPXRunAppliesSelectionBeforePrompt/codex_reasoning_effort` asserts `sessions ensure --model gpt-5.5`, then `set reasoning_effort xhigh`, then `prompt --model gpt-5.5`.
- Claude and OpenCode preparation: `TestACPXRunAppliesSelectionBeforePrompt/claude_effort` and `/opencode_effort` assert `set effort <value>` after model assignment.
- Permission and sandbox ordering: `TestACPXRunAppliesFullAccessSessionSetup` asserts model and reasoning are accepted before `set-mode` and Codex `set sandbox_mode danger-full-access`.
- Resume behavior: `TestACPXRunReappliesSelectionForFreshRunnerSessionResume` asserts a fresh runner reapplies the same concrete model and reasoning to an existing Agent Session.
- Adapter rejection: `TestACPXRunSelectionSetupErrorsPreserveAdapterFailure` asserts the failure is wrapped with `set acpx Agent Session reasoning_effort` and preserves the original `InfrastructureError` stderr.
- Custom Agent Command: `TestACPXRunCustomCommandRequiresSelectionOptionContract` asserts a command override uses `--agent`, attempts the explicit reasoning option, and fails when the ACP option contract rejects it.

Verification:

- `rtk go test ./internal/agent` passed: `Go test: 101 passed in 1 packages`.
- `rtk make verify` passed: `Go test: 1006 passed in 19 packages`; `roundfix skills check` passed; `go build` completed.
