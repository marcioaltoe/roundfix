---
task: task_01
spec: 0354-an-external-skill-is-checked-only-for-presence
status: pending
type: backend
complexity: medium
---

# Task 01: Doctor and the managed refresh stop comparing external skills with the Setup Snapshot

## Overview

Doctor appends `DR-SKILL-DIFFERS-FROM-SNAPSHOT` and `roundfix baseline update`
lists `skills.differing` for an external skill whose copy differs from its
Setup Snapshot, although the copy may be the newer one its repository chose:
Fluxus saw it for `systematic-debugging` and `ui-ux-pro-max` on 2026-11-05.
This Task removes that comparison from both commands, the `skills.differing`
field and its text lines, and `baseline.SetupSnapshotDifferences` with its
tests, as ADR-0364 records, and writes the rule into the Doctor and Baseline
guides and the README. This repository keeps its own snapshot rule through
the vendored catalog test. It answers the fluxus report of
2026-11-05, adopted as
[references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md](references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md).

## Requirements

1. MUST remove from `internal/cli/doctor.go` the `snapshotDifferences`
   member of `doctorDependencies` and its default, `manifestSnapshotDifferences`,
   `snapshotComparableSkills`, and the block of `repositorySkillsCheck` that
   appends `; snapshot comparison unavailable: ` or `;
   DR-SKILL-DIFFERS-FROM-SNAPSHOT: `, so the `skills:` line follows API
   Contract 2 and keeps the `DR-SKILL-LOCK-STALE` suffix and its next action
   byte for byte.
2. MUST remove from `internal/cli/baseline_update.go` the `Differing` member
   of `baselineUpdateSkillsResult` (JSON `differing`) and every place that
   fills or normalizes it, the `snapshotDifferences` member of
   `baselineUpdateSkillsDependencies` and its default,
   `baselineUpdateSnapshotDifferences`, the preview's call of
   `baseline.SetupSnapshotDifferences`, the skills stage's comparison, and
   the `Skills differing: <n>` and `- differs <name>: <reason>` text lines,
   as API Contract 3 states. The skills stage keeps restoring what
   `baselineUpdateExternalDrift` returns.
3. MUST delete `internal/baseline/skills_trailing.go`,
   `internal/baseline/skills_trailing_test.go`,
   `internal/baseline/repository_skills_snapshot_test.go` and
   `internal/cli/doctor_trailing_skills_test.go`; delete
   `TestSetupSnapshotDifferencesCompareARepositoryProfile` and the
   `SetupSnapshotDifferences` assertion of
   `TestSkillSnapshotProfileNamesTheRepositoryPath` in
   `internal/baseline/skills_snapshot_profile_test.go`; and leave no
   reference to `SetupSnapshotDifferences`, `snapshotComparableSkills` or
   `snapshotDifferences` under `internal` or `skills`.
4. MUST change `TestDoctorOwnedLockWarningKeepsSnapshotDetailAndNextActions`
   in `internal/cli/doctor_skills_lock_test.go` into
   `TestDoctorOwnedLockWarningKeepsReadinessDetailAndNextActions`, with no
   fake snapshot comparison and the detail checked for `;
   DR-SKILL-LOCK-STALE:`, its other assertions unchanged.
5. MUST rewrite in `internal/cli/repository_profile_docscontract_test.go`:
   - `TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile` as
     `TestBaselineUpdatePreviewIgnoresADifferenceWithARepositoryProfile`: the
     preview of the edited `testing-boss` is `current`, exit `0`, with no
     `"differing"` in stdout and no drifted skill;
   - `TestDoctorComparesARepositoryProfileWithItsSnapshot` as
     `TestDoctorChecksARepositoryProfileSkillOnlyForPresence`: the `skills`
     line is `ok` with no next action and no `DR-SKILL-DIFFERS-FROM-SNAPSHOT`
     or `snapshot comparison unavailable`, also after the manifest names the
     unresolvable profile `missing-profile`, and Doctor writes nothing;
   - `TestBaselineUpdatePreviewListsASetupSnapshotDifference` as
     `TestBaselineUpdatePreviewIgnoresASetupSnapshotDifference`: the JSON
     preview, the text preview, the `--no-skills` preview and the `--yes`
     run carry no `"differing"`, no `Skills differing` line and no restore,
     `--yes` leaves the skill's bytes and hash, and the combined subtest
     expects the outdated and retired lines with no difference line;
   - `TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile`
     without the `snapshotDifferences` dependency and the `Differing`
     assertion.
6. MUST add `internal/cli/external_skill_presence_test.go` with the helper
   `setDoctorTestManifestProfile(t, root, profile)` and two tests, each
   calling `t.Parallel()` first:
   - `TestDoctorChecksAnExternalSkillOnlyForPresence`: in
     `writeDoctorReadyRepositoryFixture` with profile `go-cli-tui` and the
     one external skill `testing-boss`, run `doctor` through `runCLI` with
     `skills.CheckRepositoryWithExternal`; rows `differs from its Setup
     Snapshot` (the fixture as written) prints the line `skills: ok
     (<n> required: <owned> Roundfix-owned, 1 external)` with `<owned>` read
     from `len(skills.Names())` and exits `0`, and `missing` (folder removed)
     prints `skills: failed (missing: testing-boss; next: bunx skills add
     marcioaltoe/skills@testing-boss)` and exits `1`; no row prints
     `DR-SKILL-DIFFERS-FROM-SNAPSHOT` or changes a repository byte;
   - `TestBaselineUpdateNeverRestoresAPresentExternalSkill`: through
     `runBaselineUpdateSkillsStageWith` with the real
     `skills.CheckRepositoryWithExternal` over the same fixture, rows `present
     copy that differs from its Setup Snapshot` (restores nothing, keeps the
     copy's bytes, and its JSON encoding has no `"differing"`) and `missing
     copy` (restores `testing-boss`).
7. MUST make in `docs/user-guide/commands/doctor.md`,
   `docs/user-guide/commands/baseline.md` and `README.md` every change API
   Contract 7 names for them, word for word, so the two guides carry the
   rule sentence and none of the three names `DR-SKILL-DIFFERS-FROM-SNAPSHOT`
   or `skills.differing`; the guides document the whole Spec in the change
   that first moves their surface.
8. MUST re-record the coverage record with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
   after the last test removal or rename, then re-record the Behavior Surface
   Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last guide edit, and never hand-edit
   either record.

## Subtasks

- [ ] Remove the snapshot comparison from Doctor and the managed refresh, and `skills.differing` with it.
- [ ] Delete `SetupSnapshotDifferences` and the tests of the comparison.
- [ ] Rewrite the repository-profile and owned-lock tests and add the two new tests.
- [ ] Write the rule into the Doctor and Baseline guides and the README.
- [ ] Re-record the coverage record, then the Behavior Surface Record.

## Acceptance Criteria

- [ ] Doctor's `skills:` line names no Setup Snapshot difference for an external skill that differs from its snapshot.
- [ ] `baseline update` prints no `Skills differing:` line and no `differing` member, and restores only what readiness reports drifted.
- [ ] No code under `internal` or `skills` names `SetupSnapshotDifferences`.
- [ ] The Doctor and Baseline guides carry the rule sentence, and neither they nor the README name `DR-SKILL-DIFFERS-FROM-SNAPSHOT` or `skills.differing`.
- [ ] The CLI and Baseline suites, with and without `docscontract`, the coverage equivalence test and the documentation contracts pass.

## Context

- instruction: `docs/adr/0364-a-required-external-skill-is-checked-only-for-presence.md`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/baseline/skills_trailing.go`
- interface: `internal/baseline/skills_trailing_test.go`
- interface: `internal/baseline/repository_skills_snapshot_test.go`
- interface: `internal/baseline/skills_snapshot_profile_test.go`
- interface: `internal/cli/doctor_trailing_skills_test.go`
- interface: `internal/cli/doctor_skills_lock_test.go`
- interface: `internal/cli/repository_profile_docscontract_test.go`
- creates: `internal/cli/external_skill_presence_test.go`
- interface: `docs/references/coverage-record.json`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `README.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- Doctor's `skills:` line — consumers: `./internal/cli`
- The `baseline update` skills result without `differing` — consumers: `./internal/cli`
- The removal of `baseline.SetupSnapshotDifferences` — consumers: `./internal/baseline`, `./internal/cli`
- The coverage record — consumers: `./internal/spec`
- The Doctor and Baseline guide fingerprints in the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDoctorChecksAnExternalSkillOnlyForPresence|TestBaselineUpdateNeverRestoresAPresentExternalSkill)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDoctorChecksAnExternalSkillOnlyForPresence TestBaselineUpdateNeverRestoresAPresentExternalSkill 'TestDoctorChecksAnExternalSkillOnlyForPresence/differs_from_its_Setup_Snapshot' 'TestBaselineUpdateNeverRestoresAPresentExternalSkill/present_copy_that_differs_from_its_Setup_Snapshot'; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; ! grep -rn 'SetupSnapshotDifferences\|snapshotComparableSkills\|snapshotDifferences' internal skills --include='*.go' || { echo 'snapshot comparison still referenced' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/cli ./internal/baseline && go test -count=1 -timeout 30m -tags docscontract ./internal/cli ./internal/baseline && go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec` — expected: exit 0. Before this Task neither new test exists, so the first pass line is missing; afterwards Doctor names no snapshot difference, the skills stage restores only the missing copy, no code names the comparison, and the CLI and Baseline suites, both tags, and the coverage equivalence pass.

- `rule='checks a required external skill only for presence: an entry in .skills-lock\.json. and a folder under .\.agents/skills.'; for f in docs/user-guide/commands/doctor.md docs/user-guide/commands/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -q -- "$rule" || { echo "missing rule: $f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < README.md | grep -qF -- 'its content is not compared' || { echo 'missing README rule' >&2; exit 1; }; for f in docs/user-guide/commands/doctor.md docs/user-guide/commands/baseline.md README.md; do ! grep -q -- 'DR-SKILL-DIFFERS-FROM-SNAPSHOT\|skills\.differing' "$f" || { echo "old comparison still named: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither guide carries the rule sentence, so the first check fails; afterwards both do, none of the three names the retired code or field, and the documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 2, 3 and 5; Success Metric 1; Success Metric 3; Success Metric 4; Success Metric 7
- `_techspec.md` → Current behavior 3, 4 and 6; API Contracts 2, 3 and 7; Invariants 2 and 3; Declared breaks; Testing Approach; Build Order 1
