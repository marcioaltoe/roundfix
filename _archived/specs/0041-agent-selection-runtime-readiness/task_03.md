---
task: task_03
spec: 0041-agent-selection-runtime-readiness
status: completed
type: backend
complexity: high
---

# Task 03: Prove exact advertised Agent Selections

## Overview

Turn advertised capabilities into a deterministic assignment plan and prove
the exact requested runtime/model/reasoning tuple through a disposable Agent
Session. The same assignment logic must prepare live Agent Sessions so
preflight cannot approve a tuple that Agent work applies differently.

## Requirements

1. MUST plan exact selections in the TechSpec order: advertised model,
   independent reasoning control, unambiguous model variant, then unsupported.
2. MUST treat an empty reasoning effort as explicit model-managed intent and
   never use it as recovery for a rejected non-empty effort.
3. MUST apply the plan, consume the complete returned configuration state, and
   compare effective model and reasoning with the requested tuple.
4. MUST use the same planner and application semantics for disposable and live
   Agent Sessions.
5. MUST classify unsupported controls, rejected selections, effective-state
   mismatches, invalid evidence, and cleanup failures distinctly.
6. MUST close every disposable Session on success, rejection, cancellation,
   timeout, malformed evidence, and joined setup/cleanup failure.
7. MUST send no Agent prompt and consume no model tokens during proof.

## Subtasks

- [x] Build deterministic canonical-to-adapter assignment plans.
- [x] Apply model and reasoning operations in the required order.
- [x] Compare the complete effective state with the requested tuple.
- [x] Reuse assignment semantics for live Agent Session setup.
- [x] Add typed selection and effective-state failures.
- [x] Prove cleanup ownership on every terminal path.
- [x] Cover exact, unsupported, rejected, cancelled, and mismatched cases.

## Acceptance Criteria

- [x] The official fixture proves `gpt-5.6-sol / high` and
      `gpt-5.5 / xhigh` exactly.
- [x] Independent-control and model-variant fixtures can prove the same
      canonical tuple without storing transport-specific IDs in profiles.
- [x] A requested non-empty effort never retries as model-managed or with a
      different model.
- [x] A zero-exit application with mismatched effective state fails proof.
- [x] Disposable and live Session tests issue the same ordered selection
      operations for equivalent requests.
- [x] Every terminal path closes the disposable Session with bounded cleanup,
      including cancellation and joined failures.
- [x] Proof records no Agent prompt and no model-token activity.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/sessions.go`

## Verification

- `rtk go test ./internal/agent -run 'Test(PlanSelectionAssignment|ProveExactSelection|ApplySessionSelection)' -count=1` — expected: independent, model-variant, model-managed, unsupported, rejected, and mismatched selections pass their exact assertions.
- `rtk go test ./internal/agent -run 'TestProveExactSelection.*(Cleanup|Cancel|Timeout|NoPrompt)' -count=1` — expected: cleanup is observed on every path and no Agent prompt is sent.
- `rtk go test -race ./internal/agent -run 'Test(PlanSelectionAssignment|ProveExactSelection|ApplySessionSelection)' -count=1` — expected: disposable and live selection paths are race-free.

## References

- `_prd.md` → User Stories 2–4; Core Features 2, 3, and 10; Success Metrics.
- `_techspec.md` → Assignment Planning; Exact Disposable-Session Proof; Error
  Taxonomy and Diagnostics; Build Order 3.
- `../../adr/0039-model-availability-preflight-uses-a-disposable-agent-session.md`
  → disposable-session proof and cleanup.
- `../../adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md`
  → exact tuple intent and adapter-specific assignment.

## Result

Added deterministic `SelectionAssignment` planning in the documented order:
advertised model, independent reasoning control, exact model variant, then a
typed unsupported result. Empty reasoning remains explicit model-managed
intent; a rejected non-empty effort is never retried with another effort or
model.

Both disposable proof and live Agent Session setup now use the same ordered
application path. Each strict ACPX operation consumes and validates the
complete returned configuration state, and the final effective model and
reasoning must equal the canonical requested tuple. The live prompt path keeps
the planned adapter model, including model variants, instead of reapplying the
canonical profile model.

Added bounded typed failures for unsupported controls, rejected selections,
effective-state mismatches, invalid capability evidence, and Session cleanup.
Disposable proof owns bounded cleanup independently of caller cancellation and
joins setup and cleanup failures when both occur. Proof invokes only Session
setup/configuration/close operations; it never invokes `prompt`.

Acceptance evidence:

- `TestProveExactSelectionOfficialFixturesNoPrompt` proves
  `gpt-5.6-sol / high` and `gpt-5.5 / xhigh` exactly and observes no prompt.
- `TestPlanSelectionAssignment` and the independent/model-variant proof tests
  prove canonical requests without transport-specific profile identifiers.
- Rejection tests assert the exact operation list, with no model-managed or
  alternate-model retry for a requested non-empty effort.
- `TestProveExactSelectionMismatchClosesSession` proves a successful ACPX
  command with mismatched effective state fails proof and still closes the
  disposable Session.
- `TestApplySessionSelectionDisposableAndLiveOrder` proves equivalent
  disposable and live requests emit the same ordered selection operations.
- Cleanup, cancellation, timeout, malformed-evidence, and joined-failure tests
  prove every terminal path closes with bounded cleanup ownership.
- No-prompt assertions and the fake ACPX action logs prove proof performs no
  Agent prompt or model-token-producing operation.

Verification:

- `rtk go test ./internal/agent -run 'Test(PlanSelectionAssignment|ProveExactSelection|ApplySessionSelection)' -count=1`: passed, 23 tests.
- `rtk go test ./internal/agent -run 'TestProveExactSelection.*(Cleanup|Cancel|Timeout|NoPrompt)' -count=1`: passed, 8 tests.
- `rtk go test -race ./internal/agent -run 'Test(PlanSelectionAssignment|ProveExactSelection|ApplySessionSelection)' -count=1`: passed.
- `rtk make verify`: passed with 1,625 Go tests, 79
  setup-context-driven tests, Roundfix skill synchronization, and the CLI
  build.

Verification Feedback repair (attempt 1):

- Replaced the timeout test's setup-wide one-second wall-clock deadline with a
  test-controlled deadline that expires only after the blocked selection
  operation reports readiness. This preserves the `context.DeadlineExceeded`
  and cleanup assertions while making the test independent of subprocess
  startup overhead under race instrumentation.
- `rtk go test -race ./internal/agent -run '^TestProveExactSelectionTimeoutCleanup$' -count=1`: passed, 1 test.
- `rtk go test -race ./internal/agent -run 'Test(PlanSelectionAssignment|ProveExactSelection|ApplySessionSelection)' -count=1`: passed, 23 tests.
