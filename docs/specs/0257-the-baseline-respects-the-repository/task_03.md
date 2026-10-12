---
task: task_03
spec: 0257-the-baseline-respects-the-repository
status: completed
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

## Result

Implemented the task_03 slice for Daemon Verification. Task status remains
Daemon-owned; the Agent did not run authored Verification or repository
Verification, commit, push or open a pull request.

### Implementation

- Renamed the comparison to `SetupSnapshotDifferences`, its callers and both
  dependency fields. The comparison's profile resolution, digest comparison,
  missing-tree handling, sorting and deduplication are unchanged.
- Doctor appends the API Contract 5 detail suffix without changing readiness
  status, next action or exit code. `snapshotComparableSkills` is unchanged.
- The preview and skills stage report the API Contract 6 reason under the
  always-present `skills.differing` array. Text prints `Skills differing` and
  its entries after the drifted list and the existing optional outdated and
  retired sections, before History. A difference alone yields
  `current`; the skills stage never includes it in the restore set. Missing
  and lock-mismatched skills still use the existing preview/confirmation
  restore request and offline source parameter.
- Renamed the six tests specified by the TechSpec. Kept the two protected
  test names and the strict repository snapshot assertion; its failure
  message now names a difference and the deliberate per-skill restore.
- Updated both command guides and the canonical Roundfix skill references,
  synced their mirrors, and generated the Roundfix skill version `0.1.65`.
  Added the ADR-0221 addendum and advanced `updated_at`; added the exact
  Setup Snapshot Difference glossary entry with its ADR-0264 link through
  domain-modeling. Setup Snapshots, lock comparison and the restore command
  were not changed.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Doctor exits `0`, names the difference, and does not warn for a difference alone | `TestDoctorNamesASetupSnapshotDifferenceWithoutWarning` passes with the exact `skills: ok` line and exit `0`. `TestDoctorSnapshotComparisonPreservesReadinessAndFiltersSkills` compares status and next action directly with readiness and covers mixed missing/lock-mismatched skills and an unavailable profile. The repository-profile Doctor test also passes. |
| Preview is `current` with `skills.differing`; `--yes` preserves the copy | `TestBaselineUpdatePreviewListsASetupSnapshotDifference` passes for JSON and text, checks the exact reason and line order, proves previews change no repository bytes, checks the empty array with `--no-skills`, and runs the real skills stage under `--yes`: the differing skill's bytes and whole-folder hash stay identical. `TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile` passes for a repository-owned profile. |
| Missing skills still restore | `TestBaselineUpdateKeepsASetupSnapshotDifference` passes for difference alone, a missing skill and a lock-mismatched skill. The latter two each restore only `golang-cli`, with the existing two-request confirmation and offline source. Its fake restore fails immediately if `testing-boss`, the differing skill, is requested. The repository-profile stage test installs every missing owned skill while keeping the external difference. |
| No live retired code, wording or function name | Local inspection found zero matches for the three retired strings and `TrailingSetupSkills` under all five named roots. Both canonical references contain the new code and are byte-identical to their mirrors; the entry-file pair also matches. |
| CONTEXT defines Setup Snapshot Difference | Local inspection found the definition beside Repository Skill Set, the required unknown-order and deliberate-restore semantics, and a resolving ADR-0264 link. ADR-0221 contains its required dated addendum. |

### Focused checks and generation

Commands below are focused implementation checks or explicitly required
record/generation commands, not the Task's authored Verification commands.

- Red starting point:
  `go test ./internal/cli -run '^TestDoctorNamesASetupSnapshotDifferenceWithoutWarning$' -count=1`
  exited `1` before the behavior change: it printed `skills: warn`, the
  retired code and `next: roundfix baseline update`, while the new test
  expected the informational line.
- Initial focused CLI check:
  `go test ./internal/cli -run '^(TestDoctorNamesASetupSnapshotDifferenceWithoutWarning|TestDoctorSnapshotComparisonPreservesReadinessAndFiltersSkills|TestBaselineUpdatePreviewListsASetupSnapshotDifference|TestBaselineUpdateKeepsASetupSnapshotDifference|TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile|TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile|TestDoctorComparesARepositoryProfileWithItsSnapshot|TestBaselineUpdateReachesCurrentWithARepositoryProfile)$' -count=1`
  exited `0` before the Verification Feedback repair below.
- Focused comparison/strictness check:
  `go test ./internal/baseline -run '^(TestSetupSnapshotDifferencesCompare|TestSkillSnapshotProfileNamesTheRepositoryPath|TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot)' -count=1`
  exited `0`.
- `make skills-sync`, then
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  then `make skills-sync` each exited `0`. The generator raised both
  Roundfix version fields and recorded the new digest.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
  generated the coverage record. The first invocation's test body passed,
  then suiteguard exited `1` for the record write; see the follow-up below.
  A second generator invocation against the generated record exited `0`,
  demonstrating byte-stable regeneration. The only removed recorded test
  names are the six intended renames. The generator also records tests
  already present before this Task that the old record did not contain;
  no record was hand-edited.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited `0` and changed only the two updated guide fingerprints in the
  Behavior Surface Record.
- `make baseline-digests` exited `0`, reported `changed:false`, and left no
  Baseline digest or golden changes.
- Local terminology/glossary/addendum/mirror inspection and
  `git diff --check` passed. Changed-file postflight restricts all changes
  to this Task's Context and this Task file; the initial dirty path was only
  the Daemon-owned `status: in_progress` change in this Task file.

### Sabotage evidence

A temporary Go overlay outside the repository added
`drifted = append(drifted, differing...)` immediately before the stage's
restore loop, deliberately putting the unordered difference back into the
restore set. The command
`go test -overlay /private/tmp/roundfix-0257-task03-sabotage/overlay.json ./internal/cli -run '^TestBaselineUpdateKeepsASetupSnapshotDifference$' -count=1`
exited `1`. All three scenarios failed with
`restore called for the differing skill`. The overlay changed no repository
source. The final focused CLI check above ran without the overlay and exited
`0`.

### Verification Feedback repair — attempt 1

Inspected the Daemon diagnostic artifact:
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261012T194228Z_1d1cb31ee8783296/verification/batch-003-attempt-1.log`.
The Daemon's `make verify-changed` run and the initial broader focused CLI
check both identified the same text-order regression in
`TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState` and
`TestBaselineUpdatePreviewReportsAnOlderOwnedSkill`. The diagnostic body is
not reproduced here.

API Contract 6 requires the differing section after the drifted list; it does
not require immediate adjacency. The renderer now prints it after the
existing optional outdated and retired sections, still after drifted and
before History. Every differing heading and entry remains present, including
`Skills differing: 0`, and both prior transcript fragments remain intact.
The repair changes only `internal/cli/baseline_update.go` and
`internal/cli/doctor_trailing_skills_test.go`, plus this Result. The initial
recommendation to revise the other test files is superseded by this repair;
neither file needs an edit.

Extended `TestBaselineUpdatePreviewListsASetupSnapshotDifference` with a
real-command preview containing a differing external skill, an older owned
skill and a retired skill together. It asserts all three contiguous report
sections and proves that the preview changes no repository bytes. An initial
fixture run exposed a missing temporary directory for the retired skill;
the fixture now creates that directory before writing its skill file.
No top-level test names changed, so the generated coverage record needs no
new entries.

Final focused repair check:
`go test ./internal/cli -run '^(TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState|TestBaselineUpdatePreviewReportsAnOlderOwnedSkill|TestBaselineUpdatePreviewListsASetupSnapshotDifference|TestBaselineUpdateKeepsASetupSnapshotDifference|TestDoctorNamesASetupSnapshotDifferenceWithoutWarning|TestDoctorSnapshotComparisonPreservesReadinessAndFiltersSkills|TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile|TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile|TestDoctorComparesARepositoryProfileWithItsSnapshot)$' -count=1`
exited `0`. This covers both diagnostic failures, both older-owned-skill
scenarios, the combined preview, informational Doctor output, repository
profiles, byte preservation and the missing/lock-mismatched restore contract.
The repair did not rerun `make verify-changed` or any authored Verification
command. The Daemon owns the next configured Verification run and Task status.

### Remaining follow-up outside this Task's bounded Context

- `internal/spec/coverage_test.go` writes its generated record without
  declaring the sanctioned regeneration command to suiteguard. This is why
  the first coverage-record generation reports a boundary violation even
  though the Spec explicitly authorizes that command and output. The record
  was generated correctly and the no-change invocation passes, but the
  generator's missing declaration still needs a separate authorized repair.
