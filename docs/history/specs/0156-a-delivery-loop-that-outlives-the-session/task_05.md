---
task: task_05
spec: 0156-a-delivery-loop-that-outlives-the-session
status: completed
type: docs
complexity: low
---

# Task 05: Describe the queue in the shipped skill and the guide

## Overview

A new command family is public CLI surface, and the repository's hard rule requires the skill update to ship with it.

## Requirements

1. MUST state that `roundfix deliver` advances a queue from Run to merge in the order Run, pre-PR review, archive on the branch, repository gate, pull request, checks, squash merge.
2. MUST state that a blocker parks its item, that a resumed queue reconciles each recorded action before retrying it, and that publication needs push, pull_request and merge in the Spec's authorization.
3. MUST regenerate the distributed mirror with `make skills-sync` rather than editing it.
4. MUST document the command family in the user guide.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Both skill files describe the command family.
- [ ] The user guide documents it.
- [ ] The skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "roundfix deliver start" .agents/skills/roundfix/SKILL.md && grep -q "roundfix deliver start" skills/roundfix/SKILL.md && grep -q "roundfix deliver" docs/user-guide/commands.md && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task none of the three names the command family, so the command fails.

## References

- [_techspec.md](_techspec.md) — Components

## Result

Updated the canonical Roundfix skill with the `roundfix deliver` command
family, its Run-to-merge order, blocker parking, resume reconciliation, and
the required `push`, `pull_request`, and `merge` authorization. Documented the
same command family and behavior in the user command guide. Regenerated
`skills/roundfix/SKILL.md` from `.agents/skills/roundfix/SKILL.md` with the
sanctioned `make skills-sync` target; the mirror was not hand-edited.

Focused checks:

- `make skills-sync` passed.
- `git diff --check` passed.
- `cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` passed;
  the canonical skill and distributed mirror match byte-for-byte.
- Manual diff inspection confirmed the guide contains `roundfix deliver start`
  and the skill sections describe the ordered stages, parking, reconciliation,
  and publication authorization.
- The Task's declared Verification commands were not run; the Daemon owns
  those checks and final settlement.

Acceptance evidence:

- Both skill files describe the command family: the canonical section and its
  synced mirror contain the four `deliver` commands and queue rules.
- The user guide documents it in the `deliver` command-reference section.
- Skill-check evidence is pending the Daemon's declared Verification; no
  terminal pass or Task status was claimed in this handoff.

## Carry-forward provenance

- Source Run: `run_20260924T121521Z_744099e82ec81b88`
- Source commit: `ea8fd7a54b289cc76faf51a8515849d1b7c20e17`
