---
task: task_08
spec: 0035-agent-selection-profiles
status: completed
type: frontend
complexity: high
---

# Task 08: Project per-work selection state in the Live Run surfaces

## Overview

Make actual Task, QA, and review selections observable through stderr, the Live Run View, Attach replay, and the Supervisor Run Event Stream. Every surface derives from the same structured lifecycle payload so fallback timing, source, role, and failure reason cannot disagree.

## Requirements

1. MUST project category, scope identity, profile source, Preferred/fallback role, attempt order, exact selection, activation status, and failure reason for each work scope.
2. MUST render the fallback notification before any fallback-active or first-prompt indication in stderr, TUI timeline, Attach replay, and Run Event Stream order.
3. MUST add stable machine-readable selection lifecycle projection without exposing raw Agent payloads, prompts, credentials, or runtime-owned configuration.
4. MUST show the effective actual selection for mixed-profile Runs rather than presenting compatibility summary columns as if one tuple drove every Task.
5. MUST preserve existing Run list columns and legacy Runs, rendering unavailable per-scope selection history explicitly rather than inventing it.
6. MUST keep compact rows bounded and put full structured detail behind the existing Detail Modal or machine event record.
7. MUST preserve text markers and ordering under no-color output and detached Console Logs.

## Subtasks

- [x] Project selection lifecycle events into the Supervisor stream.
- [x] Render scoped selection and fallback rows in stderr and Console Logs.
- [x] Add Live Run View and Detail Modal selection state.
- [x] Replay identical state through Attach.
- [x] Preserve compatibility summaries and legacy Run rendering.
- [x] Add ordering, no-color, bounded-output, and privacy tests.
- [x] Document `agent-selection` in the public `events --filter` help and user guidance.

## Acceptance Criteria

- [x] Task, QA, and review selection attempts are attributable to their category and scope in every Run surface.
- [x] Fallback notification is ordered before fallback activation and first prompt in live and replayed views.
- [x] Mixed-profile Runs display actual per-work selections while stable legacy summary columns remain available.
- [x] Legacy Runs render missing selection history explicitly and do not fail Attach or Run Browser views.
- [x] Text, TUI, Attach, and JSON projection agree on selection, role, status, and reason.
- [x] No prompt, raw ACP payload, credential, or secret appears in the new projections.
- [x] `roundfix events --help` lists `agent-selection` as an accepted filter, and the user guide documents the same value.

## Context

- instruction: `.agents/skills/tui-design/SKILL.md`
- instruction: `.agents/skills/bubbletea/SKILL.md`
- interface: `internal/runevent/stream.go`
- interface: `internal/tui/cockpit.go`
- interface: `internal/tui/tui.go`
- interface: `internal/cli/attach.go`
- interface: `internal/cli/runui.go`

## Verification

- `rtk go test ./internal/runevent ./internal/tui ./internal/cli -run 'Test(AgentSelectionStream|AgentSelectionLiveRunView|AgentSelectionAttachReplay|FallbackNotificationOrdering|LegacyRunSelectionView)' -count=1` — expected: scoped projection, ordering, replay, compatibility, no-color, and privacy cases pass.
- `rtk go test -race ./internal/runevent ./internal/tui ./internal/cli -run 'Test(AgentSelectionLiveRunView|AgentSelectionAttachReplay|FallbackNotificationOrdering)' -count=1` — expected: live and replayed selection state is race-free.
- `rtk go test ./internal/cli -run 'TestEvents.*Help|TestAgentSelection.*Event' -count=1` — expected: the public help contract includes the accepted `agent-selection` filter.

## References

- `_prd.md` → Goals 3 and 9; User Stories 3, 6, and 7; Core Features 9-10; Success Metrics.
- `_techspec.md` → Persistence and Run Events; Agent Session lifecycle notification payload; JSON profile output; Build Order 7.

## Result

- Task, QA, and review scope attribution: `internal/runevent/selection.go` projects structured lifecycle payloads; `internal/tui/tui.go`, `internal/tui/cockpit.go`, `internal/cli/attach.go`, and `internal/cli/runui.go` render those records without using compatibility summaries as actual per-work selections. Evidence: `TestAgentSelectionStreamProjectsScopedLifecycleWithoutSensitivePayload`, `TestAgentSelectionLiveRunViewRendersActualPerWorkSelections`, `TestAgentSelectionLiveRunViewDetailModalShowsScopedSelection`, and `TestAgentSelectionAttachReplayRendersPerScopeSelectionState`.
- Fallback ordering: stderr, timeline, and attach replay render `agent_selection_fallback` before fallback active and `agent_work_started`. Evidence: `TestFallbackNotificationOrderingSelectionConsoleSink` and `TestFallbackNotificationOrderingInTimeline`.
- Legacy compatibility: existing Agent/Model/Reasoning columns remain rendered, and Runs without per-scope selection history show `per-scope history: unavailable` instead of inventing records. Evidence: `TestLegacyRunSelectionViewRendersUnavailableHistory`.
- Machine-readable and privacy contract: the Supervisor stream exposes selection scope, category, source, role, status, reason, failed tuple, and next tuple fields while omitting prompt, transcript, credential, token, cookie, and secret data. Evidence: `TestAgentSelectionStreamProjectsScopedLifecycleWithoutSensitivePayload`.
- Verification: `rtk go test ./internal/runevent ./internal/tui ./internal/cli -run 'Test(AgentSelectionStream|AgentSelectionLiveRunView|AgentSelectionAttachReplay|FallbackNotificationOrdering|LegacyRunSelectionView)' -count=1` passed with 7 tests in 3 packages.
- Verification: `rtk go test -race ./internal/runevent ./internal/tui ./internal/cli -run 'Test(AgentSelectionLiveRunView|AgentSelectionAttachReplay|FallbackNotificationOrdering)' -count=1` passed with 5 tests in 3 packages.
- Full gate: `rtk make verify` passed, including `rtk go test ./...`, setup-context-driven skill tests, `roundfix skills check`, and the Roundfix build.
- Public `agent-selection` filter guidance: `roundfix events --help` now lists `task-status,batch,verification,outcome,agent-selection`, and `docs/user-guide/commands.md` documents the same accepted `--filter` value set. Evidence: `TestEventsHelpDocumentsAgentSelectionFilter` asserts command help, parser acceptance, and user guide text.
- Verification: `rtk go test ./internal/runevent ./internal/tui ./internal/cli -run 'Test(AgentSelectionStream|AgentSelectionLiveRunView|AgentSelectionAttachReplay|FallbackNotificationOrdering|LegacyRunSelectionView)' -count=1` passed with 7 tests in 3 packages.
- Verification: `rtk go test -race ./internal/runevent ./internal/tui ./internal/cli -run 'Test(AgentSelectionLiveRunView|AgentSelectionAttachReplay|FallbackNotificationOrdering)' -count=1` passed with 5 tests in 3 packages.
- Verification: `rtk go test ./internal/cli -run 'TestEvents.*Help|TestAgentSelection.*Event' -count=1` passed with 1 test in 1 package.
- Full gate: `rtk make verify` passed, including `rtk go test ./...`, setup-context-driven skill tests, `roundfix skills check`, and the Roundfix build.
