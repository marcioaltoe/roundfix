---
task: task_01
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: completed
type: backend
complexity: medium
---

# Task 01: A prompt's result carries the usage its adapter reported

## Overview

The acpx runner drops every `usage_update` notification and the `usage` object
of the prompt response. This Task keeps both and reduces them to one Turn Usage
per prompt with the counting rule of ADR-0198, returned on the prompt's result.
It is verifiable on its own: fed a recorded Codex prompt, the result carries
the sum of the context readings with basis `request-sum`; fed a recorded Claude
prompt, it carries the reported turn total, its split and the reported cost.

## Requirements

1. MUST add `TurnUsage`, `UsageBasis`, `ReportedCost` and the counting function in the new file `internal/agent/usage.go`, with the shapes and the four steps of `_techspec.md` → Interfaces and → The counting rule. The zero value of `TurnUsage` MUST mean an unreported prompt.
2. MUST add a field to `adapterLineageContract` that marks a lineage whose prompt response covers only its last model request, set for the `codex` lineage only. The counting function MUST take that field, not a runtime name compared as a string elsewhere.
3. MUST read, during a prompt, every `session/update` whose `sessionUpdate` is `usage_update` and the `usage` object of a result line, as `_techspec.md` → Reading the reports states. A `usage_update` MUST publish no Agent Run Event and MUST NOT count as Agent output. A later result without `usage` MUST NOT erase an earlier one.
4. MUST ignore a `usage_update` or `usage` payload that does not parse, and MUST NOT fail, retry or reclassify the prompt because of it.
5. MUST add `Usage TurnUsage` to `ExecuteResult` and set it on every return path of `RunPrompt` after the stream was read: a normal end, a stopped prompt and a nonzero acpx exit after a parsed result.
6. MUST add the two fixtures `internal/agent/testdata/usage/codex-turn.ndjson` and `internal/agent/testdata/usage/claude-turn.ndjson` with the lines `_techspec.md` → Testing Approach 1 lists, in acpx JSON-RPC form, with every session ID written `sess_fixture`.
7. MUST put the new tests in `internal/agent/usage_test.go` and drive the stream tests through the existing fake acpx harness of `internal/agent/acpx_runner_test.go`, which this Task reads and does not change. No test starts acpx or an adapter.
8. MUST keep every existing test in `internal/agent` passing unchanged; a result without usage MUST compare equal to the result it produced before this Task.
9. MUST prove each new gate can fail. The Result MUST record one sabotage of the counting rule (for example summing readings for every lineage) and one of the stream reading (for example letting a later result erase the usage), each with the test that failed, and that the code was restored.

## Subtasks

- [ ] Add the Turn Usage types and the counting rule.
- [ ] Mark the Codex lineage and read both reports from the stream.
- [ ] Return the Turn Usage on every result path.
- [ ] Add the fixtures and the tests, and record one sabotage per gate.

## Acceptance Criteria

- [ ] The Codex fixture yields basis `request-sum`, the sum of its three readings, no split, and three readings.
- [ ] The Claude fixture yields basis `turn`, total 4019995, its split, and the cost `3.5526034999999996 USD`.
- [ ] A last-request lineage whose reported total exceeds its last reading is counted as `turn`.
- [ ] A prompt with neither report yields the zero `TurnUsage`, and a malformed payload never fails the prompt.
- [ ] A stopped prompt returns the usage read before it stopped.

## Context

- creates: `internal/agent/usage.go`
- creates: `internal/agent/usage_test.go`
- creates: `internal/agent/testdata/usage/codex-turn.ndjson`
- creates: `internal/agent/testdata/usage/claude-turn.ndjson`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acp_stream.go`
- interface: `internal/agent/agent.go`
- instruction: `internal/agent/acpx_runner_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCountTurnUsageSumsReadingsForALastRequestLineage|TestCountTurnUsageTakesTheReportedTurnForOtherLineages|TestCountTurnUsageCountsAGrownLastRequestReportAsTurn|TestCountTurnUsageLeavesANoReportPromptUnreported|TestRunPromptReturnsTheCodexFixtureUsage|TestRunPromptReturnsTheClaudeFixtureUsageAndCost|TestRunPromptKeepsUsageWhenALaterResultHasNone|TestRunPromptIgnoresAMalformedUsagePayload|TestRunPromptReturnsUsageForAStoppedPrompt|TestOnlyTheCodexLineageReportsItsLastRequestOnly)$" ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCountTurnUsageSumsReadingsForALastRequestLineage TestCountTurnUsageTakesTheReportedTurnForOtherLineages TestCountTurnUsageCountsAGrownLastRequestReportAsTurn TestCountTurnUsageLeavesANoReportPromptUnreported TestRunPromptReturnsTheCodexFixtureUsage TestRunPromptReturnsTheClaudeFixtureUsageAndCost TestRunPromptKeepsUsageWhenALaterResultHasNone TestRunPromptIgnoresAMalformedUsagePayload TestRunPromptReturnsUsageForAStoppedPrompt TestOnlyTheCodexLineageReportsItsLastRequestOnly; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the ten named tests exists, so the command fails.
- `test -f internal/agent/testdata/usage/codex-turn.ndjson && test -f internal/agent/testdata/usage/claude-turn.ndjson && grep -q '"usage_update"' internal/agent/testdata/usage/codex-turn.ndjson && grep -q '"cost"' internal/agent/testdata/usage/claude-turn.ndjson && ! grep -rqE '01a0[0-9a-f]{4}-|019f[0-9a-f]{4}-' internal/agent/testdata/usage` — expected: exit 0; before this Task the fixtures do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1 and 4; Core Features 1, 2 and 3; Success Metric 1; Recorded limits
- [_techspec.md](_techspec.md) — Interfaces; Reading the reports; The counting rule; Testing Approach 1; Build Order 1
- ADR-0198; ADR-0017; ADR-0020


## Result

Implemented the Task 01 slice: each prompt collects `usage_update` readings
and the last successfully parsed non-null result `usage`, counts them using
`adapterLineageContract.LastRequestOnly` (true only for Codex), and returns
`ExecuteResult.Usage`. The zero value remains unreported. Turn reports retain
optional split fields, including explicit zero; request sums omit the split.
The last reported cumulative cost is retained without pricing or calculating
a session delta. Usage notifications publish no Agent Run Event and do not
mark Agent work started. Malformed usage metadata is ignored without changing
the prompt outcome. Usage is attached on both stream-return branches before
normal, stopped, failure and transport-anomaly outcomes are classified.

Added both JSON-RPC fixtures with `sess_fixture` session IDs and all new tests
in `internal/agent/usage_test.go`. Stream tests use the existing fake-acpx
helpers; `internal/agent/acpx_runner_test.go` is unchanged. No live acpx or
adapter is started. The committed baseline has no usage types, new tests or
fixtures (`git ls-tree HEAD` for these paths returned no entries).

### Acceptance evidence

The focused usage selection and the fresh complete agent-package run below
exercise these criteria:

| Criterion | Evidence |
| --- | --- |
| Codex fixture: request sum, no split, three readings | `TestRunPromptReturnsTheCodexFixtureUsage` returns 214009 = 26830 + 36330 + 150849, basis `request-sum`, three readings and nil split fields; it also requires no messages or Run Events. `TestCountTurnUsageSumsReadingsForALastRequestLineage` covers absent, equal and smaller reported totals. |
| Claude fixture: turn total, split and cost | `TestRunPromptReturnsTheClaudeFixtureUsageAndCost` returns 4019995, basis `turn`, input 19000, output 995, cached read 3900000, cached write 100000, thought 0 and cost 3.5526034999999996 USD. It also requires no Run Events. |
| Last-request report grows beyond its last reading | `TestCountTurnUsageCountsAGrownLastRequestReportAsTurn` requires `turn` and the reported split when total 200 exceeds the last reading 100; it also covers a report without readings. |
| No report stays unreported; malformed usage does not fail | `TestCountTurnUsageLeavesANoReportPromptUnreported` checks the zero value for both counting modes and an actual fake prompt. `TestRunPromptIgnoresAMalformedUsagePayload` checks invalid field types, invalid result usage shapes, malformed costs and preservation of earlier valid metadata. `TestRunPromptKeepsUsageWhenALaterResultHasNone` covers absent and null later usage. |
| Stopped prompt retains earlier usage | `TestRunPromptReturnsUsageForAStoppedPrompt` checks exit 130 with and without a parsed result: 42 request-sum tokens, one reading and cost 1.5 USD survive the StopError. Both context-cancellation and normal stream-return branches assign the collected usage before returning. |

Additional guards: `TestOnlyTheCodexLineageReportsItsLastRequestOnly` checks the
contract map and drives both Claude and an unknown lineage with the Codex
fixture; both retain the reported turn rather than sum readings.
`TestRunPromptPreservesUsageAfterTransportAnomaly` checks parsed-result success
with nonzero exit. `TestRunPromptUsageDoesNotChangeNoOutputClassification`
requires a usage-only stream without a result to remain a Selection Failure
without Agent work-started status.

### Focused checks

- `GOCACHE=/tmp/roundfix-task01-gocache rtk proxy go test -count=1 ./internal/agent -run 'TestCountTurnUsage|TestRunPrompt.*Usage|TestOnlyTheCodexLineage'` — exit 0.
- `GOCACHE=/tmp/roundfix-task01-gocache rtk proxy go test -count=1 ./internal/agent` — exit 0 after restoring both sabotages; existing tests were not changed.
- `GOCACHE=/tmp/roundfix-task01-gocache rtk make verify-incremental` — initial sandbox run exited 2. Two CLI force-stop tests could not read the process table (`operation not permitted`). `TestTaskBudgetReasonNamesTheSettlementThatRenewedIt` also missed its 200 ms deadline before reaching its expected step (233.681541 ms elapsed).
- `GOCACHE=/tmp/roundfix-task01-gocache rtk proxy go test -count=1 -run '^TestRunForceStop(LegacyRunWithoutOwnerIdentityStillStopsOwner|OwnerProcessIntegrationProvesExitBeforeStoreCompletion)$' ./internal/cli` — exit 0 with process-table access outside the sandbox.
- `GOCACHE=/tmp/roundfix-task01-gocache rtk proxy go test -count=1 -run '^TestTaskBudgetReasonNamesTheSettlementThatRenewedIt$' ./internal/daemon` — exit 0 in isolation without code or timing changes. The initial timing failure remains recorded here rather than being treated as a fix.
- `GOCACHE=/tmp/roundfix-task01-gocache rtk make verify-incremental` — rerun with process-table access exited 0: formatting, vet, package tests, skill checks and build. The incremental tier reused successful package caches; CLI and daemon suites reran successfully.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

### Sabotage evidence

1. Counting: temporarily removed `lastRequestOnly` from the summation
   condition, making every lineage eligible for a request sum. The initial
   test with a large Claude report did not catch this because the growth
   guard still selected `turn`. Added the required contrasting case of a
   non-Codex reported total below its last context reading, checked it with
   restored code (exit 0), then repeated the same sabotage.
   `go test -count=1 -run '^TestCountTurnUsageTakesTheReportedTurnForOtherLineages$' ./internal/agent`
   exited 1: got `request-sum` 1100, expected `turn` 100 with its output
   split. Restored the production file byte-for-byte.
2. Stream reading: temporarily allowed an absent or null later result usage
   to overwrite the collector's earlier report.
   `go test -count=1 -run '^TestRunPromptKeepsUsageWhenALaterResultHasNone$' ./internal/agent`
   exited 1: got unreported zero usage, expected `turn` 42 with explicit
   output zero. Restored the production file byte-for-byte.

Both sabotage commands used `rtk proxy` and the same task-scoped GOCACHE.
The complete agent-package check and incremental check passed after restoration.

### Ownership and follow-ups

The existing `status: in_progress` is Daemon-owned and was preserved. No
command from this Task's Verification section was run. No Task Graph, other
Task file, existing test, tooling configuration, commit or publication was
changed. Usage persistence, session cost deltas, usage Run Events and public
surfaces belong to later Tasks. The unrelated 200 ms budget-test failure under
the initial full-suite load is a follow-up observation; this diff does not
change that test or its deadline. Task settlement and authored Verification
remain with the Daemon.

## Carry-forward provenance

- Source Run: `run_20261001T094428Z_f4d3047c47b08ac0`
- Source commit: `dc413ede29c38b83890b08d36d4efee5ed35f841`
