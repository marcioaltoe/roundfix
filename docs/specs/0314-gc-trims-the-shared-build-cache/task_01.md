---
task: task_01
spec: 0314-gc-trims-the-shared-build-cache
status: pending
type: backend
complexity: medium
---

# Task 01: A trim engine removes only idle Go build cache entries

## Overview

GC needs to trim a Go build cache that several worktrees share, without ever
removing an entry a build is using or a file that is not a cache entry. This
Task adds the package `internal/sharedcache`: `Measure` walks a cache and
plans the removal of entries idle for at least the given time, oldest first,
until the cache is within its limit; `Apply` rereads each planned entry and
removes it only if it is still idle.

## Requirements

1. MUST answer the maintainer's answer of 2026-10-24 "roundfix gc recolhe"
   and the shared-cache finding of
   [`docs/references/2026-10-24-test-suite-audit.md`](../../references/2026-10-24-test-suite-audit.md)
   (`_prd.md` → Sources 1 and 2) by implementing API Contracts 1 and 2 of
   `_techspec.md` and ADR-0327.
2. MUST add `internal/sharedcache/sharedcache.go` with a package comment
   and the Interfaces of `_techspec.md`: `MinimumIdle` (two hours, commented
   as twice the go command's hourly refresh of a used entry), `GoBuildReadme`,
   `ErrNotGoBuildCache` (wrapped by every refusal of a directory), `Entry`,
   `Plan`, `Measure` and `Apply`, standard library only.
3. MUST add `internal/sharedcache/sharedcache_test.go` with four tests that
   each call `t.Parallel()` first, on a fixture cache built in `t.TempDir()`
   with a `README` that starts with `GoBuildReadme`, entries
   `aa/<64×a>-d` (1,000 bytes, 10 h old), `aa/<64×b>-a` (1,000 bytes, 5 h),
   `bb/<64×c>-d` (1,000 bytes, 3 h), `bb/<64×d>-d` (1,000 bytes, 30 min), a
   `trim.txt` (48 h) and `aa/notes.txt` (1,000 bytes, 48 h), times set with
   `os.Chtimes` against a fixed `now` of 2026-10-24 12:00 UTC:
   - `TestMeasurePlansIdleEntriesOldestFirstUntilWithinTheLimit`: limit
     3,500 bytes and idle 2 h plan exactly the `a` then the `b` entry;
     `Apply` removes those two (2,000 bytes); `README`, `trim.txt`,
     `aa/notes.txt` and the `c` and `d` entries remain;
   - `TestMeasureKeepsEveryEntryWithinTheLimitOrInUse`: a 1 MiB limit plans
     nothing; a zero limit plans exactly the three entries older than two
     hours and none modified less than two hours before `now`;
   - `TestMeasureRefusesWhatIsNotAGoBuildCache`: a directory without
     `README`, a symbolic link to the fixture cache, and an idle time of one
     hour are each refused;
   - `TestApplySkipsAnEntryUsedSinceItWasPlanned`: after planning with a zero
     limit, the first planned entry's time is set to `now`; `Apply` removes
     every other planned entry and keeps that one.
4. MUST NOT change any other package, any document or any configuration.

## Subtasks

- [ ] Add the engine.
- [ ] Add the four tests.

## Acceptance Criteria

- [ ] Over its limit the fixture loses its oldest idle entries only; under it, nothing.
- [ ] A directory that is not a Go build cache, a symbolic link and a short idle time are refused.
- [ ] An entry used after planning survives the apply.

## Consumers

None. The package is new and only task_02's GC section will call it.

## Context

- instruction: `internal/cli/worktree_build_cache.go`
- instruction: `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`
- instruction: `docs/references/2026-10-24-test-suite-audit.md`
- creates: `internal/sharedcache/sharedcache.go`
- creates: `internal/sharedcache/sharedcache_test.go`

## Verification

- `out="$(go test -count=1 -timeout 10m -v ./internal/sharedcache 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMeasurePlansIdleEntriesOldestFirstUntilWithinTheLimit TestMeasureKeepsEveryEntryWithinTheLimitOrInUse TestMeasureRefusesWhatIsNotAGoBuildCache TestApplySkipsAnEntryUsedSinceItWasPlanned; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/sharedcache` — expected: exit 0. Before this Task the package does not exist, so the command fails; without the README check the refusal test fails, with the idle cutoff widened the in-use test fails, without the entry-name match the oldest-first test fails on `aa/notes.txt`, and without the reread the apply test fails. Afterwards the four tests pass and the package vets for Windows.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior; API Contracts 1-2; Invariants 1-3; Testing Approach; Build Order 1
