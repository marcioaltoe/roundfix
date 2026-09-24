---
task: task_07
spec: 0158-daemon-verification-and-access-readiness
status: pending
type: backend
complexity: low
---

# Task 07: Keep every failure and name a degraded access

## Overview

Corrective Task from the pre-PR review of 2026-09-24. In independent mode a temporary failure or unknown cause discards the deterministic failures already collected, so they never reach the repair turn or the reason. And when the Codex full-access sandbox preset is unavailable, readiness records an unqualified `full-access` effective policy.

## Requirements

1. MUST keep the failures already collected in independent mode when a later command ends in a temporary failure or unknown cause, and report all of them in the outcome and its reason.
2. MUST record a degraded effective policy, with the unavailable sandbox preset named, when full access is granted without it, and surface it in the readiness output.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A deterministic failure beside a later temporary failure is kept and reported.
- [ ] A granted mode without its sandbox preset is recorded as degraded.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/engine.go`
- interface: `internal/agent/acpx_runner.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestIndependentVerificationKeepsCollectedFailuresBesideATemporaryOne|TestProofRecordsADegradedFullAccess|TestIndependentVerificationHandsEveryFailureToRepair)$" ./internal/daemon ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestIndependentVerificationKeepsCollectedFailuresBesideATemporaryOne TestProofRecordsADegradedFullAccess TestIndependentVerificationHandsEveryFailureToRepair; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task two of the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Independent Verification
