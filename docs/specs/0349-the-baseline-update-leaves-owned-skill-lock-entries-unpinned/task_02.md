---
task: task_02
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
status: completed
type: backend
complexity: medium
---

# Task 02: Doctor warns with DR-SKILL-LOCK-STALE for an owned lock entry that pins another copy

## Overview

`roundfix doctor` judges the Roundfix-owned skills by version only, so it
reported `skills: ok` on Quaestor while every owned lock entry still pinned
the older catalog copy that the next `skills update` installed. This Task
adds a read-only scan of the owned lock entries and makes the `skills:` line
warn, naming each Stale Owned Lock Entry and the command that rewrites it. It
answers the marcioaltoe-skills report of 2026-11-04 (ai-memory message
`01a1a5f7`, Suggestion C, adopted as
[references/2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md](references/2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md)).

## Requirements

1. MUST add `StaleOwnedLockEntry` and `StaleOwnedLockEntries` to
   `internal/baseline/owned_skills_lock.go` exactly as `_techspec.md` →
   Interfaces and API Contract 3 state. The scan never writes. An absent lock
   returns nothing, and an unsafe or malformed lock returns the
   `loadSkillsLock` error. It reads only the entries of the given names that
   have an installed tree, and reports an entry with a `ref` member or a
   `computedHash` other than the tree's digest.
2. MUST make `repositorySkillsCheck` in `internal/cli/doctor.go`, on the
   Setup Manifest branch and after the snapshot suffix, call
   `baseline.StaleOwnedLockEntries(ctx, root, skills.Names())` and apply API
   Contract 4:
   - the exact suffix
     `; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: <names>`;
   - `ok` turns into `warn`, and any other status is kept;
   - `ownedSkillsNextAction` is appended once, joined by ` && `;
   - on a scan error, the suffix `; skills-lock.json comparison unavailable: <error>`
     with status and next action kept.
   No new `doctorDependencies` field is added, and the exit-code rule of
   `runDoctorCommand` is unchanged.
3. MUST add `TestDoctorNamesAStaleOwnedSkillLockEntry` and
   `TestDoctorKeepsAFailedSkillsLineWithAStaleOwnedSkillLockEntry` in
   `internal/cli/doctor_skills_lock_test.go`, each asserting the whole printed
   `skills:` line or the parts `_techspec.md` → Testing Approach names. They
   include an installed external `testing-boss` folder with a pinned entry
   that is never named, and the `ok` line after one `RecordInstalledOwnedSkills`.
   Every new top-level test calls `t.Parallel()` first and uses only
   `t.TempDir()` roots.
4. MUST update `docs/user-guide/commands/doctor.md` (`_techspec.md` → API
   Contract 5):
   - the `skills:` bullet names `DR-SKILL-LOCK-STALE`;
   - a paragraph next to the Setup Snapshot Difference paragraph gives the
     exact suffix of API Contract 4 and the example line
     `skills: warn (<counts>; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: <names>; next: roundfix skills install --target project)`;
   - it says `ok` becomes `warn`, `failed` and `unversioned` stay, the exit
     code never changes, external entries are never named, and the release
     plan's `skills` check reads the same line;
   - it says that in the catalog's sync window right after a release, `skills
     update` over an unpinned entry can still install the previous copy, which
     the version check reports and one bundle install restores.
5. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last edit of `internal/cli/doctor.go`
   and of the doctor guide, and never hand-edit it.
6. MUST NOT change `roundfix skills check`, the readiness texts of
   `doctorSkillReadinessResult`, the external next action, the release plan
   code, any other guide, any skill, the Makefile, CI workflows, `go.mod` or
   `go.sum`.

## Subtasks

- [ ] Add the read-only scan.
- [ ] Append the warning to the Doctor `skills:` line.
- [ ] Add the two Doctor tests.
- [ ] Update the doctor command guide, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A pinned or mis-hashed installed owned entry makes the `skills:` line `warn` with `DR-SKILL-LOCK-STALE` and the install command; an external entry is never named.
- [ ] A `failed` line stays `failed` and names the install once.
- [ ] One bundle install over the stale lock makes the line `ok` again.

## Context

- interface: `internal/baseline/owned_skills_lock.go`
- interface: `internal/cli/doctor.go`
- creates: `internal/cli/doctor_skills_lock_test.go`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `baseline.StaleOwnedLockEntries` and the Doctor `skills:` line, which `roundfix release plan` reuses for its `skills` check — consumers: `./internal/cli`
- The doctor command guide and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDoctorNamesAStaleOwnedSkillLockEntry|TestDoctorKeepsAFailedSkillsLineWithAStaleOwnedSkillLockEntry)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDoctorNamesAStaleOwnedSkillLockEntry TestDoctorKeepsAFailedSkillsLineWithAStaleOwnedSkillLockEntry; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 40m ./internal/baseline ./internal/cli` — expected: exit 0. Before this Task neither test exists, so the first pass line is missing; afterwards they pass with both whole packages.
- `grep -q 'DR-SKILL-LOCK-STALE' internal/cli/doctor.go || { echo 'missing DR-SKILL-LOCK-STALE in Doctor' >&2; exit 1; }; grep -q 'DR-SKILL-LOCK-STALE' docs/user-guide/commands/doctor.md || { echo 'missing DR-SKILL-LOCK-STALE in the doctor guide' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither Doctor nor its guide has the finding; afterwards the documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 3-4; Success Metric 4; Success Metric 5; Glossary
- `_techspec.md` → Current behavior 4; Interfaces; API Contracts 3-5; Invariant 5; Declared breaks; Testing Approach; Build Order 2


## Result

Skills loaded: implement-task, agentic-cli-design, coding-guidelines, golang-cli, golang-context, golang-error-handling, golang-structs-interfaces, golang-testing, no-workarounds, systematic-debugging, tech-writer, testing-boss, evidence-gate.

Implemented the Task slice for Daemon Verification:

- Added `StaleOwnedLockEntry` and the read-only `StaleOwnedLockEntries` scan. It sorts and deduplicates the supplied names, considers only existing object entries with installed trees, detects any `ref` member or absent/mismatched hash, and preserves loader errors and cancellation.
- Doctor scans after the Setup Snapshot suffix on the Setup Manifest branch. It adds the exact `DR-SKILL-LOCK-STALE` suffix, changes only `ok` to `warn`, and adds the owned install action once with ` && `. Comparison errors append their diagnostic without changing readiness status or next action. No dependency field, readiness formatter, external remediation, release-plan code or exit-code logic changed.
- Updated only the Doctor command guide with the warning suffix and example, status and exit behavior, external-entry exclusion, shared release-plan line, and catalog sync-window limitation.
- Re-recorded `docs/references/behavior-surfaces.json` through its sanctioned command after the final Doctor source and guide edits; the generated diff contains only the Doctor guide fingerprint.

Acceptance evidence:

- Pinned or mis-hashed installed owned entries warn, and external entries are excluded: `TestDoctorNamesAStaleOwnedSkillLockEntry` compares the complete printed line for current-hash pinned `archive-spec` and `roundfix`, hash-only stale `write-prd`, and an installed pinned external `testing-boss` that is never named. It also checks that Doctor changes no repository bytes.
- A failed line remains failed with the install action once: `TestDoctorKeepsAFailedSkillsLineWithAStaleOwnedSkillLockEntry` lowers both installed `roundfix` version declarations, pins its entry with a zero hash, and asserts the failed prefix, exact warning before the next action, and one install command.
- One bundle recording clears the stale lock warning: both rows of `TestDoctorNamesAStaleOwnedSkillLockEntry` call `RecordInstalledOwnedSkills` once and then compare the complete `skills: ok (<counts>)` line.
- Additional focused tests cover sorted/deduplicated scan results and `Ref`/`Hash` fields, ignored uninstalled/nonobject/unlisted entries, cancellation, empty/null/nonstring refs, missing/nonstring hashes, absent/malformed/unsafe locks, preservation of failed/unversioned readiness, snapshot suffix ordering, and joining an existing external action. Every new top-level test starts with `t.Parallel()` and uses `t.TempDir()` roots.

Focused checks (all final commands exited 0):

- Red reproduction before implementation: `go test -count=1 -timeout 5m ./internal/cli -run 'TestDoctor.*StaleOwned'` exited 1. The pinned fixture printed `skills: ok` instead of the warning, and the failed fixture omitted the warning. The same focused command passed after implementation.
- `go test -count=1 -timeout 5m ./internal/cli -run 'Test(Doctor.*(StaleOwned|OwnedLock)|StaleOwnedLockEntries|DoctorNamesASetupSnapshotDifferenceWithoutWarning|DoctorSnapshotComparisonPreservesReadinessAndFiltersSkills|RunDoctorRepositorySkillReadiness|DoctorAndSkillsCheckReportSharedOwnedSkillReadiness|RunDoctorRealRepositoryCheckDoesNotMutateState|DoctorLeavesAnObsoleteLockEntryUntouched)'` — passed, including the new tests and adjacent readiness/snapshot regressions.
- `go test -count=1 -timeout 5m ./internal/baseline -run '^TestRecordInstalledOwnedSkills'` — passed.
- `go test -count=1 -timeout 5m -tags docscontract ./internal/docscontract -run '^(TestCommandReferenceLinksResolve|TestUserGuideLinksResolve)$'` — passed.
- Final implementation step: `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — passed; no generated record was hand-edited.
- `git diff --check` — passed. Changed paths are limited to this Task's declared interfaces, its new test file, and this Result section; the pre-existing Daemon status edit is preserved.

The declared Verification commands and repository gates remain Daemon-owned and were not run. No other Task, manifest, skill, guide, tooling configuration, dependency or root lock was edited. The Verification Feedback follow-up is recorded below.


### Verification Feedback — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261104T183354Z_6fb3237eae19e727/verification/batch-002-attempt-1.log`.
The Daemon's `make verify-changed` reached a suiteguard process-residue failure
in `internal/verifyselect` after that package's test assertions passed. The
remaining helper was a detached Git filesystem monitor. This is not an
assertion failure in the Task's lock scan or Doctor tests.

The selector's `TestPartitionFollowsTheMakefileRecipes` calls
`makeRecipeCLITests` (`internal/verifyselect/verifyselect_test.go`) with its
make process rooted in the real checkout. Parsing the Makefile evaluates
`BUILD_COMMIT`, which runs Git status without disabling the filesystem
monitor. A read-only query confirmed that the effective `core.fsmonitor=true`
comes from user Git configuration. A temporary `GIT_TRACE2_EVENT` capture of
the focused partition test confirmed seven unisolated status commands and
21 filesystem-monitor protocol events. That invocation started no new
monitor and exited 0; it does not establish that the Daemon's residue failure
is repaired. Trace files were removed with their temporary directory.

Fresh focused checks:

- `go test -count=1 -timeout 3m ./internal/verifyselect -run '^TestPartitionFollowsTheMakefileRecipes$'` with a temporary Git trace — exit 0; used to trace the pre-existing Makefile/Git boundary, not as a substitute for Daemon Verification.
- `go test -count=1 -timeout 5m ./internal/cli -run 'Test(Doctor.*(StaleOwned|OwnedLock)|StaleOwnedLockEntries)'` — exit 0; the Task's warning, recovery, scan, error and status/action tests still pass.
- `git diff --check` — exit 0 after the evidence update.

No Task source change is indicated by this failure. Follow-up outside this
slice: isolate Git configuration at the selector's make-process boundary so
the Makefile's build-stamp read cannot start a detached monitor under the
test binary's suiteguard mark. Task 02 expressly prohibits Makefile changes
and does not own `internal/verifyselect`, so neither was edited, nor was
suiteguard weakened or a Git/process setting changed to mask the failure.
The Behavior Surface Record remains current with the unchanged Doctor
source and guide. Task status is preserved; the declared Verification and
`make verify-changed` were not rerun by the Agent. The Daemon owns the next
full configured Verification sequence.
