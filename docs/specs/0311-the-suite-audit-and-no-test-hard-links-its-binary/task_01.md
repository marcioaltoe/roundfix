---
task: task_01
spec: 0311-the-suite-audit-and-no-test-hard-links-its-binary
status: pending
type: test
complexity: medium
---

# Task 01: The agent package adapters are copies and a fixture child makes none

## Overview

The `internal/agent` `TestMain` gives the compiled test binary three more
names, hard links under a temporary `node_modules` tree, and every
re-executed child of the binary that reaches `TestMain` does the same and
removes them when it ends; the fixture child of
`TestSpawnedFixtureDeathFailsFast` exits through `os.Exit` and never removes
them. On macOS, `syspolicyd` kills a new process whose executable path it
resolves to a name that has just disappeared, which is how the fake acpx
children of two Runs died before writing `args.json` on 2026-10-23. This
Task makes the package adapters copies of the binary, skips them in the
re-executed fixture children, and makes the fake acpx helper name the
`RunPrompt` error when the child wrote no arguments.

## Requirements

1. MUST answer the two Run failures of 2026-10-23 recorded in
   [`docs/references/2026-10-24-test-suite-audit.md`](../../references/2026-10-24-test-suite-audit.md)
   (`_prd.md` → Sources 1), by implementing API Contract 1 of `_techspec.md`
   and ADR-0324.
2. MUST change `provisionPackageAdapterLinks` in
   `internal/agent/acpx_runner_test.go` so that each of its three targets is
   a copy of `os.Args[0]` written and closed before `m.Run` starts, created
   with `os.O_CREATE|os.O_EXCL|os.O_WRONLY` and a file mode held in a named
   constant (the written-executable rule of `internal/testfixture` counts
   only integer literal modes), and removed with its temporary root after
   `m.Run` returns, as today. No `os.Link` call remains in the package.
3. MUST make `TestMain` skip `provisionPackageAdapterLinks` when
   `spawnedFixtureDeathTestEnv` is set, so the `wait` and `fixture` children
   of `TestSpawnedFixtureDeathFailsFast` create no package adapter and leave
   none behind; `packageAdapterLinks` stays nil in those children.
4. MUST make `runFakeACPXPrompt` fail, when `args.json` does not exist after
   `RunPrompt` returns, with a message that names the missing file and the
   error `RunPrompt` returned (`%v` and its type), before reading the file;
   every other path of the helper keeps its behavior.
5. MUST add the new file `internal/agent/package_adapter_test.go` with two
   tests that each call `t.Parallel()` first:
   - `TestPackageAdaptersAreCopiesOfTheTestBinary`: for each entry of
     `packageAdapterLinks` (three expected), `os.Stat` of the target and of
     `os.Args[0]` are not `os.SameFile`, the target's size equals the test
     binary's, and the target's mode has an executable bit on non-Windows
     hosts.
   - `TestASpawnedFixtureChildLeavesNoPackageAdapter`: runs `os.Args[0]`
     with `-test.run=^TestSpawnedFixtureDeathFailsFast$`, with
     `spawnedFixtureDeathTestEnv=fixture` and `TMPDIR`, `TMP` and `TEMP` set
     to a fresh `t.TempDir()`, through `environmentForTest`; requires exit
     status 23 and the output line `fixture started`; then requires that the
     temporary directory holds no entry whose name starts with
     `roundfix-agent-package-adapters-`.
   Both compile for Windows, use no Git and no real time budget, and start
   only children they wait for.
6. MUST keep every existing `internal/agent` test, golden and fixture
   unchanged apart from the lines Requirements 2 to 4 name, and keep the
   package's Sequential Test count.
7. MUST NOT change any production file, `internal/testfixture`, the
   Makefile or the CI workflows.

## Subtasks

- [ ] Copy the package adapters instead of hard-linking them, and skip them in the fixture children.
- [ ] Make the fake acpx helper name the `RunPrompt` error when no arguments were written.
- [ ] Add the two tests and run the package's adapter and fake acpx tests.

## Acceptance Criteria

- [ ] No package adapter shares the test binary's file, and no `os.Link` remains in `internal/agent`.
- [ ] A fixture child leaves no package adapter tree in its temporary directory.
- [ ] Every existing `internal/agent` test passes, and the package compiles for Windows.

## Consumers

None. The Task changes test code and documents only; no production contract or shared test helper outside its package changes.

## Context

- instruction: `internal/agent/acpx_integration_test.go`
- instruction: `internal/agent/usage_test.go`
- instruction: `internal/testfixture/written_executable_test.go`
- instruction: `internal/testfixture/parallel_tests_test.go`
- instruction: `docs/adr/0324-a-test-never-gives-a-running-executable-a-second-name.md`
- instruction: `docs/references/2026-10-24-test-suite-audit.md`
- creates: `internal/agent/package_adapter_test.go`
- interface: `internal/agent/acpx_runner_test.go`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestPackageAdaptersAreCopiesOfTheTestBinary|TestASpawnedFixtureChildLeavesNoPackageAdapter|TestSpawnedFixtureDeathFailsFast|TestCheckAdapterProvesOfficialClaudePackageAndVersion|TestFixtureBinarySurvivesConcurrentExec|TestACPXExitCodeMapping|TestRunPromptIgnoresAMalformedUsagePayload|TestNoTestWritesAnExecutableOutsideTheResidue)$' ./internal/agent ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPackageAdaptersAreCopiesOfTheTestBinary TestASpawnedFixtureChildLeavesNoPackageAdapter TestSpawnedFixtureDeathFailsFast TestACPXExitCodeMapping TestNoTestWritesAnExecutableOutsideTheResidue; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; if grep -n 'os[.]Link(' internal/agent/*_test.go; then exit 1; fi; GOOS=windows go vet ./internal/agent` — expected: exit 0. Before this Task the two new tests do not exist, so the command fails at its first pass check; with the hard links restored `TestPackageAdaptersAreCopiesOfTheTestBinary` fails on `os.SameFile` and the `os.Link` check fails, and with the fixture skip removed `TestASpawnedFixtureChildLeavesNoPackageAdapter` finds the leaked tree. Afterwards both new tests, the fixture death, adapter, fake acpx and written-executable tests pass, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior; API Contract 1; Invariants 1-2; Testing Approach; Build Order 1
