---
status: accepted
created_at: 2026-07-17T16:00:35Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Terminal Run Worktree reconciliation is explicit and proof-based

`roundfix reconcile [run-id]` classifies terminal spec Run Worktrees and Run Branches and is read-only unless `--apply` is supplied; apply removes only a clean Run Worktree whose Run Branch tip is proven reachable from its recorded target branch. Ambiguous, dirty, and unintegrated work stays intact, while a proven Integration Pending reconciliation may promote that Run to Clean through ADR-0052's guarded transition; GC remains limited to retained storage.
