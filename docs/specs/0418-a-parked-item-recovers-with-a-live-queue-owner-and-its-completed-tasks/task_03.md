---
task: task_03
spec: 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
status: completed
type: backend
complexity: high
---

# Task 03: Status reports a queue nobody advances and the Run of each running item, under ADR-0450

## Overview

This Task delivers `_techspec.md` → Build Order 3. `roundfix deliver status` gains an `Owner:` line for a queue that records no owner while an item can still advance, and one `Run:` line per item at `running` that names its Active Run and the Tasks in progress. The new accepted decision record ADR-0450 states the recovery rules, ADR-0279 and ADR-0395 gain their amendment entries and notices, the Standard Export List gains the record, and the glossary's Queue Owner entry and the four entries of this Spec cite it. It answers the finding reported by quaestor on 2026-11-23, ai-memory messages 01a20714-3a6e and 01a20714-ed79 (status printed three lines and no warning for an item at `running` with no owner and no Run), and the item reported by the Praetor on 2026-11-23 (status showed no Task and no Run of a running item, so its sweep reported a false "no progress" on fluxus Spec 0160). It is verifiable on its own through three status transcript tests, a test of the built binary on the state quaestor measured, a decision record test and the ten existing status tests that gain the new lines.

## Requirements

1. MUST change `deliveryOwnerLine` to `_techspec.md` → API Contract 5: when at least one item is neither `merged` nor `parked` and the queue records no owner, status prints `Owner: Delivery Queue has no recorded owner and <n> item(s) can still advance; run roundfix deliver resume` in the place of today's `Owner:` line, `<n>` counting those items. The line of a recorded owner proven gone, the diagnostic of an owner that cannot be checked, the silence for a proven live owner and the silence when no item can advance keep their bytes. Status stays read-only and reclaims nothing (invariant 6).
2. MUST add `internal/cli/deliver_run_line.go` with `deliveryRunLines` of `_techspec.md` → Interfaces and print its lines as API Contract 6 states: for each item at `running`, in queue order, after the `Turn:` lines and before the `Park:` lines, one line, whether or not a live owner is proven. The Active Run is the newest non-terminal Implement Run of the repository whose Spec is the item's and whose local branch is the Item Branch, read with `ListRuns`; its Tasks in progress are the Tasks whose latest `daemon.task` event of phase `started` has no later event of phase `settled` or `skipped`, in the order they started, read with `RunEventsOfKinds`, on the read-only connection status already opens. The three forms are `Run: <slug> no Active Run`, `Run: <slug> <run-id>; Tasks in progress: <task>, <task>` and `Run: <slug> <run-id>; no Task in progress`; a Run in its QA gate prints the last form (Delegated Decision of 2026-11-23 on the QA gate). A read error of the Run Database fails status as its other reads do.
3. MUST keep invariants 1, 2 and 8: every line status prints on the Delivery Base keeps its bytes and its relative order, the item row keeps its four tab-separated fields, and the help of every `deliver` command keeps its bytes. No schema, column, stage, command, flag or exit code is added.
4. MUST add `internal/cli/deliver_status_recovery_lines_test.go` with `TestDeliverStatusReproducesSurfaceTranscriptNoRecordedOwner`, `TestDeliverStatusReproducesSurfaceTranscriptTasksInProgress` and `TestDeliverStatusReproducesSurfaceTranscriptNoTaskInProgress`, which assert Surface Transcripts 1, 2 and 3 each whole: every line of standard output and standard error in order and the exit code, with only the declared variations (`<slug>`, `<item worktree>`, `<run-id>`, `<usage>`), over Runs and `daemon.task` events appended to a temporary Run Database and, for Transcripts 2 and 3, a proven live owner through the injected process controller.
5. MUST add `internal/cli/deliver_quaestor_built_binary_test.go` with `TestTheBuiltBinaryReportsTheQueueQuaestorMeasured`, which builds the binary as `buildRoundfixBinaryForMacro` does and runs it with a temporary `HOME`, a temporary repository and a temporary Run Database, and the subtests:
   - `status_of_the_measured_queue`: the state quaestor measured on 2026-11-23 (`references/sources.md`, section 2), one item at `running` with `retry_count` 1, a recorded Run that ended Failed, no recorded owner and no Active Run, written in the test from that section; the binary's `deliver status` reproduces Surface Transcript 1 whole.
   - `retry_whose_owner_cannot_start`: an item parked with the blocker `delivery-error: run Implement executor: roundfix implement failed with exit code 1`, and `<home>/.roundfix/artifacts` a regular file so the detached start cannot create its Artifact Directory (`references/2026-11-23-characterization.md`, section 3); the binary's `deliver retry` reproduces Surface Transcript 5 whole, and the stored item is `parked` with its earlier blocker and retry count and no recorded owner.
   The test reads no file under the Spec Root, reaches no network, starts no Agent and reads nothing of the real home directory.
6. MUST update the ten existing status tests that gain a line, each only by the line `references/2026-11-23-characterization.md`, section 5, lists for it: `TestDeliverStatusPrintsTheItemWorktree`, `TestDeliverStatusPrintsAnItemWarning`, `TestDeliverStatusPrintsNoWarningLineWithoutAWarning`, `TestDeliverStatusPrintsNoRunsRecorded`, `TestDeliverStatusPrintsAParkLinePerParkedItem`, `TestDeliverStatusShowsUndeclaredCorrectionSigns` and `TestDeliverStatusNamesTheMachineCapacityHolderAndWaiters` gain the `Owner:` line; `TestDeliverStatusReproducesSurfaceTranscriptQueueTurns` gains one `Run:` line; `TestDeliverStatusShowsARunSlotWaitTranscript` gains both; and `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner` is renamed `TestDeliverStatusPrintsTheOwnerLineOfAQueueWithNoRecordedOwner`, with the subtests `an_item_can_advance` (the new line) and `no_item_can_advance` (no line for a queue with no recorded owner whose items are all merged or parked). Every other existing status test keeps passing unchanged.
7. MUST create `docs/adr/0450-a-delivery-retry-ends-with-a-live-queue-owner-and-status-reports-a-queue-nobody-advances.md`, the decision record of `_techspec.md` → API Contract 7, with the title `A Delivery Retry ends with a live Queue Owner, and status reports a queue nobody advances`, the lifecycle header of `docs/agents/docs-layout.md` with `status: accepted` and, on the next line, `knowledge: true`, and the sections `## Context`, `## Decision` and `## Considered Options`. It states API Contracts 1 to 6 as rules, the two options this design gives up (a record at the Run's start and a restart of the owner by Roundfix) and that it amends ADR-0279 and ADR-0395 in part; it cites the findings of quaestor and of the Praetor of 2026-11-23 by their content, names no client and cites no Spec folder path. When 0450 is no longer the next free number on the tree this Task runs on, stop and report it; do not take another number.
8. MUST write by hand, with no `roundfix baseline update` (`_techspec.md` → Decisions):
   - in `docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`, the header line `amended_by: [ADR-0450]` after `knowledge: true`, `updated_at` set to ADR-0450's `created_at`, and one paragraph directly after the paragraph that ends with the sentence "A queue with no recorded owner adds no line.", which begins `**Amended by [ADR-0450](0450-a-delivery-retry-ends-with-a-live-queue-owner-and-status-reports-a-queue-nobody-advances.md) (<date>).**`, `<date>` being the date of ADR-0450's `created_at`, and says what changed;
   - in `docs/adr/0395-a-queue-item-stops-only-for-its-own-work.md`, the same header changes and the same notice directly after the paragraph that begins "**An owner that stops with an item left keeps its record.**";
   - in `.secondbrain-export`, the path of ADR-0450 once, inside the Standard Export List, in ascending order, directly after the line of the decision record numbered before it.
   No other line of the two ADRs changes, both stay `accepted`, and no other accepted ADR changes.
9. MUST add `internal/docscontract/delivery_recovery_decision_record_test.go`, with the `docscontract` build tag and the `verify:always` marker the package's other files carry, holding `TestTheDeliveryRecoveryDecisionAmendsTheOwnerStatusRules`. As `TestTheIndexusDecisionAmendsTheConsultationRule` does for its record, it finds the record by its title, never by its number (invariant 7), and asserts with `spec.ReadADRHeader`: accepted with a true Knowledge Mark; ADR-0279's and ADR-0395's `amended_by` hold its identifier; each notice follows its sentence and carries the date of the record's `created_at`; `.secondbrain-export` holds the record's path once between the Standard Export List markers, in ascending order.
10. MUST revise the glossary through `domain-modeling`: **Queue Owner** says that status also reports a queue with no recorded owner while an item can still advance, and the Run of each item at `running`; and each of the entries **QA-Only Failed Run**, **Task Carry-Forward**, **Delivery Retry** and **Queue Owner** cites ADR-0450 with a link to its file (`_techspec.md` → API Contract 8).
11. MUST keep every new top-level test calling `t.Parallel()` as its first statement; no new test waits on elapsed time, signals a process it did not start, reaches the network or reads the real home directory.
12. MUST NOT change any skill, user-guide page, the Behavior Surface Record, the Skill Coverage Map, any other guide, `docs/agents/setup-context.json`, or any file `_authorization.md` → Limits names.
13. MUST re-record the repository test coverage record once every test of this Task is in place, because requirement 6 renames `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner` on purpose (ADR-0450 inverts the rule it pinned) and `TestCoverageEquivalence` in `internal/spec` reports that removal as a regression: run `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1` from the tree this Task runs on and commit `docs/references/coverage-record.json`. The only removal the record may show is that one test name; every other change is an addition. (Correction of 2026-11-23: Run run_20261124T001551Z_9b40dc21bfabfe76 failed this Task on make verify-changed with that single coverage regression.)
14. MUST raise the count of `internal/docscontract` contract files that `TestRepositoryContractTestsDeclareTheirRelevance` in `internal/verifyselect/contracts_repository_test.go` pins (`len(docsFiles) != 11`) by one, for the file requirement 9 adds, and change nothing else in that test. (Correction of 2026-11-24: Run run_20261124T030803Z_0c345fdfd0f3383d failed this Task on make verify-changed with "missing declarations: ... docs files=12".)

## Subtasks

- [ ] Print the no-recorded-owner `Owner:` line and the `Run:` lines.
- [ ] Add the three status transcript tests and the built-binary test on the state quaestor measured.
- [ ] Update the ten existing status tests by the lines they gain and rename the no-owner test.
- [ ] Write ADR-0450, the two notices with their headers and the export line, and the decision record test.
- [ ] Revise Queue Owner and link ADR-0450 from the four glossary entries.

## Acceptance Criteria

- [ ] With an item that can advance and no recorded owner, status prints the new `Owner:` line naming `roundfix deliver resume`; every other owner state prints what it prints on the Delivery Base.
- [ ] Every item at `running` has one `Run:` line in one of the three forms, after the `Turn:` lines and before the `Park:` lines, and the item row and every existing line keep their bytes.
- [ ] The built binary reproduces Surface Transcript 1 on the state quaestor measured and Surface Transcript 5 on a retry whose owner cannot start.
- [ ] ADR-0450 is accepted with a true Knowledge Mark, ADR-0279 and ADR-0395 name it with dated notices, the Standard Export List holds it, and the four glossary entries cite it.

## Context

- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `internal/cli/deliver_turn.go`
- instruction: `internal/cli/implement_test.go`
- instruction: `internal/store/run_events_of_kinds.go`
- instruction: `internal/docscontract/indexus_decision_record_test.go`
- instruction: `internal/spec/adr_header.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_owner_status_test.go`
- interface: `internal/cli/deliver_test.go`
- interface: `internal/cli/deliver_revalidate_test.go`
- interface: `internal/cli/deliver_usage_test.go`
- interface: `internal/cli/deliver_park_status_test.go`
- interface: `internal/cli/deliver_corrections_status_test.go`
- interface: `internal/cli/machine_verification_status_test.go`
- interface: `internal/cli/deliver_turn_test.go`
- interface: `internal/cli/deliver_run_slot_wait_test.go`
- interface: `docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`
- interface: `docs/adr/0395-a-queue-item-stops-only-for-its-own-work.md`
- interface: `.secondbrain-export`
- interface: `CONTEXT.md`
- interface: `docs/references/coverage-record.json`
- interface: `internal/verifyselect/contracts_repository_test.go`
- creates: `internal/cli/deliver_run_line.go`
- creates: `internal/cli/deliver_status_recovery_lines_test.go`
- creates: `internal/cli/deliver_quaestor_built_binary_test.go`
- creates: `internal/docscontract/delivery_recovery_decision_record_test.go`
- creates: `docs/adr/0450-a-delivery-retry-ends-with-a-live-queue-owner-and-status-reports-a-queue-nobody-advances.md`

## Consumers

- the output of roundfix deliver status — consumers: `./internal/cli`
- this repository's decision records, the glossary and the Standard Export List — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDeliverStatusReproducesSurfaceTranscriptNoRecordedOwner|TestDeliverStatusReproducesSurfaceTranscriptTasksInProgress|TestDeliverStatusReproducesSurfaceTranscriptNoTaskInProgress|TestDeliverStatusPrintsTheOwnerLineOfAQueueWithNoRecordedOwner)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverStatusReproducesSurfaceTranscriptNoRecordedOwner TestDeliverStatusReproducesSurfaceTranscriptTasksInProgress TestDeliverStatusReproducesSurfaceTranscriptNoTaskInProgress TestDeliverStatusPrintsTheOwnerLineOfAQueueWithNoRecordedOwner/an_item_can_advance TestDeliverStatusPrintsTheOwnerLineOfAQueueWithNoRecordedOwner/no_item_can_advance; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, Surface Transcripts 1, 2 and 3 are reproduced whole and the renamed owner test passes both subtests.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheBuiltBinaryReportsTheQueueQuaestorMeasured$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheBuiltBinaryReportsTheQueueQuaestorMeasured/status_of_the_measured_queue TestTheBuiltBinaryReportsTheQueueQuaestorMeasured/retry_whose_owner_cannot_start; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the built binary reproduces Surface Transcripts 1 and 5 on the states quaestor measured.
- `out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^TestTheDeliveryRecoveryDecisionAmendsTheOwnerStatusRules$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheDeliveryRecoveryDecisionAmendsTheOwnerStatusRules ' || { echo 'missing pass: TestTheDeliveryRecoveryDecisionAmendsTheOwnerStatusRules' >&2; exit 1; }` — expected: exit 0, ADR-0450 amends ADR-0279 and ADR-0395 and joins the Standard Export List.
- `f=docs/adr/0450-a-delivery-retry-ends-with-a-live-queue-owner-and-status-reports-a-queue-nobody-advances.md; test -f "$f" && grep -qx 'status:.accepted' "$f" && grep -qx 'knowledge:.true' "$f" && grep -q '^amended_by:.*ADR.0450' docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md && grep -q '^amended_by:.*ADR.0450' docs/adr/0395-a-queue-item-stops-only-for-its-own-work.md && grep -qxF -- "$f" .secondbrain-export` — expected: exit 0, the record is accepted and marked, both amended ADRs list it and the export list holds its path.
- `grep -F -- '**Queue Owner**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Queue Owner\*\*(:| —( |$))' && awk '$0 == "**Queue Owner**:" { found=1; next } found && $0 == "" { exit } found { print }' CONTEXT.md | tr -s '[:space:]' ' ' | grep -qF -- 'no recorded owner' && for term in 'QA-Only Failed Run' 'Task Carry-Forward' 'Delivery Retry' 'Queue Owner'; do awk -v header="**$term**:" '$0 == header { found=1; next } found && $0 == "" { exit } found { print }' CONTEXT.md | grep -qF -- '0450-a-delivery-retry' || { echo "entry does not cite ADR-0450: $term" >&2; exit 1; }; done` — expected: exit 0, Queue Owner names a queue with no recorded owner and the four entries cite ADR-0450.
- `grep -qF 'func deliveryRunLines(' internal/cli/deliver_run_line.go && go test -count=1 -timeout 30m ./internal/cli && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0, the consumers of the status output and of the decision records, the glossary and the export list pass whole, the ten moved status tests among them.

## References

- `_prd.md` → Goals; User Stories 3 and 4; Core Features 5, 6, 7 and 8; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Glossary; Decisions
- `_techspec.md` → The status lines; The decision record, the glossary, the guide and the skill; Interfaces; API Contract 5; API Contract 6; API Contract 7; API Contract 8; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 5; Invariants; Coverage Map; Testing Approach; Build Order 3; Decisions
- [`references/2026-11-23-characterization.md`](references/2026-11-23-characterization.md), sections 3, 4, 5 and 8
- [`references/sources.md`](references/sources.md), sections 1, 2 and 3 (reported by quaestor and by the Praetor, 2026-11-23)
- ADR-0279, ADR-0395, ADR-0365, ADR-0304, ADR-0373

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, systematic-debugging, no-workarounds, roundfix.

Implemented this Task's status and durable decision slice. The initial worktree contained only the Daemon's Task status change; that field remains untouched by the Agent.

Acceptance evidence:

- Owner status: `deliveryOwnerLine` counts all items outside `merged` and `parked`, names `roundfix deliver resume` when none has a recorded owner, and preserves the recorded-owner proof paths. The renamed test covers both advanceable and idle queues. Existing owner-state and read-only queue assertions pass in the focused status selection.
- Run status: `deliveryRunLines` reads `ListRuns` and `RunEventsOfKinds` through status's existing reader, selecting the newest active Implement Run for the repository, Spec and Item Branch. It renders one line per running item after Turns and before Parks, ordering Tasks by their latest outstanding start. The three whole transcript tests first reproduced the missing lines on the unchanged implementation, then passed with the implementation. An additional focused test covers newer mismatched/terminal Runs, repository and branch filtering, queue order, restarted Tasks, skipped Tasks and closed-database failure. The ten existing affected tests gained their declared lines; correction signs retain their position before Owner.
- Built executable: both `status_of_the_measured_queue` and `retry_whose_owner_cannot_start` pass through the real built binary with temporary HOME, repository and database. The measured status fixture has retry count 1, a Failed Run, five completed Tasks, failed QA, nine prompts and the reported token total. The failed-start fixture uses a regular artifacts file and asserts the whole startup diagnostic plus retry failure transcript, no stdout, exit 2, original parked blocker, original nonzero retry count and no owner. Retry counts are seeded with the real retry store operation, because ordinary item updates do not write that column.
- Decision and glossary: ADR-0450 was confirmed the next free number before creation. It is accepted and knowledge-marked, records the recovery contracts and rejected alternatives, and cites the quaestor and Praetor observations without a client or Spec path. Both earlier ADRs retain all prior bytes except their authorized header changes and dated notices. The title-based docscontract test checks their amendments and the ordered, unique Standard Export List entry. All four glossary entries link the decision, and Queue Owner describes the new status behavior.

Focused checks run:

- `go test ./internal/cli -run '^(TestDeliverStatus|TestDeliveryRunLines|TestTheBuiltBinaryReportsTheQueueQuaestorMeasured)' -count=1` — exit 0.
- `go test -tags docscontract ./internal/docscontract -run 'TestTheDeliveryRecoveryDecision|TestGlossary' -count=1` — exit 0.
- `go test ./internal/verifyselect -run '^TestRepositoryContractTestsDeclareTheirRelevance$' -count=1` — exit 0; only its pinned docscontract file count changed from 11 to 12.
- Authorized regeneration: `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1` — exit 0. Compared the regenerated JSON with HEAD: the only removed test is `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner`; the other 68 entries are additions, including previously unrecorded tests already in this tree.
- `git diff --check` — exit 0. A byte comparison of each earlier ADR against HEAD, excluding the authorized timestamp, amendment header and notice, found no other changes.

No declared Verification command or repository Verification was run. Task status, Task Graph, other Task files, skills, user guide, Baseline guides and other generated records were not edited. No commit, push or pull request was made. Daemon Verification and settlement remain pending.

## Carry-forward provenance

- Source Run: `run_20261124T033912Z_c4dbe04fec2c7007`
- Source commit: `579cf57a573f74abaac66523b4f8b5cbab461262`
