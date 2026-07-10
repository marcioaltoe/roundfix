---
task: task_02
spec: 0023-agent-model-catalog-and-isolation
status: pending
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

- [ ] Extend the runtime and execution contracts with reasoning effort.
- [ ] Add the runtime-specific acpx config-option mapping.
- [ ] Apply model and reasoning in deterministic session preparation order.
- [ ] Preserve desired options across reconnect and resume paths.
- [ ] Cover all runtimes and custom Agent Commands in the acpx harness.

## Acceptance Criteria

- [ ] Codex session preparation sends the selected model followed by `reasoning_effort` and its selected value.
- [ ] Claude and OpenCode session preparation send the selected model followed by `effort` and its selected value.
- [ ] Permission mode and Codex sandbox configuration still run after the selection is accepted.
- [ ] A resumed session receives the same concrete selection instead of inheriting current runtime defaults.
- [ ] An adapter rejection identifies the failed selection operation and preserves the original error for inspection.
- [ ] A custom Agent Command that lacks the required ACP option contract fails instead of using hidden defaults.

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
