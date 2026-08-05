---
task: task_03
spec: 0065-loop-order-and-verification-honesty
status: pending
type: backend
complexity: medium
---

# Task 03: Refuse a Task that contradicts itself

## Overview

Spec 0060's `task_03` also carried requirements that contradicted each other,
so the rehearsal could not be performed as written. An Agent Session spent a
turn discovering that, and the Task settled `completed` anyway.

Two rules close it: one for requirements that cannot all hold, one for a
rehearsal Task that never declares what it must exercise.

## Requirements

1. MUST add `SC-REQUIREMENT-CONTRADICTORY`, reporting one requirement
   forbidding a state another requirement needs.
2. MUST decide contradiction from declared MUST and MUST NOT clauses over the
   same named subject, and MUST report nothing when the subject cannot be
   identified — silence is the correct answer to an undecidable pair, per
   ADR-0093.
3. MUST add `SC-REHEARSAL-UNDECLARED`, refusing a Task whose stated purpose is
   proving a gate fires but which declares no cases and no observation for
   them.
4. MUST define the authored section a rehearsal Task uses to declare each case
   and how it is observed, so the rule reads a declaration rather than guessing
   intent.
5. MUST leave every active and archived Spec checking exactly as it does today.
6. MUST keep `TestCheckCorpusBudget` passing.
7. MUST reuse the requirement parsing task_02 introduces rather than adding a
   second parser.

## Subtasks

- [ ] Add both rules and their findings.
- [ ] Replay Spec 0060's `task_03` and assert both refusals.
- [ ] Assert an undecidable pair reports nothing.
- [ ] Assert corpus non-regression and the budget test.

## Acceptance Criteria

- [ ] A MUST and a MUST NOT over the same subject are refused.
- [ ] A pair whose subject cannot be identified reports nothing.
- [ ] Spec 0060's `task_03`, replayed as written, is refused for its
      contradictory requirements.
- [ ] A rehearsal Task with no declared cases is refused.
- [ ] A rehearsal Task declaring its cases and their observation passes.
- [ ] Every Spec in the existing corpus checks as it does today.
- [ ] `TestCheckCorpusBudget` passes.

## Context

- interface: `internal/speccheck/citations.go`
- instruction: `docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'Contradict|Rehearsal' -v | grep -q -- "--- PASS"`
  — expected: exit 0; both rules' tests ran and passed.
- `go test -count=1 -parallel=1 ./internal/speccheck -run '^TestCheckCorpusBudget$'`
  — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 4 and 5; Success Metric 2.
- `_techspec.md` → Interfaces; Build Order 3.
- ADR-0093.
