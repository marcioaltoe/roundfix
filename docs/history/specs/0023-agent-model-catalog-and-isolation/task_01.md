---
task: task_01
spec: 0023-agent-model-catalog-and-isolation
status: completed
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

- [x] Add typed per-runtime defaults and strict YAML overlays.
- [x] Seed the Codex, Claude, and OpenCode built-in values.
- [x] Add the effective-selection resolver with invocation-presence tracking.
- [x] Add the ordered Model Catalog data.
- [x] Migrate the deprecated global model key to the warning path.
- [x] Cover generated config, overlay, resolver, and catalog behavior.

## Acceptance Criteria

- [x] Built-in configuration resolves Codex to `gpt-5.5`/`xhigh`, Claude to `opus`/`high`, and leaves both OpenCode values empty.
- [x] Project Config can override only one runtime key without changing the other runtime key or any User Config value it does not replace.
- [x] Explicit invocation values override all configuration layers, while explicit empty values fail resolution.
- [x] The Codex catalog has exactly seven entries and the Claude catalog exactly five entries in PRD order.
- [x] A custom non-catalog model or reasoning value survives resolution unchanged for later ACP validation.
- [x] `defaults.model` is ignored, emits one stderr warning, and never appears in generated configuration.

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

## Result

Implemented the Build Order 1 foundations without starting any Agent Session:
typed per-runtime config defaults and overlays, built-in Codex/Claude/OpenCode
selection values, an invocation-aware selection resolver, ordered Codex and
Claude Model Catalogs, and the deprecated `defaults.model` warning path.
Generated config now uses only the `runtimes.<runtime>.model` and
`runtimes.<runtime>.reasoning_effort` selection structure.

Evidence by acceptance criterion:

- Built-ins: `TestBuiltinRuntimeDefaults` verifies Codex `gpt-5.5`/`xhigh`,
  Claude `opus`/`high`, and empty OpenCode defaults;
  `TestResolveSelectionUsesBuiltInRuntimeDefaults` verifies Codex and Claude
  resolve from those defaults while OpenCode fails without explicit values.
- Per-key overlay: `TestLoadAppliesRuntimeConfigHierarchy` verifies Project
  Config can override one runtime key while unrelated User Config runtime keys
  survive.
- Invocation precedence: `TestResolveSelectionAppliesInvocationPrecedence`
  verifies explicit invocation values override config;
  `TestResolveSelectionRejectsExplicitEmptyInvocationValues` verifies explicit
  empty model and reasoning values fail resolution.
- Catalog order: `TestModelCatalogsExposeOrderedPickerData` pins the seven
  Codex entries and five Claude entries in PRD order;
  `TestModelCatalogLeavesOpenCodeWithoutBuiltInChoices` verifies OpenCode has
  no fabricated catalog.
- Custom values: `TestResolveSelectionPreservesCustomValues` verifies
  non-catalog model and reasoning strings survive resolution unchanged.
- Deprecated key: `TestLoadWarnsAndIgnoresDeprecatedDefaultsModel` verifies
  `defaults.model` is stripped, ignored, and warned once with
  `runtimes.<runtime>.model`; `TestInitCreatesUserConfig` verifies generated
  config omits `defaults.model` and contains the per-runtime structure.

Verification:

- `rtk go test ./internal/config ./internal/agent ./internal/cli` passed:
  515 tests passed in 3 packages.
- First `rtk make verify` attempt failed in `internal/daemon`
  `TestTaskCycleTaskWorktreeBootstrapFailureIsolatesIndependentTasks`. The
  test passed when rerun directly, `rtk go test ./internal/daemon -count=1`
  passed, and `rtk go test ./...` passed with 996 tests in 19 packages.
- Final `rtk make verify` passed: `rtk go test ./...`, `roundfix skills
  check`, and `go build` all exited 0.

Follow-up:

- CLI flag wiring, ACP preflight, persisted Run fields, Interactive Input
  behavior, and shipped guidance updates remain in later tasks from the
  TechSpec build order.
