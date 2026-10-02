---
task: task_01
spec: 0211-a-delivery-queue-that-finishes-without-intervention
status: pending
type: docs
complexity: low
---

# Task 01: The Roundfix Skill and the command guides describe the queue that finishes

## Overview

The repository's skill-sync rule requires a change to CLI behavior to ship
the Roundfix Skill update. This Task describes, in the `deliver` and
`runtime` references of the Roundfix Skill and in the `deliver` and `doctor`
command guides, the behavior task_02, task_03 and task_04 implement, as the
TechSpec states it: the archived retry
of a queue-started Run, the wait on GitHub's merge state, the return of a
policy-refused merge to checking, the `delivery-error` retry resumption, the
`Retry refused` output and the preload notice.

## Requirements

1. MUST describe in `.agents/skills/roundfix/references/deliver.md` and
   `docs/user-guide/commands/deliver.md`: that checking waits while GitHub's
   merge state is `BLOCKED` or `UNKNOWN`; that a merge refused with
   `base branch policy prohibits the merge` returns to checking once per head
   and parks `delivery-error` the second time; that a retry of a
   `delivery-error` park at an unchanged archived head with a Pull Request
   resumes checking; that the archived retry finds a Run the queue started;
   and the `Retry refused` output of Surface Transcript 1, using its words.
2. MUST describe in `.agents/skills/roundfix/references/runtime.md` and
   `docs/user-guide/commands/doctor.md` that Roundfix leaves a missing
   `NODE_OPTIONS` preload out of the agent environment, which preloads it
   checks, and the notice line of API Contract 4 verbatim.
3. MUST raise the Roundfix Skill's version by one patch level in both
   front-matter fields of `.agents/skills/roundfix/SKILL.md`, run
   `make skills-sync` so the mirrors equal their canonical files, and
   re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
4. MUST NOT edit the `### QA settlement` section of any skill, any other
   command reference, `commands.md` or `CONTEXT.md`.

## Subtasks

- [ ] Describe the retry, checking and merge changes in the `deliver` reference and guide.
- [ ] Describe the preload notice in the `runtime` reference and the `doctor` guide.
- [ ] Raise the version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The `deliver` reference and guide carry `Retry refused`, `merge state`
      and `base branch policy prohibits the merge`.
- [ ] The `runtime` reference and the `doctor` guide carry `NODE_OPTIONS` and
      the notice's words.
- [ ] Each mirror equals its canonical file, the version moved past `0.1.7`
      and is recorded.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/commands/doctor.md`

## Verification

- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "Retry refused" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "Retry refused" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "merge state" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "merge state" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "base branch policy prohibits the merge" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "base branch policy prohibits the merge" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- "NODE_OPTIONS" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/runtime.md "NODE_OPTIONS" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- "left it out of the agent environment" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/runtime.md "left it out of the agent environment" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "Retry refused" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "Retry refused" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "merge state" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "merge state" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- "NODE_OPTIONS" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/doctor.md "NODE_OPTIONS" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- "left it out of the agent environment" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/doctor.md "left it out of the agent environment" >&2; exit 1; }` — expected: exit 0; before this Task none of these phrases is in these files, so the command fails.
- `tr -s '[:space:]' ' ' < skills/roundfix/references/deliver.md | grep -qF -- 'Retry refused' && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && cmp .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md && ! grep -q '^version: 0.1.7$' .agents/skills/roundfix/SKILL.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the version is still `0.1.7`, so the command fails; after it the mirrors equal their canonical files and the raised version and its content digest are recorded.

## References

- `_prd.md` → Core Feature 7; User Stories 1-5
- `_techspec.md` → Vocabulary Contract; API Contract 1; API Contract 4; Surface Transcript 1; Build Order 1
- ADR-0187; ADR-0189; ADR-0211
