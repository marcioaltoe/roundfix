---
task: task_07
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: docs
complexity: medium
---

# Task 07: State the history location in every carrier that names it

## Overview

Seven skills, two catalog modules, the source-baseline corpus and manifest, and
one setup-owned guide state where retired material lives, and all of them name
the old location. This slice makes the documentation agree with the resolver, so
a reader following the guidance reaches the place the tool writes to. It is an
authorized tooling Task and may change only its bounded files.

## Requirements

1. MUST state the archive location as the resolver now answers it, in every
   carrier named in the bounded scope below.
2. MUST regenerate the derived pins and skill mirrors through their sanctioned
   commands, never by editing a pin value by hand.
3. MUST NOT change any repository path outside the bounded scope plus this Task
   file; stop and fail the Task if a changed-file check finds another path.
4. MUST NOT repair the Setup Manifest's recorded catalog digest, which this
   Task's edits leave stale. That path is outside the grant, and the drift is a
   defect owned by its own backlog entry rather than something to fix here.
5. MUST keep the archive's four families and their names unchanged; this Task
   changes where the carriers say they live, not what they are.

## Subtasks

- [ ] Update the seven skills that name the archive location.
- [ ] Update the two catalog modules and the source-baseline corpus and manifest.
- [ ] Update the setup-owned guide's managed region.
- [ ] Regenerate the derived pins and skill mirrors through their commands.

## Acceptance Criteria

- [ ] No bounded carrier names an archive location outside the documentation
      tree.
- [ ] The skill mirror is byte-identical to its source after regeneration.
- [ ] The derived pins match their canonical sources.
- [ ] The changed-file set is a subset of the bounded scope plus this Task file.
- [ ] The Setup Manifest's recorded catalog digest is unchanged by this Task.

## Bounded scope

This Task may create or modify only:

- `internal/baseline/assets/modules/spec-workflow.json`
- `internal/baseline/assets/modules/context-workflow.json`
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
- `docs/agents/docs-layout.md`
- `.agents/skills/archive-spec/SKILL.md`
- `.agents/skills/write-prd/SKILL.md`
- `.agents/skills/write-tasks/SKILL.md`
- `.agents/skills/write-techspec/SKILL.md`
- `.agents/skills/write-idea/SKILL.md`
- `.agents/skills/brainstorming/SKILL.md`
- `.agents/skills/roundfix/SKILL.md`
- `docs/specs/0094-one-history-root-under-docs/task_06.md`

Derived pins and skill mirrors rewritten by the sanctioned regeneration commands
are fallout of these edits under ADR-0081, not separate targets.

## Verification

- `! grep -rn '\b_archived/specs\b' .agents/skills docs/agents/docs-layout.md internal/baseline/assets/modules | grep -v 'docs/_archived'` — expected: exits 0, proving no bounded carrier still names the old location.
- `grep -q 'docs/history/specs' .agents/skills/archive-spec/SKILL.md && grep -q 'docs/history/specs' .agents/skills/write-prd/SKILL.md && grep -q 'docs/history/specs' .agents/skills/write-tasks/SKILL.md` — expected: exits 0, proving the new location reached the skills rather than only being removed from them.
- `make skills-sync-check > /tmp/0094-task-07.log 2>&1 || { cat /tmp/0094-task-07.log; exit 1; }` — expected: exits 0, proving the mirror matches its source.
- `git diff --name-only HEAD | grep -v -e '^internal/baseline/assets/' -e '^docs/agents/docs-layout.md$' -e '^\.agents/skills/' -e '^skills/' -e '^docs/specs/0094-one-history-root-under-docs/task_06\.md$' > /tmp/0094-task-07-scope.txt; test ! -s /tmp/0094-task-07-scope.txt || { cat /tmp/0094-task-07-scope.txt; exit 1; }` — expected: exits 0 with an empty file, proving no path outside the bounded scope changed.
- `git diff --name-only HEAD -- docs/agents/setup-context.json > /tmp/0094-task-07-manifest.txt; test ! -s /tmp/0094-task-07-manifest.txt || { cat /tmp/0094-task-07-manifest.txt; exit 1; }` — expected: exits 0 with an empty file, proving the unauthorized manifest path was left alone.

## Context

- instruction: `docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`
- instruction: `docs/workflow/baseline-digest-regeneration.md`

## References

`_techspec.md` → Build Order 7. `_prd.md` → Core Feature 10; Project Constraints:
Tooling authority. ADR-0081.
