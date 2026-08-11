# Command evidence — Spec 0092 QA

Build: `9ee875c32afbdcd538baec596b41df410fd735d1`

Environment: macOS Darwin 25.5.0 arm64; Go 1.26.5. The current source was
built as `/private/tmp/roundfix-qa-0092.05U3ww/roundfix` with
`GOCACHE=<worktree>/.gocache` and `-buildvcs=false` because sandboxed VCS
stamping could not read worktree status. The resulting CLI reports
`roundfix 0.4.0`; only VCS build metadata is absent.

## Authoring and cold static gate

- `rtk roundfix spec check 0092-a-run-that-can-hand-back-its-work --strict` —
  exit 0: `No findings.`
- `GOCACHE=<worktree>/.gocache rtk proxy go clean -testcache` — exit 0 before
  each full-gate attempt.
- First `rtk make verify` attempt — sandbox blocked `api.github.com`; no code
  verdict was assigned.
- Authorized cold `rtk make verify` — exit 2 after 112.736 seconds in
  `internal/cli`. All other listed packages passed without `(cached)`.
  `TestRunReconcileJSONMatchesTextFields` reported that the actual top-level
  JSON fields include `carryForwards`, while the exact expected field set does
  not.
- `GOCACHE=<worktree>/.gocache rtk go test ./internal/cli -run
  '^TestRunReconcileJSONMatchesTextFields$' -count=1 -v` — exit 1, reproducing
  the same single failure in isolation.
- `git blame` attributes `reconcileReport.CarryForwards` in
  `internal/cli/reconcile.go:83` to Task 06 commit `9ee875c3`; the exact field
  list in `internal/cli/cli_test.go:560` predates the Spec and omits that field.

## Built CLI sweep

- `rtk /private/tmp/roundfix-qa-0092.05U3ww/roundfix --help` — exit 0. The
  reconcile synopsis is only `roundfix reconcile [run-id] [--apply] [--format
  <text|json>]`; neither disposition flag appears.
- `rtk /private/tmp/roundfix-qa-0092.05U3ww/roundfix reconcile --help` — exit 0.
  Options list only `--apply` and `--format`; neither `--discard-superseded`
  nor `--carry-forward` appears.
- `... reconcile --carry-forward --discard-superseded` — exit 2 on stderr,
  naming all three mutually exclusive mutation flags. This proves both new
  flags reach the real parser even though help omits them.
- `... reconcile --carry-forward` without a Run ID — exit 2 on stderr:
  `--carry-forward requires one stopped Run ID`.

## Selection and fallback evidence

- `go test ./internal/agent -run '^TestWorkStartedBoundary' -count=1 -v` —
  exit 0; all three cases passed, including no-output selection failure,
  first-output publication, and inert setup.
- `go test ./internal/daemon -run '^TestFallbackEligibility' -count=1 -v` —
  exit 0; selection refusal and adapter-start failure kept fallback eligible,
  while Agent output ended eligibility.
- `go test ./internal/cli -run
  '^TestAgentSelectionProfilesMacro$/(a_selection_failure_activates_the_fallback_chain|agent_output_before_failure_keeps_the_chain_ineligible|mixed_profiles_configure_validate_fallback_persist_and_stream)$'
  -count=1 -v` — exit 0; all three public-CLI integration cases passed in
  14.525 seconds. The unavailable selection used by the blocked live row is a
  no-output model refusal in the repository harness; no live external Run ID
  was created under the no-commit QA boundary.

## Batch settlement and Run outcome evidence

- `go test ./internal/agent -run
  '^TestMarkBatchFailedKeepsAlreadySettledIssues$' -count=1 -v` — exit 0; six
  cases passed for resolved, invalid, duplicated, pending, valid, and failed
  statuses.
- `go test ./internal/daemon -run
  '^(TestRunOutcomeDerivedFromUnresolvedIssues|TestRunOutcomeDerivedStaysUnresolvedAfterAFailedBatch|TestResolveCycleVerificationFailureFailsBatchAndContinues|TestResolveCycleContinuesToNextBatchAfterFailedBatch)$'
  -count=1 -v` — exit 0; all four cases passed.
- Authorized `go test ./internal/cli -run
  '^(TestRunResolveAgentFailureMarksBatchFailed|TestRunResolveVerificationFailureDoesNotCommit|TestRunResolveAgentFailureContinuesWithLaterBatches|TestFailedVerificationJournalsFailureWithoutCommitEvents)$'
  -count=1 -v` — exit 0; all four CLI integration cases passed.

## Branch and Task dispositions

- `go test ./internal/cli -run
  '^(TestReconcileDiscard|TestReconcileWithoutTheFlag|TestCarryForward)'
  -count=1 -v` — exit 0. Eight cases passed: discard record-before-removal,
  unreachable-commit refusal, record-write refusal, no-flag preservation,
  carry-forward success, moved-input refusal, mixed-set refusal, and no-flag
  carry-forward reporting.
- `go test ./internal/daemon -run
  '^(TestSupersededBranch|TestWriteBranchDispositionRecord)' -count=1 -v` —
  exit 0; four classification and durable-record cases passed.
- `go test ./internal/spec -run
  '^TestRecordCarryForwardPreservesTaskAndRecordsSource$' -count=1 -v` — exit
  0; provenance preservation passed.

These are equivalent observed integration results for rows whose destructive
built-binary journeys require a prepared committed scratch repository and Run
store. The QA authority prohibits creating commits and the current Run
Worktree cannot safely host a nested Run.

## Glossary

- `rtk rg -n 'Branch Disposition|Selection Failure|agent_selection_failed'
  CONTEXT.md internal` — exit 0 with only `Selection Failure` at
  `CONTEXT.md:199-200` and the emitted constant at
  `internal/agent/agent.go:114`. `Branch Disposition` has no glossary entry or
  emitted occurrence.

## Commit-dependent scope audit

Raw `git -c core.fsmonitor=false diff-tree --no-commit-id --name-status -r`
was run for each Daemon Task commit and repair commit:

- Task 01: `3fc542bf` — Task file plus the characterization test.
- Task 02: `94ff10b2` — only its seven bounded paths.
- Task 03: `6cdf7298`, `1152beb8` — only its Task file and three bounded source/test paths.
- Task 05: `e1ecb02e` — only its six bounded paths.
- Task 08: `1a36a3d7` — only its three bounded paths.
- Task 04: `15fda0d0`, `ecd8feb4` — only its Task file and three bounded source/test paths.
- Task 06: `7eff0fea`, `9ee875c3` — only its Task file and four bounded source/test paths.

No Task commit changed a linter, formatter, typechecker, test-runner,
architecture checker, build tool, package manager, code generator, ignore
file, plugin declaration, version pin, or other protected tooling path. The
separate graph-authoring commit `14f174da` minted Task 08, and the separate
Daemon corpus commit `d1b5608b` re-recorded `docs/references/coverage-record.json`;
neither was represented as an implementation Task commit.
