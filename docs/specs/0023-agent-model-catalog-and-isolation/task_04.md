---
task: task_04
spec: 0023-agent-model-catalog-and-isolation
status: pending
type: frontend
complexity: high
---

# Task 04: Expose per-Run selection controls

## Overview

Expose the resolved selection contract to non-interactive callers and
Interactive Input without restricting forward-compatible custom values. The
slice is verifiable from command parsing and input collection for every
Agent-starting command and supported ACP Runtime.

## Requirements

1. MUST retain `--model` and add `--reasoning-effort` to `resolve`, `watch`, and `implement` with invocation precedence.
2. MUST distinguish an omitted flag from an explicitly empty flag and reject the latter.
3. MUST collect Agent, Agent Model, then Default Reasoning Effort during Interactive Input.
4. MUST present the exact ordered Codex and Claude catalogs while accepting a catalog number or typed custom value.
5. MUST show the concrete value behind Claude `Default` and use the configured custom value as the input default even when it is absent from the catalog.
6. MUST require typed/configured OpenCode values without fabricating an OpenCode catalog.
7. MUST keep stdout limited to requested command output and all validation guidance on stderr.

## Subtasks

- [ ] Add reasoning fields to command and Interactive Input values.
- [ ] Parse flag presence and resolve one-Run overrides.
- [ ] Render runtime-specific model and reasoning choices.
- [ ] Support catalog numbers, configured defaults, and custom typed values.
- [ ] Cover non-interactive and Interactive Input behavior for all runtimes.

## Acceptance Criteria

- [ ] Each Agent-starting command accepts both one-Run overrides and passes their concrete values to preflight.
- [ ] Omitted overrides use the selected runtime's effective configuration, while explicit empty overrides exit `2`.
- [ ] Codex Interactive Input displays seven ordered models and Claude displays five, with no extra `Custom` catalog entry.
- [ ] Claude `Default` visibly resolves to the configured concrete model.
- [ ] A custom model or reasoning value reaches ACP validation unchanged.
- [ ] OpenCode Interactive Input cannot proceed while either required value is empty.
- [ ] Buffer-captured CLI tests find no diagnostics on stdout.

## Verification

- `rtk go test ./internal/tui ./internal/cli` - expected: field ordering, catalogs, custom values, flags, precedence, and stdout/stderr tests pass.
- `rtk go run -buildvcs=false ./cmd/roundfix --help` - expected: command help renders successfully with truthful selection flags.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- interface: `internal/tui/tui.go`
- interface: `internal/agent/agent.go`

## References

`_prd.md` -> User Stories 3, 4, 5; Core Features 5-8; User Experience. `_techspec.md` -> API Contracts; Data Models: Model Catalogs; Build Order 4. ADR-0037; ADR-0039.
