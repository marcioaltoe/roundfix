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
- `rm -f /tmp/test-deletion.txt && touch /tmp/test-deletion.txt && git add /tmp/test-deletion.txt && rm /tmp/test-deletion.txt && git add --all && git diff --cached --name-status | grep -q "^D"`


## References

- Core Feature 3: Deleted File Handling

## Result
Deleted files stage without errors.
