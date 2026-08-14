---
task: task_08
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: docs
complexity: low
---

# Task 08: Write the exit-zero rule where an author reads it

## Overview

A Verification command passes only by exiting zero. The Daemon enforces it and
nothing states it, so a natural-looking assertion whose success is an empty result
fails its Task for doing the right thing. This slice writes the rule beside the
vacuity rule it belongs next to, with the working forms as the worked answer.

## Requirements

1. MUST state that a Verification command passes only by exiting zero, in the
   authoring contract, beside the vacuity rule.
2. MUST give the working forms for an assertion whose success is an empty result
   or an absent string.
3. MUST NOT change any other clause of the authoring contract.
4. MUST NOT change any repository path outside the bounded scope below plus this
   Task file; stop and fail the Task if a changed-file check finds another path.

## Subtasks

- [ ] Add the exit-zero rule beside the vacuity rule.
- [ ] Record the working forms.

## Acceptance Criteria

- [ ] The authoring contract states the exit-zero rule.
- [ ] It gives at least the three working forms an author needs for an
      empty-result assertion.
- [ ] No other clause changed.
- [ ] The changed-file set is the bounded scope plus this Task file.

## Bounded scope

This Task may create or modify only:

- `.agents/skills/write-tasks/SKILL.md`
- `docs/specs/0095-a-verification-that-ran-before-anyone-believed-it/task_08.md`

Express maintainer authorization:
`docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`,
whose per-Spec section for this Spec names exactly this file and is binding.
Skill mirrors rewritten by the sanctioned sync command are fallout under
ADR-0081.

## Verification

- `f=.agents/skills/write-tasks/SKILL.md; grep -q 'exiting zero' "$f" || { echo "FAIL: $f does not state the exit-zero rule"; exit 1; }` — expected: exits 0, printing what is missing on failure. Fails today.
- `f=.agents/skills/write-tasks/SKILL.md; for form in 'test -z' '! grep -rq' '|| { cat'; do grep -qF "$form" "$f" || { echo "FAIL: $f lacks the working form $form"; exit 1; }; done` — expected: exits 0, proving the worked answer is present rather than only the rule. Fails today.
- `make skills-sync-check > /tmp/0095-t08.log 2>&1 && grep -q 'exiting zero' skills/write-tasks/SKILL.md || { cat /tmp/0095-t08.log; exit 1; }` — expected: exits 0, proving the mirror carries the rule and matches its source. The sync check alone passes before any work, so it is anchored to the rule it guards.
- `git diff --name-only HEAD > /tmp/0095-t08-all.txt; test -s /tmp/0095-t08-all.txt || { echo 'no file changed'; exit 1; }; grep -v -e '^\.agents/skills/write-tasks/SKILL\.md$' -e '^skills/write-tasks/SKILL\.md$' -e '^docs/specs/0095-a-verification-that-ran-before-anyone-believed-it/task_08\.md$' /tmp/0095-t08-all.txt > /tmp/0095-t08-scope.txt; test ! -s /tmp/0095-t08-scope.txt || { echo 'out of bounds:'; cat /tmp/0095-t08-scope.txt; exit 1; }` — expected: exits 0, proving work happened and every changed path is in bounds.

## Context

- instruction: `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`

## References

`_techspec.md` → Build Order 8. `_prd.md` → Core Feature 6; Goal 2; User
Story 2; Project Constraints, Tooling authority. ADR-0081.
