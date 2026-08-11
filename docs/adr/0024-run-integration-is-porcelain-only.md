---
status: accepted
created_at: 2026-07-05T22:17:04Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Run integration is porcelain-only

A Run's commits reach the user's branch exclusively through porcelain that keeps ref, HEAD, index, and working tree consistent: when the user's checkout is still on the Run's branch, `git merge --ff-only <run-sha>` inside that checkout; when they have switched away, ancestry-checked `git branch -f`. Plumbing ref updates are forbidden — `git update-ref` on a checked-out branch was empirically shown to leave the user's checkout with phantom staged deletions that a routine commit would turn into a reversion of the Run's work. When both porcelain paths refuse (overlapping local changes, mid-Run user commits), the Run ends in the new terminal state Integration Pending: commits preserved on the Run Branch, worktree kept, and the exact integration command reported — never a silent divergence, never a forced move.
