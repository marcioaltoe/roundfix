---
task: task_05
spec: 0035-agent-selection-profiles
status: completed
type: backend
complexity: high
---

# Task 05: Prove relevant profiles before operational Runs

## Overview

Move Agent Selection readiness ahead of every mutating Run boundary. Implement and review commands resolve only their relevant Work Categories, apply explicit one-Run Preferred overrides, and prove every distinct preferred and fallback tuple before persistence or worktree creation.

## Requirements

1. MUST validate the complete requested Task Graph and Task Types before resolving categories; QA adds `qa`, while resolve/watch use only `review` and fetch remains Agent-free.
2. MUST resolve complete profiles for every relevant category and apply invocation agent/model/reasoning flags as Preferred Selection overrides without removing fallback chains.
3. MUST warn in text and JSON metadata when one invocation override is applied across multiple Task or QA categories.
4. MUST order proof by canonical category, Preferred Selection, then fallback index and deduplicate exact tuples for execution while retaining all category references.
5. MUST prove tuples sequentially through disposable sessions and close all sessions on every success, failure, cancellation, and panic-free early return.
6. MUST block before Run row, branch, worktree, durable selection record, or Agent work when configuration or any proof fails.
7. MUST name the category, failed tuple, affected references, adapter cause, and concrete profile remediation without proposing an unconfigured fallback.

## Subtasks

- [x] Resolve Implement and review Work Categories from requested work.
- [x] Apply complete one-Run Preferred Selection overrides.
- [x] Produce stable, deduplicated proof plans with category references.
- [x] Execute sequential disposable proofs and cleanup.
- [x] Place proof before every Run mutation boundary.
- [x] Add operational zero-side-effect and override-warning tests.

## Acceptance Criteria

- [x] A mixed backend/frontend Task Graph plus QA proves all distinct configured tuples once in stable order.
- [x] Resolve and watch prove only `review`; fetch creates no Agent Session.
- [x] One failed preferred or fallback proof prevents all Run, branch, worktree, prompt, and durable selection side effects.
- [x] Invocation overrides affect every relevant Preferred Selection, preserve each Fallback Chain, and emit the cross-category warning.
- [x] Exact duplicate tuples are probed once while results remain attributable to every category.
- [x] No benchmark rank, runtime-owned default, or dynamically discovered catalog candidate changes the configured proof plan.

## Context

- interface: `internal/cli/implement.go`
- interface: `internal/cli/selection.go`
- interface: `internal/cli/cli.go`
- interface: `internal/preflight/preflight.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/store/store.go`

## Verification

- `rtk go test ./internal/cli -run 'Test(ProfileOperationalPreflight|ImplementProfilePreflight|ReviewProfilePreflight|InvocationProfileOverride)' -count=1` — expected: relevant-category resolution, stable deduplication, zero-side-effect failure, and override warnings pass.
- `rtk go test ./internal/spec ./internal/config ./internal/agent ./internal/cli -run 'Test(TaskType|ProfileResolver|ProfileProof|ProfileOperationalPreflight)' -count=1` — expected: the complete pre-Run contract passes across owning packages.

## References

- `_prd.md` → Goals 5-7; User Stories 1-2, 5-7; Core Features 7-8 and 11; Success Metrics.
- `_techspec.md` → Configuration precedence; Operational preflight; Error and exit behavior; Build Order 5.

## Result

Implemented operational Agent Selection Profile proof before mutating Run boundaries. Implement derives Task Type categories from non-completed Tasks and adds `qa` only when requested. Resolve and watch prove only `review`; fetch remains Agent-free. Invocation-level selection flags now build a complete Preferred Selection override while preserving each resolved category's Fallback Chain and emitting cross-category warning metadata.

Evidence by acceptance criterion:

- Mixed backend/frontend plus QA: `TestProfileOperationalPreflightMixedTaskGraphWithQADeduplicatesStableOrder` proves backend, frontend, and QA tuples in stable order and deduplicates shared selections.
- Resolve/watch/fetch routing: `TestReviewProfilePreflightResolveAndWatchUseOnlyReviewProfile` proves only the `review` profile for resolve/watch, and `TestReviewProfilePreflightFetchCreatesNoAgentSession` verifies fetch performs no Agent proof or invocation.
- Zero side effects on proof failure: `TestImplementProfilePreflightFailureCreatesNoRunWorktreeOrAgentPrompt` and `TestBranchIntegrityBypassAuditFollowsProfileProof` verify failed proof blocks Run rows, branch/worktree creation, prompts, audits, and Agent work.
- Invocation overrides and fallback preservation: `TestInvocationProfileOverrideAppliesAcrossCategoriesPreservesFallbacksAndWarns` verifies the override applies to every relevant Preferred Selection, preserves configured fallbacks, and records the warning.
- Duplicate tuple attribution: `TestProfileOperationalPreflightMixedTaskGraphWithQADeduplicatesStableOrder` verifies exact duplicate tuples are probed once while retaining every category reference.
- No dynamic selection: updated review and implement failure tests verify profile proof does not read benchmark rank, probe dynamic fallback candidates, or propose unconfigured fallback reruns.

Verification run:

- `GOCACHE=/tmp/roundfix-gocache rtk go test ./internal/cli -run 'Test(ProfileOperationalPreflight|ImplementProfilePreflight|ReviewProfilePreflight|InvocationProfileOverride)' -count=1` — passed, 7 tests.
- `GOCACHE=/tmp/roundfix-gocache rtk go test ./internal/spec ./internal/config ./internal/agent ./internal/cli -run 'Test(TaskType|ProfileResolver|ProfileProof|ProfileOperationalPreflight)' -count=1` — passed, 14 tests across 4 packages.
- `GOCACHE=/tmp/roundfix-gocache rtk go test ./internal/cli -count=1` — passed, 530 tests.
- `GOCACHE=/tmp/roundfix-gocache rtk make verify` — passed, including `go test ./...`, skill checks, and build.
