---
task: task_08
spec: 0155-a-verify-that-runs-what-changed
status: pending
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
