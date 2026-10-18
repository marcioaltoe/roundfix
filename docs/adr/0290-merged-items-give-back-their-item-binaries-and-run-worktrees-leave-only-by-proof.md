---
status: accepted
created_at: 2026-10-18T17:30:00Z
updated_at: 2026-10-18T17:30:00Z
deprecated_at: null
superseded_by: null
---

# Merged items give back their Item Binaries, and Run Worktrees leave only by proof

Roundfix Home grew without a stated policy. Measured on the maintainer's
machine on 2026-10-18, `~/.roundfix/artifacts` held 412 MB, and 327 MB of it
(79%) were ten Item Binary copies of 32 MB each: five for Spec 0277 and five
for Spec 0278, both merged. ADR-0282 states that "this decision adds no
reclamation for them". `roundfix gc --dry-run` found nothing for Journal
Retention, three Runs for Run Retention with no reclaimable bytes, and no
Worktree Build Cache. The Run Worktree location held 3.4 GB, all of it other
repositories' Delivery Queue items, including one item that was still running.
This repository held no Run Worktree. The maintainer's standing rule of
2026-08-10 forbade `gc` or artifact cleanup until a Spec stated the policy.

## Decision

**What is reclaimed, by whom, and when:**

| Area | Kept while | Reclaimed by |
| --- | --- | --- |
| Run Event Journal and run artifact directory | the Run is Active or younger than Journal Retention | the GC Command and Run Retention (ADR-0033, ADR-0255) |
| Run row and its records | the Run is not terminal, is queue-referenced, keeps its Run Worktree, or is younger than Run Retention | Run Retention (ADR-0255) |
| Run Worktree and Run Branch | reconciliation does not prove them `safe` or `superseded` | `roundfix reconcile --apply` and the Delivery Queue's release after a merge, never by age (ADR-0053, ADR-0212, ADR-0289) |
| Worktree Build Cache | its Run is not terminal or its item is not `merged` | the GC Command (ADR-0266) |
| Item Binary | its Delivery Queue item exists and is not `merged` | the Delivery Queue's release after a merge, and the GC Command |
| Pre-PR Review Records, dispositions and delivery logs | always | nothing |
| Archive evidence | always | nothing: it lives in Git |

**A merged item gives back its Item Binaries.** When the Delivery Queue
releases a merged item's Runs, it also removes
`<artifact dir>/delivery/<slug>/item-binary/`. A failure joins the item's
cleanup warning and removes nothing else; the item's logs stay.

**The GC Command reclaims the Item Binaries of items that no longer run.** In
the repository it runs in, an Item Binary directory is reclaimable when the
repository's Delivery Queue has no item for its Spec or holds that item at
stage `merged`. Any other stage keeps it. GC rereads the queue immediately
before each removal, `--dry-run` lists what it would remove and writes
nothing, and a repository with no Item Binary directory prints exactly what it
prints today.

The operator previews the policy with `roundfix gc --dry-run`,
`roundfix gc compact` and `roundfix reconcile`, which are all read-only.

## Considered Options

- Reclaiming Run Worktrees by age was rejected. A dirty or unintegrated Run
  Worktree is unmerged work, and ADR-0053 releases only on proof.
- Removing older Item Binary copies while the item still runs was rejected.
  ADR-0282 copies each build to a new file because a running step may still map
  an earlier copy.
- Reclaiming Item Binaries across every repository from any checkout was
  rejected. The GC Command resolves the artifact directory of the repository
  it runs in, as the Worktree Build Cache sweep does.
- Removing a merged item's whole artifact directory was rejected. Its logs are
  the evidence a cleanup warning and a later diagnosis read.

## Consequences

This amends ADR-0282's consequence that Item Binary copies accumulate. On the
2026-10-18 measurement the policy frees the 327 MB of Specs 0277 and 0278 at
the first `roundfix gc`, and each later merge frees its own copies. It frees no
Run Worktree, because reclaiming one stays a proof.
