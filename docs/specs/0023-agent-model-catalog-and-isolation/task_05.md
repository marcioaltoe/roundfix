---
task: task_05
spec: 0023-agent-model-catalog-and-isolation
status: completed
type: data
complexity: high
---

# Task 05: Persist effective selection for Run inspection

## Overview

Persist the concrete Agent Model and Default Reasoning Effort with each Run and
surface those historical values through existing inspection views. The slice
is verifiable by migrating a schema-v6 database, creating new Runs through each
operational path, and attaching after configuration has changed.

## Requirements

1. MUST migrate the Run Database to schema version 7 with non-null model and reasoning columns and readable legacy defaults.
2. MUST persist both effective values on every new Run that starts Agent work.
3. MUST read historical selection from the Run row rather than current User or Project Config.
4. MUST display concrete Agent, Agent Model, and Default Reasoning Effort in initial progress and the Live Run View.
5. MUST render legacy empty values as `-` and never render new Agent Runs as ambiguous `auto`.
6. MUST update every Run insert, query, scan, fixture, and attach/create path consistently.

## Subtasks

- [x] Add and migrate the Run selection columns.
- [x] Extend Run creation and query contracts.
- [x] Persist values from resolve, watch, and implement.
- [x] Render stored values in progress and Run inspection.
- [x] Cover migration, round-trip, legacy, and post-config-change behavior.

## Acceptance Criteria

- [x] Opening a schema-v6 database upgrades it to v7 without losing existing Run data.
- [x] Legacy Runs return empty stored values and display `-` for model and reasoning.
- [x] New resolve, watch, and implement Runs persist the exact values accepted by preflight.
- [x] Attach displays the original stored selection after configuration changes.
- [x] New Agent Run output contains concrete model and reasoning values and no `auto` placeholder.
- [x] Store and CLI tests cover all insert/select/scan paths without column-order drift.

## Verification

- `rtk go test ./internal/store ./internal/cli ./internal/tui` - expected: schema migration, persistence, command wiring, and inspection tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/cli/attach.go`
- interface: `internal/cli/runui.go`
- interface: `internal/tui/tui.go`
- interface: `internal/tui/agent_live.go`

## References

`_prd.md` -> User Story 1; Core Feature 9; Success Metrics. `_techspec.md` -> Data Models: SQLite schema version 7; API Contracts: inspection; Testing Approach; Build Order 5. ADR-0037.

## Result

- Added schema v7 `runs.model` and `runs.reasoning_effort` columns with non-null empty-string legacy defaults; v3, v4, v5, and v6 fixtures now migrate to user_version 7 with Run rows preserved.
- Extended Run creation, selection, scan, active lookup, listing, and latest-spec lookup paths to round-trip stored Agent Model and Default Reasoning Effort.
- Wired resolve, watch, and implement Run creation to persist the concrete selection accepted by preflight, and rendered those values in initial progress and Live Run View data.
- Updated attach to read and display the stored Run selection, verified after Project Config changes that it keeps historical values instead of current config.
- Rendered legacy empty model/reasoning values as `-` and removed the Live Run View model `auto` placeholder.
- Added store, CLI, and TUI regression coverage for schema-v6 migration, insert/select/scan round trips, resolve/watch/implement persistence, attach-after-config-change behavior, legacy dash rendering, and no ambiguous `auto` output.

Verification:

- `rtk go test ./internal/store ./internal/cli ./internal/tui` passed: 601 tests in 3 packages.
- `rtk make verify` passed: `go test ./...` reported 1039 tests in 19 packages, skill checks passed, and `go build` completed.
