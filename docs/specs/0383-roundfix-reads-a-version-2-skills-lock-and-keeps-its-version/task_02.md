---
task: task_02
spec: 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
status: completed
type: backend
complexity: low
---

# Task 02: Doctor and the Baseline update name roundfix upgrade for an unsupported lock version

## Overview

Doctor's `skills:` line answered a lock it could not read with `bunx skills
experimental_install && bunx skills update -p -y`, the command that writes
the refused version again (pantheon, 2026-11-12,
[`references/2026-11-12-pantheon-skills-lock-version-2-refused.md`](references/2026-11-12-pantheon-skills-lock-version-2-refused.md),
reported by pantheon). After task_01 a version 2 lock is read; this Task makes
a version Roundfix still refuses point at `roundfix upgrade` in Doctor and in
`roundfix baseline update`, and proves at command level that Doctor, the
bundle install and the update work on a version 2 lock (vortex, 2026-11-12,
[`references/2026-11-12-vortex-skills-lock-version-2-refused.md`](references/2026-11-12-vortex-skills-lock-version-2-refused.md),
reported by vortex), as `_techspec.md` → API Contract 2 defines.

## Requirements

1. MUST make `doctorSkillReadinessResult` (`internal/cli/doctor.go`) set the
   next action to exactly `roundfix upgrade`, and nothing else, when
   `errors.Is(checkErr, skillhash.ErrUnsupportedLockVersion)`; every other
   next action, the status and the detail stay as they are. The release
   plan's `skills` check reads the same line and needs no change.
2. MUST make `roundfix baseline update` (`internal/cli/baseline_update.go`)
   report, at both failures of `baseline.InstalledRetiredSkills`, the next
   action `run roundfix upgrade, then rerun roundfix baseline update` when
   the error matches `skillhash.ErrUnsupportedLockVersion`, and `repair
   skills-lock.json and rerun roundfix baseline update` otherwise.
3. MUST create `internal/cli/skills_lock_version_test.go` with the tests of
   `_techspec.md` → Testing Approach for this Task, each top-level test
   calling `t.Parallel()` first:
   - `TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock`: a temporary project
     with the bundle installed, an external skill `testing-boss` installed
     and required, and the version 2 lock
     `{"version": 2, "skills": {"roundfix": {source, ref "old", sourceType,
     skillPath, a stale computedHash}, "testing-boss": {…}}}` written with
     two-space indentation; the Doctor `skills:` line is
     `skills: warn (<counts>; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: roundfix; next: roundfix skills install --target project)`;
     `roundfix skills install --target project` exits 0 and leaves the lock
     byte for byte as before except that the `roundfix` entry loses its
     `ref` line and its `computedHash` becomes the version 2 digest of the
     installed `.agents/skills/roundfix` folder (computed in the test with
     `skillhash.SumForLockVersion(2, …)` over the folder's files); the lock's
     `version` is still 2; and the Doctor line becomes
     `skills: ok (<counts>)`.
   - `TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`: the
     same project with the lock `{"version":3,"skills":{"testing-boss":{}}}`
     prints exactly the line of `_techspec.md` → API Contract 2, with its
     temporary path.
   - `TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`:
     an adopted Baseline repository with its required external skills
     installed and its lock's `version` set to 2 (subtest `version 2`: exit
     0, no stderr, no next action) or 3 (subtest `version 3`: exit
     `exitRunFailed`, category `execution`, next action `run roundfix
     upgrade, then rerun roundfix baseline update`, stderr holding
     `skills-lock.json version 3 is not supported; this Roundfix reads
     versions 1 and 2.`).
4. MUST reach no network, run no skills CLI, start no Agent session, and
   write only inside temporary directories in every test.

## Subtasks

- [x] Name `roundfix upgrade` in Doctor's `skills:` next action for an unsupported lock version.
- [x] Name the upgrade in both lock failures of `roundfix baseline update`.
- [x] Add the command-level tests for a version 2 and a version 3 lock.

## Acceptance Criteria

- [x] On a version 2 lock Doctor reads the lock, the bundle install repairs a stale Roundfix entry with the version 2 hash and keeps version 2, and the update succeeds.
- [x] On a version 3 lock Doctor's only next action is `roundfix upgrade` and the update names the upgrade; no next action names a skills command for it.

## Context

- interface: `internal/cli/doctor.go`
- interface: `internal/cli/baseline_update.go`
- creates: `internal/cli/skills_lock_version_test.go`
- instruction: `internal/cli/doctor_skills_lock_test.go`
- instruction: `internal/cli/owned_skills_lock_test.go`
- instruction: `internal/cli/baseline_update_retired_skills_test.go`
- instruction: `internal/cli/baseline_update_missing_skills_test.go`
- instruction: `internal/cli/releaseplan_checks.go`
- instruction: `internal/skillhash/hash.go`
- instruction: `docs/agents/cli.md`
- instruction: `docs/agents/go.md`

## Consumers

- Doctor's `skills:` next action and the update's lock next action, read by `roundfix doctor`, the release plan's `skills` check and `roundfix baseline update` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock|TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion|TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion/version_2 TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion/version_3; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the three tests do not exist, so their pass lines are missing; afterwards every named test and subtest passes.
- `grep -qF -- 'skillhash.ErrUnsupportedLockVersion' internal/cli/doctor.go && grep -qF -- 'skillhash.ErrUnsupportedLockVersion' internal/cli/baseline_update.go && go test -count=1 -timeout 30m -run 'Doctor|Skills|Lock|Retired|OwnedSkill|MissingRequired|ReleasePlan' ./internal/cli && GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task neither file reads the refusal, so the first check fails; afterwards the Doctor, skills, lock, retired-skill and release-plan tests of the package pass and it vets for Windows.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Feature 3; Success Metric 1; Success Metric 3
- `_techspec.md` → API Contract 2; Testing Approach; Build Order 2
- ADR-0393; ADR-0359

## Result

Skills loaded: implement-task, coding-guidelines

### Implementation

- `internal/cli/doctor.go`: imported `roundfix/internal/skillhash`. In
  `doctorSkillReadinessResult`, after the detail is built and before the
  owned/external next-action selection, added an early return that sets
  `result.NextAction = "roundfix upgrade"` when
  `errors.Is(checkErr, skillhash.ErrUnsupportedLockVersion)`. The status,
  detail and `Err` are unchanged; every other next action falls through to the
  existing logic. Because the release plan's `skills` check calls the same
  `repositorySkillsCheck`, it needs no change.
- `internal/cli/baseline_update.go`: imported `roundfix/internal/skillhash`.
  Added `baselineUpdateSkillsLockNextAction(err error)` returning
  `run roundfix upgrade, then rerun roundfix baseline update` when the error
  matches `skillhash.ErrUnsupportedLockVersion`, and
  `repair skills-lock.json and rerun roundfix baseline update` otherwise. Both
  `baseline.InstalledRetiredSkills` failure sites (the pre-approval listing and
  the post-skills-stage re-inspection) now use it.
- `internal/cli/skills_lock_version_test.go` (new): the three command-level
  tests of `_techspec.md` → Testing Approach. Each top-level test calls
  `t.Parallel()` first. `TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock`
  installs the bundle plus an external `testing-boss` folder, writes the
  two-space-indented version 2 lock, checks the exact `warn …
  DR-SKILL-LOCK-STALE … pins: roundfix; next: roundfix skills install --target
  project` line, runs `roundfix skills install --target project` (exit 0),
  asserts the lock is byte-for-byte the original except the `roundfix` entry
  losing its `ref` line and gaining the version 2 digest computed in the test
  with `skillhash.SumForLockVersion(2, …)`, confirms `"version": 2` is kept,
  and re-checks the `ok` line.
  `TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion` asserts API
  Contract 2's line verbatim with its temporary path.
  `TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`
  drives an adopted Baseline repository with its required external skills
  installed and the lock's `version` set to 2 (`exitOK`, no stderr, no next
  action, bytes unchanged) or 3 (`exitRunFailed`, category `execution`, the
  upgrade next action, stderr carrying the version-3 refusal text). No test
  reaches the network, runs the skills CLI, starts an Agent session, or writes
  outside its temporary directories.

### Focused checks (not the declared Verification)

- `go build ./internal/cli/` → exit 0.
- `go vet ./internal/cli/` → exit 0.
- `go test -count=1 -timeout 10m -v -run '^(TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock|TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion|TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion)$' ./internal/cli` → `ok`, with `--- PASS` for all three and for `version_2`/`version_3`.
- `go test -count=1 -timeout 20m -run 'Doctor|Skills|Lock|Retired|OwnedSkill|MissingRequired|ReleasePlan' ./internal/cli` → `ok`.
- `GOOS=windows go vet ./internal/cli` → exit 0.
- `gofmt -l` on the two changed files and the new test → no output.

### Follow-ups

- None. task_03 owns the guide/reference/glossary text and the Behavior
  Surface Record.

## Carry-forward provenance

- Source Run: `run_20261112T212907Z_a8dd596719159f5e`
- Source commit: `ddd371c956764149756e4037a545813c23609161`
