---
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
prd: _prd.md
created: 2026-10-26
---

# Runs trim the shared build cache, and doctor warns on low disk — Technical Spec

## Executive Summary

The GC Command already trims a declared Shared Build Cache safely (ADR-0327);
what is missing is a trigger nobody has to remember. This Spec makes the trim
engine stoppable through a context, wraps it in one function that trims every
declared cache under a thirty-second budget and classifies each outcome, and
calls it at the two places where new build output enters the cache: right
after an Implement Run is recorded, and before the Delivery Queue's repository
gate. Neither caller can fail because of it. `roundfix doctor` gains a
`disk-space:` warning line that reuses the free-space reader `gc compact`
already has. The trade-off accepted: an Implement Run's start may spend up to
thirty seconds trimming when its cache is far above its limit, in exchange for
never filling the disk between manual `roundfix gc` runs; a healthy cache
costs one directory walk, 1.2 s for this repository's 21.6 GB cache.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or schema
  is minted; the Run Event reuses `daemon.status`. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads and removals
  only. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0327 and ADR-0334, as `_prd.md` →
  Project Constraints states; ADR-0334 is recorded by this planning change.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — protected paths limited to the Roundfix
  skill's `deliver`, `implement`, `setup` and `storage` references with its
  version fields, and the command test harness in `internal/cli/cli_test.go`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/storage.md`,
  `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-10-26 at `d08f8557` (v0.123.0 plus the 0319 and 0320 plans).

1. `internal/sharedcache` `Measure(dir, limit, idle, now)` walks the whole
   cache and plans idle entries oldest first; `Apply(plan, idle, now)`
   rereads each entry and removes it. Neither takes a context, so a caller
   cannot bound either in time.
2. `internal/cli/gc_shared_cache.go` `sweepSharedCaches(ctx, loaded, dryRun)`
   resolves each `worktree.shared_caches` entry with
   `roundconfig.ResolveSharedCachePath(cache.Path, os.TempDir(),
   loaded.HomeDir, os.Getuid())`, measures, applies without `--dry-run`, and
   returns one `sharedCacheReport` per cache; a refusal
   (`sharedcache.ErrNotGoBuildCache`) becomes `kept`, other errors join the
   command's error. `printSharedCacheReport` prints the block of the GC guide.
   `runGC` in `internal/cli/gc.go` is its only caller.
3. `internal/cli/implement.go` `runImplementCommand` runs Run Retention, then
   records the Run with `createRunReclaimingOrphan`/`CreateRun`, calls
   `req.reportDetachedRunCreated(run.ID)` (a detached parent returns there),
   and then `createRunWorktree`. The Delivery Queue starts every item's Run
   as a child `implement --spec <slug>` (`commandDeliveryWorkflow.RunSpec` in
   `internal/cli/deliver_workflow.go`), so every Implement Run passes this
   code. Run Events written outside the Daemon use
   `runStore.AppendRunEvent` (for example the notification receipt in
   `internal/cli/cli.go`, a `daemon.status` event with a JSON payload).
4. `commandDeliveryWorkflow.Gate` runs `defaults.verification` in the item
   worktree with `daemon.ExecVerifier`, then the optional Linux Lane, and
   writes its notices to `workflow.log`, the Queue Owner's console log. A
   parked item's `deliver retry` resumes at the parked stage, so a gate that
   failed for lack of space reruns the gate without a new Run.
5. `internal/cli/doctor.go` `runDoctorCommand` prints `storage:` from the
   `storage` dependency, then `linux-lane:` only when the lane is on, then
   `codex:`. Nothing reads a filesystem's free space. `internal/store`
   `availableTemporaryCapacity(path)` (unexported; `compaction_capacity_unix.go`
   with `unix.Statfs`, `compaction_capacity_windows.go` with
   `windows.GetDiskFreeSpaceEx`) does, for `gc compact`.
6. The doctor tests build a literal `doctorDependencies` in
   `withDoctorFakeLoadedAndReadiness` (a new field stays nil there) or start
   from `defaultDoctorDependencies()` in `withDoctorLiveDeps`; the command
   tests start from `defaultCommandDependencies()` in
   `commandEnvironmentForTest` and `updateCommandDependenciesForTest`
   (`internal/cli/cli_test.go`).
7. Docs: the GC guide's `## Shared Build Caches` section and its Retention
   policy row, the configuration guide's `## Shared Build Caches` section,
   the skill's `references/storage.md` and **Shared Build Cache** in
   `CONTEXT.md` all say only the GC Command trims; the doctor guide and
   `references/setup.md` list no free-space line.

## System Architecture

- `internal/sharedcache/sharedcache.go`: `MeasureContext`, `ApplyContext`;
  `Measure` and `Apply` call them with `context.Background()`.
- `internal/cli/gc_shared_cache.go`: `trimSharedCaches`, the outcome
  classification, the console and log lines, and `sweepSharedCaches`
  rebuilt on it with GC's bytes unchanged.
- `internal/cli/gc.go`: `gcDependencies.sharedCacheTrimBudget`, default
  thirty seconds.
- `internal/cli/implement.go`: the Run-start trim and its Run Events.
- `internal/cli/deliver_workflow.go`: the gate trim.
- `internal/store/available_bytes.go`: exported `AvailableBytes`.
- `internal/cli/doctor_disk_space.go`, `internal/cli/doctor.go`,
  `internal/cli/health.go`, `internal/cli/cli.go`: the `disk-space:` check,
  its doctor dependency and the `availableBytes` command dependency.
- Docs (task_03): the GC, doctor, implement, deliver and configuration
  guides, the skill's `storage`, `setup`, `implement` and `deliver` references
  and their mirrors (the Behavior Surfaces `command gc`, `command doctor`,
  `command implement` and the `command deliver …` surfaces cover the changed
  sources), the skill's version pair,
  `skills/testdata/owned-skill-versions.json`, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json`.

## Implementation Design

### Interfaces

```go
// internal/sharedcache: stop between entries once ctx ends, returning ctx.Err()
// wrapped; ApplyContext keeps what it already removed and reports it.
func MeasureContext(ctx context.Context, dir string, limit int64, idle time.Duration, now time.Time) (Plan, error)
func ApplyContext(ctx context.Context, plan Plan, idle time.Duration, now func() time.Time) (removed int, bytes int64, err error)

// internal/cli/gc_shared_cache.go
type sharedCacheOutcome string // "trimmed", "unchanged", "kept", "failed", "stopped"
// trimSharedCaches trims every declared cache in declaration order under ctx;
// dryRun plans only (GC). Each report carries its outcome and cause.
func trimSharedCaches(ctx context.Context, loaded roundconfig.Loaded, dryRun bool) []sharedCacheReport

// internal/store/available_bytes.go: the free bytes an unprivileged user may
// write on the filesystem holding path (an existing path).
func AvailableBytes(path string) (int64, error)
```

Outcome rules, per cache: a path that does not resolve is `failed`;
`sharedcache.ErrNotGoBuildCache` is `kept`; a context that ended during
measure or apply is `stopped`; another measure or removal error is `failed`;
otherwise `trimmed` when at least one entry was removed and `unchanged` when
none was. `cause` is the error text; for `stopped` it is
`budget of <budget> elapsed` when the budget's deadline passed, with the
budget in Go duration form (`30s`), and the error text for any other
cancellation. Once the budget ends, every cache not yet trimmed is also
reported `stopped`, so each declared cache still has one report.
`size_bytes` is the measured size, `0` when measuring did not finish. The removed count and bytes are what was actually
removed, also for `failed` and `stopped`.

The Run start, in `runImplementCommand` after `reportDetachedRunCreated`
succeeds and before `createRunWorktree`: with at least one declared cache, it
runs `trimSharedCaches` under `context.WithTimeout(runCtx,
gc.sharedCacheTrimBudget)` and appends one Run Event per report through
`runStore.AppendRunEvent`; an append failure prints
`roundfix: warning: Shared Build Cache <path>: trim event not journaled:
<error>; the Run continues`. With no declared cache it does nothing. The gate,
first in `commandDeliveryWorkflow.Gate`: the same call under the same budget,
writing lines to `workflow.log`. Neither returns an error from the trim.

Console lines (stderr at Run start, the Queue Owner's log at the gate), with
`P` = `at Run start` or `before the repository gate`, `C` = `the Run
continues` or `the gate continues`, and the prefix `Delivery Queue item
<slug>: ` only at the gate; `kept` and `unchanged` print nothing:

```text
roundfix: [prefix]Shared Build Cache <path>: removed <n> idle entries (<bytes> bytes) P
roundfix: [prefix]Shared Build Cache <path>: trim stopped after its <budget> budget P; removed <n> idle entries (<bytes> bytes); the next trim continues
roundfix: warning: [prefix]Shared Build Cache <path>: trim failed P: <cause>; C
```

The `disk-space:` line, after `storage:` and before `linux-lane:`/`codex:`,
printed only with at least one finding; findings follow location order (the
worktree location, then each declared cache), and `next: roundfix gc` appears
once when any `DR-DISK-LOW` is present:

```text
disk-space: warn (DR-DISK-LOW: worktree location <path> has <n> MiB free, below 10240 MiB; DR-DISK-UNREAD: could not read free space for Shared Build Cache <path>: <error>; next: roundfix gc)
```

The worktree location is `roundconfig.ResolveWorktreeLocation(worktree.location,
GitRoot, HomeDir)`; each cache path resolves as GC resolves it. The free
space is read at the nearest existing ancestor of the path; `<path>` is the
resolved location itself; a location that does not resolve (configuration
loading already validates both, so this is close to unreachable) reports
`DR-DISK-UNREAD` naming its configured value and the error. `<n>` is the free bytes divided by 1,048,576,
rounded down; a location is low when it has fewer than 10,737,418,240 bytes.

### Documentation

task_03 writes these texts. The GC guide (`docs/user-guide/commands/gc.md`)
and the canonical skill reference `references/storage.md` each gain, as the
last part of their `## Shared Build Caches` section, a subsection
`### Automatic trims` that holds paragraph 1, the console-line block of
Interfaces in a `text` fence, and paragraph 2:

```markdown
Roundfix also trims every declared Shared Build Cache by itself, with the
rules above: at the start of every Implement Run, after the Run is recorded
and before its Run Worktree exists, and before each repository gate of the
Delivery Queue. One budget of thirty seconds covers all declared caches at
each trigger; when it runs out, the trim keeps what it removed and the next
trigger continues. A refusal, an error or the end of the budget never fails,
parks or delays the Run or the gate beyond the budget. Each Implement Run
records one `daemon.status` Run Event per declared cache whose payload holds
`event: shared_cache_trim`, `trigger: run_start`, `path`, `outcome`
(`trimmed`, `unchanged`, `kept`, `failed` or `stopped`), `size_bytes`,
`limit_bytes`, `removed_entries`, `removed_bytes` and, when there is one,
`cause`. A trim that removed entries, failed or stopped also writes a line,
on standard error at Run start and in the Queue Owner's log at the gate:
```

```markdown
Without a declared cache nothing runs and no line or event is written.
`roundfix gc` still trims on demand (ADR-0334).
```

The GC guide's Retention policy row for Shared Build Cache changes its last
cell to `the GC Command, each Implement Run start and each repository gate of
the Delivery Queue`. The configuration guide's `worktree.shared_caches` key
row ends `for GC, each Implement Run start and the queue's repository gate to
trim above their size limits.`, and its `## Shared Build Caches` section gains,
after the paragraph that names `roundfix gc --dry-run`:

```markdown
The start of every Implement Run and every repository gate of the Delivery
Queue apply the same trim by themselves under a thirty-second budget; see
[GC](commands/gc.md#automatic-trims).
```

The implement guide (`docs/user-guide/commands/implement.md`, after the
paragraph ending "Run Worktree and Run Branch for recovery.") and the
reference `references/implement.md` (after the paragraph ending "Worktree and
Run Branch for inspection and recovery.") gain:

```markdown
When `worktree.shared_caches` declares a Shared Build Cache, the Run trims it
right after the Run is recorded and before its Run Worktree exists, with the
GC Command's rules and a thirty-second budget, and records one
`shared_cache_trim` Run Event per declared cache. The trim never changes the
Run's outcome; the GC guide's Shared Build Caches section describes it.
```

The deliver guide (`docs/user-guide/commands/deliver.md`, directly before
`#### Linux Lane`) and the reference `references/deliver.md` (directly before
`## Linux Lane`) gain:

```markdown
Before the repository gate runs its Verification in an item worktree, the
queue trims each declared Shared Build Cache with the GC Command's rules under
a thirty-second budget, and writes a line to the Queue Owner's log when the
trim removed entries, failed or stopped. The gate's result still comes only
from its Verification and the Linux Lane, so a `deliver retry` that resumes at
the gate trims first.
```

The doctor guide (`docs/user-guide/commands/doctor.md`, after the `storage:`
bullet) and the reference `references/setup.md` (before the `linux-lane:`
bullet) gain the bullet below; both code tables gain the two rows after
`DR-LIGHT-ACCESS-MODE`, and the reference's sentence "the conditional
`light-tier` and `linux-lane` checks" names `disk-space` between them:

```markdown
- `disk-space:` — printed after `storage:` only when the filesystem holding
  the worktree location (`worktree.location`) or a declared Shared Build Cache
  has less than 10 GiB (10,240 MiB) free, or its free space cannot be read.
  Each location is named on its own and read at its nearest existing
  directory: `DR-DISK-LOW` gives its free MiB with `next: roundfix gc`, and
  `DR-DISK-UNREAD` gives the read error. Findings are `warn` and never change
  the exit code; with enough free space everywhere the line is absent.
```

```markdown
| `DR-DISK-LOW` | warn | Less than 10,240 MiB free where worktrees or a Shared Build Cache live; `roundfix gc`, then free space on that filesystem. |
| `DR-DISK-UNREAD` | warn | The free space could not be read; check the path, then re-run `roundfix doctor`. |
```

`CONTEXT.md` **Shared Build Cache** keeps its `_Avoid_` line and becomes:

```markdown
A build cache several worktrees share outside any one of them, declared in
`worktree.shared_caches`; the GC Command, the start of every Implement Run and
the Delivery Queue's repository gate trim it above its declared size by
removing only cache entries no build used for the declared idle time, oldest
first ([ADR-0327](docs/adr/0327-gc-trims-a-declared-shared-build-cache.md),
[ADR-0334](docs/adr/0334-every-implement-run-and-the-queue-gate-trim-a-declared-shared-build-cache.md)).
```

### Data Models

The Run Event: `Source: daemon`, `Kind: daemon.status`, `Summary`
`Shared Build Cache <path>: <outcome>` (bounded with `runevent.BoundSummary`),
payload:

```json
{"event":"shared_cache_trim","trigger":"run_start","path":"<resolved path>","outcome":"trimmed","size_bytes":0,"limit_bytes":0,"removed_entries":0,"removed_bytes":0,"cause":"<only when not empty>"}
```

`sharedCacheReport` gains `outcome` and `cause`. `gcDependencies` gains
`sharedCacheTrimBudget time.Duration`. `commandDependencies` gains
`availableBytes func(string) (int64, error)` (default `store.AvailableBytes`)
and `doctorDependencies` gains `diskSpace func(context.Context,
roundconfig.Loaded) []CheckResult`. `HealthCheckDiskSpace = "disk-space"`.

### API Contracts

1. API Contract: an Implement Run with at least one declared cache trims each
   cache before its Run Worktree exists, appends exactly one `daemon.status`
   Run Event per cache with `event: shared_cache_trim`, and its exit code and
   outcome are those it has without the cache, whatever the trim outcome.
2. API Contract: the repository gate trims each declared cache before its
   Verification starts, writes the console lines to the Queue Owner's log,
   and its `GateResult` comes only from the Verification and the Linux Lane.
3. API Contract: every trim keeps ADR-0327's rules (limit, idle, oldest
   first, reread, refusal) and stops within the budget, keeping what it
   removed.
4. API Contract: with no declared cache, implement, the gate and GC write
   exactly what they write today, and no Run Event is added.
5. API Contract: `roundfix gc` and `roundfix gc --dry-run` keep their report
   bytes and exit codes.
6. API Contract: `roundfix doctor` prints `disk-space:` only with a finding,
   always `warn`, never changes the exit code, and reads free space only
   through the `availableBytes` command dependency.

### Surface Transcripts

None. The new console, log and doctor lines are fixed in Interfaces and
asserted, line for line, by the named tests of Testing Approach; no command
gains a flag, an argument or an exit code.

### Invariants

```text
1. No automatic trim removes a file GC would not remove with the same configuration at the same moment.
2. No trim outcome changes an Implement Run's exit code or outcome, or a repository gate's result.
3. Every automatic trim of an Implement Run is recorded as one Run Event per declared cache.
4. Without a declared cache, no trim runs, no event is written and no line is printed.
5. The disk-space line never fails Doctor and never appears while every location has at least 10 GiB free.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Features 1 and 2, Success Metrics 1 and 3 → API Contracts 1-3, Interfaces (`trimSharedCaches`, Run start, gate), Invariant 1, task_01
- Goal 2, Core Feature 2, Success Metric 2 → API Contracts 1-3, Interfaces (outcome rules, budget), Invariant 2, task_01
- Goal 3, Core Feature 2, Success Metrics 1-3 → API Contracts 1-2, Data Models (Run Event), Interfaces (console lines), Invariant 3, task_01
- Goal 4, Core Feature 3, Success Metric 5 → API Contract 6, Interfaces (`disk-space:` check), Data Models (`AvailableBytes`), Invariant 5, task_02
- Goal 5, Success Metrics 2, 4 and 5 → API Contracts 4-6, Invariants 4-5, task_01, task_02
- Core Feature 4, Success Metric 6 → Documentation (guides, references, glossary), task_03

## Absence Evidence

1. A context-bounded trim — command: `grep -rn "context" internal/sharedcache/sharedcache.go`; result: no match at `d08f8557`.
2. An automatic trim trigger — command: `grep -rn "sweepSharedCaches\|sharedcache\." internal --include=*.go | grep -v _test`; result: only `runGC` and `gc_shared_cache.go` at `d08f8557`.
3. A doctor free-space line — command: `grep -rn "Statfs\|GetDiskFreeSpace\|disk-space\|DR-DISK" internal docs/user-guide .agents/skills`; result: only `internal/store/compaction_capacity_unix.go` and `compaction_capacity_windows.go`, unexported, at `d08f8557`.

## ADR Alternatives

- ADR-0327: extends it; its rejected alternatives (a Makefile target, emptying
  the cache with `go clean -cache`, discovering the cache from
  `go env GOCACHE`) stay rejected, since the automatic trim calls the same
  engine on the same declaration with the same idle and refusal rules.
- ADR-0334: records this design; its rejected options (manual only, a timer
  in a long-lived process, a trim per Task, `resolve` and `watch` starts, a
  start floor) are the ones this TechSpec does not build.

## Integration Points

None. No network, provider, forge or other repository is reached; the trim
reads and removes files under a declared cache, and doctor reads free space
through the operating system.

## Testing Approach

Every test uses temporary homes, repositories and fake caches; no test fills
a disk or reads the host's free space.

- `internal/sharedcache/sharedcache_test.go` (declared change: two tests
  added, existing tests unchanged): `TestMeasureContextStopsWhenTheContextEnds`
  (a cancelled context returns an error matching `context.Canceled` and
  removes nothing) and `TestApplyContextStopsBetweenEntries` (a `now` that
  cancels the context on its first call: exactly one entry removed, the
  error matches `context.Canceled`, the other planned entries remain).
- `internal/cli/shared_cache_trim_test.go` (new), with the fixture of
  `gc_shared_cache_test.go` (a `{home}/gocache` cache, README, three 1 MiB
  entries aged 10 h, 3 h and 0, `idle: 2h`):
  - `TestImplementStartTrimsTheDeclaredSharedBuildCache` (limit 1): exit
    `0`, the two idle entries gone, the fresh entry and README kept, the
    stderr line `removed 2 idle entries (2097152 bytes) at Run start`, one
    `daemon.status` event with outcome `trimmed`, `removed_entries` 2,
    `removed_bytes` 2097152, `limit_bytes` 1048576, sequenced before the
    Run's first `daemon.task` event;
  - `TestImplementStartRunsOnWhenTheTrimDoesNot`, subtests `refused cache`
    (no README: outcome `kept`, cause names the README, every entry kept, no
    trim line) and `budget ends` (`sharedCacheTrimBudget` 1 ns: outcome
    `stopped`, cause `budget of 1ns elapsed`, the stopped line): each exits
    `0` with the same Task lines as `TestImplementStartWithoutASharedBuildCacheRecordsNoTrim`;
  - `TestImplementStartWithoutASharedBuildCacheRecordsNoTrim`: no
    `shared_cache_trim` event and no `Shared Build Cache` text;
  - `TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification`,
    subtests `idle entries removed first` (a gate whose Verification is
    `test ! -e '<10 h entry path>'` passes, and the Queue Owner's log holds
    `roundfix: Delivery Queue item <slug>: Shared Build Cache <path>: removed
    2 idle entries (2097152 bytes) before the repository gate`) and `budget
    ends` (`sharedCacheTrimBudget` 1 ns: the gate with Verification `true`
    passes and logs the stopped line).
- `internal/cli/doctor_disk_space_test.go` (new):
  `TestDoctorWarnsWhenALocationIsLowOnDisk` through `runCLI(doctor)` with
  `withDoctorFakeLoadedAndReadiness`, the default `diskSpace` and a fake
  `availableBytes`, subtests `below the threshold` (10,239 MiB for both
  locations: the exact line with two `DR-DISK-LOW` findings and one
  `next: roundfix gc`, between `storage:` and `codex:`, exit `0`), `at the
  threshold` (exactly 10,240 MiB: output equal to the same run without the
  dependency), `unreadable` (`DR-DISK-UNREAD` with the error, no next) and
  `missing location` (the fake is asked for the nearest existing ancestor).
- `internal/store/available_bytes_test.go` (new):
  `TestAvailableBytesReadsTheFilesystemOfAPath` (positive for `t.TempDir()`,
  an error naming the path for a missing one).
- `internal/cli/cli_test.go` (declared change): `commandEnvironmentForTest`
  and `updateCommandDependenciesForTest` install an `availableBytes` fake
  that reports 1 TiB, so no existing doctor test reads the host's disk;
  `withDoctorLiveDeps` in `doctor_test.go` is the only existing helper that
  keeps the default `diskSpace`, and it then reads the fake.
- Unchanged and run: the four GC Shared Build Cache tests, the four
  `internal/sharedcache` tests, `TestRunRetentionAtStartWarnsAndNeverBlocks`,
  `TestImplementStartRunsRunRetention`, `TestGateReturnsTheLaneEnvironmentSkip`,
  `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn` and the doctor
  output tests of `doctor_test.go`.
- Sabotages the QA gate runs, each with the case whose input makes it
  observable: the Run-start call removed (`TestImplementStartTrims…` — the
  10 h entry remains); the trim outcome turned into an implement preflight
  failure (`…RunsOnWhenTheTrimDoesNot/refused_cache` — exit `2`); the budget
  context dropped (`…/budget_ends` — outcome `trimmed`); the idle check
  dropped from both `MeasureContext` and `ApplyContext`
  (`TestImplementStartTrims…` — the 0 h entry is removed; the engine refuses
  an idle time below two hours, so a zero idle alone only yields `kept`); the Run Event append skipped
  (`TestImplementStartTrims…`); the gate trim moved after the Verification
  (`TestDeliveryGateTrims…/idle_entries_removed_first` — its Verification
  sees the 10 h entry); the doctor threshold compared with `<=`
  (`TestDoctorWarns…/at_the_threshold`); the next action printed per finding
  (`…/below_the_threshold`); the configured path read instead of its
  nearest existing ancestor (`…/missing_location`).

## Build Order

1. The guides, the skill references, the glossary revision, the skill
   version and the Behavior Surface record re-recorded last (task_03; no
   dependency: its Verification reads only documents and records).
2. The bounded engine, the trim function, the Run-start and gate triggers and
   their tests (task_01; depends on: 1, which documents the `gc`,
   `implement` and `deliver` surfaces it changes; they share no file).
3. `AvailableBytes`, the `disk-space:` check and its tests (task_02; depends
   on: 1, which documents the `doctor` surface it changes; it shares no file
   with 1 or 2).
4. QA gate (task_04) (depends on: 1, 2, 3).

No Spec is listed in `requires`. Specs 0319 and 0320 deliver first and
nothing here reads their code.

## Risks & Considerations

- An Implement Run's start grows by the measure walk (about 1.2 s for 21.6 GB)
  on every Run with a declared cache, and up to thirty seconds when the trim
  removes many entries. A detached Run's parent returns before the trim.
- Two Runs starting at once may plan the same entries; `ApplyContext`
  treats a missing entry as removed by someone else and the go command
  treats it as a miss, so the overlap costs time, never correctness.
- The free-space read is advisory and reads the nearest existing ancestor;
  a location on a mount that does not exist yet reports its parent's
  filesystem.
- `deliver.md` keeps its text: the gate's log lines are documented in the
  GC guide and the skill's `storage` reference, where the trim is described.

## Glossary

- changes: **Shared Build Cache**

## Decisions

- Triggers: Implement Run start after `CreateRun`, and the repository gate
  before its Verification. See ADR-0334.
- One function, `trimSharedCaches`, serves GC, the Run start and the gate.
- Budget thirty seconds, fixed, a `gcDependencies` field so tests can shorten
  it.
- Records: `daemon.status` Run Events with a `shared_cache_trim` payload; log
  lines at the gate.
- Doctor: `disk-space:` line only with a finding, 10 GiB, reusing the
  `gc compact` reader through an exported `store.AvailableBytes`.
- No start floor, no new park class, no configuration key.
