---
task: task_01
spec: 0026-model-fallback-guardrail
status: completed
type: backend
complexity: medium
---

# Task 01: Prove a Fallback Selection with disposable sessions

## Overview

Add the agent-layer fallback probe: given the failed runtime selection and a
candidate set (Model Catalog order with the failed model excluded, plus the
runtime's effort vocabulary highest-first), prove the first functional Agent
Model with a disposable Agent Session and determine its highest functional
reasoning effort, classifying a working model with no settable effort as
model-managed. The slice is verifiable alone through the recording fake acpx
runner.

## Requirements

1. MUST walk candidate Agent Models in the given order, skip the failed
   model, and stop at the first candidate whose model assignment proves
   functional on a disposable Agent Session.
2. MUST determine the proven model's reasoning effort by trying the given
   effort values highest-first and stopping at the first accepted value.
3. MUST classify a proven model that accepts no effort value as a
   model-managed Fallback Selection with an empty reasoning effort.
4. MUST report no fallback when no candidate model proves functional,
   distinguishing that outcome from probe infrastructure errors.
5. MUST close every disposable Agent Session on every path and keep the
   unique preflight session naming.
6. MUST NOT cross to another ACP Runtime or mutate any configuration.

## Subtasks

- [x] Add the fallback candidate and Fallback Selection types and the probe
      entry point in the agent layer.
- [x] Implement candidate walking with failed-model exclusion and
      first-proof-wins ordering.
- [x] Implement highest-first effort proving with model-managed
      classification.
- [x] Guarantee disposable session cleanup on success, rejection, and error
      paths.
- [x] Cover ordering, exclusion, classification, no-candidate, and cleanup
      with unit tests on the recording fake runner.

## Acceptance Criteria

- [x] The probe returns the newest functional candidate and never proposes
      the failed model.
- [x] For a model that accepts efforts, the returned effort is the highest
      accepted value; for a model that accepts none, the returned effort is
      empty.
- [x] When no candidate proves, the probe reports no fallback without an
      infrastructure error.
- [x] Recorded acpx interactions show one disposable session per probed
      candidate, each closed, including on rejection paths.

## Verification

- `rtk go test ./internal/agent` - expected: fallback probe ordering,
  classification, no-candidate, and session-hygiene tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → Core Features 1, 6.
- `_techspec.md` → Interfaces; Build Order 1.
- ADR-0039, ADR-0040, ADR-0041.

## Result

Implemented the agent-layer fallback probe with caller-owned Agent Model and
reasoning-effort ordering. The probe skips the failed Agent Model, stops at the
first functional candidate, returns the first accepted highest-first effort,
and records an empty effort when the model manages reasoning. Selection
rejections continue probing; command, cancellation, and session-cleanup
failures remain infrastructure errors. Every attempted candidate uses one
uniquely named disposable preflight Agent Session that is closed through an
uncancelled cleanup context.

Verification:

- `rtk go test ./internal/agent`: passed — 113 tests, including fallback
  ordering, failed-model exclusion, effort ordering, model-managed
  classification, no-candidate behavior, cleanup failure, and cancellation
  cleanup.
- `rtk make verify`: passed — 1,055 tests across 19 packages, Roundfix Skill
  sync/checks, formatting checks, and the CLI build completed successfully.

Acceptance evidence:

- `TestACPXProbeFallbackUsesCandidateAndEffortOrder` records only the newest
  rejected candidate and the first functional candidate, excludes the failed
  and older models, and returns `high` after `xhigh` is rejected.
- `TestACPXProbeFallbackClassifiesModelManagedReasoning` rejects every supplied
  effort while preserving the proven model with an empty reasoning effort.
- `TestACPXProbeFallbackReportsNoFallbackWhenModelsAreRejected` returns
  `ok=false` with no error after every non-failed candidate rejects model
  assignment.
- The recording assertions prove one unique `roundfix-preflight-*` session per
  candidate and one matching close on success, rejection, no-candidate,
  cleanup-error, and cancellation paths; no probe sends a prompt.

Follow-ups: none.
