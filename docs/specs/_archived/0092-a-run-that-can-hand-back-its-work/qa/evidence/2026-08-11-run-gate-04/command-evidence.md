# Command evidence — Spec 0092 QA rerun 04

Build: `2020661682b043e19662bf9fc272ca04ef8dc9c0`

This file records fresh current-build commands and concise outcomes for
`qa-report-2026-08-11-04.md`.

## Preconditions and build

- `rtk roundfix spec check 0092-a-run-that-can-hand-back-its-work --strict`
  exited 0 with `No findings.`
- `rtk git rev-parse HEAD` returned
  `2020661682b043e19662bf9fc272ca04ef8dc9c0`.
- Before evidence writes, `rtk git -c core.fsmonitor=false status --short
  --untracked-files=all` was empty.
- The manifest names `task_07` as the `qa` node and all six direct dependencies
  are `completed`; the PRD contains no `Unreachable Acceptance` section.

## Cold static gate

- `rtk env GOCACHE="$PWD/.gocache" go clean -testcache` exited 0.
- The first `rtk make verify` attempt reached a repository integration test and
  was environment-blocked by denied access to `api.github.com`; it produced no
  code verdict.
- The same repository cache was cleared again. The authorized exact `rtk make
  verify` rerun exited 0. Every listed package was uncached; `internal/cli`,
  `internal/daemon`, `internal/spec`, the focused `skills` guard, and `roundfix
  skills check` all passed. The binary rebuilt from `20206616`.
- `rtk make verify-docs` exited 0. The docs contract, repository contract,
  corpus budget, binary build, and Spec consistency checks passed; Spec 0092
  reported `No findings.`

## Selection, fallback, settlement, and outcome

- Focused current-build Agent boundary cases passed: first-output publication,
  no-output selection failure, and inert setup exclusion.
- Focused Daemon fallback cases passed: no-output selection failure and adapter
  start failure activated the configured fallback; Agent output kept it
  ineligible.
- Focused `TestAgentSelectionProfilesMacro` plus the public CLI outcome cases
  passed, including selection-failure activation, post-output ineligibility,
  preserved Review Issue state, and terminal Run outcome reads.
- Focused `TestMarkBatchFailedKeepsAlreadySettledIssues` and
  `TestSettleAssignedIssues` passed.
- Four Daemon outcome cases and the corresponding CLI failure/continuation and
  journal cases passed. A positive unresolved count stayed `Unresolved`; a
  failed Batch with zero unresolved work could reach `Clean`.
- A live fallback Run remains unreachable: this active Run Worktree cannot host
  a nested committed disposable Spec Run, and this gate has no commit authority.
  Equivalent current-build Agent, Daemon, and public CLI evidence above passed.
- A real Review Source Batch remains unreachable because the QA prompt confirms
  there is no open Pull Request. Equivalent current-build settlement/outcome
  evidence above passed.

External acceptance origin: the repository handoff
`docs/handoffs/2026-08-08-two-quotas-out-and-a-gate-that-said-no.md` records a
real 2026-08-08 Codex quota exhaustion, nineteen-second Run failures, premature
`AGENT_WORK_STARTED`, and the Fallback Chain not activating. That pre-Spec
measurement is distinct from this Spec's authored artifacts.

## Branch disposition and carry-forward

- Eight focused current-build branch-disposition cases passed across
  `internal/cli` and `internal/daemon`: reachable and later-Run supersession,
  durable record before removal, unreachable-commit refusal, record-write
  refusal, and report-only preservation.
- Six focused current-build carry-forward cases passed across `internal/cli`
  and `internal/spec`: unchanged-input success with provenance, moved-input
  refusal, mixed-set atomic refusal, report-only preservation, full input-set
  comparison, and provenance recording.
- Direct QA journeys need prepared committed scratch Run Branches, Worktrees,
  Run stores, and stopped Runs. The no-commit gate cannot construct those Git
  histories. The focused cases above provide fresh equivalent success,
  refusal, independent-read, persistence, and atomicity evidence.

## Built CLI, docs, and vocabulary

- `rtk ./bin/roundfix --help` and `rtk ./bin/roundfix reconcile --help` exited
  0. Both publish
  `roundfix reconcile [run-id] [--apply | --discard-superseded | --carry-forward]
  [--format <text|json>]`; command help describes both disposition acts.
- `rtk ./bin/roundfix reconcile --discard-superseded --carry-forward` exited 2
  and named the three mutation flags as mutually exclusive, then pointed to
  command help.
- Two consecutive `rtk ./bin/roundfix reconcile --format json` invocations
  exited 0 and returned byte-identical `roundfix-reconcile/v1` dry-run payloads
  with zero applied mutations and `carryForwards: []`.
- The exact synopsis appears in `internal/cli/cli.go` twice,
  `internal/cli/cli_test.go`, and `docs/user-guide/commands.md:616`. The guide
  also names and describes both flags at lines 670-677. This resolves prior
  F-05.
- `CONTEXT.md:199-201` defines `Selection Failure` and
  `agent_selection_failed`; `CONTEXT.md:239-241` defines `Branch Disposition`.
  `internal/agent/agent.go` emits the documented token.

## Commit and scope audit

- A fresh `rtk proxy git rev-list --reverse main..HEAD` returned 28 commits
  through `2020661682b043e19662bf9fc272ca04ef8dc9c0`.
- Fresh `rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-status
  -r <commit>` commands exited 0 for every one of the 28 commits.
- Every attributed implementation Task commit stayed within its Task's bounded
  paths. Authoring, QA, and reopening commits changed Spec artifacts; the
  derived coverage re-record changed only
  `docs/references/coverage-record.json`. The post-QA guide repair changed only
  `docs/user-guide/commands.md`. No linter, formatter, typechecker, test-runner,
  architecture-checker, build-tool, package-manager, code-generator, ignore
  file, plugin declaration, or version pin changed.
- `rtk git -c core.fsmonitor=false diff --name-status main..HEAD` confirmed the
  assembled feature path set.

## Report contract

- Targeted checks confirmed `status: closed`, `verdict: pass`,
  `rows_blocked_environment: 7`, `rows_blocked_finding: 0`, and
  `rows_blocked_declared: 0` in `qa-report-2026-08-11-04.md`.
- Targeted checks found both `Selection Failure` and `Branch Disposition` in
  the report and confirmed this evidence file resolves.
- The report contains 17 result rows: 10 end in `pass`, 7 end in
  `blocked (environment: ...)`, and 0 end in `pending`.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.
- Final status before this evidence-only addendum listed exactly the untracked
  rerun report and this evidence file. No commit, push, or Pull Request action
  was performed.
