---
task: task_02
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
status: pending
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
