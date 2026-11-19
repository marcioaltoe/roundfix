---
task: task_01
spec: 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
status: completed
type: backend
complexity: medium
---

# Task 01: The status view counts a Delivery Queue item's Tasks on its Item Branch

## Overview

`roundfix follow-up --status` counts every Spec's Tasks in the checkout's Spec Root, so a Spec running in the Delivery Queue reads as not started: the Fiscus saw Spec 0097 at 0 of 16 Tasks while 15 were done on its Item Branch (the maintainer, 2026-11-19, [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md)). This Task gives the Task Graph loader a reader, reads a Delivery Queue item's Task Graph from its Item Branch through Git in the Spec queue and under `Pending Specs`, and states the rule in the `follow-up` guide, the glossary entry for the Spec Queue Table and the Skill Coverage Map. The reproduction is in [`references/2026-11-19-characterization.md`](references/2026-11-19-characterization.md), section 1.

## Requirements

1. MUST add `spec.LoadAt` to `internal/spec/spec.go` with the signature of `_techspec.md` → Interfaces: the PRD read of `requireActive`, the manifest read of `loadManifestNodes` and each Task file read of `loadTask` take the reader, and `spec.Load` and `spec.LoadForRecovery` pass `os.ReadFile`, so every error, path and result they return stays as it is. Remove any helper the change leaves without a caller.
2. MUST add `TestLoadAtReadsThroughTheGivenReader` to a new `internal/spec/load_at_test.go`, calling `t.Parallel()` first, with a map-backed reader that returns an error wrapping `os.ErrNotExist` for an absent path, and the subtests:
   - `same_graph_as_load`: the graph `LoadAt` returns for a Spec whose files the map holds equals, Task by Task (ID, File, Status, Type, Needs), the graph `spec.Load` returns for the same files on disk;
   - `an_absent_prd_is_not_found`: an error matching `spec.SpecNotFoundError` with `errors.As`;
   - `an_absent_task_file_is_named`: an error matching `spec.MissingTaskFileError` that names the Task.
3. MUST implement `followUpItemTaskGraph` in `internal/cli/follow_up_queue.go` with the signature of `_techspec.md` → Interfaces and invariants 1 and 2, using the item's `Branch`, `rev-parse --verify` of `refs/heads/<branch>^{commit}`, `git show <revision>:<path>` through `spec.LoadAt`, and `spec.ReadArchivedSpecAt` with the Archive Root's repository-relative path for an archived Spec. Paths given to `git show` use forward slashes. A nil runner uses `preflight.ExecGitRunner{}` and a nil context `context.TODO()`, as `followUpAuthoringPhases` does.
4. MUST make `followUpQueueRows` take a Delivery Queue item's counts from `followUpItemTaskGraph`, and `followUpPendingSpecs` in `internal/cli/follow_up_status.go` take a queued Spec's counts from the same function, as `_techspec.md` → API Contract 1 states. A planned Spec that is not queued, every authoring row and every Backlog Entry row keep today's reading. When no source loads, the counts stay 0 and 0. The text and JSON shape of the status view, its sources and its forecasts do not change, and the reader writes nothing.
5. MUST add `TestADeliveryItemCountsTasksOnItsItemBranch` to `internal/cli/follow_up_queue_test.go`, calling `t.Parallel()` first, built on `newFollowUpReadInput`, `writeImplementSpecAtRoot`, `gitImplement`, `followUpQueueDiscover` and `followUpSnapshot`. Each subtest commits a Spec `0097-widget` with two pending Tasks on `main`, then:
   - `a_running_item_counts_its_branch`: an Item Branch in a linked worktree settles `task_01` completed in a commit; a `running` item naming that branch shows `running 1/2 Tasks`;
   - `an_archived_item_counts_the_source_revision`: the branch settles three Tasks completed (a corrective `task_03` appended), then removes the folder and adds an Archive Record whose `source` is the folder and whose `source_revision` is that commit; a `reviewing` item shows `reviewing 3/3 Tasks`;
   - `an_item_without_a_branch_counts_the_spec_root`: a `queued` item with no branch shows `queued 0/2 Tasks`;
   - `a_missing_branch_counts_the_spec_root`: an item naming a branch that does not exist shows `queued 0/2 Tasks`;
   - `a_branch_without_the_spec_counts_the_spec_root`: the branch removes the folder without an Archive Record; the item shows `running 0/2 Tasks`;
   - `the_worktree_is_not_read`: the branch settles `task_01` in a commit and the worktree then sets `task_02` completed without committing; the item shows `running 1/2 Tasks`;
   - `the_pending_specs_line_counts_the_same`: `followUpPendingSpecs` gives the item of `a_running_item_counts_its_branch` 1 completed of 2;
   - `the_rows_write_nothing`: `followUpSnapshot` of the checkout and the home is equal before and after the rows are read.

   `TestTheQueueRowsComeInTableOrder`, `TestAQueueItemNamesItsWait` and `TestFollowUpStatusListsPendingSpecsInDeliveryOrder` stay unchanged and keep passing; so does every other test that pins the status view on the default branch when this Task runs.
6. MUST insert into the paragraph of `docs/user-guide/commands/follow-up.md` that begins "The status of a row is" the text of API Contract 1, word for word, right after the sentence that ends "`not yet a Spec` for a Backlog Entry.".
7. MUST, through the `domain-modeling` skill, revise the glossary entry **Spec Queue Table** in `CONTEXT.md` with the sentence of API Contract 1, word for word, after its sentence that ends "or a later counted entry" and before its ADR links. No other glossary entry changes in this Task.
8. MUST replace the `review` values of the entries `command follow-up` and `guide docs/user-guide/commands/follow-up.md` of `docs/references/skill-coverage.json` with the two values of API Contract 1, word for word, as the Coverage Review that `_prd.md` → Skills declares. No other entry changes.
9. MUST, as its last edit, after the last change to the guide, re-record the Behavior Surface Record with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, never by hand.
10. MUST NOT change any skill, the Roundfix follow-up reference, the Release Plan, `docs/references/coverage-record.json`, any Baseline file, the Makefile, a workflow, `go.mod` or a lint configuration.

## Subtasks

- [ ] Give the Task Graph loader a reader with `spec.LoadAt` and its test.
- [ ] Read a Delivery Queue item's Task Graph from its Item Branch in both counts.
- [ ] Add `TestADeliveryItemCountsTasksOnItsItemBranch` with its eight subtests.
- [ ] State the rule in the guide and the glossary entry, record the Coverage Review and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A running or archived Delivery Queue item shows the counts its Item Branch holds, in the Spec queue and under `Pending Specs`.
- [ ] An item without a readable branch shows the Spec Root's counts, an uncommitted status in its worktree is not counted, and the reader writes nothing.
- [ ] `spec.Load` returns what it returned before, and `spec.LoadAt` returns the same graph through a reader.
- [ ] The guide, the glossary entry and the Skill Coverage Map hold the prescribed text, and the Behavior Surface Record is current.

## Context

- instruction: `internal/spec/archive_reader.go`
- instruction: `internal/cli/follow_up_history_test.go`
- instruction: `internal/cli/implement_test.go`
- instruction: `internal/daemon/task_engine.go`
- interface: `internal/spec/spec.go`
- creates: `internal/spec/load_at_test.go`
- interface: `internal/cli/follow_up_queue.go`
- interface: `internal/cli/follow_up_status.go`
- interface: `internal/cli/follow_up_queue_test.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `CONTEXT.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Task Graph loader `spec.Load` and the new `spec.LoadAt` — consumers: `./internal/spec`, `./internal/cli`, `./internal/daemon`
- The Spec queue rows and the `Pending Specs` counts of the status view — consumers: `./internal/cli`
- The `follow-up` guide, the glossary, the Skill Coverage Map and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestLoadAtReadsThroughTheGivenReader$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLoadAtReadsThroughTheGivenReader/same_graph_as_load TestLoadAtReadsThroughTheGivenReader/an_absent_prd_is_not_found TestLoadAtReadsThroughTheGivenReader/an_absent_task_file_is_named; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestADeliveryItemCountsTasksOnItsItemBranch|TestTheQueueRowsComeInTableOrder|TestAQueueItemNamesItsWait|TestFollowUpStatusListsPendingSpecsInDeliveryOrder)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestADeliveryItemCountsTasksOnItsItemBranch/a_running_item_counts_its_branch TestADeliveryItemCountsTasksOnItsItemBranch/an_archived_item_counts_the_source_revision TestADeliveryItemCountsTasksOnItsItemBranch/an_item_without_a_branch_counts_the_spec_root TestADeliveryItemCountsTasksOnItsItemBranch/a_missing_branch_counts_the_spec_root TestADeliveryItemCountsTasksOnItsItemBranch/a_branch_without_the_spec_counts_the_spec_root TestADeliveryItemCountsTasksOnItsItemBranch/the_worktree_is_not_read TestADeliveryItemCountsTasksOnItsItemBranch/the_pending_specs_line_counts_the_same TestADeliveryItemCountsTasksOnItsItemBranch/the_rows_write_nothing TestTheQueueRowsComeInTableOrder TestFollowUpStatusListsPendingSpecsInDeliveryOrder; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task neither new test exists, so their pass lines are missing; afterwards every subtest and the existing tests pass.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/follow-up.md | grep -qF -- "are counted on its Item Branch as Git holds it: the Spec folder at the branch head or, once the branch holds the Spec's Archive Record" || { echo 'follow-up guide lacks the sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/follow-up.md | grep -qF -- "and no file of the item's worktree is read." || { echo 'follow-up guide lacks the worktree clause' >&2; exit 1; }; grep -F -- '**Spec Queue Table**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Spec Queue Table\*\*(:| —( |$))' || { echo 'glossary lacks the entry' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "A Delivery Queue item's Tasks are counted on its Item Branch as Git holds it" || { echo 'glossary entry lacks the sentence' >&2; exit 1; }; grep -qF -- 'Spec 0414 task_01: a Delivery Queue item' docs/references/skill-coverage.json || { echo 'coverage map lacks the review' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 30m ./internal/spec ./internal/cli ./internal/daemon` — expected: exit 0. Before this Task none of the texts exists; afterwards the guide, the glossary and the Skill Coverage Map hold them, the documentation contracts pass with the re-recorded Behavior Surface Record, and the loader's consumer suites pass.

## References

- `_prd.md` → Goals (first bullet); Core Features 1, 2 and 4; Success Metric 1; Skills; Glossary; Decisions
- `_techspec.md` → Interfaces and invariants 1 and 2; API Contract 1; Coverage Map; Testing Approach; Build Order 1
- [`references/2026-11-19-characterization.md`](references/2026-11-19-characterization.md), section 1


## Result

Skills loaded: implement-task, coding-guidelines, domain-modeling, golang-cli, golang-context, golang-error-handling, golang-safety, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

The Task Graph loader now accepts a reader through `spec.LoadAt`; `Load` and `LoadForRecovery` keep their disk readers and their existing error and recovery behavior. Both status counts use the same Delivery Queue reader: resolve the local Item Branch to a commit, load its graph through `git show`, recover an archived graph from the record's `source` at `source_revision`, and otherwise use the checkout's Spec Root. External roots also fall back. No item worktree files are opened and no new output fields or sources are added.

The guide paragraph, the Spec Queue Table definition and exactly the two prescribed Coverage Reviews now contain API Contract 1 verbatim. The Behavior Surface Record was regenerated only through the sanctioned recorder; its diff contains only the follow-up guide fingerprint. The recorder is run again after this Result as the final edit required by Requirement 9.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| Running and archived items count their Item Branch in both views | All eight subtests of `TestADeliveryItemCountsTasksOnItsItemBranch` passed in a focused run. Both views show 1/2 for the running item and 3/3 for the archived source revision with its appended corrective Task. Each fixture starts with two pending Tasks committed on main and uses a linked Item Branch worktree. |
| Fallbacks, uncommitted statuses and read-only behavior | The same test covers no branch, a missing branch and a branch holding neither folder nor record (0/2), leaves Task 02 completed but uncommitted (still 1/2), and compares snapshots of the checkout and home before and after both readers. `TestItemTaskGraphFallsBackForExternalRoots` passed for each external root and refuses any Git call in those cases. Existing table-order, wait, pending-order, no-write, no-judge and message-refusal checks passed unchanged. |
| Loader compatibility and map-backed reads | All three `TestLoadAtReadsThroughTheGivenReader` subtests passed. Graph comparison checks ID, File, Status, Type and Needs after removing the disk fixture, proving the reader supplies every file. Absent paths wrap `os.ErrNotExist`; typed errors identify the absent PRD or Task. Existing recovery, stale-gate refusal and typed-validation checks passed. |
| Prescribed documentation and current record | A read-only Python comparison against API Contract 1 confirmed exact text and insertion positions; comparison with HEAD confirmed that only the Spec Queue Table glossary entry and the two prescribed review values changed. Focused glossary phrase, accepted-ADR and rule-uniqueness checks passed. The sanctioned Behavior Surface recorder exited 0 after the final guide edit. |

### Focused checks

Commands below ran through `rtk proxy`:

- Before the implementation, `go test -count=1 -timeout 5m ./internal/cli -run '^TestADeliveryItemCountsTasksOnItsItemBranch$/a_running_item_counts_its_branch$'` exited 1: the unchanged reader reported `running 0/2 Tasks`, while the branch held 1/2.
- `go test -count=1 -timeout 5m ./internal/spec -run '^(TestLoadAtReadsThroughTheGivenReader|TestLoadForRecovery|TestLoadStillRefusesAStaleGate|TestLoadReturnsTypedValidationErrors)$'` exited 0.
- `go test -count=1 -timeout 5m ./internal/cli -run '^(TestADeliveryItemCountsTasksOnItsItemBranch|TestItemTaskGraphFallsBackForExternalRoots)$'` exited 0.
- `go test -count=1 -timeout 5m ./internal/cli -run '^(TestTheQueueRowsComeInTableOrder|TestAQueueItemNamesItsWait|TestFollowUpStatusListsPendingSpecsInDeliveryOrder|TestFollowUpStatusWritesNothing|TestFollowUpStatusNeverAsksTheJudge|TestFollowUpStatusRefusesMessages)$'` exited 0.
- `go test -count=1 -timeout 5m -tags docscontract ./internal/docscontract -run '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs|StatesEachRuleOnce)$'` exited 0.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0; it is the required implementation recorder, not the declared Verification suite.
- `git diff --check` exited 0 before this Result was appended; the final read-only postflight repeats the whitespace check and confirms the changed-path boundary.

The first combined count test run exposed an unset fixture `specs.root`; the fixture now explicitly configures `docs/specs`, and the focused count runs above passed afterward.

The Daemon-owned `status: in_progress` remains unchanged by this Agent. The Task Graph manifest, every other Task and all prohibited files remain untouched. No commit, push or pull request was made. Declared Task Verification and repository Verification were not run; their execution and Task settlement remain with the Daemon. No follow-up work was added to this slice.

## Carry-forward provenance

- Source Run: `run_20261119T151201Z_05140c8b41c958d4`
- Source commit: `ce6522d8463e7a41a44e8f0dc8a5061616718f0a`
