---
task: task_05
spec: 0035-agent-selection-profiles
status: pending
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

- [ ] Resolve Implement and review Work Categories from requested work.
- [ ] Apply complete one-Run Preferred Selection overrides.
- [ ] Produce stable, deduplicated proof plans with category references.
- [ ] Execute sequential disposable proofs and cleanup.
- [ ] Place proof before every Run mutation boundary.
- [ ] Add operational zero-side-effect and override-warning tests.

## Acceptance Criteria

- [ ] A mixed backend/frontend Task Graph plus QA proves all distinct configured tuples once in stable order.
- [ ] Resolve and watch prove only `review`; fetch creates no Agent Session.
- [ ] One failed preferred or fallback proof prevents all Run, branch, worktree, prompt, and durable selection side effects.
- [ ] Invocation overrides affect every relevant Preferred Selection, preserve each Fallback Chain, and emit the cross-category warning.
- [ ] Exact duplicate tuples are probed once while results remain attributable to every category.
- [ ] No benchmark rank, runtime-owned default, or dynamically discovered catalog candidate changes the configured proof plan.

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
