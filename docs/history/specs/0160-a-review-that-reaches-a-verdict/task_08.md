---
task: task_08
spec: 0160-a-review-that-reaches-a-verdict
status: completed
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

## Result

- Updated the canonical Roundfix skill and user guide to describe active and
  archived Spec discovery, including reading Specs under the archive root at
  `HEAD`.
- Documented that changed Specs without `## Decisions`, a PRD, or a TechSpec
  are skipped into `skippedSpecs` and do not block the review; documented the
  32 KiB per-Spec and 64 KiB total bounds, the prompt truncation marker, and
  `specContextTruncated`.
- Documented that `answerPath` is set only after the prompt reaches a reviewer.
- Regenerated the distributed mirror with `make skills-sync`.
- Focused checks after the final edit: `git diff --check` passed; targeted
  documentation searches confirmed the required contract terms are present in
  the canonical skill, mirror, and user guide. The task's declared
  Verification commands were not run; Daemon Verification remains the
  settlement authority.
