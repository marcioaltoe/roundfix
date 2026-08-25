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
- `grep -q "stageSharableChanges\|git add\|addAllChanges" internal/cli/settle.go && grep -q "integrateTask\|runGit.*commit" internal/cli/settle.go`


## References

- Core Feature 2: Settle Recovery

## Result
Settled commits integrate correctly.
