---
task: task_02
spec: 0067-derived-artifact-regeneration-boundary
status: pending
type: test
complexity: medium
---

# Task 02: Prove every declared step actually rewrites its artifacts

## Overview

A manifest can lie. task_01 closes one direction — nothing unowned — and this
slice closes the other: each `dedicated` record's command is executed and
asserted to rewrite exactly what the record claims.

This is the assertion that would have caught a wrong flag name. On 2026-08-05 a
flag was deduced from a test name, did not exist, and the mistake surfaced only
when a gate stayed red.

## Requirements

1. MUST execute each `dedicated` record's declared command in a fixture and
   assert it rewrites the artifacts under that record's directory.
2. MUST fail when a declared command does not exist, exits non-zero, or leaves
   its artifacts unchanged after a deliberate perturbation.
3. MUST restore the fixture afterwards so the test leaves no artifact modified.
4. MUST assert `sanctioned` directories are rewritten by `make baseline-digests`
   and not by any dedicated step.
5. MUST assert `frozen` directories are rewritten by nothing, including the
   sanctioned command.

## Subtasks

- [ ] Execute each dedicated command against a perturbed fixture.
- [ ] Assert the claimed artifacts are rewritten and restored.
- [ ] Assert sanctioned and frozen classes behave as declared.

## Acceptance Criteria

- [ ] Every `dedicated` command runs and rewrites its declared artifacts.
- [ ] A record whose command does not exist fails the test, proven by a fixture
      carrying a deliberately wrong flag.
- [ ] A `frozen` directory is unchanged by the sanctioned command.
- [ ] A `sanctioned` directory is rewritten by the sanctioned command.
- [ ] The repository's artifacts are byte-identical after the test run.

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -count=1 -run 'DeclaredStep|Regeneration|Frozen' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the declared-step tests ran and passed.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `git diff --quiet HEAD -- internal/baseline/testdata internal/baseline/assets`
  — expected: exit 0; the test left no artifact modified.

## References

- `_prd.md` → Core Feature 5; Success Metric 2.
- `_techspec.md` → Testing Approach; Build Order 2.
