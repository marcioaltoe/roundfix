---
task: task_02
spec: 0167-verification-and-readiness-follow-ups
status: pending
type: backend
complexity: low
---

# Task 02: Degraded access is printed

## Overview

A degraded full-access policy reaches only `profiles validate --json`; the text output and Doctor's profile readiness line print `passed`.

## Requirements

1. MUST print the effective access policy in the text output of `profiles validate` when it is degraded.
2. MUST name a degraded policy in Doctor's profile readiness line.
3. MUST print nothing extra when the policy is not degraded.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Both text surfaces name a degraded policy and stay unchanged otherwise.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/profiles_validate.go`
- interface: `internal/cli/doctor.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestProfilesValidateTextNamesADegradedPolicy|TestDoctorNamesADegradedPolicy)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProfilesValidateTextNamesADegradedPolicy TestDoctorNamesADegradedPolicy; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Degraded access
