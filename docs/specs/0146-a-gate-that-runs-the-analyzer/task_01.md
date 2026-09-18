---
task: task_01
spec: 0146-a-gate-that-runs-the-analyzer
status: completed
type: test
complexity: medium
---

# Task 01: Prove the analyzer can fail

## Overview

A gate step that cannot fail is decoration. This slice writes the control first,
against a gate that does not yet run the analyzer, so it fails for the reason it
exists before the next Task makes it pass.

## Requirements

1. MUST run the Go toolchain's analyzer over a fixture package carrying one
   known diagnostic, and MUST require a non-zero status that names it.
2. MUST place that fixture where module-wide package matching does not reach it,
   so the deliberate defect never fails the repository's own gate.
3. MUST read the repository Verification's composition and require the analyzer
   target in both the full and the incremental tier.
4. MUST run under the repository-contract build tag the other repository
   contracts use, so it stays out of the ordinary test sweep.
5. MUST NOT add a dependency, an analyzer configuration file, or any suppression
   mechanism.
6. MUST name the test `TestRepositoryGateRunsTheAnalyzer`, so the Task's own
   check and later readers agree.

## Subtasks

- [ ] Write the fixture with one known diagnostic under a testdata directory whose path names the analyzer.
- [ ] Assert the analyzer fails over it, naming the diagnostic.
- [ ] Assert the gate's composition contains the analyzer target in both tiers.

## Acceptance Criteria

- [ ] The control fails on the current tree, because the composition lacks the
      step.
- [ ] The analyzer exits non-zero over the fixture and names its diagnostic.
- [ ] The module-wide analyzer run is unaffected by the fixture.
- [ ] No dependency or configuration file is added.

## Context

- interface: `Makefile`
- interface: `internal/baseline/derived_regeneration_repocontract_test.go`

## Verification

- `grep -rq "func TestRepositoryGateRunsTheAnalyzer" internal/` — expected: exit 0; the control exists. Before this Task it does not.
- `fixture="$(git ls-files "internal/**/testdata/**" | grep "analyzer")"; test -n "$fixture" && { go vet ./... > /tmp/0146-module-vet.txt 2>&1; test ! -s /tmp/0146-module-vet.txt; }` — expected: exit 0; the fixture exists and its deliberate diagnostic stays out of module-wide matching. Before this Task the fixture does not exist, so the command fails.
- `go test -count=1 -tags repocontract -run "^TestRepositoryGateRunsTheAnalyzer$" ./... > /tmp/0146-control.txt 2>&1; grep -q "FAIL: TestRepositoryGateRunsTheAnalyzer" /tmp/0146-control.txt` — expected: exit 0; the control runs and fails against a gate that does not yet compose the analyzer, which is why it is written first. Before this Task the test does not exist, so no such line appears and the command fails.

## References

`_prd.md` → Core Features 2-3; User Story 2; Goals 2-3; Success Metrics 1 and 3;
`_techspec.md` → Implementation Design: The negative control; API Contracts 1
and 3; Build Order 1; Testing Approach 1-2 and 4.

## Result

Implemented the repository-contract control under the existing `repocontract`
build tag. The control runs `go vet` against the isolated
`internal/baseline/analyzer/testdata/printfdiagnostic` fixture, requires a
non-zero analyzer exit naming the `fmt.Printf format %d` diagnostic, then reads
the `Makefile` and requires the `vet` dependency in both `verify` tiers.

Focused evidence:

- `GOCACHE=/tmp/roundfix-0146-go-cache go vet ./internal/baseline/analyzer/testdata/printfdiagnostic`
  exited 1 and named the single deliberate `fmt.Printf format %d` diagnostic.
- `GOCACHE=/tmp/roundfix-0146-go-cache go test -tags repocontract -run '^TestRepositoryGateRunsTheAnalyzer$' ./internal/baseline`
  reached the control and exited 1 only because `verify` and
  `verify-incremental` do not yet depend on `vet`, preserving the intended red
  state for Task 02.
- `GOCACHE=/tmp/roundfix-0146-go-cache go list ./...` exited 0 and did not list
  the fixture package, showing module-wide package matching excludes it.
- `GOCACHE=/tmp/roundfix-0146-go-cache go test -run '^TestRepositoryGateRunsTheAnalyzer$' ./internal/baseline`
  exited 0 with no tests to run, showing the ordinary test sweep excludes the
  repository contract.
- `git diff --check` exited 0. The changed implementation paths are the contract
  test, its fixture, and this Task result; no dependency, analyzer configuration,
  or suppression file was added.

Acceptance evidence:

- The control fails on the current composition and names both missing tier
  dependencies.
- The analyzer exits non-zero over the fixture and names its known diagnostic.
- Module-wide package discovery excludes the fixture beneath `testdata`.
- The implementation uses only the standard library and Go toolchain and adds
  no dependency or configuration.

Verification Feedback repair, attempt 1:

- Inspected the Daemon diagnostic artifact and the redirected analyzer output.
  The fixture-discovery precondition produced no path, and the redirected file
  was empty, proving the analyzer did not run in the failed command.
- Root cause: both new implementation files were present only as untracked
  worktree files, while the authored check discovers the fixture through
  `git ls-files`.
- Staged the contract test and diagnostic fixture as the Task's intended new
  source files. A focused `git ls-files 'internal/**/testdata/**'` inspection
  now returns
  `internal/baseline/analyzer/testdata/printfdiagnostic/diagnostic.go`.
- Repeated the focused fixture analyzer check with the task-local Go cache; it
  still exits non-zero and names only the deliberate `fmt.Printf format %d`
  diagnostic. Repeated package discovery still excludes the fixture.
- The failed declared Verification command was not rerun; the Daemon owns its
  single configured rerun.
