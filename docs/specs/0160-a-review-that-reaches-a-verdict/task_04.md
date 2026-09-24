---
task: task_04
spec: 0160-a-review-that-reaches-a-verdict
status: pending
type: docs
complexity: low
---

# Task 04: Keep the shipped skill and the guide true

## Overview

The Roundfix skill says `claude` is refused and that only an exact `No findings` passes; both become false with this Spec.

## Requirements

1. MUST state in `.agents/skills/roundfix/SKILL.md` that `claude` runs, that the verdict is classified by substance, that the raw answer is kept at `answerPath`, and that a delivery's Spec decisions are reviewed.
2. MUST keep `coderabbit` described as refused.
3. MUST describe the same in `docs/user-guide/commands.md`.
4. MUST regenerate the distributed mirror with `make skills-sync` rather than editing it.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Both skill files and the guide describe the new behavior.
- [ ] The skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `grep -q "answerPath" .agents/skills/roundfix/SKILL.md && grep -q "answerPath" skills/roundfix/SKILL.md && grep -q "answerPath" docs/user-guide/commands.md && ! grep -q "refuses to execute them for now" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
