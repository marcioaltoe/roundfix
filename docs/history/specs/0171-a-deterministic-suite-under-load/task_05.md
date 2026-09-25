---
task: task_05
spec: 0171-a-deterministic-suite-under-load
status: completed
type: backend
complexity: low
---

# Task 05: An adapter fixture that survives a relative path

## Overview

`provisionFakeAdapter` in `internal/agent/acpx_runner_test.go` calls
`os.Symlink(os.Args[0], path)` and stores whatever `argv[0]` holds. Each adapter
lives in its own temporary directory, so a relative target resolves against
that directory. Under `go test` `argv[0]` is absolute and nothing fails, but
`go test -c ./internal/agent` followed by `./agent.test` creates adapters that
cannot execute.

## Requirements

1. MUST resolve `os.Args[0]` to an absolute path, against the working directory
   the test binary started in, before `provisionFakeAdapter` creates any link,
   and link that absolute path.
2. MUST keep the fixture the compiled test binary re-executed through the link,
   as ADR-0125 requires, with its behavior in the non-executable sidecar written
   by `writeFakeAdapterFixture`.
3. MUST prove the repair with
   `TestFakeAdapterRunsFromARelativeTestBinaryPath`: the test re-executes the
   compiled test binary through a relative `argv[0]`, from the package
   directory, and that child provisions a fake adapter and runs it, observing
   the fixture's configured output. Before the repair the child fails.
4. MUST prove the mirror case with
   `TestFakeAdapterRunsFromAnAbsoluteTestBinaryPath`: the same child started
   through an absolute `argv[0]` provisions and runs the adapter with the same
   output.
5. MUST leave `provisionPackageAdapterLinks`, which hard-links `os.Args[0]`
   before any test runs, unchanged unless the relative case shows it fails too,
   and MUST change no production file.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A compiled test binary started through a relative path provisions a fake
      adapter that runs and prints its configured output.
- [ ] The same holds through an absolute path.

## Context

- interface: `internal/agent/acpx_runner_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestFakeAdapterRunsFromARelativeTestBinaryPath|TestFakeAdapterRunsFromAnAbsoluteTestBinaryPath)$" ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestFakeAdapterRunsFromARelativeTestBinaryPath TestFakeAdapterRunsFromAnAbsoluteTestBinaryPath; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither test exists, so the command fails.

## References

- [_prd.md](_prd.md) — Core Feature 5; Success Metric 5
- [_techspec.md](_techspec.md) — The adapter fixture

## Result

Implementation:

- `TestMain` resolves the compiled test binary's `os.Args[0]` to an absolute
  path before any fake adapter is provisioned. `provisionFakeAdapter` uses that
  captured path for its symlink target.
- The fake adapter remains the compiled test binary with behavior stored in
  the non-executable `.fixture.json` sidecar. `provisionPackageAdapterLinks`
  remains unchanged, and no production file changed.
- `TestFakeAdapterRunsFromARelativeTestBinaryPath` and
  `TestFakeAdapterRunsFromAnAbsoluteTestBinaryPath` re-execute the compiled
  package test binary in a child, assert the child's actual `argv[0]` form,
  provision an adapter in a separate temporary directory, execute it, and
  observe the configured output.

Focused checks:

- Before the repair,
  `rtk env GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^TestFakeAdapterRunsFromARelativeTestBinaryPath$' ./internal/agent`
  failed because the adapter symlink stored `./agent.test` and could not execute
  it from the adapter directory (`no such file or directory`).
- After the repair,
  `rtk env GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^(TestFakeAdapterRunsFromA(Relative|Absolute)TestBinaryPath|TestFixtureBinarySurvivesConcurrentExec)$' ./internal/agent`
  exited 0 (`ok roundfix/internal/agent`), covering both new path cases and the
  existing compiled-binary concurrency fixture.
- The authored `## Verification` command was not run; the Daemon owns it.

Acceptance evidence:

- Relative path: `TestFakeAdapterRunsFromARelativeTestBinaryPath` passed after
  its child received relative `argv[0]`, provisioned the adapter, and observed
  `fixture-from-reexecuted-test-binary`.
- Absolute path: `TestFakeAdapterRunsFromAnAbsoluteTestBinaryPath` passed after
  its child received absolute `argv[0]`, provisioned the adapter, and observed
  the same configured output.

Follow-ups: none.
