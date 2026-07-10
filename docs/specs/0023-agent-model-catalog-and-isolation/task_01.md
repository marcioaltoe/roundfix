---
task: task_01
spec: 0023-agent-model-catalog-and-isolation
status: pending
type: backend
complexity: high
---

# Task 01: Resolve per-runtime Agent selections

## Overview

Establish the typed configuration, Model Catalog, and resolution contract that
all Agent-starting commands will consume. The slice is verifiable by loading
layered configurations and resolving concrete Codex, Claude, and OpenCode
selections without starting an Agent Session.

## Requirements

1. MUST support independent model and reasoning defaults for each supported ACP Runtime with the built-in values defined by the TechSpec.
2. MUST preserve built-in, User Config, Project Config, then invocation precedence independently for each selection key.
3. MUST treat an explicit empty effective value as missing selection rather than inheriting runtime-owned configuration.
4. MUST expose the exact ordered Codex and Claude Model Catalogs as picker data without using either catalog as an allowlist.
5. MUST deprecate and ignore `defaults.model` with exactly one stderr warning naming its per-runtime replacement.
6. MUST generate configuration examples that contain only the per-runtime selection structure.

## Subtasks

- [ ] Add typed per-runtime defaults and strict YAML overlays.
- [ ] Seed the Codex, Claude, and OpenCode built-in values.
- [ ] Add the effective-selection resolver with invocation-presence tracking.
- [ ] Add the ordered Model Catalog data.
- [ ] Migrate the deprecated global model key to the warning path.
- [ ] Cover generated config, overlay, resolver, and catalog behavior.

## Acceptance Criteria

- [ ] Built-in configuration resolves Codex to `gpt-5.5`/`xhigh`, Claude to `opus`/`high`, and leaves both OpenCode values empty.
- [ ] Project Config can override only one runtime key without changing the other runtime key or any User Config value it does not replace.
- [ ] Explicit invocation values override all configuration layers, while explicit empty values fail resolution.
- [ ] The Codex catalog has exactly seven entries and the Claude catalog exactly five entries in PRD order.
- [ ] A custom non-catalog model or reasoning value survives resolution unchanged for later ACP validation.
- [ ] `defaults.model` is ignored, emits one stderr warning, and never appears in generated configuration.

## Verification

- `rtk go test ./internal/config ./internal/agent ./internal/cli` - expected: configuration, catalog, and selection resolver tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/config/config.go`
- interface: `internal/agent/agent.go`
- interface: `internal/cli/cli.go`

## References

`_prd.md` -> User Stories 2, 3, 4, 5; Core Features 1-2, 5-8, 10; Success Metrics. `_techspec.md` -> Implementation Design: Data Models; API Contracts; Build Order 1. ADR-0027; ADR-0037.
