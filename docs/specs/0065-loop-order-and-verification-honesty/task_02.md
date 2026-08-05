---
task: task_02
spec: 0065-loop-order-and-verification-honesty
status: pending
type: backend
complexity: high
---

# Task 02: Refuse a Verification that cannot fail

## Overview

Spec 0060's `task_03` existed to prove an instruction-level gate fires. Its
Verification was `make verify` plus a clean `git status` — both of which pass
most easily when no work happened. It settled `completed` having run none of
its four cases.

This slice makes that shape refusable at authoring time, as a mechanical
`SC-VERIFY-WORK-INDEPENDENT` rule inside `roundfix spec check`, which already
fails `make verify`. A skill instruction would be advice to an Agent free to
ignore it, which is exactly how the defect happened.

## Requirements

1. MUST add `SC-VERIFY-WORK-INDEPENDENT`, reported through the existing
   `spec check` finding surface with its file, line, and fix line.
2. MUST decide the property from the Task's **declared** Verification commands,
   never from prose or intent, keeping ADR-0093's detection boundary.
3. MUST refuse a Verification composed only of repository-wide gates and
   working-tree cleanliness checks, because that sequence passes most easily
   when nothing was done.
4. MUST accept a Verification that contains a repository-wide gate **and** at
   least one command asserting the Task's own effect. The rule targets the
   composition, not the presence of any particular command.
5. MUST leave every active and archived Spec in the corpus checking exactly as
   it does today, asserted rather than assumed.
6. MUST keep `TestCheckCorpusBudget` passing, so the sweep stays within budget.
7. MUST NOT change the loop order statements or add the divergence rule; those
   are task_01 and task_04.

## Subtasks

- [ ] Add the rule and its finding.
- [ ] Replay Spec 0060's `task_03` as a fixture and assert refusal.
- [ ] Build the false-positive table over legitimate Verifications.
- [ ] Assert corpus non-regression and the budget test.

## Acceptance Criteria

- [ ] A Verification of only `make verify` plus a clean-tree check is refused.
- [ ] Spec 0060's `task_03`, replayed as written, is refused by this rule.
- [ ] A Verification with a repository gate plus an effect-asserting command
      passes.
- [ ] A Verification of only effect-asserting commands passes.
- [ ] Every Spec in the existing corpus checks as it does today.
- [ ] `TestCheckCorpusBudget` passes.

## Context

- interface: `internal/speccheck/citations.go`
- instruction: `docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'WorkIndependent|Verify' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the rule's tests ran and passed.
- `go test -count=1 -parallel=1 ./internal/speccheck -run '^TestCheckCorpusBudget$'`
  — expected: exit 0; the sweep stays within budget.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Feature 3; Success Metric 2.
- `_techspec.md` → Interfaces; Build Order 2; Risks & Considerations.
- ADR-0093.
