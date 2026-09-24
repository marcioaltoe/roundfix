---
task: task_01
spec: 0161-deliver-ready-for-real-repositories
status: completed
type: backend
complexity: low
---

# Task 01: Wait for checks that have not reported

## Overview

Right after `gh pr create`, `gh pr checks` exits 1 with `no checks reported`, which today becomes an error and parks the item as `delivery-error`.

## Requirements

1. MUST map exit 1 with empty stdout and a `no checks reported` stderr from `gh pr checks` to an empty report.
2. MUST treat an empty report as pending within the bounded wait.
3. MUST retry a check read error inside the wait until the deadline, parking only on failure, cancellation or timeout.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The exact gh output becomes an empty report.
- [ ] The engine waits through an empty report and merges once checks pass.
- [ ] A read error inside the wait is retried.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/delivery/github.go`
- interface: `internal/delivery/engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestNoChecksReportedIsAnEmptyReport|TestDeliveryWaitsThroughUnreportedChecks|TestCheckReadErrorsAreRetriedUntilTheDeadline)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestNoChecksReportedIsAnEmptyReport TestDeliveryWaitsThroughUnreportedChecks TestCheckReadErrorsAreRetriedUntilTheDeadline; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Checks

## Result

Implemented the bounded check wait without changing the Task status or running
the Daemon-owned Verification command:

- `CurrentHeadChecks` maps gh exit 1 with empty stdout and the exact
  `no checks reported on the '<branch>' branch` stderr shape to an empty check
  report, while preserving other command failures.
- The Delivery Engine already treated an empty report as pending; its regression
  test now carries the acceptance-contract name and proves the item merges after
  checks report success.
- A check read error now stays inside the existing bounded wait. A later
  successful report advances to merge, repeated errors park as `checks-timeout`,
  and context cancellation still returns through the existing cancellation path.

Focused-check evidence:

- `rtk env GOCACHE=/private/tmp/roundfix-task-0161-01-gocache go test -count=1 -run '^(TestNoChecksReportedIsAnEmptyReport|TestCheckCommandFailuresRemainErrors|TestDeliveryWaitsThroughUnreportedChecks|TestCheckReadErrorsAreRetriedUntilTheDeadline|TestPersistentCheckReadErrorsParkAtTheDeadline)$' ./internal/delivery` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task-0161-01-gocache go test -count=1 ./internal/delivery` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task-0161-01-gocache go vet ./internal/delivery` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task-0161-01-gocache make verify-incremental` — blocked when an existing integration check attempted to reach `api.github.com`; the sandbox denied network access and rejected escalation. This is not the Task's declared Verification.

Acceptance evidence:

1. `TestNoChecksReportedIsAnEmptyReport` passes against the exact gh 2.101.0
   stderr, and `TestCheckCommandFailuresRemainErrors` proves the mapping stays
   narrow.
2. `TestDeliveryWaitsThroughUnreportedChecks` passes and observes one bounded
   wait before the item merges after a passing report.
3. `TestCheckReadErrorsAreRetriedUntilTheDeadline` passes after a transient
   read error, and `TestPersistentCheckReadErrorsParkAtTheDeadline` proves
   repeated errors stop at the configured timeout.
