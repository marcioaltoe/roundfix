---
status: accepted
created_at: 2026-07-17T16:00:35Z
updated_at: 2026-11-09T15:00:00Z
deprecated_at: null
superseded_by: null
---

# Terminal Run Worktree reconciliation is explicit and proof-based

`roundfix reconcile [run-id]` classifies terminal spec Run Worktrees and Run Branches and is read-only unless `--apply` is supplied; apply removes only a clean Run Worktree whose Run Branch tip is proven reachable from its recorded target branch. Ambiguous, dirty, and unintegrated work stays intact, while a proven Integration Pending reconciliation may promote that Run to Clean through ADR-0052's guarded transition; GC remains limited to retained storage.

**Amended by [ADR-0379](0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md) (2026-11-09).** Dirty work of an ended Run whose Spec's Archive Record is on the default branch and whose commit proof holds now has one explicit exit, the Dirty Discard (`roundfix reconcile --discard-dirty <digest>`), which records a Branch Disposition naming the discarded paths before it removes anything; every other dirty, ambiguous or unintegrated work stays intact, and `--apply` still releases only proven work.
