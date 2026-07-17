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

- [x] Build the temporary-repository and fake-adapter macro fixture.
- [x] Exercise file configuration, show, and validation flows.
- [x] Execute mixed backend/frontend Tasks and QA.
- [x] Force and observe notification-first cross-runtime fallback.
- [x] Assert database and Run Event Stream selection history.
- [x] Prove invalid Task Type has zero side effects.
- [x] Run focused race and full repository verification.
- [ ] Preserve a valid signed Codex executable boundary in the macOS macro fixture.

## Acceptance Criteria

- [x] Backend, frontend, and QA actions each use their configured exact selection and independently owned session.
- [x] A pre-prompt Preferred failure continues through the configured cross-runtime fallback only after notification is durable and visible.
- [x] Stored attempts and streamed events reproduce the exact category, source, order, tuple, status, and reason history.
- [x] Invalid Task Type fails before every proof and Run side effect.
- [x] Recommendation ordering never changes configuration or fallback order during the flow.
- [x] No fallback occurs after the first prompt in the negative companion.
- [x] Focused race verification and `make verify` pass without retries or skipped assertions.
- [ ] A cold-cache macro run passes macOS Codex hygiene without copying away the executable signature.

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

## Result

Implemented the end-to-end macro proof and closed the persistence gaps it exposed:

- Added `TestAgentSelectionProfilesMacro`, which builds the real `roundfix` binary, creates temporary git repositories, runs fake pinned `acpx` plus adapter executables through `PATH`, configures project profiles from YAML, shows and validates profiles, executes backend/frontend Tasks plus QA, and exercises invalid Task Type and post-start failure companions.
- Persisted live Agent Selection lifecycle rows from owned sessions into `run_agent_selections`, including active, failed, and closed statuses, profile source, role, fallback index, exact tuple, normalized reason, and continued monotonic attempt numbering when review scope IDs repeat across watch rounds.
- Fixed Spec Run compatibility summaries to store the effective `general` Preferred Selection while leaving actual per-work Task/QA selections in the per-scope history.

Acceptance evidence:

- Backend/frontend/QA exact selections and owned sessions: the macro asserts backend `codex/macro-backend/high`, frontend failed Preferred `codex/macro-frontend-preferred/high` then cross-runtime fallback `claude/claude-fable-5/xhigh`, QA `codex/macro-qa/high`, and one close per owned session.
- Notification-first fallback: fake `acpx` refuses to create the fallback session unless the durable fallback notification already exists in SQLite, while the test also asserts the caller-visible fallback message and Run Event order before fallback active and `agent_work_started`.
- Stored and streamed history: the macro reads `run_agent_selections` and `roundfix events --filter agent-selection`, asserting category, source `project`, attempt order, role, fallback index, status, tuple, `model_unavailable`, and sanitized reason values.
- Invalid Task Type side effects: the companion invalid frontmatter run exits before proofs and asserts zero fake `acpx` calls, no Run Database, no Run Worktree root, and no `roundfix/*` branch.
- No unauthorized mutation: the macro snapshots Project Config, fake runtime-owned files, fake credential files, recommendations, and fallback order before/after the flow.
- No fallback after work starts: the post-start failure companion asserts no live fallback session, no fallback notification, and no persisted fallback attempt after the first prompt.

Verification:

- `rtk go test ./internal/cli -run 'TestAgentSelectionProfilesMacro' -count=1` — passed, 4 macro subtests.
- `rtk go test -race ./internal/config ./internal/spec ./internal/agent ./internal/store ./internal/runevent ./internal/cli ./internal/daemon ./internal/tui -run 'Test(AgentSelection|Profile|TaskType)' -count=1` — passed, 62 tests in 8 packages.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — passed.
- `rtk make verify` — passed, including `rtk go test ./...` with 1548 tests, setup-context checks, `roundfix skills check`, and build.
