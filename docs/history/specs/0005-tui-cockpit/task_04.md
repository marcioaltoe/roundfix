---
task: task_04
spec: 0005-tui-cockpit
status: completed
type: frontend
complexity: medium
---

# Task 04: Group the timeline by Batch and event kind

## Overview

Make the timeline read like the mockup: events grouped under Batch headers
with collapsible-looking sections per kind — plan, tool (with state), think,
status, daemon milestones — purely at render time, with Follow Mode behavior
guarded by tests. Verifiable through timeline rendering tests over journal
fixtures.

## Requirements

1. MUST group timeline entries under Batch headers (`BATCH nnn/mmm <state>
   <elapsed>`) in journal order.
2. MUST render kind sections within a Batch — plan, tool (name + state +
   command/summary lines), think, session/status, daemon milestones
   (verification, commit, qa, outcome) — from existing event summaries and
   raw payload conversion; no event mutation, no new Run Event vocabulary,
   unknown kinds skippable as today (ADR-0008).
3. MUST preserve Follow Mode semantics exactly: tail advances on new events,
   suspends on scrollback, resumes at bottom — guarded by the existing
   viewport tests plus regrouping-specific ones.
4. MUST render spec-Run timelines (daemon.task, daemon.qa) meaningfully in
   the same grouping.

## Subtasks

- [x] Batch header grouping over journal order
- [x] Kind-section rendering from existing summaries
- [x] Follow Mode guard tests after regrouping
- [x] Spec-Run timeline fixtures

## Acceptance Criteria

- [x] A review journal fixture renders grouped batches with plan/tool/think
      sections matching the mockup's structure.
- [x] A spec journal fixture groups task and qa milestones legibly.
- [x] Follow Mode tests pass unchanged plus new ones proving tail behavior
      across group boundaries.
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 5; Core Feature 4. `_techspec.md` → Build Order
4, Decisions (render-time only). `design/ui-redesign-plan.md` → Required
Changes (timeline grouping); `design/roundfix-01.png`. ADR-0008.

## Result

- Review grouping evidence: `TestViewportGroupsReviewTimelineByBatchAndKind`
  asserts `BATCH 001/002 executing 00:38`, plan text, tool marker and command,
  think text, session status, and a second Batch header with daemon Batch
  summaries folded into headers.
- Spec grouping evidence: `TestViewportGroupsSpecTimelineTaskAndQAMilestones`
  asserts `TASK`, `VERIFY`, and `QA` sections under Batch headers for
  `daemon.task`, `daemon.verification`, and `daemon.qa` events.
- Follow Mode evidence: existing viewport follow/scroll tests pass unchanged,
  and `TestViewportFollowModeTracksGroupedBatchBoundaries` proves scrollback
  remains frozen on new events and `End` resumes at the new Batch tail.
- Raw-payload/coalescing evidence: agent plan/tool/thought/status rendering
  uses the existing raw payload conversion, and
  `TestViewportCoalescesChunksInsideGroupedBatch` guards chunk coalescing
  inside grouped Batches.
- Verification: `rtk go test ./internal/tui/` passed with 74 tests.
- Verification: `rtk go test ./...` passed with 554 tests across 16 packages.
