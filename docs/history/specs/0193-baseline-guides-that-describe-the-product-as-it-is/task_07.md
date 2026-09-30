---
task: task_07
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: completed
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

## Result

The boundary-contract fixture lists only `clause.backend.boundary-contracts`.
The removed-clause assertion and its diagnostic now expect 13. The entry's
path and paragraph, every other fixture entry, and all restoration assertions
remain unchanged. This follows the PRD's single-paragraph contract; production
already renders the intended paragraph once.

- Fleet-sweep criterion: `GOCACHE=/private/tmp/roundfix-task07-gocache rtk proxy go test ./internal/cli -run '^TestBaselineUpdateFleetSweep$' -count=1 -json` exited 0. All four subtests passed: `manifest-predates-managed-regions`, `structural-clauses-missing`, `recorded-profile-does-not-resolve`, and `already-current`.
- Negative evidence: the focused structural-clause subtest failed before the edit with `occur 1 times in docs/agents/backend.md, want 2`. After the edit, a temporary Go overlay restored the second ID and the old total of 14. `GOCACHE=/private/tmp/roundfix-task07-gocache rtk proxy go test -overlay=/private/tmp/roundfix-task07-overlay.json ./internal/cli -run '^TestBaselineUpdateFleetSweep$/^structural-clauses-missing$' -count=1` exited 1 with the same diagnostic. The overlay changed no worktree file.
- Diff-scope criterion: inspection of `git -c core.fsmonitor=false diff -- internal/cli/baseline_update_test.go` shows only the requested ID removal and total adjustment. The only implementation path changed is `internal/cli/baseline_update_test.go`; this assigned Task file also carries the required Result evidence and its pre-existing Daemon status change. No production code, Baseline asset, golden, governed file, manifest, or other Task file changed.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.
- Incremental check: `GOCACHE=/private/tmp/roundfix-task07-gocache rtk make verify-incremental` exited 2 in the sandbox: two process-stop integration tests could not read the process table, and the suite guard detected the Result edit made during that run. The same command was rerun with process access and no concurrent edits; it exited 0, including formatting, vet, package tests, skill checks, and build. Logs: `/private/tmp/roundfix-task07-incremental.log` and `/private/tmp/roundfix-task07-incremental-retry.log`.

The default Go cache was sandbox-blocked; the focused checks used the writable
task-local cache above. Declared Verification was not run, and Task status was
not edited. The Daemon retains Verification and settlement ownership.
