---
task: task_04
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
status: pending
type: docs
complexity: low
---

# Task 04: A Task declares the one command file it changes

## Overview

After the split, a Task that changes one command edits one reference file of the Roundfix skill and one file of the command reference. The `write-tasks` skill and its Task template still describe a Task that declares a whole skill or guide. This Task states the new rule in both, with its consequence for Waves, and proves the consequence with the Wave collision detector: two Tasks in one Wave that declare different command files do not collide, and two that declare the same command file still do.

## Requirements

1. MUST state in `.agents/skills/write-tasks/SKILL.md`, in the rule on declared edits, that a Task which changes a command declares the one command file it changes: `.agents/skills/roundfix/references/<command>.md` with its mirror, and `docs/user-guide/commands/<command>.md`, never the skill's `SKILL.md` or `commands.md` unless it changes the entry file itself. The sentence MUST contain the exact words `declares the one command file it changes`.
2. MUST state the consequence in the same rule: two Tasks that declare the same file cannot share a Wave, so Tasks that change different commands can run together and Tasks that change the same command need an edge between them. The text MUST contain the exact words `cannot share a Wave`.
3. MUST show both paths in the `## Context` guidance of `.agents/skills/write-tasks/references/task-template.md`, as example `interface:` entries.
4. MUST NOT reword any other rule of the skill or the template, and MUST keep every phrase that `TestWriteTasksSkillStatesTheDeclaredPathRules` and `TestTaskAuthoringGuidanceNamesDeclarations` require.
5. MUST raise the `write-tasks` skill's version in both front-matter fields by one patch step from the value on its starting commit, following whatever version rule is in force on that commit, then run the sanctioned `make skills-sync` and `make baseline-digests`, and name in its Result every file either command rewrote.
6. MUST put the new tests in `internal/speccheck/wave_collision_command_files_test.go`, over temporary Spec fixtures only, reusing the helpers the existing Wave collision tests use. It MUST NOT change the detector.
7. MUST edit no governed file other than `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md` and `.agents/skills/write-tasks/references/task-template.md`.
8. MUST NOT name any decision by its `ADR-` identifier in this Task's `## Result`.

## Subtasks

- [ ] State the rule and its consequence in the skill, and show the paths in the template.
- [ ] Add the two Wave collision tests.
- [ ] Raise the version and regenerate the mirror.

## Acceptance Criteria

- [ ] Two Tasks in one Wave that declare `.agents/skills/roundfix/references/deliver.md` and `.agents/skills/roundfix/references/review.md` raise no `SC-WAVE-COLLISION`.
- [ ] Two Tasks in one Wave that both declare `.agents/skills/roundfix/references/deliver.md` raise `SC-WAVE-COLLISION`.
- [ ] The skill, its mirror and the template carry the rule, and the skill's version changed.
- [ ] `make skills-sync-check` passes.

## Context

- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/write-tasks/references/task-template.md`
- creates: `internal/speccheck/wave_collision_command_files_test.go`
- instruction: `internal/speccheck/surface_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestWaveCollisionAllowsDifferentCommandFiles|TestWaveCollisionRefusesTheSameCommandFile|TestWriteTasksSkillStatesTheDeclaredPathRules|TestTaskAuthoringGuidanceNamesDeclarations)$" ./internal/speccheck ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestWaveCollisionAllowsDifferentCommandFiles TestWaveCollisionRefusesTheSameCommandFile TestWriteTasksSkillStatesTheDeclaredPathRules TestTaskAuthoringGuidanceNamesDeclarations; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new named tests do not exist, so the command fails.
- `for pair in ".agents/skills/write-tasks/SKILL.md|declares the one command file it changes" ".agents/skills/write-tasks/SKILL.md|cannot share a Wave" ".agents/skills/write-tasks/SKILL.md|.agents/skills/roundfix/references/<command>.md" ".agents/skills/write-tasks/SKILL.md|docs/user-guide/commands/<command>.md" ".agents/skills/write-tasks/references/task-template.md|.agents/skills/roundfix/references/<command>.md" ".agents/skills/write-tasks/references/task-template.md|docs/user-guide/commands/<command>.md"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; add="$(git log -1 --format=%H --diff-filter=A -- internal/speccheck/wave_collision_command_files_test.go)" || exit 1; base="${add:+$add^}"; base="${base:-HEAD}"; tmp="$(mktemp -d)" || exit 1; git show "$base:.agents/skills/write-tasks/SKILL.md" > "$tmp/before.md" || exit 1; awk 'substr($0,1,8)=="version:" {print; exit}' "$tmp/before.md" > "$tmp/version-old"; awk 'substr($0,1,8)=="version:" {print; exit}' .agents/skills/write-tasks/SKILL.md > "$tmp/version-new"; if cmp -s "$tmp/version-old" "$tmp/version-new"; then printf 'the skill version did not change\n' >&2; exit 1; fi; make skills-sync-check` — expected: exit 0; before this Task the skill and the template lack the rule, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 2; Core Feature 4; Success Metric 4
- [_techspec.md](_techspec.md) — Testing Approach 5; Build Order 4
- ADR-0187

## Result
