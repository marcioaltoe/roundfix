---
spec: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
prd: _prd.md
created: 2026-11-06
---

# The build cache stays within a bound, and a full disk is named — Technical Spec

## Executive Summary

The Makefile exports a shared `GOCACHE` that every Run on this repository
fills by a full build per worktree path, and nothing trims it between Run
starts and gates, so a busy day took it to 25 GB and the disk to zero. This
fix adds the Build Cache Bound: before each target that builds or tests, the
top-level make runs `go run ./cmd/gocache-bound`, which applies the Shared
Build Cache rules of `internal/sharedcache` (idle entries only, oldest first,
at least two idle hours) above 10,240 MiB. Doctor's Free-Space Floor also
inspects the temporary directory, the Delivery Queue applies the same floor
before an item's Implement Run and its repository gate as a `low-disk`
warning, and a full filesystem during the item binary build or the gate parks
the item as a Disk-Full Park, `disk-full: <path>: <N> MiB free`, whose retry
is not counted. ADR-0370 records the decision.

The trade-off accepted: the bound is soft. An entry a build used in the last
two hours is never removed, because a clean during a concurrent build can fail
that build (golang/go issue 31948), so on a very busy hour the cache can stay
above its bound until entries age; the floor's warning and the named park
cover that window.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; `disk-full` and `low-disk` are strings in existing item fields.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0370 (new) records the decision:
  "The Makefile bounds its build cache, and the queue warns of low disk and
  names a full one". It amends two accepted records, whose notices task_01
  (ADR-0327) and task_03 (ADR-0334) write: ADR-0327, whose first Considered Option rejected "A Makefile target
  that trims the cache" and whose second rejected `go clean -cache` (kept);
  and ADR-0334, whose last Considered Option deferred a free-space floor for
  Runs and whose Consequences say `roundfix doctor` "now warns about low free
  space". ADR-0334 cites ADR-0327. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 3 answers of 2026-11-06 quoted in the authorization record;
  the Governed Paths are the Makefile, `.roundfixrc.yml` and
  `docs/agents/setup-context.json` and the repository rules guide (task_01),
  and the Roundfix skill files tasks 02 and 03 change; no workflow, module or lint
  configuration changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/setup.md`, `.roundfixrc.yml`, `Makefile`, `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-06 at `6ceb7ff5`; `_prd.md` → Validation records the runs.

1. The Makefile's first block sets `GOCACHE` to
   `$TMPDIR/roundfix-go-build-<uid>`, or `$(CURDIR)/.gocache` when that
   directory cannot be written, only when `GOCACHE` is undefined, and
   exports it. No target trims it. `.roundfixrc.yml` declares it under
   `worktree.shared_caches` with `max_size_mb: 20480` and `idle: 4h`.
2. `internal/sharedcache` (`MeasureContext`, `ApplyContext`) plans and
   applies the ADR-0327 trim; `MinimumIdle` is two hours. Its callers are
   `trimSharedCaches` (`internal/cli/gc_shared_cache.go`): `roundfix gc`, the
   Run start (`trimSharedCachesAtRunStart`) and the gate (`Gate` in
   `internal/cli/deliver_workflow.go`).
3. `defaultDoctorDiskSpaceResults` (`internal/cli/doctor_disk_space.go`)
   reads, at each location's nearest existing directory
   (`doctorDiskSpaceAncestor`), the free space of the worktree location and
   each declared Shared Build Cache through the `availableBytes` command
   dependency, and below `doctorDiskSpaceThreshold` (10,240 MiB) prints
   `DR-DISK-LOW: <label> <path> has <N> MiB free, below 10240 MiB` with
   `next: roundfix gc`. `TestDoctorWarnsWhenALocationIsLowOnDisk` pins the
   two reads and the line.
4. The Delivery Queue never reads free space.
5. `stepExecutable` builds the Item Binary through `daemon.ExecVerifier`,
   which keeps the output only in `item-binary-build.log`, and returns
   `build item binary: verification command "make build" failed;
   diagnostics: <log>: exit status <n>`. The engine's `Run` loop and
   `reviewSlots.review` park any step error as `delivery-error: <error>`,
   environment Park Class. `Gate` returns a failed `GateResult` without a
   blocker, which `gateCandidate` parks as `gate-failed`, finding Park Class;
   an error that is not a command failure becomes `delivery-error`.
6. `Retry` (`internal/delivery/engine.go`) and `RetryDeliveryQueueItem`
   (`internal/store/delivery.go`) exempt only a `runtime-infrastructure`
   blocker from the retry limit and the retry count.
7. `internal/runcause/signatures.json` already classifies `no space left on
   device` in a Run's diagnostic as the `disk-full` signature,
   environment class.

## System Architecture

New: the `cmd/gocache-bound` command (a thin `main` calling
`sharedcache.BoundMain`), the `gocache-bound` make target, the
`delivery.FreeSpaceChecker` seam with its CLI implementation, and the
`delivery.DiskFullError` type with its CLI detector. Changes by component:

- `Makefile`, `.roundfixrc.yml`, `docs/agents/setup-context.json`,
  `internal/sharedcache/bound.go`, `cmd/gocache-bound/main.go` (task_01).
- `internal/cli/doctor_disk_space.go` (the shared location inspection and
  the temporary directory), `internal/cli/deliver_free_space.go`,
  `internal/delivery/engine.go` (the seam and the `low-disk` segment),
  `internal/cli/deliver_workflow.go` (wiring) (task_02).
- `internal/delivery/disk_full.go`, `internal/delivery/park_class.go`,
  `internal/delivery/engine.go` and `internal/delivery/review_slots.go` (the
  two `delivery-error` parks and the retry limit), `internal/store/delivery.go`
  (the retry count), `internal/cli/deliver_disk_full.go` and
  `internal/cli/deliver_workflow.go` (`stepExecutable`, `Gate`) (task_03).
- Documentation: each Task documents its own behavior (API Contract 8).

## Implementation Design

### Interfaces

```go
// internal/sharedcache/bound.go
const BoundBudget = 30 * time.Second
// BoundMain parses -dir, -max-mb and -idle, trims dir as MeasureContext and
// ApplyContext plan, and returns 2 for a usage error and 0 otherwise.
func BoundMain(args []string, stderr io.Writer, now func() time.Time) int

// internal/delivery
type FreeSpaceChecker interface {
	// LowFreeSpace names each inspected location below the Free-Space Floor,
	// as "<label> <path> has <N> MiB free, below 10240 MiB".
	LowFreeSpace(ctx context.Context) []string
}
const WarningLowDisk = "low-disk"
const BlockerDiskFull = "disk-full"
type DiskFullError struct {
	Path           string
	AvailableBytes int64
	Err            error
}
func (err *DiskFullError) Error() string   // err.Err.Error()
func (err *DiskFullError) Unwrap() error
func (err *DiskFullError) Blocker() string // "disk-full: <path>: <N> MiB free"

// internal/cli
type freeSpaceLocation struct{ label, path string; available int64; err error }
func inspectFreeSpace(ctx context.Context, loaded roundconfig.Loaded) []freeSpaceLocation
func diskFullFailure(ctx context.Context, loaded roundconfig.Loaded, err error, logPath string) *delivery.DiskFullError
```

### Data Models

No schema change. The Delivery Queue item's existing `Blocker` and `Warning`
strings carry the new `disk-full` blocker and `low-disk` segment.

### API Contracts

1. API Contract: the make wiring.
   - Inside the existing `ifeq ($(origin GOCACHE),undefined)` block, after the
     cache is chosen: `GOCACHE_BOUND_MB ?= 10240`. After the block:
     `GOCACHE_BOUND_MB ?= 0` and `GOCACHE_BOUND_IDLE ?= 2h`, with a comment
     naming ADR-0370.
   - A phony target `gocache-bound`, help text `Trim the Makefile's Go build
     cache above GOCACHE_BOUND_MB, idle entries only`, whose recipe runs only
     when `$(MAKELEVEL)` is `0` and `GOCACHE_BOUND_MB` is not `0`:
     `$(GO) run $(RUN_FLAGS) ./cmd/gocache-bound -dir "$(GOCACHE)" -max-mb
     "$(GOCACHE_BOUND_MB)" -idle "$(GOCACHE_BOUND_IDLE)"`, and when that
     exits non-zero prints to standard error
     `gocache-bound: warning: the Build Cache Bound did not run (exit <n>); the build continues`
     and succeeds.
   - `gocache-bound` is the first prerequisite of `vet`, `test`,
     `test-budget`, `test-race`, `build`, `run`, `version`,
     `verify-changed-core`, `verify-changed-baseline`,
     `verify-changed-contracts`, `verify-contracts`, `docs-test`,
     `repo-test`, `spec-budget` and `baseline-digests`; `install`, `verify`,
     `verify-changed`, `verify-incremental` and `verify-docs` reach it
     through those. Every other recipe keeps its bytes.
2. API Contract: `gocache-bound`. Flags `-dir` (absolute, required),
   `-max-mb` (a positive MiB count whose bytes fit an `int64`) and `-idle`
   (a Go duration of at least `sharedcache.MinimumIdle`); any other argument,
   or a value outside these, prints one line starting `gocache-bound: ` and
   exits `2`. Otherwise it measures and trims under `BoundBudget` with
   `MeasureContext` and `ApplyContext`, exits `0`, and writes to standard
   error only:
   - when it removed entries:
     `gocache-bound: <dir>: removed <n> idle entries (<bytes> bytes)`;
   - when the cache is still above the bound after the trim:
     `gocache-bound: warning: <dir> holds <bytes> bytes, above its <limit>-byte bound, and no other entry is idle for <idle>; the build continues`,
     `<idle>` printed by `time.Duration.String`;
   - when the budget ends:
     `gocache-bound: warning: <dir>: stopped after its 30s budget; the next make continues`;
   - on a refusal or another error:
     `gocache-bound: warning: <error>; the build continues`.
   Within the bound it writes nothing.
3. API Contract: the declaration. `.roundfixrc.yml`'s
   `{tmpdir}/roundfix-go-build-{uid}` entry becomes `max_size_mb: 10240`
   and `idle: 2h`, and its comment names the Makefile's bound.
4. API Contract: Doctor's `disk-space:` line. `inspectFreeSpace` returns,
   in order, the worktree location, each declared Shared Build Cache and the
   temporary directory (`os.TempDir()`), each read at its nearest existing
   directory. `defaultDoctorDiskSpaceResults` renders them as today, with the
   label `temporary directory` for the last; its findings, threshold, codes,
   `next: roundfix gc` and exit behavior keep their bytes.
5. API Contract: the queue's `low-disk` warning. With
   `EngineDependencies.FreeSpace` set, `runCandidate` before `RunSpec` and
   `gateCandidate` before `Gate` call `LowFreeSpace`. A non-empty result
   becomes the one segment `low-disk: <locations joined by ", ">` of the
   item's `Warning`, replacing an earlier `low-disk:` segment and placed
   before a `linux-lane-skipped:` segment, which stays last; the item is
   stored, and the Queue Owner's log gets
   `roundfix: warning: Delivery Queue item <slug>: low-disk: <locations> before its Implement Run; the queue continues`
   (or `before its repository gate`). An empty result removes the segment.
   A location whose free space cannot be read adds nothing. The check never
   parks, fails or delays the item, and a nil checker changes nothing.
6. API Contract: the Disk-Full Park. `diskFullFailure` returns a
   `DiskFullError` when `errors.Is(err, syscall.ENOSPC)` or the last 64 KiB
   of `logPath` contain `no space left on device` (any case). Its `Path` is
   the last absolute path written directly before `: no space left on
   device` in that log tail, else in the error text, else the
   `inspectFreeSpace` location with the least free space; its
   `AvailableBytes` is read at the path's nearest existing directory.
   - `stepExecutable` wraps it: `build item binary: <the same text>`, so the
     error text keeps its bytes.
   - `Gate` returns `GateResult{Passed: false, Reason: <error>, Blocker:
     <Blocker()>}` for a Verification failure or an `ENOSPC` error.
   - The engine parks every step error that carries a `DiskFullError`, in
     `Run` and in `reviewSlots.review`, with its `Blocker()` instead of
     `delivery-error: <error>`.
7. API Contract: classification and retry. `ClassifyPark` puts `disk-full`
   in the environment Park Class with the next action
   `free space on the filesystem that holds the path the blocker names: run roundfix gc to trim declared Shared Build Caches and reclaim ended Runs, and remove other files there; then run roundfix deliver retry <slug>; a retry from disk-full is not counted`.
   `Retry`'s limit check and `RetryDeliveryQueueItem`'s increment exempt
   `disk-full` as they exempt `runtime-infrastructure`; the resumed stage
   follows the item's evidence as for any other blocker.
8. API Contract: the rule in writing. Line wrapping may differ; the words
   may not. task_01 writes the README, the repository rules guide, **Build
   Cache Bound** and the ADR-0327 notice; task_02 the `doctor` guide and
   `setup` reference, the `low-disk` paragraph of the `deliver` documents and
   **Free-Space Floor**; task_03 the rest of the `deliver` documents,
   **Disk-Full Park** and the ADR-0334 notice.
   - `docs/user-guide/commands/deliver.md` and the canonical
     `.agents/skills/roundfix/references/deliver.md`:
     - the `environment` row of the Park Class table lists `` `disk-full` ``
       after `` `runtime-infrastructure` ``;
     - after the paragraph starting "A failed build, a declared path Git does
       not ignore": "When the item binary build or the repository gate fails
       because a filesystem is full, with the error `ENOSPC` or a retained
       log that says `no space left on device`, the item parks as
       `disk-full: <path>: <N> MiB free` in the environment Park Class
       instead of `delivery-error` or `gate-failed`. `<path>` is the last
       absolute path the log names before that text, or else the inspected
       location with the least free space, and `<N>` is its free space when
       the item parked. A Run whose Task Verification fills the disk still
       parks `run-unresolved`; `roundfix runs causes <run>` names its
       `disk-full` signature.";
     - once each: "A `disk-full` park's next action is `<API Contract 7's
       next action>`." and, after the `runtime-infrastructure` retry
       sentence, "A Delivery Retry from `disk-full` is not counted against
       the retry limit either.";
     - after the paragraph starting "Before the repository gate runs its
       Verification": "Before an item's Implement Run and before its
       repository gate, the queue reads the free space where the worktree
       location, each declared Shared Build Cache and the temporary
       directory live. Below the Free-Space Floor of 10,240 MiB it records a
       `low-disk: <label> <path> has <N> MiB free, below 10240 MiB` segment
       in the item's warning, which `roundfix deliver status` prints, and
       writes one warning line to the Queue Owner's log; at or above the
       floor the segment goes. The warning never parks or delays the item."
   - `docs/user-guide/commands/doctor.md` and the canonical
     `.agents/skills/roundfix/references/setup.md`: in the `disk-space:`
     bullet, "the worktree location (`worktree.location`) or a declared
     Shared Build Cache has less than 10 GiB (10,240 MiB) free" becomes "the
     worktree location (`worktree.location`), a declared Shared Build Cache
     or the temporary directory has less than the Free-Space Floor of 10 GiB
     (10,240 MiB) free"; the `DR-DISK-LOW` row says "where worktrees, a
     Shared Build Cache or the temporary directory live".
   - `README.md`, the Development section, after the `make verify` block:
     "Every `make` target that builds or tests first holds the Makefile's Go
     build cache to its Build Cache Bound, `GOCACHE_BOUND_MB` (10240 MiB), by
     removing entries no build used for `GOCACHE_BOUND_IDLE` (2h); set
     `GOCACHE_BOUND_MB=0` to turn it off. To keep the cache off a small
     system disk, point `TMPDIR` or `GOCACHE` at another volume; Go's own
     user cache, which `go` commands outside `make` fill, can be declared in
     User Config under `worktree.shared_caches`."
   - `docs/agents/specific-repository.md`, the rule starting "`make` shares
     `roundfix-go-build-<uid>`" gains: "before it builds or tests, the
     top-level `make` holds the cache it chose to the Build Cache Bound with
     `cmd/gocache-bound` and never clears it (ADR-0370)."
   - `CONTEXT.md`, three entries beside **Shared Build Cache**, each
     `ADR-0370` written as the Markdown link to its file:

     ```markdown
     **Build Cache Bound**:
     The size, 10,240 MiB by default, to which this repository's Makefile holds the Go build cache it chose before each top-level make that builds or tests, removing with the Shared Build Cache rules only entries no build used for two hours; the cache stays above it while its entries are younger (ADR-0370).
     _Avoid_: cache limit, cache cleanup, gocache trim

     **Free-Space Floor**:
     The 10,240 MiB of free space below which `roundfix doctor` reports `DR-DISK-LOW` and the Delivery Queue records a `low-disk` warning for the filesystem of the worktree location, a declared Shared Build Cache or the temporary directory; it never fails a check or parks an item (ADR-0370).
     _Avoid_: disk threshold, low disk check

     **Disk-Full Park**:
     A Delivery Queue park of an item whose Item Binary build or repository gate failed on a full filesystem, as `disk-full: <path>: <N> MiB free` in the environment Park Class; its next action names what to free, and its Delivery Retry is not counted (ADR-0370).
     _Avoid_: disk error, out of space park
     ```

   - The amendment notices, each a paragraph whose bold opening links
     `0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`,
     with `updated_at: 2026-11-06T12:00:00Z`:
     - ADR-0327, after its Decision's last paragraph, which starts "`roundfix
       gc --dry-run` reports":

       ```markdown
       **Amended by ADR-0370 (2026-11-06).** This repository's Makefile also applies these rules to the cache it chooses, before each top-level make that builds or tests, as the Build Cache Bound: the Makefile target this record's first Considered Option rejected now runs beside the product's trim, which adopters keep. Emptying the cache stays rejected.
       ```

     - ADR-0334, after its Decision's last bullet:

       ```markdown
       **Amended by ADR-0370 (2026-11-06).** The Delivery Queue now reads free space before an item's Implement Run and its repository gate and warns below the Free-Space Floor, and a full filesystem during the Item Binary build or the gate parks the item as a Disk-Full Park. Neither refuses to start a Run, so the deferred floor of the Considered Options stays deferred.
       ```

### Surface Transcripts

None. Doctor prints a dozen readiness lines that depend on the host's tools,
and the queue's lines depend on its log and store, so the changed lines are
fixed by API Contracts 2, 4, 5 and 7 and asserted by name in the tests of
Testing Approach.

### Invariants

```text
1. No command removes a cache entry whose modification time is less than two hours old, and none empties a cache.
2. A make target that does not build or test, and every recipe line other than the new prerequisite, keeps its bytes.
3. The bound never changes the exit status of make.
4. Doctor's exit code and every disk-space finding other than the new temporary directory keep their bytes at 6ceb7ff5.
5. A low-disk warning never parks, fails or delays an item, and never reaches the pull request body.
6. A step error that carries no DiskFullError parks with the same blocker as at 6ceb7ff5.
```

### Declared breaks

Measured with a prototype of task_01 on a scratch copy at `6ceb7ff5`, and by
reading the tests that pin the changed lines for tasks 02 and 03.

- Behavior, task_01: every listed make target first runs the bound, about a
  second more for each top-level make; the declaration's limit falls from
  20,480 MiB and four hours to 10,240 MiB and two hours.
- Behavior, task_02: Doctor's `disk-space:` line can name `temporary
  directory <path>`; `roundfix deliver status` can print a `low-disk`
  warning.
- Behavior, task_03: a full disk during the item binary build or the gate
  parks `disk-full: ...` in the environment Park Class instead of
  `delivery-error: ...` or `gate-failed` (finding Park Class), and its retry
  is no longer counted.
- Tests changed, task_01: `TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped`
  (`internal/verifyselect/verify_docs_makefile_test.go`) counts the bound's
  `run -buildvcs=false ./cmd/gocache-bound ...` invocation and requires
  exactly one.
- Tests changed, task_02: `TestDoctorWarnsWhenALocationIsLowOnDisk`
  (`internal/cli/doctor_disk_space_test.go`) expects the third read,
  `os.TempDir()`, and its finding, and every other test that pins the
  `disk-space:` line on the default branch when the Task runs.
- No test is removed or renamed, so the coverage record is not re-recorded.

## Coverage Map

- Goal 1, Goal 2, Core Feature 1, Success Metrics 1 and 2 → API Contracts
  1 to 3, `TestMakeBoundsItsOwnGoBuildCache`,
  `TestBoundMainTrimsOnlyIdleEntriesAboveTheBound`,
  `TestBoundMainRefusesAnInvalidRequest`, Invariants 1 to 3.
- Goal 3, Core Feature 2, Success Metric 3 → API Contract 4,
  `TestDoctorWarnsWhenALocationIsLowOnDisk`, Invariant 4.
- Goal 4, Core Feature 3, Success Metric 4 → API Contract 5,
  `TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate`,
  `TestCommandFreeSpaceNamesEachLocationBelowTheFloor`, Invariant 5.
- Goal 5, Core Feature 4, Success Metric 5 → API Contracts 6 and 7,
  `TestAFullDiskParksTheItemAsDiskFull`,
  `TestADiskFullRetryIsNotCounted`,
  `TestTheItemBuildOnAFullDiskIsADiskFullFailure`,
  `TestTheRepositoryGateOnAFullDiskParksDiskFull`, Invariant 6.
- Core Feature 5, Success Metric 6 → API Contract 8,
  `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem`.

Deviations from the PRD: None.

## Absence Evidence

1. `cmd/gocache-bound`, `BoundMain` and the `gocache-bound` target —
   command: `ls cmd; grep -rn 'BoundMain\|gocache-bound' Makefile internal cmd`;
   result: `cmd` holds `roundfix` and `verify-select` only, and no match.
2. A free-space read in the Delivery Queue — command:
   `grep -rn 'availableBytes\|AvailableBytes\|Statfs' internal/delivery internal/cli/deliver*.go`;
   result: no match; only Doctor and the store's compaction read free space.
3. A disk-full blocker or ENOSPC detection in delivery — command:
   `grep -rn 'disk-full\|ENOSPC\|no space left' internal/delivery internal/cli internal/store --include='*.go'`;
   result: no match; the only classifier is the `disk-full` signature in
   `internal/runcause/signatures.json`, which reads Run diagnostics.
4. The new test files — command:
   `ls internal/sharedcache/bound_test.go internal/verifyselect/build_cache_bound_makefile_test.go internal/cli/deliver_free_space_test.go internal/delivery/low_disk_warning_test.go internal/delivery/disk_full_park_test.go internal/store/delivery_disk_full_retry_test.go internal/cli/deliver_disk_full_test.go`;
   result: no such files.

## ADR Alternatives

- ADR-0370: the decision this Spec records; it rejects clearing the cache,
  removing young entries, `-trimpath`, a shell trim, refusing a Run below
  the floor and moving the cache by default, and the design takes none of
  them.
- ADR-0327: this Spec takes the alternative ADR-0327 rejected first, a
  Makefile target that trims the cache, for this repository only and beside
  the product's trim, so ADR-0370 amends it and task_01 writes the notice;
  its rejection of emptying the cache and of discovering the cache from `go
  env GOCACHE` stays.
- ADR-0334: ADR-0334 rejected a manual-only trim, a timer, a trim per Task
  and trims at `resolve` and `watch` starts, and deferred refusing a Run
  below a floor; the design takes none of them and keeps the deferral,
  adding only a warning, so ADR-0370 amends its record of where free space is
  read and task_03 writes the notice.

## Integration Points

None external. The bound reads one directory; the floor reads filesystem
statistics; no test reaches the network, a real Run Database, a real cache or
the real disk's free space.

## Testing Approach

Every test is hermetic: `t.TempDir()` caches and repositories, the
`availableBytes` command dependency injected as today, the compiled Go stub
of `TestMakeSharesOneGoBuildCache` (`testfixture.FixtureBinary`), the
engine's fake workflow and real queue store, and no test reads or fills a real
cache or the real disk.

- `internal/sharedcache/bound_test.go` (task_01):
  `TestBoundMainTrimsOnlyIdleEntriesAboveTheBound` with an injected clock and
  600 KiB entries under a README: rows `within its bound` (nothing removed,
  no output), `above its bound with idle entries` (ten-hour, five-hour and
  thirty-minute entries with `-max-mb 1`: the two idle ones removed, the
  exact removal line), `above its bound with no idle entry` (ninety- and
  thirty-minute entries: nothing removed, the exact warning line) and `not a
  Go build cache` (no README: nothing removed, the refusal warning);
  `TestBoundMainRefusesAnInvalidRequest`: a relative `-dir`, `-max-mb 0`,
  `-idle 1h` and an extra argument each exit `2`.
- `internal/verifyselect/build_cache_bound_makefile_test.go` (task_01):
  `TestMakeBoundsItsOwnGoBuildCache`, the real Makefile with a Go stub that
  appends its arguments to a record and fails `run` on request, an
  environment without `GOCACHE`, `MAKELEVEL`, `MAKEFLAGS` and the bound
  variables: rows `default`, `once per invocation` (`make vet test`),
  `fallback cache` (unwritable `TMPDIR`: `-dir <run dir>/.gocache`),
  `sub-make` (`MAKELEVEL=1`), `explicit cache`, `explicit cache with a
  bound` (`GOCACHE_BOUND_MB=4096`), `turned off` (`GOCACHE_BOUND_MB=0`),
  `idle override` (`GOCACHE_BOUND_IDLE=3h`) and `failing bound` (the warning
  line and exit `0`), each asserting the exact list of recorded invocations.
- `internal/cli/doctor_disk_space_test.go` (task_02): the existing table
  with the third read and finding.
- `internal/cli/deliver_free_space_test.go` (task_02):
  `TestCommandFreeSpaceNamesEachLocationBelowTheFloor`, rows `below the
  floor` (10,239 MiB: three segments in order), `at the floor` (none) and
  `unreadable` (none).
- `internal/delivery/low_disk_warning_test.go` (task_02):
  `TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate` with a fake checker per
  call: low before the Run (the stored segment and the log line), healthy at
  the gate (segment gone), low with a lane-skip warning present (the
  `low-disk` segment before it, `LaneSkipBody` unchanged), and no checker
  (warning unchanged); the item reaches `merged` in each.
- `internal/delivery/disk_full_park_test.go` (task_03):
  `TestAFullDiskParksTheItemAsDiskFull` (a Run step and a review step
  returning an error that wraps a `DiskFullError`, and a gate result with
  its blocker: the exact blocker, environment class and next action; an
  error without one: `delivery-error: ...` as today);
  `TestADiskFullRetryIsNotCounted` (the retry count unchanged and a retry
  accepted at the limit).
- `internal/store/delivery_disk_full_retry_test.go` (task_03): the store
  leaves `retry_count` unchanged for `disk-full` and still counts
  `delivery-error`.
- `internal/cli/deliver_disk_full_test.go` (task_03):
  `TestTheItemBuildOnAFullDiskIsADiskFullFailure` over
  `newDeliveryItemFixture`, rows `path in the log` (build command printing
  `write /tmp/x/_pkg_.a: no space left on device`: `Path` `/tmp/x/_pkg_.a`,
  the injected free space, error text starting `build item binary: `),
  `no path in the log` (`compile: writing output: write $WORK/b001/_pkg_.a:
  no space left on device`: the location with the least injected free space)
  and `other failure` (no `DiskFullError`);
  `TestTheRepositoryGateOnAFullDiskParksDiskFull` (a Verification printing
  the text: `Blocker` `disk-full: ...`; another failure: empty `Blocker`).
- `internal/docscontract/command_documentation_test.go` (task_03):
  `checkParkActionDocumentation` requires the `disk-full` next action once
  in each document.
- Sabotages, each on a scratch copy, with the case whose input makes it
  observable (the three Makefile ones were run by the author on the task_01
  prototype):
  - the `MAKELEVEL` condition dropped: `TestMakeBoundsItsOwnGoBuildCache`
    fails on `sub-make`;
  - `GOCACHE_BOUND_MB ?= 10240` moved outside the undefined-`GOCACHE`
    block: it fails on `explicit cache`;
  - the warning fallback removed from the recipe: it fails on `failing
    bound`;
  - `BoundMain` passing a zero limit: `TestBoundMainTrimsOnlyIdleEntriesAboveTheBound`
    fails on `within its bound`;
  - the temporary directory dropped from `inspectFreeSpace`:
    `TestDoctorWarnsWhenALocationIsLowOnDisk` fails;
  - the `low-disk` segment appended after the lane-skip segment:
    `TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate` fails on the lane-skip
    row;
  - the log tail scan removed from `diskFullFailure`:
    `TestTheItemBuildOnAFullDiskIsADiskFullFailure` fails on `path in the
    log`;
  - `disk-full` dropped from the store's exemption:
    `TestADiskFullRetryIsNotCounted` fails.

## Build Order

1. The Build Cache Bound: `internal/sharedcache/bound.go`,
   `cmd/gocache-bound/main.go`, the make target and prerequisites, the
   declaration, the two Makefile tests, the README and repository rules
   text, **Build Cache Bound** and the ADR-0327 notice, and the setup
   manifest regenerated last (task_01).
2. The Free-Space Floor in Doctor and the Delivery Queue, the `doctor` guide
   and `setup` reference, the `low-disk` paragraph of the `deliver`
   documents and **Free-Space Floor**, with the mirrors synced, the
   owned-skill version recorded and the Behavior Surface Record re-recorded
   last (task_02) (depends on: 1).
3. The Disk-Full Park, which uses task_02's `inspectFreeSpace`, the rest of
   the `deliver` documents, the documentation contract, **Disk-Full Park**
   and the ADR-0334 notice, with the same sync, record and re-record last
   (task_03) (depends on: 2).
4. QA gate (task_04) (depends on: 3).

The three Tasks each write one `CONTEXT.md` entry, and tasks 02 and 03 both
change the `deliver` documents, the Roundfix skill version, the engine and
the workflow, so they run in series. task_02 depends on task_01 only for
`CONTEXT.md`. No Verification depends on another Spec; Spec 0359 delivers
first and edits `checkCandidate`, which this Spec does not touch.

## Risks & Considerations

- The bound is soft: on a very busy hour the cache can exceed it until
  entries age. The floor's warning and the Disk-Full Park cover the window;
  the operator still has `roundfix gc`.
- `go run` of the bound needs a little free space; on a full disk it fails,
  prints its warning and the build fails on its own error, which the
  Disk-Full Park then names.
- Task Verification commands that call `go` outside `make` still fill Go's
  default user cache; the README names the User Config declaration that
  trims it.
- A log path regex can name a file inside the full filesystem's temporary
  tree; its nearest existing directory is read, so the free space is the
  right filesystem's.
- Spec 0359 also edits `internal/delivery/engine.go`,
  `internal/cli/deliver_workflow.go`, the `deliver` guide and reference,
  the Roundfix skill version and `CONTEXT.md`; the edits here are local to
  the functions named above, and the version record command picks the next
  free version at the run of tasks 02 and 03.

## Glossary

- adds: **Build Cache Bound**
- adds: **Free-Space Floor**
- adds: **Disk-Full Park**

## Decisions

- The bound reuses `internal/sharedcache` through a small command; no second
  trim implementation (ADR-0370).
- Idle entries only, two hours at least; the cache is never cleared
  (ADR-0370, golang/go issue 31948).
- 10,240 MiB for the bound and the floor; the declaration takes the bound's
  values so every trim of this cache agrees.
- The bound runs once per top-level make and fails open.
- One location inspection serves Doctor, the queue warning and the park's
  fallback path.
- The Disk-Full Park covers the queue's own builds; a Run's Task
  Verification keeps `run-unresolved` and `roundfix runs causes`.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found this repository's mirrored ADRs and the Mac mini storage
  page. Exa: golang/go issues 31948, 29561, 69879 and 81830.
- Three implementation Tasks in series and a QA gate.
