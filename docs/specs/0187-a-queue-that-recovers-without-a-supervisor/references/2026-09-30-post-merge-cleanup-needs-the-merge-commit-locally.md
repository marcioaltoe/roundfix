---
type: fix
status: promoted
created: 2026-09-30
spec: 0187-a-queue-that-recovers-without-a-supervisor
reason: null
---

# Post-merge cleanup fails when the merge commit was never fetched

## Symptom

After the Delivery Queue squash-merged Spec 0183 (#282) on 2026-09-29, `deliver status` kept the warning `cleanup failed: release merged Spec "0183-…" Runs: merge commit "f3cbf0c4…" does not resolve to a commit`. It later read `is not on default branch`. The item worktree and the Spec's terminal Runs were left behind for a manual `reconcile`.

## Where

The post-merge release of a merged Spec's Runs in the delivery engine (Spec 0175, ADR-0161): it resolves the recorded merge commit before refreshing the default branch.

## Expected

The queue fetches the default branch before it releases a merged Spec's Runs, so the recorded merge commit resolves and is on the default branch. The worktree and the Runs are released without a warning.

## Evidence

`roundfix deliver status` output on 2026-09-29 from 19:26 to 22:59, and `docs/findings/2026-09-29-the-first-full-queue-trial-needed-manual-recovery.md` (Wave 5 addendum).
