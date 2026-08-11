# Command evidence — Spec 0092 QA rerun

Build: `b60deccf026136d6998b36c4ae9e65bcfba1ea25`

## Authoring precondition

- `rtk roundfix spec check 0092-a-run-that-can-hand-back-its-work --strict`
  exited 0 with `No findings.`
- `_tasks.md` names `task_07` as the terminal QA node. Tasks 01–06, 08, and
  09 are `completed`; task 07 remains Daemon-owned `pending` during this gate.
- The PRD has no `Unreachable Acceptance` declaration. The QA prompt states
  that no Pull Request is open.

## Execution evidence

### Authoritative cold static gate

- `rtk env GOCACHE="$PWD/.gocache" go clean -testcache` exited 0.
- The first `rtk make verify` attempt was blocked by the sandbox when an
  integration path reached `api.github.com`; this was classified as an
  environment boundary, not a code failure.
- The repository cache was cleared again with the same command. The authorized
  `rtk make verify` rerun exited 0. All listed Go packages passed without
  `(cached)`; `internal/cli` passed in 106.458s, the focused `skills` guard
  passed, `roundfix skills check` passed, and `bin/roundfix` rebuilt from
  `b60deccf`.

### Corrective JSON contract and built CLI discovery

- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/cli -run
  '^TestRunReconcileJSONMatchesTextFields$' -count=1 -v` exited 0. This is the
  focused reproduction that failed in the prior report; Task 09 corrected it
  while retaining exact field-set equality.
- `rtk proxy bin/roundfix --help` exited 0. Its reconcile synopsis remains
  `roundfix reconcile [run-id] [--apply] [--format <text|json>]`; neither
  disposition flag appears.
- `rtk proxy bin/roundfix reconcile --help` exited 0. The options still list
  only `--apply` and `--format`.
- `rtk proxy bin/roundfix reconcile --carry-forward --discard-superseded`
  exited 2 and named all three mutually exclusive mutation flags.
- `rtk proxy bin/roundfix reconcile --carry-forward` exited 2 and reported
  that one stopped Run ID is required. These public parser refusals confirm the
  hidden flags are real behavior, not dead implementation.

### Glossary

- `rtk rg -n 'Branch Disposition|Selection Failure|agent_selection_failed'
  CONTEXT.md internal` exited 0. It found `Selection Failure` at
  `CONTEXT.md:199-200` and the emitted status at `internal/agent/agent.go:114`.
  It found no `Branch Disposition` entry.

### Selection, fallback, Batch settlement, and Run outcome

- `go test ./internal/agent -run '^TestWorkStartedBoundary' -count=1 -v`
  exited 0; all three no-output, first-output, and inert-setup cases passed.
- `go test ./internal/daemon -run '^(TestFallbackEligibility|TestRunOutcomeDerivedFromUnresolvedIssues|TestRunOutcomeDerivedStaysUnresolvedAfterAFailedBatch|TestResolveCycleVerificationFailureFailsBatchAndContinues|TestResolveCycleContinuesToNextBatchAfterFailedBatch)' -count=1 -v`
  exited 0; three fallback and four outcome cases passed.
- `go test ./internal/agent -run
  '^TestMarkBatchFailedKeepsAlreadySettledIssues$' -count=1 -v` exited 0;
  resolved, invalid, duplicated, pending, valid, and already-failed cases passed.
- `go test ./internal/cli -run
  '^(TestAgentSelectionProfilesMacro|TestRunResolveAgentFailureMarksBatchFailed|TestRunResolveVerificationFailureDoesNotCommit|TestRunResolveAgentFailureContinuesWithLaterBatches|TestFailedVerificationJournalsFailureWithoutCommitEvents)' -count=1 -v`
  exited 0. The four Run-outcome flows passed, and the macro's four cases passed
  in 13.38s, including fallback activation after no-output refusal, persisted
  fallback attempt/notification, and ineligibility after Agent output.
- Outside-evidence origin: the pre-Spec handoff
  `docs/handoffs/2026-08-08-two-quotas-out-and-a-gate-that-said-no.md` records
  the direct 2026-08-08 measurement: exhausted Codex attempts died in nineteen
  seconds with the usage-limit text wrapped as `agent/protocol error`. The
  related pre-Spec backlog record states the configured `codex -> claude`
  fallback did not activate. A fresh live Run remains blocked because this gate
  cannot create its required committed disposable target.

### Branch and Task dispositions

- `go test ./internal/cli -run
  '^(TestReconcileDiscard|TestReconcileWithoutTheFlag|TestCarryForward)'
  -count=1 -v` exited 0. Eight public-runner cases passed: durable record before
  discard, unreachable-commit refusal, record-write refusal, no-flag exact
  preservation, carry-forward success/provenance, moved-input refusal,
  mixed-set atomicity, and no-flag carry-forward preservation.
- `go test ./internal/daemon -run
  '^(TestSupersededBranch|TestWriteBranchDispositionRecord)' -count=1 -v`
  exited 0; four classification and record-persistence cases passed.
- `go test ./internal/spec -run
  '^TestRecordCarryForwardPreservesTaskAndRecordsSource$' -count=1 -v` exited
  0; Task content and source provenance survived the public record operation.
- `rtk proxy bin/roundfix reconcile --format json` exited 0 on a fresh public
  read with schema `roundfix-reconcile/v1`, `mode: dry-run`, zero applied or
  preserved candidates, and `carryForwards: []`. A second public read is the
  report-only integration case above, which asserts exact state preservation.
- `rtk proxy bin/roundfix --version` reported `roundfix 0.4.0
  (b60deccf-dirty, built 2026-08-11 01:03:24 -0300)`.

### Commit-dependent scope and tooling audit

Fresh `git -c core.fsmonitor=false diff-tree --no-commit-id --name-status -r`
was run for every implementation and repair commit:

- Task 01 `3fc542bf`: assigned Task file and characterization test.
- Task 02 `94ff10b2`: its seven bounded source/test/Task paths.
- Task 03 `6cdf7298`, `1152beb8`: its Task file and three bounded paths.
- Task 05 `e1ecb02e`: its six bounded paths.
- Task 08 `1a36a3d7`: its three bounded test/Task paths.
- Task 04 `15fda0d0`, `ecd8feb4`: its Task file and three bounded paths.
- Task 06 `7eff0fea`, `9ee875c3`: its Task file and four bounded paths.
- Task 09 `b60deccf`: only `internal/cli/cli_test.go` and `task_09.md`.

No Task commit changed repository-tooling configuration, scripts, ignore
files, plugin declarations, or version pins. Separate chronology is traceable:
`d55eab9d` authored Task 09 and its graph edge after the prior QA report
`a4c4c75e`; `79dd4d3a` separately reopened the Daemon-owned QA Task; only then
did `b60deccf` land the bounded corrective test change.

### Report contract

- Directory existence, all three blocked-count keys, `verdict: fail`, both
  coined terms, and the evidence-file path were each checked with focused
  `test`/`grep` commands; all exited 0.
- `rtk rg -n '^\| QA-' qa-report-2026-08-11-01.md` listed exactly 17 result
  rows: 8 `pass`, 2 `fail`, 7 `blocked (environment: ...)`, and no pending,
  finding-blocked, declared-blocked, or skipped row.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.
