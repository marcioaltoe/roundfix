---
task: task_07
spec: 0035-agent-selection-profiles
status: completed
type: data
complexity: high
---

# Task 07: Persist Agent Selection attempts and events

## Overview

Create the append-only durable record for the actual selections attempted by each Task, QA action, and review action. The schema and Run Events preserve exact reasoning, ordered fallback history, normalized failures, and compatibility summaries without storing prompts or credentials.

## Requirements

1. MUST migrate the Run Database from schema 8 to 9 and create `run_agent_selections` with the fields and constraints defined by the TechSpec.
2. MUST support `task`, `qa`, and `review` scopes, Preferred and fallback roles, monotonic attempt numbers, and `attempting`, `active`, `failed`, and `closed` statuses.
3. MUST round-trip explicit empty reasoning, profile source, runtime, model, fallback index, normalized reason code, human reason, and timestamps without normalization drift.
4. MUST append structured selection-attempt, fallback, active, exhausted, and closed Run Events from the same selection lifecycle payload.
5. MUST keep existing Run-level agent/model/reasoning columns as non-null compatibility summaries using effective `general` for Spec Runs and `review` for review Runs.
6. MUST reject invalid scope, role, status, duplicate/non-monotonic attempts, and unknown Run ids without partial rows or events.
7. MUST store no prompt, transcript, credential, token, cookie, secret, or runtime-owned configuration.

## Subtasks

- [x] Add clean schema 9 creation and schema 8→9 migration.
- [x] Add typed selection-attempt persistence values.
- [x] Enforce append and monotonic attempt invariants transactionally.
- [x] Add lifecycle Run Event kinds and payloads.
- [x] Preserve Run compatibility summaries.
- [x] Cover migration, round-trip, rejection, and privacy boundaries.

## Acceptance Criteria

- [x] Existing schema 8 databases migrate without losing Runs or Run Events, and fresh databases create schema 9 directly.
- [x] Task, QA, and review histories round-trip in attempt order with exact selection and reason values.
- [x] Invalid or out-of-order appends leave the database unchanged.
- [x] Lifecycle events reference the same scope and selection data as persisted attempts.
- [x] Legacy Run readers retain non-null summary values while new readers can reconstruct actual per-scope history.
- [x] Schema and payload inspection finds none of the prohibited sensitive-content fields.

## Context

- interface: `internal/store/store.go`
- interface: `internal/store/store_test.go`
- interface: `internal/store/journal.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/stream.go`

## Verification

- `rtk go test ./internal/store ./internal/runevent -run 'Test(AgentSelectionAttempts|Schema8To9|Schema9|AgentSelectionEvent|SelectionAttemptOrdering)' -count=1` — expected: migration, clean creation, append constraints, exact round-trip, events, and privacy cases pass.
- `rtk go test -race ./internal/store ./internal/runevent -run 'Test(AgentSelectionAttempts|AgentSelectionEvent)' -count=1` — expected: concurrent readers and lifecycle append paths are race-free.

## References

- `_prd.md` → Goals 3 and 9; User Stories 3, 6, and 7; Core Feature 10; Success Metrics.
- `_techspec.md` → Persistence and Run Events; JSON profile output compatibility context; Build Order 7.

## Result

Implemented schema 9 selection persistence and lifecycle events.

Evidence:

- Schema 9 creation and v8→v9 migration are covered by `TestSchema9CreatesAgentSelectionTable` and `TestSchema8To9MigratesRunsEventsAndSelectionTable`; the migration preserves existing Runs and Run Events.
- Task, QA, and review scope histories round-trip in order through `TestAgentSelectionAttemptsRoundTripTaskQAReviewHistories`, including explicit empty reasoning, profile source, fallback index, reason code, reason, and timestamps.
- Transactional rejection is covered by `TestSelectionAttemptOrderingRejectsInvalidOrOutOfOrderAppendsWithoutPartialRows`; invalid scope, role, status, duplicate/gap attempts, and unknown Runs leave rows and Run Events unchanged.
- Lifecycle Run Events are covered by `TestAgentSelectionAttemptsRoundTripTaskQAReviewHistories`, `TestAgentSelectionAttemptsAppendExhaustedEventFromPersistedHistory`, and `TestAgentSelectionEventProjectsStablePayloadAndKeepsReasoningEmpty`.
- Compatibility summaries remain on `runs.agent`, `runs.model`, and `runs.reasoning_effort`; schema 8 migration and existing Run reader tests confirm legacy readers keep non-null summary values while `run_agent_selections` reconstructs per-scope history.
- Privacy is covered by `TestAgentSelectionAttemptsStoreNoSensitiveFields` and payload assertions that reject prompt, transcript, credential, token, cookie, secret, and runtime-owned configuration fields.

Verification:

- `rtk go test ./internal/store ./internal/runevent -run 'Test(AgentSelectionAttempts|Schema8To9|Schema9|AgentSelectionEvent|SelectionAttemptOrdering)' -count=1` — passed, 14 tests.
- `rtk go test -race ./internal/store ./internal/runevent -run 'Test(AgentSelectionAttempts|AgentSelectionEvent)' -count=1` — passed, 5 tests.
- `rtk go test ./internal/store ./internal/runevent -count=1` — passed, 120 tests.
- `rtk make verify` — passed; 1,535 Go tests, setup-context-driven unittest suite, skills check, and build completed.
