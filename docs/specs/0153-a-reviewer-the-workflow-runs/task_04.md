---
task: task_04
spec: 0153-a-reviewer-the-workflow-runs
status: pending
type: docs
complexity: low
---

# Task 04: Describe the command in the shipped skill and the guide

## Overview

A new command is public CLI surface, and this repository's hard rule requires
the skill update to ship with the pull request that changes it.

## Requirements

1. MUST state that the configured Codex pre-PR reviewer is run by the workflow
   over the current candidate, that the command hands it the candidate diff,
   and that its record names repository, base and head.
2. MUST state that explicit `none` performs no reviewer call and no readiness
   probe and records a configured omission.
3. MUST state that a runtime failure, a timeout, a transport anomaly, empty
   output and unclassifiable output each block the selected mode and never
   become a pass or an omission.
4. MUST state that `claude` and `coderabbit` are valid policy values this
   command refuses to execute for now.
5. MUST regenerate the distributed mirror with `make skills-sync` rather than
   editing it.
6. MUST document the command in the user guide beside the commands it sits with,
   and correct any statement there that says the policy is reported without
   being enforced.

## Subtasks

- [ ] Edit the canonical skill.
- [ ] Regenerate the mirror with the sanctioned command.
- [ ] Document the command in the user guide.

## Acceptance Criteria

- [ ] Both skill files describe the command, its exits and everything that
      blocks.
- [ ] The user guide documents the command, its option and its exit codes.
- [ ] The repository's skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "roundfix review" .agents/skills/roundfix/SKILL.md && grep -q "roundfix review" skills/roundfix/SKILL.md && grep -q "roundfix review" docs/user-guide/commands.md && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task none of the three names the command, so the command fails.

## References

- [_authorization.md](_authorization.md) — Approved bounded mutation
