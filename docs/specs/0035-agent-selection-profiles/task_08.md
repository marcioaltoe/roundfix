---
task: task_08
spec: 0035-agent-selection-profiles
status: pending
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

- [ ] Project selection lifecycle events into the Supervisor stream.
- [ ] Render scoped selection and fallback rows in stderr and Console Logs.
- [ ] Add Live Run View and Detail Modal selection state.
- [ ] Replay identical state through Attach.
- [ ] Preserve compatibility summaries and legacy Run rendering.
- [ ] Add ordering, no-color, bounded-output, and privacy tests.

## Acceptance Criteria

- [ ] Task, QA, and review selection attempts are attributable to their category and scope in every Run surface.
- [ ] Fallback notification is ordered before fallback activation and first prompt in live and replayed views.
- [ ] Mixed-profile Runs display actual per-work selections while stable legacy summary columns remain available.
- [ ] Legacy Runs render missing selection history explicitly and do not fail Attach or Run Browser views.
- [ ] Text, TUI, Attach, and JSON projection agree on selection, role, status, and reason.
- [ ] No prompt, raw ACP payload, credential, or secret appears in the new projections.

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

## References

- `_prd.md` → Goals 3 and 9; User Stories 3, 6, and 7; Core Features 9-10; Success Metrics.
- `_techspec.md` → Persistence and Run Events; Agent Session lifecycle notification payload; JSON profile output; Build Order 7.
