---
task: task_02
spec: 0155-a-verify-that-runs-what-changed
status: completed
type: test
complexity: medium
---

# Task 02: A partition that covers every test once

## Overview

A selector that drops a test from both sets makes verification faster by making
it blind. This Task adds the contract that makes that impossible to do quietly.

## Requirements

1. MUST prove every package listed by `go list ./...` is selected by exactly one
   set.
2. MUST prove every test function in `internal/cli` is selected by exactly one of
   the two `internal/cli` invocations.
3. MUST include negative controls showing the contract fails when an item is
   dropped from both sets and when it is added to both.
4. MUST read the set definitions from `internal/verifyselect`, the same source
   the Makefile uses.

## Subtasks

- [ ] Add the package partition check.
- [ ] Add the `internal/cli` test partition check.
- [ ] Add the negative controls.

## Acceptance Criteria

- [ ] The contract passes on the tree as delivered.
- [ ] A package removed from both sets fails the contract, naming it.
- [ ] An `internal/cli` test selected by both invocations fails the contract,
      naming it.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `Makefile`

## Verification

- `out="$(go test -count=1 -v -run "^TestPartitionCoversEveryTestExactlyOnce" ./internal/verifyselect 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestPartitionCoversEveryTestExactlyOnce"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The two sets

## Result

### Implementation

- Added one whole-tree partition contract that discovers the package inventory
  through `go list ./...` and reads the core and Baseline memberships through
  `verifyselect.Packages`.
- Added the `internal/cli` partition to the same contract. It discovers the
  runnable test functions through `go test -list` and applies the anchored
  Baseline pattern from `verifyselect.BaselineCLITestPattern`; the core
  invocation is the pattern's complement.
- Added test-only negative controls that remove a discovered package from both
  sets and place a discovered `internal/cli` test in both invocations. Each
  control requires the partition diagnostic to contain the affected item name.

### Focused checks

- Before implementation, `rtk rg -n
  "^func TestPartitionCoversEveryTestExactlyOnce" internal/verifyselect`
  exited 1 with no match, confirming that the required contract was absent.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1
  ./internal/verifyselect` — passed after the final code edit in 3.034 seconds.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go vet
  ./internal/verifyselect` — passed.
- `rtk make verify-incremental` — the initial sandboxed attempt was blocked
  when an existing check reached `cafe.github.com`; permitted reruns passed,
  including the post-edit run of format, vet, the full Go test suite, skill
  checks, and the Roundfix build.
- The Task's declared `## Verification` command was not run; the Daemon owns
  it.

### Acceptance evidence

1. `tree satisfies the partition` compares every discovered package and every
   discovered `internal/cli` test with the definitions exported by
   `internal/verifyselect`; the focused package test and incremental gate
   passed.
2. `package omitted from both sets is named` deletes one discovered package
   from copies of both memberships, requires a non-nil partition error, and
   requires that error to name the omitted package; the focused package test
   passed.
3. `CLI test selected by both invocations is named` inserts one discovered
   test into copies of both invocation memberships, requires a non-nil
   partition error, and requires that error to name the overlapping test; the
   focused package test passed.

### Follow-ups

- Task 03 owns the Makefile invocations and Daemon verification wiring; neither
  is included in this diff.

## Carry-forward provenance

- Source Run: `run_20260924T120111Z_5495c6532dc4de3f`
- Source commit: `560aa46c1ab017cd899374ab8f278053c5218204`
