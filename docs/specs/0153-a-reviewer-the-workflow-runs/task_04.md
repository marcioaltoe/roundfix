---
task: task_04
spec: 0153-a-reviewer-the-workflow-runs
status: completed
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

- [x] Edit the canonical skill.
- [x] Regenerate the mirror with the sanctioned command.
- [x] Document the command in the user guide.

## Acceptance Criteria

- [x] Both skill files describe the command, its exits and everything that
      blocks.
- [x] The user guide documents the command, its option and its exit codes.
- [ ] The repository's skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "roundfix review" .agents/skills/roundfix/SKILL.md && grep -q "roundfix review" skills/roundfix/SKILL.md && grep -q "roundfix review" docs/user-guide/commands.md && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task none of the three names the command, so the command fails.

## References

- [_authorization.md](_authorization.md) — Approved bounded mutation

## Result

Implemented the public documentation contract for `roundfix review`: the
workflow hands the current candidate diff to the configured Codex reviewer,
records repository/base/head identity, treats explicit `none` as a configured
omission without a call or readiness probe, and records every runtime,
timeout, transport, empty-output, and unclassifiable-output case as blocked.
The documentation also states that `claude` and `coderabbit` are valid policy
values that this command refuses to execute for now, and documents the
command's `--base` option and exit codes. The older user-guide wording that
described the policy as reported but unenforced was corrected.

Focused checks:

- `make skills-sync` — passed; the distributed mirror was regenerated.
- `cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — passed;
  canonical and mirror are identical.
- `rg` checks for `roundfix review` in both skill files and
  `docs/user-guide/commands.md` — passed.
- Search for stale “reported without enforcing” policy wording in the user
  guide and both skill files — no matches.
- `git diff --check` — passed.

The task's declared Verification command, including the repository skill
check, remains for the Daemon.
