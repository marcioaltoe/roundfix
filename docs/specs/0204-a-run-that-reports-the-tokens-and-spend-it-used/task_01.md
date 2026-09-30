---
task: task_01
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: pending
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
