---
task: task_02
spec: 0213-a-test-suite-that-does-not-flake
status: completed
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

## Result

Implemented this Task's test-only slice. `TestMain` records the absolute compiled
binary path and dispatches sidecars before either environment-selected helper
mode. The nine executable-write sites now create symlinks to that binary and
mode-`0o600` sidecars. The macro ACPX sidecar launches its existing Python body
through `python3 -c`, preserving arguments and prompt stdin. The Darwin Codex
fixture remains separately compiled for its existing signature check.

The new fixture test exercises the adapter-version helper, fake ACPX and both
adapters, macro ACPX and all four adapters, all four carry-forward hooks, and
the extracted settle-hook helper. It checks symlink targets and sidecar modes,
then checks output and exit status by absolute path, relative path, and shell
PATH lookup. It also proves that arguments with spaces, environment values,
and inherited CLI/detach helper-mode variables survive dispatch. The settle
fixture refuses a real staged 501-line file; the carry-forward hooks record
execution and preserve their refusing exits.

The AST guard scans literal integer modes in `os.WriteFile` and `os.OpenFile`
calls under `internal/**/*_test.go`, including aliased imports, and freezes the
TechSpec's ten residual sites across nine files. Its sabotage cases assert a
seeded `0o755` write reports `internal/seeded/fixture_test.go:3`, both one-too-high
and one-too-low inventory counts fail, and exact counts pass. Its header names
constants, computed modes, and `os.Chmod` as outside the boundary.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| No new executable writes outside the residue; seeded executable and wrong counts fail | The guard's repository, seeded-site, count-mismatch, and literal-boundary subtests all passed. The converted files no longer contain literal executable file writes. `cli_test.go` and `upgrade_test.go` are unchanged. |
| Each helper produces the compiled-binary symlink and private sidecar, with absolute and PATH execution | `TestScriptFixtureIsTheCompiledTestBinary` passed three repetitions at each of `-cpu 1,4`, also covering relative invocation. On Darwin, the separately compiled macro Codex is explicitly outside the script-fixture assertion; the non-Darwin branch uses the shared fixture and its assertion. Linux cross-compilation passed; Linux execution was not performed here. |
| Doctor, carry-forward, settle-hook and macro behavior remains intact | Both Doctor floor tests passed three repetitions at each CPU setting. The focused integration selection passed, including all carry-forward hook tests, measured settle-hook recovery, macro profiles and detached Implement completion. The authored 20-repeat selection remains for Daemon Verification. |

### Focused checks

All successful local commands below used
`GOCACHE=/tmp/roundfix-task02-gocache` after sandbox access to the host Go cache
was denied:

- `rtk proxy go test -count=1 -v -run 'TestNoTestWritesAnExecutableOutsideTheResidue/' ./internal/testfixture` — exit 0; all four guard subtests passed.
- `rtk proxy go test -count=3 -cpu 1,4 -run 'TestScriptFixtureIsTheCompiledTestBinary/|TestDoctor.*AnAdapter' ./internal/cli` — exit 0 after the final fixture-test edit.
- `rtk proxy go test -count=1 -run 'TestCarryForward.*Hook|TestSettleRecoversMeasuredHookRefusedWork|TestAgentSelectionProfilesMacro|TestRunImplementDetachPrintsReportAndCompletesRun' ./internal/cli` — exit 0, 103.803s, with the files held unchanged.
- `GOOS=linux GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test -c -o /tmp/roundfix-task02-linux-cli.test ./internal/cli` — exit 0.
- `rtk proxy go test -count=1 -run '^TestRunForceStop(OwnerProcessIntegrationProvesExitBeforeStoreCompletion|LegacyRunWithoutOwnerIdentityStillStopsOwner)$' ./internal/cli` — exit 0 with sandbox escalation for process-table access.
- `rtk make skills-sync-check skills-check build` — exit 0.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

### Incremental check and limitations

`GOCACHE=/tmp/roundfix-task02-gocache rtk make verify-incremental` exited 2.
Formatting and vet passed, but the test stage failed. This is not a passing
incremental gate:

- Several package guards detected my edit to `script_fixture_test.go` while
  that check was running. The overlapping focused integration selection also
  reported this guard violation after its assertions passed. I stopped editing
  and reran the affected focused selection successfully.
- Two owner-stop integration tests could not enumerate the process table in
  the sandbox. Both passed in the focused escalated check recorded above.
- `internal/daemon/TestTaskBudgetReasonNamesTheSettlementThatRenewedIt`
  exhausted its existing 200ms Run Budget before the stalled Task started
  (`elapsed 299.612541ms`). This is the pre-existing flake assigned to
  `task_01`; no daemon code or other Task file was changed here. The full
  incremental gate was not retried over that unresolved out-of-slice flake.
- The first fixture check lacked the fake ACPX helper's isolated home and
  attempted its config write under the host home; the sandbox refused it.
  The fixture test now sets a temporary command home, and the later focused
  fixture checks pass.

The pre-change inspection found all nine executable-write sites and neither
new test file. No production code, Task Graph, other Task file, protected CLI
test, or upgrade test was edited. Status remains Daemon-owned. No authored
Verification command, commit, push, or pull-request operation was performed.

## Carry-forward provenance

- Source Run: `run_20261002T064202Z_a22e1dd4644c4dc6`
- Source commit: `09973cf2e449af6bfbe063ee07fdeb9c702551eb`
