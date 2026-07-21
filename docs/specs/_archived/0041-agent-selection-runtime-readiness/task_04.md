---
task: task_04
spec: 0041-agent-selection-runtime-readiness
status: completed
type: backend
complexity: high
---

# Task 04: Centralize Agent Selection Profile readiness

## Overview

Create one readiness coordinator for effective Agent Selection Profiles and
make profile validation and operational Preflight Validation consume it. The
coordinator must deduplicate exact tuples, preserve every category reference,
and produce the same proof classification before any Run-side mutation.

## Requirements

1. MUST resolve effective profiles and fallback chains through one coordinator
   used by `profiles validate` and Agent-starting operational preflight.
2. MUST deduplicate identical runtime/model/reasoning tuples while retaining
   stable category, source, inheritance, role, and fallback-index references.
3. MUST prove tuples sequentially in deterministic order and stop on the first
   failure without caching proof across command invocations.
4. MUST expose additive classification, encoding, adapter command/version, and
   bounded advertised values through the existing validation JSON schema.
5. MUST make text and JSON results identify the same failed tuple, affected
   categories, classification, and next action.
6. MUST finish readiness before Run, worktree, artifact, Session prompt, or
   repository mutation in every operational command.
7. MUST preserve configured Fallback Chains and ADR-0050's notification-first
   activation boundary.

## Subtasks

- [x] Define the shared profile-readiness result and coordinator.
- [x] Centralize tuple collection, deduplication, and stable references.
- [x] Migrate `profiles validate` to the coordinator.
- [x] Migrate operational Preflight Validation to the coordinator.
- [x] Add bounded additive JSON and aligned text diagnostics.
- [x] Prove no-mutation ordering and per-invocation proof behavior.

## Acceptance Criteria

- [x] Shared tuples across categories and fallback positions are proved once
      and report every stable reference.
- [x] `profiles validate` and operational preflight return the same
      classification and category evidence for the same failing config.
- [x] Validation JSON remains `roundfix/profiles-validate/v1` and adds only
      optional bounded proof fields.
- [x] A failed tuple creates no Run, worktree, artifact, prompt, commit, or
      repository change.
- [x] Repeating the command creates fresh disposable proof rather than using a
      prior invocation's result.
- [x] Successful readiness preserves each profile's configured fallback order
      for later notification-first activation.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/cli/profiles_validate.go`
- interface: `internal/cli/profile_preflight.go`
- interface: `internal/cli/selection_test.go`
- interface: `internal/config/profiles.go`
- interface: `internal/agent/fallback.go`

## Verification

- `rtk go test ./internal/cli -run 'Test(ProveProfileSelections|ProfilesValidate|ProfileOperationalPreflight)' -count=1` — expected: tuple deduplication, references, deterministic failures, aligned diagnostics, and fresh proof pass.
- `rtk go test ./internal/cli ./internal/daemon -run 'Test(ProfileOperationalPreflight|Fallback).*' -count=1` — expected: preflight mutation guards and notification-first fallback behavior pass.
- `rtk go test -race ./internal/cli ./internal/daemon -run 'Test(ProveProfileSelections|ProfilesValidate|ProfileOperationalPreflight|Fallback)' -count=1` — expected: shared readiness consumers are race-free.

## References

- `_prd.md` → User Stories 1, 3, and 7; Core Features 4 and 10; Success Metrics.
- `_techspec.md` → Shared Readiness Coordinator; Exact Disposable-Session
  Proof; Error Taxonomy and Diagnostics; Build Order 4.
- `../../adr/0050-configured-fallbacks-activate-after-notification.md` →
  notification-first fallback boundary.

## Result

Implemented one profile-readiness coordinator that resolves effective Agent
Selection Profiles, deduplicates exact tuples in stable order, and returns the
exact disposable-session proof evidence to both `profiles validate` and
operational Preflight Validation. The existing
`roundfix/profiles-validate/v1` response now has optional classification,
encoding, adapter command/version, bounded advertised-value, and next-action
fields. Text and JSON failures use the same failed proof record.

Acceptance evidence:

- Shared tuple and stable-reference coverage: `TestProveProfileSelectionsDeduplicatesReferencesAndStartsFreshProofPass`
  proves one request per unique tuple across required categories;
  `TestProveProfileSelectionsRetainsStableFallbackPositions` preserves
  `backend` fallback index 1 and `frontend` fallback index 2 on one shared
  proof.
- Consumer alignment and schema compatibility:
  `TestProfileOperationalPreflightMatchesProfilesValidateClassifiedFailure`
  compares the complete failed proof returned to both consumers, verifies the
  `reasoning_control_not_advertised` classification and bounded advertised
  values, and confirms the schema remains
  `roundfix/profiles-validate/v1` with aligned text and JSON next actions.
- No proof cache: the fresh-proof test invokes readiness twice and observes
  two complete proof passes rather than a reused result.
- No mutation before readiness: the passing CLI regression suite includes the
  existing Run database, Run/Task worktree, artifact, prompt, bypass-audit,
  commit, and repository-mutation guards for failed operational profile proof.
- Fallback preservation: the focused operational and daemon fallback tests
  pass with the configured order unchanged and ADR-0050's notification-first,
  pre-prompt activation boundary intact.

Verification:

- `rtk go test ./internal/cli -run 'Test(ProveProfileSelections|ProfilesValidate|ProfileOperationalPreflight)' -count=1`
  — passed, 6 tests in 1 package.
- `rtk go test ./internal/cli ./internal/daemon -run 'Test(ProfileOperationalPreflight|Fallback).*' -count=1`
  — passed, 3 tests in 2 packages.
- `rtk go test -race ./internal/cli ./internal/daemon -run 'Test(ProveProfileSelections|ProfilesValidate|ProfileOperationalPreflight|Fallback)' -count=1`
  — passed, 7 tests in 2 packages.
- `rtk go test ./internal/agent -count=1` — passed, 214 tests in 1
  package.
- `rtk go test ./internal/cli ./internal/daemon -count=1` — passed, 649
  tests in 2 packages.
- `rtk make verify` — passed: 1,628 Go tests in 20 packages, 79
  setup-context-driven tests, Repository Skill Set check, and Roundfix build.

Follow-ups: none discovered within Task 04's slice.
