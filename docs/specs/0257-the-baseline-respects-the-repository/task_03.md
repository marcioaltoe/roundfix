---
task: task_03
spec: 0257-the-baseline-respects-the-repository
status: pending
type: backend
complexity: high
---

# Task 03: A skill that differs from its Setup Snapshot is kept, not called trailing

## Overview

Doctor calls a lock-matching skill that differs from its Setup Snapshot
"trailing", raises a warning, and the managed refresh restores the older
snapshot over it, though nothing orders the two revisions. This Task replaces
that classification with a Setup Snapshot Difference that Doctor names without
a warning and the managed refresh lists and keeps, retires
`DR-SKILL-TRAILS-SNAPSHOT`, and keeps the restore of missing and
lock-mismatched skills (ADR-0264). It is verifiable on its own through Doctor
and update tests with fake skill trees. It answers the Backlog Entry of
2026-10-12,
[a skill newer than the Setup Snapshot is reported as trailing](references/2026-10-12-a-skill-newer-than-the-setup-snapshot-is-reported-as-trailing.md).

This is an authorized tooling Task. It may change only the files in its
Context, the coverage record its command rewrites, and this Task file.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12 on a newer skill reported as
   trailing by replacing `TrailingSetupSkills` with `SetupSnapshotDifferences`
   of `_techspec.md` → Interfaces, same inputs and result, and renaming every
   caller and dependency field.
2. MUST implement API Contract 5 in Doctor: the exact detail suffix, the
   `skills` status and next action left as readiness set them, and an
   unchanged exit code.
3. MUST implement API Contract 6 in the update preview and the skills stage:
   the always-present `skills.differing` array with the exact reason, the
   text lines after the drifted list, no `plan_ready` for a difference alone,
   no restore of a differing skill, and the existing restore of missing and
   lock-mismatched skills.
4. MUST rename and change the tests of `_techspec.md` → Testing Approach for
   task_03, keep `TestDoctorComparesARepositoryProfileWithItsSnapshot` and
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` by name, change
   the latter's failure message so that it no longer says "trails", and then
   re-record the coverage record with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`.
   `TestBaselineUpdateKeepsASetupSnapshotDifference` MUST fail if the restore
   dependency is called for the differing skill.
5. MUST replace every live occurrence of `DR-SKILL-TRAILS-SNAPSHOT`,
   "trails its Setup Snapshot" and "trails the Setup Snapshot" under
   `internal`, `cmd`, `docs/user-guide`, `.agents/skills` and `skills` with the
   task_03 texts of `_techspec.md` → Skill and document texts, including
   `docs/user-guide/commands/doctor.md`, `docs/user-guide/commands/baseline.md`
   and the `roundfix` skill's `baseline.md` and `setup.md` references, then run
   `make skills-sync`,
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   and `make skills-sync` again. `CHANGELOG.md` and ADR texts stay.
6. MUST append the task_03 addendum of `_techspec.md` → Skill and document
   texts to ADR-0221 and raise its `updated_at`.
7. MUST re-record `docs/references/behavior-surfaces.json` with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
8. MUST add **Setup Snapshot Difference** to `CONTEXT.md` through
   `domain-modeling` with the text of `_techspec.md` → Skill and document
   texts.
9. MUST prove the keep can fail: the Result records one sabotage, for example
   adding differing skills back to the restored set, and the test that failed.
10. MUST NOT change the lock comparison, the Setup Snapshots, the restore
    command, `snapshotComparableSkills`' filtering, or this repository's
    strict snapshot test's meaning (ADR-0224).

## Subtasks

- [ ] Rename the comparison and its callers.
- [ ] Change Doctor's line and the update's preview and skills stage.
- [ ] Rename and change the tests and re-record the coverage record.
- [ ] Replace the retired wording in guides and skills, and add the ADR addendum.
- [ ] Re-record the surface record and add the glossary term.

## Acceptance Criteria

- [ ] Doctor exits `0` with `DR-SKILL-DIFFERS-FROM-SNAPSHOT` and its `skills`
      line not `warn` when the only finding is a difference.
- [ ] The update preview reports `current` with the skill under
      `skills.differing`, and `--yes` leaves its bytes identical.
- [ ] A missing skill is still restored.
- [ ] No live source, guide or skill names the retired code or wording.
- [ ] `CONTEXT.md` defines **Setup Snapshot Difference**.

## Context

- instruction: `docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md`
- instruction: `docs/adr/0224-this-repository-holds-its-upstream-skills-to-the-snapshot-and-the-lock-not-to-a-pinned-digest.md`
- interface: `internal/baseline/skills_trailing.go`
- interface: `internal/baseline/skills_trailing_test.go`
- interface: `internal/baseline/skills_snapshot_profile_test.go`
- interface: `internal/baseline/repository_skills_snapshot_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/doctor_trailing_skills_test.go`
- interface: `internal/cli/baseline_update_repository_profile_test.go`
- interface: `docs/references/coverage-record.json`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDoctorNamesASetupSnapshotDifferenceWithoutWarning|TestBaselineUpdatePreviewListsASetupSnapshotDifference|TestBaselineUpdateKeepsASetupSnapshotDifference|TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile|TestDoctorComparesARepositoryProfileWithItsSnapshot)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDoctorNamesASetupSnapshotDifferenceWithoutWarning TestBaselineUpdatePreviewListsASetupSnapshotDifference TestBaselineUpdateKeepsASetupSnapshotDifference TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile TestDoctorComparesARepositoryProfileWithItsSnapshot; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestSetupSnapshotDifferencesCompareTheInstalledTreeWithTheSnapshot|TestSetupSnapshotDifferencesCompareARepositoryProfile|TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSetupSnapshotDifferencesCompareTheInstalledTreeWithTheSnapshot TestSetupSnapshotDifferencesCompareARepositoryProfile TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run '^TestCoverageEquivalence$' ./internal/spec || exit 1` — expected: exit 0; before this Task the renamed tests do not exist, so the command fails.
- `found="$(grep -rn -e 'DR-SKILL-TRAILS-SNAPSHOT' -e 'trails its Setup Snapshot' -e 'trails the Setup Snapshot' -e 'TrailingSetupSkills' internal cmd docs/user-guide .agents/skills skills 2>&1)"; test -z "$found" || { printf '%s\n' "$found" >&2; exit 1; }; for file in docs/user-guide/commands/doctor.md docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/baseline.md .agents/skills/roundfix/references/setup.md; do grep -qF -- 'DR-SKILL-DIFFERS-FROM-SNAPSHOT' "$file" || { printf 'missing difference code in %s\n' "$file" >&2; exit 1; }; done; for name in baseline setup; do cmp -s ".agents/skills/roundfix/references/$name.md" "skills/roundfix/references/$name.md" || { printf 'the %s reference mirror differs\n' "$name" >&2; exit 1; }; done; grep -qF -- '## Addendum — 2026-10-12' docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md || { printf 'ADR-0221 lacks its addendum\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Setup Snapshot Difference**:" || { printf 'CONTEXT.md lacks **Setup Snapshot Difference**\n' >&2; exit 1; }; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' ./skills || exit 1; go test -tags docscontract -count=1 -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0; before this Task the retired code is named in the Doctor source, so the command fails at its first check.

## References

- `_prd.md` → Goal 3; User Story 4; Core Features 3 and 5; Success Metric 3
- `_techspec.md` → API Contract 5; API Contract 6; Interfaces; Current behavior; Skill and document texts; Invariants 4, 5 and 10; Testing Approach; Build Order 3
- ADR-0264; ADR-0221; ADR-0224; ADR-0256
