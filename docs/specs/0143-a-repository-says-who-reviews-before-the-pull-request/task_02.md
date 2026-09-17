---
task: task_02
spec: 0143-a-repository-says-who-reviews-before-the-pull-request
status: pending
type: backend
complexity: medium
---

# Task 02: Report the resolved policy and its source

## Overview

The diagnostic command gains one read-only check that states which reviewer runs
before a Pull Request and where that choice came from, so a maintainer can tell
a project decision from an inherited one without running a reviewer.

## Requirements

1. MUST report one check whose detail names the resolved provider and its source
   layer.
2. MUST state explicit `none` as review disabled by configuration.
3. MUST contact nothing: no provider installation, authentication, reviewer
   session, readiness probe or request happens because of this check.
4. MUST leave every existing check's name, order and detail unchanged, and MUST
   keep the command mutating nothing.
5. MUST name the check constant `HealthCheckPrePRReview`, so the Task's own
   check and later readers agree.

## Subtasks

- [ ] Add the check and its constant beside the existing ones.
- [ ] Render provider and source in its detail, with the disabled wording for
      `none`.
- [ ] Cover each source layer and the disabled case at the command's test seam.

## Acceptance Criteria

- [ ] The check reports provider and source for a project choice, a user choice
      and the built-in default.
- [ ] `none` is reported as review disabled by configuration.
- [ ] The command's other checks and their order are unchanged.
- [ ] No provider is invoked by the check.

## Context

- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`

## Verification

- `grep -q "HealthCheckPrePRReview" internal/cli/health.go` — expected: exit 0; the check exists. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestDoctorReportsPrePRReviewPolicy$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 3 and 5; User Story 4; Goals 2 and 4;
`_techspec.md` → Implementation Design: The report; API Contracts 2-3;
Build Order 2.
