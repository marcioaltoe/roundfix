---
task: task_03
spec: 0041-agent-selection-runtime-readiness
status: pending
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

- [ ] Build deterministic canonical-to-adapter assignment plans.
- [ ] Apply model and reasoning operations in the required order.
- [ ] Compare the complete effective state with the requested tuple.
- [ ] Reuse assignment semantics for live Agent Session setup.
- [ ] Add typed selection and effective-state failures.
- [ ] Prove cleanup ownership on every terminal path.
- [ ] Cover exact, unsupported, rejected, cancelled, and mismatched cases.

## Acceptance Criteria

- [ ] The official fixture proves `gpt-5.6-sol / high` and
      `gpt-5.5 / xhigh` exactly.
- [ ] Independent-control and model-variant fixtures can prove the same
      canonical tuple without storing transport-specific IDs in profiles.
- [ ] A requested non-empty effort never retries as model-managed or with a
      different model.
- [ ] A zero-exit application with mismatched effective state fails proof.
- [ ] Disposable and live Session tests issue the same ordered selection
      operations for equivalent requests.
- [ ] Every terminal path closes the disposable Session with bounded cleanup,
      including cancellation and joined failures.
- [ ] Proof records no Agent prompt and no model-token activity.

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

