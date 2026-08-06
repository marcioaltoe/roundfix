---
task: task_04
spec: 0059-run-storage-compaction-and-global-sanitation
status: pending
type: backend
complexity: medium
---

# Task 04: Give every durable table a stated lifecycle

## Overview

Tables added after the retention contract — Run summaries, Agent Selection
records — have no defined long-term lifecycle. They grow, and nothing states
who owns them or when they may go.

This slice writes the policy and makes it checkable, so a table added tomorrow
cannot quietly become permanent.

## Requirements

1. MUST state, for every durable table, its owner and its retention rule.
2. MUST preserve the compact Run index by default.
3. MUST keep Active Run locks governed by the Run lifecycle, not by retention.
4. MUST prune Agent Selection records only with their owning Run, or under an
   explicit evidence-retention rule the policy states.
5. MUST uphold the Spec 0014 promise never to delete `runs` rows or Active Run
   locks, unless this policy explicitly bounds a table with measured
   justification recorded alongside it.
6. MUST fail when a durable table exists with no stated policy, so the next
   table added cannot skip the decision.
7. MUST NOT change what retention deletes today.

## Subtasks

- [ ] Write the per-table owner and rule.
- [ ] Add the check that every durable table appears in the policy.
- [ ] Assert retention behaviour is unchanged.

## Acceptance Criteria

- [ ] Every durable table appears in the documented lifecycle policy with an
      owner and a rule.
- [ ] A durable table with no stated policy fails the check, asserted with a
      fixture table.
- [ ] The compact Run index is preserved by default.
- [ ] Agent Selection records prune only with their owning Run.
- [ ] `runs` rows and Active Run locks are not deleted, asserted.
- [ ] Retention behaviour is unchanged, asserted over the existing tests.

## Context

- interface: `internal/store/journal.go`
- interface: `internal/store/store.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/store -count=1 -run 'Lifecycle|Policy|Retention' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the lifecycle tests ran and passed.
- `go test ./internal/store -count=1` — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.

## References

- `_prd.md` → Core Feature 3; User Story 4; Success Metric 4.
- `_techspec.md` → Build Order 4.
- ADR-0033.
