---
task: task_02
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: completed
type: data
complexity: high
---

# Task 02: Every prompt of a Run leaves a usage record and a `usage` event

## Overview

task_01 puts a Turn Usage on each prompt's result, and nothing keeps it. This
Task adds the Run Database schema of `_techspec.md` → Data Models, writes one
row and one `daemon.token_usage` Run Event per prompt from the Daemon, sums
the rows per scope, per Run and per Delivery Queue, and adds the `usage`
category to the Run Event Stream. It is verifiable on its own: a Run whose
fake runner reports usage has one row per prompt, `roundfix events` prints the
`usage` records, and a failed write changes no outcome.

## Requirements

1. MUST raise the Run Database schema by one from the version on this Task's starting commit and add, in one migration, the table `run_token_usage`, the table `delivery_queue_runs` and the column `delivery_queues.max_tokens` of `_techspec.md` → Data Models, with their constraints. A database at the previous version MUST migrate with its rows intact.
2. MUST add `AppendTokenUsage`, `RunTokenUsage` and `DeliveryQueueTokenUsage` in the new file `internal/store/token_usage.go`, with the shapes of `_techspec.md` → Interfaces. `AppendTokenUsage` MUST insert the row and append the `daemon.token_usage` Run Event with the payload of `_techspec.md` → Recording in one write, as `AppendAgentSelectionAttempt` records its event.
3. MUST sum, in `RunTokenUsage` and `DeliveryQueueTokenUsage`, only reported prompts into tokens, count unreported prompts separately, keep scopes in first-recorded order, and compute each Agent Session's spend with the rule of `_techspec.md` → Data Models. A total with no reported prompt MUST carry no token number.
4. MUST insert the item's Run ID into `delivery_queue_runs` with `INSERT OR IGNORE` in the same write as `UpdateDeliveryQueueItem` and `RetryDeliveryQueueItem` whenever that Run ID is not empty, and MUST add `MaxTokens` to `DeliveryQueueLimits`, persisted by `CreateDeliveryQueueWithLimits` and read by `DeliveryQueue`. Nothing in this Task sets `MaxTokens` from a command.
5. MUST add the two tables to the durable-table lifecycle policy in `docs/user-guide/run-database-lifecycle.md`, owned by the Run lifecycle and the Delivery Queue lifecycle, never deleted by Journal Retention, and MUST leave Journal Retention and the GC Command unchanged.
6. MUST record each prompt from the Daemon as `_techspec.md` → Recording states: after every prompt of an Agent Session owner, with its scope, attempt and active candidate, and after a prompt without an owner, with scope kind `session`. A failed or stopped prompt MUST still be recorded. A failed write MUST print the warning line to the Run's progress output and MUST NOT change the prompt's result, the Task's status or the Run's outcome.
7. MUST add `KindDaemonTokenUsage` (`daemon.token_usage`) and the stream category `usage` with the fields and summary of `_techspec.md` → Surfaces, on by default and accepted by `--filter`. `IsDaemonKind` MUST stay unchanged. The stream schema MUST stay `roundfix-events/v1`.
8. MUST append `,usage` to the category list in the `events` help text in `internal/cli/cli.go`, keeping the string `task-status,batch,verification,outcome,agent-selection` whole so `internal/cli/cli_test.go` does not change, and MUST document the category and its fields in `docs/user-guide/commands/events.md` with Surface Transcript 7.
9. MUST put the new tests in `internal/store/token_usage_test.go`, `internal/daemon/token_usage_test.go`, `internal/runevent/stream_usage_test.go` and `internal/cli/events_usage_test.go`. Store tests MUST use a temporary Roundfix Home; Daemon tests MUST use a fake runner and the real store. No test opens the live Run Database under `~/.roundfix`.
10. MUST prove each new gate can fail. The Result MUST record one sabotage of the sums (for example adding an unreported prompt as zero), one of the queue links (for example linking only in `UpdateDeliveryQueueItem`) and one of the outcome guard (for example returning the write error), each with the test that failed, and that the code was restored.

Spec 0194 creates the per-command guide and skill files this Task edits
(`_prd.md` → Prerequisites). They are declared under `creates:` because they do
not exist when this Spec is authored. This Task MUST NOT create one that Spec
0194 has not created: when such a file is absent, the Task stops and reports
that Spec 0194 has not landed.

## Subtasks

- [ ] Add the migration, the store API and the queue links.
- [ ] Record every prompt from the Daemon, with the write guard.
- [ ] Add the Run Event kind and the `usage` stream category.
- [ ] Document the tables and the category.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] A database at the previous schema migrates, and the new tables and column exist with their constraints.
- [ ] Each prompt of a Run leaves one row and one `daemon.token_usage` Run Event, including a failed and a stopped prompt.
- [ ] Sums per scope, per Run and per queue count unreported prompts apart and never as zero.
- [ ] Every Run ID recorded on a queue item, including after a retry, is linked to the queue and counted by its sum.
- [ ] A failing write leaves the Task and Run outcome unchanged and prints the warning.
- [ ] `roundfix events <run-id> --filter usage` prints Surface Transcript 7.

## Context

- creates: `internal/store/token_usage.go`
- creates: `internal/store/token_usage_test.go`
- creates: `internal/daemon/token_usage_test.go`
- creates: `internal/runevent/stream_usage_test.go`
- creates: `internal/cli/events_usage_test.go`
- interface: `internal/store/store.go`
- interface: `internal/store/delivery.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/stream.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `internal/cli/cli.go`
- interface: `docs/user-guide/run-database-lifecycle.md`
- creates: `docs/user-guide/commands/events.md`
- instruction: `internal/store/agent_selection.go`
- instruction: `internal/store/journal_test.go`
- instruction: `internal/cli/cli_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestMigrationAddsTokenUsageTablesToThePreviousSchema|TestAppendTokenUsageWritesTheRowAndItsRunEvent|TestRunTokenUsageSumsScopesAndTheRun|TestRunTokenUsageKeepsUnreportedPromptsOutOfTheTotal|TestSessionSpendSumsIncreasesAndRestartsOnADrop|TestDeliveryQueueRunsLinkEveryRecordedRunID|TestDeliveryQueueTokenUsageSumsLinkedRuns|TestDeliveryQueueLimitsRoundTripMaxTokens|TestRetentionLeavesTokenUsageTables|TestDurableTableLifecyclePolicyCoversEveryTable|TestEachPromptRecordsItsUsageWithTheOwnerScope|TestAFailedAndAStoppedPromptStillRecordUsage|TestAPromptWithoutAnOwnerRecordsASessionScope|TestAFailedUsageWriteLeavesTheOutcomeUnchanged|TestUsageEventProjectsToTheUsageCategory|TestDefaultFilterIncludesUsage|TestUsageFilterSelectsOnlyUsage|TestEventsHelpNamesTheUsageCategory|TestEventsHelpDocumentsAgentSelectionFilter)$" ./internal/store ./internal/daemon ./internal/runevent ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMigrationAddsTokenUsageTablesToThePreviousSchema TestAppendTokenUsageWritesTheRowAndItsRunEvent TestRunTokenUsageSumsScopesAndTheRun TestRunTokenUsageKeepsUnreportedPromptsOutOfTheTotal TestSessionSpendSumsIncreasesAndRestartsOnADrop TestDeliveryQueueRunsLinkEveryRecordedRunID TestDeliveryQueueTokenUsageSumsLinkedRuns TestDeliveryQueueLimitsRoundTripMaxTokens TestRetentionLeavesTokenUsageTables TestDurableTableLifecyclePolicyCoversEveryTable TestEachPromptRecordsItsUsageWithTheOwnerScope TestAFailedAndAStoppedPromptStillRecordUsage TestAPromptWithoutAnOwnerRecordsASessionScope TestAFailedUsageWriteLeavesTheOutcomeUnchanged TestUsageEventProjectsToTheUsageCategory TestDefaultFilterIncludesUsage TestUsageFilterSelectsOnlyUsage TestEventsHelpNamesTheUsageCategory TestEventsHelpDocumentsAgentSelectionFilter; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task seventeen of the nineteen named tests do not exist, so the command fails.
- `for pair in "docs/user-guide/run-database-lifecycle.md|run_token_usage" "docs/user-guide/run-database-lifecycle.md|delivery_queue_runs" "docs/user-guide/commands/events.md|agent-selection,usage" "docs/user-guide/commands/events.md|token_basis"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task neither guide names the tables or the category, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1 and 2; User Story 2; Core Features 3, 4, 5 and 6; Success Metrics 4 and 6; Declared breaks
- [_techspec.md](_techspec.md) — Interfaces; Data Models; Recording; Surfaces; Surface Transcript 7; API Contract 6; API Contract 7; Testing Approach 2; Testing Approach 3; Build Order 2
- ADR-0198; ADR-0199; ADR-0008; ADR-0033; ADR-0051; ADR-0098

## Result

Implemented Task 02's persistence, recording and event-stream slice. Schema 22
adds durable prompt usage, historical Delivery Queue Run links and `max_tokens`.
Each usage row commits with its `daemon.token_usage` event. Reports preserve
NULL totals for unreported prompts, first-recorded scope order, reported split
coverage and cumulative cost increases/resets per Run, Agent Session and
currency. Both queue update paths commit Run links with the item write.

Owned and unowned prompts record returned usage even after failure or
cancellation, using `context.WithoutCancel`. A write error emits the specified
progress warning while preserving the prompt result and settlement behavior.
The `usage` stream category is included by default and accepted by `--filter`;
its projection keeps optional fields absent and retains reported zero values.
`IsDaemonKind` and `roundfix-events/v1` remain unchanged. The events guide was
already present from Spec 0194 and now includes Surface Transcript 7. The
lifecycle policy names both new tables; retention and GC code are unchanged.

### Focused implementation evidence

The following focused command exited 0 after the final implementation and
fixture changes (all four packages reported `ok`):

```bash
GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test -count=1 ./internal/store ./internal/daemon ./internal/runevent ./internal/cli -run 'TokenUsage|SessionSpend|QueueRunsLink|RoundTripMaxTokens|RetentionLeaves|EachPrompt|StoppedPrompt|WithoutAnOwner|UsageWrite|UsageRecordsPreferred|UsageLinkFailure|UsageEvent|DefaultFilter|UsageFilter|EventsHelp|EventsUsageTranscript|DurableTableLifecycle'
```

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Previous-schema migration and constraints | `TestMigrationAddsTokenUsageTablesToThePreviousSchema` preserves the existing Run and queue through schema 21 → 22, checks scope/basis/token-nullability and foreign-key constraints, and rejects duplicate queue links. |
| One row and event per prompt, including failures/stops | `TestEachPromptRecordsItsUsageWithTheOwnerScope` covers task/qa/review scopes and both runner interfaces. `TestAFailedAndAStoppedPromptStillRecordUsage` checks real persisted rows and events under canceled contexts. `TestAPromptWithoutAnOwnerRecordsASessionScope` checks the unowned path. `TestUsageRecordsPreferredAndFallbackPrompts` checks both candidates, their sessions and active owner attempts. `TestAppendTokenUsageWritesTheRowAndItsRunEvent` checks row/payload fields and uses an event-rejection trigger to prove atomic rollback. |
| Scope, Run and queue totals exclude unreported prompts | `TestRunTokenUsageSumsScopesAndTheRun`, `TestRunTokenUsageKeepsUnreportedPromptsOutOfTheTotal`, `TestSessionSpendSumsIncreasesAndRestartsOnADrop` and `TestDeliveryQueueTokenUsageSumsLinkedRuns` cover first-recorded order, missing versus reported zero, partial split coverage, cumulative increases/resets, multiple currencies and Run/session identity. |
| Every item Run ID, including retry, is linked and counted | `TestDeliveryQueueRunsLinkEveryRecordedRunID` checks both write paths, deduplication and queue-replacement cascade. The queue sum test counts prior Runs and linked Runs without prompts. `TestDeliveryQueueUsageLinkFailureRollsBackTheItem` rejects links and checks that both item writes roll back. `TestDeliveryQueueLimitsRoundTripMaxTokens` checks persisted limits. |
| Write failures preserve outcome and warn | `TestAFailedUsageWriteLeavesTheOutcomeUnchanged` checks successful, failed and stopped prompt results on both dispatch paths; the real temporary-store Task-cycle companion checks unchanged settlement and the warning. |
| Usage filter prints Surface Transcript 7 | `TestUsageEventProjectsToTheUsageCategory`, `TestDefaultFilterIncludesUsage`, `TestUsageFilterSelectsOnlyUsage` and `TestEventsUsageTranscript` cover projection and real CLI/store replay. Both events-help tests pass without changing `internal/cli/cli_test.go`. |

`TestRetentionLeavesTokenUsageTables` confirms that pruning journal entries
preserves rows and links. `TestDurableTableLifecyclePolicyCoversEveryTable`
checks the policy against the migrated database. Local documentation inspection
confirmed both table names, `agent-selection,usage`, `token_basis` and Surface
Transcript 7. `git -c core.fsmonitor=false diff --check` exited 0.

### Sabotage evidence

Each mutation was applied alone, exercised, and restored before subsequent work:

| Gate | Deliberate mutation | Observed failure | Restoration |
| --- | --- | --- | --- |
| Sums | Allocate a zero token total for an unreported prompt in the accumulator. | `go test -count=1 ./internal/store -run '^TestRunTokenUsageKeepsUnreportedPromptsOutOfTheTotal$'` exited 1: `unreported became zero`. | Restored the accumulator; the final focused command above exited 0. |
| Queue links | Remove `linkDeliveryQueueRun` from `RetryDeliveryQueueItem`, leaving it only in `UpdateDeliveryQueueItem`. | `go test -count=1 ./internal/store -run '^TestDeliveryQueueRunsLinkEveryRecordedRunID$'` exited 1: `links=1`, expected 2. | Restored the retry link write; the final focused command above exited 0. |
| Outcome guard | Return the usage-write error from recording and propagate it from both prompt dispatch paths. | `go test -count=1 ./internal/daemon -run '^TestAFailedUsageWriteLeavesTheOutcomeUnchanged$'` exited 1 on owner/session results and the Task-cycle companion: `Completed:0 Failed:1`, reason `Agent failed: usage disk failure`. | Restored warning-only recording and original prompt returns; the final focused command above exited 0. |

All sabotage commands used `GOCACHE=/private/tmp/roundfix-task02-gocache rtk
proxy`; they were focused implementation checks, not the declared Verification.

### Regression checks and execution boundary

The first incremental check found duplicate-schema creation in existing
metadata-downgrade fixtures and sandbox process-table restrictions in two CLI
process-owner tests. Migration now follows the existing idempotent migration
pattern. Historical Delivery Queue fixtures in `internal/store/delivery_test.go`
and `internal/store/delivery_limits_test.go` remove schema-22 structures before
constructing older layouts; the limit-column fixture explicitly represents
schema 20 instead of the moving latest-minus-one version. Their assertions are
unchanged. A focused migration regression selection across store and CLI exited
0 after these corrections. An intermediate incremental run was stopped after
an added fallback test exposed its incorrect candidate-based attempt expectation;
the test now checks the active owner's lifecycle attempt, and the final focused
run exited 0.

No new dependency, command-set token limit, ceiling enforcement or usage-summary
surface is added here; those surfaces remain with later Tasks. No live Run
Database is opened by the new tests. The pre-existing Daemon-written
`status: in_progress` is preserved. No authored Verification command was run,
and no commit, push or Pull Request was created. The Daemon owns Verification
and Task settlement.

The final `GOCACHE=/private/tmp/roundfix-task02-gocache rtk make
verify-incremental` exited 0 with process-table access. It ran formatting checks,
`go vet ./...`, the repository test suite, skill sync/check gates and the build.
The CLI and Daemon suites ran successfully; unchanged package results reused
Go's cache. The focused usage checks above used `-count=1` and did not rely on
cached results. The initial default-cache attempt was blocked by sandbox access
to `~/Library/Caches/go-build`; the task-scoped cache resolved that restriction.

Read-only comparison with the starting commit confirmed schema 21 → 22 is one
increment and all five new usage source/test files were absent at that revision.
The only initial worktree change was the Daemon's Task status transition. Final
scope inspection leaves `_tasks.md`, every other Task and all prior-task files
unchanged. This is implementation evidence for the Daemon's next Verification
step, not a terminal Task verdict.

### Verification Feedback — attempt 1

Inspected the Daemon diagnostic at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261001T094428Z_f4d3047c47b08ac0/verification/batch-002-attempt-1.log`
and the affected test boundaries. The configured `make verify-changed` run
exhausted available disk space while writing temporary SQLite databases, Git
fixtures and Go build outputs. The failures did not identify a usage assertion
regression. The coverage and partition tests also failed while building the
test binaries they enumerate; their assertions were not weakened.

The Task-owned temporary cache `/private/tmp/roundfix-task02-gocache` occupied
6.4 GiB. Removed only that regenerable cache. `df -h /private/tmp` showed
available space increasing from 330 MiB to 6.9 GiB. No source, fixture,
Verification configuration or unrelated data was changed by this repair.

The following focused check exited 0 using the existing shared Go cache with
filesystem/process access. It selected all eight top-level tests identified in
the diagnostic and the Task's outcome guard; Daemon, Spec and verifyselect
packages each reported `ok`:

```bash
rtk proxy go test -count=1 ./internal/daemon ./internal/spec ./internal/verifyselect -run '^(TestQASettlementAcceptsAPassWithAResultsRow|TestStaleAuditorWarningRidesARefusedGate|TestSelfAuditSeedRecordsACurrentAuditor|TestQAPromptOutsideASelfAuditOmitsTheUserFlowBinary|TestQASettlementRefusesAHollowPass|TestQASettlementRefusesARewrittenAuditingBinary|TestCoverageEquivalence|TestPartitionFollowsTheMakefileRecipes|TestAFailedUsageWriteLeavesTheOutcomeUnchanged)$'
```

Did not rerun `make verify-changed` or either declared Task Verification
command. The Daemon owns the next full configured Verification sequence and
settlement. Task status remains unchanged; no commit, push or Pull Request was
created.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/store/delivery_limits_test.go`
- `internal/store/delivery_test.go`

## Carry-forward provenance

- Source Run: `run_20261001T094428Z_f4d3047c47b08ac0`
- Source commit: `4dcd2758b8eb731fe04ff7e4fd19b787c2509e40`
