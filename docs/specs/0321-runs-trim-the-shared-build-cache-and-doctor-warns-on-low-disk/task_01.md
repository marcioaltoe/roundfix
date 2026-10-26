---
task: task_01
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
status: completed
type: backend
complexity: high
---

# Task 01: Implement Runs and the queue gate trim the declared Shared Build Caches

## Overview

Today only `roundfix gc` trims a declared Shared Build Cache. This Task makes
the trim engine stoppable through a context, wraps it in one function that
trims every declared cache under the thirty-second budget and classifies each
outcome, and calls it right after an Implement Run is recorded (with one
`daemon.status` Run Event per cache) and first in the Delivery Queue's
repository gate (with lines in the Queue Owner's log). Neither caller's
outcome ever depends on the trim, and GC's report keeps its bytes.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-the-shared-build-cache-fills-the-disk-between-manual-gc-runs.md`](references/2026-10-26-the-shared-build-cache-fills-the-disk-between-manual-gc-runs.md)
   of 2026-10-26 (`_prd.md` → Sources 1) by implementing API Contracts 1-5
   of `_techspec.md` and ADR-0334, keeping every trim rule of ADR-0327.
2. MUST add `MeasureContext` and `ApplyContext` to `internal/sharedcache` as
   `_techspec.md` → Interfaces states: both check the context between
   entries and return an error that matches the context's error with
   `errors.Is`; `ApplyContext` returns what it removed before stopping;
   `Measure` and `Apply` call them with `context.Background()` and keep
   their behavior.
3. MUST add `trimSharedCaches` and the outcome rules (`trimmed`,
   `unchanged`, `kept`, `failed`, `stopped`, with `cause`) to
   `internal/cli/gc_shared_cache.go`, rebuild `sweepSharedCaches` on it so
   `roundfix gc` and `roundfix gc --dry-run` keep their output and exit
   codes, and add `sharedCacheTrimBudget` (default thirty seconds) to
   `gcDependencies` in `internal/cli/gc.go`.
4. MUST, in `runImplementCommand`, after `reportDetachedRunCreated`
   succeeds and before `createRunWorktree`, when at least one cache is
   declared, run the trim under `context.WithTimeout` of the budget, append
   one Run Event per cache with the payload of `_techspec.md` → Data Models
   through `runStore.AppendRunEvent`, and write the console lines of
   Interfaces to standard error; an append failure prints one warning and
   the Run continues. With no declared cache nothing runs.
5. MUST, first in `commandDeliveryWorkflow.Gate`, run the same trim under
   the same budget and write the console lines, with the
   `Delivery Queue item <slug>: ` prefix, to `workflow.log`; the
   `GateResult` comes only from the Verification and the Linux Lane.
6. MUST add `TestMeasureContextStopsWhenTheContextEnds` and
   `TestApplyContextStopsBetweenEntries` to
   `internal/sharedcache/sharedcache_test.go`, keeping its existing tests
   unchanged, and add `internal/cli/shared_cache_trim_test.go` with
   `TestImplementStartTrimsTheDeclaredSharedBuildCache`,
   `TestImplementStartRunsOnWhenTheTrimDoesNot` (subtests `refused cache`
   and `budget ends`), `TestImplementStartWithoutASharedBuildCacheRecordsNoTrim`
   and `TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification`
   (subtests `idle entries removed first` and `budget ends`), with the
   inputs and assertions of `_techspec.md` → Testing Approach. Each top-level
   test calls `t.Parallel()` first or carries a `// Sequential: <reason>`
   comment.
7. MUST keep every existing GC Shared Build Cache, Run Retention at start,
   implement and repository gate test passing unchanged.
8. MUST NOT change `worktree.shared_caches` or its validation, any document,
   or the Makefile.

## Subtasks

- [ ] Make the trim engine stoppable through a context.
- [ ] Add the trim function with its outcomes and rebuild GC's sweep on it.
- [ ] Trim at the Implement Run start with its Run Events, and first in the repository gate.
- [ ] Add the engine, Run-start and gate tests.

## Acceptance Criteria

- [ ] An Implement Run above its cache limit removes exactly ADR-0327's idle entries before any Task starts and records one `shared_cache_trim` Run Event per cache.
- [ ] A refused cache or an exhausted budget leaves the Run's exit code and outcome unchanged and is recorded with its cause.
- [ ] The repository gate trims before its Verification and its result comes only from the Verification and the Linux Lane; GC's output is unchanged.

## Consumers

- `internal/sharedcache.MeasureContext` and `internal/sharedcache.ApplyContext` — consumers: `./internal/cli`
- the GC report of `sweepSharedCaches` and the `shared_cache_trim` Run Event — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`
- instruction: `docs/adr/0334-every-implement-run-and-the-queue-gate-trim-a-declared-shared-build-cache.md`
- instruction: `internal/cli/gc_shared_cache_test.go`
- instruction: `internal/cli/run_retention_start_test.go`
- instruction: `internal/cli/implement_test.go`
- instruction: `internal/cli/deliver_lane_skip_visible_test.go`
- instruction: `internal/store/run_events_of_kinds.go`
- interface: `internal/sharedcache/sharedcache.go`
- interface: `internal/sharedcache/sharedcache_test.go`
- interface: `internal/cli/gc_shared_cache.go`
- interface: `internal/cli/gc.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/shared_cache_trim_test.go`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestMeasure.*|TestApply.*|TestImplementStartTrimsTheDeclaredSharedBuildCache|TestImplementStartRunsOnWhenTheTrimDoesNot|TestImplementStartWithoutASharedBuildCacheRecordsNoTrim|TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification|TestGC.*SharedBuildCache.*|TestImplementStartRunsRunRetention|TestRunRetentionAtStartWarnsAndNeverBlocks)$' ./internal/sharedcache ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMeasureContextStopsWhenTheContextEnds TestApplyContextStopsBetweenEntries TestMeasurePlansIdleEntriesOldestFirstUntilWithinTheLimit TestApplySkipsAnEntryUsedSinceItWasPlanned TestImplementStartTrimsTheDeclaredSharedBuildCache TestImplementStartRunsOnWhenTheTrimDoesNot/refused_cache TestImplementStartRunsOnWhenTheTrimDoesNot/budget_ends TestImplementStartWithoutASharedBuildCacheRecordsNoTrim TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification/idle_entries_removed_first TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification/budget_ends TestGCDryRunReportsAnOversizedSharedBuildCacheAndChangesNothing TestGCRemovesOnlyIdleEntriesOfAnOversizedSharedBuildCache TestImplementStartRunsRunRetention; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/sharedcache ./internal/cli && GOOS=linux go build -buildvcs=false ./internal/sharedcache ./internal/cli` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check; without the Run-start call the trim test finds the 10 h entry, with the trim error returned the refused-cache subtest exits 2, without the budget the budget subtest records `trimmed`, and with the gate trim after the Verification the gate subtest's Verification fails. Afterwards the engine, Run-start, gate, GC and Run Retention tests pass and both packages build for Windows and Linux.

## References

- `_prd.md` → Goals; Core Feature 1; Core Feature 2; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior 1-4; Interfaces; Data Models; API Contracts 1-5; Invariants 1-4; Testing Approach; Build Order 2


## Result

Implemented this Task's slice; Daemon Verification and settlement remain pending.
The pre-existing `status: in_progress` is unchanged.

- Added `MeasureContext` and `ApplyContext`, with cancellation checks between
  entries, wrapped context errors and partial removal counts. The original
  `Measure` and `Apply` delegate with `context.Background()`; existing engine
  tests are unchanged.
- Added declaration-ordered trim reports with all five outcomes and causes,
  and rebuilt GC's sweep on the same helper. Automatic triggers share one
  thirty-second context budget across all declarations. A stopped measure
  reports size zero, and stopped or failed applies retain actual removal counts.
- Implement Run start trims after detached-start reporting and before Run
  Worktree creation, journals one `daemon.status` payload per cache and prints
  the specified stderr lines. Event-append errors warn without changing the
  Run outcome. No declarations produce no trim event or output.
- The Delivery Queue repository gate trims first and writes the specified
  prefixed lines to the Queue Owner's log. Its existing Verification and
  Linux Lane control the result.

Acceptance evidence:

1. `TestImplementStartTrimsTheDeclaredSharedBuildCache` observes removal of
   exactly the 10 h and 3 h entries at the Run Worktree creation boundary,
   preserves the fresh entry and README, checks the exact stderr line and all
   payload fields, and requires the trim event before the first Task event.
2. `TestImplementStartRunsOnWhenTheTrimDoesNot/refused_cache` and
   `/budget_ends` check exit zero, Clean outcome, identical Task output,
   retained entries, and the `kept`/README and `stopped`/`budget of 0s elapsed`
   event causes. `TestImplementStartWithoutASharedBuildCacheRecordsNoTrim`
   checks the same Task output and absence of trim events and lines.
   `TestImplementStartRunsOnWhenTheTrimEventCannotBeJournaled` injects a real
   SQLite append refusal and observes one warning with the Run still Clean.
3. `TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification` checks
   Verification that requires the oldest entry absent, exact log output, and
   successful Verification after budget exhaustion. Existing GC report tests,
   Run Retention tests, end-to-end implement test and Linux Lane gate matrix
   pass unchanged. Additional helper tests cover unchanged/failed outcomes and
   one stopped report for every remaining declaration after cancellation.

Focused checks (all commands below ran in this worktree):

- Initial red: `go test ./internal/sharedcache -run 'ContextStops' -count=1`
  failed to compile because `MeasureContext` and `ApplyContext` did not exist.
- `go test ./internal/sharedcache -count=1 -timeout=2m` — exit 0; all six
  engine tests passed, including cancellation and exactly one removal before
  cancellation between planned entries.
- `go test ./internal/cli -run '^(TestImplementStart.*|TestSharedCacheTrim.*|TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification|TestGC.*SharedBuildCache.*|TestRunRetentionAtStartWarnsAndNeverBlocks|TestGateReturnsTheLaneEnvironmentSkip|TestRunImplementExecutesSpecEndToEnd)$' -count=1 -timeout=5m`
  — exit 0.
- `git diff --check` — exit 0.

The cancellation helper test initially exposed a later unresolved path being
classified as failed after cancellation. Checking cancellation before cache
measurement/classification corrected the precedence; the focused checks above
include the regression.

Not run by this Agent: the declared Verification command (including its Windows and Linux
cross-builds), `make verify`, and `make verify-changed`; these remain
Daemon-owned. No other Task, manifest, configuration, Makefile or document was
edited. No commit, push or Pull Request was made. No follow-up identified.


### Verification Feedback repair — attempt 1

Inspected the Daemon diagnostic artifact
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261026T215311Z_146a5646537e95bd/verification/batch-002-attempt-1.log`.
The configured Verification stopped at
`TestEveryRealTimeBudgetIsWithinItsFileCeiling`: the two positive nanosecond
budget assignments in the new CLI tests exceeded that file's zero ceiling.

Replaced the two `1ns` test inputs with zero budgets, which expire
synchronously without depending on a timer. This is a bounded deviation from
Testing Approach's `1ns` fixture input to satisfy the repository's existing
Real-Time Budget rule. Both subtests retain exact stopped-outcome, cause and
console assertions (now `0s`), retained-file checks and successful Run/gate
outcomes. The production thirty-second default and the budget-contract test
are unchanged.

Fresh focused evidence after the repair:

- `go test ./internal/cli -run '^(TestImplementStartRunsOnWhenTheTrimDoesNot|TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification)$' -count=1 -timeout=5m`
  — exit 0.
- `go test ./internal/testfixture -run '^TestEveryRealTimeBudgetIsWithinItsFileCeiling$' -count=1 -timeout=2m`
  — exit 0.
- `git diff --check` — exit 0.

The Agent did not rerun declared Verification or `make verify-changed`.
Task status remains Daemon-owned and unchanged; no commit, push or Pull Request.
