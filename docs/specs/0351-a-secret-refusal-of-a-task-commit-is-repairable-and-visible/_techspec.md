---
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
prd: _prd.md
created: 2026-11-04
---

# A secret refusal of a Task commit is repairable and visible — Technical Spec

## Executive Summary

ADR-0292 scans each Task commit after the Task's last Verification attempt,
so a Secret Match can only fail the Task, and the refusal reaches neither
`roundfix runs causes` nor the Delivery Queue's park. The design moves the
scan into the attempt as the last Settlement Check, `settlement check: secret
scan`, run for every non-QA Task with or without a QA gate. The attempt's
existing machinery then carries a first-attempt match to the one Verification
Feedback turn and rescans on the second attempt; a final match keeps today's
reason bytes. The Causes report turns each `secret_scan` event into an item
of kind `secret-scan`, and the Implement Run result handed to the Delivery
Engine carries the failed and skipped Tasks so a `run-unresolved` park names
them. The skill, the guides, `CONTEXT.md` and ADR-0361 state the rules.

The trade-off accepted: every attempt now prepares the Task commit once more
(one worktree snapshot and one `git diff` of the prepared paths), and a Task
that fails for another reason is also scanned, so its single repair turn
receives every problem at once.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the new
  check label and item kind are fixed strings. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, endpoint or
  credential changes. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0292 owns the Task commit scan,
  ADR-0292: "a new repair class is a separate decision"; ADR-0273 keeps every
  secret from leaving the machine from a Run, ADR-0273: "it never carries the matched value or
  any part of it"; ADR-0182 and ADR-0342 own the Settlement Checks, ADR-0342:
  "The two in-process Settlement Checks, Spec Consistency and the
  authorization audit, stay scoped to a graph with a QA gate Task"; ADR-0214
  owns cause classification, ADR-0214: "classified by deterministic
  signatures"; ADR-0361 is new. Set no obligation here, because this Spec
  changes none of their rules and each reaches it only through citations of
  the ADRs above: ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0196 (a pre-PR review finding parks only after validation); ADR-0197 (a pre-PR reviewer lineage spans at most two rounds); ADR-0227 (a review through Claude keeps its verdict and fits its window); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0256 (a release waits for the skills that describe a changed surface); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0272 (a review record survives a fix, and a correction inside the archived Spec's declared paths returns to review); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0283 (a Run Marker keeps the judge silent inside a Run); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix); ADR-0290 (merged items give back their Item Binaries, and Run Worktrees leave only by proof); ADR-0291 (qA evidence holds text and images, and the QA Report commit refuses any other binary); ADR-0295 (reachability reads the fetched default branch, and origin names it when origin/HEAD is unset); ADR-0299 (the history conversion rewrites every citation it lists, and the check reports a removed history file cited without a Git Locator); ADR-0306 (archive and cancel refuse a path whose resolved target leaves the repository); ADR-0308 (an archive rebases the unproven links it records, and the Record Link Repair points a promoted target at its copy); ADR-0317 (a Result Status Note of a completed Task is dismissed without the validator); ADR-0319 (the follow-up order is code over four Jev yes/no answers); ADR-0345 (a derived-record writer names its command to the suite guard); ADR-0348 (an absent Spec Root inside the repository holds no Spec); ADR-0349 (a Spec's Release Ceiling answers the release plan's approval); ADR-0352 (the pre-PR review omits declared paths, and a large Claude candidate is reviewed in segments); ADR-0353 (a waiting item is held back from the follow-up list unless --all is asked). ADR Alternatives below relates each. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AY answers of 2026-11-04 quoted in the authorization record; the Governed
  Paths are the Roundfix skill files task_04 changes; no Makefile, workflow,
  module, lint or Baseline change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/runs.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-04 at `de579ea2`; `_prd.md` → Validation records the runs.

1. `executeTask` (`internal/daemon/task_engine.go`) runs `verifyTask`
   attempt 1, `repairTaskVerification` on a deterministic failure, attempt 2,
   then `prepareTaskCommit`, the lost-path and authorization refusals, then
   `scanTaskCommitPaths`. A match sets `failure` to `Task commit carries a
   secret: ` plus `Match.String()` joined by `; `, publishes a
   `daemon.task` event with summary `Task <id> commit refused: <n> Secret
   Match(es)` and payload `phase: secret_scan`, `matches: [{path, line,
   kind}]`, prints the summary, and settles the Task `failed`. A scan error
   returns from `executeTask` and halts the Run.
2. `verifyTask` sets `request.Checks` to `settlement check: spec consistency`
   and `settlement check: authorization` only when `plan.settlementChecks`
   (a graph with a QA gate Task); a Task with undeleted `deletes` paths
   replaces the checks with `settlement check: deletes` and no commands.
   `runVerificationAttempt` (`internal/daemon/engine.go`) runs every check
   after the commands whatever their result, publishes `started`,
   `command-passed` or `failed` for each, writes
   `<attempt>-settlement-<name>.log`, and turns a checker error into a
   failure diagnostic. `verificationTerminalReason` renders a check failure
   as `Settlement check failed: <label>: <first line>; diagnostics: <path>`.
3. `repairTaskVerification` builds one Verification Feedback prompt from every
   failure of attempt 1, each with `Failed command`, `Diagnostic artifact` and
   `Failure: verification failed: <error>`, where a check's error carries its
   whole diagnostic.
4. `ProjectStreamEvent` (`internal/runevent/stream.go`) discloses the command
   of a verification event only for the two in-process labels.
5. `runcause.Build` (`internal/runcause/report.go`) reads `daemon.verification`
   and `daemon.task` events of terminal Spec Runs, makes one `verification`
   item per failed attempt and one `corrective` item per corrective Task, and
   ignores `secret_scan` events.
6. `commandDeliveryWorkflow.runResult` (`internal/cli/deliver_workflow.go`)
   reads the `daemon.task` events of an unresolved Run only for exhausted
   runtime recoveries; `runCandidate` (`internal/delivery/engine.go`) then
   parks the bare `run-unresolved`, and `ClassifyPark`
   (`internal/delivery/park_class.go`) answers it with class `finding` and
   "resolve the blocker, then run roundfix deliver retry <slug>". A skipped
   Task publishes `phase: skipped` with `reason: needs not completed: <ids>`;
   a settled Task publishes `phase: settled`, `status` and `reason`.

## System Architecture

No package or command is added. Changes by component:

- `internal/runevent`: the label constant and its stream projection.
- `internal/daemon`: a new `settlement_secret_scan.go` holding the check; the
  check list in `verifyTask`; the halting checker error and the final reason
  in `engine.go`; the post-attempt scan removed from `executeTask`.
- `internal/runcause`: the `secret-scan` item.
- `internal/delivery` and `internal/cli/deliver_workflow.go`: the failed and
  skipped Tasks of an unresolved Run and the park detail.
- The Roundfix skill, three user guides, `CONTEXT.md` and ADR-0361.

## Implementation Design

### Interfaces

```go
// internal/runevent/event.go
const SettlementCheckSecretScan = "settlement check: secret scan"

// internal/daemon/engine.go
type verificationCheck struct {
	Label       string
	Run         func(context.Context, string) (string, error)
	HaltOnError bool // a checker error halts the Run instead of failing the attempt
}

// internal/daemon/settlement_secret_scan.go
func (engine *Engine) secretScanSettlementCheck(plan TaskPlan, task spec.Task, ordinal, attempt int, before []string) verificationCheck

// internal/runcause/report.go
type Match struct { Path string `json:"path"`; Line int `json:"line"`; Kind string `json:"kind"` }
// Item gains: Matches []Match `json:"matches"`

// internal/delivery/engine.go
type TaskFailure struct{ Task, Reason string }
// RunResult gains: FailedTasks []TaskFailure; SkippedTasks []string
func RunUnresolvedBlocker(result RunResult) string
```

### Data Models

None in the Run Database. The `secret_scan` event payload gains `attempt`.
A Causes item gains `matches`, `null` on every item that is not `secret-scan`;
the schema name stays `roundfix/runs-causes/v1`. The stored blocker of a
`run-unresolved` item may carry a detail after `run-unresolved: `.

### API Contracts

1. API Contract: the check. `secretScanSettlementCheck` returns a check with
   label `runevent.SettlementCheckSecretScan` and `HaltOnError: true`. Its Run
   calls `prepareTaskCommit(ctx, plan, task, before)` and
   `scanTaskCommitPaths(ctx, plan.WorkDir, prepared.stageable)`, both
   unchanged. With no match it returns `""`. With matches it publishes the
   `daemon.task` event of Current behavior 1 with payload `phase:
   secret_scan`, `attempt: <attempt>`, `matches: [{path, line, kind}]`,
   prints its summary to progress, and returns the two-line diagnostic:
   `Task commit carries a secret: <Match.String() joined by "; ">` and
   `Fix: replace each named value with one no Secret Pattern matches, such as letters only or fewer than 20 characters; copying a value the repository already tracks is refused the same way.`
2. API Contract: the order. `verifyTask` keeps the two in-process checks under
   `plan.settlementChecks` and then always appends the secret scan check, so
   an attempt runs its commands (the repository Verification last), `settlement
   check: spec consistency`, `settlement check: authorization`, then
   `settlement check: secret scan`; without a QA gate only the scan follows
   the commands. The `deletes` replacement is unchanged and runs no scan.
3. API Contract: the halting error. In `runVerificationAttempt`, a check whose
   `HaltOnError` is set and whose Run returns a non-stop error publishes the
   infrastructure failure for its label and returns
   `run verification attempt <n>: <error>`; the Run halts as a scan error
   does today. Other checks keep turning an error into a failure diagnostic.
4. API Contract: the final reason. `verificationTerminalReason` renders a
   failure of `runevent.SettlementCheckSecretScan` as its first diagnostic
   line only, so a final attempt whose only failure is the scan settles with
   `Task commit carries a secret: <locations>`, today's bytes; with other
   failures `verificationAttemptFailureReason` joins them with `; `, the scan
   last. `executeTask` no longer scans after `prepareTaskCommit`.
5. API Contract: the stream. `ProjectStreamEvent` also discloses the command
   `settlement check: secret scan`; every other command stays redacted.
6. API Contract: the Causes item. In `runcause.Build`, a `daemon.task` event
   with `phase: secret_scan` adds one item: `kind: secret-scan`, the Run's
   Spec, the payload `task` or else the event's Work Item, the Run ID,
   `attempt` from the payload or `null` when absent, `check` the locations
   `<path>:<line>: <kind>` joined by `; ` (through `publicCheck`), `class:
   repository-convention`, `signature` and `evidence` `null`, and `matches`.
   A failed `daemon.verification` event whose command is the scan label marks
   its attempt seen and adds no `verification` item. Text keeps its line shape:
   `secret-scan <spec>/<task> <run> attempt-<n> repository-convention -: <locations>`
   (`-` for an absent attempt). The signature table is unchanged.
7. API Contract: the Run result. `runResult` reads, from the same
   `daemon.task` events, each `phase: settled` with `status: failed` as a
   `TaskFailure{Task, Reason}` and each `phase: skipped` as a skipped Task,
   in event order, the Task taken from the payload `task` or else the event's
   Work Item.
8. API Contract: the park. `runCandidate` parks with
   `RunUnresolvedBlocker(result)` where it parked `BlockerRunUnresolved`. The
   function joins with `; ` one `<task> failed: <reason>` per failed Task (the
   reason whitespace-normalized and, past 300 characters, cut to 300 followed
   by `…`) and, when any Task was skipped, `skipped: <ids joined by ", ">`, and
   returns `run-unresolved: <parts>`, or `run-unresolved` when there are no
   parts. Runtime infrastructure and the QA environment partial keep
   precedence.
9. API Contract: the next action. `ClassifyPark` keeps class `finding` for
   every `run-unresolved` blocker. When its detail contains
   ` failed: Task commit carries a secret: `, Next is
   `if the failed Task's text prescribes the value the blocker locates, amend it in <worktree> to one no Secret Pattern matches (letters only, or fewer than 20 characters); then run roundfix deliver retry <slug>, whose rerun returns any Secret Match to the Agent as Verification Feedback`;
   every other `run-unresolved` keeps "resolve the blocker, then run roundfix
   deliver retry <slug>".
10. API Contract: the documents. task_03 writes the park paragraph into the
    deliver user guide; task_04 writes the sentences its requirements
    prescribe into the Roundfix skill's `### Settlement Checks` section and
    its implement, runs and deliver references, the implement and runs user
    guides, `CONTEXT.md` (Settlement Check, Secret Scan, Cause Class) and
    ADR-0361. The skill's entry file stays within its 20000-byte budget.

### Surface Transcripts

None. `roundfix deliver status` and `roundfix runs causes` change only for
rows a Run writes into the Run Database, so no command of a bare checkout
shows the change; the named tests of task_02 and task_03 assert their whole
standard output for seeded databases.

## Invariants

```text
1. No value, nor any eight consecutive bytes of one, reaches a prompt, a diagnostic file, an event, a progress line, a Causes item or a blocker.
2. A Task has at most two Verification attempts and one repair turn, whatever the scan finds.
3. A Task commit is created only after an attempt whose secret scan found no match on the paths that commit stages.
4. A final attempt whose only failure is the scan settles with today's reason bytes and creates no commit.
5. A scan error halts the Run and never becomes Verification Feedback.
6. A run-unresolved park without failed or skipped Tasks keeps the bare blocker and today's answer; every run-unresolved blocker keeps class finding.
```

## Declared breaks

- `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit`
  (`internal/daemon/task_secret_scan_test.go`): one refusal becomes two, with
  attempts 1 and 2, after one repair turn; the reason bytes are unchanged.
- `TestTaskCycleVerificationSequenceStopsAtFirstFailureWithOneVerdictPerAttempt`
  (`internal/daemon/task_engine_test.go`): each attempt's phases gain
  `started|command-passed` for the scan before its verdict.
- `TestTaskCycleExecutesAgentVerifySettleCommitContract` (same file): the
  scripted snapshots gain one per Task for the scan's commit preparation:
  `{nil, {"src/one.go"}, {"src/one.go"}, nil, {"src/two.go"}}`.
- `TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree`
  (`internal/cli/implement_test.go`): the exclusive retry has two `started`
  events, the command's and the scan's.
- `TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot`
  (`internal/runevent/settlement_label_projection_test.go`) gains the third
  label; its existing cases are unchanged.
- Public: every attempt shows the scan in `roundfix events`; a Task that fails
  another check is scanned too; `secret_scan` gains `attempt`; Causes JSON
  items gain `matches`; a `run-unresolved` blocker, the Pending Question and
  the `deliver retry` line may carry the detail; a light-tier Task escalates
  on a first-attempt Secret Match as on any Verification failure.

## ADR-0361

task_04 creates
`docs/adr/0361-a-task-commit-secret-match-returns-as-verification-feedback-and-names-its-cause.md`
with the lifecycle overlay (`status: accepted`, `created_at` and `updated_at`
`2026-11-05T00:00:00Z`, `deprecated_at: null`, `superseded_by: null`), the
title "A Task commit's Secret Match returns as Verification Feedback and
names its cause", and:

- Context: ADR-0292 scanned after the last attempt and deferred a repair
  turn; on 2026-11-04 Fiscus lost a Run to a `credential-assignment` fixture
  that no report, park line or Agent turn named.
- Decision: the scan is the last Settlement Check of every attempt, with or
  without a QA gate; a first-attempt match returns as Verification Feedback
  with locations and a `Fix:` line, never the value; the second attempt
  rescans; a final match keeps the reason; a scan error halts the Run; the
  Causes report lists a `secret-scan` item classed `repository-convention`
  by kind; a `run-unresolved` park names failed Tasks, reasons and skipped
  Tasks.
- Considered Options, each rejected with its reason: a second repair turn
  only for the scan; scanning only on attempt 1; keeping the post-attempt scan
  as a second pass; accepting a value tracked files already hold; a new
  signature in the cause table; a new Park Class or blocker name; a
  structured blocker column.
- Consequences: amends ADR-0292 (its rejected repair turn is now taken) and
  ADR-0214 (one item kind is classed by kind); ADR-0273's value rule extends
  to the feedback, the report and the park.

## Coverage Map

- Goal 1, Goal 2, Core Features 1-3, Success Metrics 1-3 → the check and its
  order, halting error and reason (API Contracts 1-4), Invariants 2-5,
  `TestTaskCommitSecretMatchReturnsAsVerificationFeedback`,
  `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit`,
  `TestTaskCommitSecretScanRunsAfterTheSettlementChecks`,
  `TestTaskSecretScanReadErrorHaltsTheRun`.
- Goal 3, Core Feature 4, Success Metric 4 → the stream and the Causes item
  (API Contracts 5 and 6), `TestBuildReportsASecretRefusalAsItsOwnItem`,
  `TestRunsCausesPrintsASecretRefusal`.
- Goal 4, Core Feature 5, Success Metric 5 → the Run result, the park and the
  next action (API Contracts 7-9), Invariant 6,
  `TestDeliveryRunResultReadsFailedAndSkippedTasks`,
  `TestARunUnresolvedParkNamesItsFailedAndSkippedTasks`,
  `TestDeliverStatusNamesTheSecretRefusalOfARunUnresolvedPark`.
- Goal 5 → Invariant 1 and the value checks of the tests above.
- Core Feature 6, Success Metric 6 → the documents (API Contract 10), written
  by task_03 for the deliver guide and by task_04 for the rest.

Deviations from the PRD: None.

## Absence Evidence

1. The scan as an attempt check — command: `grep -n 'secret scan\|secretScanSettlementCheck\|HaltOnError' internal/daemon/*.go internal/runevent/*.go`; result: no match; the scan runs only in `executeTask` after the attempts.
2. A Causes item for a secret refusal — command: `grep -n 'secret' internal/runcause/report.go internal/runcause/signatures.json`; result: no match.
3. The failed and skipped Tasks of an unresolved Run — command: `grep -n 'FailedTasks\|SkippedTasks\|RunUnresolvedBlocker\|"skipped"' internal/delivery/engine.go internal/cli/deliver_workflow.go`; result: no match.
4. The secret next action of a run-unresolved park — command: `grep -n 'carries a secret' internal/delivery/park_class.go`; result: no match.

## ADR Alternatives

- ADR-0292: amended; it rejected "Giving the Agent a repair turn after a Task
  Secret Match" for that change only, and this Spec takes it; its other
  rejected alternatives (exempting test files or `test`/`fixture` values,
  scanning whole files) stay rejected, and the scan keeps reading added lines.
- ADR-0273: followed; its rejected alternative of redacting a matched value
  instead of refusing stays rejected: the scan still refuses, and the new
  feedback, item and blocker carry locations only.
- ADR-0182: followed; its rejected alternative of leaving commit-dependent
  facts to the gate stays rejected, and the scan is one more fact checked on
  the commit the Daemon is about to create.
- ADR-0342: extended; the two in-process checks stay scoped to a gated graph,
  and the scan, like the repository Verification, runs in every graph.
- ADR-0214: amended for one kind; its rejected alternative of classifying by a
  model judgment stays rejected, because the `secret-scan` class is fixed by
  the event's kind, and the signature table and its digest do not change.
- ADR-0092: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0120: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0165: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0169: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0196: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0197: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0227: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0233: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0247: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0248: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0250: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0251: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0254: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0256: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0258: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0260: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0261: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0262: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0269: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0272: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0277: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0278: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0279: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0281: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0282: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0283: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0286: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0290: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0291: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0295: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0299: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0306: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0308: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0317: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0319: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0345: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0348: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0349: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0352: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0353: unaffected; it reaches this Spec only through citations of the ADRs above, and this Spec neither takes nor revisits any alternative it rejected.
- ADR-0361: new; records the options above.

## Integration Points

None. The check reads the Run Worktree and Git as the scan does today; the
report and the park read the Run Database.

## Testing Approach

All tests are hermetic: temporary Git repositories, Run Databases and homes,
the fake Agent runner (`taskFakeRunner`), the fake verifier, no network. Every
new top-level test calls `t.Parallel()` first. Fixture values are built at
runtime as `strings.Repeat("K8", 15)` behind the name `SERVICE_API_KEY`, as the
existing Task scan tests do.

- `internal/daemon/task_secret_scan_test.go`:
  - `TestTaskCommitSecretMatchReturnsAsVerificationFeedback`: the Agent writes
    the value on its first turn and `strings.Repeat("z", 30)` on its repair
    turn; the Task completes, one repair prompt names the label, the location
    and the `Fix:` line, one `secret_scan` event has `attempt` 1, and Invariant
    1 holds over prompts, `*-settlement-secret-scan.log` files, events and
    progress.
  - `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit` (changed).
  - `TestTaskCommitSecretScanRunsAfterTheSettlementChecks`: subtests `gated
    graph` and `graph without a QA gate`, reading the `started` check labels of
    attempt 1.
  - `TestTaskSecretScanReadErrorHaltsTheRun`: an untracked output made
    unreadable with mode `0o000`; `TaskCycle` returns an error, no repair
    prompt is sent and no commit is created; it skips on Windows and when the
    effective user is root.
- `internal/runcause/report_test.go`: `TestBuildReportsASecretRefusalAsItsOwnItem`.
- `internal/cli/runs_causes_test.go`: `TestRunsCausesPrintsASecretRefusal`.
- `internal/delivery/run_unresolved_park_test.go` (new):
  `TestARunUnresolvedParkNamesItsFailedAndSkippedTasks`.
- `internal/cli/deliver_run_unresolved_park_test.go` (new):
  `TestDeliveryRunResultReadsFailedAndSkippedTasks` and
  `TestDeliverStatusNamesTheSecretRefusalOfARunUnresolvedPark`.
- Sabotages, each on a temporary copy, each with the case that observes it:
  - the scan check not appended: `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit`
    creates a commit;
  - the check appended only when `plan.settlementChecks`: the `graph without a
    QA gate` subtest and the feedback test fail;
  - `HaltOnError` ignored: `TestTaskSecretScanReadErrorHaltsTheRun` sees a
    repair prompt;
  - the special final reason removed: the changed refusal test sees
    `Settlement check failed:`;
  - the `Fix:` line dropped: the feedback test fails on the prompt;
  - the scan-check skip removed from `runcause.Build`: the report test counts a
    `verification` item for the scan;
  - skipped Tasks not read in `runResult`: the run result test fails on
    `SkippedTasks`;
  - the reason cut removed: the park test's 400-character reason case fails.

## Build Order

1. The scan as a Settlement Check, its stream label and its tests (task_01).
2. The Causes item and its tests (task_02) (depends on: 1, for
   `runevent.SettlementCheckSecretScan`).
3. The Run result, the park detail, the next action, their tests and the
   deliver guide paragraph (task_03) (no dependency; it shares no file with
   steps 1 and 2).
4. The skill, references, guides, `CONTEXT.md`, ADR-0361, the skill version
   and the Behavior Surface Record last (task_04) (depends on: 1, 2, 3).
5. QA gate (task_05) (depends on: 4).

No other Spec is active, so no file is shared with another plan.

## Risks & Considerations

- Each attempt prepares the Task commit once more; on a large worktree this is
  one extra status snapshot and one diff of the prepared paths per attempt.
- A Run Database written by 0.130.1 lists its secret refusal with a null
  attempt; nothing rewrites those rows.
- A blocker with detail is longer; the reason cut bounds each Task to 300
  characters.
- Secondbrain: `_prd.md` → Validation; the mirrors agree with ADR-0292.
- Exa: GitHub's push protection pages support location-only refusals; the
  fixture-allowlist article supports declining the tracked-value exemption.

## Glossary

None.

## Decisions

- The scan is the last check of every attempt. See ADR-0361.
- A scan error halts the Run (API Contract 3).
- The final reason keeps today's bytes when the scan fails alone (API
  Contract 4).
- The Causes item kind is `secret-scan`, classed by kind. See ADR-0361.
- The park detail comes from the Run's own events, 300 characters per reason.
  See ADR-0361.
