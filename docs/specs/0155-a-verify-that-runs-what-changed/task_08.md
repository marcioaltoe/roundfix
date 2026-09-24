---
task: task_08
spec: 0155-a-verify-that-runs-what-changed
status: completed
type: backend
complexity: low
---

# Task 08: Keep the partition, widen the selection

## Overview

Corrective Task approved by the maintainer on 2026-09-24 as a one-off exception to the ceiling of two. Task 06 made `Packages(BaselineSet)` also return the core packages that import Baseline code, so `cmd/roundfix`, `internal/cli` and `internal/speccheck` belong to both sets and `TestPartitionCoversEveryTestExactlyOnce` fails, which turns `make verify-changed` red on the tree.

## Requirements

1. MUST make `Packages` a partition again: every package belongs to exactly the set that owns its directory.
2. MUST make a selection that includes the Baseline set also include the core set whenever a core package imports or test-imports a Baseline package, derived from `go list` and not from a fixed list.
3. MUST keep a core-only change selecting only the core set.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The partition contract passes.
- [ ] A Baseline-only change selects both sets while core importers of Baseline exist.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/verifyselect/verifyselect.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestPartitionCoversEveryTestExactlyOnce|TestBaselineSelectionAlsoRunsTheCoreSet|TestPackagesAreAPartition)$" ./internal/verifyselect 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestPartitionCoversEveryTestExactlyOnce TestBaselineSelectionAlsoRunsTheCoreSet TestPackagesAreAPartition; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the partition contract fails and the other two cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

- Implementation: `Packages` once again assigns packages only by their directory-owned set. `Select` now widens a Baseline-only selection to both sets when `go list -json ./...` reports that a core package imports or test-imports a Baseline package; a core-only selection bypasses that dependency scan and remains core-only. A failed dependency scan fails safe to both sets and returns the diagnostic.
- Partition evidence: before the implementation, `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test -count=1 -run '^TestPartitionCoversEveryTestExactlyOnce$' ./internal/verifyselect` exited 1 because `./cmd/roundfix` belonged to two sets. After the implementation, the same focused test with `-v` exited 0 and reported `TestPartitionCoversEveryTestExactlyOnce` plus both subtests as passing.
- Selection evidence: `TestBaselineSelectionAlsoRunsTheCoreSet` exercises production imports, external-test imports, and the negative case with no core importer. `TestPackagesAreAPartition` keeps those importers in core while Baseline roots remain in Baseline. Together with the existing core-only staged-change case, `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test -count=1 -run '^(TestPackagesAreAPartition|TestBaselineSelectionAlsoRunsTheCoreSet|TestSelectListsStagedPaths)$' ./internal/verifyselect` exited 0.
- Focused package check: `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test -count=1 ./internal/verifyselect` exited 0.
- Incremental repository check: `rtk make verify-changed` exited 0 after rerunning with network access for tests that query `api.github.com`; formatting, vet, build, selected Go tests, skill tests, and `roundfix skills check` passed.
- The authored `## Verification` command was not run; the Daemon owns it and Task settlement.

## Carry-forward provenance

- Source Run: `run_20260924T192558Z_edaf3e3e31b800e8`
- Source commit: `a747ba4c4db8684d615e1660f80e889c0191adbf`
