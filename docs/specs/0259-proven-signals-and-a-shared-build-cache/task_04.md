---
task: task_04
spec: 0259-proven-signals-and-a-shared-build-cache
status: completed
type: backend
complexity: high
---

# Task 04: doctor reports and gc reclaims Worktree Build Caches

## Overview

Nothing reports the build caches that ended worktrees still hold, and
`roundfix gc` never touches a worktree (`_techspec.md` → Current behavior).
This Task adds the `worktree.caches` Project Config key, one discovery of the
declared caches in ended worktrees, a clause on Doctor's `storage` line with a
two-second bound on its measurement, and a block in the GC Command's dry and
live reports. A repository that declares no cache sees byte-identical output.
It is verifiable on its own on temporary homes with seeded Runs, Delivery Queue
items and worktrees.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-09,
   [each worktree keeps its own multi-gigabyte Go build cache](references/2026-10-09-each-worktree-keeps-its-own-go-build-cache.md),
   whose direction asks to report cache bytes in `roundfix doctor` storage and
   reclaim them in `roundfix gc`.
2. MUST add `worktree.caches` following `_techspec.md` → API Contract 8, empty
   by default, accepted in Project and User Config, carried with a comment in
   the configuration template, and refused with a message naming the key and
   the entry for an absolute path, a `..` segment, `.git` or a path under it.
3. MUST discover ended worktrees following API Contract 9 through the existing
   Run and Delivery Queue reads, from the read-only reader in Doctor and the
   GC dry run. A worktree outside the configured worktree location, a declared
   path that is a symbolic link, and a path that resolves outside its worktree
   MUST be skipped and never measured or removed.
4. MUST extend Doctor's `storage` line following API Contract 10 and Surface
   Transcript 3. The measurement MUST stop after two seconds in total and then
   report `at least <bytes> bytes`. With no declared cache the line MUST be
   byte-identical to today's, and the check MUST still write nothing and never
   report `failed` (ADR-0172).
5. MUST extend the GC Command following API Contract 11, Surface Transcript 1
   and Surface Transcript 2. The dry run writes nothing. The live run rereads
   the Run's or item's state immediately before each removal, removes only the
   declared directory, keeps every worktree of an Active Run or of an item not
   `merged`, reports only what it removed (ADR-0171), and exits 1 after the
   other removals when one fails. With no declared cache the output MUST be
   byte-identical to today's.
6. MUST add the tests named in Verification. `TestGCDryRunListsWorktreeBuildCaches`
   asserts every line of Surface Transcript 1,
   `TestGCRemovesWorktreeBuildCachesOfEndedWorktrees` every line of Surface
   Transcript 2, and `TestDoctorStorageReportsWorktreeBuildCaches` every line
   of Surface Transcript 3, each with only the transcript's declared
   variations. The existing Doctor storage and GC tests MUST pass unchanged.
7. MUST declare `.gocache` under `worktree.caches` in `.roundfixrc.yml`, the
   bounded file of the grant for this repository's own config.
8. MUST describe the key, the `storage` clause and the GC block in
   `docs/user-guide/configuration.md`, `docs/user-guide/commands/doctor.md`,
   `docs/user-guide/commands/gc.md`, the Roundfix skill's `storage` reference
   and its `setup` reference, under their own headings and never inside a
   `### QA settlement` section. MUST run `make skills-sync`, raise and record
   the skill version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   add `internal/cli/worktree_build_cache.go` to the sources of the
   `command doctor` and `command gc` entries of the Skill Coverage Map, and
   re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   The governed files of this requirement are
   `.agents/skills/roundfix/references/storage.md`,
   `.agents/skills/roundfix/references/setup.md`,
   `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`.
9. MUST add **Worktree Build Cache** to `CONTEXT.md` and revise the **GC
   Command** entry to name the Worktree Build Caches it removes, through
   `domain-modeling`, linked to ADR-0266.

## Subtasks

- [ ] Add the configuration key, its validation and its template line.
- [ ] Discover ended worktrees and measure or remove their declared caches.
- [ ] Extend the Doctor `storage` line and the GC reports.
- [ ] Declare `.gocache` in this repository's config.
- [ ] Update the guides, the skill, the coverage records and the glossary.

## Acceptance Criteria

- [ ] With a declared cache in an ended worktree, Doctor reports it found with its bytes and GC's dry run lists it without changing a byte.
- [ ] The live GC removes it, keeps the cache of an Active Run's worktree and skips a symbolic-link escape.
- [ ] Without a declared cache, Doctor and GC print exactly what they print today.
- [ ] The guides and the skill name `worktree.caches`, and the skill version is raised and recorded.

## Context

- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- instruction: `docs/adr/0172-doctor-reports-reclaimable-run-storage-from-cheap-reads.md`
- instruction: `docs/adr/0171-a-retention-prune-reports-only-what-it-reclaimed.md`
- instruction: `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- instruction: `internal/cli/run_retention.go`
- interface: `internal/config/config.go`
- creates: `internal/config/worktree_caches_test.go`
- creates: `internal/cli/worktree_build_cache.go`
- creates: `internal/cli/worktree_build_cache_test.go`
- interface: `internal/cli/doctor_storage.go`
- interface: `internal/cli/gc.go`
- interface: `.roundfixrc.yml`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/gc.md`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Worktree Build Cache**:' || { printf 'CONTEXT.md does not define **Worktree Build Cache**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**GC Command**:' || { printf 'CONTEXT.md lost **GC Command**\n' >&2; exit 1; }; for f in docs/user-guide/configuration.md docs/user-guide/commands/doctor.md docs/user-guide/commands/gc.md .agents/skills/roundfix/references/storage.md .agents/skills/roundfix/references/setup.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'worktree.caches' || { printf 'missing worktree.caches in %s\n' "$f" >&2; exit 1; }; done; grep -Eq '^ +- \.gocache$' .roundfixrc.yml || { printf '.roundfixrc.yml does not declare .gocache\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestWorktreeCachesDefaultsToNone|TestWorktreeCachesAcceptsRepositoryRelativeDirectories|TestWorktreeCachesRefusesEscapingEntries)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestWorktreeCachesDefaultsToNone TestWorktreeCachesAcceptsRepositoryRelativeDirectories TestWorktreeCachesRefusesEscapingEntries; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/config || exit 1; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the glossary lacks the term and no guide names the key, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestDoctorStorageReportsWorktreeBuildCaches|TestDoctorStorageBoundsTheWorktreeBuildCacheWalk|TestDoctorStorageIsUnchangedWithoutDeclaredCaches|TestGCDryRunListsWorktreeBuildCaches|TestGCRemovesWorktreeBuildCachesOfEndedWorktrees|TestGCKeepsWorktreeBuildCachesOfWorkInProgress|TestGCSkipsAWorktreeBuildCacheThatEscapesItsWorktree|TestGCIsUnchangedWithoutDeclaredCaches)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDoctorStorageReportsWorktreeBuildCaches TestDoctorStorageBoundsTheWorktreeBuildCacheWalk TestDoctorStorageIsUnchangedWithoutDeclaredCaches TestGCDryRunListsWorktreeBuildCaches TestGCRemovesWorktreeBuildCachesOfEndedWorktrees TestGCKeepsWorktreeBuildCachesOfWorkInProgress TestGCSkipsAWorktreeBuildCacheThatEscapesItsWorktree TestGCIsUnchangedWithoutDeclaredCaches; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'GC|Gc|Storage|Doctor|Retention' ./internal/cli || exit 1` — expected: exit 0. Before this Task the eight tests do not exist, so the PASS check fails; afterwards they and every existing GC, storage, Doctor and retention test pass.

## References

- `_prd.md` → Goals; User Story 5; Core Feature 7; Core Feature 8; Success Metric 5; User Experience
- `_techspec.md` → API Contract 8; API Contract 9; API Contract 10; API Contract 11; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Invariants 5-6; Testing Approach; Build Order 4
- ADR-0266; ADR-0172; ADR-0171; ADR-0255

## Result

### Implementation

- Added `worktree.caches` to User and Project Config, empty by default, with
  duplicate collapse, strict path validation naming the key and offending
  entry, and a commented template declaration. This repository declares
  `.gocache`.
- Added shared discovery through the existing Run listing and Delivery Queue
  readers. Terminal Run Worktrees, retained sibling Task Worktrees (including
  when the Run Worktree is absent), and merged item worktrees are eligible.
  Any unfinished owner protects the worktree, including physical aliases.
  Lexical and physical containment checks skip outside worktrees, cache
  symlinks, resolved escapes and aliases into Git metadata.
- Doctor uses the existing read-only/immutable reader selection, so an idle
  database gains no SQLite sidecars. Its cache byte walk reads directory
  entries in batches and stops synchronously at a shared two-second deadline,
  reporting a lower bound without leaving a background walk. Storage remains
  a notice (`found` or `partial`), never a readiness failure.
- GC appends the cache block after Run Retention. The dry run writes nothing;
  live GC measures and then rereads all ownership states and revalidates the
  path immediately before removing only the declared cache directory. It
  continues other candidates after an error, prints only successful cache
  removals and their bytes, diagnoses failures on stderr and exits 1.
- Updated the three user guides, the skill's setup and storage references,
  both skill copies, the version record and coverage records. Added Worktree
  Build Cache and revised GC Command in `CONTEXT.md`, linked to ADR-0266.
  The recorded Roundfix skill version is now `0.1.67`.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Doctor reports an ended cache with bytes; dry GC lists it without changing bytes | `TestDoctorStorageReportsWorktreeBuildCaches` asserts the exact storage transcript; `TestGCDryRunListsWorktreeBuildCaches` asserts the complete dry report with only timestamp/path placeholders varying. Both compare every file's bytes and the directory entries in the temporary home before and after. |
| Live GC removes ended caches, keeps active work, skips symlink escapes | `TestGCRemovesWorktreeBuildCachesOfEndedWorktrees` asserts the full live transcript, source-file preservation and repeat-GC zero reclamation. `TestGCKeepsWorktreeBuildCachesOfWorkInProgress` covers every nonmerged queue stage plus an Active Run. Escape, physical-alias, Git-metadata, retained-Task-Worktree and post-scan ownership-change tests exercise the safety checks. `TestBuildCacheGCContinuesAfterAPathFailure` proves exit 1, stderr diagnostics and successful other removals. |
| No declared cache preserves Doctor and GC output | `TestDoctorStorageIsUnchangedWithoutDeclaredCaches` asserts the exact prior storage line; `TestGCIsUnchangedWithoutDeclaredCaches` asserts the independently authored complete prior dry and live reports. Existing Doctor storage and `TestRunGC*` checks were also included unchanged in the focused run. |
| Guides and skill describe the key; skill version raised and recorded | All five required documentation surfaces name `worktree.caches` under their own headings. The sanctioned version recorder raised both skill version fields to `0.1.67`; canonical/mirror byte comparison passed. The Behavior Surface recorder passed, and the Doctor and GC source lists now include `internal/cli/worktree_build_cache.go`. |

`TestDoctorStorageBoundsTheWorktreeBuildCacheWalk` advances the existing
command clock across the budget and asserts the exact `at least 0 bytes`
storage line. `TestBuildCacheBudgetReportsMeasuredLowerBound` interrupts after
counting real file bytes and verifies the positive lower bound.
`TestBuildCacheMeasurementHonorsCancellation` proves synchronous cancellation;
`TestBuildCacheDoctorReportsUnreadablePathsAsPartial` proves a read error keeps
Doctor's exit code at zero.

### Focused checks and required regeneration

All Go commands below used
`GOCACHE=/private/tmp/roundfix-task04-gocache`.

- `go test -count=1 -run 'WorktreeCaches|BuildCache|WithoutDeclaredCaches|^TestDoctorStorage|^TestRunGC' ./internal/config ./internal/cli`
  — exit 0, both packages passed.
- `go test -race -shuffle=20261012 -count=1 -run 'WorktreeCaches|BuildCache|WithoutDeclaredCaches|^TestDoctorStorage|^TestRunGC|^TestLoadAppliesWorktreeConfigHierarchy$|^TestInit|^TestProfileGeneratedConfigUsesCompleteProfilesSchema$' ./internal/config ./internal/cli`
  — exit 0, both packages passed with the race detector and shuffled tests.
- `make skills-sync` — exit 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  — exit 0 after mirror synchronization; recorded `0.1.67`.
- `make baseline-digests` — exit 0; reported `changed:false`, no derived
  Baseline files changed.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  — exit 0; regenerated the Behavior Surface Record through its owner.
- `git diff --check` — exit 0.
- Focused path and artifact inspection — all 21 changed paths are declared in
  this Task's context; the setup reference, storage reference and `SKILL.md`
  pairs are byte-identical; all five required documentation surfaces name
  `worktree.caches`.

Initial focused attempts exposed fixture mistakes (the `CompleteRun` result
shape and mandatory Artifact Directory), discovery's default Active-only
state filter, and idle-reader sidecar creation. They were corrected at their
source before the final focused checks. The first version-record attempt
refused an unsynchronized mirror; `make skills-sync` resolved it before the
successful record command.

### Daemon handoff

The only pre-existing worktree change was the Daemon's `status: in_progress`
transition in this Task file. That status is preserved. No other Task file,
Task Graph, Verification section or QA settlement text was edited. No command
from the Task's Verification section, selected repository Verification,
incremental Verification or full suite was run. Required artifact recording
commands were run as implementation work. No commit, push or pull request was
created. Declared Verification and Task settlement remain Daemon-owned.

### Verification Feedback — attempt 1

Inspected the Daemon's diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261012T221735Z_f00baa2097cf9b98/verification/batch-004-attempt-1.log`
and the related store controller and test. The reported failure of
`make verify-changed` is the `changed` subtest of
`TestOwnerProcessControllerReprovesBeforeForceKill` in
`internal/store/process_signal_unix_test.go`, owned by Task 01.

The controller reaches its stop deadline while waiting for graceful absence,
before the identity re-proof. This test uses the real timer path with a 20 ms
grace, a 100 ms overall stop window and a 5 ms polling interval. Its changed
identity stub is read after that wait, while its absence stub remains false.
`waitForAbsence` selects between the grace timer and the stop context, so the
observed deadline expiration prevents this subtest from reaching the behavior
it intends to assert. Scheduling sensitivity is a likely explanation; the
artifact does not establish why execution was delayed.

No cache implementation defect is identified in that diagnostic. No store
source or test is changed by this Task, and changing the controller or its
signal test would cross the explicit Task 04 boundary. Follow-up for the
Daemon/Supervisor: route the timing-sensitive identity re-proof test to its
Task 01 owner; preserve its identity-change and no-forced-signal assertions
when repairing it. A whole-gate passing result has not been established.

Fresh focused check:

- `GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 -run 'WorktreeCaches|BuildCache|WithoutDeclaredCaches' ./internal/config ./internal/cli`
  — exit 0, both changed packages passed.
- `git diff --check` — exit 0 after updating this evidence.

This feedback turn updates only this Result section. Task status stays
Daemon-owned and unchanged. No declared Verification or configured full gate
was rerun, and no commit, push or pull request was created. Task 04's
implementation is handed back with the out-of-slice gate failure recorded.

## Carry-forward provenance

- Source Run: `run_20261012T221735Z_f00baa2097cf9b98`
- Source commit: `9479226c8b49427d6dd5cea97b1cc6eedd44bb80`
