---
status: pending
type: backend
---

# Task: Handle Deleted File Staging

Ensure deleted files stage correctly.

## Work
- Verify `git add --all` handles deletions
- Test with a Task that deletes a file
- Confirm deletion commits cleanly

## Verification
- `grep -q "git add --all\|addAllChanges" internal/cli/settle.go && grep -q "name-status\|deletedPaths" internal/cli/settle.go`


## References

- Core Feature 3: Deleted File Handling

## Result
Deleted files stage without errors.
