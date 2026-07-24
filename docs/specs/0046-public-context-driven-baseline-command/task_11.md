---
task: task_11
spec: 0046-public-context-driven-baseline-command
status: completed
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

- [x] Implement the human workflow state driver and prompt adapter.
- [x] Connect adoption, update, profile, decision, and classification states.
- [x] Render consolidated plan review and exact confirmation.
- [x] Integrate approved apply and Baseline verification.
- [x] Add synchronous interaction and no-TTY macro tests.

## Acceptance Criteria

- [x] A greenfield adoption completes from preflight through verified apply in one invocation.
- [x] A preservation adoption presents one consolidated editable classification review.
- [x] A configured repository enters update with its current profile and an explicit change action.
- [x] Equivalent human and automation answers produce identical plan bytes and digest.
- [x] Declining final confirmation produces zero writes.
- [x] Redirected or absent terminal input never causes a hidden prompt or guessed answer.
- [x] The workflow uses no Bubble Tea or independent setup engine.

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

## Result

Implemented the root `roundfix baseline` human workflow as a synchronous,
numbered prompt driver over the existing deterministic `BuildPlan` and
`ApplyPlan` authority. The driver detects unconfigured, incompatible, and
configured repositories; preserves current update decisions; offers an
explicit Baseline Profile change; collects conditional decisions through the
shared catalog decision graph; presents one consolidated editable
classification review; renders `fileChanges` before the complete managed-entry
and Upgrade Retention Contract ledgers; and applies only after one confirmation
bound to the displayed Plan Digest. Redirected or absent input returns a
structured `baseline plan`/`baseline apply` next action before any prompt.

The implementation also corrected two deterministic-engine issues exposed by
the real human journeys: rejected preservation entries now retain a canonical
empty targets array, and an intact current Setup Manifest may change from one
current Baseline Profile to another without being misclassified as an
unsupported legacy transition.

### Acceptance evidence

1. `TestHumanBaselineAdoption` drove a greenfield repository through numbered
   choices, consolidated review, exact confirmation, verified apply, and an
   observed Setup Manifest write in one invocation.
2. `TestConsolidatedReview` showed one consolidated preservation proposal,
   allowed proposal editing, rendered `fileChanges` before both complete
   ledgers, and declined without mutation.
   `TestConsolidatedReviewEditsManagedClassification` exercised an edited
   managed-entry rejection through a valid complete Plan.
3. `TestHumanBaselineUpdate` started from a real verified apply, reported
   `Current Baseline Profile: go-cli-tui`, offered `Change Baseline Profile`,
   selected `rust-cli`, and proved the changed-profile Plan matched equivalent
   automation input.
4. `TestHumanAutomationPlanParity` compared the complete marshaled Plan bytes
   and digest from equivalent human and automation answers. The update
   profile-change branch performs the same byte comparison.
5. `TestHumanBaselineUpdate` and `TestConsolidatedReview` compared repository
   trees before and after declining the digest-bound final confirmation and
   observed zero writes.
6. `TestBaselineNoTTY` exercised both the injected no-terminal boundary and a
   real built binary with redirected stdin. Both exited with structured
   `interactive_input` guidance, emitted no prompt, guessed no answer, and
   wrote nothing.
7. The human driver imports no Bubble Tea package and delegates planning,
   application, rollback, and Baseline verification to the existing
   `internal/baseline` engine.

### Verification

- `rtk proxy env GOCACHE=/tmp/roundfix-task11-go-cache go test -count=1 ./internal/cli ./internal/baseline -run 'TestHumanBaselineAdoption|TestHumanBaselineUpdate|TestConsolidatedReview|TestHumanAutomationPlanParity|TestBaselineNoTTY'`
  passed.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline --help` passed and
  describes the one human workflow, exact confirmation, no-terminal refusal,
  and automation `plan`/`apply` path.
- `rtk proxy env GOCACHE=/tmp/roundfix-task11-go-cache make verify` passed:
  1,945 Go tests, 256 canonical setup-context-driven tests, 256 distributed
  setup-context-driven tests, asset validation, Roundfix skill checks, and the
  final build.
- `git -c core.fsmonitor=false diff --check` passed.

### Follow-up

- The later documentation/thin-skill slice must publish the new root human
  recipe before any CLI-changing pull request is opened, per
  `docs/agents/skill-governance.md`.
