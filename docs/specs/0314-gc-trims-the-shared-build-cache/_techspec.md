---
spec: 0314-gc-trims-the-shared-build-cache
prd: _prd.md
created: 2026-10-24
---

# GC trims the shared build cache — Technical Spec

## Executive Summary

A new package `internal/sharedcache` measures a Go build cache and plans and
applies a trim of its idle entries. A new configuration key,
`worktree.shared_caches`, declares which caches GC may trim and how. The GC
Command gains a section after its Item Binary section that resolves each
declared path, measures it, and removes or reports the planned entries. The
proof that GC never removes an entry in use rests on the go command's
documented behavior: it refreshes an entry's modification time at most once
an hour when it uses it, and the minimum idle time is two hours; the apply
rereads each entry's time just before removing it.

The trade-off accepted: an entry idle for the declared time is removed even
if a long Run will want it again later; that Run recompiles it.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files only. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0327, as `_prd.md` → Project
  Constraints states. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — protected paths limited to the Roundfix
  skill's storage reference with its version fields and this repository's
  Project Config declaration. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0314-gc-trims-the-shared-build-cache/_authorization.md`;
  bounded files: `.agents/skills/roundfix/references/storage.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.roundfixrc.yml`.

## Current behavior

Read on 2026-10-24 at `56dd1a76`.

1. `internal/cli/gc.go` `runGC` scans Worktree Build Caches when
   `worktree.caches` is declared, then sweeps Item Binaries, then applies
   Run Retention; `printGCReport` prints Run Retention, Worktree Build Cache
   and Item Binary sections in that order. A removal failure is diagnosed on
   stderr, other removals proceed, and the command exits 1.
2. `internal/config` `Worktree` holds `Caches []string`; the overlay accepts
   `concurrency`, `location`, `copy`, `caches`, `bootstrap` and
   `bootstrap_timeout` and refuses any other key; `durationValue` parses Go
   durations.
3. `docs/user-guide/commands/gc.md` says GC leaves "the shared Go build
   cache" intact. The Roundfix skill's `references/storage.md` describes
   Worktree Build Caches.
4. The Makefile cache `$TMPDIR/roundfix-go-build-<uid>` holds the go
   command's `README` ("This directory holds cached build artifacts from the
   Go build system."), `trim.txt`, `testexpire.txt` and entries
   `<2 hex>/<64 hex>-a|-d`.
5. Any change to the configuration keys makes the Behavior Surface
   fingerprint stale, and `TestTheSkillCoverageMapIsCurrent` fails until it
   is re-recorded.

## System Architecture

- `internal/sharedcache/sharedcache.go` (new): `Measure`, `Apply`,
  `MinimumIdle`, `GoBuildReadme`, `ErrNotGoBuildCache`.
- `internal/config/config.go`: `SharedCache`, the overlay, merge,
  validation, `SharedCacheMinimumIdle` and `ResolveSharedCachePath`; the
  generated Project Config template is unchanged.
- `internal/cli/gc_shared_cache.go` (new) and `internal/cli/gc.go`: the GC
  section.
- Documents: the GC and configuration guides, the skill reference, the
  glossary, `.roundfixrc.yml`, the Behavior Surface record.

## Implementation Design

### Interfaces

```go
// internal/sharedcache
const MinimumIdle = 2 * time.Hour
const GoBuildReadme = "This directory holds cached build artifacts from the Go build system."
var ErrNotGoBuildCache = errors.New("not a Go build cache")
type Entry struct { Path string; Bytes int64; ModTime time.Time }
type Plan struct { Path string; Bytes, Limit int64; Remove []Entry; RemoveSize int64 }
func Measure(dir string, limit int64, idle time.Duration, now time.Time) (Plan, error)
func Apply(plan Plan, idle time.Duration, now func() time.Time) (removed int, bytes int64, err error)

// internal/config
type SharedCache struct { Path, Kind string; MaxSizeMB int64; Idle time.Duration }
const SharedCacheMinimumIdle = 2 * time.Hour
func ResolveSharedCachePath(template, tmpdir, home string, uid int) (string, error)
```

### Data Models

The configuration key:

```yaml
worktree:
  shared_caches:
    - path: "{tmpdir}/roundfix-go-build-{uid}"
      kind: go-build
      max_size_mb: 20480
      idle: 4h
```

A Project Config list replaces a User Config list, as `worktree.copy` does.

### API Contracts

1. API Contract: `Measure` refuses an idle time below `MinimumIdle`, a path
   that is a symbolic link or not a directory, and a directory whose
   `README` does not start with `GoBuildReadme`. It sums every regular
   file's size. When the sum exceeds `limit` it plans, oldest first, the
   regular files named `<64 lowercase hex>-a` or `-d` directly inside a
   directory named by two lowercase hex digits and modified before
   `now - idle`, stopping when the remaining size is within `limit`. It
   writes nothing.
2. API Contract: `Apply` rereads each planned entry with `Lstat`; it skips a
   missing entry, a non-regular file and an entry modified at or after
   `now() - idle`; it removes the rest and returns the count, the bytes and
   every removal error joined.
3. API Contract: configuration loading refuses a `worktree.shared_caches`
   entry with a kind other than `go-build`, a path that after expanding
   `{tmpdir}`, `{home}` and `{uid}` is not a clean absolute path or still
   holds a brace, a path with `..`, `max_size_mb` not above 0, or `idle`
   below two hours; every message names `worktree.shared_caches`.
4. API Contract: `roundfix gc` and `roundfix gc --dry-run` print, after the
   Item Binary section, for each declared cache: `  Shared Build Cache:
   <resolved path>`, then either `    Kept: <reason>` or `    Size: <bytes>
   bytes of <limit> allowed` and `    Entries reclaimable|removed: <n>
   (<bytes> bytes)`. The dry run never calls `Apply`. A removal error is
   diagnosed on stderr and makes the command exit 1 after the report, as
   the other sections do. With no declaration the output is byte-identical
   to today's.

### Surface Transcripts

None. The GC output block of `_prd.md` → User Experience is asserted line by
line by task_02's GC tests next to the existing GC transcript tests; no new
command and no whole-command transcript changes.

### Invariants

1. GC never removes a file that is not an idle Go build cache entry of a
   declared cache whose directory holds the go command's README.
2. An entry modified within the idle time before the apply is never
   removed.
3. A cache within its limit loses nothing.
4. The dry run writes nothing.
5. Without a declaration, GC's output and exit code are unchanged.

## Coverage Map

- Goals 1 and 2, Core Feature 1, Success Metric 1 → API Contracts 1 and 2,
  Invariant 1, task_01.
- Success Metric 2 → API Contracts 1 and 2, Invariants 2 and 3, task_01.
- Success Metric 3 → API Contract 1, Invariant 1, task_01.
- Goal 3, Core Features 2 and 3, Success Metrics 4 and 5 → API Contracts 3
  and 4, Invariants 4 and 5, task_02.
- Goal 4, Core Feature 4 → Invariant 5, task_02 (guides), task_03
  (declaration, skill, glossary).

Deviations from the PRD: None.

## Absence Evidence

1. A shared cache key — command: `grep -rn "shared_caches" internal docs`;
   result: no match.
2. A trim of a cache outside worktrees — command:
   `grep -rln "sharedcache" internal`; result: no match.

## ADR Alternatives

- ADR-0327: this Spec implements it. Its rejected options: a Makefile trim
  (the maintainer chose the product), emptying the cache above the limit
  (three concurrent Runs share it), and discovering the cache through
  `go env GOCACHE` (the Makefile sets its own).

## Integration Points

None. No external system, network or provider is reached; the go command is
not run.

## Testing Approach

- `internal/sharedcache/sharedcache_test.go` (new), parallel tests over a
  fixture cache in `t.TempDir()` with times set by `os.Chtimes` against a
  fixed `now`:
  `TestMeasurePlansIdleEntriesOldestFirstUntilWithinTheLimit`,
  `TestMeasureKeepsEveryEntryWithinTheLimitOrInUse`,
  `TestMeasureRefusesWhatIsNotAGoBuildCache`,
  `TestApplySkipsAnEntryUsedSinceItWasPlanned`.
- `internal/config/worktree_shared_caches_test.go` (new):
  `TestWorktreeSharedCachesLoadsAndResolvesItsTemplate` and
  `TestWorktreeSharedCachesRejectsAnUnsafeEntry` (six subtests).
- `internal/cli/gc_shared_cache_test.go` (new), through the GC command with
  the package's temporary workspace:
  `TestGCDryRunReportsAnOversizedSharedBuildCacheAndChangesNothing`,
  `TestGCRemovesOnlyIdleEntriesOfAnOversizedSharedBuildCache`,
  `TestGCKeepsASharedBuildCacheWithinItsLimit`.
- The existing GC transcript tests prove Invariant 5.

## Build Order

1. The trim engine and its tests (task_01).
2. The key, the GC section, their tests, the GC and configuration guides,
   and the Behavior Surface record re-recorded after the last guide edit
   (task_02) (depends on: 1).
3. This repository's declaration, the skill reference with its version, the
   glossary, and the Behavior Surface record re-recorded after the skill
   edit (task_03) (depends on: 2, because both re-record
   `docs/references/behavior-surfaces.json`).
4. The QA gate (task_04) (depends on: 1, 2, 3).

No other Spec is required.

## Risks & Considerations

- A cache shared by several users is skipped by nothing in particular; the
  README and name checks still bound what is removed, and a file another
  user owns fails its removal, which is diagnosed.
- On Windows `{uid}` expands to `-1`; the declaration of this repository is
  Unix-only like its Makefile.
- `du`-like measurement walks about 70,000 files, about 2 s on the
  maintainer's machine; GC runs it only for declared caches.

## Glossary

- adds: **Shared Build Cache**

## Decisions

- A new package keeps the trim testable without the GC Command's store.
- Declared caches only, by placeholders (ADR-0327).
- A Project Config list replaces a User Config list.
