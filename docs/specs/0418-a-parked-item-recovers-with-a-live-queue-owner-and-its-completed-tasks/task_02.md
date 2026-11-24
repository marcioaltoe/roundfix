---
task: task_02
spec: 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
status: completed
type: backend
complexity: high
---

# Task 02: A Delivery Retry ends with a proven live Queue Owner or parks the item again and fails

## Overview

This Task delivers `_techspec.md` → Build Order 2. `roundfix deliver retry` keeps its hand-off to a recorded owner proven alive; otherwise it starts the owner with the owner's report buffered, reads the queue again and proves the owner it records. A failed proof parks the item again as it was, in one guarded store transaction, prints one `Retry failed` report on standard error and exits `2`, so the command never reports an item as retried while nobody holds its queue. It answers the finding reported by quaestor on 2026-11-23, ai-memory message 01a20714-3a6e: the retry answered `Retried ... -> running` and the item then sat at `running` for six hours with no Run and no live Roundfix process. It is verifiable on its own through a store test and four CLI tests, three of them the whole Surface Transcripts 4, 5 and 6.

## Requirements

1. MUST add `RevertDeliveryQueueItemRetry` to the store with the signature of `_techspec.md` → Interfaces and the behavior of API Contract 4: in one write transaction, and only when the queue records no owner and the stored item is at the stage the retry stored with an empty blocker, it sets the stage to `parked` and restores the blocker, the retry count and the candidate commits of the item as read before the retry; it leaves the recorded Run, the Item Branch, the item worktree and the queue's Run links; it returns `false` and changes nothing when either condition fails. No column, table or stage is added.
2. MUST change the end of `runDeliverRetry` to `_techspec.md` → API Contract 3, with `retryOwnerProof` of Interfaces:
   - a recorded owner keeps today's behavior: alive and proven, the item is handed to it and the command exits `0`; proven gone, its record is reclaimed and the command goes on; alive with an identity that cannot be checked, the command fails with today's text;
   - with no recorded owner, or after reclaiming one, the command starts an owner through the injected owner start with what the start writes to standard output buffered, then reads the queue again; when no item can advance it succeeds (Delegated Decision of 2026-11-23 on an owner that already left every item merged or parked); when an item can advance, the queue must record an owner whose PID is alive and whose start-time identity the injected process controller proves;
   - success prints, in today's order, the `Merged outside the queue` line if any, the `Carried forward from Run` lines, the `Retried` line, the `Review:` line if any, then the buffered owner report, and exits `0`;
   - failure prints nothing on standard output, releases a record proven gone and calls `RevertDeliveryQueueItemRetry`, except when the identity of a live PID cannot be checked, where it changes nothing more because that owner counts as present (ADR-0373); it then prints on standard error the report of Surface Transcripts 5 and 6, with one of the five `Reason:` texts of API Contract 3, the `Item:` line read from the Run Database after the command, the `Carried forward:` block only when the command carried Tasks, and the `Next:` line, and exits `2`.
   The start's own diagnostic stays on standard error before the report. No command, flag, exit code or help text changes (invariant 8).
3. MUST hold invariant 5: a retry that fails spends no retry, so the stored retry count after it equals the count before it, and Tasks it carried stay on the Item Branch and are named in the report (ADR-0170).
4. MUST add `internal/store/delivery_retry_revert_test.go` with `TestRevertDeliveryQueueItemRetryRestoresTheParkedItem` and its subtests `reverts_a_retried_item`, `refuses_while_an_owner_is_recorded` and `refuses_when_the_item_moved`, over a temporary Run Database.
5. MUST add `internal/cli/deliver_retry_owner_proof_test.go` with these tests, each whole and named as written:
   - `TestDeliverRetryReproducesSurfaceTranscriptAfterAQAOnlyFailedRun`: Surface Transcript 4, a QA-Only Failed Run with five completed Tasks carried through the production carry-forward, an engine fake that stores the stage as the production engine does, and an owner start fake that claims the queue with the test's own PID and identity and writes the four lines of the owner report.
   - `TestDeliverRetryReproducesSurfaceTranscriptOwnerDidNotStart`: Surface Transcript 5, an owner start fake that writes one diagnostic line to standard error and returns `2`.
   - `TestDeliverRetryReproducesSurfaceTranscriptOwnerEndedBeforeTheProof`: Surface Transcript 6, five carried Tasks and an owner start fake that claims the queue with a PID that does not exist.
   - `TestDeliverRetryProvesTheQueueOwnerItStarted`, with the subtests `different_identity_parks_the_item_again`, `unreadable_identity_leaves_the_item_and_the_record` and `owner_that_already_parked_the_item_exits_zero`.
   Each transcript test asserts every line of standard output and standard error in order and the exit code, copied from the TechSpec with only the declared variations (`<slug>`, `<run-id>`, `<blocker>`, `<owner id>`, `<path>`, `<pid>` and, in Transcript 5, the start's diagnostic lines). The last two transcript tests and the parking subtests also assert the stored item: `parked`, the earlier blocker, the earlier retry count, no recorded owner.
6. MUST update `newOutsideMergeFixture` in `internal/cli/deliver_merged_outside_test.go` so its owner start claims the queue as a live owner does, which keeps `TestDeliverRetryWithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore` passing with its assertions unchanged (`references/2026-11-23-characterization.md`, section 6). Every other existing retry test keeps passing unchanged, the successful retry transcripts among them, which keep their bytes (`_prd.md` → Success Metric 3).
7. MUST revise the glossary entry **Delivery Retry** in `CONTEXT.md` through `domain-modeling` to say that it ends with a Queue Owner proven alive, or fails and parks the item again as it was (`_techspec.md` → API Contract 8). It does not link ADR-0450 yet: task_03 adds that link.
8. MUST keep every new top-level test calling `t.Parallel()` as its first statement; no new test waits on elapsed time, signals a process, guesses a live PID, reaches the network or reads the real home directory.
9. MUST NOT change what `deliver status`, `deliver start` or `deliver resume` print, the help of any command, any skill, user-guide page, decision record or file `_authorization.md` → Limits names.

## Subtasks

- [ ] Add the guarded revert to the store with its test.
- [ ] Buffer the owner start, read the queue again and prove the owner at the end of the retry.
- [ ] Print the failure report and revert the item when the proof fails.
- [ ] Add the three transcript tests and the further proof cases, and adapt the outside-merge fixture's owner start.
- [ ] Revise Delivery Retry in the glossary.

## Acceptance Criteria

- [ ] Whatever the owner start does, `deliver retry` exits `0` only with a Queue Owner proven alive or with no item left to advance, and otherwise exits `2` with no standard output and the failure report on standard error.
- [ ] A failed retry leaves the item parked with its earlier blocker and retry count, unless a live owner's identity cannot be checked, which leaves the item and the record as they are.
- [ ] Surface Transcripts 4, 5 and 6 are reproduced whole, and every existing retry transcript keeps its bytes.
- [ ] The glossary's Delivery Retry says that the retry ends with a proven live owner or fails.

## Context

- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `internal/cli/deliver_owner_status_test.go`
- instruction: `internal/cli/deliver_retry_test.go`
- instruction: `internal/store/delivery_retry_test.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/store/delivery.go`
- interface: `internal/cli/deliver_merged_outside_test.go`
- interface: `internal/cli/deliver_retry_runs_test.go`
- interface: `CONTEXT.md`
- creates: `internal/store/delivery_retry_revert_test.go`
- creates: `internal/cli/deliver_retry_owner_proof_test.go`

## Consumers

- the output and exit code of roundfix deliver retry — consumers: `./internal/cli`
- the store's Delivery Queue item writes, with the new revert — consumers: `./internal/store`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestRevertDeliveryQueueItemRetryRestoresTheParkedItem$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRevertDeliveryQueueItemRetryRestoresTheParkedItem/reverts_a_retried_item TestRevertDeliveryQueueItemRetryRestoresTheParkedItem/refuses_while_an_owner_is_recorded TestRevertDeliveryQueueItemRetryRestoresTheParkedItem/refuses_when_the_item_moved; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the revert restores the parked item and refuses both guarded cases.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDeliverRetryReproducesSurfaceTranscriptAfterAQAOnlyFailedRun|TestDeliverRetryReproducesSurfaceTranscriptOwnerDidNotStart|TestDeliverRetryReproducesSurfaceTranscriptOwnerEndedBeforeTheProof|TestDeliverRetryProvesTheQueueOwnerItStarted|TestDeliverRetryWithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverRetryReproducesSurfaceTranscriptAfterAQAOnlyFailedRun TestDeliverRetryReproducesSurfaceTranscriptOwnerDidNotStart TestDeliverRetryReproducesSurfaceTranscriptOwnerEndedBeforeTheProof TestDeliverRetryProvesTheQueueOwnerItStarted/different_identity_parks_the_item_again TestDeliverRetryProvesTheQueueOwnerItStarted/unreadable_identity_leaves_the_item_and_the_record TestDeliverRetryProvesTheQueueOwnerItStarted/owner_that_already_parked_the_item_exits_zero TestDeliverRetryWithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, Surface Transcripts 4, 5 and 6 are reproduced whole and every proof case passes.
- `grep -F -- '**Delivery Retry**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Delivery Retry\*\*(:| —( |$))' && awk '$0 == "**Delivery Retry**:" { found=1; next } found && $0 == "" { exit } found { print }' CONTEXT.md | tr -s '[:space:]' ' ' | grep -qF -- 'proven alive'` — expected: exit 0, the revised Delivery Retry entry says the retry ends with a Queue Owner proven alive.
- `grep -qF 'func (store *Store) RevertDeliveryQueueItemRetry(' internal/store/delivery.go && grep -qF 'retryOwnerProof' internal/cli/deliver.go && go test -count=1 -timeout 30m ./internal/store ./internal/cli` — expected: exit 0, the consumers of the retry's output and of the store's item writes pass whole.

## References

- `_prd.md` → Goals; User Story 2; Core Features 4 and 7; Success Metric 3; Glossary; Decisions
- `_techspec.md` → The retry's proof; Interfaces; Data Models; API Contract 3; API Contract 4; API Contract 8; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Invariants; Coverage Map; Testing Approach; Build Order 2
- [`references/2026-11-23-characterization.md`](references/2026-11-23-characterization.md), sections 2, 3 and 6
- [`references/sources.md`](references/sources.md), section 1 (reported by quaestor, 2026-11-23)
- ADR-0373, ADR-0398, ADR-0170, ADR-0028, ADR-0395

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, systematic-debugging, no-workarounds, domain-modeling, tech-writer, evidence-gate.

Implemented the retry's post-start owner proof and buffered success output. The existing recorded-owner hand-off keeps its output and identity handling. A new start must leave either no advanceable item or a recorded owner whose live PID and identity are proved. Failed starts and absent or mismatched owners produce the failure report with empty stdout; unreadable live identities retain the owner record and item. The report reads the stored stage and blocker and names carried Tasks and the applicable next command.

Added the guarded store revert in one write transaction. It restores the earlier blocker, retry count and candidate commits only while the item still has the retry's stage, an empty blocker and no recorded owner. It preserves Run provenance, both queue Run links, the Item Branch and worktree. Production carry-forward remains on the branch when a later owner proof fails.

Acceptance evidence:

- Owner proof and exits: all three whole transcript tests and the proof subtests passed in the focused retry run below, including a missing owner record, a different identity, an unreadable identity and an owner that claims the queue, parks the item and releases its idle claim before the proof.
- Restored parked state: the failed-start, gone-owner, missing-record and mismatched-identity cases assert the earlier blocker, count and candidates with no owner. The store test also proves refusal for a recorded owner, moved stage and changed blocker, and checks that the recorded Run, worktree and both Run links survive.
- Transcript compatibility: Transcripts 4, 5 and 6 assert all stdout/stderr bytes and exit codes. Transcripts 4 and 6 carry five Tasks from a QA-Only Failed Run through production carry-forward and inspect their retained implementation and provenance. The outside-merge fixture now claims a live owner through the injected controller; its existing assertions and other existing retry tests are unchanged and passed in the focused run.
- Glossary: revised Delivery Retry through domain-modeling to require a Queue Owner proven alive or failure with the item parked again. The glossary contracts below passed. The ADR-0450 link remains for task_03.

Focused checks:

- Red starting probe: `rtk proxy go test ./internal/store -run '^TestRevertDeliveryQueueItemRetryRestoresTheParkedItem/reverts_' -count=1` failed to compile because the revert method did not exist.
- Final regression check: `rtk proxy go test ./internal/cli ./internal/store -run 'Retry|ReleaseIdleDeliveryQueueOwner' -count=1 -timeout 5m` exited 0; CLI and store reported success, with no skips reported.
- `rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs|ClosedListsEqualTheirConstants|AvoidLinesArePlainLists|StatesEachRuleOnce)$' -count=1 -timeout 3m` exited 0.
- `rtk proxy git diff --check` exited 0.

The first CLI check exposed fixture setup errors: worktree persistence requires its dedicated store method, and the outside-merge fixture's own PID requires the injected proof controller. Both fixtures were corrected without changing existing assertions. The extended store fixture likewise records its branch/worktree pair before setting its parked state.

Scope audit: only this Task's six implementation/test/glossary paths and this Result were changed. The Task status change already existed on entry and was not edited. No declared Verification command, full suite, repository Verification, commit, push or pull request was run. No Task Graph, other Task, skill, user-guide page, ADR or protected tooling file was edited. Daemon Verification and settlement remain pending.

## Carry-forward provenance

- Source Run: `run_20261124T001551Z_9b40dc21bfabfe76`
- Source commit: `88b50e5b5c5b45403566a325f5e07271dd2a639d`
