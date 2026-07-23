---
task: task_11
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 11: Guide human Baseline workflows

## Overview

Expose the complete adoption and update lifecycle through one
`roundfix baseline` invocation. Linear prompts drive the same deterministic
engine used by automation, ending in one consolidated review and exact
digest-bound apply.

## Requirements

1. MUST detect configured, unconfigured, and incompatible repositories and
   enter update or adoption without bypassing preflight, audit, decisions,
   review, apply, or Baseline verification.
2. MUST collect preservation mode, one Baseline Profile, required repository
   decisions, divergences, and editable classifications through numbered
   linear prompts.
3. MUST present `fileChanges` first while retaining the complete
   managed-entry and retention ledgers for inspection.
4. MUST require one explicit final confirmation bound to the displayed Plan
   Digest before mutation.
5. MUST offer explicit profile change during update and produce the same Plan
   Digest as equivalent non-interactive inputs.
6. MUST refuse to prompt without an interactive terminal and direct the caller
   to `baseline plan` with a structured next action.

## Subtasks

- [ ] Implement the human workflow state driver and prompt adapter.
- [ ] Connect adoption, update, profile, decision, and classification states.
- [ ] Render consolidated plan review and exact confirmation.
- [ ] Integrate approved apply and Baseline verification.
- [ ] Add synchronous interaction and no-TTY macro tests.

## Acceptance Criteria

- [ ] A greenfield adoption completes from preflight through verified apply in one invocation.
- [ ] A preservation adoption presents one consolidated editable classification review.
- [ ] A configured repository enters update with its current profile and an explicit change action.
- [ ] Equivalent human and automation answers produce identical plan bytes and digest.
- [ ] Declining final confirmation produces zero writes.
- [ ] Redirected or absent terminal input never causes a hidden prompt or guessed answer.
- [ ] The workflow uses no Bubble Tea or independent setup engine.

## Context

- instruction: `docs/adr/0068-baseline-command-uses-one-confirmation-gated-workflow.md`
- interface: `internal/cli/profiles_configure.go`
- interface: `internal/cli/setup.go`

## Verification

- `rtk go test -count=1 ./internal/cli ./internal/baseline -run 'TestHumanBaselineAdoption|TestHumanBaselineUpdate|TestConsolidatedReview|TestHumanAutomationPlanParity|TestBaselineNoTTY'` — expected: adoption, update, review, parity, confirmation, and no-TTY cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline --help` — expected: help describes one human workflow and directs automation to plan/apply subcommands.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 1 and 3–5; Core Features 1–5, 9–10, 14; User Experience.
- `_techspec.md` → System Architecture; API Contracts: Human workflow; Build Order 7.
- ADR-0068 → one confirmation-gated human workflow.
