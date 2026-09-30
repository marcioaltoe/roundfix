---
task: task_02
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: pending
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
