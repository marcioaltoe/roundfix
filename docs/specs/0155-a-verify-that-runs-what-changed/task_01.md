---
task: task_01
spec: 0155-a-verify-that-runs-what-changed
status: pending
type: backend
complexity: medium
---

# Task 01: The selector and its classification

## Overview

Add `internal/verifyselect` and a thin entry point at `cmd/verify-select` that,
given a base ref, lists the changed paths and says which test sets they can
affect: core, Baseline, both, or neither.

## Requirements

1. MUST list changed paths as those committed since the merge base with the
   base ref, plus unstaged and untracked files.
2. MUST classify `internal/baseline/`, `internal/baselineacp/`, `skills/`,
   `.agents/skills/` and `internal/cli/baseline_*` as Baseline paths.
3. MUST select both sets for `go.mod`, `go.sum` and `Makefile`.
4. MUST select the core set for any other Go source or test.
5. MUST select neither set for documentation outside the embedded directories.
6. MUST select both sets when the base cannot be resolved or Git fails.
7. MUST expose, on request, each set's package list and the test-name pattern
   for the Baseline tests in `internal/cli`, derived from the `func Test...`
   declarations in `internal/cli/baseline_*_test.go` rather than kept by hand.

## Subtasks

- [ ] Add the classification and the change listing.
- [ ] Add the set definitions and the derived test-name pattern.
- [ ] Add the entry point.
- [ ] Add tests for every classification rule and both fail-safe cases.

## Acceptance Criteria

- [ ] Core-only, Baseline-only, both, module-file and documentation-only fixture
      change sets each select the expected sets.
- [ ] An unresolvable base and a failing Git invocation each select both.
- [ ] The derived pattern matches exactly the tests declared in
      `internal/cli/baseline_*_test.go`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `Makefile`
- creates: `internal/verifyselect/verifyselect.go`
- creates: `cmd/verify-select/main.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestSelect" ./internal/verifyselect 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestSelectFailsSafeToBothSets"` — expected: exit 0; before this Task the package does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestBaselineCLITestPatternIsDerived" ./internal/verifyselect 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestBaselineCLITestPatternIsDerived"` — expected: exit 0; before this Task the package does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The selector
