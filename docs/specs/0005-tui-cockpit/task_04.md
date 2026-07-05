---
task: task_04
spec: 0005-tui-cockpit
status: pending
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

- [ ] Batch header grouping over journal order
- [ ] Kind-section rendering from existing summaries
- [ ] Follow Mode guard tests after regrouping
- [ ] Spec-Run timeline fixtures

## Acceptance Criteria

- [ ] A review journal fixture renders grouped batches with plan/tool/think
      sections matching the mockup's structure.
- [ ] A spec journal fixture groups task and qa milestones legibly.
- [ ] Follow Mode tests pass unchanged plus new ones proving tail behavior
      across group boundaries.
- [ ] Full suite passes.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 5; Core Feature 4. `_techspec.md` → Build Order
4, Decisions (render-time only). `design/ui-redesign-plan.md` → Required
Changes (timeline grouping); `design/roundfix-01.png`. ADR-0008.
