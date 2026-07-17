---
task: task_06
spec: 0035-agent-selection-profiles
status: completed
type: backend
complexity: high
---

# Task 06: Route each work action through an owned Agent Session

## Overview

Replace the Run-wide session assumption with one owned session per Task, QA action, or review action. A pre-prompt selection-start failure emits durable and caller-visible notification before automatically activating the next preflight-proven fallback, including across runtimes.

## Requirements

1. MUST give each Task its own Agent Session after worktree/bootstrap readiness, QA its own session after Tasks settle, and review work a session from the `review` profile.
2. MUST create sessions through a runtime-keyed factory so a configured fallback can cross ACP Runtime without mutating configuration.
3. MUST record a failed selection start and commit the structured `agent_selection_fallback` notification before creating, preparing, or prompting the next fallback session.
4. MUST include category, scope identity, failed and next selections, fallback index, normalized reason code, human reason, and `automatic: true` in the notification.
5. MUST mark `agent_work_started` immediately before the first prompt and forbid every fallback after that marker, including prompt, tool, verification, cancellation, rate-limit, and session-loss failures.
6. MUST exhaust fallbacks only in configured order and report every attempted tuple plus a recovery action when none starts.
7. MUST close the owned session on success, failure, cancellation, and every early return without fire-and-forget goroutines.

## Subtasks

- [x] Introduce runtime-keyed Agent Session factories.
- [x] Move Task, QA, and review work to scoped session owners.
- [x] Add pre-prompt selection-start classification.
- [x] Publish notification before fallback activation.
- [x] Add the hard `agent_work_started` boundary.
- [x] Implement ordered exhaustion and cleanup reporting.
- [x] Cover same-runtime, cross-runtime, cancellation, and post-start negatives.

## Acceptance Criteria

- [x] Mixed Task Types and QA use separate sessions with the exact resolved category selections.
- [x] The notification Run Event and caller-visible message precede fallback session creation and first prompt.
- [x] A cross-runtime fallback uses the configured runtime factory and selection verbatim.
- [x] No production or test path starts a replacement session after `agent_work_started`.
- [x] Fallback exhaustion lists every failed tuple and settles the owning Task/action through existing failure semantics.
- [x] Every created session has one owner and closes under success, failure, cancellation, and early-return tests.

## Context

- instruction: `.agents/skills/golang-concurrency/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- interface: `internal/agent/sessions.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/cli.go`

## Verification

- `rtk go test ./internal/agent ./internal/daemon ./internal/cli -run 'Test(PerWorkAgentSession|AgentSelectionFallback|CrossRuntimeFallback|NoFallbackAfterAgentWorkStarted|AgentSessionOwnerCleanup)' -count=1` — expected: ownership, ordering, fallback, exhaustion, and cleanup cases pass.
- `rtk go test -race ./internal/agent ./internal/daemon ./internal/cli -run 'Test(PerWorkAgentSession|AgentSelectionFallback|AgentSessionOwnerCleanup)' -count=1` — expected: scoped session lifecycle and notification ordering are race-free.

## References

- `_prd.md` → Goals 1-3 and 9; User Stories 1-3 and 6; Core Features 9-10; Success Metrics.
- `_techspec.md` → Agent Session lifecycle and fallback activation; System Architecture; Risks: mixed-session lifecycle and unsafe replay; Build Order 6.
- `references/openclaw-skill-analysis.md` → fallback reasons are classified and visible; proof is not work quality.

## Result

- Added scoped Agent Session owners for Task, QA, and review actions, with runtime-keyed profile activation and per-work session refs.
- Added `agent_selection_fallback`, `agent_selection_exhausted`, and `agent_work_started` event paths with ordered pre-prompt fallback activation.
- Added fallback exhaustion reporting with every attempted tuple and profile validation recovery guidance.
- Added same-runtime, cross-runtime, cancellation, post-start no-fallback, review-profile, and cleanup tests.
- Verification: `GOCACHE=/tmp/roundfix-gocache rtk go test ./internal/agent ./internal/daemon ./internal/cli -run 'Test(PerWorkAgentSession|AgentSelectionFallback|CrossRuntimeFallback|NoFallbackAfterAgentWorkStarted|AgentSessionOwnerCleanup)' -count=1` passed: 7 tests.
- Verification: `GOCACHE=/tmp/roundfix-gocache rtk go test -race ./internal/agent ./internal/daemon ./internal/cli -run 'Test(PerWorkAgentSession|AgentSelectionFallback|AgentSessionOwnerCleanup)' -count=1` passed: 5 tests.
- Full gate: `GOCACHE=/tmp/roundfix-gocache rtk make verify` passed.
