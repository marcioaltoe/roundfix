---
task: task_01
spec: 0146-a-gate-that-runs-the-analyzer
status: pending
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
- `out="$(go test -count=1 -tags repocontract -run "^TestRepositoryGateRunsTheAnalyzer$" ./... 2>&1)"; printf "%s\n" "$out" | grep -q "no tests to run" && { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q "FAIL"` — expected: exit 0; the control runs and fails against a gate that does not yet compose the analyzer, which is why it is written first.

## References

`_prd.md` → Core Features 2-3; User Story 2; Goals 2-3; Success Metrics 1 and 3;
`_techspec.md` → Implementation Design: The negative control; API Contracts 1
and 3; Build Order 1; Testing Approach 1-2 and 4.
