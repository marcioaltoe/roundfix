---
type: fix
status: promoted
created: 2026-10-01
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
reason: null
---

# Archiving or requeueing a Spec loses or breaks its delivery

## Symptom

Four failures in the delivery of Specs 0205, 0208 and 0210 on 2026-10-01 share one context: what happens to a delivery item when its Spec moves into `docs/history/specs/` or when its queue item is recreated.

1. **A test pinned to the active Spec path fails after the archive commit.** `internal/judge/questions_test.go` read `docs/specs/0205-…/_techspec.md`. The archive commit moved the Spec, and the repository gate failed (`open ../../docs/specs/0205-…/_techspec.md: no such file or directory`). The same pattern already broke Specs 0191 and 0194.
2. **Another active Spec's Task Context names the archived Spec's active path.** After 0205 was archived, `0218/task_02.md` declared `interface: docs/specs/0205-…/_techspec.md`. `TestCheckActiveCorpusHasNoErrors` failed CI on the delivery PR with `SC-REF-UNRESOLVED`.
3. **A retry after the archive refuses when the head moved.** It printed `archived item head … differs from candidate head …`, so the operator opened the PR by hand. This happened for 0193, 0203 and 0205.
4. **A new queue for a slug whose item reached its retry limit starts from the default branch.** `deliver start` created a second item worktree from `main` and would have rerun every completed Task. The first item branch held four completed Tasks and five amendments. The operator stopped the owner before any Run and merged the old branch into the new one. It happened again on 2026-10-03 for Spec 0217, after Spec 0211 had merged: its parked item had been replaced by the queue that delivered Spec 0218, and `deliver start` recreated the item from `main`; the operator again merged the old item branch by hand.

## Expected

- The archive step rewrites, or refuses on, every repository path that names the Spec's active directory: tests, other active Specs' Task Context, and their references. It reports what it changed.
- A retry after a post-archive correction resumes from the archived item's current head.
- `deliver start` for a slug whose item already exists continues that item branch. When the item cannot be continued, it refuses and names the existing branch.

## Relation

Spec 0211 adopts the archived-retry entry, which covers the start head an archived retry looks for. This entry covers the moved head (3) and the requeue (4). Group items 1 and 2 with Spec 0212's archive and reconcile work, or with 0211, within the four-Task bound (ADR-0208).
