---
task: task_04
spec: 0041-agent-selection-runtime-readiness
status: pending
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

- [ ] Define the shared profile-readiness result and coordinator.
- [ ] Centralize tuple collection, deduplication, and stable references.
- [ ] Migrate `profiles validate` to the coordinator.
- [ ] Migrate operational Preflight Validation to the coordinator.
- [ ] Add bounded additive JSON and aligned text diagnostics.
- [ ] Prove no-mutation ordering and per-invocation proof behavior.

## Acceptance Criteria

- [ ] Shared tuples across categories and fallback positions are proved once
      and report every stable reference.
- [ ] `profiles validate` and operational preflight return the same
      classification and category evidence for the same failing config.
- [ ] Validation JSON remains `roundfix/profiles-validate/v1` and adds only
      optional bounded proof fields.
- [ ] A failed tuple creates no Run, worktree, artifact, prompt, commit, or
      repository change.
- [ ] Repeating the command creates fresh disposable proof rather than using a
      prior invocation's result.
- [ ] Successful readiness preserves each profile's configured fallback order
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
