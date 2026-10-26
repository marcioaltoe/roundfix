---
task: task_01
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
status: pending
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
