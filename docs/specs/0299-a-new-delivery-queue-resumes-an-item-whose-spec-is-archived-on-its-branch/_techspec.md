---
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
prd: _prd.md
created: 2026-10-22
---

# A new Delivery Queue resumes an item whose Spec is already archived on its branch — Technical Spec

## Executive Summary

A new queue's item starts at `queued`, and `Engine.advanceItem` always runs
Delivery Revalidation after `CreateItemBranch` reuses the continued item
branch, so a branch that already archived its Spec parks `delivery-error` on
the missing Spec folder. A new optional engine dependency proves, from the
item worktree, whether the branch already archives the Spec: the commit that
removed the Spec's `_prd.md` must pass the queue's existing archive proof. A
proven archive records the item head as the candidate and enters `reviewing`,
skipping revalidation and the Implement Run; the archive stage then finds the
Spec archived and commits nothing. Anything else parks the item with the new
blocker `archived-item-unproven: <cause>`, whose retry returns it to `queued`.
The queue also resolves an item's archive paths when the item worktree has no
Specs Root directory, which happens when the archived Spec was the only active
one.

The trade-off accepted: the resumed head is reviewed again under the
configured policy, at the cost of one more provider review, because a new
queue holds no record of the earlier queue's review and an exact archive
commit does not prove one (ADR-0312).

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema, table or record
  field is minted; the blocker text lives in the existing item blocker field.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the existing push, pull request and merge
  calls carry a resumed item, and every test uses temporary repositories,
  homes and fake forges. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Orders 1, 2 and 3 follow
  ADR-0223, ADR-0229, ADR-0237, ADR-0240 and ADR-0312, and the QA gate
  follows ADR-0104, ADR-0167 and ADR-0276, with ADR-0295, all named with their titles in
  `_prd.md` →
  Project Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AO approval of 2026-10-22 ("the same rules as cycle AN"), carrying the cycle
  AJ grant "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-22 at `34b4e246` (v0.107.0). Line numbers
move when other Specs merge; the function names do not.

1. `runDeliverStart` (`internal/cli/deliver.go`) loads each Spec from the
   user checkout, finds the one item branch that holds commits the default
   branch lacks with `existingItemBranches` (ADR-0223), records a new queue
   (`store.CreateDeliveryQueueWithLimits` deletes a terminal queue, so no
   candidate, Run, or Pull Request number survives) and prints `Continuing
   item branch <branch> for <slug>`.
2. `Engine.advanceItem` (`internal/delivery/engine.go`), for a `queued` item,
   checks prerequisites, calls `CreateItemBranch`
   (`internal/cli/deliver_workflow.go`), which reuses that branch and its
   worktree, then calls `revalidator.Revalidate` without condition and moves
   the item to `running`. `commandDeliveryWorkflow.Revalidate`
   (`internal/cli/deliver_revalidate.go`) resolves the Specs Root with
   `roundconfig.ResolveSpecsRoot` and runs `strictSpecFindings`, which reads
   the active Spec folder; either error is returned, and `Engine.Run` parks
   the item `delivery-error: revalidate item: …`.
3. Reproduced through `roundfix deliver start` in a temporary repository: a
   continued branch holding the Archive Command's archive commit and a Pull
   Request Record Commit for #619 parked `delivery-error: revalidate item:
   run strict Spec Consistency Check: read Spec "0001-widget-flow": stat
   …/docs/specs/0001-widget-flow: no such file or directory`; when the
   archived Spec was the only active one, the fresh item worktree had no
   `docs/specs` and the park was `delivery-error: revalidate item: resolve
   item Specs Root: specs.root resolved to "…/docs/specs", which does not
   exist`. No Implement Run started in either case.
4. `Engine.Retry` of that park reaches `InspectItem`, which reports the item
   `Archived`; with no recorded candidate and no QA override it refuses
   `candidate head is missing` (ADR-0229).
5. The archive proof already exists: `reconcileArchiveCommit` adopts an
   archive commit whose parent is the reviewed head when
   `archiveCommitIsExact(parent, head, source, destination)` holds; that
   function reads only the change (Archive Record retirement through
   `archiveRetirementIsExact`, or a legacy folder move), never the subject.
6. `commandDeliveryWorkflow.archivePaths` resolves the Specs Root with
   `roundconfig.ResolveSpecsRoot`, which fails when the directory is absent.
   `Publication`, `Authorization`, `RecordPullRequest`, `Archive` and
   `qaEnvironmentPartial` call it, so every stage after review fails in an
   item worktree whose only Spec was archived. `InspectItem` already uses
   `resolveHistoricalSpecsRoot` (`internal/cli/archive_reader.go`), which
   returns the configured root without requiring the directory.
7. Once an item holds a candidate at `reviewing`, the stages already handle an
   archived head: `Archive` returns `AlreadyArchived` when the reviewed head
   holds the archive and not the folder, and `archiveCandidate` moves to
   `gating` with no commit; `RecordPullRequest` makes no commit when the
   record already names the Pull Request; `FindOrCreatePullRequest` finds an
   open Pull Request by head branch.

## System Architecture

- `internal/delivery/engine.go`: `ArchivedItemProof`, `ArchivedItemProver`,
  `EngineDependencies.ArchivedItems`, the `queued` branch of `advanceItem`,
  the blocker `BlockerArchivedItemUnproven` and its `Retry` rule.
- `internal/delivery/park_class.go`: the blocker's Park Class and next action.
- `internal/cli/deliver_workflow.go`: `commandDeliveryWorkflow.ProveArchivedItem`,
  its wiring in `newCommandDeliveryEngine`, and `archivePaths` through
  `resolveHistoricalSpecsRoot`.
- `docs/user-guide/commands/deliver.md`, the Roundfix `deliver` reference and
  its mirror, `CONTEXT.md`: the resume, the blocker and the glossary.

No package, import edge or file outside those changes in production.

## Implementation Design

### Interfaces

```go
// ArchivedItemProof is what a queued item's worktree proves about an archive
// its Item Branch already holds.
type ArchivedItemProof struct {
	Archived      bool   // the head holds a proven archive of the Spec
	ArchiveCommit string // the commit that removed the Spec's _prd.md
	Head          string // the item worktree head
	Refusal       string // non-empty: the cause the item cannot be proven
}

type ArchivedItemProver interface {
	ProveArchivedItem(ctx context.Context, workDir, specSlug string) (ArchivedItemProof, error)
}
```

### Data Models

None. The Delivery Queue item, its blocker field, the Archive Record and the
review record keep their shapes.

### API Contracts

1. API Contract: the engine dependency. `EngineDependencies.ArchivedItems`
   holds an `ArchivedItemProver`; a nil value keeps today's `queued` path
   byte for byte. `newCommandDeliveryEngine` wires the command workflow, and
   `var _ delivery.ArchivedItemProver = (*commandDeliveryWorkflow)(nil)`
   holds the implementation to it.
2. API Contract: the `queued` path. After `CreateItemBranch` records the
   branch and worktree, and before `Revalidate`, `advanceItem` calls
   `ProveArchivedItem(ctx, item.Worktree, item.SpecSlug)`:
   - an error returns `prove archived item: <error>`, which parks
     `delivery-error` as today;
   - a non-empty `Refusal` parks the item
     `archived-item-unproven: <Refusal>` and returns;
   - `Archived` with an empty `Head` returns `prove archived item: archived
     item head is empty`;
   - `Archived` writes `roundfix: Delivery Queue item <slug>: item branch
     <branch> already archives the Spec at <commit>; resuming at reviewing`
     to the engine log, sets the candidate commits to exactly `[Head]` and
     sets the stage `reviewing`; neither `Revalidate` nor `RunSpec` is
     called, and the stage loop continues from `reviewing`;
   - otherwise `Revalidate` and the move to `running` run as today.
3. API Contract: the blocker. `BlockerArchivedItemUnproven =
   "archived-item-unproven"`. `ClassifyPark` gives it `ParkClassFinding` and
   the next action `in <worktree>, make the item branch hold the Spec folder
   or one exact archive commit of the Spec, then run roundfix deliver retry
   <slug>`, where `<worktree>` is the item's recorded worktree.
4. API Contract: its retry. `Engine.Retry` of an item parked
   `archived-item-unproven` or `archived-item-unproven: <cause>` keeps every
   check that precedes the dependency-park rule (merge observation, token
   ceiling, queue deadline, retry limit) and then, like a
   `prerequisite-unmerged` park, stores the item at `queued` through
   `RetryDeliveryQueueItem`, which counts one retry; it uses no workspace,
   recovery or revalidation.
5. API Contract: the proof. `commandDeliveryWorkflow.ProveArchivedItem`
   resolves `source` and `destination` with `archivePaths(workDir, slug)`,
   reads `head` with `git rev-parse HEAD`, and, with `<h12>` and `<c12>` the
   first 12 characters of a commit as `deliveryCommitName` gives them:
   - when `<head>:<source>/_prd.md` exists, returns the zero proof;
   - when neither `<head>:<destination>.md` nor
     `<head>:<destination>/_prd.md` exists, refuses `item head <h12> holds
     neither <source> nor an archive of the Spec`;
   - when `git rev-list -n 1 <head> -- <source>/_prd.md` prints nothing,
     refuses `no commit on the item branch removed <source>/_prd.md`;
   - otherwise, with `C` that commit and `P` the full hash of `C^`, refuses
     `commit <c12> that removed <source> is not an exact archive of the Spec`
     unless `archiveCommitIsExact(ctx, workDir, P, C, source, destination)`
     holds, and then returns `Archived`, `ArchiveCommit: C`, `Head: head`.
   A Git failure is returned as an error. It writes nothing.
6. API Contract: the Specs Root of an archived item. `archivePaths` resolves
   the root with `resolveHistoricalSpecsRoot`, so an item worktree without
   its Specs Root directory still yields the configured source and archive
   paths; with the directory present the result is unchanged.
7. API Contract: documentation. The `deliver` guide and the Roundfix
   `deliver` reference state the resume, the blocker with its next action,
   and its retry; `CONTEXT.md` adds **Archived Item Resume** and revises
   **Delivery Revalidation**.

### Surface Transcripts

None. `roundfix deliver start` prints the same lines; `roundfix deliver status`
prints the existing `Park: <slug> <class>: <next command>` line, whose text
for the new blocker API Contract 3 fixes and the tests assert.

### Invariants

```text
1. A new queue never starts an Implement Run or makes an archive commit for an item whose head holds a proven archive of its Spec.
2. A resumed item's first candidate is its worktree head, and it is reviewed under the configured policy before any push.
3. An item whose head holds the Spec folder is revalidated and run exactly as before.
4. An item whose head lacks the Spec folder and holds no proven archive parks archived-item-unproven before any Run, review or push, with the cause in the blocker.
5. The archive proof reads the change of the commit that removed the Spec's _prd.md, never its subject.
6. A retry of archived-item-unproven returns the item to queued and counts one retry.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1, Invariants 1-2 → `advanceItem` `queued` path, `ProveArchivedItem` (API Contracts 1, 2, 5)
- Core Feature 3, Success Metric 2 → `archivePaths` (API Contract 6)
- Goal 2, Core Feature 1, Invariant 2 → the unchanged `reviewing`, `archiving`, `gating`, `publishing`, `checking` and `merging` stages (API Contract 2)
- Goal 3, Core Feature 2, Success Metrics 3-4, Invariants 4-6 → the blocker, its Park Class and its retry (API Contracts 3-5)
- Goal 4, Success Metric 5, Invariant 3 → the zero proof (API Contracts 2, 5)
- Core Feature 4 → the guide, the reference and the glossary (API Contract 7)
- Success Metric 6 → the PRD's Unreachable Acceptance

## Absence Evidence

1. `ArchivedItemProof`, `ArchivedItemProver`, `ProveArchivedItem` and `BlockerArchivedItemUnproven` — command: `grep -rn -e ArchivedItem -e archived-item-unproven internal`; result: no match at `34b4e246`; the `queued` path calls `Revalidate` unconditionally and no blocker names an archived item.
2. A resume of an archived item in a new queue — command: `grep -n -e Revalidate -e InspectItem internal/delivery/engine.go`; result: at `34b4e246` `InspectItem` is called only by `Retry`, and the `queued` branch of `advanceItem` calls `Revalidate` after `CreateItemBranch` with no archive check.

## ADR Alternatives

- ADR-0104: no rejected alternative is taken; the outside-evidence row rests on this repository's history and published documentation, recorded with their sources.
- ADR-0167: no rejected alternative is taken; the pre-PR Pull Request row stays blocked and never decides a qualifying partial.
- ADR-0223: no rejected alternative is taken; the archive commit stays an exact move that no archive rewrites, the continuation rule is unchanged, and a post-archive correction is still reviewed before publication.
- ADR-0229: no rejected alternative is taken; the queue never archives a partial, and the retry rule for an archived item without a candidate is unchanged — a new queue's resume rests on the archive's exactness, not on a QA override.
- ADR-0237: no rejected alternative is taken; a Delivery Retry still asks first whether the item merged, no owner pass probes a parked item, and an `archived-item-unproven` park is answered only by an explicit retry or a new queue.
- ADR-0240: no rejected alternative is taken; the QA partial policy is untouched, and the resume runs no QA gate and reads no QA verdict.
- ADR-0276: no rejected alternative is taken; the row the authorization keeps offline is recorded as such, and every byte invariant of the QA gate compares with the Delivery Base.
- ADR-0295: no rejected alternative is taken; the archive proof reads only the item branch's own history, and no locator or default-branch reachability rule changes.
- ADR-0312: no rejected alternative is taken; the resumed head is reviewed rather than sent to the gate, nothing is revalidated from the record, and no QA override is required.

## Integration Points

Git, read in the item worktree for the proof; the forge, through the
unchanged push, pull request and merge calls of a resumed item. No test,
Verification command or QA row reads the real `~/.roundfix`, another
repository or a forge.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary repositories, homes and Git configuration, and waits on no real
time.

- `internal/delivery/archived_item_resume_test.go` (new), with the package's
  fake workflow and boundary and a fake prover:
  - `TestANewQueueResumesAProvenArchivedItemAtReview`: the prover returns
    `Archived` with a head; the item reaches `merged` through the fakes,
    `RunSpec` and `Revalidate` are never called, the first review and the
    archive call receive that head, the candidate commits start with it, and
    the log holds the resume line of API Contract 2.
  - `TestANewQueueParksAnUnprovenArchivedItem`: a refusal parks
    `archived-item-unproven: <refusal>` with no Run, revalidation or review,
    and `ClassifyPark` gives `finding` and the next action of API Contract 3.
  - `TestANewQueueRevalidatesAnItemWithoutAnArchive`: the zero proof calls
    `Revalidate` once and reaches `RunSpec`; a nil prover behaves the same.
  - `TestANewQueueParksAnArchivedItemProofError`: a prover error parks
    `delivery-error: prove archived item: <error>`.
  - `TestARetryReturnsAnUnprovenArchivedItemToQueued`: the parked item is
    stored at `queued` with its retry count raised by one, without a call to
    the workspace, recovery or revalidator.
  - `TestParkClassesAndNextActions` (existing, `park_class_test.go`) gains
    the blocker's row.
- `internal/cli/deliver_archived_resume_test.go` (new), through `roundfix
  deliver start` with `deliveryReadiness` and `startDeliveryOwner` replaced,
  a bare origin, the Archive Command's real archive, and an owner whose
  workspace, recovery, history, revalidator, archiver, recorder and prover
  are the command workflow while the runner, reviewer, gate, authorizer,
  publication and pull requests are a fake forge holding open Pull Request
  #619 for the item branch:
  - `TestDeliverStartResumesAnArchivedItemBranchAtReview`, subtests
    `another_active_Spec` and `only_active_Spec`: a branch with a completed
    Task and passing QA, the archive commit and a Pull Request Record Commit
    for #619 ends `merged` with #619, no Implement call, one review of the
    record commit's head, no Pull Request created, and an item head equal to
    that commit.
  - `TestDeliverStartParksAnUnprovenArchivedItemBranch`, subtests
    `inexact_archive_commit` (the archive commit also adds an unrelated
    file), `folder_removed_without_archive` (the Spec folder removed by
    `git rm` with no archive) and `no_removing_commit` (an orphan branch
    whose only commit holds the archived tree): each parks
    `archived-item-unproven: ` followed by the cause of API Contract 5 with
    no Implement call and no review, and `roundfix deliver status` prints
    `Park: <slug> finding: ` followed by the next action of API Contract 3.
  - `TestDeliverStartRevalidatesAContinuedBranchWithoutAnArchive`: a
    continued branch whose head still holds the Spec folder reaches the
    Implement call once, as on the Delivery Base.
- Unchanged and run by name: `TestResumeAcceptsARealArchiveCommit`,
  `TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord`,
  `TestInspectItemReportsAnArchivedSpec`,
  `TestDeliverStartContinuesTheItemBranchThatHoldsWork`,
  `TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview`,
  `TestDeliveryPublicationRequestsNoCodeRabbitReview` and
  `TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug`, the readers
  of `archivePaths` and of an archived item.

## Build Order

1. The `deliver` guide, the Roundfix `deliver` reference, the skill version,
   the Behavior Surface Record and the glossary (no dependency: the text
   states the contracts of this TechSpec, and its Verification reads only
   documents).
2. The engine dependency, the `queued` path, the blocker, its Park Class and
   its retry, with their tests (depends on: 1, which documents the `deliver`
   surface this step changes).
3. `ProveArchivedItem`, its wiring, `archivePaths` and the end-to-end tests
   through `roundfix deliver start` (depends on: 2, whose interface it
   implements).
4. QA gate (depends on: 1, 2, 3).

No Spec is listed in `requires`. The cycle AO Specs 0297 and 0298 deliver
first; 0298 edits other functions of `internal/delivery/engine.go`, and every
documentation Task raises the skill version over what the tree holds when it
starts, so this plan rebases on them.

## Risks & Considerations

- A resumed head costs one pre-PR review; with the earlier queue's review
  record still beside the item worktree, the review lineage makes it a round
  2 of the delta since that head.
- A review that raises findings on the archived Spec parks
  `corrective-spec-required`, as ADR-0223 provides for any archived
  candidate.
- `deliver start` still loads the Spec from the user checkout, so a Spec
  already archived on that checkout is refused before a queue exists; a Pull
  Request merged by hand is recorded by a Delivery Retry (ADR-0237).
- The proof's commit is the newest commit that touched `<source>/_prd.md`
  reachable from the head; a merge that removed the folder is judged against
  its first parent and is refused unless that change is an exact archive.

## Glossary

- adds: **Archived Item Resume**
- changes: **Delivery Revalidation**

## Decisions

- A proven archive on a continued item branch resumes at review; see
  ADR-0312.
- The proof reuses the queue's archive-commit proof on the commit that
  removed the Spec's `_prd.md`; no subject is read.
- The unproven park is a new blocker in the `finding` Park Class whose retry
  returns the item to `queued`.
- `archivePaths` tolerates an absent Specs Root, as `InspectItem` already does.
