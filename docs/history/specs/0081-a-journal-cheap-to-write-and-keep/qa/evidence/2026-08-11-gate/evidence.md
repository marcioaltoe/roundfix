# QA evidence — Spec 0081 — 2026-08-11

Build: `7dea0bd3f25594bfbc0673c0558c2e67db8c6179`

This log records fresh command evidence for rows R01–R06. It is updated as
the seeded resumable matrix closes.

## Authoring precondition

Command:

```text
rtk roundfix spec check 0081-a-journal-cheap-to-write-and-keep --strict
```

Exit: 0.

```text
Spec 0081-a-journal-cheap-to-write-and-keep
No findings.
Skipped:
  SC-VOCABULARY-UNDOCUMENTED: missing docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md Vocabulary Contract
```

## Static gate

### Cold repository Verification

Commands:

```text
rtk env GOCACHE=/private/tmp/roundfix-0081-qa-gocache go clean -testcache
rtk env GOCACHE=/private/tmp/roundfix-0081-qa-gocache make verify
```

The cache clear exited 0. Verification exited 2 at `fmt-check` before tests:

```text
Go files need formatting:
./internal/store/journal.go
./internal/store/journal_consumer_corpus_test.go
./internal/store/journal_header_test.go
make: *** [fmt-check] Error 1
```

### Repository documentation contracts

Command: `rtk make verify-docs`.

Exit: 0. It rebuilt `bin/roundfix`, passed `internal/docscontract`, the
measured sanctioned-ownership and frozen-boundary checks, the corpus budget,
and Spec consistency. Spec 0081 had no findings and retained the expected
`SC-VOCABULARY-UNDOCUMENTED` skip.

## R01 — Run-start cost

Command:

```text
rtk proxy env GOCACHE=/private/tmp/roundfix-0081-qa-gocache go test -count=2 ./internal/store -run '^TestJournalMeasurementHarness$' -v
```

Exit: 0; both executions passed. Fresh Run-start retention observations:

| Journal events | Database bytes | Run A p50/p95 us | Run B p50/p95 us |
| ---: | ---: | ---: | ---: |
| 0 | 61,440 | 51 / 99 | 59 / 79 |
| 1,000 | 2,162,688 | 41 / 68 | 40 / 69 |
| 10,000 | 21,135,360 | 37 / 66 | 38 / 67 |

The test used fresh migrated temporary Roundfix Homes, fixed 1,800-byte raw
payloads, and the production retention sequence. The repeat run independently
confirmed that Run-start cost no longer follows retained journal size.

## R02 — Parallel Runs

Command:

```text
rtk proxy env GOCACHE=/private/tmp/roundfix-0081-qa-gocache go test -count=1 ./internal/store -run '^TestParallelRuns$' -v
```

Exit: 0; all three rehearsal subtests passed.

```text
PARALLEL_RUNS runs=6 events_per_run=256 total_events=1536 busy_timeout_ms=5000 sqlite_busy=0 completed_runs=6 elapsed_us=58577 writer_latency_p50_us=25176 writer_latency_p95_us=58539 concurrent_wall_per_1000_events_us=38136
```

The second subtest read all six journals back and asserted contiguous cursors
plus publisher order. The cancellation subtest held the advisory lock,
cancelled its owner, and proved the other Runs proceeded without a busy error.

## R03 — Event-write cost

Command:

```text
rtk proxy env GOCACHE=/private/tmp/roundfix-0081-qa-gocache go test -count=1 ./internal/store -run 'Batch|AmbiguousCommit' -v
```

Exit: 0. The current build passed count, linger, immediate, explicit-flush,
ordering, raw-payload, close, failure-preservation, and all ambiguous-commit
rehearsals. The R02 run persisted all 1,536 events and measured 38,136 us per
1,000. Current source pins `journalBatchSize=128` and the scenario publishes
256 events per Run, so the six Runs close twelve count batches without
reducing event count.

The required improvement is not established because the pre-change report
contains only direct `AppendRunEvent` latency percentiles. It records no
pre-change wall cost per 1,000 events. The repaired same-harness report states:

```text
The production JournalWriter also amortizes agent output across batches; this
direct-AppendRunEvent harness intentionally preserves task_01's measurement
boundary and does not claim a batched-path latency speedup.
```

The available before and after values therefore use different measurement
boundaries. R03 fails rather than inferring an improvement from incomparable
numbers.

## R04 — Bytes and retention shape

Preservation command:

```text
rtk git -c core.fsmonitor=false diff --exit-code 685d201b658cc46e944634a3c072da2a7d1d83c3..HEAD -- docs/adr/0008-run-event-payload-stores-raw-producer-json.md docs/adr/0033-the-run-event-journal-is-pruned-by-retention.md docs/user-guide/run-database-lifecycle.md
```

Exit: 0 with no diff. The binding raw-payload decision, terminal-only age
retention decision, and operator lifecycle guide are byte-preserved from the
pre-Spec production commit.

Focused command:

```text
rtk proxy env GOCACHE=/private/tmp/roundfix-0081-qa-gocache go test -count=1 ./internal/store ./internal/cli -run '^(TestBatchPreservesPayloadBytes|TestRunEventPayloadRoundTripsByteExact|TestPruneTerminalRunsDeletesOnlyEligibleJournalRows|TestPruneTerminalRunsNoOpsWhenCutoffSelectsNothing|TestRetentionPreservesRunLifecycleRecords|TestRunGCDryRunListsEligibleRunsAndChangesNothing|TestRunGCPrunesEligibleJournalsArtifactsAndOrphans|TestRunGCSkipsWhenJournalRetentionIsZero)$' -v
```

Exit: 0; all eight selected tests passed.

Before and fresh after database bytes are identical:

| Journal events | Before bytes | Fresh after bytes |
| ---: | ---: | ---: |
| 0 | 61,440 | 61,440 |
| 1,000 | 2,162,688 | 2,162,688 |
| 10,000 | 21,135,360 | 21,135,360 |

The unchanged byte curve is intentional: retained Runs keep every raw payload
until the existing terminal-only age boundary.

Outside acceptance evidence:
`docs/findings/2026-08-06-three-gigabytes-of-event-journal-inside-the-retention-window.md`
was authored before and outside this Spec. It records 2,915 MB of
`run_events`, 1,645,457 events, and a largest Run of 42,000, providing the
independent observed workload that the hermetic fixture models.

## R05 — Cockpit incremental refresh and TUI sweep

Command:

```text
rtk proxy env GOCACHE=/private/tmp/roundfix-0081-qa-gocache go test -count=1 ./internal/tui -run '^(TestCockpitRefreshCostTracksNewEvents|TestCockpitTaskJournalRefreshUsesForwardCursorAndHeaderProjection|TestCockpitTaskJournalForwardCursorKeepsSummaryFallback|TestCockpitRenderSnapshots|TestCockpitResponsiveFallbackAndStableSizes|TestCockpitDegenerateSizesRenderEmptyWithoutPanic|TestCockpitTabSwitchesFocusAndArrowsMoveSelection|TestCockpitAttachReplaysFinishedSpecRunThroughRedesignedCockpit|TestCockpitScrollFreezesFollowAndStatusBarNarratesStates|TestCockpitHeaderNoColorRendersMarkerOnlyText|TestCockpitWorkItemCardNoColorKeepsMarkerDistinctions|TestViewportScrolledFreezesAndCountsNewEventsBelow)$' -v
```

Exit: 0; all twelve selected tests and every subtest passed.

`TestCockpitRefreshCostTracksNewEvents` opened a 20-event and a 400-event
backlog, appended the same agent line plus one daemon Task event, and required
equal costs: exactly two header rows and one full row in both cases. It also
required identical rendered Task states. The rest of the sweep covered
completed-Run attach replay; ten unchanged 88x24/120x40 snapshots; 80x24,
120x40, and 200x50 responsive layouts; zero/negative dimensions; keyboard
focus; follow freeze/resume; new-event counts; missing-payload summary
fallback; and no-color markers.

## R06

### Current behavior observations

The four current Store corpus tests passed:

```text
rtk proxy env GOCACHE=/private/tmp/roundfix-0081-qa-gocache go test -count=1 ./internal/store -run '^(TestConsumerCorpusFullReadReplaysIdentically|TestConsumerCorpusEventsStreamReplaysIdentically|TestReplayCorpusHeaderMatchesFullRead|TestReplayCorpusBatchClockMatchesFullEvents)$' -v
```

Four public `events` runner tests passed, covering default/filter JSONL,
legacy verification, immediate terminal replay, and malformed payload failure:

```text
rtk proxy env GOCACHE=/private/tmp/roundfix-0081-qa-gocache go test -count=1 ./internal/cli -run '^(TestEventsReplayDefaultAndFilterJSONLRecordsOnly|TestEventsReplayLegacyVerificationEvent|TestEventsTerminalRunReplaysAndExitsImmediately|TestEventsMalformedRelevantPayloadFailsNoStdout)$' -v
```

`TestAttachSpecRunReplaysTaskAndVerificationCapacity` and
`TestAttachFollowerAppendsOnlyNewerEventsWithoutDuplicates` each passed in
fresh runs. Two reconcile payload-coverage paths passed:
`TestCarryForwardSettlesATaskWhoseInputsAreUnchanged` and
`TestCarryForwardWithoutTheFlagReportsAndChangesNothing`. R04 records passing
GC dry-run/apply semantics; R05 records passing Cockpit replay and snapshots.

Built public help exited 0 for `bin/roundfix events --help`, `attach --help`,
`gc --help`, and `reconcile --help`. It retains the canonical terms and
unchanged command shapes. Optional review-Run Attach fixtures tried to reach
`api.github.com` during their resolve preflight and were denied by the managed
sandbox; no customer credential or Pull Request path was bypassed.

### Missing regression proof

Commands:

```text
rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-status -r a8af25b7
rtk git -c core.fsmonitor=false log --diff-filter=A --format='%H %s' -- internal/store/journal_consumer_corpus_test.go
```

The output proves task_05 commit
`a8af25b7fd8093637480d46e5ee5da7aaa6f80bc` simultaneously modified
`internal/store/journal.go` and `internal/tui/cockpit.go` and added
`internal/store/journal_consumer_corpus_test.go`. The corpus therefore was not
recorded before the production change.

Inspection of that file finds exactly four consumer tests: full-read paging,
`events` projection, header/full subset equality, and batch-clock equivalence.
It contains no attach, full Cockpit render, reconcile payload parser, or GC
replay. Separate current-fixture passes prove current behavior only; they do
not prove byte-identical pre/post behavior for one pre-change journal. R06
fails on that missing provenance and surface coverage.

## Scope and Non-Goal audit

`rtk git -c core.fsmonitor=false diff --name-status
685d201b658cc46e944634a3c072da2a7d1d83c3..HEAD` names only the two baseline
reports, task_01 through task_08, journal/store implementation and tests, CLI
journal wiring, and Cockpit implementation/tests. An exact diff of
`.roundfixrc.yml`, `go.mod`, and `go.sum` exits 0 with no output.

Current and pre-Spec `internal/store/store.go` both declare
`schemaVersion = 12`. No migration file exists in the changed inventory.
`rtk git -c core.fsmonitor=false diff --check` exits 0; this checks whitespace
errors but does not override F-01's `gofmt` failure.

The audit confirms no second store, dependency change, Project Config edit,
schema migration, stream-schema change, compaction change, second retention
window, or early payload shedding. Fixed event counts in R02/R03 confirm the
measurements do not reduce event production.
