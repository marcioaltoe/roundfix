# QA-09 — Active Run guard

Status: pass.

Fresh independent fixtures passed:

- `TestAuditPreservesActiveRunSurvivors` injects a pushed-and-merged Active Run
  owning a branch and worktree. Both stay non-residue, name the Run in their
  evidence, and carry no reclaim command.
- `TestInspectTerminalRunActiveRunDoesNotInspectDeletedTargetContent` proves an
  Active Run bypasses the deleted-target content resolver.

The public real-repository audit had no Agent-visible Run bound to this
worktree, so it safely classified it `preserved`; it did not touch a live
writer. The full repository gate and both named selections exited 0.
