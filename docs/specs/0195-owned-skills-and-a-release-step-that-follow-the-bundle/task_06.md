---
task: task_06
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
status: completed
type: chore
complexity: low
---

# Task 06: The Roundfix skill declares the version its merged content needs

## Overview

task_01 raised the Roundfix skill to `0.0.3`, and task_02 recorded that version's digest in `skills/testdata/owned-skill-versions.json`. Spec 0192 then landed on main and changed the Roundfix skill's text. After this branch was rebased onto main, the skill carries new content under `0.0.3`, and the check this Spec added refuses it: "roundfix: content changed under version 0.0.3; raise the version". The record also lacks the versions Spec 0192 raised (for example `archive-spec` `0.0.3`). task_04 could not complete because it may not change a skill. This Task applies this Spec's own rule to the merged skill.

## Requirements

1. MUST raise both version fields of `.agents/skills/roundfix/SKILL.md` (`metadata.version` and the top-level `version`) from `0.0.3` to `0.0.4`, and change nothing else in the skill.
2. MUST regenerate `skills/roundfix/SKILL.md` only with `make skills-sync`.
3. MUST re-record `skills/testdata/owned-skill-versions.json` only with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, keeping every entry already recorded.
4. MUST NOT edit any other skill, test, Baseline asset or generated guide.

## Subtasks

- [ ] Raise the canonical version and sync the bundle.
- [ ] Re-record the owned skill versions and prove the check passes.

## Acceptance Criteria

- [ ] `TestEveryOwnedSkillVersionIsRecorded` passes, and `make skills-version-check` passes.
- [ ] The Task's diff touches only the three files it declares.

## Result

Updated the canonical Roundfix skill's metadata and top-level version from
`0.0.3` to `0.0.4`, synchronized the embedded copy, and recorded the current
owned-skill versions and digests while preserving the existing entries.

Focused checks:

- `rtk make skills-sync` — passed.
- `GOCACHE=/tmp/roundfix-task06-gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed.
- `git diff --check` — passed.
- The implementation diff contains only `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, and `skills/testdata/owned-skill-versions.json`; this Result is the required Task-file evidence.

The declared Verification commands remain for the Daemon; no terminal Task
status was changed.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills && make skills-version-check && grep -q '^version: 0.0.4$' skills/roundfix/SKILL.md && cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: exit 0; before this Task the check fails with "roundfix: content changed under version 0.0.3".

## References

- task_01, task_02, task_04
- ADR-0189
- Run `run_20260930T203751Z_2b4c20bff8fbb572` (task_04 Result and verification logs)
