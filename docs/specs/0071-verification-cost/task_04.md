---
task: task_04
spec: 0071-verification-cost
status: completed
type: backend
complexity: high
---

# Task 04: Free the Baseline package from process state

## Overview

The Baseline package is the second heaviest at 83.9s, with 16 environment
mutations and one working-directory change blocking parallelism. This Task
applies the same two moves as the CLI package — remove the process-global
dependency, then declare parallelism — on a smaller surface. It reduces the sum
of package times and speeds a single-package run; it cannot move the suite
floor, which the CLI package sets.

## Requirements

1. MUST make functions reading process environment or working directory, for
   values their callers know, receive them as parameters, with the default
   resolved once at the command boundary.
2. MUST declare parallelism on every test that no longer mutates process state
   and owns the filesystem it touches.
3. MUST leave a one-line reason on every test that stays sequential.
4. MUST fix, not silence, tests that fail only under parallel execution.
5. MUST prove no races and no cross-test leakage.
6. MUST leave the coverage record and the package's observable behavior
   unchanged.

## Subtasks

- [ ] Give the process-reading functions their parameters.
- [ ] Declare parallelism where process state is no longer touched.
- [ ] State a reason on every test left sequential.
- [ ] Fix shared-state failures and prove no races.

## Acceptance Criteria

- [ ] Production functions no longer read the process for values callers
      supply.
- [ ] The package's parallel-declaring test count rises above its recorded 28.
- [ ] Every test left sequential carries a one-line reason.
- [ ] The package passes with race detection and with its tests run twice in
      one invocation.
- [ ] The coverage record from task 01 is unchanged.
- [ ] The package completes measurably faster than 83.9s on the same machine.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`,
      `internal/`, and this task file.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/catalog_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `go test ./internal/baseline -count=2` — expected: exit 0.
- `go test ./internal/baseline -race -count=1` — expected: exit 0.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0.
- `go vet ./internal/baseline` — expected: exit 0.

## References

- `_prd.md` → Core Features 1 and 2.
- `_techspec.md` → Build Order 4; Risks (task 04 helps the sum, not the floor).
- ADR-0089.

## Result

### Implementation

- Replaced Baseline executable discovery's process `PATH` and implicit working-
  directory reads with explicit absolute search directories on Plan, Profile
  Alignment, and Capability Recheck requests. The CLI now captures `PATH` and
  the working directory once in `commandEnvironment`, resolves empty and
  relative entries there, and passes the normalized directories through both
  non-interactive and human Baseline flows.
- Removed all 16 Baseline test environment mutations and the working-directory
  mutation. Executable-probe fixtures now pass their private directories per
  call, while Git characterization dates are set on each test subprocess.
- Classified all 176 runnable top-level Baseline tests. 170 isolated tests now
  declare unconditional parallelism. The six regeneration-capable tests carry
  a first-line `// Sequential:` reason and declare parallelism only when their
  shared-artifact update flag is disabled.
- Daemon Verification exposed one shared-state defect missed by the initial
  focused groups: the catalog load-count test replaced the package-global
  embedded assets while unrelated plan tests were running. The test now passes
  its counting filesystem to the existing parameterized regeneration loader,
  leaving the immutable embedded catalog untouched without removing
  parallelism or weakening assertions.

### Focused checks

- `rtk env GOCACHE=/private/tmp/roundfix-0071-task04-gocache go test
  ./internal/baseline -run
  '^(TestCapabilityRecheckMatchesFullPlan|TestDivergenceCarriesProbeEvidence|TestExecutableCandidateResolution|TestExecutableCandidateNeverExecutes|TestExecutableEvidenceDistinguishesFailureFromAbsence|TestCapabilityAuditNoExecution|TestDivergenceRendersProbe|TestCapabilityTextRendersProbe|TestCapabilityTextAndJSONAgree|TestDivergenceGroupsByRequirement|TestBaselinePlanCharacterization)$'
  -race -count=2 -parallel=12` — passed (`ok`, 20.061s).
- `rtk env GOCACHE=/private/tmp/roundfix-0071-task04-gocache go test
  ./internal/baseline -run
  '^(TestPlan|TestFormatter|TestInstruction|TestADR|TestFinding|TestGreenfield|TestRepositoryExtension|TestPreservation|TestResidual|TestNestedCarrier)'
  -race -count=2 -parallel=12` — passed (`ok`, 39.499s).
- `rtk env GOCACHE=/private/tmp/roundfix-0071-task04-gocache go test
  ./internal/baseline -run
  '^(TestApply|TestTransaction|TestSkills|TestAssets|TestBaselineVerification|TestEmptyReapply|TestImmutable|TestManagedRoot|TestProfileDraft)'
  -race -count=2 -parallel=12` — passed (`ok`, 137.420s).
- `rtk env GOCACHE=/private/tmp/roundfix-0071-task04-gocache go test
  ./internal/baseline -run
  '^(TestCatalog|TestSource|TestGuidance|TestProjectDecision|TestCarrier|TestClassification|TestSegmentation|TestRevision|TestRepository|TestInventory|TestBounded|TestInstructionAlias|TestBaselineRelease|TestBaselineFinding|TestTooling)'
  -race -count=2 -parallel=12` — passed (`ok`, 37.362s).
- `rtk env GOCACHE=/private/tmp/roundfix-0071-task04-gocache go test
  ./internal/cli -run
  '^(TestBaseline|TestHuman|TestProjectDecision|TestGuidance|TestProfileAdaptation|TestConsolidated|TestDivergence|TestBetterAuth|TestRepeatedPlan|TestCapability|TestToolingAuthority)'
  -count=1 -parallel=12` — passed (`ok`, 31.646s), exercising the CLI
  command-boundary propagation.
- Static accounting reported `tests=176 parallel=170 sequential=6 missing=`
  and 188 total `t.Parallel()` declarations. A source sweep found no
  `t.Setenv`, `t.Chdir`, `os.Getenv`, `os.LookupEnv`, or `os.Getwd` in
  `internal/baseline`.
- `rtk git diff --exit-code --
  docs/specs/0071-verification-cost/coverage-record.json` — passed; the task 01
  record is byte-unchanged.
- `rtk git diff --check` — passed.
- `rtk git -c core.fsmonitor=false status --short` — only `internal/baseline/`,
  `internal/cli/`, and this Task file are changed.

### Acceptance evidence

- Explicit production inputs: supported by the request plumbing, the zero-hit
  process-read sweep, and the focused executable/CLI checks.
- Parallel-declaring count above 28: supported by 188 declarations, including
  170 unconditional top-level declarations.
- Sequential rationale: supported by the six first-line reasons and zero
  unclassified top-level tests.
- Race and repeated-run isolation: supported for the focused groups above; the
  complete package race and repeated invocations remain for Daemon
  Verification.
- Coverage-record preservation: supported by the clean path-specific diff.
- Observable behavior: supported by the unchanged plan-characterization
  goldens and the focused CLI Baseline suite.
- Whole-package runtime below 83.9s: remains for Daemon Verification on the
  same machine; focused timings do not support that package-wide claim.
- Changed-path boundary: supported by the final short status output.

### Verification feedback repair — attempt 1

- The attempt 1 diagnostic identified `TestRegenerationLoadsCatalogOnce` and
  three concurrently running plan tests. A focused four-test reproduction
  failed on all three repetitions, confirming that the regeneration test's
  assignment to `embeddedAssets` leaked its stale-digest fixture into normal
  `BuildPlan` catalog loads.
- Removed the global asset swap and cleanup. The load-count fixture is now
  injected directly through `loadCatalogForRegeneration`, the existing
  filesystem-parameterized seam, so the test still proves one catalog load
  while owning all mutable state it touches.
- `GOCACHE=/private/tmp/roundfix-0071-task04-gocache rtk go test
  ./internal/baseline -run
  '^(TestRegenerationLoadsCatalogOnce|TestInstructionHierarchyRendersActivePointersOnce|TestInstructionHierarchyPreservesPlanAndResultSchemas|TestADRLifecycleContract)$'
  -race -count=10` — passed (80 test executions), proving the reported
  interaction no longer races or leaks across repeated runs.
- A source sweep now finds the sole `embeddedAssets` assignment at its package
  declaration; no test replaces the shared embedded catalog.
