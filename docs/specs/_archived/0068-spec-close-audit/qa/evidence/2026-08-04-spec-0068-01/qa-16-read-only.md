# QA-16 — Read-only boundary

Status: pass.

Before and after the built help, repeated real-Spec audits, Spec Check, and
Reconcile dry runs:

- `git show-ref` was identical.
- `git worktree list --porcelain` was identical.
- `git status --short` was identical apart from the same QA rerun paths present
  at both snapshots.
- `/Users/marcio/.roundfix/roundfix.db` retained SHA-1
  `7a0fa0d5540b442e177fca7bcd468668c00d1fda`.

The public scratch fixture also retained refs and worktrees across two audits.
`TestAuditScratchWorktreeGitStateUnchanged` and
`TestAuditDeliveryCheckLeavesGitStateByteIdentical` passed. A production sweep
for `net/http`, network calls, and `os.Create`, `os.WriteFile`, or
`os.RemoveAll` returned no match; production imports contain no transport
package. Reclaim operations remain returned strings.
