---
task: task_01
spec: 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
status: completed
type: backend
complexity: high
---

# Task 01: The item records its Failed Run and a Delivery Retry carries the Tasks of a QA-Only Failed Run

## Overview

This Task delivers `_techspec.md` → Build Order 1. The Queue Owner writes on the Delivery Queue item the Run its Implement executor reports as soon as the executor returns, a Run that ended Failed included, and a Delivery Retry accepts a QA-Only Failed Run as one more source of Task Carry-Forward while the Reconcile Command and the Implement Command's Preflight keep their three outcomes. It answers two findings reported by quaestor on 2026-11-23, ai-memory messages 01a20714-3a6e and 01a20714-ed79: the item kept the earlier Run after its Run ended Failed, and the retry after a Run that failed only in its QA gate carried none of its five completed Tasks. It documents the QA-Only Failed Run in the retry section of the `deliver` guide, the one guide edit of this slice. It is verifiable on its own through three new tests over the production Delivery Engine with a fake runner and the production carry-forward over temporary repositories and a temporary Run Database.

## Requirements

1. MUST make the executor's result name the Run its child created for every exit of the child, as `_techspec.md` → API Contract 1 states: the failure branch of `deliveryRunResult` returns the created Run with the failure, whose text does not change, and returns no Run when the child created none. A Run that ended Unresolved or Budget Exceeded keeps today's result.
2. MUST add `recordRun` to the Delivery Engine with the signature of `_techspec.md` → Interfaces and call it first thing after the Implement executor returns, in `runCandidate` and in the Run Slot Wait, before the error or the outcome is read. When the reported Run differs from the item's, it writes the Run on the item in memory and in the Run Database under a context the owner's stop does not cancel; a store error is returned joined with the executor's error. The park that follows keeps its blocker, word for word `delivery-error: run Implement executor: roundfix implement failed with exit code 1` for the case quaestor reported, and its Park Class.
3. MUST add `qaOnlyFailedRun` with the signature of `_techspec.md` → Interfaces and classify a QA-Only Failed Run exactly as API Contract 2 states: outcome Failed, Run Branch `roundfix/run-<run-id>` present, the Spec's Task Graph loaded from the commit at that branch's tip with the loader the Implement Command uses for a committed graph, a QA gate's Task in the graph that is not `completed` there, and every other Task `completed` there. A Run whose branch is gone or whose graph does not load there is skipped as today. It reads local Git objects only, never the Run Worktree's files.
4. MUST make `commandDeliveryWorkflow.CarryForward` treat a QA-Only Failed Run as it treats an accepted outcome, through the one new `qaOnlyFailed` field of `reconcileRunSelection`: the same proofs for every Task, the whole remaining set carried or refused, every terminal Run of the Spec considered newest first, each carried Task stamped with its Run and commit, one `Carried forward from Run <run-id>: <task>, <task>` line per Run. The evidence of the QA gate's Task is removed before inspection, so that Task is never a candidate (invariant 4). The refusal of a QA-Only Failed Run is whole and reported as today's refusal, and its next action never suggests `roundfix reconcile <run> --carry-forward`, which does not accept a Failed Run.
5. MUST keep `runworktree.CarryForwardAcceptedOutcomes` at its three outcomes (invariant 3). The Reconcile Command and the Implement Command's Preflight never set the new field, so they accept and refuse exactly what they do on the Delivery Base (`_prd.md` → Core Feature 3, Delegated Decision of 2026-11-23 on the routes to Task Carry-Forward).
6. MUST add `internal/delivery/failed_run_record_test.go` with `TestTheItemRecordsTheRunOfAFailedImplementExecutor`, over the production Delivery Engine and the package's fake runner and a temporary Run Database, with these subtests, each named as written:
   - `failed_executor_parks_with_its_run`: an item that records `run-earlier` and a runner that returns a Run with the executor's error; the stored item is `parked` with the reported Run and the unchanged blocker.
   - `stopped_owner_still_records_the_run`: the runner cancels the owner's context before it returns the Run; the stored item records the Run though nothing parked it.
7. MUST add `internal/cli/deliver_qa_only_failed_run_test.go` with:
   - `TestAFailedImplementChildNamesTheRunItCreated`, with the subtests `failed_run_is_returned_with_the_failure` (exit `1` with a new Run in state Failed returns that Run and the unchanged failure text) and `no_run_created_returns_none`.
   - `TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun`, built on `createDeliveryRetryRun` with the outcome Failed, a Run Branch holding the committed graph, five settled Tasks and a QA gate's Task, through the production `CarryForward`. Its subtests, each named as written: `carries_the_settled_tasks` (the five Tasks are completed on the Item Branch and stamped with the Run and commit, the gate's Task stays pending, one `Carried forward from Run` line names the five); `another_task_not_completed_carries_nothing`; `run_branch_gone_carries_nothing`; `unresolved_run_carries_as_before`; `preflight_reports_no_candidate` (`inspectSpecCarryForwards` reports no candidate for the QA-Only Failed Run); and `reconcile_keeps_refusing_the_failed_run` (the Reconcile Command's carry-forward of that Run refuses it as on the Delivery Base).
   When `createDeliveryRetryRun` needs a parameter for the Run Branch or the gate's Task, extend it so that every existing caller keeps its behavior.
8. MUST revise the glossary through `domain-modeling`: add **QA-Only Failed Run** to `CONTEXT.md` with the definition of `_techspec.md` → API Contract 2, next to Task Carry-Forward and in the file's form (the bold term and its colon alone on one line, the definition on the next line, then its `_Avoid_` line), and revise Task Carry-Forward to say that a Delivery Retry also accepts a QA-Only Failed Run and that the Reconcile Command and the Preflight do not. Neither entry links ADR-0450 yet: task_03 adds that link when it creates the record.
9. MUST record the Coverage Review of the Behavior Surface `command reconcile`: in `docs/references/skill-coverage.json`, prepend to the `review` value of that entry one dated sentence, dated the day this Task writes it, that begins `<YYYY-MM-DD> — Spec 0418 task_01` and says that a field only a Delivery Retry sets was added and that the Reconcile Command's behavior, flags, output and exit codes are unchanged; keep the existing text after it. No other entry of the map changes (`_prd.md` → Skills).
10. MUST document this slice in `docs/user-guide/commands/deliver.md`, where the retry section says that the retry carries the remaining settled Tasks from every terminal Implement Run of the item's Spec: add that a Run that ended Failed is a source only when it is a QA-Only Failed Run, that its QA gate's Task is never carried, that the Reconcile Command and the Implement Command's Preflight do not accept it, and that the item records a Run that ended Failed before it parks. No other paragraph of the guide changes; the status lines and the retry's failure report are task_04's. Then run the `-record-skill-coverage` step of `_authorization.md` → Sanctioned regeneration, `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, which re-records the guide page's fingerprint in `docs/references/behavior-surfaces.json`; hand-edit none of its output. The Roundfix skill's files and their sync and version record stay task_04's.
11. MUST keep every new top-level test calling `t.Parallel()` as its first statement, as the Parallel Test Package rule requires; no new test waits on elapsed time, reaches the network, starts an Agent or reads the real home directory.
12. MUST keep passing unchanged `TestCarryForwardAcceptsBudgetExceededRun`, the item recovery tests and every reconcile and Preflight carry-forward test. When an existing test of `./internal/delivery` or `./internal/cli` pinned the stale Run on the item, update only that assertion to the new contract and name the test in the Result.
13. MUST NOT change the `deliver` output contracts of task_02 and task_03, any skill, any other user-guide page, any decision record, `.secondbrain-export`, the Makefile, `go.mod` or any file `_authorization.md` → Limits names.

## Subtasks

- [ ] Return the created Run from a failed Implement child and record it on the item in `recordRun`.
- [ ] Classify a QA-Only Failed Run from the committed graph at its Run Branch and route it through the retry's carry-forward.
- [ ] Add the delivery test and the two CLI tests with their controls.
- [ ] Add QA-Only Failed Run and revise Task Carry-Forward in the glossary.
- [ ] Write the Coverage Review of `command reconcile`.
- [ ] Name the QA-Only Failed Run in the guide's retry section and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] For every exit of the Implement executor that created a Run, Failed included, the stored item names that Run, with the park's blocker and Park Class unchanged.
- [ ] A Delivery Retry after a QA-Only Failed Run carries its completed Tasks, never the gate's Task; a Failed Run with another Task not completed or with no Run Branch carries nothing; an Unresolved Run carries as before.
- [ ] The Reconcile Command, the Preflight and the accepted-outcome list behave as on the Delivery Base.
- [ ] The glossary defines QA-Only Failed Run and Task Carry-Forward names it; the guide's retry section names it; the reconcile surface carries its dated Coverage Review; the Behavior Surface Record is current.
- [ ] The whole `./internal/delivery`, `./internal/cli` and `./internal/worktree` packages pass.

## Context

- instruction: `internal/worktree/worktree.go`
- instruction: `internal/delivery/engine_test.go`
- instruction: `internal/cli/carryforward_completed_target_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/carryforward.go`
- interface: `internal/cli/reconcile.go`
- interface: `internal/cli/deliver_retry_runs_test.go`
- interface: `CONTEXT.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/delivery/failed_run_record_test.go`
- creates: `internal/cli/deliver_qa_only_failed_run_test.go`

## Consumers

- the Run the Delivery Engine records on a Delivery Queue item and the result of the Implement executor — consumers: `./internal/delivery`, `./internal/cli`
- the carry-forward selection the Delivery Retry, the Reconcile Command and the Preflight share — consumers: `./internal/cli`, `./internal/worktree`
- the deliver guide page, the Skill Coverage Map and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheItemRecordsTheRunOfAFailedImplementExecutor$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheItemRecordsTheRunOfAFailedImplementExecutor/failed_executor_parks_with_its_run TestTheItemRecordsTheRunOfAFailedImplementExecutor/stopped_owner_still_records_the_run TestTheItemRecordsTheRunOfAFailedImplementExecutor; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the item records the Run of a failed executor in both subtests.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAFailedImplementChildNamesTheRunItCreated|TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAFailedImplementChildNamesTheRunItCreated/failed_run_is_returned_with_the_failure TestAFailedImplementChildNamesTheRunItCreated/no_run_created_returns_none TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun/carries_the_settled_tasks TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun/another_task_not_completed_carries_nothing TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun/run_branch_gone_carries_nothing TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun/unresolved_run_carries_as_before TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun/preflight_reports_no_candidate TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun/reconcile_keeps_refusing_the_failed_run; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the failed child names its Run and the retry carries the QA-Only Failed Run with every control passing.
- `grep -F -- '**QA-Only Failed Run**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*QA-Only Failed Run\*\*(:| —( |$))'` — expected: exit 0, the glossary defines the new term.
- `grep -F -- '**Task Carry-Forward**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Task Carry-Forward\*\*(:| —( |$))' && awk '$0 == "**Task Carry-Forward**:" { found=1; next } found && $0 == "" { exit } found { print }' CONTEXT.md | tr -s '[:space:]' ' ' | grep -qF -- 'QA-Only Failed Run'` — expected: exit 0, the revised Task Carry-Forward entry names the new term.
- `grep -qF 'Spec 0418 task_01' docs/references/skill-coverage.json && tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "QA-Only Failed Run" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "QA-Only Failed Run" >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0, the reconcile surface carries its dated Coverage Review, the guide names the QA-Only Failed Run, and the documentation contracts pass with the re-recorded Behavior Surface Record.
- `grep -qF 'func (engine *Engine) recordRun(' internal/delivery/engine.go && grep -qF 'qaOnlyFailed' internal/cli/reconcile.go && go test -count=1 -timeout 30m ./internal/delivery ./internal/cli ./internal/worktree` — expected: exit 0, the consumers of the recorded Run and of the carry-forward selection pass whole.

## References

- `_prd.md` → Goals; User Stories 1 and 5; Core Features 1, 2, 3 and 7; Success Metric 1; Success Metric 2; Skills; Glossary; Decisions
- `_techspec.md` → The item's latest Run; The QA-Only Failed Run as a source; Interfaces; Data Models; API Contract 1; API Contract 2; API Contract 8; Invariants; Coverage Map; Testing Approach; Build Order 1
- [`references/2026-11-23-characterization.md`](references/2026-11-23-characterization.md), sections 1, 2 and 3
- [`references/sources.md`](references/sources.md), sections 1 and 2 (reported by quaestor, 2026-11-23)
- ADR-0170, ADR-0158, ADR-0284, ADR-0304

## Result

Skills loaded: implement-task, coding-guidelines, domain-modeling, golang-cli, golang-context, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, no-workarounds, systematic-debugging, tech-writer, evidence-gate.

Implemented this Task's slice:

- The failed-child result now includes the created Run without changing the existing failure diagnostic. Both executor call sites record a nonempty changed Run before handling errors or outcomes. Recording uses `context.WithoutCancel`; a recording failure joins the executor error. The engine regression covers the unchanged blocker and environment Park Class, and recording after owner cancellation without parking.
- Delivery Retry classifies a Failed Implement Run from local Git objects at its Run Branch tip using the committed-graph loader. Only a graph whose sole unfinished Task is its QA gate qualifies. The retry removes the gate's evidence and sets the new selection field before using the existing proofs and whole-set application. Failed-Run refusals name Delivery Retry as recovery and never suggest Reconcile.
- The new CLI regressions carry five Tasks with Run/commit provenance, leave the gate pending, and assert exactly one per-Run output line. Controls cover an unfinished implementation Task (even with an uncommitted completed status), a missing branch, an unloadable committed graph, an already completed gate, stale completed gate evidence, whole-set refusal on moved inputs, unchanged Unresolved recovery, Preflight exclusion, and Reconcile's unchanged exit and accepted-outcome diagnostic.
- `CarryForwardAcceptedOutcomes` remains unchanged. No existing test, fixture, or stale-Run assertion needed modification; `createDeliveryRetryRun` already supports the required committed graph and QA fixture.
- Added QA-Only Failed Run next to Task Carry-Forward and revised the latter's route distinction through domain-modeling. Only the retry paragraph of the deliver guide changed. Prepended the dated review only to `command reconcile`. The sanctioned recorder changed only the deliver guide fingerprint in the Behavior Surface Record. ADR-0450 links and all other Tasks' work remain outside this diff.

Focused evidence (all commands below exited 0 unless explicitly described as the initial reproduction):

1. Initial regression run before production edits: `go test ./internal/delivery ./internal/cli -run 'Test(TheItemRecordsTheRunOfAFailedImplementExecutor|AFailedImplementChildNamesTheRunItCreated|ADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun)' -count=1 -timeout=4m` reproduced both stale-Run cases, the dropped child Run, and missing five-Task carry. The initial child-error expectation also exposed the existing trailing `: exit code 1` diagnostic; the test was corrected to preserve that existing contract. The same focused command passed after implementation.
2. `go test ./internal/cli ./internal/delivery ./internal/worktree -run 'Test(ADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun|AFailedImplementChildNamesTheRunItCreated|TheItemRecordsTheRunOfAFailedImplementExecutor|CarryForward|ItemRecovery|ReconcileCarryForward|DeliveryEngine|DeliveryActions|RunSlot|DeliveryRunResult)' -count=1 -timeout=6m` passed, including unchanged Budget Exceeded acceptance, item recovery, and carry-forward compatibility checks.
3. After strengthening the exact Reconcile refusal assertion: `go test ./internal/cli -run 'Test(ADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun|ImplementRefusesWhenCarryForwardIsAvailable|ImplementCarryForwardRefusesOnlyWhenTheSetWouldCarry|InspectSpecCarryForwards)' -count=1 -timeout=4m` passed.
4. `go test ./internal/delivery -run 'Test(ARunCeilingRefusalWaitsAndStartsWhenASlotFrees|ARunSlotWaitParksAsTodayAtItsDeadline|TheQueueDeadlineEndsARunSlotWait|AZeroRunSlotWaitParksAtOnce|OtherRunSpecFailuresDoNotWait|AStoppedOwnerClearsTheRunSlotWait|TwoQueuesWaitingForOneSlotStartOneRun)' -count=1 -timeout=4m` passed for the second executor call site.
5. Required sanctioned regeneration: `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` passed and generated the record; no generated output was edited by hand.
6. `go test -tags docscontract ./internal/docscontract -run '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs|ClosedListsEqualTheirConstants|AvoidLinesArePlainLists|StatesEachRuleOnce)$' -count=1` passed. `git diff --check` passed.

Acceptance evidence: the first criterion is covered by the child/engine regressions and inspection of both recording call sites; the second by all retry cases above; the third by unchanged accepted-outcome source plus compatibility, Reconcile and Preflight controls; the fourth by glossary checks, scoped diff inspection and sanctioned regeneration. The fifth criterion (whole delivery, CLI and worktree packages), full documentation contracts, and every declared Verification command remain for Daemon Verification. No terminal Task verdict is claimed. The incoming Daemon-owned `in_progress` status was preserved; no manifest, other Task, skill, decision record, module or tooling file was edited. No commit, push, or pull request was made.

## Carry-forward provenance

- Source Run: `run_20261124T001551Z_9b40dc21bfabfe76`
- Source commit: `12c6950a9717b36e22ad855560f1a247b0b9c60b`
