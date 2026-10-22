---
type: fix
status: promoted
created: 2026-10-22
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
reason: null
---

# A new Delivery Queue parks an item whose Spec is already archived on its branch

## Symptom

On 2026-10-21 the Delivery Queue item of Spec 0292 used up its five retries on
`checks-timeout` parks while its Pull Request #619 was green: the item branch
already carried the Spec's archive commit and the Pull Request Record Commit.
The operator then started a new queue with `roundfix deliver start
0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it`. It printed
`Continuing item branch roundfix/deliver-0292-…` and the owner parked the item
`delivery-error` before any Run: the revalidation's strict Spec Consistency
Check stat-ed `docs/specs/<slug>` in the item worktree and failed with
`no such file or directory`, because the archive commit had removed the Spec
folder. The operator merged #619 by hand and ran `roundfix deliver stop`.

## Where

`roundfix deliver start` continues the one item branch that holds commits the
default branch lacks (`existingItemBranches`, `internal/cli/deliver_workflow.go`,
ADR-0223). A new queue's item starts at `queued`; `Engine.advanceItem`
(`internal/delivery/engine.go`) calls `CreateItemBranch`, which reuses that
branch, and then always calls `Revalidate` (`internal/cli/deliver_revalidate.go`),
whose `strictSpecFindings` (`internal/cli/deliver_plan.go`) reads the active
Spec folder. The error is returned, so `Engine.Run` parks the item
`delivery-error: revalidate item: run strict Spec Consistency Check: read Spec
"<slug>": stat <worktree>/docs/specs/<slug>: no such file or directory`. When
the archived Spec was the only active one, the item worktree has no
`docs/specs` directory at all and the error is `resolve item Specs Root:
specs.root resolved to "<worktree>/docs/specs", which does not exist`.

A Delivery Retry cannot recover either park: `Engine.Retry` sees the archived
item and refuses `candidate head is missing`, because the new queue recorded
no candidate and the archive carries no QA override (ADR-0229).

## Expected

When the item branch already carries the Spec's archive commit, adopted by its
change as the queue adopts any archive commit, a new queue resumes the item
from the stage after the archive: it never runs the Implement Command or the
archive again, the pre-PR review runs per the configured policy, and the
repository gate, publication (finding the open Pull Request), the required
checks and the merge follow. A branch whose archive commit is not exact, or
whose Spec folder is missing for another reason, still parks, with a reason
that names the cause and the manual path.

## Evidence

The operator's record of 2026-10-21 in the cycle AO briefing for this Spec. On
2026-10-22 a temporary repository reproduced it through `roundfix deliver
start` at `34b4e246` (v0.107.0): an item branch with an exact archive commit
and a recorded Pull Request parked
`delivery-error: revalidate item: run strict Spec Consistency Check: read Spec
"0001-widget-flow": stat …/docs/specs/0001-widget-flow: no such file or
directory`, and, without another active Spec, `delivery-error: revalidate
item: resolve item Specs Root: specs.root resolved to "…/docs/specs", which
does not exist`. No Finding.
