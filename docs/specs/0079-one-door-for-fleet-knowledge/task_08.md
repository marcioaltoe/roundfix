---
task: task_08
spec: 0079-one-door-for-fleet-knowledge
status: pending
type: test
complexity: low
---

# Task 08: Re-record the maintained Source Baseline expectation

## Overview

Corrective slice opened by the QA gate's finding F-001, which blocked seven
rows behind one red repository gate. The clauses this Spec authored added
Source Baseline entries in task_02 and task_06; the maintained
compatibility fixture's hand-declared expectation still names the pre-Spec
count, so `TestReadoptionCompatibilityMaintainedFixture` fails while the
invariants it protects — identity agreeing with entries, and the accounting
shape — hold.

The expectation exists precisely so a legitimate corpus change moves one
declared value instead of hunting literals. This task moves it through the
sanctioned re-recording workflow and proves the package green.

## Requirements

1. MUST confirm the corpus change is legitimate before moving anything: the
   added entries trace to the clauses task_02 and task_06 authored under the
   recorded authorization, and the digest chain is converged.
2. MUST run the sanctioned re-recording workflow rather than hand-editing
   any generated artifact; only the named expectation is authored by hand.
3. MUST move the maintained Source Baseline expectation to the value the
   regenerated fixture actually carries, keeping the identity-agrees-with-
   entries invariant and the accounting expectation intact as separate
   assertions.
4. MUST NOT weaken the assertion — no tolerance, no computing the expected
   value from the fixture under test, no skip. The expectation stays a
   declared constant a reviewer can read.
5. MUST leave every other package untouched; this is a one-expectation
   correction, not a refactor.

## Subtasks

- [ ] Confirm provenance of the added entries and chain convergence.
- [ ] Re-record through the sanctioned workflow; move the expectation.
- [ ] Prove the package and the repository gate green.

## Acceptance Criteria

- [ ] `TestReadoptionCompatibilityMaintainedFixture` passes with the moved
      expectation.
- [ ] The identity-agrees-with-entries and accounting assertions remain
      present and unweakened.
- [ ] The whole `internal/baseline` package passes.
- [ ] No file outside `internal/baseline` and this task file changed.

## Context

- interface: internal/baseline/preservation_test.go

## Verification

- `output="$(go test -count=1 ./internal/baseline -run '^TestReadoptionCompatibilityMaintainedFixture$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the named fixture test is selected and passes — the
  constant now agrees with the regenerated corpus.
- `go test -count=1 ./internal/baseline/...`
  — expected: exit 0; the package the correction lives in is green.
- `grep -q 'sourceBaseline.Identity.EntryCount != len(sourceBaseline.Entries)' internal/baseline/preservation_test.go && grep -q 'maintainedSourceBaselineAccounting' internal/baseline/preservation_test.go`
  — expected: exit 0; the invariant and accounting assertions survived the
  correction rather than being relaxed.
- `output="$(git status --porcelain | awk '{print $NF}' | grep -vE '^(internal/baseline/|docs/specs/0079-one-door-for-fleet-knowledge/task_08\.md$)')"; [ -z "$output" ]`
  — expected: exit 0; the correction stayed inside its package.

## References

- `qa/qa-report-2026-08-06.md` → Findings F-001 and its blocked rows.
- `_techspec.md` → Testing Approach (module choreography); Build Order 1
  and 5.
- ADR-0081.
