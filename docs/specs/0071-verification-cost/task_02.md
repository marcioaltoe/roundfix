---
task: task_02
spec: 0071-verification-cost
status: completed
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

## Result

Implementation-ready for the Daemon-owned Verification step:

- `Run` and `RunContext` now capture home, working directory, terminal/color
  settings, detach values, child environment, Codex path, and branch actor once
  in `commandEnvironmentFromProcess`. Command handlers receive that snapshot;
  config loading, repository discovery, release planning, Setup/Doctor checks,
  detach, and rendering no longer resolve those defaults themselves.
- CLI tests now pass explicit home and working-directory values through
  concurrency-safe test helpers. The package has 9 process-state mutations,
  down from the recorded 38: 7 exercise TUI/color defaults and 2 exercise the
  capability check's PATH default. Every remaining mutation has an adjacent
  one-line reason.

Acceptance evidence:

1. `rg` over non-test `internal/cli` sources found `os.UserHomeDir`, `os.Getwd`,
   `os.Getenv`, and `os.Environ` only in `commandEnvironmentFromProcess`.
   `roundconfig.Load` now has one call site behind `loadCommandConfig`.
2. Focused behavior check:
   `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/cli -run '^(TestRunInitCreatesProjectConfigWithExplicitScope|TestRunSetupHealthyMachineIsIdempotent|TestRunDoctorDerivesExternalSkillRequirementFromSetupManifest|TestRunGCDryRunListsEligibleRunsAndChangesNothing|TestRunArchiveMovesCompletedSpecAndStampsMetadata|TestRunSkillsInstallCopiesArtifactsToProjectByDefault|TestRunStopByRunIDRecordsStopRequest|TestEventsReplayDefaultAndFilterJSONLRecordsOnly|TestRunReconcileDryRunReadOnly|TestRunSettleRequiresSpecAndTask|TestReleasePlanCommandMatchesPRDOutcomes|TestProfilesConfigureFileWritesProjectProfileJSON|TestBaselineProfileCommandInitShowAndValidate|TestCapabilityRecheck|TestRunImplementValidationFailures|TestRunRunsWithoutSubcommandHonorsInteractivity|TestRunPreflightFailureColorsOutputWhenForced|TestAttachRunBrowserCancelExitsZeroWithoutAttaching)$' -count=1`
   exercised 51 cases successfully across the affected command surfaces.
3. `rg -c 't\.(Setenv|Chdir)\(' internal/cli -g '*_test.go'` reports 7 in
   `cli_test.go` and 2 in `baseline_profile_test.go`; inspection with one line
   of context confirms a reason immediately precedes every mutation.
4. `git diff --exit-code -- docs/specs/0071-verification-cost/coverage-record.json`
   exits 0; the task 01 coverage record is unchanged.
5. `rg -n 't\.Parallel\(' internal/cli -g '*_test.go'` finds no matches. The
   pre-existing declaration was removed and this Task added none.
6. `git status --short` lists only this task file and paths under
   `internal/cli/`. `git diff --check` exits 0.

Additional focused compile evidence:
`GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/cli -run '^$' -count=1`
exits 0. The commands under `## Verification` were not run; the Daemon owns
that gate.

### Verification feedback repair — attempt 1

The Daemon check exposed two gaps in the explicit test environment. A parent
test's home/work override did not reach its subtests, and subprocess-based
detach tests did not project the explicit home and working directory into the
child process. This made those commands read the real Run Database or a
different temporary home.

The test environment registry now resolves overrides through the test-name
ancestry and cleans each registration at its owning test boundary. Explicit
home values are also projected into the child environment, and detached
commands receive the command boundary's resolved working directory through
`exec.Cmd.Dir`.

Focused repair evidence:

- Before the repair,
  `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/cli -run '^TestRunReconcileInvalidSelectorsMutateNothing$' -count=1`
  reproduced four failing cases, and the equivalent focused detach check
  reproduced its child-environment failure.
- After the repair,
  `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/cli -run '^(TestRunReconcileInvalidSelectorsMutateNothing|TestRunRunsListStateFlagFiltersAndNotes|TestRunRunsListTerminalAndAllReportRetainedWorktreesByRepository|TestRunRunsListLimitBoundsNewestMatches|TestRunImplementDetachPrintsReportAndCompletesRun|TestRunImplementDetachSurvivesCallerProcessGroupKill)$' -count=1`
  exercised all diagnostic test groups successfully: 20 cases.
- `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/cli -run '^TestRunDetachedCommand' -count=1`
  exercised 7 detach handshake cases successfully.
- The mutation audit remains 9 documented process-default cases, the coverage
  record remains unchanged, and `git diff --check` exits 0.

The failed declared Verification command was not rerun; the Daemon owns its
next full attempt.
