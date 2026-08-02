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
- `go test ./internal/baseline -run '^TestClauseDeltaRendersBeforeLedger$' -count=1 -v | grep -q -- "--- PASS: TestClauseDeltaRendersBeforeLedger"`
  — expected: exit 0; ordering, dispositions, and counts hold.
- `go test ./internal/baseline -run '^TestSameIdentityDriftRequiresRetention$' -count=1 -v | grep -q -- "--- PASS: TestSameIdentityDriftRequiresRetention"`
  — expected: exit 0; the gate from task 08 still holds.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 1; Core Features 3; User Experience.
- `_techspec.md` → Build Order 7.

## Result

Implemented a baseline-owned consolidated-review projection that consumes the
accounted `ClauseDelta` directly. It renders the clause total, all seven
disposition counts in canonical order, and every previous clause in lexical ID
order before appending the caller's already-rendered file ledger byte-for-byte.
It does not infer a disposition from retention evidence, files, or any other
plan field.

The existing retention gate remains authoritative for apply availability. The
new negative-path coverage builds a same-identity drift plan with a provably
unaccounted clause, asserts that planning returns no `Plan`, and renders that
clause's stable identity with the `unaccounted` disposition.

Verification Feedback diagnosis:

- The attempt-1 diagnostic artifact was empty, and repository search found no
  `TestClauseDeltaRendersBeforeLedger` or clause-delta renderer. The failed
  pipeline therefore came from the missing Task 09 implementation rather than
  a test assertion or environment failure.

Focused checks:

- Before implementation, `rtk proxy env
  GOCACHE=/private/tmp/roundfix-task09-gocache rtk go test
  ./internal/baseline -run
  'TestClauseDeltaRendersBeforeLedger/accounted_dispositions_precede_unchanged_file_ledger'
  -count=1` failed to compile because `RenderClauseDeltaBeforeLedger` was
  undefined, establishing the Task's red signal.
- `rtk gofmt -w internal/baseline/plan.go internal/baseline/plan_test.go`
  exited 0.
- The same focused accounted-dispositions subtest passed after implementation
  (`Go test: 2 passed in 1 packages`).
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-gocache rtk go test
  ./internal/baseline -run
  'TestClauseDeltaRendersBeforeLedger/unaccounted_clause_withholds_apply'
  -count=1` passed (`Go test: 2 passed in 1 packages`).
- After the final code edit, `rtk proxy env
  GOCACHE=/private/tmp/roundfix-task09-gocache rtk go test
  ./internal/baseline -run
  'TestClauseDeltaRendersBeforeLedger/(accounted|unaccounted)' -count=1`
  passed (`Go test: 3 passed in 1 packages`).
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-gocache rtk go test
  ./internal/baseline -run
  'Test(ReadyPlanNeverCarriesEmptyLedger|PlanDocumentStrictCodecs)$' -count=1`
  passed (`Go test: 2 passed in 1 packages`).
- `rtk git diff --check` passed.

Acceptance evidence:

- `TestClauseDeltaRendersBeforeLedger` supplies one clause for every supported
  disposition and requires every ID/disposition pair exactly once in the
  rendered review.
- The test requires one displayed count for each of the seven dispositions and
  a displayed total of seven clauses, proving the counts sum to the total.
- The negative subtest requires the exact unaccounted clause ID in the rendered
  delta, not only its count.
- The test requires the clause-delta heading to precede the file ledger and
  requires the complete supplied ledger to remain an exact output suffix.
- The negative subtest exercises the existing planning gate and requires
  `outcome.Plan == nil` when the delta contains an unaccounted clause.
- `rtk git -c core.fsmonitor=false status --porcelain` showed only this Task
  file and files under `internal/baseline/`.

The Task's declared `## Verification` commands were not rerun; the Daemon owns
the remaining Verification attempt and Task settlement.
