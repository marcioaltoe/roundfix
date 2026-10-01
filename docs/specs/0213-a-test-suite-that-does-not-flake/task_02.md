---
task: task_02
spec: 0213-a-test-suite-that-does-not-flake
status: pending
type: test
complexity: high
---

# Task 02: Script fixtures are the compiled test binary, and a guard holds the class

## Overview

`TestDoctorAcceptsAnAdapterAtTheFloor/codex` failed in CI because the test
writes a shell script and executes it at once. On Linux, a child forked
concurrently can keep the script open for writing, and `execve` refuses with
`ETXTBSY`. ADR-0125 forbids that shape, and eight more written executables
in non-governed `internal/cli` test files share it. This Task adds one script
fixture that re-executes the compiled test binary, and converts the nine call
sites. It also adds a repository guard that fails on any new written
executable outside the named residual inventory.

## Requirements

1. MUST add `writeScriptFixture` and `runScriptFixture` in the new
   `internal/cli/script_fixture_test.go`, as the TechSpec Interfaces and "The
   script fixture and its call sites" describe. The fixture path MUST be a
   symlink to the absolute compiled test binary. The body MUST live in a
   sidecar of mode `0o600`, which `/bin/sh` runs with the original arguments
   and environment.
2. MUST make `TestMain` in `internal/cli/implement_test.go` record the
   absolute test binary path and call `runScriptFixture` before any
   environment-selected mode. A sidecar MUST be found when the fixture is
   invoked by absolute path, by relative path, and by bare name through
   `PATH`.
3. MUST convert the nine write sites, in five helpers, listed in the TechSpec:
   - `doctorWithVersionFixture`, split into
     `writeAdapterVersionFixture(t, pkg, version) string`;
   - `fakeACPXCommand`'s three executables;
   - `newMacroFakeACPX`'s five executables and the non-Darwin branch of
     `buildMacroCodexExecutable`;
   - the carry-forward hooks;
   - the settle pre-commit hook.
4. MUST add `TestScriptFixtureIsTheCompiledTestBinary` in `internal/cli`. For
   each fixture helper, it asserts that the path is a symlink to the absolute
   test binary, that the sidecar has mode `0o600`, and that the fixture runs
   and prints its expected output, both by absolute path and through `PATH`.
5. MUST add `TestNoTestWritesAnExecutableOutsideTheResidue` in the new
   `internal/testfixture/written_executable_test.go`. It scans
   `internal/**/*_test.go` for `os.WriteFile` and `os.OpenFile` calls whose
   literal mode sets any `0o111` bit, and compares the sites with the exact
   inventory in the TechSpec section "The residual inventory". It MUST fail,
   naming `file:line`, on a temporary tree whose test file writes `0o755`, and
   MUST fail when an inventory count differs by one. Its suite header MUST
   name constants and `os.Chmod` as outside its boundary.
6. MUST NOT edit `internal/cli/cli_test.go`, `internal/cli/upgrade_test.go`
   or any file outside `internal/cli` and `internal/testfixture`, and MUST NOT
   change production code.
7. MUST keep every converted test's assertions, including the Doctor floor
   and below-floor results and the hook refusals.

## Subtasks

- [ ] Add the script fixture and its `TestMain` dispatch.
- [ ] Convert the nine write sites.
- [ ] Prove the fixture shape with `TestScriptFixtureIsTheCompiledTestBinary`.
- [ ] Add the guard with its residual inventory and its sabotage cases.

## Acceptance Criteria

- [ ] No non-governed `internal/cli` test file writes a file with an execute
      bit except `upgrade_test.go`. The guard passes on the tree, fails on a
      seeded written executable, and fails on a wrong inventory count.
- [ ] Each fixture helper yields a symlink to the test binary with a `0o600`
      sidecar, and the fixture runs by absolute path and through `PATH`.
- [ ] The Doctor floor tests, the carry-forward hook tests, the settle hook
      test and the macro profile test pass. The floor tests also pass 20
      repeated runs at `-cpu 1,4`.

## Context

- creates: `internal/cli/script_fixture_test.go`
- creates: `internal/testfixture/written_executable_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/adapter_floor_test.go`
- interface: `internal/cli/carryforward_hooks_test.go`
- interface: `internal/cli/settle_test.go`
- instruction: `docs/adr/0125-a-spawned-fixture-is-a-compiled-binary-never-a-written-script.md`
- instruction: `internal/worktree/adminlock_guard_test.go`
- instruction: `internal/agent/acpx_runner_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestNoTestWritesAnExecutableOutsideTheResidue$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestNoTestWritesAnExecutableOutsideTheResidue' || { printf '%s\n' "$out"; exit 1; }` — expected: exit 0; before this Task the guard does not exist, so the command fails.
- `out="$(go test -count=20 -cpu 1,4 -v -run '^(TestScriptFixtureIsTheCompiledTestBinary|TestDoctorAcceptsAnAdapterAtTheFloor|TestDoctorRefusesAnAdapterBelowTheFloor)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestScriptFixtureIsTheCompiledTestBinary TestDoctorAcceptsAnAdapterAtTheFloor/codex TestDoctorAcceptsAnAdapterAtTheFloor/claude; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task `TestScriptFixtureIsTheCompiledTestBinary` does not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestCarryForwardAppliesThroughRefusingCommitHooks|TestCarryForwardProofIgnoresRefusingCommitHooks|TestCarryForwardStagingRunsNoRepositoryHook|TestSettleRecoversMeasuredHookRefusedWork|TestAgentSelectionProfilesMacro|TestRunImplementDetachPrintsReportAndCompletesRun)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCarryForwardAppliesThroughRefusingCommitHooks TestSettleRecoversMeasuredHookRefusedWork TestAgentSelectionProfilesMacro TestRunImplementDetachPrintsReportAndCompletesRun; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && ! grep -nE 'WriteFile\(.*0o?7[0-7][0-7]\)' internal/cli/adapter_floor_test.go internal/cli/implement_test.go internal/cli/carryforward_hooks_test.go internal/cli/settle_test.go` — expected: exit 0; before this Task those files write executables, so the final check fails.

## References

- [_prd.md](_prd.md) — Goals 1-2; Core Feature 3; Success Metric 3; Acceptance evidence
- [_techspec.md](_techspec.md) — Interfaces; The script fixture and its call sites; The residual inventory; Testing Approach 2; Build Order 2
- ADR-0125; ADR-0126
