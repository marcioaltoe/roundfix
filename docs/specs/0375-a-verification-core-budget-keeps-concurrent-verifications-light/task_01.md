---
task: task_01
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
status: completed
type: backend
complexity: medium
---

# Task 01: The Verification Core Budget is computed and a Verification command receives it

## Overview

A new `internal/corebudget` package computes the Verification Core Budget
from the logical CPUs, `runs.max_verifying`, `runs.max_active` and the
process environment, describes where it came from, and rewrites an
environment so it carries the budget as `GOMAXPROCS` and
`ROUNDFIX_VERIFICATION_CORES`. The Daemon's `ExecVerifier` applies that
rewrite to every command whose request carries a budget. No caller sets a
budget yet, so no command changes behavior in this Task.

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   `ExecVerifier.Verify` builds the command environment from
   `judge.WithoutKeyVariables(os.Environ())` and the Run Marker only, and no
   non-test file under `internal/` mentions `GOMAXPROCS` or `NumCPU`.
2. MUST create `internal/corebudget/corebudget.go` with exactly the
   constants, types and functions of `_techspec.md` → Interfaces for
   `internal/corebudget`, behaving as API Contract 1 says, including the four
   `Describe` texts word for word. The package imports nothing outside the
   standard library and reads nothing but its inputs.
3. MUST add `Cores int` to `daemon.VerifyRequest` in
   `internal/daemon/daemon.go` and make `ExecVerifier.Verify` set the command
   environment to `corebudget.Environment(<the environment it builds today>, req.Cores)`
   after the Run Marker is applied, as API Contract 2 says. With `Cores` 0
   the environment is byte-identical to today's. No other file changes
   behavior.
4. MUST create `internal/corebudget/corebudget_test.go` with
   `TestTheVerificationCoreBudgetIsHalfTheLogicalCPUsDividedByTheMachineCapacity`,
   `TestAnExplicitGOMAXPROCSWinsTheVerificationCoreBudget` and
   `TestTheCoreBudgetEnvironmentReplacesBothVariables`, table tests with the
   named rows of `_techspec.md` → Testing Approach: rows `default on twelve`
   (12 CPUs, capacity 1 → 6), `two at once` (12, 2 → 3), `odd count rounds
   down` (13, 1 → 6), `capacity off divides by max active` (12, 0, `max_active`
   4 → 1), `both off is undivided` (12, 0, 0 → 6), `one CPU keeps one`,
   `never below one` (4 CPUs, capacity 5 → 1), each also asserting the
   divisor, the source and the `Describe` text; rows `positive whole number`,
   `larger than the share`, `last entry wins`, `zero is ignored`, `negative is
   ignored`, `signed is ignored` (`+4`), `text is ignored`, `empty is ignored`
   and `the Roundfix variable alone is ignored`; and an environment whose two
   `GOMAXPROCS` entries and one `ROUNDFIX_VERIFICATION_CORES` entry become one
   of each, other entries in order, with the input unchanged and a budget of 0
   returning an equal copy.
5. MUST create `internal/daemon/verification_cores_test.go` with
   `TestTheVerificationCommandReceivesTheCoreBudget`: an `ExecVerifier`
   command in a temporary directory writes
   `${GOMAXPROCS-unset}|${ROUNDFIX_VERIFICATION_CORES-unset}|<number of GOMAXPROCS entries in env>`
   to a file; row `budget four` reads `4|4|1`; row `no budget` reads the
   test process's own two values, unchanged.
6. MUST call `t.Parallel()` first in every new top-level test and subtest
   that can run in parallel, use no real-time wait, and set no process
   environment variable.

## Subtasks

- [ ] Characterize the verifier's environment.
- [ ] Add the `corebudget` package.
- [ ] Apply the budget in `ExecVerifier`.
- [ ] Add the package and verifier tests.

## Acceptance Criteria

- [ ] On twelve logical CPUs the budget is 6 at capacity 1, 3 at 2, and never below 1; an explicit positive whole-number `GOMAXPROCS` wins and any other value is ignored.
- [ ] A Verification command with a budget sees exactly one `GOMAXPROCS` and one `ROUNDFIX_VERIFICATION_CORES`, both the budget; without one it sees today's environment.

## Context

- instruction: `docs/adr/0385-each-verification-receives-a-verification-core-budget-in-its-environment.md`
- instruction: `internal/daemon/verification_environment_test.go`
- instruction: `internal/judge/environment.go`
- creates: `internal/corebudget/corebudget.go`
- creates: `internal/corebudget/corebudget_test.go`
- interface: `internal/daemon/daemon.go`
- creates: `internal/daemon/verification_cores_test.go`

## Consumers

- `corebudget.Compute`, `Budget.Describe`, `Budget.Concurrent` and `corebudget.Environment` — consumers: `./internal/corebudget`
- `daemon.VerifyRequest.Cores` and `ExecVerifier` — consumers: `./internal/daemon`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheVerificationCoreBudgetIsHalfTheLogicalCPUsDividedByTheMachineCapacity|TestAnExplicitGOMAXPROCSWinsTheVerificationCoreBudget|TestTheCoreBudgetEnvironmentReplacesBothVariables)$' ./internal/corebudget 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheVerificationCoreBudgetIsHalfTheLogicalCPUsDividedByTheMachineCapacity TestAnExplicitGOMAXPROCSWinsTheVerificationCoreBudget TestTheCoreBudgetEnvironmentReplacesBothVariables; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the package does not exist; afterwards its three tests pass.

- `out="$(go test -count=1 -timeout 30m -v -run '^TestTheVerificationCommandReceivesTheCoreBudget$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheVerificationCommandReceivesTheCoreBudget ' || { echo 'missing pass: TestTheVerificationCommandReceivesTheCoreBudget' >&2; exit 1; }; go test -count=1 -timeout 30m -run 'Verification|Verifier' ./internal/daemon` — expected: exit 0. Before this Task the test does not exist; afterwards it passes, and the existing Verification and verifier tests of the Daemon pass unchanged.

## References

- `_prd.md` → Core Feature 1; Core Feature 2; Core Feature 3; User Story 4; Goal 2; Goal 4; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior 1; Interfaces; API Contract 1; API Contract 2; Invariants 1-3; Testing Approach; Build Order 1
- ADR-0385

## Result

Skills loaded: implement-task, coding-guidelines, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate.

Starting-tree characterization, recorded before implementation:

- `ExecVerifier.Verify` in `internal/daemon/daemon.go` builds `cmd.Env`
  from `judge.WithoutKeyVariables(os.Environ())`, then replaces the Run
  Marker only when `req.RunID` is nonempty. No core budget is applied.
- `rg -n 'GOMAXPROCS|NumCPU' internal --glob '*.go' --glob '!**/*_test.go'`
  returned no matches (exit 1). `internal/corebudget` does not exist.
- Initial `git status --short` showed only this Task file changed; its
  pre-existing change is the Daemon-owned `pending` to `in_progress` status.
  That field is preserved.

Implementation and acceptance evidence:

- Budget computation: added the exact `corebudget` interface from API
  Contract 1, using only standard-library imports and explicit inputs.
  The named table rows assert all Budget fields and the exact descriptions:
  twelve CPUs produce 6 at capacity 1 and 3 at 2; rounding, fallback
  capacities and the minimum of 1 are covered. The override rows cover
  last-entry precedence and every required invalid value, plus ASCII-only
  parsing, the 2147483647 ceiling and leading zeros. Additional rows cover
  CPU normalization and `Concurrent`'s capacity/active-Run fallback.
- Command environment: added `VerifyRequest.Cores` and apply
  `corebudget.Environment` after the Run Marker. The package test verifies
  replacement of duplicate variables, ordering of other entries, unchanged
  input and independent copies with zero/negative budgets. The real
  `ExecVerifier` command writes `4|4|1` for budget four and the process's
  inherited values/count for no budget. The existing Judge-key isolation
  test also passed. No caller sets a budget in this slice.
- Every new top-level test and subtest starts with `t.Parallel()`. The new
  tests set no process environment variable and add no real-time wait.

Focused implementation checks on this working tree:

- `rtk go test -race -count=1 -timeout 2m ./internal/corebudget` — exit 0;
  RTK reported 38 passing test entries, including the required tables and
  additional concurrency cases.
- `rtk go test -count=1 -timeout 2m -run '^(TestTheVerificationCommandReceivesTheCoreBudget|TestVerificationReceivesNoJudgeKeyVariable)$' ./internal/daemon`
  — exit 0; RTK reported 4 passing test entries.
- Both test commands initially stopped because the sandbox denied a write
  to Go's existing build cache; rerunning with approved cache access produced
  the results above. No cache or toolchain configuration was changed.
- `git diff --check` — exit 0.

The declared Verification commands and repository gates were not run; the
Daemon owns them and Task settlement. Task status, checkboxes, the Task
Graph and other Task files are unchanged by this Agent. Caller wiring and
user documentation remain with the subsequent Tasks.

## Carry-forward provenance

- Source Run: `run_20261111T064558Z_6d985452394229fa`
- Source commit: `9f0b955bb9cf6fa91f0c1f1569ed02af0d7dd104`
