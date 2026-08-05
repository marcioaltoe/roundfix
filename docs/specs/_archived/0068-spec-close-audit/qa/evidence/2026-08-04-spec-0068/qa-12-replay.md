# QA-12 — Motivating-session replay

Status: pass.

The source finding confirms four end-state categories: two Supervisor scratch
worktrees, an orphaned Run Worktree after target deletion, a stale remote
backup branch, and two branches held by unmerged Pull Requests.

`TestAuditReplaysMotivatingSessionResidue` passed in the fresh 12-test
selection. Its real-Git fixture asserts the source path, both preserved
unpushed scratch worktrees, content-resolved orphan residue, exact remote
backup reclaim command, and Pull Request #58 and #68 branches with empty
reclaim. Every survivor requires evidence. QA-18 adds the distinct pushed and
merged scratch state absent from this replay and records F-001.
