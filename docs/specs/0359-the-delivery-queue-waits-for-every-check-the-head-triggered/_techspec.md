---
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
prd: _prd.md
created: 2026-11-06
---

# The Delivery Queue waits for every check the head triggered — Technical Spec

## Executive Summary

Four slices. First, the documents: ADR-0369, the amendment notices in
ADR-0311 and ADR-0365, `CONTEXT.md`, the `deliver` and configuration guides
and the Roundfix skill's `deliver` reference, which state the rule the code
slices then implement. Second, the GitHub boundary: `CurrentHeadChecks` also
reads the GitHub Actions workflow runs of the head commit and the head's
Expected Workflows, and returns them in the `CheckReport`; the engine ignores
them in this slice. Third, the checking stage applies the completion rule to
the Head Check Set: it waits while any member is pending or any Expected
Workflow has no run, parks `checks-failed` naming every failed member with its
link, names the pending members on the console log and in the
`checks-timeout` blocker, and never moves on the merge state alone. Fourth,
the Queue Owner records the pending members on the Turn Record and `deliver
status` prints them on the checking item's `Turn:` line. The trade-off this
design accepts is that a workflow file filtered by `paths`, branch globs or
`branches-ignore` creates no expectation, so a filtered workflow whose run
GitHub has not created in the seconds after the event can still be missed
when no unfiltered workflow exists; in exchange there is no settle delay and
no re-implementation of GitHub's filter semantics, and a wrong guess can never
hold an item until the Check Wait ends.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; workflow
  runs keep GitHub's numeric ids as decimal strings and items keep their Spec
  slugs. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint or credential of
  Roundfix changes; the queue reads `repos/{owner}/{repo}/actions/runs`
  through the operator's authenticated `gh api`, as `InspectFailedCheck`
  already reads check-run annotations, and reads workflow files with local
  `git`. Source: `docs/agents/go.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0311: "**One key bounds every
  Check Wait.**"; the wait keeps `delivery.check_timeout`, its default and
  its restart rule, and ADR-0369 amends what it waits for and what its
  blocker names. ADR-0236: "A failure without the annotation keeps the rules
  it had before"; the stale
  and outside-change re-runs keep their code path. ADR-0365: "**The turn is a
  file.**"; the pending members join the Turn Record and no stage or column
  is added, and ADR-0369 amends the `Turn:` line of an item in `checking`.
  ADR-0304: "A new `waiting` stage or an item column was rejected"; nothing
  is stored in the Run Database. ADR-0199: "Tokens outside Runs, such as the pre-PR
  review, do not count."; it carries no obligation here, listed
  because ADR-0365 cites it. ADR-0237: "A Delivery Retry records a merge made
  outside the queue"; it carries no obligation either, listed because it cites
  ADR-0199, and the retry's merge observation is unchanged. ADR-0369, written by task_01, records the
  decision below. Source: `docs/agents/domain.md`, `docs/agents/cli.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ batch 3 answers of 2026-11-06 in the authorization record; bounded files
  are the Roundfix skill files task_01 changes. No Makefile, workflow,
  `go.mod`, `go.sum`, lint, formatter or Baseline module change;
  `gopkg.in/yaml.v3` is already a direct dependency. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0359-the-delivery-queue-waits-for-every-check-the-head-triggered/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read at `6ceb7ff5`.

1. `Engine.checkCandidate` (`internal/delivery/engine.go`) polls
   `PullRequestBoundary.CurrentHeadChecks` every `checkInterval` (15 s) until
   the Check Wait deadline. `GitHubCLI.CurrentHeadChecks`
   (`internal/delivery/github.go`) runs `gh pr view <n> --json
   pullRequestJSONFields`, returns early on `CONFLICTING`, runs `gh pr checks
   <n> --json bucket,link,name,state,workflow` (every check run and commit
   status, without `--required`), and `gh pr view` again to prove the head did
   not move.
2. In `checkCandidate`, an empty check list or `Mergeable` `UNKNOWN` keeps
   waiting; a `pass` or `skipping` bucket passes; a `fail` bucket goes through
   the ADR-0236 re-run rules (`InspectFailedCheck`, `RerunFailedCheck`) and
   otherwise parks the bare blocker `checks-failed`; `cancel` parks bare
   `checks-failed`; any other bucket keeps waiting; a merge state `BLOCKED` or
   `UNKNOWN` keeps waiting; otherwise the item moves to `merging`. The only
   log line is `roundfix: checks: Delivery Queue item <slug>: merge state
   <state>`, printed when the state changes. A read error is retried silently.
3. At the deadline the item parks `checks-timeout: current-head checks not
   settled within delivery.check_timeout <duration>` (`checksTimeoutBlocker`).
   `ClassifyPark` (`internal/delivery/park_class.go`) reads only the text
   before the first `:`.
4. On Fluxus #466, #469 and #477 (references) the three workflow runs of each
   head existed eight seconds after the pull request was created, while the
   job check runs appeared more than a minute later; the queue read one
   passing check, `CLEAN`, and merged. The queue never read workflow runs.
5. The Turn Record (`internal/delivery/turn.go`, `internal/cli/deliver_turn.go`)
   records `advancing` for the item the owner advances; `TurnLines` prints
   `Turn: <slug> advancing at <stage>` for it, including in `checking`.
6. A prototype on a scratch copy of `6ceb7ff5` implemented API Contracts 1 to
   4. In `./internal/delivery` the tests Declared breaks names failed, each on
   blocker or log text only, and the four scripted `CurrentHeadChecks` tests
   failed on the two new commands; in `./internal/cli` run with `-run
   'Deliver|Delivery|Queue|FollowUp|Status'` only
   `TestCommandDeliveryEngineRecordsThePullRequest` failed, because its fake
   `gh` refuses `gh api`.

## System Architecture

- `internal/delivery/head_check_set.go` (new, task_02): `WorkflowRun` and
  `ExpectedWorkflows`.
- `internal/delivery/head_check_wait.go` (new, task_03): the Head Check Set
  evaluation and the `checks-failed` text.
- `internal/delivery/github.go`: the two new reads in `CurrentHeadChecks`,
  `baseRefName` in the pull request read.
- `internal/delivery/engine.go`: the completion rule, the log lines and the
  `checks-timeout` text in `checkCandidate`.
- `internal/delivery/turn.go`: the pending members on the Turn Record and the
  checking `Turn:` line.
- Documents: ADR-0369, notices in ADR-0311 and ADR-0365, `CONTEXT.md`, the
  `deliver` and configuration guides, the Roundfix `deliver` reference.

## Implementation Design

### Interfaces

```go
// internal/delivery/github.go
const pullRequestJSONFields = "number,url,state,baseRefName,headRefName,headRefOid,mergedAt,mergeCommit,mergeable,mergeStateStatus"
type PullRequest struct { BaseBranch string /* + existing fields */ }
type CheckReport struct { /* existing */; WorkflowRuns []WorkflowRun; ExpectedWorkflows []string }

// internal/delivery/head_check_set.go
type WorkflowRun struct { ID, Name, Path, Event, Status, Conclusion string; Attempt int; URL string }
func ExpectedWorkflows(files map[string][]byte, base string) []string

// internal/delivery/head_check_wait.go
type headCheckSet struct { Pending, Failed []string; RunFailures, Members int }
func evaluateHeadCheckSet(report CheckReport, rerunAttempts map[string]int) headCheckSet
func checksFailedBlocker(report CheckReport) string

// internal/delivery/engine.go
func checksTimeoutBlocker(timeout time.Duration, pending []string, mergeState string) string
// internal/delivery/turn.go
type ItemTurn struct { /* existing */; Checks []string `json:"checks,omitempty"` }
func (engine *Engine) checkTurn(ctx context.Context, gitRoot, slug string, pending []string)
```

### Data Models

No Run Database change. The Turn Record item of an item in `checking` gains
an optional `checks` list, omitted when empty:

```json
{"spec":"0149-example","stage":"checking","turn":"advancing","checks":["CI | Backend / Lint, Type Check, Test & Build",".github/workflows/ci-frontends.yml (no run yet)"]}
```

### API Contracts

1. API Contract: the reader. `GitHubCLI.CurrentHeadChecks` runs, in this
   order: `gh pr view <n> --json <pullRequestJSONFields>`; on `CONFLICTING`
   it returns as today; `gh pr checks <n> --json bucket,link,name,state,workflow`
   as today; `gh api repos/{owner}/{repo}/actions/runs?head_sha=<head>&per_page=100`;
   `git ls-tree -z --name-only <head> -- .github/workflows/`; for each listed
   path ending in `.yml` or `.yaml`, in listed order, `git show <head>:<path>`;
   then `gh pr view` again as today. `PullRequest.BaseBranch` is the trimmed
   `baseRefName`. A non-zero exit of the runs read is an error starting `read
   workflow runs:`; a `total_count` greater than the number of runs returned is
   the error `read workflow runs: <total> runs reported for <head>, more than
   one page of <n>`; a failed `git ls-tree` or `git show` is an error starting
   `list workflow files:` or `read workflow file:`. Runs whose `head_sha` is
   not the head are dropped; of runs with the same `path` and `event`, only the
   one with the greatest `run_number` is kept, at the place of the first one
   seen. `ID` is the decimal `id`, `Status` and `Conclusion` are lower-cased,
   `Attempt` is `run_attempt` and `URL` is `html_url`. `ExpectedWorkflows` is
   `ExpectedWorkflows(files, BaseBranch)`.
2. API Contract: Expected Workflows. `ExpectedWorkflows` returns, sorted, each
   path whose YAML document is a mapping with an `on` key that triggers on
   every pull request into `base`:
   - `on` is the string `pull_request`, `pull_request_target` or `push`, or a
     sequence containing one of them;
   - or `on` is a mapping with `pull_request` or `pull_request_target` whose
     value is empty or a mapping whose only keys are `branches` (a list or
     string that contains `base` exactly and no entry starting with `!`) and
     `types` (containing both `opened` and `synchronize`);
   - or `on` is a mapping with `push` whose value is empty.

   Any other key under the trigger (`paths`, `paths-ignore`,
   `branches-ignore`, `tags`, …), an unparsable file and a file without `on`
   make no expectation.
3. API Contract: the Head Check Set. `evaluateHeadCheckSet` walks
   `report.Checks` in order: each is a member; its label is `<workflow> /
   <name>`, or `<name>` when `workflow` is empty; a `pass` or `skipping`
   bucket passes; a `fail`, `cancel` or `cancelled` bucket is pending as
   `<label> (re-run requested)` when `rerunAttempts` holds the run its `link`
   names, and is failed otherwise; any other bucket is pending with its
   label. It then walks `report.WorkflowRuns`: a run whose `ID` no check
   `link` names (`checkRunID`) is a member; a run with a pending check adds
   nothing more; otherwise a run not `completed` is pending as `<name>
   (<status>)`, and a completed run whose conclusion is not `success`,
   `neutral` or `skipped` is pending as `<name> (re-run requested)` when
   `rerunAttempts` holds its `ID` with an attempt at least its `Attempt`, and
   otherwise, unless one of its checks is failed, counts in `RunFailures`.
   Each Expected Workflow that no run's `Path` equals is a member, pending as
   `<path> (no run yet)`. With no member at all, `Pending` is
   `[no check reported]`. `<name>` is the run's `name`, or its `path` when the
   name is empty.
4. API Contract: the completion rule in `checkCandidate`. The existing loop
   over `report.Checks` is unchanged except that each park it makes with
   `checks-failed` uses `checksFailedBlocker(report)`. Then the stage
   evaluates the set: `RunFailures > 0` parks `checksFailedBlocker(report)`;
   pending members keep waiting; a merge state `BLOCKED` or `UNKNOWN` or
   `Mergeable` `UNKNOWN` keeps waiting; otherwise the item moves to `merging`.
   `checksFailedBlocker` is `checks-failed: ` followed by, joined by `; `,
   each check in a `fail`, `cancel` or `cancelled` bucket as `<label>
   (<state>) <link>` (`<state>` lower-cased, or the bucket when `state` is
   empty; ` <link>` omitted when empty) and then each completed run with a
   failing conclusion that no failed check names as `<name> (<conclusion>)
   <url>`; with none, the bare `checks-failed`.
5. API Contract: the log and the deadline. After the merge state line, which
   is unchanged, the stage writes:
   - `roundfix: checks: Delivery Queue item <slug>: waiting for <n> check(s): <member>; <member>`
     each time the pending list differs from the last one written;
   - `roundfix: checks: Delivery Queue item <slug>: <m> check(s) completed for <head>`
     right before it moves to `merging`, `<m>` the member count;
   - `roundfix: checks: Delivery Queue item <slug>: read failed: <error>` when
     a read fails with an error text different from the last one written; a
     successful read resets it.

   At the deadline the blocker is `checksTimeoutBlocker(timeout, pending,
   state)`: the current text, then `; pending: <member>; <member>` when the
   last successful read had pending members, else `; merge state <state>` when
   it waited only for a merge state `BLOCKED` or `UNKNOWN` (`UNKNOWN` also for
   `Mergeable` `UNKNOWN`), else nothing.
6. API Contract: the Turn. `checkCandidate` calls `checkTurn` with the
   pending list each time it writes the waiting line and with `nil` when it
   returns; `checkTurn` keeps the list per slug under `turnMu` and records the
   turn, and does nothing without a recorder. `recordTurnLocked` sets
   `ItemTurn.Checks` for the advancing item from that list. `TurnLines`
   prints, for an `advancing` entry whose item is in `checking` with
   `Checks`, `Turn: <slug> advancing at checking; waiting for <n> check(s):
   <member>; <member>`; every other form is unchanged.

### Surface Transcripts

1. Surface Transcript: `deliver status` for a live owner whose Turn Record
   lists two pending members for an item in `checking` (controlled
   variations: none; the item worktree is `/worktrees/0149-example`):

   ```transcript
   $ roundfix deliver status
   stdout:
   0149-example	checking	-	/worktrees/0149-example
   Turn: 0149-example advancing at checking; waiting for 2 check(s): CI | Backend / Lint, Type Check, Test & Build; .github/workflows/ci-frontends.yml (no run yet)
   Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
   Usage: no Runs recorded
   stderr:
   exit: 0
   ```

## Invariants

1. The item never moves from `checking` to `merging` while a member of the
   Head Check Set is pending, an Expected Workflow has no run, or the set is
   empty, whatever the merge state.
2. Every `checks-failed` park made in `checking` names each failed member of
   that read; `ClassifyPark` still gives it the `finding` class, and
   `checks-timeout` keeps the `environment` class.
3. The ADR-0236 stale and outside-change re-runs are made exactly as today,
   and a run the queue re-ran is pending, not failed, until a newer attempt
   reports.
4. No Run Database schema, Delivery Stage, blocker name, Park Class,
   configuration key or exit code changes, and no check suite is read.
5. The engine with a boundary that returns no workflow run and no Expected
   Workflow behaves as today except for the blocker and log text in Declared
   breaks.

## Declared breaks

- task_02 adds two commands to `GitHubCLI.CurrentHeadChecks`, so every test
  that scripts its commands on the default branch when the Task runs gains
  the `gh api` and `git ls-tree` steps of API Contract 1, answering
  `{"total_count":0,"workflow_runs":[]}` and an empty listing, and nothing
  else changes: at `6ceb7ff5` `TestPullRequestBoundaryReportsCurrentHeadChecks`,
  `TestNoChecksReportedIsAnEmptyReport`,
  `TestPullRequestBoundaryReportsFailedCurrentHeadChecks` and
  `TestPullRequestBoundaryRejectsChecksFromAMovedHead`
  (`internal/delivery/github_test.go`), `TestPullRequestReadCarriesTheMergeState`
  (`internal/delivery/merge_state_test.go`), and the fake `gh` of
  `TestCommandDeliveryEngineRecordsThePullRequest`
  (`internal/cli/deliver_record_pull_request_test.go`), which answers `api`
  with the empty runs payload.
- task_03 changes the `checks-failed` and `checks-timeout` text and adds log
  lines, so every test that pins them on the default branch when the Task
  runs changes only that expected text: at `6ceb7ff5`
  `TestAFailedCheckInAChangedPackageParksWithoutARerun`,
  `TestAnUnattributableCheckFailureParksWithoutARerun`,
  `TestACheckPastItsFirstAttemptIsNotRerun`,
  `TestNilCheckRecoveryKeepsFailedCheckBehavior`,
  `TestCancelledChecksAreNotRerun`, `TestCheckInspectionErrorsDoNotRerun`,
  `TestCheckRerunErrorsParkTheItem`,
  `TestSecondCheckFailureInChangedPackageIsNotCalledFlaky` and
  `TestAnOldCheckAttemptCannotKeepResettingTheTimeout`
  (`check_rerun_test.go`), `TestAStaleCheckWhoseRerunIsRefusedParksChecksFailed`
  and `TestStaleRerunsAreBoundedPerTipAndRestartTheTimeoutOnlyOnce`
  (`stale_check_test.go`), where the timeout names `CI / verify (re-run
  requested)`, `TestTheCheckWaitDefaultsToSixtyMinutes`,
  `TestTheCheckWaitUsesTheConfiguredDuration` and
  `TestEachCheckWaitStartsItsOwnDeadline` (`check_wait_test.go`),
  `TestDeliveryEngineParksDeclaredBlockersAndContinues` and
  `TestChecksTimeoutParksTheItem` (`engine_test.go`),
  `TestABlockedMergeStateParksAtTheChecksTimeout` (`merge_state_test.go`) and
  `TestNoChangedPremiseRecordsNoWarning` and
  `TestAChangedPremiseWarnsAndContinues` (`revalidate_test.go`), all in
  `internal/delivery`. No poll count, elapsed time, stage or re-run count in
  them changes.

## ADR-0369

task_01 creates `docs/adr/0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md`
with exactly this content, apart from line wrapping:

```markdown
---
status: accepted
created_at: 2026-11-06T12:00:00Z
updated_at: 2026-11-06T12:00:00Z
deprecated_at: null
superseded_by: null
---

# The checking stage waits for the Head Check Set and parks on any failed check

A Delivery Queue item in `checking` moves to `merging` only when every member
of its Head Check Set is completed and none failed: every check run and commit
status the pull request head lists, the latest GitHub Actions workflow run per
workflow file and event of the head commit, and every Expected Workflow of the
head's own workflow files. A merge state of `CLEAN` alone never moves it.

## Context

On 2026-11-05 the Delivery Queue of Fluxus merged pull requests #466 and #469
with the `CI | Backend` job about to fail, and #477 before its CI ran. The
queue read every check `gh pr checks` listed, required or not, and parked on a
listed failure, but the three workflow runs of each head existed seconds after
the pull request was opened while their job check runs appeared more than a
minute later; the list held one passing check, and the merge state was
`CLEAN` because the base branch requires no check. The maintainer chose:
"na etapa de checks, a fila espera todos os workflows que o head disparou
terminarem, até o prazo de sempre. Se algum falhar, mesmo sem ser obrigatório,
o item para com o nome do check e o link."

## Decision

**The set.** At each poll the queue reads the head's check runs and commit
statuses as before, the workflow runs of the head commit through `gh api
repos/{owner}/{repo}/actions/runs?head_sha=<head>`, keeping the latest run per
workflow file and event, and the head's workflow files through local `git`.
A check run that belongs to a listed run counts once. Check suites are not
read: a suite with no check run never completes.

**Expected Workflows close the first seconds.** A workflow file whose trigger
makes GitHub run it for every pull request into the base branch is expected:
`pull_request` or `pull_request_target` without a path filter or
`branches-ignore`, whose `branches` names the base branch literally without a
`!` pattern and whose `types` name `opened` and `synchronize`, when present;
or a `push` trigger with no filter. Until such a file has a run for the head,
the set is not complete. A filtered workflow is never expected, but its run is
a member once GitHub creates it.

**Conclusions.** `success`, `neutral` and `skipped` pass. Any other
conclusion fails: `failure`, `timed_out`, `cancelled`, `action_required`,
which never completes without a person, `startup_failure` and `stale`, which
did not test the head. An empty set keeps waiting.

**Parks keep their names.** A failed check still goes through the stale and
outside-change re-runs of ADR-0236, and a run the queue re-ran is pending
until a newer attempt reports. Every other failure parks `checks-failed:
<check> (<conclusion>) <link>`, each failed member joined by `; `. The deadline
stays `delivery.check_timeout` and parks `checks-timeout: current-head checks
not settled within delivery.check_timeout <duration>; pending: <check>;
<check>`. Both keep their Park Classes.

**The wait is visible.** The console log names the pending members when they
change and the member count when the set completes. The Queue Owner records
the pending members on the item's Queue Turn, and `roundfix deliver status`
prints `Turn: <slug> advancing at checking; waiting for <n> check(s): <check>;
<check>`.

## Considered Options

- A settle delay before trusting the list was rejected: it adds a fixed wait
  to every merge and still guesses how long GitHub takes.
- Evaluating `paths`, branch globs and `branches-ignore` to expect filtered
  workflows was rejected: it re-implements GitHub's filter semantics, and a
  wrong guess waits for a run that never comes until the Check Wait ends.
- Waiting for every check suite was rejected: an app can create a suite that
  never gets a check run.
- Requiring only the branch's required checks, or warning about failed
  optional checks while still merging, was rejected by the maintainer's
  answer.
- Reading `UNSTABLE` from the merge state was rejected: GitHub sets it only
  for checks that already reported.

## Consequences

A pull request now merges no earlier than every workflow its head triggered,
which costs nothing when CI already takes minutes. A repository whose only
workflows are filtered still relies on GitHub creating their runs within the
first poll interval. An operator whose `gh` token cannot read Actions sees a
logged read error and a `checks-timeout` park. A blocker recorded before this
change keeps its bare text; status and retry read both forms.
```

### Amendment notices

task_01 writes, with `updated_at: 2026-11-06T12:00:00Z`:

- in ADR-0311, after the paragraph that starts `**The park names what it
  waited.**`, the paragraph

  ```markdown
  **Amended by [ADR-0369](0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md) (2026-11-06).** The Check Wait of an item in `checking` waits for its Head Check Set: every check run and commit status of the head, the latest workflow run per workflow file and event of the head commit, and every Expected Workflow. The `checks-timeout` blocker adds `; pending: <check>; <check>` for the members still pending, or `; merge state <state>` when only the merge state was awaited. The key, its default, its bounds, its User Config scope and the one restart after a re-run still hold.
  ```

- in ADR-0365, after the paragraph that starts `**The turn is a file.**`,
  the paragraph

  ```markdown
  **Amended by [ADR-0369](0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md) (2026-11-06).** The Turn Record of the item in `checking` also lists the members of its Head Check Set still pending, and its `Turn:` line reads `advancing at checking; waiting for <n> check(s): <check>; <check>`. Every other turn form, the file, its readers and the Review Capacity still hold.
  ```

## Coverage Map

- Goal 1, Core Features 1–3, Success Metrics 1–2 → API Contract 4, with API Contracts 1 to 3: the reader of `CurrentHeadChecks` (task_02) collects the workflow runs and Expected Workflows of the head, and `evaluateHeadCheckSet` with the completion rule of `checkCandidate` (task_03) keeps the item in `checking` until every member is completed.
- Goal 2, Core Feature 4, Success Metric 3 → API Contract 4: `checksFailedBlocker` and the conclusion rule (task_03) park `checks-failed` naming each failed member, the re-run path stays as it is, and the `run_number` rule of API Contract 1 (task_02) drops a superseded run.
- Goal 3 → API Contract 2, with API Contracts 1 and 3: the Expected Workflows and the runs read (task_02) and the `no run yet` and `no check reported` members (task_03) keep a head whose workflows have not registered in `checking`.
- Goal 4, Core Feature 5, Success Metric 4 → API Contract 5, with API Contract 6 and Surface Transcript 1: the log lines and the `checks-timeout` text (task_03) and the Turn Record and `Turn:` line (task_04) name the pending members.
- Goal 5, Core Feature 6, Success Metric 5 → Amendment notices, with ADR-0369 and Build Order 1: task_01 writes ADR-0369, the notices in ADR-0311 and ADR-0365, the four `CONTEXT.md` entries, the `deliver` and configuration guides and both copies of the Roundfix `deliver` reference, each stating the Head Check Set rule.

Deviations from the PRD: None.

## Absence Evidence

1. The workflow-run read and `WorkflowRun` — command: `grep -rn 'actions/runs\|WorkflowRun\|head_sha' internal/delivery`; result: no match.
2. `ExpectedWorkflows` and any read of `.github/workflows` by the queue — command: `grep -rn 'github/workflows\|ExpectedWorkflow\|yaml' internal/delivery`; result: no match.
3. The Head Check Set evaluation and a named `checks-failed` text — command: `grep -n 'BlockerChecksFailed' internal/delivery/engine.go`; result: six `engine.park(ctx, gitRoot, item, BlockerChecksFailed)` calls in `checkCandidate`, all bare.
4. The waiting and completed log lines — command: `grep -n 'roundfix: checks:' internal/delivery/engine.go`; result: only the `merge state` line.
5. Pending checks on the Turn Record — command: `grep -n 'Checks' internal/delivery/turn.go`; result: no match.
6. ADR-0369 and the glossary terms — command: `ls docs/adr | grep '^0369'; grep -n 'Head Check Set\|Expected Workflow' CONTEXT.md`; result: no match.

## ADR Alternatives

- ADR-0311: it rejected a constant, a Project Config key, `0` as no limit
  and deriving the wait from `timeout-minutes`; this design takes none of
  them and keeps the key, and its decision to parse no workflow for the wait
  holds, because the Expected Workflows decide what to wait for, never how
  long.
- ADR-0236: it rejected reopening the pull request, checking out the dynamic
  merge ref and comparing start times; this design takes none of them and
  keeps its re-run path.
- ADR-0365: it rejected a `waiting` stage or column, deriving the turn from
  the Run Database, a Project Config key and a default of `1`; this design
  takes none of them and adds the pending members to its file.
- ADR-0304: it rejected a distinct exit code, a `waiting` stage or column, a
  Project Config key and an unbounded wait; this design takes none of them.
- ADR-0199: it records no rejected alternative this design touches.
- ADR-0237: it records no rejected alternative this design touches; the
  retry's merge observation is unchanged.
- ADR-0369: written here; it takes none of its own rejected options.

## Integration Points

GitHub's REST API through the operator's `gh api`, one more request per poll
(every 15 seconds), well inside the authenticated rate limit; local `git
ls-tree` and `git show` in the item worktree. No test reaches GitHub: the
reader is tested with the scripted command runner and the engine with fake
boundaries.

## Testing Approach

The reader is tested with `newScriptedCommandRunner`, which fails on any
command out of order. The engine tests build on `newCheckWaitFixture` and
`recoveryReport` with the fake clock and sleeper, so no test waits in real
time; the Fluxus replays return a scripted report per poll. The `Turn:` line
is tested through `TurnLines`, through `checkTurn` with a recording fake, and
through `deliver status` with a Turn Record written at `turnRecordPath`, as
`deliver_turn_test.go` does.

Sabotages, each on a temporary copy, with the test and the input that sees
it:

1. Drop the workflow runs from the set (treat every run as passed):
   `TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered` fails on
   its first poll, where one passing check is listed beside two runs in
   progress; the item reaches `merging`.
2. Drop the Expected Workflows from the set:
   `TestAnExpectedWorkflowWithoutARunKeepsTheItemWaiting` fails on its report
   with one passing commit status and one Expected Workflow without a run.
3. Treat `action_required` as passing:
   `TestAFailedWorkflowRunWithoutAFailingCheckParks` fails on its
   `action_required` case.
4. Accept a `paths` filter as expected:
   `TestExpectedWorkflowsFollowTheTriggerRule` fails on its `paths` case.
5. Keep the first run of a workflow file instead of the greatest
   `run_number`: `TestCurrentHeadChecksReadsTheWorkflowRunsOfTheHead` fails on
   its superseded `cancelled` run 9 listed before run 10.
6. Park bare `checks-failed`:
   `TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered` fails on its
   failure poll.
7. Print `advancing at checking` without the checks:
   `TestDeliverStatusPrintsTheChecksAnItemWaitsFor` fails.

## Build Order

1. ADR-0369, the notices in ADR-0311 and ADR-0365, `CONTEXT.md`, the
   `deliver` and configuration guides and the Roundfix `deliver` reference
   (task_01), first, so every later Task implements what the documents state.
2. The reader: workflow runs, Expected Workflows and `baseRefName` (task_02;
   depends on: 1, for the documents).
3. The completion rule, the blocker texts and the log lines (task_03;
   depends on: 2, for the `CheckReport` fields and `WorkflowRun`).
4. The pending members on the Turn Record and the `Turn:` line (task_04;
   depends on: 3, for the pending list and the shared `engine.go`).
5. QA gate (depends on: 4).

No `requires`. Spec 0360, delivered after this Spec, may change
`internal/delivery/park_class.go`, the Roundfix `deliver` reference and its
version record, the `deliver` guide, `CONTEXT.md` and the Behavior Surface
Record; this Spec does not change `park_class.go`, and the record commands
recompute versions and fingerprints for whichever Spec merges second.

## Risks & Considerations

- A token without Actions read access makes every read fail; the log names
  the error once and the item parks `checks-timeout`. Every adopter uses a
  `gh` login with the `repo` scope, which reads Actions.
- A head with more than 100 workflow runs fails the read; one head has a
  handful.
- A workflow file the default branch deleted after the item branched stays
  expected and waits until the deadline; the park names it.
- A workflow disabled in the repository settings stays expected and waits
  until the deadline; the park names it.
- `pull_request` workflows use the merge commit's files, not the head's; a
  workflow added only on the default branch is not expected but is a member
  once its run exists.
- One more API request per poll; the queue polls one item at a time.

## Glossary

- adds: **Head Check Set**
- adds: **Expected Workflow**
- changes: **Check Wait**
- changes: **Queue Turn**

## Decisions

- Workflow runs are read by head commit, not by branch or event, so push and
  pull request runs of the head both count. See ADR-0369.
- Only filter-free triggers, a literal base branch and default-compatible
  types create an expectation. See ADR-0369.
- Check suites are not read. See ADR-0369.
- The blocker names keep their Park Classes; only the detail grows. See
  ADR-0369.
- The pending members live on the Turn Record, not in the Run Database. See
  ADR-0365 and ADR-0369.
