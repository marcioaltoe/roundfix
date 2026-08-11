# QA evidence — 0081 rerun 01

Build: `5af85c4785e9f1c322e554fdd420600ce163d8e9`.

## Static gate

- `rtk roundfix spec check 0081-a-journal-cheap-to-write-and-keep --strict`
  exited 0 with `No findings.` The checker skipped
  `SC-VOCABULARY-UNDOCUMENTED` because `_techspec.md` has no Vocabulary
  Contract.
- `rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache go clean
  -testcache` exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache make verify`
  passed formatting and started `go test -parallel 16 ./...`; the sandbox
  terminated the process when it attempted network access to
  `api.github.com`. The exact denial was `Network access to "api.github.com"
  was blocked: domain is not on the allowlist for the current sandbox mode.`
  This is an environment block, not a red test result; it was not retried.
- `rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache make verify-docs`
  exited 0. It built `bin/roundfix`, passed `docscontract`, all three
  `repocontract` tests, the Spec corpus budget, and the full Spec consistency
  sweep.
- `rtk make fmt-check` exited 0. `rtk env
  GOCACHE=/private/tmp/roundfix-0081-qa01-gocache make skills-sync-check
  skills-check build` exited 0. These cover every non-test prerequisite of
  `make verify` after the network-denied test target.
- `rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache go test -count=1
  ./internal/store ./internal/tui ./internal/runevent` exited 0. These are the
  complete feature-owning package suites; the focused CLI replay below covers
  the named CLI consumers.

## R03 — same-boundary write cost

The pre-change build was extracted from
`a2a4c86b7570e4ef782ccc2ff390033f586d47fd`. Both builds ran the same shared
harness through Go overlays; only the API adapter differed.

Pre-change command:

```text
rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-old-gocache go test -overlay=/private/tmp/roundfix-0081-old-overlay.json -count=3 ./internal/store -run '^TestTask10SameBoundaryMeasurement$' -v
```

Observed pre-change wall cost per 1,000 events: 130,068, 128,901, and 128,737
µs; median 128,901 µs. Every sample persisted 1,536 events from six Runs with
256 events per Run at 5,000 ms.

Current-build command:

```text
rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache go test -overlay=/private/tmp/roundfix-0081-current-overlay.json -count=3 ./internal/store -run '^TestTask10SameBoundaryMeasurement$' -v
```

Observed current wall cost per 1,000 events: 37,571, 37,045, and 37,707 µs;
median 37,571 µs. Every sample persisted the same 1,536 events at the same
publish-through-flush boundary. The median fell 70.9%, so the improvement is
not attributable to fewer events.

The adjacent current-build boundary canary also exited 0:
`TestBatchClosesOnCountLingerAndImmediate`, all four
`TestBatchAmbiguousCommit` outcomes, and `TestBatchPreservesPayloadBytes`
passed. This confirms the current path closes batches at its declared
boundaries and preserves raw bytes.

## R06 — pre-change consumer replay

`rtk shasum -a 256` matched every recorded digest:

```text
1c38ed391769fb2c3ba5e48875ed868737bf0bc5659c12bdbf7ee8b19f5c5dec  2026-08-11-prechange-roundfix.db
3c4eb9511ee29dda3116ac8863e8b7f0b5a74fcaf594470ce6f0788b7f5ab04b  2026-08-11-prechange-consumer-expectations.json
16eb1d01eb12a55728e0343c12768f04e40809359acc8214c368a42e5a8d9700  2026-08-11-prechange-cockpit.golden
```

Command:

```text
rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache go test -count=1 ./internal/store -run '^TestJournalConsumerCorpusReplaysEveryConsumer$' -v
```

The command exited 0. Its CLI subtest passed the current production `events`,
Attach, reconcile replay detection, and fixed-cutoff `gc` observations against
the recorded pre-change expectations. Its TUI subtest passed a byte-for-byte
comparison of the full color-free 120x40 Cockpit frame. The source database
was reopened and copied for mutating GC observation, so the recorded fixture
remained byte-identical.

The built current CLI's `events --help`, `attach --help`, `gc --help`, and
`reconcile --help` each exited 0. They retain the glossary's `Supervisor event
stream`, `Run Database`, read-only `Attach`, `Run Event Journal`, `Journal
Retention`, `Run rows`, and `Active Run locks` vocabulary and unchanged command
shapes.

## R01 — Run-start cost

Command:

```text
rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache go test -count=2 ./internal/store -run '^TestJournalMeasurementHarness$' -v
```

The command exited 0 twice. Run A p50 was 104 µs with 0 events, 59 µs with
1,000, and 36 µs with 10,000. Run B was 49, 39, and 39 µs respectively. The
10,000-event database occupied 21,135,360 bytes versus 61,440 bytes empty, yet
neither fresh observation followed the former size curve. Both observations
used newly migrated temporary Run Databases and the production retention path.

## R02 — parallel Runs

Command:

```text
rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache go test -count=1 ./internal/store -run '^TestParallelRuns$' -v
```

The command exited 0. Six production writer Stores at the exact 5,000 ms
pre-raise timeout persisted 1,536 of 1,536 events with `sqlite_busy=0` and all
six Runs complete. Wall time was 55,951 µs (36,426 µs per 1,000 events), with
writer p50/p95 24,900/55,933 µs. Production-reader read-back confirmed cursor
contiguity and publisher order for every Run. The cancellation rehearsal
confirmed a cancelled lock holder releases the advisory lock and the remaining
writers proceed.

## R04 — bytes and retention shape

The external acceptance source is
`docs/findings/2026-08-06-three-gigabytes-of-event-journal-inside-the-retention-window.md`.
It predates the Spec and records 2,915 MB of `run_events`, 1,645,457 events,
and a 42,000-event largest Run.

The fresh R01 harness reported 61,440 / 2,162,688 / 21,135,360 database bytes
at 0 / 1,000 / 10,000 retained events, identical to the recorded before. Exact
`a2a4c86b..HEAD` diffs of ADR-0008, ADR-0033, and
`docs/user-guide/run-database-lifecycle.md` exited 0 with no output. The same
was true for `.roundfixrc.yml`, `go.mod`, and `go.sum`.

Command:

```text
rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache go test -count=1 ./internal/store ./internal/cli -run 'TestBatchPreservesPayloadBytes|TestPruneTerminalRunsDeletesOnlyEligibleJournalRows|TestPruneTerminalRunsNoOpsWhenCutoffSelectsNothing|TestRetentionPreservesRunLifecycleRecords|TestRunGCDryRunListsEligibleRunsAndChangesNothing|TestRunGCPrunesEligibleJournalsArtifactsAndOrphans|TestRunGCSkipsWhenJournalRetentionIsZero' -v
```

All seven selected tests passed. Raw producer JSON round-tripped byte-for-byte;
only eligible terminal journal rows were deleted; Run lifecycle rows survived;
dry-run did not mutate; actual GC pruned its disposable eligible fixture; and
zero retention kept everything.

## Scope and Non-Goals

The exact `a2a4c86b..HEAD` changed-path inventory contains Spec artifacts,
journal/store implementation and tests, CLI journal wiring, and Cockpit paths.
It contains no `.roundfixrc.yml`, dependency file, schema migration, second
store, or repository-tooling path. `internal/store/store.go` still declares
schema version 12. Exact diffs of `.roundfixrc.yml`, `go.mod`, and `go.sum`
were empty. `git diff --check` exited 0.

## R05 — Cockpit refresh and TUI sweep

The 15 selected production-model checks passed for forward-cursor cost,
header projection, missing-payload summary fallback, all ten unchanged
snapshots, 80x24/120x40/200x50 layouts, degenerate dimensions, keyboard detail
toggle, Follow Mode freeze/resume, completed review/spec Attach replay,
incremental/full replay identity, stale-event resistance, and no-color marker
distinctions.

The PRD's 10,000-event boundary was exercised through the reproducible
QA-only overlay harness at
`r05-10000-event-harness_test.go.txt`:

```text
rtk env GOCACHE=/private/tmp/roundfix-0081-qa01-gocache go test -overlay=/private/tmp/roundfix-0081-r05-overlay.json -count=2 ./internal/tui -run '^TestQACockpitRefreshCostAtTenThousandEvents$' -v
```

Both executions passed. After both 20- and 10,000-event backlogs, the same two
new events cost exactly two header rows, one full row, and 67 payload bytes.
The opening fold read each backlog once, so the comparison measures the next
refresh instead of skipping initial replay cost.
