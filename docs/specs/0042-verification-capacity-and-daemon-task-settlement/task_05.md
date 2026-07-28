---
task: task_05
spec: 0042-verification-capacity-and-daemon-task-settlement
status: completed
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

- [x] Extend Live Run View capacity data and attach replay.
- [x] Replace the generic spec capacity header with both canonical labels.
- [x] Fold Task and Verification journal events into a per-Task phase model.
- [x] Render accessible working, waiting, verifying, and terminal labels.
- [x] Protect review Runs, legacy attach, narrow layout, and no-color behavior.
- [x] Drive Bubble Tea v2 updates synchronously in focused tests.

## Acceptance Criteria

- [x] A spec Run header renders `Task Capacity: 2` and
      `Verification Capacity: 1`; a review Run renders neither.
- [x] Attach replays new event fields, and a legacy Task-cycle-start event
      falls back without crashing or inventing unknown capacity.
- [x] Two interleaved Tasks can simultaneously render `Agent working` and
      `Verifying`, or `Waiting for Verification` and `Verifying`, independent
      of aggregate Run state.
- [x] Verification Feedback returns only the affected Task to `Agent working`;
      its next waiting and started events advance it deterministically.
- [x] Terminal settlement cannot regress to an earlier phase when later stale
      or replayed events are encountered.
- [x] Styled and `NO_COLOR` views contain the same status words, and phase
      meaning remains readable without ANSI sequences.
- [x] Existing review cockpit, keybinding, detail, and viewport tests remain
      unchanged or pass with behavior-equivalent expectations.

## Context

- instruction: `.agents/skills/bubbletea/SKILL.md`
- instruction: `.agents/skills/tui-design/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/cli/attach.go`
- interface: `internal/cli/cli_test.go`
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

## Result

### What changed

- `internal/tui/tui.go` — `LiveRunView` keeps `Concurrency` as the effective
  Task Capacity and adds `VerificationConcurrency`; the plain header renders
  `Task Capacity:` and `Verification Capacity:` for spec Runs only, replacing
  the generic `Concurrency:` line. `WorkItem` gains `Phase`, the row slot
  Review Issues spend on `Severity`.
- `internal/tui/cockpit.go` — per-Task phase model folded from the journal:
  `taskPhase` / `taskJournalState`, `applyDaemonTaskEvent`,
  `applyDaemonVerificationEvent`, and `advanceTaskPhase`. `taskStatusLabel`
  now returns the exact labels `Agent working`, `Waiting for Verification`,
  `Verifying`, `Completed`, `Failed`, `Skipped`, `Paused`, `Waiting`; new
  `[queued]` / `[verify]` markers join `[run]`; the Run header shows both
  canonical capacity labels.
- `internal/cli/attach.go`, `internal/cli/runbrowser.go` —
  `attachRunConcurrency` became `attachRunCapacities`, reading
  `task_capacity` / `verification_capacity` from the Task-cycle-start
  `daemon.status` event, with the legacy `concurrency` alias as Task Capacity
  and the configured `verification.concurrency` standing in for a Verification
  Capacity the Run never recorded.
- `internal/cli/implement.go` — the owning Live Run View carries the
  configured Verification Capacity.
- `internal/daemon/task_engine.go` — `repairTaskVerification` publishes the
  `daemon.task` phase `verification_feedback` the TechSpec's Data Models
  section specifies. No prior Task emitted it, and without it the Live Run
  View has no per-Task signal for the Verification Feedback turn; it is the
  producer half of Requirement 4 and is covered by a daemon-side ordering
  assertion.

### Acceptance criteria

1. **Capacity header** — `TestRenderLiveRunViewSpecRunShowsTaskAndVerificationCapacity`
   asserts a spec Run renders `Task Capacity: 2` and `Verification Capacity: 1`,
   a review Run renders neither label (nor `Concurrency`), and an unresolved
   capacity is omitted rather than invented.
   `TestCockpitSpecRunHeaderShowsTaskAndVerificationCapacity` and
   `TestCockpitReviewRunHeaderShowsNoCapacity` assert the same for the
   interactive header.
2. **Attach replay and legacy fallback** —
   `TestAttachSpecRunReplaysTaskAndVerificationCapacity` replays
   `{"concurrency":3,"task_capacity":3,"verification_capacity":2}` and renders
   `Task Capacity: 3` / `Verification Capacity: 2`.
   `TestAttachSpecRunLegacyCapacityEventFallsBackDeterministically` attaches a
   Run whose journal has only `{"concurrency":4}`: exit 0, `Task Capacity: 4`,
   and `Verification Capacity: 1` from the configured built-in default.
3. **Interleaved phases** —
   `TestCockpitSpecRunTaskPhasesAgentWorkingWaitingForVerificationAndVerifying`
   runs with the aggregate Run state pinned to `Verifying` and asserts
   `[run] Agent working` on task_01 beside `[verify] Verifying` on task_02,
   then `[queued] Waiting for Verification` beside `[verify] Verifying`, then
   `[done] Completed` after settlement — proving the aggregate state never
   overwrites per-Task evidence.
4. **Verification Feedback** —
   `TestCockpitSpecRunTaskVerificationFeedbackReturnsOneTaskToAgentWorking`
   asserts the `verification_feedback` event returns only task_01 to
   `Agent working` while task_02 stays `Waiting for Verification`, and that
   task_01's attempt-2 `waiting` then `started` events advance it to
   `Waiting for Verification` then `Verifying`.
   `TestTaskCycleRepairReacquiresVerificationCapacityAfterFeedback` now also
   asserts the Daemon publishes exactly one `verification_feedback` event, for
   the repaired Task only, ordered after its failed verdict and before its
   next waiting event.
5. **Terminal precedence** —
   `TestCockpitSpecRunTaskSettlementResistsStaleAndReplayedEvents` replays
   `started` and verification `started` events *after* settlement and skip, and
   asserts the rows stay `[fail] Failed` and `[skip] Skipped`.
6. **Styled / NO_COLOR parity** —
   `TestCockpitSpecRunTaskPhaseLabelsReadTheSameUnderNoColor` asserts the two
   renders are identical after stripping ANSI and that each phase label appears
   verbatim in *both* raw renders, so no label is split by escape sequences.
7. **Existing behavior** — the third declared command (14 tests over review
   cockpit, keys, detail, viewport) passes unchanged. The two review snapshot
   goldens are byte-identical; only the two `spec_run_pane_*` goldens changed,
   adding the phase word inside the existing row (`[run]` → `[run] Agent
   working`) with no line, panel, or boundary change.
   `TestCockpitSpecRunTaskVerifyingLabelsKeepNarrowLayoutStable` compares a
   phase-carrying render against a baseline at 88/100/120 columns and fails on
   any changed line count or line width.

### Verification evidence

- `rtk go test ./internal/cli -run 'TestAttach.*(Capacit|Legacy|Verification)' -count=1` → `Go test: 2 passed in 1 packages`
- `rtk go test ./internal/tui -run 'Test(RenderLiveRunView.*Capacit|Cockpit.*Task.*(AgentWorking|WaitingForVerification|Verifying|Settlement|NoColor))' -count=1` → `Go test: 6 passed in 1 packages`
- `rtk go test ./internal/tui -run 'Test(Cockpit.*Review|Cockpit.*Key|Cockpit.*Detail|Cockpit.*Viewport|RenderLiveRunView.*Review)' -count=1` → `Go test: 14 passed in 1 packages`
- `rtk go test -race ./internal/cli ./internal/tui -run 'Test(Attach.*(Capacit|Legacy|Verification)|RenderLiveRunView.*Capacit|Cockpit.*Task.*(AgentWorking|WaitingForVerification|Verifying|Settlement|NoColor))' -count=1` → `Go test: 8 passed in 2 packages`
- `rtk go test -race ./... -count=1` → `Go test: 2721 passed in 23 packages`
- `make verify` → fmt-check, full test run, skills-sync-check, skills-check, and build all passed.

### Notes

- At the minimum two-pane width (88 columns) the longest label truncates
  inside the card — `[queued] Waiting for Veri…` — which is the truncation the
  TechSpec's Risks section sanctions. The `[queued]` / `[verify]` / `[run]`
  markers never truncate, so the phase stays distinguishable at every
  supported width, and the label reads in full from ~100 columns up. No
  panel was resized and no card gained a line.
- Follow-up for Task 07: the `verification_feedback` phase is now part of the
  public `roundfix-events/v1` task-status projection (it flows through the
  existing required `phase` field); operator docs listing Task phases should
  mention it.
