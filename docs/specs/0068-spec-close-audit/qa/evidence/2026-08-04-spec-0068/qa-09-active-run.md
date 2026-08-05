# QA-09 — Active Run guard

Status: pass.

- `TestAuditPreservesActiveRunSurvivors` passed in the fresh 12-test assembled
  selection. It injects one Active Run owning a branch and worktree, requires
  both to avoid `residue`, requires evidence naming the Run ID, and requires
  empty reclaim commands.
- `TestInspectTerminalRunActiveRunDoesNotInspectDeletedTargetContent` passed in
  the six-test deleted-target selection and proves an Active Run bypasses
  content comparison.
- The complete worktree and CLI packages passed 1,041 tests.

The real worktree could not serve as Active Run proof because the Agent-visible
Run Database does not bind that path; it safely classified `preserved`. No live
writer was touched.
