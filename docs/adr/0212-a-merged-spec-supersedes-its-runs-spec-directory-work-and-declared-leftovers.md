---
status: accepted
created_at: 2026-10-01T00:00:00Z
updated_at: 2026-10-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A merged Spec supersedes its Runs' Spec-directory work and declared leftovers

ADR-0161 releases a terminal Run of a merged Spec when every commit its Run
Branch holds and the merged head lacks is represented at that head. Two kinds
of work were never represented after a squash merge and an archive: commits
that changed only the Spec's own directory, whose files the archive moved and
restamped, and the uncommitted leftovers of a Task that a later Run settled.
Reconcile classified those Runs `unintegrated` or `dirty` and kept them, and
on 2026-10-01 this repository held 32 Run Worktrees of merged Specs, 19.1 GB.

When the merged head holds the Spec archived, its archived directory now
represents every path a Run changed under the Spec's directory, and a Run
Worktree's uncommitted changes are superseded when each changed path is one
the Spec's Tasks declared or recorded, or lies under the Spec's directory.
Any other path, committed or not, keeps the Run preserved, and a Task commit
is still represented only by its Task being completed at the merged head.

## Consequences

This refines ADR-0053 and ADR-0161: dirty work is still preserved unless the
merged head holds the archived Spec and every dirty path belongs to the Spec's
declared or recorded scope. Releasing remains an explicit act, through
`roundfix reconcile --apply` or the delivery owner's release after a merge,
and `--discard-superseded` keeps its own Branch Disposition (ADR-0115). An
uncommitted change outside that scope, such as a hand edit to an unrelated
file, still keeps the Run.
