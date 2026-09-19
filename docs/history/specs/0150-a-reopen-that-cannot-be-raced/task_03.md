---
task: task_03
spec: 0150-a-reopen-that-cannot-be-raced
status: completed
type: docs
complexity: low
---

# Task 03: Name the recheck refusal in the shipped skill

## Overview

The recheck adds an exit-2 condition. Exit codes are public API here, and the
repository's hard rule requires the skill update to ship with the pull request
that changes CLI behavior.

## Requirements

1. MUST state, in the canonical skill, that the reopen command rechecks the
   gate's staleness immediately before it writes and refuses when the gate
   changed since preflight, using that exact phrase so the check has an anchor.
2. MUST regenerate the distributed mirror with `make skills-sync` rather than
   editing it.
3. MUST NOT change the skill's version or any behavior it documents beyond this
   Spec's own.

## Subtasks

- [ ] Edit the canonical skill.
- [ ] Regenerate the mirror with the sanctioned command.

## Acceptance Criteria

- [ ] Both skill files name the recheck refusal.
- [ ] The repository's skill check passes.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `skills/roundfix/SKILL.md`

## Verification

- `grep -q "gate changed since preflight" .agents/skills/roundfix/SKILL.md && grep -q "gate changed since preflight" skills/roundfix/SKILL.md && go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0; before this Task neither skill file carries this phrase, so the command fails. The word "recheck" alone already appears in an unrelated paragraph, so it is not an anchor.

## References

- [_authorization.md](_authorization.md) — Approved bounded mutation

## Result

Implemented the bounded skill documentation slice. The canonical Roundfix
skill now states that reopen rechecks the gate's staleness immediately before
writing and refuses with exit 2 when the gate changed since preflight. Ran
`make skills-sync` to regenerate the distributed mirror; no skill version or
unrelated behavior text was changed.

Focused checks:

- `rg -n -F "gate changed since preflight" .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — found the required anchor in both files.
- `diff -u .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — canonical and mirror are identical.
- `git diff --check` — passed.

The declared repository skill verification command remains for Daemon-owned
verification and was not run in this handoff.
