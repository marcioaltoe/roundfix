---
task: task_03
spec: 0151-a-supersession-the-archive-can-see
status: completed
type: docs
complexity: low
---

# Task 03: Describe both in the shipped skill and the guide

## Overview

A new command and a changed archive precondition are public CLI behavior. The
repository's hard rule requires the skill update to ship with the pull request
that changes it.

## Requirements

1. MUST state, in the canonical skill, that a Spec whose content another Spec
   delivered records a supersession rather than editing its own status or
   writing a note.
2. MUST state that archive accepts a recorded supersession in place of completed
   Tasks and a passing gate, and keeps every other precondition.
3. MUST state that an unknown superseding Spec, a self-supersession and a second
   supersession are refused.
4. MUST regenerate the distributed mirror with `make skills-sync` rather than
   editing it.
5. MUST document the command in the user guide beside the commands it sits with.
6. MUST NOT change the skill's version or any behavior it documents beyond this
   Spec's own.

## Subtasks

- [ ] Edit the canonical skill.
- [ ] Regenerate the mirror with the sanctioned command.
- [ ] Document the command in the user guide.

## Acceptance Criteria

- [ ] Both skill files describe the command and the changed archive precondition.
- [ ] The user guide documents the command, its options and its exit codes.
- [ ] The repository's skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "roundfix supersede --spec" .agents/skills/roundfix/SKILL.md && grep -q "roundfix supersede --spec" skills/roundfix/SKILL.md && grep -q "roundfix supersede --spec" docs/user-guide/commands.md && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task none of the three carries the invocation, so the command fails. The bare word "supersede" is not an anchor: it already appears as "superseded" in the reconcile vocabulary.

## References

- [_authorization.md](_authorization.md) — Approved bounded mutation

## Result

Implemented the documentation slice:

- The canonical Roundfix skill now documents `roundfix supersede`, including
  the amendment-vs-status/note rule, options, exit codes, refusal cases, and
  archive's recorded-supersession proof with all other preconditions retained.
- The distributed `skills/roundfix/SKILL.md` mirror was regenerated with
  `make skills-sync`.
- The user guide documents the command beside the Spec lifecycle commands,
  including its options, exit codes, refusal cases, and archive interaction.

Focused checks after the edits:

- `make skills-sync` — passed.
- `cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — passed;
  canonical and distributed skill files match.
- `git diff --check` — passed.

The declared Verification command was not run; the Daemon owns that gate.
