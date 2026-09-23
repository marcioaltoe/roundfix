---
task: task_03
spec: 0154-evidence-when-the-target-branch-is-gone
status: pending
type: docs
complexity: low
---

# Task 03: Describe the fallback in the shipped skill

## Overview

`reconcile` reporting is public CLI behavior the skill documents, and this Spec
changes what it can prove.

## Requirements

1. MUST state that a Run whose target branch is absent is checked for content
   evidence against the default branch, using the rule already accepted for a
   missed ancestry.
2. MUST state that a Run is released only on positive content evidence, and that
   an unreachable default branch, an unarchived Spec, or evidence naming another
   Spec each leave it preserved.
3. MUST state that the preserved reason names the proof that was missing.
4. MUST regenerate the distributed mirror with `make skills-sync` rather than
   editing it.
5. MUST NOT change the skill's version or any behavior it documents beyond this
   Spec's own.

## Subtasks

- [ ] Edit the canonical skill.
- [ ] Regenerate the mirror with the sanctioned command.

## Acceptance Criteria

- [ ] Both skill files describe the fallback and what still preserves.
- [ ] The repository's skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "default branch" .agents/skills/roundfix/SKILL.md && grep -q "absent target" .agents/skills/roundfix/SKILL.md && grep -q "absent target" skills/roundfix/SKILL.md && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task neither skill file carries the phrase "absent target", so the command fails.

## References

- [_authorization.md](_authorization.md) — Approved bounded mutation
