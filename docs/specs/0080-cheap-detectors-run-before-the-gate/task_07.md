---
task: task_07
spec: 0080-cheap-detectors-run-before-the-gate
status: pending
type: test
complexity: low
---

# Task 07: Re-record the expectation task_06 invalidates

## Overview

The consequent fix, landing as its own commit after its cause. task_06's new
clauses grow the Source Baseline, and the maintained compatibility fixture
declares its expected shape by hand so a legitimate corpus change moves one
declared value instead of hunting literals.

This node exists because folding the correction into task_06 would fail the
tooling-authority gate: `docs/agents/specific-repository.md` requires a
consequent fix — one that only became necessary because the authorized change
made something stale — to be its own commit landing after the Task commit.
Spec 0079 paid a full gate round to learn that, and this graph spends a node
instead.

## Requirements

1. MUST confirm the corpus change is legitimate before moving anything: the
   added entries trace to the clauses task_06 authored under the recorded
   authorization, and the digest chain is converged.
2. MUST run the sanctioned re-recording workflow rather than hand-editing any
   generated artifact; only the declared expectation is authored by hand.
3. MUST move the maintained Source Baseline expectation to the value the
   regenerated fixture carries, keeping the identity-agrees-with-entries
   invariant and the accounting expectation intact as separate assertions.
4. MUST NOT weaken the assertion: no tolerance, no computing the expected
   value from the fixture under test, no skip. The expectation stays a
   declared constant a reviewer can read.
5. MUST leave every other package untouched; this is a one-expectation
   correction, not a refactor.

## Subtasks

- [ ] Confirm provenance and chain convergence.
- [ ] Re-record through the sanctioned workflow; move the expectation.
- [ ] Prove the package and the repository gate green.

## Acceptance Criteria

- [ ] The maintained-fixture test passes with the moved expectation.
- [ ] The invariant and accounting assertions remain present and unweakened.
- [ ] The whole `internal/baseline` package passes.
- [ ] No file outside `internal/baseline` and this task file changed.

## Context

- interface: internal/baseline/preservation_test.go

## Verification

- `output="$(go test -count=1 ./internal/baseline -run '^TestReadoptionCompatibilityMaintainedFixture$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the named fixture test is selected and passes.
- `go test -count=1 ./internal/baseline/...`
  — expected: exit 0; the package the correction lives in is green.
- `grep -q 'sourceBaseline.Identity.EntryCount != len(sourceBaseline.Entries)' internal/baseline/preservation_test.go && grep -q 'maintainedSourceBaselineAccounting' internal/baseline/preservation_test.go && ! grep -q 'maintainedSourceBaselineEntries    = 132' internal/baseline/preservation_test.go` — expected: exits 0. The first two halves prove the invariant and the named
constant survive as separate assertions; the third proves the constant actually
moved off `132`, the entry count Task 06's clauses invalidated. The Agent proved
on 2026-08-11 that `maintainedSourceBaselineAccounting` legitimately stays at
`51` — the sanctioned regeneration reports 134 entries and 51 accounting rows —
and refused to make an assertion true that the regenerated artifacts say is
false. This command named the wrong constant; the entry count is the one that
moves. Without it the command
asserts only that two strings this file already contains are still there, and
the pre-work probe refused exactly that on 2026-08-11.
  — expected: exit 0; the invariant and accounting assertions survived rather
  than being relaxed.
  — expected: exit 0; the correction stayed inside its package.

## References

- `_techspec.md` → Build Order 7.
- ADR-0081.
- `docs/findings/_archived/` → the 2026-08-06 rollup covering gate and
  verification evidence, whose members record what folding this fix costs.
