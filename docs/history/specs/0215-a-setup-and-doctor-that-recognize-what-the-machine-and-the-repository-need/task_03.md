---
task: task_03
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
status: completed
type: backend
complexity: medium
---

# Task 03: A required upstream skill that trails its Setup Snapshot is reported and restored

## Overview

A required upstream skill that matches its lock passes Doctor even when it
differs from the tree the Setup Snapshot pins; on 2026-10-02 eleven skills in
this repository did. This Task adds the comparison of the TechSpec's "Skills
that trail the snapshot" (ADR-0221): Doctor's `skills` line names trailing
skills with `DR-SKILL-TRAILS-SNAPSHOT`, and the managed refresh restores them
through the restore path it already runs. It is verifiable alone with a
temporary repository and the embedded snapshots.

## Requirements

1. MUST add `TrailingSetupSkills` in the new file
   `internal/baseline/skills_trailing.go`, using the restore command's own
   profile loading, installed-tree reading and portable tree digest, and
   returning the sorted names of required external skills whose installed
   tree exists and differs from their snapshot `treeDigest`.
2. MUST add the `trailingSkills` field to `doctorDependencies`, defaulting to
   `TrailingSetupSkills` with the Setup Manifest's profile, and MUST call it
   only for skills that are neither missing nor outdated.
3. MUST report trailing skills on the `skills` line as "Skills that trail the
   snapshot" states: `warn` with `DR-SKILL-TRAILS-SNAPSHOT`, the names and
   `next: roundfix baseline update` when nothing else fails; appended to a
   `failed` line otherwise; and `snapshot comparison unavailable` with an
   unchanged status when the profile does not resolve.
4. MUST keep `skills.RepositoryReadiness.Ready()` and every existing
   `skills` line text unchanged when no skill trails (Invariant 5).
5. MUST add the trailing names to the external drift `roundfix baseline
   update` restores, so its preview lists them and its confirmed run restores
   them through the existing restore request.
6. MUST describe the comparison and its code in
   `docs/user-guide/commands/doctor.md` and the restore of a trailing skill
   in `docs/user-guide/commands/baseline.md`.
7. MUST NOT change the Setup Snapshots, `skills-lock.json`, this
   repository's `.agents/skills/` trees or the restore plan's format.
8. MUST add the tests named in Verification: the comparison in the new file
   `internal/baseline/skills_trailing_test.go`, with a temporary repository
   holding one installed external skill that matches its lock and differs
   from its snapshot tree and one that matches both; and Surface Transcript 4
   and the `baseline update` preview in the new file
   `internal/cli/doctor_trailing_skills_test.go`.

## Subtasks

- [ ] Add the snapshot comparison.
- [ ] Report trailing skills on the `skills` line.
- [ ] Add trailing skills to the managed refresh's external drift.
- [ ] Describe both in the guides.
- [ ] Add the tests.

## Acceptance Criteria

- [ ] A skill that matches its lock and differs from its snapshot tree is
      returned; one that matches both is not; a missing skill is not.
- [ ] Doctor prints `skills: warn` with `DR-SKILL-TRAILS-SNAPSHOT`, the name
      and `roundfix baseline update`, and exits `0` when nothing else fails.
- [ ] The `baseline update` preview lists the trailing skill for restore.
- [ ] `TestThisRepositoryHoldsEveryRequiredExternalSkill` passes unedited.

## Context

- creates: `internal/baseline/skills_trailing.go`
- creates: `internal/baseline/skills_trailing_test.go`
- creates: `internal/cli/doctor_trailing_skills_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/baseline_update_test.go`
- interface: `internal/cli/baseline_update_outdated_skills_test.go`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/baseline.md`
- instruction: `internal/baseline/skills_restore.go`
- instruction: `skills/repository.go`
- instruction: `internal/cli/this_repository_skill_set_test.go`
- instruction: `docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md`
- instruction: `docs/specs/0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need/references/2026-09-30-an-installed-skill-that-matches-its-lock-can-trail-the-snapshot.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTrailingSetupSkillsComparesTheInstalledTreeWithTheSnapshot)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestTrailingSetupSkillsComparesTheInstalledTreeWithTheSnapshot" || { printf 'missing pass: %s\n' TestTrailingSetupSkillsComparesTheInstalledTreeWithTheSnapshot >&2; exit 1; }` — expected: exit 0; before this Task the test does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestDoctorWarnsOnASkillThatTrailsItsSnapshot|TestBaselineUpdatePreviewListsATrailingSkill|TestThisRepositoryHoldsEveryRequiredExternalSkill)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDoctorWarnsOnASkillThatTrailsItsSnapshot TestBaselineUpdatePreviewListsATrailingSkill TestThisRepositoryHoldsEveryRequiredExternalSkill; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the readiness of this repository's skill set still passes beside two new tests, which do not exist before this Task.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- "DR-SKILL-TRAILS-SNAPSHOT" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/doctor.md "DR-SKILL-TRAILS-SNAPSHOT" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- "trails its Setup Snapshot" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/baseline.md "trails its Setup Snapshot" >&2; exit 1; }` — expected: exit 0; before this Task neither guide carries these phrases, so the command fails.

## References

- `_prd.md` → User Story 5; Core Feature 6; Success Metric 5; Acceptance evidence
- `_techspec.md` → Interfaces; Invariant 5; Skills that trail the snapshot; Finding codes; API Contract 4; Surface Transcript 4; Testing Approach 5; Build Order 3
- ADR-0221; ADR-0191; ADR-0189

## Result

Implemented the Task 03 slice for Daemon Verification; status remains
Daemon-owned.

- Added `TrailingSetupSkills` using `loadRestoreProfile`,
  `inspectRestoreTarget`, and `portableRestoreDigest`. It compares required
  external snapshot contracts, ignores missing trees and owned skills, and
  returns sorted, deduplicated trailing names.
- Doctor defaults its new comparison dependency to the Setup Manifest's
  profile. It excludes missing and outdated external skills and avoids
  comparing incompletely inspected trees after a readiness read error.
  Trailing skills add `DR-SKILL-TRAILS-SNAPSHOT` and the managed refresh next
  action: warning alone preserves exit zero, existing failures remain failed,
  and an unavailable comparison preserves the readiness status. An empty
  comparison preserves the existing skills text.
- The managed refresh preview lists trailing skills without acquiring sources
  or changing repository bytes, including when guidance is current. Its
  confirmed skills stage adds trailing names to sorted external drift and
  uses the existing per-skill restore preview and exact confirmation request.
  `--no-skills` skips the comparison. The restore plan format is unchanged.
- Updated the Doctor and Baseline command guides. No Setup Snapshot,
  repository-installed skill tree, lock, Task Graph, other Task file, or
  Repository Skill Set readiness implementation was changed.

### Acceptance evidence from focused implementation checks

1. Lock-intact differences, matching trees, and absent trees:
   `TestTrailingSetupSkillsComparesTheInstalledTreeWithTheSnapshot` uses a
   temporary repository, an embedded profile and snapshot, a snapshot-matching
   repository-held control, and a changed installed skill with its own matching
   lock hash. It proves readiness still holds, returns only the difference,
   excludes a missing tree and owned skill, rejects an unknown profile, and
   checks sorting and deduplication.
2. Doctor warning, name, next action, and exit zero:
   `TestDoctorWarnsOnASkillThatTrailsItsSnapshot` exercises the CLI with a real
   lock-intact installed external skill and the default snapshot comparison.
   It asserts the Surface Transcript 4 skills line, exit zero, empty stderr,
   and unchanged repository bytes. Companion cases cover unchanged matching
   text, missing/outdated filtering, appended failed findings, and unavailable
   comparison with unchanged status.
3. Restore preview:
   `TestBaselineUpdatePreviewListsATrailingSkill` exercises text and JSON
   previews on an adopted temporary repository with an intentionally changed
   external tree and a matching lock hash. It asserts the trailing restore
   entry, approval exit, no mutations, and `--no-skills` behavior.
   `TestBaselineUpdateRestoresTrailingSkillsThroughExistingRequest` proves the
   confirmed stage selects the trailing skill and forwards profile, offline
   source directory, and the restore preview's exact confirmation digest.
4. Existing repository skill readiness:
   `TestThisRepositoryHoldsEveryRequiredExternalSkill` was included unedited
   in the focused CLI checks and passed. `skills.RepositoryReadiness.Ready()`
   remains byte-identical.

Focused commands and observed outcomes:

- Before implementation:
  `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/baseline -run TestTrailingSetupSkills -count=1`
  exited 1 with `undefined: TrailingSetupSkills`, demonstrating the absent API.
- After the final baseline source/test changes:
  `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/baseline -run 'TestTrailingSetupSkills|TestSkillsRestoreOfflinePreviewApplyAndIdempotence' -count=1`
  exited 0 (`ok roundfix/internal/baseline`).
- After the final production CLI changes:
  `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'TestDoctor|TestRunDoctor|TestBaselineUpdate|TestThisRepositoryHolds' -count=1`
  exited 0 (`ok roundfix/internal/cli`, 23.695s).
- After the final CLI test edit:
  `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'TestDoctor.*Snapshot|TestBaselineUpdate.*Trailing|TestThisRepositoryHolds' -count=1`
  exited 0 (`ok roundfix/internal/cli`, 1.692s).
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

Declared Verification commands and repository-wide gates were not run in this
child turn; the Daemon owns declared Verification and settlement. No commit,
push, or Pull Request was made. The shipped Roundfix Skill's corresponding
reference update remains in the later Task's owned slice.

### Verification feedback — attempt 1

Inspected the Daemon diagnostic artifact
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261002T124121Z_c56fb26225ec8dc1/verification/batch-003-attempt-1.log`
and the related detached-process fixture and teardown code. The configured
`make verify-changed` sequence stopped at the existing
`TestRunImplementDetachSurvivesCallerProcessGroupKill`: the operating system
refused the fixture process-group cleanup and zero-signal probe. This is outside
Task 03's snapshot-comparison behavior and its changed paths.

Focused environment diagnosis:

- `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run '^TestRunImplementDetachSurvivesCallerProcessGroupKill$' -count=1`,
  run with approved sandbox escalation for fixture process control, exited 0
  (`ok roundfix/internal/cli`, 4.363s). The same unmodified test and production
  code work with the required process-group permission.
- `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/baseline ./internal/cli -run 'TestTrailingSetupSkills|TestDoctor.*Snapshot|TestBaselineUpdate.*Trailing|TestThisRepositoryHolds' -count=1`
  exited 0 for both packages, refreshing the evidence for all four acceptance
  criteria without executing any declared Verification command.

No source, test, or Verification configuration repair was warranted by this
permission failure. The Daemon's full configured rerun needs an execution
environment that permits signaling and inspecting its own fixture process
groups. Status remains Daemon-owned; no commit, push, or Pull Request was made.

## Carry-forward provenance

- Source Run: `run_20261002T124121Z_c56fb26225ec8dc1`
- Source commit: `b1fea4129c0268b740428d8fd734cc72036e8f6d`
