---
task: task_04
spec: 0259-proven-signals-and-a-shared-build-cache
status: pending
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
