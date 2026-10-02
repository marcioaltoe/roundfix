---
type: fix
status: promoted
created: 2026-10-01
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
reason: null
---

# Reconcile releases the Runs of a squash-merged Spec

## Opportunity

`roundfix reconcile` preserves the Runs of Specs merged by squash. Each one is classified `unintegrated` ("default branch has no superseding QA Report after target branch disappeared", or a content comparison that never closes after the archive move) or `dirty` (leftovers of an interrupted Task). In Vortex on 2026-09-25 that held 13 Runs and 11 GB. In this repository on 2026-10-01 it held 12 Runs, and with them the Run worktrees reached 16 GB, contributing to a disk-full stop of a delivery. The repository rules forbid deleting Roundfix worktrees by hand.

## Value

Storage stays bounded without manual cleanup, and a merged Spec's Runs stop occupying the disk.

## Shape

Treat a Spec merged into the default branch, whose archive is in `docs/history/specs/<slug>` and whose PR is merged, as superseding every terminal Run of that Spec, including `dirty` Runs whose Tasks the merged history settled. Keep refusing any Run whose Task commits are not represented. Evidence: Secondbrain `inbox/roundfix/2026-09-25-reconcile-nunca-libera-run-de-spec-mergeada-por-squash.md`.
