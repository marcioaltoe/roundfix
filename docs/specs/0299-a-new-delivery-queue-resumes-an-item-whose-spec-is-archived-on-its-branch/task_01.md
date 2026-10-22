---
task: task_01
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
status: pending
type: docs
complexity: low
---

# Task 01: The deliver guide, the skill reference and the glossary describe the resume of an archived item branch

## Overview

The `deliver` user guide and the Roundfix skill's `deliver` reference say a
continued item branch starts at `queued` and is always revalidated before its
first Run, which is what parks an item whose branch already archived its Spec.
This Task states the resume, the `archived-item-unproven` blocker and its
retry in those documents, raises the Roundfix skill version, re-records the
Behavior Surface Record, adds **Archived Item Resume** to the glossary and
revises **Delivery Revalidation** (`_prd.md` → Core Feature 4).

## Requirements

1. MUST add to `docs/user-guide/commands/deliver.md`, right after the
   paragraph that ends "With no branch holding work, it creates a new item
   branch from the refreshed default branch.", this paragraph:
   "When the continued item branch already archives the Spec, the owner
   resumes the item at review instead of revalidating it. Its worktree head
   holds the Spec's Archive Record, or a legacy archive folder, and not the
   Spec folder, and the commit that removed the Spec's `_prd.md` is an exact
   archive of the Spec, proven by its change as the queue proves its own
   archive commit. The owner records that head as the candidate and moves the
   item to `reviewing`: it runs neither Delivery Revalidation nor the
   Implement Command, the pre-PR review runs per the configured policy, the
   archive step finds the Spec archived and commits nothing, and an open Pull
   Request of the item branch is found rather than opened again. When the head
   lacks the Spec folder and no such archive is proven, the item parks as
   `archived-item-unproven: <cause>` before any Run, where the cause is
   `item head <commit> holds neither <Spec folder> nor an archive of the
   Spec`, `no commit on the item branch removed <Spec folder>/_prd.md` or
   `commit <commit> that removed <Spec folder> is not an exact archive of the
   Spec`. A Delivery Retry of that park returns the item to `queued` and
   counts a retry, so the owner proves the branch again."
2. MUST add the same paragraph to the canonical Roundfix reference
   `.agents/skills/roundfix/references/deliver.md`, right after its paragraph
   that ends "With no branch holding work, it creates a new item branch from
   the refreshed default branch.".
3. MUST append, in both documents, to the paragraph that begins "After the
   item worktree is created or continued and before the first Run, Roundfix
   runs the strict Spec Consistency Check in the worktree.", the sentence "An
   item resumed from an item branch that already archives its Spec is not
   revalidated."
4. MUST insert `` `archived-item-unproven` `` between
   `` `revalidation-failed` `` and `` `secret-found` `` in the `finding` row of
   the Park Class table of both documents, and add, right after the
   paragraph that follows each table and begins "An unrecovered Provider
   System Error parks the item as", the sentence "An `archived-item-unproven`
   park's next action is `in <worktree>, make the item branch hold the Spec
   folder or one exact archive commit of the Spec, then run roundfix deliver
   retry <slug>`."
5. MUST add **Archived Item Resume** to `CONTEXT.md` through
   `domain-modeling`, right after the **Delivery Revalidation** entry, as
   "**Archived Item Resume**:" followed on the next line by "The start of a
   new Delivery Queue item whose continued Item Branch already archives its
   Spec in a commit proven to be an exact archive of it: the item enters
   review with its head as the candidate, without Delivery Revalidation, an
   Implement Run or a new archive. An item branch whose archive cannot be
   proven parks as `archived-item-unproven` (ADR-0312)." and the line
   "_Avoid_: Archived retry, requeue". The ADR reference takes the glossary's
   link form, to
   `docs/adr/0312-a-new-delivery-queue-resumes-an-archived-item-branch-at-review.md`.
6. MUST revise **Delivery Revalidation** in `CONTEXT.md` through
   `domain-modeling` by appending to its definition the sentence "It does not
   run for an Archived Item Resume." and keep its `_Avoid_` line.
7. MUST run these commands, because this Task edits an owned skill and a
   guide. Before them, raise both version fields (`metadata.version` and the
   top-level `version`) of the canonical `.agents/skills/roundfix/SKILL.md`
   by one patch over the value the tree holds when this Task starts, because
   the record command refuses changed content under an unchanged version.
   The Agent runs them itself before handing back, because the Daemon never
   regenerates a record; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     `deliver` guide's fingerprint changes.
8. MUST search the guides, the agent guides and the skill copies for wording
   that says every queued or continued item is revalidated, or starts its
   first Run, before review, and report what it found and changed.
9. MUST NOT change any Go source, ADR, `docs/references/skill-coverage.json`,
   any other guide page or any text inside a `### QA settlement` section.

## Subtasks

- [ ] Add the resume paragraph and the revalidation sentence to the guide and the reference.
- [ ] Add the blocker to both Park Class tables with its next action.
- [ ] Add **Archived Item Resume** and revise **Delivery Revalidation** in the glossary.
- [ ] Raise the skill version, sync the mirrors, record the version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The guide and both copies of the reference describe the resume, the three causes, the retry, and the blocker's Park Class and next action.
- [ ] `CONTEXT.md` defines **Archived Item Resume**, links ADR-0312, and says **Delivery Revalidation** does not run for it.
- [ ] The skill version is raised and recorded, the mirrors are in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirror, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `docs/adr/0312-a-new-delivery-queue-resumes-an-archived-item-branch-at-review.md`
- instruction: `docs/adr/0223-a-delivery-keeps-its-work-across-archive-requeue-and-review.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do n="$(tr -s '[:space:]' ' ' < "$f")"; for p in 'the owner resumes the item at review instead of revalidating it' 'is not an exact archive of the Spec' 'returns the item to .queued. and counts a retry' 'is not revalidated.' 'revalidation-failed., .archived-item-unproven., .secret-found' 'make the item branch hold the Spec folder or one exact archive commit of the Spec'; do printf '%s\n' "$n" | grep -q -- "$p" || { printf 'missing %s in %s\n' "$p" "$f" >&2; exit 1; }; done; done; c="$(tr -s '[:space:]' ' ' < CONTEXT.md)"; printf '%s\n' "$c" | grep -qF -- '**Archived Item Resume**:' || { echo 'missing glossary term' >&2; exit 1; }; printf '%s\n' "$c" | grep -qF -- 'It does not run for an Archived Item Resume.' || { echo 'missing revalidation revision' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document describes the resume and the glossary lacks **Archived Item Resume**, so the command fails at its first phrase check. Afterwards the guide and both reference copies describe the resume, the causes, the retry and the Park Class, the glossary defines the term and revises **Delivery Revalidation**, the mirrors are in sync, the skill version is recorded, and the documentation contracts, the Skill Coverage Map among them, pass.

## References

- `_prd.md` → Core Feature 4; Goals; Glossary
- `_techspec.md` → API Contract 7; Build Order 1
- ADR-0223
- ADR-0312
