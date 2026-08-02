---
task: task_02
spec: 0071-verification-cost
status: pending
type: backend
complexity: high
---

# Task 02: Let the CLI package take its environment

## Overview

The suite's floor is the CLI package at 113.2s of essentially sequential work,
and the reason is not oversight: 20 environment mutations and 18
working-directory changes each make Go refuse `t.Parallel()`, because both
mutate state the whole process shares. This Task removes the reason. Functions
that read the process for values their caller already knows take them as
parameters instead, and the process default resolves once at the command
boundary. No test declares parallelism yet — that is the next slice.

## Requirements

1. MUST make functions that read process environment variables or the process
   working directory, for values a caller already knows, receive those values
   as parameters instead.
2. MUST resolve the process default exactly once, at the command boundary,
   preserving today's behavior for every real invocation.
3. MUST leave the CLI's observable surface unchanged: same commands, flags,
   output, and exit codes.
4. MUST reduce the count of process-state mutations in the package's tests,
   leaving only those whose subject is the process-level default itself.
5. MUST leave every remaining process-state mutation accompanied by a one-line
   reason stating why the test needs it.
6. MUST NOT declare parallelism in this Task.

## Subtasks

- [ ] Identify the functions reading process environment or working directory.
- [ ] Give them parameters and resolve the default at the command boundary.
- [ ] Convert the tests that mutated process state incidentally.
- [ ] State a reason on every mutation that remains.

## Acceptance Criteria

- [ ] Production functions no longer read the process for values their callers
      supply; the default resolves at the command boundary.
- [ ] The CLI's commands, flags, output, and exit codes are unchanged.
- [ ] The package's process-state mutations are reduced, and each remaining one
      carries a one-line reason.
- [ ] The coverage record from task 01 is unchanged.
- [ ] No test in the package declares parallelism yet.
- [ ] `git status --porcelain` shows no path outside `internal/cli/`,
      `internal/`, and this task file.

## Context

- interface: `internal/cli/cli.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -count=1` — expected: exit 0; behavior is unchanged
  while still sequential.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0; no test disappeared.
- `go vet ./internal/cli` — expected: exit 0.

## References

- `_prd.md` → Core Features 1; Goals (uses the machine it runs on).
- `_techspec.md` → Implementation Design: Interfaces; Build Order 2; Risks.
- ADR-0089.
