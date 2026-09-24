---
task: task_06
spec: 0163-baseline-decisions-and-regeneration
status: pending
type: docs
complexity: low
---

# Task 06: Shipped skills describe the changes

## Overview

The repository's hard rule requires the shipped skills to match the CLI behaviour this Spec changes: the HTTP mode change, the Greenfield refusal and the new reconciliation command.

## Requirements

1. MUST add `roundfix baseline skills reconcile` to the explicit maintenance operations of `.agents/skills/roundfix/SKILL.md`, stating that it removes only lock entries absent at the selected immutable commit.
2. MUST state in `.agents/skills/setup-context-driven/SKILL.md` the sentence `Changing only the HTTP mode retains every exception and the source.`
3. MUST state in `.agents/skills/setup-context-driven/SKILL.md` that Greenfield refuses when retained managed source needs classification, and names `preservation.mode=preservation`.
4. MUST regenerate both mirrors with `make skills-sync` rather than editing them.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Both copies of each skill carry the new text.
- [ ] Each mirror equals its canonical copy.
- [ ] The skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`
- instruction: `.agents/skills/setup-context-driven/SKILL.md`
- instruction: `skills/setup-context-driven/SKILL.md`

## Verification

- `grep -q "roundfix baseline skills reconcile" .agents/skills/roundfix/SKILL.md && grep -q "Changing only the HTTP mode retains every exception and the source." .agents/skills/setup-context-driven/SKILL.md && grep -q "preservation.mode=preservation" .agents/skills/setup-context-driven/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix && diff -r .agents/skills/setup-context-driven skills/setup-context-driven && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task neither skill carries the new text, so the command fails.

## References

- [_techspec.md](_techspec.md) — Testing Approach 5
