---
task: task_08
spec: 0160-a-review-that-reaches-a-verdict
status: pending
type: docs
complexity: low
---

# Task 08: Describe the Spec-context contract the correctives delivered

## Overview

QA finding F-001 of 2026-09-24: Tasks 06 and 07 changed what `roundfix review` records and sends, but the shipped Roundfix skill, its mirror and `docs/user-guide/commands.md` still describe only active changed Specs under the Spec Root.

## Requirements

1. MUST state in `.agents/skills/roundfix/SKILL.md` and `docs/user-guide/commands.md` that Specs archived within the candidate are read from the archive root at `HEAD`.
2. MUST state that a changed Spec without a `## Decisions` section, a PRD or a TechSpec is skipped and listed in the record's `skippedSpecs`, and the review proceeds.
3. MUST state that carried Spec context is bounded per Spec and in total, that truncation is marked in the prompt, and that the record's `specContextTruncated` says so.
4. MUST state that `answerPath` is set only when the prompt reached a reviewer.
5. MUST regenerate the distributed mirror with `make skills-sync` rather than editing it.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The skill, its mirror and the guide describe archived discovery, skipped Specs, bounded context and when `answerPath` is set.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `grep -q "skippedSpecs" .agents/skills/roundfix/SKILL.md && grep -q "specContextTruncated" .agents/skills/roundfix/SKILL.md && grep -q "archive root" .agents/skills/roundfix/SKILL.md && grep -q "skippedSpecs" skills/roundfix/SKILL.md && grep -q "specContextTruncated" skills/roundfix/SKILL.md && grep -q "archive root" skills/roundfix/SKILL.md && grep -q "skippedSpecs" docs/user-guide/commands.md && grep -q "specContextTruncated" docs/user-guide/commands.md && grep -q "archive root" docs/user-guide/commands.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task none of the three surfaces names `skippedSpecs`, so the command fails.

## References

- [_techspec.md](_techspec.md) — Spec-aware prompt
