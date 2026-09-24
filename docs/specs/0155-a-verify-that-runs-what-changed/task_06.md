---
task: task_06
spec: 0155-a-verify-that-runs-what-changed
status: pending
type: backend
complexity: medium
---

# Task 06: Select every set a change can affect

## Overview

Corrective Task from the pre-PR review of 2026-09-24. Three selector gaps let a broken change pass `make verify-changed`: a non-`.go` path that is not documentation, such as a `testdata/` fixture or a `.golden` snapshot, selects no set; staged but uncommitted paths are never listed; and a Baseline-only change skips the core packages that import Baseline code (`internal/speccheck`, `internal/cli`).

## Requirements

1. MUST send a path under a package directory, including its `testdata/`, to the set that owns that package.
2. MUST select no set only for an explicit documentation allowlist (`docs/**` and Markdown outside package and embedded directories); every other unknown path fails safe to both sets.
3. MUST list staged changes as well as committed, unstaged and untracked ones.
4. MUST, when the Baseline set is selected, also run the core packages whose imports or test imports reach a Baseline package, derived from `go list` rather than a fixed list.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A fixture-only change in a core package selects the core set.
- [ ] An unknown non-documentation path selects both sets.
- [ ] A staged-only change is selected.
- [ ] A Baseline-only change also runs the core importers of Baseline packages.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/verifyselect/verifyselect.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestFixtureChangesSelectTheOwningSet|TestUnknownNonDocumentationPathsFailSafe|TestSelectListsStagedPaths|TestBaselineChangesAlsoRunCoreImporters|TestClassifyPath|TestSelectListsCommittedUnstagedAndUntrackedPaths)$" ./internal/verifyselect 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestFixtureChangesSelectTheOwningSet TestUnknownNonDocumentationPathsFailSafe TestSelectListsStagedPaths TestBaselineChangesAlsoRunCoreImporters TestClassifyPath TestSelectListsCommittedUnstagedAndUntrackedPaths; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task four of the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
