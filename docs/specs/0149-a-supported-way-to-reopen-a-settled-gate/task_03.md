---
task: task_03
spec: 0149-a-supported-way-to-reopen-a-settled-gate
status: pending
type: docs
complexity: low
---

# Task 03: Keep the shipped skill and the guide true

## Overview

The Roundfix skill is the shipped description of the CLI surface. A new command
the skill does not name is a skill that is no longer true, and this repository's
hard rule requires the skill update to ship with the pull request that changes
CLI behavior.

## Requirements

1. MUST state, in the canonical skill, that a settled QA gate whose dependencies
   are no longer completed is cleared with the reopen command, never by editing
   the Task file.
2. MUST state that the command refuses when the gate is not settled or not
   stale, and that the prior QA Report and Result are retained.
3. MUST regenerate the distributed mirror with `make skills-sync` rather than
   editing it.
4. MUST document the command in the user guide beside the commands it sits with.
5. MUST NOT change the skill's version or any behavior it documents beyond this
   Spec's own.

## Subtasks

- [ ] Edit the canonical skill.
- [ ] Regenerate the mirror with the sanctioned command.
- [ ] Document the command in the user guide.

## Acceptance Criteria

- [ ] Both skill files name the command and its refusals.
- [ ] The user guide documents the command, its options and its exit codes.
- [ ] The repository's skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "reopen" .agents/skills/roundfix/SKILL.md && grep -q "reopen" skills/roundfix/SKILL.md && grep -q "reopen" docs/user-guide/commands.md && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task no skill file or guide names the command, so the command fails.

## References

- [_authorization.md](_authorization.md) — Approved bounded mutation
