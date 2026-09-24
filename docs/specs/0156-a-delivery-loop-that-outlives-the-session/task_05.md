---
task: task_05
spec: 0156-a-delivery-loop-that-outlives-the-session
status: pending
type: docs
complexity: low
---

# Task 05: Describe the queue in the shipped skill and the guide

## Overview

A new command family is public CLI surface, and the repository's hard rule requires the skill update to ship with it.

## Requirements

1. MUST state that `roundfix deliver` advances a queue from Run to merge in the order Run, pre-PR review, archive on the branch, repository gate, pull request, checks, squash merge.
2. MUST state that a blocker parks its item, that a resumed queue reconciles each recorded action before retrying it, and that publication needs push, pull_request and merge in the Spec's authorization.
3. MUST regenerate the distributed mirror with `make skills-sync` rather than editing it.
4. MUST document the command family in the user guide.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Both skill files describe the command family.
- [ ] The user guide documents it.
- [ ] The skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "roundfix deliver start" .agents/skills/roundfix/SKILL.md && grep -q "roundfix deliver start" skills/roundfix/SKILL.md && grep -q "roundfix deliver" docs/user-guide/commands.md && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task none of the three names the command family, so the command fails.

## References

- [_techspec.md](_techspec.md) — Components
