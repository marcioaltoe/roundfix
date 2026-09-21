---
task: task_04
spec: 0152-one-declared-acceptance-policy
status: completed
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

## Result

Implemented the shipped-skill wording for the one declared-acceptance policy.
The canonical skill now states that settlement and archive use the same policy,
that a qualifying `partial` settles the terminal `qa` Task, and that `fail`,
undeclared `partial`, missing or unparseable reports, and disallowed blocked
rows on `pass` refuse. The existing environment-blocked-row exception remains
documented. The distributed mirror was regenerated from the canonical skill
with the sanctioned command; the skill version was unchanged.

Focused checks:

- `make skills-sync` — passed.
- `cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — passed;
  the mirror is identical to the canonical file.
- Focused text checks for the policy, qualifying `partial`, and refusal wording
  in both skill files — passed.
- `git diff --check` — passed.

Acceptance evidence:

- Both skill files contain the one declared-acceptance policy and settlement /
  archive behavior, with the focused mirror and text checks passing.
- The repository skill check remains for Daemon Verification; it was not run
  here because it is the Task's declared Verification command.
