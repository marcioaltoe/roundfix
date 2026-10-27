---
task: task_02
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
status: completed
type: backend
complexity: medium
---

# Task 02: The bundle install rewrites the lock hash of each installed owned skill

## Overview

After `roundfix baseline update` or `roundfix skills install --target
project` installs the Roundfix-owned skills the binary carries, each existing
`skills-lock.json` entry that names one of them gets the `computedHash` of
the installed folder, and nothing else in the lock changes. This answers the
adopted Backlog Entry "Baseline update leaves the skills lock hash of the
catalog" (Quaestor, 2026-10-26, repeated on v0.125.0 on 2026-10-27).

## Requirements

1. MUST add `RecordInstalledOwnedSkills(ctx, root, names)` in a new
   `internal/baseline/owned_skills_lock.go`, implementing `_techspec.md` →
   API Contract 1: no lock means no write and no new file; the lock is
   loaded through `loadSkillsLock`, so a symbolic link, a non-regular file
   or a malformed document returns an error naming `skills-lock.json` with
   the bytes unchanged; only the `computedHash` field of an existing object
   entry whose `.agents/skills/<name>` folder exists is set, to
   `externalSkillsLockDigest` of the files `inspectRestoreTarget` reads;
   every other field, value and the key order are kept; the lock is written
   only when an entry changed, through a temporary file in the repository
   root renamed over it, keeping its file mode; the rewritten names return
   sorted.
2. MUST add the `recordLock` dependency and `LockUpdated`
   (`json:"lockUpdated,omitempty"`) to `internal/cli/baseline_update.go` and
   call it right after a successful install in
   `runBaselineUpdateSkillsStageWith`, implementing API Contract 2: an error
   fails the stage with `record installed Roundfix skills in
   skills-lock.json: <error>`; text prints `Skills lock updated: <count>`
   and one `- lock updated <name>` line per name after the installed skills,
   only when the list is not empty. `runBaselineUpdateSkillsStage` wires
   `baseline.RecordInstalledOwnedSkills`. The preview and the Plan Digest
   are unchanged.
3. MUST call the function from `runSkillsInstall` (`internal/cli/cli.go`)
   for `--target project` without `--dir`, implementing API Contract 3:
   `Updated skills-lock.json computedHash: <name>` per rewritten name on
   stdout; an error prints `roundfix: skills install failed: <error>` and
   exits `2`. Other targets and `--dir` never read a lock.
4. MUST give every existing test that builds `baselineUpdateSkillsDependencies`
   in `internal/cli/baseline_update_test.go` a `recordLock` fake, keeping
   each existing assertion.
5. MUST add the tests `_techspec.md` → Testing Approach names for task_02:
   `TestRecordInstalledOwnedSkillsRewritesOnlyTheHashOfAnInstalledOwnedSkill`,
   `TestRecordInstalledOwnedSkillsLeavesAMatchingLockUntouched`,
   `TestRecordInstalledOwnedSkillsAddsNoEntryAndNoLock` and
   `TestRecordInstalledOwnedSkillsRefusesAnUnsafeLock` in
   `internal/baseline/owned_skills_lock_test.go`;
   `TestSkillsInstallRewritesTheOwnedSkillLockHash`,
   `TestBaselineUpdateSkillsStageReportsTheRewrittenLockEntries` and
   `TestBaselineUpdateSkillsStageFailsOnAnUnsafeLock` in
   `internal/cli/owned_skills_lock_test.go`. Each expected hash is computed
   from the installed folder in the test, never written as a literal; the
   stale hash in the fixture lock is any other valid 64-hex value. The
   rewrite test's lock carries an extra field in the owned entry and one
   external entry, so a rewrite that rebuilds the entry or touches another
   entry is observable.

## Subtasks

- [ ] Characterize: a test that installs the bundle over a lock with a stale owned entry and fails on the current tree.
- [ ] Add the Baseline function and its tests.
- [ ] Wire it into the update's skills stage and `skills install`.
- [ ] Add the CLI tests and update the stage fakes.

## Acceptance Criteria

- [ ] After either command, an owned-skill lock entry's `computedHash` equals its installed folder's digest, and every other lock byte is unchanged.
- [ ] A second run writes nothing and reports no name.
- [ ] No lock is created, no entry is added, and an unsafe lock fails both commands without a write.

## Context

- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `docs/user-guide/commands/baseline.md`
- instruction: `docs/user-guide/commands/skills.md`
- instruction: `.agents/skills/roundfix/references/baseline.md`
- instruction: `internal/baseline/skills_restore_lock.go`
- instruction: `internal/baseline/skills_restore.go`
- instruction: `skills/skills.go`
- instruction: `skills/repository.go`
- creates: `internal/baseline/owned_skills_lock.go`
- creates: `internal/baseline/owned_skills_lock_test.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/baseline_update_test.go`
- creates: `internal/cli/owned_skills_lock_test.go`

## Consumers

- `RecordInstalledOwnedSkills` — consumers: `./internal/cli`
- `baselineUpdateSkillsResult` JSON and text — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRecordInstalledOwnedSkills.*|TestSkillsInstallRewritesTheOwnedSkillLockHash|TestBaselineUpdateSkillsStage.*|TestBaselineUpdateSkillStage.*|TestBaselineUpdatePassesOfflineSourceToSkillStage|TestBaselineUpdateSkipsSkillStageAndPreservesSkillDirectory)$' ./internal/baseline ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRecordInstalledOwnedSkillsRewritesOnlyTheHashOfAnInstalledOwnedSkill TestRecordInstalledOwnedSkillsLeavesAMatchingLockUntouched TestRecordInstalledOwnedSkillsAddsNoEntryAndNoLock TestRecordInstalledOwnedSkillsRefusesAnUnsafeLock TestSkillsInstallRewritesTheOwnedSkillLockHash TestBaselineUpdateSkillsStageReportsTheRewrittenLockEntries TestBaselineUpdateSkillsStageFailsOnAnUnsafeLock TestBaselineUpdateSkillsStageUsesPreviewThenConfirmation TestBaselineUpdateSkillStageReportsInstalledAndRestoredSkills; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the seven new tests do not exist, so the first of them is missing from the passes. Afterwards the new tests and the existing skills-stage tests that consume the changed dependencies pass.

## References

- `_prd.md` → Goals; Core Feature 1; User Experience; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior 1-3; Interfaces; API Contract 1; API Contract 2; API Contract 3; Invariants 1-2; Testing Approach; Build Order 2
- ADR-0189

## Result

Implemented the owned bundle's lock recording in both project install paths.
`RecordInstalledOwnedSkills` loads through the existing safe lock loader,
reads installed folders through `inspectRestoreTarget`, and computes their
external skills digest. JSON value offsets preserve every byte outside the
changed `computedHash` values, including whitespace, escapes, numeric spelling,
extra fields, external entries and key order. Changed documents are written
through a repository-local temporary file and renamed with the original mode.
No lock or entry is created; unchanged locks are not written. Rewritten names
are deduplicated and sorted.

The Baseline skills stage records the lock after installation, exposes optional
`lockUpdated` JSON and the specified text lines, and propagates recording
failures. Project `skills install` reports rewritten names and returns exit 2
on unsafe locks. Explicit `--dir` bypasses lock recording. Preview planning and
Plan Digest code are unchanged.

Characterization: before the implementation,
`go test ./internal/cli -run '^TestSkillsInstallRewritesTheOwnedSkillLockHash$' -count=1`
failed because the installed folder's computed digest differed from the stale
hash retained in the lock.

Focused checks after implementation (all exit 0):

- `go test ./internal/baseline -run '^TestRecordInstalledOwnedSkills' -count=1 -v`
  — all five tests passed, including the four required tests and coverage for
  missing/non-string hash values and sorted multi-entry results.
- `go test ./internal/cli -run '^(TestSkillsInstallRewritesTheOwnedSkillLockHash|TestBaselineUpdateSkillsStageReportsTheRewrittenLockEntries|TestBaselineUpdateSkillsStageFailsOnAnUnsafeLock|TestBaselineUpdateKeepsASetupSnapshotDifference|TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile|TestBaselineUpdateSkillsStageUsesPreviewThenConfirmation|TestBaselineUpdateSkillsStageDegradesUnreachableSourcePerSkill|TestBaselineUpdateSkillStageReportsInstalledAndRestoredSkills|TestBaselineUpdatePassesOfflineSourceToSkillStage|TestBaselineUpdateSkipsSkillStageAndPreservesSkillDirectory)$' -count=1 -v`
  — nine untagged tests passed. The repository-profile test is build-tagged,
  so it was exercised separately below.
- `go test -tags=docscontract ./internal/cli -run '^TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile$' -count=1 -v`
  — passed.
- `git diff --check` — no whitespace errors.

Acceptance evidence:

- Hashes match installed folders through both entry points; tests derive each
  expected hash from installed files. Full-byte comparisons prove unrelated
  lock bytes remain unchanged; the Baseline unit test also checks file mode.
- Matching-lock and second-stage tests verify unchanged file identity and
  modification time; both CLI paths report no rewritten name on repeat runs,
  and empty `lockUpdated` is omitted from JSON.
- Absent-lock/entry tests prove neither is added. Symlink, directory and
  malformed-lock cases reject recording through both CLI paths without
  replacing or changing the lock. The Baseline unit test additionally rejects
  an invalid schema; explicit directory installation succeeds beside unsafe
  locks without reading them.

Dependency compatibility also required one constructor update in each of
`internal/cli/doctor_trailing_skills_test.go` and
`internal/cli/repository_profile_docscontract_test.go`; both consume the changed
skills-stage dependency and retain their existing assertions. These are the
only additional implementation paths beyond the Task's listed interfaces.

The pre-existing `status: in_progress` edit was left untouched. No Task Graph
or other Task file was edited. Declared Task Verification and repository
Verification were not run; settlement remains Daemon-owned. No commit, push
or Pull Request was made.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/doctor_trailing_skills_test.go`
- `internal/cli/repository_profile_docscontract_test.go`

## Carry-forward provenance

- Source Run: `run_20261027T182924Z_39232f1ce158fe4d`
- Source commit: `129d114b66688394e54284a2e58445e2cb3afcee`
