---
task: task_03
spec: 0071-verification-cost
status: completed
type: test
complexity: high
---

# Task 03: Run the CLI tests in parallel

## Overview

With the process-global dependency removed, the CLI package's tests can declare
parallelism — and this is the only slice that can move the headline number,
because packages already overlap and this package is the floor. Parallel
execution also surfaces shared state that sequential execution hid; those are
defects to fix, not reasons to revert.

## Requirements

1. MUST declare parallelism on every test in the package that no longer mutates
   process state and owns whatever filesystem it touches.
2. MUST leave a one-line reason on every test that stays sequential, so
   sequential execution is a recorded decision rather than an omission.
3. MUST fix, not silence, any test that fails only under parallel execution: a
   test that passes alone and fails alongside others has found shared state.
4. MUST prove the absence of races and cross-test leakage by running the
   package with race detection and repeated execution.
5. MUST leave the coverage record unchanged.
6. MUST make the package measurably faster than its recorded 113.2s baseline.

## Subtasks

- [ ] Declare parallelism where process state is no longer touched.
- [ ] State a reason on every test left sequential.
- [ ] Fix shared-state failures surfaced by overlapping execution.
- [ ] Prove no races and no cross-test leakage.

## Acceptance Criteria

- [ ] The package's parallel-declaring test count rises from one to a
      substantial share of its tests.
- [ ] Every test left sequential carries a one-line reason.
- [ ] The package passes with race detection enabled.
- [ ] The package passes when its tests run twice in the same invocation,
      proving no state leaks between overlapping tests.
- [ ] The coverage record from task 01 is unchanged.
- [ ] The package completes measurably faster than 113.2s on the same machine.
- [ ] `git status --porcelain` shows no path outside `internal/cli/` and this
      task file.

## Context

- interface: `internal/cli/cli_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -count=2` — expected: exit 0; no state leaks between
  overlapping repeated tests.
- `go test ./internal/cli -race -count=1` — expected: exit 0; no data race.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0; no test disappeared.
- `grep -rc 't.Parallel()' internal/cli | grep -qv ':1$'` — expected: exit 0;
  the package no longer has a single parallel declaration.

## References

- `_prd.md` → Core Features 1 and 2; Success Metrics (`internal/cli` faster
  than 113.2s).
- `_techspec.md` → Build Order 3; Risks (parallelising surfaces real defects).
- ADR-0089.

## Result

### Implementation

- Classified all 487 runnable top-level `internal/cli` tests by direct and
  transitive process/shared-state use. Added `t.Parallel()` to 207 isolated
  tests (42.5% of the package) that own their filesystem and do not override a
  process or package-level seam.
- Kept 280 tests sequential because they override package-level test seams,
  verify a process-level environment default, or coordinate process-wide
  signal handling. Every retained sequential test now carries its reason as
  the first line of the test body.
- Relied on Go's top-level parallel-test scheduling boundary: the sequential
  seam-overriding tests finish before the parallel group resumes, while the
  isolated tests overlap with one another.
- No shared-state failure surfaced in the initial focused overlapping runs.
  Daemon Verification later exposed a caller-process lifecycle race under the
  broader package overlap; the focused repair evidence is recorded below.

### Focused checks

- `rtk env GOCACHE=/private/tmp/roundfix-go-cache go test ./internal/cli -run
  '^(TestBaseline|TestHuman|TestProjectDecision|TestGuidance|TestProfileAdaptation|TestConsolidated|TestDivergence|TestBetterAuth|TestRejectedPlan|TestRepeatedPlan|TestCapability|TestToolingAuthority)'
  -race -count=2` — passed (`ok`, 68.620s).
- `rtk env GOCACHE=/private/tmp/roundfix-go-cache go test ./internal/cli -run
  '^(TestRunReconcile|TestCommandUsage|TestEventsHelp|TestProfilesShow|TestProfileProof|TestProveProfile|TestInvocationProfile|TestHealthChecker|TestRunSkills|TestStageable|TestReleasePlan|TestResolveSelection|TestRunSettle)'
  -race -count=2` — passed (`ok`, 12.716s).
- `rtk env GOCACHE=/private/tmp/roundfix-go-cache go test ./internal/cli -run
  '^(TestLoadCommitted|TestRunImplementHelp|TestRunImplementVerificationCapacityDoesNotAddFlag|TestRunImplementDetach|TestRunHelpListsImplementCommand|TestRunImplementValidationFailures|TestRunImplementRejectsInvalidVerificationCapacity|TestRenderImplementTaskLines|TestRunImplementRejectsExplicitEmptySelection|TestRunImplementSelectionOverride|TestAgentSelectionProfilesMacro)'
  -race -count=2` — passed (`ok`, 21.097s).
- Static accounting over `internal/cli/*_test.go` reported `tests=488
  parallel=207 sequential=280`; the remaining function is `TestMain`.
- A source-shape check found no runnable top-level test whose first body line
  lacked either `t.Parallel()` or a `// Sequential:` reason.
- `git diff --exit-code --
  docs/specs/0071-verification-cost/coverage-record.json` — passed; the task 01
  coverage record is unchanged.
- `git diff --check` — passed.
- `git -c core.fsmonitor=false status --short` — only `internal/cli/` test
  files and this Task file are changed.

### Acceptance evidence

- Parallel-declaring share: supported by 207 declarations across 487 runnable
  top-level tests, up from the recorded single declaration.
- Sequential rationale: supported by 280 first-line reasons and the zero-output
  source-shape check.
- Race safety and repeated-run isolation: supported for the three focused
  overlapping groups above; the full-package race and repeated invocations
  remain for Daemon Verification.
- Coverage-record preservation: supported by the clean path-specific diff.
- Package runtime below 113.2s: remains for Daemon Verification on the same
  machine; focused checks do not support a whole-package timing claim.
- Changed-path boundary: supported by the final short status output.

### Verification feedback repair — attempt 1

- Inspected the Daemon diagnostic for the failed repeated package invocation.
  `TestRunImplementDetachSurvivesCallerProcessGroupKill` reached its process
  group kill after the caller could exit naturally, so the test had no
  lifecycle barrier guaranteeing that its kill target remained live. Broader
  parallel load exposed that race as a process-group permission failure and
  incomplete temporary-directory cleanup.
- Added an opt-in test-helper lifecycle barrier: the original CLI helper
  registers signal readiness before `Run`, then remains alive after `Run`
  returns until the test kills its isolated process group. A Detached Run child
  is excluded by its handshake environment, so the test still proves that the
  child survives caller-group termination. No delay, retry, accepted error, or
  sequential fallback was added.
- `rtk env GOCACHE=/private/tmp/roundfix-go-cache go test ./internal/cli -run
  '^TestRunImplementDetach(PrintsReportAndCompletesRun|ReportsAndRelaysPreflightFailure|SurvivesCallerProcessGroupKill)$'
  -count=10 -parallel=3` — passed (`ok`, 10.642s).
- `rtk env GOCACHE=/private/tmp/roundfix-go-cache go test ./internal/cli -run
  '^TestRunImplementDetach(PrintsReportAndCompletesRun|ReportsAndRelaysPreflightFailure|SurvivesCallerProcessGroupKill)$'
  -race -count=5 -parallel=3` — passed (`ok`, 11.071s).
