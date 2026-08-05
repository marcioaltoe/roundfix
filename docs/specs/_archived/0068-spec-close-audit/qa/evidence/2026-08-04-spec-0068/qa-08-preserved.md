# QA-08 — Unmatched worktree preservation

Status: pass.

The built CLI on the real Spec classified the current worktree `preserved`,
named the missing matching Run, and emitted no reclaim command. A fresh
worktree-list read after repeat execution proved it remained registered.

`TestAuditPreservesUnmatchedWorktree` passed in the fresh 12-test selection and
independently requires `preserved`, non-empty missing-Run evidence, and empty
reclaim. QA-18 separately tests the pushed-and-merged Supervisor scratch
promise and fails it; this row covers Task 02's general unmatched safety rule.
