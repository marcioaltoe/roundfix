# Command evidence — Spec 0092 QA rerun 03

Build: `efe682367289d313d89ba2179f7a2917da3346ed`

## Authoring precondition

- `rtk roundfix spec check 0092-a-run-that-can-hand-back-its-work --strict`
  exited 0: `No findings.`

## Build and static gates

- `rtk env GOCACHE="$PWD/.gocache" go clean -testcache && rtk make verify`
  reached the integration suite, then the managed sandbox denied
  `api.github.com`.
- The same cold-cache command was restarted with that integration boundary
  authorized. It exited 0. Every listed Go package was uncached;
  `roundfix/internal/cli` passed in 104.153s, the Roundfix skill check passed,
  and `bin/roundfix` rebuilt from `efe68236`.
- `rtk ./bin/roundfix --version` reported
  `roundfix 0.4.0 (efe68236-dirty, built 2026-08-11 01:57:32 -0300)`; `dirty`
  is the uncommitted QA report and evidence only.
- `rtk make verify-docs` exited 0. It ran the tagged docs contract, sanctioned
  ownership checks, corpus budget, and `bin/roundfix spec check`; Spec 0092
  reported `No findings.`

## Selection, settlement, and Run outcome

- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/agent
  ./internal/daemon -run
  '^(TestWorkStartedBoundaryPublishesOnFirstAgentOutput|TestWorkStartedBoundaryReportsSelectionFailureWithoutOutput|TestWorkStartedBoundaryIgnoresInertSessionSetup|TestFallbackEligibilitySurvivesASelectionFailure|TestFallbackEligibilityEndsAfterAnyAgentOutput|TestFallbackEligibilitySurvivesAdapterStartFailure|TestMarkBatchFailedKeepsAlreadySettledIssues|TestRunOutcomeDerivedFromUnresolvedIssues|TestRunOutcomeDerivedStaysUnresolvedAfterAFailedBatch|TestResolveCycleVerificationFailureFailsBatchAndContinues|TestResolveCycleContinuesToNextBatchAfterFailedBatch)$'
  -count=1` exited 0 for both packages.
- The focused CLI journey command selecting
  `TestAgentSelectionProfilesMacro`, the two Agent-failure flows, the
  Verification-failure flow, and the failure-journal flow was initially denied
  at `api.github.com`; its authorized rerun exited 0 in 13.363s.
- The passed macro includes both named guards: no-output selection failure
  activates the configured fallback, while Agent output before failure keeps
  it ineligible. The outcome cases preserve resolved Review Issues and derive
  `Clean` or `Unresolved` from remaining work.
- No live fallback Run was created. The active per-Run worktree and no-commit
  QA contract do not provide an authorized committed disposable Spec target.
  External origin remains the pre-Spec quota observation in
  `docs/handoffs/2026-08-08-two-quotas-out-and-a-gate-that-said-no.md:16-20`
  and `:164-171`.
- No live multi-issue Pull Request Batch was created. The QA prompt states that
  no Pull Request is open; the focused current-build Agent, Daemon, and CLI
  cases are equivalent observed evidence.

## Branch disposition and carry-forward

- The authorized focused command over `internal/cli`, `internal/daemon`, and
  `internal/spec` selected 17 named cases covering the repaired exact help
  contract, JSON field contract, read-only reconcile, branch classification,
  durable record ordering, unreachable-commit refusal, record-write refusal,
  carry-forward success, moved-input refusal, mixed-set atomicity, no-flag
  preservation, input coverage, and provenance. It exited 0 for all three
  packages.
- Direct destructive journeys remain environment-blocked: branch disposition
  needs prepared committed scratch Git history, and carry-forward needs a
  prepared committed stopped Run. The no-commit contract forbids constructing
  those histories. The named integration cases provide equivalent
  current-build evidence for success, refusal, persistence, and atomicity.
- `rtk ./bin/roundfix reconcile --format json` exited 0 in `dry-run` mode with
  `applied: 0`, `results: []`, and `carryForwards: []`.
- Two fresh built-binary JSON reports were captured under a temporary directory
  and compared with `cmp`; the comparison exited 0 and each file was 593
  bytes, confirming stable report-only output.

## CLI and documentation surface

- `rtk ./bin/roundfix --help` and
  `rtk ./bin/roundfix reconcile --help` both show
  `roundfix reconcile [run-id] [--apply | --discard-superseded |
  --carry-forward] [--format <text|json>]` and describe both new flags.
- Invoking both disposition flags together exited 2 with the diagnostic
  `--apply, --discard-superseded, and --carry-forward are mutually exclusive`;
  no mutation path ran.
- `rtk ./bin/roundfix implement --help` and
  `rtk ./bin/roundfix resolve --help` expose only the existing complete
  one-Run selection tuple. Neither exposes a one-Run Fallback Chain override.
- Focused source inspection found the complete reconcile synopsis in
  `internal/cli/cli.go:51`, `internal/cli/cli.go:5296`, and
  `internal/cli/cli_test.go:432`.
- The same inspection found the public command guide still publishing
  `roundfix reconcile [run-id] [--apply] [--format <text|json>]` at
  `docs/user-guide/commands.md:616`; neither new flag appears in that guide.
  This is F-05. The prior QA reports contain no `commands.md`, pre-Spec
  synopsis, or F-05 history, so the finding is new.
- `CONTEXT.md:199-201` defines `Selection Failure` and its
  `agent_selection_failed` token. `CONTEXT.md:239-241` defines
  `Branch Disposition`; current source emits the documented vocabulary.

## Commit and scope audit

- `rtk git log --reverse --format='%H %s' main..HEAD` listed 25 commits through
  `efe682367289d313d89ba2179f7a2917da3346ed`.
- Fresh `rtk git diff-tree --no-commit-id --name-status -r <commit>` inspection
  covered every commit. One first-pass hash was mistyped; the corrected
  `6cdf7298f55336043f46ae5e990b2acc7493ba5b` inspection changed only
  `task_03.md`.
- Every implementation Task commit stayed within its bounded paths. Authoring,
  QA, and repair commits changed only Spec artifacts or the derived
  `docs/references/coverage-record.json`. No protected tooling configuration,
  script, ignore file, plugin declaration, or version pin changed.
- `rtk git -c core.fsmonitor=false diff --name-status main..HEAD` confirmed the
  assembled path set. `rtk git -c core.fsmonitor=false status
  --short --untracked-files=all` listed only this QA report and evidence as
  untracked. `rtk git diff --check` exited 0.

## Report contract

- Report-specific checks confirmed `status: closed`, `verdict: fail`,
  `rows_blocked_environment: 7`, `rows_blocked_finding: 0`,
  `rows_blocked_declared: 0`, both coined glossary terms, and the resolving
  evidence file.
- Result-row counts are exact: 9 `pass`, 1 `fail`, and 7
  `blocked (environment: ...)`.
- A focused search for pending result rows, `status: in-progress`, or
  `verdict: pending` exited 1 with no matches.
- Final `rtk roundfix spec check
  0092-a-run-that-can-hand-back-its-work --strict` exited 0 with
  `No findings.`; `rtk git diff --check` also exited 0.
