---
task: task_01
spec: 0318-the-baseline-reports-a-required-skill-that-is-not-installed
status: pending
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
