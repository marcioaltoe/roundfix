---
status: pending
type: backend
---

# Task: Integrate Settled Commits

Commit and integrate settled work onto Run Branch.

## Work
- Stage changes with `git add --all`
- Create standard Task commit
- Integrate onto Run Branch
- Remove Task Worktree on success

## Verification
- `git log --oneline -n 1 | grep -q "^[a-f0-9]*"`

## Result
Settled commits integrate correctly.
