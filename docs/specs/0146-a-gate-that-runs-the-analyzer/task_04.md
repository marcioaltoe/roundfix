---
task: task_04
spec: 0146-a-gate-that-runs-the-analyzer
status: pending
type: chore
complexity: low
---

# Task 04: Let a gate run the control

## Overview

Independent review found the control unreachable. It carries the
repository-contract build tag, which keeps it out of the ordinary test sweep,
and the target that runs those contracts selects a fixed list of three older
names. So the test written to prove the analyzer step can fail is itself never
run by any gate — which is the same defect it exists to prevent, one level up.
This slice puts it in that list, and is the first of the two corrective Tasks
the contract allows.

## Requirements

1. MUST include the analyzer control in the set of repository contracts the
   pull-request gate runs.
2. MUST keep every other contract in that set running, with its name unchanged.
3. MUST keep the control out of the ordinary test sweep, as the other repository
   contracts are.
4. MUST NOT change any path outside this Spec's bounded governed files, its own
   Task file and ordinary source.

## Subtasks

- [ ] Add the control to the repository-contract selection.
- [ ] Confirm the pull-request gate runs it and that it passes.
- [ ] Confirm the ordinary sweep still does not run it.

## Acceptance Criteria

- [ ] The pull-request gate runs the analyzer control.
- [ ] The three pre-existing contracts still run.
- [ ] `go test ./...` still does not run the control.
- [ ] No path outside the bounded list changes.

## Context

- interface: `Makefile`

## Verification

- `grep -q "TestRepositoryGateRunsTheAnalyzer" Makefile` — expected: exit 0; the gate's contract selection names the control. Before this Task it does not.
- `make repo-test > /tmp/0146-repo-test.txt 2>&1; status=$?; grep -q "TestRepositoryGateRunsTheAnalyzer" /tmp/0146-repo-test.txt || { cat /tmp/0146-repo-test.txt; exit 1; }; exit $status` — expected: exit 0; the gate actually executes the control, named in its output, and passes. Before this Task the control is not selected, so it is absent from the output.
- `grep -q "TestRepositoryGateRunsTheAnalyzer" Makefile && go test -count=1 -run "^TestRepositoryGateRunsTheAnalyzer$" ./internal/baseline > /tmp/0146-sweep.txt 2>&1 && grep -q "no tests to run" /tmp/0146-sweep.txt` — expected: exit 0; the gate names the control while the untagged sweep still does not run it. Before this Task the Makefile does not name it, so the command fails.

## References

`_prd.md` → Core Feature 2; User Story 2; Goal 2; Success Metric 1;
Project Constraints: Tooling authority;
`_techspec.md` → Implementation Design: The negative control; Testing Approach 1;
`_authorization.md`.
