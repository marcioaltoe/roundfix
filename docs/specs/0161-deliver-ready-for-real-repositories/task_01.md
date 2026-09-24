---
task: task_01
spec: 0161-deliver-ready-for-real-repositories
status: pending
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
