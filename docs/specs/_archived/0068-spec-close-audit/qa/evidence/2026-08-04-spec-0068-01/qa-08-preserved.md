# QA-08 — Unmatched worktree preservation

Status: pass.

The public CLI classified the current Run Worktree `preserved`, named the
missing matching Run and five unintegrated changed paths, and emitted no
reclaim command. A fresh JSON invocation agreed, and the worktree remained
registered.

The fresh named audit selection passed
`TestAuditPreservesUnpushedScratchWorktree`,
`TestAuditPreservesUnmergedScratchWorktree`, and
`TestAuditPreservesIndeterminateScratchWorktree`. Each fixture requires
`preserved`, evidence naming the failed proof, and an empty reclaim command.
QA-18 separately proves the pushed-and-merged exception.
