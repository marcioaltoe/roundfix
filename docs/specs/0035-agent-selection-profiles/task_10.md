---
task: task_10
spec: 0035-agent-selection-profiles
status: pending
type: test
complexity: high
---

# Task 10: Prove mixed-profile fallback behavior end to end

## Overview

Add the cross-boundary macro harness that proves the individually tested profile components cooperate through the real CLI, files, SQLite store, Run Event Stream, and fake pinned adapters. This task adds no routing behavior; it exercises mixed Task Types, QA, automatic fallback, invalid authoring, race safety, and the full repository gate.

## Requirements

1. MUST build the real binary in a temporary repository and use fake pinned ACP adapters through the production process boundary without production-only hooks.
2. MUST configure complete profiles from a file, show and validate them, and run a backend/frontend Task Graph plus QA with the exact expected profile per action.
3. MUST force one Preferred Selection to fail before its first prompt and prove caller-visible plus structured notification precede cross-runtime fallback creation and continued execution.
4. MUST assert persisted selection attempts, profile sources, roles, exact tuples, fallback reasons, compatibility summaries, and Run Event Stream order.
5. MUST run an invalid Task Type companion and prove zero disposable probes, Run rows, branches, worktrees, and Agent invocations.
6. MUST snapshot configuration, runtime-owned state, credentials boundary, and recommendation data to prove no unauthorized mutation.
7. MUST use event/state synchronization rather than wall-clock sleeps or retries and pass the focused race and full verification gates.

## Subtasks

- [ ] Build the temporary-repository and fake-adapter macro fixture.
- [ ] Exercise file configuration, show, and validation flows.
- [ ] Execute mixed backend/frontend Tasks and QA.
- [ ] Force and observe notification-first cross-runtime fallback.
- [ ] Assert database and Run Event Stream selection history.
- [ ] Prove invalid Task Type has zero side effects.
- [ ] Run focused race and full repository verification.

## Acceptance Criteria

- [ ] Backend, frontend, and QA actions each use their configured exact selection and independently owned session.
- [ ] A pre-prompt Preferred failure continues through the configured cross-runtime fallback only after notification is durable and visible.
- [ ] Stored attempts and streamed events reproduce the exact category, source, order, tuple, status, and reason history.
- [ ] Invalid Task Type fails before every proof and Run side effect.
- [ ] Recommendation ordering never changes configuration or fallback order during the flow.
- [ ] No fallback occurs after the first prompt in the negative companion.
- [ ] Focused race verification and `make verify` pass without retries or skipped assertions.

## Context

- instruction: `.agents/skills/testing-boss/SKILL.md`
- instruction: `.agents/skills/qa-gate/SKILL.md`
- interface: `internal/cli/implement_test.go`
- interface: `internal/agent/acpx_integration_test.go`
- interface: `internal/store/store_test.go`
- interface: `internal/runevent/stream.go`

## Verification

- `rtk go test ./internal/cli -run 'TestAgentSelectionProfilesMacro' -count=1` — expected: complete configure/show/validate, mixed Task, QA, fallback, persistence, stream, and invalid-type flows pass.
- `rtk go test -race ./internal/config ./internal/spec ./internal/agent ./internal/store ./internal/runevent ./internal/cli ./internal/daemon ./internal/tui -run 'Test(AgentSelection|Profile|TaskType)' -count=1` — expected: selection, lifecycle, persistence, and presentation contracts are race-free.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: shipped authorial workflow skills are valid and synchronized.
- `rtk make verify` — expected: the complete repository gate passes with no format, test, skill, or build failure.

## References

- `_prd.md` → all Goals; User Stories 1-7; Core Features 1-11; Success Metrics; Non-Goals.
- `_techspec.md` → Testing Approach: Macro QA; Required verification; Build Order 8; Risks and Decisions.
- `references/model-ranking.md` → ranking is advisory and never routing input.
- `references/openclaw-skill-analysis.md` → proof, review, and QA remain distinct concerns.
