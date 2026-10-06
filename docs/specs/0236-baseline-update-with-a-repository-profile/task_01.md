---
task: task_01
spec: 0236-baseline-update-with-a-repository-profile
status: pending
type: backend
complexity: medium
---

# Task 01: The skill readers resolve a repository profile as baseline update does

## Overview

The snapshot comparison, skill restore and lock reconciliation stop loading
only built-in Baseline Profiles. A new loader keeps the built-in path
byte-for-byte and resolves any other ID through the same resolver `baseline
update` and `baseline profile validate` use; a repository profile then
requires the skills of its own modules, with external contracts from the
agreeing embedded Setup Snapshots. The slice is verifiable through the
baseline package alone: the comparison, restore and reconcile accept a
repository profile, and an unresolvable one is named by its repository path.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-06, "`baseline update` refuses
   every repository Baseline Profile"
   ([references/2026-10-06-baseline-update-refuses-a-repository-profile.md](references/2026-10-06-baseline-update-refuses-a-repository-profile.md)),
   by adding `loadSkillSnapshotProfile` and `repositorySnapshotContracts`
   beside `loadRestoreProfile` with the signatures and Invariants 1 to 5 of
   the TechSpec's Interfaces.
2. MUST make `TrailingSetupSkills`, `RestoreSkills` and
   `ReconcileSkillsLock` use `loadSkillSnapshotProfile` per Invariants 6 and
   7. Restore and reconcile resolve a non-built-in ID against the Git top
   level of the request's repository; their profile-ID consistency messages
   no longer say "built-in".
3. MUST add, in a new baseline test file, the tests
   `TestSkillSnapshotProfileKeepsEveryBuiltInProfile` (for every built-in ID
   the new loader equals `loadRestoreProfile`),
   `TestSkillSnapshotProfileResolvesARepositoryProfile` (a profile written by
   `InitCustomProfile` from `go-cli-tui` in a temporary repository has the
   union of its modules' required skills, the `go` snapshot's contract for
   each external one, its own ID and an empty setup),
   `TestSkillSnapshotProfileNamesTheRepositoryPath` (a missing ID returns
   `restore.profile-unresolved` whose message contains
   `.roundfix/baseline/profiles/<id>.json` and not "built-in Baseline
   Profile"), `TestRepositorySnapshotContractsRefusesDisagreement` (two
   snapshots with different tree digests for one required skill return
   `restore.snapshot-conflict` naming the skill and both snapshot IDs, while
   agreeing snapshots and a skill only one lists succeed),
   `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill` (Invariant 8 over
   the embedded catalog), and
   `TestTrailingSetupSkillsComparesARepositoryProfile` (the existing trailing
   test's matching and edited trees, under a repository profile, report
   exactly the edited skill).
4. MUST add `TestRestoreSkillsAcceptsARepositoryProfile` and
   `TestReconcileSkillsLockAcceptsARepositoryProfile`: on a temporary Git
   repository holding a repository profile and a missing `--source-dir`
   path, each returns a finding other than a profile finding, its payload
   names the repository profile, and its `setup` is null.
5. MUST keep `TestTrailingSetupSkillsComparesTheInstalledTreeWithTheSnapshot`
   and every existing restore and reconcile test passing unchanged, which
   proves Invariant 1 for built-in profiles.
6. MUST NOT change `baseline update`, Doctor, the repository profile schema,
   draft binding, the embedded catalog or any derived Baseline file.

## Subtasks

- [ ] Add the loader and the snapshot-agreement merge.
- [ ] Switch the comparison, restore and reconcile to the loader.
- [ ] Print a null setup for a repository profile.
- [ ] Add the loader, comparison, restore and reconcile tests.

## Acceptance Criteria

- [ ] A built-in profile's loaded contracts are identical to today's.
- [ ] A repository profile's comparison reports exactly its trailing skill.
- [ ] Restore and reconcile accept a repository profile and reach their
      source checks.
- [ ] An unresolvable profile is named by its repository path.
- [ ] Disagreeing snapshots are refused by skill and snapshot names.

## Context

- instruction: `docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md`
- instruction: `internal/baseline/custom_profile.go`
- instruction: `internal/baseline/skills_trailing_test.go`
- interface: `internal/baseline/skills_restore.go`
- interface: `internal/baseline/skills_reconcile.go`
- interface: `internal/baseline/skills_trailing.go`
- creates: `internal/baseline/skills_snapshot_profile_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSkillSnapshotProfileKeepsEveryBuiltInProfile|TestSkillSnapshotProfileResolvesARepositoryProfile|TestSkillSnapshotProfileNamesTheRepositoryPath|TestRepositorySnapshotContractsRefusesDisagreement|TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill|TestTrailingSetupSkillsComparesARepositoryProfile|TestRestoreSkillsAcceptsARepositoryProfile|TestReconcileSkillsLockAcceptsARepositoryProfile|TestTrailingSetupSkillsComparesTheInstalledTreeWithTheSnapshot)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSkillSnapshotProfileKeepsEveryBuiltInProfile TestSkillSnapshotProfileResolvesARepositoryProfile TestSkillSnapshotProfileNamesTheRepositoryPath TestRepositorySnapshotContractsRefusesDisagreement TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill TestTrailingSetupSkillsComparesARepositoryProfile TestRestoreSkillsAcceptsARepositoryProfile TestReconcileSkillsLockAcceptsARepositoryProfile TestTrailingSetupSkillsComparesTheInstalledTreeWithTheSnapshot; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the new tests do not exist, so their pass lines are missing and the command fails.

## References

- `_prd.md` → Goals; Core Features 1, 2, 3; Success Metric 5; Success Metric 6
- `_techspec.md` → Interfaces; Invariants 1 to 8; API Contract 2; API Contract 3; Vocabulary Contract; Testing Approach; Build Order 1
- ADR-0241; ADR-0221; ADR-0067
