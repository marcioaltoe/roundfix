---
task: task_04
spec: 0065-loop-order-and-verification-honesty
status: pending
type: backend
complexity: low
---

# Task 04: Check that the order restatements agree

## Overview

task_01 states the loop order once. This slice makes a second divergence
mechanically impossible, so the next edit to any one statement cannot silently
recreate the contradiction.

It depends on task_01 for a reason that is not cosmetic: this rule fails while
the statements disagree, and `spec check` runs inside `make verify`. Landing it
first would leave the repository red between two Tasks, and the Daemon runs the
configured Verification command as a precondition — so the Task meant to repair
that state would be settled without ever starting. Spec 0075 lost a Run to that
exact shape on 2026-08-05.

## Requirements

1. MUST add `SC-LOOP-ORDER-DIVERGENT`, failing when the loop order differs
   between the places that state it.
2. MUST cover all three sources: the shipped clause, `docs/agents/autonomous-work.md`,
   and the Baseline module asset. A check reading only two recreates the defect
   in a smaller form.
3. MUST report which sources disagree and how, so the fix does not require
   diffing three files by hand.
4. MUST pass on the corrected statements task_01 produced, proving the rule and
   the sources agree at the moment it lands.
5. MUST keep `TestCheckCorpusBudget` passing.
6. MUST NOT restate or edit the order itself; task_01 owns the content.

## Subtasks

- [ ] Add the rule reading all three sources.
- [ ] Fixture a divergence in each source and assert detection.
- [ ] Assert the corrected repository passes.

## Acceptance Criteria

- [ ] A divergence in the shipped clause is detected and named.
- [ ] A divergence in the repository guide is detected and named.
- [ ] A divergence in the Baseline module asset is detected and named.
- [ ] The repository as task_01 left it passes the rule.
- [ ] `TestCheckCorpusBudget` passes.

## Context

- interface: `internal/speccheck/citations.go`
- instruction: `docs/agents/autonomous-work.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'LoopOrder|Divergent' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the divergence tests ran and passed.
- `go run -buildvcs=false ./cmd/roundfix spec check > /dev/null` — expected:
  exit 0; the corrected repository passes its own new rule.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Feature 1; Success Metric 1.
- `_techspec.md` → Build Order 1; Risks & Considerations.
- ADR-0093.
