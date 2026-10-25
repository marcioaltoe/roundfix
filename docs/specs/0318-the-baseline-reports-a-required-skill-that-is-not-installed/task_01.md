---
task: task_01
spec: 0318-the-baseline-reports-a-required-skill-that-is-not-installed
status: completed
type: backend
complexity: medium
---

# Task 01: The baseline update preview names a missing required skill

## Overview

The `roundfix baseline update` preview inspects the external Repository Skill
Set only for `skills.differing` and reports `current`, exit `0`, while a
required external skill is absent or lock-mismatched and `--yes` would
restore it. This Task adds the helper that names each Missing Required Skill
with its restore command, makes the preview list them under `skills.missing`
and report `plan_ready`, exit `3`, pins Doctor's existing failure for both
removal shapes, and documents the result.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-25-the-baseline-update-preview-reports-current-while-a-required-skill-is-missing.md`](references/2026-10-25-the-baseline-update-preview-reports-current-while-a-required-skill-is-missing.md)
   of 2026-10-25 (`_prd.md` → Sources 1) by implementing API Contracts 1-4
   and 6 of `_techspec.md`.
2. MUST add `internal/cli/required_skills.go` with `missingRequiredSkills`
   and a reader of the Setup Manifest's `profile`, with the sources, order and
   exact reasons of `_techspec.md` → Interfaces. An inspection error that
   names no skill adds nothing.
3. MUST, in `internal/cli/baseline_update.go`, add the `Missing` field and the
   `missing` status of `_techspec.md` → Data Models, fill the list in the
   preview branch from the inspection it already makes (no second inspection,
   nothing under `--no-skills`, `--yes` or `--confirm-plan`), return the
   state, category, message, next action and exit `3` of Interfaces before
   the outdated-owned branch when guidance is unchanged, and print the
   `Skills missing` block of Interfaces in text. An empty list leaves every
   byte unchanged.
4. MUST add `internal/cli/baseline_update_missing_skills_test.go` with the
   `Suite:`/`Invariant:`/`Boundary` header, the helpers
   `installRequiredExternalSkills` and `removeInstalledSkill`, and
   `TestBaselineUpdatePreviewNamesAMissingRequiredSkill` and
   `TestDoctorNamesAMissingRequiredSkill` exactly as `_techspec.md` →
   Testing Approach describes, with those subtest names. Each top-level test
   calls `t.Parallel()` first. The tests pick the skill they remove from the
   required list the helper returns, never by a literal skill name, and never
   read this repository's skills or the host's.
5. MUST change `adoptedRepositoryWithRetiredSkills` in
   `internal/cli/baseline_update_retired_skills_test.go` to call
   `installRequiredExternalSkills` after adoption, and nothing else in that
   file: the fixture's lock otherwise lacks the required skills and its
   preview would no longer be `current`.
6. MUST append the guide paragraph of `_techspec.md` → Interfaces to
   `docs/user-guide/commands/baseline.md` and to the canonical
   `.agents/skills/roundfix/references/baseline.md`, each directly after the
   sentence ending "is still restored under `skills.drifted`.", and add the
   `CONTEXT.md` entry of Interfaces before `**Park Class**:` through the
   `domain-modeling` skill. Never touch a `### QA settlement` section.
7. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   (it raises both version fields of the Roundfix skill and records them),
   then `make skills-sync` again so the mirror matches.
8. MUST, as the Task's last step, after every guide and skill edit,
   re-record the Behavior Surface fingerprint ("re-record after the last
   guide edit") with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST NOT change Doctor's output, `skills.differing`, the applied path,
   `internal/cli/cli.go`, any Baseline module or derived Baseline file,
   `docs/agents/`, the Makefile or the CI workflows.

## Subtasks

- [ ] Add the helper and the preview list, status, state and text.
- [ ] Add the preview and Doctor tests and give the retired fixture its required skills.
- [ ] Add the guide and reference paragraph and the glossary term.
- [ ] Sync and record the Roundfix skill, then re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] Each removal shape makes the preview exit `3` with `plan_ready` and the skill under `skills.missing`, and the repository bytes are unchanged.
- [ ] A complete Repository Skill Set and `--no-skills` keep the preview `current`.
- [ ] Doctor fails with the skill named for both removal shapes; the owned and retired preview tests pass.

## Consumers

None. `skills.missing` is a new optional JSON list and `missing` a new
`skills.status` value of the `baseline update` result; no Go package parses
that result, and task_02 calls the helper.

## Context

- instruction: `internal/cli/doctor.go`
- instruction: `skills/repository.go`
- instruction: `internal/cli/baseline_update_outdated_skills_test.go`
- instruction: `internal/cli/baseline_update_test.go`
- instruction: `internal/cli/baseline_plan_test.go`
- instruction: `internal/cli/doctor_trailing_skills_test.go`
- creates: `internal/cli/required_skills.go`
- creates: `internal/cli/baseline_update_missing_skills_test.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/baseline_update_retired_skills_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -timeout 15m -count=1 -v -run '^(TestBaselineUpdatePreviewNamesAMissingRequiredSkill|TestDoctorNamesAMissingRequiredSkill|TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState|TestBaselineUpdatePreviewReportsAnOlderOwnedSkill|TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineUpdatePreviewNamesAMissingRequiredSkill/folder_and_lock_entry_removed TestBaselineUpdatePreviewNamesAMissingRequiredSkill/folder_removed_with_a_stale_lock_entry TestBaselineUpdatePreviewNamesAMissingRequiredSkill/folder_no_longer_matches_its_lock_entry TestBaselineUpdatePreviewNamesAMissingRequiredSkill/every_required_skill_installed TestDoctorNamesAMissingRequiredSkill/folder_and_lock_entry_removed TestDoctorNamesAMissingRequiredSkill/folder_removed_with_a_stale_lock_entry TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState TestBaselineUpdatePreviewReportsAnOlderOwnedSkill TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -timeout 10m -count=1 ./skills && go test -timeout 15m -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Missing Required Skill**' && for guide in docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$guide" | grep -qF -- 'each external skill the adopted modules require and dispatch that is absent from' || { printf 'missing paragraph: %s\n' "$guide" >&2; exit 1; }; done` — expected: exit 0. Before this Task the two new tests do not exist, so the command fails at its first pass check, and the glossary term and guide paragraph are absent. Afterwards every removal shape and the complete set pass, Doctor's failure is pinned, the retired and owned preview tests pass, the skill version, mirror and Behavior Surface record are current, and the term and paragraph exist.

## References

- `_prd.md` → Sources; Goals; Core Features 1, 3 and 4; Success Metrics 1, 2, 4, 5 and 6; Glossary
- `_techspec.md` → Current behavior; Interfaces; Data Models; API Contracts 1-4 and 6; Invariants 1-4; Testing Approach; Build Order 1

## Result

Implemented only task_01's slice. Status and declared Verification remain
Daemon-owned; no commit, push, Task Graph edit or other Task edit was made.
The Task status was already `in_progress` on entry and was preserved.

- Added `missingRequiredSkills`, merging readiness missing names, wrapped
  readiness-error names and lock mismatches in that priority, deduplicated and
  sorted, with the exact restore reasons. Added the Setup Manifest profile
  reader for task_02's consumer.
- The preview reuses its external inspection to emit `skills.missing` and
  status `missing`; unchanged guidance returns `plan_ready`, approval, exit
  `3`, the named message and restore action before the outdated-owned branch.
  Text adds the conditional missing block. Doctor, the applied path and
  `skills.differing` were not changed.
- Added hermetic temporary-repository helpers and the specified preview and
  Doctor subtests. The removed skill comes from the manifest-derived required
  list. The retired fixture's only change installs those required externals
  while retaining its existing retired entries.
- Added the exact guide/reference paragraph and the Missing Required Skill
  glossary entry before Park Class. The Roundfix skill was synced, recorded
  as `0.1.137` by the sanctioned recorder, and synced again. No QA settlement
  section was edited.

Focused evidence by acceptance criterion:

1. Before implementation, the focused folder-and-lock removal subtest failed:
   the preview returned exit `0`, `ok:true`, `current`, with no missing list.
   After implementation, all three removal-shape subtests passed assertions
   for JSON exit `3`, `ok:false`, approval, exact list/reasons/message/action,
   text output, and unchanged repository bytes.
2. The complete-set subtest passed `current`, exit `0`, and no missing key.
   Each removal-shape subtest also passed its `--no-skills` current/exit-0
   assertion and unchanged-tree check.
3. Both Doctor removal subtests passed with exit `1`, the failed skills line,
   skill name and existing `bunx skills add` action. Focused retired-skill and
   owned-skill preview compatibility tests passed as part of the same run.
   Additional helper tests passed sorted/deduplicated precedence, a wrapped
   readiness error, an unnamed inspection error, and manifest parsing/errors.

Commands and outcomes:

- `go test -timeout 5m ./internal/cli -run '^TestBaselineUpdatePreviewNamesAMissingRequiredSkill$/folder_and_lock_entry_removed$' -count=1`:
  failed before implementation with the reproduced current/exit-0 bug.
- Initial focused run exposed a redundant fixture commit in the complete-set
  case; corrected by committing only removal/mismatch mutations.
- A subsequent run's assertions passed, but suiteguard rejected an overlapping
  skill-sync write. This was an execution sequencing error; generation and
  tests were then run sequentially, without weakening the guard.
- `make skills-sync`, then
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  then `make skills-sync`: all exit `0`.
- `go test -timeout 5m -count=1 ./internal/cli -run '^(TestBaselineUpdatePreviewNamesAMissingRequiredSkill|TestDoctorNamesAMissingRequiredSkill|TestMissingRequiredSkillsCollectsNamedEvidence|TestRequiredSkillsManifestProfileReadsOnlyProfile|TestBaselineUpdate.*Retired.*|TestBaselineUpdatePreview.*OwnedSkill.*)$'`:
  exit `0`, `ok roundfix/internal/cli`, after the final skill sync, with no
  concurrent mutation.
- `git diff --check`: exit `0`. Diff review confirms only declared task_01
  paths changed, plus its own Result; the pre-existing Task status diff was
  preserved. Canonical and embedded Roundfix content are synced.

The sanctioned Behavior Surface recorder runs after this Result as the last
implementation step, after every guide and skill edit; its outcome is reported
in the Agent handback. The declared Verification and full repository checks
were not run, as required by the Daemon-assigned execution contract.
