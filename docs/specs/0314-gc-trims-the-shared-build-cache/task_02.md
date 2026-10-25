---
task: task_02
spec: 0314-gc-trims-the-shared-build-cache
status: completed
type: backend
complexity: high
---

# Task 02: GC trims each declared Shared Build Cache

## Overview

This Task adds the configuration key `worktree.shared_caches`, documents it
in the GC and configuration guides, and adds a GC section that resolves each declared cache, measures it with task_01's
engine, and removes (or, with `--dry-run`, reports) its idle entries above
the declared size. A repository without a declaration sees GC unchanged. The
new key changes a Behavior Surface, so the Task re-records its fingerprint.

## Requirements

1. MUST implement API Contracts 3 and 4 of `_techspec.md` and ADR-0327,
   answering the maintainer's answer "roundfix gc recolhe" (`_prd.md` →
   Sources 1).
2. MUST change `internal/config/config.go`:
   - add `SharedCaches []SharedCache` to `Worktree` after `Caches`, the
     exported `SharedCache` type of the Interfaces, a YAML overlay entry with
     `path`, `kind`, `max_size_mb` and `idle` (parsed with the existing
     `durationValue`), and `shared_caches` to the accepted `worktree` keys;
   - merge it so that a later scope's list replaces an earlier one;
   - validate each entry as API Contract 3 states, with
     `SharedCacheMinimumIdle` and the exported `ResolveSharedCachePath`
     (`{tmpdir}`, `{home}`, `{uid}`), and leave the built-in default empty
     and the generated Project Config template unchanged.
3. MUST add `internal/cli/gc_shared_cache.go` and call it from
   `internal/cli/gc.go` after the Item Binary sweep: resolve each path with
   `os.TempDir()`, the loaded home directory and `os.Getuid()`, `Measure`
   with `max_size_mb` MiB as the limit, `Apply` only without `--dry-run`,
   and print the block of API Contract 4 after the Item Binary section. A
   refusal prints `    Kept: <reason>` and is not an error; a removal error
   joins the command's error so it exits 1 after the report.
4. MUST add `internal/config/worktree_shared_caches_test.go` with
   `TestWorktreeSharedCachesLoadsAndResolvesItsTemplate` (a Project Config
   entry `{tmpdir}/roundfix-go-build-{uid}`, `go-build`, 20,480, `4h` loads;
   `ResolveSharedCachePath` with `/var/tmp/`, `/Users/me` and 501 gives
   `/var/tmp/roundfix-go-build-501`; the built-in list is empty) and
   `TestWorktreeSharedCachesRejectsAnUnsafeEntry` with the subtests unknown
   kind, relative path, parent segment, no limit, short idle (`1h59m`) and
   unknown token, each failing `Load` with an error naming
   `worktree.shared_caches`; both call `t.Parallel()` first.
5. MUST add `internal/cli/gc_shared_cache_test.go` with a fixture that
   declares `{home}/gocache` (`go-build`, `idle: 2h`) in the workspace's
   Project Config and builds a cache with the go command's README and three
   1 MiB entries aged 10 h, 3 h and 0, and three tests that each call
   `t.Parallel()` first:
   - `TestGCDryRunReportsAnOversizedSharedBuildCacheAndChangesNothing`
     (limit 1): exit 0, empty stderr, the block names the cache and
     `    Entries reclaimable: 2 (2097152 bytes)`, every entry remains;
   - `TestGCRemovesOnlyIdleEntriesOfAnOversizedSharedBuildCache` (limit 1):
     `    Entries removed: 2 (2097152 bytes)`, the two old entries are gone,
     the fresh entry and the README remain;
   - `TestGCKeepsASharedBuildCacheWithinItsLimit` (limit 9):
     `    Entries removed: 0 (0 bytes)` and every entry remains.
6. MUST add a `## Shared Build Caches` section to
   `docs/user-guide/commands/gc.md` after `## Item Binaries`, describing the
   key, the placeholders, the oldest-first trim to the limit, the minimum
   idle time and why (the go command refreshes a used entry at most hourly),
   the README and entry-name checks, the dry run, the `Kept:` line, and the
   output block of API Contract 4 in a `text` fence; and change the sentence
   of `## Worktree Build Caches` that says the shared Go build cache stays
   intact so that it says a shared cache is reclaimed only as a declared
   Shared Build Cache. Add a row to the `## Retention policy` table: "Shared
   Build Cache | its entries were used within `idle`, or the cache is within
   `max_size_mb` | the GC Command".
7. MUST add a `worktree.shared_caches` row (default `[]`) to the key table of
   `docs/user-guide/configuration.md` after `worktree.caches`, a
   `## Shared Build Caches` section after `## Worktree Build Caches` with the
   entry fields, the placeholders, the validation rules and a YAML example,
   and change that guide's sentence that puts the shared Go build cache
   outside reclamation to point at the new section.
8. MUST re-record the Behavior Surface fingerprint as the last step, after
   every source edit and after the last guide edit ("re-record after the
   last guide edit"), with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   and commit `docs/references/behavior-surfaces.json`.
9. MUST keep every existing GC, Worktree Build Cache, Item Binary and
   configuration test passing unchanged, and change no document other than
   the two guides and the Behavior Surface record.

## Subtasks

- [ ] Add the key, its merge and its validation.
- [ ] Add the GC section and wire it after the Item Binary sweep.
- [ ] Add the configuration and GC tests, describe the key and the section in the GC and configuration guides, then re-record the Behavior Surface fingerprint.

## Acceptance Criteria

- [ ] A declared cache over its limit loses only its idle entries with `roundfix gc` and nothing with `--dry-run`.
- [ ] Every unsafe declaration fails configuration loading naming the key.
- [ ] Without a declaration GC's output is unchanged, and the Behavior Surface fingerprint is current.

## Consumers

- `internal/config.Worktree.SharedCaches` and `internal/config.ResolveSharedCachePath` — consumers: `./internal/cli`

## Context

- instruction: `internal/cli/gc_item_binary.go`
- instruction: `internal/cli/gc_item_binary_test.go`
- instruction: `internal/cli/worktree_build_cache_test.go`
- instruction: `internal/config/worktree_caches_test.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`
- creates: `internal/cli/gc_shared_cache.go`
- creates: `internal/cli/gc_shared_cache_test.go`
- creates: `internal/config/worktree_shared_caches_test.go`
- interface: `internal/config/config.go`
- interface: `internal/cli/gc.go`
- interface: `docs/user-guide/commands/gc.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestWorktreeSharedCachesLoadsAndResolvesItsTemplate|TestWorktreeSharedCachesRejectsAnUnsafeEntry|TestWorktreeCaches.*|TestGCDryRunReportsAnOversizedSharedBuildCacheAndChangesNothing|TestGCRemovesOnlyIdleEntriesOfAnOversizedSharedBuildCache|TestGCKeepsASharedBuildCacheWithinItsLimit|TestGC.*ItemBinar.*)$' ./internal/config ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestWorktreeSharedCachesLoadsAndResolvesItsTemplate TestWorktreeSharedCachesRejectsAnUnsafeEntry TestWorktreeSharedCachesRejectsAnUnsafeEntry/short_idle TestGCDryRunReportsAnOversizedSharedBuildCacheAndChangesNothing TestGCRemovesOnlyIdleEntriesOfAnOversizedSharedBuildCache TestGCKeepsASharedBuildCacheWithinItsLimit TestGCDryRunListsReclaimableItemBinaries; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/gc.md docs/user-guide/configuration.md; do grep -q '^## Shared Build Caches$' "$f" || { printf 'missing section: %s\n' "$f" >&2; exit 1; }; done; go test -count=1 -timeout 20m -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract && GOOS=windows go build ./internal/config ./internal/cli` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check; with the dry run applying the plan the dry-run test fails on its entries, without the minimum idle the short-idle subtest fails, with the fingerprint not re-recorded after the guide edits the coverage test fails. Afterwards the configuration, GC and Item Binary tests pass, the Behavior Surface fingerprint is current, and the packages build for Windows.

## References

- `_prd.md` → Goals; Core Feature 2; Core Feature 3; User Experience; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior; Data Models; API Contracts 3-4; Invariants 4-5; Testing Approach; Build Order 2


## Result

Implemented this Task's configuration, GC section, tests and two guide updates.
The existing `status: in_progress` was left unchanged. At entry, only this
Task file was already modified (the Daemon-owned status); no other pre-existing
changes were present. Initial inspection found no shared-cache configuration
key or GC integration, while task_01's engine was available.

### Acceptance evidence

- **Declared cache trimming and dry run:** GC resolves declarations and measures
  with the configured MiB limit, calls Apply only for live GC, and prints each
  block after Item Binaries. The new command tests assert the complete block,
  two idle entries totaling 2,097,152 bytes removed or reclaimable, preservation
  of fresh entries and README, unchanged dry-run contents, and no removal below
  the limit. An additional command test proves a missing README produces
  `Kept:` with exit 0 and no deletion. Removal errors join the GC error and the
  report is printed before exit 1.
- **Unsafe declarations:** the configuration tests cover unsupported kind,
  relative and unclean paths, parent segments, missing/nonpositive limit,
  short idle, unknown placeholder, malformed duration and byte-count overflow;
  every error names `worktree.shared_caches`. Tests also cover template
  resolution, the empty built-in default, the unchanged generated Project
  Config template, and replacement across scopes (including an empty list).
- **No declaration and current fingerprint:** existing GC and Worktree Build
  Cache transcript tests ran unchanged in the focused regression selection.
  No declaration adds no report block. The two guides describe the key,
  validation and reclamation rules. The sanctioned recorder ran after all
  source and guide edits and updated only the configuration-key fingerprint
  and those two guide fingerprints.

### Focused checks

- `go test ./internal/config ./internal/cli -run 'Shared(Build)?Cache|SharedCaches' -count=1 -timeout=5m`
  — exit 0; both packages passed the new shared-cache checks.
- `go test ./internal/config -count=1 -timeout=5m`
  — exit 0; configuration regression tests passed after the final source edit.
- `go test ./internal/cli -run '^Test(GC|WorktreeBuildCache)' -count=1 -timeout=5m`
  — exit 0; focused GC, Item Binary and Worktree Build Cache regressions passed.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  — exit 0; required artifact generation after the last source and guide edits.
- `git diff --check` — exit 0 before recording this Result.

The shell commands used `rtk proxy` for unfiltered output. The executor rejected
its symlinked writable-root configuration; approved execution accessed the same
assigned checkout at its resolved `/Volumes/Brain-SSD/cache/.roundfix/worktrees/`
location. No sandbox configuration was changed.

Declared Verification, the Windows cross-build and repository Verification were
not run; they remain Daemon-owned. No Task status, Task Graph, other Task file,
repository declaration, skill or glossary was changed. No commit, push or pull
request was made. No follow-up implementation was added.
