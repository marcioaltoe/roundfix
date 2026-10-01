---
task: task_06
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: pending
type: chore
complexity: low
---

# Task 06: The two skills task_04 changed declare new versions

## Overview

task_04 added text to the Roundfix skill and the write-tasks skill but left both at `0.0.4`. The repository's owned-skill version rule refuses content that changes under a recorded version, so the first delivery Run's QA gate failed: "roundfix: content changed under version 0.0.4; raise the version" (Run `run_20260930T230750Z_197aa9aff15490a9`, `batch-005-attempt-1.log`).

## Requirements

1. MUST raise both version fields (`metadata.version` and the top-level `version`) of `.agents/skills/roundfix/SKILL.md` and `.agents/skills/write-tasks/SKILL.md` from `0.0.4` to `0.0.5`, and change nothing else in either skill.
2. MUST regenerate `skills/roundfix/SKILL.md` and `skills/write-tasks/SKILL.md` only with `make skills-sync`.
3. MUST re-record `skills/testdata/owned-skill-versions.json` only with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
4. MUST NOT edit any other file.

## Subtasks

- [ ] Raise the two canonical versions and sync the bundle.
- [ ] Re-record the owned skill versions.

## Acceptance Criteria

- [ ] `TestEveryOwnedSkillVersionIsRecorded` and `make skills-version-check` pass, and the two bundle copies equal their canonical files.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills && make skills-version-check && grep -q '^version: 0.0.5$' skills/roundfix/SKILL.md && grep -q '^version: 0.0.5$' skills/write-tasks/SKILL.md && cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp -s .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md` — expected: exit 0; before this Task the version check fails with "roundfix: content changed under version 0.0.4".

## References

- task_04
