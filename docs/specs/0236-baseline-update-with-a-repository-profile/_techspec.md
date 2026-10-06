---
spec: 0236-baseline-update-with-a-repository-profile
prd: _prd.md
created: 2026-10-06
---

# Baseline update with a repository profile — Technical Spec

## Executive Summary

The skill readers of `internal/baseline` — the snapshot comparison, skill
restore and lock reconciliation — share one loader, `loadRestoreProfile`,
which resolves only built-in profiles. A new loader,
`loadSkillSnapshotProfile`, keeps that path for a built-in ID and otherwise
resolves the profile through `ResolveProfile`, the resolver `update`, plan and
`profile validate` already use. A repository profile then requires the
skills of its own modules, and each external contract comes from the
embedded Setup Snapshots that list the skill, which must agree (ADR-0241).
The three readers switch to the new loader; `baseline update` and Doctor
need no code change because they already pass the repository root and the
manifest's profile ID. The primary trade-off is a contract source that spans
every embedded snapshot instead of one named snapshot: it works for any
repository profile without a schema change, at the cost of a catalog
invariant (all snapshots agree on shared skills) that a test now holds and a
refusal enforces.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  only new names are the finding codes `restore.profile-unresolved` and
  `restore.snapshot-conflict` in the existing `restore.` family. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — profile resolution and the
  comparison read local files and the embedded catalog only; restore keeps
  its existing source acquisition, and every test uses temporary
  repositories, `--source-dir` values and no network. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0241 (this Spec): "A repository
  profile requires the skills its own modules require". Repository profiles
  stay repository-owned and catalog-bound, ADR-0067: "may compose only the
  modules, decisions, Repository Capabilities, and templates". The
  comparison keeps its rule, ADR-0221: "A skill that differs is reported as
  trailing the snapshot". ADR-0204's composed snapshot already refuses two
  different entries under one skill name, ADR-0204: "would merge two
  different entries under one skill name". ADR-0219 is not changed: draft
  binding still chooses a plan's source. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_02 changes the Governed Paths
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`; task_01
  changes none. Express maintainer authorization: the standing grant of
  2026-09-30, "considere autorizado a ajustar todas as skills se
  necessário", and "Concedo" of 2026-10-06; bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `skills/roundfix/SKILL.md`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0236-baseline-update-with-a-repository-profile/_authorization.md`.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| Skill snapshot profile | `loadSkillSnapshotProfile` beside `loadRestoreProfile` in `internal/baseline/skills_restore.go` | Built-in path unchanged; repository path through `ResolveProfile` and the agreeing snapshots |
| Snapshot comparison | `TrailingSetupSkills` in `internal/baseline/skills_trailing.go` | Uses the new loader |
| Skill restore | `RestoreSkills` in `internal/baseline/skills_restore.go` | Uses the new loader against the repository's Git top level; `setup` null for a repository profile |
| Lock reconciliation | `ReconcileSkillsLock` in `internal/baseline/skills_reconcile.go` | Same as restore |
| `baseline update`, Doctor | `internal/cli/baseline_update.go`, `internal/cli/doctor.go` | No code change; covered by new regression tests |
| Records | `docs/user-guide/commands/baseline.md`, `doctor.md`, the Roundfix Skill's baseline reference | Repository profiles in restore, reconcile, update and Doctor |

## Implementation Design

### Interfaces

```go
// internal/baseline/skills_restore.go, beside loadRestoreProfile
func loadSkillSnapshotProfile(repoRoot, profileID string, catalog *Catalog) (restoreProfile, error)

type setupSnapshotSkills struct {
	ID     string
	Skills map[string]restoreSkillContract // external (github) entries only
}

func repositorySnapshotContracts(
	profileID string,
	required map[string]struct{},
	snapshots []setupSnapshotSkills, // in Catalog.SetupIDs order
) (map[string]restoreSkillContract, error)
```

```text
1. A built-in ID returns exactly what loadRestoreProfile returns for it, so built-in restore, reconcile and comparison results and Plan Digests are unchanged.
2. Any other ID resolves through ResolveProfile(repoRoot, id, catalog). Failure returns a SkillsRestoreError, category invalid, code restore.profile-unresolved, message: Baseline Profile "<id>" is neither built-in nor resolvable at .roundfix/baseline/profiles/<id>.json. Its action: Run roundfix baseline profile validate <id> to see why, restore that file, or choose a built-in profile id. The resolution error stays wrapped; the message never contains "built-in Baseline Profile".
3. For a repository profile, RequiredSkills is the union of requiredSkills over the profile's own Modules.
4. Skills holds, for each required skill, the external entry of the embedded snapshots that list it as external. A skill no snapshot lists as external is absent, as a repository-held skill is today. Two different entries (provider, repository, ref, path or tree digest) for one skill return code restore.snapshot-conflict naming the skill and the disagreeing snapshot IDs in order; the first conflicting skill in lexical order is named.
5. A repository profile's restoreProfile has ID = the profile ID and Setup = "". Restore and reconcile payloads print "setup": null for it, on success and failure.
6. TrailingSetupSkills, RestoreSkills and ReconcileSkillsLock call loadSkillSnapshotProfile. RestoreSkills and ReconcileSkillsLock resolve a non-built-in ID against the Git top level of request.Repository; their profile-ID consistency checks keep comparing IDs, and their messages no longer say "built-in".
7. TrailingSetupSkills keeps wrapping failures as "load profile for snapshot comparison: <message>", so Doctor's "snapshot comparison unavailable: " detail now names the searched path.
8. Every embedded Setup Snapshot agrees on every skill it lists as external with every other snapshot listing it; a test asserts this over the embedded catalog.
```

### Data Models

None. The repository profile schema, the Setup Manifest and the embedded
catalog are unchanged. The restore and reconcile payloads keep their schema;
`setup` is already nullable and is null for a repository profile.

### API Contracts

1. API Contract: `roundfix baseline update` on a repository adopted with a
   repository profile runs the snapshot comparison over that profile; it
   exits 0 with `state: current` when guidance and skills match (Surface
   Transcript 1), exits 3 with `plan_ready` and the trailing skill under
   `skills.drifted` when a locked external skill differs from its snapshot,
   and its `--yes` skills stage installs the owned skills and restores a
   trailing skill with the repository profile's ID.
2. API Contract: `roundfix baseline skills restore` and `roundfix baseline
   skills reconcile` accept a repository profile ID; restorable skills are
   the external skills its modules require, and the payload's `setup` is
   null (Surface Transcript 3).
3. API Contract: an ID that resolves neither as built-in nor from the
   repository yields `restore.profile-unresolved` with the message and action
   of Invariant 2, exit 2 (Surface Transcript 2).
4. API Contract: `roundfix doctor`'s `skills:` line compares a repository
   profile's skills with their snapshot, adding `DR-SKILL-TRAILS-SNAPSHOT`
   for a trailing skill; for an unresolvable profile it says `snapshot
   comparison unavailable` and names the searched path. Doctor's exit code
   depends on the machine, so this contract is proven by tests rather than
   a transcript.

### Surface Transcripts

1. Surface Transcript: a repository adopted with a repository profile, with
   no locked external skill installed, is current.

   ```transcript
   $ roundfix baseline update --repo <repository> --format text
   stdout:
   Baseline update: current
   Result: the repository already matches the current Baseline catalog
   ...
   stderr:
   exit: 0
   ```

2. Surface Transcript: an unresolvable profile names the repository path.

   ```transcript
   $ roundfix baseline skills restore --profile missing-profile --repo <repository> --format text
   stdout:
   Baseline skills restore: blocked
   restore.profile-unresolved: Baseline Profile "missing-profile" is neither built-in nor resolvable at .roundfix/baseline/profiles/missing-profile.json.
   Next action: Run roundfix baseline profile validate missing-profile to see why, restore that file, or choose a built-in profile id.
   stderr:
   roundfix: baseline skills restore failed: Baseline Profile "missing-profile" is neither built-in nor resolvable at .roundfix/baseline/profiles/missing-profile.json.
   exit: 2
   ```

3. Surface Transcript: restore accepts a repository profile and reaches its
   source checks.

   ```transcript
   $ roundfix baseline skills restore --profile <repository-profile> --skill context7-cli --source-dir <missing-directory> --repo <repository> --format text
   stdout:
   Baseline skills restore: blocked
   restore.source-dir-invalid: Offline Git object store is not a directory: <path>.
   Next action: Pass an existing Git checkout or bare object store to --source-dir.
   stderr:
   roundfix: baseline skills restore failed: Offline Git object store is not a directory: <path>.
   exit: 2
   ```

## Coverage Map

- Goal 1 → `loadSkillSnapshotProfile`, `TrailingSetupSkills` (Invariants 2 to 4, 6; API Contract 1; Surface Transcript 1)
- Goal 2 → the unchanged owned install in `runBaselineUpdateSkillsStageWith`, which now completes (API Contract 1)
- Goal 3 → `RestoreSkills`, `ReconcileSkillsLock`, Doctor through `TrailingSetupSkills` (Invariants 6, 7; API Contracts 2, 4)
- Goal 4 → Invariant 2; API Contract 3; Surface Transcript 2
- Core Feature 1 → `loadSkillSnapshotProfile` (Invariants 1, 6)
- Core Feature 2 → `repositorySnapshotContracts` (Invariants 3, 4, 8)
- Core Feature 3 → Invariants 2, 7
- Core Feature 4 → the CLI regression tests of task_02
- Core Feature 5 → the baseline and doctor references and the Roundfix Skill's baseline reference
- Success Metric 1 → `TestBaselineUpdateReachesCurrentWithARepositoryProfile`
- Success Metric 2 → `TestBaselineUpdatePreviewListsATrailingSkillWithARepositoryProfile`
- Success Metric 3 → `TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile`
- Success Metric 4 → `TestDoctorComparesARepositoryProfileWithItsSnapshot`
- Success Metric 5 → task_01's restore tests; Surface Transcripts 2, 3
- Success Metric 6 → the repository Verification at each Task's settlement

## Integration Points

- Git, through the existing repository identity inspection that restore and
  reconcile already run.
- The upstream skills source, reached only by a real restore without
  `--source-dir`; never in tests.

## Testing Approach

- `internal/baseline/skills_snapshot_profile_test.go` tests the loader
  directly: every built-in ID equals `loadRestoreProfile`; a repository
  profile written by `InitCustomProfile` from `go-cli-tui` has the module
  union as required skills and the `go` snapshot's contracts; a missing
  profile returns `restore.profile-unresolved` naming the path; the pure
  `repositorySnapshotContracts` refuses two disagreeing snapshots; and the
  embedded catalog satisfies Invariant 8.
- The comparison is tested as the existing trailing test does, with the
  repository's own `golang-cli` tree as the matching control and an edited
  tree as the trailing one, under a repository profile.
- `RestoreSkills` and `ReconcileSkillsLock` are tested on a temporary Git
  repository with a repository profile and a missing `--source-dir`: the
  finding is the source-directory refusal and the payload names the profile
  with a null `setup`.
- `internal/cli/baseline_update_repository_profile_test.go` adopts a
  temporary repository through the real `baseline profile init`,
  `baseline plan --profile <repository-id>` and `baseline apply`, copies this
  repository's skill set and lock as `TestBaselineUpdatePreviewListsATrailingSkill`
  does, and runs the real `baseline update` and Doctor's skills check. The
  `--yes` skills stage runs through `runBaselineUpdateSkillsStageWith` with
  the real install and comparison and a fake restore, so no test fetches a
  source.

## Build Order

1. Add `loadSkillSnapshotProfile` and `repositorySnapshotContracts`, switch
   the comparison, restore and reconcile to it, and test them in the
   baseline package.
2. Add the `baseline update`, `--yes` skills stage and Doctor regression
   tests for a repository profile, and describe repository profiles in the
   baseline and doctor references and the Roundfix Skill's baseline
   reference (depends on: 1).
3. Final QA gate (depends on: 1, 2).

## Risks & Considerations

- A future upstream refresh that moves one snapshot to a new commit before
  another would make shared skills disagree; Invariant 8's test fails first,
  and the restore refusal names the skill, so no repository silently
  compares against the wrong tree.
- A repository profile's restorable set is narrower than a built-in
  profile's (required skills only); a restore without `--skill` installs no
  other stack's skills, which is the intended difference.
- Specs 0235 and 0237 may also raise the Roundfix Skill's version; the record
  command picks a free version and the conflict resolves at merge
  (ADR-0233).
- `baseline profile init --from go-cli-tui` writes a profile that is not a
  valid adaptation of its own source (the `autonomous-work` module requires
  `runtime.backend`); this Spec does not depend on adaptation and leaves
  that observation to its own entry.

## Vocabulary Contract

- emits: `internal/baseline/skills_restore.go`
  pattern: `restore.profile-unresolved`
  documented-in: `docs/user-guide/commands/baseline.md`

No glossary term is adopted; "Setup Snapshot" and "Baseline Profile" keep
their glossary meanings.

## Decisions

- A repository profile's external contracts come from the agreeing embedded
  snapshots, not from a bound built-in profile. See ADR-0241.
- Compare a repository profile rather than skip the comparison with a
  warning.
- Leave `baseline update` and Doctor code unchanged; the loader fixes both.
