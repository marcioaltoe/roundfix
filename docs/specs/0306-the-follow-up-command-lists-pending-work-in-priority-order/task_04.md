---
task: task_04
spec: 0306-the-follow-up-command-lists-pending-work-in-priority-order
status: completed
type: backend
complexity: medium
---

# Task 04: roundfix follow-up --status shows the work in progress, the pending Specs, the next release and what waits on the maintainer

## Overview

This Task adds the status view of the Follow-up Command: `--status` reads the
in-progress sources task_03 built, the active Specs with their Task counts,
the Release Plan, the parked items and the queue's Pending Question, and
never asks the judge (`_prd.md` → Core Feature 6; `_techspec.md` → API
Contract 8). It finishes the command's help and root usage, so it asserts
Surface Transcripts 1 and 3.

## Requirements

1. MUST add `--status` to the command: the help text becomes exactly
   `_techspec.md` → Surface Transcript 1, the root usage line becomes
   `  roundfix follow-up [--status] [--messages <path>] [--format <text|json>]`,
   and `--status` together with `--messages` is refused before any read with
   the two standard error lines and exit `2` of Surface Transcript 3.
2. MUST implement the status view of API Contract 8 in a new
   `internal/cli/follow_up_status.go`, reusing task_03's `runs`,
   `delivery-queue`, `pull-requests`, `git` and `specs` readers,
   `spec.Load` for the Task counts, `releaseplan.Build` with an empty
   `Request` over `newReleasePlanGitSource` and the `releasePlanGitRunner`
   dependency, `errors.Is(err, releaseplan.ErrInvalidReleaseRange)` for
   `no_unreleased_changes`, `delivery.ClassifyPark` for each parked item, and
   `delivery.PendingQuestionFor` for the queue's question. Any other Release
   Plan error is `unavailable` with its reason and never changes the exit
   code. The view never calls the judge and never reads the Secondbrain.
3. MUST order `pending_specs` as the queue's unmerged items by position, then
   the active Specs not in the queue by slug, and list under `maintainer` the
   Pending Question and, for a plan in `approval_required` or
   `manual_classification_required`, one release entry; and compute `state`
   as API Contract 8 states.
4. MUST add `internal/cli/follow_up_status.go` to the `sources` of the
   `command follow-up` entry in `docs/references/skill-coverage.json`, and run
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
   last, because the command's help changes. The record is never
   hand-edited.
5. MUST add `internal/cli/follow_up_status_test.go`, under the same harness
   rules as task_03's tests (every test calls `t.Parallel()` first, no
   real-time constant below ten seconds, no process environment variable, no
   real `gh`, no network, temporary repositories and homes, the CLI test
   binary never started as `roundfix`), with:
   - `TestFollowUpStatusListsPendingSpecsInDeliveryOrder` (Success Metric 6):
     a queue whose positions put `0009-second` before `0003-first`, plus an
     unqueued active `0005-unqueued`, lists `0009-second`, `0003-first`,
     `0005-unqueued` in that order with their stages and Task counts;
   - `TestFollowUpStatusReportsTheNextRelease` (Success Metric 6), with
     subtests: `HEAD` on the last release tag gives `no_unreleased_changes`
     and the text line `Next release: no unreleased changes`; a later `feat:`
     commit gives the plan's state, base, proposed version and that commit's
     subject; an uncommitted change gives `unavailable` with the plan's
     reason and exit 0;
   - `TestFollowUpStatusListsBlockersAndTheMaintainersQuestion` (Success
     Metric 6): two parked items give two blockers with their Park Class and
     next action, and the lowest-position one is the one `pending_question`
     under `maintainer`, with the answer `deliver status` prints;
   - `TestFollowUpStatusNeverAsksTheJudge`: with a Judge Key Variable in the
     command environment and a fake transport, `--status` sends no request;
   - `TestFollowUpStatusWritesNothing` (Success Metric 7): the repository,
     the Roundfix Home and `git status --porcelain` are byte-identical before
     and after;
   - `TestFollowUpHelpTranscript`: asserts Surface Transcript 1 whole —
     every standard output line in order, empty standard error, exit 0;
   - `TestFollowUpStatusRefusesMessages`: asserts Surface Transcript 3 whole,
     with no `messages.json` on disk.
6. MUST NOT change the list view's output, `internal/store`,
   `internal/delivery`, `internal/releaseplan` or `roundfix release plan`.

## Subtasks

- [ ] Add `--status` to the parser, the help and the root usage, with its refusal.
- [ ] Build the status view from the existing readers and the Release Plan.
- [ ] Render the JSON and text views with the stop state.
- [ ] Update the map sources, add the tests, and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] `--status` lists queued Specs by position before unqueued ones, the next release or "no unreleased changes", each park with its Park Class, and the Pending Question.
- [ ] A Release Plan failure is reported in the view and never fails the command.
- [ ] The status view sends nothing to the judge and writes nothing.
- [ ] The help and both refusals match their Surface Transcripts.

## Consumers

- the `roundfix follow-up` command surface — consumers: `./internal/docscontract`

## Context

- instruction: `internal/cli/releaseplan_command.go`
- instruction: `internal/cli/releaseplan_git_source.go`
- instruction: `internal/cli/deliver.go`
- instruction: `internal/cli/deliver_limits_test.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/follow_up.go`
- creates: `internal/cli/follow_up_status.go`
- creates: `internal/cli/follow_up_status_test.go`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestFollowUpStatusListsPendingSpecsInDeliveryOrder|TestFollowUpStatusReportsTheNextRelease|TestFollowUpStatusListsBlockersAndTheMaintainersQuestion|TestFollowUpStatusNeverAsksTheJudge|TestFollowUpStatusWritesNothing|TestFollowUpHelpTranscript|TestFollowUpStatusRefusesMessages|TestFollowUpSaysNothingToDo|TestEveryCommandPathPrintsItsHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFollowUpStatusListsPendingSpecsInDeliveryOrder TestFollowUpStatusReportsTheNextRelease TestFollowUpStatusListsBlockersAndTheMaintainersQuestion TestFollowUpStatusNeverAsksTheJudge TestFollowUpStatusWritesNothing TestFollowUpHelpTranscript TestFollowUpStatusRefusesMessages TestFollowUpSaysNothingToDo TestEveryCommandPathPrintsItsHelp; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestEveryCommandIsNamedInTheRoundfixSkill)$' ./internal/docscontract || exit 1; grep -qF 'internal/cli/follow_up_status.go' docs/references/skill-coverage.json || { echo 'status source not mapped' >&2; exit 1; }; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the status tests do not exist, so the command fails at its first pass check. Afterwards the status view, its release, blocker and maintainer sections, its read-only and judge-free boundary, the help and the refusal pass with the list view's stop rule and every command's help, the map is current, and the package's tests compile for Windows.

## References

- `_prd.md` → Goal 4; User Story 4; Core Feature 6; Success Metric 6; Success Metric 7
- `_techspec.md` → API Contracts 1 and 8; Surface Transcripts 1 and 3; Invariants 1-3 and 6; Testing Approach; Build Order 4
- ADR-0143
- ADR-0319


## Result

Implemented the `--status` text and JSON views without changing the list
report. Status branches before the list inventory, reuses the in-progress
readers and active-Spec discovery, loads Task counts with `spec.Load`, and
builds the Release Plan with an empty request over the existing Git adapter
and its command dependency. The queue snapshot supplies the same
`PendingQuestionFor` answer as `deliver status`; parked rows retain
`ClassifyPark`'s Park Class and next action.

Acceptance evidence from focused implementation checks:

- Pending Specs, release, parks and question:
  `TestFollowUpStatusListsPendingSpecsInDeliveryOrder` checks position order
  (`0009-second`, `0003-first`) before `0005-unqueued`, stages, null queue
  fields, Task counts and text. Merged items are excluded.
  `TestFollowUpStatusListsBlockersAndTheMaintainersQuestion` checks both
  parks, the lowest-position question and the public `deliver status`
  answer. Release subtests cover no unreleased changes, feature approval,
  patch readiness, maintenance-only `no_release`, manual classification,
  dirty worktree and absent release tag, including both stop states.
- Release errors: dirty-worktree and missing-tag subtests check
  `unavailable`, its reason and exit 0; `errors.Is` recognizes the wrapped
  invalid-range error as `no_unreleased_changes`.
- Judge-free and read-only: `TestFollowUpStatusNeverAsksTheJudge` supplies a
  Judge Key Variable through the command environment, pending source items,
  a malformed Secondbrain entry and the existing rejecting fake transport;
  neither format sends a request or lists a Secondbrain source.
  `TestFollowUpStatusWritesNothing` compares repository, Roundfix Home and
  Secondbrain file bytes and Git porcelain output before and after both
  formats, with absent database, closed database, live WAL, clean release
  and dirty worktree cases.
- Help and refusals: `TestFollowUpHelpTranscript` checks Surface Transcript
  1 in full and the root usage line. `TestFollowUpStatusRefusesMessages`
  checks Surface Transcript 3, both flag orders and an explicitly empty
  messages flag against invalid configuration and no messages file.
  The existing unknown-format transcript and list-help test also pass;
  the latter expectation now follows the intentionally changed help.

Focused checks:

- Initial `go test ./internal/cli -run '^TestFollowUpHelpTranscript$'
  -count=1`: failed against the old help, establishing the pre-change gap.
- `go test ./internal/cli -run '^TestFollowUp' -count=1`: passed after the
  implementation and final test additions; includes the existing list
  behavior and all seven required status/help/refusal tests.
- `go test -race ./internal/cli -run '^TestFollowUp(Status|Help)' -count=1`:
  passed after the final test additions.
- Windows CLI test compilation through `go test -c ./internal/cli` with
  `GOOS=windows`, writing only a temporary executable: passed after the
  final test additions.
- `git diff --check`: passed after Behavior Surface Record regeneration.
- Added `internal/cli/follow_up_status.go` to the command's Skill Coverage
  Map sources. Ran the required generator last:
  `go test -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`:
  passed. The generated Behavior Surface Record changed only the
  `command follow-up` digest and was never hand-edited.

The only pre-existing changed path was this Task file, whose daemon-owned
`status: in_progress` is preserved. No changes were made to `internal/store`,
`internal/delivery`, `internal/releaseplan`, other Task files or `_tasks.md`.
Declared Task Verification and repository Verification were not run; the
Daemon owns those checks and Task settlement. No commit, push or Pull
Request was performed.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/follow_up_inventory.go`
- `internal/cli/follow_up_test.go`

## Carry-forward provenance

- Source Run: `run_20261024T004042Z_a30dec8c359819a6`
- Source commit: `f6dba5dc257d8dee69b71be34031afa4bb73d575`
