---
task: task_06
spec: 0254-a-faster-make-test
status: completed
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
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/carryforward.go`
- interface: `internal/cli/deliver_revalidate.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/readiness_forge.go`
- interface: `internal/cli/reconcile.go`
- interface: `internal/cli/orphan_unix_test.go`
- interface: `internal/cli/readiness_forge_test.go`
- interface: `internal/testfixture/parallel_tests_test.go`
- interface: `docs/references/skill-coverage.json`
- interface: `internal/cli/deliver_item_binary_test.go`
- interface: `internal/store/journal_test.go`

## Result

Implemented the task_06 slice with eight remaining Sequential Tests and an
`internal/cli` ceiling of eight. The starting checkout had 37 explicit
`// Sequential:` tests under a ceiling of 40; 29 now call `t.Parallel()` first.

Added `getenv`, `environ`, `tempDir`, and `auditor` to `commandDependencies`.
Their production defaults are `os.Getenv`, `os.Environ`, `os.TempDir`, and
`app.Auditor`. Command environment construction, delivery children,
reconciliation children, carry-forward staging directories, and readiness
reads use those dependencies. Tests override them through
`updateCommandDependenciesForTest`; `setCommandEnvForTest` keeps single-value
reads and subprocess environment lists consistent, and direct workflow tests
use the existing `commandContextForTest` helper.

Helper subprocess variables travel through the specific command's `Env`.
The archive advice test already supplied an explicit keyless environment, so
its redundant process-wide key clearing was removed. Existing outcome
assertions, expected values, golden data, and fixture scripts retain their
meaning. Every top-level test name in every changed test file matches HEAD;
the Coverage Record consequently requires no regeneration and is unchanged.

### Sequential-test inventory and read paths

The following inventory includes every starting Sequential Test. Retained
rows keep their existing reasons.

| Test | Production read or retained reason |
| --- | --- |
| `TestArchivePlanWithoutAKeyPrintsNoAdvice` | Existing explicit commandEnvironment.environ clears provider keys; removed redundant process mutation |
| `TestSupersedeAcceptsAnArchiveRecord` | Calls an already-parallel test on the same testing.T; retained sequential |
| `TestCapabilityRecheck` | commandEnvironmentWithDependencies → dependencies.getenv(PATH) → executableDirectories |
| `TestCapabilityTextRendersProbe` | commandEnvironmentWithDependencies → dependencies.getenv(PATH) → executableDirectories |
| `TestCarryForwardAppliesThroughRefusingCommitHooks` | reconcile Git child Env → dependencies.environ; staging directory → dependencies.tempDir |
| `TestCarryForwardProofIgnoresRefusingCommitHooks` | reconcile Git child Env → dependencies.environ; staging directory → dependencies.tempDir |
| `TestCarryForwardStagingRunsNoRepositoryHook` | reconcile Git child Env → dependencies.environ; staging directory → dependencies.tempDir |
| `TestCarryForwardLeavesTheCheckoutHooksPathUnchanged` | reconcile Git child Env → dependencies.environ; staging directory → dependencies.tempDir |
| `TestCarryForwardRemovesItsEmptyHooksDirectory` | reconcile Git child Env → dependencies.environ; staging directory → dependencies.tempDir |
| `TestRunRunsWithoutSubcommandHonorsInteractivity` | commandEnvironmentWithDependencies → dependencies.getenv(ROUNDFIX_TUI) |
| `TestRunResolveDetachedChildReportsProfileProofFailure` | Raw detached-child descriptor and environment cleanup; retained sequential |
| `TestRunNoAgentConsoleRejectsInteractiveCockpit` | commandEnvironmentWithDependencies → dependencies.getenv(ROUNDFIX_TUI) |
| `TestRunPreflightFailureColorsOutputWhenForced` | commandEnvironmentWithDependencies → dependencies.getenv(ROUNDFIX_TUI) |
| `TestAttachRunBrowserLoopOpensCockpitAndRefreshes` | commandEnvironmentWithDependencies → dependencies.getenv(ROUNDFIX_TUI) |
| `TestAttachRunBrowserCancelExitsZeroWithoutAttaching` | commandEnvironmentWithDependencies → dependencies.getenv(ROUNDFIX_TUI) |
| `TestDeliveryArchiveStageCommitsTheArchiveRecord` | runRoundfix → dependencies.environ → child Env (ROUNDFIX_CLI_TEST_HELPER) |
| `TestDeliveryStepRunsTheItemBinary` | runRoundfix → dependencies.environ → child Env (DELIVERY_ITEM_* / helper marker) |
| `TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments` | runRoundfix → dependencies.environ → child Env (DELIVERY_ITEM_* / helper marker) |
| `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` | runRoundfix → dependencies.environ → child Env (DELIVERY_ITEM_* / helper marker) |
| `TestDeliveryStepParksWhenTheItemBuildFails` | runRoundfix → dependencies.environ → child Env (DELIVERY_ITEM_* / helper marker) |
| `TestDeliveryStepParksWhenTheItemBinaryPathIsNotIgnored` | runRoundfix → dependencies.environ → child Env (DELIVERY_ITEM_* / helper marker) |
| `TestDeliveryStepWithoutADeclarationRunsTheOwnerExecutable` | runRoundfix → dependencies.environ → child Env (DELIVERY_ITEM_* / helper marker) |
| `TestDeliveryItemBinaryThatCannotStartReturnsAParkError` | runRoundfix → dependencies.environ → child Env (DELIVERY_ITEM_* / helper marker) |
| `TestRevalidateWarnsWhenTheOwnerPredatesSourceOnMain` | deliveryOwnerWarning → dependencies.auditor (default app.Auditor / app.BuildCommit) |
| `TestRevalidateStaysQuietForDocsOnlyMain` | deliveryOwnerWarning → dependencies.auditor (default app.Auditor / app.BuildCommit) |
| `TestRevalidateStaysQuietForACurrentOwner` | deliveryOwnerWarning → dependencies.auditor (default app.Auditor / app.BuildCommit) |
| `TestRevalidateStaysQuietWithoutTheBuildCommit` | deliveryOwnerWarning → dependencies.auditor (default app.Auditor / app.BuildCommit) |
| `TestRevalidateDoesNotRecomputeTheOwnerWarningOnRetry` | deliveryOwnerWarning → dependencies.auditor (default app.Auditor / app.BuildCommit) |
| `TestRevalidateStaysQuietWhenTheOwnerSourceDiffFails` | deliveryOwnerWarning → dependencies.auditor (default app.Auditor / app.BuildCommit) |
| `TestDeliveryOwnerRunsAgainForAnItemRetriedDuringItsPass` | Raw detached-child descriptor and environment cleanup; retained sequential |
| `TestDeliveryOwnerReleasesAnIdleQueueAfterOnePass` | Raw detached-child descriptor and environment cleanup; retained sequential |
| `TestDeliveryOwnerFailsWhenIdleReleaseFails` | Raw detached-child descriptor and environment cleanup; retained sequential |
| `TestDeliveryOwnerReleasesItsClaimAfterEngineFailure` | Raw detached-child descriptor and environment cleanup; retained sequential |
| `TestCLIForceStopOwnerProcessHelper` | Process signal handler in the helper child; retained sequential |
| `TestForgeProbesAreBoundedAndNeverPromptOrPrintAToken` | Readiness resolver → dependencies.getenv(PATH); probe child Env → dependencies.environ |
| `TestScriptFixtureIsTheCompiledTestBinary` | Hook fixture child Env carries its private marker |
| `TestSpecJudgeReadsTheKeyFromItsOwnEnvironment` | Deliberately poisoned provider-key process environment; retained sequential |

### Focused checks and acceptance evidence

- At most 12 Sequential Tests: a Python inventory over all
  `internal/cli/*_test.go` files counted 37 before, 29 converted with
  first-statement `t.Parallel()`, and eight remaining comments; the map cap
  is eight. A separate comparison of test declarations against HEAD confirmed
  that no test name changed. `git -c core.fsmonitor=false diff --check` passed.
- Race and shuffle: the focused selection of the 29 converted tests passed
  with `-race -short` (10.786 s) and `-shuffle=on` (2.789 s). A broader focused
  selection covering affected existing command families passed with
  `-race -short` after the final code edit (16.557 s).
- Unchanged production behavior: the focused checks exercise real disposable
  repositories, compiled helper children, readiness cancellation, delivery
  binary selection, owner staleness, carry-forward hooks and cleanup, archive
  output, and TUI/color dispatch with their existing outcome assertions.
  All existing package tests and the full-package race/shuffle acceptance
  remain for Daemon Verification; these focused results make no full-package
  claim.

Commands run (all used `GOCACHE=/tmp/roundfix-task06-gocache`):

```sh
# converted is the anchored alternation of the 29 converted names in the inventory.
go test -count=1 -race -short -run "$converted" ./internal/cli
go test -count=1 -shuffle=on -run "$converted" ./internal/cli
# Broader check of the affected dependency consumers:
go test -count=1 -race -short -run '^(Test(Doctor|Delivery|CarryForward|Revalidate|Forge|GitReadiness|Readiness|Capability|Archive|ScriptFixture|History)|TestRunRunsWithoutSubcommandHonorsInteractivity$|TestRunNoAgentConsoleRejectsInteractiveCockpit$|TestRunPreflightFailureColorsOutputWhenForced$|TestAttachRunBrowser)' ./internal/cli
```

The first focused build exposed a duplicate context helper introduced during
implementation. It was removed in favor of the existing two-argument
`commandContextForTest`; subsequent checks above passed. No declared
Verification command, repository Verification, commit, push, or pull request
operation was run. The initial worktree change was the Daemon's
`status: in_progress`; that field and all other Task/manifest files were left
untouched. No follow-up outside this slice was implemented.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/archive_output_test.go`
- `internal/cli/baseline_profile_test.go`
- `internal/cli/carryforward.go`
- `internal/cli/carryforward_hooks_test.go`
- `internal/cli/deliver_archive_record_test.go`
- `internal/cli/deliver_item_binary_test.go`
- `internal/cli/deliver_owner_staleness_test.go`
- `internal/cli/deliver_revalidate.go`
- `internal/cli/deliver_workflow.go`
- `internal/cli/doctor.go`
- `internal/cli/readiness_forge.go`
- `internal/cli/readiness_forge_test.go`
- `internal/cli/reconcile.go`
- `internal/cli/script_fixture_test.go`

## Carry-forward provenance

- Source Run: `run_20261009T005858Z_c134f9dbf1a064cc`
- Source commit: `84a095f9152b55ecd880eef20f74473facc86b30`
