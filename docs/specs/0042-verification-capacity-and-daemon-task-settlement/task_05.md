---
task: task_05
spec: 0042-verification-capacity-and-daemon-task-settlement
status: pending
type: frontend
complexity: medium
---

# Task 05: Render per-Task Verification phases in the Live Run View

## Overview

Make Attach and the Live Run View distinguish Agent work, Waiting for
Verification, active Verification, and terminal Task state under concurrent
execution. The existing layout, navigation, panel positions, and review Run
experience stay stable while spec Runs show both effective capacities and
derive Task truth from journal events.

## Requirements

1. MUST add Verification Capacity beside Task Capacity in the spec Run header
   and use the canonical labels rather than the generic `Concurrency` label.
2. MUST replay effective capacities from the Task-cycle-start event for
   attached Runs and provide a deterministic legacy fallback when new fields
   are absent.
3. MUST fold interleaved per-Task `daemon.task` and `daemon.verification`
   events into the exact labels `Agent working`, `Waiting for Verification`,
   `Verifying`, `Completed`, `Failed`, `Skipped`, `Paused`, and `Waiting`.
4. MUST map initial Task start and Verification Feedback to Agent working,
   `waiting` to Waiting for Verification, acquired attempt phases to Verifying,
   and settlement to terminal status.
5. MUST treat per-Task journal evidence as authoritative when concurrent Tasks
   occupy different phases; the aggregate Run state must not overwrite it.
6. MUST keep review Run rendering, keyboard/mouse behavior, detail panes,
   spatial layout, terminal restoration, and public Attach command behavior
   unchanged.
7. MUST communicate phase through text and structure rather than color alone,
   preserve styled/no-color text parity, and prevent longer labels from
   wrapping or moving panel boundaries at supported widths.

## Subtasks

- [ ] Extend Live Run View capacity data and attach replay.
- [ ] Replace the generic spec capacity header with both canonical labels.
- [ ] Fold Task and Verification journal events into a per-Task phase model.
- [ ] Render accessible working, waiting, verifying, and terminal labels.
- [ ] Protect review Runs, legacy attach, narrow layout, and no-color behavior.
- [ ] Drive Bubble Tea v2 updates synchronously in focused tests.

## Acceptance Criteria

- [ ] A spec Run header renders `Task Capacity: 2` and
      `Verification Capacity: 1`; a review Run renders neither.
- [ ] Attach replays new event fields, and a legacy Task-cycle-start event
      falls back without crashing or inventing unknown capacity.
- [ ] Two interleaved Tasks can simultaneously render `Agent working` and
      `Verifying`, or `Waiting for Verification` and `Verifying`, independent
      of aggregate Run state.
- [ ] Verification Feedback returns only the affected Task to `Agent working`;
      its next waiting and started events advance it deterministically.
- [ ] Terminal settlement cannot regress to an earlier phase when later stale
      or replayed events are encountered.
- [ ] Styled and `NO_COLOR` views contain the same status words, and phase
      meaning remains readable without ANSI sequences.
- [ ] Existing review cockpit, keybinding, detail, and viewport tests remain
      unchanged or pass with behavior-equivalent expectations.

## Context

- instruction: `.agents/skills/bubbletea/SKILL.md`
- instruction: `.agents/skills/tui-design/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/cli/attach.go`
- interface: `internal/cli/attach_test.go`
- interface: `internal/cli/runbrowser.go`
- interface: `internal/tui/tui.go`
- interface: `internal/tui/tui_test.go`
- interface: `internal/tui/cockpit.go`
- interface: `internal/tui/cockpit_test.go`
- interface: `internal/tui/cockpit_fidelity_test.go`
- interface: `internal/runevent/event.go`

## Verification

- `rtk go test ./internal/cli -run 'TestAttach.*(Capacit|Legacy|Verification)' -count=1` — expected: attached spec Runs replay both capacities and legacy events degrade deterministically.
- `rtk go test ./internal/tui -run 'Test(RenderLiveRunView.*Capacit|Cockpit.*Task.*(AgentWorking|WaitingForVerification|Verifying|Settlement|NoColor))' -count=1` — expected: exact accessible labels track interleaved per-Task events and terminal precedence.
- `rtk go test ./internal/tui -run 'Test(Cockpit.*Review|Cockpit.*Key|Cockpit.*Detail|Cockpit.*Viewport|RenderLiveRunView.*Review)' -count=1` — expected: review layout and interactions retain existing behavior.
- `rtk go test -race ./internal/cli ./internal/tui -run 'Test(Attach.*(Capacit|Legacy|Verification)|RenderLiveRunView.*Capacit|Cockpit.*Task.*(AgentWorking|WaitingForVerification|Verifying|Settlement|NoColor))' -count=1` — expected: attach replay and synchronous Bubble Tea phase projection are race-free.

## References

- `_prd.md` → Goal 4; User Story 3; Core Feature 8; User Experience; Success Metrics.
- `_techspec.md` → System Architecture; Implementation Design: Data Models and API Contracts; Integration Points; Testing Approach; Build Orders 5–6.
- `../../adr/0056-spec-runs-separate-task-and-verification-capacity.md` → durable capacity and per-Task observability decision.
