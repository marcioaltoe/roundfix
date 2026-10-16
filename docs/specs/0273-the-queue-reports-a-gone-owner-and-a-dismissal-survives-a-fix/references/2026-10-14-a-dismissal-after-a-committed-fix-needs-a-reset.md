---
type: fix
status: promoted
created: 2026-10-14
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
reason: null
---

# A dismissal after a committed fix needs a reset

## Symptom

A parked item had two review findings: one real, one to dismiss. After the
fix was committed, `roundfix review dispose F1 --dismiss` refused because
HEAD no longer equals the reviewed head, and checking out the reviewed head
detached is refused too (`HEAD is detached`). The only route was to reset the
branch to the reviewed head, dismiss, and reset forward again.

Separately, `roundfix deliver retry <slug>` run inside the item worktree says
`Delivery Queue for repository "<item worktree>" does not exist` instead of
naming the checkout that owns the queue.

## Where

`roundfix review dispose` and `roundfix deliver retry`.

## Expected

A dismissal is accepted when the reviewed head is an ancestor of HEAD and the
finding's anchor is unchanged, or the `deliver` guide says to dismiss before
committing fixes. The retry refusal names the main checkout.

## Evidence

Spec 0268 delivery on 2026-10-14 at 12:22 BRT: fix `825406ac` committed on
reviewed head `bf2da46e`; dismissal recorded after the reset round-trip.
