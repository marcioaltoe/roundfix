---
task: task_01
spec: 0155-a-verify-that-runs-what-changed
status: completed
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

## Result

### Implementation

- Added a context-aware selector that combines committed paths since the merge
  base with unstaged and untracked paths, deduplicates them, and classifies the
  result into the core and Baseline sets. Any Git or base-resolution error
  returns both sets with the diagnostic preserved.
- Added shared package-set definitions and an anchored Baseline CLI test pattern
  derived by parsing `internal/cli/baseline_*_test.go` declarations.
- Added the `verify-select` entry point. Its default stdout is one selected set
  per line; `-packages core|baseline` and `-baseline-cli-pattern` expose the
  definitions consumed by later Tasks. Diagnostics stay on stderr.
- Added pure classification cases and real temporary-repository coverage for
  committed, unstaged, untracked, unresolvable-base, and failing-Git behavior.

### Focused checks

- `rtk env GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 ./internal/verifyselect`
  — passed.
- `rtk env GOCACHE=/tmp/roundfix-task01-gocache go vet ./internal/verifyselect ./cmd/verify-select`
  — passed.
- `rtk env GOCACHE=/tmp/roundfix-task01-gocache go build -o /tmp/roundfix-verify-select-task01 ./cmd/verify-select`
  — passed.
- `rtk env GOCACHE=/tmp/roundfix-task01-gocache go run ./cmd/verify-select -base refs/heads/does-not-exist`
  — exited 0, printed `core` and `baseline` to stdout, and reported the
  unresolved base on stderr.
- `rtk make verify-incremental` — passed after rerunning with the permission its
  existing GitHub-dependent check requires; formatting, vet, the full Go test
  suite, skill checks, and the Roundfix build passed.
- The Task's declared `## Verification` commands were not run; the Daemon owns
  them.

### Acceptance evidence

1. `TestSelectClassifiesFixtureChangeSets` covers core-only, Baseline-only,
   mixed, module-file, and documentation-only fixtures; the focused package
   test passed.
2. `TestSelectFailsSafeToBothSets` covers an unresolvable base and a Git
   invocation in a non-repository directory; the focused package test passed,
   and the live entry-point probe also selected both sets.
3. `TestBaselineCLITestPatternIsDerived` compares the derived names with an
   independent declaration scan and checks the pattern against every declared
   `internal/cli` test; the focused package test passed.

### Follow-ups

- Task 02 owns the whole-suite partition contract. Task 03 owns Makefile and
  Daemon wiring; neither is included in this diff.

## Carry-forward provenance

- Source Run: `run_20260924T120111Z_5495c6532dc4de3f`
- Source commit: `2552f2776fb451cee16579b99a0a04a667cad6cb`
