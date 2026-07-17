---
task: task_01
spec: 0035-agent-selection-profiles
status: pending
type: backend
complexity: high
---

# Task 01: Resolve atomic Agent Selection Profiles

## Overview

Establish the immutable selection, profile, category, catalog, and configuration contracts that every later routing surface consumes. This slice is independently verifiable through configuration and Agent boundary tests and performs no Run creation or profile probing.

## Requirements

1. MUST define complete Agent Selection and Agent Selection Profile values with one Preferred Selection and a non-empty ordered Fallback Chain.
2. MUST validate runtime, model, explicitly present reasoning effort, duplicate tuples, empty fallbacks, and empty or partial profile overlays without turning user-supplied custom models into an allowlist.
3. MUST resolve required `general`, `backend`, `frontend`, `qa`, and `review` profiles from built-ins and optional Task Type profiles from the effective `general` profile when absent.
4. MUST implement atomic Project Config over User Config over built-in precedence and preserve the resolved fallback chain when an invocation overrides only the Preferred Selection.
5. MUST ship the exact official built-in identifiers and selections defined by the PRD, including `gpt-5.6-sol`, `gpt-5.6-terra`, `gpt-5.6-luna`, and `claude-fable-5`.
6. MUST distinguish a missing reasoning key from an explicitly empty model-managed value and render only values that the runtime will actually receive.
7. MUST support the documented legacy-runtime-default compatibility window and reject same-scope legacy/new schema conflicts with a migration action.

## Subtasks

- [ ] Define Agent Selection, Profile, Work Category, and source values.
- [ ] Normalize official built-in Model Catalog entries and labels.
- [ ] Add strict profile YAML decoding and validation.
- [ ] Add atomic precedence and optional-category inheritance.
- [ ] Preserve fallbacks under one-Run Preferred Selection overrides.
- [ ] Add legacy conversion and conflict diagnostics.
- [ ] Cover positive, negative, inheritance, and migration cases with table tests.

## Acceptance Criteria

- [ ] Every required category resolves to the specified built-in Preferred Selection and Fallback Chain when no config exists.
- [ ] An absent optional profile is byte-equivalent to the effective `general` profile and reports inheritance without duplicating stored config.
- [ ] A present higher-precedence profile replaces the lower profile as one object; no tuple field or fallback entry merges across scopes.
- [ ] Missing keys, partial profiles, empty fallbacks, duplicate tuples, and mixed legacy/new schemas fail with actionable diagnostics.
- [ ] Official built-ins render official identifiers, while explicit custom model strings survive unchanged for ACP proof.
- [ ] Invocation overrides replace only the relevant Preferred Selection and retain the configured Fallback Chain.

## Context

- instruction: `docs/adr/0040-reasoning-effort-is-assigned-only-when-configured.md`
- interface: `internal/config/config.go`
- interface: `internal/config/config_test.go`
- interface: `internal/agent/catalog.go`
- interface: `internal/agent/agent.go`

## Verification

- `rtk go test ./internal/config ./internal/agent -run 'Test(AgentSelectionProfile|ProfileResolver|ProfileLegacyMigration|ModelCatalog)' -count=1` — expected: built-ins, atomic precedence, inheritance, strict validation, migration, and official catalog cases pass.
- `rtk go test ./internal/config ./internal/agent -count=1` — expected: complete configuration and Agent package suites pass together.

## References

- `_prd.md` → Goals 1-4; User Stories 1-2 and 6-7; Core Features 1-4 and 11; Decisions.
- `_techspec.md` → Domain types; Configuration schema and precedence; Official Model Catalog and adapter proof; Build Order 1.
- `references/model-ranking.md` → official model identifiers used by built-in selections.
- ADR-0040 → reasoning effort must be applied exactly as displayed and persisted.
