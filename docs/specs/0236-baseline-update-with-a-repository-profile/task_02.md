---
task: task_02
spec: 0236-baseline-update-with-a-repository-profile
status: pending
type: backend
complexity: medium
---

# Task 02: Baseline update, Doctor and skill restore work end to end with a repository profile

## Overview

With task_01's loader in place, a repository adopted with its own Baseline
Profile reaches `state: current` through the real `roundfix baseline update`,
its preview lists a trailing external skill, its `--yes` skills stage
refreshes the Roundfix-owned skills, Doctor compares its skills, and
`baseline skills restore` accepts its profile. This Task proves those paths
with regression tests through the real commands, and describes repository
profiles in the baseline and doctor references and the Roundfix Skill's
baseline reference.

## Requirements

1. MUST add, in a new CLI test file, a fixture that adopts a temporary Git
   repository through the real `baseline profile init --id <id> --from
   go-cli-tui`, `baseline plan --profile <id>` with the decisions
   `baselineApplyTestPlan` passes, and `baseline apply`, then copies this
   repository's `.agents/skills` and `skills-lock.json` as
   `TestBaselineUpdatePreviewListsATrailingSkill` does. It closes the
   Backlog Entry of 2026-10-06, "`baseline update` refuses every repository
   Baseline Profile".
2. MUST add `TestBaselineUpdateReachesCurrentWithARepositoryProfile`: on the
   fixture, `baseline update --format json` exits 0 with `state: current`
   and no stderr, and the text form matches Surface Transcript 1; on a
   fixture without the copied skill set the text form also matches it.
3. MUST add `TestBaselineUpdatePreviewListsATrailingSkillWithARepositoryProfile`:
   with one locked external skill edited and its lock hash updated, the
   preview exits 3 with `plan_ready`, lists that skill under
   `skills.drifted` with "trails its Setup Snapshot", and changes no
   repository bytes.
4. MUST add `TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile`:
   `runBaselineUpdateSkillsStageWith` on the fixture with the real owned
   install, external resolution, readiness check and `TrailingSetupSkills`,
   and a fake restore, reports `verified`, installs every owned skill, and
   asks restore for the edited skill with the repository profile's ID.
5. MUST add `TestDoctorComparesARepositoryProfileWithItsSnapshot`: Doctor's
   repository skills check on the fixture with the edited skill reports
   `DR-SKILL-TRAILS-SNAPSHOT` naming it and no "snapshot comparison
   unavailable"; with the manifest naming a missing repository profile it
   says "snapshot comparison unavailable" and names
   `.roundfix/baseline/profiles/<id>.json`, not "built-in Baseline Profile".
6. MUST add `TestBaselineSkillsRestoreNamesTheProfilePath` and
   `TestBaselineSkillsRestoreAcceptsARepositoryProfile`, asserting the
   stdout, stderr and exit of Surface Transcript 2 and Surface Transcript 3
   through the real command.
7. MUST describe in `docs/user-guide/commands/baseline.md` that restore,
   reconcile and the managed refresh accept a repository profile, that its
   restorable skills are the external skills its modules require with
   contracts the embedded Setup Snapshots agree on, that its `setup` is
   null, and the findings `restore.profile-unresolved` and
   `restore.snapshot-conflict`; and in `docs/user-guide/commands/doctor.md`
   that an unresolvable profile's comparison is unavailable with a detail
   that says "neither built-in nor resolvable" and names the searched path.
8. MUST change `.agents/skills/roundfix/references/baseline.md` so the
   restore and reconcile examples take `--profile <profile-id>`, and state in
   one sentence that a repository profile is accepted and that
   `restore.profile-unresolved` names its path; then run `make
   skills-sync`, re-record the owned skill versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   after the last skill edit, and run `make baseline-digests`, as
   implementation steps outside Verification. Never edit a `### QA
   settlement` section or a version or digest by hand.
9. MUST NOT change `baseline update` or Doctor production code unless a test
   above fails for a reason task_01 did not cover; MUST NOT change
   `CONTEXT.md` or `CHANGELOG.md`.

## Subtasks

- [ ] Build the repository-profile adoption fixture.
- [ ] Add the update, preview and skills-stage tests.
- [ ] Add the Doctor and restore tests.
- [ ] Describe repository profiles in the two references and the skill.
- [ ] Sync the skill and record its version.

## Acceptance Criteria

- [ ] A repository-profile repository reaches `current` through the real
      `baseline update`.
- [ ] Its preview lists a trailing skill, and its `--yes` skills stage
      refreshes the owned skills and restores with its profile ID.
- [ ] Doctor compares its skills and names the path of a missing profile.
- [ ] Both restore transcripts match through the real command.
- [ ] The references and the skill name the new findings; each mirror equals
      its canonical file and the raised version is recorded.

## Context

- instruction: `docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md`
- instruction: `internal/cli/doctor_trailing_skills_test.go`
- instruction: `internal/cli/baseline_update_test.go`
- instruction: `internal/cli/baseline_apply_test.go`
- instruction: `internal/cli/baseline_update.go`
- instruction: `internal/cli/doctor.go`
- creates: `internal/cli/baseline_update_repository_profile_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestBaselineUpdateReachesCurrentWithARepositoryProfile|TestBaselineUpdatePreviewListsATrailingSkillWithARepositoryProfile|TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile|TestDoctorComparesARepositoryProfileWithItsSnapshot|TestBaselineSkillsRestoreNamesTheProfilePath|TestBaselineSkillsRestoreAcceptsARepositoryProfile|TestBaselineUpdatePreviewListsATrailingSkill)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineUpdateReachesCurrentWithARepositoryProfile TestBaselineUpdatePreviewListsATrailingSkillWithARepositoryProfile TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile TestDoctorComparesARepositoryProfileWithItsSnapshot TestBaselineSkillsRestoreNamesTheProfilePath TestBaselineSkillsRestoreAcceptsARepositoryProfile TestBaselineUpdatePreviewListsATrailingSkill; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the new tests do not exist, so their pass lines are missing and the command fails.
- `for phrase in 'restore.profile-unresolved' 'restore.snapshot-conflict'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- "$phrase" || { printf 'missing phrase in baseline guide: %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- 'neither built-in nor resolvable' || { printf 'missing phrase in doctor guide\n' >&2; exit 1; }; for phrase in 'skills restore --repo . --profile <profile-id>' 'restore.profile-unresolved'; do tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- "$phrase" || { printf 'missing phrase in skill reference: %s\n' "$phrase" >&2; exit 1; }; done; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md && out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass: TestEveryOwnedSkillVersionIsRecorded\n' >&2; exit 1; }` — expected: exit 0; before this Task neither guide nor the skill reference names the new findings, so the command fails; after it the mirrors equal their canonical files and the raised version is recorded.

## References

- `_prd.md` → Goals; Core Features 3, 4, 5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Vocabulary Contract; Testing Approach; Build Order 2
- ADR-0241; ADR-0221; ADR-0189
