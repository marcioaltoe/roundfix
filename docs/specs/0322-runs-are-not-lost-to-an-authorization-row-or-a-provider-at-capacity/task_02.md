---
task: task_02
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
status: pending
type: backend
complexity: high
---

# Task 02: A Provider System Error falls through the Fallback Chain during any Task or QA turn

## Overview

A provider that answers "Selected model is at capacity" ends a QA turn that
already wrote a verdict, or a Task's repair turn after its First Handoff,
with recovery `exhausted` although the profile names a next selection, and
the Delivery Queue parks the item `runtime-infrastructure`. This Task lets
the next Fallback Selection take any Task or QA turn after a Provider System
Error, resends a failed repair turn's prompt to the fallback, and reopens a
QA verdict the failed turn wrote, recording the verdict it discarded.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-a-provider-at-capacity-parks-a-turn-the-fallback-chain-could-finish.md`](references/2026-10-26-a-provider-at-capacity-parks-a-turn-the-fallback-chain-could-finish.md)
   of 2026-10-26 (`_prd.md` → Sources 2) by implementing API Contracts 3-5 of
   `_techspec.md`.
2. MUST, in `internal/daemon/agent_session_owner.go`
   `recoverProviderSystemError`, choose `fallback` whenever the scope is
   `task` or `qa` and a next candidate exists, with no First Handoff or
   report-verdict condition, and `exhausted` otherwise; the review scope's
   early return in `Run`, the event payload keys, the attempt record, the
   progress line and `providerSystemErrorExhaustedError` keep their bytes.
3. MUST, in `Run`, keep the failed turn's prompt and resend it after
   `providerSystemErrorNotice` when it differs from `owner.taskPrompt`, as
   API Contract 4 states.
4. MUST add `resetUnfinishedQAVerdict` with the signature and comment of
   `_techspec.md` → Interfaces and call it for every QA fallback that has a
   report path, before appending the `## Agent runtime fallback` section,
   which carries `- discarded_verdict: <verdict>` between the
   `provider_system_error` and `fallback` lines only when a verdict was
   reopened. It rewrites only the frontmatter `verdict:` line, and a report
   without one keeps its bytes.
5. MUST NOT change `recoverLostRollout`, `reportPending` or its use there,
   the Light Tier code (`applyLightTier`, `canYieldLight`, `yieldLight`,
   `escalate`), `activate`, `fallbackAfterSelectionFailure`,
   `internal/agent`, `internal/cli` or the Delivery Queue's park decision.
6. MUST change `internal/daemon/provider_system_error_recovery_test.go` and
   create `internal/daemon/provider_system_error_unfinished_qa_test.go`
   exactly as `_techspec.md` → Testing Approach describes, keeping every
   existing top-level test name so `TestCoverageEquivalence` reports no
   regression. Each new top-level test and subtest calls `t.Parallel()`
   first; the tests use the package's fake runners and runtime factories and
   never start an ACP adapter or reach the network.
7. MUST NOT change any document, skill, Baseline module, the Makefile or the
   CI workflows.

## Subtasks

- [ ] Remove the First Handoff and pending-verdict window for Task and QA scopes.
- [ ] Resend the failed repair prompt and reopen a QA verdict on fallback.
- [ ] Update the recovery tests and add the unfinished-QA tests.

## Acceptance Criteria

- [ ] A repair turn after the First Handoff and a QA turn that wrote `partial`, `pass` or `fail` both complete on the fallback after a Provider System Error.
- [ ] The QA fallback receives `verdict: pending` and the report records `- discarded_verdict: <verdict>`; a pending report records no such line.
- [ ] Without a next selection the Task or QA Task still settles `failed` as runtime infrastructure; review, Lost Rollout and Light Tier tests pass unchanged.

## Consumers

- the `daemon.task` `provider_system_error` event with `recovery: exhausted` — consumers: `internal/cli`
- the QA report's `## Agent runtime fallback` section — consumers: `internal/daemon`

## Context

- instruction: `internal/daemon/lost_rollout_recovery_test.go`
- instruction: `internal/daemon/light_tier_test.go`
- instruction: `internal/daemon/light_tier_access_test.go`
- instruction: `internal/daemon/task_engine.go`
- instruction: `internal/cli/deliver_provider_system_error_test.go`
- instruction: `internal/spec/coverage_test.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `internal/daemon/provider_system_error_recovery_test.go`
- creates: `internal/daemon/provider_system_error_unfinished_qa_test.go`

## Verification

- `out="$(go test -timeout 15m -count=1 -v -run '^(TestProviderSystemErrorAfterFirstHandoffActivatesTheFallback|TestProviderSystemErrorAfterFirstHandoffSettlesRuntimeInfrastructure|TestProviderSystemErrorBeforeFirstHandoffActivatesTheFallback|TestProviderSystemErrorWithoutFallbackSettlesRuntimeInfrastructure|TestQAProviderSystemErrorBeforeTheReportFallsBackAndRecordsIt|TestQAProviderSystemErrorExhaustionSettlesInsteadOfHalting|TestProviderSystemErrorRecoveryEligibility|TestProviderSystemErrorDoesNotUseLostRolloutBudget|TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt|TestQAProviderSystemErrorAfterAVerdictWithoutFallbackSettlesRuntimeInfrastructure|TestLostRolloutBeforeFirstHandoffActivatesTheFallback|TestLostRolloutAfterFirstHandoffResumesInANewSession|TestLostRolloutWithoutFallbackResumesTheSameSelection|TestThirdLostRolloutSettlesRuntimeInfrastructure|TestQALostRolloutBeforeTheReportFallsBackAndRecordsIt|TestQALostRolloutExhaustionSettlesInsteadOfHalting|TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess|TestLightTierYieldsToTheCategoryProfileBeforeAgentWork|TestLightTierDoesNotYieldAfterWorkOrOutsideTheLightTier|TestLightTierDispatchOrder|TestLightTierSkipsWithAWarning|TestLightTierEscalatesOnce)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestProviderSystemErrorAfterFirstHandoffActivatesTheFallback TestProviderSystemErrorRecoveryEligibility/QA_with_a_finished_report_falls_back TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt/partial TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt/pass TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt/fail TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt/a_pending_report_records_no_discarded_verdict TestQAProviderSystemErrorAfterAVerdictWithoutFallbackSettlesRuntimeInfrastructure; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -timeout 15m -count=1 -run '^TestCoverageEquivalence$' ./internal/spec` — expected: exit 0. Before this Task the new tests and the renamed eligibility case do not exist, so the first pass check fails. Afterwards a repair turn and a QA turn that wrote a verdict fall through to the next selection, the QA fallback starts from a reopened verdict, an exhausted chain still settles as runtime infrastructure, the Lost Rollout and Light Tier tests pass unchanged, and no top-level test left the coverage record.
- `go test -timeout 15m -count=1 -run '^TestDeliveryRunResultReadsAnExhaustedProviderSystemError$' ./internal/cli && out="$(go test -timeout 15m -count=1 -v -run '^TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt$' ./internal/daemon 2>&1)" && printf '%s\n' "$out" | grep -q -- '--- PASS: TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt/partial '` — expected: exit 0. Before this Task the QA test does not exist, so the pass check fails. Afterwards the Delivery Queue still parks an exhausted Provider System Error as runtime infrastructure and the Daemon reopens a partial verdict.

## References

- `_prd.md` → Sources 2; Goals 3-4; Core Feature 2; Success Metrics 4-6; Validation
- `_techspec.md` → Current behavior 4-7; Interfaces; Data Models; API Contracts 3-5; Invariants 3-6; Testing Approach; Build Order 2
- ADR-0335; ADR-0274
