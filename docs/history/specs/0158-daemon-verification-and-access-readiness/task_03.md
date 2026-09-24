---
task: task_03
spec: 0158-daemon-verification-and-access-readiness
status: completed
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

## Result

Runtime specifications now distinguish a requested access policy from the mode
that can honour it. Codex and Claude report their supported full-access modes;
OpenCode and custom commands report the policy unavailable. The runtime-default
policy remains valid without a mode, so readiness does no extra access work
unless full access was explicitly requested.

The disposable Exact Agent Selection Proof now applies the requested mode after
model and reasoning assignment and before session cleanup. A successful proof
records `full-access` as its effective access policy, including in the profile
proof report. A missing mode fails the `runtime_has_access_mode` predicate; an
adapter refusal fails `adapter_accepts_access_mode`. Both errors name the exact
selection through profile preflight and prescribe disabling full access or
selecting a supporting runtime. Because operational preflight still runs before
the Run Database and Run Worktree boundaries, either refusal creates neither.

Acceptance evidence:

- A runtime without a mode and an adapter refusal: `TestProfilePreflightRefusesUnhonouredAccessPolicy` passed. Its OpenCode case failed before invoking the prover; its Codex case surfaced the adapter refusal. Both named the failed predicate and remedy and asserted that no Run Database, Run Worktree or Agent work was created. `TestDisposableProofAppliesRequestedAccessMode` also exercised an actual disposable acpx command sequence whose `set-mode` refusal returned `access_policy_rejected` and still closed the session.
- Supported mode application and recording: `TestDisposableProofAppliesRequestedAccessMode` passed after observing one `set-mode bypassPermissions` call on the disposable session and `EffectiveAccessPolicy == full-access`. The supported profile-preflight case proved every preferred and fallback report for the backend category and recorded `full-access` on each report.

Focused-check evidence:

- Pre-change: `rtk rg -n "TestRuntimeReportsAccessPolicySupport|TestDisposableProofAppliesRequestedAccessMode|TestProfilePreflightRefusesUnhonouredAccessPolicy" internal/agent internal/cli --glob '*_test.go'` exited 1 because none of the required cases existed. After adding the tests first, the focused agent test failed to compile because the access-policy fields, method and proof record did not exist.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -run '^(TestRuntimeReportsAccessPolicySupport|TestDisposableProofAppliesRequestedAccessMode)$' ./internal/agent` passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -run '^TestProfilePreflightRefusesUnhonouredAccessPolicy$' ./internal/cli` passed outside the restricted network sandbox required by the CLI freshness-check harness.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 ./internal/agent` passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache make verify-incremental` passed after the final implementation and test edits, including formatting, vet, the repository Go suite, skill checks and build.

The Daemon-owned `## Verification` command was not run in this Agent turn.
