---
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
status: active
created: 2026-11-06
surfaces: [cli, docs]
type: fix
---

# The Delivery Queue waits for every check the head triggered

On 2026-11-05 the Delivery Queue of Fluxus merged pull requests #466 and #469 while the job `Lint, Type Check, Test & Build` of `CI | Backend` had not yet started on their heads, and both later failed, leaving the default branch red; it merged #477 the same way before `CI | Backend` and `CI | Frontend` had run, as reported by fluxus. The queue already reads every check `gh pr checks` lists, required or not, and parks on a listed failure, but it trusted a list that did not yet hold the jobs of workflow runs GitHub had created for the head more than a minute before those jobs registered their check runs, and it read the merge state `CLEAN` as the rest of the answer. After this fix the checking stage waits until every check run and commit status of the head, every GitHub Actions workflow run created for the head, and every workflow the head's own workflow files make GitHub start for its pull request are completed, within the existing Check Wait; it parks on any of them that ends in failure, naming each one with its link, and `deliver status` and the console log name the checks it still waits for.

## Sources

One source, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-05-fluxus-delivery-queue-merged-before-checks.md`](references/2026-11-05-fluxus-delivery-queue-merged-before-checks.md),
   reported by fluxus on 2026-11-05: two ai-memory messages, 11:05 BRT (pull
   requests #466 and #469 merged with `CI | Backend` in failure) and 12:35
   BRT (pull request #477 merged before its CI ran), copied verbatim.

The timing this Spec relies on was measured from GitHub, read only, on
2026-11-06 and is recorded in
[`references/2026-11-06-fluxus-check-timing.md`](references/2026-11-06-fluxus-check-timing.md).
No open Backlog Entry or unresolved Finding on `main` shares this context;
the only open entry, the repository build cache bound, belongs to Spec 0360.

## Prerequisites

None. Authored against `6ceb7ff5` (v0.136.0 plus a backlog commit), with no
active Spec on `main`. Spec 0360 is authored in parallel and delivers after
this Spec; it may change `internal/delivery/park_class.go`, the Roundfix
skill's `deliver` reference and its version record, the `deliver` guide,
`CONTEXT.md` and the Behavior Surface Record. This Spec does not change
`park_class.go`, keeps each edit to the shared documents narrow, and the
record commands recompute versions and fingerprints on whichever Spec merges
second.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; workflow
  runs keep GitHub's numeric ids and items keep their Spec slugs. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint or credential of
  Roundfix changes; the queue reads one more GitHub REST resource through the
  operator's authenticated `gh`, as it already does for check annotations.
  Source: `docs/agents/go.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0311 bounds every Check Wait by
  `delivery.check_timeout`, ADR-0311: "**One key bounds every Check Wait.**",
  and this Spec keeps that bound while it widens what the wait waits for and
  what the `checks-timeout` blocker names; ADR-0369 amends it. ADR-0236 keeps
  the re-run of a stale or outside-change failure, ADR-0236: "A failure
  without the annotation keeps the rules it had before", and the re-run rules
  stay as they are. ADR-0365
  records each item's Queue Turn in a file, ADR-0365: "**The turn is a
  file.**", and the Check Wait's pending checks join that record; ADR-0369
  amends it. ADR-0304 keeps waits out of the Run Database, ADR-0304: "A new
  `waiting` stage or an item column was rejected", so no stage or column is
  added. ADR-0199, which ADR-0365 cites, carries no obligation here: no limit
  on starts or tokens changes, ADR-0199: "Tokens outside Runs, such as the pre-PR
  review, do not count.". ADR-0237, which cites ADR-0199, carries no
  obligation either: the retry's merge observation does not change, ADR-0237:
  "A Delivery Retry records a merge made outside the queue". ADR-0369, written
  by this Spec, records the decision. Source: `docs/agents/domain.md`, `docs/agents/cli.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 3 answers of 2026-11-06 quoted in the authorization record;
  the governed files are the Roundfix skill files the documentation Task
  changes. No Makefile, workflow, Go module, lint, formatter or Baseline
  module change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0359-the-delivery-queue-waits-for-every-check-the-head-triggered/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Goals

- The Delivery Queue never merges a pull request while a check, a workflow
  run or an Expected Workflow of its head is not completed, whatever the
  branch rules require and whatever the merge state says.
- A check or workflow run of the head that ends in failure, a timeout,
  cancellation, a startup failure, a required action or staleness parks the
  item as `checks-failed`, naming each failed one and its link, with the
  re-run rules for stale and outside-change failures unchanged.
- The queue does not mistake a head whose workflows have not registered yet
  for a head with nothing to wait for.
- While it waits, `deliver status` and the console log name the checks still
  pending; at the deadline the `checks-timeout` blocker names them.
- The guides, the Roundfix skill, `CONTEXT.md` and ADR-0369 state the rule.

## Core Features

1. The **Head Check Set**. At each poll of `checking` the queue reads, for
   the pull request head, every check run and commit status `gh pr checks`
   lists (as today), the latest GitHub Actions workflow run per workflow file
   and event whose head commit is the head, and the head's Expected
   Workflows; a check run that belongs to a listed workflow run counts once.
2. The **Expected Workflow**. A workflow file under `.github/workflows/` at
   the head whose trigger makes GitHub run it for every pull request into the
   base branch: a `pull_request` or `pull_request_target` trigger without a
   path filter or `branches-ignore`, whose `branches` (if any) names the base
   branch literally without a `!` pattern, and whose `types` (if any) name
   both `opened` and `synchronize`; or a `push` trigger with no filter. Until
   a workflow run of that file exists for the head, the set is not complete.
3. The completion rule. The item moves to `merging` only when the set has at
   least one member, every member is completed, every Expected Workflow has a
   run, and the merge state is neither `BLOCKED` nor `UNKNOWN`; a merge state
   of `CLEAN` alone never moves it.
4. The conclusions. `success`, `neutral` and `skipped` pass. `failure`,
   `timed_out`, `cancelled`, `action_required`, `startup_failure` and
   `stale`, and any other conclusion, fail. A failed check keeps today's
   stale and outside-change re-run path; any other failure parks
   `checks-failed: <check> (<conclusion>) <link>`, every failed member of
   that read joined by `; `.
5. The wait is visible. When the pending members change, the console log
   writes `roundfix: checks: Delivery Queue item <slug>: waiting for <n>
   check(s): <check>; <check>`; when the set completes it writes `roundfix:
   checks: Delivery Queue item <slug>: <n> check(s) completed for <head>`; a
   read error is logged once per distinct error. `deliver status` prints the
   pending checks on the item's `Turn:` line, and at the deadline the item
   parks as `checks-timeout: current-head checks not settled within
   delivery.check_timeout <duration>; pending: <check>; <check>`.
6. The rules in writing. The Roundfix skill's `deliver` reference, the
   `deliver` and configuration guides, `CONTEXT.md`, ADR-0369 and the
   amendment notices in ADR-0311 and ADR-0365 state the rule.

## User Experience

The item of Fluxus pull request #466 at its first poll after the change,
then at the failure:

```text
roundfix: checks: Delivery Queue item 0149-o-fluxus-serve-a-stock-attention-ao-oraculum: merge state CLEAN
roundfix: checks: Delivery Queue item 0149-o-fluxus-serve-a-stock-attention-ao-oraculum: waiting for 2 check(s): CI | Frontend (in_progress); CI | Backend (in_progress)
roundfix: checks: Delivery Queue item 0149-o-fluxus-serve-a-stock-attention-ao-oraculum: waiting for 1 check(s): CI | Backend / Lint, Type Check, Test & Build
```

and `deliver status` once the item parked:

```text
0149-o-fluxus-serve-a-stock-attention-ao-oraculum	parked	checks-failed: CI | Backend / Lint, Type Check, Test & Build (failure) https://github.com/gesttione-solutions/fluxus/actions/runs/<run>/job/<job>	<worktree>
```

## Non-Goals / Out of Scope

- Evaluating `paths`, `paths-ignore`, branch globs or `branches-ignore`: a
  filtered workflow is never expected, but once GitHub creates its run, the
  run is a member and must complete.
- Waiting for check suites: a suite with no check run never completes (the
  `cursor` suite on Fluxus), so suites are not read.
- A settle delay, a new configuration key or a new Delivery Stage, Park
  Class, blocker name or Run Database column.
- The re-run rules of ADR-0236 and the flaky-check rule, which stay as they
  are, and the Baseline clause that says the queue verifies required checks,
  which stays true.
- Spec 0360's build-cache bound, disk check and disk-full park.
- The Makefile, CI workflows, `go.mod`, `go.sum`, lint or formatter
  configuration, any Baseline module and any `### QA settlement` section.

## Success Metrics

1. Success Metric: replaying the three Fluxus heads (one passing check
   listed while two workflow runs are in progress), the item stays in
   `checking`, and parks `checks-failed` naming `CI | Backend / Lint, Type
   Check, Test & Build (failure)` with its link when that check fails, or
   moves to `merging` when every member succeeds.
2. Success Metric: a head with an Expected Workflow and no workflow run yet
   stays in `checking` and parks at the deadline as `checks-timeout` naming
   `<path> (no run yet)`; a workflow file filtered by `paths` creates no
   expectation.
3. Success Metric: a workflow run concluded `startup_failure`,
   `action_required`, `cancelled` or `timed_out` with no failing check of its
   own parks `checks-failed` naming the workflow and its run link; `neutral`
   and `skipped` pass; an earlier run of a workflow superseded by a later run
   of the same file and event is ignored.
4. Success Metric: the console log lines and the `Turn:` line read exactly
   as the TechSpec's Surface Transcripts, and every existing Delivery Engine
   and `deliver` test passes, except the tests the TechSpec's Declared breaks
   name, which change only the blocker and log text stated there.
5. Success Metric: the `deliver` reference of the Roundfix skill and its
   mirror, the `deliver` and configuration guides, `CONTEXT.md` (Head Check
   Set, Expected Workflow, the revised Check Wait and Queue Turn), ADR-0369
   and the notices in ADR-0311 and ADR-0365 state the rule.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Fluxus report, adopted verbatim in `references/` before this Spec
  existed: "no momento do merge só "Validate PR title" tinha rodado; "CI |
  Backend" e "CI | Frontend" ainda estavam pending e terminaram success
  depois do merge".
- GitHub's REST documentation for workflow runs
  (https://docs.github.com/en/rest/actions/workflow-runs, read through Exa
  on 2026-11-06), which defines the filter the queue reads: "`head_sha`
  (string) Only returns workflow runs that are associated with the specified
  head_sha."
- GitHub's documentation of the events that trigger workflows
  (https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows,
  read through Exa on 2026-11-06), on which the Expected Workflow rule rests:
  "By default, a workflow only runs when a `pull_request` event's activity
  type is `opened`, `synchronize`, or `reopened`."

## Unreachable Acceptance

- criterion: a Fluxus pull request delivered by the released queue is merged only after every workflow its head triggered completed
  reason: a Run cannot reach the Fluxus repository, GitHub or the MacBook that runs its queue; the QA gate proves the rule on synthetic reports through the Delivery Engine and on a scripted gh through the GitHub boundary
  satisfied-by: the Fluxus console log and pull request timeline after the release that carries this Spec, reported back through ai-memory

## Validation

The decisions were checked on 2026-11-06 against these sources.

- Reading at `6ceb7ff5`: `checkCandidate` (`internal/delivery/engine.go`)
  polls `GitHubCLI.CurrentHeadChecks` (`internal/delivery/github.go`), which
  runs `gh pr checks <n> --json bucket,link,name,state,workflow` without
  `--required`; a `fail` bucket parks bare `checks-failed` (after the re-run
  rules), a `cancel` bucket parks bare `checks-failed`, any other bucket and
  an empty list keep waiting, `BLOCKED` and `UNKNOWN` keep waiting, and
  otherwise the item moves to `merging`. The only log line is `roundfix:
  checks: Delivery Queue item <slug>: merge state <state>`, the line Fluxus
  quoted.
- The GitHub timing recorded in `references/2026-11-06-fluxus-check-timing.md`:
  the three workflow runs of each head existed about eight seconds after the
  pull request was created, and their job check runs appeared more than a
  minute later, after the merge.
- A prototype on a scratch copy of `6ceb7ff5` added the workflow-run read,
  the Expected Workflows and the completion rule. In `./internal/delivery`
  only the tests the TechSpec's Declared breaks name failed, each on blocker
  or log text; the scripted GitHub boundary tests failed on the two new
  commands.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix delivery queue checks merge state CLEAN fluxus" --all --files --min-score 0.3`,
  which returned the Roundfix mirrors of the `deliver` guide, the Roundfix
  skill's `deliver` reference, `CONTEXT.md` and the changelog, and Fiscus's
  mirror of `docs/agents/autonomous-work.md`; none records a decision on
  which checks the queue waits for, and the Fluxus entity page
  (`wiki/entities/fluxus.md`) records no CI decision. The Secondbrain
  confirmed the current contract and changed no decision.
- Exa: the two GitHub pages quoted above. They led to reading workflow runs
  by head commit instead of check suites, and to the Expected Workflow rule
  instead of a settle delay.

## Glossary

- adds: **Head Check Set**
- adds: **Expected Workflow**
- changes: **Check Wait**
- changes: **Queue Turn**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ batch 3 brief for Spec 0359 of 2026-11-06 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-06), verbatim:

- Question "O Fluxus relatou que a Delivery Queue fez merge dos PRs #466 e
  #469 com o CI vermelho e do #477 antes de o CI terminar. O ruleset do main
  não exige checks, então o GitHub dava o merge como CLEAN, e a fila confia
  só nisso. O que a fila deve fazer?" → "Esperar todos e parar
  (Recommended)" — option text: "Spec agora: na etapa de checks, a fila
  espera todos os workflows que o head disparou terminarem, até o prazo de
  sempre. Se algum falhar, mesmo sem ser obrigatório, o item para com o nome
  do check e o link."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, merge
  pela fila, conserto e retry quando a entrega parar. Release assim que cada
  Spec ficar pronta, com notas aos adotantes e resposta ao Fluxus e ao
  Oraculum."
- Release, standing instruction of 2026-11-05: "Quero que faça os releases
  hoje assim que estiverem prontos. E se existe no baseline regra para um
  release por dia por algum motivo deve ser removida." No such rule exists in
  the repository or the Baseline, so this Spec changes nothing for it.

Taken from those answers, the reading, the measured timing and the prototype
above, with no further interview:

- The set is completed by reading the workflow runs of the head commit
  (GitHub creates every run of one event in one batch, seconds after the
  event, before any job registers a check run) and by the Expected
  Workflows of the head's own files, which close the seconds before that
  batch. A settle delay was rejected: it would add a fixed wait to every
  merge and to every test that drives `checking` on a real clock, and it
  still guesses. Evaluating path filters was rejected: GitHub's filter
  semantics would be re-implemented, and a wrong guess would wait for a run
  that never comes until the Check Wait ends. See ADR-0369.
- `neutral` and `skipped` pass because GitHub itself treats them as not
  failing; `action_required` parks because it never completes without a
  person; `startup_failure` and `stale` park because the workflow did not
  test the head. See ADR-0369.
- The park keeps `checks-failed` and `checks-timeout` and their Park Classes
  (`finding` and `environment`); only the text after the colon grows, which
  `ClassifyPark` already ignores. See ADR-0369.
- Four implementation Tasks (the documents with ADR-0369 first, the GitHub
  reader, the completion rule, the `Turn:` line) and a QA gate.

## Open Questions

None.
