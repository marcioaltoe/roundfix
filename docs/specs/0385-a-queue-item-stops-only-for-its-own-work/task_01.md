---
task: task_01
spec: 0385-a-queue-item-stops-only-for-its-own-work
status: completed
type: backend
complexity: high
---

# Task 01: The Queue Owner logs its capacity wait and keeps its record when it stops with an item left

## Overview

The production Delivery Queue workflow gets the owner's log, so a repository
gate that waits for the Machine Verification Capacity writes its wait line to
the owner's console log instead of panicking on a nil writer, and every
workflow write goes through one accessor that never returns nil. The owner
then releases its record on exit only when no item can still advance, so an
owner that fails or panics with an item left is reported by `roundfix deliver
status`. This answers the Backlog Entry recorded by the operator on 2026-11-12
(`references/2026-11-12-the-delivery-owner-crashes-when-the-repository-gate-waits-for-machine-verification.md`)
and the message reported by quaestor on 2026-11-12
(`references/2026-11-12-quaestor-the-delivery-owner-panics-at-the-gate.md`).

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   `commandDeliveryEngineDependencies` builds the workflow with a nil `log`,
   and an owner whose engine returns an error while an item is at `gating`
   leaves the queue with owner PID `0`.
2. MUST change `internal/cli/deliver_workflow.go` as `_techspec.md` → API
   Contract 1 says: `commandDeliveryEngineDependencies(runStore, loaded, log
   io.Writer)` sets the workflow's `log` and the engine's `Log` to `log`;
   `newCommandDeliveryEngine` passes `os.Stderr`; add `func (workflow
   *commandDeliveryWorkflow) console() io.Writer`, which returns
   `workflow.log`, or `os.Stderr` when it is nil; and replace every read of
   `workflow.log` in `Gate`, in the item-binary probe and in both Linux Lane
   writers of `internal/cli/deliver_linux_lane.go` with `workflow.console()`,
   removing their own nil fallbacks. After the change, no non-test file of
   `internal/cli` reads `workflow.log` outside `console()`. The wait line keeps
   its text word for word.
3. MUST change every caller of `commandDeliveryEngineDependencies` on the
   default branch when the Task runs to pass a writer: the five test files the
   TechSpec's Declared breaks names pass `io.Discard` and gain the `io` import
   where they lack it; MUST change nothing else in them, and MUST say in the
   `## Result` if it found another caller.
4. MUST change `runDeliveryOwner` in `internal/cli/deliver.go` as API Contract
   2 says: its deferred owner release calls
   `runStore.ReleaseIdleDeliveryQueueOwner(context.WithoutCancel(ctx),
   loaded.GitRoot, pid, identity)` and ignores its result, with a comment
   saying that an owner that fails or panics while an item can still advance
   keeps its record so status, resume, retry and start prove it gone. MUST NOT
   change `deliveryOwnerLine`, `deliver status`, `resume`, `retry`, `start` or
   `stop`.
5. MUST create `internal/cli/deliver_owner_log_test.go` with
   `TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate`,
   `TestADeliveryWorkflowWithoutALogWritesToStandardError` and
   `TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord` (subtests
   `panic`, `error`, `idle`), each calling `t.Parallel()` first, with the
   fixtures and assertions `_techspec.md` → Testing Approach lists. The first
   test asserts the wiring before it runs the gate, so a missing writer fails
   it at once and never waits; its capacity holder is a lease the test takes
   with `store.AcquireMachineVerification` and releases from the log writer's
   first write, never by elapsed time. The second calls `runDeliveryOwner`
   with `&detachChild{}` and a fake engine installed through
   `updateCommandDependenciesForTest`, and recovers the `panic` subtest's
   panic inside the test.
6. MUST change `docs/user-guide/commands/deliver.md` with the two texts of API
   Contract 6, word for word, at the anchors it names.
7. MUST change the **Queue Owner** entry of `CONTEXT.md` through
   `domain-modeling` with the sentence of API Contract 9, word for word, at the
   end of its definition paragraph; MUST add no other capitalized phrase.
8. MUST, as its last step and after the last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST NOT change the wait line's text, the capacity entries, the gate's
   step order, any exit code, the Delivery Engine package, or any test other
   than the calls of requirement 3; no test reads the host's home or real
   configuration, and no test waits on elapsed time.

## Subtasks

- [ ] Record the characterization.
- [ ] Pass the owner's log to the workflow and route every write through `console()`.
- [ ] Keep the record of an owner that stops with an item left.
- [ ] Add the three tests and update the five callers.
- [ ] Add the guide text and the glossary sentence, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A gate built by the production dependencies logs its capacity wait to the owner's log and passes when the unit frees.
- [ ] No workflow write can reach a nil writer.
- [ ] An owner that panics or fails with an item left keeps its record; one that exits idle releases it.
- [ ] The `deliver` guide and the **Queue Owner** entry say so.

## Context

- instruction: `docs/adr/0395-a-queue-item-stops-only-for-its-own-work.md`
- instruction: `docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`
- instruction: `internal/store/machine_verification.go`
- instruction: `internal/store/delivery.go`
- instruction: `internal/cli/machine_verification_status_test.go`
- instruction: `internal/cli/deliver_owner_status_test.go`
- instruction: `internal/cli/deliver_item_binary_test.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- creates: `internal/cli/deliver_owner_log_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_linux_lane.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_archived_resume_test.go`
- interface: `internal/cli/deliver_check_timeout_test.go`
- interface: `internal/cli/deliver_turn_test.go`
- interface: `internal/cli/deliver_review_capacity_test.go`
- interface: `internal/cli/deliver_free_space_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `commandDeliveryEngineDependencies`'s signature, `console()` and the owner's record release — consumers: `./internal/cli`
- The `deliver` guide, the **Queue Owner** entry and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 40m -v -run 'QueueOwner|DeliveryWorkflow|DeliveryEngine|MachineCapacity|RepositoryGate|DeliveryGate|LinuxLane|GateWithoutTheLane|GateSkipsTheLane|CommandFreeSpace|DeliverStart|DeliveryStep|DeliveryItemBinary|DeliverStatus|DeliverResume|DeliverRetry|DeliverStop' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate TestADeliveryWorkflowWithoutALogWritesToStandardError TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the three tests do not exist; afterwards they pass, and every existing Delivery Queue, gate, Linux Lane, item binary and owner test of `./internal/cli` still passes.

- `grep -F -- '**Queue Owner**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Queue Owner\*\*(:| —( |$))' || { echo 'missing term: Queue Owner' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A Queue Owner that stops while an item can still advance keeps that record, so the report names it' || { echo 'Queue Owner not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'because its engine returned an error or the process panicked, keeps its owner record' || { echo 'deliver guide: owner record text missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qE -- 'The queue log is the Queue Owner.s console log' || { echo 'deliver guide: queue log text missing' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the sentences are absent; afterwards the glossary entry and the guide carry them, and the documentation contracts, the glossary phrase contract and the Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; Core Features 1, 2 and 5; Success Metric 1; Success Metric 2; Success Metric 5; Glossary
- `_techspec.md` → Current behavior 1-4, 7; Interfaces; API Contracts 1, 2, 6 and 9; Invariants 1-3; Declared breaks; Testing Approach; Build Order 1
- ADR-0395
- [The adopted Backlog Entry](references/2026-11-12-the-delivery-owner-crashes-when-the-repository-gate-waits-for-machine-verification.md)
- [The adopted quaestor message](references/2026-11-12-quaestor-the-delivery-owner-panics-at-the-gate.md)

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-context, golang-safety, golang-concurrency, golang-error-handling, golang-testing, testing-boss, systematic-debugging, no-workarounds, agentic-cli-design, domain-modeling, tech-writer, evidence-gate.

Starting-tree characterization (before implementation):

- `commandDeliveryEngineDependencies` initializes the workflow without `log`,
  while the engine's `Log` is `os.Stderr`; the capacity wait callback writes
  directly to the nil workflow writer.
- `runDeliveryOwner` defers unconditional `ReleaseDeliveryQueueOwner`, which
  clears the PID and identity even when an engine error leaves an item at
  `gating`. The focused starting-tree regression run below observed PID `0`
  and an empty identity in both `error` and `panic`; `idle` already passed.
- The only pre-existing change was the Daemon's `pending` to `in_progress`
  status transition in this Task file; it is preserved.

Pre-fix focused check:

- `go test -count=1 -timeout 5m -v -run '^(TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate|TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord)$' ./internal/cli`
  exited `1`: the production-wiring assertion found `log:<nil>` before
  running the gate, and `panic` and `error` found PID `0`, empty identity,
  stage `gating`. The `idle` subtest passed. This used the original
  two-argument builder, before its signature or owner cleanup changed.

Implementation and acceptance evidence:

- Capacity wait and gate continuation: the production dependency builder now
  gives the workflow and engine the same writer; the owner constructor passes
  `os.Stderr`. `TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate`
  checks that wiring before entering the gate, acquires a real capacity lease,
  releases it synchronously on the writer's first write, asserts the exact
  wait line naming the specified Run, and finds no capacity entries afterwards.
  Its lane is disabled in the committed fixture and Verification is `exit 0`.
- Nil safety: `console()` supplies the configured writer or `os.Stderr`;
  the gate, item-binary selection and both Linux Lane writers use it.
  `TestADeliveryWorkflowWithoutALogWritesToStandardError` checks both branches.
  A focused source assertion found no production read of `workflow.log` outside
  that accessor, and verified that the capacity wait text is unchanged.
- Owner lifetime: the deferred cleanup uses
  `ReleaseIdleDeliveryQueueOwner(context.WithoutCancel(ctx), ...)`, with the
  specified explanatory comment. The new owner regression invokes
  `runDeliveryOwner` with `&detachChild{}` and an injected engine; `panic`
  and `error` keep the claimed PID and exact identity with the item at `gating`,
  while `idle` merges the item, clears both fields and exits `0`. The panic
  is recovered only in the test. The error still exits `exitPreflight`.
- Documentation: the two guide texts were inserted word for word at their
  prescribed anchors. Through `domain-modeling`, only the Queue Owner
  definition gained the specified sentence and accepted ADR-0395 link.
  Exact-text assertions and the focused glossary phrase/link contracts passed.
- Caller inventory: six existing test calls in exactly the five declared test
  files now pass `io.Discard`; only those calls and the four missing `io`
  imports changed. No additional existing caller was found anywhere in the
  repository; the production constructor and the new regression supply their
  specified writers.

Post-change focused checks (all exited `0`):

- `go test -count=1 -timeout 5m -v -run '^(TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate|TestADeliveryWorkflowWithoutALogWritesToStandardError|TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord)$' ./internal/cli`
  passed all three tests and all three owner subtests.
- `go test -count=1 -timeout 5m -run '^(TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate|TestADeliveryWorkflowWithoutALogWritesToStandardError|TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord|TestRepositoryGateNamesCapacityHoldersBeforeACancelledWait|TestTheDeliveryEngineReceivesTheConfiguredCheckWait|TestTheDeliveryEngineGetsTheTurnRecorder|TestTheDeliveryEngineGetsTheReviewCapacity|TestCommandFreeSpaceNamesEachLocationBelowTheFloor|TestDeliveryStepRunsTheItemBinary|TestGateWithoutTheLaneStartsNoDocker|TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead|TestLinuxLaneSkipsWithANamedReason|TestDeliverStartResumesAnArchivedItemBranchAtReview)$' ./internal/cli`
  passed the 13 selected top-level regressions and compatibility checks after
  the final code and guide edits.
- `go test -count=1 -timeout 5m -tags docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs)$' ./internal/docscontract`
  passed both selected contracts.
- Focused Python source-scope assertions passed for accessor-only log reads,
  call/import-only changes in the five existing tests, unchanged wait text
  and exact guide/glossary wording. `git diff --check` passed.

The new tests use temporary homes/repositories and injected command
environments; no live configuration, provider, Agent session, network or
elapsed-time wait is used. Task status, other Task files and `_tasks.md` were
not edited. Declared Verification and repository Verification remain the
Daemon's responsibility; no commit, push, PR or cross-project message was made.
No follow-up work was added to this slice.

Final implementation step:

- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited `0` after the last guide edit and regenerated
  `docs/references/behavior-surfaces.json` through its sanctioned writer.
