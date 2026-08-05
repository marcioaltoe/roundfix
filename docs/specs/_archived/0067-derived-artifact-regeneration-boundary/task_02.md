---
task: task_02
spec: 0067-derived-artifact-regeneration-boundary
status: completed
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

## Result

### Implementation

- Extended the canonical derived-ownership suite with an isolated temporary
  repository fixture. The test resolves every `dedicated` record, perturbs all
  regular artifacts it governs, executes the record's exact shell command, and
  requires every claimed byte to return to its pre-perturbation value while all
  other derived bytes remain unchanged.
- Added negative declared-command fixtures for a nonexistent command, a valid
  command carrying a deliberately wrong Go test flag, and a successful command
  that leaves the perturbed artifacts unchanged. Each case must produce the
  corresponding execution or unchanged-artifact failure.
- Added one perturbation probe for every `sanctioned` and `frozen` record. The
  sanctioned command must restore all sanctioned probes, while every frozen
  probe must remain perturbed; dedicated steps also leave frozen probes
  untouched.
- The exercise helper restores the complete temporary derived tree after every
  command, including command failures, before the next subtest runs.

### Focused checks

- Red signal before implementation:
  `rtk rg -n 'TestDeclaredStepRegenerationAndFrozenBoundaries|runDeclaredRegenerationStep' internal/baseline`
  exited 1 because no declared-step execution test existed.
- After implementation,
  `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/wt67run-dee297f2/run_20260805T050602Z_bfee71697528c7d7/.gocache go test ./internal/baseline -run '^TestDeclaredStepRegenerationAndFrozenBoundaries$' -count=1 -v`
  exited 0. Both real dedicated commands, all three negative command fixtures,
  and the sanctioned/frozen boundary subtest passed.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/wt67run-dee297f2/run_20260805T050602Z_bfee71697528c7d7/.gocache go vet ./internal/baseline`
  exited 0.
- `rtk git diff --name-only -- internal/baseline/testdata internal/baseline/assets`
  exited 0 with no output; the focused run changed no repository artifact.

### Acceptance evidence

1. `dedicated/*` subtests derive their artifact sets from the validated
   ownership resolution, run each exact declared command, and passed only after
   every deliberately perturbed claimed artifact returned byte-for-byte.
2. `failure/command does not exist`, `failure/declared flag is wrong`, and
   `failure/command leaves artifacts unchanged` passed by observing the three
   required failure modes. The wrong-flag case is read back from a temporary
   `_ownership.yml` before execution.
3. The sanctioned-command subtest perturbs one probe for each frozen record and
   passed only because `make baseline-digests` left all three perturbed bytes
   untouched.
4. The same subtest perturbs one probe for every sanctioned record and passed
   only because one `make baseline-digests` invocation restored all five.
5. Every command path defers full fixture restoration and the test asserts the
   clean snapshot after each subtest. The repository derived-tree diff remained
   empty after the focused run.

The Task's declared `## Verification` commands were not run; the Daemon owns
that gate and terminal settlement.
