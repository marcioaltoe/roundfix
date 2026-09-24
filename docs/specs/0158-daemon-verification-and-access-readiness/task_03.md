---
task: task_03
spec: 0158-daemon-verification-and-access-readiness
status: pending
type: backend
complexity: medium
---

# Task 03: Prove access policy in readiness

## Overview

Profile readiness proves model and effort on a disposable session but never applies the requested access mode, so an adapter that refuses it fails every Task after the Run starts, and a runtime without the mode runs silently without it.

## Requirements

1. MUST let a runtime report whether it has an access mode for a requested policy.
2. MUST apply the requested access mode on the disposable proof session, for the preferred selection and each fallback readiness already proves, and record the effective policy in the proof.
3. MUST refuse before Run creation when a runtime has no mode for the requested policy or its adapter refuses the mode, naming the selection, the failed predicate and the remedy.
4. MUST leave readiness unchanged when full access is not requested.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A runtime without the mode and an adapter that refuses it both fail readiness, and no Run is created.
- [ ] A supported mode is applied during the proof and its effective policy is recorded.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/cli/profile_preflight.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRuntimeReportsAccessPolicySupport|TestDisposableProofAppliesRequestedAccessMode|TestProfilePreflightRefusesUnhonouredAccessPolicy)$" ./internal/agent ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestRuntimeReportsAccessPolicySupport TestDisposableProofAppliesRequestedAccessMode TestProfilePreflightRefusesUnhonouredAccessPolicy; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the three cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Access policy in readiness
