---
task: task_04
spec: 0168-deliver-one-worktree-per-item
status: pending
type: docs
complexity: low
---

# Task 04: Describe item worktrees in the shipped skill and the guide

## Overview

Tasks 02 and 03 change what `roundfix deliver` does to the checkout and what `deliver status` prints, and the repository's hard rule requires the shipped skill to stay true to CLI behaviour.

## Requirements

1. MUST state, in the skill's delivery queue section and in the user guide's `deliver` section, that each queued Spec runs in its own linked worktree under `worktree.location`, created from the refreshed default branch.
2. MUST state in both that `roundfix deliver` never switches, resets or cleans your checkout, and does not need it clean.
3. MUST state in both that a parked item keeps its worktree, that `deliver status` prints it, that a resume recreates a missing worktree from its branch or parks the item as `item-worktree-missing`, and that a merged item's worktree and local branch are removed.
4. MUST regenerate the distributed mirror with `make skills-sync` rather than editing it.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Both skill files and the user guide describe item worktrees and the untouched checkout.
- [ ] The mirror is identical to the canonical skill.
- [ ] The skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `grep -q "runs in its own linked worktree" .agents/skills/roundfix/SKILL.md && grep -q "never switches, resets or cleans your checkout" .agents/skills/roundfix/SKILL.md && grep -q "item-worktree-missing" .agents/skills/roundfix/SKILL.md && grep -q "runs in its own linked worktree" docs/user-guide/commands.md && grep -q "never switches, resets or cleans your checkout" docs/user-guide/commands.md && grep -q "item-worktree-missing" docs/user-guide/commands.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task neither the skill nor the guide describes item worktrees, so the command fails.

## References

- [_techspec.md](_techspec.md) — API Contracts
