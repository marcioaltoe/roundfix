---
task: task_02
spec: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
status: completed
type: backend
complexity: high
---

# Task 02: A proof whose session answered passes when only its close fails, and the session is recorded

## Overview

The verdict of `ProveExactSelection` and of the Fallback Chain candidate
probe rests on their setup alone. The close of the disposable Agent Session
runs under the attempt's Proof Setup Deadline instead of 5 seconds; a failed
close prints one warning and writes an Unclosed Proof Session record; a setup
step that reaches its deadline is named. Baseline semantic analysis keeps
its stop on an unclosed proof. This answers the adopted Backlog Entry of
2026-10-27 (a profile proof fails when closing its disposable session under
load), the Vortex report of 2026-10-26.

## Requirements

1. MUST characterize first (lesson 9): before changing production code,
   write `TestAProofWhoseSessionAnsweredPassesWhenItsCloseHangs` against the
   current tree with the fake acpx harness, a passing setup (the model
   variant fixture of `TestProveExactSelectionModelVariant` suffices) and
   `sessions close` blocked through the harness's block-command setting,
   observe the `*AgentSessionCleanupError` wrapping
   `context.DeadlineExceeded` that `_prd.md` → Validation records, then turn
   it into the new contract. The close deadline in this test is the budget
   it asserts: set `cleanupTimeout` on the harness runner, and add the
   ceiling `"internal/agent/proof_close_test.go": {1, …}` to
   `realTimeBudgetCeilings` in `internal/testfixture/real_time_budget_test.go`
   with a reason naming the blocked close. No other real-time budget is
   added, and no other ceiling moves.
2. MUST implement `_techspec.md` → API Contract 1: when every setup step
   succeeded and the close fails with anything but a proven absence,
   `ProveExactSelection` returns the proof, a nil error and
   `SelectionProof.Unclosed`; a failed setup returns its own error with no
   `*AgentSessionCleanupError` in the chain; `probeFallbackCandidate` keeps
   the candidate's result whatever the close does. A session that never
   answered still fails (Invariant 1).
3. MUST implement API Contract 2: the proof closes under the attempt's
   resolved setup deadline and the candidate probe under
   `acpxPreflightSetupTimeout`; the close context still derives from
   `context.WithoutCancel(ctx)`; the `cleanupTimeout` seam still overrides
   when positive; `RunSealedPrompt` (`internal/agent/sealed.go`), `EndSession`
   and Run session cleanup keep their deadlines.
4. MUST implement API Contract 3: wrap a setup step error that satisfies
   `errors.Is(err, context.DeadlineExceeded)` as `profile proof step
   "<step>": %w` with the step names the contract lists; every other error
   keeps its bytes.
5. MUST implement API Contract 4 in the new
   `internal/agent/unclosed_proof_sessions.go`: the record at
   `<HOME>/.roundfix/unclosed-proof-sessions/<session>.json` with `HOME`
   read from `runner.baseEnv()`, directory `0700`, file `0600`, written
   through a temporary file and rename, holding only the Data Models fields;
   and the two warning lines, printed once per failed close through
   `runner.warningf`. A record never holds an environment value (Invariant
   5).
6. MUST implement API Contract 8 in `Analyzer.proveSelection`
   (`internal/baselineacp/analyzer.go`) and add
   `TestAProofThatLeftItsSessionOpenStopsSemanticAnalysis` to
   `internal/baselineacp/analyzer_test.go`: a fake runtime whose preferred
   proof carries `Unclosed` makes segmentation, classification and revision
   fall back to their manual proposals without starting a sealed prompt, as
   `TestPreferredFallbackDoesNotStartAfterProofCleanupFailure` does for the
   close error.
7. MUST add the tests `_techspec.md` → Testing Approach names for task_02 in
   the new `internal/agent/proof_close_test.go`, each failing when its
   behavior is removed: `TestAProofWhoseSessionAnsweredPassesWhenItsCloseHangs`,
   `TestTheProofCloseDeadlineFollowsTheSetupDeadline` (with `cleanupTimeout`
   zero and `SetupTimeout` two hours, a close that exits `2` yields a warning
   naming `2h0m0s`; a 5-second close deadline would name `5s`),
   `TestASessionThatNeverAnsweredStillFailsWhenItsCloseFails`,
   `TestATimedOutProofStepIsNamed` (block the first `sessions show`, expire
   the parent context the test controls, assert `profile proof step "read
   runtime catalogue"` and `errors.Is(err, context.DeadlineExceeded)`),
   `TestAFallbackCandidateKeepsItsResultWhenItsCloseFails` (a sequential test
   with `t.Chdir`, as `TestACPXProbeFallbackReportsCleanupInfrastructureError`
   is), and `TestAnUnclosedProofSessionRecordHoldsNoEnvironmentValue` (a
   sentinel value built at run time by concatenation in the runner
   environment, under an OpenRouter key variable name, never appears in the
   record). Each top-level test calls `t.Parallel()` first or carries a
   `// Sequential: <reason>` comment.
8. MUST remove the four tests that assert the old join, declared as breaks in
   `_techspec.md` → Testing Approach, from
   `internal/agent/acpx_runner_test.go`:
   `TestProveExactSelectionCleanupJoinedFailure`,
   `TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose`,
   `TestACPXProbeCleanupFailureJoinsSelectionError` and
   `TestACPXProbeFallbackReportsCleanupInfrastructureError`. Their cases
   move to the new contract in the tests of requirement 7. No other existing
   test changes.
9. MUST re-record `docs/references/coverage-record.json` after removing those
   four tests, only with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
   (`docs/agents/specific-repository.md`); `TestCoverageEquivalence` runs in
   the repository Verification at settlement and fails while the record still
   lists them. The record's only change is the removal of those four names and
   the addition of the tests this Task and earlier Specs added.

## Subtasks

- [ ] Characterize the hanging close on the current tree.
- [ ] Make the setup decide the verdict and give the close the setup deadline.
- [ ] Name the timed-out step.
- [ ] Write the record and the warnings.
- [ ] Keep the analyzer's stop.
- [ ] Add the named tests, remove the four join tests, add the ceiling.
- [ ] Re-record the coverage record with its command.

## Acceptance Criteria

- [ ] A passing setup with a hanging or failing close proves the tuple, prints one warning and writes one record.
- [ ] A rejected setup fails with its own error and no close error in the chain.
- [ ] A setup deadline names its step.
- [ ] Baseline semantic analysis stops on an unclosed proof as it did on the close error.

## Context

- instruction: `docs/adr/0344-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails-and-the-next-proof-closes-it.md`
- instruction: `docs/adr/0298-a-timed-out-profile-proof-gets-a-third-attempt-and-its-setup-deadline-is-configurable.md`
- instruction: `docs/user-guide/configuration.md`
- instruction: `internal/agent/sealed.go`
- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/fallback.go`
- interface: `internal/agent/acpx_runner.go`
- creates: `internal/agent/unclosed_proof_sessions.go`
- creates: `internal/agent/proof_close_test.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/baselineacp/analyzer.go`
- interface: `internal/baselineacp/analyzer_test.go`
- interface: `internal/testfixture/real_time_budget_test.go`
- interface: `docs/references/coverage-record.json`

## Consumers

- `SelectionProof.Unclosed` and the proof verdict — consumers: `./internal/agent`, `./internal/baselineacp`
- The real-time budget ceilings — consumers: `./internal/testfixture`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAProof.*|TestTheProofCloseDeadlineFollowsTheSetupDeadline|TestASessionThatNeverAnswered.*|TestATimedOutProofStepIsNamed|TestAFallbackCandidate.*|TestAnUnclosedProofSession.*|TestProveExactSelection.*|TestACPXProbe.*|TestExactSelectionProofSetup.*|TestTimedOutProofSetup.*|TestMissingSession.*|TestPreferredFallback.*|TestManualClassification.*|TestRuleSegmentation.*|TestRevisionManualFallback|TestEveryRealTimeBudgetIsWithinItsFileCeiling)$' ./internal/agent ./internal/baselineacp ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAProofWhoseSessionAnsweredPassesWhenItsCloseHangs TestTheProofCloseDeadlineFollowsTheSetupDeadline TestASessionThatNeverAnsweredStillFailsWhenItsCloseFails TestATimedOutProofStepIsNamed TestAFallbackCandidateKeepsItsResultWhenItsCloseFails TestAnUnclosedProofSessionRecordHoldsNoEnvironmentValue TestAProofThatLeftItsSessionOpenStopsSemanticAnalysis TestPreferredFallbackDoesNotStartAfterProofCleanupFailure TestExactSelectionProofSetupUsesTheRequestedDeadline TestEveryRealTimeBudgetIsWithinItsFileCeiling; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; if printf '%s\n' "$out" | grep -q -- '--- PASS: TestProveExactSelectionCleanupJoinedFailure '; then echo 'the old join test still runs' >&2; exit 1; fi; GOOS=windows go build -buildvcs=false ./internal/agent ./internal/baselineacp` — expected: exit 0. Before this Task the six new agent tests and the analyzer test do not exist, so the first pass check fails. Afterwards they pass with the existing proof, probe, setup-deadline, analyzer and real-time budget tests that consume the changed contracts, the old join test is gone, and the packages build for Windows.

## References

- `_prd.md` → Goals; Core Features 1-3; Core Feature 5; User Experience; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior 1-4, 6; Interfaces; Data Models; API Contracts 1-4; API Contract 8; Invariants 1-3, 5; Testing Approach; Build Order 1
- ADR-0344; ADR-0298; ADR-0069

## Result

Implemented this Task's proof-close contract. Exact proofs and fallback
candidates keep their setup verdict; their close contexts derive from
`context.WithoutCancel(ctx)` and use the resolved setup deadline, with the
existing positive `cleanupTimeout` override. Sealed prompts, EndSession and
Run cleanup were not changed. Deadline errors name the contracted setup
step; other setup diagnoses retain their bytes.

Failed closes emit one warning and atomically record only the specified
schema, session, working directory, runtime fields, timestamp and cause
under the runner environment's HOME. The new directory is `0700` and the
temporary file, renamed into place, is `0600`. Missing HOME and write
failures emit the alternative warning with the acpx close arguments.
Baseline analysis converts `SelectionProof.Unclosed` into the existing
cleanup-unproven error before any sealed prompt or fallback proof starts.

Acceptance evidence:

- Passing setup with a blocked or failing close: the hanging-close and
  setup-deadline tests retain a proven tuple, one warning and one record;
  the latter observes `2h0m0s` with the cleanup override disabled. The
  fallback test covers functional, rejected and error outcomes. The record
  test checks environment-value exclusion, the variable name and private
  permissions.
- Rejected setup: `TestASessionThatNeverAnsweredStillFailsWhenItsCloseFails`
  covers both exact proof and Probe, asserts the original rejection bytes,
  rejects any cleanup error in the chain and checks the warning and record.
- Named deadline: `TestATimedOutProofStepIsNamed` blocks the first catalogue
  inspection, expires the controlled parent context, and asserts both
  `profile proof step "read runtime catalogue"` and `errors.Is` matching.
- Analyzer stop: `TestAProofThatLeftItsSessionOpenStopsSemanticAnalysis`
  covers segmentation, classification and revision, asserts their manual
  proposals and that neither a fallback proof nor a sealed prompt starts;
  it also checks the cleanup error's session and warning cause.

Focused checks:

- Before production edits,
  `go test ./internal/agent -run '^TestAProofWhoseSessionAnsweredPassesWhenItsCloseHangs$' -count=1 -v`
  characterized the old behavior: `*AgentSessionCleanupError` wrapping
  `context.DeadlineExceeded` after the model-variant setup succeeded, with
  `cleanupTimeout = 100 * time.Millisecond`. That characterization was
  then changed to assert the new contract.
- The six named agent tests plus the analyzer test passed together with an
  exact-name `go test ./internal/agent ./internal/baselineacp -run ... -count=1 -v`
  selection.
- `go test ./internal/agent -run '^(TestProveExactSelection|TestACPXProbe|TestExactSelectionProofSetup|TestTimedOutProofSetup|TestMissingSession|TestAProofWhose|TestTheProofClose|TestASessionThat|TestATimedOutProof|TestAFallbackCandidate|TestAnUnclosedProof)' -count=1`
  exited 0.
- `go test ./internal/baselineacp -run '^(TestAProofThatLeftItsSessionOpenStopsSemanticAnalysis|TestPreferredFallback|TestManualClassification|TestRuleSegmentation|TestRevisionManualFallback)' -count=1`
  exited 0.
- `go test ./internal/testfixture -run '^(TestEveryRealTimeBudgetIsWithinItsFileCeiling|TestEveryTestInAParallelTestPackageRunsInParallel)$' -count=1`
  exited 0. Only the required new ceiling of one was added; no existing
  ceiling changed.
- Seven temporary `go test -overlay` mutations independently removed the
  unclosed proof field assignment, setup-derived close deadline, setup-only
  error verdict, deadline step name, fallback verdict preservation,
  environment-value exclusion and analyzer stop. Each named regression
  test failed on its corresponding mutation. Overlays remained outside the
  repository; production files were not mutated for these probes.

Removed exactly the four declared old join tests, with no edits to other
existing tests. Generated the coverage record only through
`go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`.
The recorder wrote the record and its test body passed, but its process
exited 1 because suiteguard reported
`modified: docs/references/coverage-record.json`. Inspection found that the
recorder does not call `suiteguard.DeclareSanctionedRegeneration`, although
this Spec's authorization already declares its command and output. No
suiteguard bypass or recorder code change was made.

The subsequent read-only
`go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` exited 0.
A structural comparison of the generated record against HEAD found exactly
the four authorized removals, 318 common-test additions (including this
Task and earlier Specs), three platform-test additions and an unchanged
platform list.

Follow-up outside this slice: make the existing coverage recorder declare
its sanctioned regeneration before writing, so a first write can exit 0
under suiteguard. Reaping and CLI warning propagation remain task_03's
work. Task status and checkboxes are unchanged; declared Verification,
repository Verification, commits, pushes and pull requests were not run.
