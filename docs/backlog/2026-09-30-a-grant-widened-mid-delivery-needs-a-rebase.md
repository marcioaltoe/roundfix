---
type: fix
status: open
created: 2026-09-30
spec: null
reason: null
---

# A grant widened mid-delivery covers a Task only after the item is rebased

## Symptom

Spec 0181's grant was widened on main (#277) and cherry-picked onto the item branch. task_07 then ran on top of it. The QA mechanical authorization audit still refused task_07's commit (`QA-AUTH-PATHS`), reporting the grant at revision `6784210b`, where the path was absent. Merging `origin/main` into the item did not help. Rebasing the item branch onto main did.

## Where

The mechanical authorization audit's revision for each Task commit (`internal/speccheck/mechanical.go` `detectMechanicalAuthPaths` and `internal/daemon/task_engine.go` around the QA mechanical request), and how the Delivery Base or the authorization revision is chosen for a commit made before a merge.

## Expected

A grant that is on the delivery target's default branch when the gate runs covers every Task commit of the candidate, whether the item branch merged main or was rebased onto it. A path outside that grant is still refused. The rule is stated in the QA gate guidance.

## Evidence

The 0181 QA reports on Run branches `run_20260929T200238Z_ca4a5877c9dc9d5c` and `run_20260929T201521Z_3d1a54f26096ba81`, and `docs/findings/2026-09-29-the-first-full-queue-trial-needed-manual-recovery.md`, section 1.
