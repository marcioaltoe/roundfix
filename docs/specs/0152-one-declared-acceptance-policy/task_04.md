---
task: task_04
spec: 0152-one-declared-acceptance-policy
status: pending
type: docs
complexity: low
---

# Task 04: State the one policy in the shipped skill

## Overview

The derived Verification command is public contract surface the skill
describes, and this Spec changes what it accepts.

## Requirements

1. MUST state, in the canonical skill, that one eligibility policy decides
   whether the newest QA Report is acceptable, and that settlement and archive
   both apply it.
2. MUST state that a `partial` verdict whose blocked rows are declared
   unreachable settles the terminal `qa` Task.
3. MUST state that `fail`, an undeclared partial, a missing or unparseable
   report, and a `pass` carrying blocked rows each still refuse.
4. MUST regenerate the distributed mirror with `make skills-sync` rather than
   editing it.
5. MUST NOT change the skill's version or any behavior it documents beyond this
   Spec's own.

## Subtasks

- [ ] Edit the canonical skill.
- [ ] Regenerate the mirror with the sanctioned command.

## Acceptance Criteria

- [ ] Both skill files state the one policy and what settles under it.
- [ ] The repository's skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "one declared-acceptance eligibility policy" .agents/skills/roundfix/SKILL.md && grep -q "one declared-acceptance eligibility policy" skills/roundfix/SKILL.md && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task neither skill file carries the phrase, so the command fails.

## References

- [_authorization.md](_authorization.md) — Approved bounded mutation
