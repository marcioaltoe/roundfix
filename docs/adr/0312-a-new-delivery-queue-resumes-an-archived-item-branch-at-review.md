---
status: accepted
created_at: 2026-10-22T06:00:00Z
updated_at: 2026-10-22T06:00:00Z
deprecated_at: null
superseded_by: null
---

# A new Delivery Queue resumes an archived item branch at review

ADR-0223 lets `deliver start` continue the one item branch of a Spec that
holds commits the default branch lacks, and a new queue starts that item at
`queued`, whose first step is Delivery Revalidation of the active Spec folder.
On 2026-10-21 Spec 0292's item had already archived its Spec and recorded its
green Pull Request #619 when its queue ran out of retries; the new queue
parked it `delivery-error` because the Spec folder no longer existed on the
item branch, and a Delivery Retry could not help, because ADR-0229 accepts an
archived item without a recorded candidate only under a QA override.

## Decision

**An archive on the item branch is proven by its change.** When a queued
item's worktree head holds the Spec's archive (its Archive Record, or a legacy
archive folder) and not its folder, the commit that removed the Spec's
`_prd.md` must be an exact archive of the Spec, judged against its first
parent by the same proof the queue applies to its own archive commit; its
subject is never read.

**A proven archive resumes at review.** The item records its head as its
candidate and enters `reviewing`. It runs neither Delivery Revalidation nor
the Implement Command. The pre-PR review runs per the configured policy on
that head, the archive step finds the Spec already archived and commits
nothing, and the repository gate, publication, the required checks and the
merge follow as for any candidate. An open Pull Request of the item branch is
found, not opened again.

**Anything else parks with its cause.** A head that holds neither the folder
nor an archive, a branch with no commit that removed the folder, or a removing
commit that is not an exact archive parks the item as `archived-item-unproven:
<cause>` in the `finding` Park Class. A Delivery Retry of that park returns the
item to `queued`, counted, so the owner proves the branch again.

**An archived item's Specs Root may be absent.** The queue resolves an item's
archive paths from the configured Specs Root even when the item worktree has
no such directory, because archiving the only active Spec removes it.

## Considered Options

- Resume at the repository gate, after the archive, without a review:
  rejected. The new queue holds no evidence that any review saw the item
  head, and an exact archive commit does not prove one, because an operator
  archive produces one too; publication still needs a review of the current
  candidate under the configured policy.
- Revalidate the archived Spec from its record at `source_revision`:
  rejected. Revalidation guards a first Run against a moved premise, and no
  Run starts.
- Require a QA override, as ADR-0229 does for a retry without a candidate:
  rejected. The archive's exactness anchors the head, and the Archive Command
  already refused an ineligible QA when it wrote that archive.

## Consequences

A new queue finishes a delivery whose archive was already made without
re-running implementation, the archive or a hand-made Pull Request. The
resumed head is reviewed again, so a review that raises findings on the
archived Spec parks `corrective-spec-required` as ADR-0223 already provides.
ADR-0223's continuation rule and ADR-0229's retry rule are otherwise
unchanged: an item branch without an archive still revalidates and runs.
