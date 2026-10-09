---
task: task_06
spec: 0254-a-faster-make-test
status: pending
type: test
complexity: high
---

# Task 06: internal/cli's sequential tests stop sharing process-global state

## Overview

Corrective Task for QA finding F2 in `qa/qa-report-2026-10-08-01.md`. Two
back-to-back measurements put the suite at 72% and 69% of the starting
revision, against a target of at most 65%. They put the `internal/cli`
sequential phase at 43% and 44%, against a target of at most 25%.
`sequential-contributors.json` shows that the cost is the 40 tests that still
run one at a time. Under the load of the now-parallel packages they took 10 to
35 times longer than at the start, for example
`TestRevalidateStaysQuietWhenTheOwnerSourceDiffFails`, which went from 0.39 s
to 14.46 s.

These tests stay sequential because they change process-wide state:
- `DELIVERY_ITEM_*` and `ROUNDFIX_CLI_TEST_HELPER` (7 tests);
- `app.BuildCommit` (6);
- `TMPDIR` and the carry-forward hook marker (5);
- `ROUNDFIX_TUI` (4);
- `PATH` (2);
- several single tests that change other globals.

This Task gives each read a per-test path through `commandDependencies`, the
seam `updateCommandDependenciesForTest` already provides. With that path, the
tests can call `t.Parallel()`.

## Requirements

1. MUST route every production read of process-global state that a
   `// Sequential:` test in `internal/cli` changes through a field of
   `commandDependencies`, defaulting to today's process read, so production
   behavior is unchanged. This covers the environment variables named in the
   reasons and `app.BuildCommit`. Production changes are limited to
   `internal/cli` and to adding those seams.
2. MUST convert each test whose only reason is such state: replace `os.Setenv`,
   `t.Setenv` or the global swap with `updateCommandDependenciesForTest`, make
   `t.Parallel()` its first statement, and remove its `// Sequential:` comment.
   A test whose helper subprocess needs the variable MUST pass it through that
   command's `Env` rather than the process environment.
3. MUST keep a test sequential, with its existing reason, when the state cannot
   move into a dependency. The examples are the detached-child tests that own a
   raw process file descriptor, and the regression that deliberately poisons the
   process environment.
4. MUST lower the `internal/cli` cap in
   `internal/testfixture/parallel_tests_test.go` to the number of tests that
   remain sequential. That number MUST be at most 12.
5. MUST NOT change any assertion, expected value, golden or fixture meaning.
6. MUST re-record the Coverage Record through its declared command if a test
   name changes, and MUST NOT edit it by hand.

## Subtasks

- [ ] List each `// Sequential:` test in `internal/cli` with the production read it needs.
- [ ] Add the `commandDependencies` fields and route the reads through them.
- [ ] Convert the tests and lower the cap.
- [ ] Run `internal/cli` with `-race -short` and with `-shuffle=on`.

## Acceptance Criteria

- [ ] At most 12 `// Sequential:` tests remain in `internal/cli`.
- [ ] `internal/cli` passes with `-race -short` and with `-shuffle=on`.
- [ ] Production behavior is unchanged: every existing `internal/cli` test passes.

## Verification

- `out="$(go test -count=1 -v -run '^TestEveryTestInAParallelTestPackageRunsInParallel$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; grep -qE '"internal/cli": *([0-9]|1[0-2]),' internal/testfixture/parallel_tests_test.go || { printf 'internal/cli cap is above 12\n' >&2; exit 1; }; test "$(grep -rh '// Sequential:' internal/cli/*_test.go | wc -l | tr -d ' ')" -le 12 || { printf 'more than 12 Sequential tests remain\n' >&2; exit 1; }; go test -count=1 -race -short ./internal/cli && go test -count=1 -shuffle=on ./internal/cli`

## References

- `qa/qa-report-2026-10-08-01.md` → F2; `qa/evidence/*/sequential-contributors.json`
- `_prd.md` → Goal 1; `_techspec.md` → Invariant 1
- ADR-0259

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/testfixture/parallel_tests_test.go`
