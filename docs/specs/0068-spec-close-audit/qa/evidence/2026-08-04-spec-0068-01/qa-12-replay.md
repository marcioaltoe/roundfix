# QA-12 — Motivating-session replay

Status: pass.

`TestAuditReplaysMotivatingSessionResidue` passed in the fresh named audit
selection. Its real-Git fixture asserts the source finding path, two unpushed
scratch worktrees, the content-resolved orphaned Run Worktree, the stale remote
backup branch with an exact reclaim command, and Pull Request #58 and #68
branches without reclaim commands. Every survivor must carry evidence.

QA-18 separately exercises the pushed-and-merged scratch state added by the
corrective Task. The complete 17-test audit selection exited 0.
