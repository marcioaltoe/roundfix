---
task: task_02
spec: 0167-verification-and-readiness-follow-ups
status: completed
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

## Result

### Implementation

- `profiles validate` now appends the effective access policy to each text
  proof line only when that policy identifies itself as degraded.
- Doctor now deduplicates degraded policies across successful profile proofs
  and appends them to the profile readiness line. Runtime-default and ordinary
  `full-access` proofs retain their prior text on both surfaces.

### Focused checks

- Before the production changes,
  `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 -run '^TestProfilesValidateTextNamesADegradedPolicy$' ./internal/cli`
  exited 1 because the degraded proof line printed only `passed`.
- Before the production changes,
  `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 -run '^TestDoctorNamesADegradedPolicy$' ./internal/cli`
  exited 1 because Doctor omitted the degraded policy from its profile
  readiness line.
- After the production changes, both focused commands above exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 -run '^(TestProfilesValidateDeduplicatesProofsAndReportsEveryReference|TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts)$' ./internal/cli`
  exited 0 and preserved the adjacent JSON and ordinary Doctor output
  contracts.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache make verify-incremental`
  exited 0 with network access. The first restricted run was blocked when a
  GitHub integration check could not reach `api.github.com`, so it produced no
  repository verdict.
- The Task's declared `## Verification` command was not run; the Daemon owns
  that check.

### Acceptance evidence

- `TestProfilesValidateTextNamesADegradedPolicy` exercises the text command
  with degraded, ordinary full-access and runtime-default proofs. It verifies
  the degraded policy is named and compares the two non-degraded outputs
  exactly with their prior text.
- `TestDoctorNamesADegradedPolicy` exercises Doctor with the same three policy
  states. It verifies the degraded policy appears once on the aggregate
  profile readiness line and that both non-degraded lines remain unchanged.
