---
task: task_02
spec: 0250-a-run-database-that-keeps-only-recent-runs
status: pending
type: backend
complexity: high
---

# Task 02: roundfix gc applies Run Retention and reports the Runs, rows and bytes it removes

## Overview

Adds the User Config window `store.run_retention_days` and the Run Retention
Sweep engine in `internal/cli`, and wires it into the GC Command. The dry run
reports what the window would remove and keep; the live run removes the Runs,
their artifact directories under proven roots, compacts the database and
records the sweep. The GC Command and configuration guides describe both. It
is verifiable on its own through `roundfix gc` on fixture homes.

## Requirements

1. MUST add `store.run_retention_days` to `internal/config` as
   `_techspec.md` → API Contract 1 states: `Store.RunRetentionDays`, Builtin
   30, the values 7, 15 and 30 only, the exact refusal
   `store.run_retention_days must be 7, 15 or 30` inside the existing
   `parse config "<path>": ` prefix, the Project Config removal with the
   exact ignored-setting warning, and the key with a comment under `store:` in
   the User Config template only. `store.journal_retention` keeps its parsing,
   validation, default and template line.
2. MUST add `sweepRunDatabase` in `internal/cli/run_retention.go` with the
   shape of `_techspec.md` → Interfaces, implementing API Contract 2 and API
   Contract 3: the cutoff is the GC Command's `now` dependency minus the
   window; it keeps a candidate that is queue-referenced, whose recorded Run
   Worktree exists, or whose `runs/<run-id>` directory exists under an
   unproven root, counting each reason; it removes a removable Run's
   directory under its proven root first and then calls `RemoveRetainedRun`;
   it treats `RunRetentionKeptError` as kept; it compacts and records
   completion as API Contract 3 states. With a budget it checks the clock
   before each Run removal and each compaction slice and records nothing when
   it stops.
3. MUST make `roundfix gc --dry-run` print the Run Retention section of
   `_techspec.md` → Surface Transcript 1 after the existing Journal Retention
   lines and lists, opening only the reader and writing nothing, and make
   `roundfix gc` print the section of Surface Transcript 2, as API Contract 5
   states, including the five compaction wordings and page-count bytes.
4. MUST, with `journal_retention: 0`, print the normal `GC dry-run` or
   `GC complete` header, `Journal Retention: 0` and
   `No journal pruning performed.`, and then run Run Retention. MUST update
   `TestRunGCSkipsWhenJournalRetentionIsZero` in `internal/cli/gc_test.go` to
   that contract, keeping its assertions that the journal and directories of
   the 400-hour-old Run stay. No other existing gc assertion changes.
5. MUST make an unsupported window fail `roundfix gc` before any work with the
   output of `_techspec.md` → Surface Transcript 3 and exit 2.
6. MUST update the `gc` usage text in `internal/cli/cli.go` so it no longer
   says that gc never deletes Run rows and says that it removes terminal Runs
   past Run Retention, and keep the words `Journal Retention` it already has.
7. MUST describe the key in `docs/user-guide/configuration.md` (its example
   block, its defaults table and its retention paragraph) and the Run
   Retention section in `docs/user-guide/commands/gc.md`, naming the kept
   reasons and that the dry run's bytes are an estimate.
8. MUST add these tests:
   - in `internal/config/run_retention_test.go`:
     `TestRunRetentionDaysDefaultsToThirty`,
     `TestRunRetentionDaysAcceptsOnlySevenFifteenOrThirty` (7, 15 and 30
     load; 0, 10, 31, -7, `30d` and `"30"` fail with the exact message),
     `TestRunRetentionDaysIsUserConfigOnly` (a Project Config value is
     ignored with the exact warning and the User Config value applies), and
     `TestInitUserConfigCarriesRunRetentionDays` (the User template carries
     the key with 30, the Project template does not);
   - in `internal/cli/gc_run_retention_test.go`, through `runCLIContext`
     with `withGCNow`:
     `TestGCDryRunReportsRunRetention` asserts Surface Transcript 1 against a
     fixture with one removable Run, one queue-referenced Run and one Run
     with an existing Run Worktree past the cutoff, and that the database and
     artifact directories are byte-identical afterwards;
     `TestGCRemovesRunsPastRunRetention` asserts Surface Transcript 2 and that
     the removable Run left every table and its directory while the kept and
     recent Runs kept every row;
     `TestGCKeepsProtectedRuns` covers an Active Run created before the
     cutoff and the three kept reasons;
     `TestGCKeepsARunUnderAnUnprovenArtifactRoot` records an
     `artifact_dir` outside Roundfix Home whose directory exists and proves
     the Run and the directory stay and are counted as unproven;
     `TestGCCompactsAfterRunRetention` proves `incremental (<n> pages)` on a
     fresh database and `full (converted to incremental)` on a mode-0
     database, each with fewer bytes after than before, and
     `skipped (<reason>)` with exit 0 when an Active Run blocks the full
     compaction;
     `TestGCRefusesAnUnsupportedRunRetention` asserts Surface Transcript 3.

## Subtasks

- [ ] Add the configuration key, its refusal, its warning and its template line.
- [ ] Add the sweep engine with the kept reasons, the artifact proof and the compaction.
- [ ] Add the Run Retention section to the dry run and the live run.
- [ ] Update the zero-window test, the usage text and the two guides.
- [ ] Write the ten tests and run the config and gc tests.

## Acceptance Criteria

- [ ] `store.run_retention_days` loads 7, 15 or 30 from User Config, defaults to 30, and refuses everything else.
- [ ] `roundfix gc --dry-run` reports the Runs, rows and estimated bytes and changes nothing.
- [ ] `roundfix gc` removes exactly the removable Runs, keeps every protected Run, and compacts.

## Context

- interface: `internal/config/config.go`
- interface: `internal/cli/gc.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/gc_test.go`
- interface: `docs/user-guide/commands/gc.md`
- interface: `docs/user-guide/configuration.md`
- creates: `internal/config/run_retention_test.go`
- creates: `internal/cli/run_retention.go`
- creates: `internal/cli/gc_run_retention_test.go`
- instruction: `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- instruction: `docs/adr/0171-a-retention-prune-reports-only-what-it-reclaimed.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestRunRetentionDaysDefaultsToThirty|TestRunRetentionDaysAcceptsOnlySevenFifteenOrThirty|TestRunRetentionDaysIsUserConfigOnly|TestInitUserConfigCarriesRunRetentionDays)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestRunRetentionDaysDefaultsToThirty TestRunRetentionDaysAcceptsOnlySevenFifteenOrThirty TestRunRetentionDaysIsUserConfigOnly TestInitUserConfigCarriesRunRetentionDays; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestGCDryRunReportsRunRetention|TestGCRemovesRunsPastRunRetention|TestGCKeepsProtectedRuns|TestGCKeepsARunUnderAnUnprovenArtifactRoot|TestGCCompactsAfterRunRetention|TestGCRefusesAnUnsupportedRunRetention|TestRunGCSkipsWhenJournalRetentionIsZero)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestGCDryRunReportsRunRetention TestGCRemovesRunsPastRunRetention TestGCKeepsProtectedRuns TestGCKeepsARunUnderAnUnprovenArtifactRoot TestGCCompactsAfterRunRetention TestGCRefusesAnUnsupportedRunRetention TestRunGCSkipsWhenJournalRetentionIsZero; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; for f in docs/user-guide/commands/gc.md docs/user-guide/configuration.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'store.run_retention_days' || { printf 'missing store.run_retention_days in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/gc.md | grep -qF -- 'Database bytes reclaimable (estimated)' || { printf 'missing the dry-run estimate in the gc guide\n' >&2; exit 1; }; go test -count=1 ./internal/config && go test -count=1 -run 'GC|Gc|Storage|Retention' ./internal/cli` — expected: exit 0. Before this Task the ten new tests do not exist and neither guide names the key, so the first check fails. After it, the new tests and the updated zero-window test pass, both guides describe the key and the report, and every `./internal/config` test and every gc and storage test of `./internal/cli` passes.

## References

- `_prd.md` → Core Feature 1; Core Feature 2; Core Feature 4; Core Feature 5; Core Feature 6; User Story 2; User Story 4; Success Metric 1; Success Metric 2; Success Metric 4
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 5; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Invariants 1-7; Build Order 2
- ADR-0255; ADR-0171
