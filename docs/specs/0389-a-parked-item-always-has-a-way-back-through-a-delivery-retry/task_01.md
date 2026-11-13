---
task: task_01
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
status: completed
type: backend
complexity: medium
---

# Task 01: A Delivery Retry requeues an item parked before it has an Item Branch

## Overview

`roundfix deliver retry` returns a parked item whose recorded Item Branch is
empty to `queued`, counting a retry, as it already does for a prerequisite
park, and the requeue clears a worktree recorded without its branch, so the
Queue Owner creates the branch as on a first start. This answers the Backlog
Entry reported by fiscus on 2026-11-13
([`references/2026-11-13-a-retry-refuses-an-item-parked-before-it-has-a-branch.md`](references/2026-11-13-a-retry-refuses-an-item-parked-before-it-has-a-branch.md)).

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   the production engine's `Retry` of an item parked `delivery-error` with no
   branch refuses with `use item branch: branch is required` and leaves it
   parked with no retry counted.
2. MUST change `Engine.Retry` in `internal/delivery/engine.go` as
   `_techspec.md` → API Contract 1 says: the condition that returns an item to
   `queued` also holds when `strings.TrimSpace(item.Branch)` is empty, with the
   comment it names. A retry of an item with a recorded branch MUST take the
   paths it takes today.
3. MUST change `RetryDeliveryQueueItem` in `internal/store/delivery.go` as API
   Contract 2 says: inside its transaction, when the target stage is `queued`,
   clear `worktree` and `worktree_provisioned` only where the stored branch is
   empty; MUST NOT change the retry count rule or any row with a branch.
4. MUST create `internal/delivery/branchless_retry_test.go` with
   `TestRetryRequeuesAnItemParkedBeforeItHasABranch` and
   `internal/cli/deliver_branchless_retry_test.go` with
   `TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch` (subtests `no
   workspace` and `worktree without branch`), each calling `t.Parallel()`
   first, with the fixtures and assertions `_techspec.md` → Testing Approach
   lists. The CLI test uses the production engine built by
   `commandDeliveryEngineDependencies` over `newItemBranchStartRepository`, whose
   origin is a temporary bare repository; it never reaches the network.
5. MUST change `docs/user-guide/commands/deliver.md` with the table row and
   Text A of API Contract 6 (a), word for word, at the anchors it names.
6. MUST change the **Delivery Retry** entry of `CONTEXT.md` through
   `domain-modeling` with the sentence of API Contract 8, word for word, at the
   end of its definition paragraph; MUST add no other capitalized phrase.
7. MUST, as its last step and after the last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
8. MUST NOT change any exit code, any other park, `UseItemBranch`,
   `CreateItemBranch`, `RecordDeliveryQueueItemWorktree` or any existing test;
   no test reads the host's home or real configuration, reaches the network or
   waits on elapsed time.

## Subtasks

- [ ] Record the characterization.
- [ ] Requeue a branchless item and clear a half-recorded workspace.
- [ ] Add the two tests.
- [ ] Add the guide row, Text A and the glossary sentence, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A retry of an item parked with no Item Branch returns it to `queued` with one counted retry and no recorded workspace.
- [ ] The owner's `CreateItemBranch` then creates a new Item Branch at the default branch, also after a worktree was recorded without its branch.
- [ ] The `deliver` guide and the **Delivery Retry** entry say so.

## Context

- instruction: `docs/adr/0398-a-parked-item-always-has-a-way-back-through-a-delivery-retry.md`
- instruction: `docs/adr/0193-a-spec-names-its-prerequisites-and-the-owner-waits-for-them.md`
- instruction: `internal/delivery/retry_test.go`
- instruction: `internal/cli/deliver_item_branch_test.go`
- instruction: `internal/cli/deliver_recovery_test.go`
- instruction: `internal/cli/deliver_workflow.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- creates: `internal/delivery/branchless_retry_test.go`
- creates: `internal/cli/deliver_branchless_retry_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/store/delivery.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The requeue of `Engine.Retry` and `RetryDeliveryQueueItem` — consumers: `./internal/delivery`, `./internal/store`, `./internal/cli`
- The `deliver` guide, the **Delivery Retry** entry and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run 'Retry|Requeue|DeliveryQueue' ./internal/delivery ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestRetryRequeuesAnItemParkedBeforeItHasABranch " || { echo 'missing pass: TestRetryRequeuesAnItemParkedBeforeItHasABranch' >&2; exit 1; }` — expected: exit 0. Before this Task the test does not exist; afterwards it passes, and every retry and queue test of `./internal/delivery` and `./internal/store` still passes.

- `out="$(go test -count=1 -timeout 40m -v -run 'DeliveryRetry|DeliverRetry|DeliverStart|ItemBranch|Prerequisite' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch/no_workspace TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch/worktree_without_branch; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the test does not exist; afterwards both subtests pass with the production engine, and the existing retry, start and Item Branch tests of `./internal/cli` still pass.

- `grep -F -- '**Delivery Retry**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Delivery Retry\*\*(:| —( |$))' || { echo 'missing term: Delivery Retry' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'An item parked before it has an Item Branch returns to' || { echo 'Delivery Retry not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'An item parked before it has an Item Branch re-enters at' || { echo 'deliver guide: Text A missing' >&2; exit 1; }; grep -qF -- '| Item parked before it has an Item Branch |' docs/user-guide/commands/deliver.md || { echo 'deliver guide: table row missing' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the sentences are absent; afterwards the glossary entry and the guide carry them, and the documentation contracts, the glossary phrase contract and the Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals; User Story 1; Core Features 1 and 5; Success Metric 1; Success Metric 5; Glossary
- `_techspec.md` → Current behavior 1-3; API Contracts 1, 2, 6 and 8; Invariants 1 and 2; Testing Approach; Build Order 1
- ADR-0398
- [The adopted fiscus Backlog Entry](references/2026-11-13-a-retry-refuses-an-item-parked-before-it-has-a-branch.md)

## Result

Skills loaded: implement-task, coding-guidelines, systematic-debugging, no-workarounds, golang-cli, golang-testing, testing-boss, golang-error-handling, domain-modeling, tech-writer, evidence-gate.

Characterization before production edits:

- `go test -count=1 -v ./internal/cli -run '^TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch$'` exited 1. Both `no_workspace` and `worktree_without_branch` called the production engine built by `commandDeliveryEngineDependencies` over a temporary local repository with a bare local origin. Both returned `retry Delivery Queue item "0001-widget-flow": use item branch "": use item branch: branch is required`. The stored item stayed `parked`, blocker unchanged, retry count 0, branch empty. The first had no worktree; the second retained its recorded worktree and `worktree_provisioned=true`.
- `go test -count=1 -v ./internal/delivery -run '^TestRetryRequeuesAnItemParkedBeforeItHasABranch$'` exposed the same premature workspace path through the fake, returning `item worktree is empty` instead of requeueing.

Implementation and acceptance evidence:

1. A branchless parked item requeues with one counted retry and no workspace. `Engine.Retry` now includes the trimmed empty-branch condition after the existing merge and limit checks, bypassing workspace/recovery calls. `RetryDeliveryQueueItem` clears `worktree` and `worktree_provisioned` inside the transaction before the stage update, only for a `queued` target whose stored branch is empty. Its counter rule is unchanged. The two new regression checks below passed after the change; the engine check also asserts no recovery event or workspace use.
2. The owner creates a new Item Branch at the default branch. Both production-engine subtests passed, including the fixture that records a pair, marks it provisioned, then parks with an empty branch. Each asserts the cleared stored workspace and calls the unchanged `CreateItemBranch`, which returns a fresh branch and a worktree whose head equals `origin/main`. All repositories, homes and worktree locations are disposable; origin is local and bare.
3. The guide and glossary describe the new retry path. Added the exact retry table row before the unfinished-Task row, Text A immediately after the table, and the exact definition sentence with the accepted ADR-0398 link at the end of **Delivery Retry** through `domain-modeling`. Focused glossary contracts passed, and the diff was inspected for wording and anchors. No new domain term was added.

Focused checks after the production edits (all exited 0):

- `go test -count=1 -v ./internal/delivery -run '^TestRetryRequeuesAnItemParkedBeforeItHasABranch$'`.
- `go test -count=1 -v ./internal/cli -run '^TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch$'`: both `no_workspace` and `worktree_without_branch` passed.
- `go test -count=1 ./internal/delivery ./internal/store -run '^(TestRetryReturnsADependencyParkToTheQueue|TestRetryCarriesForwardBeforeReenteringTheRun|TestRetryAfterArchiveReentersCheckingWithAPullRequest|TestRetryRefusesAnItemWhoseBranchIsGone|TestRetryDeliveryQueueItemReentersAParkedItem|TestRetryDeliveryQueueItemRefusesAChangedBlocker|TestRetryIncrementsTheItemRetryCount|TestRetryAtTheLimitIsRefusedAndLeavesTheItemUnchanged|TestTheStoreDoesNotCountADiskFullRetry)$'`: existing prerequisite, recorded-branch, archive, blocker, retry-count and retry-limit behavior passed.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs|StatesEachRuleOnce)$'`.
- `go test -count=1 -v -overlay /var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/roundfix-task01-workspace-dbc89q6o/overlay.json ./internal/delivery -run '^TestTask01RequeueKeepsARecordedWorkspace$'`: a temporary test overlay, with no extra repository file, proved a `queued` retry preserves a stored branch/worktree/provisioning pair even when the caller supplies an empty branch; it still counts one retry.
- `git diff --check`.

Scope: only this Task's declared implementation, tests, guide, glossary and derived record, plus this Result. The starting task-file difference was the Daemon's `pending` → `in_progress` transition and was preserved. No existing test or branch/worktree method changed. The skill reference update remains assigned to task_04. Declared Verification, repository Verification and terminal Task status remain Daemon-owned; no commit, push or Pull Request was performed.


Final implementation step: after the last guide edit, ran `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`; exit 0 (`ok roundfix/internal/docscontract`). It regenerated `docs/references/behavior-surfaces.json` through the sanctioned writer. Only this evidence note follows the regeneration; no guide or implementation edit follows it.
