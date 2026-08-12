---
status: accepted
created_at: 2026-07-15T16:53:06Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Review Runs execute in the user's checkout

Review Runs (fetch, resolve, watch) no longer create a Run Worktree: they operate directly on the checked-out PR Head Branch, guarded by a deterministic Branch Integrity Preflight that requires zero unintegrated Run Branch commits and no other Run bound to the branch before any work starts, with only an explicit, PR-comment-audited bypass. This partially supersedes ADR-0023, which keeps applying to spec Runs (implement), where Task concurrency needs isolation. Observed motivation: a review fix is by definition a delta over the pull request's published HEAD; the worktree indirection let review Runs start from a HEAD missing stranded Run Branch work and produced Integration Pending outcomes and Final Pushes that silently omitted completed work.
