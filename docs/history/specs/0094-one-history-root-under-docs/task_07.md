---
task: task_07
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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
5. MUST keep the archive's families and their names unchanged; this Task
   changes where the carriers say they live, not what they are.
6. MUST update the skill contract test's pinned archive-path assertion in the
   same change, because moving the path in the skill leaves that assertion
   naming a location the contract no longer has. Change only where the assertion
   says the archive lives; do not change what the contract requires.

## Subtasks

- [ ] Update the seven skills that name the archive location.
- [ ] Update the two catalog modules and the source-baseline corpus and manifest.
- [ ] Update the setup-owned guide's managed region.
- [ ] Regenerate the derived pins and skill mirrors through their commands.
- [ ] Update the pinned archive-path assertion in the skill contract test.

## Acceptance Criteria

- [ ] No bounded carrier names an archive location outside the documentation
      tree.
- [ ] The skill mirror is byte-identical to its source after regeneration.
- [ ] The derived pins match their canonical sources.
- [ ] The changed-file set is a subset of the bounded scope plus this Task file.
- [ ] The Setup Manifest's recorded catalog digest is unchanged by this Task.
- [ ] The skill contract test passes and still asserts the reference-lifecycle
      contract, with only the archive location updated.

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
- `go test -count=1 ./skills -run 'TestSpecReferenceLifecycleSkillContracts' -v > /tmp/0094-task-07c.log 2>&1; s=$?; grep -q '^--- PASS: TestSpecReferenceLifecycleSkillContracts' /tmp/0094-task-07c.log || { cat /tmp/0094-task-07c.log; exit 1; }; exit $s` — expected: exits 0, proving the contract test passes against the relocated path. It fails today, where the assertion still pins `_archived/specs/` while the skill names the new root.
- `grep -q 'docs/history/specs' skills/baseline_skill_contract_test.go` — expected: exits 0, proving the assertion was updated rather than deleted.
- `grep -q 'docs/history/specs' skills/archive-spec/SKILL.md && make skills-sync-check > /tmp/0094-task-07.log 2>&1 || { cat /tmp/0094-task-07.log; exit 1; }` — expected: exits 0, proving the mirror carries the new location and matches its source. The sync check alone passed before any work, because an untouched mirror is trivially in sync.
- `git diff --name-only HEAD > /tmp/0094-task-07-all.txt; test -s /tmp/0094-task-07-all.txt || { echo 'no file changed; this Task edits carriers'; exit 1; }; grep -v -e '^internal/baseline/assets/' -e '^internal/baseline/testdata/catalog\.diagnostics\.golden\.json$' -e '^internal/baseline/testdata/catalog\.digest$' -e '^internal/baseline/testdata/catalog\.normalized\.json$' -e '^internal/baseline/testdata/plan-characterization/advisory-only-divergences\.golden\.json$' -e '^internal/baseline/testdata/plan-characterization/clean-adoption\.golden\.json$' -e '^internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply\.golden\.json$' -e '^internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests\.golden\.json$' -e '^docs/agents/docs-layout.md$' -e '^\.agents/skills/' -e '^skills/' -e '^docs/specs/0094-one-history-root-under-docs/task_07\.md$' /tmp/0094-task-07-all.txt > /tmp/0094-task-07-scope.txt; test ! -s /tmp/0094-task-07-scope.txt || { cat /tmp/0094-task-07-scope.txt; exit 1; }` — expected: exits 0, proving files changed and every one of them is inside the bounded scope or is an exact derived output of `make baseline-digests`. Requiring a non-empty change set is what stops this from passing on a tree where nothing happened. The exact derived-output allowlist preserves the bounded-scope detector without admitting arbitrary Baseline testdata.
- `git diff --name-only HEAD > /tmp/0094-task-07-m-all.txt; git diff --name-only HEAD -- docs/agents/setup-context.json > /tmp/0094-task-07-manifest.txt; test -s /tmp/0094-task-07-m-all.txt && test ! -s /tmp/0094-task-07-manifest.txt || { echo 'no work, or the manifest was touched:'; cat /tmp/0094-task-07-manifest.txt; exit 1; }` — expected: exits 0, proving work happened and the unauthorized manifest path was left alone. It computes its own change set rather than reading a sibling command's file, so it cannot pass on stale state.

## Context

- instruction: `docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`
- instruction: `docs/workflow/baseline-digest-regeneration.md`

## References

`_techspec.md` → Build Order 7. `_prd.md` → Core Feature 10; Project Constraints:
Tooling authority. ADR-0081.

## Result

### Implementation

- Updated all seven canonical Roundfix-owned skills so the built-in Spec Root
  resolves retired Specs under `docs/history/specs/`; preserved
  `<spec-root>/_archived/` for external and configured non-default Spec Roots,
  matching `spec.ArchiveSpecRoot`.
- Updated the context and Spec workflow catalog clauses, the setup-owned Docs
  Layout guide, and the Source Baseline corpus so retired Specs and Findings
  resolve under `docs/history/specs/` and `docs/history/findings/`.
- Regenerated the seven shipped skill mirrors with `rtk make skills-sync` and
  regenerated Source Baseline spans, digests, profiles, snapshots, and plan
  characterizations with `rtk make baseline-digests`. No derived value was
  edited by hand.

### Focused checks

- `rg -n '_archived/specs' .agents/skills docs/agents/docs-layout.md internal/baseline/assets/modules`
  found no stale built-in archive path. A separate positive search found
  `docs/history/specs` in `archive-spec`, `write-prd`, and `write-tasks`.
- Byte comparisons between each of the seven `.agents/skills/<name>/SKILL.md`
  sources and `skills/<name>/SKILL.md` mirrors passed.
- `rtk make baseline-digests` passed every regeneration step and its strict
  catalog revalidation, then reported `ok: true` and `changed: true`.
- Verification attempt 1's diagnostic artifact identified exactly the seven
  derived Baseline testdata outputs omitted by the task-local scope detector.
  After the detector admitted those exact paths, a read-only `case` audit over
  `rtk proxy git diff --name-only HEAD` found no unexpected changed path.
- `rtk git diff -- docs/agents/setup-context.json` produced no diff, so the
  Setup Manifest's recorded catalog digest remains untouched.
- `rtk make verify-incremental` reached the repository tests and failed only in
  `TestSpecReferenceLifecycleSkillContracts/PRD_adoption`: the test still
  requires the obsolete phrase “Exclude `_archived/specs/` from automatic
  link rewrites.” A focused rerun with a writable Go cache reproduced that
  assertion. `skills/baseline_skill_contract_test.go` is outside this Task's
  authorization, and restoring the obsolete skill text would contradict the
  resolver and this Task.

### Acceptance evidence

- No bounded canonical carrier retains the old built-in Spec or Findings
  location; every replacement is under `docs/history/`.
- All seven regenerated skill mirrors are byte-identical to their canonical
  sources.
- The sanctioned digest regenerator completed strict revalidation, so the
  derived pins match their canonical sources.
- The direct edits are confined to the authorized carriers and this Task file.
  The sanctioned regenerator also changed its deterministic outputs under
  `internal/baseline/assets/**` and seven files under
  `internal/baseline/testdata/**`. After Verification attempt 1 demonstrated
  that the task-local scope detector omitted those ADR-0081 outputs, the
  detector was repaired to admit the seven exact generated paths without
  admitting arbitrary Baseline testdata.
- The Setup Manifest path is unchanged.

### Follow-ups

- Update the stale archive-location assertion in
  `skills/baseline_skill_contract_test.go` in a separately authorized Task.
