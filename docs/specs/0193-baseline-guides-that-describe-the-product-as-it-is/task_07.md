---
task: task_07
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: pending
type: test
complexity: low
---

# Task 07: The fleet sweep expects the boundary paragraph once

## Overview

The second delivery Run's QA gate failed on `TestBaselineUpdateFleetSweep/structural-clauses-missing` in `internal/cli/baseline_update_test.go`. Its fixture `fleetStructuralClauses` still pairs `clause.backend.boundary-contracts` with `rule.backend.boundary-contracts` on the same line of `docs/agents/backend.md`, so the builder expects that line twice and removes 14 structural clauses in total. task_01 made the line render once, which is this Spec's intended result. This Task moves the fixture to that fact.

## Requirements

1. MUST change the `fleetStructuralClauses` entry whose `ids` are `clause.backend.boundary-contracts` and `rule.backend.boundary-contracts` so that its `ids` hold only `clause.backend.boundary-contracts`, and change nothing else in that entry or in any other entry.
2. MUST change the removed-clause total that `newFleetStructuralClauseRepository` asserts from 14 to 13.
3. MUST keep every other expectation of `TestBaselineUpdateFleetSweep` unchanged: the plan still restores every carrier, and each restored line still occurs once per listed ID.
4. MUST NOT edit production code, Baseline assets, goldens or any governed file.

## Subtasks

- [ ] Move the fixture entry and the total.
- [ ] Prove the fleet sweep passes, and that restoring the second ID makes it fail.

## Acceptance Criteria

- [ ] `TestBaselineUpdateFleetSweep` passes with all its subtests.
- [ ] `git diff` of this Task touches only `internal/cli/baseline_update_test.go`.

## Context

- interface: `internal/cli/baseline_update_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestBaselineUpdateFleetSweep$' ./internal/cli 2>&1)" || { printf "%s\n" "$out" | grep -E -- "--- FAIL|_test.go:"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestBaselineUpdateFleetSweep/structural-clauses-missing" && ! grep -q '"clause.backend.boundary-contracts", "rule.backend.boundary-contracts"' internal/cli/baseline_update_test.go` — expected: exit 0; before this Task the subtest fails with "occur 1 times in docs/agents/backend.md, want 2".

## References

- task_01, task_06
- QA Report of Run `run_20260930T184402Z_6aa3e4970df8def4` (verification log `batch-002-attempt-1.log`)
