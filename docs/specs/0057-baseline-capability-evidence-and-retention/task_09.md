---
task: task_09
spec: 0057-baseline-capability-evidence-and-retention
status: completed
type: backend
complexity: medium
---

# Task 09: Render the clause-level delta before apply

## Overview

Retention is now accounted, but a maintainer confirming an update still reads a
file ledger and cannot see which rules survived. This Task renders the
clause-level semantic delta before final confirmation, so the decision is made
against meaning rather than against bytes.

## Requirements

1. MUST render, before final confirmation, every previous clause with its
   disposition and a count per disposition.
2. MUST place the clause-level delta ahead of the file ledger, which remains
   for machine review.
3. MUST make an unaccounted clause visible in the delta with its identity, not
   only in a count.
4. MUST NOT offer apply while any clause is unaccounted, matching the gate.
5. MUST render from the accounted dispositions rather than re-deriving them, so
   the delta and the gate cannot disagree.
6. MUST leave the file ledger's content and format unchanged.

## Subtasks

- [ ] Render each clause with its disposition.
- [ ] Render counts per disposition.
- [ ] Place the delta ahead of the file ledger.
- [ ] Name unaccounted clauses individually.

## Acceptance Criteria

- [ ] The consolidated review shows every previous clause with its disposition.
- [ ] Counts per disposition are shown and sum to the clause total.
- [ ] An unaccounted clause is named individually, not only counted.
- [ ] The clause delta appears before the file ledger, and the ledger is
      unchanged.
- [ ] Apply is not offered when the delta contains an unaccounted clause.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/plan.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestClauseDeltaRendersBeforeLedger -count=1`
  — expected: exit 0; ordering, dispositions, and counts hold.
- `go test ./internal/baseline -run TestSameIdentityDriftRequiresRetention -count=1`
  — expected: exit 0; the gate from task 08 still holds.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 1; Core Features 3; User Experience.
- `_techspec.md` → Build Order 7.
