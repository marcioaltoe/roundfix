---
type: fix
status: promoted
created: 2026-10-04
spec: 0227-reconcile-releases-the-runs-of-merged-specs
---

# Reconcile keeps the Runs of Specs that are already merged

## Problem

On 2026-10-04 the maintainer found 12 `roundfix/run-…` and item branches
in the main checkout. `roundfix reconcile --apply` released only the 4
classified `superseded`. It preserved 6 Runs as `unintegrated` or `dirty`,
from Specs 0181, 0184 (two), 0200, 0204 and 0207, every one archived and
merged on main. It also never touched two item branches, for 0205 and 0217,
that the queue had recreated from scratch. The reasons it gave were "21
differing shared files", "1 Run-only file" (a skill main later renamed),
"no superseding QA Report", and uncommitted changes of a Task that was
later redone. The operator removed them by hand
(`~/.roundfix-operator/queue-interventions.md`, entries 144 and 145).

Every queue merge on 2026-10-04 also ended with
`warning: cleanup failed: release merged Spec "<slug>" Runs: candidate head
"<sha>" is not represented by merge commit "<sha>"` (0222, 0224, 0226). In
each case the only difference was other Specs' commits that reached main
while the item ran, so the Runs of a delivered Spec are left for a later
manual reconcile.

## Expected

A Run whose Spec is archived on main, and whose Spec delivery merged, is
releasable even when shared files diverged later, because main supersedes
it. An item branch whose item was recreated, and whose Spec merged, is
releasable too. The post-merge cleanup accepts a merge commit that contains
the candidate's changes plus later main commits.

## Notes

Non-binding. The content-proof rule that Spec 0212 strengthened should stay
for Specs that are not merged. Merge-evidence (the squash commit that names
the Spec and is reachable from main) could stand in for content
representation.
