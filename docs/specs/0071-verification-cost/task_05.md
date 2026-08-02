---
task: task_05
spec: 0071-verification-cost
status: pending
type: docs
complexity: medium
---

# Task 05: Stop charging every Task for the whole suite

## Overview

Spec 0057's fourteen Tasks each carried a whole-package suite command as the
last line of their Verification — about 28 minutes per pass, repaid on every
retry, proving something the Run-level gate already proves. This Task removes
those commands from every active Spec and records the authoring rule that keeps
them from coming back.

## Requirements

1. MUST remove whole-package suite commands from the Verification of every
   active Spec's Task files, leaving each Task's focused checks that prove its
   own effect.
2. MUST NOT weaken any focused check, scope assertion, or build command; only
   the redundant whole-suite line goes.
3. MUST record the authoring rule in the Task-authoring skill: a Task proves
   its own effect, and the Run-level gate proves nothing else regressed.
4. MUST record alongside it the rule that a Verification command must be able
   to fail when no work was done, so a command naming a missing test cannot
   pass.
5. MUST leave archived Specs byte-identical.
6. MUST change only the owned skill pair among protected tooling, per this
   Spec's Tooling authority row; any other tooling path is out of scope.

## Subtasks

- [ ] Remove whole-package suite lines from active Specs' Task Verification.
- [ ] Confirm each Task keeps a focused check proving its own effect.
- [ ] Record both authoring rules in the Task-authoring skill.
- [ ] Confirm archived Specs are untouched.

## Acceptance Criteria

- [ ] No active Spec's Task file carries a whole-package suite command in its
      Verification.
- [ ] Every active Task file still carries at least one focused check naming a
      specific test or asserting a specific effect.
- [ ] The Task-authoring skill states that a Task proves its own effect and the
      Run-level gate owns regression.
- [ ] The Task-authoring skill states that a Verification command must be able
      to fail when no work was done.
- [ ] Archived Spec artifacts are byte-identical.
- [ ] `git status --porcelain` shows no path outside `docs/specs/`,
      `.agents/skills/`, `skills/`, and this task file.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `.agents/skills/write-tasks/SKILL.md`

## Verification

- `grep -rl -- "go test ./internal/" docs/specs --include='task_*.md' | grep -v _archived | xargs -r grep -l -- '-count=1$' | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no active Task ends a Verification line with a bare
  whole-package suite command.
- `grep -q 'prove its own effect' .agents/skills/write-tasks/SKILL.md` —
  expected: exit 0; the authoring rule is recorded.
- `grep -q 'able to fail when no work was done' .agents/skills/write-tasks/SKILL.md`
  — expected: exit 0; the vacuous-Verification rule is recorded.
- `diff -r .agents/skills/write-tasks skills/write-tasks` — expected: exit 0;
  the embedded mirror does not drift.
- `git diff --name-only HEAD -- docs/specs/_archived | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no archived artifact changed.

## References

- `_prd.md` → Core Features 3 and 4; Goals (verification proportional to what
  changed).
- `_techspec.md` → Build Order 5; Project Constraints: Tooling authority.
